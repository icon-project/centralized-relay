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

// Todo need to remove: used just for testing purpose only
type NewMessage struct {
	Sn  uint64
	Src string
	Dst string
}

func (m *NewMessage) Convert(val xdr.ScVal) error {
	for _, entry := range **val.Map {
		switch entry.Key.String() {
		case "sn":
			snU32, ok := entry.Val.GetU32()
			if ok {
				m.Sn = uint64(snU32)
			} else {
				return fmt.Errorf("%s: got unexpected sn type: %s", ScValConversionErr, entry.Val.Type)
			}
		case "src":
			srcStr, ok := entry.Val.GetStr()
			if ok {
				m.Src = string(srcStr)
			} else {
				return fmt.Errorf("%s: got unexpected src type: %s", ScValConversionErr, entry.Val.Type)
			}
		case "dst":
			dstStr, ok := entry.Val.GetStr()
			if ok {
				m.Dst = string(dstStr)
			} else {
				return fmt.Errorf("%s: got unexpected dst type: %s", ScValConversionErr, entry.Val.Type)
			}
		}
	}
	return nil
}
