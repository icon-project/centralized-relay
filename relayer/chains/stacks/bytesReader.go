package stacks

import (
	"encoding/binary"
	"fmt"
	"math/big"
)

type BytesReader struct {
	source   []byte
	consumed int
}

func NewBytesReader(arr []byte) *BytesReader {
	return &BytesReader{
		source:   arr,
		consumed: 0,
	}
}

func (br *BytesReader) ReadBytes(length int) []byte {
	if br.consumed+length > len(br.source) {
		panic("Attempt to read past end of buffer")
	}
	view := br.source[br.consumed : br.consumed+length]
	br.consumed += length
	return view
}

func (br *BytesReader) ReadUint32BE() uint32 {
	return binary.BigEndian.Uint32(br.ReadBytes(4))
}

func (br *BytesReader) ReadUint8() uint8 {
	return br.ReadBytes(1)[0]
}

func (br *BytesReader) ReadUint16BE() uint16 {
	return binary.BigEndian.Uint16(br.ReadBytes(2))
}

func (br *BytesReader) ReadBigUintLE(length int) *big.Int {
	bytes := br.ReadBytes(length)
	for i := 0; i < len(bytes)/2; i++ {
		bytes[i], bytes[len(bytes)-1-i] = bytes[len(bytes)-1-i], bytes[i]
	}
	return new(big.Int).SetBytes(bytes)
}

func (br *BytesReader) ReadBigUintBE(length int) *big.Int {
	bytes := br.ReadBytes(length)
	return new(big.Int).SetBytes(bytes)
}

func (br *BytesReader) ReadOffset() int {
	return br.consumed
}

func (br *BytesReader) SetReadOffset(val int) {
	if val < 0 || val > len(br.source) {
		panic("Invalid read offset")
	}
	br.consumed = val
}

func (br *BytesReader) ReadUint8Enum(enumChecker func(uint8) bool) (uint8, error) {
	num := br.ReadUint8()
	if enumChecker(num) {
		return num, nil
	}
	return 0, fmt.Errorf("invalid enum value: %d", num)
}
