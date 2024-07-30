package stacks

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSTXTokenTransferTransactionSerializationAndDeserialization(t *testing.T) {
	transactionVersion := TransactionVersionTestnet
	chainID := DefaultChainID
	anchorMode := AnchorModeAny
	postConditionMode := PostConditionModeDeny

	address := "SP3FGQ8Z7JY9BWYZ5WM53E0M9NK7WHJF0691NZ159"
	recipientAddress, err := NewAddress(address)
	assert.NoError(t, err)

	amount := uint64(2500000)
	memo := "memo (not included)"

	payload := NewTokenTransferPayload(recipientAddress, amount, memo)

	addressHashMode := AddressHashModeSerializeP2PKH
	nonce := uint64(0)
	fee := uint64(0)

	pubKey := "03ef788b3830c00abe8f64f62dc32fc863bc0b2cafeb073b6c8e1c7657d9c2c3ab"
	pubKeyBytes, err := hex.DecodeString(pubKey)
	assert.NoError(t, err)

	secretKey := "edf9aee84d9b7abc145504dde6726c64f369d37ee34ded868fabd876c26570bc01"
	secretKeyBytes, err := hex.DecodeString(secretKey)
	assert.NoError(t, err)

	spendingCondition := NewSingleSigSpendingCondition(addressHashMode, pubKeyBytes, nonce, fee)
	authorization := NewStandardAuthorization(spendingCondition)

	transaction := NewTransaction(
		transactionVersion,
		chainID,
		authorization,
		anchorMode,
		postConditionMode,
		[]PostCondition{},
		payload,
	)

	// Sign the transaction
	err = transaction.Sign(secretKeyBytes)
	assert.NoError(t, err)

	// Verify the transaction
	err = transaction.Verify()
	assert.NoError(t, err)

	serialized, err := transaction.Serialize()
	assert.NoError(t, err)

	deserialized, err := DeserializeTransaction(serialized)
	assert.NoError(t, err)

	// Verify deserialized transaction
	assert.Equal(t, transactionVersion, deserialized.Version)
	assert.Equal(t, chainID, deserialized.ChainID)
	assert.Equal(t, AuthTypeStandard, deserialized.Authorization.AuthType)
	assert.Equal(t, addressHashMode, deserialized.Authorization.SpendingCondition.HashMode)
	assert.Equal(t, nonce, deserialized.Authorization.SpendingCondition.Nonce)
	assert.Equal(t, fee, deserialized.Authorization.SpendingCondition.Fee)
	assert.Equal(t, anchorMode, deserialized.AnchorMode)
	assert.Equal(t, postConditionMode, deserialized.PostConditionMode)
	assert.Empty(t, deserialized.PostConditions)

	deserializedPayload, ok := deserialized.Payload.(*TokenTransferPayload)
	assert.True(t, ok)
	assert.Equal(t, recipientAddress, deserializedPayload.Recipient)
	assert.Equal(t, amount, deserializedPayload.Amount)

	// Verify the deserialized transaction's signature
	err = deserialized.Verify()
	assert.NoError(t, err)

	// Test serialization of the deserialized transaction
	reserializedBytes, err := deserialized.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, serialized, reserializedBytes)
}