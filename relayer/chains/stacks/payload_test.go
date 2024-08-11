// relayer/chains/stacks/payload_test.go
package stacks

import (
	"testing"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/clarity"
	"github.com/stretchr/testify/assert"
)

func TestTokenTransferPayloadSerializationDeserialization(t *testing.T) {
	recipient := "SP3FGQ8Z7JY9BWYZ5WM53E0M9NK7WHJF0691NZ159"
	amount := uint64(2500000)
	memo := "test memo"

	payload, err := NewTokenTransferPayload(recipient, amount, memo)
	assert.NoError(t, err)

	serialized, err := payload.Serialize()
	assert.NoError(t, err)

	deserialized := &TokenTransferPayload{}
	_, err = deserialized.Deserialize(serialized)
	assert.NoError(t, err)

	assert.Equal(t, payload.Recipient, deserialized.Recipient)
	assert.Equal(t, payload.Amount, deserialized.Amount)
	assert.Equal(t, payload.Memo, deserialized.Memo)
}

func TestTokenTransferPayloadToContractAddress(t *testing.T) {
	recipient := "SP3FGQ8Z7JY9BWYZ5WM53E0M9NK7WHJF0691NZ159.contract-name"
	amount := uint64(2500000)
	memo := "test memo"

	payload, err := NewTokenTransferPayload(recipient, amount, memo)
	assert.NoError(t, err)

	serialized, err := payload.Serialize()
	assert.NoError(t, err)

	deserialized := &TokenTransferPayload{}
	_, err = deserialized.Deserialize(serialized)
	assert.NoError(t, err)

	assert.Equal(t, payload.Recipient, deserialized.Recipient)
	assert.Equal(t, payload.Amount, deserialized.Amount)
	assert.Equal(t, payload.Memo, deserialized.Memo)
}

func TestContractCallPayloadSerializationDeserialization(t *testing.T) {
	contractAddress := "SP3FGQ8Z7JY9BWYZ5WM53E0M9NK7WHJF0691NZ159"
	contractName := "contract_name"
	functionName := "function_name"
	args := []clarity.ClarityValue{
		clarity.NewBool(true),
		clarity.NewBool(false),
	}

	payload := &ContractCallPayload{
		ContractAddress: contractAddress,
		ContractName:    contractName,
		FunctionName:    functionName,
		FunctionArgs:    args,
	}

	serialized, err := payload.Serialize()
	assert.NoError(t, err)

	deserialized := &ContractCallPayload{}
	_, err = deserialized.Deserialize(serialized)
	assert.NoError(t, err)

	assert.Equal(t, payload.ContractAddress, deserialized.ContractAddress)
	assert.Equal(t, payload.ContractName, deserialized.ContractName)
	assert.Equal(t, payload.FunctionName, deserialized.FunctionName)
	assert.Equal(t, len(payload.FunctionArgs), len(deserialized.FunctionArgs))
	for i, arg := range payload.FunctionArgs {
		assert.Equal(t, arg.Type(), deserialized.FunctionArgs[i].Type())
	}
}
