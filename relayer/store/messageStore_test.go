package store

import (
	"testing"

	"github.com/icon-project/centralized-relay/relayer/lvldb"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/stretchr/testify/assert"
)

func TestMessageStoreSet(t *testing.T) {

	testdb, err := lvldb.NewLvlDB(testDBName)
	if err != nil {
		assert.Fail(t, "error while creating test db ", err)
	}

	if err := testdb.ClearStore(); err != nil {
		assert.Fail(t, "failed to clear db ", err)
	}

	prefix := "block"
	chainId := "icon"
	Sn := uint64(1)
	messageStore := NewMessageStore(testdb, prefix)

	storeMessage := provider.RouteMessage{
		RelayMessage: provider.RelayMessage{
			Src:    chainId,
			Target: "archway",
			Sn:     Sn,
			Data:   []byte("test message"),
		},
		Retry: 2,
	}

	t.Run("store message", func(t *testing.T) {
		// storing the message
		if err := messageStore.StoreMessage(&storeMessage); err != nil {
			assert.Fail(t, "Failed to store message ", err)
		}

	})

	t.Run("getCount", func(t *testing.T) {

		// checking count
		count, err := messageStore.TotalCountByChain(chainId)
		if err != nil {
			assert.Fail(t, "failed to get message count ", err)
		}
		assert.Equal(t, count, uint64(1))

		count, err = messageStore.TotalCountByChain("archway")
		if err != nil {
			assert.Fail(t, "failed to get message count ", err)
		}
		assert.Equal(t, count, uint64(0))

	})

	t.Run("getMessage", func(t *testing.T) {
		getMessage, err := messageStore.GetMessage(chainId, Sn)
		assert.NoError(t, err, " error occured while getting message")
		assert.Equal(t, getMessage, &storeMessage)

		if err := testdb.ClearStore(); err != nil {
			assert.Fail(t, "failed to clear db ", err)
		}

		// getMessageFail case
		getMessage, err = messageStore.GetMessage("archway", Sn)
		assert.Error(t, err)

		// getMessageFail case
		getMessage, err = messageStore.GetMessage(chainId, Sn+1)
		assert.Error(t, err)

	})

	t.Run("deleteMessage", func(t *testing.T) {

		err := messageStore.DeleteMessage(chainId, Sn)
		assert.NoError(t, err)

		_, err = messageStore.GetMessage(chainId, Sn)
		assert.Error(t, err)
	})

	t.Run("GetMessages", func(t *testing.T) {

		t.Run("GetMessages empty", func(t *testing.T) {
			msg, err := messageStore.GetMessages(chainId, false, 0, 10)
			assert.NoError(t, err, "error occured when fetching messages")
			assert.Equal(t, len(msg), 0)

		})

		storeMessage1 := provider.RouteMessage{
			RelayMessage: provider.RelayMessage{
				Src:    chainId,
				Target: "archway",
				Sn:     uint64(1),
				Data:   []byte("test message"),
			},
			Retry: 2,
		}
		storeMessage2 := provider.RouteMessage{
			RelayMessage: provider.RelayMessage{
				Src:    chainId,
				Target: "archway",
				Sn:     uint64(2),
				Data:   []byte("test message"),
			},
			Retry: 2,
		}
		storeMessage3 := provider.RouteMessage{
			RelayMessage: provider.RelayMessage{
				Src:    chainId,
				Target: "archway",
				Sn:     uint64(3),
				Data:   []byte("test message"),
			},
			Retry: 2,
		}
		messageStore.StoreMessage(&storeMessage1)
		messageStore.StoreMessage(&storeMessage2)
		messageStore.StoreMessage(&storeMessage3)

		t.Run("GetMessages all", func(t *testing.T) {
			msgs, err := messageStore.GetMessages(chainId, true, 0, 0)
			assert.NoError(t, err, "error occured when fetching messages")
			assert.Equal(t, 3, len(msgs))
		})

		t.Run("GetMessages pagination by limit & offset", func(t *testing.T) {
			msgs, err := messageStore.GetMessages(chainId, false, 1, 2)
			assert.NoError(t, err, "error occured when fetching messages")
			assert.Equal(t, 2, len(msgs))
			assert.Equal(t, []*provider.RouteMessage{&storeMessage2, &storeMessage3}, msgs)
		})

		t.Run("GetMessages when offset is greater than total element", func(t *testing.T) {
			_, err := messageStore.GetMessages(chainId, false, 4, 1)
			assert.Error(t, err, "error occured when fetching messages")
		})

	})

	if err := testdb.ClearStore(); err != nil {
		assert.Fail(t, "failed to clear db ", err)
	}

}
