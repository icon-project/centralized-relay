package lvldb

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LVLDBSnaphot struct {
	snapshot *leveldb.Snapshot
}

func NewDBShapshot(db *LVLDB) (*LVLDBSnaphot, error) {
	if db == nil {
		return nil, fmt.Errorf("failed to create DBSnapshot: db is nil")
	}
	snapshot, err := db.SnapShot()
	if err != nil {
		return nil, fmt.Errorf("failed to create DBSnapshot: %v \n", err)
	}

	return &LVLDBSnaphot{
		snapshot: snapshot,
	}, nil

}

func (db *LVLDBSnaphot) GetByKey(key []byte) ([]byte, error) {
	return db.snapshot.Get(key, nil)
}

func (db *LVLDBSnaphot) NewIterator(prefix []byte) iterator.Iterator {
	return db.snapshot.NewIterator(util.BytesPrefix(prefix), nil)
}
