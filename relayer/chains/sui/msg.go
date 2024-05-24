package sui

import (
	jsoniter "github.com/json-iterator/go"
)

type SuiCallArg struct {
	Val  interface{}
	Type string
}
type SuiMessage struct {
	Params    []SuiCallArg
	Method    string
	PackageId string
	Module    string
}

func (m *SuiMessage) Type() string {
	return m.Method
}

func (m *SuiMessage) MsgBytes() ([]byte, error) {
	return jsoniter.Marshal(m.Params)
}

func (p *Provider) NewSuiMessage(params []SuiCallArg, packageId, module, method string) *SuiMessage {
	return &SuiMessage{
		Params:    params,
		PackageId: packageId,
		Module:    module,
		Method:    method,
	}
}