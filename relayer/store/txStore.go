// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package store

type LastProcessedTxStore struct {
	db     Store
	prefix string
}

func NewLastProcessedTxStore(db Store, prefix string) *LastProcessedTxStore {
	return &LastProcessedTxStore{
		db:     db,
		prefix: prefix,
	}
}

func (s *LastProcessedTxStore) getKey(nId string) []byte {
	return GetKey([]string{s.prefix, nId})
}

func (s *LastProcessedTxStore) Set(nId string, txInfo []byte) error {
	return s.db.SetByKey(s.getKey(nId), txInfo)
}

func (s *LastProcessedTxStore) Get(nId string) ([]byte, error) {
	return s.db.GetByKey(s.getKey(nId))
}
