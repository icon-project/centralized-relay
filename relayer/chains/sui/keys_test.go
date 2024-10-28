package sui

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/coming-chat/go-sui/v2/sui_types"
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

	sign, _ := key.SignSecureWithoutEncode([]byte("helloworld"), sui_types.DefaultIntent())
	signs := sign.Ed25519SuiSignature.Signature
	// ms, _ := sign.MarshalJSON()
	fmt.Println(hex.EncodeToString(key.KeyPair.PublicKey()))
	fmt.Println(key.Address)

	fmt.Println("message", hex.EncodeToString(signs[:]))
}
