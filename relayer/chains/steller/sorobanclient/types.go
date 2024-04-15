package sorobanclient

type LatestLedgerResponse struct {
	ID              string `json:"id"`
	ProtocolVersion uint64 `json:"protocolVersion"`
	Sequence        uint64 `json:"sequence"`
}
