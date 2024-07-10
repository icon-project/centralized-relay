package multisig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/btcsuite/btcd/chaincfg"
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
	MsgTx  string  `json:"msgTx"`
	UTXOs  []*UTXO  `json:"UTXOs"`
	TapSigInfo TapSigParams `json:"tapSigInfo"`
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
	privKeys, _ := randomMultisigInfo(3, 3, chainParam, []int{0, 1, 2})
	msgTx, _ := ParseTx(input.MsgTx)
	sigs, _ := PartSignOnRawExternalTx(privKeys[1], msgTx, input.UTXOs, input.TapSigInfo, chainParam, false)

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
	privKeys, _ := randomMultisigInfo(3, 3, chainParam, []int{0, 1, 2})
	msgTx, _ := ParseTx(input.MsgTx)
	sigs, _ := PartSignOnRawExternalTx(privKeys[2], msgTx, input.UTXOs,  input.TapSigInfo, chainParam, false)

	c.IndentedJSON(http.StatusOK, sigs)
}
