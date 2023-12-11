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
	BuildWallets(ctx context.Context, keyName string) (ibc.Wallet, error)
	GetRelayConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error)
	SetupXCall(ctx context.Context, keyName string) error
	FindTargetXCallMessage(ctx context.Context, target Chain, height uint64, to string) (*XCallResponse, error)
	SendPacketXCall(ctx context.Context, keyName, _to string, data, rollback []byte) (context.Context, error)
	IsPacketReceived(ctx context.Context, params map[string]interface{}, order ibc.Order) bool
	XCall(ctx context.Context, targetChain Chain, keyName, _to string, data, rollback []byte) (*XCallResponse, error)
	CheckForTimeout(ctx context.Context, src Chain, params map[string]interface{}, listener EventListener) (context.Context, error)
	ExecuteCall(ctx context.Context, reqId, data string) (context.Context, error)
	ExecuteRollback(ctx context.Context, sn string) (context.Context, error)
	FindCallMessage(ctx context.Context, startHeight uint64, from, to, sn string) (string, string, error)
	FindCallResponse(ctx context.Context, startHeight uint64, sn string) (string, error)
	GetContractAddress(key string) string
	DeployXCallMockApp(ctx context.Context, keyName string, connections []XCallConnection) error
	PauseNode(context.Context) error
	UnpauseNode(context.Context) error
	InitEventListener(ctx context.Context, contract string) EventListener

	BackupConfig() ([]byte, error)
	RestoreConfig([]byte) error
	//integration test specific
	SendPacketMockDApp(ctx context.Context, targetChain Chain, keyName string, params map[string]interface{}) (PacketTransferResponse, error)
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

//
//func HexBytesToProtoUnmarshal(encoded types.HexBytes, v proto.Message) error {
//	inputBytes, err := encoded.Value()
//	if err != nil {
//		return fmt.Errorf("error unmarshalling HexByte: %s", err)
//	}
//
//	if bytes.Equal(inputBytes, make([]byte, 0)) {
//		return fmt.Errorf("encoded hexbyte is empty: %s", inputBytes)
//	}
//
//	if err := proto.Unmarshal(inputBytes, v); err != nil {
//		return err
//
//	}
//	return nil
//}
//
//func ProtoMarshalToHexBytes(v proto.Message) (string, error) {
//	value, err := proto.Marshal(v)
//	if err != nil {
//		return "", fmt.Errorf("error marshalling proto.Message: %s", err)
//	}
//	return fmt.Sprintf("0x%s", hex.EncodeToString(value)), nil
//}
