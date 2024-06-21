package provider

import (
	"context"
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
}

type ChainQuery interface {
	QueryLatestHeight(ctx context.Context) (uint64, error)
	QueryTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error)
}

type ChainProvider interface {
	ChainQuery
	NID() string
	Name() string
	Init(context.Context, string, kms.KMS) error
	Type() string
	Config() Config
	Listener(ctx context.Context, lastProcessedTxInfo interface{}, blockInfo chan *types.BlockInfo) error
	Route(ctx context.Context, message *types.Message, callback types.TxResponseFunc) error
	ShouldReceiveMessage(ctx context.Context, message *types.Message) (bool, error)
	ShouldSendMessage(ctx context.Context, message *types.Message) (bool, error)
	MessageReceived(ctx context.Context, key *types.MessageKey) (bool, error)
	SetAdmin(context.Context, string) error

	FinalityBlock(ctx context.Context) uint64
	GenerateMessages(ctx context.Context, messageKey *types.MessageKeyWithMessageHeight) ([]*types.Message, error)
	QueryBalance(ctx context.Context, addr string) (*types.Coin, error)

	NewKeystore(string) (string, error)
	RestoreKeystore(context.Context) error
	ImportKeystore(context.Context, string, string) (string, error)
	RevertMessage(context.Context, *big.Int) error
	GetFee(context.Context, string, bool) (uint64, error)
	SetFee(context.Context, string, uint64, uint64) error
	ClaimFee(context.Context) error
}
