package keys

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecp256k1Key(t *testing.T) {
	kp, err := NewKeyPair(Secp256k1)
	assert.NoError(t, err)
	msgHash := Sha256Hash([]byte("hello"))
	signature := kp.Sign(msgHash)

	fmt.Println("lenght of signature: ", len(signature))
	fmt.Println("Signature: ", hex.EncodeToString(signature))

	_, err = kp.VerifySignature(Sha256Hash([]byte("hello")), signature)
	assert.NoError(t, err)

}
