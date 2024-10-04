package bitcoin

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"

	"github.com/icon-project/icon-bridge/common/codec"
)

func GetRuneTxIndex(endpoint, method, bearToken, txId string, index int) (*RuneTxIndexResponse, error) {
	client := &http.Client{}
	endpoint = endpoint + "/runes/utxo/" + txId + "/" + strconv.FormatUint(uint64(index), 10) + "/balance"
	fmt.Println(endpoint)
	req, err := http.NewRequest(method, endpoint, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Authorization", bearToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var resp *RuneTxIndexResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp, nil
}

func XcallFormat(callData []byte, from, to string, sn uint, protocols []string, messType uint8) ([]byte, error) {
	//
	csV2 := CSMessageRequestV2{
		From:        from,
		To:          to,
		Sn:          big.NewInt(int64(sn)).Bytes(),
		MessageType: messType,
		Data:        callData,
		Protocols:   protocols,
	}

	//
	cvV2EncodeMsg, err := codec.RLP.MarshalToBytes(csV2)
	if err != nil {
		return nil, err
	}

	message := CSMessage{
		MsgType: big.NewInt(int64(CS_REQUEST)).Bytes(),
		Payload: cvV2EncodeMsg,
	}

	//
	finalMessage, err := codec.RLP.MarshalToBytes(message)
	if err != nil {
		return nil, err
	}

	return finalMessage, nil
}

func mulDiv(a, nNumerator, nDenominator *big.Int) *big.Int {
	return big.NewInt(0).Div(big.NewInt(0).Mul(a, nNumerator), nDenominator)
}
