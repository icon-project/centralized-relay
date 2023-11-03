package evm

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"time"

	eth "github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/icon-project/goloop/module"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/icon-project/centralized-relay/relayer/chains/evm/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	defaultReadTimeout = 50 * time.Second //
	RPCCallRetry       = 5

	DefaultGetTransactionResultPollingInterval = 500 * time.Millisecond // 1.5sec
	BlockFinalityConfirmations                 = 10
	BlockInterval                              = 3 * time.Second
	BlockHeightPollInterval                    = BlockInterval * 5
	SyncConcurrency                            = 10
)

var txSerializeExcludes = map[string]bool{"signature": true}

type Wallet interface {
	Sign(data []byte) ([]byte, error)
	Address() string
}

func NewClient(url string, l *zap.Logger) (IClient, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	eth := ethclient.NewClient(rpcClient)
	return &Client{log: l, rpc: rpcClient, eth: eth}, nil
}

// Client grouped rpc api clients
type Client struct {
	endpoint string
	rpc      *rpc.Client
	eth      *ethclient.Client
	log      *zap.Logger
	chainID  string
}

type IClient interface {
	Log() *zap.Logger
	GetBalance(ctx context.Context, hexAddr string) (*big.Int, error)
	GetBlockNumber() (uint64, error)
	GetBlockByHash(hash common.Hash) (*ethTypes.Block, error)
	GetHeaderByHeight(ctx context.Context, height *big.Int) (*ethTypes.Header, error)
	GetBlockReceipts(hash common.Hash) (ethTypes.Receipts, error)
	GetChainID() string
	GetClient() *ethclient.Client

	FilterLogs(ctx context.Context, q eth.FilterQuery) ([]ethTypes.Log, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	TransactionByHash(ctx context.Context, blockHash common.Hash) (tx *ethTypes.Transaction, isPending bool, err error)
	Call(ctx context.Context, msg eth.CallMsg, blockNumber *big.Int) ([]byte, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*ethTypes.Receipt, error)
	TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
	TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*ethTypes.Transaction, error)
}

func (c *Client) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return c.eth.NonceAt(ctx, account, blockNumber)
}

func (c *Client) GetClient() *ethclient.Client {
	return c.eth
}

func (c *Client) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	return c.eth.TransactionCount(ctx, blockHash)
}

func (c *Client) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*ethTypes.Transaction, error) {
	return c.eth.TransactionInBlock(ctx, blockHash, index)
}

func (c *Client) TransactionByHash(ctx context.Context, blockHash common.Hash) (tx *ethTypes.Transaction, isPending bool, err error) {
	return c.eth.TransactionByHash(ctx, blockHash)
}

func (c *Client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*ethTypes.Receipt, error) {
	return c.eth.TransactionReceipt(ctx, txHash)
}

func (c *Client) Call(ctx context.Context, msg eth.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return c.eth.CallContract(ctx, msg, blockNumber)
}

func (c *Client) GetBalance(ctx context.Context, hexAddr string) (*big.Int, error) {
	if !common.IsHexAddress(hexAddr) {
		return nil, fmt.Errorf("invalid hex address: %v", hexAddr)
	}
	return c.eth.BalanceAt(ctx, common.HexToAddress(hexAddr), nil)
}

func (c *Client) FilterLogs(ctx context.Context, q eth.FilterQuery) ([]ethTypes.Log, error) {
	return c.eth.FilterLogs(ctx, q)
}

func (c *Client) GetBlockNumber() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	bn, err := c.eth.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return bn, nil
}

type ReceiverOptions struct {
	SyncConcurrency uint64           `json:"syncConcurrency"`
	Verifier        *VerifierOptions `json:"verifier"`
}

type BnOptions struct {
	StartHeight uint64
	Concurrency uint64
}

func (c *Client) GetBlockByHash(hash common.Hash) (*ethTypes.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	block, err := c.eth.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (c *Client) GetHeaderByHeight(ctx context.Context, height *big.Int) (*ethTypes.Header, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()
	return c.eth.HeaderByNumber(ctx, height)
}

func (c *Client) SignTransaction(w module.Wallet, p *ethTypes.Transaction) error {
	return nil
}

func (c *Client) SendTransaction(p *types.TransactionParam) (*types.HexBytes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	res := new(types.HexBytes)
	if err := c.rpc.CallContext(ctx, res, "eth_sendTransaction", p); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) SendTransactionAndWait(p *types.TransactionParam) (*types.HexBytes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	res := new(types.HexBytes)
	if err := c.rpc.CallContext(ctx, res, "eth_sendTransaction", p); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetTransactionResult(tx common.Hash) (*ethTypes.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	receipt, err := c.eth.TransactionReceipt(ctx, tx)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (c *Client) WaitTransactionResult(p *types.TransactionHashParam) (*types.TransactionResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	res := new(types.TransactionResult)
	if err := c.rpc.CallContext(ctx, res, "eth_waitTransactionResult", p); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) WaitForResults(ctx context.Context, hash common.Hash) (*ethTypes.Receipt, error) {
	ticker := time.NewTicker(time.Duration(DefaultGetTransactionResultPollingInterval) * time.Nanosecond)
	retryLimit := 20
	retryCounter := 0
	for {
		defer ticker.Stop()
		select {
		case <-ctx.Done():
			err := errors.New("Context Cancelled ReceiptWait Exiting ")
			return nil, err
		case <-ticker.C:
			if retryCounter >= retryLimit {
				err := errors.New("Retry Limit Exceeded while waiting for results of transaction")
				return nil, err
			}
			retryCounter++
			tx, isPendng, err := c.eth.TransactionByHash(ctx, hash)
			if err != nil {
				return nil, err
			}
			if isPendng {
				continue
			}
			receipt, err := c.TransactionReceipt(ctx, tx.Hash())
			if err != nil {
				return nil, err
			}
			return receipt, nil
		}
	}
}

func (c *Client) GetBlockByHeight(p *types.BlockHeightParam) (*types.Block, error) {
	block, err := c.eth.BlockByNumber(context.Background(), p.Height)
	if err != nil {
		return nil, err
	}
	blockInfo := &types.Block{
		Height:    block.Number().Uint64(),
		Timestamp: block.Time(),
	}
	return blockInfo, nil
}

func (c *Client) GetBlockHeaderBytesByHeight(p *types.BlockHeightParam) (*ethTypes.Header, error) {
	header, err := c.eth.HeaderByNumber(context.Background(), p.Height)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (c *Client) GetDataByHash(p *types.DataHashParam) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	var res []byte
	if err := c.rpc.CallContext(ctx, &res, "eth_getBlockByHash", p); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetProofForResult(p *types.ProofResultParam) ([][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	var res [][]byte
	if err := c.rpc.CallContext(ctx, &res, "eth_getProofForResult", p); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetProofForEvents(p *types.ProofEventsParam) ([][][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	var res [][][]byte
	if err := c.rpc.CallContext(ctx, &res, "eth_getProofForEvents", p); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) MonitorBlock(ctx context.Context, opts *BnOptions, callback func(v *types.BlockNotification) error) error {
	if opts == nil {
		return errors.New("receiveLoop: invalid options: <nil>")
	}

	// block a notification channel
	// (buffered: to avoid deadlock)
	// increase concurrency parameter for faster sync
	bnch := make(chan *types.BlockNotification, SyncConcurrency)

	heightTicker := time.NewTicker(BlockInterval)
	defer heightTicker.Stop()

	heightPoller := time.NewTicker(BlockHeightPollInterval)
	defer heightPoller.Stop()

	latestHeight := func() uint64 {
		height, err := c.GetBlockNumber()
		if err != nil {
			return 0
		}
		return height - BlockFinalityConfirmations
	}
	next, latest := opts.StartHeight, latestHeight()

	// last unverified block notification
	var lbn *types.BlockNotification
	// start monitor loop

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-heightTicker.C:
			latest++

		case <-heightPoller.C:
			if height := latestHeight(); height > 0 {
				latest = height
			}

		case bn := <-bnch:
			// process all notifications
			for ; bn != nil; next++ {
				if lbn != nil {
					if bn.Height.Cmp(lbn.Height) == 0 {
						if bn.Header.ParentHash != lbn.Header.ParentHash {
							break
						}
					} else {
						if vr != nil {
							if err := vr.Verify(lbn.Header, bn.Header, bn.Receipts); err != nil {
								next--
								break
							}
							if err := vr.Update(lbn.Header); err != nil {
								return errors.Wrapf(err, "receiveLoop: vr.Update: %v", err)
							}
						}
						if err := callback(lbn); err != nil {
							return errors.Wrapf(err, "receiveLoop: callback: %v", err)
						}
					}
				}
				if lbn, bn = bn, nil; len(bnch) > 0 {
					bn = <-bnch
				}
			}
			// remove unprocessed notifications
			for len(bnch) > 0 {
				<-bnch
			}

		default:
			if next >= latest {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			type bnq struct {
				h     uint64
				v     *types.BlockNotification
				err   error
				retry int
			}
			qch := make(chan *bnq, cap(bnch))
			for i := next; i < latest &&
				len(qch) < cap(qch); i++ {
				qch <- &bnq{i, nil, nil, RPCCallRetry} // fill bch with requests
			}
			if len(qch) == 0 {
				c.log.Error("Fatal: Zero length of query channel. Avoiding deadlock")
				continue
			}
			bns := make([]*types.BlockNotification, 0, len(qch))
			for q := range qch {
				switch {
				case q.err != nil:
					if q.retry > 0 {
						q.retry--
						q.v, q.err = nil, nil
						qch <- q
						continue
					}
					bns = append(bns, nil)
					if len(bns) == cap(bns) {
						close(qch)
					}

				case q.v != nil:
					bns = append(bns, q.v)
					if len(bns) == cap(bns) {
						close(qch)
					}
				default:
					go func(q *bnq) {
						defer func() {
							time.Sleep(500 * time.Millisecond)
							qch <- q
						}()

						if q.v == nil {
							q.v = &types.BlockNotification{}
						}

						q.v.Height = (&big.Int{}).SetUint64(q.h)

						if q.v.Header == nil {
							header, err := c.GetHeaderByHeight(ctx, q.v.Height)
							if err != nil {
								q.err = errors.Wrapf(err, "GetHeaderByHeight: %v", err)
								return
							}
							q.v.Header = header
							q.v.Hash = q.v.Header.Hash()
						}
						if q.v.Header.GasUsed > 0 {
							if q.v.HasBTPMessage == nil {
								hasBTPMessage, err := r.hasBTPMessage(ctx, q.v.Height)
								if err != nil {
									q.err = errors.Wrapf(err, "hasBTPMessage: %v", err)
									return
								}
								q.v.HasBTPMessage = &hasBTPMessage
							}
							if !*q.v.HasBTPMessage {
								return
							}
							// TODO optimize retry of GetBlockReceipts()
							q.v.Receipts, q.err = c.GetBlockReceipts(q.v.Hash)
							if q.err != nil {
								q.err = errors.Wrapf(q.err, "GetBlockReceipts: %v", q.err)
								return
							}
						}
					}(q)
				}
			}
			// filter nil
			_bns_, bns := bns, bns[:0]
			for _, v := range _bns_ {
				if v != nil {
					bns = append(bns, v)
				}
			}
			// sort and forward notifications
			if len(bns) > 0 {
				sort.SliceStable(bns, func(i, j int) bool {
					return bns[i].Height.Uint64() < bns[j].Height.Uint64()
				})
				for i, v := range bns {
					if v.Height.Uint64() == next+uint64(i) {
						bnch <- v
					}
				}
			}
		}
	}
}

func (c *Client) MonitorEvent(ctx context.Context, q eth.FilterQuery, cb func(*ethclient.Client, *types.EventNotification) error, errCb func(*ethclient.Client, error)) error {
	ch := make(chan ethTypes.Log)
	sub, err := c.eth.SubscribeFilterLogs(ctx, q, ch)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-sub.Err():
			errCb(c.eth, err)
			return err
		case v := <-ch:
			data := &types.EventNotification{
				Hash:   types.HexBytes(v.BlockHash.Bytes()),
				Height: v.BlockNumber,
				Index:  types.HexInt(v.Index),
			}
			if err := cb(c.eth, data); err != nil {
				return err
			}
		}
	}
}

func (c *Client) GetBlockHeaderByHeight(height int64) (*ethTypes.Header, error) {
	p := &types.BlockHeightParam{Height: big.NewInt(height)}
	header, err := c.GetBlockHeaderBytesByHeight(p)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (c *Client) GetBlockReceipts(hash common.Hash) (ethTypes.Receipts, error) {
	hb, err := c.GetBlockByHash(hash)
	if err != nil {
		return nil, err
	}
	if len(hb.NormalTransactions) == 0 {
		return nil, nil
	}
	txhs := hb.NormalTransactions
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

func (c *Client) GetChainID() string {
	return c.chainID
}

func (c *Client) GetEthClient() *ethclient.Client {
	return c.eth
}

func (c *Client) Log() *zap.Logger {
	return c.log
}
