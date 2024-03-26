package sui

import (
	"context"

	"github.com/icon-project/centralized-relay/relayer/types"
)

func (p Provider) QueryTransactionReceipt(ctx context.Context, txDigest string) (*types.Receipt, error) {
	//Todo
	return nil, nil
}

func (p Provider) Route(ctx context.Context, message *types.Message, callback types.TxResponseFunc) error {
	//Todo
	return nil
}

func (p Provider) ShouldReceiveMessage(ctx context.Context, message *types.Message) (bool, error) {
	//Todo
	return false, nil
}

func (p Provider) ShouldSendMessage(ctx context.Context, message *types.Message) (bool, error) {
	//Todo
	return false, nil
}

func (p Provider) MessageReceived(ctx context.Context, key *types.MessageKey) (bool, error) {
	//Todo
	return false, nil
}
