package sui

import (
	"context"
	"sync"

	suisdkClient "github.com/coming-chat/go-sui/v2/client"
	"github.com/icon-project/centralized-relay/relayer/provider"

	"go.uber.org/zap"
)

type Config struct {
	ChainName string `yaml:"-" json:"-"`
	ChainID   string `yaml:"chain-id" json:"chain-id"`
	RPCUrl    string `yaml:"rpc-url" json:"rpc-url"`
	Address   string `yaml:"address" json:"address"`
	NID       string `yaml:"nid" json:"nid"`

	XcallPkgIDs    []string `yaml:"xcall-package-ids" json:"xcall-package-ids"`
	XcallStorageID string   `yaml:"xcall-storage-id" json:"xcall-storage-id"`

	ConnectionID    string `yaml:"connection-id" json:"connection-id"`
	ConnectionCapID string `yaml:"connection-cap-id" json:"connection-cap-id"`

	DappPkgID   string       `yaml:"dapp-package-id" json:"dapp-package-id"`
	DappModules []DappModule `yaml:"dapp-modules" json:"dapp-modules"`

	SpokeTokenTreasuryCapID string `yaml:"spoke-token-treasury-cap-id" json:"spoke-token-treasury-cap-id"`

	HomeDir  string `yaml:"home-dir" json:"home-dir"`
	GasLimit uint64 `yaml:"gas-limit" json:"gas-limit"`
	Disabled bool   `json:"disabled" yaml:"disabled"`
}

type DappModule struct {
	Name     string `yaml:"name" json:"name"`
	CapID    string `yaml:"cap-id" json:"cap-id"`
	ConfigID string `yaml:"config-id" json:"config-id"`
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
func (c *Config) Enabled() bool {
	return !c.Disabled
}
