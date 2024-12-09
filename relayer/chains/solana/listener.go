package solana

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/icon-project/centralized-relay/relayer/chains/solana/types"
	relayerevents "github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/sorter"
	"github.com/near/borsh-go"
	"go.uber.org/zap"
)

var (
	MaxSupportedTxVersion = 0
)

func (p *Provider) Listener(ctx context.Context, lastProcessedTx relayertypes.LastProcessedTx, blockInfo chan *relayertypes.BlockInfo) error {
	txInfo := new(types.TxInfo)

	if lastProcessedTx.Info != nil {
		if err := txInfo.Deserialize(lastProcessedTx.Info); err != nil {
			p.log.Error("failed to deserialize last processed tx digest", zap.Error(err))
			return err
		}
	}

	if p.cfg.StartTxSign != "" {
		txInfo.TxSign = p.cfg.StartTxSign
	}

	if txInfo.TxSign == "" {
		latestTxSign, err := p.getLatestXcallTxSignature()
		if err != nil {
			p.log.Error("failed to get latest xcall tx signature", zap.Error(err))
			return err
		}
		if latestTxSign != nil {
			txInfo.TxSign = latestTxSign.Signature.String()
		}
	}

	p.log.Info("started querying", zap.String("from-signature", txInfo.TxSign))

	return p.listenByPolling(ctx, txInfo.TxSign, blockInfo)
}

func (p *Provider) listenByPolling(ctx context.Context, fromSignature string, blockInfo chan *relayertypes.BlockInfo) error {
	ticker := time.NewTicker(3 * time.Second)

	startSignature := fromSignature
	var startSignatureSlot uint64

	if startSignature != "" {
		fromSign, err := solana.SignatureFromBase58(startSignature)
		if err != nil {
			return err
		}
		txn, err := p.processTxSignature(ctx, fromSign, blockInfo)
		if err != nil {
			p.log.Error("failed to process tx signature", zap.String("signature", fromSign.String()), zap.Error(err))
		} else {
			startSignatureSlot = txn.Slot
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			//fetch txSigns from most recent to oldest. 0th index is the most recent and last index is oldest
			txSigns, err := p.getSignatures(ctx, startSignature)
			if err != nil {
				p.log.Error("failed to get signatures", zap.Error(err))
				break
			}

			// sort tx signatures in descending order
			sorter.Sort(txSigns, func(s1, s2 *solrpc.TransactionSignature) bool {
				return s1.Slot > s2.Slot
			})

			//start processing from last index i.e oldest signature
			for i := len(txSigns) - 1; i >= 0; i-- {
				// do not process signature that are older than the start tx signature slot
				if txSigns[i].Slot < startSignatureSlot {
					p.log.Error(
						"found tx signature older than the current start signature slot",
						zap.String("signature", txSigns[i].Signature.String()),
						zap.Uint64("slot", txSigns[i].Slot),
						zap.String("start-signature", startSignature),
						zap.Uint64("start-signature-slot", startSignatureSlot),
						zap.Error(err))
					continue
				}
				sign := txSigns[i].Signature
				time.Sleep(500 * time.Millisecond)
				if _, err := p.processTxSignature(ctx, sign, blockInfo); err != nil {
					p.log.Error("failed to process tx signature", zap.String("signature", sign.String()), zap.Error(err))
				}
			}

			if len(txSigns) > 0 {
				//next query start from most recent signature.
				if txSigns[0].Slot > startSignatureSlot {
					startSignature = txSigns[0].Signature.String()
					startSignatureSlot = txSigns[0].Slot
				}
			}
		}
	}
}

func (p *Provider) processTxSignature(ctx context.Context, sign solana.Signature, blockInfo chan *relayertypes.BlockInfo) (*solrpc.GetTransactionResult, error) {
	txVersion := uint64(0)
	timeoutCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	txn, err := p.client.GetTransaction(timeoutCtx, sign, &solrpc.GetTransactionOpts{MaxSupportedTransactionVersion: &txVersion})
	if err != nil {
		return nil, fmt.Errorf("failed to get txn with sign %s: %w", sign, err)
	}

	if txn.Meta != nil && txn.Meta.Err != nil {
		return nil, fmt.Errorf("failed txn with sign %s: %v", sign, txn.Meta.Err)
	}

	if txn.Meta != nil && len(txn.Meta.LogMessages) > 0 {
		event := types.SolEvent{Slot: txn.Slot, Signature: sign, Logs: txn.Meta.LogMessages}
		messages, err := p.parseMessagesFromEvent(event)
		if err != nil {
			return txn, fmt.Errorf("failed to parse messages from event [%+v]: %w", event, err)
		}
		if len(messages) > 0 {
			for _, msg := range messages {
				p.log.Info("Detected event log: ",
					zap.Uint64("height", msg.MessageHeight),
					zap.String("event-type", msg.EventType),
					zap.Any("sn", msg.Sn),
					zap.Any("req-id", msg.ReqID),
					zap.String("src", msg.Src),
					zap.String("dst", msg.Dst),
					zap.Any("data", hex.EncodeToString(msg.Data)),
				)
			}
			blockInfo <- &relayertypes.BlockInfo{
				Height:   event.Slot,
				Messages: messages,
			}
		}
	}
	return txn, nil
}

func (p *Provider) parseMessagesFromEvent(solEvent types.SolEvent) ([]*relayertypes.Message, error) {
	messages := []*relayertypes.Message{}

	txInfo := types.TxInfo{TxSign: solEvent.Signature.String()}
	txInfoBytes, err := txInfo.Serialize()
	if err != nil {
		return nil, err
	}

	for _, log := range solEvent.Logs {
		if strings.HasPrefix(log, types.EventLogPrefix) {
			eventLog := strings.Replace(log, types.EventLogPrefix, "", 1)
			eventLogBytes, err := base64.StdEncoding.DecodeString(eventLog)
			if err != nil {
				return nil, err
			}

			if len(eventLogBytes) < 8 {
				return nil, fmt.Errorf("decoded bytes too short to contain discriminator: %v", eventLogBytes)
			}

			discriminator := eventLogBytes[:8]
			eventBytes := eventLogBytes[8:]

			allEvents := p.connIdl.Events
			if len(p.cfg.Dapps) > 0 {
				allEvents = append(allEvents, p.xcallIdl.Events...)
			}

			for _, ev := range allEvents {
				if slices.Equal(ev.Discriminator, discriminator) {
					switch ev.Name {
					case types.EventSendMessage:
						smEvent := types.SendMessageEvent{}
						if err := borsh.Deserialize(&smEvent, eventBytes); err != nil {
							return nil, fmt.Errorf("failed to decode send message event: %w", err)
						}

						messages = append(messages, &relayertypes.Message{
							EventType:     relayerevents.EmitMessage,
							Sn:            &smEvent.ConnSn,
							Src:           p.NID(),
							Dst:           smEvent.TargetNetwork,
							Data:          smEvent.Msg,
							MessageHeight: solEvent.Slot,
							TxInfo:        txInfoBytes,
						})

					case types.EventCallMessage:
						cmEvent := types.CallMessageEvent{}
						if err := borsh.Deserialize(&cmEvent, eventBytes); err != nil {
							return nil, fmt.Errorf("failed to decode call message event: %w", err)
						}
						fromNID := strings.Split(cmEvent.FromNetworkAddress, "/")[0]
						connProgram := solana.PublicKeyFromBytes(cmEvent.ConnProgram[:]).String()
						if connProgram != "" {
							messages = append(messages, &relayertypes.Message{
								EventType:      relayerevents.CallMessage,
								Sn:             &cmEvent.ConnSn,
								XcallSn:        &cmEvent.Sn,
								ReqID:          &cmEvent.ReqId,
								Src:            fromNID,
								Dst:            p.NID(),
								Data:           cmEvent.Data,
								MessageHeight:  solEvent.Slot,
								TxInfo:         txInfoBytes,
								DstConnAddress: connProgram,
							})
						}

					case types.EventRollbackMessage:
						rmEvent := types.RollbackMessageEvent{}
						if err := borsh.Deserialize(&rmEvent, eventBytes); err != nil {
							return nil, fmt.Errorf("failed to decode rollback message event: %w", err)
						}
						messages = append(messages, &relayertypes.Message{
							EventType:     relayerevents.RollbackMessage,
							XcallSn:       &rmEvent.Sn,
							Src:           p.NID(),
							Dst:           p.NID(),
							MessageHeight: solEvent.Slot,
							TxInfo:        txInfoBytes,
						})
					}

					break
				}
			}
		}
	}

	return messages, nil
}

func (p *Provider) getLatestXcallTxSignature() (*solrpc.TransactionSignature, error) {
	progId := p.xcallIdl.GetProgramID()

	limit := 1
	opts := &solrpc.GetSignaturesForAddressOpts{
		Limit: &limit,
	}

	txSigns, err := p.client.GetSignaturesForAddress(context.Background(), progId, opts)
	if err != nil {
		return nil, err
	}

	if len(txSigns) > 0 {
		return txSigns[0], nil
	}

	return nil, nil
}

func (p *Provider) getSignatures(ctx context.Context, fromSignature string) ([]*solrpc.TransactionSignature, error) {
	progId := p.xcallIdl.GetProgramID()

	limit := 1000
	opts := &solrpc.GetSignaturesForAddressOpts{
		Limit: &limit,
	}

	if fromSignature != "" {
		initialFromSign, err := solana.SignatureFromBase58(fromSignature)
		if err != nil {
			return nil, err
		}
		opts.Until = initialFromSign
	}

	ticker := time.NewTicker(3 * time.Second)

	signatureList := []*solrpc.TransactionSignature{}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			timeoutCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
			defer cancel()
			txSigns, err := p.client.GetSignaturesForAddress(timeoutCtx, progId, opts)
			if err != nil {
				p.log.Error("failed to get signatures for address",
					zap.String("account", progId.String()),
					zap.String("before", opts.Before.String()),
					zap.String("until", opts.Until.String()),
					zap.Error(err),
				)
				break
			}

			p.log.Debug("signature query successful",
				zap.Int("received-count", len(txSigns)),
				zap.Any("tx-signatures", txSigns),
				zap.String("account", progId.String()),
				zap.String("before", opts.Before.String()),
				zap.String("until", opts.Until.String()),
			)

			if len(txSigns) > 0 {
				opts.Before = txSigns[len(txSigns)-1].Signature
				signatureList = append(signatureList, txSigns...)
				if len(txSigns) < limit || opts.Before == opts.Until {
					return signatureList, nil
				}
			} else {
				return signatureList, nil
			}
		}
	}
}
