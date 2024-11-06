package stacks_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks"
	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/provider"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupTestStacksProvider(t *testing.T) *stacks.Provider {
	logger, _ := zap.NewDevelopment()
	cfg := &stacks.Config{
		CommonConfig: provider.CommonConfig{
			RPCUrl: "https://stacks-node-api.testnet.stacks.co",
			Contracts: providerTypes.ContractConfigMap{
				"XcallContract":      "ST15C893XJFJ6FSKM020P9JQDB5T7X6MQTXMBPAVH.xcall-proxy",
				"ConnectionContract": "ST15C893XJFJ6FSKM020P9JQDB5T7X6MQTXMBPAVH.centralized-connection",
			},
			NID: "stacks_testnet",
		},
	}

	p, err := cfg.NewProvider(context.Background(), logger, "/tmp/relayer", false, "stacks_testnet")
	assert.NoError(t, err)
	assert.NotNil(t, p)

	return p.(*stacks.Provider)
}

func TestGenerateMessages(t *testing.T) {
	p := setupTestStacksProvider(t)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Trigger an event on the contract (you'll need to do this manually or automate it)
	// For example, call a function on your XCall contract that emits an event

	key := &providerTypes.MessageKeyWithMessageHeight{
		Height: 12345, // Use an appropriate block height
	}

	messages, err := p.GenerateMessages(ctx, key)
	assert.NoError(t, err)
	assert.NotEmpty(t, messages)

	for _, msg := range messages {
		t.Logf("Generated message: %+v", msg)
		// Add more specific assertions based on the expected event data
	}
}

func TestRoute(t *testing.T) {
	p := setupTestStacksProvider(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	message := &providerTypes.Message{
		Dst:       "stacks_testnet",
		Src:       "icon_testnet",
		Sn:        big.NewInt(12345),
		EventType: events.EmitMessage,
		Data:      []byte("Hello, Stacks!"),
	}

	callback := func(key *providerTypes.MessageKey, response *providerTypes.TxResponse, err error) {
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, providerTypes.Success, response.Code)
	}

	err := p.Route(ctx, message, callback)
	assert.NoError(t, err)
}
