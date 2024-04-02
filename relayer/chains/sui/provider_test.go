package sui

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/coming-chat/go-sui/v2/account"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
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

func TestSendTransactionErrors(t *testing.T) {
	pro, err := GetSuiProvider()
	pro.wallet = &account.Account{
		Address: "0xe847098636459aa93f4da105414edca4790619b291ffdac49419f5adc19c4d21",
		KeyPair: sui_types.SuiKeyPair{},
	}
	assert.NoError(t, err)
	suiMessage := pro.NewSuiMessage([]interface{}{},
		"connectionContractAddress", "ConnectionModule", "MethodClaimFee")
	_, err = pro.SendTransaction(context.TODO(), suiMessage)
	assert.ErrorContains(t, err, "invalid packageId")

	pro.client = &mockClient{}
	pro.cfg.GasLimit = 10
	suiMessage = pro.NewSuiMessage([]interface{}{},
		"0xe847098636459aa93f4da105414edca4790619b291ffdac49419f5adc19c4d21", "ConnectionModule", "MethodClaimFee")
	_, err = pro.SendTransaction(context.TODO(), suiMessage)
	assert.ErrorContains(t, err, "gas requirement is too high")

}
