package types

import (
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/relayer/v2/relayer/codecs/injective"
)

const (
	CodeTypeOK  uint32 = 0
	CodeTypeErr uint32 = 1
	ChainType          = "cosmos"
)

var (
	SupportedAlgorithms       = keyring.SigningAlgoList{hd.Secp256k1, injective.EthSecp256k1}
	SupportedAlgorithmsLedger = keyring.SigningAlgoList{hd.Secp256k1, injective.EthSecp256k1}

	// Default parameters for RPC
	RPCMaxRetryAttempts = 5
	BaseRPCRetryDelay   = 3 * time.Second
	MaxRPCRetryDelay    = 60 * time.Second
)
