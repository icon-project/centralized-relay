package store

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"

	"github.com/icon-project/centralized-relay/relayer/types"
)

type FinalityStore struct {
	db     Store
	prefix string
}

func NewFinalityStore(db Store, prefix string) *FinalityStore {
	return &FinalityStore{
		db:     db,
		prefix: prefix,
	}
}

func (ms *FinalityStore) TotalCount() (uint64, error) {
	return ms.getCountByKey(GetKey([]string{ms.prefix}))
}

func (ms *FinalityStore) TotalCountByChain(nId string) (uint64, error) {
	return ms.getCountByKey(GetKey([]string{ms.prefix, nId}))
}

func (ms *FinalityStore) getCountByKey(key []byte) (uint64, error) {
	iter := ms.db.NewIterator(key)
	count := 0
	for iter.Next() {
		count++
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return 0, err
	}
	return uint64(count), nil
}

// message will be stored based on destination nId
func (ms *FinalityStore) StoreTxObject(message *types.TransactionObject) error {
	if message == nil {
		return fmt.Errorf("error while storingMessage: message cannot be nil")
	}

	key := GetKey([]string{ms.prefix, message.Dst, message.Sn.String()})

	msgByte, err := ms.Encode(message)
	if err != nil {
		return err
	}
	return ms.db.SetByKey(key, msgByte)
}

func (ms *FinalityStore) GetTxObject(messageKey *types.MessageKey) (*types.TransactionObject, error) {
	v, err := ms.db.GetByKey(GetKey([]string{ms.prefix, messageKey.Dst, messageKey.Sn.String()}))
	if err != nil {
		return nil, err
	}

	var msg types.TransactionObject
	if err := ms.Decode(v, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (ms *FinalityStore) GetTxObjects(nId string, p *Pagination) ([]*types.TransactionObject, error) {
	var messages []*types.TransactionObject

	iter := ms.db.NewIterator(GetKey([]string{ms.prefix, nId}))
	defer iter.Release()

	if p.All {
		for iter.Next() {
			var msg types.TransactionObject
			if err := ms.Decode(iter.Value(), &msg); err != nil {
				return nil, err
			}
			messages = append(messages, &msg)
		}
	}

	for i := uint(0); i < p.Limit && iter.Next(); i++ {
		var msg types.TransactionObject
		if err := ms.Decode(iter.Value(), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, iter.Error()
}

func (ms *FinalityStore) DeleteTxObject(messageKey *types.MessageKey) error {
	return ms.db.DeleteByKey(GetKey([]string{ms.prefix, messageKey.Dst, messageKey.Sn.String()}))
}

func (ms *FinalityStore) Encode(d interface{}) ([]byte, error) {
	return jsoniter.Marshal(d)
}

func (ms *FinalityStore) Decode(data []byte, output interface{}) error {
	return jsoniter.Unmarshal(data, output)
}
