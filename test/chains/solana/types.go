package solana

import (
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"
)

var (
	EventLogPrefix        = "Program data: "
	EventSendMessage      = "SendMessage"
	EventCallMessage      = "CallMessage"
	EventRollbackMessage  = "RollbackMessage"
	EventCallMessageSent  = "CallMessageSent"
	EventRollbackExecuted = "RollbackExecuted"
	EventResponseMessage  = "ResponseMessage"
)

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

type SolEvent struct {
	Slot      uint64
	Signature solana.Signature
	Logs      []string
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

type CallMessageSent struct {
	From solana.PublicKey
	To   string
	Sn   big.Int
}

type RollbackExecuted struct {
	Sn big.Int
}

type ResponseMessage struct {
	Code uint8
	Sn   big.Int
}

type CallMessageEvent struct {
	FromNetworkAddress string
	To                 string
	Sn                 big.Int
	ReqId              big.Int
	Data               []byte
}
