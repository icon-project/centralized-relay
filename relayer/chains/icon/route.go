package icon

import (
	"context"

	"github.com/icon-project/centralized-relay/relayer/provider"
)

func (icp *IconProvider) Route(ctx context.Context, message *provider.RouteMessage, callback func(response provider.ExecuteMessageResponse)) error {
	return nil
}
