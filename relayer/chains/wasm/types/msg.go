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
			Sn: message.Sn,
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

// ClaimFee
type ClaimFee struct{}

type ExecClaimFee struct {
	ClaimFee *ClaimFee `json:"claim_fees"`
}

func NewExecClaimFee() *ExecClaimFee {
	return &ExecClaimFee{
		ClaimFee: &ClaimFee{},
	}
}

// SetFee
type SetFee struct {
	NetworkID   string `json:"network_id"`
	MessageFee  string `json:"message_fee"`
	ResponseFee string `json:"response_fee"`
}

type ExecSetFee struct {
	SetFee *SetFee `json:"set_fee"`
}

func NewExecSetFee(networkID string, msgFee, resFee uint64) *ExecSetFee {
	return &ExecSetFee{
		SetFee: &SetFee{
			NetworkID:   networkID,
			MessageFee:  strconv.FormatUint(msgFee, 10),
			ResponseFee: strconv.FormatUint(resFee, 10),
		},
	}
}

// GetFee
type GetFee struct {
	NetworkID string `json:"network_id"`
	Response  bool   `json:"response"`
}

type ExecGetFee struct {
	GetFee *GetFee `json:"get_fee"`
}

func NewExecGetFee(networkID string, response bool) *ExecGetFee {
	return &ExecGetFee{
		GetFee: &GetFee{
			NetworkID: networkID,
			Response:  response,
		},
	}
}

type QueryGetFeeResponse struct {
	Total uint64 `json:"message_fee"`
}
