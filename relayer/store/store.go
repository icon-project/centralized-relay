package store

import (
	"errors"

	"github.com/syndtr/goleveldb/leveldb/iterator"
)

var (
	ErrNotFound = errors.New("key not found")
)

type Store interface {
	KeyValueReader
	KeyValueWriter
}

type KeyValueReader interface {
	GetByKey(key []byte) ([]byte, error)
	NewIterator(prefix []byte) iterator.Iterator
}

type KeyValueWriter interface {
	ClearStore() error
	SetByKey(key []byte, value []byte) error
	DeleteByKey(key []byte) error
}
