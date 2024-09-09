package evm

import (
	"context"
	"math"
	"math/big"
	"runtime"
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
	defaultReadTimeout         = 60 * time.Second
	websocketReadTimeout       = 10 * time.Second
	monitorBlockMaxConcurrency = 10 // number of concurrent requests to synchronize older blocks from source chain
	maxBlockRange              = 50
	maxBlockQueryFailedRetry   = 5
	DefaultFinalityBlock       = 10
	BaseRetryInterval          = 3 * time.Second
	MaxRetryInterval           = 5 * time.Minute
	MaxRetryCount              = 5
	ClientReconnectDelay       = 5 * time.Second
)

type BnOptions struct {
	StartHeight uint64
	Concurrency uint64
}

type blockReq struct {
	start, end uint64
	err        error
}

func (r *Provider) latestHeight(ctx context.Context) uint64 {
	height, err := r.client.GetBlockNumber(ctx)
	if err != nil {
		r.log.Error("Evm listener: failed to GetBlockNumber", zap.Error(err))
		return 0
	}
	return height
}

func (p *Provider) Listener(ctx context.Context, lastProcessedTx relayertypes.LastProcessedTx, blockInfoChan chan *relayertypes.BlockInfo) error {
	lastSavedHeight := lastProcessedTx.Height

	startHeight, err := p.startFromHeight(ctx, lastSavedHeight)
	if err != nil {
		return err
	}
	p.log.Info("Start from height ", zap.Uint64("height", startHeight), zap.Uint64("finality block", p.FinalityBlock(ctx)))

	var (
		subscribeStart = time.NewTicker(time.Second * 1)
		errChan        = make(chan error)
	)

	for {
		select {
		case <-ctx.Done():
			p.log.Debug("evm listener: done")
			return ctx.Err()
		case err := <-errChan:
			p.log.Error("connection error", zap.Error(err))
			clientReconnected := false
			for !clientReconnected {
				p.log.Info("reconnecting client")
				client, err := p.client.Reconnect()
				if err == nil {
					clientReconnected = true
					p.log.Info("client reconnected")
					p.client = client
				} else {
					p.log.Error("failed to re-connect", zap.Error(err))
					time.Sleep(ClientReconnectDelay)
				}
			}
			startHeight = p.GetLastSavedBlockHeight()
			subscribeStart.Reset(time.Second * 1)
		case <-subscribeStart.C:
			subscribeStart.Stop()
			go p.Subscribe(ctx, blockInfoChan, errChan)

			latestHeight := p.latestHeight(ctx)
			if startHeight == 0 {
				startHeight = latestHeight
			}

			var blockReqs []*blockReq
			for start := startHeight; start <= latestHeight; start += p.cfg.BlockBatchSize {
				end := min(start+p.cfg.BlockBatchSize-1, latestHeight)
				blockReqs = append(blockReqs, &blockReq{start, end, nil})
			}

			for _, br := range blockReqs {
				filter := ethereum.FilterQuery{
					FromBlock: new(big.Int).SetUint64(br.start),
					ToBlock:   new(big.Int).SetUint64(br.end),
					Addresses: p.blockReq.Addresses,
					Topics:    p.blockReq.Topics,
				}
				p.log.Info("syncing", zap.Uint64("start", br.start), zap.Uint64("end", br.end), zap.Uint64("latest", latestHeight), zap.Uint64("delta", latestHeight-br.end))
				logs, err := p.getLogsRetry(ctx, filter)
				if err != nil {
					p.log.Warn("failed to fetch blocks", zap.Uint64("from", br.start), zap.Uint64("to", br.end), zap.Error(err))
					continue
				}
				p.log.Info("synced", zap.Uint64("start", br.start), zap.Uint64("end", br.end), zap.Uint64("latest", latestHeight), zap.Uint64("delta", latestHeight-br.end))
				for _, log := range logs {
					message, err := p.getRelayMessageFromLog(log)
					if err != nil {
						p.log.Error("failed to get relay message from log", zap.Error(err))
						continue
					}
					p.log.Info("Detected eventlog",
						zap.String("dst", message.Dst),
						zap.Uint64("sn", message.Sn.Uint64()),
						zap.Any("req_id", message.ReqID),
						zap.String("event_type", message.EventType),
						zap.String("tx_hash", log.TxHash.String()),
						zap.Uint64("height", log.BlockNumber),
					)
					blockInfoChan <- &relayertypes.BlockInfo{
						Height:   log.BlockNumber,
						Messages: []*relayertypes.Message{message},
					}
				}
			}
		}
	}
}

func (p *Provider) getLogsRetry(ctx context.Context, filter ethereum.FilterQuery) ([]ethTypes.Log, error) {
	var (
		logs     []ethTypes.Log
		err      error
		attempts = 0
	)

	for attempts < MaxRetryCount {
		logs, err = p.client.FilterLogs(ctx, filter)
		if err == nil {
			return logs, nil
		}
		attempts++
		delay := time.Duration(math.Pow(2, float64(attempts))) * BaseRetryInterval
		if delay > MaxRetryInterval {
			delay = MaxRetryInterval
		}
		p.log.Error("failed to get logs", zap.Error(err), zap.Int("attempts", attempts), zap.Duration("delay", delay), zap.Uint64("from", filter.FromBlock.Uint64()), zap.Uint64("to", filter.ToBlock.Uint64()))
		time.Sleep(delay)
	}
	return nil, err
}

func (p *Provider) isConnectionError(err error) bool {
	return strings.Contains(err.Error(), "tcp") || errors.Is(err, context.DeadlineExceeded)
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
			zap.String("dst", message.Dst),
			zap.Uint64("sn", message.Sn.Uint64()),
			zap.Any("req_id", message.ReqID),
			zap.String("event_type", message.EventType),
			zap.String("tx_hash", log.TxHash.String()),
			zap.Uint64("height", log.BlockNumber),
		)
		messages = append(messages, message)
	}
	return messages, nil
}

func (p *Provider) GetConcurrency(ctx context.Context, startHeight, currentHeight uint64) int {
	diff := int((currentHeight-startHeight)/p.cfg.BlockBatchSize) + 1
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
func (p *Provider) Subscribe(ctx context.Context, blockInfoChan chan *relayertypes.BlockInfo, resetCh chan error) error {
	ch := make(chan ethTypes.Log, 10)
	subContext, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	sub, err := p.client.Subscribe(subContext, ethereum.FilterQuery{
		Addresses: p.blockReq.Addresses,
		Topics:    p.blockReq.Topics,
	}, ch)
	if err != nil {
		p.log.Error("failed to subscribe", zap.Error(err))
		resetCh <- err
		return err
	}
	defer sub.Unsubscribe()
	defer close(ch)
	p.log.Info("Subscribed to new blocks", zap.Any("address", p.blockReq.Addresses))
	for {
		select {
		case <-ctx.Done():
			p.log.Debug("subscriptions stopped")
			return ctx.Err()
		case err := <-sub.Err():
			p.log.Warn("subscription error", zap.Error(err))
			resetCh <- err
			return err
		case log := <-ch:
			message, err := p.getRelayMessageFromLog(log)
			if err != nil {
				p.log.Error("failed to get relay message from log", zap.Error(err))
				continue
			}
			p.log.Info("Detected eventlog",
				zap.String("dst", message.Dst),
				zap.Uint64("sn", message.Sn.Uint64()),
				zap.Any("req_id", message.ReqID),
				zap.String("event_type", message.EventType),
				zap.String("tx_hash", log.TxHash.String()),
				zap.Uint64("height", log.BlockNumber),
			)
			blockInfoChan <- &relayertypes.BlockInfo{
				Height:   log.BlockNumber,
				Messages: []*relayertypes.Message{message},
			}
		case <-time.After(time.Minute * 2):
			ctx, cancel := context.WithTimeout(ctx, websocketReadTimeout)
			defer cancel()
			if _, err := p.QueryLatestHeight(ctx); err != nil {
				resetCh <- err
				return err
			}
		}
	}
}
