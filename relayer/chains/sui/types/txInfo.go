package types

import "encoding/json"

type TxInfo struct {
	TxDigest string `json:"tx_digest"`
}

func (txi *TxInfo) Serialize() ([]byte, error) {
	return json.Marshal(txi)
}

func (txi *TxInfo) Deserialize(bytesVal []byte) error {
	return json.Unmarshal(bytesVal, &txi)
}
