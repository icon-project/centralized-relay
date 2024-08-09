package solana

import (
	"context"
	"fmt"
	"slices"
	"sync"

	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/icon-project/centralized-relay/relayer/chains/solana/types"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"go.uber.org/zap"
)

type Config struct {
	Disabled  bool   `yaml:"disabled" json:"disabled"`
	ChainName string `yaml:"-"`

	RPCUrl  string `yaml:"rpc-url"`
	Address string `yaml:"address"`

	XcallProgram string `yaml:"xcall-program"`

	ConnectionProgram string   `yaml:"connection-program"`
	OtherConnections  []string `yaml:"other-connections"`

	Dapps []types.Dapp `yaml:"dapps"`

	CpNIDs []string `yaml:"cp-nids"` //counter party NIDs Eg: ["0x2.icon", "0x7.icon"]

	AltAddress string `yaml:"alt-address"` // address lookup table address

	NID         string `yaml:"nid"`
	HomeDir     string `yaml:"home-dir"`
	GasLimit    uint64 `yaml:"gas-limit"`
	StartHeight uint64 `yaml:"start-height"`
}

func (pc *Config) NewProvider(ctx context.Context, logger *zap.Logger, homePath string, debug bool, chainName string) (provider.ChainProvider, error) {
	pc.HomeDir = homePath
	pc.ChainName = chainName

	if err := pc.Validate(); err != nil {
		return nil, err
	}

	client := NewClient(solrpc.New(pc.RPCUrl))

	xcallIdl := IDL{}
	if pc.XcallProgram != "" {
		if err := client.FetchIDL(ctx, pc.XcallProgram, &xcallIdl); err != nil {
			return nil, err
		}
	}

	connIdl := IDL{}
	if pc.ConnectionProgram != "" {
		if err := client.FetchIDL(ctx, pc.ConnectionProgram, &connIdl); err != nil {
			return nil, err
		}
	}

	pdaRegistry := types.NewPDARegistry(xcallIdl.GetProgramID(), connIdl.GetProgramID())

	return &Provider{
		log:         logger.With(zap.String("nid ", pc.NID), zap.String("name", pc.ChainName)),
		cfg:         pc,
		client:      client,
		txmut:       &sync.Mutex{},
		xcallIdl:    &xcallIdl,
		connIdl:     &connIdl,
		pdaRegistry: pdaRegistry,
		staticAlts:  make(types.AddressTables),
	}, nil
}

func (pc *Config) SetWallet(addr string) {
	pc.Address = addr
}

func (pc *Config) GetWallet() string {
	return pc.Address
}

func (pc *Config) Validate() error {
	for _, dapp := range pc.Dapps {
		if !slices.Contains(types.DappsEnabled, dapp.Name) {
			return fmt.Errorf("invalid dapp name %s; should be one of %+v", dapp.Name, types.DappsEnabled)
		}
	}
	return nil
}

func (pc *Config) Enabled() bool {
	return !pc.Disabled
}
