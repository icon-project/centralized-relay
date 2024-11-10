package stacks

import (
	"context"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func (p *Provider) Listener(ctx context.Context, lastProcessedTx providerTypes.LastProcessedTx, blockInfoChan chan *providerTypes.BlockInfo) error {
	p.log.Info("Starting Stacks event listener")

	wsURL := p.client.GetWebSocketURL()
	p.log.Debug("Using WebSocket URL", zap.String("url", wsURL))

	eventSystem := events.NewEventSystem(ctx, wsURL, p.log)

	eventSystem.OnEvent(func(event *events.Event) error {
		msg, err := p.getRelayMessageFromEvent(event.Type, event.Data)
		if err != nil {
			p.log.Error("Failed to parse relay message from event",
				zap.Error(err),
				zap.String("eventType", event.Type))
			return err
		}

		msg.MessageHeight = event.BlockHeight

		blockInfo := &providerTypes.BlockInfo{
			Height:   event.BlockHeight,
			Messages: []*providerTypes.Message{msg},
		}

		select {
		case blockInfoChan <- blockInfo:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	if err := eventSystem.Start(); err != nil {
		return fmt.Errorf("failed to start event system: %w", err)
	}

	p.log.Info("Stacks event listener started successfully")

	<-ctx.Done()
	p.log.Info("Stopping Stacks event listener")
	eventSystem.Stop()
	return ctx.Err()
}
