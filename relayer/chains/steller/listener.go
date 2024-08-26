package steller

import (
	"context"
	"fmt"
	"time"

	"github.com/icon-project/centralized-relay/relayer/chains/steller/types"
	relayerevents "github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

var (
	hzContextTimeout = 2 * time.Minute
)

func (p *Provider) Listener(ctx context.Context, lastProcessedTx relayertypes.LastProcessedTx,
	blockInfo chan *relayertypes.BlockInfo) error {
	if err := p.RestoreKeystore(ctx); err != nil {
		return fmt.Errorf("failed to restore key: %w", err)
	}
	lastSavedLedgerSeq := lastProcessedTx.Height
	latestLedger, err := p.client.GetLatestLedger(ctx)
	if err != nil {
		return err
	}

	latestSeq := latestLedger.Sequence

	startSeq := latestSeq
	if lastSavedLedgerSeq != 0 && lastSavedLedgerSeq < latestSeq {
		startSeq = lastSavedLedgerSeq
	}

	if p.cfg.StartHeight != 0 && p.cfg.StartHeight < startSeq {
		startSeq = p.cfg.StartHeight
	}

	reconnectCh := make(chan struct{}, 1) // reconnect channel

	reconnect := func() {
		select {
		case reconnectCh <- struct{}{}:
		default:
		}
	}
	p.log.Info("start querying from ledger seq", zap.Uint64("start-seq", startSeq))
	eventChannel := make(chan types.Event, 10)
	reconnect()
	hzStreamCtx, cancel := context.WithTimeout(ctx, hzContextTimeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-hzStreamCtx.Done():
			hzStreamCtx, cancel = context.WithTimeout(ctx, hzContextTimeout)
			defer cancel()
			reconnect()
		case ev, ok := <-eventChannel:
			if ok {
				if ev.ContractEvent != nil {
					messages := p.parseMessagesFromEvent(ev)
					if messages != nil {
						blockInfo <- &relayertypes.BlockInfo{
							Height: messages.MessageHeight, Messages: []*relayertypes.Message{messages},
						}
					}
				}
				startSeq = ev.LedgerSeq
			}
		case <-reconnectCh:
			p.log.Info("Query started.", zap.Uint64("from-seq", startSeq))
			eventFilter := types.EventFilter{
				LedgerSeq:   startSeq + 1,
				ContractIds: []string{p.cfg.Contracts[relayertypes.ConnectionContract], p.cfg.Contracts[relayertypes.XcallContract]},
				Topics:      []string{"Message", "CallMessage", "RollbackMessage"},
			}
			go p.client.StreamEvents(hzStreamCtx, eventFilter, eventChannel)

		}
	}
}

func (p *Provider) fetchLedgerMessages(ctx context.Context, ledgerSeq uint64) ([]*relayertypes.Message, error) {
	eventFilter := types.EventFilter{
		LedgerSeq:   ledgerSeq,
		ContractIds: []string{p.cfg.Contracts[relayertypes.ConnectionContract], p.cfg.Contracts[relayertypes.XcallContract]},
		Topics:      []string{"Message", "CallMessage", "RollbackMessage"},
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
		msg := p.parseMessagesFromEvent(ev)
		if msg != nil {
			messages = append(messages, msg)
		}
	}
	return messages, nil
}

func (p *Provider) parseMessagesFromEvent(ev types.Event) *relayertypes.Message {
	var eventType string
	for _, topic := range ev.Body.V0.Topics {
		switch topic.String() {
		case "Message":
			eventType = relayerevents.EmitMessage
		case "CallMessage":
			eventType = relayerevents.CallMessage
		case "RollbackMessage":
			eventType = relayerevents.RollbackMessage
		}
	}
	if eventType == "" {
		return nil
	}

	msg := &relayertypes.Message{
		EventType:     eventType,
		MessageHeight: ev.LedgerSeq,
	}
	scval := ev.Body.V0.Data
	scMap, ok := scval.GetMap()
	if !ok {
		return nil
	}
	for _, mapItem := range *scMap {
		switch mapItem.Key.String() {
		case "connSn", "sn":
			sn, ok := mapItem.Val.GetU128()
			if !ok {
				p.log.Warn("failed to decode sn", zap.Any("value", mapItem.Val))
				return nil
			}
			msg.Sn = types.Uint128ToBigInt(sn)
		case "reqId":
			reqId, ok := mapItem.Val.GetU128()
			if !ok {
				p.log.Warn("failed to decode req_id", zap.Any("value", mapItem.Val))
				return nil
			}
			msg.ReqID = types.Uint128ToBigInt(reqId)
		case "from":
			msg.Src = mapItem.Val.String()
		case "targetNetwork", "to":
			msg.Dst = mapItem.Val.String()
		case "msg", "data":
			data, ok := mapItem.Val.GetBytes()
			if !ok {
				p.log.Warn("failed to decode data", zap.Any("value", mapItem.Val))
				return nil
			}
			msg.Data = data
		}
	}
	switch eventType {
	case relayerevents.EmitMessage:
		msg.Src = p.cfg.NID
	case relayerevents.CallMessage:
		msg.Dst = p.cfg.NID
	case relayerevents.RollbackMessage:
		msg.Src = p.cfg.NID
		msg.Dst = p.cfg.NID
	}
	p.log.Info("Detected eventlog:", zap.Any("event", *msg))
	return msg
}
