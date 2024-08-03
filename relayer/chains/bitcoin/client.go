package bitcoin

import (
	"os"
	"context"
	"github.com/btcsuite/btcd/rpcclient"
	"go.uber.org/zap"
	// "github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/btcjson"
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
}

// grouped rpc api clients
type Client struct {
	log *zap.Logger
	client *rpcclient.Client
}

// create new client 
func newClient(rpcUrl, user, pass string, httpPostMode, disableTLS bool, l *zap.Logger) (IClient, error) {
	// Connect to the Bitcoin Core RPC server
	connConfig := &rpcclient.ConnConfig{
		Host:         rpcUrl,
		User:         user,
		Pass:         pass,
		HTTPPostMode: httpPostMode,
		DisableTLS:   disableTLS,
	}

	client, err := rpcclient.New(connConfig, nil)
	if err != nil {
		return nil, err
	}

	// ws 

	return &Client {
		log: l,
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

//
func (c *Client) GetBalance(ctx context.Context, addr string) (uint64, error) {
	return 0, nil
}

// 
func (c *Client) Subscribe(ctx context.Context, _, query string) error {
	// 

	return nil
}

//
func (c *Client) Unsubscribe(ctx context.Context, _, query string) error {
	return nil	
}

// test data
func (c *Client) GetFee(ctx context.Context) (uint64, error) {
	return 10, nil
}