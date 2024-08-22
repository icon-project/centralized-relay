package clarity

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestCrockfordDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Decode 'hello world'",
			input:    "38CNP6RVS0EXQQ4V34",
			expected: "68656c6c6f20776f726c64",
			wantErr:  false,
		},
		{
			name:     "Decode 'ST1PQHQKV0RJXZFY1DGX8MNSNYVE3VGZJSRTPGZGM'",
			input:    "ST1PQHQKV0RJXZFY1DGX8MNSNYVE3VGZJSRTPGZGM",
			expected: "19d06d78de7b0625dfbfc16c3a8a5735f6dc3dc3f2ce35687e14",
			wantErr:  false,
		},
		{
			name:     "Decode 'SP2J6ZY48GV1EZ5V2V5RB9MP66SW86PYKKNRV9EJ7'",
			input:    "SP2J6ZY48GV1EZ5V2V5RB9MP66SW86PYKKNRV9EJ7",
			expected: "19b0a46ff88886c2ef9762d970b4d2c63678835bd39d71b4ba47",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoded, err := CrockfordDecode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CrockfordDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotHex := hex.EncodeToString(decoded)

			if !reflect.DeepEqual(gotHex, tt.expected) {
				t.Errorf("CrockfordDecode() = %v, want %v", gotHex, tt.expected)
			}
		})
	}
}

func TestDecodeC32Address(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedVer  byte
		expectedHash string
		wantErr      bool
	}{
		{
			name:         "Decode Stacks address",
			input:        "ST1PQHQKV0RJXZFY1DGX8MNSNYVE3VGZJSRTPGZGM",
			expectedVer:  26,
			expectedHash: "6d78de7b0625dfbfc16c3a8a5735f6dc3dc3f2ce",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, hash160, err := DecodeC32Address(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeC32Address() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if version != tt.expectedVer {
				t.Errorf("DecodeC32Address() version = %v, want %v", version, tt.expectedVer)
			}

			gotHash := hex.EncodeToString(hash160[:])
			if gotHash != tt.expectedHash {
				t.Errorf("DecodeC32Address() hash160 = %v, want %v", gotHash, tt.expectedHash)
			}
		})
	}
}
