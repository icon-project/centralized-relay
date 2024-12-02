package bitcoin

import (
	"context"
	"path"
)

func (p *Provider) RestoreKeystore(ctx context.Context) error {

	return nil
}

func (p *Provider) NewKeystore(password string) (string, error) {

	return "", nil
}

// ImportKeystore imports a keystore from a file
func (p *Provider) ImportKeystore(ctx context.Context, keyPath, passphrase string) (string, error) {
	return "", nil
}

// keystorePath is the path to the keystore file
func (p *Provider) keystorePath(addr string) string {
	return path.Join(p.cfg.HomeDir, "keystore", p.NID(), addr)
}
