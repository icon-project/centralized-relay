package store

import (
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/types"
	jsoniter "github.com/json-iterator/go"
)

type MessageStore struct {
	db     Store
	prefix string
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

func NewMessageStore(db Store, prefix string) *MessageStore {
	return &MessageStore{
		db:     db,
		prefix: prefix,
	}
}

func (ms *MessageStore) TotalCount() (uint, error) {
	return ms.getCountByKey(GetKey([]string{ms.prefix}))
}

func (ms *MessageStore) TotalCountByChain(nId string) (uint, error) {
	return ms.getCountByKey(GetKey([]string{ms.prefix, nId}))
}

func (ms *MessageStore) getCountByKey(key []byte) (uint, error) {
	iter := ms.db.NewIterator(key)
	var count uint
	for iter.Next() {
		count++
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ms *MessageStore) StoreMessage(message *types.RouteMessage) error {
	if message == nil {
		return fmt.Errorf("error while storingMessage: message cannot be nil")
	}

	key := GetKey([]string{ms.prefix, message.Src, message.Sn.String()})

	msgByte, err := ms.Encode(message)
	if err != nil {
		return err
	}
	return ms.db.SetByKey(key, msgByte)
}

func (ms *MessageStore) GetMessage(messageKey *types.MessageKey) (*types.RouteMessage, error) {
	v, err := ms.db.GetByKey(GetKey([]string{ms.prefix, messageKey.Src, messageKey.Sn.String()}))
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

	iter := ms.db.NewIterator(GetKey([]string{ms.prefix, nId}))
	defer iter.Release()

	// return all the messages
	for iter.Next() {
		msg := new(types.RouteMessage)
		if err := ms.Decode(iter.Value(), msg); err != nil {
			return nil, err
		}

		messages = append(messages, msg)
		if uint(len(messages)) == p.Limit {
			break
		}
	}

	return messages, iter.Error()
}

func (ms *MessageStore) DeleteMessage(messageKey *types.MessageKey) error {
	return ms.db.DeleteByKey(GetKey([]string{ms.prefix, messageKey.Src, messageKey.Sn.String()}))
}

func (ms *MessageStore) Encode(d interface{}) ([]byte, error) {
	return jsoniter.Marshal(d)
}

func (ms *MessageStore) Decode(data []byte, output interface{}) error {
	return jsoniter.Unmarshal(data, output)
}
