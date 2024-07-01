package solana

import (
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"
)

var (
	EventLogPrefix       = "Program data: "
	EventSendMessage     = "SendMessage"
	EventCallMessage     = "CallMessage"
	EventRollbackMessage = "RollbackMessage"
)

type LatestLedgerResponse struct {
	ID              string `json:"id"`
	ProtocolVersion uint64 `json:"protocolVersion"`
	Sequence        uint64 `json:"sequence"`
}

type EventResponseEvent struct {
	Type         string   `json:"type"`
	ContractId   string   `json:"contractID"`
	Topic        []string `json:"topic"`
	Value        string   `json:"value"`
	ValueDecoded map[string]interface{}
}

type EventResponse struct {
	Events []EventResponseEvent `json:"events"`
}

type EventQueryFilter struct {
	StartLedger uint64     `json:"startLedger"`
	Pagination  Pagination `json:"pagination"`
	Filters     []Filter   `json:"filters"`
}

type Pagination struct {
	Limit uint64 `json:"limit"`
}

type Filter struct {
	Type        string   `json:"type"`
	ContractIDS []string `json:"contractIds"`
}

type SolEvent struct {
	Slot      uint64
	Signature solana.Signature
	Logs      []string
}

type SendMessageEvent struct {
	TargetNetwork string
	ConnSn        big.Int
	Msg           []byte
}

type IDL struct {
	Address      string           `json:"address"`
	Metadata     IdlMetadata      `json:"metadata"`
	Instructions []IdlInstruction `json:"instructions"`
	Accounts     []IdlAccount     `json:"accounts"`
	Events       []IdlEvent       `json:"events"`
	Types        []IdlType        `json:"types"`
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

// IdlDS DS means data structure
type IdlDS struct {
	Kind   string        `json:"kind"`
	Fields []interface{} `json:"fields"`
}

type IdlType struct {
	Name string `json:"name"`
	Type IdlDS  `json:"type"`
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

func (idl *IDL) GetInstructionDiscriminator(name string) ([]byte, error) {
	for _, ins := range idl.Instructions {
		if ins.Name == name {
			return ins.Discriminator, nil
		}
	}
	return nil, fmt.Errorf("instruction not found")
}

func (idl *IDL) GetProgramID() (solana.PublicKey, error) {
	return solana.PublicKeyFromBase58(idl.Address)
}
