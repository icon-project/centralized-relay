package wasm

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/cometbft/cometbft/rpc/client/http"
	sdkClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/client"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"path/filepath"
	"time"
)

type ProviderConfig struct {
	RpcUrl  string `json:"rpc-url" yaml:"rpc-url"`
	GrpcUrl string `json:"grpc-url" yaml:"grpc-url"`
	ChainID string `json:"chain-id" yaml:"chain-id"`
	NID     string `json:"nid" yaml:"nid"`

	HomeDir string `json:"home-dir" yaml:"home-dir"`

	KeyringBackend string `json:"keyring-backend" yaml:"keyring-backend"`
	KeyName        string `json:"key-name" yaml:"key-name"`
	KeyringDir     string `json:"keyring-dir" yaml:"keyring-dir"`

	AccountPrefix string `json:"account-prefix" yaml:"account-prefix"`

	ContractAddress string `json:"contract-address" yaml:"contract-address"`

	Denomination string `json:"denomination" yaml:"denomination"`

	GasPrices     string  `json:"gas-prices" yaml:"gas-prices"`
	GasAdjustment float64 `json:"gas-adjustment" yaml:"gas-adjustment"`
	MinGasAmount  uint64  `json:"min-gas-amount" yaml:"min-gas-amount"`
	MaxGasAmount  uint64  `json:"max-gas-amount" yaml:"max-gas-amount"`

	BlockInterval          string `json:"block-interval" yaml:"block-interval"`
	TxConfirmationInterval string `json:"tx-wait-interval" yaml:"tx-confirmation-interval"`

	BroadcastMode string `json:"broadcast-mode" yaml:"broadcast-mode"` //sync, async and block. Recommended: sync
	SignModeStr   string `json:"sign-mode" yaml:"sign-mode"`

	Simulate bool `json:"simulate" yaml:"simulate"`

	StartHeight uint64 `json:"start-height" yaml:"start-height"`

	ChainName                  string
	FromAddress                string
	BlockIntervalTime          time.Duration
	TxConfirmationIntervalTime time.Duration
}

func (pc ProviderConfig) NewProvider(logger *zap.Logger, homePath string, _ bool, chainName string) (provider.ChainProvider, error) {
	if chainName != "" {
		pc.ChainName = chainName
	}

	if pc.HomeDir == "" {
		pc.HomeDir = homePath
	}

	if pc.KeyringDir == "" {
		pc.KeyringDir = filepath.Join(pc.HomeDir, fmt.Sprintf(".%s", pc.ChainName))
	}

	if err := pc.Validate(); err != nil {
		return nil, err
	}

	pc, err := pc.sanitize()
	if err != nil {
		return nil, err
	}

	clientContext, err := newClientContext(pc)
	if err != nil {
		return nil, err
	}

	wClient := client.New(clientContext)
	senderInfo, err := wClient.GetAccountInfo(context.Background(), clientContext.FromAddress.String())
	if err != nil {
		return nil, err
	}

	pc.FromAddress = clientContext.FromAddress.String()

	accounts := map[string]AccountInfo{
		senderInfo.GetAddress().String(): {
			AccountNumber: senderInfo.GetAccountNumber(),
			Sequence:      senderInfo.GetSequence(),
		},
	}

	return &Provider{
		logger:         logger,
		config:         pc,
		client:         wClient,
		seqTracker:     NewSeqTracker(accounts),
		memPoolTracker: &MemPoolInfo{isBlocked: false},
	}, nil
}

func (pc ProviderConfig) Validate() error {
	if _, err := time.ParseDuration(pc.BlockInterval); err != nil {
		return fmt.Errorf("invalid block-interval: %w", err)
	}

	if _, err := time.ParseDuration(pc.TxConfirmationInterval); err != nil {
		return fmt.Errorf("invalid tx-confirmation-interval: %w", err)
	}

	if pc.ChainName == "" {
		return fmt.Errorf("chain-name cannot be empty")
	}

	if pc.HomeDir == "" {
		return fmt.Errorf("home-dir cannot be empty")
	}
	return nil
}

func (pc ProviderConfig) sanitize() (ProviderConfig, error) {
	blockIntervalTime, err := time.ParseDuration(pc.BlockInterval)
	if err != nil {
		return pc, fmt.Errorf("invalid block-interval: %w", err)
	}
	pc.BlockIntervalTime = blockIntervalTime

	txConfirmationIntervalTime, err := time.ParseDuration(pc.TxConfirmationInterval)
	if err != nil {
		return pc, fmt.Errorf("invalid tx-confirmation-interval: %w", err)
	}
	pc.TxConfirmationIntervalTime = txConfirmationIntervalTime

	return pc, nil
}

func newClientContext(pc ProviderConfig) (sdkClient.Context, error) {
	clientContext := sdkClient.Context{}

	codecCfg := GetCodecConfig(pc)

	keyRing, err := keyring.New(
		pc.ChainName,
		pc.KeyringBackend,
		pc.KeyringDir,
		nil,
		codecCfg.Codec,
		func(options *keyring.Options) {
			options.SupportedAlgos = types.SupportedAlgorithms
			options.SupportedAlgosLedger = types.SupportedAlgorithmsLedger
		},
	)
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

	cometRPCClient, err := http.New(pc.RpcUrl, "/websocket")
	if err != nil {
		return clientContext, err
	}

	grpcClient, err := grpc.Dial(pc.GrpcUrl, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		return clientContext, err
	}

	clientContext = clientContext.
		WithTxConfig(codecCfg.TxConfig).
		WithInterfaceRegistry(codecCfg.InterfaceRegistry).
		WithCodec(codecCfg.Codec).
		WithFromAddress(fromAddress).
		WithFrom(keyRecord.Name).
		WithFromName(keyRecord.Name).
		WithKeyring(keyRing).
		WithKeyringDir(pc.KeyringDir).
		WithNodeURI(pc.RpcUrl).
		WithChainID(pc.ChainID).
		WithHomeDir(pc.HomeDir).
		WithBroadcastMode(pc.BroadcastMode).
		WithSignModeStr(pc.SignModeStr).
		WithSimulation(pc.Simulate).
		WithFeePayerAddress(fromAddress).
		WithFeeGranterAddress(fromAddress).
		WithClient(cometRPCClient).
		WithGRPCClient(grpcClient)

	return clientContext, nil
}
