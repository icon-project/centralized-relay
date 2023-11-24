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

type Pagination struct {
	Limit  uint64
	Offset uint64
	All    bool
}

func NewPagination() Pagination {
	return Pagination{}
}
func (p Pagination) GetAll() Pagination {
	p.All = true
	return p
}

func (p Pagination) WithLimit(l uint64) Pagination {
	p.Limit = l
	return p
}

func (p Pagination) WithOffset(o uint64) Pagination {
	p.Offset = o
	return p
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

func (ms *MessageStore) StoreMessage(message *types.RouteMessage) error {

	if message == nil {
		return fmt.Errorf("error while storingMessage: message cannot be nil")
	}

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

func (ms *MessageStore) GetMessage(messageKey types.MessageKey) (*types.RouteMessage, error) {
	v, err := ms.db.GetByKey(GetKey([]string{ms.prefix,
		messageKey.Src,
		fmt.Sprintf("%d", messageKey.Sn),
	}))
	if err != nil {
		return nil, err
	}

	var msg types.RouteMessage
	if err := ms.Decode(v, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (ms *MessageStore) GetMessages(chainId string, p Pagination) ([]*types.RouteMessage, error) {
	var messages []*types.RouteMessage

	keyPrefixList := []string{ms.prefix}
	if chainId != "" {
		keyPrefixList = append(keyPrefixList, chainId)
	}
	iter := ms.db.NewIterator(GetKey(keyPrefixList))

	// return all the messages
	if p.All {
		for iter.Next() {
			var msg types.RouteMessage
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

		var msg types.RouteMessage
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
