package chains

import (
	"context"

	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/blockdb"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"

	"os"
)

const (
	DefaultNumValidators = 1
	DefaultNumFullNodes  = 1
)

type Chain interface {
	ibc.Chain
	DeployContract(ctx context.Context, keyName string) (context.Context, error)
	QueryContract(ctx context.Context, contractAddress, methodName string, params map[string]interface{}) (context.Context, error)
	ExecuteContract(ctx context.Context, contractAddress, keyName, methodName string, param map[string]interface{}) (context.Context, error)
	GetLastBlock(ctx context.Context) (context.Context, error)
	GetBlockByHeight(ctx context.Context) (context.Context, error)
	FindTxs(ctx context.Context, height uint64) ([]blockdb.Tx, error)
	GetRelayConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error)
	SetupXCall(ctx context.Context) error
	SetupConnection(ctx context.Context, target Chain) error
	FindTargetXCallMessage(ctx context.Context, target Chain, height uint64, to string) (*XCallResponse, error)
	SendPacketXCall(ctx context.Context, keyName, _to string, data, rollback []byte) (context.Context, error)
	XCall(ctx context.Context, targetChain Chain, keyName, _to string, data, rollback []byte) (*XCallResponse, error)
	ExecuteCall(ctx context.Context, reqId, data string) (context.Context, error)
	ExecuteRollback(ctx context.Context, sn string) (context.Context, error)
	FindCallMessage(ctx context.Context, startHeight uint64, from, to, sn string) (string, string, error)
	FindCallResponse(ctx context.Context, startHeight uint64, sn string) (string, error)
	GetContractAddress(key string) string
	DeployXCallMockApp(ctx context.Context, keyName string, connections []XCallConnection) error
	InitEventListener(ctx context.Context, contract string) EventListener
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
