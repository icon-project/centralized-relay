package alt

import (
	"encoding/binary"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/icon-project/centralized-relay/test/chains/solana/bincode"
)

const (
	InstructionCreateLookupTable uint32 = iota
	InstructionFreezeLookupTable
	InstructionExtendLookupTable
	InstructionDeactivateLookupTable
	InstructionCloseLookupTable
)

var (
	lookupTableProgramID, _ = solana.PublicKeyFromBase58("AddressLookupTab1e1111111111111111111111111")
)

// Constructs an instruction to create a table account and returns
// the instruction and the table account's derived address.
func CreateLookupTable(
	authority, payer solana.PublicKey,
	recentSlot uint64,
) (solana.Instruction, solana.PublicKey, error) {
	slotBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(slotBytes, recentSlot)
	acSeeds := [][]byte{
		authority.Bytes(),
		slotBytes,
	}

	address, bump, err := solana.FindProgramAddress(acSeeds, lookupTableProgramID)
	if err != nil {
		return nil, solana.PublicKey{}, err
	}

	recentSlotBytes, err := bincode.Serialize(recentSlot)
	if err != nil {
		return nil, solana.PublicKey{}, err
	}

	bumpBytes, err := bincode.Serialize(bump)
	if err != nil {
		return nil, solana.PublicKey{}, err
	}

	discriminantBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(discriminantBytes, InstructionCreateLookupTable)
	instructionData := discriminantBytes
	instructionData = append(instructionData, recentSlotBytes...)
	instructionData = append(instructionData, bumpBytes...)

	keys := solana.AccountMetaSlice{
		{PublicKey: address, IsSigner: false, IsWritable: true},
		{PublicKey: authority, IsSigner: true, IsWritable: false},
		{PublicKey: payer, IsSigner: true, IsWritable: true},
		{PublicKey: system.ProgramID, IsSigner: false, IsWritable: false},
	}

	return &solana.GenericInstruction{
		ProgID:        lookupTableProgramID,
		AccountValues: keys,
		DataBytes:     instructionData,
	}, address, nil
}

// Constructs an instruction which extends an address lookup
// table account with new addresses.
func ExtendLookupTable(
	tableAddr, authorityAddr solana.PublicKey,
	payerAddr *solana.PublicKey,
	addresses solana.PublicKeySlice,
) solana.Instruction {
	discriminantBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(discriminantBytes, InstructionExtendLookupTable)
	instructionData := discriminantBytes
	addressesBytes, err := bincode.Serialize(addresses)
	if err != nil {
		panic(err)
	}
	instructionData = append(instructionData, addressesBytes...)

	keys := solana.AccountMetaSlice{
		{PublicKey: tableAddr, IsWritable: true, IsSigner: false},
		{PublicKey: authorityAddr, IsWritable: false, IsSigner: true},
	}

	if payerAddr != nil {
		keys = append(keys, solana.AccountMetaSlice{
			{PublicKey: *payerAddr, IsSigner: true, IsWritable: true},
			{PublicKey: system.ProgramID, IsSigner: false, IsWritable: false},
		}...)
	}

	return &solana.GenericInstruction{
		ProgID:        lookupTableProgramID,
		AccountValues: keys,
		DataBytes:     instructionData,
	}
}

// Constructs an instruction that freezes an address lookup
// table so that it can never be closed or extended again. Empty
// lookup tables cannot be frozen.
func FreezeLookupTable(
	tableAddr, authorityAddr solana.PublicKey,
) solana.Instruction {
	discriminantBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(discriminantBytes, InstructionFreezeLookupTable)
	return &solana.GenericInstruction{
		ProgID: lookupTableProgramID,
		AccountValues: solana.AccountMetaSlice{
			&solana.AccountMeta{PublicKey: tableAddr, IsWritable: true},
			&solana.AccountMeta{PublicKey: authorityAddr, IsSigner: true},
		},
		DataBytes: discriminantBytes,
	}
}

// Constructs an instruction that deactivates an address lookup
// table so that it cannot be extended again and will be unusable
// and eligible for closure after a short amount of time.
func DeactivateLookupTable(
	tableAddr, authorityAddr solana.PublicKey,
) solana.Instruction {
	discriminantBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(discriminantBytes, InstructionDeactivateLookupTable)
	return &solana.GenericInstruction{
		ProgID: lookupTableProgramID,
		AccountValues: solana.AccountMetaSlice{
			&solana.AccountMeta{PublicKey: tableAddr, IsWritable: true},
			&solana.AccountMeta{PublicKey: authorityAddr, IsSigner: true},
		},
		DataBytes: discriminantBytes,
	}
}

// Returns an instruction that closes an address lookup table
// account. The account will be deallocated and the lamports
// will be drained to the recipient address.
func CloseLookupTable(
	tableAddr, authorityAddr, recipientAddr solana.PublicKey,
) solana.Instruction {
	discriminantBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(discriminantBytes, InstructionCloseLookupTable)
	return &solana.GenericInstruction{
		ProgID: lookupTableProgramID,
		AccountValues: solana.AccountMetaSlice{
			&solana.AccountMeta{PublicKey: tableAddr, IsWritable: true},
			&solana.AccountMeta{PublicKey: authorityAddr, IsSigner: true},
			&solana.AccountMeta{PublicKey: recipientAddr},
		},
		DataBytes: discriminantBytes,
	}
}
