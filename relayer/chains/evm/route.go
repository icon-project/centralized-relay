package evm

import (
	"context"

	"github.com/icon-project/centralized-relay/relayer/types"
)

func (p *EVMProvider) Route(ctx context.Context, message *types.RouteMessage, callback types.TxResponseFunc) error {
	return nil
}
