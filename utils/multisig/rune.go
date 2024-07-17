package multisig

import (
	"math/big"

	"github.com/btcsuite/btcd/txscript"
	"github.com/multiformats/go-varint"
)

type Rune struct {
	BlockNumber uint64
	TxIndex     uint32
}

// Varint encoder, ported from
// https://github.com/ordinals/ord/blob/1e6cb641faf3b1eb0aba501a7a2822d7a3dc8643/crates/ordinals/src/varint.rs#L3-L39
// Using big.Int since go doesn't support u128
func encodeToSlice(n *big.Int) []byte {
	var result []byte
	var oneTwentyEight = big.NewInt(128)

	for n.Cmp(oneTwentyEight) >= 0 {
		temp := new(big.Int).Mod(n, oneTwentyEight)
		tempByte := byte(temp.Uint64()) | 0x80
		result = append(result, tempByte)
		n.Div(n, oneTwentyEight)
	}
	result = append(result, byte(n.Uint64()))
	return result
}

func CreateRuneTransferScript(rune Rune, amount *big.Int, output uint64) ([]byte, error) {
	builder := txscript.NewScriptBuilder()

	builder.AddOp(txscript.OP_RETURN)
	builder.AddOp(txscript.OP_13)

	data := varint.ToUvarint(0)
	data = append(data, encodeToSlice(big.NewInt(int64(rune.BlockNumber)))...)
	data = append(data, encodeToSlice(big.NewInt(int64(rune.TxIndex)))...)
	data = append(data, encodeToSlice(amount)...)
	data = append(data, encodeToSlice(big.NewInt(int64(output)))...)

	return builder.AddData(data).Script()
}