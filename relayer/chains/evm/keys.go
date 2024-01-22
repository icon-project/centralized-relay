package evm

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/icon-project/centralized-relay/relayer/kms"
)

func (p *EVMProvider) RestoreKeyStore(ctx context.Context, homepath string, client kms.KMS) error {
	path := path.Join(homepath, "keystore", p.NID(), p.cfg.Keystore)
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
	key, err := keystore.DecryptKey(keystoreJson, string(secret))
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
	path := path.Join(dir, fmt.Sprintf("%s.json", key.Address.Hex()))
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return key.Address.Hex(), os.Remove(key.URL.Path)
}
