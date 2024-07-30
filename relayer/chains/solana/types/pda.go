package types

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
)

const ()

type PDA struct {
	SeedPrefix string
	ProgramID  solana.PublicKey

	// if address is temporary then it can be deactivated from address lookup table onced used.
	IsTemp bool
}

type AddressTables map[solana.PublicKey]solana.PublicKeySlice

func GetPDA(progID solana.PublicKey, prefix string, additionalSeeds ...[]byte) (solana.PublicKey, error) {
	seeds := [][]byte{[]byte(prefix)}

	seeds = append(seeds, additionalSeeds...)

	addr, _, err := solana.FindProgramAddress(seeds, progID)
	if err != nil {
		return solana.PublicKey{}, err
	}

	return addr, nil
}

func (pda PDA) GetAddress(additionalSeeds ...[]byte) (solana.PublicKey, error) {
	seeds := [][]byte{[]byte(pda.SeedPrefix)}

	seeds = append(seeds, additionalSeeds...)

	addr, _, err := solana.FindProgramAddress(seeds, pda.ProgramID)
	if err != nil {
		return solana.PublicKey{}, err
	}

	return addr, nil
}

type PDARegistry struct {
	XcallConfig          PDA
	XcallRollback        PDA
	XcallProxyRequest    PDA //temp
	XcallPendingRequest  PDA //temp
	XcallPendingResponse PDA //temp

	ConnConfig     PDA
	ConnNetworkFee PDA
	ConnClaimFees  PDA
	ConnReceipt    PDA //temp
}

func NewPDARegistry(xcallProgramID, connProgramID solana.PublicKey) *PDARegistry {
	return &PDARegistry{
		XcallConfig:          PDA{SeedPrefix: "config", ProgramID: xcallProgramID},
		XcallRollback:        PDA{SeedPrefix: "rollback", ProgramID: xcallProgramID},
		XcallProxyRequest:    PDA{SeedPrefix: "proxy", ProgramID: xcallProgramID, IsTemp: true},
		XcallPendingRequest:  PDA{SeedPrefix: "req", ProgramID: xcallProgramID, IsTemp: true},
		XcallPendingResponse: PDA{SeedPrefix: "res", ProgramID: xcallProgramID, IsTemp: true},

		ConnConfig:     PDA{SeedPrefix: "config", ProgramID: connProgramID},
		ConnNetworkFee: PDA{SeedPrefix: "fee", ProgramID: connProgramID},
		ConnClaimFees:  PDA{SeedPrefix: "claim_fees", ProgramID: connProgramID},
		ConnReceipt:    PDA{SeedPrefix: "receipt", ProgramID: connProgramID, IsTemp: true},
	}
}

type ContractRollback struct {
	From      solana.PublicKey
	To        string
	Protocols []string
	Rollback  []byte
	Enabled   bool
}

type RollbackAccount struct {
	Rollback ContractRollback
	Owner    solana.PublicKey
	Bump     uint8
}

type XcallConfigAccount struct {
	Admin       solana.PublicKey
	FeeHandler  solana.PublicKey
	NetworkID   string
	ProtocolFee uint64
	SequenceNo  big.Int
	LastReqID   big.Int
	Bump        uint8
}
