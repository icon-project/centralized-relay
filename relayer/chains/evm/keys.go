package evm

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func (p *EVMProvider) RestoreKeyStore(path string, secret string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	key, err := keystore.DecryptKey(data, secret)
	if err != nil {
		return err
	}
	p.wallet = key
	return nil
}

// AddressFromKeyStore returns the address of the key stored in the given keystore file.
func (p *EVMProvider) AddressFromKeyStore(keystoreFile, auth string) (string, error) {
	data, err := os.ReadFile(keystoreFile)
	if err != nil {
		return "", err
	}
	key, err := keystore.DecryptKey(data, auth)
	if err != nil {
		return "", err
	}
	return key.Address.Hex(), nil
}

func (p *EVMProvider) NewKeyStore(dir, password string) (string, error) {
	key, err := keystore.StoreKey(os.TempDir(), password, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(key.URL.Path)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(fmt.Sprintf("%s/%s.json", dir, key.Address.Hex()), data, 0o644); err != nil {
		return "", err
	}
	return key.Address.Hex(), os.Remove(key.URL.Path)
}
