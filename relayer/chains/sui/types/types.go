package types

import (
	"context"

	suimodels "github.com/block-vision/sui-go-sdk/models"
)

const (
	ChainType = "sui"

	WsConnReadError = "ws_conn_read_err"

	QUERY_MAX_RESULT_LIMIT = 50
)

type EventNotification struct {
	suimodels.SuiEventResponse
	Error error
}

type TimeRange struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type TxDigests struct {
	FromCheckpoint uint64
	ToCheckpoint   uint64
	Digests        []string
}

type IClient interface {
	GetLatestCheckpointSeq(ctx context.Context) (uint64, error)
	GetCheckpoints(ctx context.Context, req suimodels.SuiGetCheckpointsRequest) (suimodels.PaginatedCheckpointsResponse, error)
	GetBalance(ctx context.Context, addr string) ([]suimodels.CoinData, error)

	SubscribeEventNotification(done chan interface{}, wsUrl string, eventFilters []interface{}) (<-chan EventNotification, error)
	GetEventsFromTxBlocks(ctx context.Context, digests []string) ([]suimodels.SuiEventResponse, error)
}
