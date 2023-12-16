package store

import (
	"encoding/json"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/types"
)

var (
	PrefixMessageStore = "message"
)

type MessageStore struct {
	dbReader KeyValueReader
	dbWriter KeyValueWriter
	prefix   string
}

type Pagination struct {
	Limit  uint
	Offset uint
	All    bool
}

func NewPagination() *Pagination {
	return new(Pagination)
}

func (p *Pagination) GetAll() *Pagination {
	p.All = true
	return p
}

func (p *Pagination) WithLimit(l uint) *Pagination {
	p.Limit = l
	return p
}

// WithPage sets the page and calculates the offset
func (p *Pagination) WithPage(page, limit uint) *Pagination {
	p.Limit = limit
	p.Offset = p.CalculateOffset(page)
	return p
}

// CalculateTotalPages calculates the total pages based on the limit and total count
func (p *Pagination) CalculateTotalPages(total int) uint {
	page := uint(total) / p.Limit
	if uint(total)%p.Limit != 0 {
		page++
	}
	return page
}

// CalculateOffset calculates the offset based on the page and limit
func (p *Pagination) CalculateOffset(page uint) uint {
	if page <= 1 {
		return 0
	}
	return page * p.Limit
}

func (p *Pagination) WithOffset(o uint) *Pagination {
	p.Offset = o
	return p
}

func NewMessageStore(db Store) *MessageStore {
	return &MessageStore{
		dbReader: db,
		dbWriter: db,
		prefix:   PrefixMessageStore,
	}
}

func NewMessageStoreReadOnly(keyValueReader KeyValueReader) *MessageStore {
	return &MessageStore{
		dbReader: keyValueReader,
		prefix:   PrefixMessageStore,
	}
}

func (ms *MessageStore) TotalCount() (uint64, error) {
	return ms.getCountByKey(GetKey([]string{ms.prefix}))
}

func (ms *MessageStore) TotalCountByChain(nId string) (uint64, error) {
	return ms.getCountByKey(GetKey([]string{ms.prefix, nId}))
}

func (ms *MessageStore) getCountByKey(key []byte) (uint64, error) {
	iter := ms.dbReader.NewIterator(key)
	if iter == nil {
		return 0, fmt.Errorf("both db and snapshot object cannot be empty")
	}

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

	key := GetKey([]string{ms.prefix, message.Src, fmt.Sprintf("%d", message.Sn)})

	msgByte, err := ms.Encode(message)
	if err != nil {
		return err
	}
	return ms.dbWriter.SetByKey(key, msgByte)
}

func (ms *MessageStore) GetMessage(messageKey types.MessageKey) (*types.RouteMessage, error) {

	v, err := ms.dbReader.GetByKey(GetKey([]string{ms.prefix, messageKey.Src, fmt.Sprintf("%d", messageKey.Sn)}))
	if err != nil {
		return nil, err
	}

	msg := new(types.RouteMessage)
	if err := ms.Decode(v, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (ms *MessageStore) GetMessages(nId string, p *Pagination) ([]*types.RouteMessage, error) {
	var messages []*types.RouteMessage

	keyPrefixList := []string{ms.prefix}
	if nId != "" {
		keyPrefixList = append(keyPrefixList, nId)
	}
	iter := ms.dbReader.NewIterator(GetKey(keyPrefixList))

	// return all the messages
	if p.All {
		for iter.Next() {
			msg := new(types.RouteMessage)
			if err := ms.Decode(iter.Value(), msg); err != nil {
				return nil, err
			}

			messages = append(messages, msg)
		}
		iter.Release()
		return messages, iter.Error()
	}

	// if not all use the offset logic
	for i := 0; i < int(p.Offset); i++ {
		if !iter.Next() {
			return nil, fmt.Errorf("no message after offset")
		}
	}

	for i := uint(0); i < p.Limit; i++ {
		if !iter.Next() {
			break
		}

		msg := new(types.RouteMessage)
		if err := ms.Decode(iter.Value(), msg); err != nil {
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
	return ms.dbWriter.DeleteByKey(GetKey([]string{ms.prefix, messageKey.Src, fmt.Sprintf("%d", messageKey.Sn)}))
}

func (ms *MessageStore) Encode(d interface{}) ([]byte, error) {
	return json.Marshal(d)
}

func (ms *MessageStore) Decode(data []byte, output interface{}) error {
	return json.Unmarshal(data, output)
}
