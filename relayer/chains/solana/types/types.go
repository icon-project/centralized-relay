package types

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
)

const (
	MethodSetAdmin = "set_admin"

	MethodSetFee    = "set_fee"
	MethodClaimFees = "claim_fees"

	MethodRevertMessage = "revert_message"

	MethodDecodeCsMessage = "decode_cs_message"

	MethodSendMessage     = "send_message"
	MethodRecvMessage     = "recv_message"
	MethodExecuteCall     = "execute_call"
	MethodExecuteRollback = "execute_rollback"

	MethodQueryRecvMessageAccounts     = "query_recv_message_accounts"
	MethodQueryExecuteCallAccounts     = "query_execute_call_accounts"
	MethodQueryExecuteRollbackAccounts = "query_execute_rollback_accounts"
	MethodQueryRevertMessageAccounts   = "query_revert_message_accounts"

	ChainType = "solana"

	EventLogPrefix      = "Program data: "
	ProgramReturnPrefix = "Program return: "

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
	FromNetworkAddress string
	To                 string
	Sn                 big.Int
	ReqId              big.Int
	Data               []byte
	ConnProgram        solana.PublicKey
	ConnSn             big.Int
}

type RollbackMessageEvent struct {
	Sn big.Int
}

type QueryAccountsResponse struct {
	Accounts      []solana.AccountMeta
	TotalAccounts uint8
	Limit         uint8
	Page          uint8
	HasNextPage   bool
}

type Dapp struct {
	Name         string   `yaml:"name"`
	ProgramID    string   `yaml:"program-id"`
	ConfigPrefix string   `yaml:"config-prefix"`
	OtherPrefix  []string `yaml:"other-prefix"`
}
