package bitcoin

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/multisig"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

func TestDecodeWithdrawToMessage(t *testing.T) {
	input := "+QEdAbkBGfkBFrMweDIuaWNvbi9jeGZjODZlZTc2ODdlMWJmNjgxYjU1NDhiMjY2Nzg0NDQ4NWMwZTcxOTK4PnRiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1gh6FAbhS+FCKV2l0aGRyYXdUb4MwOjC4PnRiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1ZPhIuEYweDIuYnRjL3RiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1"
	// Decode base64
	decodedBytes, _ := base64.StdEncoding.DecodeString(input)

	result, data, err := decodeWithdrawToMessage(decodedBytes)

	fmt.Println("data", data)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su", result.To)
	assert.Equal(t, big.NewInt(100).Bytes(), result.Amount)
	assert.Equal(t, "WithdrawTo", result.Action)
	assert.Equal(t, "0:0", result.TokenAddress)
}

func TestCreateBitcoinMultisigTx(t *testing.T) {
	data := "+QEdAbkBGfkBFrMweDIuaWNvbi9jeGZjODZlZTc2ODdlMWJmNjgxYjU1NDhiMjY2Nzg0NDQ4NWMwZTcxOTK4PnRiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1gh6FAbhS+FCKV2l0aGRyYXdUb4MwOjC4PnRiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1ZPhIuEYweDIuYnRjL3RiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1"
	// Decode base64
	dataBytes, _ := base64.StdEncoding.DecodeString(data)

	chainParam := &chaincfg.TestNet3Params
	_, relayersMultisigInfo := // multisig.RandomMultisigInfo(3, 3, chainParam, []int{0, 1, 2}, 0, 1)
	multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	relayersMultisigWallet, _ := multisig.BuildMultisigWallet(relayersMultisigInfo)

	_, _, hexRawTx, _, err := CreateBitcoinMultisigTx(dataBytes, 5000, relayersMultisigWallet, chainParam, UNISAT_DEFAULT_TESTNET)
	fmt.Println("err: ", err)
	fmt.Println("hexRawTx: ", hexRawTx)
}

func TestProvider_Route(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "bitcoin_provider_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := leveldb.OpenFile(dbPath, nil)
	assert.NoError(t, err)
	defer db.Close()

	logger, _ := zap.NewDevelopment()
	provider := &Provider{
		logger: logger,
		cfg:    &Config{Mode: SlaveMode},
		db:     db,
	}

	// Create a test message
	testMessage := &types.Message{
		Dst:           "destination",
		Src:           "source",
		Sn:            big.NewInt(123),
		Data:          []byte("test data"),
		MessageHeight: 456,
		EventType:     events.EmitMessage,
	}

	// Test storing the message
	err = provider.Route(context.Background(), testMessage, nil)
	assert.NoError(t, err)

	// Test retrieving the message
	key := []byte(fmt.Sprintf("bitcoin_message_%s", testMessage.Sn.String()))
	storedData, err := db.Get(key, nil)
	assert.NoError(t, err)

	var retrievedMessage types.Message
	err = json.Unmarshal(storedData, &retrievedMessage)
	assert.NoError(t, err)

	assert.Equal(t, testMessage.Dst, retrievedMessage.Dst)
	assert.Equal(t, testMessage.Src, retrievedMessage.Src)
	assert.Equal(t, testMessage.Sn.String(), retrievedMessage.Sn.String())
	assert.Equal(t, testMessage.Data, retrievedMessage.Data)
	assert.Equal(t, testMessage.MessageHeight, retrievedMessage.MessageHeight)
	assert.Equal(t, testMessage.EventType, retrievedMessage.EventType)

	// Test deleting the message
	err = db.Delete(key, nil)
	assert.NoError(t, err)

	_, err = db.Get(key, nil)
	assert.Error(t, err) // Should return an error as the key no longer exists
}
