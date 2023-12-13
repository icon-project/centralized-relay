package types

import (
	"fmt"
	"sync"
	"time"
)

var (
	DefaultTxRetry = 3
	// message is stale after TotalMaxRetryTx
	TotalMaxRetryTx = DefaultTxRetry * 5
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

func (m *Message) MessageKey() *MessageKey {
	return NewMessageKey(m.Sn, m.Src, m.Dst, m.EventType)
}

type RouteMessage struct {
	*Message
	Retry        uint64
	IsProcessing bool
	Time         int64
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

func (r *RouteMessage) SetTime() {
	r.Time = time.Now().Unix()
}

// Parse time to human readable format
func (r *RouteMessage) GetTime() string {
	return time.Unix(r.Time, 0).String()
}

type TxResponseFunc func(key *MessageKey, response TxResponse, err error)

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

func NewMessageKey(sn uint64, src string, dst string, eventType string) *MessageKey {
	return &MessageKey{sn, src, dst, eventType}
}

type MessageCache struct {
	Messages map[*MessageKey]*RouteMessage
	sync.Mutex
}

func NewMessageCache() *MessageCache {
	return &MessageCache{
		Messages: make(map[*MessageKey]*RouteMessage),
	}
}

func (m *MessageCache) Add(r *RouteMessage) {
	key := NewMessageKey(r.Sn, r.Src, r.Dst, r.EventType)

	m.Lock()
	defer m.Unlock()
	m.Messages[key] = r
}

func (m *MessageCache) Len() uint64 {
	return uint64(len(m.Messages))
}

func (m *MessageCache) Remove(key *MessageKey) {
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
