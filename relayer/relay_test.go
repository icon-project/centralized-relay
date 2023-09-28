package relayer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/icon-project/centralized-relay/relayer/chains/mockchain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRelayer(t *testing.T) {

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
	errorchan := Start(ctx, logger, chains, 3*time.Second, true)

	for {
		select {
		case err := <-errorchan:
			fmt.Println("error occured: ", err)
			break
		}
	}

}
