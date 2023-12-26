package tui

import (
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/blockdb"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestModel_RootView(t *testing.T) {
	m := NewModel(&mockQueryService{}, "test.db", "abc123", time.Now(), make([]blockdb.TestCaseResult, 1))
	view := m.RootView()
	require.NotNil(t, view)
	require.Greater(t, view.GetItemCount(), 0)
}
