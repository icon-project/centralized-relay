package bitcoin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetRuneTxIndex(endpoint, method, bearToken, txId string, index int) (*RuneTxIndexResponse, error) {
	client := &http.Client{}
	endpoint = endpoint + "/utxo/" + txId + "/" + strconv.FormatUint(uint64(index), 10) + "/balance"
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

	body, err := ioutil.ReadAll(res.Body)
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
