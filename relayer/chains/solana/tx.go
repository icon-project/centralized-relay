package solana

import (
	"context"

	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
)

func (p *Provider) Route(ctx context.Context, message *relayertypes.Message, callback relayertypes.TxResponseFunc) error {
	return nil
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txDigest string) (*relayertypes.Receipt, error) {
	return nil, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayertypes.MessageKey) (bool, error) {
	return false, nil
}
