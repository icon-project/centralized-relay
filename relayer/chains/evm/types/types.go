package types

import (
	"context"
	"math/big"
	"sync"
	"time"

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

type NonceValue struct {
	LastUpdated time.Time
	Previous    *big.Int
	Current     *big.Int
}

type NonceTracker struct {
	address map[common.Address]*NonceValue
	Getter  func(context.Context, common.Address, *big.Int) (*big.Int, error)
	*sync.Mutex
}

var NonceUpdateInterval = 3 * time.Minute

type NonceTrackerI interface {
	Get(common.Address) *big.Int
	Set(common.Address, *big.Int)
}

// NewNonceTracker
func NewNonceTracker(getter func(context.Context, common.Address, *big.Int) (*big.Int, error)) NonceTrackerI {
	return &NonceTracker{
		address: make(map[common.Address]*NonceValue),
		Mutex:   &sync.Mutex{},
		Getter:  getter,
	}
}

// Ugly hack to fix the nonce issue
func (n *NonceTracker) Get(addr common.Address) *big.Int {
	n.Lock()
	defer n.Unlock()
	nonce := n.address[addr]
	if time.Since(nonce.LastUpdated) > NonceUpdateInterval {
		n, err := n.Getter(context.Background(), addr, nil)
		if err == nil {
			nonce.LastUpdated = time.Now()
			nonce.Current = n
			nonce.Previous = n.Sub(n, big.NewInt(1))
		}
	} else if nonce.Current.Cmp(nonce.Previous) != 1 {
		nonce.Current = nonce.Current.Add(nonce.Current, big.NewInt(1))
		nonce.Previous = nonce.Current
	}
	return nonce.Current
}

func (n *NonceTracker) Set(addr common.Address, nonce *big.Int) {
	n.Lock()
	defer n.Unlock()
	n.address[addr] = &NonceValue{
		Previous:    nonce.Sub(nonce, big.NewInt(1)),
		Current:     nonce,
		LastUpdated: time.Now(),
	}
}

type ErrorMessageRpc struct {
	RPCUrl string
	Error  error
}
