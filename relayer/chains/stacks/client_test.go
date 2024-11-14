package stacks_test

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks"
	"github.com/icon-project/stacks-go-sdk/pkg/clarity"
	stacksSdk "github.com/icon-project/stacks-go-sdk/pkg/stacks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestClient_GetAccountNonce(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	network := stacksSdk.NewStacksTestnet()
	client, err := stacks.NewClient(logger, network)
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

func TestClient_CallReadOnlyFunction(t *testing.T) {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()
	network := stacksSdk.NewStacksTestnet()
	client, err := stacks.NewClient(logger, network)
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
	client, err := stacks.NewClient(logger, network)
	assert.NoError(t, err, "Failed to create client")

	contractAddress := "ST15C893XJFJ6FSKM020P9JQDB5T7X6MQTXMBPAVH"

	impl, err := client.GetCurrentImplementation(ctx, contractAddress)
	assert.NoError(t, err, "Failed to get current implementation")
	assert.NotEmpty(t, impl, "Implementation address should not be empty")

	t.Logf("Current implementation: %s", impl)
}
