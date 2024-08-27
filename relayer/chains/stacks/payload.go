// relayer/chains/stacks/payload.go
package stacks

import (
	"bytes"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/clarity"
)

type Payload interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) (int, error)
}

type TokenTransferPayload struct {
	Recipient clarity.ClarityValue // Can be either StandardPrincipal or ContractPrincipal
	Amount    uint64
	Memo      string
}

type ContractCallPayload struct {
	ContractAddress string
	ContractName    string
	FunctionName    string
	FunctionArgs    []clarity.ClarityValue
}

func NewTokenTransferPayload(recipient string, amount uint64, memo string) (*TokenTransferPayload, error) {
	principalCV, err := clarity.StringToPrincipal(recipient)
	if err != nil {
		return nil, err
	}

	return &TokenTransferPayload{
		Recipient: principalCV,
		Amount:    amount,
		Memo:      memo,
	}, nil
}

func (p *TokenTransferPayload) Serialize() ([]byte, error) {
	buf := make([]byte, 0, 128)

	buf = append(buf, byte(PayloadTypeTokenTransfer))

	recipientBytes, err := p.Recipient.Serialize()
	if err != nil {
		return nil, err
	}
	buf = append(buf, recipientBytes...)

	amountBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(amountBytes, p.Amount)
	buf = append(buf, amountBytes...)

	memoBytes := make([]byte, MemoMaxLengthBytes)
	copy(memoBytes, p.Memo)
	buf = append(buf, memoBytes...)

	return buf, nil
}

func (p *TokenTransferPayload) Deserialize(data []byte) (int, error) {
	if len(data) < 1 || PayloadType(data[0]) != PayloadTypeTokenTransfer {
		return 0, errors.New("invalid token transfer payload")
	}

	offset := 1

	recipient, n, err := deserializePrincipal(data[offset:])
	if err != nil {
		return 0, err
	}
	p.Recipient = recipient
	offset += n

	if len(data[offset:]) < 8 {
		return 0, errors.New("insufficient data for amount")
	}
	p.Amount = binary.BigEndian.Uint64(data[offset : offset+8])
	offset += 8

	if len(data[offset:]) < MemoMaxLengthBytes {
		return 0, errors.New("insufficient data for memo")
	}
	p.Memo = string(bytes.TrimRight(data[offset:offset+MemoMaxLengthBytes], "\x00"))
	offset += MemoMaxLengthBytes

	return offset, nil
}

func (p *ContractCallPayload) Serialize() ([]byte, error) {
	buf := make([]byte, 0, 128)

	buf = append(buf, byte(PayloadTypeContractCall))

	contractAddressBytes, err := serializeAddress(p.ContractAddress)
	if err != nil {
		return nil, err
	}
	buf = append(buf, contractAddressBytes...)

	contractNameBytes, err := serializeString(p.ContractName, MaxStringLengthBytes)
	if err != nil {
		return nil, err
	}
	buf = append(buf, contractNameBytes...)

	functionNameBytes, err := serializeString(p.FunctionName, MaxStringLengthBytes)
	if err != nil {
		return nil, err
	}
	buf = append(buf, functionNameBytes...)

	argCountBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(argCountBytes, uint32(len(p.FunctionArgs)))
	buf = append(buf, argCountBytes...)

	for _, arg := range p.FunctionArgs {
		argBytes, err := arg.Serialize()
		if err != nil {
			return nil, err
		}
		buf = append(buf, argBytes...)
	}

	return buf, nil
}

func (p *ContractCallPayload) Deserialize(data []byte) (int, error) {
	if len(data) < 1 || PayloadType(data[0]) != PayloadTypeContractCall {
		return 0, errors.New("invalid contract call payload")
	}

	offset := 1

	contractAddress, n, err := deserializeAddress(data[offset:])
	if err != nil {
		return 0, err
	}
	p.ContractAddress = contractAddress
	offset += n

	contractName, n, err := deserializeString(data[offset:], MaxStringLengthBytes)
	if err != nil {
		return 0, err
	}
	p.ContractName = contractName
	offset += n

	functionName, n, err := deserializeString(data[offset:], MaxStringLengthBytes)
	if err != nil {
		return 0, err
	}
	p.FunctionName = functionName
	offset += n

	if len(data[offset:]) < 4 {
		return 0, errors.New("insufficient data for function args count")
	}
	argCount := binary.BigEndian.Uint32(data[offset : offset+4])
	offset += 4

	p.FunctionArgs = make([]clarity.ClarityValue, argCount)
	for i := uint32(0); i < argCount; i++ {
		arg, err := clarity.DeserializeClarityValue(data[offset:])
		if err != nil {
			return 0, err
		}
		p.FunctionArgs[i] = arg
		serialized, _ := arg.Serialize()
		offset += len(serialized)
	}

	return offset, nil
}

func deserializePrincipal(data []byte) (clarity.ClarityValue, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("insufficient data for principal")
	}

	principal, err := clarity.DeserializeClarityValue(data)
	if err != nil {
		return nil, 0, fmt.Errorf("error deserializing principal: %w", err)
	}

	serialized, _ := principal.Serialize()
	return principal, len(serialized), nil
}

func serializeAddress(address string) ([]byte, error) {
	if len(address) != 1+AddressHashLength*2 { // 'S' + version + 40 hex chars
		return nil, fmt.Errorf("invalid address length: %d", len(address))
	}

	var version AddressVersion
	switch address[0] {
	case 'S':
		version = AddressVersionMainnetSingleSig
	case 'T':
		version = AddressVersionTestnetSingleSig
	default:
		return nil, fmt.Errorf("invalid address version: %c", address[0])
	}

	hashBytes, err := c32Decode(address[1:])
	if err != nil {
		return nil, fmt.Errorf("invalid address hash: %v", err)
	}

	result := make([]byte, 1+len(hashBytes))
	result[0] = byte(version)
	copy(result[1:], hashBytes)

	return result, nil
}

func deserializeAddress(data []byte) (string, int, error) {
	if len(data) < 1+AddressHashLength {
		return "", 0, errors.New("insufficient data for address")
	}

	version := AddressVersion(data[0])
	var prefix string
	switch version {
	case AddressVersionMainnetSingleSig:
		prefix = "S"
	case AddressVersionTestnetSingleSig:
		prefix = "T"
	default:
		return "", 0, fmt.Errorf("invalid address version: %d", version)
	}

	c32hash := c32Encode(data[1 : 1+AddressHashLength+5])
	address := fmt.Sprintf("%s%s", prefix, c32hash)

	return address, 1 + AddressHashLength + 5, nil
}

func serializeString(s string, maxLength int) ([]byte, error) {
	if len(s) > maxLength {
		return nil, errors.New("string exceeds maximum length")
	}
	buf := make([]byte, 1+len(s))
	buf[0] = byte(len(s))
	copy(buf[1:], s)
	return buf, nil
}

func deserializeString(data []byte, maxLength int) (string, int, error) {
	if len(data) < 1 {
		return "", 0, errors.New("insufficient data for string length")
	}
	length := int(data[0])
	if length > maxLength || len(data) < 1+length {
		return "", 0, errors.New("invalid string length")
	}
	return string(data[1 : 1+length]), 1 + length, nil
}

func c32Encode(input []byte) string {
	alphabet := "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	encoder := base32.NewEncoding(alphabet).WithPadding(base32.NoPadding)
	return encoder.EncodeToString(input)
}

func c32Decode(input string) ([]byte, error) {
	alphabet := "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	decoder := base32.NewEncoding(alphabet).WithPadding(base32.NoPadding)
	return decoder.DecodeString(input)
}
