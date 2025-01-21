package solana

import (
	"context"
	"fmt"
	"sync"
	"time"

	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/icon-project/centralized-relay/relayer/chains/solana/types"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

type Config struct {
	provider.CommonConfig
	Disabled  bool   `yaml:"disabled" json:"disabled"`
	ChainName string `yaml:"-"`

	RPCUrl  string `yaml:"rpc-url" json:"rpc-url"`
	Address string `yaml:"address" json:"address"`

	XcallProgram string `yaml:"xcall-program" json:"xcall-program"`

	ConnectionProgram string   `yaml:"connection-program" json:"connection-program"`
	OtherConnections  []string `yaml:"other-connections" json:"other-connections"`

	Dapps []types.Dapp `yaml:"dapps" json:"dapps"`

	CpNIDs []string `yaml:"cp-nids" json:"cp-nids"` //counter party NIDs Eg: ["0x2.icon", "0x7.icon"]

	AltAddress string `yaml:"alt-address" json:"alt-address"` // address lookup table address

	NID         string `yaml:"nid" json:"nid"`
	HomeDir     string `yaml:"home-dir" json:"home-dir"`
	StartTxSign string `yaml:"start-tx-sign" json:"start-tx-sign"`

	TxConfirmationTime time.Duration `yaml:"tx-confirmation-time" json:"tx-confirmation-time"`

	ComputeUnitLimit      uint64 `yaml:"compute-unit-limit" json:"compute-unit-limit"`
	Adjustment            uint64 `yaml:"adjustment" json:"adjustment"`
	PriorityFeePercentile uint64 `yaml:"priority-fee-percentile" json:"priority-fee-percentile"`
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
			return nil, fmt.Errorf("failed to fetch xcall idl: %w", err)
		}
	}

	connIdl := IDL{}
	if pc.ConnectionProgram != "" {
		if err := client.FetchIDL(ctx, pc.ConnectionProgram, &connIdl); err != nil {
			return nil, fmt.Errorf("failed to fetch conn idl: %w", err)
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
	return nil
}

func (pc *Config) Enabled() bool {
	return !pc.Disabled
}

func (pc *Config) GetConnContract() string {
	return pc.ConnectionProgram
}

func (pc *Config) ContractsAddress() relayertypes.ContractConfigMap {
	addresses := relayertypes.ContractConfigMap{
		relayertypes.ConnectionContract: pc.ConnectionProgram,
		relayertypes.XcallContract:      pc.XcallProgram,
	}
	for _, dapp := range pc.Dapps {
		addresses[dapp.Name] = dapp.ProgramID
	}

	return addresses
}
