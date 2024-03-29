package sui

import (
	"context"

	suisdkClient "github.com/coming-chat/go-sui/v2/client"
	providerTypes "github.com/icon-project/centralized-relay/relayer/chains/sui/types"
	"github.com/icon-project/centralized-relay/relayer/provider"

	"go.uber.org/zap"
)

type Config struct {
	ChainID   string                          `yaml:"chain-id"`
	ChainName string                          `yaml:"-"`
	RPCUrl    string                          `yaml:"rpc-url"`
	WsUrl     string                          `yaml:"ws-url"`
	Address   string                          `yaml:"address"`
	Contracts providerTypes.ContractConfigMap `yaml:"contracts"`
	NID       string                          `json:"nid" yaml:"nid"`
	PackageID string                          `yaml:"package-id"`
	HomeDir   string                          `yaml:"home-dir"`
	GasPrice  uint64                          `yaml:"gas-price"`
	GasMin    uint64                          `yaml:"gas-min"`
	GasLimit  uint64                          `yaml:"gas-limit"`
}

func (pc *Config) NewProvider(ctx context.Context, logger *zap.Logger, homePath string, debug bool, chainName string) (provider.ChainProvider, error) {
	pc.HomeDir = homePath
	pc.ChainName = chainName

	if err := pc.Validate(); err != nil {
		return nil, err
	}
	rpcClient, err := suisdkClient.Dial(pc.RPCUrl)
	if err != nil {
		return nil, err
	}
	client := NewClient(rpcClient, logger)

	return &Provider{
		log:    logger.With(zap.String("nid ", pc.NID), zap.String("name", pc.ChainName)),
		cfg:    pc,
		client: client,
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
