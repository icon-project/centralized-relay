package lvldb

import (
	"os"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LVLDB struct {
	db *leveldb.DB
}

func NewLvlDB(path string) (*LVLDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	return &LVLDB{db: db}, errors.Wrap(err, "levelDB.OpenFile fail")
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

func (db *LVLDB) RemoveDbFile(filepath string) error {
	return errors.Wrapf(os.Remove(filepath), "unable to remove db file")
}

func (db *LVLDB) ClearStore() error {
	for {
		batch := new(leveldb.Batch)

		iter := db.db.NewIterator(nil, nil)
		keysFound := false
		for iter.Next() {
			keysFound = true
			batch.Delete(iter.Key())
		}
		iter.Release()
		if err := iter.Error(); err != nil {
			return err
		}

		if !keysFound {
			break
		}

		if err := db.db.Write(batch, nil); err != nil {
			return err
		}
	}
	return nil
}

func (db *LVLDB) Close() error {
	return db.db.Close()
}
