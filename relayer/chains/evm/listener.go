package evm

import (
	"context"
	"sort"
	"time"

	"github.com/icon-project/centralized-relay/relayer/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ListenToEvents goes block by block of a network and executes event handlers that are
// configured for the listener.
func (p *EVMProvider) Listener(ctx context.Context, startHeight uint64, blockInfo chan types.BlockInfo) error {
	errCh := make(chan error)                                            // error channel
	reconnectCh := make(chan struct{}, 1)                                // reconnect channel
	btpBlockNotifCh := make(chan *types.BlockNotification, 100)          // block notification channel
	btpBlockRespCh := make(chan *btpBlockResponse, cap(btpBlockNotifCh)) // block result channel

	reconnect := func() {
		select {
		case reconnectCh <- struct{}{}:
		default:
		}
		for len(btpBlockRespCh) > 0 || len(btpBlockNotifCh) > 0 {
			select {
			case <-btpBlockRespCh: // clear block result channel
			case <-btpBlockNotifCh: // clear block notification channel
			}
		}
	}

	processedheight, err := p.startFromHeight(ctx, startHeight)
	if err != nil {
		return errors.Wrapf(err, "failed to calculate start height")
	}

	p.log.Info("Start querying from height", zap.Int64("height", processedheight))
	// subscribe to monitor block
	ctxMonitorBlock, cancelMonitorBlock := context.WithCancel(ctx)
	reconnect()

	blockReq := &types.BlockRequest{
		Height:       types.NewHexInt(int64(processedheight)),
		EventFilters: GetMonitorEventFilters(p.cfg.ContractAddress),
	}

loop:
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			return err

		case <-reconnectCh:
			cancelMonitorBlock()
			ctxMonitorBlock, cancelMonitorBlock = context.WithCancel(ctx)

			go func(ctx context.Context, cancel context.CancelFunc) {
				blockReq.Height = types.NewHexInt(int64(processedheight))
				p.log.Debug("Try to reconnect from", zap.Int64("height", processedheight))
				err := p.client.eth(ctx, blockReq, func(conn *websocket.Conn, v *types.BlockNotification) error {
					if !errors.Is(ctx.Err(), context.Canceled) {
						btpBlockNotifCh <- v
					}
					return nil
				}, func(conn *websocket.Conn) {
				}, func(conn *websocket.Conn, err error) {})
				if err != nil {
					if errors.Is(err, context.Canceled) {
						return
					}
					time.Sleep(time.Second * 5)
					reconnect()
					p.log.Warn("Error occured during monitor block", zap.Error(err))
				}
			}(ctxMonitorBlock, cancelMonitorBlock)
		case br := <-btpBlockRespCh:
			for ; br != nil; processedheight++ {
				p.log.Debug("Verified block ",
					zap.Int64("height", int64(processedheight)))

				message := parseMessagesFromEventlogs(p.log, br.EventLogs, uint64(br.Height))

				// TODO: check for the concurrency
				incoming <- providerTypes.BlockInfo{
					Messages: message,
					Height:   uint64(br.Height),
				}

				if br = nil; len(btpBlockRespCh) > 0 {
					br = <-btpBlockRespCh
				}
			}
			// remove unprocessed blockResponses
			for len(btpBlockRespCh) > 0 {
				<-btpBlockRespCh
			}

		default:
			select {
			default:
			case bn := <-btpBlockNotifCh:
				requestCh := make(chan *btpBlockRequest, cap(btpBlockNotifCh))
				for i := int64(0); bn != nil; i++ {
					height, err := bn.Height.Value()

					if err != nil {
						return err
					} else if height != processedheight+i {
						p.log.Warn("Reconnect: missing block notification",
							zap.Int64("got", height),
							zap.Int64("expected", processedheight+i),
						)
						reconnect()
						continue loop
					}

					requestCh <- &btpBlockRequest{
						height:  height,
						hash:    bn.Hash,
						indexes: bn.Indexes,
						events:  bn.Events,
						retry:   maxRetires,
					}
					if bn = nil; len(btpBlockNotifCh) > 0 && len(requestCh) < cap(requestCh) {
						bn = <-btpBlockNotifCh
					}
				}

				brs := make([]*btpBlockResponse, 0, len(requestCh))
				for request := range requestCh {
					switch {
					case request.err != nil:
						if request.retry > 0 {
							request.retry--
							request.response, request.err = nil, nil
							requestCh <- request
							continue
						}
						p.log.Info("Request error ",
							zap.Any("height", request.height),
							zap.Error(request.err))
						brs = append(brs, nil)
						if len(brs) == cap(brs) {
							close(requestCh)
						}
					case request.response != nil:
						brs = append(brs, request.response)
						if len(brs) == cap(brs) {
							close(requestCh)
						}
					default:
						go p.handleBTPBlockRequest(request, requestCh)
					}
				}
				// filter nil
				_brs, brs := brs, brs[:0]
				for _, v := range _brs {
					if v != nil {
						brs = append(brs, v)
					}
				}

				// sort and forward notifications
				if len(brs) > 0 {
					sort.SliceStable(brs, func(i, j int) bool {
						return brs[i].Height < brs[j].Height
					})
					for i, d := range brs {
						if d.Height == processedheight+int64(i) {
							btpBlockRespCh <- d
						}
					}
				}
			}
		}
	}
}

func (p *EVMProvider) startFromHeight(ctx context.Context, lastSavedHeight uint64) (int64, error) {
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
