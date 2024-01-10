package client

import (
	"context"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdkClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	txTypes "github.com/cosmos/cosmos-sdk/types/tx"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
)

type IClient interface {
	Context() sdkClient.Context
	GetLatestBlockHeight(ctx context.Context) (uint64, error)
	GetTransactionReceipt(ctx context.Context, txHash string) (*txTypes.GetTxResponse, error)
	GetBalance(ctx context.Context, addr string, denomination string) (*sdkTypes.Coin, error)
	GetMessages(ctx context.Context, param types.TxSearchParam) ([]*relayTypes.Message, error)

	QuerySmartContract(ctx context.Context, contractAddress string, queryData []byte) (*wasmTypes.QuerySmartContractStateResponse, error)
	SendTx(ctx context.Context, tf tx.Factory, messages []sdkTypes.Msg) (*sdkTypes.TxResponse, error)
}

type Client struct {
	context sdkClient.Context
}

func New(clientCtx sdkClient.Context) Client {
	return Client{clientCtx}
}

func (cl Client) Context() sdkClient.Context {
	return cl.context
}

func (cl Client) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	nodeStatus, err := cl.context.Client.Status(ctx)
	if err != nil {
		return 0, err
	}

	return uint64(nodeStatus.SyncInfo.LatestBlockHeight), nil
}

func (cl Client) GetTransactionReceipt(ctx context.Context, txHash string) (*txTypes.GetTxResponse, error) {
	serviceClient := txTypes.NewServiceClient(cl.context.GRPCClient)
	return serviceClient.GetTx(ctx, &txTypes.GetTxRequest{Hash: txHash})
}

func (cl Client) GetBalance(ctx context.Context, addr string, denomination string) (*sdkTypes.Coin, error) {
	queryClient := bankTypes.NewQueryClient(cl.context.GRPCClient)

	res, err := queryClient.Balance(ctx, &bankTypes.QueryBalanceRequest{
		Address: addr,
		Denom:   denomination,
	})
	if err != nil {
		return nil, err
	}
	return res.Balance, nil
}

func (cl Client) GetMessages(ctx context.Context, param types.TxSearchParam) ([]*relayTypes.Message, error) {
	result, err := cl.context.Client.TxSearch(ctx, param.Query, param.Prove, param.Page, param.PerPage, param.OrderBy)
	if err != nil {
		return nil, err
	}

	messages := make([]*relayTypes.Message, len(result.Txs))

	//Todo message parse from tx

	return messages, nil
}

func (cl Client) QuerySmartContract(ctx context.Context, contractAddress string, queryData []byte) (*wasmTypes.QuerySmartContractStateResponse, error) {
	queryClient := wasmTypes.NewQueryClient(cl.context)
	return queryClient.SmartContractState(ctx, &wasmTypes.QuerySmartContractStateRequest{
		Address:   contractAddress,
		QueryData: queryData,
	})
}

func (cl Client) SendTx(ctx context.Context, tf tx.Factory, messages []sdkTypes.Msg) (*sdkTypes.TxResponse, error) {
	txBuilder, err := tf.BuildUnsignedTx(messages...)
	if err != nil {
		return nil, err
	}

	if err := tx.Sign(ctx, tf, cl.context.GetFromName(), txBuilder, true); err != nil {
		return nil, err
	}

	txBytes, err := cl.context.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	return cl.context.BroadcastTx(txBytes)
}
