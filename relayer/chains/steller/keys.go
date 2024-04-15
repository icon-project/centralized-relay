package steller

import (
	"context"
)

func (p *Provider) RestoreKeystore(ctx context.Context) error {
	//Todo
	return nil
}

func (p *Provider) NewKeystore(password string) (string, error) {
	//Todo
	return "", nil
}

func (p *Provider) ImportKeystore(ctx context.Context, keyPath, passphrase string) (string, error) {
	//Todo
	return "", nil
}
