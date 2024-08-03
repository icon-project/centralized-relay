package bitcoin

type TxSearchParam struct {
	StartHeight, EndHeight uint64
	BitcoinAddress string
	OPReturnPrefix byte
}

// HightRange is a struct to represent a range of heights
type HeightRange struct {
	Start uint64
	End   uint64
}