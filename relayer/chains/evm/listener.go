package evm

import (
	"context"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/types"
)

// ListenToEvents goes block by block of a network and executes event handlers that are
// configured for the listener.
func (l *EVMProvider) Listener(ctx context.Context, startHeight uint64, blockInfo chan types.BlockInfo) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Get the latest block height
			latestHeight, err := l.QueryLatestHeight(ctx)
			if err != nil {
				return err
			}

			// If the start height is greater than the latest height, return an error
			if startHeight > latestHeight {
				return fmt.Errorf("start height %d is greater than latest height %d", startHeight, latestHeight)
			}

			// If the start height is equal to the latest height, return
			if startHeight == latestHeight {
				return nil
			}

			// If the start height is 0, set it to the latest height
			if startHeight == 0 {
				startHeight = latestHeight
			}

			// Start at the start height and go until the latest height
			for i := startHeight; i <= latestHeight; i++ {
				// Get the block at the current height
				block, err := l.QueryBlockByHeight(ctx, i)
				if err != nil {
					return err
				}

				// If the block is empty, continue
				if block == nil {
					continue
				}

				// Execute the event handlers for the block
				if err := l.ExecuteEventHandlers(ctx, block); err != nil {
					return err
				}
			}
		}
		return nil
	}
}
