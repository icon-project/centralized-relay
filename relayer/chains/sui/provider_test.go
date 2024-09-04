package sui

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/coming-chat/go-sui/v2/account"
	suisdkClient "github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	suitypes "github.com/icon-project/centralized-relay/relayer/chains/sui/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	privateKeyEncodedWithFlag = "ALWS4mKTtggWc8gH+a5bFLFQ0AeNbZpUdDI//3OpAVys"
	expectedDecodedAddr       = "0xe847098636459aa93f4da105414edca4790619b291ffdac49419f5adc19c4d21"
	expectedDecodedPrivKey    = "b592e26293b6081673c807f9ae5b14b150d0078d6d9a5474323fff73a9015cac"
)

type mockKms struct {
}

func (*mockKms) Init(context.Context) (*string, error) {
	initcompleted := "yes"
	return &initcompleted, nil
}

func (*mockKms) Encrypt(ctx context.Context, input []byte) ([]byte, error) {
	return input, nil
}

func (*mockKms) Decrypt(ctx context.Context, input []byte) ([]byte, error) {
	return input, nil
}

type mockClient struct {
}

func (*mockClient) GetLatestCheckpointSeq(ctx context.Context) (uint64, error) {
	panic("not implemented")
}
func (*mockClient) GetTotalBalance(ctx context.Context, addr string) (uint64, error) {
	panic("not implemented")
}
func (*mockClient) EstimateGas(ctx context.Context, txBytes lib.Base64Data) (*types.DryRunTransactionBlockResponse, int64, error) {
	return &types.DryRunTransactionBlockResponse{
		Effects: lib.TagJson[types.SuiTransactionBlockEffects]{
			Data: types.SuiTransactionBlockEffects{
				V1: &types.SuiTransactionBlockEffectsV1{
					Status: types.ExecutionStatus{Status: "success"},
				},
			},
		},
	}, 100, nil
}
func (*mockClient) ExecuteContract(ctx context.Context, suiMessage *SuiMessage, address string, gasBudget uint64) (*types.TransactionBytes, error) {
	return &types.TransactionBytes{
		TxBytes: []byte("txbytes"),
	}, nil
}
func (*mockClient) CommitTx(ctx context.Context, wallet *account.Account, txBytes lib.Base64Data, signatures []any) (*types.SuiTransactionBlockResponse, error) {
	return &types.SuiTransactionBlockResponse{}, nil
}
func (*mockClient) GetTransaction(ctx context.Context, txDigest string) (*types.SuiTransactionBlockResponse, error) {
	panic("not implemented")
}
func (*mockClient) QueryContract(ctx context.Context, suiMessage *SuiMessage, address string, gasBudget uint64) (any, error) {
	panic("not implemented")
}

func (*mockClient) GetEventsFromTxBlocks(ctx context.Context, packageID string, digests []string) ([]suitypes.EventResponse, error) {
	panic("not implemented")
}

func GetSuiProvider() (*Provider, error) {
	pc := Config{
		NID:     "sui.testnet",
		Address: "0xe847098636459aa93f4da105414edca4790619b291ffdac49419f5adc19c4d21",
		RPCUrl:  "https://sui-devnet-endpoint.blockvision.org",
		HomeDir: "./tmp/",
	}
	logger := zap.NewNop()
	ctx := context.Background()
	prov, err := pc.NewProvider(ctx, logger, "./tmp/", true, "sui.testnet")
	if err != nil {
		return nil, err
	}

	suiProvider, ok := prov.(*Provider)
	if !ok {
		return nil, fmt.Errorf("unable to type cast to sui chain provider")
	}
	suiProvider.Init(ctx, "./tmp/", &mockKms{})
	err = os.MkdirAll(suiProvider.keystorePath("test"), 0777)
	if err != nil {
		fmt.Println(err)
	}
	return suiProvider, nil
}

func TestNewKeystore(t *testing.T) {
	pro, err := GetSuiProvider()
	assert.NoError(t, err)
	generatedKeyStore, err := pro.NewKeystore("password")
	assert.NoError(t, err)
	pro.cfg.Address = generatedKeyStore
	err = pro.RestoreKeystore(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, generatedKeyStore, pro.wallet.Address)
}

func TestImportKeystore(t *testing.T) {
	pro, err := GetSuiProvider()
	assert.NoError(t, err)
	data := []byte("[\"" + privateKeyEncodedWithFlag + "\"]")
	os.WriteFile("./tmp/ks.keystore", data, 0644)
	restoredKeyStore, err := pro.ImportKeystore(context.TODO(), "./tmp/ks.keystore", "passphrase")
	assert.NoError(t, err)
	pro.cfg.Address = restoredKeyStore
	err = pro.RestoreKeystore(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedDecodedAddr, pro.wallet.Address)
	assert.Equal(t, expectedDecodedPrivKey, hex.EncodeToString(pro.wallet.KeyPair.PrivateKey()[:32]))
}

func newRootLogger() *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.UTC().Format("2006-01-02T15:04:05.000000Z07:00"))
	}
	config.LevelKey = "lvl"

	enc := zapcore.NewJSONEncoder(config)
	level := zap.InfoLevel

	core := zapcore.NewTee(zapcore.NewCore(enc, os.Stderr, level))

	return zap.New(core)
}
func TestQueryTxBlocks(t *testing.T) {
	rpcClient, err := suisdkClient.Dial("https://fullnode.testnet.sui.io:443")
	assert.NoError(t, err)

	client := NewClient(rpcClient, newRootLogger())

	inputObj, err := sui_types.NewObjectIdFromHex("0xde9b85d02710e2651f530e83ba2fb1daea45fd05e72f1c0dd2398239b917b46c")
	assert.NoError(t, err)

	query := types.SuiTransactionBlockResponseQuery{
		Filter: &types.TransactionFilter{
			InputObject: inputObj,
		},
		Options: &types.SuiTransactionBlockResponseOptions{
			ShowEvents: true,
		},
	}

	cursor, err := sui_types.NewDigest("HGjoM1cL5dzR7B6radZ1woPFESe2cvwAaN6z46rBN9c5")
	assert.NoError(t, err)

	limit := uint(100)
	res, err := client.QueryTxBlocks(context.Background(), query, cursor, &limit, false)
	assert.NoError(t, err)

	count := 0
	for _, blockRes := range res.Data {
		checkpoint := ""
		if blockRes.Checkpoint != nil {
			checkpoint = blockRes.Checkpoint.Decimal().String()
		}
		for _, ev := range blockRes.Events {
			count++
			client.log.Info("event",
				zap.String("checkpoint", checkpoint),
				zap.String("package-id", ev.PackageId.String()),
				zap.String("module", ev.TransactionModule),
				zap.String("event_type", ev.Type),
				zap.String("tx-digest", ev.Id.TxDigest.String()),
			)
		}
	}

	fmt.Println("Total Events: ", count)
}
