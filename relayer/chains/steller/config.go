package steller

import (
	"context"
	"net/http"

	"github.com/icon-project/centralized-relay/relayer/provider"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stellar/go/clients/horizonclient"
	"go.uber.org/zap"
)

type Config struct {
	ChainID    string                         `yaml:"chain-id"`
	ChainName  string                         `yaml:"-"`
	HorizonUrl string                         `yaml:"horizon-url"`
	Address    string                         `yaml:"address"`
	Contracts  relayertypes.ContractConfigMap `yaml:"contracts"`
	NID        string                         `json:"nid" yaml:"nid"`
	HomeDir    string                         `yaml:"home-dir"`
	GasPrice   uint64                         `yaml:"gas-price"`
	GasMin     uint64                         `yaml:"gas-min"`
	GasLimit   uint64                         `yaml:"gas-limit"`
}

func (pc *Config) NewProvider(ctx context.Context, logger *zap.Logger, homePath string, debug bool, chainName string) (provider.ChainProvider, error) {
	pc.HomeDir = homePath
	pc.ChainName = chainName

	if err := pc.Validate(); err != nil {
		return nil, err
	}

	httpClient := http.Client{}
	horizonClient := &horizonclient.Client{
		HorizonURL: pc.HorizonUrl,
		HTTP:       &httpClient,
		AppName:    "centralized-relay",
	}

	client := NewClient(horizonClient)

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
