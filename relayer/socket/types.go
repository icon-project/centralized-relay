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
	Event Event  `json:"event"`
	Data  []byte `json:"data"`
}

type Server struct {
	listener net.Listener
	rly      *relayer.Relayer
}

type ReqMessageList struct {
	Chain      string            `json:"chain"`
	Pagination *store.Pagination `json:"pagination"`
}

type ReqGetBlock struct {
	Chain string `json:"chain"`
	All   bool   `json:"all"`
}

type ReqRelayMessage struct {
	Chain  string   `json:"chain"`
	Sn     *big.Int `json:"sn"`
	Height uint64   `json:"height"`
}

type ReqRelayRangeMessage struct {
	Chain      string `json:"chain"`
	FromHeight uint64 `json:"from_height"`
	ToHeight   uint64 `json:"to_height"`
}

type ResRelayRangeMessage struct {
	Chain    string                `json:"chain"`
	Messages []*types.RouteMessage `json:"messages"`
}

type ReqMessageRemove struct {
	Chain string   `json:"chain"`
	Sn    *big.Int `json:"sn"`
}

type ResMessageRemove struct {
	Sn     *big.Int `json:"sn"`
	Chain  string   `json:"chain"`
	Dst    string   `json:"dst"`
	Height uint64   `json:"height"`
	Event  string   `json:"event"`
}

type ResMessageList struct {
	Messages []*types.RouteMessage `json:"messages"`
	Total    int                   `json:"total"`
}

type ResGetBlock struct {
	Chain  string `json:"chain"`
	Height uint64 `json:"height"`
}

type ResRelayMessage struct {
	*types.RouteMessage
}

type ReqPruneDB struct {
	Chain string `json:"chain"`
}

type ResPruneDB struct {
	Status string `json:"status"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

type ReqRevertMessage struct {
	Chain string `json:"chain"`
	Sn    uint64 `json:"sn"`
}

type ResRevertMessage struct {
	Sn uint64 `json:"sn"`
}

type ReqGetFee struct {
	Chain    string `json:"chain"`
	Network  string `json:"network"`
	Response bool   `json:"response"`
}

type ResGetFee struct {
	Chain    string `json:"chain"`
	Fee      uint64 `json:"fee"`
	Response bool   `json:"response"`
}

// ReqSetFee sends SetFee event to socket
type ReqSetFee struct {
	Chain   string   `json:"chain"`
	Network string   `json:"network"`
	MsgFee  *big.Int `json:"msg_fee"`
	ResFee  *big.Int `json:"res_fee"`
}

// ResSetFee sends SetFee event to socket
type ResSetFee struct {
	Status string `json:"status"`
}

// ReqClaimFee sends ClaimFee event to socket
type ReqClaimFee struct {
	Chain string `json:"chain"`
}

// ResClaimFee sends ClaimFee event to socket
type ResClaimFee struct {
	Status string `json:"status"`
}

type ReqChainHeight struct {
	Chain string `json:"chain"`
}

type ResChainHeight struct {
	Chain  string `json:"chain"`
	Height uint64 `json:"height"`
}

type ReqProcessedBlock struct {
	Chain string `json:"chain"`
}

type ResProcessedBlock struct {
	Chain  string `json:"chain"`
	Height uint64 `json:"height"`
}

type ReqRangeBlockQuery struct {
	Chain      string `json:"chain"`
	FromHeight uint64 `json:"from_height"`
	ToHeight   uint64 `json:"to_height"`
}

type ResRangeBlockQuery struct {
	Chain string           `json:"chain"`
	Msgs  []*types.Message `json:"messages"`
}
