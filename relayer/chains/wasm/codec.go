package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
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
}

func GetCodecConfig() CodecConfig {
	ifr := types.NewInterfaceRegistry()

	std.RegisterInterfaces(ifr)

	basicManager := module.NewBasicManager(moduleBasics...)
	basicManager.RegisterInterfaces(ifr)

	return CodecConfig{
		InterfaceRegistry: ifr,
		Codec:             codec.NewProtoCodec(ifr),
	}
}
