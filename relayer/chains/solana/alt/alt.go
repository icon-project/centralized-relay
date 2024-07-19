package alt

import (
	"encoding/binary"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/near/borsh-go"
)

const (
	InstructionCreateLookupTable uint8 = iota
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

	slotBorsh, err := borsh.Serialize(recentSlot)
	if err != nil {
		return nil, solana.PublicKey{}, err
	}

	bumpBorsh, err := borsh.Serialize(bump)
	if err != nil {
		return nil, solana.PublicKey{}, err
	}

	instructionData := []byte{InstructionCreateLookupTable}
	instructionData = append(instructionData, slotBorsh...)
	instructionData = append(instructionData, bumpBorsh...)

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
	instructionData := []byte{InstructionExtendLookupTable}
	addressesBytes, _ := borsh.Serialize(addresses)
	instructionData = append(instructionData, addressesBytes...)

	keys := solana.AccountMetaSlice{
		{PublicKey: tableAddr, IsWritable: true},
		{PublicKey: authorityAddr, IsSigner: true},
	}

	if payerAddr != nil {
		keys = append(keys, solana.AccountMetaSlice{
			{PublicKey: *payerAddr, IsSigner: true},
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
	return &solana.GenericInstruction{
		ProgID: lookupTableProgramID,
		AccountValues: solana.AccountMetaSlice{
			&solana.AccountMeta{PublicKey: tableAddr, IsWritable: true},
			&solana.AccountMeta{PublicKey: authorityAddr, IsSigner: true},
		},
		DataBytes: []byte{InstructionFreezeLookupTable},
	}
}

// Constructs an instruction that deactivates an address lookup
// table so that it cannot be extended again and will be unusable
// and eligible for closure after a short amount of time.
func DeactivateLookupTable(
	tableAddr, authorityAddr solana.PublicKey,
) solana.Instruction {
	return &solana.GenericInstruction{
		ProgID: lookupTableProgramID,
		AccountValues: solana.AccountMetaSlice{
			&solana.AccountMeta{PublicKey: tableAddr, IsWritable: true},
			&solana.AccountMeta{PublicKey: authorityAddr, IsSigner: true},
		},
		DataBytes: []byte{InstructionDeactivateLookupTable},
	}
}

// Returns an instruction that closes an address lookup table
// account. The account will be deallocated and the lamports
// will be drained to the recipient address.
func CloseLookupTable(
	tableAddr, authorityAddr, recipientAddr solana.PublicKey,
) solana.Instruction {
	return &solana.GenericInstruction{
		ProgID: lookupTableProgramID,
		AccountValues: solana.AccountMetaSlice{
			&solana.AccountMeta{PublicKey: tableAddr, IsWritable: true},
			&solana.AccountMeta{PublicKey: authorityAddr, IsSigner: true},
			&solana.AccountMeta{PublicKey: recipientAddr},
		},
		DataBytes: []byte{InstructionCloseLookupTable},
	}
}
