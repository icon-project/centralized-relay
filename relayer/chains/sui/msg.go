package sui

import (
	jsoniter "github.com/json-iterator/go"
)

type SuiMessage struct {
	Params          []interface{}
	Method          string
	PackageObjectId string
	Module          string
}

func (m *SuiMessage) Type() string {
	return m.Method
}

func (m *SuiMessage) MsgBytes() ([]byte, error) {
	return jsoniter.Marshal(m.Params)
}

func (p *Provider) NewSuiMessage(params []interface{}, packageId, module, method string) *SuiMessage {
	return &SuiMessage{
		Params:          params,
		PackageObjectId: packageId,
		Module:          module,
		Method:          method,
	}
}
