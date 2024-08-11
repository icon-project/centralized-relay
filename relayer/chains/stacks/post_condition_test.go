package stacks

import (
	"encoding/hex"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyPostConditionsSerialization(t *testing.T) {
	emptyPostConditions := []PostCondition{}
	serialized := SerializePostConditions(emptyPostConditions)

	// Check that the serialized data is 4 bytes long (uint32 for length)
	assert.Equal(t, 4, len(serialized), "Serialized empty post conditions should be 4 bytes long")

	// Check that the serialized length is 0
	length := binary.BigEndian.Uint32(serialized)
	assert.Equal(t, uint32(0), length, "Serialized length of empty post conditions should be 0")
}

func TestEmptyPostConditionsDeserialization(t *testing.T) {
	serialized := make([]byte, 4)
	binary.BigEndian.PutUint32(serialized, 0)

	deserialized, bytesRead, err := DeserializePostConditions(serialized)

	assert.NoError(t, err, "Deserialization of empty post conditions should not produce an error")

	assert.Equal(t, 4, bytesRead, "Deserialization should have read 4 bytes")

	assert.Empty(t, deserialized, "Deserialized post conditions should be empty")
}

func TestEmptyPostConditionsRoundTrip(t *testing.T) {
	originalPostConditions := []PostCondition{}

	serialized := SerializePostConditions(originalPostConditions)

	deserialized, bytesRead, err := DeserializePostConditions(serialized)

	assert.NoError(t, err, "Round trip should not produce an error")
	assert.Equal(t, 4, bytesRead, "Round trip should read 4 bytes")
	assert.Empty(t, deserialized, "Round trip should result in empty post conditions")
	assert.Equal(t, len(originalPostConditions), len(deserialized), "Original and deserialized lengths should match")
}

// func TestSTXPostConditionSerializationAndDeserialization(t *testing.T) {
// 	address := "SP2JXKMSH007NPYAQHKJPQMAQYAD90NQGTVJVQ02B"
// 	principal := createStandardPrincipal(address)
// 	conditionCode := FungibleConditionCodeGreaterEqual
// 	amount := uint64(1000000)

// 	postCondition := createSTXPostCondition(principal, conditionCode, amount)

// 	serialized, err := postCondition.Serialize()
// 	assert.NoError(t, err)

// 	deserialized := &STXPostCondition{}
// 	_, err = deserialized.Deserialize(serialized)
// 	assert.NoError(t, err)

// 	assert.Equal(t, PostConditionTypeSTX, deserialized.GetType())
// 	assert.Equal(t, PostConditionPrincipalIDStandard, deserialized.PrincipalID)
// 	assert.Equal(t, address, principalToString(deserialized.Principal))
// 	assert.Equal(t, conditionCode, deserialized.ConditionCode)
// 	assert.Equal(t, amount, deserialized.Amount)
// }

// func TestFungiblePostConditionSerializationAndDeserialization(t *testing.T) {
// 	address := "SP2JXKMSH007NPYAQHKJPQMAQYAD90NQGTVJVQ02B"
// 	principal := createStandardPrincipal(address)
// 	conditionCode := FungibleConditionCodeGreaterEqual
// 	amount := uint64(1000000)

// 	assetAddress := "SP2ZP4GJDZJ1FDHTQ963F0292PE9J9752TZJ68F21"
// 	assetContractName := "contract_name"
// 	assetName := "asset_name"
// 	info := createAssetInfo(assetAddress, assetContractName, assetName)

// 	postCondition := createFungiblePostCondition(principal, conditionCode, amount, info)

// 	serialized, err := postCondition.Serialize()
// 	assert.NoError(t, err)

// 	deserialized := &FungiblePostCondition{}
// 	_, err = deserialized.Deserialize(serialized)
// 	assert.NoError(t, err)

// 	assert.Equal(t, PostConditionTypeFungible, deserialized.GetType())
// 	assert.Equal(t, PostConditionPrincipalIDStandard, deserialized.PrincipalID)
// 	assert.Equal(t, address, principalToString(deserialized.Principal))
// 	assert.Equal(t, conditionCode, deserialized.ConditionCode)
// 	assert.Equal(t, amount, deserialized.Amount)
// 	assert.Equal(t, assetAddress, addressToString(deserialized.AssetInfo.Address))
// 	assert.Equal(t, assetContractName, deserialized.AssetInfo.ContractName)
// 	assert.Equal(t, assetName, deserialized.AssetInfo.AssetName)
// }

// func TestNonFungiblePostConditionSerializationAndDeserialization(t *testing.T) {
// 	address := "SP2JXKMSH007NPYAQHKJPQMAQYAD90NQGTVJVQ02B"
// 	contractName := "contract-name"
// 	principal := createContractPrincipal(address, contractName)
// 	conditionCode := NonFungibleConditionCodeDoesNotSend

// 	assetAddress := "SP2ZP4GJDZJ1FDHTQ963F0292PE9J9752TZJ68F21"
// 	assetContractName := "contract_name"
// 	assetName := "asset_name"
// 	info := createAssetInfo(assetAddress, assetContractName, assetName)

// 	nftAssetName := "nft_asset_name"
// 	assetValue := []byte(nftAssetName)

// 	postCondition := createNonFungiblePostCondition(principal, conditionCode, info, assetValue)

// 	serialized, err := postCondition.Serialize()
// 	assert.NoError(t, err)

// 	deserialized := &NonFungiblePostCondition{}
// 	_, err = deserialized.Deserialize(serialized)
// 	assert.NoError(t, err)

// 	assert.Equal(t, PostConditionTypeNonFungible, deserialized.GetType())
// 	assert.Equal(t, PostConditionPrincipalIDContract, deserialized.PrincipalID)
// 	assert.Equal(t, address, principalToString(deserialized.Principal[:21])) // First 21 bytes are the address
// 	assert.Equal(t, contractName, string(deserialized.Principal[22:]))       // Rest is the contract name
// 	assert.Equal(t, conditionCode, deserialized.ConditionCode)
// 	assert.Equal(t, assetAddress, addressToString(deserialized.AssetInfo.Address))
// 	assert.Equal(t, assetContractName, deserialized.AssetInfo.ContractName)
// 	assert.Equal(t, assetName, deserialized.AssetInfo.AssetName)
// 	assert.Equal(t, nftAssetName, string(deserialized.AssetName))
// }

// func TestNonFungiblePostConditionWithStringIDsSerializationAndDeserialization(t *testing.T) {
// 	address := "SP2JXKMSH007NPYAQHKJPQMAQYAD90NQGTVJVQ02B"
// 	contractName := "contract-name"
// 	conditionCode := NonFungibleConditionCodeDoesNotSend

// 	assetAddress := "SP2ZP4GJDZJ1FDHTQ963F0292PE9J9752TZJ68F21"
// 	assetContractName := "contract_name"
// 	assetName := "asset_name"

// 	nftAssetName := "nft_asset_name"
// 	assetValue := []byte(nftAssetName)

// 	principal := createContractPrincipalString(address + "." + contractName)
// 	assetIdentifier := assetAddress + "." + assetContractName + "::" + assetName
// 	postCondition := createNonFungiblePostConditionFromStrings(principal, conditionCode, assetIdentifier, assetValue)

// 	serialized, err := postCondition.Serialize()
// 	assert.NoError(t, err)

// 	deserialized := &NonFungiblePostCondition{}
// 	_, err = deserialized.Deserialize(serialized)
// 	assert.NoError(t, err)

// 	assert.Equal(t, PostConditionTypeNonFungible, deserialized.GetType())
// 	assert.Equal(t, PostConditionPrincipalIDContract, deserialized.PrincipalID)
// 	assert.Equal(t, address, principalToString(deserialized.Principal[:21])) // First 21 bytes are the address
// 	assert.Equal(t, contractName, string(deserialized.Principal[22:]))       // Rest is the contract name
// 	assert.Equal(t, conditionCode, deserialized.ConditionCode)
// 	assert.Equal(t, assetAddress, addressToString(deserialized.AssetInfo.Address))
// 	assert.Equal(t, assetContractName, deserialized.AssetInfo.ContractName)
// 	assert.Equal(t, assetName, deserialized.AssetInfo.AssetName)
// 	assert.Equal(t, nftAssetName, string(deserialized.AssetName))
// }

// Helper functions

func createStandardPrincipal(address string) []byte {
	decoded, _ := hex.DecodeString(address[2:]) // Remove "SP" prefix and decode
	return append([]byte{0x02}, decoded...)     // 0x02 for standard principal
}

func createContractPrincipal(address string, contractName string) []byte {
	decoded, _ := hex.DecodeString(address[2:]) // Remove "SP" prefix and decode
	principal := append([]byte{0x03}, decoded...)
	return append(principal, byte(len(contractName)))
}

func createContractPrincipalString(fullAddress string) string {
	return fullAddress
}

func createAssetInfo(address, contractName, assetName string) AssetInfo {
	decoded, _ := hex.DecodeString(address[2:]) // Remove "SP" prefix and decode
	var addr [20]byte
	copy(addr[:], decoded)
	return AssetInfo{
		Address:      addr,
		ContractName: contractName,
		AssetName:    assetName,
	}
}

func createSTXPostCondition(principal []byte, conditionCode FungibleConditionCode, amount uint64) *STXPostCondition {
	return &STXPostCondition{
		PrincipalID:   PostConditionPrincipalIDStandard,
		Principal:     principal,
		ConditionCode: conditionCode,
		Amount:        amount,
	}
}

func createFungiblePostCondition(principal []byte, conditionCode FungibleConditionCode, amount uint64, assetInfo AssetInfo) *FungiblePostCondition {
	return &FungiblePostCondition{
		PrincipalID:   PostConditionPrincipalIDStandard,
		Principal:     principal,
		AssetInfo:     assetInfo,
		ConditionCode: conditionCode,
		Amount:        amount,
	}
}

// func createNonFungiblePostCondition(principal []byte, conditionCode NonFungibleConditionCode, assetInfo AssetInfo, assetValue []byte) *NonFungiblePostCondition {
// 	return &NonFungiblePostCondition{
// 		PrincipalID:   PostConditionPrincipalIDContract,
// 		Principal:     principal,
// 		AssetInfo:     assetInfo,
// 		AssetName:     assetValue,
// 		ConditionCode: conditionCode,
// 	}
// }

// func createNonFungiblePostConditionFromStrings(principal string, conditionCode NonFungibleConditionCode, assetIdentifier string, assetValue []byte) *NonFungiblePostCondition {
// 	// Parse principal
// 	parts := bytes.Split([]byte(principal), []byte("."))
// 	address := string(parts[0])
// 	contractName := string(parts[1])
// 	principalBytes := createContractPrincipal(address, contractName)

// 	// Parse asset identifier
// 	assetParts := bytes.Split([]byte(assetIdentifier), []byte("::"))
// 	assetAddressParts := bytes.Split(assetParts[0], []byte("."))
// 	assetAddress := string(assetAddressParts[0])
// 	assetContractName := string(assetAddressParts[1])
// 	assetName := string(assetParts[1])
// 	assetInfo := createAssetInfo(assetAddress, assetContractName, assetName)

// 	return createNonFungiblePostCondition(principalBytes, conditionCode, assetInfo, assetValue)
// }

func principalToString(principal []byte) string {
	return "SP" + hex.EncodeToString(principal[1:]) // Skip the first byte (principal type) and add "SP" prefix
}

func addressToString(address [20]byte) string {
	return "SP" + hex.EncodeToString(address[:])
}