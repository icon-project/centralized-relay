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
	// TODO: verify msg content

	outputs := []*wire.TxOut{
		// bitcoin send to receiver
		{
			Value: new(big.Int).SetBytes(msg.Message.Amount).Int64(),
			PkScript: receiverPkScript,
		},
	}

	bridgeScripts, err := CreateBridgeMessageScripts(msg)
	if err != nil {
		return nil, err
	}
	for _, script := range bridgeScripts {
		outputs = append(outputs, &wire.TxOut{
			Value: 0,
			PkScript: script,
		})
	}

	return CreateTx(inputs, outputs, senderPkScript, txFee, 0)
}

func CreateBridgeTxSendRune(
	msg *BridgeDecodedMsg,
	inputs []*Input,
	senderPkScript []byte,
	receiverPkScript []byte,
	txFee int64,
) (*wire.MsgTx, error) {
	// TODO: verify msg content

	// create runestone OP_RETURN
	runeId, err := runestone.RuneIdFromString(msg.Message.TokenAddress)
	if err != nil {
		return nil, err
	}
	relayerChangeOutput := uint32(1)
	runeOutput := &runestone.Runestone{
		Edicts: []runestone.Edict{
			{
				ID:		*runeId,
				Amount:	uint128.FromBig(new(big.Int).SetBytes(msg.Message.Amount)),
				Output: 0,
			},
		},
		Pointer: &relayerChangeOutput,
	}
	runeScript, _ := runeOutput.Encipher()

	outputs := []*wire.TxOut{
		// rune send to receiver
		{
			Value: DUST_UTXO_AMOUNT,
			PkScript: receiverPkScript,
		},
		// rune change output
		{
			Value: DUST_UTXO_AMOUNT,
			PkScript: senderPkScript,
		},
		// rune OP_RETURN
		{
			Value: 0,
			PkScript: runeScript,
		},
	}

	bridgeScripts, err := CreateBridgeMessageScripts(msg)
	if err != nil {
		return nil, err
	}
	for _, script := range bridgeScripts {
		outputs = append(outputs, &wire.TxOut{
			Value: 0,
			PkScript: script,
		})
	}

	return CreateTx(inputs, outputs, senderPkScript, txFee, 0)
}
