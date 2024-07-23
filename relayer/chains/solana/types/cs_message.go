package types

import (
	"math/big"
)

type CsMessage struct {
	MessageType CsMessageType
	Request     *CsMessageRequestType
	Result      *CsMessageResultType
}

type CsMessageRequestType struct {
	From       struct{ Address string }
	To         string
	SequenceNo big.Int
	MsgType    MessageType
	Data       []byte
	Protocols  []string
}

type CsMessageResultType struct {
	SequenceNo   big.Int
	ResponseCode CsResponseType
	Message      []byte
}

type CsResponseType uint8

const (
	CsResponseFailure CsResponseType = iota
	CsResponseSuccess
)

type MessageType uint8

const (
	CallMessage MessageType = iota
	CallMessageWithRollback
	CallMessagePersisted
)

type CsMessageType uint8

const (
	CsMessageRequest CsMessageType = iota
	CsMessageResult
)

type ProxyRequestAccount struct {
	ReqMessage CsMessageRequestType
	Bump       uint8
}
