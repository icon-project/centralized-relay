package provider

import (
	"context"
	"math/big"

	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

type ProviderConfig interface {
	NewProvider(context.Context, *zap.Logger, string, bool, string) (ChainProvider, error)
	SetWallet(string)
	GetWallet() string
	Validate() error
}

type ChainQuery interface {
	QueryLatestHeight(ctx context.Context) (uint64, error)
	QueryTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error)
}

type ChainProvider interface {
	ChainQuery
	NID() string
	ChainName() string
	Init(context.Context, string, kms.KMS) error
	Type() string
	ProviderConfig() ProviderConfig
	Listener(ctx context.Context, lastSavedHeight uint64, blockInfo chan *types.BlockInfo) error
	Route(ctx context.Context, message *types.Message, callback types.TxResponseFunc) error
	ShouldReceiveMessage(ctx context.Context, message *types.Message) (bool, error)
	ShouldSendMessage(ctx context.Context, message *types.Message) (bool, error)
	MessageReceived(ctx context.Context, key types.MessageKey) (bool, error)
	SetAdmin(context.Context, string) error

	FinalityBlock(ctx context.Context) uint64
	GenerateMessage(ctx context.Context, messageKey *types.MessageKeyWithMessageHeight) (*types.Message, error)
	QueryBalance(ctx context.Context, addr string) (*types.Coin, error)

	NewKeyStore(string, string) (string, error)
	RestoreKeyStore(context.Context, string, kms.KMS) error
	AddressFromKeyStore(string, string) (string, error)
	RevertMessage(ctx context.Context, sn *big.Int) error
	ExecuteCall(context.Context, *big.Int, []byte) ([]byte, error)
}
