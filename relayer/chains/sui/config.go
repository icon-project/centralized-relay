package sui

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	suisdkClient "github.com/coming-chat/go-sui/v2/client"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/centralized-relay/relayer/types"

	"go.uber.org/zap"
)

type Config struct {
	ChainName string `yaml:"-" json:"-"`
	ChainID   string `yaml:"chain-id" json:"chain-id"`
	RPCUrl    string `yaml:"rpc-url" json:"rpc-url"`
	Address   string `yaml:"address" json:"address"`
	NID       string `yaml:"nid" json:"nid"`

	XcallPkgID     string `yaml:"xcall-package-id" json:"xcall-package-id"`
	XcallStorageID string `yaml:"xcall-storage-id" json:"xcall-storage-id"`

	ConnectionModule string `yaml:"connection-module" json:"connection-module"`
	ConnectionID     string `yaml:"connection-id" json:"connection-id"`
	ConnectionCapID  string `yaml:"connection-cap-id" json:"connection-cap-id"`

	Dapps []Dapp `yaml:"dapps" json:"dapps"`

	HomeDir  string `yaml:"home-dir" json:"home-dir"`
	GasLimit uint64 `yaml:"gas-limit" json:"gas-limit"`
	Disabled bool   `json:"disabled" yaml:"disabled"`

	// Start tx-digest cursor to begin querying for events.
	// Should be empty if we want to query using last saved tx-digest
	// from database.
	StartTxDigest string `json:"start-tx-digest" yaml:"start-tx-digest"`

	PollInterval time.Duration `json:"poll-interval" yaml:"poll-interval"`
}

type DappModule struct {
	Name     string `yaml:"name" json:"name"`
	CapID    string `yaml:"cap-id" json:"cap-id"`
	ConfigID string `yaml:"config-id" json:"config-id"`
}

type Dapp struct {
	PkgID string `json:"package-id" yaml:"package-id"`

	// DappConstant is a map of name of sui constant to object id.
	Constants map[string]string `json:"constants" yaml:"constants"`

	Modules []DappModule `json:"modules" yaml:"modules"`
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

func (c *Config) ContractsAddress() types.ContractConfigMap {
	dapps, _ := json.Marshal(c.Dapps)

	return types.ContractConfigMap{
		"xcall-package-id":  c.XcallPkgID,
		"xcall-storage-id":  c.XcallStorageID,
		"connection-id":     c.ConnectionID,
		"connection-cap-id": c.ConnectionCapID,
		"dapps":             string(dapps),
	}
}
