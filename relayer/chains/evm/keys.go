package evm

import (
	"context"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func (p *Provider) RestoreKeystore(ctx context.Context) error {
	path := p.keystorePath(p.cfg.Address)
	keystoreCipher, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	keystoreJson, err := p.kms.Decrypt(ctx, keystoreCipher)
	if err != nil {
		return err
	}
	authCipher, err := os.ReadFile(path + ".pass")
	if err != nil {
		return err
	}
	secret, err := p.kms.Decrypt(ctx, authCipher)
	if err != nil {
		return err
	}
	key, err := keystore.DecryptKey(keystoreJson, string(secret))
	if err != nil {
		return err
	}
	p.wallet = key
	return nil
}

func (p *Provider) NewKeystore(password string) (string, error) {
	key, err := keystore.StoreKey(os.TempDir(), password, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(key.URL.Path)
	if err != nil {
		return "", err
	}
	keystoreEncrypted, err := p.kms.Encrypt(context.Background(), data)
	if err != nil {
		return "", err
	}
	passwordEncrypted, err := p.kms.Encrypt(context.Background(), []byte(password))
	if err != nil {
		return "", err
	}
	path := p.keystorePath(key.Address.Hex())
	if err := os.WriteFile(path, keystoreEncrypted, 0o644); err != nil {
		return "", err
	}
	if err := os.WriteFile(path+".pass", passwordEncrypted, 0o644); err != nil {
		return "", err
	}
	return key.Address.Hex(), os.Remove(key.URL.Path)
}

// ImportKeystore imports a keystore from a file
func (p *Provider) ImportKeystore(ctx context.Context, keyPath, passphrase string) (string, error) {
	keystoreContent, err := os.ReadFile(keyPath)
	if err != nil {
		return "", err
	}
	key, err := keystore.DecryptKey(keystoreContent, passphrase)
	if err != nil {
		return "", err
	}
	keystoreEncrypted, err := p.kms.Encrypt(context.Background(), keystoreContent)
	if err != nil {
		return "", err
	}
	passwordEncrypted, err := p.kms.Encrypt(context.Background(), []byte(passphrase))
	if err != nil {
		return "", err
	}
	path := p.keystorePath(p.cfg.Address)
	if err := os.WriteFile(path, keystoreEncrypted, 0o644); err != nil {
		return "", err
	}
	if err := os.WriteFile(path+".pass", passwordEncrypted, 0o644); err != nil {
		return "", err
	}
	return key.Address.Hex(), nil
}

// keystorePath is the path to the keystore file
func (p *Provider) keystorePath(addr string) string {
	return path.Join(p.cfg.HomeDir, "keystore", "wallets", p.NID(), addr)
}
