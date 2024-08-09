package socket

import (
	"math/big"
	"net"

	"github.com/icon-project/centralized-relay/relayer"
	"github.com/icon-project/centralized-relay/relayer/store"
	"github.com/icon-project/centralized-relay/relayer/types"
)

type Event string

type Message struct {
	Event Event
	Data  []byte
}

type Server struct {
	listener net.Listener
	rly      *relayer.Relayer
}

type ReqMessageList struct {
	Chain      string
	Pagination *store.Pagination
}

type ReqGetBlock struct {
	Chain string
	All   bool
}

type ReqRelayMessage struct {
	Chain  string
	Sn     *big.Int
	Height uint64
}

type ReqMessageRemove struct {
	Chain string
	Sn    *big.Int
}

type ResMessageRemove struct {
	Sn     *big.Int
	Chain  string
	Dst    string
	Height uint64
	Event  string
}

type ResMessageList struct {
	Messages []*types.RouteMessage
	Total    int
}

type ResGetBlock struct {
	Chain  string
	Height uint64
}

type ResRelayMessage struct {
	*types.RouteMessage
}

type ReqPruneDB struct {
	Chain string
}

type ResPruneDB struct {
	Status string
}

type ErrResponse struct {
	Error string
}

type ReqRevertMessage struct {
	Chain string
	Sn    uint64
}

type ResRevertMessage struct {
	Sn uint64
}

type ReqGetFee struct {
	Chain    string
	Network  string
	Response bool
}

type ResGetFee struct {
	Chain    string
	Fee      uint64
	Response bool
}

// ReqSetFee sends SetFee event to socket
type ReqSetFee struct {
	Chain   string
	Network string
	MsgFee  *big.Int
	ResFee  *big.Int
}

// ResSetFee sends SetFee event to socket
type ResSetFee struct {
	Status string
}

// ReqClaimFee sends ClaimFee event to socket
type ReqClaimFee struct {
	Chain string
}

// ResClaimFee sends ClaimFee event to socket
type ResClaimFee struct {
	Status string
}
