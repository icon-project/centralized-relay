package testsuite

import (
	"context"
	"fmt"
	setup "github.com/icon-project/centralized-relay/test"
	"github.com/icon-project/centralized-relay/test/chains/evm"
	"github.com/icon-project/centralized-relay/test/chains/icon"
	"github.com/icon-project/centralized-relay/test/interchaintest"
	"github.com/icon-project/centralized-relay/test/interchaintest/testutil"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"strconv"
	"time"

	"strings"

	"github.com/icon-project/centralized-relay/test/chains"

	dockerclient "github.com/docker/docker/client"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/interchaintest/testreporter"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// E2ETestSuite has methods and functionality which can be shared among all test suites.
type E2ETestSuite struct {
	suite.Suite
	relayer ibc.Relayer
	cfg     *testconfig.TestConfig
	//grpcClients    map[string]GRPCClients
	paths          map[string]path
	relayers       ibc.RelayerMap
	logger         *zap.Logger
	DockerClient   *dockerclient.Client
	network        string
	startRelayerFn func(relayer ibc.Relayer) error

	// pathNameIndex is the latest index to be used for generating paths
	pathNameIndex   int64
	CurrentPathName string
	pathNames       []string
}

func (s *E2ETestSuite) SetCfg() error {
	tc, err := testconfig.New()
	if err != nil {
		return err
	}
	s.cfg = tc
	return nil
}

// path is a pairing of two chains which will be used in a test.
type path struct {
	chainA, chainB chains.Chain
}

// newPath returns a path built from the given chains.
func newPath(chainA, chainB chains.Chain) path {
	return path{
		chainA: chainA,
		chainB: chainB,
	}
}

// SetupRelayer sets up the relayer, creates interchain networks, builds chains, and starts the relayer.
// It returns a Relayer interface and an error if any.
func (s *E2ETestSuite) SetupRelayer(ctx context.Context) (ibc.Relayer, error) {
	chainA, chainB := s.GetChains()
	r := interchaintest.New(s.T(), s.cfg.RelayerConfig, s.logger, s.DockerClient, s.network)
	//pathName := s.GeneratePathName()
	ic := interchaintest.NewInterchain().
		AddChain(chainA).
		AddChain(chainB).
		AddRelayer(r, "r").
		AddLink(interchaintest.InterchainLink{
			Chain1:  chainA,
			Chain2:  chainB,
			Relayer: r,
		})

	eRep := s.GetRelayerExecReporter()
	buildOptions := interchaintest.InterchainBuildOptions{
		TestName:          s.T().Name(),
		Client:            s.DockerClient,
		NetworkID:         s.network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation:  true,
	}
	if err := ic.BuildChains(ctx, eRep, buildOptions); err != nil {
		return nil, err
	}

	var err error
	err = s.buildWallets(ctx, chainA, chainB)
	if err != nil {
		return nil, err
	}

	if err := s.SetupXCall(ctx); err != nil {
		return nil, err
	}

	if err := ic.BuildRelayer(ctx, eRep, buildOptions); err != nil {
		return nil, err
	}
	s.startRelayerFn = func(relayer ibc.Relayer) error {
		if err := relayer.StartRelayer(ctx, eRep); err != nil {
			return fmt.Errorf("failed to start relayer: %s", err)
		}
		_ctx, cancel := context.WithTimeout(ctx, time.Second*20)
		defer cancel()
		if err := testutil.WaitForBlocks(_ctx, 2, chainA.(ibc.Chain), chainB.(ibc.Chain)); err != nil {
			return fmt.Errorf("failed to wait for blocks: %v", err)
		}
		return nil
	}
	s.relayer = r
	return r, err
}

func (s *E2ETestSuite) buildWallets(ctx context.Context, chainA chains.Chain, chainB chains.Chain) error {
	if _, err := chainA.BuildWallets(ctx, setup.IBCOwnerAccount); err != nil {
		return err
	}
	if _, err := chainB.BuildWallets(ctx, setup.IBCOwnerAccount); err != nil {
		return err
	}
	if _, err := chainA.BuildWallets(ctx, setup.UserAccount); err != nil {
		return err
	}
	if _, err := chainB.BuildWallets(ctx, setup.UserAccount); err != nil {
		return err
	}
	if _, err := chainA.BuildWallets(ctx, setup.XCallOwnerAccount); err != nil {
		return err
	}
	if _, err := chainB.BuildWallets(ctx, setup.XCallOwnerAccount); err != nil {
		return err
	}
	return nil
}

func (s *E2ETestSuite) DeployXCallMockApp(ctx context.Context, port string) error {
	//testcase := ctx.Value("testcase").(string)

	chainA, chainB := s.GetChains()
	if err := chainA.DeployXCallMockApp(ctx, setup.XCallOwnerAccount, []chains.XCallConnection{{
		Nid:         chainB.(ibc.Chain).Config().ChainID,
		Destination: chainB.GetContractAddress("connection"),
		Connection:  "connection",
	}}); err != nil {
		return err
	}
	if err := chainB.DeployXCallMockApp(ctx, setup.XCallOwnerAccount, []chains.XCallConnection{{
		Nid:         chainA.(ibc.Chain).Config().ChainID,
		Destination: chainA.GetContractAddress("connection"),
		Connection:  "connection",
	}}); err != nil {
		return err
	}
	return nil
}

// GetChains returns two chains that can be used in a test. The pair returned
// is unique to the current test being run. Note: this function does not create containers.
func (s *E2ETestSuite) GetChains(chainOpts ...testconfig.ChainOptionConfiguration) (chains.Chain, chains.Chain) {
	if s.paths == nil {
		s.paths = map[string]path{}
	}

	path, ok := s.paths[s.T().Name()]
	if ok {
		return path.chainA, path.chainB
	}

	chainOptions, err := testconfig.DefaultChainOptions()
	s.Require().NoError(err)
	for _, opt := range chainOpts {
		opt(chainOptions)
	}

	chainA, chainB := s.createChains(chainOptions)
	path = newPath(chainA, chainB)
	s.paths[s.T().Name()] = path
	return path.chainA, path.chainB
}

// GetRelayerWallets returns the relayer wallets associated with the chains.
func (s *E2ETestSuite) GetRelayerWallets(relayer ibc.Relayer) (ibc.Wallet, ibc.Wallet, error) {
	chainA, chainB := s.GetChains()
	chainARelayerWallet, ok := relayer.GetWallet(chainA.(ibc.Chain).Config().ChainID)
	if !ok {
		return nil, nil, fmt.Errorf("unable to find chain A relayer wallet")
	}

	chainBRelayerWallet, ok := relayer.GetWallet(chainB.(ibc.Chain).Config().ChainID)
	if !ok {
		return nil, nil, fmt.Errorf("unable to find chain B relayer wallet")
	}
	return chainARelayerWallet, chainBRelayerWallet, nil
}

// StartRelayer starts the given relayer.
func (s *E2ETestSuite) StartRelayer(relayer ibc.Relayer) error {
	if s.startRelayerFn == nil {
		return fmt.Errorf("cannot start relayer before it is created: %v", relayer)
	}
	return s.startRelayerFn(relayer)
}

// StopRelayer stops the given relayer.
func (s *E2ETestSuite) StopRelayer(ctx context.Context, relayer ibc.Relayer) error {
	err := relayer.StopRelayer(ctx, s.GetRelayerExecReporter())
	return err
}

// createChains creates two separate chains in docker containers.
// test and can be retrieved with GetChains.
func (s *E2ETestSuite) createChains(chainOptions *testconfig.ChainOptions) (chains.Chain, chains.Chain) {
	client, network := interchaintest.DockerSetup(s.T())
	t := s.T()

	s.logger = zap.NewExample()
	s.DockerClient = client
	s.network = network

	logger := zaptest.NewLogger(t)

	chainA, _ := buildChain(logger, t.Name(), chainOptions.ChainAConfig)

	chainB, _ := buildChain(logger, t.Name(), chainOptions.ChainBConfig)

	// this is intentionally called after the setup.DockerSetup function. The above function registers a
	// cleanup task which deletes all containers. By registering a cleanup function afterwards, it is executed first
	// this allows us to process the logs before the containers are removed.
	//t.Cleanup(func() {
	//	diagnostics.Collect(t, s.DockerClient, chainOptions)
	//})

	return chainA, chainB
}

func buildChain(log *zap.Logger, testName string, cfg *testconfig.Chain) (chains.Chain, error) {
	var (
		chain chains.Chain
	)
	ibcChainConfig := cfg.ChainConfig.GetIBCChainConfig(&chain)
	switch cfg.ChainConfig.Type {
	case "icon":
		chain = icon.NewIconLocalnet(testName, log, ibcChainConfig, chains.DefaultNumValidators, chains.DefaultNumFullNodes, cfg.Contracts)
		return chain, nil
	case "evm":
		chain = evm.NewEVMLocalnet(testName, log, ibcChainConfig, chains.DefaultNumValidators, chains.DefaultNumFullNodes, cfg.Contracts)
		return chain, nil
	default:
		return nil, fmt.Errorf("unexpected error, unknown chain type: %s for chain: %s", cfg.ChainConfig.Type, cfg.Name)
	}
}

// GetRelayerExecReporter returns a testreporter.RelayerExecReporter instances
// using the current test's testing.T.
func (s *E2ETestSuite) GetRelayerExecReporter() *testreporter.RelayerExecReporter {
	rep := testreporter.NewNopReporter()
	return rep.RelayerExecReporter(s.T())
}

func (s *E2ETestSuite) ConvertToPlainString(input string) (string, error) {
	var plainString []byte
	if strings.HasPrefix(input, "[") && strings.HasSuffix(input, "]") {
		input = input[1 : len(input)-1]
		for _, part := range strings.Split(input, ", ") {
			value, err := strconv.Atoi(part)
			if err != nil {
				return "", err
			}
			plainString = append(plainString, byte(value))
		}
		return string(plainString), nil
	} else if strings.HasPrefix(input, "0x") {
		input = input[2:]
		for i := 0; i < len(input); i += 2 {
			value, err := strconv.ParseUint(input[i:i+2], 16, 8)
			if err != nil {
				return "", err
			}
			plainString = append(plainString, byte(value))
		}
		return string(plainString), nil
	}
	return "", fmt.Errorf("invalid input length")
}
