package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/bank"
	ethermintcodecs "github.com/cosmos/relayer/v2/relayer/codecs/ethermint"
	injectivecodecs "github.com/cosmos/relayer/v2/relayer/codecs/injective"
)

var moduleBasics = []module.AppModuleBasic{
	wasm.AppModuleBasic{},
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
}

type Codec struct {
	InterfaceRegistry types.InterfaceRegistry
	TxConfig          client.TxConfig
	Codec             codec.Codec
	Amino             *codec.LegacyAmino
}

func (c *Config) MakeCodec(moduleBasics []module.AppModuleBasic) Codec {
	encodingConfig := c.makeCodecConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	basicManager := module.NewBasicManager(moduleBasics...)
	basicManager.RegisterLegacyAminoCodec(encodingConfig.Amino)
	basicManager.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ethermintcodecs.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	injectivecodecs.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	encodingConfig.Amino.RegisterConcrete(&injectivecodecs.PubKey{}, injectivecodecs.PubKeyName, nil)
	encodingConfig.Amino.RegisterConcrete(&injectivecodecs.PrivKey{}, injectivecodecs.PrivKeyName, nil)
	return encodingConfig
}

func (c *Config) makeCodecConfig() Codec {
	// Set the global configuration for address prefixes
	config := sdkTypes.GetConfig()

	valAddrPrefix := c.AccountPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixOperator
	valAddrPrefixPub := valAddrPrefix + sdkTypes.PrefixPublic

	consensusNodeAddrPrefix := c.AccountPrefix + sdkTypes.PrefixConsensus + sdkTypes.PrefixOperator
	consensusNodeAddrPrefixPub := consensusNodeAddrPrefix + sdkTypes.PrefixPublic

	config.SetBech32PrefixForAccount(c.AccountPrefix, c.AccountPrefix+sdkTypes.PrefixPublic)
	config.SetBech32PrefixForValidator(valAddrPrefix, valAddrPrefixPub)
	config.SetBech32PrefixForConsensusNode(consensusNodeAddrPrefix, consensusNodeAddrPrefixPub)
	interfaceRegistry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	return Codec{
		InterfaceRegistry: interfaceRegistry,
		Codec:             cdc,
		TxConfig:          tx.NewTxConfig(cdc, tx.DefaultSignModes),
		Amino:             codec.NewLegacyAmino(),
	}
}
