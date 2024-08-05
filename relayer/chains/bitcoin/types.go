package bitcoin

import "github.com/btcsuite/btcd/wire"

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
