package provider

import (
	"context"

	"github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

var messageMaxRetry = 5

type ProviderConfig interface {
	NewProvider(log *zap.Logger, homepath string, debug bool, chainName string) (ChainProvider, error)
	Validate() error
}

type ChainQuery interface {
	QueryLatestHeight(ctx context.Context) (uint64, error)
}

type ChainProvider interface {
	ChainQuery
	NID() string
	Init(ctx context.Context) error
	Listener(ctx context.Context, lastSavedHeight uint64, blockInfo chan types.BlockInfo) error
	Route(ctx context.Context, message *types.Message, callback types.TxResponseFunc) error
	ShouldReceiveMessage(ctx context.Context, messagekey types.Message) (bool, error)
	ShouldSendMessage(ctx context.Context, messageKey types.Message) (bool, error)

	QueryBalance(ctx context.Context, addr string) (*types.Coin, error)
}
