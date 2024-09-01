package bitcoin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
)

const (
	UNISAT_DEFAULT_MAINNET	= "https://open-api.unisat.io"
	UNISAT_DEFAULT_TESTNET	= "https://open-api-testnet.unisat.io"
)

type DataBlockchainInfo struct {
	Chain         string `json:"chain"`
	Blocks        int64  `json:"blocks"`
	Headers       int64  `json:"headers"`
	BestBlockHash string `json:"bestBlockHash"`
	PrevBlockHash string `json:"prevBlockHash"`
	Difficulty    string `json:"difficulty"`
	MedianTime    int64  `json:"medianTime"`
	ChainWork     string `json:"chainwork"`
}

type ResponseBlockchainInfo struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data DataBlockchainInfo `json:"data"`
}

type Tx struct {
	TxId             string   `json:"txid"`
	Ins              int      `json:"nIn"`
	Outs             int      `json:"nOut"`
	Size             int      `json:"size"`
	WitOffset        int      `json:"witOffset"`
	Locktime         int      `json:"locktime"`
	InSatoshi        *big.Int `json:"inSatoshi"`
	OutSatoshi       *big.Int `json:"outSatoshi"`
	NewInscriptions  int      `json:"nNewInscription"`
	InInscriptions   int      `json:"nInInscription"`
	OutInscriptions  int      `json:"nOutInscription"`
	LostInscriptions int      `json:"nLostInscription"`
	Timestamp        int64    `json:"timestamp"`
	Height           int64    `json:"height"`
	BlockId          string   `json:"blkid"`
	Index            int      `json:"idx"`
	Confirmations    int      `json:"confirmations"`
}

type ResponseBlockTransactions struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data []Tx `json:"data"`
}

type ResponseTxInfo struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data Tx `json:"data"`
}

type Inscription struct {
	InscriptionNumber int64  `json:"inscriptionNumber"`
	InscriptionId     string `json:"inscriptionId"`
	Offset            int    `json:"offset"`
	Moved             bool   `json:"moved"`
	IsBRC20           bool   `json:"isBRC20"`
}

type Input struct {
	Height       int64         `json:"height"`
	TxId         string        `json:"txid"`
	Index        int           `json:"idx"`
	ScriptSig    string        `json:"scriptSig"`
	ScriptWits   string        `json:"scriptWits"`
	Sequence     int           `json:"sequence"`
	HeightTxo    int64         `json:"heightTxo"`
	Utxid        string        `json:"utxid"`
	Vout         int           `json:"vout"`
	Address      string        `json:"address"`
	CodeType     int           `json:"codeType"`
	Satoshi      *big.Int      `json:"satoshi"`
	ScriptType   string        `json:"scriptType"`
	ScriptPk     string        `json:"scriptPk"`
	Inscriptions []Inscription `json:"inscriptions"`
}

type ResponseTxInputs struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data []Input `json:"data"`
}

type Output struct {
	TxId         string        `json:"txid"`
	Vout         int           `json:"vout"`
	Address      string        `json:"address"`
	CodeType     int           `json:"codeType"`
	Satoshi      *big.Int      `json:"satoshi"`
	ScriptType   string        `json:"scriptType"`
	ScriptPk     string        `json:"scriptPk"`
	Height       int64         `json:"height"`
	Index        int           `json:"idx"`
	Inscriptions []Inscription `json:"inscriptions"`
	TxSpent      string        `json:"txidSpent"`
	HeightSpent  int64         `json:"heightSpent"`
}

type ResponseTxOutputs struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data []Output `json:"data"`
}

type Balance struct {
	Address string `json:"address"`

	Satoshi        *big.Int `json:"satoshi"`
	PendingSatoshi *big.Int `json:"pendingSatoshi"`
	UtxoCount      int64    `json:"utxoCount"`

	BtcSatoshi        *big.Int `json:"btcSatoshi"`
	BtcPendingSatoshi *big.Int `json:"btcPendingSatoshi"`
	BtcUtxoCount      int64    `json:"btcUtxoCount"`

	InscriptionSatoshi        *big.Int `json:"inscriptionSatoshi"`
	InscriptionPendingSatoshi *big.Int `json:"inscriptionPendingSatoshi"`
	InscriptionUtxoCount      int64    `json:"inscriptionUtxoCount"`
}

type ResponseAddressBalance struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data Balance `json:"data"`
}

type UTXO struct {
	TxId         string        `json:"txid"`
	Vout         int           `json:"vout"`
	Satoshi      *big.Int      `json:"satoshi"`
	ScriptType   string        `json:"scriptType"`
	ScriptPk     string        `json:"scriptPk"`
	CodeType     int           `json:"codeType"`
	Address      string        `json:"address"`
	Height       int64         `json:"height"`
	Index        int           `json:"idx"`
	IsOpInRBF    bool          `json:"isOpInRBF"`
	Inscriptions []Inscription `json:"inscriptions"`
}

type DataUtxoList struct {
	Cursor                int    `json:"cursor"`
	Total                 int    `json:"total"`
	TotalConfirmed        int    `json:"totalConfirmed"`
	TotalUnconfirmed      int    `json:"totalUnconfirmed"`
	TotalUnconfirmedSpent int    `json:"totalUnconfirmedSpent"`
	Utxo                  []UTXO `json:"utxo"`
}

type ResponseBtcUtxo struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data DataUtxoList `json:"data"`
}

type RuneDetail struct {
	Amount       string        `json:"amount"`
	RuneId       string        `json:"runeid"`
	Rune         string        `json:"rune"`
	SpacedRune   string        `json:"spacedRune"`
	Symbol       string        `json:"symbol"`
	Divisibility int           `json:"divisibility"`
}

type RuneUTXO struct {
	Address      string        `json:"address"`
	Satoshi      *big.Int      `json:"satoshi"`
	ScriptPk     string        `json:"scriptPk"`
	TxId         string        `json:"txid"`
	Vout         int           `json:"vout"`
	Runes 		[]RuneDetail   `json:"runes"`
}

type DataRuneUtxoList struct {
	Utxo                  []RuneUTXO `json:"utxo"`
}

type ResponseRuneUtxo struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data DataRuneUtxoList `json:"data"`
}

func BtcUtxoUrl(server, address string, offset, limit int64) string {
	return fmt.Sprintf("%s/v1/indexer/address/%s/utxo-data?cursor=%d&size=%d", server, address, offset, limit)
}

func RuneUtxoUrl(server, address, runeId string) string {
	return fmt.Sprintf("%s/v1/indexer/address/%s/runes/%s/utxo", server, address, runeId)
}

func GetWithHeader(ctx context.Context, url string, header map[string]string, response interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, response); err != nil {
		return err
	}
	return nil
}

func GetWithBear(ctx context.Context, url, bear string, response interface{}) error {
	header := make(map[string]string)
	header["Authorization"] = fmt.Sprintf("Bearer %s", bear)
	return GetWithHeader(ctx, url, header, response)
}

func GetBtcUtxo(ctx context.Context, server, bear, address string, offset, limit int64) (ResponseBtcUtxo, error) {
	var resp ResponseBtcUtxo
	url := BtcUtxoUrl(server, address, offset, limit)
	err := GetWithBear(ctx, url, bear, &resp)
	return resp, err
}

func GetRuneUtxo(ctx context.Context, server, address, runeId string) (ResponseRuneUtxo, error) {
	var resp ResponseRuneUtxo
	url := RuneUtxoUrl(server, address, runeId)
	err := GetWithHeader(ctx, url, make(map[string]string), &resp)
	return resp, err
}
