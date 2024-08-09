package store

import (
	"math/big"
	"os"
	"testing"

	"github.com/icon-project/centralized-relay/relayer/lvldb"
	"github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stretchr/testify/assert"
)

func TestMessageStoreSet(t *testing.T) {
	tmpDir := os.TempDir()
	testdb, err := lvldb.NewLvlDB(tmpDir + "/tt")
	if err != nil {
		assert.Fail(t, "error while creating test db ", err)
	}

	if err := testdb.ClearStore(); err != nil {
		assert.Fail(t, "failed to clear db ", err)
	}

	prefix := "block"
	nId := "icon"
	Sn := big.NewInt(1)
	messageStore := NewMessageStore(testdb, prefix)

	storeMessage := &types.Message{
		Src:  nId,
		Dst:  "archway",
		Sn:   Sn,
		Data: []byte("test message"),
	}

	t.Run("store message", func(t *testing.T) {
		// storing the message

		if err := messageStore.StoreMessage(types.NewRouteMessage(storeMessage)); err != nil {
			assert.Fail(t, "Failed to store message ", err)
		}
	})

	t.Run("getCount", func(t *testing.T) {
		// checking count
		count, err := messageStore.TotalCountByChain(nId)
		if err != nil {
			assert.Fail(t, "failed to get message count ", err)
		}
		assert.Equal(t, count, uint(1))

		count, err = messageStore.TotalCountByChain("archway")
		if err != nil {
			assert.Fail(t, "failed to get message count ", err)
		}
		assert.Equal(t, count, uint(0))
	})

	t.Run("getMessage", func(t *testing.T) {
		getMessage, err := messageStore.GetMessage(types.NewMessageKey(Sn, nId, "", "emitMessage"))
		assert.NoError(t, err, " error occured while getting message")
		assert.Equal(t, getMessage, types.NewRouteMessage(storeMessage))

		if err := testdb.ClearStore(); err != nil {
			assert.Fail(t, "failed to clear db ", err)
		}

		// getMessageFail case
		getMessage, err = messageStore.GetMessage(types.NewMessageKey(Sn, "archway", "", "emitMessage"))
		assert.Error(t, err)

		// getMessageFail case
		newSn := big.NewInt(0)
		newSn.Add(Sn, big.NewInt(1))
		getMessage, err = messageStore.GetMessage(types.NewMessageKey(newSn, "archway", "", "emitMessage"))
		assert.Error(t, err)
	})

	t.Run("deleteMessage", func(t *testing.T) {
		err := messageStore.DeleteMessage(types.NewMessageKey(Sn, nId, "", "emitMessage"))
		assert.NoError(t, err)

		_, err = messageStore.GetMessage(types.NewMessageKey(Sn, nId, "", "emitMessage"))
		assert.Error(t, err)
	})

	t.Run("GetMessages", func(t *testing.T) {
		t.Run("GetMessages empty", func(t *testing.T) {
			p := NewPagination().
				WithLimit(10).
				WithOffset(0)
			msg, err := messageStore.GetMessages(nId, p)
			assert.NoError(t, err, "error occured when fetching messages")
			assert.Equal(t, len(msg), 0)
		})

		storeMessage1 := &types.Message{
			Src:  nId,
			Dst:  "archway",
			Sn:   big.NewInt(1),
			Data: []byte("test message"),
		}
		storeMessage2 := &types.Message{
			Src:  nId,
			Dst:  "archway",
			Sn:   big.NewInt(2),
			Data: []byte("test message"),
		}
		storeMessage3 := &types.Message{
			Src:  nId,
			Dst:  "archway",
			Sn:   big.NewInt(3),
			Data: []byte("test message"),
		}
		messageStore.StoreMessage(types.NewRouteMessage(storeMessage1))
		messageStore.StoreMessage(types.NewRouteMessage(storeMessage2))
		messageStore.StoreMessage(types.NewRouteMessage(storeMessage3))

		t.Run("GetMessages all", func(t *testing.T) {
			p := NewPagination().GetAll()
			msgs, err := messageStore.GetMessages(nId, p)
			assert.NoError(t, err, "error occured when fetching messages")
			assert.Equal(t, 3, len(msgs))
		})

		t.Run("GetMessages pagination by limit & offset", func(t *testing.T) {
			p := NewPagination().
				WithLimit(2).
				WithOffset(1)
			msgs, err := messageStore.GetMessages(nId, p)
			assert.NoError(t, err, "error occured when fetching messages")
			assert.Equal(t, 2, len(msgs))
			assert.Equal(t, storeMessage1.Sn, msgs[0].Sn)
			assert.Equal(t, storeMessage2.Sn, msgs[1].Sn)
		})

		t.Run("GetMessages when offset is greater than total element", func(t *testing.T) {
			p := NewPagination().
				WithLimit(1).
				WithOffset(14)
			_, err := messageStore.GetMessages(nId, p)
			// TODO  fix the issue with offset
			assert.NoError(t, err, "error occured when fetching messages")
		})
	})

	if err := testdb.ClearStore(); err != nil {
		assert.Fail(t, "failed to clear db ", err)
	}
}
