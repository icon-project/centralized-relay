package stacks

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/clarity"
)

type Payload interface {
	Type() PayloadType
	Serialize() ([]byte, error)
}

type TokenTransferPayload struct {
	Recipient clarity.PrincipalCV
	Amount    uint64
	Memo      MemoString
}

func (p *TokenTransferPayload) Type() PayloadType {
	return PayloadTypeTokenTransfer
}

type StacksTransaction struct {
	Version           TransactionVersion
	ChainID           uint32
	Auth              Authorization
	AnchorMode        AnchorMode
	PostConditionMode PostConditionMode
	PostConditions    []PostCondition
	Payload           Payload
}

type Authorization interface {
	Serialize() ([]byte, error)
}

type StandardAuthorization struct {
	SpendingCondition SpendingCondition
}

type SponsoredAuthorization struct {
	SpendingCondition        SpendingCondition
	SponsorSpendingCondition SpendingCondition
}

type SpendingCondition interface {
	Serialize() ([]byte, error)
}

type SingleSigSpendingCondition struct {
	HashMode    SingleSigHashMode
	Signer      []byte
	Nonce       uint64
	Fee         uint64
	KeyEncoding PubKeyEncoding
	Signature   MessageSignature
}

type MultiSigSpendingCondition struct {
	HashMode           MultiSigHashMode
	Signer             []byte
	Nonce              uint64
	Fee                uint64
	Fields             []TransactionAuthField
	SignaturesRequired uint16
}

type TransactionAuthFieldContents interface {
	Serialize() ([]byte, error)
}

type TransactionAuthField struct {
	Type           StacksMessageType
	PubKeyEncoding PubKeyEncoding
	Contents       TransactionAuthFieldContents
}

type PostCondition interface {
	Serialize() ([]byte, error)
}

type StacksPublicKey struct {
	Type StacksMessageType
	Data []byte
}

type LengthPrefixedList struct {
	LengthPrefixBytes int
	Values            []interface{}
}

func (l *LengthPrefixedList) Type() StacksMessageType {
	return StacksMessageTypeLengthPrefixedList
}

func (l *LengthPrefixedList) Serialize() ([]byte, error) {
	buffer := new(bytes.Buffer)

	lengthBytes := make([]byte, l.LengthPrefixBytes)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(l.Values)))
	if _, err := buffer.Write(lengthBytes); err != nil {
		return nil, fmt.Errorf("failed to write length prefix: %w", err)
	}

	for _, value := range l.Values {
		var serializedValue []byte
		var err error

		switch v := value.(type) {
		case StacksMessage:
			serializedValue, err = v.Serialize()
		case PostCondition:
			serializedValue, err = v.Serialize()
		case TransactionAuthField:
			serializedValue, err = v.Serialize()
		default:
			return nil, fmt.Errorf("unsupported type in LengthPrefixedList")
		}

		if err != nil {
			return nil, fmt.Errorf("failed to serialize list item: %w", err)
		}
		if _, err := buffer.Write(serializedValue); err != nil {
			return nil, fmt.Errorf("failed to write serialized list item: %w", err)
		}
	}

	return buffer.Bytes(), nil
}

func createLPList(values interface{}, lengthPrefixBytes int) *LengthPrefixedList {
	if lengthPrefixBytes == 0 {
		lengthPrefixBytes = 4
	}

	lpList := &LengthPrefixedList{
		LengthPrefixBytes: lengthPrefixBytes,
		Values:            make([]interface{}, 0),
	}

	switch v := values.(type) {
	case []PostCondition:
		for _, item := range v {
			lpList.Values = append(lpList.Values, item)
		}
	case []TransactionAuthField:
		for _, item := range v {
			lpList.Values = append(lpList.Values, item)
		}
	case []StacksMessage:
		for _, item := range v {
			lpList.Values = append(lpList.Values, item)
		}
	default:
		panic("Unsupported type for createLPList")
	}

	return lpList
}

type StacksMessage interface {
	Type() StacksMessageType
	Serialize() ([]byte, error)
}
