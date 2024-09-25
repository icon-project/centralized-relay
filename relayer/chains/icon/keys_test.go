package icon

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/goloop/common/wallet"
	"github.com/stretchr/testify/assert"
)

func TestWallet(t *testing.T) {
	priv, _ := crypto.GenerateKeyPair()
	wallet, _ := wallet.NewFromPrivateKey(priv)
	fmt.Println("Wallet Address: ", wallet.Address())
	fmt.Println("Private Key: ", priv.String())

	data := []byte("hello")
	dataHash := sha256.Sum256(data)
	signature, err := wallet.Sign(dataHash[:])
	assert.NoError(t, err)

	sig, err := crypto.NewSignature(dataHash[:], priv)
	assert.NoError(t, err)

	signature1, err := sig.SerializeRSV()
	assert.NoError(t, err)

	assert.Equal(t, signature, signature1)

	fmt.Println("Signaure length: ", len(signature))
	fmt.Println("Signaure1 length: ", len(signature1))

}
