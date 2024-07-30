package icon

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/goloop/common"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type btpBlockResponse struct {
	Height    int64
	Hash      common.HexHash
	Header    *types.BlockHeader
	EventLogs []*types.EventLog
}

type btpBlockRequest struct {
	height   int64
	hash     types.HexBytes
	indexes  [][]types.HexInt
	events   [][][]types.HexInt
	err      error
	retry    uint8
	response *btpBlockResponse
}

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, incoming chan *providerTypes.BlockInfo) error {
	reconnectCh := make(chan struct{}, 1) // reconnect channel

	reconnect := func() {
		select {
		case reconnectCh <- struct{}{}:
		default:
		}
	}

	processedheight, err := p.StartFromHeight(ctx, lastSavedHeight)
	if err != nil {
		return errors.Wrapf(err, "failed to calculate start height")
	}

	p.log.Info("Start from height", zap.Int64("height", processedheight), zap.Uint64("finality block", p.FinalityBlock(ctx)))
	// subscribe to monitor block
	reconnect()

	eventReq := &types.EventRequest{
		Height:           types.NewHexInt(processedheight),
		EventFilter:      p.GetMonitorEventFilters(),
		Logs:             types.NewHexInt(1),
		ProgressInterval: types.NewHexInt(25),
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-reconnectCh:
			ctxMonitorBlock, cancelMonitorBlock := context.WithCancel(ctx)
			go func(ctx context.Context, cancel context.CancelFunc) {
				err := p.client.MonitorEvent(ctx, eventReq, incoming, func(v *types.EventNotification, outgoing chan *providerTypes.BlockInfo) error {
					if !errors.Is(ctx.Err(), context.Canceled) {
						p.log.Debug("event notification received", zap.Any("event", v))
						if v.Progress != "" {
							height, err := v.Progress.Uint64()
							if err != nil {
								p.log.Error("failed to get progress height", zap.Error(err))
								return err
							}
							p.SetLastProcessedHeight(height)
							return nil
						}
						msgs, err := p.parseMessageEvent(v)
						if err != nil {
							p.log.Error("failed to parse message event", zap.Error(err))
							return err
						}
						for _, msg := range msgs {
							p.log.Info("Detected eventlog",
								zap.Uint64("height", msg.MessageHeight),
								zap.String("target_network", msg.Dst),
								zap.Uint64("sn", msg.Sn.Uint64()),
								zap.String("tx_hash", v.Hash.String()),
								zap.String("event_type", msg.EventType),
							)
							outgoing <- &providerTypes.BlockInfo{
								Messages: []*providerTypes.Message{msg},
								Height:   msg.MessageHeight,
							}
						}
					}
					return err
				}, func(conn *websocket.Conn, err error) {})
				if err != nil {
					if errors.Is(err, context.Canceled) {
						return
					}
					eventReq.Height = types.NewHexInt(int64(p.GetCheckpoint()))
					time.Sleep(time.Second * 3)
					reconnect()
					p.log.Warn("error occured during monitor event", zap.Error(err))
				}
			}(ctxMonitorBlock, cancelMonitorBlock)
		}
	}
}

func (p *Provider) StartFromHeight(ctx context.Context, lastSavedHeight uint64) (int64, error) {
	latestHeight, err := p.QueryLatestHeight(ctx)
	if err != nil {
		return 0, err
	}

	if p.cfg.StartHeight > latestHeight {
		p.log.Error("start height provided on config cannot be greater than latest height",
			zap.Uint64("start-height", p.cfg.StartHeight),
			zap.Int64("latest-height", int64(latestHeight)),
		)
	}

	// priority1: startHeight from config
	if p.cfg.StartHeight != 0 && p.cfg.StartHeight < latestHeight {
		return int64(p.cfg.StartHeight), nil
	}

	// priority2: lastsaveheight from db
	if lastSavedHeight != 0 && lastSavedHeight < latestHeight {
		return int64(lastSavedHeight), nil
	}

	// priority3: latest height
	return int64(latestHeight), nil
}
