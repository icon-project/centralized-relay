package types

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
)

const (
	MethodSetAdmin = "set_admin"

	MethodSetFee    = "set_fee"
	MethodClaimFees = "claim_fees"

	MethodSendMessage = "send_message"
	MethodRecvMessage = "recv_message"

	ChainType = "solana"

	EventLogPrefix = "Program data: "

	EventSendMessage     = "SendMessage"
	EventCallMessage     = "CallMessage"
	EventRollbackMessage = "RollbackMessage"

	SolanaDenom = "lamport"
)

type SolEvent struct {
	Slot      uint64
	Signature solana.Signature
	Logs      []string
}

type SendMessageEvent struct {
	TargetNetwork string
	ConnSn        big.Int
	Msg           []byte
}

type CallMessageEvent struct {
	From  string
	To    string
	Sn    big.Int
	ReqId big.Int
	Data  []byte
}

type RollbackMessageEvent struct {
	Sn big.Int
}
