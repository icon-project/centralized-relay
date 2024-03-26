package types

import (
	"context"

	suimodels "github.com/block-vision/sui-go-sdk/models"
)

const (
	ChainType = "sui"
)

type IClient interface {
	GetLatestCheckpointSeq(ctx context.Context) (uint64, error)
	GetCheckpoints(ctx context.Context, req suimodels.SuiGetCheckpointsRequest) (suimodels.PaginatedCheckpointsResponse, error)
}
