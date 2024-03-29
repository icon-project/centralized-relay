package types

import (
	"context"

	suimodels "github.com/block-vision/sui-go-sdk/models"
)

const (
	ChainType = "sui"

	WsConnReadError = "ws_conn_read_err"
)

type EventNotification struct {
	suimodels.SuiEventResponse
	Error error
}

type TimeRange struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type IClient interface {
	GetLatestCheckpointSeq(ctx context.Context) (uint64, error)
	GetCheckpoints(ctx context.Context, req suimodels.SuiGetCheckpointsRequest) (suimodels.PaginatedCheckpointsResponse, error)
	GetBalance(ctx context.Context, addr string) ([]suimodels.CoinData, error)

	SubscribeEventNotification(done chan interface{}, wsUrl string, eventFilters []interface{}) (<-chan EventNotification, error)
}
