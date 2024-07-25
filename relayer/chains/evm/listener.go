package evm

import (
	"context"
	"math/big"
	"runtime"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/icon-project/centralized-relay/relayer/chains/evm/types"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/pkg/errors"
)

const (
	defaultReadTimeout         = 60 * time.Second
	monitorBlockMaxConcurrency = 10 // number of concurrent requests to synchronize older blocks from source chain
	maxBlockRange              = 50
	maxBlockQueryFailedRetry   = 3
	DefaultFinalityBlock       = 10
	pollerTime                 = 5 * time.Second
	PollFetch                  = "poll"
	WsFetch                    = "ws"
	RPCRedundancy              = "rpc-verify"
	WsRedundancy               = "ws-verify"
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
	height, err := r.GetClient().GetBlockNumber(ctx)
	if err != nil {
		r.log.Error("Evm listener: failed to GetBlockNumber", zap.Error(err))
		return 0
	}
	return height
}

func GetBlockInfoKey(bInfo *relayertypes.BlockInfo) string {
	return bInfo.Messages[0].Src + "-" + bInfo.Messages[0].Dst +
		"-" + bInfo.Messages[0].EventType + "-" + bInfo.Messages[0].Sn.String()
}

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfoChan chan *relayertypes.BlockInfo) error {
	startHeight, err := p.startFromHeight(ctx, lastSavedHeight)
	if err != nil {
		return err
	}

	p.log.Info("Start from height ", zap.Uint64("height", startHeight), zap.Uint64("finality block", p.FinalityBlock(ctx)))
	if p.cfg.Redundancy == WsRedundancy {
		return p.listenWithWsRedundancy(ctx, startHeight, blockInfoChan)
	} else {
		return p.listenNormalWsNPoll(ctx, startHeight, blockInfoChan)
	}
}

func (p *Provider) listenWithWsRedundancy(ctx context.Context, startHeight uint64, blockInfoChan chan *relayertypes.BlockInfo) error {
	cache := expirable.NewLRU[string, any](1000, nil, time.Hour*6)
	internalChan := make(chan *relayertypes.BlockInfo)
	var (
		subscribeStart = time.NewTicker(time.Second * 1)
		errChan        = make(chan types.ErrorMessageRpc)
	)
	var restartLists []string
	for {
		select {
		case <-ctx.Done():
			p.log.Debug("evm listener: done")
			return nil
		case pkt := <-internalChan:
			cacheKey := GetBlockInfoKey(pkt)
			if cache.Contains(cacheKey) {
				cache.Remove(cacheKey)
				blockInfoChan <- pkt
			} else {
				cache.Add(cacheKey, nil)
			}
		case err := <-errChan:
			if p.isConnectionError(err.Error) {
				p.log.Error("connection error", zap.Error(err.Error))
				index := -1
				for lindex, lclient := range p.clients {
					if lclient.GetRPCUrl() == err.RPCUrl {
						index = lindex
					}
				}
				if index != -1 {
					nclient, cnErr := p.clients[index].Reconnect()
					if cnErr != nil {
						p.log.Error("failed to reconnect", zap.Error(cnErr))
					} else {
						p.log.Info("client reconnected")
						p.clients[index] = nclient
					}
					restartLists = append(restartLists, p.clients[index].GetRPCUrl())
				}
			} else {
				index := -1
				for lindex, lclient := range p.clients {
					if lclient.GetRPCUrl() == err.RPCUrl {
						index = lindex
					}
				}
				restartLists = append(restartLists, p.clients[index].GetRPCUrl())
			}
			startHeight = p.GetLastSavedBlockHeight()
			subscribeStart.Reset(time.Second * 1)
		case <-subscribeStart.C:
			subscribeStart.Stop()
			latestHeight := p.latestHeight(ctx)
			restartIndex := -1
			for _, client := range p.clients {
				if len(restartLists) == 0 {
					go p.Subscribe(ctx, client, internalChan, errChan, latestHeight)
				}
				for idx, rpc := range restartLists {
					if rpc == client.GetRPCUrl() {
						restartIndex = idx
						go p.Subscribe(ctx, client, internalChan, errChan, latestHeight)
					}
				}
			}
			restartLists = restartLists[:0]
			if restartIndex == 1 || len(restartLists) == 0 {
				concurrency := p.GetConcurrency(ctx, startHeight, latestHeight)
				var blockReqs []*blockReq
				if startHeight == 0 {
					startHeight = latestHeight
				}
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
							p.log.Info("syncing", zap.Uint64("start", br.start), zap.Uint64("end", br.end), zap.Uint64("latest", latestHeight), zap.Any("host", p.GetClient().GetRPCUrl()))
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
			p.backlogProcessing = false
		}
	}
}

func (p *Provider) listenNormalWsNPoll(ctx context.Context, startHeight uint64, blockInfoChan chan *relayertypes.BlockInfo) error {
	var (
		subscribeStart = time.NewTicker(time.Second * 1)
		pollerStart    = time.NewTicker(pollerTime)
		errChan        = make(chan types.ErrorMessageRpc)
	)

	if p.cfg.Fetch == PollFetch {
		if p.cfg.Redundancy == RPCRedundancy {
			panic("cannot set rpc-verify on poll fetch strategy")
		}
		subscribeStart.Stop()
	} else {
		pollerStart.Stop()
	}

	for {
		select {
		case <-ctx.Done():
			p.log.Debug("evm listener: done")
			return nil
		case err := <-errChan:
			if p.isConnectionError(err.Error) {
				p.log.Error("connection error", zap.Error(err.Error))
				index := -1
				for lindex, lclient := range p.clients {
					if lclient.GetRPCUrl() == err.RPCUrl {
						index = lindex
					}
				}
				if index != -1 {
					nclient, cnErr := p.clients[index].Reconnect()
					if cnErr != nil {
						p.log.Error("failed to reconnect", zap.Error(cnErr))
					} else {
						p.log.Info("client reconnected")
						p.clients[index] = nclient
					}
				}
			}
			startHeight = p.GetLastSavedBlockHeight()
			subscribeStart.Reset(time.Second * 1)
		case <-subscribeStart.C:
			subscribeStart.Stop()
			latestHeight := p.latestHeight(ctx)
			for _, client := range p.clients {
				go p.Subscribe(ctx, client, blockInfoChan, errChan, latestHeight)
			}
			concurrency := p.GetConcurrency(ctx, startHeight, latestHeight)
			var blockReqs []*blockReq
			for start := startHeight; start <= latestHeight; start += p.cfg.BlockBatchSize {
				end := min(start+p.cfg.BlockBatchSize-1, latestHeight)
				blockReqs = append(blockReqs, &blockReq{start, end, nil, maxBlockQueryFailedRetry})
			}
			totalReqs := len(blockReqs)
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
			p.backlogProcessing = false
		case <-pollerStart.C:
			latestHeight := p.latestHeight(ctx)
			var blockReqs []*blockReq
			if startHeight < latestHeight {
				for start := startHeight; start <= latestHeight; start += p.cfg.BlockBatchSize {
					end := min(start+p.cfg.BlockBatchSize-1, latestHeight)
					blockReqs = append(blockReqs, &blockReq{start, end, nil, maxBlockQueryFailedRetry})
				}
				for _, br := range blockReqs {
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
				p.saveHeightFunc(latestHeight)
				startHeight = latestHeight
			}
			pollerStart.Reset(pollerTime)
		}
	}
}

func (p *Provider) getLogsRetry(ctx context.Context, filter ethereum.FilterQuery, retry int) ([]ethTypes.Log, error) {
	var logs []ethTypes.Log
	var err error
	for i := 0; i < retry; i++ {
		logs, err = p.GetClient().FilterLogs(ctx, filter)
		if err == nil {
			return logs, nil
		}
		p.log.Error("failed to get logs", zap.Error(err), zap.Int("retry", i+1))
		time.Sleep(time.Second * 30)
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
func (p *Provider) Subscribe(ctx context.Context, client IClient,
	blockInfoChan chan *relayertypes.BlockInfo, resetCh chan types.ErrorMessageRpc, latestHeight uint64) error {
	ch := make(chan ethTypes.Log, 10)
	wsEventsFound := false
	sub, err := client.Subscribe(ctx, ethereum.FilterQuery{
		Addresses: p.blockReq.Addresses,
		Topics:    p.blockReq.Topics,
		//TODO: required for some rpcs, review
		// FromBlock: new(big.Int).SetUint64(latestHeight),
	}, ch)
	if err != nil {
		p.log.Error("failed to subscribe", zap.Error(err), zap.Any("host", client.GetRPCUrl()))
		resetCh <- types.ErrorMessageRpc{Error: err, RPCUrl: client.GetRPCUrl()}
		return err
	}
	defer sub.Unsubscribe()
	defer close(ch)
	p.log.Info("Subscribed to new blocks", zap.Any("address", p.blockReq.Addresses), zap.Any("rpc", client.GetRPCUrl()))
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-sub.Err():
			p.log.Error("subscription error", zap.Error(err))
			resetCh <- types.ErrorMessageRpc{Error: err, RPCUrl: client.GetRPCUrl()}
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
				zap.Any("rpc", client.GetRPCUrl()),
			)
			blockInfoChan <- &relayertypes.BlockInfo{
				Height:   log.BlockNumber,
				Messages: []*relayertypes.Message{message},
			}
		case <-time.After(time.Minute * 2):
			if _, err := client.GetHeaderByHeight(ctx, big.NewInt(1)); err != nil {
				resetCh <- types.ErrorMessageRpc{Error: err, RPCUrl: client.GetRPCUrl()}
				return err
			}
			if p.cfg.Redundancy == RPCRedundancy {
				p.handleRpcVerification(ctx, wsEventsFound, latestHeight, resetCh, blockInfoChan)
			}
		}

	}
}

func (p *Provider) handleRpcVerification(ctx context.Context, wsEventsFound bool, latestHeight uint64,
	resetCh chan types.ErrorMessageRpc, blockInfoChan chan *relayertypes.BlockInfo) {
	if !wsEventsFound {
		p.log.Info("No events found in ws,verifying rpc")
		currentLatestHeight := p.latestHeight(ctx)
		var blockReqs []*blockReq
		if latestHeight < currentLatestHeight {
			for start := latestHeight; start <= currentLatestHeight; start += p.cfg.BlockBatchSize {
				end := min(start+p.cfg.BlockBatchSize-1, currentLatestHeight)
				blockReqs = append(blockReqs, &blockReq{start, end, nil, maxBlockQueryFailedRetry})
			}
			if len(blockReqs) > 0 {
				p.log.Info("Need to reset ws connections")
				resetCh <- types.ErrorMessageRpc{Error: errors.New("ws stale connection"), RPCUrl: p.GetClient().GetRPCUrl()}
			}
			for _, br := range blockReqs {
				filter := ethereum.FilterQuery{
					FromBlock: new(big.Int).SetUint64(br.start),
					ToBlock:   new(big.Int).SetUint64(br.end),
					Addresses: p.blockReq.Addresses,
					Topics:    p.blockReq.Topics,
				}
				p.log.Info("syncing missing", zap.Uint64("start", br.start), zap.Uint64("end", br.end), zap.Uint64("latest", latestHeight))
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
					p.log.Info("Detected missing eventlog",
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
		}
		p.saveHeightFunc(currentLatestHeight)
	}
}
