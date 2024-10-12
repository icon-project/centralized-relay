package interfaces

import (
	"context"
	"math/big"

	blockchainApiClient "github.com/icon-project/stacks-go-sdk/pkg/stacks_blockchain_api_client"
	"go.uber.org/zap"
)

type IClient interface {
	GetAccountBalance(ctx context.Context, address string) (*big.Int, error)
	GetAccountNonce(ctx context.Context, address string) (uint64, error)
	GetBlockByHeightOrHash(ctx context.Context, height uint64) (*blockchainApiClient.GetBlocks200ResponseResultsInner, error)
	GetLatestBlock(ctx context.Context) (*blockchainApiClient.GetBlocks200ResponseResultsInner, error)
	CallReadOnlyFunction(ctx context.Context, contractAddress string, contractName string, functionName string, functionArgs []string) (*string, error)
	SubscribeToEvents(ctx context.Context, eventTypes []string, callback EventCallback) error
	Log() *zap.Logger
}

type EventCallback func(eventType string, data interface{}) error
