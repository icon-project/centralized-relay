package icon

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/goloop/common/wallet"
)

func (cp *IconProvider) RestoreKeyStore(keystorePath string, auth string) error {
	ksByte, err := os.ReadFile(keystorePath)
	if err != nil {
		return err
	}
	wallet, err := wallet.NewFromKeyStore(ksByte, []byte(auth))
	if err != nil {
		return err
	}
	cp.wallet = wallet
	return nil
}

type OnlyAddr struct {
	Address string `json:"address"`
}

func (p *IconProvider) AddressFromKeyStore(keystorePath, auth string) (string, error) {
	ksByte, err := os.ReadFile(keystorePath)
	if err != nil {
		return "", err
	}
	wallet, err := wallet.NewFromKeyStore(ksByte, []byte(auth))
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
	addr, err := p.AddressFromKeyStore(tempKey, password)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(fmt.Sprintf("%s/%s.json", dir, addr), data, 0o644); err != nil {
		return "", err
	}
	return addr, os.Remove(tempKey)
}
