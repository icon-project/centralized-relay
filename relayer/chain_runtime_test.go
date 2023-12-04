package relayer

import (
	"context"
	"testing"
	"time"

	"github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestChainRuntime(t *testing.T) {
	ctx := context.Background()

	logger := zap.NewNop()

	mockProvider, err := GetMockChainProvider(logger, 1*time.Second, "mock", "mock-2", 10, 20)
	assert.NoError(t, err)

	runtime, err := NewChainRuntime(logger, NewChain(&zap.Logger{}, mockProvider, true))
	assert.NoError(t, err)

	m1 := &types.Message{
		Dst: "mock-2",
		Src: "mock-1",
		Sn:  1,
	}
	m2 := &types.Message{
		Dst: "mock-2",
		Src: "mock-1",
		Sn:  2,
	}
	info := types.BlockInfo{
		Height:   15,
		Messages: []*types.Message{m1, m2},
	}

	t.Run("merge messages", func(t *testing.T) {
		runtime.mergeMessages(ctx, info.Messages)
		assert.Equal(t, len(runtime.MessageCache.Messages), len(info.Messages))
	})

	t.Run("clear messages", func(t *testing.T) {
		runtime.clearMessageFromCache([]*types.MessageKey{m1.MessageKey()})
		assert.Equal(t, len(runtime.MessageCache.Messages), len(info.Messages)-1)
		assert.Equal(t, runtime.MessageCache.Messages[m2.MessageKey()], types.NewRouteMessage(m2))
	})
}
