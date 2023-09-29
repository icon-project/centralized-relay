// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package store

import (
	"encoding/json"
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

func (bs *BlockStore) GetKey(chainId string) []byte {
	return GetKey([]string{bs.prefix, chainId})
}

// StoreBlock stores block number per domainID into blockstore
func (bs *BlockStore) StoreBlock(height uint64, chainId string) error {
	heightByte, err := bs.Encode(height)
	if err != nil {
		return err
	}
	return bs.db.SetByKey(bs.GetKey(chainId), heightByte)
}

// GetLastStoredBlock queries the blockstore and returns latest known block
func (bs *BlockStore) GetLastStoredBlock(chainId string) (uint64, error) {
	v, err := bs.db.GetByKey(bs.GetKey(chainId))
	if err != nil {
		return 0, err
	}

	var height uint64
	if err := json.Unmarshal(v, &height); err != nil {
		return 0, err
	}
	return height, nil
}

func (ms *BlockStore) Encode(d interface{}) ([]byte, error) {
	return json.Marshal(d)
}

func (ms *BlockStore) Decode(data []byte, output interface{}) error {
	return json.Unmarshal(data, output)
}
