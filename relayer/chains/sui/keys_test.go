package sui

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	expectedAddr    = "0xe847098636459aa93f4da105414edca4790619b291ffdac49419f5adc19c4d21"
	expectedPrivKey = "b592e26293b6081673c807f9ae5b14b150d0078d6d9a5474323fff73a9015cac"
)

func TestRestoreKey(t *testing.T) {
	encodedWithFlag := "ALWS4mKTtggWc8gH+a5bFLFQ0AeNbZpUdDI//3OpAVys"
	key, err := fetchKeyPair(encodedWithFlag)
	assert.NoError(t, err)
	assert.Equal(t, expectedAddr, key.Address)
	assert.Equal(t, expectedPrivKey, hex.EncodeToString(key.KeyPair.PrivateKey()[:32]))

}
