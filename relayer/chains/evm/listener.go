package evm

import (
	"context"
	"math/big"
	"sort"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/icon-project/centralized-relay/relayer/chains/evm/types"
)

// Listener goes block by block of a network and executes event handlers that are
// configured for the listener.
func (p *EVMProvider) Listener(ctx context.Context, opts *BnOptions, callback func(v *types.BlockNotification) error) error {
	if opts == nil {
		return errors.New("receiveLoop: invalid options: <nil>")
	}

	// block a notification channel
	// (buffered: to avoid deadlock)
	// increase concurrency parameter for faster sync
	bnch := make(chan *types.BlockNotification, SyncConcurrency)

	heightTicker := time.NewTicker(BlockInterval)
	defer heightTicker.Stop()

	heightPoller := time.NewTicker(BlockHeightPollInterval)
	defer heightPoller.Stop()

	latestHeight := func() uint64 {
		height, err := c.GetBlockNumber()
		if err != nil {
			return 0
		}
		return height - BlockFinalityConfirmations
	}
	next, latest := opts.StartHeight, latestHeight()

	// last unverified block notification
	var lbn *types.BlockNotification
	// start monitor loop

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-heightTicker.C:
			latest++

		case <-heightPoller.C:
			if height := latestHeight(); height > 0 {
				latest = height
			}

		case bn := <-bnch:
			// process all notifications
			for ; bn != nil; next++ {
				if lbn != nil {
					if bn.Height.Cmp(lbn.Height) == 0 {
						if bn.Header.ParentHash != lbn.Header.ParentHash {
							break
						}
					} else {
						if vr != nil {
							if err := vr.Verify(lbn.Header, bn.Header, bn.Receipts); err != nil {
								next--
								break
							}
							if err := vr.Update(lbn.Header); err != nil {
								return errors.Wrapf(err, "receiveLoop: vr.Update: %v", err)
							}
						}
						if err := callback(lbn); err != nil {
							return errors.Wrapf(err, "receiveLoop: callback: %v", err)
						}
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
				qch <- &bnq{i, nil, nil, RPCCallRetry} // fill bch with requests
			}
			if len(qch) == 0 {
				c.log.Error("Fatal: Zero length of query channel. Avoiding deadlock")
				continue
			}
			bns := make([]*types.BlockNotification, 0, len(qch))
			for q := range qch {
				switch {
				case q.err != nil:
					if q.retry > 0 {
						q.retry--
						q.v, q.err = nil, nil
						qch <- q
						continue
					}
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
							q.v = &types.BlockNotification{}
						}

						q.v.Height = (&big.Int{}).SetUint64(q.h)

						if q.v.Header == nil {
							header, err := c.GetHeaderByHeight(ctx, q.v.Height)
							if err != nil {
								q.err = errors.Wrapf(err, "GetHeaderByHeight: %v", err)
								return
							}
							q.v.Header = header
							q.v.Hash = q.v.Header.Hash()
						}
						if q.v.Header.GasUsed > 0 {
							if q.v.HasBTPMessage == nil {
								hasBTPMessage, err := r.hasBTPMessage(ctx, q.v.Height)
								if err != nil {
									q.err = errors.Wrapf(err, "hasBTPMessage: %v", err)
									return
								}
								q.v.HasBTPMessage = &hasBTPMessage
							}
							if !*q.v.HasBTPMessage {
								return
							}
							// TODO optimize retry of GetBlockReceipts()
							q.v.Receipts, q.err = c.GetBlockReceipts(q.v.Hash)
							if q.err != nil {
								q.err = errors.Wrapf(q.err, "GetBlockReceipts: %v", q.err)
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
