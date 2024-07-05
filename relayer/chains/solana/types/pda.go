package types

import "github.com/gagliardetto/solana-go"

const ()

type PDA struct {
	SeedPrefix string
	ProgramID  solana.PublicKey
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
	XcallReply           PDA
	XcallRollback        PDA
	XcallDefaultConn     PDA
	XcallPendingResponse PDA
	XcallProxyRequest    PDA
	XcallSuccessRes      PDA

	ConnConfig     PDA
	ConnNetworkFee PDA
	ConnClaimFees  PDA
	ConnReceipt    PDA
}

func NewPDARegistry(xcallProgramID, connProgramID solana.PublicKey) *PDARegistry {
	return &PDARegistry{
		XcallConfig:          PDA{SeedPrefix: "config", ProgramID: xcallProgramID},
		XcallReply:           PDA{SeedPrefix: "reply", ProgramID: xcallProgramID},
		XcallRollback:        PDA{SeedPrefix: "rollback", ProgramID: xcallProgramID},
		XcallDefaultConn:     PDA{SeedPrefix: "conn", ProgramID: xcallProgramID},
		XcallPendingResponse: PDA{SeedPrefix: "res", ProgramID: xcallProgramID},
		XcallProxyRequest:    PDA{SeedPrefix: "proxy", ProgramID: xcallProgramID},
		XcallSuccessRes:      PDA{SeedPrefix: "success", ProgramID: xcallProgramID},

		ConnConfig:     PDA{SeedPrefix: "config", ProgramID: connProgramID},
		ConnNetworkFee: PDA{SeedPrefix: "fee", ProgramID: connProgramID},
		ConnClaimFees:  PDA{SeedPrefix: "claim_fees", ProgramID: connProgramID},
		ConnReceipt:    PDA{SeedPrefix: "receipt", ProgramID: connProgramID},
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
