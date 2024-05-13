package solana

import (
	"context"
	"sync"

	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"go.uber.org/zap"
)

type Config struct {
	ChainName string `yaml:"-"`

	RPCUrl            string `yaml:"rpc-url"`
	Address           string `yaml:"address"`
	XcallAddress      string `yaml:"xcall-address"`
	ConnectionAddress string `yaml:"connection-address"`
	NID               string `yaml:"nid"`
	HomeDir           string `yaml:"home-dir"`
	GasLimit          uint64 `yaml:"gas-limit"`
	StartHeight       uint64 `yaml:"start-height"`
}

func (pc *Config) NewProvider(ctx context.Context, logger *zap.Logger, homePath string, debug bool, chainName string) (provider.ChainProvider, error) {
	pc.HomeDir = homePath
	pc.ChainName = chainName

	if err := pc.Validate(); err != nil {
		return nil, err
	}

	client := NewClient(solrpc.New(pc.RPCUrl))

	return &Provider{
		log:    logger.With(zap.String("nid ", pc.NID), zap.String("name", pc.ChainName)),
		cfg:    pc,
		client: client,
		txmut:  &sync.Mutex{},
	}, nil
}

func (pc *Config) SetWallet(addr string) {
	pc.Address = addr
}

func (pc *Config) GetWallet() string {
	return pc.Address
}

func (pc *Config) Validate() error {
	//Todo
	return nil
}
