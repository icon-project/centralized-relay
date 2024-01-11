package icon

import (
	"fmt"
	"os"
	"path"
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

// Decrypt the keystore file
func (p *IconProvider) DecryptKeyStore() (string, error) {
	addr := p.PCfg.GetWallet()
	if addr == "" {
		return "", fmt.Errorf("no wallet address")
	}
	// TODO: get homepath from config
	return "", nil
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
	keystorePath := path.Join(dir, fmt.Sprintf("%s.json", addr))
	if err := os.WriteFile(keystorePath, data, 0o644); err != nil {
		return "", err
	}
	return addr, os.Remove(tempKey)
}
