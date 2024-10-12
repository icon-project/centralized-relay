package stacks_test

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks"
	"github.com/icon-project/stacks-go-sdk/pkg/clarity"
	"go.uber.org/zap"
)

func TestClient_GetAccountBalance(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	client, err := stacks.NewClient("https://stacks-node-api.testnet.stacks.co", logger)
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
	client, err := stacks.NewClient("https://stacks-node-api.testnet.stacks.co", logger)
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
	client, err := stacks.NewClient("https://stacks-node-api.testnet.stacks.co", logger)
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
	client, err := stacks.NewClient("https://stacks-node-api.testnet.stacks.co", logger)
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
	client, err := stacks.NewClient("https://stacks-node-api.testnet.stacks.co", logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	contractAddress := "ST15C893XJFJ6FSKM020P9JQDB5T7X6MQTXMBPAVH"
	contractName := "contract_name"
	functionName := "address-string-to-principal"

	strArg, _ := clarity.NewStringASCII("test")
	encodedStrArg, _ := strArg.Serialize()
	hexEncodedStrArg := hex.EncodeToString(encodedStrArg)

	functionArgs := []string{hexEncodedStrArg}

	result, err := client.CallReadOnlyFunction(ctx, contractAddress, contractName, functionName, functionArgs)
	if err != nil {
		t.Fatalf("Failed to call read-only function: %v", err)
	}

	t.Logf("Result of calling %s::%s: %s", contractName, functionName, *result)
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
