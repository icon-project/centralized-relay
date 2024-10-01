package alt

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/gagliardetto/solana-go"
)

const (
	LOOKUP_TABLE_MAX_ADDRESSES = 256
	LOOKUP_TABLE_META_SIZE     = 56
)

type ProgramStateType uint32

const (
	ProgramStateUninitialized ProgramStateType = iota
	ProgramStateLookupTable
)

type LookupTableMeta struct {
	// Lookup tables cannot be closed until the deactivation slot is
	// no longer "recent" (not accessible in the `SlotHashes` sysvar).
	DeactivationSlot uint64
	//The slot that the table was last extended. Address tables may
	//only be used to lookup addresses that were extended before
	//the current bank's slot.
	LastExtendedSlot uint64
	//The start index where the table was last extended from during
	//the `last_extended_slot`.
	LastExtendedSlotStartIndex uint8
	//Authority address which must sign for each modification.
	Authority *solana.PublicKey
	// Padding to keep addresses 8-byte aligned
	padding uint16
	// Raw list of addresses follows this serialized structure in
	// the account's data, starting from `LOOKUP_TABLE_META_SIZE`.
}

func (ac *LookupTableAccount) IsActive() bool {
	return ac.Meta.DeactivationSlot == math.MaxUint64
}

type LookupTableAccount struct {
	ProgramState ProgramStateType
	Meta         LookupTableMeta
	Addresses    solana.PublicKeySlice
}

func DeserializeLookupTable(data []byte) (*LookupTableAccount, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("invalid account data")
	}

	current := 0
	programState := ProgramStateType(binary.LittleEndian.Uint32(data[current : current+4]))
	current += 4

	if programState == ProgramStateUninitialized {
		return &LookupTableAccount{ProgramState: programState}, nil
	}

	if len(data) < LOOKUP_TABLE_META_SIZE {
		return nil, fmt.Errorf("invalid account data")
	}

	ac := LookupTableAccount{
		ProgramState: programState,
	}

	ac.Meta.DeactivationSlot = binary.LittleEndian.Uint64(data[current : current+8])
	current += 8

	ac.Meta.LastExtendedSlot = binary.LittleEndian.Uint64(data[current : current+8])
	current += 8

	ac.Meta.LastExtendedSlotStartIndex = data[current]
	current += 1

	some := bool(data[current] == 1)
	current += 1
	if some {
		pubkey := solana.PublicKeyFromBytes(data[current : current+32])
		current += 32
		ac.Meta.Authority = &pubkey
	}

	ac.Meta.padding = binary.LittleEndian.Uint16(data[current : current+2])
	current += 2

	addressCount := (len(data) - current) / 32
	addresses := make([]solana.PublicKey, 0, addressCount)
	for i := 0; i < addressCount; i++ {
		addresses = append(addresses, solana.PublicKeyFromBytes(data[current:current+32]))
		current += 32
	}
	ac.Addresses = addresses

	return &ac, nil
}
