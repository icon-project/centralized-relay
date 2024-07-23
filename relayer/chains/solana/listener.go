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
	"github.com/near/borsh-go"
	"go.uber.org/zap"
)

var (
	MaxSupportedTxVersion = 0
)

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfo chan *relayertypes.BlockInfo) error {
	// fromSignature := "2D1RxZptfGrfHwqfzMD3Z6TdnGb5Ugff3ZHuDmNT7xMTfpRtaBSpdwrC2R8qrioDDFvkU3TU5yTSukv2iByoGcuN"
	fromSignature := ""

	if err := p.RestoreKeystore(ctx); err != nil {
		p.log.Error("failed to restore wallet", zap.Error(err))
	}

	p.log.Info("started querying from height", zap.String("from-signature", fromSignature))

	return p.listenByPolling(ctx, fromSignature, blockInfo)
}

func (p *Provider) listenByPolling(ctx context.Context, fromSignature string, blockInfo chan *relayertypes.BlockInfo) error {
	ticker := time.NewTicker(3 * time.Second)

	startSignature := fromSignature

	if startSignature != "" {
		fromSign, err := solana.SignatureFromBase58(startSignature)
		if err != nil {
			return err
		}
		if err := p.processTxSignature(ctx, fromSign, blockInfo); err != nil {
			p.log.Error("failed to process tx signature", zap.String("signature", fromSign.String()), zap.Error(err))
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
			if len(txSigns) > 0 {
				//next query start from most recent signature.
				startSignature = txSigns[0].Signature.String()
			}
			//start processing from last index i.e oldest signature
			for i := len(txSigns) - 1; i >= 0; i-- {
				sign := txSigns[i].Signature
				if err := p.processTxSignature(ctx, sign, blockInfo); err != nil {
					p.log.Error("failed to process tx signature", zap.String("signature", sign.String()), zap.Error(err))
				}
			}
		}
	}
}

func (p *Provider) processTxSignature(ctx context.Context, sign solana.Signature, blockInfo chan *relayertypes.BlockInfo) error {
	txVersion := uint64(0)
	txn, err := p.client.GetTransaction(ctx, sign, &solrpc.GetTransactionOpts{MaxSupportedTransactionVersion: &txVersion})
	if err != nil {
		return fmt.Errorf("failed to get txn with sign %s: %w", sign, err)
	}

	if txn.Meta != nil && len(txn.Meta.LogMessages) > 0 {
		event := types.SolEvent{Slot: txn.Slot, Signature: sign, Logs: txn.Meta.LogMessages}
		messages, err := p.parseMessagesFromEvent(event)
		if err != nil {
			return fmt.Errorf("failed to parse messages from event [%+v]: %w", event, err)
		}
		if len(messages) > 0 {
			for _, msg := range messages {
				p.log.Info("Detected event log: ",
					zap.Uint64("height", msg.MessageHeight),
					zap.String("event-type", msg.EventType),
					zap.Uint64("sn", msg.Sn),
					zap.Uint64("req-id", msg.ReqID),
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
	return nil
}

func (p *Provider) parseMessagesFromEvent(solEvent types.SolEvent) ([]*relayertypes.Message, error) {
	messages := []*relayertypes.Message{}

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

			for _, ev := range p.xcallIdl.Events {
				if slices.Equal(ev.Discriminator, discriminator) {
					switch ev.Name {
					case types.EventSendMessage:
						smEvent := types.SendMessageEvent{}
						if err := borsh.Deserialize(&smEvent, eventBytes); err != nil {
							return nil, fmt.Errorf("failed to decode send message event: %w", err)
						}
						messages = append(messages, &relayertypes.Message{
							EventType:     relayerevents.EmitMessage,
							Sn:            smEvent.ConnSn.Uint64(),
							Src:           p.NID(),
							Dst:           smEvent.TargetNetwork,
							Data:          smEvent.Msg,
							MessageHeight: solEvent.Slot,
						})

					case types.EventCallMessage:
						cmEvent := types.CallMessageEvent{}
						if err := borsh.Deserialize(&cmEvent, eventBytes); err != nil {
							return nil, fmt.Errorf("failed to decode call message event: %w", err)
						}
						fromNID := strings.Split(cmEvent.FromNetworkAddress, "/")[0]
						messages = append(messages, &relayertypes.Message{
							EventType:     relayerevents.CallMessage,
							Sn:            cmEvent.Sn.Uint64(),
							ReqID:         cmEvent.ReqId.Uint64(),
							Src:           fromNID,
							Dst:           cmEvent.To,
							Data:          cmEvent.Data,
							MessageHeight: solEvent.Slot,
						})

					case types.EventRollbackMessage:
						rmEvent := types.RollbackMessageEvent{}
						if err := borsh.Deserialize(&rmEvent, eventBytes); err != nil {
							return nil, fmt.Errorf("failed to decode rollback message event: %w", err)
						}
						messages = append(messages, &relayertypes.Message{
							EventType:     relayerevents.RollbackMessage,
							Sn:            rmEvent.Sn.Uint64(),
							Src:           p.NID(),
							Dst:           p.NID(),
							MessageHeight: solEvent.Slot,
						})
					}

					break
				}
			}
		}
	}

	return messages, nil
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
			txSigns, err := p.client.GetSignaturesForAddress(context.Background(), progId, opts)
			if err != nil {
				p.log.Error("failed to get signatures for address",
					zap.String("account", progId.String()),
					zap.String("before", opts.Before.String()),
					zap.String("until", opts.Until.String()),
				)
				break
			}

			p.log.Debug("signature query successful",
				zap.Int("received-count", len(txSigns)),
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
