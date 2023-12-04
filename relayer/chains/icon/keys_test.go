package icon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testKeyAddr     = "../../../example/icon/keystore.json"
	testKeyPassword = "gochain"
	expectedAddr    = "hxb6b5791be0b5ef67063b3c10b840fb81514db2fd"
)

func TestRestoreIconKey(t *testing.T) {

	iconProvider, err := GetMockIconProvider()
	assert.NoError(t, err)

	addr, err := iconProvider.GetWalletAddress()
	assert.NoError(t, err)
	assert.Equal(t, expectedAddr, addr)
}

func TestGetAddrFromKeystore(t *testing.T) {

	addr, err := getAddrFromKeystore(testKeyAddr)
	assert.NoError(t, err)
	assert.Equal(t, expectedAddr, addr)
}
