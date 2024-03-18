package icon

import (
	"encoding/json"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
)

const defaultStepLimit = 13610920010

type IconMessage struct {
	Address types.Address
	Params  interface{}
	Method  string
}

func (m *IconMessage) Type() string {
	return m.Method
}

func (m *IconMessage) MsgBytes() ([]byte, error) {
	return json.Marshal(m.Params)
}

func (p *Provider) NewIconMessage(address types.Address, msg interface{}, method string) *IconMessage {
	return &IconMessage{Address: address, Params: msg, Method: method}
}
