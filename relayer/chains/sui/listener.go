package sui

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/coming-chat/go-sui/v2/sui_types"
	cctypes "github.com/coming-chat/go-sui/v2/types"
	"github.com/icon-project/centralized-relay/relayer/chains/sui/types"
	relayerEvents "github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/sorter"
	"go.uber.org/zap"
)

func (p *Provider) Listener(ctx context.Context, lastProcessedTx relayertypes.LastProcessedTx, blockInfo chan *relayertypes.BlockInfo) error {
	txInfo := new(types.TxInfo)

	if lastProcessedTx.Info != nil {
		if err := txInfo.Deserialize(lastProcessedTx.Info); err != nil {
			p.log.Error("failed to deserialize last processed tx digest", zap.Error(err))
			return err
		}
	}

	if p.cfg.StartTxDigest != "" {
		txInfo.TxDigest = p.cfg.StartTxDigest
	}

	if txInfo.TxDigest == "" {
		latestTx, err := p.getLatestXcallTransaction()
		if err != nil {
			p.log.Error("failed to get latest xcall transaction", zap.Error(err))
			return err
		}
		if latestTx != nil {
			txInfo.TxDigest = latestTx.Digest.String()
		}
	}

	return p.listenByPolling(ctx, txInfo.TxDigest, blockInfo)
}

func (p *Provider) getLatestXcallTransaction() (*cctypes.SuiTransactionBlockResponse, error) {
	inputObj, err := sui_types.NewObjectIdFromHex(p.cfg.XcallStorageID)
	if err != nil {
		return nil, err
	}
	query := cctypes.SuiTransactionBlockResponseQuery{
		Filter: &cctypes.TransactionFilter{
			InputObject: inputObj,
		},
	}
	limit := uint(1)

	res, err := p.client.QueryTxBlocks(context.Background(), query, nil, &limit, true)
	if err != nil {
		return nil, err
	}

	if len(res.Data) > 0 {
		return &res.Data[0], nil
	}

	return nil, nil
}

func (p *Provider) shouldSkipMessage(msg *relayertypes.Message) bool {
	// if relayer is not an executor then skip CallMessage and RollbackMessage events.
	if len(p.cfg.Dapps) == 0 &&
		(msg.EventType == relayerEvents.CallMessage ||
			msg.EventType == relayerEvents.RollbackMessage) {
		return true
	}
	return false
}

func (p *Provider) parseMessagesFromEvents(events []types.EventResponse) ([]relayertypes.BlockInfo, error) {
	checkpointMessages := make(map[uint64][]*relayertypes.Message)
	for _, ev := range events {
		msg, err := p.parseMessageFromEvent(ev)
		if err != nil {
			if err.Error() == types.InvalidEventError {
				continue
			}
			return nil, err
		}

		if p.shouldSkipMessage(msg) {
			continue
		}

		p.log.Info("Detected event log: ",
			zap.Uint64("checkpoint", msg.MessageHeight),
			zap.String("event_type", msg.EventType),
			zap.Any("sn", msg.Sn),
			zap.String("dst", msg.Dst),
			zap.Any("req_id", msg.ReqID),
			zap.String("tx_hash", ev.Id.TxDigest.String()),
		)
		checkpointMessages[ev.Checkpoint.Uint64()] = append(checkpointMessages[ev.Checkpoint.Uint64()], msg)
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
		MessageHeight: ev.Checkpoint.Uint64(),
		Src:           p.cfg.NID,
	}

	txInfo := types.TxInfo{TxDigest: ev.Id.TxDigest.String()}
	txInfoBytes, err := txInfo.Serialize()
	if err != nil {
		return nil, err
	}
	msg.TxInfo = txInfoBytes

	eventBytes, err := json.Marshal(ev.ParsedJson)
	if err != nil {
		return nil, err
	}

	eventParts := strings.Split(ev.Type, "::")
	eventType := eventParts[len(eventParts)-1]

	switch eventType {
	case "Message":
		msg.EventType = relayerEvents.EmitMessage
		var emitEvent types.EmitEvent
		if err := json.Unmarshal(eventBytes, &emitEvent); err != nil {
			return nil, err
		}
		if emitEvent.ConnectionID != p.cfg.ConnectionID {
			return nil, fmt.Errorf(types.InvalidEventError)
		}
		sn, ok := new(big.Int).SetString(emitEvent.Sn, 10)
		if !ok {
			return nil, err
		}
		msg.Sn = sn
		msg.Data = emitEvent.Msg
		msg.Src = p.cfg.NID
		msg.Dst = emitEvent.To

	case "CallMessage":
		msg.EventType = relayerEvents.CallMessage
		var callMsgEvent types.CallMsgEvent
		if err := json.Unmarshal(eventBytes, &callMsgEvent); err != nil {
			return nil, err
		}
		msg.Src = callMsgEvent.From.NetID
		msg.Data = callMsgEvent.Data

		sn, ok := new(big.Int).SetString(callMsgEvent.Sn, 10)
		if !ok {
			return nil, err
		}

		reqID, ok := new(big.Int).SetString(callMsgEvent.ReqId, 10)
		if !ok {
			return nil, err
		}

		msg.Sn = sn
		msg.ReqID = reqID
		msg.DappModuleCapID = callMsgEvent.DappModuleCapId
		msg.Dst = p.cfg.NID

	case "RollbackMessage":
		msg.EventType = relayerEvents.RollbackMessage
		var rollbackMsgEvent types.RollbackMsgEvent
		if err := json.Unmarshal(eventBytes, &rollbackMsgEvent); err != nil {
			return nil, err
		}

		sn, ok := new(big.Int).SetString(rollbackMsgEvent.Sn, 10)
		if !ok {
			return nil, fmt.Errorf("failed to parse sn from rollback event")
		}
		msg.Sn = sn
		msg.DappModuleCapID = rollbackMsgEvent.DappModuleCapId
		msg.Src = p.cfg.NID
		msg.Dst = p.cfg.NID
		msg.Data = rollbackMsgEvent.Data

	default:
		return nil, fmt.Errorf(types.InvalidEventError)
	}

	return &msg, nil
}

func (p *Provider) handleEventNotification(ctx context.Context, ev types.EventResponse, blockStream chan *relayertypes.BlockInfo) {
	for ev.Checkpoint == nil {
		p.log.Warn("checkpoint not found for transaction", zap.String("tx-digest", ev.Id.TxDigest.String()))
		time.Sleep(2 * time.Second)
		txRes, err := p.client.GetTransaction(ctx, ev.Id.TxDigest.String())
		if err != nil {
			p.log.Error("failed to get transaction while handling event notification",
				zap.Error(err), zap.Any("event", ev))
			return
		}
		ev.Checkpoint = txRes.Checkpoint
	}

	msg, err := p.parseMessageFromEvent(ev)
	if err != nil {
		if err.Error() == types.InvalidEventError {
			return
		}
		p.log.Error("failed to parse message from event while handling event notification",
			zap.Error(err),
			zap.Any("event", ev))
		return
	}

	if p.shouldSkipMessage(msg) {
		return
	}

	p.log.Info("Detected event log: ",
		zap.Uint64("checkpoint", msg.MessageHeight),
		zap.String("event_type", msg.EventType),
		zap.Any("sn", msg.Sn),
		zap.String("dst", msg.Dst),
		zap.Any("req_id", msg.ReqID),
		zap.String("tx_hash", ev.Id.TxDigest.String()),
	)

	blockStream <- &relayertypes.BlockInfo{
		Height:   msg.MessageHeight,
		Messages: []*relayertypes.Message{msg},
	}
}

func (p *Provider) listenByPolling(ctx context.Context, startTxDigestStr string, blockStream chan *relayertypes.BlockInfo) error {
	done := make(chan interface{})
	defer close(done)

	var startTxDigest *sui_types.TransactionDigest

	if startTxDigestStr != "" { //process probably unexplored events of last saved tx digest
		var err error
		startTxDigest, err = sui_types.NewDigest(startTxDigestStr)
		if err != nil {
			return err
		}
		currentEvents, err := p.client.GetEvents(ctx, *startTxDigest)
		if err != nil {
			return err
		}

		for _, ev := range currentEvents {
			p.handleEventNotification(ctx, types.EventResponse{SuiEvent: ev}, blockStream)
		}
	}

	eventStream := p.getObjectEventStream(done, p.cfg.XcallStorageID, startTxDigest)

	p.log.Info("event query started", zap.String("start-tx-digest", startTxDigestStr))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev, ok := <-eventStream:
			if ok {
				p.handleEventNotification(ctx, ev, blockStream)
			}
		}
	}
}

func (p *Provider) getObjectEventStream(done chan interface{}, objectID string, afterTxDigest *sui_types.TransactionDigest) <-chan types.EventResponse {
	eventStream := make(chan types.EventResponse)

	go func() {
		defer close(eventStream)

		inputObj, err := sui_types.NewObjectIdFromHex(objectID)
		if err != nil {
			p.log.Panic("failed to create object from hex string", zap.Error(err))
		}

		query := cctypes.SuiTransactionBlockResponseQuery{
			Filter: &cctypes.TransactionFilter{
				InputObject: inputObj,
			},
			Options: &cctypes.SuiTransactionBlockResponseOptions{
				ShowEvents: true,
			},
		}

		cursor := afterTxDigest

		limit := uint(15)

		pollInterval := 6 * time.Second
		if p.cfg.PollInterval != 0 {
			pollInterval = p.cfg.PollInterval
		}

		ticker := time.NewTicker(pollInterval)

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				res, err := p.client.QueryTxBlocks(context.Background(), query, cursor, &limit, false)
				if err != nil {
					p.log.Error("failed to query tx blocks", zap.Error(err), zap.Any("cursor", cursor))
					break
				}

				p.log.Debug("tx block query successful", zap.Any("cursor", cursor))

				for _, blockRes := range res.Data {
					for _, ev := range blockRes.Events {
						eventStream <- types.EventResponse{
							SuiEvent:   ev,
							Checkpoint: blockRes.Checkpoint,
						}
					}
				}

				cursor = res.NextCursor
			}
		}
	}()

	return eventStream
}
