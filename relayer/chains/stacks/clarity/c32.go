package clarity

import (
	"fmt"
	"math/big"
	"strings"
)

var crockfordAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

func CrockfordDecode(input string) ([]byte, error) {
	input = strings.ToUpper(strings.ReplaceAll(input, "-", ""))
	input = strings.ReplaceAll(input, "O", "0")
	input = strings.ReplaceAll(input, "I", "1")
	input = strings.ReplaceAll(input, "L", "1")

	bi := big.NewInt(0)
	for _, char := range input {
		bi.Mul(bi, big.NewInt(32))
		index := strings.IndexRune(crockfordAlphabet, char)
		if index == -1 {
			return nil, fmt.Errorf("invalid character: %c", char)
		}
		bi.Add(bi, big.NewInt(int64(index)))
	}

	bytes := bi.Bytes()

	for len(bytes) > 0 && bytes[0] == 0 {
		bytes = bytes[1:]
	}

	return bytes, nil
}

func DecodeC32Address(address string) (version byte, hash160 [20]byte, err error) {
	if len(address) < 5 || address[0] != 'S' {
		return 0, [20]byte{}, fmt.Errorf("invalid C32 address: must start with 'S' and be at least 5 characters long")
	}

	versionChar := address[1]
	version = byte(strings.IndexRune(crockfordAlphabet, rune(versionChar)))

	decoded, err := CrockfordDecode(address[2:])
	if err != nil {
		return 0, [20]byte{}, fmt.Errorf("failed to decode address: %v", err)
	}

	if len(decoded) != 24 { // 20 bytes hash160 + 4 bytes checksum
		return 0, [20]byte{}, fmt.Errorf("invalid decoded length: expected 24, got %d", len(decoded))
	}

	copy(hash160[:], decoded[:20])

	return version, hash160, nil
}
