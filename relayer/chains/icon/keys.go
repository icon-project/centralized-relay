package icon

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/goloop/common/wallet"
)

func (p *IconProvider) RestoreKeyStore(ctx context.Context, homePath string, client kms.KMS) error {
	path := path.Join(homePath, "keystore", p.NID(), p.cfg.KeyStore)
	keystoreJson, err := os.ReadFile(fmt.Sprintf("%s.json", path))
	if err != nil {
		return err
	}
	authCipher, err := os.ReadFile(fmt.Sprintf("%s.password", path))
	if err != nil {
		return err
	}
	secret, err := client.Decrypt(ctx, authCipher)
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

type OnlyAddr struct {
	Address string `json:"address"`
}

func (p *IconProvider) AddressFromKeyStore(keystorePath, password string) (string, error) {
	data, err := os.ReadFile(keystorePath)
	if err != nil {
		return "", err
	}
	wallet, err := wallet.NewFromKeyStore(data, []byte(password))
	if err != nil {
		return "", err
	}
	return wallet.Address().String(), nil
}

func (p *IconProvider) NewKeyStore(dir, password string) (string, error) {
	priv, _ := crypto.GenerateKeyPair()
	data, err := wallet.EncryptKeyAsKeyStore(priv, []byte(password))
	if err != nil {
		return "", err
	}
	tempKey := filepath.Join(os.TempDir(), time.Now().Format("20060102150405"))
	if err := os.WriteFile(tempKey, data, 0o644); err != nil {
		return "", err
	}
	wallet, err := wallet.NewFromKeyStore(data, []byte(password))
	if err != nil {
		return "", err
	}
	addr := wallet.Address().String()
	keystorePath := path.Join(dir, fmt.Sprintf("%s.json", addr))
	if err := os.WriteFile(keystorePath, data, 0o644); err != nil {
		return "", err
	}
	return addr, os.Remove(tempKey)
}
