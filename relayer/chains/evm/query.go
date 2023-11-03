package evm

import (
	"context"
	"math/big"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/centralized-relay/relayer/types"
)

func (p *EVMProvider) QueryLatestHeight(ctx context.Context) (height uint64, err error) {
	height, err = p.client.eth.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return
}

func (p *EVMProvider) QueryBlockByHeight(ctx context.Context, height uint64) (*ethTypes.Header, error) {
	return p.client.eth.HeaderByNumber(ctx, big.NewInt(int64(height)))
}

func (p *EVMProvider) QueryBlockByNumber(ctx context.Context, height uint64) (*ethTypes.Block, error) {
	return p.client.eth.BlockByNumber(ctx, big.NewInt(int64(height)))
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
