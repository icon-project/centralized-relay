package icon

import (
	"encoding/json"
)

const defaultStepLimit = 13610920010

type IconMessage struct {
	Params interface{}
	Method string
}

func (m *IconMessage) Type() string {
	return m.Method
}

func (m *IconMessage) MsgBytes() ([]byte, error) {
	return json.Marshal(m.Params)
}

func (p *IconProvider) NewIconMessage(msg interface{}, method string) *IconMessage {
	return &IconMessage{Params: msg, Method: method}
}
