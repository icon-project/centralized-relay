package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockNotification struct {
	Hash   common.Hash
	Height *big.Int
	Header *types.Header
	Logs   []types.Log
}
