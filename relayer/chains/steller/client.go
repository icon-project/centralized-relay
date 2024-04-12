package steller

import (
	"context"

	"github.com/icon-project/centralized-relay/relayer/chains/steller/sorobanclient"
	"github.com/stellar/go/clients/horizonclient"
)

type IClient interface {
	GetLatestLedger(ctx context.Context) (*sorobanclient.LatestLedgerResponse, error)
}

type Client struct {
	horizon *horizonclient.Client
	soroban *sorobanclient.Client
}

func NewClient(hClient *horizonclient.Client, srbClient *sorobanclient.Client) IClient {
	return &Client{horizon: hClient, soroban: srbClient}
}

func (cl *Client) GetLatestLedger(ctx context.Context) (*sorobanclient.LatestLedgerResponse, error) {
	return cl.soroban.GetLatestLedger(ctx)
}
