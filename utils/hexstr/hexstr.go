package hexstr

import (
	"encoding/hex"
	"strings"
)

const (
	HexPrefix string = "0x"
)

type HexString string

func NewFromString(val string) HexString {
	if strings.HasPrefix(val, HexPrefix) {
		return HexString(val)
	} else {
		return HexString(HexPrefix + val)
	}
}

func NewFromByte(val []byte) HexString {
	return HexString(HexPrefix + hex.EncodeToString(val))
}

func (h HexString) ToByte() ([]byte, error) {
	if h == "" {
		return nil, nil
	}
	cleanString := strings.TrimPrefix(string(h), HexPrefix)
	return hex.DecodeString(cleanString)
}
