package stacks

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/clarity"
	"github.com/stretchr/testify/assert"
)

func TestSTXTokenTransferTransactionSerializationAndDeserialization(t *testing.T) {
	transactionVersion := TransactionVersionTestnet
	chainID := ChainIDTestnet
	anchorMode := AnchorModeOnChainOnly
	postConditionMode := PostConditionModeDeny

	recipientAddress := "SP3FGQ8Z7JY9BWYZ5WM53E0M9NK7WHJF0691NZ159"
	amount := uint64(2500000)
	memo := "memo (not included)"

	payload, err := NewTokenTransferPayload(recipientAddress, amount, memo)
	assert.NoError(t, err)

	addressHashMode := AddressHashModeSerializeP2PKH
	nonce := uint64(0)
	fee := uint64(0)

	pubKey := "03ef788b3830c00abe8f64f62dc32fc863bc0b2cafeb073b6c8e1c7657d9c2c3ab"
	pubKeyBytes, err := hex.DecodeString(pubKey)
	assert.NoError(t, err)

	secretKey := "edf9aee84d9b7abc145504dde6726c64f369d37ee34ded868fabd876c26570bc01"
	secretKeyBytes, err := hex.DecodeString(secretKey)
	assert.NoError(t, err)

	spendingCondition := SpendingCondition{
		HashMode:    addressHashMode,
		Signer:      [20]byte{},
		Nonce:       nonce,
		Fee:         fee,
		KeyEncoding: PubKeyEncodingCompressed,
		Signature:   [65]byte{},
	}

	auth := TransactionAuth{
		AuthType:   AuthTypeStandard,
		OriginAuth: spendingCondition,
	}

	transaction := &TokenTransferTransaction{
		BaseTransaction: BaseTransaction{
			Version:           transactionVersion,
			ChainID:           chainID,
			Auth:              auth,
			AnchorMode:        anchorMode,
			PostConditionMode: postConditionMode,
			PostConditions:    []PostCondition{},
		},
		Payload: *payload,
	}

	err = SignTransaction(transaction, secretKeyBytes)
	assert.NoError(t, err)

	isValid, err := VerifyTransaction(transaction, pubKeyBytes)
	assert.NoError(t, err)

	assert.True(t, isValid)

	serialized, err := transaction.Serialize()
	assert.NoError(t, err)

	deserialized, err := DeserializeTransaction(serialized)
	assert.NoError(t, err)

	tokenTx, ok := deserialized.(*TokenTransferTransaction)
	assert.True(t, ok, "Deserialized transaction is not a TokenTransferTransaction")

	assert.Equal(t, transactionVersion, tokenTx.Version)
	assert.Equal(t, chainID, tokenTx.ChainID)
	assert.Equal(t, AuthTypeStandard, tokenTx.Auth.AuthType)
	assert.Equal(t, addressHashMode, tokenTx.Auth.OriginAuth.HashMode)
	assert.Equal(t, nonce, tokenTx.Auth.OriginAuth.Nonce)
	assert.Equal(t, fee, tokenTx.Auth.OriginAuth.Fee)
	assert.Equal(t, anchorMode, tokenTx.AnchorMode)
	assert.Equal(t, postConditionMode, tokenTx.PostConditionMode)
	assert.Empty(t, tokenTx.PostConditions)

	recipientPrincipal, _ := clarity.StringToPrincipal(recipientAddress)

	assert.Equal(t, recipientPrincipal, tokenTx.Payload.Recipient)
	assert.Equal(t, amount, tokenTx.Payload.Amount)
	assert.Equal(t, memo, tokenTx.Payload.Memo)

	isValid, err = VerifyTransaction(tokenTx, pubKeyBytes)
	assert.NoError(t, err)
	assert.True(t, isValid)

	reserializedBytes, err := deserialized.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, serialized, reserializedBytes)
}

func TestSingleSpendingConditionSerializationAndDeserialization(t *testing.T) {
	addressHashMode := AddressHashModeSerializeP2PKH
	nonce := uint64(0)
	fee := uint64(0)
	pubKey := "03ef788b3830c00abe8f64f62dc32fc863bc0b2cafeb073b6c8e1c7657d9c2c3ab"

	spendingCondition := createSingleSigSpendingCondition(addressHashMode, pubKey, nonce, fee)
	emptySignature := emptyMessageSignature()

	serialized, err := spendingCondition.SerializeSpendingCondition()
	assert.NoError(t, err, "Failed to serialize spending condition")

	deserialized := SpendingCondition{}
	_, err = deserialized.DeserializeSpendingCondition(serialized)
	assert.NoError(t, err, "Failed to deserialize spending condition")

	assert.Equal(t, addressHashMode, deserialized.HashMode, "HashMode mismatch")
	assert.Equal(t, nonce, deserialized.Nonce, "Nonce mismatch")
	assert.Equal(t, fee, deserialized.Fee, "Fee mismatch")
	assert.True(t, bytes.Equal(deserialized.Signature[:], emptySignature[:]), "Signature mismatch")
}

func TestSingleSigP2PKHSpendingCondition(t *testing.T) {
	spCompressed := createSingleSigSpendingCondition(AddressHashModeSerializeP2PKH, "", 345, 456)
	spCompressed.KeyEncoding = PubKeyEncodingCompressed
	spCompressed.Signature = createMessageSignature("fe")
	spCompressed.Signer = [20]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}

	serialized, err := spCompressed.SerializeSpendingCondition()
	assert.NoError(t, err, "Failed to serialize compressed P2PKH spending condition")

	deserialized := SpendingCondition{}
	_, err = deserialized.DeserializeSpendingCondition(serialized)
	assert.NoError(t, err, "Failed to deserialize compressed P2PKH spending condition")

	assert.Equal(t, spCompressed, deserialized, "Compressed P2PKH spending condition mismatch")

	spUncompressed := createSingleSigSpendingCondition(AddressHashModeSerializeP2PKH, "", 123, 456)
	spUncompressed.KeyEncoding = PubKeyEncodingUncompressed
	spUncompressed.Signature = createMessageSignature("ff")
	spUncompressed.Signer = [20]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}

	serialized, err = spUncompressed.SerializeSpendingCondition()
	assert.NoError(t, err, "Failed to serialize uncompressed P2PKH spending condition")

	deserialized = SpendingCondition{}
	_, err = deserialized.DeserializeSpendingCondition(serialized)
	assert.NoError(t, err, "Failed to deserialize uncompressed P2PKH spending condition")

	assert.Equal(t, spUncompressed, deserialized, "Uncompressed P2PKH spending condition mismatch")
}

func TestSingleSigP2WPKHSpendingCondition(t *testing.T) {
	sp := createSingleSigSpendingCondition(AddressHashModeSerializeP2WPKH, "", 345, 567)
	sp.KeyEncoding = PubKeyEncodingCompressed
	sp.Signature = createMessageSignature("fe")
	sp.Signer = [20]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}

	serialized, err := sp.SerializeSpendingCondition()
	assert.NoError(t, err, "Failed to serialize P2WPKH spending condition")

	deserialized := SpendingCondition{}
	_, err = deserialized.DeserializeSpendingCondition(serialized)
	assert.NoError(t, err, "Failed to deserialize P2WPKH spending condition")

	assert.Equal(t, sp, deserialized, "P2WPKH spending condition mismatch")
}

func TestInvalidSpendingConditions(t *testing.T) {
	// Test invalid hash mode
	invalidHashMode := []byte{
		0xff,                                                                                                                   // Invalid hash mode
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // Signer
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xc8, // Nonce
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x37, // Fee
		0x00, // Key encoding
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, // Signature
	}

	sp := SpendingCondition{}
	_, err := sp.DeserializeSpendingCondition(invalidHashMode)
	assert.Error(t, err, "Expected error for invalid hash mode")

	// Test incompatible hash mode and key encoding (P2WPKH with uncompressed key)
	sp = createSingleSigSpendingCondition(AddressHashModeSerializeP2WPKH, "", 123, 567)
	sp.KeyEncoding = PubKeyEncodingUncompressed
	sp.Signature = createMessageSignature("ff")
	sp.Signer = [20]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}

	serialized, _ := sp.SerializeSpendingCondition()
	_, err = sp.DeserializeSpendingCondition(serialized)
	assert.Error(t, err, "Expected error for incompatible hash mode and key encoding")
}

func createSingleSigSpendingCondition(hashMode AddressHashMode, pubKey string, nonce, fee uint64) SpendingCondition {
	return SpendingCondition{
		HashMode: hashMode,
		Signer:   [20]byte{},
		Nonce:    nonce,
		Fee:      fee,
	}
}

func createMessageSignature(hexString string) [65]byte {
	var signature [65]byte
	sigBytes, _ := hex.DecodeString(hexString)
	copy(signature[:], sigBytes)
	return signature
}

func emptyMessageSignature() [65]byte {
	return [65]byte{}
}
