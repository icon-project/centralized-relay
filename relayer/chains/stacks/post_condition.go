package stacks

import (
	"encoding/binary"
	"errors"
)

type PostCondition interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) (int, error)
	GetType() PostConditionType
}

// Assume empty post condition
func SerializePostConditions(postConditions []PostCondition) []byte {
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(postConditions)))
	return lenBytes
}

// Assume empty post condition
func DeserializePostConditions(data []byte) ([]PostCondition, int, error) {
	if len(data) < 4 {
		return nil, 0, errors.New("insufficient data for post conditions length")
	}
	length := binary.BigEndian.Uint32(data[:4])
	return make([]PostCondition, length), 4, nil
}

// func (t *BaseTransaction) SerializePostConditions() ([]byte, error) {
// 	buf := make([]byte, 0, 128)
// 	countBytes := make([]byte, 4)
// 	binary.BigEndian.PutUint32(countBytes, uint32(len(t.PostConditions)))
// 	buf = append(buf, countBytes...)
// 	for _, pc := range t.PostConditions {
// 		pcBytes, err := pc.Serialize()
// 		if err != nil {
// 			return nil, err
// 		}
// 		buf = append(buf, pcBytes...)
// 	}
// 	return buf, nil
// }

// func (t *BaseTransaction) DeserializePostConditions(data []byte) (int, error) {
// 	if len(data) < 4 {
// 		return 0, errors.New("not enough data to deserialize post conditions")
// 	}
// 	count := binary.BigEndian.Uint32(data[:4])
// 	offset := 4
// 	t.PostConditions = make([]PostCondition, 0, count)
// 	for i := uint32(0); i < count; i++ {
// 		if len(data[offset:]) < 1 {
// 			return 0, errors.New("not enough data to deserialize post condition")
// 		}
// 		var pc PostCondition
// 		switch PostConditionType(data[offset]) {
// 		case PostConditionTypeSTX:
// 			pc = &STXPostCondition{}
// 		case PostConditionTypeFungible:
// 			pc = &FungiblePostCondition{}
// 		// case PostConditionTypeNonFungible:
// 		// 	pc = &NonFungiblePostCondition{}
// 		default:
// 			return 0, errors.New("unknown post condition type")
// 		}
// 		size, err := pc.Deserialize(data[offset:])
// 		if err != nil {
// 			return 0, err
// 		}
// 		offset += size
// 		t.PostConditions = append(t.PostConditions, pc)
// 	}
// 	return offset, nil
// }

type STXPostCondition struct {
	PrincipalID PostConditionPrincipalID
	Principal   []byte
	ConditionCode FungibleConditionCode
	Amount      uint64
}

func (pc *STXPostCondition) GetType() PostConditionType {
	return PostConditionTypeSTX
}

func (pc *STXPostCondition) Serialize() ([]byte, error) {
	buf := make([]byte, 0, 64)
	buf = append(buf, byte(PostConditionTypeSTX))
	buf = append(buf, byte(pc.PrincipalID))
	buf = append(buf, pc.Principal...)
	buf = append(buf, byte(pc.ConditionCode))
	amountBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(amountBytes, pc.Amount)
	buf = append(buf, amountBytes...)
	return buf, nil
}

func (pc *STXPostCondition) Deserialize(data []byte) (int, error) {
	if len(data) < 11 { // 1 + 1 + 1 + 8 (minimum size for principal)
		return 0, errors.New("not enough data to deserialize STXPostCondition")
	}
	offset := 1 // Skip type byte
	pc.PrincipalID = PostConditionPrincipalID(data[offset])
	offset++

	principalSize, err := getPrincipalSize(pc.PrincipalID, data[offset:])
	if err != nil {
		return 0, err
	}
	if len(data[offset:]) < principalSize + 9 { // + 9 for condition code and amount
		return 0, errors.New("not enough data for principal, condition code, and amount")
	}
	pc.Principal = make([]byte, principalSize)
	copy(pc.Principal, data[offset:offset+principalSize])
	offset += principalSize

	pc.ConditionCode = FungibleConditionCode(data[offset])
	offset++

	pc.Amount = binary.BigEndian.Uint64(data[offset : offset+8])
	offset += 8

	return offset, nil
}

type FungiblePostCondition struct {
	PrincipalID PostConditionPrincipalID
	Principal   []byte
	AssetInfo   AssetInfo
	ConditionCode FungibleConditionCode
	Amount      uint64
}

func (pc *FungiblePostCondition) GetType() PostConditionType {
	return PostConditionTypeFungible
}

func (pc *FungiblePostCondition) Serialize() ([]byte, error) {
	buf := make([]byte, 0, 128)
	buf = append(buf, byte(PostConditionTypeFungible))
	buf = append(buf, byte(pc.PrincipalID))
	buf = append(buf, pc.Principal...)
	assetInfoBytes, err := pc.AssetInfo.Serialize()
	if err != nil {
		return nil, err
	}
	buf = append(buf, assetInfoBytes...)
	buf = append(buf, byte(pc.ConditionCode))
	amountBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(amountBytes, pc.Amount)
	buf = append(buf, amountBytes...)
	return buf, nil
}

func (pc *FungiblePostCondition) Deserialize(data []byte) (int, error) {
	if len(data) < 12 { // 1 + 1 + 1 + 1 + 8 (minimum size)
		return 0, errors.New("not enough data to deserialize FungiblePostCondition")
	}
	offset := 1 // Skip type byte
	pc.PrincipalID = PostConditionPrincipalID(data[offset])
	offset++
	principalSize, err := getPrincipalSize(pc.PrincipalID, data[offset:])
	pc.Principal = data[offset : offset+principalSize]
	offset += principalSize
	assetInfoSize, err := pc.AssetInfo.Deserialize(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += assetInfoSize
	pc.ConditionCode = FungibleConditionCode(data[offset])
	offset++
	pc.Amount = binary.BigEndian.Uint64(data[offset : offset+8])
	offset += 8
	return offset, nil
}

// type NonFungiblePostCondition struct {
// 	PrincipalID PostConditionPrincipalID
// 	Principal   []byte
// 	AssetInfo   AssetInfo
// 	AssetName   clarity.ClarityValue
// 	ConditionCode NonFungibleConditionCode
// }

// func (pc *NonFungiblePostCondition) GetType() PostConditionType {
// 	return PostConditionTypeNonFungible
// }

// func (pc *NonFungiblePostCondition) Serialize() ([]byte, error) {
// 	buf := make([]byte, 0, 128)
// 	buf = append(buf, byte(PostConditionTypeNonFungible))
// 	buf = append(buf, byte(pc.PrincipalID))
// 	buf = append(buf, pc.Principal...)
// 	assetInfoBytes, err := pc.AssetInfo.Serialize()
// 	if err != nil {
// 		return nil, err
// 	}
// 	buf = append(buf, assetInfoBytes...)
// 	assetNameBytes, err := pc.AssetName.Serialize()
// 	if err != nil {
// 		return nil, err
// 	}
// 	buf = append(buf, assetNameBytes...)
// 	buf = append(buf, byte(pc.ConditionCode))
// 	return buf, nil
// }

// func (pc *NonFungiblePostCondition) Deserialize(data []byte) (int, error) {
// 	if len(data) < 4 { // 1 + 1 + 1 + 1 (minimum size)
// 		return 0, errors.New("not enough data to deserialize NonFungiblePostCondition")
// 	}
// 	offset := 1 // Skip type byte
// 	pc.PrincipalID = PostConditionPrincipalID(data[offset])
// 	offset++
// 	principalSize := getPrincipalSize(pc.PrincipalID, data[offset:])
// 	pc.Principal = data[offset : offset+principalSize]
// 	offset += principalSize
// 	assetInfoSize, err := pc.AssetInfo.Deserialize(data[offset:])
// 	if err != nil {
// 		return 0, err
// 	}
// 	offset += assetInfoSize
// 	assetNameValue, assetNameSize, err := clarity.DeserializeClarityValue(data[offset:])
// 	if err != nil {
// 		return 0, err
// 	}
// 	pc.AssetName = assetNameValue
// 	offset += assetNameSize
// 	pc.ConditionCode = NonFungibleConditionCode(data[offset])
// 	offset++
// 	return offset, nil
// }

type AssetInfo struct {
	Address     [20]byte
	ContractName string
	AssetName    string
}

func (ai *AssetInfo) Serialize() ([]byte, error) {
	buf := make([]byte, 0, 64)
	buf = append(buf, ai.Address[:]...)
	buf = append(buf, byte(len(ai.ContractName)))
	buf = append(buf, []byte(ai.ContractName)...)
	buf = append(buf, byte(len(ai.AssetName)))
	buf = append(buf, []byte(ai.AssetName)...)
	return buf, nil
}

func (ai *AssetInfo) Deserialize(data []byte) (int, error) {
	if len(data) < 22 { // 20 + 1 + 1 (minimum size)
		return 0, errors.New("not enough data to deserialize AssetInfo")
	}
	offset := 0
	copy(ai.Address[:], data[:20])
	offset += 20
	contractNameLen := int(data[offset])
	offset++
	if len(data[offset:]) < contractNameLen+1 {
		return 0, errors.New("not enough data for contract name")
	}
	ai.ContractName = string(data[offset : offset+contractNameLen])
	offset += contractNameLen
	assetNameLen := int(data[offset])
	offset++
	if len(data[offset:]) < assetNameLen {
		return 0, errors.New("not enough data for asset name")
	}
	ai.AssetName = string(data[offset : offset+assetNameLen])
	offset += assetNameLen
	return offset, nil
}

// Helper functions

func getPrincipalSize(principalID PostConditionPrincipalID, data []byte) (int, error) {
	switch principalID {
	case PostConditionPrincipalIDOrigin:
		return 0, nil
	case PostConditionPrincipalIDStandard:
		return 21, nil // 1 byte version + 20 bytes hash
	case PostConditionPrincipalIDContract:
		if len(data) < 22 {
			return 0, errors.New("not enough data for contract principal")
		}
		return 21 + 1 + int(data[21]), nil // 21 bytes standard principal + 1 byte name length + name
	default:
		return 0, errors.New("invalid principal ID")
	}
}