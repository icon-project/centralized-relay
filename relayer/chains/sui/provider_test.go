package sui

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

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
		return nil, fmt.Errorf("unbale to type case to icon chain provider")
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
	assert.Equal(t, expectedDecodedPrivKey, hex.EncodeToString(pro.wallet.PrivateKey))
}
