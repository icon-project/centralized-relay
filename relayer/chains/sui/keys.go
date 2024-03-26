package sui

import (
	"context"

	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
)

func (p Provider) QueryBalance(ctx context.Context, addr string) (*relayertypes.Coin, error) {
	//Todo
	return nil, nil
}

func (p Provider) NewKeystore(string) (string, error) {
	//Todo
	return "", nil
}

func (p Provider) RestoreKeystore(context.Context) error {
	//Todo
	return nil
}

func (p Provider) ImportKeystore(context.Context, string, string) (string, error) {
	//Todo
	return "", nil
}
