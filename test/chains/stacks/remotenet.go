package stacks

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	stacksClient "github.com/icon-project/centralized-relay/relayer/chains/stacks"
	"github.com/icon-project/centralized-relay/relayer/chains/stacks/events"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/test/chains"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer/centralized"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"github.com/icon-project/stacks-go-sdk/pkg/clarity"
	"github.com/icon-project/stacks-go-sdk/pkg/crypto"
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
	client, err := stacksClient.NewClient(log, network, testconfig.Contracts["xcall-abi"])
	if err != nil {
		log.Error("Failed to create Stacks client", zap.Error(err))
		return nil
	}

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
	contracts["connection"] = s.GetContractAddress("connection")

	config := &centralized.StacksRelayerChainConfig{
		Type: "stacks",
		Value: centralized.StacksRelayerChainConfigValue{
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

func (s *StacksLocalnet) SetupXCall(ctx context.Context) error {
	if s.testconfig.Environment == "preconfigured" {
		testcase := ctx.Value("testcase").(string)
		s.IBCAddresses["xcall-proxy"] = "STXCALLPROXYADDRESS"
		s.IBCAddresses["connection"] = "STXCONNECTIONADDRESS"
		s.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = "STXDAPPADDRESS"
		return nil
	}

	privateKey, senderAddress, err := s.loadPrivateKey(s.testconfig.KeystoreFile, s.testconfig.KeystorePassword)
	if err != nil {
		return fmt.Errorf("failed to load deployer's private key: %w", err)
	}

	deployments := []struct {
		name       string                             // Contract name for deployment
		contract   string                             // Key in testconfig.Contracts
		wait       bool                               // Whether to wait for confirmation
		postDeploy func(contractAddress string) error // Optional post-deployment initialization
	}{
		{"xcall-common-trait", "common-trait", true, nil},
		{"xcall-receiver-trait", "receiver-trait", true, nil},
		{"xcall-impl-trait", "impl-trait", true, nil},
		{"xcall-proxy-trait", "proxy-trait", true, nil},

		{"util", "util", true, nil},
		{"rlp-encode", "rlp-encode", true, nil},
		{"rlp-decode", "rlp-decode", true, nil},

		{"xcall-proxy", "xcall-proxy", true, nil},

		{"centralized-connection", "connection", true, nil},

		{"xcall-impl-v3", "xcall-impl", true, func(proxyAddr string) error {
			implAddr := senderAddress + ".xcall-impl-v3"
			implPrincipal, err := clarity.StringToPrincipal(implAddr)
			if err != nil {
				return fmt.Errorf("failed to convert implementation address to principal: %w", err)
			}

			upgradeTx, err := transaction.MakeContractCall(
				senderAddress,
				"xcall-proxy",
				"upgrade",
				[]clarity.ClarityValue{
					implPrincipal,
					clarity.NewOptionNone(),
				},
				*s.network,
				senderAddress,
				privateKey,
				nil,
				nil,
			)
			if err != nil {
				return fmt.Errorf("failed to create upgrade transaction: %w", err)
			}

			txID, err := transaction.BroadcastTransaction(upgradeTx, s.network)
			if err != nil {
				return fmt.Errorf("failed to broadcast upgrade transaction: %w", err)
			}

			err = s.waitForTransactionConfirmation(ctx, txID)
			if err != nil {
				return fmt.Errorf("failed to confirm upgrade transaction: %w", err)
			}

			err = s.initializeXCallImpl(ctx, privateKey, senderAddress)
			if err != nil {
				return fmt.Errorf("failed to initialize xcall-impl: %w", err)
			}

			err = s.setAdminXCallImpl(ctx, privateKey, senderAddress)
			if err != nil {
				return fmt.Errorf("failed to set admin for xcall-impl: %w", err)
			}

			err = s.setDefaultConnection(ctx, privateKey, senderAddress)
			if err != nil {
				return fmt.Errorf("failed to set default connection: %w", err)
			}

			err = s.setProtocolFeeHandler(ctx, privateKey, senderAddress)
			if err != nil {
				return fmt.Errorf("failed to set protocol fee handler: %w", err)
			}

			err = s.setConnectionFees(ctx, privateKey, senderAddress)
			if err != nil {
				return fmt.Errorf("failed to set connection fees: %w", err)
			}

			err = s.setProtocolFee(ctx, privateKey, senderAddress)
			if err != nil {
				return fmt.Errorf("failed to set protocol fee: %w", err)
			}

			s.log.Info("Initialized proxy with implementation",
				zap.String("proxy", proxyAddr),
				zap.String("implementation", implAddr))

			return nil
		}},
	}

	connectionContractName := "centralized-connection"
	implContractName := "xcall-impl-v3"
	proxyContractName := "xcall-impl-v3"

	s.IBCAddresses["xcall-proxy"] = senderAddress + "." + proxyContractName
	s.IBCAddresses["xcall-impl"] = senderAddress + "." + implContractName
	s.IBCAddresses["connection"] = senderAddress + "." + connectionContractName

	deployedContracts := make(map[string]string)

	for _, deployment := range deployments {
		contractAddress := senderAddress + "." + deployment.name
		deployedContracts[deployment.name] = contractAddress

		contract, err := s.client.GetContractById(ctx, contractAddress)
		if err != nil {
			return fmt.Errorf("failed to check contract existence for %s: %w", deployment.name, err)
		}
		if contract != nil {
			txResp, err := s.client.GetTransactionById(ctx, contract.TxId)
			if err == nil {
				receipt, err := stacksClient.GetReceipt(txResp)
				if err == nil && receipt.Status {
					s.log.Info("Contract already successfully deployed, skipping",
						zap.String("contract", deployment.name),
						zap.String("address", contractAddress))
					continue
				}
			}
		}

		codeBody, err := os.ReadFile(s.testconfig.Contracts[deployment.contract])
		if err != nil {
			return fmt.Errorf("failed to read contract code for %s: %w", deployment.name, err)
		}

		tx, err := transaction.MakeContractDeploy(
			deployment.name,
			string(codeBody),
			*s.network,
			senderAddress,
			privateKey,
			nil,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to create contract deploy transaction for %s: %w", deployment.name, err)
		}

		txID, err := transaction.BroadcastTransaction(tx, s.network)
		if err != nil {
			return fmt.Errorf("failed to broadcast transaction for %s: %w", deployment.name, err)
		}

		s.log.Info("Deployed contract",
			zap.String("contract", deployment.name),
			zap.String("txID", txID))

		if deployment.wait {
			err = s.waitForTransactionConfirmation(ctx, txID)
			if err != nil {
				return fmt.Errorf("failed to confirm transaction for %s: %w", deployment.name, err)
			}
		}

		if deployment.postDeploy != nil {
			if err := deployment.postDeploy(contractAddress); err != nil {
				return fmt.Errorf("post-deployment initialization failed for %s: %w", deployment.name, err)
			}
		}
	}

	s.IBCAddresses["xcall-proxy"] = deployedContracts["xcall-proxy"]
	s.IBCAddresses["xcall-impl"] = deployedContracts["xcall-impl-v3"]
	s.IBCAddresses["connection"] = deployedContracts["centralized-connection"]

	return nil
}

func (s *StacksLocalnet) setDefaultConnection(ctx context.Context, privateKey []byte, senderAddress string) error {
	nid := "test"
	connectionAddress := senderAddress + "." + "centralized-connection"

	nidArg, err := clarity.NewStringASCII(nid)
	if err != nil {
		return fmt.Errorf("failed to create nid argument: %w", err)
	}

	connArg, err := clarity.NewStringASCII(connectionAddress)
	if err != nil {
		return fmt.Errorf("failed to create connection address argument: %w", err)
	}

	implAddr := senderAddress + "." + "xcall-impl-v3"
	implPrincipal, err := clarity.StringToPrincipal(implAddr)
	if err != nil {
		return fmt.Errorf("failed to convert implementation address to principal: %w", err)
	}

	args := []clarity.ClarityValue{
		nidArg,
		connArg,
		implPrincipal,
	}

	txCall, err := transaction.MakeContractCall(
		senderAddress,
		"xcall-proxy",
		"set-default-connection",
		args,
		*s.network,
		senderAddress,
		privateKey,
		nil,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to create set-default-connection transaction: %w", err)
	}

	txID, err := transaction.BroadcastTransaction(txCall, s.network)
	if err != nil {
		return fmt.Errorf("failed to broadcast set-default-connection transaction: %w", err)
	}

	s.log.Info("Set default connection in xcall-impl", zap.String("txID", txID))

	err = s.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return err
	}

	return nil
}

func (s *StacksLocalnet) initializeXCallImpl(ctx context.Context, privateKey []byte, senderAddress string) error {
	nid := s.cfg.ChainID
	addr := senderAddress

	nidArg, err := clarity.NewStringASCII(nid)
	if err != nil {
		return fmt.Errorf("failed to create nid argument: %w", err)
	}

	addrArg, err := clarity.NewStringASCII(addr)
	if err != nil {
		return fmt.Errorf("failed to create addr argument: %w", err)
	}

	args := []clarity.ClarityValue{nidArg, addrArg}

	txCall, err := transaction.MakeContractCall(
		senderAddress,
		"xcall-impl-v3",
		"init",
		args,
		*s.network,
		senderAddress,
		privateKey,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create init transaction: %w", err)
	}

	txID, err := transaction.BroadcastTransaction(txCall, s.network)
	if err != nil {
		return fmt.Errorf("failed to broadcast init transaction: %w", err)
	}

	s.log.Info("Initialized xcall-impl contract", zap.String("txID", txID))

	err = s.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return err
	}

	return nil
}

func (s *StacksLocalnet) setConnectionFees(ctx context.Context, privateKey []byte, senderAddress string) error {
	networks := []string{"stacks_testnet", "test"}
	messageFees := map[string]uint64{"stacks_testnet": 500000, "test": 1000000}
	responseFees := map[string]uint64{"stacks_testnet": 250000, "test": 500000}

	for _, nid := range networks {
		nidArg, err := clarity.NewStringASCII(nid)
		if err != nil {
			return fmt.Errorf("failed to create nid argument: %w", err)
		}

		messageFeeArg, _ := clarity.NewUInt(messageFees[nid])
		responseFeeArg, _ := clarity.NewUInt(responseFees[nid])

		args := []clarity.ClarityValue{
			nidArg,
			messageFeeArg,
			responseFeeArg,
		}

		txCall, err := transaction.MakeContractCall(
			senderAddress,
			"centralized-connection",
			"set-fee",
			args,
			*s.network,
			senderAddress,
			privateKey,
			nil,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to create set-fee transaction: %w", err)
		}

		txID, err := transaction.BroadcastTransaction(txCall, s.network)
		if err != nil {
			return fmt.Errorf("failed to broadcast set-fee transaction: %w", err)
		}

		s.log.Info("Set fee in centralized-connection", zap.String("txID", txID), zap.String("nid", nid))

		err = s.waitForTransactionConfirmation(ctx, txID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *StacksLocalnet) setProtocolFee(ctx context.Context, privateKey []byte, senderAddress string) error {
	protocolFee := uint64(100000)
	protocolFeeClarity, _ := clarity.NewUInt(protocolFee)

	implAddr := s.IBCAddresses["xcall-impl"]
	implPrincipal, err := clarity.StringToPrincipal(implAddr)
	if err != nil {
		return fmt.Errorf("failed to convert implPrincipal to principal: %w", err)
	}

	args := []clarity.ClarityValue{
		protocolFeeClarity,
		implPrincipal,
	}

	txCall, err := transaction.MakeContractCall(
		senderAddress,
		"xcall-proxy",
		"set-protocol-fee",
		args,
		*s.network,
		senderAddress,
		privateKey,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create set-protocol-fee transaction: %w", err)
	}

	txID, err := transaction.BroadcastTransaction(txCall, s.network)
	if err != nil {
		return fmt.Errorf("failed to broadcast set-protocol-fee transaction: %w", err)
	}

	s.log.Info("Set protocol fee in xcall-proxy", zap.String("txID", txID))

	err = s.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return err
	}

	return nil
}

func (s *StacksLocalnet) setAdminXCallImpl(ctx context.Context, privateKey []byte, senderAddress string) error {
	senderPrincipal, err := clarity.StringToPrincipal(senderAddress)
	if err != nil {
		return fmt.Errorf("failed to convert senderAddress to principal: %w", err)
	}

	args := []clarity.ClarityValue{
		senderPrincipal,
	}

	txCall, err := transaction.MakeContractCall(
		senderAddress,
		"xcall-impl-v3",
		"set-admin",
		args,
		*s.network,
		senderAddress,
		privateKey,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create set-admin transaction: %w", err)
	}

	txID, err := transaction.BroadcastTransaction(txCall, s.network)
	if err != nil {
		return fmt.Errorf("failed to broadcast set-admin transaction: %w", err)
	}

	s.log.Info("Set admin for xcall-impl", zap.String("txID", txID))

	err = s.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return err
	}

	return nil
}

func (s *StacksLocalnet) setProtocolFeeHandler(ctx context.Context, privateKey []byte, senderAddress string) error {
	connAddr := s.IBCAddresses["connection"]
	connPrincipal, err := clarity.StringToPrincipal(connAddr)
	if err != nil {
		return fmt.Errorf("failed to create connection principal: %w", err)
	}

	implAddr := s.IBCAddresses["xcall-impl"]

	implPrincipal, err := clarity.StringToPrincipal(implAddr)
	if err != nil {
		return fmt.Errorf("failed to convert implPrincipal to principal: %w", err)
	}

	args := []clarity.ClarityValue{
		connPrincipal,
		implPrincipal,
	}

	txCall, err := transaction.MakeContractCall(
		senderAddress,
		"xcall-proxy",
		"set-protocol-fee-handler",
		args,
		*s.network,
		senderAddress,
		privateKey,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create set-protocol-fee-handler transaction: %w", err)
	}

	txID, err := transaction.BroadcastTransaction(txCall, s.network)
	if err != nil {
		return fmt.Errorf("failed to broadcast set-protocol-fee-handler transaction: %w", err)
	}

	s.log.Info("Set protocol fee handler in xcall-proxy", zap.String("txID", txID))

	err = s.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return err
	}

	return nil
}

func (s *StacksLocalnet) SetupConnection(ctx context.Context, target chains.Chain) error {
	if s.testconfig.Environment == "preconfigured" {
		return nil
	}

	privateKey, senderAddress, err := s.loadPrivateKey(s.testconfig.KeystoreFile, s.testconfig.KeystorePassword)
	if err != nil {
		return fmt.Errorf("failed to load deployer's private key: %w", err)
	}

	connectionContractName := "centralized-connection"
	contractAddress := senderAddress + "." + connectionContractName

	contract, err := s.client.GetContractById(ctx, contractAddress)
	if err == nil && contract != nil {
		txResp, err := s.client.GetTransactionById(ctx, contract.TxId)
		if err == nil {
			receipt, err := stacksClient.GetReceipt(txResp)
			if err == nil && receipt.Status {
				s.log.Info("Connection contract already successfully deployed, skipping deployment",
					zap.String("address", contractAddress))

				s.IBCAddresses["connection"] = contractAddress
				return s.initializeConnection(ctx, privateKey, senderAddress)
			}
		}
	}

	codeBody, err := os.ReadFile(s.testconfig.Contracts["connection"])
	if err != nil {
		return fmt.Errorf("failed to read connection contract code: %w", err)
	}

	tx, err := transaction.MakeContractDeploy(
		connectionContractName,
		string(codeBody),
		*s.network,
		senderAddress,
		privateKey,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create contract deploy transaction: %w", err)
	}

	txID, err := transaction.BroadcastTransaction(tx, s.network)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	s.log.Info("Deployed centralized-connection contract", zap.String("txID", txID))
	s.IBCAddresses["connection"] = contractAddress

	err = s.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return err
	}

	return s.initializeConnection(ctx, privateKey, senderAddress)
}

func (s *StacksLocalnet) isConnectionInitialized(ctx context.Context, contractAddress string) (bool, error) {
	parts := strings.Split(contractAddress, ".")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid contract ID format: %s", contractAddress)
	}
	contractAddress = parts[0]
	contractName := parts[1]

	result, err := s.client.CallReadOnlyFunction(
		ctx,
		contractAddress,
		contractName,
		"get-xcall",
		[]string{},
	)
	if err != nil {
		return false, fmt.Errorf("failed to check connection initialization: %w", err)
	}

	byteResult, err := hex.DecodeString(strings.TrimPrefix(*result, "0x"))
	if err != nil {
		return false, fmt.Errorf("failed to hex decode get-xcall response: %w", err)
	}

	clarityValue, err := clarity.DeserializeClarityValue(byteResult)
	if err != nil {
		return false, fmt.Errorf("failed to deserialize get-xcall response: %w", err)
	}

	responseValue, ok := clarityValue.(*clarity.ResponseOk)
	if !ok {
		return false, fmt.Errorf("unexpected response type: %T", clarityValue)
	}

	_, ok = responseValue.Value.(*clarity.OptionNone)
	if ok {
		return false, nil
	}

	return true, nil
}

func (s *StacksLocalnet) initializeConnection(ctx context.Context, privateKey []byte, senderAddress string) error {
	contractAddress := s.IBCAddresses["connection"]
	initialized, err := s.isConnectionInitialized(ctx, contractAddress)
	if err != nil {
		return fmt.Errorf("failed to check connection initialization status: %w", err)
	}

	if initialized {
		s.log.Info("Connection contract already initialized, skipping initialization")
		return nil
	}

	xcallAddress := s.IBCAddresses["xcall-proxy"]
	relayerAddress := s.testconfig.RelayWalletAddress

	xcallPrincipal, err := clarity.StringToPrincipal(xcallAddress)
	if err != nil {
		return fmt.Errorf("invalid xcall address: %w", err)
	}
	relayerPrincipal, err := clarity.StringToPrincipal(relayerAddress)
	if err != nil {
		return fmt.Errorf("invalid relayer address: %w", err)
	}

	args := []clarity.ClarityValue{xcallPrincipal, relayerPrincipal}

	txCall, err := transaction.MakeContractCall(
		senderAddress,
		"centralized-connection",
		"initialize",
		args,
		*s.network,
		senderAddress,
		privateKey,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create contract call transaction: %w", err)
	}

	txID, err := transaction.BroadcastTransaction(txCall, s.network)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	s.log.Info("Initialized centralized-connection contract", zap.String("txID", txID))

	err = s.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return err
	}

	return nil
}

func shortenContractName(testcase string) string {
	// Stacks has a contract name limit defined in SIP-003
	maxLength := 30
	prefix := "x-dapp-"

	cleaned := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' {
			return r
		}
		return '-'
	}, strings.ToLower(testcase))

	totalLen := len(prefix) + len(cleaned)
	if totalLen > maxLength {
		remaining := maxLength - len(prefix)
		if remaining > 0 {
			cleaned = cleaned[:remaining]
		} else {
			cleaned = ""
		}
	}
	return prefix + cleaned
}

func (s *StacksLocalnet) DeployXCallMockApp(ctx context.Context, keyName string, connections []chains.XCallConnection) error {
	if s.testconfig.Environment == "preconfigured" {
		return nil
	}

	testcase := ctx.Value("testcase").(string)
	appContractName := shortenContractName(testcase)

	privateKey, senderAddress, err := s.loadPrivateKey(s.testconfig.KeystoreFile, s.testconfig.KeystorePassword)
	if err != nil {
		return fmt.Errorf("failed to load deployer's private key: %w", err)
	}

	contractAddress := senderAddress + "." + appContractName

	contract, err := s.client.GetContractById(ctx, contractAddress)
	if err == nil && contract != nil {
		txResp, err := s.client.GetTransactionById(ctx, contract.TxId)
		if err == nil {
			receipt, err := stacksClient.GetReceipt(txResp)
			if err == nil && receipt.Status {
				s.log.Info("XCall mock app contract already successfully deployed, skipping deployment",
					zap.String("address", contractAddress))

				s.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = contractAddress
				return nil
			}
		}
	}

	codeBody, err := os.ReadFile(s.testconfig.Contracts["dapp"])
	if err != nil {
		return fmt.Errorf("failed to read dapp contract code: %w", err)
	}

	tx, err := transaction.MakeContractDeploy(
		appContractName,
		string(codeBody),
		*s.network,
		senderAddress,
		privateKey,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create contract deploy transaction: %w", err)
	}

	txID, err := transaction.BroadcastTransaction(tx, s.network)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	s.log.Info("Deployed xcall mock app contract",
		zap.String("txID", txID),
		zap.String("contractName", appContractName))

	s.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = contractAddress

	err = s.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return err
	}

	xcallAddress := s.IBCAddresses["xcall-proxy"]
	xcallPrincipal, err := clarity.StringToPrincipal(xcallAddress)
	if err != nil {
		return fmt.Errorf("invalid xcall address: %w", err)
	}

	args := []clarity.ClarityValue{xcallPrincipal}

	txCall, err := transaction.MakeContractCall(
		senderAddress,
		appContractName,
		"initialize",
		args,
		*s.network,
		senderAddress,
		privateKey,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create contract call transaction: %w", err)
	}

	txID, err = transaction.BroadcastTransaction(txCall, s.network)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	s.log.Info("Initialized xcall mock app contract", zap.String("txID", txID))

	err = s.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return err
	}

	for _, connection := range connections {
		s.log.Debug("Setting up connection",
			zap.String("nid", connection.Nid),
			zap.String("destination", connection.Destination))
		err := s.addConnection(ctx, senderAddress, appContractName, privateKey, connection)
		if err != nil {
			s.log.Error("Failed to add connection",
				zap.Error(err),
				zap.String("nid", connection.Nid))
			continue
		}
	}

	return nil
}

func (s *StacksLocalnet) isConnectionExists(ctx context.Context, contractAddress, contractName, nid, source, destination string) (bool, error) {
	sourcesResult, err := s.client.CallReadOnlyFunction(
		ctx,
		contractAddress,
		contractName,
		"get-sources",
		[]string{
			fmt.Sprintf("0x%s", hex.EncodeToString([]byte(nid))),
		},
	)
	if err != nil {
		return false, fmt.Errorf("failed to check sources: %w", err)
	}

	destsResult, err := s.client.CallReadOnlyFunction(
		ctx,
		contractAddress,
		contractName,
		"get-destinations",
		[]string{
			fmt.Sprintf("0x%s", hex.EncodeToString([]byte(nid))),
		},
	)
	if err != nil {
		return false, fmt.Errorf("failed to check destinations: %w", err)
	}

	sourceBytes, err := hex.DecodeString(strings.TrimPrefix(*sourcesResult, "0x"))
	if err != nil {
		return false, fmt.Errorf("failed to decode sources response: %w", err)
	}

	sourcesValue, err := clarity.DeserializeClarityValue(sourceBytes)
	if err != nil {
		return false, fmt.Errorf("failed to deserialize sources: %w", err)
	}

	destBytes, err := hex.DecodeString(strings.TrimPrefix(*destsResult, "0x"))
	if err != nil {
		return false, fmt.Errorf("failed to decode destinations response: %w", err)
	}

	destsValue, err := clarity.DeserializeClarityValue(destBytes)
	if err != nil {
		return false, fmt.Errorf("failed to deserialize destinations: %w", err)
	}

	sourcesList, ok := sourcesValue.(*clarity.List)
	if !ok {
		return false, fmt.Errorf("unexpected sources type: %T", sourcesValue)
	}

	destsList, ok := destsValue.(*clarity.List)
	if !ok {
		return false, fmt.Errorf("unexpected destinations type: %T", destsValue)
	}

	sourceExists := false
	for _, item := range sourcesList.Values {
		if stringVal, ok := item.(*clarity.StringASCII); ok {
			if stringVal.Data == source {
				sourceExists = true
				break
			}
		}
	}

	destExists := false
	for _, item := range destsList.Values {
		if stringVal, ok := item.(*clarity.StringASCII); ok {
			if stringVal.Data == destination {
				destExists = true
				break
			}
		}
	}

	s.log.Debug("Connection check result",
		zap.String("nid", nid),
		zap.String("source", source),
		zap.String("destination", destination),
		zap.Bool("sourceExists", sourceExists),
		zap.Bool("destExists", destExists))

	return sourceExists && destExists, nil
}

func (s *StacksLocalnet) addConnection(ctx context.Context, senderAddress, appContractName string, privateKey []byte, connection chains.XCallConnection) error {
	connAddress := s.IBCAddresses[connection.Connection]

	connArg, err := clarity.NewStringASCII(connAddress)
	if err != nil {
		return fmt.Errorf("failed to create connection argument: %w", err)
	}

	destArg, err := clarity.NewStringASCII(connection.Destination)
	if err != nil {
		return fmt.Errorf("failed to create destination argument: %w", err)
	}

	nidArg, err := clarity.NewStringASCII(connection.Nid)
	if err != nil {
		return fmt.Errorf("failed to create nid argument: %w", err)
	}

	args := []clarity.ClarityValue{
		nidArg,
		connArg,
		destArg,
	}

	parts := strings.Split(connAddress, ".")
	if len(parts) != 2 {
		return fmt.Errorf("invalid connection address format: %s", connAddress)
	}
	contractName := parts[1]

	exists, err := s.isConnectionExists(ctx, senderAddress, appContractName,
		connection.Nid, contractName, connection.Destination)
	if err != nil {
		s.log.Warn("Failed to check existing connection",
			zap.Error(err),
			zap.String("nid", connection.Nid))
	}

	if exists {
		s.log.Info("Connection already exists, skipping",
			zap.String("nid", connection.Nid),
			zap.String("source", contractName),
			zap.String("destination", connection.Destination))
		return nil
	}

	txCall, err := transaction.MakeContractCall(
		senderAddress,
		appContractName,
		"add-connection",
		args,
		*s.network,
		senderAddress,
		privateKey,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create contract call transaction: %w", err)
	}

	txID, err := transaction.BroadcastTransaction(txCall, s.network)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	s.log.Info("Adding new connection to xcall mock app",
		zap.String("txID", txID),
		zap.String("nid", connection.Nid),
		zap.String("source", contractName),
		zap.String("destination", connection.Destination))

	err = s.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return fmt.Errorf("failed to confirm transaction: %w", err)
	}

	s.log.Info("Successfully added new connection",
		zap.String("nid", connection.Nid),
		zap.String("source", contractName),
		zap.String("destination", connection.Destination))

	return nil
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
	implResult, err := s.client.CallReadOnlyFunction(
		ctx,
		"ST1FTM84RHDZ3AB21MNYKA3XQEFB090HZBB81DSFE",
		"xcall-proxy",
		"get-current-implementation",
		[]string{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get current implementation: %w", err)
	}
	s.log.Debug("Current implementation", zap.String("result", *implResult))

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
	contractAddr := parts[0]

	network := strings.Split(_to, "/")[0]

	s.log.Debug("Checking network sources/destinations",
		zap.String("network", network))

	networkArg, err := clarity.NewStringASCII(network)
	if err != nil {
		return nil, fmt.Errorf("failed to create network argument: %w", err)
	}

	networkBytes, err := networkArg.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize network argument: %w", err)
	}

	hexEncodedNetwork := hex.EncodeToString(networkBytes)

	sourcesResult, err := s.client.CallReadOnlyFunction(
		ctx,
		contractAddr,
		contractName,
		"get-sources",
		[]string{hexEncodedNetwork},
	)
	if err != nil {
		s.log.Error("Failed to get sources",
			zap.Error(err),
			zap.String("contractAddr", contractAddr),
			zap.String("contractName", contractName),
			zap.String("network", network),
			zap.String("networkArg", hexEncodedNetwork))
	} else {
		sourceBytes, err := hex.DecodeString(strings.TrimPrefix(*sourcesResult, "0x"))
		if err != nil {
			s.log.Error("Failed to decode sources result", zap.Error(err))
		} else {
			sourcesValue, err := clarity.DeserializeClarityValue(sourceBytes)
			if err != nil {
				s.log.Error("Failed to deserialize sources", zap.Error(err))
			} else {
				s.log.Debug("Sources for network",
					zap.String("network", network),
					zap.Any("sources", sourcesValue))
			}
		}
	}

	destsResult, err := s.client.CallReadOnlyFunction(
		ctx,
		contractAddr,
		contractName,
		"get-destinations",
		[]string{hexEncodedNetwork},
	)
	if err != nil {
		s.log.Error("Failed to get destinations",
			zap.Error(err),
			zap.String("contractAddr", contractAddr),
			zap.String("contractName", contractName),
			zap.String("network", network),
			zap.String("networkArg", hexEncodedNetwork))
	} else {
		destBytes, err := hex.DecodeString(strings.TrimPrefix(*destsResult, "0x"))
		if err != nil {
			s.log.Error("Failed to decode destinations result", zap.Error(err))
		} else {
			destsValue, err := clarity.DeserializeClarityValue(destBytes)
			if err != nil {
				s.log.Error("Failed to deserialize destinations", zap.Error(err))
			} else {
				s.log.Debug("Destinations for network",
					zap.String("network", network),
					zap.Any("destinations", destsValue))
			}
		}
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

	ctx = context.WithValue(ctx, snContextKey, sn)
	return ctx, nil
}

func (s *StacksLocalnet) extractSnFromTransaction(ctx context.Context, txID string) (string, error) {
	tx, err := s.client.GetTransactionById(ctx, txID)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction by ID: %w", err)
	}

	if confirmed := tx.GetTransactionList200ResponseResultsInner; confirmed != nil {
		if contractCall := confirmed.ContractCallTransaction; contractCall != nil {
			for _, event := range contractCall.Events {
				if event.SmartContractLogTransactionEvent != nil {
					contractLog := event.SmartContractLogTransactionEvent.ContractLog
					if contractLog.Topic == "print" {
						repr := contractLog.Value.Repr
						if strings.Contains(repr, "CallMessageSent") {
							startIdx := strings.Index(repr, "(sn u")
							if startIdx != -1 {
								startIdx += 5 // Move past "(sn u"
								endIdx := strings.Index(repr[startIdx:], ")")
								if endIdx != -1 {
									sn := repr[startIdx : startIdx+endIdx]
									return sn, nil
								}
							}
						}
					}
				}
			}
		}
	}

	return "", fmt.Errorf("serial number 'sn' not found in transaction events")
}

func (s *StacksLocalnet) FindCallMessage(ctx context.Context, startHeight uint64, from, to, sn string) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	resultChan := make(chan struct {
		txID string
		data string
	}, 1)

	wsURL := s.client.GetWebSocketURL()
	eventSystem := events.NewEventSystem(ctx, wsURL, s.log)

	eventSystem.OnEvent(func(event *events.Event) error {
		if event.Type != stacksClient.CallMessage {
			return nil
		}

		if data, ok := event.Data.(events.CallMessageData); ok {
			if data.Sn == sn {
				select {
				case resultChan <- struct {
					txID string
					data string
				}{
					txID: data.ReqID,
					data: data.Data,
				}:
				case <-ctx.Done():
				}
			}
		}
		return nil
	})

	if err := eventSystem.Start(); err != nil {
		return "", "", fmt.Errorf("failed to start event system: %w", err)
	}
	defer eventSystem.Stop()

	select {
	case result := <-resultChan:
		return result.txID, result.data, nil
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return "", "", fmt.Errorf("find call message timed out")
		}
		return "", "", ctx.Err()
	}
}

func (s *StacksLocalnet) FindCallResponse(ctx context.Context, startHeight uint64, sn string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	resultChan := make(chan string, 1)

	wsURL := s.client.GetWebSocketURL()
	eventSystem := events.NewEventSystem(ctx, wsURL, s.log)

	eventSystem.OnEvent(func(event *events.Event) error {
		if event.Type != stacksClient.CallMessage {
			return nil
		}

		if data, ok := event.Data.(events.CallMessageData); ok {
			if data.Sn == sn {
				select {
				case resultChan <- data.ReqID:
				case <-ctx.Done():
				}
			}
		}
		return nil
	})

	if err := eventSystem.Start(); err != nil {
		return "", fmt.Errorf("failed to start event system: %w", err)
	}
	defer eventSystem.Stop()

	select {
	case txID := <-resultChan:
		return txID, nil
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("find call response timed out")
		}
		return "", ctx.Err()
	}
}

func (s *StacksLocalnet) FindRollbackExecutedMessage(ctx context.Context, startHeight uint64, sn string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	resultChan := make(chan string, 1)

	wsURL := s.client.GetWebSocketURL()
	eventSystem := events.NewEventSystem(ctx, wsURL, s.log)

	eventSystem.OnEvent(func(event *events.Event) error {
		if event.Type != stacksClient.RollbackMessage {
			return nil
		}

		if data, ok := event.Data.(events.RollbackMessageData); ok {
			if data.Sn == sn {
				select {
				case resultChan <- data.Sn:
				case <-ctx.Done():
				}
			}
		}
		return nil
	})

	if err := eventSystem.Start(); err != nil {
		return "", fmt.Errorf("failed to start event system: %w", err)
	}
	defer eventSystem.Stop()

	select {
	case txID := <-resultChan:
		return txID, nil
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("find rollback message timed out")
		}
		return "", ctx.Err()
	}
}

func (s *StacksLocalnet) FindTargetXCallMessage(ctx context.Context, target chains.Chain, height uint64, to string) (*chains.XCallResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	resultChan := make(chan *chains.XCallResponse, 1)

	wsURL := s.client.GetWebSocketURL()
	eventSystem := events.NewEventSystem(ctx, wsURL, s.log)

	eventSystem.OnEvent(func(event *events.Event) error {
		var response *chains.XCallResponse

		switch event.Type {
		case stacksClient.EmitMessage:
			if data, ok := event.Data.(events.EmitMessageData); ok {
				if data.TargetNetwork == to {
					response = &chains.XCallResponse{
						SerialNo: data.Sn,
						Data:     data.Msg,
					}
				}
			}

		case stacksClient.CallMessage:
			if data, ok := event.Data.(events.CallMessageData); ok {
				response = &chains.XCallResponse{
					SerialNo:  data.Sn,
					RequestID: data.ReqID,
					Data:      data.Data,
				}
			}
		}

		if response != nil {
			select {
			case resultChan <- response:
			case <-ctx.Done():
			}
		}
		return nil
	})

	if err := eventSystem.Start(); err != nil {
		return nil, fmt.Errorf("failed to start event system: %w", err)
	}
	defer eventSystem.Stop()

	select {
	case response := <-resultChan:
		return response, nil
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("find target message timed out")
		}
		return nil, ctx.Err()
	}
}

func (s *StacksLocalnet) loadPrivateKey(keystoreFile, password string) ([]byte, string, error) {
	if s.kms == nil {
		return nil, "", fmt.Errorf("KMS not initialized")
	}

	encryptedKey, err := os.ReadFile(keystoreFile)
	if err != nil {
		s.log.Error("Failed to read keystore file",
			zap.String("path", keystoreFile),
			zap.Error(err))
		return nil, "", fmt.Errorf("failed to read keystore file: %w", err)
	}

	privateKey, err := s.kms.Decrypt(context.Background(), encryptedKey)
	if err != nil {
		s.log.Error("Failed to decrypt keystore", zap.Error(err))
		return nil, "", fmt.Errorf("failed to decrypt keystore: %w", err)
	}

	network, err := stacksClient.MapNIDToChainID(s.cfg.ChainID)
	if err != nil {
		s.log.Error("Chain id not found. Check the NID config", zap.Error(err))
		return nil, "", fmt.Errorf("chain id not found: %w", err)
	}

	address, err := crypto.GetAddressFromPrivateKey(privateKey, network)
	if err != nil {
		s.log.Error("Failed to derive address from private key", zap.Error(err))
		return nil, "", fmt.Errorf("failed to derive address: %w", err)
	}

	return privateKey, address, nil
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

func (s *StacksLocalnet) waitForTransactionConfirmation(ctx context.Context, txID string) error {
	timeout := time.After(MAX_WAIT_TIME)
	ticker := time.NewTicker(BLOCK_TIME)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			res, err := s.client.GetTransactionById(ctx, txID)
			if err != nil {
				s.log.Warn("Failed to query transaction receipt", zap.Error(err))
				continue
			}

			receipt, err := stacksClient.GetReceipt(res)
			if err != nil {
				s.log.Warn("Failed to extract transaction receipt", zap.Error(err))
				continue
			}

			if receipt.Status {
				s.log.Info("Transaction confirmed",
					zap.String("txID", txID),
					zap.Uint64("height", receipt.Height))
				return nil
			}
			s.log.Debug("Transaction not yet confirmed", zap.String("txID", txID))

		case <-timeout:
			return fmt.Errorf("transaction confirmation timed out after %v seconds", MAX_WAIT_TIME)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
