package sui

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/icon-project/centralized-relay/relayer/chains/sui/types"
	relayerEvents "github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/sorter"
	"go.uber.org/zap"
)

func (p *Provider) Listener(ctx context.Context, lastSavedCheckpointSeq uint64, blockInfo chan *relayertypes.BlockInfo) error {
	latestCheckpointSeq, err := p.client.GetLatestCheckpointSeq(ctx)
	if err != nil {
		return err
	}

	startCheckpointSeq := latestCheckpointSeq
	if lastSavedCheckpointSeq != 0 && lastSavedCheckpointSeq < latestCheckpointSeq {
		startCheckpointSeq = lastSavedCheckpointSeq
	}

	return p.listenByPolling(ctx, startCheckpointSeq, blockInfo)
}

func (p *Provider) listenByPolling(ctx context.Context, startCheckpointSeq uint64, blockStream chan *relayertypes.BlockInfo) error {
	done := make(chan interface{})
	defer close(done)

	txDigestsStream := p.getTxDigestsStream(done, strconv.Itoa(int(startCheckpointSeq)-1))

	p.log.Info("Started to query sui from", zap.Uint64("checkpoint", startCheckpointSeq))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case txDigests, ok := <-txDigestsStream:
			if ok {
				p.log.Debug("executing query",
					zap.Any("from-checkpoint", txDigests.FromCheckpoint),
					zap.Any("to-checkpoint", txDigests.ToCheckpoint),
					zap.Any("tx-digests", txDigests.Digests),
				)

				eventResponse, err := p.client.GetEventsFromTxBlocks(ctx, p.allowedEventTypes(), txDigests.Digests)
				if err != nil {
					p.log.Error("failed to query events", zap.Error(err))
				}

				blockInfoList, err := p.parseMessagesFromEvents(eventResponse)
				if err != nil {
					p.log.Error("failed to parse messages from events", zap.Error(err))
				}

				for _, blockMsg := range blockInfoList {
					blockStream <- &blockMsg
				}
			}
		}
	}
}

func (p *Provider) allowedEventTypes() []string {
	return []string{
		fmt.Sprintf("%s::%s::%s", p.cfg.XcallPkgID, ModuleConnection, "Message"),
		fmt.Sprintf("%s::%s::%s", p.cfg.XcallPkgID, ModuleMain, "CallMessage"),
		fmt.Sprintf("%s::%s::%s", p.cfg.XcallPkgID, ModuleMain, "RollbackMessage"),
	}
}

func (p *Provider) parseMessagesFromEvents(events []types.EventResponse) ([]relayertypes.BlockInfo, error) {
	checkpointMessages := make(map[uint64][]*relayertypes.Message)
	for _, ev := range events {
		msg, err := p.parseMessageFromEvent(ev)
		if err != nil {
			return nil, err
		}
		p.log.Info("Detected event log: ",
			zap.Uint64("checkpoint", msg.MessageHeight),
			zap.String("event-type", msg.EventType),
			zap.Uint64("sn", msg.Sn),
			zap.String("dst", msg.Dst),
			zap.Uint64("req-id", msg.ReqID),
			zap.Any("data", hex.EncodeToString(msg.Data)),
		)
		checkpointMessages[ev.Checkpoint] = append(checkpointMessages[ev.Checkpoint], msg)
	}

	var blockInfoList []relayertypes.BlockInfo
	for checkpoint, messages := range checkpointMessages {
		blockInfoList = append(blockInfoList, relayertypes.BlockInfo{
			Height:   checkpoint,
			Messages: messages,
		})
	}

	sorter.Sort(blockInfoList, func(bi1, bi2 relayertypes.BlockInfo) bool {
		return bi1.Height < bi2.Height //ascending order
	})

	return blockInfoList, nil
}

func (p *Provider) parseMessageFromEvent(ev types.EventResponse) (*relayertypes.Message, error) {
	msg := relayertypes.Message{
		MessageHeight: ev.Checkpoint,
		Src:           p.cfg.NID,
	}

	eventBytes, err := json.Marshal(ev.ParsedJson)
	if err != nil {
		return nil, err
	}

	switch ev.Type {
	case fmt.Sprintf("%s::%s::%s", p.cfg.XcallPkgID, ModuleConnection, "Message"):
		msg.EventType = relayerEvents.EmitMessage
		var emitEvent types.EmitEvent
		if err := json.Unmarshal(eventBytes, &emitEvent); err != nil {
			return nil, err
		}
		sn, err := strconv.Atoi(emitEvent.Sn)
		if err != nil {
			return nil, err
		}
		msg.Sn = uint64(sn)
		msg.Data = emitEvent.Msg
		msg.Dst = emitEvent.To

	case fmt.Sprintf("%s::%s::%s", p.cfg.XcallPkgID, ModuleMain, "CallMessage"):
		msg.EventType = relayerEvents.CallMessage
		var callMsgEvent types.CallMsgEvent
		if err := json.Unmarshal(eventBytes, &callMsgEvent); err != nil {
			return nil, err
		}
		msg.Data = callMsgEvent.Data
		reqID, err := strconv.Atoi(callMsgEvent.ReqId)
		if err != nil {
			return nil, err
		}
		msg.ReqID = uint64(reqID)
		msg.DappModuleCapID = callMsgEvent.DappModuleCapId
		msg.Dst = p.cfg.NID

	case fmt.Sprintf("%s::%s::%s", p.cfg.XcallPkgID, ModuleMain, "RollbackMessage"):
		msg.EventType = relayerEvents.ExecuteRollback
		var rollbackMsgEvent types.RollbackMsgEvent
		if err := json.Unmarshal(eventBytes, &rollbackMsgEvent); err != nil {
			return nil, err
		}
		sn, err := strconv.Atoi(rollbackMsgEvent.Sn)
		if err != nil {
			return nil, err
		}
		msg.Sn = uint64(sn)
		msg.DappModuleCapID = rollbackMsgEvent.DappModuleCapId
		msg.Dst = p.cfg.NID

	default:
		return nil, fmt.Errorf("invalid event type")
	}

	msg.Src = p.cfg.NID

	return &msg, nil
}

// GenerateTxDigests forms the packets of txDigests from the list of checkpoint responses such that each packet
// contains as much as possible number of digests but not exceeding max limit of maxDigests value
func (p *Provider) GenerateTxDigests(checkpointResponses []types.CheckpointResponse, maxDigestsPerItem int) []types.TxDigests {
	// stage-1: split checkpoint to multiple checkpoints if number of transactions is greater than maxDigests
	var checkpoints []types.CheckpointResponse
	for _, cp := range checkpointResponses {
		if len(cp.Transactions) > maxDigestsPerItem {
			totalBatches := len(cp.Transactions) / maxDigestsPerItem
			if (len(cp.Transactions) % maxDigestsPerItem) != 0 {
				totalBatches = totalBatches + 1
			}
			for i := 0; i < totalBatches; i++ {
				fromIndex := i * maxDigestsPerItem
				toIndex := fromIndex + maxDigestsPerItem
				if i == totalBatches-1 {
					toIndex = len(cp.Transactions)
				}
				subCheckpoint := types.CheckpointResponse{
					SequenceNumber: cp.SequenceNumber,
					Transactions:   cp.Transactions[fromIndex:toIndex],
				}
				checkpoints = append(checkpoints, subCheckpoint)
			}
		} else {
			checkpoints = append(checkpoints, cp)
		}
	}

	// stage-2: form packets of txDigests
	var txDigestsList []types.TxDigests

	digests := []string{}
	fromCheckpoint, _ := strconv.Atoi(checkpoints[0].SequenceNumber)
	for i, cp := range checkpoints {
		if (len(digests) + len(cp.Transactions)) > maxDigestsPerItem {
			toCheckpoint, _ := strconv.Atoi(checkpoints[i-1].SequenceNumber)
			if len(digests) < maxDigestsPerItem {
				toCheckpoint, _ = strconv.Atoi(cp.SequenceNumber)
			}
			for i, tx := range cp.Transactions {
				if len(digests) == maxDigestsPerItem {
					txDigestsList = append(txDigestsList, types.TxDigests{
						FromCheckpoint: uint64(fromCheckpoint),
						ToCheckpoint:   uint64(toCheckpoint),
						Digests:        digests,
					})
					digests = cp.Transactions[i:]
					fromCheckpoint, _ = strconv.Atoi(cp.SequenceNumber)
					break
				} else {
					digests = append(digests, tx)
				}
			}
		} else {
			digests = append(digests, cp.Transactions...)
		}
	}

	lastCheckpointSeq := checkpoints[len(checkpoints)-1].SequenceNumber
	lastCheckpoint, _ := strconv.Atoi(lastCheckpointSeq)
	txDigestsList = append(txDigestsList, types.TxDigests{
		FromCheckpoint: uint64(fromCheckpoint),
		ToCheckpoint:   uint64(lastCheckpoint),
		Digests:        digests,
	})

	return txDigestsList
}

func (p *Provider) getTxDigestsStream(done chan interface{}, afterSeq string) <-chan types.TxDigests {
	txDigestsStream := make(chan types.TxDigests)

	go func() {
		nextCursor := afterSeq
		checkpointTicker := time.NewTicker(3 * time.Second) //todo need to decide this interval

		for {
			select {
			case <-done:
				return
			case <-checkpointTicker.C:
				req := types.SuiGetCheckpointsRequest{
					Cursor:          nextCursor,
					Limit:           types.QUERY_MAX_RESULT_LIMIT,
					DescendingOrder: false,
				}
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				paginatedRes, err := p.client.GetCheckpoints(ctx, req)
				if err != nil {
					p.log.Error("failed to fetch checkpoints", zap.Error(err))
					continue
				}

				if len(paginatedRes.Data) > 0 {
					for _, txDigests := range p.GenerateTxDigests(paginatedRes.Data, types.QUERY_MAX_RESULT_LIMIT) {
						txDigestsStream <- types.TxDigests{
							FromCheckpoint: uint64(txDigests.FromCheckpoint),
							ToCheckpoint:   uint64(txDigests.ToCheckpoint),
							Digests:        txDigests.Digests,
						}
					}

					nextCursor = paginatedRes.Data[len(paginatedRes.Data)-1].SequenceNumber
				}
			}
		}
	}()

	return txDigestsStream
}
