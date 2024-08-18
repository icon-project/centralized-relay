package stacks

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeriveStxPrivateKey(t *testing.T) {
	testCases := []struct {
		name     string
		mnemonic string
		index    uint32
		expected string
	}{
		{
			name:     "Test Vector 1",
			mnemonic: "sound idle panel often situate develop unit text design antenna vendor screen opinion balcony share trigger accuse scatter visa uniform brass update opinion media",
			index:    0,
			expected: "8721c6a5237f5e8d361161a7855aa56885a3e19e2ea6ee268fb14eabc5e2ed9001",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			derivedKey, err := DeriveStxPrivateKey(tc.mnemonic, tc.index)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, hex.EncodeToString(derivedKey))
		})
	}
}
