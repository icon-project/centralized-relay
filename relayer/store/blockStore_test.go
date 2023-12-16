package store

import (
	"testing"

	"github.com/icon-project/centralized-relay/relayer/lvldb"
	"github.com/stretchr/testify/assert"
)

var (
	testDBName = "./testdb"
)

func TestBlockStore(t *testing.T) {

	testdb, err := lvldb.NewLvlDB(testDBName, false)
	if err != nil {
		assert.Fail(t, "error while creating test db ", err)
	}

	if err := testdb.ClearStore(); err != nil {
		assert.Fail(t, "failed to clear db ", err)
	}

	nId := "icon"
	blockStore := NewBlockStore(testdb)

	key := blockStore.GetKey(nId)
	assert.Equal(t, key, []byte("block-icon"), "key computation looks good")

	saveHeight := uint64(2000)
	if err := blockStore.StoreBlock(saveHeight, nId); err != nil {
		assert.Fail(t, "error occured when storing Block ", err)
	}

	getHeight, err := blockStore.GetLastStoredBlock(nId)
	assert.NoError(t, err)
	assert.Equal(t, saveHeight, getHeight)

	replaceHeight := uint64(3000)
	if err := blockStore.StoreBlock(replaceHeight, nId); err != nil {
		assert.Fail(t, "error occured when storing Block ", err)
	}

	getHeight, err = blockStore.GetLastStoredBlock(nId)
	assert.NoError(t, err)
	assert.Equal(t, replaceHeight, getHeight)

	if err := testdb.ClearStore(); err != nil {
		assert.Fail(t, "failed to clear db ", err)
	}

}
