package types

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockNotification struct {
	Hash   common.Hash
	Height *big.Int
	Header *types.Header
	Logs   []types.Log
}

type Block struct {
	Transactions []string `json:"transactions"`
	GasUsed      string   `json:"gasUsed"`
}

type NonceTracker struct {
	address map[common.Address]*big.Int
	*sync.Mutex
}

type NonceTrackerI interface {
	Get(common.Address) *big.Int
	Set(common.Address, *big.Int)
	Inc(common.Address)
}

// NewNonceTracker
func NewNonceTracker() NonceTrackerI {
	return &NonceTracker{
		address: make(map[common.Address]*big.Int),
		Mutex:   &sync.Mutex{},
	}
}

func (n *NonceTracker) Get(addr common.Address) *big.Int {
	n.Lock()
	defer n.Unlock()
	return n.address[addr]
}

func (n *NonceTracker) Set(addr common.Address, nonce *big.Int) {
	n.Lock()
	defer n.Unlock()
	n.address[addr] = nonce
}

func (n *NonceTracker) Inc(addr common.Address) {
	n.Lock()
	defer n.Unlock()
	n.address[addr] = n.address[addr].Add(n.address[addr], big.NewInt(1))
}
