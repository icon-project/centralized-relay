package clarity

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClarityInt(t *testing.T) {
	testCases := []struct {
		input    interface{}
		expected string
		hexValue string
	}{
		{1, "1", "0x00000000000000000000000000000001"},
		{-1, "-1", "0xffffffffffffffffffffffffffffffff"},
		{-10, "-10", "0xfffffffffffffffffffffffffffffff6"},
		{big.NewInt(-10), "-10", "0xfffffffffffffffffffffffffffffff6"},
		{"-10", "-10", "0xfffffffffffffffffffffffffffffff6"},
		{"0xfff6", "-10", "0xfffffffffffffffffffffffffffffff6"},
		{[]byte{0xff, 0xf6}, "-10", "0xfffffffffffffffffffffffffffffff6"},
		{200, "200", "0x000000000000000000000000000000c8"},
		{10, "10", "0x0000000000000000000000000000000a"},
		{"10", "10", "0x0000000000000000000000000000000a"},
		{"0x0a", "10", "0x0000000000000000000000000000000a"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			i, err := NewInt(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, i.Value.String())

			serialized, err := i.Serialize()
			require.NoError(t, err)
			assert.Equal(t, tc.hexValue, "0x"+hex.EncodeToString(serialized[1:]))

			deserialized, err := DeserializeClarityValue(serialized)
			require.NoError(t, err)
			assert.Equal(t, i.Value, deserialized.(*Int).Value)
		})
	}

	maxInt, err := NewInt(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 127), big.NewInt(1)))
	require.NoError(t, err)
	assert.Equal(t, "170141183460469231731687303715884105727", maxInt.Value.String())

	minInt, err := NewInt(new(big.Int).Neg(new(big.Int).Lsh(big.NewInt(1), 127)))
	require.NoError(t, err)
	assert.Equal(t, "-170141183460469231731687303715884105728", minInt.Value.String())

	_, err = NewInt(new(big.Int).Add(MaxInt128, big.NewInt(1)))
	assert.Error(t, err)

	_, err = NewInt(new(big.Int).Sub(MinInt128, big.NewInt(1)))
	assert.Error(t, err)
}

func TestClarityUInt(t *testing.T) {
	testCases := []struct {
		input    interface{}
		expected string
		hexValue string
	}{
		{200, "200", "0x000000000000000000000000000000c8"},
		{10, "10", "0x0000000000000000000000000000000a"},
		{"10", "10", "0x0000000000000000000000000000000a"},
		{"0x0a", "10", "0x0000000000000000000000000000000a"},
		{big.NewInt(10), "10", "0x0000000000000000000000000000000a"},
		{[]byte{0x0a}, "10", "0x0000000000000000000000000000000a"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			u, err := NewUInt(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, u.Value.String())

			serialized, err := u.Serialize()
			require.NoError(t, err)
			assert.Equal(t, tc.hexValue, "0x"+hex.EncodeToString(serialized[1:]))

			deserialized, err := DeserializeClarityValue(serialized)
			require.NoError(t, err)
			assert.Equal(t, u.Value, deserialized.(*UInt).Value)
		})
	}

	maxUint, err := NewUInt(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1)))
	require.NoError(t, err)
	assert.Equal(t, "340282366920938463463374607431768211455", maxUint.Value.String())

	minUint, err := NewUInt(0)
	require.NoError(t, err)
	assert.Equal(t, "0", minUint.Value.String())

	_, err = NewUInt(new(big.Int).Add(MaxUint128, big.NewInt(1)))
	assert.Error(t, err)

	_, err = NewUInt(-1)
	assert.Error(t, err)
}

func TestClarityBuffer(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected string
	}{
		{[]byte{0xde, 0xad, 0xbe, 0xef}, "deadbeef"},
		{[]byte{0x12, 0x34, 0x56, 0x78}, "12345678"},
	}

	for _, tc := range testCases {
		t.Run(hex.EncodeToString(tc.input), func(t *testing.T) {
			b := NewBuffer(tc.input)
			serialized, err := b.Serialize()
			require.NoError(t, err)

			assert.Equal(t, byte(ClarityTypeBuffer), serialized[0])
			length := uint32(len(tc.input))
			assert.Equal(t, length, binary.BigEndian.Uint32(serialized[1:5]))
			assert.Equal(t, tc.expected, hex.EncodeToString(serialized[5:]))

			deserialized, err := DeserializeClarityValue(serialized)
			require.NoError(t, err)
			assert.Equal(t, b.Data, deserialized.(*Buffer).Data)
		})
	}
}

func TestClarityBool(t *testing.T) {
	testCases := []struct {
		input    bool
		expected string
	}{
		{true, "03"},
		{false, "04"},
	}

	for _, tc := range testCases {
		b := NewBool(tc.input)
		serialized, err := b.Serialize()
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, hex.EncodeToString(serialized))

		deserialized, err := DeserializeClarityValue(serialized)
		assert.NoError(t, err)
		assert.Equal(t, b, deserialized.(Bool))
	}
}

func TestClarityStringASCII(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"hello world", "68656c6c6f20776f726c64"},
		{"Clarity", "436c6172697479"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			s := NewStringASCII(tc.input)
			serialized, err := s.Serialize()
			require.NoError(t, err)

			assert.Equal(t, byte(ClarityTypeStringASCII), serialized[0])
			length := uint32(len(tc.input))
			assert.Equal(t, length, binary.BigEndian.Uint32(serialized[1:5]))
			assert.Equal(t, tc.expected, hex.EncodeToString(serialized[5:]))

			deserialized, err := DeserializeClarityValue(serialized)
			require.NoError(t, err)
			assert.Equal(t, s.Data, deserialized.(*StringASCII).Data)
		})
	}
}

func TestClarityStringUTF8(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"hello world 🌍", "68656c6c6f20776f726c6420f09f8c8d"},
		{"Clarity 💡", "436c617269747920f09f92a1"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			s := NewStringUTF8(tc.input)
			serialized, err := s.Serialize()
			require.NoError(t, err)

			assert.Equal(t, byte(ClarityTypeStringUTF8), serialized[0])
			length := uint32(len(tc.input))
			assert.Equal(t, length, binary.BigEndian.Uint32(serialized[1:5]))
			assert.Equal(t, tc.expected, hex.EncodeToString(serialized[5:]))

			deserialized, err := DeserializeClarityValue(serialized)
			require.NoError(t, err)
			assert.Equal(t, s.Data, deserialized.(*StringUTF8).Data)
		})
	}
}

func TestClarityList(t *testing.T) {
	intValue1, _ := NewInt(1)
	intValue2, _ := NewInt(2)
	intValue3, _ := NewInt(3)
	intValue4, _ := NewInt(-4)

	list := NewList([]ClarityValue{intValue1, intValue2, intValue3, intValue4})

	serialized, err := list.Serialize()
	require.NoError(t, err)

	expected := "070000000400000000000000000000000000000000010000000000000000000000000000000002000000000000000000000000000000000300fffffffffffffffffffffffffffffffc"
	assert.Equal(t, expected, hex.EncodeToString(serialized))

	deserialized, err := DeserializeClarityValue(serialized)
	require.NoError(t, err)

	deserializedList, ok := deserialized.(*List)

	assert.True(t, ok)
	assert.Equal(t, 4, len(deserializedList.Values))
	assert.Equal(t, intValue1.Value, deserializedList.Values[0].(*Int).Value)
	assert.Equal(t, intValue2.Value, deserializedList.Values[1].(*Int).Value)
	assert.Equal(t, intValue3.Value, deserializedList.Values[2].(*Int).Value)
	assert.Equal(t, intValue4.Value, deserializedList.Values[3].(*Int).Value)
}

func TestClarityTuple(t *testing.T) {
	intValue, _ := NewInt(1)
	tuple := NewTuple(map[string]ClarityValue{
		"a": intValue,
		"b": NewBool(true),
		"c": NewStringASCII("hello"),
	})

	serialized, err := tuple.Serialize()
	require.NoError(t, err)

	deserialized, err := DeserializeClarityValue(serialized)
	require.NoError(t, err)

	deserializedTuple, ok := deserialized.(*Tuple)
	require.True(t, ok)
	assert.Equal(t, 3, len(deserializedTuple.Data))
	assert.Equal(t, intValue.Value, deserializedTuple.Data["a"].(*Int).Value)
	assert.Equal(t, true, deserializedTuple.Data["b"].(Bool).Bool())
	assert.Equal(t, "hello", deserializedTuple.Data["c"].(*StringASCII).Data)
}
