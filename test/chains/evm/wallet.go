package evm

import (
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
)

var _ ibc.Wallet = &EvmWallet{}

type EvmWallet struct {
	mnemonic string
	address  []byte
	keyName  string
	chainCfg ibc.ChainConfig
}

func NewWallet(keyname string, address []byte, mnemonic string, chainCfg ibc.ChainConfig) ibc.Wallet {
	return &EvmWallet{
		mnemonic: mnemonic,
		address:  address,
		keyName:  keyname,
		chainCfg: chainCfg,
	}
}

func (w *EvmWallet) KeyName() string {
	return w.keyName
}

// Get formatted address, passing in a prefix
func (w *EvmWallet) FormattedAddress() string {
	return string(w.address)
}

// Get mnemonic, only used for relayer wallets
func (w *EvmWallet) Mnemonic() string {
	return w.mnemonic
}

// Get Address with chain's prefix
func (w *EvmWallet) Address() []byte {
	return w.address
}

func (w *EvmWallet) FormattedAddressWithPrefix(prefix string) string {
	return string(w.address)
}
