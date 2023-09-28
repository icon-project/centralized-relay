package store

import (
	"errors"

	"github.com/syndtr/goleveldb/leveldb/iterator"
)

var (
	ErrNotFound = errors.New("key not found")
)

type KeyValueReaderWriter interface {
	KeyValueReader
	KeyValueWriter
	NewIterator(prefix []byte) iterator.Iterator
	ClearStore() error
	DeleteByKey(key []byte) error
}

type KeyValueReader interface {
	GetByKey(key []byte) ([]byte, error)
}

type KeyValueWriter interface {
	SetByKey(key []byte, value []byte) error
}
