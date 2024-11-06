package chains

import (
	"context"

	"os"
)

const (
	DefaultNumValidators = 1
	DefaultNumFullNodes  = 1
)

type Chain interface {
	// ibc.Chain
	Height(ctx context.Context) (uint64, error)
	Config() ChainConfig
	GetRelayConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error)
	SetupXCall(ctx context.Context) error
	SetupConnection(ctx context.Context, target Chain) error
	FindTargetXCallMessage(ctx context.Context, target Chain, height uint64, to string) (*XCallResponse, error)
	SendPacketXCall(ctx context.Context, keyName, _to string, data, rollback []byte) (context.Context, error)
	XCall(ctx context.Context, targetChain Chain, keyName, _to string, data, rollback []byte) (*XCallResponse, error)
	FindCallMessage(ctx context.Context, startHeight uint64, from, to, sn string) (string, string, error)
	FindCallResponse(ctx context.Context, startHeight uint64, sn string) (string, error)
	FindRollbackExecutedMessage(ctx context.Context, startHeight uint64, sn string) (string, error)
	GetContractAddress(key string) string
	DeployXCallMockApp(ctx context.Context, keyName string, connections []XCallConnection) error
	DeployNSetupClusterContracts(ctx context.Context, targetChains []Chain) error
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
