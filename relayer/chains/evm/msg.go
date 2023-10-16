package evm

import (
	"encoding/json"
)

const defaultStepLimit = 13610920010

type Message struct {
	Params interface{}
	Method string
}

func (m *Message) Type() string {
	return m.Method
}

func (m *Message) MsgBytes() ([]byte, error) {
	return json.Marshal(m.Params)
}

func (p *EVMProvider) NewMessage(msg interface{}, method string) *Message {
	return &Message{Params: msg, Method: method}
}
