package stacks

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/clarity"
)

type StacksTransaction interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

type BaseTransaction struct {
	Version           TransactionVersion
	ChainID           ChainID
	Auth              TransactionAuth
	AnchorMode        AnchorMode
	PostConditionMode PostConditionMode
	PostConditions    []PostCondition
}

type TransactionAuth struct {
	AuthType    AuthType
	OriginAuth  SpendingCondition
	SponsorAuth *SpendingCondition
}

type SpendingCondition struct {
	HashMode    AddressHashMode
	Signer      [20]byte
	Nonce       uint64
	Fee         uint64
	KeyEncoding PubKeyEncoding
	Signature   [RecoverableECDSASigLengthBytes]byte
}

type TokenTransferTransaction struct {
	BaseTransaction
	Payload TokenTransferPayload
}

type ContractCallTransaction struct {
	BaseTransaction
	Payload ContractCallPayload
}

func NewTokenTransferTransaction(recipient string, amount uint64, memo string) (*TokenTransferTransaction, error) {
	payload, err := NewTokenTransferPayload(recipient, amount, memo)
	if err != nil {
		return nil, err
	}
	return &TokenTransferTransaction{
		BaseTransaction: BaseTransaction{
			Version:           TransactionVersionMainnet,
			ChainID:           ChainIDMainnet,
			AnchorMode:        AnchorModeOnChainOnly,
			PostConditionMode: PostConditionModeAllow,
			PostConditions:    []PostCondition{}, // Empty post condition
		},
		Payload: *payload,
	}, nil
}

func NewContractCallTransaction(contractAddress, contractName, functionName string, functionArgs []clarity.ClarityValue) *ContractCallTransaction {
	return &ContractCallTransaction{
		BaseTransaction: BaseTransaction{
			Version:           TransactionVersionMainnet,
			ChainID:           ChainIDMainnet,
			AnchorMode:        AnchorModeOnChainOnly,
			PostConditionMode: PostConditionModeAllow,
		},
		Payload: ContractCallPayload{
			ContractAddress: contractAddress,
			ContractName:    contractName,
			FunctionName:    functionName,
			FunctionArgs:    functionArgs,
		},
	}
}

func (t *TokenTransferTransaction) Serialize() ([]byte, error) {
	buf := make([]byte, 0, 128)

	buf = append(buf, byte(t.Version))
	chainIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(chainIDBytes, uint32(t.ChainID))
	buf = append(buf, chainIDBytes...)

	authBytes, err := t.Auth.SerializeAuth()
	if err != nil {
		return nil, err
	}
	buf = append(buf, authBytes...)

	buf = append(buf, byte(t.AnchorMode))
	buf = append(buf, byte(t.PostConditionMode))

	// assumes post condition is empty
	postConditionsLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(postConditionsLenBytes, uint32(len(t.PostConditions)))
	buf = append(buf, postConditionsLenBytes...)

	buf = append(buf, byte(PayloadTypeTokenTransfer))

	payloadBytes, err := t.Payload.Serialize()
	if err != nil {
		return nil, err
	}
	buf = append(buf, payloadBytes...)

	return buf, nil
}

func (auth *TransactionAuth) SerializeAuth() ([]byte, error) {
	buf := make([]byte, 0, 256)

	buf = append(buf, byte(auth.AuthType))

	originAuthBytes, err := auth.OriginAuth.SerializeSpendingCondition()
	if err != nil {
		return nil, err
	}
	buf = append(buf, originAuthBytes...)

	if auth.AuthType == AuthTypeSponsored {
		if auth.SponsorAuth == nil {
			return nil, errors.New("sponsor auth is required for sponsored transactions")
		}
		sponsorAuthBytes, err := auth.SponsorAuth.SerializeSpendingCondition()
		if err != nil {
			return nil, err
		}
		buf = append(buf, sponsorAuthBytes...)
	}

	return buf, nil
}

func (t *TokenTransferTransaction) Deserialize(data []byte) error {
	if len(data) < 7 {
		return errors.New("insufficient data for transaction")
	}

	t.Version = TransactionVersion(data[0])
	t.ChainID = ChainID(binary.BigEndian.Uint32(data[1:5]))
	authType := AuthType(data[5])

	offset := 5 // Start after version and chain ID

	auth := TransactionAuth{AuthType: authType}
	authLen, err := auth.DeserializeAuth(data[offset:])
	if err != nil {
		return fmt.Errorf("failed to deserialize auth: %w", err)
	}
	offset += authLen
	t.Auth = auth

	if len(data) < offset+1 {
		return errors.New("insufficient data for anchor mode")
	}
	t.AnchorMode = AnchorMode(data[offset])
	offset++

	if len(data) < offset+1 {
		return errors.New("insufficient data for post condition mode")
	}
	t.PostConditionMode = PostConditionMode(data[offset])
	offset++

	postConditions, postConditionsLen, err := DeserializePostConditions(data[offset:])
	if err != nil {
		return fmt.Errorf("failed to deserialize post conditions: %w", err)
	}
	t.PostConditions = postConditions
	offset += postConditionsLen

	if len(data) < offset+1 {
		return errors.New("insufficient data for payload type")
	}
	if PayloadType(data[offset]) != PayloadTypeTokenTransfer {
		return errors.New("invalid payload type for token transfer transaction")
	}
	offset++

	payloadLen, err := t.Payload.Deserialize(data[offset:])
	if err != nil {
		return fmt.Errorf("failed to deserialize token transfer payload: %w", err)
	}
	offset += payloadLen

	return nil
}

func (t *ContractCallTransaction) Deserialize(data []byte) error {
	if len(data) < 7 {
		return errors.New("insufficient data for transaction")
	}

	t.Version = TransactionVersion(data[0])
	t.ChainID = ChainID(binary.BigEndian.Uint32(data[1:5]))
	authType := AuthType(data[5])

	offset := 5 // Start after version and chain ID

	auth := TransactionAuth{AuthType: authType}
	authLen, err := auth.DeserializeAuth(data[offset:])
	if err != nil {
		return fmt.Errorf("failed to deserialize auth: %w", err)
	}
	offset += authLen
	t.Auth = auth

	if len(data) < offset+1 {
		return errors.New("insufficient data for anchor mode")
	}
	t.AnchorMode = AnchorMode(data[offset])
	offset++

	if len(data) < offset+1 {
		return errors.New("insufficient data for post condition mode")
	}
	t.PostConditionMode = PostConditionMode(data[offset])
	offset++

	postConditions, postConditionsLen, err := DeserializePostConditions(data[offset:])
	if err != nil {
		return fmt.Errorf("failed to deserialize post conditions: %w", err)
	}
	t.PostConditions = postConditions
	offset += postConditionsLen

	if len(data) < offset+1 {
		return errors.New("insufficient data for payload type")
	}
	if PayloadType(data[offset]) != PayloadTypeContractCall {
		return errors.New("invalid payload type for contract call transaction")
	}
	offset++

	payloadLen, err := t.Payload.Deserialize(data[offset:])
	if err != nil {
		return fmt.Errorf("failed to deserialize contract call payload: %w", err)
	}
	offset += payloadLen

	return nil
}

func (auth *TransactionAuth) DeserializeAuth(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, errors.New("invalid auth data length")
	}

	auth.AuthType = AuthType(data[0])
	offset := 1

	originAuthLen, err := auth.OriginAuth.DeserializeSpendingCondition(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += originAuthLen

	if auth.AuthType == AuthTypeSponsored {
		auth.SponsorAuth = &SpendingCondition{}
		sponsorAuthLen, err := auth.SponsorAuth.DeserializeSpendingCondition(data[offset:])
		if err != nil {
			return 0, err
		}
		offset += sponsorAuthLen
	}

	return offset, nil
}

func (sc *SpendingCondition) SerializeSpendingCondition() ([]byte, error) {
	buf := make([]byte, 0, 103) // 1 + 20 + 8 + 8 + 1 + 65

	buf = append(buf, byte(sc.HashMode))
	buf = append(buf, sc.Signer[:]...)

	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, sc.Nonce)
	buf = append(buf, nonceBytes...)

	feeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(feeBytes, sc.Fee)
	buf = append(buf, feeBytes...)

	buf = append(buf, byte(sc.KeyEncoding))
	buf = append(buf, sc.Signature[:]...)

	return buf, nil
}

func (sc *SpendingCondition) DeserializeSpendingCondition(data []byte) (int, error) {
	if len(data) < 103 {
		return 0, errors.New("invalid spending condition data length")
	}

	hashMode := AddressHashMode(data[0])
	if !isValidAddressHashMode(hashMode) {
		return 0, fmt.Errorf("invalid address hash mode: %d", hashMode)
	}
	sc.HashMode = hashMode

	copy(sc.Signer[:], data[1:21])
	sc.Nonce = binary.BigEndian.Uint64(data[21:29])
	sc.Fee = binary.BigEndian.Uint64(data[29:37])

	keyEncoding := PubKeyEncoding(data[37])
	if !isValidPubKeyEncoding(keyEncoding) {
		return 0, fmt.Errorf("invalid public key encoding: %d", keyEncoding)
	}
	sc.KeyEncoding = keyEncoding

	if !isCompatibleHashModeAndKeyEncoding(sc.HashMode, sc.KeyEncoding) {
		return 0, fmt.Errorf("incompatible hash mode (%d) and key encoding (%d)", sc.HashMode, sc.KeyEncoding)
	}

	copy(sc.Signature[:], data[38:103])

	return 103, nil
}

func isValidAddressHashMode(mode AddressHashMode) bool {
	return mode == AddressHashModeSerializeP2PKH ||
		mode == AddressHashModeSerializeP2WPKH
}

func isValidPubKeyEncoding(encoding PubKeyEncoding) bool {
	return encoding == PubKeyEncodingCompressed || encoding == PubKeyEncodingUncompressed
}

func isCompatibleHashModeAndKeyEncoding(hashMode AddressHashMode, keyEncoding PubKeyEncoding) bool {
	// P2WPKH and P2WSH only support compressed public keys
	if (hashMode == AddressHashModeSerializeP2WPKH) &&
		keyEncoding != PubKeyEncodingCompressed {
		return false
	}
	return true
}

func (t *ContractCallTransaction) Serialize() ([]byte, error) {
	buf := make([]byte, 0, 256) // Initial capacity, adjust as needed

	buf = append(buf, byte(t.Version))

	chainIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(chainIDBytes, uint32(t.ChainID))
	buf = append(buf, chainIDBytes...)

	authBytes, err := t.Auth.SerializeAuth()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize auth: %w", err)
	}
	buf = append(buf, authBytes...)

	buf = append(buf, byte(t.AnchorMode))

	buf = append(buf, byte(t.PostConditionMode))

	postConditionsBytes := SerializePostConditions(t.PostConditions)
	buf = append(buf, postConditionsBytes...)

	buf = append(buf, byte(PayloadTypeContractCall))

	payloadBytes, err := t.Payload.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize contract call payload: %w", err)
	}
	buf = append(buf, payloadBytes...)

	return buf, nil
}

func DeserializeTransaction(data []byte) (StacksTransaction, error) {
	if len(data) < 7 { // Minimum length for version, chain ID, and auth type
		return nil, errors.New("insufficient data for transaction")
	}

	version := TransactionVersion(data[0])
	chainID := ChainID(binary.BigEndian.Uint32(data[1:5]))

	offset := 5 // Start after version and chain ID

	baseTx := BaseTransaction{
		Version: version,
		ChainID: chainID,
	}

	authLen, err := baseTx.Auth.DeserializeAuth(data[offset:])
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize auth: %w", err)
	}
	offset += authLen

	if len(data) < offset+1 {
		return nil, errors.New("insufficient data for anchor mode")
	}
	baseTx.AnchorMode = AnchorMode(data[offset])
	offset++

	if len(data) < offset+1 {
		return nil, errors.New("insufficient data for post condition mode")
	}
	baseTx.PostConditionMode = PostConditionMode(data[offset])
	offset++

	postConditions, postConditionsLen, err := DeserializePostConditions(data[offset:])
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize post conditions: %w", err)
	}
	baseTx.PostConditions = postConditions
	offset += postConditionsLen

	if len(data) < offset+1 {
		return nil, errors.New("insufficient data for payload type")
	}
	payloadType := PayloadType(data[offset])
	offset++

	var tx StacksTransaction

	switch payloadType {
	case PayloadTypeTokenTransfer:
		tokenTx := &TokenTransferTransaction{BaseTransaction: baseTx}
		payloadLen, err := tokenTx.Payload.Deserialize(data[offset:])
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize token transfer payload: %w", err)
		}
		offset += payloadLen
		tx = tokenTx
	case PayloadTypeContractCall:
		contractTx := &ContractCallTransaction{BaseTransaction: baseTx}
		payloadLen, err := contractTx.Payload.Deserialize(data[offset:])
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize contract call payload: %w", err)
		}
		offset += payloadLen
		tx = contractTx
	default:
		return nil, fmt.Errorf("unsupported payload type: %d", payloadType)
	}

	return tx, nil
}
