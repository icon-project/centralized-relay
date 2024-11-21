package events

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID          string
	Type        string
	Data        interface{}
	BlockHeight uint64
	Timestamp   time.Time
	Raw         []byte
}

type SmartContractLogEvent struct {
	EventType     string        `json:"event_type"`
	ContractEvent ContractEvent `json:"contract_event"`
	TxID          string        `json:"tx_id"`
	BlockHeight   uint64        `json:"block_height"`
}

type ContractEvent struct {
	ContractID string          `json:"contract_id"`
	Topic      string          `json:"topic"`
	Value      json.RawMessage `json:"value"`
}

type WSMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type WSEventParams struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

type WSRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int64       `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type WSResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type CallMessageSentData struct {
	From         string   `json:"from"`
	To           string   `json:"to"`
	Sn           uint64   `json:"sn"`
	Data         string   `json:"data"`
	Sources      []string `json:"sources"`
	Destinations []string `json:"destinations"`
}

type CallMessageData struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Sn    uint64 `json:"sn"`
	ReqID uint64 `json:"req-id"`
	Data  string `json:"data"`
}

type ResponseMessageData struct {
	Sn   uint64 `json:"sn"`
	Code uint64 `json:"code"`
	Msg  string `json:"msg"`
}

type RollbackMessageData struct {
	Sn uint64 `json:"sn"`
}

const (
	CallMessageSent = "CallMessageSent"
	CallMessage     = "CallMessage"
	ResponseMessage = "ResponseMessage"
	RollbackMessage = "RollbackMessage"
)

type EventHandler func(event *Event) error

type EventStore interface {
	SaveEvent(event *Event) error
	GetEvents(fromHeight uint64) ([]*Event, error)
	MarkProcessed(eventID string) error
	GetLastProcessedHeight() (uint64, error)
}
