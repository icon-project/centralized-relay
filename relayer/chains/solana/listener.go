package solana

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfo chan *relayertypes.BlockInfo) error {
	latestHeight, err := p.client.GetLatestBlockHeight(ctx)
	if err != nil {
		return err
	}

	startHeight := latestHeight
	if lastSavedHeight != 0 && lastSavedHeight < latestHeight {
		startHeight = lastSavedHeight
	}

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

	p.log.Info("started querying from height", zap.Uint64("height", startHeight))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (p *Provider) GetSignatures(fromSignature string) ([]*solrpc.TransactionSignature, error) {
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
