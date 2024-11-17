package bitcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/icon-project/centralized-relay/utils/multisig"
	"go.uber.org/zap"

	// "github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type MempoolFeeResponse struct {
	FastestFee  uint64 `json:"fastestFee"`
	HalfHourFee uint64 `json:"halfHourFee"`
	HourFee     uint64 `json:"hourFee"`
	EconomyFee  uint64 `json:"economyFee"`
	MinimumFee  uint64 `json:"minimumFee"`
}

type IClient interface {
	// IsConnected() bool
	// Reconnect() error
	GetLatestBlockHeight(ctx context.Context) (uint64, error)
	GetTransactionReceipt(ctx context.Context, tx string) (*btcjson.TxRawResult, error)
	GetBalance(ctx context.Context, addr string) (uint64, error)

	Subscribe(ctx context.Context, _, query string) error
	Unsubscribe(ctx context.Context, _, query string) error
	GetFee(ctx context.Context) (uint64, error)
	DecodeAddress(btcAddr string) ([]byte, error)
	TxSearch(ctx context.Context, param TxSearchParam) ([]*TxSearchRes, error)
	SendRawTransaction(url string, rawMsg []json.RawMessage) (string, error)
}

// grouped rpc api clients
type Client struct {
	log        *zap.Logger
	client     *rpcclient.Client
	chainParam *chaincfg.Params
}

// create new client
func newClient(ctx context.Context, rpcUrl string, user string, password string, httpPostMode, disableTLS bool, l *zap.Logger) (IClient, error) {
	// Connect to the Bitcoin Core RPC server
	connConfig := &rpcclient.ConnConfig{
		Host:         rpcUrl,
		User:         user,
		Pass:         password,
		HTTPPostMode: httpPostMode,
		DisableTLS:   disableTLS,
	}

	client, err := rpcclient.New(connConfig, nil)
	if err != nil {
		return nil, err
	}

	// ws

	return &Client{
		log:    l,
		client: client,
	}, nil
}

// query block height
func (c *Client) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	// Get the current block height
	blockCount, err := c.client.GetBlockCount()
	if err != nil {
		return 0, err
	}

	return uint64(blockCount), nil
}

// get transaction reciept
func (c *Client) GetTransactionReceipt(ctx context.Context, tx string) (*btcjson.TxRawResult, error) {
	// convert to chain hash type
	txHash, err := chainhash.NewHashFromStr(tx)
	if err != nil {
		return nil, err
	}

	// query transaction
	txVerbose, err := c.client.GetRawTransactionVerbose(txHash)
	if err != nil {
		return nil, err
	}

	return txVerbose, nil
}

func (c *Client) GetBalance(ctx context.Context, addr string) (uint64, error) {
	return 0, nil
}

func (c *Client) Subscribe(ctx context.Context, _, query string) error {
	return nil
}

func (c *Client) Unsubscribe(ctx context.Context, _, query string) error {
	return nil
}

// test data
func (c *Client) GetFee(ctx context.Context) (uint64, error) {
	return 10, nil
}

func (c *Client) TxSearch(ctx context.Context, param TxSearchParam) ([]*TxSearchRes, error) {
	//
	res := []*TxSearchRes{}

	for i := param.StartHeight; i <= param.EndHeight; i++ {

		blockHash, err := c.client.GetBlockHash(int64(i))
		if err != nil {
			c.log.Error("Failed to get block hash", zap.Error(err))
			return nil, err
		}

		block, err := c.client.GetBlock(blockHash)
		if err != nil {
			c.log.Error("Failed to get block", zap.Error(err))
			return nil, err
		}
		// loop thru transactions
		for j, tx := range block.Transactions {
			bridgeMessage, err := multisig.ReadBridgeMessage(tx)
			if err != nil {
				continue
			}
			res = append(res, &TxSearchRes{Height: i, Tx: tx, TxIndex: uint64(j), BridgeMessage: bridgeMessage})
		}
	}

	return res, nil
}

func (c *Client) DecodeAddress(btcAddr string) ([]byte, error) {
	// return bitcoin script value
	decodedAddr, err := btcutil.DecodeAddress(btcAddr, c.chainParam)
	if err != nil {
		return nil, err
	}
	destinationAddrByte, err := txscript.PayToAddrScript(decodedAddr)
	if err != nil {
		return nil, err
	}

	return destinationAddrByte, nil
}

func (c *Client) SendRawTransaction(url string, rawMsg []json.RawMessage) (string, error) {
	if len(rawMsg) == 0 {
		return "", fmt.Errorf("empty raw message")
	}

	resp, err := http.Post(url, "text/plain", bytes.NewReader(rawMsg[0]))
	if err != nil {
		c.log.Error("failed to send transaction", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.log.Error("failed to broadcast transaction", zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return "", fmt.Errorf("broadcast failed: %v", err)
	}

	return string(body), nil
}
