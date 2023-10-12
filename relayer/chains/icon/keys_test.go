package icon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testKeyAddr     = "../../../example/icon/keystore.json"
	testKeyPassword = "x"
	expectedAddr    = "hxf36d99db01ef599d8117cbdd4036c4a598fdb2f9"
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
