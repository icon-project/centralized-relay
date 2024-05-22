package solana

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

func fetchIDL(idlAccountAddr string, cl IClient) (*IDL, error) {
	idlPubKey, err := solana.PublicKeyFromBase58(idlAccountAddr)
	if err != nil {
		return nil, err
	}
	idlAccount, err := cl.GetAccountInfo(context.Background(), idlPubKey)
	if err != nil {
		return nil, err
	}

	data := idlAccount.Data.GetBinary()[8:] //skip discriminator

	idlAcInfo := struct {
		Authority solana.PublicKey
		DataLen   uint32
	}{}
	if err := borsh.Deserialize(&idlAcInfo, data); err != nil {
		return nil, err
	}

	compressedBytes := data[36 : 36+idlAcInfo.DataLen] //skip authority and unwanted trailing bytes

	decompressedBytes, err := decompress(compressedBytes)
	if err != nil {
		return nil, err
	}

	var idlData IDL
	if err = json.Unmarshal(decompressedBytes, &idlData); err != nil {
		return nil, err
	}

	return &idlData, nil

}

func decompress(compressedData []byte) ([]byte, error) {
	// Create a new bytes reader from the compressed data
	b := bytes.NewReader(compressedData)

	// Create a new zlib reader
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// Create a buffer to hold the decompressed data
	var out bytes.Buffer

	// Copy the decompressed data into the buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return nil, err
	}

	// Return the decompressed data
	return out.Bytes(), nil
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
	Name string `json:"name"`
	Type string `json:"type"`
}

// IdlDS DS means data structure
type IdlDS struct {
	Kind   string     `json:"kind"`
	Fields []IdlField `json:"fields"`
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
