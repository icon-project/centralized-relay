package steller

import (
	"context"

	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func (p *Provider) Listener(ctx context.Context, lastSavedCheckpointSeq uint64, blockInfo chan *relayertypes.BlockInfo) error {
	//Todo
	latestLedger, err := p.client.GetLatestLedger(ctx)
	if err != nil {
		return err
	}

	p.log.Info("steller listener started from", zap.Uint64("height", latestLedger.Sequence))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
