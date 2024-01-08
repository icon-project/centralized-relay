package client

import (
	"context"
	abiTypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	txTypes "github.com/cosmos/cosmos-sdk/types/tx"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type IClient interface {
	GetLatestBlock(ctx context.Context) (uint64, error)
	GetTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error)
	GetBalance(ctx context.Context, addr string) (*types.Coin, error)
}

type Client struct {
	logger          *zap.Logger
	grpcConn        *grpc.ClientConn
	contractAddress string
	chainID         string
}

func (cl *Client) GetLatestBlock(ctx context.Context) (uint64, error) {
	serviceClient := cmtservice.NewServiceClient(cl.grpcConn)
	res, err := serviceClient.GetLatestBlock(ctx, &cmtservice.GetLatestBlockRequest{})
	if err != nil {
		return 0, err
	}

	return uint64(res.GetBlock().Header.Height), nil
}

func (cl *Client) GetTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error) {
	serviceClient := txTypes.NewServiceClient(cl.grpcConn)
	res, err := serviceClient.GetTx(ctx, &txTypes.GetTxRequest{Hash: txHash})
	if err != nil {
		return nil, err
	}
	return &types.Receipt{
		TxHash: txHash,
		Height: uint64(res.TxResponse.Height),
		Status: abiTypes.CodeTypeOK == res.TxResponse.Code,
	}, nil
}

func (cl *Client) GetBalance(ctx context.Context, addr string) (*types.Coin, error) {
	queryClient := bankTypes.NewQueryClient(cl.grpcConn)

	res, err := queryClient.Balance(ctx, &bankTypes.QueryBalanceRequest{
		Address: addr,
		Denom:   "s",
	})
	if err != nil {
		return nil, err
	}

	return &types.Coin{
		Denom:  res.GetBalance().Denom,
		Amount: res.GetBalance().Amount.Uint64(),
	}, nil
}
