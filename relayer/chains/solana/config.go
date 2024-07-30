package solana

import (
	"context"
	"fmt"
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

	CpNIDs []string `yaml:"cp-nids"`

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

	p := &Provider{
		log:         logger.With(zap.String("nid ", pc.NID), zap.String("name", pc.ChainName)),
		cfg:         pc,
		client:      client,
		txmut:       &sync.Mutex{},
		xcallIdl:    &xcallIdl,
		connIdl:     &connIdl,
		pdaRegistry: pdaRegistry,
		staticAlts:  make(types.AddressTables),
	}

	if err := p.initStaticAlts(); err != nil {
		return nil, fmt.Errorf("failed to init static address lookup tables: %w", err)
	}

	return p, nil
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

func (pc *Config) Enabled() bool {
	return !pc.Disabled
}
