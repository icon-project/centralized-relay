package types

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

const (
	CodeTypeOK uint32 = 0

	ChainType string = "wasm"
)

var (
	SupportedAlgorithms       = keyring.SigningAlgoList{hd.Secp256k1}
	SupportedAlgorithmsLedger = keyring.SigningAlgoList{hd.Secp256k1}
)
