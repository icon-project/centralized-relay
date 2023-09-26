package icon

import (
	"context"

	"github.com/icon-project/centralized-relay/relayer/provider"
)

// starting listener
func (icp *IconProvider) Listener(ctx context.Context, incoming <-chan provider.BlockInfo) error {
	return nil
}
