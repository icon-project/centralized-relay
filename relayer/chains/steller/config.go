package steller

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/icon-project/centralized-relay/relayer/chains/steller/sorobanclient"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/centralized-relay/relayer/types"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stellar/go/clients/horizonclient"
	"go.uber.org/zap"
)

type Config struct {
	provider.CommonConfig
	ChainID           string                         `json:"chain-id" yaml:"chain-id"`
	ChainName         string                         `json:"-t" yaml:"-"`
	HorizonUrl        string                         `json:"horizon-url" yaml:"horizon-url"`
	SorobanUrl        string                         `json:"soroban-url" yaml:"soroban-url"`
	Address           string                         `json:"address" yaml:"address"`
	Contracts         relayertypes.ContractConfigMap `json:"contracts" yaml:"contracts"`
	NID               string                         `json:"nid" yaml:"nid"`
	HomeDir           string                         `json:"home-dir" yaml:"home-dir"`
	MaxInclusionFee   uint64                         `json:"max-inclusion-fee" yaml:"max-inclusion-fee"` // in stroop: the smallest unit of a lumen, one ten-millionth of a lumen (.0000001 XLM).
	NetworkPassphrase string                         `json:"network-passphrase" yaml:"network-passphrase"`
	StartHeight       uint64                         `json:"start-height" yaml:"start-height"` // would be of highest priority
	Disabled          bool                           `json:"disabled" yaml:"disabled"`
	PollInterval      time.Duration                  `json:"poll-interval" yaml:"poll-interval"`
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

// Enabled returns true if the chain is enabled
func (pc *Config) Enabled() bool {
	return !pc.Disabled
}

func (pc *Config) ContractsAddress() types.ContractConfigMap {
	return pc.Contracts
}
