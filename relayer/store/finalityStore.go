package store

import (
	"encoding/json"
	"fmt"

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

func (ms *FinalityStore) TotalCountByChain(chainId string) (uint64, error) {
	return ms.getCountByKey(GetKey([]string{ms.prefix, chainId}))
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

// message will be stored based on destination chainId
func (ms *FinalityStore) StoreTxObject(message *types.TransactionObject) error {
	if message == nil {
		return fmt.Errorf("error while storingMessage: message cannot be nil")
	}

	key := GetKey([]string{
		ms.prefix,
		message.Dst,
		fmt.Sprintf("%d", message.Sn),
	})

	msgByte, err := ms.Encode(message)
	if err != nil {
		return err
	}
	return ms.db.SetByKey(key, msgByte)
}

func (ms *FinalityStore) GetTxObject(messageKey *types.MessageKey) (*types.TransactionObject, error) {
	v, err := ms.db.GetByKey(GetKey([]string{
		ms.prefix,
		messageKey.Dst,
		fmt.Sprintf("%d", messageKey.Sn),
	}))
	if err != nil {
		return nil, err
	}

	var msg types.TransactionObject
	if err := ms.Decode(v, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (ms *FinalityStore) GetTxObjects(chainId string, p Pagination) ([]*types.TransactionObject, error) {
	var messages []*types.TransactionObject

	keyPrefixList := []string{ms.prefix}
	if chainId != "" {
		keyPrefixList = append(keyPrefixList, chainId)
	}
	iter := ms.db.NewIterator(GetKey(keyPrefixList))

	// return all the messages
	if p.All {
		for iter.Next() {
			var msg types.TransactionObject
			if err := ms.Decode(iter.Value(), &msg); err != nil {
				return nil, err
			}

			messages = append(messages, &msg)
		}
		iter.Release()
		err := iter.Error()
		if err != nil {
			return nil, err
		}
		return messages, nil
	}

	// if not all use the offset logic
	for i := 0; i < int(p.Offset); i++ {
		if !iter.Next() {
			return nil, fmt.Errorf("no message after offset")
		}
	}

	for i := uint64(0); i < p.Limit; i++ {
		if !iter.Next() {
			break
		}

		var msg types.TransactionObject
		if err := ms.Decode(iter.Value(), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (ms *FinalityStore) DeleteTxObject(messageKey *types.MessageKey) error {
	return ms.db.DeleteByKey(
		GetKey([]string{ms.prefix, messageKey.Dst, fmt.Sprintf("%d", messageKey.Sn)}))
}

func (ms *FinalityStore) Encode(d interface{}) ([]byte, error) {
	return json.Marshal(d)
}

func (ms *FinalityStore) Decode(data []byte, output interface{}) error {
	return json.Unmarshal(data, output)
}
