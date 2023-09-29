package types

type BlockInfo struct {
	Height   uint64
	Messages []RelayMessage
}

type RelayMessage struct {
	Target        string
	Src           string
	Sn            uint64
	Data          []byte
	MessageHeight uint64
}

type RouteMessage struct {
	RelayMessage
	Retry uint64
}

func NewRouteMessage(m RelayMessage) *RouteMessage {
	return &RouteMessage{
		RelayMessage: m,
		Retry:        0,
	}
}

func (r *RouteMessage) IncrementRetry() {
	r.Retry += 1
}
func (r *RouteMessage) GetRetry() uint64 {
	return r.Retry
}

type ExecuteMessageResponse struct {
	RouteMessage
	TxResponse
}

type TxResponse struct {
	Height    int64
	TxHash    string
	Codespace string
	Code      ResponseCode
	Data      string
}

type ResponseCode uint8

const (
	Success ResponseCode = 0
	Failed  ResponseCode = 1
)
