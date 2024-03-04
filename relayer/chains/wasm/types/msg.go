package types

import (
	"encoding/json"
	"fmt"
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
	ReqID string `json:"request_id"`
	Data  []int  `json:"data"`
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
	exec := &ExecMessage{
		ReqID: strconv.FormatUint(message.ReqID, 10),
	}
	if err := json.Unmarshal(message.Data, &exec.Data); err != nil {
		return nil
	}
	return &ExecExecMsg{
		ExecMessage: exec,
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

// SetAdmin
type SetAdmin struct {
	Address string `json:"address"`
}

type ExecSetAdmin struct {
	SetAdmin *SetAdmin `json:"set_admin"`
}

func NewExecSetAdmin(address string) *ExecSetAdmin {
	return &ExecSetAdmin{
		SetAdmin: &SetAdmin{
			Address: address,
		},
	}
}
