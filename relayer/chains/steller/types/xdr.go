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
