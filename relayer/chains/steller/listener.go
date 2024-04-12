package steller

import (
	"context"

	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
)

func (p *Provider) Listener(ctx context.Context, lastSavedCheckpointSeq uint64, blockInfo chan *relayertypes.BlockInfo) error {
	//
	return nil
}
