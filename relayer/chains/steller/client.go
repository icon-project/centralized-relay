package steller

import (
	"context"

	"github.com/stellar/go/clients/horizonclient"
)

type IClient interface {
	GetLatestLedgerSeq(ctx context.Context) (uint64, error)
}

type Client struct {
	horizon *horizonclient.Client
}

func NewClient(hClient *horizonclient.Client) IClient {
	return &Client{horizon: hClient}
}

func (cl *Client) GetLatestLedgerSeq(ctx context.Context) (uint64, error) {
	//Todo
	return 0, nil
}
