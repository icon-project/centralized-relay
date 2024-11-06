package utils

import (
	"math/big"
)

func ToTruncatedBE(num *big.Int) []byte {
	return num.Bytes()
}

func ToTruncatedLE(num *big.Int) []byte {
	beBytes := num.Bytes()
	leBytes := make([]byte, len(beBytes))
	for i := range beBytes {
		leBytes[i] = beBytes[len(beBytes)-1-i]
	}
	return leBytes
}
