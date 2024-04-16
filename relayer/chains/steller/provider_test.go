package steller

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stellar/go/strkey"
	"github.com/stretchr/testify/assert"
)

func TestBase58(t *testing.T) {
	hexString := "70e9bcd9996297b6e5efd00ed53dd3459d920c4d883023f9a661b513470f75af"

	// Decode hexadecimal string into bytes
	hexBytes, err := hex.DecodeString(hexString)
	assert.NoError(t, err)

	val, err := strkey.Encode(strkey.VersionByteContract, hexBytes)
	assert.NoError(t, err)

	fmt.Println("Result: ", val)
}
