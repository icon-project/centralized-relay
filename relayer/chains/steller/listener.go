package steller

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/url"
	"slices"
	"time"

	"github.com/icon-project/centralized-relay/relayer/chains/steller/types"
	relayerevents "github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/strkey"
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
	eventFilter := p.getEventFilter(0)
	ledger, err := p.client.LedgerDetail(uint32(startSeq))
	if err != nil {
		p.log.Error("error getting ledger details", zap.Error(err))
		return err
	}
	ledgerCursor := ledger.PagingToken()
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
			hasNext := true
			for hasNext {
				p.log.Debug("Querying for ledger ", zap.Any("cursor", ledgerCursor))
				trRequest := horizonclient.TransactionRequest{
					Cursor:        ledgerCursor,
					Limit:         uint(limit),
					IncludeFailed: false,
					Order:         horizonclient.OrderAsc,
				}
				txns, err := p.client.Transactions(trRequest)
				if err != nil {
					p.log.Warn("error occurred while fetching transactions", zap.Error(err))
					break
				}
				p.log.Debug("got ledger result", zap.Any("size", len(txns.Embedded.Records)))
				p.processTxns(txns, eventFilter, blockInfo)
				newledgerCursor := ""
				if txns.Links.Next.Href != "" {
					parsedURL, err := url.Parse(txns.Links.Next.Href)
					if err != nil {
						p.log.Warn("error occurred while parsing cursor url", zap.Error(err))
						break
					}
					queryParams := parsedURL.Query()
					cursor := queryParams.Get("cursor")
					newledgerCursor = cursor
				}
				if newledgerCursor == ledgerCursor || len(txns.Embedded.Records) < limit {
					hasNext = false
				}
				if newledgerCursor != "" {
					ledgerCursor = newledgerCursor
				}
			}
		}
	}
}

func (p *Provider) processTxns(txns horizon.TransactionsPage, eventFilter types.EventFilter,
	blockInfo chan *relayertypes.BlockInfo) {
	for _, txn := range txns.Embedded.Records {
		if !txn.Successful {
			continue
		}
		var txnMeta xdr.TransactionMeta
		if err := xdr.SafeUnmarshalBase64(txn.ResultMetaXdr, &txnMeta); err != nil {
			continue
		}
		if txnMeta.V3 == nil || txnMeta.V3.SorobanMeta == nil {
			continue
		}
		if len(txnMeta.V3.SorobanMeta.Events) == 0 {
			continue
		}
		for _, ev := range txnMeta.V3.SorobanMeta.Events {
			hexBytes, err := hex.DecodeString(ev.ContractId.HexString())
			if err != nil {
				continue
			}
			contractID, err := strkey.Encode(strkey.VersionByteContract, hexBytes)
			if err != nil {
				continue
			}
			if slices.Contains(eventFilter.ContractIds, contractID) {
				for _, topic := range ev.Body.V0.Topics {
					if slices.Contains(eventFilter.Topics, topic.String()) {
						event := types.Event{
							ContractEvent: &ev,
							LedgerSeq:     uint64(txn.Ledger),
						}
						messages := p.parseMessagesFromEvent(event)
						if messages != nil {
							blockInfo <- &relayertypes.BlockInfo{
								Height: messages.MessageHeight, Messages: []*relayertypes.Message{messages},
							}
						}
					}
				}
			}
		}
	}
}

func (p *Provider) fetchLedgerMessages(ctx context.Context, ledgerSeq uint64) ([]*relayertypes.Message, error) {
	eventFilter := p.getEventFilter(ledgerSeq)
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

func (p *Provider) getEventFilter(ledgerSeq uint64) types.EventFilter {
	contractIds := []string{p.cfg.Contracts[relayertypes.ConnectionContract]}
	topics := []string{"Message"}

	addr, ok := p.cfg.Contracts[relayertypes.XcallContract]
	isExecutor := ok && addr != ""
	if isExecutor {
		contractIds = append(contractIds, p.cfg.Contracts[relayertypes.XcallContract])
		topics = append(topics, []string{"CallMessage", "RollbackMessage"}...)
	}

	if ledgerSeq <= 0 {
		return types.EventFilter{
			ContractIds: contractIds,
			Topics:      topics,
		}
	}
	return types.EventFilter{
		LedgerSeq:   ledgerSeq,
		ContractIds: contractIds,
		Topics:      topics,
	}
}
