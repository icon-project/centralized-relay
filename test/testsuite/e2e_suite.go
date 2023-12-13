package testsuite

import (
	"context"
	interchaintest "github.com/icon-project/centralized-relay/test"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
)

func (s *E2ETestSuite) SetupXCall(ctx context.Context) error {
	chainA, chainB := s.GetChains()
	if err := chainA.SetupXCall(ctx, interchaintest.XCallOwnerAccount); err != nil {
		return err
	}
	if err := chainA.SetupConnection(ctx, interchaintest.XCallOwnerAccount, chainB); err != nil {
		return err
	}

	if err := chainB.SetupXCall(ctx, interchaintest.XCallOwnerAccount); err != nil {
		return err
	}
	if err := chainB.SetupConnection(ctx, interchaintest.XCallOwnerAccount, chainA); err != nil {
		return err
	}
	return nil
}

// SetupChainsAndRelayer create two chains, a relayer, establishes a connection and creates a channel
// using the given channel options. The relayer returned by this function has not yet started. It should be started
// with E2ETestSuite.StartRelayer if needed.
// This should be called at the start of every test, unless fine grained control is required.
func (s *E2ETestSuite) SetupChainsAndRelayer(ctx context.Context) ibc.Relayer {
	relayer, err := s.SetupRelayer(ctx)
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
