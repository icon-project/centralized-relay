package client

import (
	"context"
	"fmt"

	suimodels "github.com/block-vision/sui-go-sdk/models"
	suisdk "github.com/block-vision/sui-go-sdk/sui"
	"go.uber.org/zap"
)

type Client struct {
	rpc suisdk.ISuiAPI
	log *zap.Logger
}

func NewClient(rpcClient suisdk.ISuiAPI, l *zap.Logger) *Client {
	return &Client{
		rpc: rpcClient,
		log: l,
	}
}

func (c Client) GetLatestCheckpointSeq(ctx context.Context) (uint64, error) {
	return c.rpc.SuiGetLatestCheckpointSequenceNumber(ctx)
}

func (c Client) GetCheckpoints(ctx context.Context, req suimodels.SuiGetCheckpointsRequest) (suimodels.PaginatedCheckpointsResponse, error) {
	return c.rpc.SuiGetCheckpoints(ctx, req)
}

func (c *Client) GetBalance(ctx context.Context, addr string) ([]suimodels.CoinData, error) {
	result, err := c.rpc.SuiXGetAllCoins(ctx, suimodels.SuiXGetAllCoinsRequest{
		Owner: addr,
	})
	if err != nil {
		c.log.Error(fmt.Sprintf("error getting balance for address %s", addr), zap.Error(err))
		return nil, err
	}
	return result.Data, nil
}
