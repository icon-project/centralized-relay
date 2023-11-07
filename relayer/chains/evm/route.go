package evm

import (
	"context"

	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

// func (p *transactionTx) SendTransaction() error {
// 	return nil
// }

func (p *EVMProvider) Route(ctx context.Context, message providerTypes.Message, callback providerTypes.TxResponseFunc) error {
	return nil
}
