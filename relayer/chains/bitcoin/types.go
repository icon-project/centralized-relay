package bitcoin

import (
	"github.com/btcsuite/btcd/wire"
)

type MessageType int

const (
	CS_REQUEST MessageType = iota + 1
	CS_RESPONSE
	CS_RESULT
)

type CallMessageType int

const (
	CALL_MESSAGE_TYPE MessageType = iota
	CALL_MESSAGE_ROLLBACK_TYPE
	PERSISTENT_MESSAGE_TYPE
)

type TxSearchParam struct {
	StartHeight, EndHeight uint64
	BitcoinScript          []byte
	OPReturnPrefix         int
}

type TxSearchRes struct {
	Tx     *wire.MsgTx
	Height uint64
}

// HightRange is a struct to represent a range of heights
type HeightRange struct {
	Start uint64
	End   uint64
}

type XCallMessage struct {
	Action       string
	TokenAddress string
	From         string
	To           string
	Amount       []byte
	Data         []byte
}

type RuneInfo struct {
	Rune         string `json:"rune"`
	RuneId       string `json:"runeId"`
	SpaceRune    string `json:"spaceRune"`
	Amount       string `json:"amount"`
	Symbol       string `json:"symbol"`
	Divisibility int    `json:"divisibility"`
}

type RuneTxIndexResponse struct {
	Code int        `json:"code"`
	Data []RuneInfo `json:"data"`
}

type CSMessageRequestV2 struct {
	From        string   `json:"from"`
	To          string   `json:"to"`
	Sn          []byte   `json:"sn"`
	MessageType uint8    `json:"messageType"`
	Data        []byte   `json:"data"`
	Protocols   []string `json:"protocols"`
}

type CSMessage struct {
	MsgType []byte `json:"msgType"`
	Payload []byte `json:"payload"`
}
