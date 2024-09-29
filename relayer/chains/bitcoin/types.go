package bitcoin

import (
	"github.com/btcsuite/btcd/wire"
)

type MessageType int

const (
	CS_REQUEST MessageType = iota + 1
	CS_RESPONSE
	CS_RESULT
)

type CallMessageType int

const (
	CALL_MESSAGE_TYPE CallMessageType = iota
	CALL_MESSAGE_ROLLBACK_TYPE
	PERSISTENT_MESSAGE_TYPE
)

type TxSearchParam struct {
	StartHeight, EndHeight uint64
	BitcoinScript          []byte
	OPReturnPrefix         int
}

type TxSearchRes struct {
	Tx      *wire.MsgTx
	Height  uint64
	TxIndex uint64
}

// HightRange is a struct to represent a range of heights
type HeightRange struct {
	Start uint64
	End   uint64
}

type XCallMessage struct {
	Action       string
	TokenAddress string
	From         string
	To           string
	Amount       []byte
	Data         []byte
}

type RuneInfo struct {
	Rune         string `json:"rune"`
	RuneId       string `json:"runeId"`
	SpaceRune    string `json:"spaceRune"`
	Amount       string `json:"amount"`
	Symbol       string `json:"symbol"`
	Divisibility int    `json:"divisibility"`
}

type RuneTxIndexResponse struct {
	Code int        `json:"code"`
	Data []RuneInfo `json:"data"`
}

type CSMessageRequestV2 struct {
	From        string   `json:"from"`
	To          string   `json:"to"`
	Sn          []byte   `json:"sn"`
	MessageType uint8    `json:"messageType"`
	Data        []byte   `json:"data"`
	Protocols   []string `json:"protocols"`
}

type CSMessage struct {
	MsgType []byte `json:"msgType"`
	Payload []byte `json:"payload"`
}

type Block struct {
	Hash   string        `json:"hash"`
	Height int           `json:"height"`
	Tx     []Transaction `json:"tx"`
}

// Transaction represents the structure of a Bitcoin transaction.
type Transaction struct {
	TxID string `json:"txid"`
	Vout []Vout `json:"vout"`
}

// Vout represents the transaction output.
type Vout struct {
	Value        float64      `json:"value"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

// ScriptPubKey represents the script public key.
type ScriptPubKey struct {
	Hex       string   `json:"hex"`
	Addresses []string `json:"addresses"`
	Type      string   `json:"type"`
	Asm       string   `json:"asm"`
}

// BlockchainInfo represents the information returned by getblockchaininfo.
type BlockchainInfo struct {
	Chain                string  `json:"chain"`
	Blocks               int     `json:"blocks"`
	Headers              int     `json:"headers"`
	BestBlockHash        string  `json:"bestblockhash"`
	Difficulty           float64 `json:"difficulty"`
	Mediantime           int64   `json:"mediantime"`
	VerificationProgress float64 `json:"verificationprogress"`
	InitialBlockDownload bool    `json:"initialblockdownload"`
	Chainwork            string  `json:"chainwork"`
	SizeOnDisk           int64   `json:"size_on_disk"`
	Pruned               bool    `json:"pruned"`
	Warnings             string  `json:"warnings"`
}

// QuicknodeRequest represents a JSON-RPC request payload.
type QuicknodeRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}
