package sui

import (
	"context"
	"fmt"

	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func (p Provider) Listener(ctx context.Context, lastSavedCheckpointSeq uint64, blockInfo chan *relayertypes.BlockInfo) error {
	//Todo
	latestCheckpointSeq, err := p.client.GetLatestCheckpointSeq(ctx)
	if err != nil {
		return err
	}
	p.log.Info(fmt.Sprintf("Start sui listener from checkpoint %d", latestCheckpointSeq), zap.Uint64("checkpoint", latestCheckpointSeq))
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
