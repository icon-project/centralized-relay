package sui

import (
	"context"
	"strings"

	"github.com/icon-project/centralized-relay/relayer/chains/sui/types"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func (p Provider) Listener(ctx context.Context, lastSavedCheckpointSeq uint64, blockInfo chan *relayertypes.BlockInfo) error {
	return p.listenRealtime(ctx, lastSavedCheckpointSeq, blockInfo)
}

func (p Provider) listenRealtime(ctx context.Context, _ uint64, blockInfo chan *relayertypes.BlockInfo) error {
	eventFilters := []interface{}{
		map[string]interface{}{
			"Package": p.cfg.PackageID,
		},
	}

	done := make(chan interface{})
	defer close(done)
	eventStream, err := p.client.SubscribeEventNotification(done, p.cfg.WsUrl, eventFilters)
	if err != nil {
		p.log.Error("failed to subscribe event notification", zap.Error(err))
		return err
	}

	reconnectCh := make(chan bool)

	p.log.Info("started realtime checkpoint listener")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case en, ok := <-eventStream:
			if ok {
				if en.Error != nil {
					p.log.Error("failed to read event notification", zap.Error(en.Error))
					if strings.Contains(en.Error.Error(), types.WsConnReadError) {
						go func() {
							reconnectCh <- true
						}()
					}
				} else {
					p.log.Info("received new event notification", zap.Any("event", en))
				}
			}
		case val := <-reconnectCh:
			if val {
				p.log.Warn("something went wrong while reading from conn: reconnecting...")
				eventStream, err = p.client.SubscribeEventNotification(done, p.cfg.WsUrl, eventFilters)
				if err != nil {
					return err
				}
				p.log.Warn("connection restablished: listener restarted")
			}
		}
	}
}
