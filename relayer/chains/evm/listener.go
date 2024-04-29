package evm

import (
	"context"
	"math/big"
	"runtime"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
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

type bnq struct {
	h     uint64
	v     *types.BlockNotification
	err   error
	retry int
}

func (r *Provider) latestHeight() uint64 {
	height, err := r.client.GetBlockNumber()
	if err != nil {
		r.log.Error("Evm listener: failed to GetBlockNumber", zap.Error(err))
		return 0
	}
	return height
}

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfoChan chan *relayertypes.BlockInfo) error {
	startHeight, err := p.startFromHeight(ctx, lastSavedHeight)
	if err != nil {
		return err
	}

	p.log.Info("Start from height ", zap.Uint64("height", startHeight), zap.Uint64("finality block", p.FinalityBlock(ctx)))

	subscribeStart := time.NewTicker(time.Second * 1)

	next, latest := startHeight, p.latestHeight()
	concurrency := p.GetConcurrency(ctx, startHeight, latest)
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
			p.log.Debug("evm listener: done")
			return nil
		case <-subscribeStart.C:
			subscribeStart.Stop()
			go p.Subscribe(ctx, subscribeStart, blockInfoChan)

		case bn := <-bnch:
			// process all notifications
			for ; bn != nil; next++ {
				if lbn != nil {
					p.log.Debug("block-notification received", zap.Uint64("height", lbn.Height.Uint64()),
						zap.Int64("gas-used", int64(lbn.Header.GasUsed)))

					messages, err := p.FindMessages(ctx, lbn)
					if err != nil {
						return errors.Wrapf(err, "receiveLoop: callback: %v", err)
					}
					blockInfoChan <- &relayertypes.BlockInfo{
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
				continue
			}

			qch := make(chan *bnq, cap(bnch))
			for i := next; i < latest && len(qch) < cap(qch); i++ {
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
							qch <- q
						}()
						if q.v == nil {
							q.v = new(types.BlockNotification)
						}
						q.v.Height = new(big.Int).SetUint64(q.h)
						q.v.Header, q.err = p.client.GetHeaderByHeight(ctx, q.v.Height)
						if q.err != nil {
							q.err = errors.Wrapf(q.err, "GetEvmHeaderByHeight %v", q.err)
							return
						}
						ht := big.NewInt(q.v.Height.Int64())

						if q.v.Header.GasUsed > 0 {
							p.blockReq.FromBlock = ht
							p.blockReq.ToBlock = ht
							q.v.Logs, q.err = p.client.FilterLogs(ctx, p.blockReq)
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

func (p *Provider) FindMessages(ctx context.Context, lbn *types.BlockNotification) ([]*relayertypes.Message, error) {
	if lbn == nil && lbn.Logs == nil {
		return nil, nil
	}
	var messages []*relayertypes.Message
	for _, log := range lbn.Logs {
		message, err := p.getRelayMessageFromLog(log)
		if err != nil {
			return nil, err
		}
		p.log.Info("Detected eventlog",
			zap.Uint64("height", lbn.Height.Uint64()),
			zap.String("target_network", message.Dst),
			zap.Uint64("sn", message.Sn),
			zap.String("event_type", message.EventType),
		)
		messages = append(messages, message)
	}
	return messages, nil
}

func (p *Provider) GetConcurrency(ctx context.Context, startHeight, currentHeight uint64) int {
	diff := int(currentHeight - startHeight)
	cpu := runtime.NumCPU()
	if diff <= cpu {
		return diff
	}
	return cpu
}

func (p *Provider) startFromHeight(ctx context.Context, lastSavedHeight uint64) (uint64, error) {
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

	// priority1: startHeight from config
	if p.cfg.StartHeight != 0 && p.cfg.StartHeight < latestQueryHeight {
		return p.cfg.StartHeight, nil
	}

	// priority2: lastsaveheight from db
	if lastSavedHeight != 0 && lastSavedHeight < latestQueryHeight {
		return lastSavedHeight, nil
	}

	// priority3: latest height
	return latestQueryHeight, nil
}

// Subscribe listens to new blocks and sends them to the channel
func (p *Provider) Subscribe(ctx context.Context, ticker *time.Ticker, blockInfoChan chan *relayertypes.BlockInfo) {
	ch := make(chan ethTypes.Log)
	sub, err := p.client.Subscribe(ctx, ethereum.FilterQuery{
		Addresses: p.blockReq.Addresses,
		Topics:    p.blockReq.Topics,
	}, ch)
	if err != nil {
		p.log.Error("failed to subscribe", zap.Error(err))
		ticker.Reset(time.Second * 3)
		return
	}
	defer sub.Unsubscribe()
	p.log.Info("Subscribed to new blocks", zap.Uint64("from", p.blockReq.FromBlock.Uint64()), zap.Any("address", p.blockReq.Addresses))
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-sub.Err():
			p.log.Error("subscription error", zap.Error(err))
			ticker.Reset(time.Second * 3)
			return
		case log := <-ch:
			message, err := p.getRelayMessageFromLog(log)
			if err != nil {
				p.log.Error("failed to get relay message from log", zap.Error(err))
				continue
			}
			p.log.Info("Detected eventlog",
				zap.String("target_network", message.Dst),
				zap.Uint64("sn", message.Sn),
				zap.String("event_type", message.EventType),
				zap.String("tx_hash", log.TxHash.String()),
				zap.Uint64("block_number", log.BlockNumber),
			)
			blockInfoChan <- &relayertypes.BlockInfo{
				Height:   log.BlockNumber,
				Messages: []*relayertypes.Message{message},
			}
		}
	}
}
