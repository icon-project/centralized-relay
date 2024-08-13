package keystore

import (
	"fmt"
	"testing"

	bftRand "github.com/cometbft/cometbft/libs/rand"

	"github.com/coming-chat/go-sui/v2/account"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/stretchr/testify/assert"
)

const (
	Ed25519Flag            byte = 0
	ed25519PublicKeyLength      = 32
)

func TestJSONkeystore(t *testing.T) {
	signatureScheme, err := sui_types.NewSignatureScheme(byte(Ed25519Flag))
	assert.NoError(t, err)

	account := account.NewAccount(signatureScheme, bftRand.Bytes(32))
	privateKeyWithFlag := append([]byte{byte(Ed25519Flag)}, account.KeyPair.PrivateKey()[:ed25519PublicKeyLength]...)

	password := "password"

	jsonBytes, err := EncryptToJSONKeystore(privateKeyWithFlag, account.Address, password)
	assert.NoError(t, err)

	fmt.Println("JSON Keystore: ", string(jsonBytes))

	decryptedPrivateKey, err := DecryptFromJSONKeystore(jsonBytes, password)
	assert.NoError(t, err)

	assert.Equal(t, privateKeyWithFlag, decryptedPrivateKey)

	wrongPassword := "wrongPassword"
	_, err = DecryptFromJSONKeystore(jsonBytes, wrongPassword)
	assert.Contains(t, err.Error(), "cipher: message authentication failed")
}
