package clarity

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	CLARITY_INT_SIZE      = 128
	CLARITY_INT_BYTE_SIZE = CLARITY_INT_SIZE / 8
)

var (
	MaxInt128  = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 127), big.NewInt(1))
	MinInt128  = new(big.Int).Neg(new(big.Int).Lsh(big.NewInt(1), 127))
	MaxUint128 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1))
)

type ClarityType byte

const (
	ClarityTypeInt               ClarityType = 0x00
	ClarityTypeUInt              ClarityType = 0x01
	ClarityTypeBuffer            ClarityType = 0x02
	ClarityTypeBoolTrue          ClarityType = 0x03
	ClarityTypeBoolFalse         ClarityType = 0x04
	ClarityTypeStandardPrincipal ClarityType = 0x05
	ClarityTypeContractPrincipal ClarityType = 0x06
	ClarityTypeResponseOk        ClarityType = 0x07
	ClarityTypeResponseErr       ClarityType = 0x08
	ClarityTypeOptionNone        ClarityType = 0x09
	ClarityTypeOptionSome        ClarityType = 0x0a
	ClarityTypeList              ClarityType = 0x0b
	ClarityTypeTuple             ClarityType = 0x0c
	ClarityTypeStringASCII       ClarityType = 0x0d
	ClarityTypeStringUTF8        ClarityType = 0x0e
)

type ClarityValue interface {
	Type() ClarityType
	Serialize() ([]byte, error)
}

type Int struct {
	Value *big.Int
}

func NewInt(value interface{}) (*Int, error) {
	bigInt, err := toBigInt(value, true)
	if err != nil {
		return nil, err
	}
	if bigInt.Cmp(MaxInt128) > 0 || bigInt.Cmp(MinInt128) < 0 {
		return nil, fmt.Errorf("value out of range for 128-bit signed integer")
	}
	if bigInt.Sign() == 0 {
		return &Int{Value: &big.Int{}}, nil
	}
	return &Int{Value: bigInt}, nil
}

func (i *Int) Type() ClarityType {
	return ClarityTypeInt
}

func (i *Int) Serialize() ([]byte, error) {
	bytes := make([]byte, CLARITY_INT_BYTE_SIZE)
	twosComplement := new(big.Int).Set(i.Value)
	if i.Value.Sign() < 0 {
		twosComplement.Add(twosComplement, new(big.Int).Lsh(big.NewInt(1), CLARITY_INT_SIZE))
	}
	twosComplement.FillBytes(bytes)
	return append([]byte{byte(ClarityTypeInt)}, bytes...), nil
}

type UInt struct {
	Value *big.Int
}

func NewUInt(value interface{}) (*UInt, error) {
	bigInt, err := toBigInt(value, false)
	if err != nil {
		return nil, err
	}
	if bigInt.Cmp(big.NewInt(0)) < 0 || bigInt.Cmp(MaxUint128) > 0 {
		return nil, fmt.Errorf("value out of range for 128-bit unsigned integer")
	}
	if bigInt.Sign() == 0 {
		return &UInt{Value: &big.Int{}}, nil
	}
	return &UInt{Value: bigInt}, nil
}

func (u *UInt) Type() ClarityType {
	return ClarityTypeUInt
}

func (u *UInt) Serialize() ([]byte, error) {
	bytes := make([]byte, CLARITY_INT_BYTE_SIZE)
	u.Value.FillBytes(bytes)
	return append([]byte{byte(ClarityTypeUInt)}, bytes...), nil
}

func toBigInt(value interface{}, signed bool) (*big.Int, error) {
	switch v := value.(type) {
	case int:
		return big.NewInt(int64(v)), nil
	case int64:
		return big.NewInt(v), nil
	case uint64:
		return new(big.Int).SetUint64(v), nil
	case string:
		if strings.HasPrefix(v, "0x") {
			bigInt, success := new(big.Int).SetString(v[2:], 16)
			if !success {
				return nil, fmt.Errorf("invalid hex string: %s", v)
			}
			if signed && len(v) > 2 && v[2] >= '8' {
				bigInt.Sub(bigInt, new(big.Int).Lsh(big.NewInt(1), uint(len(v[2:])*4)))
			}
			return bigInt, nil
		}
		bigInt, success := new(big.Int).SetString(v, 10)
		if !success {
			return nil, fmt.Errorf("invalid integer string: %s", v)
		}
		return bigInt, nil
	case []byte:
		if len(v) > CLARITY_INT_BYTE_SIZE {
			return nil, fmt.Errorf("byte array too long for 128-bit integer")
		}
		bigInt := new(big.Int).SetBytes(v)
		if signed && len(v) > 0 && v[0]&0x80 != 0 {
			bigInt.Sub(bigInt, new(big.Int).Lsh(big.NewInt(1), uint(len(v)*8)))
		}
		return bigInt, nil
	case *big.Int:
		return new(big.Int).Set(v), nil
	default:
		return nil, fmt.Errorf("unsupported type for conversion to bigint: %T", value)
	}
}

type Buffer struct {
	Data []byte
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{Data: data}
}

func (b *Buffer) Type() ClarityType {
	return ClarityTypeBuffer
}

func (b *Buffer) Serialize() ([]byte, error) {
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(b.Data)))
	return append(append([]byte{byte(ClarityTypeBuffer)}, lenBytes...), b.Data...), nil
}

type Bool bool

func NewBool(value bool) Bool {
	return Bool(value)
}

func (b Bool) Type() ClarityType {
	if b {
		return ClarityTypeBoolTrue
	}
	return ClarityTypeBoolFalse
}

func (b Bool) Serialize() ([]byte, error) {
	if b {
		return []byte{byte(ClarityTypeBoolTrue)}, nil
	}
	return []byte{byte(ClarityTypeBoolFalse)}, nil
}

func (b Bool) Bool() bool {
	return bool(b)
}

type StandardPrincipal struct {
	Version byte
	Hash160 [20]byte
}

func (p *StandardPrincipal) Serialize() ([]byte, error) {
	buf := make([]byte, 22)
	buf[0] = byte(ClarityTypeStandardPrincipal)
	buf[1] = p.Version
	copy(buf[2:], p.Hash160[:])
	return buf, nil
}

func NewStandardPrincipal(version byte, hash160 [20]byte) *StandardPrincipal {
	return &StandardPrincipal{Version: version, Hash160: hash160}
}

func (p *StandardPrincipal) Type() ClarityType {
	return ClarityTypeStandardPrincipal
}

type ContractPrincipal struct {
	StandardPrincipal
	ContractName string
}

func NewContractPrincipal(version byte, hash160 [20]byte, contractName string) (*ContractPrincipal, error) {
	if len(contractName) > 128 {
		return nil, fmt.Errorf("contract name too long (max 128 characters)")
	}
	return &ContractPrincipal{
		StandardPrincipal: StandardPrincipal{Version: version, Hash160: hash160},
		ContractName:      contractName,
	}, nil
}

func (p *ContractPrincipal) Type() ClarityType {
	return ClarityTypeContractPrincipal
}

func (p *ContractPrincipal) Serialize() ([]byte, error) {
	result := []byte{byte(ClarityTypeContractPrincipal), p.Version}
	result = append(result, p.Hash160[:]...)
	result = append(result, byte(len(p.ContractName)))
	result = append(result, []byte(p.ContractName)...)
	return result, nil
}

type ResponseOk struct {
	Value ClarityValue
}

func NewResponseOk(value ClarityValue) *ResponseOk {
	return &ResponseOk{Value: value}
}

func (r *ResponseOk) Type() ClarityType {
	return ClarityTypeResponseOk
}

func (r *ResponseOk) Serialize() ([]byte, error) {
	valueBytes, err := r.Value.Serialize()
	if err != nil {
		return nil, err
	}
	return append([]byte{byte(ClarityTypeResponseOk)}, valueBytes...), nil
}

type ResponseErr struct {
	Value ClarityValue
}

func NewResponseErr(value ClarityValue) *ResponseErr {
	return &ResponseErr{Value: value}
}

func (r *ResponseErr) Type() ClarityType {
	return ClarityTypeResponseErr
}

func (r *ResponseErr) Serialize() ([]byte, error) {
	valueBytes, err := r.Value.Serialize()
	if err != nil {
		return nil, err
	}
	return append([]byte{byte(ClarityTypeResponseErr)}, valueBytes...), nil
}

type OptionNone struct{}

func NewOptionNone() *OptionNone {
	return &OptionNone{}
}

func (o *OptionNone) Type() ClarityType {
	return ClarityTypeOptionNone
}

func (o *OptionNone) Serialize() ([]byte, error) {
	return []byte{byte(ClarityTypeOptionNone)}, nil
}

type OptionSome struct {
	Value ClarityValue
}

func NewOptionSome(value ClarityValue) *OptionSome {
	return &OptionSome{Value: value}
}

func (o *OptionSome) Type() ClarityType {
	return ClarityTypeOptionSome
}

func (o *OptionSome) Serialize() ([]byte, error) {
	valueBytes, err := o.Value.Serialize()
	if err != nil {
		return nil, err
	}
	return append([]byte{byte(ClarityTypeOptionSome)}, valueBytes...), nil
}

type List struct {
	Values []ClarityValue
}

func NewList(values []ClarityValue) *List {
	return &List{Values: values}
}

func (l *List) Type() ClarityType {
	return ClarityTypeList
}

func (l *List) Serialize() ([]byte, error) {
	result := []byte{byte(ClarityTypeList)}
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(l.Values)))
	result = append(result, lenBytes...)

	for _, value := range l.Values {
		serialized, err := value.Serialize()
		if err != nil {
			return nil, fmt.Errorf("error serializing list item: %w", err)
		}
		result = append(result, serialized...)
	}

	return result, nil
}

type Tuple struct {
	Data map[string]ClarityValue
}

func NewTuple(data map[string]ClarityValue) *Tuple {
	return &Tuple{Data: data}
}

func (t *Tuple) Type() ClarityType {
	return ClarityTypeTuple
}

func (t *Tuple) Serialize() ([]byte, error) {
	result := []byte{byte(ClarityTypeTuple)}
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(t.Data)))
	result = append(result, lenBytes...)

	keys := make([]string, 0, len(t.Data))
	for k := range t.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := t.Data[k]
		keyBytes := []byte(k)
		result = append(result, byte(len(keyBytes)))
		result = append(result, keyBytes...)
		serialized, err := v.Serialize()
		if err != nil {
			return nil, fmt.Errorf("error serializing tuple value: %w", err)
		}
		result = append(result, serialized...)
	}
	return result, nil
}

type StringASCII struct {
	Data string
}

func NewStringASCII(data string) (*StringASCII, error) {
	if err := validateASCII(data); err != nil {
		return nil, err
	}
	return &StringASCII{Data: data}, nil
}

func (s *StringASCII) Type() ClarityType {
	return ClarityTypeStringASCII
}

func (s *StringASCII) Serialize() ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("cannot serialize nil StringASCII")
	}
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(s.Data)))
	return append(append([]byte{byte(ClarityTypeStringASCII)}, lenBytes...), []byte(s.Data)...), nil
}

func validateASCII(s string) error {
	for i, r := range s {
		if r > 127 || !(unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsPunct(r) || unicode.IsSpace(r)) {
			return fmt.Errorf("invalid character in ASCII string at position %d: %q", i, r)
		}
	}
	return nil
}

type StringUTF8 struct {
	Data string
}

func NewStringUTF8(data string) (*StringUTF8, error) {
	if err := validateUTF8(data); err != nil {
		return nil, err
	}
	return &StringUTF8{Data: data}, nil
}

func (s *StringUTF8) Type() ClarityType {
	return ClarityTypeStringUTF8
}

func (s *StringUTF8) Serialize() ([]byte, error) {
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(s.Data)))
	return append(append([]byte{byte(ClarityTypeStringUTF8)}, lenBytes...), []byte(s.Data)...), nil
}

func validateUTF8(s string) error {
	if !utf8.ValidString(s) {
		return fmt.Errorf("invalid UTF-8 string")
	}
	return nil
}

func DeserializeClarityValue(data []byte) (ClarityValue, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	switch ClarityType(data[0]) {
	case ClarityTypeInt:
		if len(data) < CLARITY_INT_BYTE_SIZE+1 {
			return nil, fmt.Errorf("invalid int data length")
		}
		value := new(big.Int).SetBytes(data[1 : CLARITY_INT_BYTE_SIZE+1])
		if data[1]&0x80 != 0 {
			value.Sub(value, new(big.Int).Lsh(big.NewInt(1), CLARITY_INT_SIZE))
		}
		return &Int{Value: value}, nil
	case ClarityTypeUInt:
		if len(data) < CLARITY_INT_BYTE_SIZE+1 {
			return nil, fmt.Errorf("invalid uint data length")
		}
		value := new(big.Int).SetBytes(data[1 : CLARITY_INT_BYTE_SIZE+1])
		return &UInt{Value: value}, nil
	case ClarityTypeBuffer:
		if len(data) < 5 {
			return nil, fmt.Errorf("invalid buffer data length")
		}
		length := binary.BigEndian.Uint32(data[1:5])
		if len(data) < int(5+length) {
			return nil, fmt.Errorf("invalid buffer length")
		}
		return &Buffer{Data: data[5 : 5+length]}, nil
	case ClarityTypeBoolTrue:
		return NewBool(true), nil
	case ClarityTypeBoolFalse:
		return NewBool(false), nil
	case ClarityTypeStandardPrincipal:
		if len(data) < 22 {
			return nil, fmt.Errorf("invalid standard principal data length")
		}
		var hash160 [20]byte
		copy(hash160[:], data[2:22])
		return NewStandardPrincipal(data[1], hash160), nil
	case ClarityTypeContractPrincipal:
		if len(data) < 23 {
			return nil, fmt.Errorf("invalid contract principal data length")
		}
		var hash160 [20]byte
		copy(hash160[:], data[2:22])
		contractNameLength := int(data[22])
		if len(data) < 23+contractNameLength {
			return nil, fmt.Errorf("invalid contract principal name length")
		}
		contractName := string(data[23 : 23+contractNameLength])
		return NewContractPrincipal(data[1], hash160, contractName)
	case ClarityTypeResponseOk:
		value, err := DeserializeClarityValue(data[1:])
		if err != nil {
			return nil, fmt.Errorf("error deserializing ResponseOk value: %w", err)
		}
		return NewResponseOk(value), nil
	case ClarityTypeResponseErr:
		value, err := DeserializeClarityValue(data[1:])
		if err != nil {
			return nil, fmt.Errorf("error deserializing ResponseErr value: %w", err)
		}
		return NewResponseErr(value), nil
	case ClarityTypeOptionNone:
		return NewOptionNone(), nil
	case ClarityTypeOptionSome:
		value, err := DeserializeClarityValue(data[1:])
		if err != nil {
			return nil, fmt.Errorf("error deserializing OptionSome value: %w", err)
		}
		return NewOptionSome(value), nil
	case ClarityTypeList:
		return deserializeList(data[1:])
	case ClarityTypeTuple:
		return deserializeTuple(data[1:])
	case ClarityTypeStringASCII:
		if len(data) < 5 {
			return nil, fmt.Errorf("invalid string ASCII data length")
		}
		length := binary.BigEndian.Uint32(data[1:5])
		if len(data) < int(5+length) {
			return nil, fmt.Errorf("invalid string ASCII length")
		}
		str := string(data[5 : 5+length])
		if err := validateASCII(str); err != nil {
			return nil, err
		}
		return &StringASCII{Data: str}, nil
	case ClarityTypeStringUTF8:
		if len(data) < 5 {
			return nil, fmt.Errorf("invalid string UTF8 data length")
		}
		length := binary.BigEndian.Uint32(data[1:5])
		if len(data) < int(5+length) {
			return nil, fmt.Errorf("invalid string UTF8 length")
		}
		str := string(data[5 : 5+length])
		if err := validateUTF8(str); err != nil {
			return nil, err
		}
		return &StringUTF8{Data: str}, nil
	default:
		return nil, fmt.Errorf("unknown Clarity type: %d", data[0])
	}
}

func deserializeList(data []byte) (*List, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("invalid list data: too short")
	}

	length := binary.BigEndian.Uint32(data[:4])
	data = data[4:]

	values := make([]ClarityValue, 0, length)
	for i := uint32(0); i < length; i++ {
		if len(data) == 0 {
			return nil, fmt.Errorf("invalid list data: unexpected end")
		}

		value, err := DeserializeClarityValue(data)
		if err != nil {
			return nil, fmt.Errorf("error deserializing list item: %w", err)
		}

		values = append(values, value)
		serialized, _ := value.Serialize()
		data = data[len(serialized):]
	}

	return NewList(values), nil
}

func deserializeTuple(data []byte) (*Tuple, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("invalid tuple data: too short")
	}

	length := binary.BigEndian.Uint32(data[:4])
	data = data[4:]

	tupleData := make(map[string]ClarityValue)
	for i := uint32(0); i < length; i++ {
		if len(data) == 0 {
			return nil, fmt.Errorf("invalid tuple data: unexpected end")
		}

		keyLen := int(data[0])
		data = data[1:]
		if len(data) < keyLen {
			return nil, fmt.Errorf("invalid tuple data: key length exceeds remaining data")
		}

		key := string(data[:keyLen])
		data = data[keyLen:]

		value, err := DeserializeClarityValue(data)
		if err != nil {
			return nil, fmt.Errorf("error deserializing tuple value: %w", err)
		}

		tupleData[key] = value
		serialized, _ := value.Serialize()
		data = data[len(serialized):]
	}

	return NewTuple(tupleData), nil
}

func HexToClarityValue(hexStr string) (ClarityValue, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")

	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %v", err)
	}

	return DeserializeClarityValue(data)
}
