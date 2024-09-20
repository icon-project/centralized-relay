package stellar

import (
	jsoniter "github.com/json-iterator/go"
)

type StellarMessage struct {
	Params          []string
	Method          string
	PackageObjectId string
	Module          string
}

func (m *StellarMessage) Type() string {
	return m.Method
}

func (m *StellarMessage) MsgBytes() ([]byte, error) {
	return jsoniter.Marshal(m.Params)
}
