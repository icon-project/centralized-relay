package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	relayerTypes "github.com/icon-project/centralized-relay/relayer/types"
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

type TxResultWaitResponse struct {
	Height int64 `json:"height"`
	Result struct {
		Code      int    `json:"code"`
		Codespace string `json:"codespace"`
		Data      []byte `json:"data"`
	} `json:"result"`
}

type TxResultChan struct {
	TxResult *relayerTypes.TxResponse
	Error    error
}
