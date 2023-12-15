package evm

import (
	"context"
	"math/big"

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

func (p *EVMProvider) MessageReceived(ctx context.Context, messageKey types.MessageKey) (bool, error) {
	snBigInt := big.NewInt(int64(messageKey.Sn))
	return p.client.MessageReceived(nil, messageKey.Src, snBigInt)
}

// func (p *EVMProvider) QueryBalance(ctx context.Context, addr string) (*types.Coin, error) {
// 	balance, err := p.client.GetBalance(ctx, addr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &types.Coin{Amount: balance.Uint64(), Denom: "eth"}, nil
// }
