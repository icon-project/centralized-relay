package clarity

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strings"
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
	ClarityTypeInt ClarityType = iota
	ClarityTypeUInt
	ClarityTypeBuffer
	ClarityTypeBoolTrue
	ClarityTypeBoolFalse
	ClarityTypeStringASCII
	ClarityTypeStringUTF8
	ClarityTypeList
	ClarityTypeTuple
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

type StringASCII struct {
	Data string
}

func NewStringASCII(data string) *StringASCII {
	return &StringASCII{Data: data}
}

func (s *StringASCII) Type() ClarityType {
	return ClarityTypeStringASCII
}

func (s *StringASCII) Serialize() ([]byte, error) {
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(s.Data)))
	return append(append([]byte{byte(ClarityTypeStringASCII)}, lenBytes...), []byte(s.Data)...), nil
}

type StringUTF8 struct {
	Data string
}

func NewStringUTF8(data string) *StringUTF8 {
	return &StringUTF8{Data: data}
}

func (s *StringUTF8) Type() ClarityType {
	return ClarityTypeStringUTF8
}

func (s *StringUTF8) Serialize() ([]byte, error) {
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(s.Data)))
	return append(append([]byte{byte(ClarityTypeStringUTF8)}, lenBytes...), []byte(s.Data)...), nil
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
	case ClarityTypeStringASCII:
		if len(data) < 5 {
			return nil, fmt.Errorf("invalid string ASCII data length")
		}
		length := binary.BigEndian.Uint32(data[1:5])
		if len(data) < int(5+length) {
			return nil, fmt.Errorf("invalid string ASCII length")
		}
		return &StringASCII{Data: string(data[5 : 5+length])}, nil
	case ClarityTypeStringUTF8:
		if len(data) < 5 {
			return nil, fmt.Errorf("invalid string UTF8 data length")
		}
		length := binary.BigEndian.Uint32(data[1:5])
		if len(data) < int(5+length) {
			return nil, fmt.Errorf("invalid string UTF8 length")
		}
		return &StringUTF8{Data: string(data[5 : 5+length])}, nil
	case ClarityTypeList:
		return deserializeList(data[1:])
	case ClarityTypeTuple:
		return deserializeTuple(data[1:])
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
