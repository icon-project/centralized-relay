package types

import (
	"fmt"

	"github.com/stellar/go/xdr"
)

const (
	ScValConversionErr = "ScVal conversion failed"
)

// any type that needs conversion from ScVal should implement this interface
type ScValConverter interface {
	Convert(val xdr.ScVal) error
}

type ScvBool bool

func (v *ScvBool) Convert(val xdr.ScVal) error {
	b, ok := val.GetB()
	if ok {
		*v = ScvBool(b)
	} else {
		return fmt.Errorf("%s: got unexpected type: %s", ScValConversionErr, val.Type)
	}
	return nil
}

type ScvU64F128 uint64

func (v *ScvU64F128) Convert(val xdr.ScVal) error {
	parts, ok := val.GetU128()
	_ = parts
	if ok {
		*v = ScvU64F128(parts.Lo)
	} else {
		return fmt.Errorf("%s: got unexpected type: %s", ScValConversionErr, val.Type)
	}
	return nil
}
