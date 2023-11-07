package evm

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/icon-project/centralized-relay/relayer/chains/evm/abi"
	"go.uber.org/zap"
)

func NewClient(url string, l *zap.Logger) (cls *Client, err error) {
	clrpc, err := rpc.Dial(url)
	if err != nil {
		l.Error("failed to create evm rpc client: url=%v, %v",
			zap.String("url", url),
			zap.Error(err))
		return nil, err
	}
	cleth := ethclient.NewClient(clrpc)
	return &Client{
		log: l,
		rpc: clrpc,
		eth: cleth,
	}, nil
}

// grouped rpc api clients
type Client struct {
	log *zap.Logger
	rpc *rpc.Client
	eth *ethclient.Client
	abi *abi.Storage
}
