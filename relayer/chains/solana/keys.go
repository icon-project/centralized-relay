package solana

import (
	"context"
)

func (p *Provider) RestoreKeystore(ctx context.Context) error {
	return nil
}

func (p *Provider) NewKeystore(password string) (string, error) {
	return "", nil
}

func (p *Provider) ImportKeystore(ctx context.Context, keyPath, passphrase string) (string, error) {
	return "", nil
}
