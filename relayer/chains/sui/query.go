package sui

import (
	"context"
	"strconv"

	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
)

const suiCurrencyDenom = "MIST"

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*relayertypes.Coin, error) {
	balance, err := p.client.GetBalance(ctx, addr)
	if err != nil {
		return nil, err
	}
	var totalBalance uint64 = 0
	for _, bal := range balance {
		balance64, err := strconv.ParseUint(bal.Balance, 10, 64)
		if err != nil {
			return nil, err
		}
		totalBalance += (balance64)
	}
	return &relayertypes.Coin{Amount: totalBalance, Denom: suiCurrencyDenom}, nil
}
