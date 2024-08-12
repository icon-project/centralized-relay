package types

import (
	cctypes "github.com/coming-chat/go-sui/v2/types"
)

const (
	ChainType          = "sui"
	XcallContract      = "xcall"
	ConnectionContract = "connection"

	InvalidEventError = "invalid_event_err"
)

type EpochRollingGasCostSummary struct {
	ComputationCost         string `json:"computationCost"`
	StorageCost             string `json:"storageCost"`
	StorageRebate           string `json:"storageRebate"`
	NonRefundableStorageFee string `json:"nonRefundableStorageFee"`
}

type CheckpointResponse struct {
	Epoch                      string                     `json:"epoch"`
	SequenceNumber             string                     `json:"sequenceNumber"`
	Digest                     string                     `json:"digest"`
	NetworkTotalTransactions   string                     `json:"networkTotalTransactions"`
	PreviousDigest             string                     `json:"previousDigest"`
	EpochRollingGasCostSummary EpochRollingGasCostSummary `json:"epochRollingGasCostSummary"`
	TimestampMs                string                     `json:"timestampMs"`
	Transactions               []string                   `json:"transactions"`
	CheckpointCommitments      []interface{}              `json:"checkpointCommitments"`
	ValidatorSignature         string                     `json:"validatorSignature"`
}

type EventResponse struct {
	cctypes.SuiEvent
	Checkpoint *cctypes.SafeSuiBigInt[uint64]
}

type EmitEvent struct {
	Sn           string `json:"conn_sn"`
	Msg          []byte `json:"msg"`
	To           string `json:"to"`
	ConnectionID string `json:"connection_id"`
}

type NetworkAddress struct {
	Addr  string `json:"addr"`
	NetID string `json:"net_id"`
}
type CallMsgEvent struct {
	Sn              string         `json:"sn"`
	From            NetworkAddress `json:"from"`
	ReqId           string         `json:"req_id"`
	Data            []byte         `json:"data"`
	DappModuleCapId string         `json:"to"`
}

type RollbackMsgEvent struct {
	Sn              string `json:"sn"`
	Data            []byte `json:"data"`
	DappModuleCapId string `json:"dapp"`
}

type SuiMethod string

func (sm SuiMethod) String() string {
	return string(sm)
}
