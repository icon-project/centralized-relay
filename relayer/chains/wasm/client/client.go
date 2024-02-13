package client

import (
	"context"
	"sync"

	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	sdkClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	txTypes "github.com/cosmos/cosmos-sdk/types/tx"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	"github.com/spf13/pflag"
)

type IClient interface {
	HTTP(rpcUrl string) (http.HTTP, error)
	GetLatestBlockHeight(ctx context.Context) (uint64, error)
	GetTransactionReceipt(ctx context.Context, txHash string) (*txTypes.GetTxResponse, error)
	GetBalance(ctx context.Context, addr string, denomination string) (*sdkTypes.Coin, error)

	BuildTxFactory() (tx.Factory, error)
	CalculateGas(txf tx.Factory, msgs []sdkTypes.Msg) (*txTypes.SimulateResponse, uint64, error)
	PrepareTx(ctx context.Context, txf tx.Factory, msgs []sdkTypes.Msg) ([]byte, error)
	BroadcastTx(txBytes []byte) (*sdkTypes.TxResponse, error)

	TxSearch(ctx context.Context, param types.TxSearchParam) (*coretypes.ResultTxSearch, error)

	GetAccountInfo(ctx context.Context, accountAddr string) (sdkTypes.AccountI, error)

	QuerySmartContract(ctx context.Context, contractAddress string, queryData []byte) (*wasmTypes.QuerySmartContractStateResponse, error)
}

type Client struct {
	context sdkClient.Context
	txMutex *sync.Mutex
}

func New(clientCtx sdkClient.Context) Client {
	return Client{clientCtx, &sync.Mutex{}}
}

func (cl Client) BuildTxFactory() (tx.Factory, error) {
	txf, err := tx.NewFactoryCLI(cl.context, &pflag.FlagSet{})
	if err != nil {
		return tx.Factory{}, err
	}

	txf = txf.
		WithTxConfig(cl.context.TxConfig).
		WithKeybase(cl.context.Keyring).
		WithFeePayer(cl.context.FeePayer).
		WithChainID(cl.context.ChainID).
		WithSimulateAndExecute(cl.context.Simulate)

	return txf, nil
}

func (cl *Client) CalculateGas(txf tx.Factory, msgs []sdkTypes.Msg) (*txTypes.SimulateResponse, uint64, error) {
	return tx.CalculateGas(cl.context, txf, msgs...)
}

func (cl *Client) PrepareTx(ctx context.Context, txf tx.Factory, msgs []sdkTypes.Msg) ([]byte, error) {
	txBuilder, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, err
	}

	if err = tx.Sign(ctx, txf, cl.context.FromName, txBuilder, true); err != nil {
		return nil, err
	}

	return cl.context.TxConfig.TxEncoder()(txBuilder.GetTx())
}

func (cl Client) BroadcastTx(txBytes []byte) (*sdkTypes.TxResponse, error) {
	return cl.context.BroadcastTx(txBytes)
}

func (cl *Client) HTTP(rpcUrl string) (http.HTTP, error) {
	httpClient, err := http.New(rpcUrl, "/websocket")
	return *httpClient, err
}

func (cl *Client) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	nodeStatus, err := cl.context.Client.Status(ctx)
	if err != nil {
		return 0, err
	}

	return uint64(nodeStatus.SyncInfo.LatestBlockHeight), nil
}

func (cl *Client) GetTransactionReceipt(ctx context.Context, txHash string) (*txTypes.GetTxResponse, error) {
	serviceClient := txTypes.NewServiceClient(cl.context.GRPCClient)
	return serviceClient.GetTx(ctx, &txTypes.GetTxRequest{Hash: txHash})
}

func (cl *Client) GetBalance(ctx context.Context, addr string, denomination string) (*sdkTypes.Coin, error) {
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

// Create new AccountI
func (cl *Client) CreateAccount(ctx context.Context, name string, mnemonic string) (sdkTypes.AccountI, error) {
	// Create a new key
	key, err := cl.context.Keyring.NewAccount(name, mnemonic, "", "", sdkClient.DefaultKeyPass, sdkClient.DefaultAlgo)
	if err != nil {
		return nil, err
	}
}

func (cl *Client) GetAccountInfo(ctx context.Context, accountAddr string) (sdkTypes.AccountI, error) {
	qc := authTypes.NewQueryClient(cl.context.GRPCClient)

	res, err := qc.Account(
		ctx,
		&authTypes.QueryAccountRequest{Address: accountAddr},
	)
	if err != nil {
		return nil, err
	}

	var account sdkTypes.AccountI

	if err := cl.context.InterfaceRegistry.UnpackAny(res.Account, &account); err != nil {
		return nil, err
	}

	return account, nil
}

func (cl *Client) TxSearch(ctx context.Context, param types.TxSearchParam) (*coretypes.ResultTxSearch, error) {
	return cl.context.Client.TxSearch(ctx, param.BuildQuery(), param.Prove, param.Page, param.PerPage, param.OrderBy)
}

func (cl *Client) QuerySmartContract(ctx context.Context, contractAddress string, queryData []byte) (*wasmTypes.QuerySmartContractStateResponse, error) {
	queryClient := wasmTypes.NewQueryClient(cl.context)
	return queryClient.SmartContractState(ctx, &wasmTypes.QuerySmartContractStateRequest{
		Address:   contractAddress,
		QueryData: queryData,
	})
}
