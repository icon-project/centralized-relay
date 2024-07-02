package multisig

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// create unsigned multisig transaction
// input: UTXO, output, tx fee, chain config, change receiver
// output: unsigned tx message, previous output fetcher
func CreateMultisigTx(
    inputs []*UTXO,
    outputs []*OutputTx,
    feePerOutput uint64,
    chainParam *chaincfg.Params,
    changeReceiverAddress string,
    relayersPKScript []byte,
	userPKScript []byte,
) (*wire.MsgTx, *txscript.MultiPrevOutFetcher, error) {
	msgTx := wire.NewMsgTx(wire.TxVersion)

	// add TxIns into raw tx
	// totalInputAmount in external unit
	totalInputAmount := uint64(0)
	prevOuts := txscript.NewMultiPrevOutFetcher(nil)
	for _, in := range inputs {
		utxoHash, err := chainhash.NewHashFromStr(in.TxHash)
		if err != nil {
			return nil, nil, err
		}
		outPoint := wire.NewOutPoint(utxoHash, in.OutputIdx)
		txIn := wire.NewTxIn(outPoint, nil, nil)
		txIn.Sequence = uint32(feePerOutput)
		msgTx.AddTxIn(txIn)
		totalInputAmount += in.OutputAmount

		var pkScript []byte
		if (in.IsRelayersMultisig) {
			pkScript = relayersPKScript
		} else {
			pkScript = userPKScript
		}

		prevOuts.AddPrevOut(*outPoint, &wire.TxOut{
			Value:    int64(in.OutputAmount),
			PkScript: pkScript,
		})
	}

	// add TxOuts into raw tx
	// totalOutputAmount in external unit
	totalOutputAmount := uint64(0)
	for _, out := range outputs {
		// adding the output to tx
		decodedAddr, err := btcutil.DecodeAddress(out.ReceiverAddress, chainParam)
		if err != nil {
			return nil, nil, fmt.Errorf("CreateRawExternalTx - Error when decoding receiver address: %v - %v", err, out.ReceiverAddress)
		}
		destinationAddrByte, err := txscript.PayToAddrScript(decodedAddr)
		if err != nil {
			return nil, nil, err
		}

		// adding the destination address and the amount to the transaction
		if out.Amount <= feePerOutput || out.Amount-feePerOutput < MIN_SAT {
			return nil, nil, fmt.Errorf("[CreateRawExternalTx-BTC] Output amount %v must greater than fee %v", out.Amount, feePerOutput)
		}
		redeemTxOut := wire.NewTxOut(int64(out.Amount-feePerOutput), destinationAddrByte)

		msgTx.AddTxOut(redeemTxOut)
		totalOutputAmount += out.Amount
	}

	// check amount of input coins and output coins
	if totalInputAmount < totalOutputAmount {
		return nil, nil, fmt.Errorf("[CreateRawExternalTx-BTC] Total input amount %v is less than total output amount %v", totalInputAmount, totalOutputAmount)
	}

	// calculate the change output
	changeAmt := uint64(0)
	if totalInputAmount > totalOutputAmount {
		// adding the output to tx
		decodedAddr, err := btcutil.DecodeAddress(changeReceiverAddress, chainParam)
		if err != nil {
			return nil, nil, err
		}
		destinationAddrByte, err := txscript.PayToAddrScript(decodedAddr)
		if err != nil {
			return nil, nil, err
		}

		// adding the destination address and the amount to the transaction
		changeAmt = totalInputAmount - totalOutputAmount
		if changeAmt >= MIN_SAT {
			redeemTxOut := wire.NewTxOut(int64(changeAmt), destinationAddrByte)
			msgTx.AddTxOut(redeemTxOut)
		} else {
			changeAmt = 0
		}
	}

	var rawTxBytes bytes.Buffer
	err := msgTx.Serialize(&rawTxBytes)
	if err != nil {
		return nil, nil, err
	}

	return msgTx, prevOuts, nil
}

func PartSignOnRawExternalTx(
	privKey string,
	msgTx *wire.MsgTx,
	inputs []*UTXO,
	relayersPKScript []byte,
	relayersTapLeaf txscript.TapLeaf,
	userPKScript []byte,
	userTapLeaf txscript.TapLeaf,
	chainParam *chaincfg.Params,
	isMasterRelayer bool,
) ([][]byte, error) {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return nil, fmt.Errorf("[PartSignOnRawExternalTx] Error when generate btc private key from seed: %v", err)
	}

	// sign on each TxIn
	if len(inputs) != len(msgTx.TxIn) {
		return nil, fmt.Errorf("[PartSignOnRawExternalTx] Len of Public seeds %v and len of TxIn %v are not correct", len(inputs), len(msgTx.TxIn))
	}

	prevOuts := txscript.NewMultiPrevOutFetcher(nil)
	for _, in := range inputs {
		utxoHash, err := chainhash.NewHashFromStr(in.TxHash)
		if err != nil {
			return nil, err
		}
		outPoint := wire.NewOutPoint(utxoHash, in.OutputIdx)

		var pkScript []byte
		if (in.IsRelayersMultisig) {
			pkScript = relayersPKScript
		} else {
			pkScript = userPKScript
		}

		prevOuts.AddPrevOut(*outPoint, &wire.TxOut{
			Value:    int64(in.OutputAmount),
			PkScript: pkScript,
		})
	}
	txSigHashes := txscript.NewTxSigHashes(msgTx, prevOuts)

	sigs := [][]byte{}
	for i := range msgTx.TxIn {
		if (inputs[i].IsRelayersMultisig) {
			sig, err := txscript.RawTxInTapscriptSignature(
				msgTx, txSigHashes, i, int64(inputs[i].OutputAmount), relayersPKScript, relayersTapLeaf, txscript.SigHashAll, wif.PrivKey)
			if err != nil {
				return nil, fmt.Errorf("[PartSignOnRawExternalTx] Error when signing on raw btc tx: %v", err)
			}

			sigs = append(sigs, sig)
		} else if (isMasterRelayer) {
			sig, err := txscript.RawTxInTapscriptSignature(
				msgTx, txSigHashes, i, int64(inputs[i].OutputAmount), userPKScript, userTapLeaf, txscript.SigHashAll, wif.PrivKey)
			if err != nil {
				return nil, fmt.Errorf("[PartSignOnRawExternalTx] Error when signing on raw btc tx: %v", err)
			}

			sigs = append(sigs, sig)
		} else {
			sigs = append(sigs, []byte{})
		}
	}

	return sigs, nil
}

func TransposeSigs(sigs [][][]byte) [][][]byte {
	xl := len(sigs[0])
    yl := len(sigs)
    result := make([][][]byte, xl)

    for i := range result {
        result[i] = make([][]byte, yl)
    }
    for i := 0; i < xl; i++ {
        for j := 0; j < yl; j++ {
            result[i][j] = sigs[j][i]
        }
    }
    return result
}

// combine all the signatures to create the signed tx
// input: unsigned tx message, signatures, multisig struct
// output: signed tx message, tx hash
func CombineMultisigSigs(
	relayersMultisigTapLeafScript []byte,
	relayersMultisigControlBlock []byte,
	userMultisigTapLeafScript []byte,
	userMultisigControlBlock []byte,
    msgTx *wire.MsgTx,
	inputs []*UTXO,
    transposedSigs [][][]byte,
) (*wire.MsgTx, error) {
	for idxInput, v := range transposedSigs {
		reverseV := [][]byte{}
		for i := len(v) - 1; i >= 0; i-- {
			if (len(v[i]) != 0) {
				reverseV = append(reverseV, v[i])
			}
		}

		witness := append([][]byte{}, reverseV...)

		if (inputs[idxInput].IsRelayersMultisig) {
			witness = append(witness, relayersMultisigTapLeafScript, relayersMultisigControlBlock)
		} else {
			witness = append(witness, userMultisigTapLeafScript, userMultisigControlBlock)
		}

		msgTx.TxIn[idxInput].Witness = witness
	}

	return msgTx, nil
}

func ParseTx(data string) (*wire.MsgTx, error) {
	fmt.Printf("ParseTx data: %v\n", string(data))
	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	fmt.Printf("ParseTx dataBytes: %v\n", string(dataBytes))
	tx := &wire.MsgTx{}
	err = tx.Deserialize(strings.NewReader(string(dataBytes)))
	if err != nil {
		return nil, err
	}
	return tx, nil
}
