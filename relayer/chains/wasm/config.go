package wasm

import (
	"fmt"
	"github.com/cometbft/cometbft/rpc/client/http"
	sdkClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/client"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

type ProviderConfig struct {
	ChainName       string `json:"-" yaml:"-"`
	ChainID         string `json:"chain_id" yaml:"chain-id"`
	NID             string `json:"nid" yaml:"nid"`
	KeyringBackend  string `json:"keyring_backend" yaml:"keyring-backend"`
	KeyringDir      string `json:"keyring_dir" yaml:"keyring-dir"`
	KeyName         string `json:"key_name" yaml:"key-name"`
	Codec           codec.Codec
	NodeURI         string  `json:"node_uri" yaml:"node-uri"`
	ContractAddress string  `json:"contract-address" yaml:"contract-address"`
	AccountPrefix   string  `json:"account-prefix" yaml:"account-prefix"`
	GasAdjustment   float64 `json:"gas-adjustment" yaml:"gas-adjustment"`
	GasPrices       string  `json:"gas-prices" yaml:"gas-prices"`
	MinGasAmount    uint64  `json:"min-gas-amount" yaml:"min-gas-amount"`
	MaxGasAmount    uint64  `json:"max-gas-amount" yaml:"max-gas-amount"`
	BlockInterval   string  `json:"block_interval" yaml:"block-interval"`
	BroadcastMode   string  `json:"broadcast_mode" json:"broadcast-mode"`
	SignModeStr     string  `json:"sign-mode" yaml:"sign-mode"`
	SkipConfirm     bool    `json:"skip_confirm" yaml:"skip-confirm"`
	Simulate        bool    `json:"simulate" yaml:"simulate"`
	Debug           bool    `json:"debug"`
	HomePath        string  `json:"home_path"`
}

func (pc ProviderConfig) NewProvider(logger *zap.Logger, homePath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := pc.Validate(); err != nil {
		return nil, err
	}

	pc.ChainName = chainName
	pc.HomePath = homePath
	pc.Debug = debug

	clientContext, err := newClientContext(pc)
	if err != nil {
		return nil, err
	}

	cp := &Provider{
		logger: logger.With(zap.String("module", "wasm")),
		config: pc,
		client: client.New(clientContext),
	}

	return cp, nil
}

func (pc ProviderConfig) Validate() error {
	if _, err := time.ParseDuration(pc.BlockInterval); err != nil {
		return fmt.Errorf("invalid block-interval: %w", err)
	}
	return nil
}

func newClientContext(pc ProviderConfig) (sdkClient.Context, error) {
	clientContext := sdkClient.Context{}

	keyRing, err := keyring.New("myApp", pc.KeyringBackend, pc.KeyringDir, nil, pc.Codec)
	if err != nil {
		return clientContext, err
	}

	keyRecord, err := keyRing.Key(pc.KeyName)
	if err != nil {
		return clientContext, err
	}
	fromAddress, err := keyRecord.GetAddress()
	if err != nil {
		return clientContext, err
	}

	cometRPCClient, err := http.New(pc.NodeURI, "/websocket")
	if err != nil {
		return clientContext, err
	}

	grpcClient, err := grpc.Dial(pc.NodeURI) //Todo use secured rpc channel for production
	if err != nil {
		return clientContext, err
	}

	clientContext.FromAddress = fromAddress
	clientContext.From = keyRecord.Name
	clientContext.FromName = keyRecord.Name
	clientContext.Keyring = keyRing

	clientContext.Client = cometRPCClient
	clientContext.GRPCClient = grpcClient

	clientContext.KeyringDir = pc.KeyringDir
	clientContext.NodeURI = pc.NodeURI
	clientContext.ChainID = pc.ChainID

	clientContext.HomeDir = pc.HomePath

	clientContext.BroadcastMode = pc.BroadcastMode
	clientContext.SignModeStr = pc.SignModeStr
	clientContext.SkipConfirm = pc.SkipConfirm
	clientContext.Simulate = pc.Simulate

	clientContext.FeePayer = fromAddress
	clientContext.FeeGranter = fromAddress

	return clientContext, nil
}
