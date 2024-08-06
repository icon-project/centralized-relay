package interchaintest

import (
	chain "github.com/icon-project/centralized-relay/test/chains"

	"go.uber.org/zap"
)

// chainSet is an unordered collection of chains.Chain,
// to group methods that apply actions against all chains in the set.
//
// The main purpose of the chainSet is to unify test setup when working with any number of chains.
type chainSet struct {
	log *zap.Logger

	chains map[chain.Chain]struct{}
}

func newChainSet(log *zap.Logger, chains []chain.Chain) *chainSet {
	cs := &chainSet{
		log: log,

		chains: make(map[chain.Chain]struct{}, len(chains)),
	}

	for _, chain := range chains {
		cs.chains[chain] = struct{}{}
	}

	return cs
}
