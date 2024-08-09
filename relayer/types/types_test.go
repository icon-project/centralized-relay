package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteMessage(t *testing.T) {
	m1 := &Message{
		Dst: "mock-2",
		Src: "mock-1",
		Sn:  big.NewInt(1),
	}

	routeMessage := NewRouteMessage(m1)

	t.Run("route message set processing", func(t *testing.T) {
		routeMessage.ToggleProcessing()
		assert.Equal(t, true, routeMessage.IsProcessing())
	})

	t.Run("route message increment retry", func(t *testing.T) {
		routeMessage.IncrementRetry()
		assert.Equal(t, uint8(1), routeMessage.GetRetry())
	})

	t.Run("getMessage from routeMessage", func(t *testing.T) {
		assert.Equal(t, m1, routeMessage.GetMessage())
	})
}

func TestMessageCache(t *testing.T) {
	messageCache := NewMessageCache()

	m1 := &Message{
		Dst: "mock-2",
		Src: "mock-1",
		Sn:  big.NewInt(1),
	}

	t.Run("message cache add", func(t *testing.T) {
		messageCache.Add(NewRouteMessage(m1))
		assert.Equal(t, messageCache.Len(), int(1))
	})

	t.Run("message removed from cache", func(t *testing.T) {
		messageCache.Remove(m1.MessageKey())
		assert.Equal(t, messageCache.Len(), int(0))
	})
}
