package types

import (
	"math/big"

	"github.com/near/borsh-go"
)

type CsMessage struct {
	Variant borsh.Enum `borsh_enum:"true"`
	Request CsMessageRequestType
	Result  CsMessageResultType
}

type CsMessageRequestType struct {
	From       string
	To         string
	SequenceNo *big.Int
	MsgType    MessageType
	Data       []byte
	Protocols  []string
}

type CsMessageResultType struct {
	SequenceNo   *big.Int
	ResponseCode CsResponseType
	Message      []byte
}

type CsResponseType borsh.Enum

const (
	CsResponseFailure CsResponseType = iota
	CsResponseSuccess
)

type MessageType borsh.Enum

const (
	CallMessage MessageType = iota
	CallMessageWithRollback
	CallMessagePersisted
)

type CsMessageType borsh.Enum

const (
	CsMessageRequest CsMessageType = iota
	CsMessageResult
)
