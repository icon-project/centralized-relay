package solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
)

func (idl *IDL) GetInstructionDiscriminator(name string) ([]byte, error) {
	for _, ins := range idl.Instructions {
		if ins.Name == name {
			return ins.Discriminator, nil
		}
	}
	return nil, fmt.Errorf("instruction not found")
}

func (idl *IDL) GetProgramID() solana.PublicKey {
	id, _ := solana.PublicKeyFromBase58(idl.Address)
	return id
}

type IDL struct {
	Address      string           `json:"address"`
	Metadata     IdlMetadata      `json:"metadata"`
	Instructions []IdlInstruction `json:"instructions"`
	Accounts     []IdlAccount     `json:"accounts"`
	Events       []IdlEvent       `json:"events"`
}

type IdlInstruction struct {
	Name          string       `json:"name"`
	Discriminator []byte       `json:"discriminator"`
	Accounts      []IdlAccount `json:"accounts"`
	Args          []IdlField   `json:"args"`
}

type IdlAccount struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	Writable      bool   `json:"writeable"`
	Signer        bool   `json:"signer"`
	Discriminator []byte `json:"discriminator"`
}

type IdlField struct {
	Name string      `json:"name"`
	Type interface{} `json:"type"`
}

type IdlMetadata struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Spec        string `json:"spec"`
	Description string `json:"description"`
}

type IdlEvent struct {
	Name          string `json:"name"`
	Discriminator []byte `json:"discriminator"`
}
