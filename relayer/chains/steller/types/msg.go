package types

import (
	"math/big"

	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stellar/go/xdr"
)

type StellerMsg struct {
	relayertypes.Message
}

// Convert big.Int to Uint128
func BigIntToUint128(b *big.Int) xdr.UInt128Parts {
	lowMask := new(big.Int).SetUint64(0xFFFFFFFFFFFFFFFF)
	low := new(big.Int).And(b, lowMask).Uint64()
	high := new(big.Int).Rsh(b, 64).Uint64()
	return xdr.UInt128Parts{
		Hi: xdr.Uint64(high),
		Lo: xdr.Uint64(low),
	}
}

func Uint128ToBigInt(u xdr.UInt128Parts) *big.Int {
	high := new(big.Int).SetUint64(uint64(u.Hi))
	high.Lsh(high, 64)
	low := new(big.Int).SetUint64(uint64(u.Lo))
	high.Add(high, low)
	return high
}

func (m StellerMsg) ScvSn() xdr.ScVal {
	scVal, _ := xdr.NewScVal(xdr.ScValTypeScvU128, BigIntToUint128(m.Sn))
	return scVal
}

func (m StellerMsg) ScvMessageHeight() xdr.ScVal {
	scVal, _ := xdr.NewScVal(xdr.ScValTypeScvU64, xdr.Uint64(m.MessageHeight))
	return scVal
}

func (m StellerMsg) ScvReqID() xdr.ScVal {
	scVal, _ := xdr.NewScVal(xdr.ScValTypeScvU128, BigIntToUint128(m.ReqID))
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

func (m StellerMsg) ScvSignatures() xdr.ScVal {
	var scSignatures xdr.ScVec
	for _, sign := range m.Signatures {
		scVal, _ := xdr.NewScVal(xdr.ScValTypeScvBytes, xdr.ScBytes(sign))
		scSignatures = append(scSignatures, scVal)
	}

	scValSignatures, _ := xdr.NewScVal(xdr.ScValTypeScvVec, &scSignatures)
	return scValSignatures
}
