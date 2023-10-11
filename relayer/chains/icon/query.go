package icon

import (
	"context"

	"github.com/icon-project/centralized-relay/relayer/types"
)

// TODO:
func (icp *IconProvider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	return 0, nil
}

func (icp *IconProvider) ShouldReceiveMessage(ctx context.Context, messagekey types.Message) (bool, error) {
	return true, nil

}
func (icp *IconProvider) ShouldSendMessage(ctx context.Context, messageKey types.Message) (bool, error) {
	return true, nil

}
