package client

import (
	"context"

	suimodels "github.com/block-vision/sui-go-sdk/models"
	suisdk "github.com/block-vision/sui-go-sdk/sui"
)

type Client struct {
	rpc suisdk.ISuiAPI
}

func NewClient(rpcClient suisdk.ISuiAPI) *Client {
	return &Client{
		rpc: rpcClient,
	}
}

func (c Client) GetLatestCheckpointSeq(ctx context.Context) (uint64, error) {
	return c.rpc.SuiGetLatestCheckpointSequenceNumber(ctx)
}

func (c Client) GetCheckpoints(ctx context.Context, req suimodels.SuiGetCheckpointsRequest) (suimodels.PaginatedCheckpointsResponse, error) {
	return c.rpc.SuiGetCheckpoints(ctx, req)
}
