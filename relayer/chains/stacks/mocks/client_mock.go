package mocks

import (
	"context"
	"math/big"

	"github.com/icon-project/stacks-go-sdk/pkg/clarity"
	blockchainApiClient "github.com/icon-project/stacks-go-sdk/pkg/stacks_blockchain_api_client"
	"github.com/icon-project/stacks-go-sdk/pkg/transaction"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) GetAccountBalance(ctx context.Context, address string) (*big.Int, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockClient) GetAccountNonce(ctx context.Context, address string) (uint64, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockClient) GetTransactionById(ctx context.Context, id string) (*blockchainApiClient.GetTransactionById200Response, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*blockchainApiClient.GetTransactionById200Response), args.Error(1)
}

func (m *MockClient) GetBlockByHeightOrHash(ctx context.Context, height uint64) (*blockchainApiClient.GetBlocks200ResponseResultsInner, error) {
	args := m.Called(ctx, height)
	return args.Get(0).(*blockchainApiClient.GetBlocks200ResponseResultsInner), args.Error(1)
}

func (m *MockClient) GetLatestBlock(ctx context.Context) (*blockchainApiClient.GetBlocks200ResponseResultsInner, error) {
	args := m.Called(ctx)
	return args.Get(0).(*blockchainApiClient.GetBlocks200ResponseResultsInner), args.Error(1)
}

func (m *MockClient) CallReadOnlyFunction(ctx context.Context, contractAddress string, contractName string, functionName string, functionArgs []string) (*string, error) {
	args := m.Called(ctx, contractAddress, contractName, functionName, functionArgs)
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockClient) GetCurrentImplementation(ctx context.Context, contractAddress string) (clarity.ClarityValue, error) {
	args := m.Called(ctx, contractAddress)
	return args.Get(0).(clarity.ClarityValue), args.Error(1)
}

func (m *MockClient) SetAdmin(ctx context.Context, contractAddress string, newAdmin string, currentImplementation clarity.ClarityValue, senderAddress string, senderKey []byte) (string, error) {
	args := m.Called(ctx, contractAddress, newAdmin, currentImplementation, senderAddress, senderKey)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockClient) GetReceipt(ctx context.Context, contractAddress string, srcNetwork string, connSnIn *big.Int) (bool, error) {
	args := m.Called(ctx, contractAddress, srcNetwork, connSnIn)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockClient) ClaimFee(ctx context.Context, contractAddress string, senderAddress string, senderKey []byte) (string, error) {
	args := m.Called(ctx, contractAddress, senderAddress, senderKey)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockClient) SetFee(ctx context.Context, contractAddress string, networkID string, messageFee *big.Int, responseFee *big.Int, senderAddress string, senderKey []byte) (string, error) {
	args := m.Called(ctx, contractAddress, networkID, messageFee, responseFee, senderAddress, senderKey)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockClient) GetFee(ctx context.Context, contractAddress string, networkID string, responseFee bool) (uint64, error) {
	args := m.Called(ctx, contractAddress, networkID, responseFee)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockClient) SendCallMessage(ctx context.Context, contractAddress string, args []clarity.ClarityValue, senderAddress string, senderKey []byte) (string, error) {
	mockArgs := m.Called(ctx, contractAddress, args, senderAddress, senderKey)
	return mockArgs.String(0), mockArgs.Error(1)
}

func (m *MockClient) ExecuteCall(ctx context.Context, contractAddress string, args []clarity.ClarityValue, senderAddress string, senderKey []byte) (string, error) {
	mockArgs := m.Called(ctx, contractAddress, args, senderAddress, senderKey)
	return mockArgs.String(0), mockArgs.Error(1)
}

func (m *MockClient) ExecuteRollback(ctx context.Context, contractAddress string, args []clarity.ClarityValue, senderAddress string, senderKey []byte) (string, error) {
	mockArgs := m.Called(ctx, contractAddress, args, senderAddress, senderKey)
	return mockArgs.String(0), mockArgs.Error(1)
}

func (m *MockClient) MakeContractCall(
	ctx context.Context,
	contractAddress string,
	contractName string,
	functionName string,
	args []clarity.ClarityValue,
	senderAddress string,
	senderKey []byte,
) (*transaction.ContractCallTransaction, error) {
	mockArgs := m.Called(ctx, contractAddress, contractName, functionName, args, senderAddress, senderKey)
	return mockArgs.Get(0).(*transaction.ContractCallTransaction), mockArgs.Error(1)
}

func (m *MockClient) BroadcastTransaction(ctx context.Context, tx transaction.StacksTransaction) (string, error) {
	mockArgs := m.Called(ctx, tx)
	return mockArgs.String(0), mockArgs.Error(1)
}

func (m *MockClient) GetContractEvents(ctx context.Context, contractId string, limit, offset int32) (*blockchainApiClient.GetContractEventsById200Response, error) {
	args := m.Called(ctx, contractId, limit, offset)
	return args.Get(0).(*blockchainApiClient.GetContractEventsById200Response), args.Error(1)
}

func (m *MockClient) GetWebSocketURL() string {
	mockArgs := m.Called()
	return mockArgs.String(0)
}

func (m *MockClient) Log() *zap.Logger {
	args := m.Called()
	return args.Get(0).(*zap.Logger)
}
