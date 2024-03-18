package wasm

import (
	"context"
	"crypto/tls"
	"fmt"
	"path/filepath"
	"time"

	"github.com/cometbft/cometbft/rpc/client/http"
	sdkClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	"github.com/icon-project/centralized-relay/relayer/provider"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ProviderConfig struct {
	RpcUrl  string `json:"rpc-url" yaml:"rpc-url"`
	GrpcUrl string `json:"grpc-url" yaml:"grpc-url"`
	ChainID string `json:"chain-id" yaml:"chain-id"`
	NID     string `json:"nid" yaml:"nid"`

	HomeDir string `json:"home-dir" yaml:"home-dir"`

	KeyringBackend string `json:"keyring-backend" yaml:"keyring-backend"`
	KeyringDir     string `json:"keyring-dir" yaml:"keyring-dir"`
	AccountPrefix  string `json:"account-prefix" yaml:"account-prefix"`

	Contracts providerTypes.ContractConfigMap `json:"contracts" yaml:"contracts"`
	Address   string                          `json:"address" yaml:"address"`

	Denomination string `json:"denomination" yaml:"denomination"`

	GasPrices     string  `json:"gas-prices" yaml:"gas-prices"`
	GasAdjustment float64 `json:"gas-adjustment" yaml:"gas-adjustment"`
	MinGasAmount  uint64  `json:"min-gas-amount" yaml:"min-gas-amount"`
	MaxGasAmount  uint64  `json:"max-gas-amount" yaml:"max-gas-amount"`

	BlockInterval          time.Duration `json:"block-interval" yaml:"block-interval"`
	TxConfirmationInterval time.Duration `json:"tx-confirmation-interval" yaml:"tx-confirmation-interval"`

	BroadcastMode string `json:"broadcast-mode" yaml:"broadcast-mode"` // sync, async and block. Recommended: sync
	SignModeStr   string `json:"sign-mode" yaml:"sign-mode"`

	Simulate bool `json:"simulate" yaml:"simulate"`

	StartHeight uint64 `json:"start-height" yaml:"start-height"`

	FinalityBlock uint64 `json:"finality-block" yaml:"finality-block"`

	ChainName string `json:"-" yaml:"-"`
}

func (pc *ProviderConfig) NewProvider(ctx context.Context, log *zap.Logger, homePath string, _ bool, chainName string) (provider.ChainProvider, error) {
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
		logger:    log.With(zap.Stringp("nid", &pc.NID), zap.Stringp("name", &pc.ChainName)),
		cfg:       pc,
		client:    ws,
		contracts: contracts,
		eventList: pc.GetMonitorEventFilters(contracts),
	}, nil
}

func (pc *ProviderConfig) SetWallet(addr string) {
	pc.Address = addr
}

func (pc *ProviderConfig) GetWallet() string {
	return pc.Address
}

func (pc *ProviderConfig) Validate() error {
	if pc.ChainName == "" {
		return fmt.Errorf("chain-name cannot be empty")
	}

	if pc.HomeDir == "" {
		return fmt.Errorf("home-dir cannot be empty")
	}
	return nil
}

func (pc *ProviderConfig) sanitize() (*ProviderConfig, error) {
	return pc, nil
}

func (c *ProviderConfig) newClientContext(ctx context.Context) (*sdkClient.Context, error) {
	codec := GetCodecConfig(c)

	keyRing, err := keyring.New(
		c.ChainName,
		c.KeyringBackend,
		c.KeyringDir,
		nil,
		codec.Codec,
		func(options *keyring.Options) {
			options.SupportedAlgos = types.SupportedAlgorithms
			options.SupportedAlgosLedger = types.SupportedAlgorithmsLedger
		},
	)
	if err != nil {
		return nil, err
	}

	cometRPCClient, err := http.New(c.RpcUrl, "/websocket")
	if err != nil {
		return nil, err
	}

	grpcClient, err := grpc.DialContext(ctx, c.GrpcUrl, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	if err != nil {
		return nil, err
	}

	return &sdkClient.Context{
		ChainID:           c.ChainID,
		Client:            cometRPCClient,
		NodeURI:           c.RpcUrl,
		Codec:             codec.Codec,
		Keyring:           keyRing,
		KeyringDir:        c.KeyringDir,
		TxConfig:          codec.TxConfig,
		HomeDir:           c.HomeDir,
		BroadcastMode:     c.BroadcastMode,
		SignModeStr:       c.SignModeStr,
		Simulate:          c.Simulate,
		GRPCClient:        grpcClient,
		InterfaceRegistry: codec.InterfaceRegistry,
	}, cometRPCClient.Start()
}
