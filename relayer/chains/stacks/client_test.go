package stacks_test

import (
	"context"
	"encoding/hex"
	"path/filepath"
	"strings"
	"testing"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks"
	"github.com/icon-project/stacks-go-sdk/pkg/clarity"
	"github.com/icon-project/stacks-go-sdk/pkg/crypto"
	stacksSdk "github.com/icon-project/stacks-go-sdk/pkg/stacks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestClient_GetAccountBalance(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	network := stacksSdk.NewStacksTestnet()
	xcallAbiPath := filepath.Join("abi", "xcall-proxy-abi.json")
	client, err := stacks.NewClient(logger, network, xcallAbiPath)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	address := "ST15C893XJFJ6FSKM020P9JQDB5T7X6MQTXMBPAVH"
	balance, err := client.GetAccountBalance(ctx, address)
	if err != nil {
		t.Fatalf("Failed to get account balance: %v", err)
	}

	t.Logf("Balance for address %s: %s", address, balance.String())
}

func TestClient_GetAccountNonce(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	network := stacksSdk.NewStacksTestnet()
	xcallAbiPath := filepath.Join("abi", "xcall-proxy-abi.json")
	client, err := stacks.NewClient(logger, network, xcallAbiPath)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	address := "ST15C893XJFJ6FSKM020P9JQDB5T7X6MQTXMBPAVH"
	nonce, err := client.GetAccountNonce(ctx, address)
	if err != nil {
		t.Fatalf("Failed to get account nonce: %v", err)
	}

	t.Logf("Nonce for address %s: %d", address, nonce)
}

func TestClient_GetBlockByHeightOrHash(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	network := stacksSdk.NewStacksTestnet()
	xcallAbiPath := filepath.Join("abi", "xcall-proxy-abi.json")
	client, err := stacks.NewClient(logger, network, xcallAbiPath)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	block, err := client.GetLatestBlock(ctx)
	if err != nil {
		t.Fatalf("Failed to get latest blocks: %v", err)
	}
	if block == nil {
		t.Fatalf("No blocks found")
	}

	blockHeight := block.Height

	block, err = client.GetBlockByHeightOrHash(ctx, uint64(blockHeight))
	if err != nil {
		t.Fatalf("Failed to get block by height: %v", err)
	}

	t.Logf("Block at height %d: %+v", blockHeight, block)
}

func TestClient_GetLatestBlock(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	network := stacksSdk.NewStacksTestnet()
	xcallAbiPath := filepath.Join("abi", "xcall-proxy-abi.json")
	client, err := stacks.NewClient(logger, network, xcallAbiPath)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	block, err := client.GetLatestBlock(ctx)
	if err != nil {
		t.Fatalf("Failed to get latest blocks: %v", err)
	}

	t.Logf("Latest block: %+v", block)
}

func TestClient_CallReadOnlyFunction(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	network := stacksSdk.NewStacksTestnet()
	xcallAbiPath := filepath.Join("abi", "xcall-proxy-abi.json")
	client, err := stacks.NewClient(logger, network, xcallAbiPath)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	contractAddress := "ST15C893XJFJ6FSKM020P9JQDB5T7X6MQTXMBPAVH"
	contractName := "xcall-proxy"
	functionName := "get-current-implementation"

	functionArgs := []string{}

	result, err := client.CallReadOnlyFunction(ctx, contractAddress, contractName, functionName, functionArgs)
	if err != nil {
		t.Fatalf("Failed to call read-only function: %v", err)
	}

	t.Logf("Result of calling %s::%s: %s", contractName, functionName, *result)

	decodedResult, err := hex.DecodeString(strings.TrimPrefix(*result, "0x"))
	assert.NoError(t, err, "Failed to decode hex string")

	cv, err := clarity.DeserializeClarityValue(decodedResult)
	assert.NoError(t, err, "Failed to deserialize clarity value")
	resp, ok := cv.(*clarity.ResponseOk)
	assert.True(t, ok, "Expected result to be ResponseOk")

	principalType := resp.Value.Type()
	assert.Equal(t, principalType, clarity.ClarityTypeStandardPrincipal)
}

func TestClient_GetCurrentImplementation(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	network := stacksSdk.NewStacksTestnet()
	xcallAbiPath := filepath.Join("abi", "xcall-proxy-abi.json")
	client, err := stacks.NewClient(logger, network, xcallAbiPath)
	assert.NoError(t, err, "Failed to create client")

	contractAddress := "ST15C893XJFJ6FSKM020P9JQDB5T7X6MQTXMBPAVH"

	impl, err := client.GetCurrentImplementation(ctx, contractAddress)
	assert.NoError(t, err, "Failed to get current implementation")
	assert.NotEmpty(t, impl, "Implementation address should not be empty")

	t.Logf("Current implementation: %s", impl)
}

func TestClient_SetAdmin(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	network := stacksSdk.NewStacksTestnet()
	xcallAbiPath := filepath.Join("abi", "xcall-proxy-abi.json")
	client, err := stacks.NewClient(logger, network, xcallAbiPath)
	assert.NoError(t, err, "Failed to create client")

	contractAddress := "ST15C893XJFJ6FSKM020P9JQDB5T7X6MQTXMBPAVH.xcall-proxy"
	newAdmin := "ST15C893XJFJ6FSKM020P9JQDB5T7X6MQTXMBPAVH"

	currentImplementation, _ := client.GetCurrentImplementation(ctx, contractAddress)
	senderAddress := "ST15C893XJFJ6FSKM020P9JQDB5T7X6MQTXMBPAVH"
	mnemonic := "vapor unhappy gather snap project ball gain puzzle comic error avocado bounce letter anxiety wheel provide canyon promote sniff improve figure daughter mansion baby"
	senderKey, err := crypto.DeriveStxPrivateKey(mnemonic, 0)
	if err != nil {
		t.Fatalf("Failed to derive sender key: %v", err)
	}

	txID, err := client.SetAdmin(ctx, contractAddress, newAdmin, currentImplementation, senderAddress, senderKey)
	assert.NoError(t, err, "Failed to set admin")
	assert.NotEmpty(t, txID, "Transaction ID should not be empty")

	t.Logf("SetAdmin transaction ID: %s", txID)
}

// func TestClient_SubscribeToEvents(t *testing.T) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	logger, _ := zap.NewDevelopment()
// 	client, err := stacks.NewClient("https://api.testnet.hiro.so", logger)
// 	if err != nil {
// 		t.Fatalf("Failed to create client: %v", err)
// 	}

// 	var wg sync.WaitGroup
// 	wg.Add(1)

// 	callback := func(eventType string, data interface{}) error {
// 		t.Logf("Received event: %s, Data: %+v", eventType, data)
// 		wg.Done()
// 		return nil
// 	}

// 	err = client.SubscribeToEvents(ctx, []string{"block"}, callback)
// 	if err != nil {
// 		t.Fatalf("Failed to subscribe to events: %v", err)
// 	}

// 	wg.Wait()
// }
