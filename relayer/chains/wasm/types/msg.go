package types

import (
	"fmt"

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
	ReqID string `json:"request_id"`
	Data  []byte `json:"data"`
}

type GetReceiptMsg struct {
	SrcNetwork string `json:"src_network"`
	ConnSn     string `json:"conn_sn"`
}

type ExecSendMsg struct {
	SendMessage *SendMessage `json:"send_message"`
}

type ExecRecvMsg struct {
	RecvMessage *ReceiveMessage `json:"recv_message"`
}

type ExecExecMsg struct {
	ExecMessage *ExecMessage `json:"execute_call"`
}

func NewExecRecvMsg(message *relayTypes.Message) *ExecRecvMsg {
	return &ExecRecvMsg{
		RecvMessage: &ReceiveMessage{
			SrcNetwork: message.Src,
			ConnSn:     fmt.Sprintf("%d", message.Sn),
			Msg:        hexstr.NewFromByte(message.Data),
		},
	}
}

func NewExecExecMsg(message *relayTypes.Message) *ExecExecMsg {
	return &ExecExecMsg{
		ExecMessage: &ExecMessage{
			ReqID: fmt.Sprintf("%d", message.ReqID),
			Data:  message.Data,
		},
	}
}

type QueryReceiptMsg struct {
	GetReceipt *GetReceiptMsg `json:"get_receipt"`
}

type QueryReceiptMsgResponse struct {
	Status bool `json:"status"`
}

type ExecRevertMessge struct {
	ExecMessage *RevertMessage `json:"revert_message"`
}

type RevertMessage struct {
	Sn uint64 `json:"sn"`
}

func NewExecRevertMsg(message *relayTypes.Message) *ExecRevertMessge {
	return &ExecRevertMessge{
		ExecMessage: &RevertMessage{
			Sn: message.ReqID,
		},
	}
}
