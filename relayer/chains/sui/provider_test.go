package sui

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/coming-chat/go-sui/v2/account"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/fardream/go-bcs/bcs"
	suitypes "github.com/icon-project/centralized-relay/relayer/chains/sui/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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
func (*mockClient) GetCheckpoints(ctx context.Context, req suitypes.SuiGetCheckpointsRequest) (*suitypes.PaginatedCheckpointsResponse, error) {
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

func TestGenerateTxDigests(t *testing.T) {
	type test struct {
		name           string
		input          []suitypes.CheckpointResponse
		expectedOutput []suitypes.TxDigests
	}

	p, err := GetSuiProvider()
	assert.NoError(t, err)

	maxDigests := 5

	tests := []test{
		{
			name: "case-1",
			input: []suitypes.CheckpointResponse{
				{
					SequenceNumber: "1",
					Transactions:   []string{"tx1", "tx2"},
				},
			},
			expectedOutput: []suitypes.TxDigests{
				{
					FromCheckpoint: 1,
					ToCheckpoint:   1,
					Digests:        []string{"tx1", "tx2"},
				},
			},
		},
		{
			name: "case-2",
			input: []suitypes.CheckpointResponse{
				{
					SequenceNumber: "1",
					Transactions:   []string{"tx1", "tx2"},
				},
				{
					SequenceNumber: "2",
					Transactions:   []string{"tx3", "tx4"},
				},
			},
			expectedOutput: []suitypes.TxDigests{
				{
					FromCheckpoint: 1,
					ToCheckpoint:   2,
					Digests:        []string{"tx1", "tx2", "tx3", "tx4"},
				},
			},
		},
		{
			name: "case-3",
			input: []suitypes.CheckpointResponse{
				{
					SequenceNumber: "1",
					Transactions:   []string{"tx1", "tx2"},
				},
				{
					SequenceNumber: "2",
					Transactions:   []string{"tx3", "tx4", "tx5", "tx6"},
				},
			},
			expectedOutput: []suitypes.TxDigests{
				{
					FromCheckpoint: 1,
					ToCheckpoint:   2,
					Digests:        []string{"tx1", "tx2", "tx3", "tx4", "tx5"},
				},
				{
					FromCheckpoint: 2,
					ToCheckpoint:   2,
					Digests:        []string{"tx6"},
				},
			},
		},
		{
			name: "case-4",
			input: []suitypes.CheckpointResponse{
				{
					SequenceNumber: "1",
					Transactions:   []string{"tx1", "tx2"},
				},
				{
					SequenceNumber: "2",
					Transactions:   []string{"tx3", "tx4", "tx5", "tx6", "tx7", "tx8", "tx9", "tx10", "tx11"},
				},
			},
			expectedOutput: []suitypes.TxDigests{
				{
					FromCheckpoint: 1,
					ToCheckpoint:   2,
					Digests:        []string{"tx1", "tx2", "tx3", "tx4", "tx5"},
				},
				{
					FromCheckpoint: 2,
					ToCheckpoint:   2,
					Digests:        []string{"tx6", "tx7", "tx8", "tx9", "tx10"},
				},
				{
					FromCheckpoint: 2,
					ToCheckpoint:   2,
					Digests:        []string{"tx11"},
				},
			},
		},
		{
			name: "case-5",
			input: []suitypes.CheckpointResponse{
				{
					SequenceNumber: "1",
					Transactions:   []string{},
				},
				{
					SequenceNumber: "2",
					Transactions:   []string{"tx1", "tx2", "tx3", "tx4", "tx5", "tx6", "tx7", "tx8", "tx9", "tx10", "tx11"},
				},
			},
			expectedOutput: []suitypes.TxDigests{
				{
					FromCheckpoint: 1,
					ToCheckpoint:   2,
					Digests:        []string{"tx1", "tx2", "tx3", "tx4", "tx5"},
				},
				{
					FromCheckpoint: 2,
					ToCheckpoint:   2,
					Digests:        []string{"tx6", "tx7", "tx8", "tx9", "tx10"},
				},
				{
					FromCheckpoint: 2,
					ToCheckpoint:   2,
					Digests:        []string{"tx11"},
				},
			},
		},
	}

	for _, eachTest := range tests {
		t.Run(eachTest.name, func(subTest *testing.T) {
			assert.Equal(subTest, eachTest.expectedOutput, p.GenerateTxDigests(eachTest.input, maxDigests))
		})
	}
}

func TestEncodeBCSByteArray(t *testing.T) {
	hexStr := "f89c8a5769746864726177546fb84a303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030323a3a7375693a3a535549b84230786566396432393635326639623236343831626662373664643931383930353736396661623134663165633363623863303464383834376664356232323364336264"
	hexByte, _ := hex.DecodeString(hexStr)

	bcsBytes, _ := bcs.Marshal(hexByte)

	hexBcs := hex.EncodeToString(bcsBytes)

	resultStrList := strings.Split(hexBcs, hexStr)

	fmt.Println("Prefix: ", resultStrList[0])

	fmt.Println("Lenght of hexbyte:", len(hexByte))
	fmt.Println("Lenght of hexbcs:", len(bcsBytes))
}
