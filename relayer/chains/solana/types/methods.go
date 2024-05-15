package types

const (
	MethodSendMessage = 0
	MethodRecvMessage = 1
	MethodExecuteCall = 2
)

type SendMessageParams struct {
	To   string `borsh:"to"`
	Data []byte `borsh:"data"`
}

type RecvMessageParams struct {
	Sn   uint64 `borsh:"sn"`
	Src  string `borsh:"src"`
	Data []byte `borsh:"data"`
}

type ExecuteCallParams struct {
	ReqId uint64 `borsh:"req_id"`
	Data  []byte `borsh:"data"`
}
