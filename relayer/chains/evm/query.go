package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	provider "github.com/icon-project/centralized-relay/relayer/chains/evm/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/types"
)

func (p *Provider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	return p.client.GetBlockNumber(ctx)
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, messagekey *types.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, messageKey *types.Message) (bool, error) {
	return true, nil
}

func (p *Provider) MessageReceived(ctx context.Context, messageKey *types.MessageKey) (bool, error) {
	switch messageKey.EventType {
	case events.EmitMessage:
		ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
		defer cancel()
		return p.client.MessageReceived(&bind.CallOpts{Context: ctx}, messageKey.Src, messageKey.Sn)
	case events.CallMessage:
		return false, nil
	case events.RollbackMessage:
		return false, nil
	default:
		return true, fmt.Errorf("unknown event type")
	}
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*types.Coin, error) {
	balance, err := p.client.GetBalance(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &types.Coin{Amount: balance.Uint64(), Denom: "eth"}, nil
}

// TODO: may not be need anytime soon so its ok to implement later on
func (p *Provider) GenerateMessages(ctx context.Context, key *types.MessageKeyWithMessageHeight) ([]*types.Message, error) {
	header, err := p.client.GetHeaderByHeight(ctx, new(big.Int).SetUint64(key.Height))
	if err != nil {
		return nil, err
	}
	p.blockReq.FromBlock = new(big.Int).SetUint64(key.Height)
	p.blockReq.ToBlock = new(big.Int).SetUint64(key.Height)
	logs, err := p.client.FilterLogs(ctx, p.blockReq)
	if err != nil {
		return nil, err
	}
	return p.FindMessages(ctx, &provider.BlockNotification{Height: new(big.Int).SetUint64(key.Height), Header: header, Logs: logs, Hash: header.Hash()})
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error) {
	receipt, err := p.client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return nil, fmt.Errorf("queryTransactionReceipt: %v", err)
	}

	finalizedReceipt := &types.Receipt{
		TxHash: txHash,
		Height: receipt.BlockNumber.Uint64(),
		Status: receipt.Status == ethTypes.ReceiptStatusSuccessful,
	}

	return finalizedReceipt, nil
}
