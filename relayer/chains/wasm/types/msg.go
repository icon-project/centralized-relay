package types

import (
	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
	"strconv"
)

type SendMessage struct {
	To  string `json:"to"`
	Svc uint8  `json:"svc"`
	Sn  uint64 `json:"sn"`
	Msg []byte `json:"msg"`
}

type ReceiveMessage struct {
	SrcNetwork string   `json:"src_network"`
	ConnSn     string   `json:"conn_sn"`
	Msg        HexBytes `json:"msg"`
}

type GetReceiptMsg struct {
	SrcNetwork string `json:"src_network"`
	ConnSn     string `json:"conn_sn"`
}

type ExecSendMsg struct {
	SendMessage SendMessage `json:"send_message"`
}

type ExecRecvMsg struct {
	RecvMessage ReceiveMessage `json:"recv_message"`
}

func NewExecRecvMsg(message *relayTypes.Message) ExecRecvMsg {
	return ExecRecvMsg{
		RecvMessage: ReceiveMessage{
			SrcNetwork: message.Src,
			ConnSn:     strconv.Itoa(int(message.Sn)),
			Msg:        NewHexBytes(message.Data),
		},
	}
}

type QueryReceiptMsg struct {
	GetReceipt GetReceiptMsg `json:"get_receipt"`
}

type QueryReceiptMsgResponse struct {
	Status bool `json:"status"`
}
