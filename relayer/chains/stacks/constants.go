package stacks

type ChainID uint32

const (
	ChainIDTestnet ChainID = 0x80000000
	ChainIDMainnet ChainID = 0x00000001
)

type PayloadType byte

const (
	PayloadTypeTokenTransfer PayloadType = 0x00
	PayloadTypeContractCall  PayloadType = 0x02
)

type AddressType byte

const (
	AddressTypeStandard AddressType = 0x05
	AddressTypeContract AddressType = 0x06
)

type AnchorMode uint8

const (
	AnchorModeOnChainOnly AnchorMode = 0x01
)

type TransactionVersion uint8

const (
	TransactionVersionMainnet TransactionVersion = 0x00
	TransactionVersionTestnet TransactionVersion = 0x80
)

type PostConditionMode uint8

const (
	PostConditionModeAllow PostConditionMode = 0x01
	PostConditionModeDeny  PostConditionMode = 0x02
)

type PostConditionType uint8

const (
	PostConditionTypeSTX         PostConditionType = 0x00
	PostConditionTypeFungible    PostConditionType = 0x01
	PostConditionTypeNonFungible PostConditionType = 0x02
)

type FungibleConditionCode uint8

const (
	FungibleConditionCodeEqual        FungibleConditionCode = 0x01
	FungibleConditionCodeGreater      FungibleConditionCode = 0x02
	FungibleConditionCodeGreaterEqual FungibleConditionCode = 0x03
	FungibleConditionCodeLess         FungibleConditionCode = 0x04
	FungibleConditionCodeLessEqual    FungibleConditionCode = 0x05
)

type NonFungibleConditionCode uint8

const (
	NonFungibleConditionCodeSends       NonFungibleConditionCode = 0x10
	NonFungibleConditionCodeDoesNotSend NonFungibleConditionCode = 0x11
)

type PostConditionPrincipalID uint8

const (
	PostConditionPrincipalIDOrigin   PostConditionPrincipalID = 0x01
	PostConditionPrincipalIDStandard PostConditionPrincipalID = 0x02
	PostConditionPrincipalIDContract PostConditionPrincipalID = 0x03
)

type AuthType uint8

const (
	AuthTypeStandard  AuthType = 0x04
	AuthTypeSponsored AuthType = 0x05
)

type AddressHashMode uint8

const (
	AddressHashModeSerializeP2PKH  AddressHashMode = 0x00
	AddressHashModeSerializeP2WPKH AddressHashMode = 0x02
)

type AddressVersion uint8

const (
	AddressVersionMainnetSingleSig AddressVersion = 22
	AddressVersionTestnetSingleSig AddressVersion = 26
)

type PubKeyEncoding uint8

const (
	PubKeyEncodingCompressed   PubKeyEncoding = 0x00
	PubKeyEncodingUncompressed PubKeyEncoding = 0x01
)

const (
	MaxStringLengthBytes           = 128
	ClarityIntSize                 = 128
	ClarityIntByteSize             = 16
	RecoverableECDSASigLengthBytes = 65
	CompressedPubkeyLengthBytes    = 32
	UncompressedPubkeyLengthBytes  = 64
	MemoMaxLengthBytes             = 34
	AddressHashLength              = 20
)
