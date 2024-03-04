package types

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/types"
	relayerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

type TxSearchParam struct {
	BlockHeight uint64
	Events      []types.Event
	Prove       bool
	Page        *int
	PerPage     *int
	OrderBy     string
}

func (param *TxSearchParam) BuildQuery() string {
	heightQuery := &Query{
		Field: "tx.height", Value: param.BlockHeight,
	}

	var attribQueries []QueryExpression

	for _, event := range param.Events {
		for _, attrib := range event.Attributes {
			field := fmt.Sprintf("%s.%s", event.Type, attrib.Key)
			attribQueries = append(attribQueries, &Query{Field: field, Value: attrib.Value})
		}
	}

	eventQuery := &CompositeQuery{Or: false, Queries: attribQueries}

	finalQuery := &CompositeQuery{
		Or:      false,
		Queries: []QueryExpression{heightQuery, eventQuery},
	}

	return finalQuery.GetQuery()
}

type TxResultWaitResponse struct {
	Height int64 `json:"height"`
	Result struct {
		Code      int    `json:"code"`
		Codespace string `json:"codespace"`
		Data      []byte `json:"data"`
		Log       string `json:"log"`
	} `json:"result"`
}

type TxResultChan struct {
	TxResult *relayerTypes.TxResponse
	Error    error
}

// HexBytes
type HexBytes []byte

// NewHexBytes returns a new HexBytes
func NewHexBytes(bz []byte) HexBytes {
	return HexBytes(bz)
}

// MarshalJSON marshals the HexBytes to JSON
func (bz HexBytes) MarshalJSON() ([]byte, error) {
	s := strings.ToUpper(hex.EncodeToString(bz))
	jbz := make([]byte, len(s)+2)
	jbz[0] = '"'
	copy(jbz[1:], s)
	jbz[len(jbz)-1] = '"'
	return jbz, nil
}
