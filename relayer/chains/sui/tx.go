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
		Height: txBlock.Checkpoint.Uint64(),
		Status: txBlock.Effects.Data.IsSuccess(),
	}
	return receipt, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayerTypes.MessageKey) (bool, error) {
	suiMessage := p.NewSuiMessage([]SuiCallArg{
		{Type: CallArgObject, Val: p.cfg.XcallStorageID},
		{Type: CallArgPure, Val: key.Src},
		{Type: CallArgPure, Val: key.Sn},
	}, p.cfg.XcallPkgID, EntryModule, MethodGetReceipt)
	var msgReceived bool
	wallet, err := p.Wallet()
	if err != nil {
		return msgReceived, err
	}
	if err := p.client.QueryContract(ctx, suiMessage, wallet.Address, p.cfg.GasLimit, &msgReceived); err != nil {
		return msgReceived, err
	}
	return msgReceived, nil
}
