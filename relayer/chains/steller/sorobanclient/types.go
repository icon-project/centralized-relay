package sorobanclient

type LatestLedgerResponse struct {
	ID              string `json:"id"`
	ProtocolVersion uint64 `json:"protocolVersion"`
	Sequence        uint64 `json:"sequence"`
}

type TxSimulationResult struct {
	LatestLedger       uint64 `json:"latestLedger"`
	MinResourceFee     string `json:"minResourceFee"`  // Recommended minimum resource fee to add when submitting the transaction. This fee is to be added on top of the Stellar network fee.
	TransactionDataXDR string `json:"transactionData"` // The recommended Soroban Transaction Data to use when submitting the simulated transaction. This data contains the refundable fee and resource usage information such as the ledger footprint and IO access data (serialized in a base64 string).
}
