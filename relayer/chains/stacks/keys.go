package stacks

import (
	"fmt"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func DeriveStxPrivateKey(mnemonic string, index uint32) ([]byte, error) {
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to create master key: %w", err)
	}

	path := []uint32{
		44 + bip32.FirstHardenedChild,   // Purpose
		5757 + bip32.FirstHardenedChild, // Coin type (Stacks)
		0 + bip32.FirstHardenedChild,    // Account
		0,                               // Change (external chain)
		index,                           // Address index
	}

	key := masterKey
	for _, childIndex := range path {
		key, err = key.NewChildKey(childIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to derive child key: %w", err)
		}
	}

	privateKey := key.Key

	compressedPrivKey := make([]byte, 33)
	copy(compressedPrivKey, privateKey)
	compressedPrivKey[32] = 0x01

	return compressedPrivKey, nil
}
