package types

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
)

const (
	PrefixConfig       = "config"
	PrefixRollback     = "rollback"
	PrefixProxyRequest = "proxy"
	PrefixPendingReq   = "req"
	PrefixPendingRes   = "res"
	PrefixNetworkFee   = "fee"
	PrefixClaimFees    = "claim_fees"
	PrefixReceipt      = "receipt"

	PrefixState = "state"

	PrefixVaultNative    = "vault_native"
	PrefixBnUSDAuthority = "bnusd_authority"
)

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
		XcallConfig:          PDA{SeedPrefix: PrefixConfig, ProgramID: xcallProgramID},
		XcallRollback:        PDA{SeedPrefix: PrefixRollback, ProgramID: xcallProgramID},
		XcallProxyRequest:    PDA{SeedPrefix: PrefixProxyRequest, ProgramID: xcallProgramID, IsTemp: true},
		XcallPendingRequest:  PDA{SeedPrefix: PrefixPendingReq, ProgramID: xcallProgramID, IsTemp: true},
		XcallPendingResponse: PDA{SeedPrefix: PrefixPendingRes, ProgramID: xcallProgramID, IsTemp: true},

		ConnConfig:     PDA{SeedPrefix: PrefixConfig, ProgramID: connProgramID},
		ConnNetworkFee: PDA{SeedPrefix: PrefixNetworkFee, ProgramID: connProgramID},
		ConnClaimFees:  PDA{SeedPrefix: PrefixClaimFees, ProgramID: connProgramID},
		ConnReceipt:    PDA{SeedPrefix: PrefixReceipt, ProgramID: connProgramID, IsTemp: true},
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
