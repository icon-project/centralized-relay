package stacks

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	stacksClient "github.com/icon-project/centralized-relay/relayer/chains/stacks"
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
	network := stacks.NewStacksLocalnet()
	client, err := stacksClient.NewClient(log, network, testconfig.Contracts["xcall"])
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
	contracts["xcall"] = s.GetContractAddress("xcall")
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
		s.IBCAddresses["xcall"] = "STXCALLPROXYADDRESS"
		s.IBCAddresses["connection"] = "STXCONNECTIONADDRESS"
		s.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = "STXDAPPADDRESS"
		return nil
	}

	privateKey, senderAddress, err := s.loadPrivateKey(s.testconfig.KeystoreFile, s.testconfig.KeystorePassword)
	if err != nil {
		return fmt.Errorf("failed to load deployer's private key: %w", err)
	}

	xcallContractName := "xcall-proxy"
	codeBody, err := os.ReadFile(s.testconfig.Contracts["xcall"])
	if err != nil {
		return fmt.Errorf("failed to read xcall contract code: %w", err)
	}

	tx, err := transaction.MakeContractDeploy(
		xcallContractName,
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

	s.log.Info("Deployed xcall-proxy contract", zap.String("txID", txID))

	contractAddress := senderAddress + "." + xcallContractName
	s.IBCAddresses["xcall"] = contractAddress

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

	contractAddress := senderAddress + "." + connectionContractName
	s.IBCAddresses["connection"] = contractAddress

	err = s.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return err
	}

	xcallAddress := s.IBCAddresses["xcall"]
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
		connectionContractName,
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

	s.log.Info("Initialized centralized-connection contract", zap.String("txID", txID))

	err = s.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return err
	}

	return nil
}

func (s *StacksLocalnet) DeployXCallMockApp(ctx context.Context, keyName string, connections []chains.XCallConnection) error {
	if s.testconfig.Environment == "preconfigured" {
		return nil
	}

	testcase := ctx.Value("testcase").(string)
	appContractName := "xcall-mock-app-" + testcase

	codeBody, err := os.ReadFile(s.testconfig.Contracts["dapp"])
	if err != nil {
		return fmt.Errorf("failed to read dapp contract code: %w", err)
	}

	privateKey, senderAddress, err := s.loadPrivateKey(s.testconfig.KeystoreFile, s.testconfig.KeystorePassword)
	if err != nil {
		return fmt.Errorf("failed to load deployer's private key: %w", err)
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

	s.log.Info("Deployed xcall mock app contract", zap.String("txID", txID))

	contractAddress := senderAddress + "." + appContractName
	s.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = contractAddress

	err = s.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return err
	}

	xcallAddress := s.IBCAddresses["xcall"]
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
		connPrincipal, err := clarity.StringToPrincipal(s.IBCAddresses[connection.Connection])
		if err != nil {
			s.log.Error("Invalid connection address", zap.Error(err), zap.String("address", s.IBCAddresses[connection.Connection]))
			continue
		}

		nidArg, err := clarity.NewStringASCII(connection.Nid)
		if err != nil {
			s.log.Error("Failed to create nid argument", zap.Error(err))
			continue
		}

		destArg, err := clarity.NewStringASCII(connection.Destination)
		if err != nil {
			s.log.Error("Failed to create destination argument", zap.Error(err))
			continue
		}

		args := []clarity.ClarityValue{
			nidArg,
			connPrincipal,
			destArg,
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
			s.log.Error("Failed to create contract call transaction", zap.Error(err))
			continue
		}

		txID, err = transaction.BroadcastTransaction(txCall, s.network)
		if err != nil {
			s.log.Error("Failed to broadcast transaction", zap.Error(err))
			continue
		}

		s.log.Info("Added connection to xcall mock app", zap.String("txID", txID))

		err = s.waitForTransactionConfirmation(ctx, txID)
		if err != nil {
			s.log.Error("Failed to confirm transaction", zap.Error(err))
			continue
		}
	}

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
	rollbackArg := clarity.NewBuffer(rollback)

	args := []clarity.ClarityValue{
		toArg,
		dataArg,
		rollbackArg,
	}

	txCall, err := transaction.MakeContractCall(
		senderAddress,
		dappAddress,
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

	// TODO: Extract 'sn' (serial number) from transaction events
	// For now, we return the txID in context
	ctx = context.WithValue(ctx, "txID", txID)
	return ctx, nil
}

func (s *StacksLocalnet) FindCallMessage(ctx context.Context, startHeight uint64, from, to, sn string) (string, string, error) {
	foundChan := make(chan struct {
		txID string
		data string
	}, 1)

	callback := func(eventType string, data interface{}) error {
		if eventType == stacksClient.CallMessage {
			if callMsg, ok := data.(stacksClient.CallMessageEvent); ok {
				if callMsg.Sn == sn {
					foundChan <- struct {
						txID string
						data string
					}{
						txID: callMsg.ReqID,
						data: callMsg.Data,
					}
				}
			}
		}
		return nil
	}

	err := s.client.SubscribeToEvents(ctx, []string{stacksClient.CallMessage}, callback)
	if err != nil {
		return "", "", fmt.Errorf("failed to subscribe to events: %w", err)
	}

	select {
	case found := <-foundChan:
		return found.txID, found.data, nil
	case <-time.After(2 * time.Minute):
		return "", "", fmt.Errorf("find call message timed out")
	case <-ctx.Done():
		return "", "", ctx.Err()
	}
}

func (s *StacksLocalnet) FindCallResponse(ctx context.Context, startHeight uint64, sn string) (string, error) {
	foundChan := make(chan string, 1)

	callback := func(eventType string, data interface{}) error {
		if eventType == stacksClient.CallMessage {
			if callMsg, ok := data.(stacksClient.CallMessageEvent); ok {
				if callMsg.Sn == sn {
					foundChan <- callMsg.ReqID
				}
			}
		}
		return nil
	}

	err := s.client.SubscribeToEvents(ctx, []string{stacksClient.CallMessage}, callback)
	if err != nil {
		return "", fmt.Errorf("failed to subscribe to events: %w", err)
	}

	select {
	case txID := <-foundChan:
		return txID, nil
	case <-time.After(2 * time.Minute):
		return "", fmt.Errorf("find call response timed out")
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (s *StacksLocalnet) FindRollbackExecutedMessage(ctx context.Context, startHeight uint64, sn string) (string, error) {
	foundChan := make(chan string, 1)

	callback := func(eventType string, data interface{}) error {
		if eventType == stacksClient.RollbackMessage {
			if rollbackMsg, ok := data.(stacksClient.RollbackMessageEvent); ok {
				if rollbackMsg.Sn == sn {
					foundChan <- rollbackMsg.Sn
				}
			}
		}
		return nil
	}

	err := s.client.SubscribeToEvents(ctx, []string{stacksClient.RollbackMessage}, callback)
	if err != nil {
		return "", fmt.Errorf("failed to subscribe to events: %w", err)
	}

	select {
	case txID := <-foundChan:
		return txID, nil
	case <-time.After(2 * time.Minute):
		return "", fmt.Errorf("find rollback message timed out")
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (s *StacksLocalnet) FindTargetXCallMessage(ctx context.Context, target chains.Chain, height uint64, to string) (*chains.XCallResponse, error) {
	foundChan := make(chan *chains.XCallResponse, 1)

	callback := func(eventType string, data interface{}) error {
		if eventType == stacksClient.EmitMessage {
			if emitMsg, ok := data.(stacksClient.EmitMessageEvent); ok {
				if emitMsg.TargetNetwork == to {
					foundChan <- &chains.XCallResponse{
						SerialNo: emitMsg.Sn,
						Data:     emitMsg.Msg,
						// RequestID isn't in EmitMessageEvent, will be empty
					}
				}
			}
		} else if eventType == stacksClient.CallMessage {
			if callMsg, ok := data.(stacksClient.CallMessageEvent); ok {
				foundChan <- &chains.XCallResponse{
					SerialNo:  callMsg.Sn,
					RequestID: callMsg.ReqID,
					Data:      callMsg.Data,
				}
			}
		}
		return nil
	}

	err := s.client.SubscribeToEvents(ctx, []string{stacksClient.EmitMessage, stacksClient.CallMessage}, callback)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to events: %w", err)
	}

	select {
	case response := <-foundChan:
		return response, nil
	case <-time.After(2 * time.Minute):
		return nil, fmt.Errorf("find target message timed out")
	case <-ctx.Done():
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

	return s.FindTargetXCallMessage(ctx, targetChain, height, strings.Split(_to, "/")[1])
}

func (s *StacksLocalnet) waitForTransactionConfirmation(ctx context.Context, txID string) error {
	timeout := time.After(2 * time.Minute)
	ticker := time.NewTicker(2 * time.Second)
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
			return fmt.Errorf("transaction confirmation timed out after 2 minutes")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
