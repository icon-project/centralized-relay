package bitcoin

import (
    "github.com/btcsuite/btcd/wire"
)

type Action int

const (
    AddLiquidity Action = iota
    Swap
    WithdrawLiquidity
    CollectFee
    IncreaseLiquidity
    DecreaseLiquidity
)

type TxSearchParam struct {
    StartHeight, EndHeight uint64
    BitcoinScript          []byte
    OPReturnPrefix         int
}

type TxSearchRes struct {
    Tx     *wire.MsgTx
    Height uint64
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
