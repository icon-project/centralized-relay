package multisig

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/icon-project/goloop/common/codec"
)

const (
	OP_BRIDGE_IDENT = txscript.OP_14
)

type XCallMessage struct {
	MessageType  uint8
	Action       string
	TokenAddress string
	From         string
	To           string
	Amount       []byte
	Data         []byte
}

type BridgeDecodedMsg struct {
	Message    *XCallMessage
	ChainId    uint8
	Receiver   string
	Connectors []string
}

func AddressToPayload(address string) ([]byte, error) {
	prefix := address[0:2]
	var prefixType uint8 // max number of supported type is 7
	switch prefix {
	case "0x":
		prefixType = 1
	case "hx":
		prefixType = 2
	case "cx":
		prefixType = 3
	default:
		return nil, fmt.Errorf("address type not supported")
	}
	addressBytes, err := hex.DecodeString(address[2:])
	if err != nil {
		return nil, fmt.Errorf("could decode string address - Error %v", err)
	}

	addressBytesLen := uint8(len(addressBytes))
	if addressBytesLen > 32 {
		return nil, fmt.Errorf("address length not supported")
	}
	prefixByte := byte((prefixType << 5) ^ (addressBytesLen - 1))

	return append([]byte{prefixByte}, addressBytes...), nil
}

func PayloadToAddress(payload []byte) (string, []byte, error) {
	prefixByte := payload[0]
	prefixType := uint8(prefixByte >> 5)
	var prefix string
	switch prefixType {
	case 1:
		prefix = "0x"
	case 2:
		prefix = "hx"
	case 3:
		prefix = "cx"
	default:
		return "", nil, fmt.Errorf("prefix type not supported")
	}
	addressBytesLen := uint8((prefixByte << 3 >> 3) + 1)
	address := prefix + hex.EncodeToString(payload[1:addressBytesLen+1])
	remainPayload := payload[addressBytesLen+1:]

	return address, remainPayload, nil
}

func CreateBridgePayload(msg *BridgeDecodedMsg) ([]byte, error) {
	payload, err := codec.RLP.MarshalToBytes(msg.Message)
	if err != nil {
		return nil, fmt.Errorf("could not marshal message - Error %v", err)
	}

	payload = append(payload, msg.ChainId)

	receiverBytes, err := AddressToPayload(msg.Receiver)
	if err != nil {
		return nil, err
	}
	payload = append(payload, receiverBytes...)

	for _, connector := range msg.Connectors {
		connectorBytes, err := AddressToPayload(connector)
		if err != nil {
			return nil, err
		}
		payload = append(payload, connectorBytes...)
	}

	return payload, nil
}

func EncodePayloadToScripts(payload []byte) ([][]byte, error) {
	// divide []byte payload to parts
	var chunk []byte
	chunks := make([][]byte, 0, len(payload)/PART_LIMIT+1)
	for len(payload) >= PART_LIMIT {
		chunk, payload = payload[:PART_LIMIT], payload[PART_LIMIT:]
		chunks = append(chunks, chunk)
	}
	if len(payload) > 0 {
		chunks = append(chunks, payload)
	}
	// turn parts to OP_RETURN script
	scripts := [][]byte{}
	for _, part := range chunks {
		builder := txscript.NewScriptBuilder()

		builder.AddOp(txscript.OP_RETURN)
		builder.AddOp(OP_BRIDGE_IDENT)
		builder.AddData(part)

		script, err := builder.Script()
		if err != nil {
			return nil, fmt.Errorf("could not build script - Error %v", err)
		}

		scripts = append(scripts, script)
	}

	return scripts, nil
}

func CreateBridgeMessageScripts(msg *BridgeDecodedMsg) ([][]byte, error) {
	payload, err := CreateBridgePayload(msg)
	if err != nil {
		return nil, err
	}

	return EncodePayloadToScripts(payload)
}

func ReadBridgeMessage(transaction *wire.MsgTx) (*BridgeDecodedMsg, error) {
	payload := []byte{}
	for _, output := range transaction.TxOut {
		tokenizer := txscript.MakeScriptTokenizer(0, output.PkScript)
		if !tokenizer.Next() || tokenizer.Err() != nil || tokenizer.Opcode() != txscript.OP_RETURN {
			// Check for OP_RETURN
			continue
		}
		if !tokenizer.Next() || tokenizer.Err() != nil || tokenizer.Opcode() != OP_BRIDGE_IDENT {
			// Check to ignore non Bridge protocol identifier (Rune or RadFi)
			continue
		}

		// Construct the payload by concatenating remaining data pushes
		for tokenizer.Next() {
			if tokenizer.Err() != nil {
				return nil, tokenizer.Err()
			}
			payload = append(payload, tokenizer.Data()...)
		}
	}

	if len(payload) == 0 {
		return nil, fmt.Errorf("no Bridge message found")
	}

	var message XCallMessage
	remainData, err := codec.RLP.UnmarshalFromBytes(payload, &message)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal message - Error %v", err)
	}
	chainId := uint8(remainData[0])
	receiver, remainData, err := PayloadToAddress(remainData[1:])
	if err != nil {
		return nil, err
	}
	connectors := []string{}
	var connector string
	for len(remainData) > 0 {
		connector, remainData, err = PayloadToAddress(remainData)
		if err != nil {
			return nil, err
		}
		connectors = append(connectors, connector)
	}

	return &BridgeDecodedMsg{
		Message:    &message,
		ChainId:    chainId,
		Receiver:   receiver,
		Connectors: connectors,
	}, nil
}
