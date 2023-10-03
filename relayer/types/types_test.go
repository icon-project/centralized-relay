package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageCache(t *testing.T) {

	messageCache := NewMessageCache()

	m1 := Message{
		Dst: "mock-2",
		Src: "mock-1",
		Sn:  1,
	}

	t.Run("message cache add", func(t *testing.T) {
		messageCache.Add(NewRouteMessage(m1))
		assert.Equal(t, messageCache.Len(), uint64(1))
	})

	t.Run("message removed from cache", func(t *testing.T) {
		messageCache.Remove(m1.MessageKey())
		assert.Equal(t, messageCache.Len(), uint64(0))
	})
}
