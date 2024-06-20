package types

import (
	cctypes "github.com/coming-chat/go-sui/v2/types"
)

const (
	ChainType          = "sui"
	XcallContract      = "xcall"
	ConnectionContract = "connection"

	ConnectionIDMismatchError = "connection_id_mismatch_error"
	WsConnReadError           = "ws_conn_read_err"

	QUERY_MAX_RESULT_LIMIT = 50
)

type ContractConfigMap map[string]string

type SuiGetCheckpointsRequest struct {
	// optional paging cursor
	Cursor interface{} `json:"cursor"`
	// maximum number of items per page
	Limit uint64 `json:"limit" validate:"lte=50"`
	// query result ordering, default to false (ascending order), oldest record first
	DescendingOrder bool `json:"descendingOrder"`
}

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

type PaginatedCheckpointsResponse struct {
	Data        []CheckpointResponse `json:"data"`
	NextCursor  string               `json:"nextCursor"`
	HasNextPage bool                 `json:"hasNextPage"`
}

type TxDigests struct {
	FromCheckpoint uint64
	ToCheckpoint   uint64
	Digests        []string
}

type EventResponse struct {
	cctypes.SuiEvent
	Checkpoint *cctypes.SafeSuiBigInt[uint64]
}

type SuiMultiGetTransactionBlocksRequest struct {
	Digests []string                                   `json:"digests"`
	Options cctypes.SuiTransactionBlockResponseOptions `json:"options"`
}

type EmitEvent struct {
	Sn           string `json:"conn_sn"`
	Msg          []byte `json:"msg"`
	To           string `json:"to"`
	ConnectionID string `json:"connection_id"`
}

type CallMsgEvent struct {
	ReqId           string `json:"req_id"`
	Data            []byte `json:"data"`
	DappModuleCapId string `json:"to"`
}

type RollbackMsgEvent struct {
	Sn              string `json:"sn"`
	Data            []byte `json:"data"`
	DappModuleCapId string `json:"dapp"`
}

type EventNotification struct {
	cctypes.SuiEvent
	Error error
}
type WsSubscriptionResp struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  int64  `json:"result"`
	Id      int64  `json:"id"`
}

type JsonRPCRequest struct {
	Version string      `json:"jsonrpc,omitempty"`
	ID      interface{} `json:"id,omitempty"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
}

type EventQueryFilter struct {
	FromCheckpoint uint64
	ToCheckpoint   uint64
	Packages       []string
	EventModule    string
}

type EventQueryResponse struct {
	Data        []cctypes.SuiEvent `json:"data"`
	NextCursor  cctypes.EventId    `json:"nextCursor"`
	HasNextPage bool               `json:"hasNextPage"`
}

type EventQueryRequest struct {
	EventFilter interface{}
	Cursor      cctypes.EventId
	Limit       uint64
	Descending  bool
}

type SuiMethod string

func (sm SuiMethod) String() string {
	return string(sm)
}
