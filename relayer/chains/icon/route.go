package icon

import (
	"context"

	"github.com/icon-project/centralized-relay/relayer/types"
)

func (icp *IconProvider) Route(ctx context.Context, message *types.RouteMessage, callback func(response types.ExecuteMessageResponse)) error {
	return nil
}
