package multisig

import (
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
	Fee			[]byte
	UpperTick	[]byte
	LowerTick	[]byte
	Min0		[]byte
	Min1		[]byte
}

type RadFiDecodedMsg struct {
	Flag				[]byte
	ProvideLiquidityMsg	RadFiProvideLiquidityMsg
}

// func createProvideLiquidityScript(fee uint8, upperTick uint32) ([]byte, error) {
// 	builder := txscript.NewScriptBuilder()

// 	builder.AddOp(txscript.OP_RETURN)
// 	builder.AddOp(OP_RADFI_IDENT)

// 	return builder.Script()
// }

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

func readRelayMessage(transaction *wire.MsgTx, isRadFi bool) ([]byte, error) {
	var payload []byte
	for _, output := range transaction.TxOut {
		tokenizer := txscript.MakeScriptTokenizer(0, output.PkScript)
		if !tokenizer.Next() || tokenizer.Err() != nil || tokenizer.Opcode() != txscript.OP_RETURN {
			// Check for OP_RETURN
			continue
		}
		if !tokenizer.Next() || tokenizer.Err() != nil || tokenizer.Opcode() == OP_RUNE_IDENT {
			// Check to ignore Rune protocol identifier (Runestone::MAGIC_NUMBER)
			continue
		}

		if (isRadFi && tokenizer.Opcode() != OP_RADFI_IDENT) || ((!isRadFi && tokenizer.Opcode() != OP_BRIDGE_IDENT)){
			// Check for Relayer protocol identifier (RadFi or ICON Bridge)
			continue
		}

		// Construct the payload by concatenating remaining data pushes
		for tokenizer.Next() {
			if tokenizer.Err() != nil {
				return nil, tokenizer.Err()
			}
			payload = append(payload, tokenizer.Data()...)
		}

		// only read 1 message output for radfi protocol
		if isRadFi {
			break
		}
	}

	return payload, nil
}

// func decodeRadFiMessage(payload []byte) (*DecodedMsg, error) {
// 	tokenizer := txscript.MakeScriptTokenizer(0, payload)

// 	// take the flag
// 	if !tokenizer.Next() || tokenizer.Err() != nil {
// 		return nil, fmt.Errorf("decodeRadFiMessage could not read the flag - Error %v", tokenizer.Err())
// 	}
// 	flag := tokenizer.Opcode()

// 	switch flag {
//     case OP_RADFI_PROVIDE_LIQUIDITY:
// 		if !tokenizer.Next() || tokenizer.Err() != nil {
// 			return nil, fmt.Errorf("decodeRadFiMessage could not read the ProvideLiquidityMsg Fee - Error %v", tokenizer.Err())
// 		}
// 		fee :=
//         return &DecodedMsg {
// 			Flag		: flag,
// 			ProvideLiquidityMsg: RadFiProvideLiquidityMsg {
// 				Fee			: ,
// 				UpperTick	: ,
// 				LowerTick	: ,
// 				Min0		: ,
// 				Min1		: ,
// 			}
// 		}
//     case 2:
//         fmt.Println("two")
//     case 3:
//         fmt.Println("three")
//     default:
//         fmt.Println("It's after noon")
//     }


// 	return payload, nil
// }

// func nextToken(tokenizer *txscript.ScriptTokenizer) ([]byte, error) {
// 	if !tokenizer.Next() || tokenizer.Err() != nil {
// 		return nil, fmt.Errorf("decodeRadFiMessage could not read the flag - Error %v", tokenizer.Err())
// 	}

// 	return tokenizer.Opcode()
// }