package solana

import (
	"context"

	"github.com/gagliardetto/solana-go"
	relayerevents "github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
)

func (p *Provider) MakeCallInstructions(msg *relayertypes.Message) ([]solana.GenericInstruction, error) {
	switch msg.EventType {
	case relayerevents.EmitMessage:

	}
	return nil, nil
}

func (p *Provider) Route(ctx context.Context, message *relayertypes.Message, callback relayertypes.TxResponseFunc) error {

	return nil
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txDigest string) (*relayertypes.Receipt, error) {
	return nil, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayertypes.MessageKey) (bool, error) {
	return false, nil
}
