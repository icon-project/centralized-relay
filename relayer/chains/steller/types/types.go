package types

import "github.com/stellar/go/xdr"

const (
	ChainType = "steller"
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
