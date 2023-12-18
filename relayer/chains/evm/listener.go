package evm

import (
	"context"
	"math/big"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/icon-project/centralized-relay/relayer/chains/evm/types"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/pkg/errors"
)

const (
	BlockInterval              = 2 * time.Second
	BlockHeightPollInterval    = 60 * time.Second
	defaultReadTimeout         = 15 * time.Second
	monitorBlockMaxConcurrency = 10 // number of concurrent requests to synchronize older blocks from source chain
	DefaultFinalityBlock       = 10
)

type BnOptions struct {
	StartHeight uint64
	Concurrency uint64
}

func (r *EVMProvider) latestHeight() uint64 {
	height, err := r.client.GetBlockNumber()
	if err != nil {
		// TODO:
		// r.Log.WithFields(log.Fields{"error": err}).Error("receiveLoop: failed to GetBlockNumber")
		return 0
	}
	return height
}

func (r *EVMProvider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfoChan chan relayertypes.BlockInfo) error {
	startHeight, err := r.startFromHeight(ctx, lastSavedHeight)
	if err != nil {
		return err
	}

	r.log.Info("Start query from height ", zap.Uint64("start-height", startHeight))

	heightTicker := time.NewTicker(BlockInterval)
	defer heightTicker.Stop()

	heightPoller := time.NewTicker(BlockHeightPollInterval)
	defer heightPoller.Stop()

	next, latest := startHeight, r.latestHeight()
	concurrency := r.GetConcurrency(ctx, startHeight, latest)
	// block notification channel
	// (buffered: to avoid deadlock)
	// increase concurrency parameter for faster sync
	bnch := make(chan *types.BlockNotification, concurrency)
	// last unverified block notification
	var lbn *types.BlockNotification
	// Loop started
	for {
		select {
		case <-ctx.Done():
			r.log.Debug("evm listener: context done")
			return nil

		case <-heightTicker.C:
			r.log.Debug("receiveLoop: heightTicker", zap.Uint64("latest", latest))
			latest++

		case <-heightPoller.C:
			if height := r.latestHeight(); height > latest {
				latest = height
				if next > latest {
					// TODO:
					r.log.Debug("receiveLoop: skipping; ", zap.Uint64("latest", latest), zap.Uint64("next", next))
				}
			}

		case bn := <-bnch:
			// process all notifications
			for ; bn != nil; next++ {
				if lbn != nil {
					r.log.Debug("block-notification received", zap.Uint64("height", lbn.Height.Uint64()),
						zap.Int64("gas-used", int64(lbn.Header.GasUsed)))

					messages, err := r.FindMessages(ctx, lbn)
					if err != nil {
						return errors.Wrapf(err, "receiveLoop: callback: %v", err)
					}
					blockInfoChan <- relayertypes.BlockInfo{
						Height:   lbn.Height.Uint64(),
						Messages: messages,
					}
				}

				if lbn, bn = bn, nil; len(bnch) > 0 {
					bn = <-bnch
				}
			}
			// remove unprocessed notifications
			for len(bnch) > 0 {
				<-bnch
			}

		default:
			if next >= latest {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			type bnq struct {
				h     uint64
				v     *types.BlockNotification
				err   error
				retry int
			}

			qch := make(chan *bnq, cap(bnch))
			for i := next; i < latest &&
				len(qch) < cap(qch); i++ {
				qch <- &bnq{i, nil, nil, 3} // fill bch with requests
			}
			bns := make([]*types.BlockNotification, 0, len(qch))
			for q := range qch {
				switch {
				case q.err != nil:
					if q.retry > 0 {
						if !strings.HasSuffix(q.err.Error(), "requested block number greater than current block number") {
							q.retry--
							q.v, q.err = nil, nil
							qch <- q
							continue
						}
						if latest >= q.h {
							latest = q.h - 1
						}
					}
					// r.Log.Debugf("receiveLoop: bnq: h=%d:%v, %v", q.h, q.v.Header.Hash(), q.err)
					bns = append(bns, nil)
					if len(bns) == cap(bns) {
						close(qch)
					}

				case q.v != nil:
					bns = append(bns, q.v)
					if len(bns) == cap(bns) {
						close(qch)
					}
				default:
					go func(q *bnq) {
						defer func() {
							time.Sleep(500 * time.Millisecond)
							qch <- q
						}()
						if q.v == nil {
							q.v = new(types.BlockNotification)
						}
						q.v.Height = new(big.Int).SetUint64(q.h)
						q.v.Header, q.err = r.client.GetHeaderByHeight(ctx, q.v.Height)
						if q.err != nil {
							q.err = errors.Wrapf(q.err, "GetEvmHeaderByHeight %v", q.err)
							return
						}
						ht := big.NewInt(q.v.Height.Int64())

						if q.v.Header.GasUsed > 0 {
							r.blockReq.FromBlock = ht
							r.blockReq.ToBlock = ht
							q.v.Logs, q.err = r.client.FilterLogs(context.TODO(), r.blockReq)
							if q.err != nil {
								q.err = errors.Wrapf(q.err, "FilterLogs: %v", q.err)
								return
							}
						}
					}(q)
				}
			}
			// filter nil
			_bns_, bns := bns, bns[:0]
			for _, v := range _bns_ {
				if v != nil {
					bns = append(bns, v)
				}
			}
			// sort and forward notifications
			if len(bns) > 0 {
				sort.SliceStable(bns, func(i, j int) bool {
					return bns[i].Height.Uint64() < bns[j].Height.Uint64()
				})
				for i, v := range bns {
					if v.Height.Uint64() == next+uint64(i) {
						bnch <- v
					}
				}
			}
		}
	}
}

func (p *EVMProvider) FindMessages(ctx context.Context, lbn *types.BlockNotification) ([]*relayertypes.Message, error) {
	if lbn == nil && lbn.Logs == nil {
		return nil, nil
	}
	messages := make([]*relayertypes.Message, 0)
	for _, log := range lbn.Logs {
		message, err := p.getRelayMessageFromLog(log)
		if err != nil {
			return nil, err
		}
		p.log.Debug("detected eventlog ", zap.Uint64("height", lbn.Height.Uint64()),
			zap.String("target-network", message.Dst),
			zap.Uint64("sn", message.Sn),
			zap.String("event-type", message.EventType),
		)
		messages = append(messages, message)
	}
	return messages, nil
}

func (p *EVMProvider) GetConcurrency(ctx context.Context, startHeight, currentHeight uint64) int {
	concurrency := p.cfg.Concurrency
	if concurrency == 0 {
		concurrency = monitorBlockMaxConcurrency
	}
	// we calculate concurrency based on the height to sync
	// so that we avoid duplicate block number is picked up by multiple workers
	heightTosync := currentHeight - startHeight
	if heightTosync < 1 {
		concurrency = 1 // we don't want to span multiple workers for 1 block
	} else if heightTosync < concurrency {
		concurrency = heightTosync // we don't want to span more workers than the height to sync
	}
	return int(concurrency)
}

func (p *EVMProvider) startFromHeight(ctx context.Context, lastSavedHeight uint64) (uint64, error) {
	latestHeight, err := p.QueryLatestHeight(ctx)
	if err != nil {
		return 0, err
	}

	latestQueryHeight := latestHeight - p.cfg.FinalityBlock

	if p.cfg.StartHeight > latestQueryHeight {
		p.log.Error("start height provided on config cannot be greater than latest query height",
			zap.Uint64("start-height", p.cfg.StartHeight),
			zap.Uint64("latest-height", latestQueryHeight),
		)
	}

	// priority1: lastsaveheight from db
	if lastSavedHeight != 0 && lastSavedHeight < latestQueryHeight {
		return lastSavedHeight, nil
	}

	// priority2: startHeight from config
	if p.cfg.StartHeight != 0 && p.cfg.StartHeight < latestQueryHeight {
		return p.cfg.StartHeight, nil
	}

	// priority3: latest height
	return latestHeight, nil
}
