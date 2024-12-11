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
	dstNetwork := "archway"
	sn := new(big.Int).SetUint64(128)
	data := []byte("hello")

	msg := types.Message{
		Src:  srcNetwork,
		Dst:  dstNetwork,
		Sn:   sn,
		Data: data,
	}

	msgBytes := msg.SignableMsgV1()

	fmt.Println("String encoded bytes: ", hex.EncodeToString(msgBytes))

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

	_, err = kp.VerifySignature(msgHash, signature)
	assert.NoError(t, err)
}

func TestIntToString(t *testing.T) {
	val := big.NewInt(128)
	fmt.Println("String val: ", val.String())
}

func TestSigVerificationAgain(t *testing.T) {
	signatureHex := "04fdb08012c0909afe508595d03718822ae691299459625b3cebd0237cd3b8643e6564a4806ffbcd344b09b61dea4a920f83edeb1ddde293344682684e4eaa3d01"
	signatureBytes, err := hex.DecodeString(signatureHex)
	assert.NoError(t, err)

	fmt.Println("Signature lenght: ", len(signatureBytes))

	recoveryKey := signatureBytes[64]
	vrsSig := make([]byte, 65)
	for i := 64; i > 0; i-- {
		vrsSig[i] = signatureBytes[i-1]
	}

	if recoveryKey == 0 {
		vrsSig[0] = 27
	} else {
		vrsSig[0] = 28
	}

	srcNetwork := "0x1.icon"
	dstNetwork := "archway-1"
	sn := new(big.Int).SetUint64(24)
	data, err := hex.DecodeString("f8c801b8c5f8c3b33078312e69636f6e2f637863356434306664373439393562656434373365356431623235396462633630313532373366666335b84261726368776179316c647a6a647134307a796634373773616b73797078303238687a7468726b7964366b6766727171333971356d3679377278656773647a3075306e8234030080f844b84261726368776179316c766d783275366634376e3879723064673766616e677572326c37326e7778786b6c61737179616c3266687470797739757866716d7564656c38")
	assert.NoError(t, err)

	msg := types.Message{
		Src:  srcNetwork,
		Dst:  dstNetwork,
		Sn:   sn,
		Data: data,
	}

	msgBytes := msg.SignableMsgV1()
	msgHash := Keccak256Hash(msgBytes)

	fmt.Println("msg: ", hex.EncodeToString(msgBytes))
	fmt.Println("hash: ", hex.EncodeToString(msgHash))

	err = VerifySignaturee(msgHash, vrsSig)
	assert.NoError(t, err)
}
