package solana

import (
	jsoniter "github.com/json-iterator/go"
)

type SolanaMessage struct {
	Params          []string
	Method          string
	PackageObjectId string
	Module          string
}

func (m *SolanaMessage) Type() string {
	return m.Method
}

func (m *SolanaMessage) MsgBytes() ([]byte, error) {
	return jsoniter.Marshal(m.Params)
}
