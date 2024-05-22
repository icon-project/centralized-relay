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

	if err := p.RestoreKeystore(ctx); err != nil {
		p.log.Error("failed to restore keystore", zap.Error(err))
	} else {
		p.log.Info("key restore successful: ", zap.String("public-key", p.wallet.PublicKey().String()))
	}

	// if err := p.Route(ctx, &relayertypes.Message{
	// 	Dst:       "0x3.icon",
	// 	EventType: "sendMessage",
	// 	Data:      []byte("hello"),
	// }, func(key *relayertypes.MessageKey, response *relayertypes.TxResponse, err error) {
	// 	if err != nil {
	// 		p.log.Info("message route successful", zap.String("tx-hash", response.TxHash))
	// 	} else {
	// 		p.log.Error("message route failed: ", zap.String("tx-hash", response.TxHash), zap.Error(err))
	// 	}
	// }); err != nil {
	// 	p.log.Error("failed to route message: ", zap.Error(err))
	// }

	// if err := p.InitXcall(ctx); err != nil {
	// 	p.log.Error("failed to init xcall", zap.Error(err))
	// }

	p.log.Info("started querying from height", zap.Uint64("height", startHeight))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
