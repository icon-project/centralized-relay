package sui

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/sui_types"
	cctypes "github.com/coming-chat/go-sui/v2/types"
	"github.com/icon-project/centralized-relay/relayer/chains/sui/types"
	relayerEvents "github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/sorter"
	"go.uber.org/zap"
)

func (p *Provider) Listener(ctx context.Context, lastProcessedTxDigest interface{}, blockInfo chan *relayertypes.BlockInfo) error {
	txInfo := new(types.TxInfo)

	lastProcessedTxDigestBytes, _ := lastProcessedTxDigest.([]byte)

	if lastProcessedTxDigestBytes != nil {
		if err := txInfo.Deserialize(lastProcessedTxDigestBytes); err != nil {
			p.log.Error("failed to deserialize last processed tx digest", zap.Error(err))
			return err
		}
	}

	return p.listenByPolling(ctx, txInfo.TxDigest, blockInfo)
}

func (p *Provider) allowedEventTypes() []string {
	allowedEvents := []string{}
	for _, xcallPkgId := range p.cfg.XcallPkgIDs {
		allowedEvents = append(allowedEvents, []string{
			fmt.Sprintf("%s::%s::%s", xcallPkgId, ModuleConnection, "Message"),
			fmt.Sprintf("%s::%s::%s", xcallPkgId, ModuleMain, "CallMessage"),
			fmt.Sprintf("%s::%s::%s", xcallPkgId, ModuleMain, "RollbackMessage"),
		}...)
	}
	return allowedEvents
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

		p.log.Info("Detected event log: ",
			zap.Uint64("checkpoint", msg.MessageHeight),
			zap.String("event-type", msg.EventType),
			zap.Uint64("sn", msg.Sn),
			zap.String("dst", msg.Dst),
			zap.Uint64("req-id", msg.ReqID),
			zap.Any("data", hex.EncodeToString(msg.Data)),
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
	eventSuffix := strings.Join(eventParts[1:], "::")

	switch eventSuffix {
	case fmt.Sprintf("%s::%s", ModuleConnection, "Message"):
		msg.EventType = relayerEvents.EmitMessage
		var emitEvent types.EmitEvent
		if err := json.Unmarshal(eventBytes, &emitEvent); err != nil {
			return nil, err
		}
		if emitEvent.ConnectionID != p.cfg.ConnectionID {
			return nil, fmt.Errorf(types.InvalidEventError)
		}

		sn, err := strconv.Atoi(emitEvent.Sn)
		if err != nil {
			return nil, err
		}
		msg.Sn = uint64(sn)
		msg.Data = emitEvent.Msg
		msg.Dst = emitEvent.To

	case fmt.Sprintf("%s::%s", ModuleMain, "CallMessage"):
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

	case fmt.Sprintf("%s::%s", ModuleMain, "RollbackMessage"):
		msg.EventType = relayerEvents.RollbackMessage
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
		msg.Data = rollbackMsgEvent.Data

	default:
		return nil, fmt.Errorf(types.InvalidEventError)
	}

	msg.Src = p.cfg.NID

	return &msg, nil
}

func (p *Provider) handleEventNotification(ctx context.Context, ev types.EventResponse, blockStream chan *relayertypes.BlockInfo) {
	if ev.Checkpoint == nil {
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

	p.log.Info("Detected event log: ",
		zap.Uint64("checkpoint", msg.MessageHeight),
		zap.String("event-type", msg.EventType),
		zap.Uint64("sn", msg.Sn),
		zap.String("dst", msg.Dst),
		zap.Uint64("req-id", msg.ReqID),
		zap.Any("data", hex.EncodeToString(msg.Data)),
	)

	blockStream <- &relayertypes.BlockInfo{
		Height:   msg.MessageHeight,
		Messages: []*relayertypes.Message{msg},
	}
}

func (p *Provider) listenByPolling(ctx context.Context, lastSavedTxDigestStr string, blockStream chan *relayertypes.BlockInfo) error {
	done := make(chan interface{})
	defer close(done)

	var lastSavedTxDigest *sui_types.TransactionDigest

	if lastSavedTxDigestStr != "" { //process probably unexplored events of last saved tx digest
		lastSavedTxDigest, err := sui_types.NewDigest(lastSavedTxDigestStr)
		if err != nil {
			return err
		}
		currentEvents, err := p.client.GetEvents(ctx, *lastSavedTxDigest)
		if err != nil {
			return err
		}

		for _, ev := range currentEvents {
			p.handleEventNotification(ctx, types.EventResponse{SuiEvent: ev}, blockStream)
		}
	}

	eventStream := p.getObjectEventStream(done, p.cfg.XcallStorageID, lastSavedTxDigest)

	p.log.Info("event query started from last saved tx digest", zap.String("last-saved-tx-digest", lastSavedTxDigestStr))

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

		limit := uint(100)

		ticker := time.NewTicker(3 * time.Second)

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

				if len(res.Data) > 0 {
					var nextCursor *lib.Base58
					for _, blockRes := range res.Data {
						for _, ev := range blockRes.Events {
							eventStream <- types.EventResponse{
								SuiEvent:   ev,
								Checkpoint: blockRes.Checkpoint,
							}
							nextCursor = &ev.Id.TxDigest
						}
					}
					if nextCursor != nil {
						cursor = nextCursor
					}
				}
			}
		}
	}()

	return eventStream
}
