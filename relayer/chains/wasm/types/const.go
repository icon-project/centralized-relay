package types

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/relayer/v2/relayer/codecs/injective"
)

const (
	CodeTypeOK       uint32 = 0
	CodeTypeErr      uint32 = 1
	ethereumCoinType        = uint32(60)
	ChainType               = "cosmos"
)

var (
	SupportedAlgorithms       = keyring.SigningAlgoList{hd.Secp256k1, injective.EthSecp256k1}
	SupportedAlgorithmsLedger = keyring.SigningAlgoList{hd.Secp256k1, injective.EthSecp256k1}
)
