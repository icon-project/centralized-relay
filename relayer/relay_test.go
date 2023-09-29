package relayer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/icon-project/centralized-relay/relayer/chains/mockchain"
	"github.com/icon-project/centralized-relay/relayer/lvldb"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// func TestNewRelayer(t *testing.T) {
// 	logger := zap.NewNop()
// 	// Create a map of dummy chains for testing.
// 	chains := map[string]*Chain{
// 		"chain1": &Chain{debug: false},
// 		"chain2": &Chain{debug: false, log: logger},
// 	}

// 	// Create a new Relayer instance.
// 	relayer := NewRelayer(chains, logger)

// 	// Check if the chains and logger are correctly assigned.
// 	if reflect.DeepEqual(relayer.chains, chains) || relayer.log != logger {
// 		t.Errorf("NewRelayer did not initialize the Relayer struct correctly")
// 	}

// 	// Check if listener channels are created for each chain.
// 	for chainID, ch := range relayer.listenerChans {
// 		if _, ok := chains[chainID]; !ok {
// 			t.Errorf("NewRelayer created an unexpected listener channel for chain %s", chainID)
// 		}
// 		if cap(ch) != listenerChannelBufferSize {
// 			t.Errorf("NewRelayer created a listener channel with the wrong capacity for chain %s", chainID)
// 		}
// 	}
// }

func TestListenerRelayer(t *testing.T) {

	chains := make(map[string]*Chain, 0)

	logger := zap.NewNop()

	// adding mock-1
	mock1ChainId := "mock-1"
	mock2ChainId := "mock-2"
	mock1ProviderConfig := mockchain.MockProviderConfig{
		ChainId:       mock1ChainId,
		BlockDuration: 2 * time.Second,
		TargetChains:  []string{mock2ChainId},
	}
	mock1Provider, err := mock1ProviderConfig.NewProvider(logger, "empty", false, mock1ChainId)
	assert.NoError(t, err)

	chains[mock1ChainId] = NewChain(logger, mock1Provider, true)

	mock2ProviderConfig := mockchain.MockProviderConfig{
		ChainId:       mock2ChainId,
		BlockDuration: 6 * time.Second,
		TargetChains:  []string{mock2ChainId},
	}
	mock2Provider, err := mock2ProviderConfig.NewProvider(logger, "empty", false, mock2ChainId)
	assert.NoError(t, err)

	chains[mock2ChainId] = NewChain(logger, mock2Provider, true)

	ctx := context.Background()

	db, err := lvldb.NewLvlDB("./testdb")

	if err != nil {
		assert.Fail(t, "unable to create database", err)

	}

	errorchan, err := Start(ctx, logger, chains, 3*time.Second, true, db)
	if err != nil {
		assert.Fail(t, "unable to start the relayer ", err)
	}

	for {
		select {
		case err := <-errorchan:
			fmt.Println("error occured: ", err)
			break
		}
	}

}
