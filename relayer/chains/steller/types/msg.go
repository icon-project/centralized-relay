package types

import (
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stellar/go/xdr"
)

type StellerMsg struct {
	relayertypes.Message
}

func (m StellerMsg) ScvSn() xdr.ScVal {
	scVal, _ := xdr.NewScVal(xdr.ScValTypeScvU64, xdr.Uint64(m.Sn))
	return scVal
}

func (m StellerMsg) ScvMessageHeight() xdr.ScVal {
	scVal, _ := xdr.NewScVal(xdr.ScValTypeScvU64, xdr.Uint64(m.MessageHeight))
	return scVal
}

func (m StellerMsg) ScvReqID() xdr.ScVal {
	scVal, _ := xdr.NewScVal(xdr.ScValTypeScvU64, xdr.Uint64(m.ReqID))
	return scVal
}

func (m StellerMsg) ScvSrc() xdr.ScVal {
	scVal, _ := xdr.NewScVal(xdr.ScValTypeScvString, xdr.ScString(m.Src))
	return scVal
}

func (m StellerMsg) ScvDst() xdr.ScVal {
	scVal, _ := xdr.NewScVal(xdr.ScValTypeScvString, xdr.ScString(m.Dst))
	return scVal
}

func (m StellerMsg) ScvData() xdr.ScVal {
	scVal, _ := xdr.NewScVal(xdr.ScValTypeScvBytes, xdr.ScBytes(m.Data))
	return scVal
}
