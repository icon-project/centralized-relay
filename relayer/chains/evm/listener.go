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
	pollerTime                 = 5 * time.Second
	PollFetch                  = "poll"
	WsFetch                    = "ws"
	RPCRedundancy              = "rpc-verify"
	BaseRetryInterval          = 3 * time.Second
	MaxRetryInterval           = 5 * time.Minute
	MaxRetryCount              = 5
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

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfoChan chan *relayertypes.BlockInfo) error {
	startHeight, err := p.startFromHeight(ctx, lastSavedHeight)
	if err != nil {
		return err
	}
	p.log.Info("Start from height ", zap.Uint64("height", startHeight), zap.Uint64("finality block", p.FinalityBlock(ctx)))
	return p.listenNormalWsNPoll(ctx, startHeight, blockInfoChan)
}

func (p *Provider) listenNormalWsNPoll(ctx context.Context, startHeight uint64, blockInfoChan chan *relayertypes.BlockInfo) error {
	var (
		subscribeStart = time.NewTicker(time.Second * 1)
		pollerStart    = time.NewTicker(pollerTime)
		errChan        = make(chan error)
	)

	if p.cfg.Fetch == PollFetch {
		if p.cfg.Redundancy == RPCRedundancy {
			panic("cannot set rpc-verify on poll fetch strategy")
		}
		subscribeStart.Stop()
		p.backlogProcessing = true
	} else {
		pollerStart.Stop()
		if p.cfg.Redundancy == RPCRedundancy {
			p.backlogProcessing = true
		}
	}

	for {
		select {
		case <-ctx.Done():
			p.log.Debug("evm listener: done")
			return ctx.Err()
		case err := <-errChan:
			if p.isConnectionError(err) {
				p.log.Error("connection error", zap.Error(err))
				client, err := p.client.Reconnect()
				if err != nil {
					p.log.Error("failed to reconnect", zap.Error(err))
				} else {
					p.log.Info("client reconnected")
					p.client = client
				}
			}
			startHeight = p.GetLastSavedBlockHeight()
			subscribeStart.Reset(time.Second * 1)
		case <-subscribeStart.C:
			subscribeStart.Stop()
			latestHeight := p.latestHeight(ctx)
			go p.Subscribe(ctx, blockInfoChan, errChan, latestHeight)
			p.backlogProcessing = true
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
				p.log.Info("synced", zap.Uint64("start", br.start), zap.Uint64("end", br.end), zap.Uint64("latest", latestHeight))
				p.log.Info("synced", zap.Uint64("start", br.start), zap.Uint64("end", br.end), zap.Uint64("latest", latestHeight), zap.Uint64("delta", latestHeight-br.end))
				for _, log := range logs {
					message, err := p.getRelayMessageFromLog(log)
					if err != nil {
						p.log.Error("failed to get relay message from log", zap.Error(err))
						continue
					}
					p.log.Info("Detected eventlog",
						zap.String("target_network", message.Dst),
						zap.Uint64("sn", message.Sn.Uint64()),
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
			p.saveHeightFunc(latestHeight)
			p.backlogProcessing = false
		case <-pollerStart.C:
			latestHeight := p.latestHeight(ctx)
			var blockReqs []*blockReq
			if startHeight < latestHeight {
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
					p.log.Info("syncing", zap.Uint64("start", br.start), zap.Uint64("end", br.end), zap.Uint64("latest", latestHeight))
					logs, err := p.getLogsRetry(ctx, filter)
					if err != nil {
						p.log.Warn("failed to fetch blocks", zap.Uint64("from", br.start), zap.Uint64("to", br.end), zap.Error(err))
						continue
					}
					p.log.Info("synced", zap.Uint64("start", br.start), zap.Uint64("end", br.end), zap.Uint64("latest", latestHeight))
					for _, log := range logs {
						message, err := p.getRelayMessageFromLog(log)
						if err != nil {
							p.log.Error("failed to get relay message from log", zap.Error(err))
							continue
						}
						p.log.Info("Detected eventlog",
							zap.String("target_network", message.Dst),
							zap.Uint64("sn", message.Sn.Uint64()),
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
				p.saveHeightFunc(latestHeight)
				startHeight = latestHeight
			}
			pollerStart.Reset(pollerTime)
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
			zap.Uint64("height", lbn.Height.Uint64()),
			zap.String("target_network", message.Dst),
			zap.Uint64("sn", message.Sn.Uint64()),
			zap.String("event_type", message.EventType),
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
func (p *Provider) Subscribe(ctx context.Context,
	blockInfoChan chan *relayertypes.BlockInfo, resetCh chan error,
	latestHeight uint64) error {
	ch := make(chan ethTypes.Log, 10)
	wsEventsFound := false
	sub, err := p.client.Subscribe(ctx, ethereum.FilterQuery{
		Addresses: p.blockReq.Addresses,
		Topics:    p.blockReq.Topics,
		//TODO: required for some rpcs, review
		FromBlock: new(big.Int).SetUint64(latestHeight),
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
			return nil
		case err := <-sub.Err():
			p.log.Error("subscription error", zap.Error(err))
			resetCh <- err
			return err
		case log := <-ch:
			wsEventsFound = true
			message, err := p.getRelayMessageFromLog(log)
			if err != nil {
				p.log.Error("failed to get relay message from log", zap.Error(err))
				continue
			}
			p.log.Info("Detected eventlog",
				zap.String("target_network", message.Dst),
				zap.Uint64("sn", message.Sn.Uint64()),
				zap.String("event_type", message.EventType),
				zap.String("tx_hash", log.TxHash.String()),
				zap.Uint64("block_number", log.BlockNumber),
			)
			blockInfoChan <- &relayertypes.BlockInfo{
				Height:   log.BlockNumber,
				Messages: []*relayertypes.Message{message},
			}
		case <-time.After(time.Minute * 2):
			ctx, cancel := context.WithTimeout(ctx, websocketReadTimeout)
			defer cancel()
			if _, err := p.client.GetHeaderByHeight(ctx, big.NewInt(1)); err != nil {
				resetCh <- err
				return err
			}
			if p.cfg.Redundancy == RPCRedundancy {
				lastSavedHeight := p.GetLastSavedBlockHeight()
				if lastSavedHeight == 0 {
					lastSavedHeight = latestHeight
				}
				p.handleRpcVerification(ctx, wsEventsFound, lastSavedHeight, resetCh, blockInfoChan)
			}
		}

	}
}

func (p *Provider) handleRpcVerification(ctx context.Context, wsEventsFound bool, latestHeight uint64,
	resetCh chan error, blockInfoChan chan *relayertypes.BlockInfo) {
	if !wsEventsFound {
		p.log.Info("No events found in ws,verifying rpc")
		currentLatestHeight := p.latestHeight(ctx)
		var blockReqs []*blockReq
		syncBatchSize := p.cfg.BlockBatchSize
		if latestHeight < currentLatestHeight {
			for start := latestHeight; start <= currentLatestHeight; start += syncBatchSize {
				end := min(start+syncBatchSize-1, currentLatestHeight)
				blockReqs = append(blockReqs, &blockReq{start, end, nil})
			}
			messageFound := false
			for _, br := range blockReqs {
				filter := ethereum.FilterQuery{
					FromBlock: new(big.Int).SetUint64(br.start),
					ToBlock:   new(big.Int).SetUint64(br.end),
					Addresses: p.blockReq.Addresses,
					Topics:    p.blockReq.Topics,
				}
				p.log.Info("syncing missing", zap.Uint64("start", br.start), zap.Uint64("end", br.end), zap.Uint64("latest", latestHeight))
				logs, err := p.getLogsRetry(ctx, filter)
				if err != nil {
					p.log.Warn("failed to fetch missing blocks", zap.Uint64("from", br.start), zap.Uint64("to", br.end), zap.Error(err))
					continue
				}
				if len(logs) > 0 {
					p.log.Info("got missing logs", zap.Uint64("start", br.start), zap.Uint64("end", br.end), zap.Uint64("latest", latestHeight))
				}
				p.log.Info("synced", zap.Uint64("start", br.start), zap.Uint64("end", br.end), zap.Uint64("latest", latestHeight))
				for _, log := range logs {
					message, err := p.getRelayMessageFromLog(log)
					if err != nil {
						p.log.Error("failed to get relay message from log", zap.Error(err))
						continue
					}
					p.log.Info("Detected missing eventlog",
						zap.String("target_network", message.Dst),
						zap.Uint64("sn", message.Sn.Uint64()),
						zap.String("event_type", message.EventType),
						zap.String("tx_hash", log.TxHash.String()),
						zap.Uint64("block_number", log.BlockNumber),
					)
					messageFound = true
					blockInfoChan <- &relayertypes.BlockInfo{
						Height:   log.BlockNumber,
						Messages: []*relayertypes.Message{message},
					}
				}
			}
			if messageFound {
				p.log.Info("Need to reset ws connections")
				resetCh <- errors.New("ws stale connection")
			}
		}
		p.saveHeightFunc(currentLatestHeight)
	}
}
