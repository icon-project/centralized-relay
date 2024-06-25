package solana

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/icon-project/centralized-relay/relayer/chains/solana/types"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfo chan *relayertypes.BlockInfo) error {
	if err := p.RestoreKeystore(ctx); err != nil {
		p.log.Error("failed to restore keystore", zap.Error(err))
	} else {
		p.log.Info("key restore successful: ", zap.String("public-key", p.wallet.PublicKey().String()))
	}

	// if err := p.InitXcall(ctx); err != nil {
	// 	p.log.Error("failed to init xcall", zap.Error(err))
	// }

	if err := p.SendMessage(ctx, &relayertypes.Message{
		Dst:  "0x3.icon",
		Data: []byte("hello"),
	}); err != nil {
		p.log.Error("failed to send message", zap.Error(err))
	}

	fromSignature := ""

	p.log.Info("started querying from height", zap.String("from-signature", fromSignature))

	return p.listenByPolling(ctx, fromSignature, blockInfo)
}

func (p *Provider) listenByPolling(ctx context.Context, fromSignature string, blockInfo chan *relayertypes.BlockInfo) error {
	ticker := time.NewTicker(3 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			txSigns, err := p.getSignatures(ctx, fromSignature)
			if err != nil {
				p.log.Error("failed to get signatures", zap.Error(err))
				break
			}
			if len(txSigns) > 0 {
				fromSignature = txSigns[len(txSigns)-1].Signature.String()
			}
			for i := len(txSigns) - 1; i >= 0; i-- {
				sign := txSigns[i].Signature
				txn, err := p.client.GetTransaction(ctx, sign, nil)
				if err != nil {
					p.log.Error("failed to get transaction", zap.String("signature", sign.String()))
					continue
				}

				if txn.Meta != nil && len(txn.Meta.LogMessages) > 0 {
					event := types.SolEvent{Slot: txn.Slot, Signature: sign, Logs: txn.Meta.LogMessages}
					message, err := p.parseMessageFromEvent(event)
					if err != nil {
						p.log.Error("failed to parse message from event", zap.Any("event", event))
						continue
					}
					if message != nil {
						blockInfo <- &relayertypes.BlockInfo{
							Height:   message.MessageHeight,
							Messages: []*relayertypes.Message{message},
						}
					}
				}
			}
		}
	}
}

func (p *Provider) parseMessageFromEvent(ev types.SolEvent) (*relayertypes.Message, error) {
	return nil, nil
}

func (p *Provider) getSignatures(ctx context.Context, fromSignature string) ([]*solrpc.TransactionSignature, error) {
	progId, err := p.xcallIdl.GetProgramID()
	if err != nil {
		return nil, err
	}

	initialFromSign, err := solana.SignatureFromBase58(fromSignature)
	if err != nil {
		return nil, err
	}

	limit := 1000
	opts := &solrpc.GetSignaturesForAddressOpts{
		Limit: &limit,
		Until: initialFromSign,
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
			if len(txSigns) > 0 {
				opts.Before = txSigns[len(txSigns)-1].Signature
				signatureList = append(signatureList, txSigns...)
				if len(txSigns) < limit || opts.Before == initialFromSign {
					return signatureList, nil
				}
			} else {
				return signatureList, nil
			}
		}
	}
}
