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

	KeyName        string `json:"key-name" yaml:"key-name"`
	KeyringBackend string `json:"keyring-backend" yaml:"keyring-backend"`
	AccountPrefix  string `json:"account-prefix" yaml:"account-prefix"`

	Contracts providerTypes.ContractConfigMap `json:"contracts" yaml:"contracts"`

	Denomination string `json:"denomination" yaml:"denomination"`

	GasPrices     string  `json:"gas-prices" yaml:"gas-prices"`
	GasAdjustment float64 `json:"gas-adjustment" yaml:"gas-adjustment"`
	MinGasAmount  uint64  `json:"min-gas-amount" yaml:"min-gas-amount"`
	MaxGasAmount  uint64  `json:"max-gas-amount" yaml:"max-gas-amount"`

	BlockInterval          string `json:"block-interval" yaml:"block-interval"`
	TxConfirmationInterval string `json:"tx-wait-interval" yaml:"tx-confirmation-interval"`

	BroadcastMode string `json:"broadcast-mode" yaml:"broadcast-mode"` // sync, async and block. Recommended: sync
	SignModeStr   string `json:"sign-mode" yaml:"sign-mode"`

	Simulate bool `json:"simulate" yaml:"simulate"`

	StartHeight uint64 `json:"start-height" yaml:"start-height"`

	ChainName                  string
	FromAddress                string
	BlockIntervalTime          time.Duration
	TxConfirmationIntervalTime time.Duration
}

func (pc *ProviderConfig) NewProvider(ctx context.Context, log *zap.Logger, homePath string, _ bool, chainName string) (provider.ChainProvider, error) {
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

	clientContext, err := pc.newClientContext()
	if err != nil {
		return nil, err
	}

	wClient := newClient(clientContext)
	senderInfo, err := wClient.GetAccountInfo(ctx, clientContext.FromAddress.String())
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
		logger:         log.With(zap.String("nid", pc.NID), zap.String("chain", pc.ChainName)),
		cfg:            pc,
		client:         wClient,
		seqTracker:     NewSeqTracker(accounts),
		memPoolTracker: &MemPoolInfo{isBlocked: false},
		contracts:      pc.eventMap(),
	}, nil
}

func (c *Provider) GetWallet() string {
	return c.cfg.FromAddress
}

func (pc ProviderConfig) SetWallet(addr string) {
	// Todo set wallet
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

func (pc *ProviderConfig) sanitize() (*ProviderConfig, error) {
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

func (pc *ProviderConfig) newClientContext() (*sdkClient.Context, error) {
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
		return nil, err
	}

	keyRecord, err := keyRing.Key(pc.KeyName)
	if err != nil {
		return nil, err
	}

	fromAddress, err := keyRecord.GetAddress()
	if err != nil {
		return nil, err
	}

	cometRPCClient, err := http.New(pc.RpcUrl, "/websocket")
	if err != nil {
		return nil, err
	}

	grpcClient, err := grpc.Dial(pc.GrpcUrl, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		return nil, err
	}

	return &sdkClient.Context{
		ChainID:           pc.ChainID,
		Client:            cometRPCClient,
		NodeURI:           pc.RpcUrl,
		Codec:             codecCfg.Codec,
		From:              keyRecord.Name,
		FromName:          keyRecord.Name,
		FromAddress:       fromAddress,
		Keyring:           keyRing,
		KeyringDir:        pc.KeyringDir,
		TxConfig:          codecCfg.TxConfig,
		HomeDir:           pc.HomeDir,
		BroadcastMode:     pc.BroadcastMode,
		SignModeStr:       pc.SignModeStr,
		Simulate:          pc.Simulate,
		FeePayer:          fromAddress,
		FeeGranter:        fromAddress,
		GRPCClient:        grpcClient,
		InterfaceRegistry: codecCfg.InterfaceRegistry,
	}, nil
}
