package testsuite

import (
	"context"

	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
)

func (s *E2ETestSuite) SetupXCall(ctx context.Context) error {
	createdChains := s.GetChains()
	// chainA, chainB := createdChains[0], createdChains[1]
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
	// if err := chainA.SetupXCall(ctx); err != nil {
	// 	return err
	// }

	// if err := chainA.SetupConnection(ctx, chainB); err != nil {
	// 	return err
	// }

	// if err := chainB.SetupXCall(ctx); err != nil {
	// 	return err
	// }
	// if err := chainB.SetupConnection(ctx, chainA); err != nil {
	// 	return err
	// }
	return nil
}

// SetupChainsAndRelayer create two chains, a relayer, establishes a connection and creates a channel
// using the given channel options. The relayer returned by this function has not yet started. It should be started
// with E2ETestSuite.StartRelayer if needed.
// This should be called at the start of every test, unless fine grained control is required.
func (s *E2ETestSuite) SetupChainsAndRelayer(ctx context.Context) ibc.Relayer {
	relayer, err := s.SetupRelayer(ctx, "centralized")

	s.Require().NoErrorf(err, "Error while configuring relayer %v", err)
	//eRep := s.GetRelayerExecReporter()

	//pathName := s.GeneratePathName()
	//chainA, chainB := s.GetChains()

	//s.Require().NoErrorf(relayer.GeneratePath(ctx, eRep, chainA.(ibc.Chain).Config().ChainID, chainB.(ibc.Chain).Config().ChainID, pathName), "Error on generating path, %v", err)
	//err = relayer.CreateClients(ctx, eRep, pathName, ibc.CreateClientOptions{
	//	TrustingPeriod: "100000m",
	//})
	//s.Require().NoErrorf(err, "Error while creating client relayer : %s, %v", pathName, err)

	//s.Require().NoError(relayer.CreateConnections(ctx, eRep, pathName))
	s.Require().NoError(s.StartRelayer(relayer))
	return relayer
}
