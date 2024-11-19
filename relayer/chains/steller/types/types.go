package types

import "github.com/stellar/go/xdr"

const (
	ChainType          = "stellar"
	LedgerSeqBatchSize = 50 // the number of ledger sequences to query concurrently for listener
)

type EventFilter struct {
	LedgerSeq   uint64
	ContractIds []string
	Topics      []string
}

type Event struct {
	*xdr.ContractEvent
	LedgerSeq uint64
}

type LedgerSeqBatch struct {
	FromSeq uint64
	ToSeq   uint64
}

type GetEventFilter struct {
	StartLedger uint64     `json:"startLedger"`
	Pagination  Pagination `json:"pagination"`
	Filters     []Filter   `json:"filters"`
}

type Filter struct {
	Type        string     `json:"type,omitempty"`
	ContractIDS []string   `json:"contractIds"`
	Topics      [][]string `json:"topics,omitempty"`
}

type Pagination struct {
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}
