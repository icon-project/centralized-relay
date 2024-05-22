package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
)

type IClient interface {
	GetLatestBlockHeight(ctx context.Context) (uint64, error)
	GetLatestBlockHash(ctx context.Context) (*solana.Hash, error)

	GetAccountInfo(ctx context.Context, accountId solana.PublicKey) (*solrpc.Account, error)

	GetSignatureStatus(
		ctx context.Context,
		searchTxHistory bool,
		sign solana.Signature,
	) (*solrpc.SignatureStatusesResult, error)

	GetMinBalanceForRentExemption(
		ctx context.Context,
		dataSize uint64,
	) (uint64, error)

	SimulateTx(
		ctx context.Context,
		tx *solana.Transaction,
		opts *solrpc.SimulateTransactionOpts,
	) (*solrpc.SimulateTransactionResult, error)

	SendTx(
		ctx context.Context,
		tx *solana.Transaction,
		opts *solrpc.TransactionOpts,
	) (solana.Signature, error)
}

type Client struct {
	rpc *solrpc.Client
}

func NewClient(rpcCl *solrpc.Client) IClient {
	return Client{rpc: rpcCl}
}

func (cl Client) GetAccountInfo(ctx context.Context, accountId solana.PublicKey) (*solrpc.Account, error) {
	res, err := cl.rpc.GetAccountInfo(ctx, accountId)
	if err != nil {
		return nil, err
	}
	return res.Value, nil
}

func (cl Client) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	return cl.rpc.GetBlockHeight(ctx, solrpc.CommitmentFinalized)
}

func (cl Client) GetLatestBlockHash(ctx context.Context) (*solana.Hash, error) {
	hashRes, err := cl.rpc.GetLatestBlockhash(ctx, solrpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	return &hashRes.Value.Blockhash, nil
}

func (cl Client) SimulateTx(
	ctx context.Context,
	tx *solana.Transaction,
	opts *solrpc.SimulateTransactionOpts,
) (*solrpc.SimulateTransactionResult, error) {
	res, err := cl.rpc.SimulateTransactionWithOpts(ctx, tx, opts)
	if err != nil {
		return nil, err
	}

	if res.Value.Err != nil {
		return nil, fmt.Errorf("failed to simulate tx: %v", res.Value.Err)
	}

	return res.Value, nil
}

func (cl Client) SendTx(
	ctx context.Context,
	tx *solana.Transaction,
	opts *solrpc.TransactionOpts,
) (solana.Signature, error) {
	if opts != nil {
		return cl.rpc.SendTransactionWithOpts(ctx, tx, *opts)
	}
	return cl.rpc.SendTransaction(ctx, tx)
}

func (cl Client) GetSignatureStatus(
	ctx context.Context,
	searchTxHistory bool,
	sign solana.Signature,
) (*solrpc.SignatureStatusesResult, error) {
	res, err := cl.rpc.GetSignatureStatuses(ctx, searchTxHistory, sign)
	if err != nil {
		return nil, err
	}
	if len(res.Value) > 0 {
		return res.Value[0], nil
	}
	return nil, fmt.Errorf("tx signature result not found")
}

func (cl Client) GetMinBalanceForRentExemption(
	ctx context.Context,
	dataSize uint64,
) (uint64, error) {
	return cl.rpc.GetMinimumBalanceForRentExemption(ctx, dataSize, "")
}
