package evm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testKeyStore    = "../../../example/wallets/evm/keystore.json"
	testKeyPassword = "secret"
	expectedAddr    = "0x33768aEdeAF1D4d3634C93d551dEbB69Eb4104a5"
)

func TestRestoreKey(t *testing.T) {
	key, err := RestoreKey(testKeyStore, testKeyPassword)
	assert.NoError(t, err)
	assert.Equal(t, key.Address.String(), expectedAddr)
}
