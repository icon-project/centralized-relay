package types

import "github.com/gagliardetto/solana-go"

const ()

type PDA struct {
	SeedPrefix string
	Seeds      []string
	ProgramID  solana.PublicKey
}

func (pda PDA) GetAddress() (solana.PublicKey, error) {
	seeds := [][]byte{[]byte(pda.SeedPrefix)}
	for _, sd := range pda.Seeds {
		seeds = append(seeds, []byte(sd))
	}

	addr, _, err := solana.FindProgramAddress(seeds, pda.ProgramID)
	if err != nil {
		return solana.PublicKey{}, err
	}

	return addr, nil
}

type PDARegistry struct {
	XcallConfig PDA
	XcallReply  PDA

	ConnConfig PDA
}

func NewPDARegistry(xcallProgramID, connProgramID solana.PublicKey) *PDARegistry {
	return &PDARegistry{
		XcallConfig: PDA{SeedPrefix: "config", ProgramID: xcallProgramID},
		XcallReply:  PDA{SeedPrefix: "reply", ProgramID: xcallProgramID},

		ConnConfig: PDA{SeedPrefix: "config", ProgramID: connProgramID},
	}
}
