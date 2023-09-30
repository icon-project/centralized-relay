package lvldb

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LVLDB struct {
	db *leveldb.DB

	dbMu sync.Mutex
}

func NewLvlDB(path string) (*LVLDB, error) {
	ldb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "levelDB.OpenFile fail")
	}
	return &LVLDB{db: ldb}, nil
}

func (db *LVLDB) GetByKey(key []byte) ([]byte, error) {
	return db.db.Get(key, nil)
}

func (db *LVLDB) SetByKey(key []byte, value []byte) error {
	return db.db.Put(key, value, nil)
}

func (db *LVLDB) DeleteByKey(key []byte) error {
	return db.db.Delete(key, nil)
}

func (db *LVLDB) NewIterator(prefix []byte) iterator.Iterator {
	return db.db.NewIterator(util.BytesPrefix(prefix), nil)
}

func (db *LVLDB) ClearStore() error {
	iter := db.db.NewIterator(nil, nil)
	batch := new(leveldb.Batch)

	for iter.Next() {
		key := iter.Key()
		batch.Delete(key)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil
	}
	return db.db.Write(batch, nil)
}

func (db *LVLDB) Close() error {
	return db.db.Close()
}
