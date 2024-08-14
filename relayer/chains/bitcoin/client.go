package bitcoin

import (
	"bytes"
	"context"
	"os"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"go.uber.org/zap"

	// "github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func RunApp() {
	goEnv := os.Getenv("GO_ENV")
	if goEnv == "master" {
		startMaster()
	} else {
		startSlave()
	}
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
}

// grouped rpc api clients
type Client struct {
	log        *zap.Logger
	client     *rpcclient.Client
	chainParam *chaincfg.Params
}

// create new client
func newClient(ctx context.Context, rpcUrl string, httpPostMode, disableTLS bool, l *zap.Logger) (IClient, error) {
	// Connect to the Bitcoin Core RPC server
	connConfig := &rpcclient.ConnConfig{
		Host:         rpcUrl,
		User:         "123",
		Pass:         "123",
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
	meetRequirement1 := 0
	meetRequirement2 := 0

	for i := param.StartHeight; i <= param.EndHeight; i++ {
		blockHash, err := c.client.GetBlockHash(int64(i))
		if err != nil {
			return nil, err
		}

		block, err := c.client.GetBlock(blockHash)
		// loop thru transactions
		for _, tx := range block.Transactions {
			// loop thru tx output
			for _, txOutput := range tx.TxOut {
				if len(txOutput.PkScript) > 2 {
					// check OP_RETURN
					if txOutput.PkScript[0] == txscript.OP_RETURN && txOutput.PkScript[1] == byte(param.OPReturnPrefix) {
						meetRequirement1++
					}

					// check EQUAL to multisig script
					if bytes.Equal(param.BitcoinScript, txOutput.PkScript) {
						meetRequirement2++
					}

					if meetRequirement2*meetRequirement1 != 0 {
						res = append(res, &TxSearchRes{Height: i, Tx: tx})
						break
					}
				}
			}
			meetRequirement2 = 0
			meetRequirement1 = 0
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
