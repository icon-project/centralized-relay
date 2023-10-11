package icon

import (
	"context"

	"github.com/icon-project/centralized-relay/relayer/types"
)

// starting listener
func (icp *IconProvider) Listener(ctx context.Context, lastSavedHeight uint64, incoming chan types.BlockInfo) error {
	return nil
}
