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

func (p *Provider) Listener(ctx context.Context, lastSavedCheckpointSeq uint64, blockInfo chan *relayertypes.BlockInfo) error {
	latestCheckpointSeq, err := p.client.GetLatestCheckpointSeq(ctx)
	if err != nil {
		return err
	}

	startCheckpointSeq := latestCheckpointSeq
	if lastSavedCheckpointSeq != 0 && lastSavedCheckpointSeq < latestCheckpointSeq {
		startCheckpointSeq = lastSavedCheckpointSeq
	}

	// go p.listenRealtime(ctx, blockInfo)

	return p.listenByPollingV1(ctx, startCheckpointSeq, blockInfo)
}

func (p *Provider) listenByPolling(ctx context.Context, startCheckpointSeq, endCheckpointSeq uint64, blockStream chan *relayertypes.BlockInfo) error {
	done := make(chan interface{})
	defer close(done)

	txDigestsStream := p.getTxDigestsStream(done, startCheckpointSeq, endCheckpointSeq)

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

func (p *Provider) getTxDigestsStream(done chan interface{}, fromSeq, toSeq uint64) <-chan types.TxDigests {
	txDigestsStream := make(chan types.TxDigests, 50)

	go func() {
		afterSeq := strconv.Itoa(int(fromSeq) - 1)
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

func (p *Provider) listenRealtime(ctx context.Context, blockStream chan *relayertypes.BlockInfo) error {
	eventTypes := []map[string]interface{}{}
	for _, evType := range p.allowedEventTypes() {
		eventTypes = append(eventTypes, map[string]interface{}{
			"MoveEventType": evType,
		})
	}
	eventFilters := map[string]interface{}{
		"Any": eventTypes,
	}

	done := make(chan interface{})
	defer close(done)

	wsUrl := strings.Replace(p.cfg.RPCUrl, "http", "ws", 1)

	eventStream, err := p.client.SubscribeEventNotification(done, wsUrl, eventFilters)
	if err != nil {
		p.log.Error("failed to subscribe event notification", zap.Error(err))
		return err
	}

	reconnectCh := make(chan bool)

	p.log.Info("started realtime checkpoint listener")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case en, ok := <-eventStream:
			if ok {
				if en.Error != nil {
					p.log.Error("failed to read event notification", zap.Error(en.Error))
					go func() {
						reconnectCh <- true
					}()
				} else {
					event := types.EventResponse{
						SuiEvent: en.SuiEvent,
					}
					go p.handleEventNotification(ctx, event, blockStream)
				}
			}
		case val := <-reconnectCh:
			if val {
				p.log.Warn("something went wrong while reading from websocket conn: reconnecting...")
				eventStream, err = p.client.SubscribeEventNotification(done, wsUrl, eventFilters)
				if err != nil {
					return err
				}
				p.log.Warn("websocket conn restablished")

			}
		}
	}
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

func (p *Provider) listenByPollingV1(ctx context.Context, fromCheckpointSeq uint64, blockStream chan *relayertypes.BlockInfo) error {
	prevCheckpoint, err := p.client.GetCheckpoint(ctx, fromCheckpointSeq-1)
	if err != nil {
		return fmt.Errorf("failed to get previous checkpoint[%d]: %w", fromCheckpointSeq-1, err)
	}

	done := make(chan interface{})
	defer close(done)

	afterTxDigest := prevCheckpoint.Transactions[len(prevCheckpoint.Transactions)-1]
	eventStream := p.getObjectEventStream(done, p.cfg.XcallStorageID, afterTxDigest)

	p.log.Info("event query started", zap.Uint64("checkpoint", fromCheckpointSeq))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev, ok := <-eventStream:
			if ok {
				go p.handleEventNotification(ctx, ev, blockStream)
			}
		}
	}
}

func (p *Provider) listenByEventPolling(ctx context.Context, fromCheckpointSeq, toCheckpointSeq uint64, blockStream chan *relayertypes.BlockInfo) error {
	prevCheckpoint, err := p.client.GetCheckpoint(ctx, fromCheckpointSeq-1)
	if err != nil {
		return fmt.Errorf("failed to get from-checkpoint: %w", err)
	}

	done := make(chan interface{})
	defer close(done)

	eventPkgId := p.cfg.XcallPkgIDs[len(p.cfg.XcallPkgIDs)-1]
	afterTxDigest := prevCheckpoint.Transactions[len(prevCheckpoint.Transactions)-1]
	eventStream := p.getPollEventStream(done, eventPkgId, ModuleMain, afterTxDigest)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev, ok := <-eventStream:
			if ok {
				go p.handleEventNotification(ctx, types.EventResponse{SuiEvent: ev}, blockStream)
			}
		}
	}
}

func (p *Provider) getObjectEventStream(done chan interface{}, objectID string, afterTxDigest string) <-chan types.EventResponse {
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

		cursor, err := sui_types.NewDigest(afterTxDigest)
		if err != nil {
			p.log.Panic("failed to create new tx digest from base58 string", zap.Error(err))
		}

		limit := uint(100)

		ticker := time.NewTicker(3 * time.Second)

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				res, err := p.client.QueryTxBlocks(context.Background(), query, cursor, &limit, false)
				if err != nil {
					p.log.Error("failed to query tx blocks", zap.Error(err), zap.String("cursor", cursor.String()))
					break
				}

				p.log.Debug("tx block query successful", zap.String("cursor", cursor.String()))

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

func (p *Provider) getPollEventStream(done chan interface{}, packageId string, eventModule string, afterTxDigest string) <-chan cctypes.SuiEvent {
	eventStream := make(chan cctypes.SuiEvent)

	go func() {
		defer close(eventStream)

		req := types.EventQueryRequest{
			EventFilter: map[string]interface{}{
				"MoveEventModule": map[string]interface{}{
					"package": packageId,
					"module":  eventModule,
				},
			},
			Cursor: cctypes.EventId{
				TxDigest: lib.Base58(afterTxDigest),
			},
			Limit:      100,
			Descending: false,
		}

		ticker := time.NewTicker(3 * time.Second)

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				res, err := p.client.QueryEvents(context.Background(), req)
				if err != nil {
					p.log.Error("failed to query events", zap.Error(err))
					break
				}

				if len(res.Data) > 0 {
					for _, ev := range res.Data {
						eventStream <- ev
					}
					lastEvent := res.Data[len(res.Data)-1]
					req.Cursor = lastEvent.Id
				}
			}
		}
	}()

	return eventStream
}
