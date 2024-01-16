package evm

import (
	"context"
	"fmt"
	"math"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	bridgeContract "github.com/icon-project/centralized-relay/relayer/chains/evm/abi"
	types "github.com/icon-project/centralized-relay/relayer/chains/evm/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	RPCCallRetry             = 5
	MaxGasPriceInceremtRetry = 5
	GasPriceRatio            = 10.0
)

func newClient(url string, contractAddress string, l *zap.Logger) (IClient, error) {
	clrpc, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	cleth := ethclient.NewClient(clrpc)

	bridgeContract, err := bridgeContract.NewAbi(common.HexToAddress(contractAddress), cleth)
	if err != nil {
		return nil, fmt.Errorf("error occured when creating eth client: %v ", err)
	}

	// getting the chain id
	evmChainId, err := cleth.ChainID(context.TODO())
	if err != nil {
		return nil, err
	}

	return &Client{
		log:            l,
		rpc:            clrpc,
		eth:            cleth,
		EVMChainID:     evmChainId,
		bridgeContract: bridgeContract,
	}, nil
}

// grouped rpc api clients
type Client struct {
	log      *zap.Logger
	rpc      *rpc.Client
	eth      *ethclient.Client
	verifier *Client
	// evm chain ID
	EVMChainID     *big.Int
	bridgeContract *bridgeContract.Abi
}

type IClient interface {
	Log() *zap.Logger
	GetBalance(ctx context.Context, hexAddr string) (*big.Int, error)
	GetBlockNumber() (uint64, error)
	GetBlockByHash(hash common.Hash) (*types.Block, error)
	GetHeaderByHeight(ctx context.Context, height *big.Int) (*ethTypes.Header, error)
	GetBlockReceipts(hash common.Hash) (ethTypes.Receipts, error)
	GetMedianGasPriceForBlock(ctx context.Context) (gasPrice *big.Int, gasHeight *big.Int, err error)
	GetChainID() *big.Int

	// ethClient
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	TransactionByHash(ctx context.Context, blockHash common.Hash) (tx *ethTypes.Transaction, isPending bool, err error)
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*ethTypes.Receipt, error)
	TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
	TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*ethTypes.Transaction, error)

	// transaction
	SendTransaction(ctx context.Context, tx *ethTypes.Transaction) error

	// abiContract
	ParseMessage(log ethTypes.Log) (*bridgeContract.AbiMessage, error)
	SendMessage(opts *bind.TransactOpts, _to string, _svc string, _sn *big.Int, _msg []byte) (*ethTypes.Transaction, error)
	ReceiveMessage(opts *bind.TransactOpts, srcNID string, sn *big.Int, msg []byte) (*ethTypes.Transaction, error)
	MessageReceived(opts *bind.CallOpts, srcNetwork string, _connSn *big.Int) (bool, error)
	SetAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*ethTypes.Transaction, error)
	RevertMessage(opts *bind.TransactOpts, sn *big.Int) (*ethTypes.Transaction, error)
}

func (cl *Client) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return cl.eth.NonceAt(ctx, account, blockNumber)
}

func (cl *Client) ParseMessage(log ethTypes.Log) (*bridgeContract.AbiMessage, error) {
	return cl.bridgeContract.ParseMessage(log)
}

func (cl *Client) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	return cl.eth.TransactionCount(ctx, blockHash)
}

func (cl *Client) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*ethTypes.Transaction, error) {
	return cl.eth.TransactionInBlock(ctx, blockHash, index)
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

func (cl *Client) GetBlockByHash(hash common.Hash) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	var hb types.Block
	err := cl.rpc.CallContext(ctx, &hb, "eth_getBlockByHash", hash, false)
	if err != nil {
		return nil, err
	}
	return &hb, nil
}

func (cl *Client) GetHeaderByHeight(ctx context.Context, height *big.Int) (*ethTypes.Header, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()
	return cl.eth.HeaderByNumber(ctx, height)
}

func (cl *Client) GetBlockReceipts(hash common.Hash) (ethTypes.Receipts, error) {
	c := IClient(cl)

	hb, err := c.GetBlockByHash(hash)
	if err != nil {
		return nil, err
	}
	if hb.GasUsed == "0x0" || len(hb.Transactions) == 0 {
		return nil, nil
	}
	txhs := hb.Transactions
	// fetch all txn receipts concurrently
	type rcq struct {
		txh   string
		v     *ethTypes.Receipt
		err   error
		retry int
	}
	qch := make(chan *rcq, len(txhs))
	for _, txh := range txhs {
		qch <- &rcq{txh, nil, nil, RPCCallRetry}
	}
	rmap := make(map[string]*ethTypes.Receipt)
	for q := range qch {
		switch {
		case q.err != nil:
			if q.retry == 0 {
				return nil, q.err
			}
			q.retry--
			q.err = nil
			qch <- q
		case q.v != nil:
			rmap[q.txh] = q.v
			if len(rmap) == cap(qch) {
				close(qch)
			}
		default:
			go func(q *rcq) {
				defer func() { qch <- q }()
				ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
				defer cancel()
				if q.v == nil {
					q.v = &ethTypes.Receipt{}
				}
				q.v, err = c.TransactionReceipt(ctx, common.HexToHash(q.txh))
				if q.err != nil {
					q.err = errors.Wrapf(q.err, "getTranasctionReceipt: %v", q.err)
				}
			}(q)
		}
	}
	receipts := make(ethTypes.Receipts, 0, len(txhs))
	for _, txh := range txhs {
		if r, ok := rmap[txh]; ok {
			receipts = append(receipts, r)
		}
	}
	return receipts, nil
}

func (cl *Client) GetMedianGasPriceForBlock(ctx context.Context) (gasPrice *big.Int, gasHeight *big.Int, err error) {
	c := IClient(cl)

	gasPrice = big.NewInt(0)
	header, err := c.GetHeaderByHeight(ctx, nil)
	if err != nil {
		err = errors.Wrapf(err, "GetHeaderByNumber(height:latest) Err: %v", err)
		return
	}
	height := header.Number
	txnCount, err := c.TransactionCount(ctx, header.Hash())
	if err != nil {
		err = errors.Wrapf(err, "GetTransactionCount(height:%v, headerHash: %v) Err: %v", height, header.Hash(), err)
		return
	} else if err == nil && txnCount == 0 {
		return nil, nil, fmt.Errorf("TransactionCount is zero for height(%v, headerHash %v)", height, header.Hash())
	}
	// txnF, err := c.eth.TransactionInBlock(ctx, header.Hash(), 0)
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "GetTransactionInBlock(headerHash: %v, height: %v Index: %v) Err: %v", header.Hash(), height, 0, err)
	// }
	txnS, err := c.TransactionInBlock(ctx, header.Hash(), uint(math.Floor(float64(txnCount)/2)))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "GetTransactionInBlock(headerHash: %v, height: %v Index: %v) Err: %v", header.Hash(), height, txnCount-1, err)
	}

	gasPrice = txnS.GasPrice()
	gasHeight = header.Number
	return
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

func (c *Client) SendMessage(opts *bind.TransactOpts, _to string, _svc string, _sn *big.Int, _msg []byte) (*ethTypes.Transaction, error) {
	return c.bridgeContract.SendMessage(opts, _to, _svc, _sn, _msg)
}

func (c *Client) ReceiveMessage(opts *bind.TransactOpts, srcNID string, sn *big.Int, msg []byte) (*ethTypes.Transaction, error) {
	return c.bridgeContract.RecvMessage(opts, srcNID, sn, msg)
}

func (c *Client) SendTransaction(ctx context.Context, tx *ethTypes.Transaction) error {
	return c.eth.SendTransaction(ctx, tx)
}

func (c *Client) MessageReceived(opts *bind.CallOpts, srcNetwork string, _connSn *big.Int) (bool, error) {
	return c.bridgeContract.GetReceipt(opts, srcNetwork, _connSn)
}

func (c *Client) SetAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*ethTypes.Transaction, error) {
	return c.bridgeContract.SetAdmin(opts, newAdmin)
}

func (c *Client) RevertMessage(opts *bind.TransactOpts, sn *big.Int) (*ethTypes.Transaction, error) {
	return c.bridgeContract.RevertMessage(opts, sn)
}
