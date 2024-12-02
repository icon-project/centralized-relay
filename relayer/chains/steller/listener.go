package steller

import (
	"context"
	"fmt"
	"time"

	"github.com/icon-project/centralized-relay/relayer/chains/steller/sorobanclient"
	"github.com/icon-project/centralized-relay/relayer/chains/steller/types"
	relayerevents "github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stellar/go/xdr"
	"go.uber.org/zap"
)

var (
	limit = 100
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
	p.log.Info("start querying from ledger seq", zap.Uint64("start-seq", startSeq))
	pollInterval := 6 * time.Second
	if p.cfg.PollInterval != 0 {
		pollInterval = p.cfg.PollInterval
	}
	ticker := time.NewTicker(pollInterval)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			eventFilter := p.getEventFilter(startSeq, "")
			response, err := p.client.GetEvents(ctx, eventFilter)
			if err != nil {
				p.log.Warn("error occurred while fetching transactions", zap.Error(err))
				break
			}
			for _, ev := range response.Events {
				msg := p.parseMessagesFromSorobanEvent(ev)
				if msg != nil {
					blockInfo <- &relayertypes.BlockInfo{
						Height: msg.MessageHeight, Messages: []*relayertypes.Message{msg},
					}
				}
			}
			startSeq = response.LatestLedger
		}
	}
}

func (p *Provider) getEventFilter(ledgerSeq uint64, cursor string) types.GetEventFilter {
	getEventFilter := types.GetEventFilter{}
	contractIds := []string{p.cfg.Contracts[relayertypes.ConnectionContract]}
	getEventFilter.Pagination.Limit = limit
	if cursor == "" {
		getEventFilter.StartLedger = ledgerSeq
	} else {
		getEventFilter.Pagination.Cursor = cursor
	}
	addr, ok := p.cfg.Contracts[relayertypes.XcallContract]
	isExecutor := ok && addr != ""
	if isExecutor {
		contractIds = append(contractIds, p.cfg.Contracts[relayertypes.XcallContract])
	}
	getEventFilter.Filters = append(getEventFilter.Filters, types.Filter{
		ContractIDS: contractIds,
		Type:        "contract",
	})

	return getEventFilter
}

func (p *Provider) parseMessagesFromSorobanEvent(ev sorobanclient.LedgerEvents) *relayertypes.Message {
	var eventType string
	for _, topic := range ev.Topic {
		var decodedRes xdr.ScVal
		err := xdr.SafeUnmarshalBase64(topic, &decodedRes)
		if err != nil {
			return nil
		}
		switch decodedRes.String() {
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
		MessageHeight: uint64(ev.Ledger),
	}
	var scval xdr.ScVal
	err := xdr.SafeUnmarshalBase64(ev.Value, &scval)
	if err != nil {
		return nil
	}
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

func (p *Provider) fetchLedgerMessages(ctx context.Context, ledgerSeq uint64) ([]*relayertypes.Message, error) {
	eventFilter := p.getEventFilter(ledgerSeq, "")
	response, err := p.client.GetEvents(ctx, eventFilter)
	if err != nil {
		p.log.Warn("error occurred while fetching transactions", zap.Error(err))
		return nil, err
	}
	var messages []*relayertypes.Message
	for _, ev := range response.Events {
		msg := p.parseMessagesFromSorobanEvent(ev)
		if msg != nil {
			p.log.Info("detected event log:", zap.Any("event", *msg))
			messages = append(messages, msg)
		}
	}
	p.log.Debug("query successful", zap.Uint64("ledger-seq", ledgerSeq))
	return messages, err
}
