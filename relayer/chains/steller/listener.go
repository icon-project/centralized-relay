package steller

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"runtime"
	"time"

	"github.com/icon-project/centralized-relay/relayer/chains/steller/types"
	relayerevents "github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/concurrency"
	"github.com/icon-project/centralized-relay/utils/sorter"
	xdr "github.com/stellar/go-xdr/xdr3"
	"go.uber.org/zap"
)

func (p *Provider) Listener(ctx context.Context, lastSavedLedgerSeq uint64, blockInfo chan *relayertypes.BlockInfo) error {
	go func() { //Todo remove: used temporarily for testing purpose only
		time.Sleep(5 * time.Second)
		if err := p.RestoreKeystore(ctx); err != nil {
			p.log.Error("error restoring keystore: ", zap.Error(err))
		}
		if err := p.Route(ctx, &relayertypes.Message{
			Dst:  "icon",
			Data: []byte("hello"),
		}, func(key *relayertypes.MessageKey, response *relayertypes.TxResponse, err error) {
			if err != nil {
				p.log.Info("message relay failed", zap.String("src", "steller"), zap.String("dst", "icon"), zap.Int64("height", response.Height), zap.String("hash", response.TxHash), zap.Error(err))
			} else {
				p.log.Info("message relay successfull", zap.String("src", "steller"), zap.String("dst", "icon"), zap.Int64("height", response.Height), zap.String("hash", response.TxHash))
			}
		}); err != nil {
			p.log.Error("error sending tx: ", zap.Error(err))
		}
	}()

	latestLedger, err := p.client.GetLatestLedger(ctx)
	if err != nil {
		return err
	}

	latestSeq := latestLedger.Sequence

	startSeq := latestSeq
	if lastSavedLedgerSeq != 0 && lastSavedLedgerSeq < latestSeq {
		startSeq = lastSavedLedgerSeq
	}

	blockIntervalTicker := time.NewTicker(p.cfg.BlockInterval)
	defer blockIntervalTicker.Stop()

	p.log.Info("start querying from ledger seq", zap.Uint64("start-seq", startSeq))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-blockIntervalTicker.C:
			newLatestLedger, err := p.client.GetLatestLedger(ctx)
			if err != nil {
				p.log.Error("failed to query latest ledger", zap.Error(err))
			} else if newLatestLedger.Sequence > latestSeq {
				latestSeq = newLatestLedger.Sequence
			}
		default:
			if startSeq < latestSeq {
				p.log.Info("Query started.", zap.Uint64("from-seq", startSeq), zap.Uint64("to-seq", latestSeq))
				for _, seqBatch := range getSeqBatches(startSeq, latestSeq, types.LedgerSeqBatchSize) {
					p.runLedgerQuery(blockInfo, seqBatch.FromSeq, seqBatch.ToSeq)
					p.log.Info("Query completed.", zap.Uint64("from-seq", seqBatch.FromSeq), zap.Uint64("to-seq", seqBatch.ToSeq))
				}
				startSeq = latestSeq + 1
			}
		}
	}
}

func (p *Provider) runLedgerQuery(blockInfoChan chan *relayertypes.BlockInfo, fromSeq, toSeq uint64) {
	done := make(chan interface{})
	defer close(done)

	seqStream := getLedgerSeqStream(done, fromSeq, toSeq)

	numOfPipelines := int(toSeq - fromSeq + 1)
	if numOfPipelines > runtime.NumCPU() {
		numOfPipelines = runtime.NumCPU()
	}

	pipelines := make([]<-chan interface{}, numOfPipelines)

	for i := 0; i < numOfPipelines; i++ {
		pipelines[i] = p.getLedgerInfoStream(done, seqStream)
	}

	var blockInfoList []relayertypes.BlockInfo
	for bn := range concurrency.Take(done, concurrency.FanIn(done, pipelines...), int(toSeq-fromSeq+1)) {
		if bn != nil {
			block := bn.(relayertypes.BlockInfo)
			blockInfoList = append(blockInfoList, block)
		}
	}

	sorter.Sort(
		blockInfoList,
		func(p1, p2 relayertypes.BlockInfo) bool {
			return p1.Height < p2.Height //ascending order
		},
	)

	for _, blockInfo := range blockInfoList {
		blockInfoChan <- &relayertypes.BlockInfo{
			Height: blockInfo.Height, Messages: blockInfo.Messages,
		}
	}
}

func getSeqBatches(fromSeq, toSeq, batchSize uint64) []types.LedgerSeqBatch {
	seqBatches := []types.LedgerSeqBatch{}
	seq := fromSeq
	for seq+batchSize <= toSeq {
		seqBatches = append(seqBatches, types.LedgerSeqBatch{FromSeq: seq, ToSeq: seq + batchSize - 1})
		seq = seq + batchSize
	}
	if seq <= toSeq {
		seqBatches = append(seqBatches, types.LedgerSeqBatch{FromSeq: seq, ToSeq: toSeq})
	}
	return seqBatches
}

func getLedgerSeqStream(done <-chan interface{}, fromSeq, toSeq uint64) <-chan uint64 {
	seqStream := make(chan uint64)
	seq := fromSeq
	go func() {
		defer close(seqStream)
		for seq <= toSeq {
			select {
			case <-done:
				return
			default:
				seqStream <- seq
				seq++
			}
		}
	}()
	return seqStream
}

func (p *Provider) getLedgerInfoStream(done <-chan interface{}, seqStream <-chan uint64) <-chan interface{} {
	ledgerInfoStream := make(chan interface{})
	go func() {
		defer close(ledgerInfoStream)
	Loop:
		for {
			select {
			case <-done:
				return
			case seq, ok := <-seqStream:
				if ok {
					for { // will block until and unless ledger messages are fetched so that we are not skipping/missing this ledger seq.
						messages, err := p.fetchLedgerMessages(context.Background(), seq)
						if err != nil {
							p.log.Error("failed to fetch ledger messages: ", zap.Error(err), zap.Uint64("ledger seq", seq))
							time.Sleep(1 * time.Second)
						} else {
							ledgerInfoStream <- relayertypes.BlockInfo{
								Height:   seq,
								Messages: messages,
							}
							break
						}
					}
				} else {
					break Loop // break out of the outer loop
				}
			}
		}
	}()
	return ledgerInfoStream
}

func (p *Provider) fetchLedgerMessages(ctx context.Context, ledgerSeq uint64) ([]*relayertypes.Message, error) {
	eventFilter := types.EventFilter{
		LedgerSeq:   ledgerSeq,
		ContractIds: []string{p.cfg.Contracts[relayertypes.ConnectionContract]},
		Topics:      []string{"new_message"},
	}
	events, err := p.client.FetchEvents(ctx, eventFilter)
	if err != nil {
		return nil, err
	}

	messages, err := p.parseMessagesFromEvents(events)
	for _, msg := range messages {
		p.log.Info("detected event log:", zap.Any("event", *msg))
	}
	p.log.Debug("query successful", zap.Uint64("ledger-seq", ledgerSeq))
	return messages, err
}

func (p *Provider) parseMessagesFromEvents(events []types.Event) ([]*relayertypes.Message, error) {
	messages := []*relayertypes.Message{}
	for _, ev := range events {
		var eventType string
		for _, topic := range ev.Body.V0.Topics {
			switch topic.String() {
			case "new_message": //used only for testing; need to remove
				eventType = "new_message"
			case "emitMessage":
				eventType = relayerevents.EmitMessage
			case "callMessage":
				eventType = relayerevents.CallMessage
			}
		}

		// if event type is not matched then skip this event
		if eventType == "" {
			continue
		}

		msg := &relayertypes.Message{
			EventType:     eventType,
			MessageHeight: ev.LedgerSeq,
		}

		scval := ev.Body.V0.Data
		scMap, ok := scval.GetMap()
		if !ok {
			continue
		}

		for _, mapItem := range *scMap {
			valBytes, err := mapItem.Val.MarshalBinary()
			if err != nil {
				return nil, err
			}
			decoder := xdr.NewDecoder(bytes.NewBuffer(valBytes))
			switch mapItem.Key.String() {
			case "sn":
				intVal, _, err := decoder.DecodeInt()
				if err != nil {
					return nil, fmt.Errorf("failed to decode sn: %v", err)
				}
				msg.Sn = uint64(intVal)
			case "reqId":
				intVal, _, err := decoder.DecodeInt()
				if err != nil {
					return nil, fmt.Errorf("failed to decode reqId: %v", err)
				}
				msg.ReqID = uint64(intVal)
			case "src":
				msg.Src = mapItem.Val.String()
			case "dst":
				msg.Dst = mapItem.Val.String()
			case "data":
				data, err := hex.DecodeString(mapItem.Val.String())
				if err != nil {
					return nil, fmt.Errorf("failed to decode data: %v", err)
				}
				msg.Data = data
			}
		}

		//skip invalid message
		if msg.Sn == 0 || msg.Src == "" || msg.Dst == "" {
			p.log.Warn("detected invalid message: ", zap.Any("msg", msg))
			continue
		}

		messages = append(messages, msg)
	}

	return messages, nil
}
