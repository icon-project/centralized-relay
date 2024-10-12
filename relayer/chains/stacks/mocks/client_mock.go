package mocks

import (
	"context"
	"math/big"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/interfaces"
	blockchainApiClient "github.com/icon-project/stacks-go-sdk/pkg/stacks_blockchain_api_client"
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

func (m *MockClient) SubscribeToEvents(ctx context.Context, eventTypes []string, callback interfaces.EventCallback) error {
	args := m.Called(ctx, eventTypes, callback)
	return args.Error(0)
}

func (m *MockClient) Log() *zap.Logger {
	args := m.Called()
	return args.Get(0).(*zap.Logger)
}
