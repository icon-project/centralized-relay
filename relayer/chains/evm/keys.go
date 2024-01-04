package evm

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func (p *EVMProvider) RestoreKeyStore(ctx context.Context, path string, secret string) (*keystore.Key, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return keystore.DecryptKey(data, secret)
}

// AddressFromKeyStore returns the address of the key stored in the given keystore file.
func (p *EVMProvider) AddressFromKeyStore(keystoreFile string) (string, error) {
	key, err := p.RestoreKeyStore(context.TODO(), keystoreFile, p.cfg.Password)
	if err != nil {
		return "", err
	}
	return key.Address.Hex(), nil
}

func (p *EVMProvider) NewKeyStore(ctx context.Context, dir, password string) ([]byte, error) {
	key, err := keystore.StoreKey(dir, password, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%s/%s", dir, key.URL.Path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, os.Remove(path)
}
