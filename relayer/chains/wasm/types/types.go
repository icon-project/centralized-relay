package types

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"

	abiTypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/types"

	relayerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

type TxSearchParam struct {
	StartHeight, EndHeight uint64
	Events                 []types.Event
	Prove                  bool
	Page                   *int
	PerPage                *int
	OrderBy                string
}

func (param *TxSearchParam) BuildQuery() string {
	var queries []QueryExpression

	if param.EndHeight-param.StartHeight == 0 { // if diff is 0, then it is a single height
		queries = append(queries, &Query{Field: "tx.height", Value: param.StartHeight, Operator: QueryOperator.Eq})
	} else {
		startHeight := &Query{
			Field: "tx.height", Value: param.StartHeight,
			Operator: QueryOperator.Gte,
		}
		endHeight := &Query{
			Field: "tx.height", Value: param.EndHeight,
			Operator: QueryOperator.Lte,
		}
		queries = append(queries, startHeight, endHeight)
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
		Queries: append(queries, eventQuery),
	}

	return finalQuery.GetQuery()
}

type TxResultResponse struct {
	Height int64 `json:"height"`
	Result struct {
		Code      int              `json:"code"`
		Codespace string           `json:"codespace"`
		Data      []byte           `json:"data"`
		Log       string           `json:"log"`
		Events    []abiTypes.Event `json:"events"`
	} `json:"result"`
}

type TxResult struct {
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

// SubscribeOpts
type SubscribeOpts struct {
	Height  uint64
	Address string
	Method  string
}

// HightRange is a struct to represent a range of heights
type HeightRange struct {
	Start uint64
	End   uint64
}

func Retry(ctx context.Context, maxAttempts uint8, baseDelay, maxDelay time.Duration, logger *zap.Logger, operation func() error) error {
	var attempts uint8
	for attempts < maxAttempts {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := operation()
		if err == nil {
			return nil
		}

		attempts++
		delay := time.Duration(math.Pow(2, float64(attempts))) * baseDelay
		if delay > maxDelay {
			delay = maxDelay
		}
		if attempts >= maxAttempts {
			logger.Error("operation failed", zap.Uint8("attempt", attempts), zap.Error(err))
			return err
		}
		logger.Warn("operation failed, retrying...", zap.Uint8("attempt", attempts), zap.Duration("retrying_in", delay))
		time.Sleep(delay)
	}
	return fmt.Errorf("operation failed after %d attempts", maxAttempts)
}
