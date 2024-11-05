package steller

import (
	"context"
	"os"
	"path"

	"github.com/stellar/go/keypair"
	"github.com/stellar/go/strkey"
)

func (p *Provider) RestoreKeystore(ctx context.Context) error {
	encryptedPkSeed, err := os.ReadFile(p.keystorePath(p.cfg.Address))
	if err != nil {
		return err
	}

	rawPkSeed, err := p.kms.Decrypt(ctx, encryptedPkSeed)
	if err != nil {
		return err
	}

	seed, err := strkey.Encode(strkey.VersionByteSeed, rawPkSeed)
	if err != nil {
		return err
	}

	fkp, err := keypair.ParseFull(seed)
	if err != nil {
		return err
	}

	p.wallet = fkp

	return nil
}

func (p *Provider) NewKeystore(password string) (string, error) {
	kp, err := keypair.Random()
	if err != nil {
		return "", err
	}

	rawSeed, err := strkey.Decode(strkey.VersionByteSeed, kp.Seed())
	if err != nil {
		return "", err
	}

	keyStoreContent, err := p.kms.Encrypt(context.Background(), rawSeed)
	if err != nil {
		return "", err
	}
	passphraseCipher, err := p.kms.Encrypt(context.Background(), []byte(password))
	if err != nil {
		return "", err
	}

	path := p.keystorePath(kp.Address())
	if err = os.WriteFile(path, keyStoreContent, 0o644); err != nil {
		return "", err
	}
	if err = os.WriteFile(path+".pass", passphraseCipher, 0o644); err != nil {
		return "", err
	}

	return kp.Address(), nil
}

func (p *Provider) ImportKeystore(ctx context.Context, keyPath, passphrase string) (string, error) {
	pkSeedFile, err := os.ReadFile(keyPath)
	if err != nil {
		return "", err
	}
	pkSeed := string(pkSeedFile)

	fullKeyPair, err := keypair.ParseFull(pkSeed)
	if err != nil {
		return "", err
	}

	rawSeed, err := strkey.Decode(strkey.VersionByteSeed, fullKeyPair.Seed())
	if err != nil {
		return "", err
	}

	keyStoreContent, err := p.kms.Encrypt(ctx, rawSeed)
	if err != nil {
		return "", err
	}
	passphraseCipher, err := p.kms.Encrypt(ctx, []byte(passphrase))
	if err != nil {
		return "", err
	}

	path := p.keystorePath(fullKeyPair.Address())
	if err = os.WriteFile(path, keyStoreContent, 0o644); err != nil {
		return "", err
	}
	if err = os.WriteFile(path+".pass", passphraseCipher, 0o644); err != nil {
		return "", err
	}

	return fullKeyPair.Address(), nil
}

func (p *Provider) keystorePath(addr string) string {
	return path.Join(p.cfg.HomeDir, "keystore", p.NID(), addr)
}
