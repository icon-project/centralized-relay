// Package testsuite provides a suite of end-to-end tests for the IBC relayer.
// This file contains the implementation of the E2ETestSuite struct and its methods.
// The E2ETestSuite struct provides methods for setting up the relayer, creating clients, connections, and channels,
// and executing packet flows between chains.
// It also provides methods for retrieving client, connection, and channel states and sequences.
// All methods in this file use the relayer package to interact with the relayer and the interchaintest package to build and manage interchain networks.
package testsuite

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	interchaintest "github.com/icon-project/centralized-relay/test"
	"github.com/icon-project/centralized-relay/test/chains"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
)

func (s *E2ETestSuite) SetupMockDApp(ctx context.Context, portId string) error {
	chainA, chainB := s.GetChains()
	ctx = context.WithValue(ctx, chains.ContractName{}, chains.ContractName{ContractName: "mockdapp"})
	ibcHostChainA := chainA.GetContractAddress("ibc")
	ctx = context.WithValue(ctx, chains.InitMessageKey("init-msg"), chains.InitMessage{
		Message: map[string]interface{}{
			"ibc_host": ibcHostChainA,
		},
	})
	var err error
	ctx, err = chainA.DeployContract(ctx, interchaintest.XCallOwnerAccount)

	if err != nil {
		return err
	}

	ctx, err = chainA.ExecuteContract(ctx, ibcHostChainA, interchaintest.IBCOwnerAccount, chains.BindPort, map[string]interface{}{
		"port_id": portId,
		"address": chainA.GetContractAddress(GetAppKey(ctx, "mockdapp")),
	})

	if err != nil {
		return err
	}
	ibcHostChainB := chainB.GetContractAddress("ibc")
	ctx = context.WithValue(ctx, chains.InitMessageKey("init-msg"), chains.InitMessage{
		Message: map[string]interface{}{
			"ibc_host": ibcHostChainB,
		},
	})
	ctx, err = chainB.DeployContract(ctx, interchaintest.XCallOwnerAccount)

	if err != nil {
		return err
	}

	ctx, err = chainB.ExecuteContract(ctx, ibcHostChainB, interchaintest.IBCOwnerAccount, chains.BindPort, map[string]interface{}{
		"port_id": portId,
		"address": chainB.GetContractAddress(GetAppKey(ctx, "mockdapp")),
	})

	return err
}

// SendPacket sends a packet from src to dst
func (s *E2ETestSuite) SendPacket(ctx context.Context, src, target chains.Chain, msg string, timeout uint64) (chains.PacketTransferResponse, error) {
	height, _ := src.(ibc.Chain).Height(ctx)
	params := map[string]interface{}{
		"msg":            chains.BufferArray(msg),
		"timeout_height": height + timeout,
	}
	return src.SendPacketMockDApp(ctx, target, interchaintest.UserAccount, params)
}

// CrashRelayer Node
func (s *E2ETestSuite) CrashNode(ctx context.Context, chain chains.Chain) error {
	return chain.PauseNode(ctx)
}

// Resume Node
func (s *E2ETestSuite) ResumeNode(ctx context.Context, chain chains.Chain) error {
	return chain.UnpauseNode(ctx)
}

func (s *E2ETestSuite) CrashRelayer(ctx context.Context, callbacks ...func() error) error {
	eRep := s.GetRelayerExecReporter()
	s.logger.Info("crashing relayer")
	now := time.Now()
	if len(callbacks) > 0 {
		var eg errgroup.Group
		for _, cb := range callbacks {
			eg.Go(cb)
		}
		if err := eg.Wait(); err != nil {
			return err
		}
	}
	err := s.relayer.StopRelayerContainer(ctx, eRep)
	s.logger.Info("relayer crashed", zap.Duration("elapsed", time.Since(now)))
	return err
}

// WriteBlockHeight writes the block height to the given file.
func (s *E2ETestSuite) WriteCurrentBlockHeight(ctx context.Context, chain chains.Chain) func() error {
	return func() error {
		height, err := chain.(ibc.Chain).Height(ctx)
		if err != nil {
			return err
		}
		chanID := chain.(ibc.Chain).Config().ChainID
		return s.WriteBlockHeight(ctx, chanID, height-1)
	}
}

func (s *E2ETestSuite) WriteBlockHeight(ctx context.Context, chainID string, height uint64) error {
	s.T().Logf("updating latest height of %s to %d", chainID, height)
	return s.relayer.WriteBlockHeight(ctx, chainID, height)
}

// Recover recovers a relay and waits for the relay to catch up to the current height of the stopped chains.
// This is because relay needs to sync with the counterchain network when it was on crashed state.
func (s *E2ETestSuite) Recover(ctx context.Context, waitDuration time.Duration) error {
	s.logger.Info("waiting for relayer to restart")
	now := time.Now()
	if err := s.relayer.RestartRelayerContainer(ctx); err != nil {
		return err
	}
	time.Sleep(waitDuration)
	s.logger.Info("relayer restarted", zap.Duration("elapsed", time.Since(now)))
	return nil
}

func GetAppKey(ctx context.Context, contract string) string {
	testcase := ctx.Value("testcase").(string)
	return fmt.Sprintf("%s-%s", contract, testcase)
}
