package sorobanclient

type LedgerEventResponse struct {
	Events       []LedgerEvents `json:"events"`
	LatestLedger uint64         `json:"latestLedger"`
}
type LedgerEvents struct {
	Type                     string   `json:"type"`
	Ledger                   int64    `json:"ledger"`
	ContractID               string   `json:"contractId"`
	ID                       string   `json:"id"`
	PagingToken              string   `json:"pagingToken"`
	Topic                    []string `json:"topic"`
	Value                    string   `json:"value"`
	InSuccessfulContractCall bool     `json:"inSuccessfulContractCall"`
	TxHash                   string   `json:"txHash"`
}

type LatestLedgerResponse struct {
	ID              string `json:"id"`
	ProtocolVersion uint64 `json:"protocolVersion"`
	Sequence        uint64 `json:"sequence"`
}

type CallResult struct {
	Xdr  string   `json:"xdr"`
	Auth []string `json:"auth"`
}

type TxSimulationResult struct {
	LatestLedger       uint64           `json:"latestLedger"`
	Results            []CallResult     `json:"results"`
	MinResourceFee     string           `json:"minResourceFee"`  // Recommended minimum resource fee to add when submitting the transaction. This fee is to be added on top of the Stellar network fee.
	TransactionDataXDR string           `json:"transactionData"` // The recommended Soroban Transaction Data to use when submitting the simulated transaction. This data contains the refundable fee and resource usage information such as the ledger footprint and IO access data (serialized in a base64 string).
	RestorePreamble    *RestorePreamble `json:"restorePreamble,omitempty"`
}

type RestorePreamble struct {
	TransactionData string `json:"transactionData"` // SorobanTransactionData XDR in base64
	MinResourceFee  int64  `json:"minResourceFee,string"`
}

type ResourceConfig struct {
	InstructionLeeway uint64 `json:"instructionLeeway"`
}

type TransactionResponse struct {
	Status                string `json:"status"`
	LatestLedger          int64  `json:"latestLedger"`
	LatestLedgerCloseTime string `json:"latestLedgerCloseTime"`
	OldestLedger          int64  `json:"oldestLedger"`
	OldestLedgerCloseTime string `json:"oldestLedgerCloseTime"`
	ApplicationOrder      int64  `json:"applicationOrder"`
	EnvelopeXdr           string `json:"envelopeXdr"`
	ResultXdr             string `json:"resultXdr"`
	ResultMetaXdr         string `json:"resultMetaXdr"`
	Ledger                int64  `json:"ledger"`
	CreatedAt             string `json:"createdAt"`
}
