package bitcoin

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

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
