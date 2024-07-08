package multisig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/btcsuite/btcd/chaincfg"
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

	chainParam := &chaincfg.SigNetParams
	privKeys, multisigInfo := randomMultisigInfo(3, 3, chainParam, []int{0, 1, 2})
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

	chainParam := &chaincfg.SigNetParams
	privKeys, multisigInfo := randomMultisigInfo(3, 3, chainParam, []int{0, 1, 2})
	multisigWallet, _ := GenerateMultisigWallet(multisigInfo)

	sigs, _ := PartSignOnRawExternalTx(privKeys[2], input.MsgTx, input.Inputs, multisigWallet.PKScript, multisigWallet.TapLeaves[0], nil, txscript.TapLeaf{}, chainParam, false)
    c.IndentedJSON(http.StatusOK, sigs)
}
