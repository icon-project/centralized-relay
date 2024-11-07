package types

import (
	"fmt"
	"math/big"

	jsoniter "github.com/json-iterator/go"

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

type ClusterReceiveMessage struct {
	SrcNetwork string           `json:"src_network"`
	ConnSn     string           `json:"conn_sn"`
	Msg        hexstr.HexString `json:"msg"`
	Signatures [][]byte         `json:"signatures"`
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

type ExecClusterRecvMsg struct {
	RecvMessage *ClusterReceiveMessage `json:"cluster_recv_message"`
}

type ExecExecMsg struct {
	ExecMessage *ExecMessage `json:"execute_call"`
}

func NewExecRecvMsg(message *relayTypes.Message) *ExecRecvMsg {
	return &ExecRecvMsg{
		RecvMessage: &ReceiveMessage{
			SrcNetwork: message.Src,
			ConnSn:     message.Sn.String(),
			Msg:        hexstr.NewFromByte(message.Data),
		},
	}
}

func NewExecClusterRecvMsg(message *relayTypes.Message) *ExecClusterRecvMsg {
	return &ExecClusterRecvMsg{
		RecvMessage: &ClusterReceiveMessage{
			SrcNetwork: message.Src,
			ConnSn:     message.Sn.String(),
			Msg:        hexstr.NewFromByte(message.Data),
			Signatures: message.Signatures,
		},
	}
}

func NewExecExecMsg(message *relayTypes.Message) *ExecExecMsg {
	exec := &ExecMessage{
		ReqID: message.ReqID.String(),
	}
	if err := jsoniter.Unmarshal(message.Data, &exec.Data); err != nil {
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
	Sn string `json:"sn"`
}

func NewExecRevertMsg(message *relayTypes.Message) *ExecRevertMessge {
	return &ExecRevertMessge{
		ExecMessage: &RevertMessage{
			Sn: fmt.Sprintf("%d", message.Sn),
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

func NewExecSetFee(networkID string, msgFee, resFee *big.Int) *ExecSetFee {
	return &ExecSetFee{
		SetFee: &SetFee{
			NetworkID:   networkID,
			MessageFee:  msgFee.String(),
			ResponseFee: resFee.String(),
		},
	}
}

// GetFee
type ExecGetFee struct {
	GetFee *GetFee `json:"get_fee"`
}

type GetFee struct {
	NetworkID string `json:"nid"`
	Response  bool   `json:"response"`
}

func NewExecGetFee(networkID string, response bool) *ExecGetFee {
	return &ExecGetFee{
		GetFee: &GetFee{
			NetworkID: networkID,
			Response:  response,
		},
	}
}

// ExecuteRollback
type ExecExecuteRollback struct {
	ExecuteRollback *ExecuteRollback `json:"execute_rollback"`
}

type ExecuteRollback struct {
	Sn string `json:"sequence_no"`
}

func NewExecExecuteRollback(sn *big.Int) *ExecExecuteRollback {
	return &ExecExecuteRollback{
		ExecuteRollback: &ExecuteRollback{
			Sn: sn.String(),
		},
	}
}
