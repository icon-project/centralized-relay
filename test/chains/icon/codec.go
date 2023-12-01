package icon

//
//func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
//	registry.RegisterImplementations(
//		(*exported.ClientState)(nil),
//		&tendermint.ClientState{},
//		&icon.ClientState{},
//	)
//	registry.RegisterImplementations(
//		(*exported.ConsensusState)(nil),
//		&tendermint.ConsensusState{},
//		&icon.ConsensusState{},
//	)
//	registry.RegisterImplementations(
//		(*exported.ClientMessage)(nil),
//		&icon.SignedHeader{},
//	)
//}
//
//func MakeCodec() *codec.ProtoCodec {
//	interfaceRegistry := types.NewInterfaceRegistry()
//	RegisterInterfaces(interfaceRegistry)
//	return codec.NewProtoCodec(interfaceRegistry)
//}
