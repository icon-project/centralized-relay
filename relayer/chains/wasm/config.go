package wasm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cometbft/cometbft/rpc/client/http"
	sdkClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayTypes "github.com/icon-project/centralized-relay/relayer/types"

	"go.uber.org/zap"
)

type Config struct {
	provider.CommonConfig  `json:",inline" yaml:",inline"`
	KeyringBackend         string        `json:"keyring-backend" yaml:"keyring-backend"`
	KeyringDir             string        `json:"keyring-dir" yaml:"keyring-dir"`
	AccountPrefix          string        `json:"account-prefix" yaml:"account-prefix"`
	Denomination           string        `json:"denomination" yaml:"denomination"`
	GasPrices              string        `json:"gas-prices" yaml:"gas-prices"`
	GasAdjustment          float64       `json:"gas-adjustment" yaml:"gas-adjustment"`
	MinGasAmount           uint64        `json:"min-gas-amount" yaml:"min-gas-amount"`
	MaxGasAmount           uint64        `json:"max-gas-amount" yaml:"max-gas-amount"`
	TxConfirmationInterval time.Duration `json:"tx-confirmation-interval" yaml:"tx-confirmation-interval"`
	BroadcastMode          string        `json:"broadcast-mode" yaml:"broadcast-mode"` // sync, async and block. Recommended: sync
	SignModeStr            string        `json:"sign-mode" yaml:"sign-mode"`
	Simulate               bool          `json:"simulate" yaml:"simulate"`
	ExtraCodec             string        `json:"extra-codecs" yaml:"extra-codecs"`
	BlockBatchSize         uint64        `json:"block-batch-size" yaml:"block-batch-size"`
}

func (pc *Config) NewProvider(ctx context.Context, log *zap.Logger, homePath string, _ bool, chainName string) (provider.ChainProvider, error) {
	pc.HomeDir = homePath
	pc.ChainName = chainName

	if pc.KeyringDir == "" {
		pc.KeyringDir = filepath.Join(pc.HomeDir, pc.NID)
	}

	if err := pc.Validate(); err != nil {
		return nil, err
	}

	pc, err := pc.sanitize()
	if err != nil {
		return nil, err
	}

	clientContext, err := pc.newClientContext(ctx)
	if err != nil {
		return nil, err
	}

	contracts := pc.eventMap()

	ws := newClient(clientContext)

	return &Provider{
		logger:      log.With(zap.Stringp("nid", &pc.NID), zap.Stringp("name", &pc.ChainName)),
		cfg:         pc,
		client:      ws,
		contracts:   contracts,
		eventList:   pc.GetMonitorEventFilters(contracts),
		routerMutex: new(sync.Mutex),
	}, nil
}

func (pc *Config) SetWallet(addr string) {
	pc.Address = addr
}

func (pc *Config) GetWallet() string {
	return pc.Address
}

func (pc *Config) Validate() error {
	if pc.ChainName == "" {
		return fmt.Errorf("chain-name cannot be empty")
	}

	if pc.HomeDir == "" {
		return fmt.Errorf("home-dir cannot be empty")
	}
	return nil
}

func (pc *Config) sanitize() (*Config, error) {
	if pc.BlockBatchSize == 0 {
		pc.BlockBatchSize = 50
	}
	if pc.Decimals == 0 {
		pc.Decimals = relayTypes.DefaultCoinDecimals
	}
	return pc, nil
}

func (c *Config) newClientContext(ctx context.Context) (sdkClient.Context, error) {
	codec := c.MakeCodec(moduleBasics, strings.Split(c.ExtraCodec, ",")...)

	keyRing, err := keyring.New(
		c.ChainName,
		c.KeyringBackend,
		c.KeyringDir,
		os.Stdin,
		codec.Codec,
		func(options *keyring.Options) {
			options.SupportedAlgos = types.SupportedAlgorithms
			options.SupportedAlgosLedger = types.SupportedAlgorithmsLedger
		},
	)
	if err != nil {
		return sdkClient.Context{}, err
	}

	cometRPCClient, err := http.New(c.RPCUrl, "/websocket")
	if err != nil {
		return sdkClient.Context{}, err
	}

	networkInfo, err := cometRPCClient.Status(ctx)
	if err != nil {
		return sdkClient.Context{}, err
	}

	return sdkClient.Context{
		CmdContext:        ctx,
		ChainID:           networkInfo.NodeInfo.Network,
		Client:            cometRPCClient,
		NodeURI:           c.RPCUrl,
		Codec:             codec.Codec,
		Keyring:           keyRing,
		KeyringDir:        c.KeyringDir,
		TxConfig:          codec.TxConfig,
		HomeDir:           c.HomeDir,
		BroadcastMode:     c.BroadcastMode,
		SignModeStr:       c.SignModeStr,
		Simulate:          c.Simulate,
		InterfaceRegistry: codec.InterfaceRegistry,
		AccountRetriever:  authtypes.AccountRetriever{},
	}, cometRPCClient.Start()
}

func (p *Config) GetConnContract() string {
	return p.Contracts[relayTypes.ConnectionContract]
}
