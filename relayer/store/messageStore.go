package store

import (
	"encoding/json"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/types"
)

type MessageStore struct {
	db     Store
	prefix string
}

func NewMessageStore(db Store, prefix string) *MessageStore {
	return &MessageStore{
		db:     db,
		prefix: prefix,
	}
}

func (ms *MessageStore) TotalCount() (uint64, error) {
	return ms.getCountByKey(GetKey([]string{ms.prefix}))
}

func (ms *MessageStore) TotalCountByChain(chainId string) (uint64, error) {
	return ms.getCountByKey(GetKey([]string{ms.prefix, chainId}))
}

func (ms *MessageStore) getCountByKey(key []byte) (uint64, error) {
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

func (ms *MessageStore) StoreMessage(message types.Message) error {
	key := GetKey([]string{ms.prefix,
		message.Src,
		fmt.Sprintf("%d", message.Sn),
	})

	msgByte, err := ms.Encode(message)
	if err != nil {
		return err
	}
	return ms.db.SetByKey(key, msgByte)

}

func (ms *MessageStore) GetMessage(messageKey types.MessageKey) (types.Message, error) {
	v, err := ms.db.GetByKey(GetKey([]string{ms.prefix,
		messageKey.Src,
		fmt.Sprintf("%d", messageKey.Sn),
	}))
	if err != nil {
		return types.Message{}, err
	}

	var msg types.Message
	if err := ms.Decode(v, &msg); err != nil {
		return types.Message{}, err
	}
	return msg, nil
}

func (ms *MessageStore) GetMessages(chainId string, all bool, offset int, limit int) ([]types.Message, error) {
	var messages []types.Message

	keyPrefixList := []string{ms.prefix}
	if chainId != "" {
		keyPrefixList = append(keyPrefixList, chainId)
	}
	iter := ms.db.NewIterator(GetKey(keyPrefixList))

	// return all the messages
	if all {
		for iter.Next() {
			var msg types.Message
			if err := ms.Decode(iter.Value(), &msg); err != nil {
				return nil, err
			}
			messages = append(messages, msg)
		}
		iter.Release()
		err := iter.Error()
		if err != nil {

			return nil, err
		}
		return messages, nil
	}

	// if not all use the offset logic
	for i := 0; i < int(offset); i++ {
		if !iter.Next() {
			return nil, fmt.Errorf("no message after offset")
		}
	}

	for i := 0; i < limit; i++ {
		if !iter.Next() {
			break
		}

		var msg types.Message
		if err := ms.Decode(iter.Value(), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (ms *MessageStore) DeleteMessage(messageKey types.MessageKey) error {
	return ms.db.DeleteByKey(
		GetKey([]string{ms.prefix, messageKey.Src, fmt.Sprintf("%d", messageKey.Sn)}))
}

func (ms *MessageStore) Encode(d interface{}) ([]byte, error) {
	return json.Marshal(d)
}

func (ms *MessageStore) Decode(data []byte, output interface{}) error {
	return json.Unmarshal(data, output)
}
