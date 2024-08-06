package multisig

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const (
	OP_RADFI_IDENT				= txscript.OP_12
	OP_RUNE_IDENT				= txscript.OP_13
	OP_BRIDGE_IDENT				= txscript.OP_14

	OP_RADFI_PROVIDE_LIQUIDITY	= txscript.OP_1
	OP_RADFI_SWAP				= txscript.OP_2
	OP_RADFI_WITHDRAW_LIQUIDITY	= txscript.OP_3
	OP_RADFI_COLLECT_FEE		= txscript.OP_4
)

type RadFiProvideLiquidityMsg struct {
	Fee			uint8
	UpperTick	int32
	LowerTick	int32
	Min0		uint16
	Min1		uint16
}

type RadFiDecodedMsg struct {
	Flag				byte
	ProvideLiquidityMsg	*RadFiProvideLiquidityMsg
}

func CreateBridgeMessageScripts(payload []byte, partLimit int) ([][]byte, error) {
	var chunk []byte
	chunks := make([][]byte, 0, len(payload)/partLimit+1)
	for len(payload) >= partLimit {
		chunk, payload = payload[:partLimit], payload[partLimit:]
		chunks = append(chunks, chunk)
	}
	if len(payload) > 0 {
		chunks = append(chunks, payload)
	}

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

func CreateProvideLiquidityScript(fee uint8, upperTick int32, lowerTick int32, min0 uint16, min1 uint16) ([]byte, error) {
	builder := txscript.NewScriptBuilder()

	builder.AddOp(txscript.OP_RETURN)
	builder.AddOp(OP_RADFI_IDENT)
	builder.AddOp(OP_RADFI_PROVIDE_LIQUIDITY)
	// encode message content
	buf := new(bytes.Buffer)
	var data = []any{ fee, upperTick, lowerTick, min0, min1 }
	for _, v := range data {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			fmt.Println("CreateProvideLiquidityScript encode data failed:", err)
		}
	}

	return builder.AddData(buf.Bytes()).Script()
}

func ReadBridgeMessage(transaction *wire.MsgTx) ([]byte, error) {
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
		return nil, fmt.Errorf("ReadBridgeMessage - no Bridge message found")
	}

	return payload, nil
}

func ReadRadFiMessage(transaction *wire.MsgTx) (*RadFiDecodedMsg, error) {
	var flag byte
	var payload []byte
	for _, output := range transaction.TxOut {
		tokenizer := txscript.MakeScriptTokenizer(0, output.PkScript)
		if !tokenizer.Next() || tokenizer.Err() != nil || tokenizer.Opcode() != txscript.OP_RETURN {
			// Check for OP_RETURN
			continue
		}
		if !tokenizer.Next() || tokenizer.Err() != nil || tokenizer.Opcode() != OP_RADFI_IDENT {
			// Check to ignore non RadFi protocol identifier (Rune or Bridge)
			continue
		}

		if tokenizer.Next() && tokenizer.Err() == nil {
			flag = tokenizer.Opcode()
		}

		// Construct the payload by concatenating remaining data pushes
		for tokenizer.Next() {
			if tokenizer.Err() != nil {
				return nil, tokenizer.Err()
			}
			payload = append(payload, tokenizer.Data()...)
		}

		// only read 1 OP_RETURN output for RadFi protocol
		break
	}

	// Decode RadFi message
	r := bytes.NewReader(payload)
	switch flag {
		case OP_RADFI_PROVIDE_LIQUIDITY:
			var provideLiquidityMsg RadFiProvideLiquidityMsg
			if err := binary.Read(r, binary.LittleEndian, &provideLiquidityMsg); err != nil {
				fmt.Println("OP_RADFI_PROVIDE_LIQUIDITY Read failed:", err)
			}

			return &RadFiDecodedMsg {
				Flag		: flag,
				ProvideLiquidityMsg: &provideLiquidityMsg,
			}, nil

		case OP_RADFI_SWAP:
			fmt.Println("OP_RADFI_SWAP")

		case OP_RADFI_WITHDRAW_LIQUIDITY:
			fmt.Println("OP_RADFI_WITHDRAW_LIQUIDITY")

		case OP_RADFI_COLLECT_FEE:
			fmt.Println("OP_RADFI_COLLECT_FEE")

		default:
			return nil, fmt.Errorf("ReadRadFiMessage - invalid flag")
	}

	return nil, fmt.Errorf("ReadRadFiMessage - no RadFi message found")
}
