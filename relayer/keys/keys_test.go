package keys

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/icon-project/centralized-relay/relayer/types"
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

func TestSigVerification(t *testing.T) {
	kp, err := NewKeyPair(Secp256k1)
	assert.NoError(t, err)

	srcNetwork := "0x2.icon"
	sn := new(big.Int).SetUint64(456456)
	data := []byte("hello")

	msg := types.Message{
		Src:  srcNetwork,
		Sn:   sn,
		Data: data,
	}

	msgBytes := msg.SignableMsg()

	msgHash := Keccak256Hash(msgBytes)
	signature := kp.Sign(msgHash)
	recoveryKey := signature[0]
	rsvSig := make([]byte, 65)
	for i := 0; i < 64; i++ {
		rsvSig[i] = signature[i+1]
	}
	rsvSig[64] = recoveryKey

	fmt.Println("65 Bytes Signature: ", hex.EncodeToString(rsvSig))
	fmt.Println("Public key: ", kp.PublicKey().String())
	fmt.Println("Recovery param: ", recoveryKey)

	_, err = kp.VerifySignature(msgHash, signature)
	assert.NoError(t, err)
}
