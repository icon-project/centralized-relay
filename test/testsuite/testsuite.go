package testsuite

import (
	"context"
	"fmt"
	"strconv"

	"github.com/icon-project/centralized-relay/test/chains/cosmos"
	"github.com/icon-project/centralized-relay/test/chains/evm"
	"github.com/icon-project/centralized-relay/test/chains/icon"
	"github.com/icon-project/centralized-relay/test/interchaintest"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"

	"strings"

	"github.com/icon-project/centralized-relay/test/chains"

	dockerclient "github.com/docker/docker/client"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/interchaintest/testreporter"
	ibcv8 "github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// E2ETestSuite has methods and functionality which can be shared among all test suites.
type E2ETestSuite struct {
	suite.Suite

	cfg *testconfig.TestConfig
	//grpcClients    map[string]GRPCClients
	paths           map[string]path
	Relayers        map[string]ibc.Relayer
	RelayersWallets ibc.RelayerMap
	logger          *zap.Logger
	DockerClient    *dockerclient.Client
	network         string
	startRelayerFn  func(relayer ibc.Relayer) error

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
func (s *E2ETestSuite) SetupRelayer(ctx context.Context, name string) (ibc.Relayer, error) {
	createdChains := s.GetChains()
	r := interchaintest.New(s.T(), s.cfg.RelayerConfig, s.logger, s.DockerClient, s.network)
	ic := interchaintest.NewInterchain()

	for index := range createdChains {
		ic.AddChain(createdChains[index])
	}
	ic.AddRelayer(r, "r").
		AddLink(
			interchaintest.InterchainLink{
				Chains:  createdChains,
				Relayer: r,
			},
		)

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

	if err := s.SetupXCall(ctx); err != nil {
		return nil, err
	}

	if err := ic.BuildRelayer(ctx, eRep, buildOptions, s.cfg.RelayerConfig.KMS_ID); err != nil {
		return nil, err
	}
	s.startRelayerFn = func(relayer ibc.Relayer) error {
		if err := relayer.StartRelayer(ctx, eRep); err != nil {
			return fmt.Errorf("failed to start relayer: %s", err)
		}
		return nil
	}
	if s.Relayers == nil {
		s.Relayers = make(map[string]ibc.Relayer)
	}
	s.Relayers[name] = r
	return r, err
}

func (s *E2ETestSuite) DeployXCallMockApp(ctx context.Context, port string) error {
	createdChains := s.GetChains()
	// chainA, chainB := createdChains[0], createdChains[1]
	for idx, chain := range createdChains {
		var connections []chains.XCallConnection
		for id, cn := range createdChains {
			if id != idx {
				connections = append(connections, chains.XCallConnection{
					Nid:         cn.(ibc.Chain).Config().ChainID,
					Destination: cn.GetContractAddress("connection"),
					Connection:  "connection",
				})
			}
		}
		if err := chain.DeployXCallMockApp(ctx, "cn.(ibc.Chain).Config().Name", connections); err != nil {
			return err
		}
	}
	return nil
}

// GetChains returns two chains that can be used in a test. The pair returned
// is unique to the current test being run. Note: this function does not create containers.
func (s *E2ETestSuite) GetChains(chainOpts ...testconfig.ChainOptionConfiguration) []chains.Chain {
	if s.paths == nil {
		s.paths = map[string]path{}
	}
	preCreatedChains := []chains.Chain{}
	for i := 0; i <= 10; i++ {
		pathKey := fmt.Sprintf("%s-%d", s.T().Name(), i)
		path, ok := s.paths[pathKey]
		if ok {
			if len(preCreatedChains) == 0 {
				preCreatedChains = append(preCreatedChains, path.chainA, path.chainB)
			} else {
				preCreatedChains = append(preCreatedChains, path.chainB)
			}
		} else {
			if len(preCreatedChains) != 0 {
				return preCreatedChains
			}
		}
	}
	chainOptions, err := testconfig.DefaultChainOptions()
	s.Require().NoError(err)
	for _, opt := range chainOpts {
		opt(chainOptions)
	}

	createdChains := s.createChains(chainOptions)
	for index := range createdChains {
		if index < len(createdChains)-1 {
			path := newPath(createdChains[index], createdChains[index+1])
			pathKey := fmt.Sprintf("%s-%d", s.T().Name(), index)
			s.paths[pathKey] = path
		}
	}
	// chainA, chainB := createdChains[0], createdChains[1]
	// path = newPath(chainA, chainB)
	// s.paths[s.T().Name()] = path
	return createdChains
}

// GetRelayerWallets returns the relayer wallets associated with the chains.
func (s *E2ETestSuite) GetRelayerWallets(relayer ibc.Relayer) (ibc.Wallet, ibc.Wallet, error) {
	chains := s.GetChains()
	chainA, chainB := chains[0], chains[1]
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
func (s *E2ETestSuite) createChains(chainOptions *testconfig.ChainOptions) []chains.Chain {
	client, network := interchaintest.DockerSetup(s.T())
	t := s.T()

	s.logger = zap.NewExample()
	s.DockerClient = client
	s.network = network

	logger := zaptest.NewLogger(t)
	chains := []chains.Chain{}

	for _, config := range *chainOptions.ChainConfig {
		chain, _ := buildChain(logger, t.Name(), s, &config)
		chains = append(chains,
			chain,
		)
	}

	// chainA, _ := buildChain(logger, t.Name(), s, chainOptions.ChainAConfig)

	// chainB, _ := buildChain(logger, t.Name(), s, chainOptions.ChainBConfig)

	// this is intentionally called after the setup.DockerSetup function. The above function registers a
	// cleanup task which deletes all containers. By registering a cleanup function afterwards, it is executed first
	// this allows us to process the logs before the containers are removed.
	//t.Cleanup(func() {
	//	diagnostics.Collect(t, s.DockerClient, chainOptions)
	//})

	return chains
}

func toInterchantestConfig(config ibc.ChainConfig) ibcv8.ChainConfig {

	images := []ibcv8.DockerImage{
		ibcv8.NewDockerImage(
			config.Images[0].Repository, config.Images[0].Version, config.Images[0].UidGid),
	}
	decimals := int64(6)
	return ibcv8.ChainConfig{
		Type:           config.Type,
		Name:           config.Name,
		ChainID:        config.ChainID,
		Images:         images,
		Bin:            config.Bin,
		Bech32Prefix:   config.Bech32Prefix,
		Denom:          config.Denom,
		SkipGenTx:      config.SkipGenTx,
		CoinType:       config.CoinType,
		GasPrices:      config.GasPrices,
		GasAdjustment:  config.GasAdjustment,
		TrustingPeriod: config.TrustingPeriod,
		NoHostMount:    config.NoHostMount,
		CoinDecimals:   &decimals,
	}
}

func buildChain(log *zap.Logger, testName string, s *E2ETestSuite, cfg *testconfig.Chain) (chains.Chain, error) {
	var (
		chain chains.Chain
	)
	ibcChainConfig := cfg.ChainConfig.GetIBCChainConfig(&chain)
	switch cfg.ChainConfig.Type {
	case "icon":
		chain = icon.NewIconRemotenet(testName, log, ibcChainConfig, s.DockerClient, s.network, cfg)
		return chain, nil
	case "evm":
		chain = evm.NewEVMRemotenet(testName, log, ibcChainConfig, s.DockerClient, s.network, cfg)
		return chain, nil
	case "wasm", "cosmos":
		interchainTestConfig := toInterchantestConfig(ibcChainConfig)
		chain, err := cosmos.NewCosmosRemotenet(testName, log, interchainTestConfig, s.DockerClient, s.network, cfg)
		return chain, err
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

func (s *E2ETestSuite) ConvertToPlainString(input string) string {
	var plainString []byte
	if strings.HasPrefix(input, "[") && strings.HasSuffix(input, "]") {
		input = input[1 : len(input)-1]
		for _, part := range strings.Split(input, ", ") {
			value, err := strconv.Atoi(part)
			if err != nil {
				return ""
			}
			plainString = append(plainString, byte(value))
		}
		return string(plainString)
	} else if strings.HasPrefix(input, "0x") {
		input = input[2:]
		for i := 0; i < len(input); i += 2 {
			value, err := strconv.ParseUint(input[i:i+2], 16, 8)
			if err != nil {
				return ""
			}
			plainString = append(plainString, byte(value))
		}
		return string(plainString)
	}
	return input
}
