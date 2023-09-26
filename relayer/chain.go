package relayer

import (
	"encoding/json"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/provider"
	"go.uber.org/zap"
)

type Chain struct {
	log *zap.Logger

	ChainProvider provider.ChainProvider
	Chainid       string `yaml:"chain-id" json:"chain-id"`
	RPCAddr       string `yaml:"rpc-addr" json:"rpc-addr"`

	debug bool
}

// Chains is a collection of Chain (mapped by chain_name)
type Chains map[string]*Chain

func NewChain(log *zap.Logger, prov provider.ChainProvider, debug bool) *Chain {
	return &Chain{
		log:           log,
		ChainProvider: prov,
		debug:         debug,
	}
}

func (c *Chain) String() string {
	out, _ := json.Marshal(c)
	return string(out)
}

func (c *Chain) ChainID() string {
	return c.ChainProvider.ChainId()
}

// Get returns the configuration for a given chain
func (c Chains) Get(chainID string) (*Chain, error) {
	for _, chain := range c {
		if chainID == chain.ChainProvider.ChainId() {
			return chain, nil
		}
	}
	return nil, fmt.Errorf("chain with ID %s is not configured", chainID)
}

// MustGet returns the chain and panics on any error
func (c Chains) MustGet(chainID string) *Chain {
	out, err := c.Get(chainID)
	if err != nil {
		panic(err)
	}
	return out
}

// Gets returns a map chainIDs to their chains
func (c Chains) Gets(chainIDs ...string) (map[string]*Chain, error) {
	out := make(map[string]*Chain)
	for _, cid := range chainIDs {
		chain, err := c.Get(cid)
		if err != nil {
			return out, err
		}
		out[cid] = chain
	}
	return out, nil
}
