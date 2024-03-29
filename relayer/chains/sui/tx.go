package sui

import (
	"context"

	relayerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txDigest string) (*relayerTypes.Receipt, error) {
	txBlock, err := p.client.GetTransaction(ctx, txDigest)
	if err != nil {
		return nil, err
	}
	receipt := &relayerTypes.Receipt{
		TxHash: txDigest,
		Height: txBlock.TimestampMs.Uint64(),
		Status: txBlock.Effects.Data.IsSuccess(),
	}
	return receipt, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayerTypes.MessageKey) (bool, error) {
	suiMessage := p.NewSuiMessage([]interface{}{
		key.Src,
		key.Sn,
	}, p.cfg.Contracts[relayerTypes.ConnectionContract], ConnectionModule, MethodGetReceipt)
	msgReceived, err := p.GetReturnValuesFromCall(ctx, suiMessage)
	if err != nil {
		return false, err
	}
	return msgReceived.(bool), nil

}
