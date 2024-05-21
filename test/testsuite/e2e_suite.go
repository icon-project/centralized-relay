package testsuite

import (
	"context"

	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
)

func (s *E2ETestSuite) SetupXCall(ctx context.Context) error {
	createdChains := s.GetChains()
	for index, chain := range createdChains {
		if err := chain.SetupXCall(ctx); err != nil {
			return err
		}
		for ind, cn := range createdChains {
			if ind != index {
				if err := chain.SetupConnection(ctx, cn); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// SetupChainsAndRelayer create two chains, a relayer, establishes a connection and creates a channel
// using the given channel options. The relayer returned by this function has not yet started. It should be started
// with E2ETestSuite.StartRelayer if needed.
// This should be called at the start of every test, unless fine grained control is required.
func (s *E2ETestSuite) SetupChainsAndRelayer(ctx context.Context) ibc.Relayer {
	relayer, err := s.SetupRelayer(ctx, "centralized")

	s.Require().NoErrorf(err, "Error while configuring relayer %v", err)

	s.Require().NoError(s.StartRelayer(relayer))
	return relayer
}
