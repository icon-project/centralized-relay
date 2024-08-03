package stacks

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/big"

	"golang.org/x/crypto/sha3"
)

func Equals(a, b []byte) bool {
	return bytes.Equal(a, b)
}

func Alloc(length int, value byte) []byte {
	a := make([]byte, length)
	for i := 0; i < length; i++ {
		a[i] = value
	}
	return a
}

func ReadUint16BE(source []byte, offset int) uint16 {
	return binary.BigEndian.Uint16(source[offset:])
}

func WriteUint16BE(destination []byte, value uint16, offset int) []byte {
	binary.BigEndian.PutUint16(destination[offset:], value)
	return destination
}

func ReadUint8(source []byte, offset int) uint8 {
	return source[offset]
}

func WriteUint8(destination []byte, value uint8, offset int) []byte {
	destination[offset] = value
	return destination
}

func ReadUint16LE(source []byte, offset int) uint16 {
	return binary.LittleEndian.Uint16(source[offset:])
}

func WriteUint16LE(destination []byte, value uint16, offset int) []byte {
	binary.LittleEndian.PutUint16(destination[offset:], value)
	return destination
}

func ReadUint32BE(source []byte, offset int) uint32 {
	return binary.BigEndian.Uint32(source[offset:])
}

func WriteUint32BE(destination []byte, value uint32, offset int) []byte {
	binary.BigEndian.PutUint32(destination[offset:], value)
	return destination
}

func ReadUint32LE(source []byte, offset int) uint32 {
	return binary.LittleEndian.Uint32(source[offset:])
}

func WriteUint32LE(destination []byte, value uint32, offset int) []byte {
	binary.LittleEndian.PutUint32(destination[offset:], value)
	return destination
}

func TxIDFromData(data []byte) string {
	hash := sha512_256(data)
	return hex.EncodeToString(hash)
}

func sha512_256(data []byte) []byte {
	hash := sha3.New512()
	hash.Write(data)
	return hash.Sum(nil)[:32]
}

func With0x(value string) string {
	if len(value) >= 2 && value[:2] == "0x" {
		return value
	}
	return "0x" + value
}

func IntToBigInt(value interface{}, signed bool) *big.Int {
	var result *big.Int

	switch v := value.(type) {
	case int:
		result = big.NewInt(int64(v))
	case int64:
		result = big.NewInt(v)
	case uint64:
		result = new(big.Int).SetUint64(v)
	case string:
		result = new(big.Int)
		result.SetString(v, 10)
	case *big.Int:
		result = new(big.Int).Set(v)
	case []byte:
		if signed {
			result = new(big.Int).SetBytes(v)
			if v[0]&0x80 != 0 {
				// If the first bit is set, it's a negative number in two's complement
				result.Sub(result, new(big.Int).Lsh(big.NewInt(1), uint(len(v)*8)))
			}
		} else {
			result = new(big.Int).SetBytes(v)
		}
	default:
		panic("Unsupported type for IntToBigInt")
	}

	return result
}

func hexToBytes(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

func IntToBytes(value interface{}, signed bool, byteLength int) []byte {
	bigInt := IntToBigInt(value, signed)
	return BigIntToBytes(bigInt, byteLength)
}

func BigIntToBytes(value *big.Int, byteLength int) []byte {
	bytes := value.Bytes()
	if len(bytes) > byteLength {
		panic("BigInt too large for specified byte length")
	}
	result := make([]byte, byteLength)
	copy(result[byteLength-len(bytes):], bytes)
	return result
}

func ConcatBytes(arrays ...[]byte) []byte {
	var totalLen int
	for _, arr := range arrays {
		totalLen += len(arr)
	}
	result := make([]byte, totalLen)
	var i int
	for _, arr := range arrays {
		i += copy(result[i:], arr)
	}
	return result
}
