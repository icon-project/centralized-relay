package wasm

import (
	"sync"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmos/cosmos-sdk/codec/legacy"

	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/relayer/v2/relayer/codecs/injective"
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

func (c *Config) MakeCodec(moduleBasics []module.AppModuleBasic, extraCodecs ...string) *Codec {
	encodingConfig := c.makeCodecConfig()
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	basicManager := module.NewBasicManager(moduleBasics...)
	basicManager.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	for _, codec := range extraCodecs {
		switch codec {
		case "injective":
			injective.RegisterInterfaces(encodingConfig.InterfaceRegistry)
			legacy.Cdc.RegisterConcrete(&injective.PubKey{}, injective.PubKeyName, nil)
			legacy.Cdc.RegisterConcrete(&injective.PrivKey{}, injective.PrivKeyName, nil)
		}
	}
	return encodingConfig
}

func (c *Config) makeCodecConfig() *Codec {
	interfaceRegistry := types.NewInterfaceRegistry()
	done := SetSDKConfigContext(c.AccountPrefix)
	defer done()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	return &Codec{
		InterfaceRegistry: interfaceRegistry,
		Codec:             cdc,
		TxConfig:          tx.NewTxConfig(cdc, tx.DefaultSignModes),
	}
}

// This file is cursed and this mutex is too
// you don't want none of this dewey cox.
var sdkConfigMutex sync.Mutex

// SetSDKContext sets the SDK config to the proper bech32 prefixes.
// Don't use this unless you know what you're doing.
// TODO: :dagger: :knife: :chainsaw: remove this function
func (cc *Provider) SetSDKContext() func() {
	return SetSDKConfigContext(cc.cfg.AccountPrefix)
}

// SetSDKContext sets the SDK config to the given bech32 prefixes
func SetSDKConfigContext(prefix string) func() {
	sdkConfigMutex.Lock()
	sdkConf := sdkTypes.GetConfig()
	sdkConf.SetBech32PrefixForAccount(prefix, prefix+"pub")
	sdkConf.SetBech32PrefixForValidator(prefix+"valoper", prefix+"valoperpub")
	sdkConf.SetBech32PrefixForConsensusNode(prefix+"valcons", prefix+"valconspub")
	return sdkConfigMutex.Unlock
}
