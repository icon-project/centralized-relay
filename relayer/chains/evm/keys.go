package evm

import (
	"io"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func RestoreKey(keystoreFile string, secret string) (*keystore.Key, error) {
	file, err := os.Open(keystoreFile)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	key, err := keystore.DecryptKey(data, secret)
	if err != nil {
		return nil, err
	}
	return key, nil
}
