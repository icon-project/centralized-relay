package solana

import (
	"context"

	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfo chan *relayertypes.BlockInfo) error {
	latestHeight, err := p.client.GetLatestBlockHeight(ctx)
	if err != nil {
		return err
	}

	startHeight := latestHeight
	if lastSavedHeight != 0 && lastSavedHeight < latestHeight {
		startHeight = lastSavedHeight
	}

	p.log.Info("started querying from height", zap.Uint64("height", startHeight))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
