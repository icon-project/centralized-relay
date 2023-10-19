package evm

import (
	"context"
	"math/big"

	"github.com/avast/retry-go/v4"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func (p *EVMProvider) QueryLatestHeight(ctx context.Context) (height uint64, err error) {
	err = retry.Do(func() error {
		height, err = p.client.GetBlockNumber()
		return err
	}, retry.Context(ctx),
		retry.Attempts(RPCCallRetry), // TODO: set max retry count
		retry.OnRetry(func(n uint, err error) {
			p.log.Warn("retrying failed latestHeight query", zap.String("Chain Id", p.ChainId()))
		}))
	return
}

func (p *EVMProvider) QueryBlockByHeight(ctx context.Context, height uint64) (*ethTypes.Header, error) {
	p.log.Info("QueryBlockByHeight", zap.Uint64("height", height))
	return p.client.GetHeaderByHeight(ctx, new(big.Int).SetUint64(height))
}

func (p *EVMProvider) ExecuteEventHandlers(ctx context.Context, block *ethTypes.Header) error {
	return nil
}

func (p *EVMProvider) ShouldReceiveMessage(ctx context.Context, messagekey types.Message) (bool, error) {
	return true, nil
}

func (p *EVMProvider) ShouldSendMessage(ctx context.Context, messageKey types.Message) (bool, error) {
	return true, nil
}

func (p *EVMProvider) QueryBalance(ctx context.Context, addr string) (*types.Coin, error) {
	balance, err := p.client.GetBalance(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &types.Coin{Amount: balance.Uint64(), Denom: "eth"}, nil
}
