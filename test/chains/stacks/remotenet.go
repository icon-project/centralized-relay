package stacks

import (
	"context"
	"fmt"
	"strings"
	"time"

	stacksClient "github.com/icon-project/centralized-relay/relayer/chains/stacks"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/test/chains"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer/centralized"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"github.com/icon-project/stacks-go-sdk/pkg/clarity"
	"github.com/icon-project/stacks-go-sdk/pkg/stacks"
	"github.com/icon-project/stacks-go-sdk/pkg/transaction"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const BLOCK_TIME = 5 * time.Second
const MAX_WAIT_TIME = 200 * BLOCK_TIME

type contextKey string

const (
	snContextKey contextKey = "sn"
)

type StacksLocalnet struct {
	log          *zap.Logger
	testName     string
	cfg          chains.ChainConfig
	IBCAddresses map[string]string
	Wallets      map[string]string // map of keyName to private key (as hex string)
	network      *stacks.StacksNetwork
	client       *stacksClient.Client
	testconfig   *testconfig.Chain
	kms          kms.KMS
}

func NewStacksLocalnet(testName string, log *zap.Logger, chainConfig chains.ChainConfig, testconfig *testconfig.Chain, kms kms.KMS) chains.Chain {
	network := stacks.NewStacksTestnet()
	client, err := stacksClient.NewClient(log, network)
	if err != nil {
		log.Error("Failed to create Stacks client", zap.Error(err))
		return nil
	}

	log.Debug("Creating Stacks chain",
		zap.String("chainID", chainConfig.ChainID),
		zap.String("name", chainConfig.Name),
		zap.String("type", chainConfig.Type))

	return &StacksLocalnet{
		testName:     testName,
		cfg:          chainConfig,
		log:          log,
		network:      network,
		client:       client,
		testconfig:   testconfig,
		IBCAddresses: make(map[string]string),
		Wallets:      make(map[string]string),
		kms:          kms,
	}
}

func (s *StacksLocalnet) Config() chains.ChainConfig {
	return s.cfg
}

func (s *StacksLocalnet) Height(ctx context.Context) (uint64, error) {
	block, err := s.client.GetLatestBlock(ctx)
	if err != nil {
		return 0, err
	}
	return uint64(block.Height), nil
}

func (s *StacksLocalnet) GetRelayConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error) {
	contracts := make(map[string]string)
	contracts["xcall-proxy"] = s.GetContractAddress("xcall-proxy")
	contracts["xcall"] = s.GetContractAddress("xcall-proxy")
	contracts["connection"] = s.GetContractAddress("connection")

	config := &centralized.StacksRelayerChainConfig{
		Type: "stacks",
		Value: centralized.StacksRelayerChainConfigValue{
			NID:           s.testconfig.NID,
			RPCURL:        s.testconfig.RPCUri,
			StartHeight:   0,
			Contracts:     contracts,
			BlockInterval: "6s",
			Address:       s.testconfig.RelayWalletAddress,
			FinalityBlock: 10,
		},
	}
	return yaml.Marshal(config)
}

func (s *StacksLocalnet) GetContractAddress(key string) string {
	value, exist := s.IBCAddresses[key]
	if !exist {
		panic(fmt.Sprintf("IBC address does not exist: %s", key))
	}
	return value
}

func (s *StacksLocalnet) SendPacketXCall(ctx context.Context, keyName, _to string, data, rollback []byte) (context.Context, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	dappAddress := s.IBCAddresses[dappKey]

	privateKey, senderAddress, err := s.loadPrivateKey(s.testconfig.KeystoreFile, s.testconfig.KeystorePassword)
	if err != nil {
		return nil, fmt.Errorf("failed to load user's private key: %w", err)
	}

	toArg, err := clarity.NewStringASCII(_to)
	if err != nil {
		return nil, fmt.Errorf("failed to create 'to' argument: %w", err)
	}

	dataArg := clarity.NewBuffer(data)
	var rollbackArg clarity.ClarityValue
	if len(rollback) > 0 {
		rollbackArg = clarity.NewOptionSome(clarity.NewBuffer(rollback))
	} else {
		rollbackArg = clarity.NewOptionNone()
	}

	implAddr := s.IBCAddresses["xcall-impl"]
	implPrincipal, err := clarity.StringToPrincipal(implAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create implementation argument: %w", err)
	}

	args := []clarity.ClarityValue{
		toArg,
		dataArg,
		rollbackArg,
		implPrincipal,
	}

	parts := strings.Split(dappAddress, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid dapp address format: %s", dappAddress)
	}
	contractName := parts[1]

	s.log.Debug("Contract addresses",
		zap.String("dapp", dappAddress),
		zap.String("impl", s.IBCAddresses["xcall-impl"]),
		zap.String("proxy", s.IBCAddresses["xcall-proxy"]),
		zap.String("connection", s.IBCAddresses["connection"]))

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid dapp address format: %s", dappAddress)
	}

	txCall, err := transaction.MakeContractCall(
		senderAddress,
		contractName,
		"send-message",
		args,
		*s.network,
		senderAddress,
		privateKey,
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract call transaction: %w", err)
	}

	txID, err := transaction.BroadcastTransaction(txCall, s.network)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	s.log.Info("Sent message via xcall mock app", zap.String("txID", txID))

	err = s.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return nil, err
	}

	sn, err := s.extractSnFromTransaction(ctx, txID)
	if err != nil {
		return nil, fmt.Errorf("failed to extract serial number: %w", err)
	}

	s.log.Info("Successfully sent packet",
		zap.String("txID", txID),
		zap.String("sn", sn),
		zap.String("to", _to))

	ctx = context.WithValue(ctx, snContextKey, sn)
	return ctx, nil
}

func (s *StacksLocalnet) XCall(ctx context.Context, targetChain chains.Chain, keyName, _to string, data, rollback []byte) (*chains.XCallResponse, error) {
	height, err := targetChain.Height(ctx)
	if err != nil {
		return nil, err
	}

	ctx, err = s.SendPacketXCall(ctx, keyName, _to, data, rollback)
	if err != nil {
		return nil, err
	}

	sn := ctx.Value(snContextKey).(string)
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	from := s.cfg.ChainID + "/" + s.GetContractAddress(dappKey)
	toAddress := strings.Split(_to, "/")[1]

	reqID, destData, err := targetChain.FindCallMessage(ctx, height, from, toAddress, sn)
	if err != nil {
		return nil, err
	}

	return &chains.XCallResponse{
		SerialNo:  sn,
		RequestID: reqID,
		Data:      destData,
	}, nil
}
