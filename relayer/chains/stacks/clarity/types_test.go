package clarity

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializeThenDeserialize(t *testing.T) {
	t.Run("IntCV", func(t *testing.T) {
		testCases := []struct {
			name  string
			value *big.Int
		}{
			{"Zero", big.NewInt(0)},
			{"Positive", big.NewInt(12345)},
			{"Negative", big.NewInt(-67890)},
			{"Max", new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 127), big.NewInt(1))},
			{"Min", new(big.Int).Neg(new(big.Int).Lsh(big.NewInt(1), 127))},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				cv, err := NewInt(tc.value)
				assert.NoError(t, err)

				serialized, err := cv.Serialize()
				assert.NoError(t, err)

				deserialized, err := DeserializeClarityValue(serialized)
				assert.NoError(t, err)
				deserializedInt, ok := deserialized.(*Int)
				assert.True(t, ok)
				assert.Equal(t, 0, cv.Value.Cmp(deserializedInt.Value))
			})
		}
	})

	t.Run("UIntCV", func(t *testing.T) {
		testCases := []struct {
			name  string
			value *big.Int
		}{
			{"Zero", big.NewInt(0)},
			{"Positive", big.NewInt(12345)},
			{"Max", new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1))},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				cv, err := NewUInt(tc.value)
				assert.NoError(t, err)

				serialized, err := cv.Serialize()
				assert.NoError(t, err)

				deserialized, err := DeserializeClarityValue(serialized)
				assert.NoError(t, err)
				deserializedUInt, ok := deserialized.(*UInt)
				assert.True(t, ok)
				assert.Equal(t, 0, cv.Value.Cmp(deserializedUInt.Value))
			})
		}
	})

	t.Run("BufferCV", func(t *testing.T) {
		buffer := []byte("this is a test")
		cv := NewBuffer(buffer)

		serialized, err := cv.Serialize()
		assert.NoError(t, err)

		deserialized, err := DeserializeClarityValue(serialized)
		assert.NoError(t, err)
		deserializedBuffer, ok := deserialized.(*Buffer)
		assert.True(t, ok)
		assert.Equal(t, cv.Data, deserializedBuffer.Data)
	})

	t.Run("StandardPrincipal", func(t *testing.T) {
		cv, err := StringToPrincipal("SP2JXKMSH007NPYAQHKJPQMAQYAD90NQGTVJVQ02B")
		assert.NoError(t, err)

		serialized, err := cv.Serialize()
		assert.NoError(t, err)

		deserialized, err := DeserializeClarityValue(serialized)
		assert.NoError(t, err)
		deserializedPrincipal, ok := deserialized.(*StandardPrincipal)
		assert.True(t, ok)
		assert.Equal(t, cv, deserializedPrincipal)
	})

	t.Run("ContractPrincipal", func(t *testing.T) {
		cv, err := StringToPrincipal("SP2JXKMSH007NPYAQHKJPQMAQYAD90NQGTVJVQ02B.test-contract")
		assert.NoError(t, err)

		serialized, err := cv.Serialize()
		assert.NoError(t, err)

		deserialized, err := DeserializeClarityValue(serialized)
		assert.NoError(t, err)
		deserializedContractPrincipal, ok := deserialized.(*ContractPrincipal)
		assert.True(t, ok)
		assert.Equal(t, cv, deserializedContractPrincipal)
	})

	t.Run("ListCV", func(t *testing.T) {
		cv := NewList([]ClarityValue{
			&Int{Value: big.NewInt(1)},
			&Int{Value: big.NewInt(2)},
			&Int{Value: big.NewInt(3)},
		})

		serialized, err := cv.Serialize()
		assert.NoError(t, err)

		deserialized, err := DeserializeClarityValue(serialized)
		assert.NoError(t, err)
		deserializedList, ok := deserialized.(*List)
		assert.True(t, ok)
		assert.Equal(t, cv, deserializedList)
	})

	t.Run("TupleCV", func(t *testing.T) {
		cv := NewTuple(map[string]ClarityValue{
			"a": &Int{Value: big.NewInt(1)},
			"b": &UInt{Value: big.NewInt(2)},
			"c": NewBuffer([]byte("test")),
		})

		serialized, err := cv.Serialize()
		assert.NoError(t, err)

		deserialized, err := DeserializeClarityValue(serialized)
		assert.NoError(t, err)
		deserializedTuple, ok := deserialized.(*Tuple)
		assert.True(t, ok)
		assert.Equal(t, cv, deserializedTuple)
	})
}

func TestSerializationTestVectors(t *testing.T) {
	t.Run("Int 1 Vector", func(t *testing.T) {
		cv := &Int{Value: big.NewInt(1)}
		serialized, err := cv.Serialize()
		assert.NoError(t, err)
		assert.Equal(t, "0000000000000000000000000000000001", hex.EncodeToString(serialized))
	})

	t.Run("Int -1 Vector", func(t *testing.T) {
		cv := &Int{Value: big.NewInt(-1)}
		serialized, err := cv.Serialize()
		assert.NoError(t, err)
		assert.Equal(t, "00ffffffffffffffffffffffffffffffff", hex.EncodeToString(serialized))
	})

	t.Run("UInt 1 Vector", func(t *testing.T) {
		cv := &UInt{Value: big.NewInt(1)}
		serialized, err := cv.Serialize()
		assert.NoError(t, err)
		assert.Equal(t, "0100000000000000000000000000000001", hex.EncodeToString(serialized))
	})

	t.Run("Buffer Vector", func(t *testing.T) {
		cv := NewBuffer([]byte{0xde, 0xad, 0xbe, 0xef})
		serialized, err := cv.Serialize()
		assert.NoError(t, err)
		assert.Equal(t, "0200000004deadbeef", hex.EncodeToString(serialized))
	})

	t.Run("True Vector", func(t *testing.T) {
		cv := Bool(true)
		serialized, err := cv.Serialize()
		assert.NoError(t, err)
		assert.Equal(t, "03", hex.EncodeToString(serialized))
	})

	t.Run("False Vector", func(t *testing.T) {
		cv := Bool(false)
		serialized, err := cv.Serialize()
		assert.NoError(t, err)
		assert.Equal(t, "04", hex.EncodeToString(serialized))
	})

	t.Run("Standard Principal Vector", func(t *testing.T) {
		addressBytes, _ := hex.DecodeString("11deadbeef11ababffff11deadbeef11ababffff")
		var hash160 [20]byte
		copy(hash160[:], addressBytes)
		cv := &StandardPrincipal{Version: 0x00, Hash160: hash160}
		serialized, err := cv.Serialize()
		assert.NoError(t, err)
		assert.Equal(t, "050011deadbeef11ababffff11deadbeef11ababffff", hex.EncodeToString(serialized))
	})

	t.Run("Standard Principal Vector 2", func(t *testing.T) {
		cv, err := StringToPrincipal("ST1PQHQKV0RJXZFY1DGX8MNSNYVE3VGZJSRTPGZGM")
		assert.NoError(t, err)
		serialized, err := cv.Serialize()
		assert.NoError(t, err)
		assert.Equal(t, "051a6d78de7b0625dfbfc16c3a8a5735f6dc3dc3f2ce", hex.EncodeToString(serialized))
	})

	t.Run("Contract Principal Vector", func(t *testing.T) {
		addressBytes, _ := hex.DecodeString("11deadbeef11ababffff11deadbeef11ababffff")
		var hash160 [20]byte
		copy(hash160[:], addressBytes)
		cv, err := NewContractPrincipal(0x00, hash160, "abcd")
		assert.NoError(t, err)
		serialized, err := cv.Serialize()
		assert.NoError(t, err)
		assert.Equal(t, "060011deadbeef11ababffff11deadbeef11ababffff0461626364", hex.EncodeToString(serialized))
	})

	t.Run("List Vector", func(t *testing.T) {
		list := NewList([]ClarityValue{
			&Int{Value: big.NewInt(1)},
			&Int{Value: big.NewInt(2)},
			&Int{Value: big.NewInt(3)},
			&Int{Value: big.NewInt(-4)},
		})
		serialized, err := list.Serialize()
		assert.NoError(t, err)
		expected := "0b0000000400000000000000000000000000000000010000000000000000000000000000000002000000000000000000000000000000000300fffffffffffffffffffffffffffffffc"
		assert.Equal(t, expected, hex.EncodeToString(serialized))
	})

	t.Run("Tuple Vector", func(t *testing.T) {
		tuple := NewTuple(map[string]ClarityValue{
			"baz":    &OptionNone{},
			"foobar": Bool(true),
		})
		serialized, err := tuple.Serialize()
		assert.NoError(t, err)
		assert.Equal(t, "0c000000020362617a0906666f6f62617203", hex.EncodeToString(serialized))
	})

	t.Run("StringAscii Vector", func(t *testing.T) {
		cv, err := NewStringASCII("hello world")
		assert.NoError(t, err)
		serialized, err := cv.Serialize()
		assert.NoError(t, err)
		assert.Equal(t, "0d0000000b68656c6c6f20776f726c64", hex.EncodeToString(serialized))
	})

	t.Run("StringUtf8 Vector", func(t *testing.T) {
		cv, err := NewStringUTF8("hello world")
		assert.NoError(t, err)
		serialized, err := cv.Serialize()
		assert.NoError(t, err)
		assert.Equal(t, "0e0000000b68656c6c6f20776f726c64", hex.EncodeToString(serialized))
	})
}

func TestIntBounds(t *testing.T) {
	t.Run("Max Int", func(t *testing.T) {
		maxInt := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 127), big.NewInt(1))
		cv, err := NewInt(maxInt)
		assert.NoError(t, err)
		assert.Equal(t, maxInt, cv.Value)
	})

	t.Run("Min Int", func(t *testing.T) {
		minInt := new(big.Int).Neg(new(big.Int).Lsh(big.NewInt(1), 127))
		cv, err := NewInt(minInt)
		assert.NoError(t, err)
		assert.Equal(t, minInt, cv.Value)
	})

	t.Run("Overflow Max Int", func(t *testing.T) {
		overflowInt := new(big.Int).Add(new(big.Int).Lsh(big.NewInt(1), 127), big.NewInt(1))
		_, err := NewInt(overflowInt)
		assert.Error(t, err)
	})

	t.Run("Underflow Min Int", func(t *testing.T) {
		underflowInt := new(big.Int).Neg(new(big.Int).Add(new(big.Int).Lsh(big.NewInt(1), 127), big.NewInt(1)))
		_, err := NewInt(underflowInt)
		assert.Error(t, err)
	})
}

func TestUIntBounds(t *testing.T) {
	t.Run("Max UInt", func(t *testing.T) {
		maxUInt := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1))
		cv, err := NewUInt(maxUInt)
		assert.NoError(t, err)
		assert.Equal(t, maxUInt, cv.Value)
	})

	t.Run("Overflow Max UInt", func(t *testing.T) {
		overflowUInt := new(big.Int).Lsh(big.NewInt(1), 128)
		_, err := NewUInt(overflowUInt)
		assert.Error(t, err)
	})

	t.Run("Negative UInt", func(t *testing.T) {
		negativeUInt := big.NewInt(-1)
		_, err := NewUInt(negativeUInt)
		assert.Error(t, err)
	})
}
