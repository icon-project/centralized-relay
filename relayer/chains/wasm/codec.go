package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

var moduleBasics = []module.AppModuleBasic{
	wasm.AppModuleBasic{},
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
}

type CodecConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Codec             codec.Codec
	TxConfig          client.TxConfig
}

func GetCodecConfig(pc *ProviderConfig) *CodecConfig {
	// Set the global configuration for address prefixes
	config := sdkTypes.GetConfig()

	valAddrPrefix := pc.AccountPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixOperator
	valAddrPrefixPub := valAddrPrefix + sdkTypes.PrefixPublic

	consensusNodeAddrPrefix := pc.AccountPrefix + sdkTypes.PrefixConsensus + sdkTypes.PrefixOperator
	consensusNodeAddrPrefixPub := consensusNodeAddrPrefix + sdkTypes.PrefixPublic

	config.SetBech32PrefixForAccount(pc.AccountPrefix, pc.AccountPrefix+sdkTypes.PrefixPublic)
	config.SetBech32PrefixForValidator(valAddrPrefix, valAddrPrefixPub)
	config.SetBech32PrefixForConsensusNode(consensusNodeAddrPrefix, consensusNodeAddrPrefixPub)

	ifr := types.NewInterfaceRegistry()

	std.RegisterInterfaces(ifr)

	basicManager := module.NewBasicManager(moduleBasics...)
	basicManager.RegisterInterfaces(ifr)

	cdc := codec.NewProtoCodec(ifr)

	txConfig := tx.NewTxConfig(cdc, tx.DefaultSignModes)

	return &CodecConfig{
		InterfaceRegistry: ifr,
		Codec:             cdc,
		TxConfig:          txConfig,
	}
}
