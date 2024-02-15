package wasm

import (
	"context"
	"path"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256r1"
	"github.com/icon-project/centralized-relay/relayer/kms"
)

func (p *Provider) RestoreKeyStore(ctx context.Context, homePath string, client kms.KMS) error {
	path := path.Join(homePath, "keystore", p.NID(), p.cfg.KeyStore)
	// TODO: Restore
}

func (p *Provider) NewKeyStore(dir, password string) (string, error) {
	priv, err := secp256r1.GenPrivKey()
	if err != nil {
		return "", err
	}

	// TODO: Create keystore
}
