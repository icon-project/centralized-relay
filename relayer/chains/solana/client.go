package solana

import (
	"context"

	solrpc "github.com/gagliardetto/solana-go/rpc"
)

type IClient interface {
	GetLatestBlockHeight(ctx context.Context) (uint64, error)
}

type Client struct {
	rpc *solrpc.Client
}

func NewClient(rpcCl *solrpc.Client) IClient {
	return Client{rpc: rpcCl}
}

func (cl Client) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	return cl.rpc.GetBlockHeight(ctx, solrpc.CommitmentFinalized)
}
