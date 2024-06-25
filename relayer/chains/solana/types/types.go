package types

import "github.com/gagliardetto/solana-go"

const (
	MethodSendMessage = "send_message"
	MethodRecvMessage = "recv_message"

	ChainType = "solana"
)

type SolEvent struct {
	Slot      uint64
	Signature solana.Signature
	Logs      []string
}
