package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	bridgeContract "github.com/icon-project/centralized-relay/relayer/chains/evm/abi"
	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	DefaultPollingInterval = time.Second * 30
	MaximumPollTry         = 15
)

func newClient(ctx context.Context, connectionContract, XcallContract common.Address, rpcUrl, websocketUrl string, l *zap.Logger) (IClient, error) {
	cleth, err := ethclient.DialContext(ctx, websocketUrl)
	if err != nil {
		return nil, err
	}

	// debug: use rpc instead of websocket
	cleths, err := ethclient.DialContext(ctx, rpcUrl)
	if err != nil {
		return nil, err
	}

	connection, err := bridgeContract.NewConnection(connectionContract, cleths)
	if err != nil {
		return nil, fmt.Errorf("error occured when creating connection cobtract: %v ", err)
	}

	xcall, err := bridgeContract.NewXcall(XcallContract, cleths)
	if err != nil {
		return nil, fmt.Errorf("error occured when creating eth client: %v ", err)
	}

	// getting the chain id
	evmChainId, err := cleth.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:        l,
		eth:        cleth,
		EVMChainID: evmChainId,
		connection: connection,
		xcall:      xcall,
	}, nil
}

// grouped rpc api clients
type Client struct {
	log        *zap.Logger
	eth        *ethclient.Client
	EVMChainID *big.Int
	connection *bridgeContract.Connection
	xcall      *bridgeContract.Xcall
}

type IClient interface {
	Log() *zap.Logger
	GetBalance(ctx context.Context, hexAddr string) (*big.Int, error)
	GetBlockNumber() (uint64, error)
	GetHeaderByHeight(ctx context.Context, height *big.Int) (*ethTypes.Header, error)
	GetLogs(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error)
	GetChainID() *big.Int

	// ethClient
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	SuggestGasTip(ctx context.Context) (*big.Int, error)
	PendingNonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	TransactionByHash(ctx context.Context, blockHash common.Hash) (tx *ethTypes.Transaction, isPending bool, err error)
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*ethTypes.Receipt, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	SendTransaction(ctx context.Context, tx *ethTypes.Transaction) error
	Subscribe(ctx context.Context, q ethereum.FilterQuery, ch chan<- ethTypes.Log) (ethereum.Subscription, error)

	// abiContract for connection
	ParseConnectionMessage(log ethTypes.Log) (*bridgeContract.ConnectionMessage, error)
	SendMessage(opts *bind.TransactOpts, _to string, _svc string, _sn *big.Int, _msg []byte) (*ethTypes.Transaction, error)
	ReceiveMessage(opts *bind.TransactOpts, srcNID string, sn *big.Int, msg []byte) (*ethTypes.Transaction, error)
	MessageReceived(opts *bind.CallOpts, srcNetwork string, _connSn *big.Int) (bool, error)
	SetAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*ethTypes.Transaction, error)
	RevertMessage(opts *bind.TransactOpts, sn *big.Int) (*ethTypes.Transaction, error)
	GetFee(opts *bind.CallOpts, networkID string) (*big.Int, error)
	SetFee(opts *bind.TransactOpts, src string, msg, res *big.Int) (*ethTypes.Transaction, error)
	ClaimFee(opts *bind.TransactOpts) (*ethTypes.Transaction, error)

	// abiContract for xcall
	ParseXcallMessage(log ethTypes.Log) (*bridgeContract.XcallCallMessage, error)
	ExecuteCall(opts *bind.TransactOpts, reqID *big.Int, data []byte) (*ethTypes.Transaction, error)
	ExecuteRollback(opts *bind.TransactOpts, sn *big.Int) (*ethTypes.Transaction, error)
}

func (c *Client) PendingNonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	nonce, err := c.eth.PendingNonceAt(ctx, account)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetUint64(nonce), nil
}

func (cl *Client) TransactionByHash(ctx context.Context, blockHash common.Hash) (tx *ethTypes.Transaction, isPending bool, err error) {
	return cl.eth.TransactionByHash(ctx, blockHash)
}

func (cl *Client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*ethTypes.Receipt, error) {
	return cl.eth.TransactionReceipt(ctx, txHash)
}

func (cl *Client) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return cl.eth.CallContract(ctx, msg, blockNumber)
}

func (cl *Client) GetBalance(ctx context.Context, hexAddr string) (*big.Int, error) {
	if !common.IsHexAddress(hexAddr) {
		return nil, fmt.Errorf("invalid hex address: %v", hexAddr)
	}
	return cl.eth.BalanceAt(ctx, common.HexToAddress(hexAddr), nil)
}

func (cl *Client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
	return cl.eth.FilterLogs(ctx, q)
}

func (cl *Client) GetBlockNumber() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	bn, err := cl.eth.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return bn, nil
}

func (cl *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return cl.eth.SuggestGasPrice(ctx)
}

func (cl *Client) SuggestGasTip(ctx context.Context) (*big.Int, error) {
	return cl.eth.SuggestGasTipCap(ctx)
}

func (cl *Client) GetLogs(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
	return cl.eth.FilterLogs(ctx, q)
}

func (cl *Client) GetHeaderByHeight(ctx context.Context, height *big.Int) (*ethTypes.Header, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()
	return cl.eth.HeaderByNumber(ctx, height)
}

func (c *Client) GetChainID() *big.Int {
	return c.EVMChainID
}

func (c *Client) GetEthClient() *ethclient.Client {
	return c.eth
}

func (c *Client) Log() *zap.Logger {
	return c.log
}

func (cl *Client) ParseConnectionMessage(log ethTypes.Log) (*bridgeContract.ConnectionMessage, error) {
	return cl.connection.ParseMessage(log)
}

func (c *Client) SendMessage(opts *bind.TransactOpts, _to string, _svc string, _sn *big.Int, _msg []byte) (*ethTypes.Transaction, error) {
	return c.connection.SendMessage(opts, _to, _svc, _sn, _msg)
}

func (c *Client) ReceiveMessage(opts *bind.TransactOpts, srcNID string, sn *big.Int, msg []byte) (*ethTypes.Transaction, error) {
	return c.connection.RecvMessage(opts, srcNID, sn, msg)
}

func (c *Client) SendTransaction(ctx context.Context, tx *ethTypes.Transaction) error {
	return c.eth.SendTransaction(ctx, tx)
}

func (c *Client) MessageReceived(opts *bind.CallOpts, srcNetwork string, _connSn *big.Int) (bool, error) {
	return c.connection.GetReceipt(opts, srcNetwork, _connSn)
}

func (c *Client) SetAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*ethTypes.Transaction, error) {
	return c.connection.SetAdmin(opts, newAdmin)
}

func (c *Client) RevertMessage(opts *bind.TransactOpts, sn *big.Int) (*ethTypes.Transaction, error) {
	return c.connection.RevertMessage(opts, sn)
}

func (c *Client) ParseXcallMessage(log ethTypes.Log) (*bridgeContract.XcallCallMessage, error) {
	return c.xcall.ParseCallMessage(log)
}

func (c *Client) ExecuteCall(opts *bind.TransactOpts, reqID *big.Int, data []byte) (*ethTypes.Transaction, error) {
	return c.xcall.ExecuteCall(opts, reqID, data)
}

func (c *Client) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	return c.eth.EstimateGas(ctx, msg)
}

// GetFee
func (c *Client) GetFee(opts *bind.CallOpts, networkID string) (*big.Int, error) {
	return c.connection.GetFee(opts, networkID, true)
}

// SetFee
func (c *Client) SetFee(opts *bind.TransactOpts, src string, msg, res *big.Int) (*ethTypes.Transaction, error) {
	return c.connection.SetFee(opts, src, msg, res)
}

// ClaimFee
func (c *Client) ClaimFee(opts *bind.TransactOpts) (*ethTypes.Transaction, error) {
	return c.connection.ClaimFees(opts)
}

// ExecuteRollback
func (c *Client) ExecuteRollback(opts *bind.TransactOpts, sn *big.Int) (*ethTypes.Transaction, error) {
	return c.xcall.ExecuteRollback(opts, sn)
}

// Subscribe
func (c *Client) Subscribe(ctx context.Context, q ethereum.FilterQuery, ch chan<- ethTypes.Log) (ethereum.Subscription, error) {
	return c.eth.SubscribeFilterLogs(ctx, q, ch)
}
