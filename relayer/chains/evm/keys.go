package evm

import (
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func RestoreKey(keystoreFile string, secret string) (*keystore.Key, error) {
	file, err := os.Open(keystoreFile)
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	data := make([]byte, fileSize)

	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	key, err := keystore.DecryptKey(data, secret)
	if err != nil {
		return nil, err
	}
	return key, nil
}
