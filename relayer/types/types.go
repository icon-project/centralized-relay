package types

import (
	"fmt"
	"sync"
)

var (
	DefaultTxRetry uint8 = 3
	// message is stale after TotalMaxRetryTx
	TotalMaxRetryTx uint8 = DefaultTxRetry * 5
)

type BlockInfo struct {
	Height   uint64
	Messages []*Message
}

type Message struct {
	Dst           string `json:"dst"`
	Src           string `json:"src"`
	Sn            uint64 `json:"sn"`
	Data          []byte `json:"data"`
	MessageHeight uint64 `json:"messageHeight"`
	EventType     string `json:"eventType"`
}

func (m *Message) MessageKey() MessageKey {
	return NewMessageKey(m.Sn, m.Src, m.Dst, m.EventType)
}

type RouteMessage struct {
	*Message
	Retry        uint64
	IsProcessing bool
}

func NewRouteMessage(m *Message) *RouteMessage {
	return &RouteMessage{
		Message:      m,
		Retry:        0,
		IsProcessing: false,
	}
}

func (r *RouteMessage) GetMessage() *Message {
	return r.Message
}

func (r *RouteMessage) IncrementRetry() {
	r.Retry += 1
}

func (r *RouteMessage) GetRetry() uint64 {
	return r.Retry
}

func (r *RouteMessage) SetIsProcessing(isProcessing bool) {
	r.IsProcessing = isProcessing
}

func (r *RouteMessage) GetIsProcessing() bool {
	return r.IsProcessing
}

// stale means message which is expired
func (r *RouteMessage) IsStale() bool {
	return r.Retry >= uint64(TotalMaxRetryTx)
}

type TxResponseFunc func(key MessageKey, response TxResponse, err error)

type TxResponse struct {
	Height    int64
	TxHash    string
	Codespace string
	Code      ResponseCode
	Data      string
}

type ResponseCode uint8

const (
	Failed  ResponseCode = 0
	Success ResponseCode = 1
)

type MessageKey struct {
	Sn        uint64
	Src       string
	Dst       string
	EventType string
}

func NewMessageKey(sn uint64, src string, dst string, eventType string) MessageKey {
	return MessageKey{sn, src, dst, eventType}
}

type MessageKeyWithMessageHeight struct {
	MessageKey
	MsgHeight uint64
}

func NewMessagekeyWithMessageHeight(key MessageKey, height uint64) *MessageKeyWithMessageHeight {
	return &MessageKeyWithMessageHeight{key, height}
}

type MessageCache struct {
	Messages map[MessageKey]*RouteMessage
	sync.Mutex
}

func NewMessageCache() *MessageCache {
	return &MessageCache{
		Messages: make(map[MessageKey]*RouteMessage),
	}
}

func (m *MessageCache) Add(r *RouteMessage) {
	m.Lock()
	defer m.Unlock()
	m.Messages[r.MessageKey()] = r
}

func (m *MessageCache) Len() uint64 {
	return uint64(len(m.Messages))
}

func (m *MessageCache) Remove(key MessageKey) {
	m.Lock()
	defer m.Unlock()
	delete(m.Messages, key)
}

type Coin struct {
	Denom  string
	Amount uint64
}

func NewCoin(denom string, amount uint64) Coin {
	return Coin{denom, amount}
}

func (c *Coin) String() string {
	return fmt.Sprintf("%d%s", c.Amount, c.Denom)
}

type TransactionObject struct {
	MessageKeyWithMessageHeight
	TxHash   string
	TxHeight uint64
}

func NewTransactionObject(messageKey MessageKeyWithMessageHeight, txHash string, height uint64) *TransactionObject {
	return &TransactionObject{messageKey, txHash, height}
}

type Receipt struct {
	TxHash string
	Height uint64
	Status bool
}
