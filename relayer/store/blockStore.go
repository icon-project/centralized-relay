// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package store

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/syndtr/goleveldb/leveldb"
)

type BlockStore struct {
	db     Store
	prefix string
}

func NewBlockStore(db Store, prefix string) *BlockStore {
	return &BlockStore{
		db:     db,
		prefix: prefix,
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
	return bs.db.SetByKey(bs.GetKey(nId), heightByte)
}

// GetLastStoredBlock queries the blockstore and returns latest known block
func (bs *BlockStore) GetLastStoredBlock(nId string) (uint64, error) {
	v, err := bs.db.GetByKey(bs.GetKey(nId))
	if err != nil {
		if err == leveldb.ErrNotFound {
			return 0, nil
		}
		return 0, err
	}
	var height uint64
	return height, jsoniter.Unmarshal(v, &height)
}

func (ms *BlockStore) Encode(d interface{}) ([]byte, error) {
	return jsoniter.Marshal(d)
}

func (ms *BlockStore) Decode(data []byte, output interface{}) error {
	return jsoniter.Unmarshal(data, output)
}
