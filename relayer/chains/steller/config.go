package steller

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/icon-project/centralized-relay/relayer/chains/steller/sorobanclient"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stellar/go/clients/horizonclient"
	"go.uber.org/zap"
)

type Config struct {
	ChainID           string                         `yaml:"chain-id"`
	ChainName         string                         `yaml:"-"`
	HorizonUrl        string                         `yaml:"horizon-url"`
	SorobanUrl        string                         `yaml:"soroban-url"`
	Address           string                         `yaml:"address"`
	Contracts         relayertypes.ContractConfigMap `yaml:"contracts"`
	NID               string                         `json:"nid" yaml:"nid"`
	HomeDir           string                         `yaml:"home-dir"`
	MaxInclusionFee   uint64                         `yaml:"max-inclusion-fee"` // in stroop: the smallest unit of a lumen, one ten-millionth of a lumen (.0000001 XLM).
	BlockInterval     time.Duration                  `yaml:"block-interval"`
	NetworkPassphrase string                         `yaml:"network-passphrase"`
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

	sorobanclient, err := sorobanclient.New(pc.SorobanUrl, &httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create soroban client: %v", err)
	}

	client := NewClient(horizonClient, sorobanclient)

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
