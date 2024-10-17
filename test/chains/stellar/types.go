package stellar

type LatestLedgerResponse struct {
	ID              string `json:"id"`
	ProtocolVersion uint64 `json:"protocolVersion"`
	Sequence        uint64 `json:"sequence"`
}

type EventResponseEvent struct {
	Type         string   `json:"type"`
	ContractId   string   `json:"contractID"`
	Topic        []string `json:"topic"`
	Value        string   `json:"value"`
	ValueDecoded map[string]interface{}
}

type EventResponse struct {
	Events []EventResponseEvent `json:"events"`
}

type EventQueryFilter struct {
	StartLedger uint64     `json:"startLedger"`
	Pagination  Pagination `json:"pagination"`
	Filters     []Filter   `json:"filters"`
}

type Pagination struct {
	Limit uint64 `json:"limit"`
}

type Filter struct {
	Type        string   `json:"type"`
	ContractIDS []string `json:"contractIds"`
}
