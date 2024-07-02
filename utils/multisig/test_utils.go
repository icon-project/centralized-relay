package multisig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine{
    router := gin.Default()
	router.POST("/requestSign1", postRequestSignSlaveRelayer1)
	router.POST("/requestSign2", postRequestSignSlaveRelayer2)

    return router
}

func requestSign(url string, requestJson []byte, router *gin.Engine) [][]byte{
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestJson))
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    responseData, _ := io.ReadAll(w.Body)

	sigs := [][]byte{}
	err := json.Unmarshal(responseData, &sigs)
	if err != nil {
		fmt.Println("err Unmarshal: ", err)
	}

	return sigs
}

type requestSignInput struct {
    MsgTx  *wire.MsgTx  `json:"msgTx"`
    Inputs  []*UTXO  `json:"inputs"`
}

func postRequestSignSlaveRelayer1(c *gin.Context) {
    var input requestSignInput

	if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

	chainParam := &chaincfg.RegressionNetParams
	privKeys, multisigInfo := randomMultisigInfo(3, 2, chainParam, []int{0, 1, 2})
	multisigWallet, _ := GenerateMultisigWallet(multisigInfo)

	sigs, _ := PartSignOnRawExternalTx(privKeys[1], input.MsgTx, input.Inputs, multisigWallet.PKScript, multisigWallet.TapLeaves[0], nil, txscript.TapLeaf{}, chainParam, false)
    c.IndentedJSON(http.StatusOK, sigs)
}

func postRequestSignSlaveRelayer2(c *gin.Context) {
    var input requestSignInput

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

	chainParam := &chaincfg.RegressionNetParams
	privKeys, multisigInfo := randomMultisigInfo(3, 2, chainParam, []int{0, 1, 2})
	multisigWallet, _ := GenerateMultisigWallet(multisigInfo)

	sigs, _ := PartSignOnRawExternalTx(privKeys[2], input.MsgTx, input.Inputs, multisigWallet.PKScript, multisigWallet.TapLeaves[0], nil, txscript.TapLeaf{}, chainParam, false)
    c.IndentedJSON(http.StatusOK, sigs)
}

func UserSignTx(
	privKey string,
	msgTx *wire.MsgTx,
	inputs []*UTXO,
	multisigWallet *MultisigWallet,
	indexTapLeaf int,
	chainParam *chaincfg.Params,
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

		prevOuts.AddPrevOut(*outPoint, &wire.TxOut{
			Value:    int64(in.OutputAmount),
			PkScript: multisigWallet.PKScript,
		})
	}
	txSigHashes := txscript.NewTxSigHashes(msgTx, prevOuts)

	sigs := [][]byte{}
	for i := range msgTx.TxIn {
		if (!inputs[i].IsRelayersMultisig) {
			sig, err := txscript.RawTxInTapscriptSignature(
				msgTx, txSigHashes, i, int64(inputs[i].OutputAmount), multisigWallet.PKScript, multisigWallet.TapLeaves[indexTapLeaf], txscript.SigHashAll, wif.PrivKey)
			if err != nil {
				return nil, fmt.Errorf("[PartSignOnRawExternalTx] Error when signing on raw btc tx: %v", err)
			}
			fmt.Printf("PartSignOnRawExternalTx sig len : %v\n", len(sig))

			sigs = append(sigs, sig)
		} else {
			sigs = append(sigs, []byte{})
		}
	}

	return sigs, nil
}