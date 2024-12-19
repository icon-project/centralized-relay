package provider

import (
	"context"
	"fmt"
	"math/big"

	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

type Config interface {
	NewProvider(context.Context, *zap.Logger, string, bool, string) (ChainProvider, error)
	SetWallet(string)
	GetWallet() string
	Validate() error
	Enabled() bool
	GetConnContract() string
	ContractsAddress() types.ContractConfigMap
}

type ClusterConfig interface {
	SetClusterMode(bool)
	GetClusterMode() bool
}

type ChainQuery interface {
	QueryLatestHeight(ctx context.Context) (uint64, error)
	QueryTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error)
}

type ClusterChainProvider interface {
	SubmitClusterMessage(ctx context.Context, message *types.Message, callback types.TxResponseFunc) error
	ClusterMessageReceived(ctx context.Context, message *types.Message) (bool, error)
	ClusterMessageAcknowledged(ctx context.Context, message *types.Message) (bool, error)
}

type ClusterChainVerifier interface {
	VerifyMessage(ctx context.Context, messageKey *types.MessageKeyWithMessageHeight) ([]*types.Message, error)
}

type ChainProvider interface {
	ChainQuery
	NID() string
	Name() string
	Type() string
	Init(context.Context, string, kms.KMS) error
	Config() Config
	Listener(ctx context.Context, lastProcessedTx types.LastProcessedTx, blockInfo chan *types.BlockInfo) error
	Route(ctx context.Context, message *types.Message, callback types.TxResponseFunc) error
	ShouldReceiveMessage(ctx context.Context, message *types.Message) (bool, error)
	ShouldSendMessage(ctx context.Context, message *types.Message) (bool, error)
	SetLastSavedHeightFunc(func() uint64)
	MessageReceived(ctx context.Context, key *types.Message) (bool, error)
	SetAdmin(context.Context, string) error

	FinalityBlock(ctx context.Context) uint64
	GenerateMessages(ctx context.Context, fromHeight, toHeight uint64) ([]*types.Message, error)
	FetchTxMessages(ctx context.Context, txHash string) ([]*types.Message, error)
	QueryBalance(ctx context.Context, addr string) (*types.Coin, error)

	NewKeystore(string) (string, error)
	RestoreKeystore(context.Context) error
	ImportKeystore(context.Context, string, string) (string, error)
	RevertMessage(context.Context, *big.Int) error
	GetFee(context.Context, string, bool) (uint64, error)
	SetFee(context.Context, string, *big.Int, *big.Int) error
	ClaimFee(context.Context) error
}

// CommonConfig is the common configuration for all chain providers
type CommonConfig struct {
	ChainName     string                  `json:"-" yaml:"-"`
	RPCUrl        string                  `json:"rpc-url" yaml:"rpc-url"`
	StartHeight   uint64                  `json:"start-height" yaml:"start-height"`
	Address       string                  `json:"address" yaml:"address"`
	Contracts     types.ContractConfigMap `json:"contracts" yaml:"contracts"`
	FinalityBlock uint64                  `json:"finality-block" yaml:"finality-block"`
	NID           string                  `json:"nid" yaml:"nid"`
	Decimals      int                     `json:"decimals" yaml:"decimals"`
	HomeDir       string                  `json:"-" yaml:"-"`
	Disabled      bool                    `json:"disabled" yaml:"disabled"`
	ClusterMode   bool                    `yaml:"cluster-mode" json:"cluster-mode"`
}

// Enabled returns true if the provider is enabled
func (c *CommonConfig) Enabled() bool {
	return !c.Disabled
}

func (pc *CommonConfig) SetWallet(addr string) {
	pc.Address = addr
}

func (pc *CommonConfig) GetWallet() string {
	return pc.Address
}

func (pc *CommonConfig) Validate() error {
	if pc.ChainName == "" {
		return fmt.Errorf("chain-name cannot be empty")
	}

	if pc.HomeDir == "" {
		return fmt.Errorf("home-dir cannot be empty")
	}
	return nil
}

func (c *CommonConfig) SetClusterMode(mode bool) {
	c.ClusterMode = mode
}

func (c *CommonConfig) GetClusterMode() bool {
	return c.ClusterMode
}

func (pc *CommonConfig) ContractsAddress() types.ContractConfigMap {
	return pc.Contracts
}
