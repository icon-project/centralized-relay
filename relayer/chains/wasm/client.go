package wasm

import (
	"context"
	"strconv"

	jsoniter "github.com/json-iterator/go"

	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	sdkClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/go-bip39"

	txTypes "github.com/cosmos/cosmos-sdk/types/tx"

	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	"github.com/spf13/pflag"
)

type IClient interface {
	HTTP(rpcUrl string) (*http.HTTP, error)
	IsConnected() bool
	Reconnect() error
	GetLatestBlockHeight(ctx context.Context) (uint64, error)
	GetTransactionReceipt(ctx context.Context, txHash string) (*txTypes.GetTxResponse, error)
	GetBalance(ctx context.Context, addr string, denomination string) (*sdkTypes.Coin, error)
	BuildTxFactory() (tx.Factory, error)
	EstimateGas(txf tx.Factory, msgs ...sdkTypes.Msg) (*txTypes.SimulateResponse, uint64, error)
	PrepareTx(ctx context.Context, txf tx.Factory, msgs ...sdkTypes.Msg) ([]byte, error)
	BroadcastTx(txBytes []byte) (*sdkTypes.TxResponse, error)
	TxSearch(ctx context.Context, param types.TxSearchParam) (*coretypes.ResultTxSearch, error)
	GetAccountInfo(ctx context.Context, addr string) (sdkTypes.AccountI, error)
	QuerySmartContract(ctx context.Context, address string, queryData []byte) (*wasmTypes.QuerySmartContractStateResponse, error)
	CreateAccount(name, pass string) (string, string, error)
	ImportArmor(uid string, armor []byte, passphrase string) error
	GetArmor(uid, passphrase string) (string, error)
	GetKey(uid string) (*keyring.Record, error)
	GetKeyByAddr(addr sdkTypes.Address) (*keyring.Record, error)
	SetAddress(account sdkTypes.AccAddress) sdkTypes.AccAddress
	Subscribe(ctx context.Context, _, query string) (<-chan coretypes.ResultEvent, error)
	Unsubscribe(ctx context.Context, _, query string) error
	GetFee(ctx context.Context, addr string, queryData []byte) (uint64, error)
}

type Client struct {
	ctx sdkClient.Context
}

func newClient(ctx *sdkClient.Context) *Client {
	return &Client{*ctx}
}

func (c *Client) BuildTxFactory() (tx.Factory, error) {
	txf, err := tx.NewFactoryCLI(c.ctx, &pflag.FlagSet{})
	if err != nil {
		return txf, err
	}
	return txf.WithSimulateAndExecute(c.ctx.Simulate), nil
}

func (c *Client) EstimateGas(txf tx.Factory, msgs ...sdkTypes.Msg) (*txTypes.SimulateResponse, uint64, error) {
	return tx.CalculateGas(c.ctx, txf, msgs...)
}

func (c *Client) PrepareTx(ctx context.Context, txf tx.Factory, msgs ...sdkTypes.Msg) ([]byte, error) {
	txBuilder, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, err
	}

	if err = tx.Sign(ctx, txf, c.ctx.FromName, txBuilder, true); err != nil {
		return nil, err
	}
	return c.ctx.TxConfig.TxEncoder()(txBuilder.GetTx())
}

func (c *Client) BroadcastTx(txBytes []byte) (*sdkTypes.TxResponse, error) {
	return c.ctx.BroadcastTx(txBytes)
}

func (c *Client) HTTP(rpcUrl string) (*http.HTTP, error) {
	return http.New(rpcUrl, "/websocket")
}

func (c *Client) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	nodeStatus, err := c.ctx.Client.Status(ctx)
	if err != nil {
		return 0, err
	}

	return uint64(nodeStatus.SyncInfo.LatestBlockHeight), nil
}

func (c *Client) GetTransactionReceipt(ctx context.Context, txHash string) (*txTypes.GetTxResponse, error) {
	serviceClient := txTypes.NewServiceClient(c.ctx.GRPCClient)
	return serviceClient.GetTx(ctx, &txTypes.GetTxRequest{Hash: txHash})
}

func (c *Client) GetBalance(ctx context.Context, addr string, denomination string) (*sdkTypes.Coin, error) {
	queryClient := bankTypes.NewQueryClient(c.ctx.GRPCClient)

	res, err := queryClient.Balance(ctx, &bankTypes.QueryBalanceRequest{
		Address: addr,
		Denom:   denomination,
	})
	if err != nil {
		return nil, err
	}
	return res.Balance, nil
}

func (c *Client) GetAccountInfo(ctx context.Context, addr string) (sdkTypes.AccountI, error) {
	qc := authTypes.NewQueryClient(c.ctx.GRPCClient)
	res, err := qc.Account(ctx, &authTypes.QueryAccountRequest{Address: addr})
	if err != nil {
		return nil, err
	}

	var account sdkTypes.AccountI

	if err := c.ctx.InterfaceRegistry.UnpackAny(res.Account, &account); err != nil {
		return nil, err
	}

	return account, nil
}

// Create new AccountI
func (c *Client) CreateAccount(uid, passphrase string) (string, string, error) {
	// create hdpath
	kb := keyring.NewInMemory(c.ctx.Codec)
	hdPath := hd.CreateHDPath(sdkTypes.CoinType, 0, 0).String()
	bip39seed, err := bip39.NewEntropy(256)
	if err != nil {
		return "", "", err
	}
	mnemonic, err := bip39.NewMnemonic(bip39seed)
	if err != nil {
		return "", "", err
	}
	algos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), algos)
	if err != nil {
		return "", "", err
	}
	key, err := kb.NewAccount(uid, mnemonic, passphrase, hdPath, algo)
	if err != nil {
		return "", "", err
	}
	armor, err := kb.ExportPrivKeyArmor(uid, passphrase)
	if err != nil {
		return "", "", err
	}
	addr, err := key.GetAddress()
	if err != nil {
		return "", "", err
	}
	return armor, addr.String(), nil
}

// Load private key from keyring
func (c *Client) ImportArmor(uid string, armor []byte, passphrase string) error {
	if _, err := c.ctx.Keyring.Key(uid); err == nil {
		return nil
	}
	return c.ctx.Keyring.ImportPrivKey(uid, string(armor), passphrase)
}

// GetPrivateKey returns private key from keyring
func (c *Client) GetArmor(uid, passphrase string) (string, error) {
	return c.ctx.Keyring.ExportPrivKeyArmor(uid, passphrase)
}

// GetKey returns key from keyring
func (c *Client) GetKey(uid string) (*keyring.Record, error) {
	return c.ctx.Keyring.Key(uid)
}

// GetAccount returns account from keyring
func (c *Client) GetKeyByAddr(addr sdkTypes.Address) (*keyring.Record, error) {
	return c.ctx.Keyring.KeyByAddress(addr)
}

func (c *Client) TxSearch(ctx context.Context, param types.TxSearchParam) (*coretypes.ResultTxSearch, error) {
	return c.ctx.Client.TxSearch(ctx, param.BuildQuery(), param.Prove, param.Page, param.PerPage, param.OrderBy)
}

// Set the address to be used for the transactions
func (c *Client) SetAddress(addr sdkTypes.AccAddress) sdkTypes.AccAddress {
	key, err := c.ctx.Keyring.KeyByAddress(addr)
	if err != nil {
		return nil
	}
	c.ctx = c.ctx.WithFromAddress(addr).WithFeePayerAddress(addr).WithFrom(addr.String()).WithFromName(key.Name).WithFeeGranterAddress(addr)
	return addr
}

func (c *Client) QuerySmartContract(ctx context.Context, address string, queryData []byte) (*wasmTypes.QuerySmartContractStateResponse, error) {
	queryClient := wasmTypes.NewQueryClient(c.ctx)
	return queryClient.SmartContractState(ctx, &wasmTypes.QuerySmartContractStateRequest{
		Address:   address,
		QueryData: queryData,
	})
}

// GetFee returns the fee set for the network
func (c *Client) GetFee(ctx context.Context, addr string, queryData []byte) (uint64, error) {
	res, err := c.QuerySmartContract(ctx, addr, queryData)
	if err != nil {
		return 0, err
	}

	var feeStr string

	if err := jsoniter.Unmarshal(res.Data, &feeStr); err != nil {
		return 0, err
	}
	fee, err := strconv.ParseUint(feeStr, 10, strconv.IntSize)
	if err != nil {
		return 0, err
	}
	return fee, nil
}

// Subscribe
func (c *Client) Subscribe(ctx context.Context, _, query string) (<-chan coretypes.ResultEvent, error) {
	return c.ctx.Client.(*http.HTTP).Subscribe(ctx, "subscribe", query)
}

// Unsubscribe
func (c *Client) Unsubscribe(ctx context.Context, _, query string) error {
	return c.ctx.Client.(*http.HTTP).Unsubscribe(ctx, "unsubscribe", query)
}

// IsConnected returns if the client is connected to the network
func (c *Client) IsConnected() bool {
	return c.ctx.Client.(*http.HTTP).IsRunning()
}

// RestartClient restarts the client
func (c *Client) Reconnect() error {
	client, err := c.HTTP(c.ctx.NodeURI)
	if err != nil {
		return err
	}
	c.ctx.Client = client
	return nil
}
