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

	fmt.Println("65 Bytes Signature: ", hex.EncodeToString(signature))
	fmt.Println("Public key: ", kp.PublicKey().String())
	fmt.Println("Recovery param: ", signature[0])

	//65 Bytes Signature:  1b30b073c7fdb2b2fac752f40c13a5dbadaa904fe4b7c1b19fe1a31778cdacd88f7388d5d990f99428d185a2de65428979e6890ba66f7dec2f64e4b3d9d3433d5b
	//Public key:  02e27e3817bf0b6d451004609c2a5d29fe315dc1d1017500399fab540785958b7a
	//Recovery param:  27 => 0, 28 => 1

	_, err = kp.VerifySignature(msgHash, signature)
	assert.NoError(t, err)
}
