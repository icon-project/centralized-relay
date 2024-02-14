package types

import (
	"strconv"

	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/hexstr"
)

type SendMessage struct {
	To  string `json:"to"`
	Svc uint8  `json:"svc"`
	Sn  uint64 `json:"sn"`
	Msg []byte `json:"msg"`
}

type ReceiveMessage struct {
	SrcNetwork string           `json:"src_network"`
	ConnSn     string           `json:"conn_sn"`
	Msg        hexstr.HexString `json:"msg"`
}

type ExecMessage struct {
	SrcNetwork string           `json:"src_network"`
	ReqID      string           `json:"reqId"`
	Msg        hexstr.HexString `json:"msg"`
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

type ExecExecMsg struct {
	ExecMessage ExecMessage `json:"exec_message"`
}

func NewExecRecvMsg(message *relayTypes.Message) ExecRecvMsg {
	return ExecRecvMsg{
		RecvMessage: ReceiveMessage{
			SrcNetwork: message.Src,
			ConnSn:     strconv.Itoa(int(message.Sn)),
			Msg:        hexstr.NewFromByte(message.Data),
		},
	}
}

func NewExecExecMsg(message *relayTypes.Message) ExecExecMsg {
	return ExecExecMsg{
		ExecMessage: ExecMessage{
			SrcNetwork: message.Src,
			ReqID:      strconv.Itoa(int(message.ReqID)),
			Msg:        hexstr.NewFromByte(message.Data),
		},
	}
}

type QueryReceiptMsg struct {
	GetReceipt GetReceiptMsg `json:"get_receipt"`
}

type QueryReceiptMsgResponse struct {
	Status bool `json:"status"`
}
