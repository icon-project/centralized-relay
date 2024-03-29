package sui

import (
	"context"

	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
)

const suiCurrencyDenom = "SUI"

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*relayertypes.Coin, error) {
	balance, err := p.client.GetTotalBalance(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &relayertypes.Coin{Amount: balance, Denom: suiCurrencyDenom}, nil
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, messagekey *relayertypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, messageKey *relayertypes.Message) (bool, error) {
	return true, nil
}
