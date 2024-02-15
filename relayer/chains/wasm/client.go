package wasm

import (
	"context"
	"sync"

	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	sdkClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	txTypes "github.com/cosmos/cosmos-sdk/types/tx"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	"github.com/spf13/pflag"
)

type IClient interface {
	HTTP(rpcUrl string) (*http.HTTP, error)
	GetLatestBlockHeight(ctx context.Context) (uint64, error)
	GetTransactionReceipt(ctx context.Context, txHash string) (*txTypes.GetTxResponse, error)
	GetBalance(ctx context.Context, addr string, denomination string) (*sdkTypes.Coin, error)
	BuildTxFactory() (*tx.Factory, error)
	CalculateGas(txf tx.Factory, msgs []sdkTypes.Msg) (*txTypes.SimulateResponse, uint64, error)
	PrepareTx(ctx context.Context, txf tx.Factory, msgs []sdkTypes.Msg) ([]byte, error)
	BroadcastTx(txBytes []byte) (*sdkTypes.TxResponse, error)
	TxSearch(ctx context.Context, param types.TxSearchParam) (*coretypes.ResultTxSearch, error)
	GetAccountInfo(ctx context.Context, addr string) (*sdkTypes.AccountI, error)
	QuerySmartContract(ctx context.Context, contractAddress string, queryData []byte) (*wasmTypes.QuerySmartContractStateResponse, error)
	Context() *sdkClient.Context
}

type Client struct {
	context *sdkClient.Context
	txMutex *sync.Mutex
}

func newClient(ctx *sdkClient.Context) *Client {
	return &Client{ctx, new(sync.Mutex)}
}

func (c *Client) BuildTxFactory() (*tx.Factory, error) {
	txf, err := tx.NewFactoryCLI(*c.context, &pflag.FlagSet{})
	if err != nil {
		return nil, err
	}

	txf = txf.
		WithTxConfig(c.context.TxConfig).
		WithKeybase(c.context.Keyring).
		WithFeePayer(c.context.FeePayer).
		WithChainID(c.context.ChainID).
		WithSimulateAndExecute(c.context.Simulate)

	return &txf, nil
}

func (c *Client) CalculateGas(txf tx.Factory, msgs []sdkTypes.Msg) (*txTypes.SimulateResponse, uint64, error) {
	return tx.CalculateGas(c.context, txf, msgs...)
}

func (c *Client) PrepareTx(ctx context.Context, txf tx.Factory, msgs []sdkTypes.Msg) ([]byte, error) {
	txBuilder, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, err
	}

	if err = tx.Sign(ctx, txf, c.context.FromName, txBuilder, true); err != nil {
		return nil, err
	}

	return c.context.TxConfig.TxEncoder()(txBuilder.GetTx())
}

func (c *Client) BroadcastTx(txBytes []byte) (*sdkTypes.TxResponse, error) {
	return c.context.BroadcastTx(txBytes)
}

func (c *Client) HTTP(rpcUrl string) (*http.HTTP, error) {
	return http.New(rpcUrl, "/websocket")
}

func (c *Client) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	nodeStatus, err := c.context.Client.Status(ctx)
	if err != nil {
		return 0, err
	}

	return uint64(nodeStatus.SyncInfo.LatestBlockHeight), nil
}

func (c *Client) GetTransactionReceipt(ctx context.Context, txHash string) (*txTypes.GetTxResponse, error) {
	serviceClient := txTypes.NewServiceClient(c.context.GRPCClient)
	return serviceClient.GetTx(ctx, &txTypes.GetTxRequest{Hash: txHash})
}

func (c *Client) GetBalance(ctx context.Context, addr string, denomination string) (*sdkTypes.Coin, error) {
	queryClient := bankTypes.NewQueryClient(c.context.GRPCClient)

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
func (c *Client) CreateAccount(ctx context.Context, name string) (*sdkTypes.AccountI, error) {
	record, t, err := c.context.Keyring.NewMnemonic(uid, language, hdPath, bip39Passphrase, algo)
	if err != nil {
		return nil, err
	}
	key, err := c.context.Keyring.NewAccount(name, mnemonic, "", "", keyring.SigningAlgoList{}, sdkClient.DefaultAlgo)
	if err != nil {
		return nil, err
	}
}

func (c *Client) GetAccountInfo(ctx context.Context, addr string) (*sdkTypes.AccountI, error) {
	qc := authTypes.NewQueryClient(c.context.GRPCClient)
	res, err := qc.Account(ctx, &authTypes.QueryAccountRequest{Address: addr})
	if err != nil {
		return nil, err
	}

	account := new(sdkTypes.AccountI)

	if err := c.context.InterfaceRegistry.UnpackAny(res.Account, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (c *Client) TxSearch(ctx context.Context, param types.TxSearchParam) (*coretypes.ResultTxSearch, error) {
	return c.context.Client.TxSearch(ctx, param.BuildQuery(), param.Prove, param.Page, param.PerPage, param.OrderBy)
}

func (c *Client) QuerySmartContract(ctx context.Context, contractAddress string, queryData []byte) (*wasmTypes.QuerySmartContractStateResponse, error) {
	queryClient := wasmTypes.NewQueryClient(cl.context)
	return queryClient.SmartContractState(ctx, &wasmTypes.QuerySmartContractStateRequest{
		Address:   contractAddress,
		QueryData: queryData,
	})
}

func (c *Client) Context() *sdkClient.Context {
	return c.context
}
