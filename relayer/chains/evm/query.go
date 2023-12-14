package evm

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/centralized-relay/relayer/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

func (p *EVMProvider) QueryLatestHeight(ctx context.Context) (height uint64, err error) {
	height, err = p.client.GetBlockNumber()
	if err != nil {
		return 0, err
	}
	return
}

func (p *EVMProvider) QueryBalance(ctx context.Context, addr string) (*providerTypes.Coin, error) {
	//TODO:
	return nil, nil
}

func (p *EVMProvider) ShouldReceiveMessage(ctx context.Context, messagekey types.Message) (bool, error) {
	return true, nil
}

func (p *EVMProvider) ShouldSendMessage(ctx context.Context, messageKey types.Message) (bool, error) {
	return true, nil
}

// func (p *EVMProvider) QueryBalance(ctx context.Context, addr string) (*types.Coin, error) {
// 	balance, err := p.client.GetBalance(ctx, addr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &types.Coin{Amount: balance.Uint64(), Denom: "eth"}, nil
// }

// TODO: may not be need anytime soon so its ok to implement later on
func (ip *EVMProvider) GenerateMessage(ctx context.Context, key *providerTypes.MessageKey) (*providerTypes.Message, error) {
	return nil, nil
}

func (icp *EVMProvider) QueryTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error) {
	receipt, err := icp.client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return nil, fmt.Errorf("queryTransactionReceipt: %v", err)
	}

	finalizedReceipt := types.Receipt{
		TxHash: txHash,
		Height: receipt.BlockNumber.Uint64(),
	}

	if receipt.Status == 1 {
		finalizedReceipt.Status = true
	}

	return &finalizedReceipt, nil
}
