package icon

import (
	"fmt"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"strings"
)

var _ ibc.Wallet = &IconWallet{}

type IconWallet struct {
	mnemonic string
	address  []byte
	keyName  string
	chainCfg ibc.ChainConfig
}

func NewWallet(keyname string, address []byte, mnemonic string, chainCfg ibc.ChainConfig) ibc.Wallet {
	return &IconWallet{
		mnemonic: mnemonic,
		address:  address,
		keyName:  keyname,
		chainCfg: chainCfg,
	}
}

func (w *IconWallet) KeyName() string {
	return w.keyName
}

// Get formatted address, passing in a prefix
func (w *IconWallet) FormattedAddress() string {
	return strings.ReplaceAll(string(w.address), `"`, "")
}

// Get mnemonic, only used for relayer wallets
func (w *IconWallet) Mnemonic() string {
	return w.mnemonic
}

// Get Address with chain's prefix
func (w *IconWallet) Address() []byte {
	return w.address
}

func (w *IconWallet) FormattedAddressWithPrefix(prefix string) string {
	return fmt.Sprintf("%s%s", prefix, strings.ReplaceAll(string(w.address), `"`, ""))
}
