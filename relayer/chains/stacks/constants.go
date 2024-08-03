package stacks

type SingleSigHashMode byte

type MultiSigHashMode byte

type PubKeyEncoding byte

const (
	PubKeyEncodingCompressed PubKeyEncoding = iota
	PubKeyEncodingUncompressed
)

const (
	SerializeP2PKH SingleSigHashMode = iota
	SerializeP2WPKH
)

const (
	SerializeP2SH MultiSigHashMode = iota
	SerializeP2WSH
	SerializeP2SHNonSequential
	SerializeP2WSHNonSequential
)

type StacksMessageType byte

const (
	StacksMessageTypeAddress StacksMessageType = iota
	StacksMessageTypePrincipal
	StacksMessageTypeLengthPrefixedString
	StacksMessageTypeLengthPrefixedList
	StacksMessageTypePayload
	StacksMessageTypeMemoString
	StacksMessageTypeAssetInfo
	StacksMessageTypePostCondition
	StacksMessageTypePublicKey
	StacksMessageTypeTransactionAuthField
	StacksMessageTypeMessageSignature
)

const RECOVERABLE_ECDSA_SIG_LENGTH_BYTES = 65

type AnchorMode byte

const (
	AnchorModeOnChainOnly  AnchorMode = 0x01
	AnchorModeOffChainOnly AnchorMode = 0x02
	AnchorModeAny          AnchorMode = 0x03
)

type PostConditionMode byte

const (
	PostConditionModeAllow PostConditionMode = 0x01
	PostConditionModeDeny  PostConditionMode = 0x02
)

type TransactionVersion byte

const (
	TransactionVersionMainnet TransactionVersion = 0x00
	TransactionVersionTestnet TransactionVersion = 0x80
)

type AuthType byte

const (
	AuthTypeStandard  AuthType = 0x04
	AuthTypeSponsored AuthType = 0x05
)

type AuthFieldType byte

const (
	AuthFieldTypePublicKeyCompressed   AuthFieldType = 0x00
	AuthFieldTypePublicKeyUncompressed AuthFieldType = 0x01
	AuthFieldTypeSignatureCompressed   AuthFieldType = 0x02
	AuthFieldTypeSignatureUncompressed AuthFieldType = 0x03
)

type PayloadType byte

const (
	PayloadTypeTokenTransfer PayloadType = iota
	PayloadTypeContractCall
	PayloadTypeSmartContract
	PayloadTypeVersionedSmartContract
	PayloadTypePoisonMicroblock
	PayloadTypeCoinbase
	PayloadTypeCoinbaseToAltRecipient
	PayloadTypeNakamotoCoinbase
	PayloadTypeTenureChange
)
