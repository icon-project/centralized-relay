package icon

import (
	"context"
	"os"
	"path"

	"github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/goloop/common/wallet"
)

func (p *IconProvider) RestoreKeystore(ctx context.Context) error {
	path := path.Join(p.cfg.HomeDir, "keystore", p.NID(), p.cfg.Address)
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

func (p *IconProvider) NewKeystore(password string) (string, error) {
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
	keystorePath := path.Join(p.cfg.HomeDir, p.NID(), addr)
	if err := os.WriteFile(keystorePath, keystoreEncrypted, 0o644); err != nil {
		return "", err
	}
	if err := os.WriteFile(keystorePath+".pass", passwordEncrypted, 0o644); err != nil {
		return "", err
	}
	return addr, nil
}
