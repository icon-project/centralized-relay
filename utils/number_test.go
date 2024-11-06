package utils

import (
	"math/big"
	"testing"
)

func TestToTruncatedBE(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"0", []byte{}},
		{"1", []byte{0x01}},
		{"255", []byte{0xFF}},
		{"256", []byte{0x01, 0x00}},
		{"65535", []byte{0xFF, 0xFF}},
		{"4294967295", []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, test := range tests {
		num := new(big.Int)
		num.SetString(test.input, 10)

		result := ToTruncatedBE(num)
		if !equal(result, test.expected) {
			t.Errorf("ToTruncatedBE(%s) = %x; expected %x", test.input, result, test.expected)
		}
	}
}

func TestToTruncatedLE(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"0", []byte{}},
		{"1", []byte{0x01}},
		{"255", []byte{0xFF}},
		{"256", []byte{0x00, 0x01}},
		{"65535", []byte{0xFF, 0xFF}},
		{"4294967295", []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, test := range tests {
		num := new(big.Int)
		num.SetString(test.input, 10)

		result := ToTruncatedLE(num)
		if !equal(result, test.expected) {
			t.Errorf("ToTruncatedLE(%s) = %x; expected %x", test.input, result, test.expected)
		}
	}
}

// Helper function to compare two byte slices
func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
