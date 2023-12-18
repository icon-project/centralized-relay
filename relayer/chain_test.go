package relayer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestChain(t *testing.T) {

	logger := zap.NewNop()

	nid := "mock-1"

	mockProvider, err := GetMockChainProvider(logger, 1*time.Second, nid, "mock-2", 10, 20)
	assert.NoError(t, err)

	chain := NewChain(logger, mockProvider, true)
	assert.Equal(t, chain.NID(), nid)
}
