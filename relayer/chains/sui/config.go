package sui

import (
	"context"

	suisdk "github.com/block-vision/sui-go-sdk/sui"
	suiclient "github.com/icon-project/centralized-relay/relayer/chains/sui/client"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"go.uber.org/zap"
)

type Config struct {
	ChainID   string `yaml:"chain-id"`
	ChainName string `yaml:"-"`
	RPCUrl    string `yaml:"rpc-url"`
	WsUrl     string `yaml:"ws-url"`
	Address   string `json:"address" yaml:"address"`
	NID       string `json:"nid" yaml:"nid"`
	PackageID string `yaml:"package-id"`
	HomeDir   string `yaml:"home-dir"`
}

func (pc *Config) NewProvider(ctx context.Context, logger *zap.Logger, homePath string, debug bool, chainName string) (provider.ChainProvider, error) {
	pc.HomeDir = homePath
	pc.ChainName = chainName

	if err := pc.Validate(); err != nil {
		return nil, err
	}

	rpcClient := suisdk.NewSuiClient(pc.RPCUrl)

	client := suiclient.NewClient(rpcClient, logger)

	return &Provider{
		log:    logger.With(zap.String("nid ", pc.ChainID), zap.String("name", pc.ChainName)),
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
