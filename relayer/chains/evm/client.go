package evm

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/websocket"
	"github.com/icon-project/centralized-relay/relayer/chains/evm/types"
	"github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/goloop/module"
	"github.com/icon-project/goloop/server/jsonrpc"
	"github.com/icon-project/goloop/service/transaction"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	txMaxDataSize        = 8 * 1024 // 8 KB
	txOverheadScale      = 0.01     // base64 encoding overhead 0.36, rlp and other fields 0.01
	defaultTxSizeLimit   = txMaxDataSize / (1 + txOverheadScale)
	defaultSendTxTimeout = 15 * time.Second
	defaultGasPrice      = 18000000000
	maxGasPriceBoost     = 10.0
	defaultReadTimeout   = 50 * time.Second //
	DefaultGasLimit      = 25000000
	RPCCallRetry         = 5

	DefaultSendTransactionRetryInterval        = 3 * time.Second        // 3sec
	DefaultGetTransactionResultPollingInterval = 500 * time.Millisecond // 1.5sec
	JsonrpcApiVersion                          = 3
)

var txSerializeExcludes = map[string]bool{"signature": true}

type Wallet interface {
	Sign(data []byte) ([]byte, error)
	Address() string
}

func newClient(url string, l *zap.Logger) (IClient, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	eth := ethclient.NewClient(rpcClient)
	if err != nil {
		return nil, err
	}
	cl := &Client{
		log: l,
		rpc: rpcClient,
		eth: eth,
	}
	return cl, nil
}

// grouped rpc api clients
type Client struct {
	log     *zap.Logger
	rpc     *rpc.Client
	eth     *ethclient.Client
	chainID string
}

type IClient interface {
	Log() *zap.Logger
	GetBalance(ctx context.Context, hexAddr string) (*big.Int, error)
	GetBlockNumber() (uint64, error)
	GetBlockByHash(hash common.Hash) (*types.Block, error)
	GetHeaderByHeight(ctx context.Context, height *big.Int) (*ethTypes.Header, error)
	GetBlockReceipts(hash common.Hash) (ethTypes.Receipts, error)
	GetChainID() string

	// ethClient
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	TransactionByHash(ctx context.Context, blockHash common.Hash) (tx *ethTypes.Transaction, isPending bool, err error)
	Call(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*ethTypes.Receipt, error)
	TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
	TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*ethTypes.Transaction, error)
}

func (c *Client) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return c.eth.NonceAt(ctx, account, blockNumber)
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

func (c *Client) Call(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return c.eth.CallContract(ctx, msg, blockNumber)
}

func (c *Client) GetBalance(ctx context.Context, hexAddr string) (*big.Int, error) {
	if !common.IsHexAddress(hexAddr) {
		return nil, fmt.Errorf("invalid hex address: %v", hexAddr)
	}
	return c.eth.BalanceAt(ctx, common.HexToAddress(hexAddr), nil)
}

func (c *Client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
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

func (c *Client) GetBlockByHash(hash common.Hash) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()
	block, err := c.eth.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return &types.Block{
		Height:    block.Number().Uint64(),
		Timestamp: block.Time(),
		Header:    block.Header,
	}, nil
}

func (c *Client) GetHeaderByHeight(ctx context.Context, height *big.Int) (*ethTypes.Header, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()
	return c.eth.HeaderByNumber(ctx, height)
}

func (c *Client) SignTransaction(w module.Wallet, p *types.TransactionParam) error {
	p.Timestamp = types.NewHexInt(time.Now().UnixNano() / int64(time.Microsecond))
	js, err := json.Marshal(p)
	if err != nil {
		return err
	}

	bs, err := transaction.SerializeJSON(js, nil, txSerializeExcludes)
	if err != nil {
		return err
	}
	bs = append([]byte("icx_sendTransaction."), bs...)
	txHash := crypto.SHA3Sum256(bs)
	p.TxHash = types.NewHexBytes(txHash)
	sig, err := w.Sign(txHash)
	if err != nil {
		return err
	}
	p.Signature = base64.StdEncoding.EncodeToString(sig)
	return nil
}

func (c *Client) SendTransaction(p *types.TransactionParam) (*types.HexBytes, error) {
	var result types.HexBytes
	if _, err := c.Do("icx_sendTransaction", p, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SendTransactionAndWait(p *types.TransactionParam) (*types.HexBytes, error) {
	var result types.HexBytes
	if _, err := c.Do("icx_sendTransactionAndWait", p, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetTransactionResult(p *types.TransactionHashParam) (*types.TransactionResult, error) {
	tr := &types.TransactionResult{}
	if _, err := c.Do("icx_getTransactionResult", p, tr); err != nil {
		return nil, err
	}
	return tr, nil
}

func (c *Client) WaitTransactionResult(p *types.TransactionHashParam) (*types.TransactionResult, error) {
	tr := &types.TransactionResult{}
	if _, err := c.Do("icx_waitTransactionResult", p, tr); err != nil {
		return nil, err
	}
	return tr, nil
}

func (c *Client) WaitForResults(ctx context.Context, thp *types.TransactionHashParam) (txh *types.HexBytes, txr *types.TransactionResult, err error) {
	ticker := time.NewTicker(time.Duration(DefaultGetTransactionResultPollingInterval) * time.Nanosecond)
	retryLimit := 20
	retryCounter := 0
	txh = &thp.Hash
	for {
		defer ticker.Stop()
		select {
		case <-ctx.Done():
			err = errors.New("Context Cancelled ReceiptWait Exiting ")
			return
		case <-ticker.C:
			if retryCounter >= retryLimit {
				err = errors.New("Retry Limit Exceeded while waiting for results of transaction")
				return
			}
			retryCounter++
			txr, err = c.GetTransactionResult(thp)
			if err != nil {
				switch re := err.(type) {
				case *jsonrpc.Error:
					switch re.Code {
					case jsonrpc.ErrorCodePending, jsonrpc.ErrorCodeNotFound, jsonrpc.ErrorCodeExecuting:
						continue
					}
				}
			}
			return
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

func (c *Client) GetVotesByHeight(p *types.BlockHeightParam) ([]byte, error) {
	var result []byte
	if _, err := c.Do("icx_getVotesByHeight", p, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetDataByHash(p *types.DataHashParam) ([]byte, error) {
	var result []byte
	_, err := c.Do("icx_getDataByHash", p, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetProofForResult(p *types.ProofResultParam) ([][]byte, error) {
	var result [][]byte
	if _, err := c.Do("icx_getProofForResult", p, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetProofForEvents(p *types.ProofEventsParam) ([][][]byte, error) {
	var result [][][]byte
	if _, err := c.Do("icx_getProofForEvents", p, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) MonitorBlock(ctx context.Context, p *types.BlockRequest, cb func(conn *websocket.Conn, v *types.BlockNotification) error, scb func(conn *websocket.Conn), errCb func(*websocket.Conn, error)) error {
	resp := &types.BlockNotification{}
	return c.Monitor(ctx, "/block", p, resp, func(conn *websocket.Conn, v interface{}) error {
		switch t := v.(type) {
		case *types.BlockNotification:
			if err := cb(conn, t); err != nil {
				// c.log.Debugf("MonitorBlock callback return err:%+v", err)
				return err
			}
		case types.WSEvent:
			switch t {
			case types.WSEventInit:
				if scb != nil {
					scb(conn)
				} else {
					return errors.New("Second Callback function (scb) is nil ")
				}
			}
		case error:
			errCb(conn, t)
			return t
		default:
			errCb(conn, fmt.Errorf("not supported type %T", t))
			return errors.New("Not supported type")
		}
		return nil
	})
}

func (c *Client) MonitorEvent(ctx context.Context, p *types.EventRequest, cb func(conn *websocket.Conn, v *types.EventNotification) error, errCb func(*websocket.Conn, error)) error {
	resp := &types.EventNotification{}
	return c.Monitor(ctx, "/event", p, resp, func(conn *websocket.Conn, v interface{}) error {
		switch t := v.(type) {
		case *types.EventNotification:
			if err := cb(conn, t); err != nil {
				c.log.Debug(fmt.Sprintf("MonitorEvent callback return err:%+v", err))
			}
		case error:
			errCb(conn, t)
		default:
			errCb(conn, fmt.Errorf("not supported type %T", t))
		}
		return nil
	})
}

func (c *Client) Monitor(ctx context.Context, reqUrl string, reqPtr, respPtr interface{}, cb types.WsReadCallback) error {
	if cb == nil {
		return fmt.Errorf("callback function cannot be nil")
	}
	conn, err := c.wsConnect(reqUrl, nil)
	if err != nil {
		return err
	}
	defer func() {
		c.log.Debug(fmt.Sprintf("Monitor finish %s", conn.LocalAddr().String()))
		c.wsClose(conn)
	}()
	if err = c.wsRequest(conn, reqPtr); err != nil {
		return err
	}
	if err := cb(conn, types.WSEventInit); err != nil {
		return err
	}
	return c.wsReadJSONLoop(ctx, conn, respPtr, cb)
}

func (c *Client) CloseMonitor(conn *websocket.Conn) {
	c.log.Debug(fmt.Sprintf("CloseMonitor %s", conn.LocalAddr().String()))
	c.wsClose(conn)
}

func (c *Client) CloseAllMonitor() {
	for _, conn := range c.conns {
		c.log.Debug(fmt.Sprintf("CloseAllMonitor %s", conn.LocalAddr().String()))
		c.wsClose(conn)
	}
}

func (c *Client) _addWsConn(conn *websocket.Conn) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	la := conn.LocalAddr().String()
	c.conns[la] = conn
}

func (c *Client) _removeWsConn(conn *websocket.Conn) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	la := conn.LocalAddr().String()
	_, ok := c.conns[la]
	if ok {
		delete(c.conns, la)
	}
}

type wsConnectError struct {
	error
	httpResp *http.Response
}

func (c *Client) wsConnect(reqUrl string, reqHeader http.Header) (*websocket.Conn, error) {
	wsEndpoint := strings.Replace(c.Endpoint, "http", "ws", 1)
	conn, httpResp, err := websocket.DefaultDialer.Dial(wsEndpoint+reqUrl, reqHeader)
	if err != nil {
		wsErr := wsConnectError{error: err}
		wsErr.httpResp = httpResp
		return nil, wsErr
	}
	c._addWsConn(conn)
	return conn, nil
}

type wsRequestError struct {
	error
	wsResp *types.WSResponse
}

func (c *Client) wsRequest(conn *websocket.Conn, reqPtr interface{}) error {
	if reqPtr == nil {
		log.Panicf("reqPtr cannot be nil")
	}
	var err error
	wsResp := &types.WSResponse{}
	if err = conn.WriteJSON(reqPtr); err != nil {
		return wsRequestError{fmt.Errorf("fail to WriteJSON err:%+v", err), nil}
	}

	if err = conn.ReadJSON(wsResp); err != nil {
		return wsRequestError{fmt.Errorf("fail to ReadJSON err:%+v", err), nil}
	}

	if wsResp.Code != 0 {
		return wsRequestError{
			fmt.Errorf("invalid WSResponse code:%d, message:%s", wsResp.Code, wsResp.Message),
			wsResp,
		}
	}
	return nil
}

func (c *Client) wsClose(conn *websocket.Conn) {
	c._removeWsConn(conn)
	if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		c.log.Debug(fmt.Sprintf("fail to WriteMessage CloseNormalClosure err:%+v", err))
	}
	if err := conn.Close(); err != nil {
		c.log.Debug(fmt.Sprintf("fail to Close err:%+v", err))
	}
}

func (c *Client) wsRead(conn *websocket.Conn, respPtr interface{}) error {
	mt, r, err := conn.NextReader()
	if err != nil {
		return err
	}
	if mt == websocket.CloseMessage {
		return io.EOF
	}
	return json.NewDecoder(r).Decode(respPtr)
}

func (c *Client) wsReadJSONLoop(ctx context.Context, conn *websocket.Conn, respPtr interface{}, cb types.WsReadCallback) error {
	elem := reflect.ValueOf(respPtr).Elem()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			v := reflect.New(elem.Type())
			ptr := v.Interface()
			if _, ok := c.conns[conn.LocalAddr().String()]; !ok {
				c.log.Debug(fmt.Sprintf("wsReadJSONLoop c.conns[%s] is nil", conn.LocalAddr().String()))
				return errors.New("wsReadJSONLoop c.conns is nil")
			}
			if err := c.wsRead(conn, ptr); err != nil {
				c.log.Debug(fmt.Sprintf("wsReadJSONLoop c.conns[%s] ReadJSON err:%+v", conn.LocalAddr().String(), err))
				if cErr, ok := err.(*websocket.CloseError); !ok || cErr.Code != websocket.CloseNormalClosure {
					cb(conn, err)
				}
				return err
			}
			if err := cb(conn, ptr); err != nil {
				return err
			}
		}
	}
}

func (c *Client) GetBlockHeaderByHeight(height int64) (*types.BlockHeader, error) {
	p := &types.BlockHeightParam{Height: eth}
	b, err := c.GetBlockHeaderBytesByHeight(p)
	if err != nil {
		return nil, err
	}
	var blockHeader types.BlockHeader
	_, err = codec.RLP.UnmarshalFromBytes(b, &blockHeader)
	if err != nil {
		return nil, err
	}
	return &blockHeader, nil
}

func (c *Client) GetValidatorsByHash(hash common.HexHash) ([]common.Address, error) {
	data, err := c.GetDataByHash(&types.DataHashParam{Hash: types.NewHexBytes(hash.Bytes())})
	if err != nil {
		return nil, errors.Wrapf(err, "GetDataByHash; %v", err)
	}
	if !bytes.Equal(hash, crypto.SHA3Sum256(data)) {
		return nil, errors.Errorf(
			"invalid data: hash=%v, data=%v", hash, common.HexBytes(data))
	}
	var validators []common.Address
	_, err = codec.BC.UnmarshalFromBytes(data, &validators)
	if err != nil {
		return nil, errors.Wrapf(err, "Unmarshal Validators: %v", err)
	}
	return validators, nil
}

func (c *Client) GetBalance(param *types.AddressParam) (*big.Int, error) {
	var result types.HexInt
	_, err := c.Do("icx_getBalance", param, &result)
	if err != nil {
		return nil, err
	}
	bInt, err := result.BigInt()
	if err != nil {
		return nil, err
	}
	return bInt, nil
}

const (
	HeaderKeyIconOptions = "Icon-Options"
	IconOptionsDebug     = "debug"
	IconOptionsTimeout   = "timeout"
)

type IconOptions map[string]string

func (opts IconOptions) Set(key, value string) {
	opts[key] = value
}

func (opts IconOptions) Get(key string) string {
	if opts == nil {
		return ""
	}
	v := opts[key]
	if len(v) == 0 {
		return ""
	}
	return v
}

func (opts IconOptions) Del(key string) {
	delete(opts, key)
}

func (opts IconOptions) SetBool(key string, value bool) {
	opts.Set(key, strconv.FormatBool(value))
}

func (opts IconOptions) GetBool(key string) (bool, error) {
	return strconv.ParseBool(opts.Get(key))
}

func (opts IconOptions) SetInt(key string, v int64) {
	opts.Set(key, strconv.FormatInt(v, 10))
}

func (opts IconOptions) GetInt(key string) (int64, error) {
	return strconv.ParseInt(opts.Get(key), 10, 64)
}

func (opts IconOptions) ToHeaderValue() string {
	if opts == nil {
		return ""
	}
	strs := make([]string, len(opts))
	i := 0
	for k, v := range opts {
		strs[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}
	return strings.Join(strs, ",")
}

func NewIconOptionsByHeader(h http.Header) IconOptions {
	s := h.Get(HeaderKeyIconOptions)
	if s != "" {
		kvs := strings.Split(s, ",")
		m := make(map[string]string)
		for _, kv := range kvs {
			if kv != "" {
				idx := strings.Index(kv, "=")
				if idx > 0 {
					m[kv[:idx]] = kv[(idx + 1):]
				} else {
					m[kv] = ""
				}
			}
		}
		return m
	}
	return nil
}

func (c *Client) EstimateStep(param *types.TransactionParamForEstimate) (*types.HexInt, error) {
	if len(c.DebugEndPoint) == 0 {
		return nil, errors.New("UnavailableDebugEndPoint")
	}
	currTime := time.Now().UnixNano() / time.Hour.Microseconds()
	param.Timestamp = types.NewHexInt(currTime)
	var result types.HexInt
	if _, err := c.DoURL(c.DebugEndPoint,
		"debug_estimateStep", param, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func NewClient(uri string, l *zap.Logger) (*Client, error) {
	rpcClient, err := rpc.Dial(uri)
	if err != nil {
		return nil, err
	}
	eth := ethclient.NewClient(rpcClient)
	if err != nil {
		return nil, err
	}
	cl := &Client{
		log: l,
		rpc: rpcClient,
		eth: eth,
	}
	return cl, nil
}

func guessDebugEndpoint(endpoint string) string {
	uo, err := url.Parse(endpoint)
	if err != nil {
		return ""
	}
	ps := strings.Split(uo.Path, "/")
	for i, v := range ps {
		if v == "api" {
			if len(ps) > i+1 && ps[i+1] == "v3" {
				ps[i+1] = "v3d"
				uo.Path = strings.Join(ps, "/")
				return uo.String()
			}
			break
		}
	}
	return ""
}

func (c *Client) GetBlockReceipts(hash common.Hash) (ethTypes.Receipts, error) {
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

func (c *Client) GetChainID() string {
	return c.chainID
}

func (c *Client) GetEthClient() *ethclient.Client {
	return c.eth
}

func (c *Client) Log() *zap.Logger {
	return c.log
}
