package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	provider "github.com/icon-project/centralized-relay/relayer/chains/evm/types"
	"github.com/icon-project/centralized-relay/relayer/types"
)

func (p *Provider) QueryLatestHeight(ctx context.Context) (height uint64, err error) {
	height, err = p.client.GetBlockNumber()
	if err != nil {
		return 0, err
	}
	return
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, messagekey *types.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, messageKey *types.Message) (bool, error) {
	return true, nil
}

func (p *Provider) MessageReceived(ctx context.Context, messageKey *types.MessageKey) (bool, error) {
	return p.client.MessageReceived(&bind.CallOpts{Context: ctx}, messageKey.Src, big.NewInt(0).SetUint64(messageKey.Sn))
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
	header, err := p.client.GetHeaderByHeight(ctx, big.NewInt(0).SetUint64(key.Height))
	if err != nil {
		return nil, err
	}
	p.blockReq.FromBlock = big.NewInt(0).SetUint64(key.Height)
	p.blockReq.ToBlock = big.NewInt(0).SetUint64(key.Height)
	logs, err := p.client.FilterLogs(ctx, p.blockReq)
	if err != nil {
		return nil, err
	}
	return p.FindMessages(ctx, &provider.BlockNotification{Height: big.NewInt(0).SetUint64(key.Height), Header: header, Logs: logs, Hash: header.Hash()})
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error) {
	receipt, err := p.client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return nil, fmt.Errorf("queryTransactionReceipt: %v", err)
	}

	finalizedReceipt := &types.Receipt{
		TxHash: txHash,
		Height: receipt.BlockNumber.Uint64(),
	}

	if receipt.Status == 1 {
		finalizedReceipt.Status = true
	}

	return finalizedReceipt, nil
}
