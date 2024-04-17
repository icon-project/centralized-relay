package types

import "github.com/stellar/go/xdr"

const (
	ChainType          = "steller"
	LedgerSeqBatchSize = 50 // the number of ledger sequences to query concurrently for listener
)

type EventFilter struct {
	LedgerSeq   uint64
	ContractIds []string
	Topics      []string
}

type Event struct {
	xdr.ContractEvent
	LedgerSeq uint64
}

type LedgerSeqBatch struct {
	FromSeq uint64
	ToSeq   uint64
}
