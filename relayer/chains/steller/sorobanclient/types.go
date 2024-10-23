package sorobanclient

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
