// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package store

import (
	"encoding/json"
)

var (
	PrefixBlockStore = "block"
)

type BlockStore struct {
	dbReader KeyValueReader
	dbWriter KeyValueWriter
	prefix   string
}

func NewBlockStore(db Store) *BlockStore {
	return &BlockStore{
		dbReader: db,
		dbWriter: db,
		prefix:   PrefixBlockStore,
	}
}

func NewBlockStoreReadOnly(reader KeyValueReader) *BlockStore {
	return &BlockStore{
		dbReader: reader,
		prefix:   PrefixBlockStore,
	}
}

func (bs *BlockStore) GetKey(nId string) []byte {
	return GetKey([]string{bs.prefix, nId})
}

// StoreBlock stores block number per domainID into blockstore
func (bs *BlockStore) StoreBlock(height uint64, nId string) error {
	heightByte, err := bs.Encode(height)
	if err != nil {
		return err
	}
	return bs.dbWriter.SetByKey(bs.GetKey(nId), heightByte)
}

// GetLastStoredBlock queries the blockstore and returns latest known block
func (bs *BlockStore) GetLastStoredBlock(nId string) (uint64, error) {
	v, err := bs.dbReader.GetByKey(bs.GetKey(nId))
	if err != nil {
		return 0, err
	}
	var height uint64
	return height, json.Unmarshal(v, &height)
}

func (ms *BlockStore) Encode(d interface{}) ([]byte, error) {
	return json.Marshal(d)
}

func (ms *BlockStore) Decode(data []byte, output interface{}) error {
	return json.Unmarshal(data, output)
}
