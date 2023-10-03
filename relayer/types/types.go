package types

type BlockInfo struct {
	Height   uint64
	Messages []Message
}

type Message struct {
	Dst           string `json:"dst"`
	Src           string `json:"src"`
	Sn            uint64 `json:"sn"`
	Data          []byte `json:"data"`
	MessageHeight uint64 `json:"messageHeight"`
	EventType     string `json:"eventType"`
}

func (m Message) MessageKey() MessageKey {
	return NewMessageKey(m.Sn, m.Src, m.Dst, m.EventType)
}

type RouteMessage struct {
	Message
	Retry        uint64
	IsProcessing bool
}

func NewRouteMessage(m Message) *RouteMessage {
	return &RouteMessage{
		Message:      m,
		Retry:        0,
		IsProcessing: false,
	}
}

func (r *RouteMessage) GetMessage() Message {
	return r.Message
}

func (r *RouteMessage) IncrementRetry() {
	r.Retry += 1
}
func (r *RouteMessage) GetRetry() uint64 {
	return r.Retry
}

func (r *RouteMessage) SetIsProcessing(isProcessing bool) {
	r.IsProcessing = isProcessing
}

func (r *RouteMessage) GetIsProcessing() bool {
	return r.IsProcessing
}

type ExecuteMessageResponse struct {
	// *RouteMessage
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

type MessageKey struct {
	Sn        uint64
	SrcChain  string
	dstChain  string
	EventType string
}

func NewMessageKey(Sn uint64, SrcChain string, DstChain string, EventType string) MessageKey {
	return MessageKey{Sn, SrcChain, DstChain, EventType}
}

type MessageCache map[MessageKey]*RouteMessage

func (m MessageCache) Add(r *RouteMessage) {
	key := NewMessageKey(r.Sn, r.Src, r.Dst, r.EventType)
	m[key] = r
}

func (m MessageCache) Len() uint64 {
	return uint64(len(m))
}

func (m MessageCache) Remove(key MessageKey) {
	delete(m, key)
}
