package interfaces

import (
	"context"
	"math/big"

	"github.com/icon-project/stacks-go-sdk/pkg/clarity"
	blockchainApiClient "github.com/icon-project/stacks-go-sdk/pkg/stacks_blockchain_api_client"
	"go.uber.org/zap"
)

type IClient interface {
	GetAccountBalance(ctx context.Context, address string) (*big.Int, error)
	GetAccountNonce(ctx context.Context, address string) (uint64, error)
	GetTransactionById(ctx context.Context, id string) (*blockchainApiClient.GetTransactionById200Response, error)
	GetBlockByHeightOrHash(ctx context.Context, height uint64) (*blockchainApiClient.GetBlocks200ResponseResultsInner, error)
	GetCurrentImplementation(ctx context.Context, contractAddress string) (clarity.ClarityValue, error)
	SetAdmin(ctx context.Context, contractAddress string, newAdmin string, currentImplementation clarity.ClarityValue, senderAddress string, senderKey []byte) (string, error)
	GetReceipt(ctx context.Context, contractAddress string, srcNetwork string, connSnIn *big.Int) (bool, error)
	ClaimFee(ctx context.Context, contractAddress string, senderAddress string, senderKey []byte) (string, error)
	SetFee(ctx context.Context, contractAddress string, networkID string, messageFee *big.Int, responseFee *big.Int, senderAddress string, senderKey []byte) (string, error)
	GetFee(ctx context.Context, contractAddress string, networkID string, responseFee bool) (uint64, error)
	SendCallMessage(ctx context.Context, contractAddress string, args []clarity.ClarityValue, senderAddress string, senderKey []byte) (string, error)
	ExecuteCall(ctx context.Context, contractAddress string, args []clarity.ClarityValue, senderAddress string, senderKey []byte) (string, error)
	ExecuteRollback(ctx context.Context, contractAddress string, args []clarity.ClarityValue, senderAddress string, senderKey []byte) (string, error)
	GetLatestBlock(ctx context.Context) (*blockchainApiClient.GetBlocks200ResponseResultsInner, error)
	CallReadOnlyFunction(ctx context.Context, contractAddress string, contractName string, functionName string, functionArgs []string) (*string, error)
	SubscribeToEvents(ctx context.Context, eventTypes []string, callback EventCallback) error
	Log() *zap.Logger
}

type EventCallback func(eventType string, data interface{}) error
