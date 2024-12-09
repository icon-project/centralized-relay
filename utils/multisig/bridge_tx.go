package multisig

import (
	"math/big"

	"github.com/btcsuite/btcd/wire"
	"github.com/bxelab/runestone"
	"lukechampine.com/uint128"
)

func CreateBridgeTxSendBitcoin(
	msg *BridgeDecodedMsg,
	inputs []*Input,
	senderPkScript []byte,
	receiverPkScript []byte,
	txFee int64,
) (*wire.MsgTx, error) {
	outputs := []*wire.TxOut{
		// bitcoin send to receiver
		{
			Value:    new(big.Int).SetBytes(msg.Message.Amount).Int64(),
			PkScript: receiverPkScript,
		},
	}

	bridgeScripts, err := CreateBridgeMessageScripts(msg)
	if err != nil {
		return nil, err
	}
	outputs = BuildBridgeScriptsOutputs(outputs, bridgeScripts)

	return CreateTx(inputs, outputs, senderPkScript, txFee, 0)
}

func BuildBridgeScriptsOutputs(outputs []*wire.TxOut, bridgeScripts [][]byte) []*wire.TxOut {
	for _, script := range bridgeScripts {
		outputs = append(outputs, &wire.TxOut{
			Value:    OP_MIN_DUST_AMOUNT,
			PkScript: script,
		})
	}
	return outputs
}

func CreateBridgeTxSendRune(
	msg *BridgeDecodedMsg,
	inputs []*Input,
	senderPkScript []byte,
	receiverPkScript []byte,
	txFee int64,
) (*wire.MsgTx, error) {
	// create runestone OP_RETURN
	runeId, err := runestone.RuneIdFromString(msg.Message.TokenAddress)
	if err != nil {
		return nil, err
	}
	relayerChangeOutput := uint32(1)
	runeOutput := &runestone.Runestone{
		Edicts: []runestone.Edict{
			{
				ID:     *runeId,
				Amount: uint128.FromBig(new(big.Int).SetBytes(msg.Message.Amount)),
				Output: 0,
			},
		},
		Pointer: &relayerChangeOutput,
	}
	runeScript, _ := runeOutput.Encipher()

	outputs := []*wire.TxOut{
		// rune send to receiver
		{
			Value:    RUNE_DUST_UTXO_AMOUNT,
			PkScript: receiverPkScript,
		},
		// rune change output
		{
			Value:    RUNE_DUST_UTXO_AMOUNT,
			PkScript: senderPkScript,
		},
		// rune OP_RETURN
		{
			Value:    0,
			PkScript: runeScript,
		},
	}

	bridgeScripts, err := CreateBridgeMessageScripts(msg)
	if err != nil {
		return nil, err
	}
	outputs = BuildBridgeScriptsOutputs(outputs, bridgeScripts)
	return CreateTx(inputs, outputs, senderPkScript, txFee, 0)
}
