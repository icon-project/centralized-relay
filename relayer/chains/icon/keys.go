package icon

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/goloop/common/wallet"
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
	wallet, err := wallet.NewFromKeyStore(keystoreJson, secret)
	if err != nil {
		return err
	}
	p.wallet = wallet
	return nil
}

// keystorePath is the path to the keystore file
func (p *Provider) keystorePath(addr string) string {
	return path.Join(p.cfg.HomeDir, "keystore", p.NID(), addr)
}

func (p *Provider) NewKeystore(password string) (string, error) {
	priv, _ := crypto.GenerateKeyPair()
	data, err := wallet.EncryptKeyAsKeyStore(priv, []byte(password))
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
	wallet, err := wallet.NewFromKeyStore(data, []byte(password))
	if err != nil {
		return "", err
	}
	addr := wallet.Address().String()
	keystorePath := p.keystorePath(addr)
	if err := os.WriteFile(keystorePath, keystoreEncrypted, 0o644); err != nil {
		return "", err
	}
	if err := os.WriteFile(keystorePath+".pass", passwordEncrypted, 0o644); err != nil {
		return "", err
	}
	return addr, nil
}

func (p *Provider) ImportKeystore(ctx context.Context, keyPath, passphrase string) (string, error) {
	keystore, err := os.ReadFile(keyPath)
	if err != nil {
		return "", err
	}
	wallet, err := wallet.NewFromKeyStore(keystore, []byte(passphrase))
	if err != nil {
		return "", err
	}
	addr := wallet.Address().String()
	keyStoreEncrypted, err := p.kms.Encrypt(ctx, keystore)
	if err != nil {
		return "", err
	}
	keystorePath := p.keystorePath(addr)
	if err := os.WriteFile(keystorePath, keyStoreEncrypted, 0o644); err != nil {
		return "", err
	}
	passwordEncrypted, err := p.kms.Encrypt(ctx, []byte(passphrase))
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(keystorePath+".pass", passwordEncrypted, 0o644); err != nil {
		return "", err
	}
	return addr, nil
}

func (p *Provider) ConvertPrivateKey(ctx context.Context, keyPath, passphrase string) (string, error) {
	return "", fmt.Errorf("not applicable for icon chain")
}
