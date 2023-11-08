package icon

import (
	"encoding/json"
	"os"

	"github.com/icon-project/goloop/common/wallet"
	"github.com/icon-project/goloop/module"
)

func (cp *IconProvider) RestoreIconKeyStore() (module.Wallet, error) {
	ksByte, err := os.ReadFile(cp.PCfg.KeyStore)
	if err != nil {
		return nil, err
	}
	w, err := wallet.NewFromKeyStore(ksByte, []byte(cp.PCfg.Password))
	if err != nil {
		return nil, err
	}
	return w, nil
}

type OnlyAddr struct {
	Address string `json:"address"`
}

func getAddrFromKeystore(keystorePath string) (string, error) {
	ksFile, err := os.ReadFile(keystorePath)
	if err != nil {
		return "", err
	}

	var a OnlyAddr
	err = json.Unmarshal(ksFile, &a)
	if err != nil {
		return "", err
	}
	return a.Address, nil

}
