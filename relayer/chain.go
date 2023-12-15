package relayer

import (
	"encoding/json"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/provider"
	"go.uber.org/zap"
)

type Chain struct {
	log           *zap.Logger
	ChainProvider provider.ChainProvider
	debug         bool
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

func (c *Chain) NID() string {
	return c.ChainProvider.NID()
}

// Get returns the configuration for a given chain
func (c Chains) Get(nid string) (*Chain, error) {
	for _, chain := range c {
		if nid == chain.ChainProvider.NID() {
			return chain, nil
		}
	}
	return nil, fmt.Errorf("chain with NID %s is not configured", nid)
}

func (c Chains) GetAll() map[string]*Chain {
	out := make(map[string]*Chain)

	for _, chain := range c {
		out[chain.NID()] = chain
	}
	return out
}
