package evm

import (
	"context"

	"github.com/avast/retry-go/v4"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

func (p *EVMProvider) QueryLatestHeight(ctx context.Context) (height uint64, err error) {
	err = retry.Do(func() error {
		height, err = p.Client.GetBlockNumber()
		return err
	}, retry.Context(ctx),
		retry.Attempts(3), // TODO: set max retry count
		retry.OnRetry(func(n uint, err error) {
			p.log.Warn("retrying failed latestHeight query", zap.String("Chain Id", p.ChainId()))
		}))
	return
}

// not implemented yet
func (p *EVMProvider) QueryBlockByHeight(ctx context.Context, height uint64) (*ethTypes.Header, error) {
	return nil, nil
}

func (p *EVMProvider) ExecuteEventHandlers(ctx context.Context, block *ethTypes.Header) error {
	return nil
}
