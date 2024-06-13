package evm

import (
	"context"
	"math/big"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/centralized-relay/relayer/chains/evm/types"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
)

const (
	defaultReadTimeout         = 60 * time.Second
	monitorBlockMaxConcurrency = 10 // number of concurrent requests to synchronize older blocks from source chain
	maxBlockRange              = 50
	maxBlockQueryFailedRetry   = 3
	DefaultFinalityBlock       = 10
)

type BnOptions struct {
	StartHeight uint64
	Concurrency uint64
}

type blockReq struct {
	start, end uint64
	err        error
	retry      int
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

	var (
		subscribeStart = time.NewTicker(time.Second * 1)
		isSubError     bool
		latestHeight   = p.latestHeight(ctx)
		concurrency    = p.GetConcurrency(ctx, startHeight, latestHeight)
		resetFunc      = func() {
			isSubError = true
			subscribeStart.Reset(time.Second * 3)
			client, err := p.client.Reconnect()
			if err != nil {
				p.log.Error("failed to reconnect", zap.Error(err))
			} else {
				p.client = client
			}
		}
	)

	for {
		select {
		case <-ctx.Done():
			p.log.Debug("evm listener: done")
			return nil
		case <-subscribeStart.C:
			subscribeStart.Stop()
			go p.Subscribe(ctx, blockInfoChan, resetFunc)

			if isSubError {
				startHeight = p.GetLastSavedBlockHeight()
			}

			var blockReqs []*blockReq
			for start := startHeight; start <= latestHeight; start += p.cfg.BlockBatchSize {
				end := min(start+p.cfg.BlockBatchSize-1, latestHeight)
				blockReqs = append(blockReqs, &blockReq{start, end, nil, maxBlockQueryFailedRetry})
			}
			totalReqs := len(blockReqs)
			// Calculate the size of each chunk
			chunkSize := (totalReqs + concurrency - 1) / concurrency

			var wg sync.WaitGroup

			for i := 0; i < totalReqs; i += chunkSize {
				wg.Add(1)

				go func(blockReqsChunk []*blockReq, wg *sync.WaitGroup) {
					defer wg.Done()
					for _, br := range blockReqsChunk {
						filter := ethereum.FilterQuery{
							FromBlock: new(big.Int).SetUint64(br.start),
							ToBlock:   new(big.Int).SetUint64(br.end),
							Addresses: p.blockReq.Addresses,
							Topics:    p.blockReq.Topics,
						}
						p.log.Info("syncing", zap.Uint64("start", br.start), zap.Uint64("end", br.end), zap.Uint64("latest", latestHeight))
						logs, err := p.getLogsRetry(ctx, filter, br.retry)
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
				}(blockReqs[i:min(i+chunkSize, totalReqs)], &wg)
			}
			go func() {
				wg.Wait()
			}()
		}
	}
}

func (p *Provider) getLogsRetry(ctx context.Context, filter ethereum.FilterQuery, retry int) ([]ethTypes.Log, error) {
	var logs []ethTypes.Log
	var err error
	for i := 0; i < retry; i++ {
		logs, err = p.client.FilterLogs(ctx, filter)
		if err == nil {
			return logs, nil
		}
		p.log.Error("failed to get logs", zap.Error(err), zap.Int("retry", i+1))
		time.Sleep(time.Second * 15)
	}
	return nil, err
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
func (p *Provider) Subscribe(ctx context.Context, blockInfoChan chan *relayertypes.BlockInfo, resetFunc func()) error {
	ch := make(chan ethTypes.Log, 10)
	sub, err := p.client.Subscribe(ctx, ethereum.FilterQuery{
		Addresses: p.blockReq.Addresses,
		Topics:    p.blockReq.Topics,
	}, ch)
	if err != nil {
		p.log.Error("failed to subscribe", zap.Error(err))
		resetFunc()
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
			resetFunc()
			return err
		case log := <-ch:
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
			ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
			defer cancel()
			if _, err := p.client.GetHeaderByHeight(ctx, big.NewInt(1)); err != nil {
				p.log.Error("connection error", zap.Error(err))
				resetFunc()
				return err
			}
		}
	}
}
