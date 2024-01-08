package wasm

import (
	"context"
	"fmt"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/client"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	ChainType string = "wasm"
)

type Provider struct {
	logger *zap.Logger
	config *ProviderConfig
	client client.IClient
	txMu   sync.Mutex
}

type ProviderConfig struct {
	ChainName string `json:"-" yaml:"-"`
	ChainID   string `json:"chain_id" yaml:"chain-id"`
	NID       string `json:"nid" yaml:"nid"`

	KeyringBackend  string `json:"keyring_backend" yaml:"keyring-backend"`
	KeyringFilePath string `json:"keyring_file_path" yaml:"keyring-file-path"`
	KeyName         string `json:"key_name" yaml:"key-name"`

	RPCUrl string `json:"rpc-url" yaml:"rpc-url"`

	ContractAddress string `json:"contract-address" yaml:"contract-address"`

	AccountPrefix string `json:"account-prefix" yaml:"account-prefix"`

	GasAdjustment float64 `json:"gas-adjustment" yaml:"gas-adjustment"`
	GasPrices     string  `json:"gas-prices" yaml:"gas-prices"`
	MinGasAmount  uint64  `json:"min-gas-amount" yaml:"min-gas-amount"`
	MaxGasAmount  uint64  `json:"max-gas-amount" yaml:"max-gas-amount"`

	BlockInterval string `json:"block_interval" yaml:"block-interval"`

	SignModeStr      string `json:"sign-mode" yaml:"sign-mode"`
	SigningAlgorithm string `json:"signing-algorithm" yaml:"signing-algorithm"`

	Debug    bool   `json:"debug"`
	HomePath string `json:"home_path"`
}

func (pc ProviderConfig) NewProvider(logger *zap.Logger, homePath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := pc.Validate(); err != nil {
		return nil, err
	}

	pc.ChainName = chainName

	cp := &Provider{
		logger: logger,
	}

	return cp, nil
}

func (pc ProviderConfig) Validate() error {
	if _, err := time.ParseDuration(pc.BlockInterval); err != nil {
		return fmt.Errorf("invalid block-interval: %w", err)
	}
	return nil
}

func (p *Provider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	return p.client.GetLatestBlock(ctx)
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error) {
	return p.client.GetTransactionReceipt(ctx, txHash)
}

func (p *Provider) NID() string {
	return p.config.NID
}

func (p *Provider) ChainName() string {
	return p.config.ChainName
}

func (p *Provider) Init(ctx context.Context) error {
	return nil
}

func (p *Provider) Type() string {
	return ChainType
}

func (p *Provider) ProviderConfig() provider.ProviderConfig {
	return *p.config
}

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfo chan types.BlockInfo) error {
	return nil
}

func (p *Provider) Route(ctx context.Context, message *types.Message, callback types.TxResponseFunc) error {
	return nil
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, message types.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, message types.Message) (bool, error) {
	return true, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key types.MessageKey) (bool, error) {

	return false, nil
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*types.Coin, error) {
	return p.client.GetBalance(ctx, addr)
}

func (p *Provider) GenerateMessage(ctx context.Context, messageKey *types.MessageKeyWithMessageHeight) (*types.Message, error) {
	return nil, nil
}

func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	return 0
}
