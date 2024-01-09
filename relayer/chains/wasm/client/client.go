package client

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client/grpc/node"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	txTypes "github.com/cosmos/cosmos-sdk/types/tx"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	wasmTypes "github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
	tmHttp "github.com/tendermint/tendermint/rpc/client/http"
	"google.golang.org/grpc"
)

type IClient interface {
	GetLatestBlockHeight(ctx context.Context) (uint64, error)
	GetTransactionReceipt(ctx context.Context, txHash string) (*txTypes.GetTxResponse, error)
	GetBalance(ctx context.Context, addr string) (*sdkTypes.Coin, error)
	GetMessages(ctx context.Context, param wasmTypes.TxSearchParam) ([]*relayTypes.Message, error)
}

type Client struct {
	grpcConn        *grpc.ClientConn
	tmHttpClient    *tmHttp.HTTP
	contractAddress string
}

func (cl *Client) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	nodeStatus, err := cl.getNodeStatus(ctx)
	if err != nil {
		return 0, err
	}
	return nodeStatus.Height, nil
}

func (cl *Client) GetTransactionReceipt(ctx context.Context, txHash string) (*txTypes.GetTxResponse, error) {
	serviceClient := txTypes.NewServiceClient(cl.grpcConn)
	return serviceClient.GetTx(ctx, &txTypes.GetTxRequest{Hash: txHash})
}

func (cl *Client) GetBalance(ctx context.Context, addr string) (*sdkTypes.Coin, error) {
	queryClient := bankTypes.NewQueryClient(cl.grpcConn)

	res, err := queryClient.Balance(ctx, &bankTypes.QueryBalanceRequest{
		Address: addr,
		Denom:   "s",
	})
	if err != nil {
		return nil, err
	}
	return res.Balance, nil
}

func (cl *Client) getNodeStatus(ctx context.Context) (*node.StatusResponse, error) {
	serviceClient := node.NewServiceClient(cl.grpcConn)
	return serviceClient.Status(ctx, &node.StatusRequest{})
}

func (cl *Client) GetMessages(ctx context.Context, param wasmTypes.TxSearchParam) ([]*relayTypes.Message, error) {
	result, err := cl.tmHttpClient.TxSearch(ctx, param.Query, param.Prove, param.Page, param.PerPage, param.OrderBy)
	if err != nil {
		return nil, err
	}

	messages := make([]*relayTypes.Message, len(result.Txs))

	//for _, tx := range result.Txs {
	//	tx.TxResult.Log
	//}

	return messages, nil
}
