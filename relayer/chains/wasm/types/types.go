package types

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	"io"
)

type TxSearchParam struct {
	BlockHeight uint64
	Events      []types.Event
	Prove       bool
	Page        *int
	PerPage     *int
	OrderBy     string
}

func (param TxSearchParam) BuildQuery() string {
	heightQuery := Query{
		Field: "tx.height", Value: param.BlockHeight,
	}

	var attribQueries []QueryExpression

	for _, event := range param.Events {
		for _, attrib := range event.Attributes {
			field := fmt.Sprintf("%s.%s", event.Type, attrib.Key)
			attribQueries = append(attribQueries, Query{Field: field, Value: attrib.Value})
		}
	}

	eventQuery := CompositeQuery{
		Or: false, Queries: attribQueries,
	}

	finalQuery := CompositeQuery{
		Or:      false,
		Queries: []QueryExpression{heightQuery, eventQuery},
	}

	return finalQuery.GetQuery()
}

type KeyringPassword string

func (kp KeyringPassword) Read(p []byte) (n int, err error) {
	copy(p, kp)
	return len(kp), io.EOF
}

type HexBytes string

func (hs HexBytes) Value() ([]byte, error) {
	if hs == "" {
		return nil, nil
	}
	return hex.DecodeString(string(hs[2:]))
}

func NewHexBytes(b []byte) HexBytes {
	return HexBytes("0x" + hex.EncodeToString(b))
}

type Base64Str string

func (bs Base64Str) Decode() ([]byte, error) {
	if bs == "" {
		return nil, nil
	}
	return base64.StdEncoding.DecodeString(string(bs))
}
