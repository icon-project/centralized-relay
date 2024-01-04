package icon

import (
	"context"
	"os"

	"github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/goloop/common/wallet"
	"github.com/icon-project/goloop/module"
)

func (cp *IconProvider) RestoreKeyStore() (module.Wallet, error) {
	ksByte, err := os.ReadFile(cp.PCfg.KeyStore)
	if err != nil {
		return nil, err
	}
	return wallet.NewFromKeyStore(ksByte, []byte(cp.PCfg.Password))
}

type OnlyAddr struct {
	Address string `json:"address"`
}

func (p *IconProvider) AddressFromKeyStore(keystorePath string) (string, error) {
	ksByte, err := os.ReadFile(keystorePath)
	if err != nil {
		return "", err
	}
	a, err := wallet.ReadAddressFromKeyStore(ksByte)
	if err != nil {
		return "", err
	}
	return string(a.Bytes()), nil
}

func (p *IconProvider) NewKeyStore(ctx context.Context, dir, password string) ([]byte, error) {
	priv, _ := crypto.GenerateKeyPair()
	return wallet.EncryptKeyAsKeyStore(priv, []byte(password))
}
