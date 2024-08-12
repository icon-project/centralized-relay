package types

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"
)

var (
	MaxTxRetry         uint8 = 5
	StaleMarkCount           = MaxTxRetry * 3
	RouteDuration            = 3 * time.Second
	XcallContract            = "xcall"
	ConnectionContract       = "connection"
	SupportedContracts       = []string{XcallContract, ConnectionContract}
	RetryInterval            = 3*time.Second + RouteDuration
)

type BlockInfo struct {
	Height   uint64
	Messages []*Message
}

type Message struct {
	Dst           string   `json:"dst"`
	Src           string   `json:"src"`
	Sn            *big.Int `json:"sn"`
	Data          []byte   `json:"data"`
	MessageHeight uint64   `json:"messageHeight"`
	EventType     string   `json:"eventType"`
	ReqID         *big.Int `json:"reqID,omitempty"`
}

type ContractConfigMap map[string]string

type EventMap struct {
	Address      string
	ContractName string
	SigType      map[string]string
}

func (c ContractConfigMap) Validate() error {
	for _, contract := range SupportedContracts {
		val, ok := (c)[contract]
		if !ok {
			continue
		}
		if val == "" {
			continue
		}
	}
	return nil
}

// GetWasmMsgType returns the wasm message type
func (m *EventMap) GetWasmMsgType() string {
	for wasmType := range m.SigType {
		return wasmType
	}
	return ""
}

func (m *Message) MessageKey() *MessageKey {
	return NewMessageKey(m.Sn, m.Src, m.Dst, m.EventType)
}

type RouteMessage struct {
	*Message
	Retry      uint8
	Processing bool
	LastTry    time.Time
}

func NewRouteMessage(m *Message) *RouteMessage {
	return &RouteMessage{
		Message: m,
	}
}

func (r *RouteMessage) GetMessage() *Message {
	return r.Message
}

func (r *RouteMessage) IncrementRetry() {
	r.Retry++
	r.AddNextTry()
}

func (r *RouteMessage) ToggleProcessing() {
	r.Processing = !r.Processing
}

func (r *RouteMessage) GetRetry() uint8 {
	return r.Retry
}

// ResetLastTry resets the last try time to the current time plus the retry interval
func (r *RouteMessage) AddNextTry() {
	r.LastTry = time.Now().Add(RetryInterval * time.Duration(math.Pow(2, float64(r.Retry)))) // exponential backoff
}

func (r *RouteMessage) IsProcessing() bool {
	return r.Processing || !(r.LastTry.IsZero() || r.LastTry.Before(time.Now()))
}

// stale means message which is expired
func (r *RouteMessage) IsStale() bool {
	return r.Retry >= StaleMarkCount
}

// IsElasped checks if the last try is elasped by the duration
func (r *RouteMessage) IsElasped(duration time.Duration) bool {
	return r.LastTry.Add(duration).Before(time.Now())
}

type TxResponseFunc func(key *MessageKey, response *TxResponse, err error)

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
	Sn        *big.Int
	Src       string
	Dst       string
	EventType string
}

func NewMessageKey(sn *big.Int, src string, dst string, eventType string) *MessageKey {
	return &MessageKey{sn, src, dst, eventType}
}

type MessageKeyWithMessageHeight struct {
	*MessageKey
	Height uint64
}

func NewMessagekeyWithMessageHeight(key *MessageKey, height uint64) *MessageKeyWithMessageHeight {
	return &MessageKeyWithMessageHeight{key, height}
}

type MessageCache struct {
	Messages map[string]*RouteMessage
	*sync.RWMutex
}

func NewMessageCache() *MessageCache {
	return &MessageCache{
		Messages: make(map[string]*RouteMessage),
		RWMutex:  new(sync.RWMutex),
	}
}

func (m *MessageCache) Add(r *RouteMessage) {
	m.Lock()
	defer m.Unlock()
	cacheKey := m.GetCacheKey(r.MessageKey())
	if !m.HasCacheKey(cacheKey) {
		m.Messages[cacheKey] = r
	}
}

func (m *MessageCache) Len() int {
	return len(m.Messages)
}

func (m *MessageCache) Remove(key *MessageKey) {
	m.Lock()
	defer m.Unlock()
	cacheKey := m.GetCacheKey(key)
	delete(m.Messages, cacheKey)
}

// Get returns the message from the cache
func (m *MessageCache) Get(key *MessageKey) (*RouteMessage, bool) {
	m.RLock()
	defer m.RUnlock()
	cacheKey := m.GetCacheKey(key)
	msg, ok := m.Messages[cacheKey]
	return msg, ok
}

func (m *MessageCache) GetCacheKey(key *MessageKey) string {
	return key.Src + "-" + key.Dst + "-" + key.EventType + "-" + key.Sn.String()
}

func (m *MessageCache) HasCacheKey(cacheKey string) bool {
	_, ok := m.Messages[cacheKey]
	return ok
}

type Coin struct {
	Denom  string
	Amount uint64
}

func NewCoin(denom string, amount uint64) *Coin {
	return &Coin{strings.ToLower(denom), amount}
}

func (c *Coin) String() string {
	return fmt.Sprintf("%d%s", c.Amount, c.Denom)
}

func (c *Coin) Calculate() string {
	balance := new(big.Float).SetUint64(c.Amount)
	amount := balance.Quo(balance, big.NewFloat(1e18))
	value, _ := amount.Float64()
	return fmt.Sprintf("%.18f %s", value, c.Denom)
}

type TransactionObject struct {
	*MessageKeyWithMessageHeight
	TxHash   string
	TxHeight uint64
}

func NewTransactionObject(messageKey *MessageKeyWithMessageHeight, txHash string, height uint64) *TransactionObject {
	return &TransactionObject{messageKey, txHash, height}
}

type Receipt struct {
	TxHash string
	Height uint64
	Status bool
}
