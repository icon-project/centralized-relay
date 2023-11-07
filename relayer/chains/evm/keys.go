package evm

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func (p *EVMProvider) RestoreKey() (string, error) {
	privateKey, err := crypto.LoadECDSA(p.cfg.Keystore)
	if err != nil {
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address, nil
}

func (p *EVMProvider) GetWallet() {

}

type Addr struct {
	Address string `json:"address"`
}
