package stacks

import (
	"context"
	"time"

	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func (p *Provider) Listener(ctx context.Context, lastProcessedTx providerTypes.LastProcessedTx, blockInfoChan chan *providerTypes.BlockInfo) error {
	eventTypes := p.getSubscribedEventTypes()

	errChan := make(chan error, 1)

	callback := func(eventType string, data interface{}) error {
		msg, err := p.getRelayMessageFromEvent(eventType, data)
		if err != nil {
			p.log.Error("Failed to parse relay message from event", zap.Error(err))

			select {
			case errChan <- err:
			default:
			}
			return err
		}

		blockHeight := uint64(0)

		blockInfo := &providerTypes.BlockInfo{
			Height:   blockHeight,
			Messages: []*providerTypes.Message{msg},
		}

		select {
		case blockInfoChan <- blockInfo:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	for {
		select {
		case <-ctx.Done():
			p.log.Info("Listener context canceled")
			return ctx.Err()
		default:
			p.log.Info("Subscribing to Stacks contract events")
			err := p.client.SubscribeToEvents(ctx, eventTypes, callback)
			if err != nil {
				p.log.Error("Failed to subscribe to events", zap.Error(err))
				select {
				case <-time.After(5 * time.Second):
					p.log.Info("Retrying subscription to events")
				case <-ctx.Done():
					p.log.Info("Listener context canceled during retry wait")
					return ctx.Err()
				}
				continue
			}

			select {
			case err := <-errChan:
				p.log.Error("Error received in event subscription", zap.Error(err))
				select {
				case <-time.After(5 * time.Second):
					p.log.Info("Re-subscribing to events after error")
				case <-ctx.Done():
					p.log.Info("Listener context canceled during resubscription wait")
					return ctx.Err()
				}
			case <-ctx.Done():
				p.log.Info("Listener context canceled")
				return ctx.Err()
			}
		}
	}
}

func (p *Provider) getSubscribedEventTypes() []string {
	eventTypeSet := make(map[string]struct{})
	for _, eventMap := range p.contracts {
		for _, eventType := range eventMap.SigType {
			eventTypeSet[eventType] = struct{}{}
		}
	}
	eventTypes := make([]string, 0, len(eventTypeSet))
	for et := range eventTypeSet {
		eventTypes = append(eventTypes, et)
	}
	return eventTypes
}
