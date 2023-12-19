package socket

import (
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

type dbServer struct {
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
	Sn     uint64
	Height uint64
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
	Hash string
}

type ReqPruneDB struct {
	Chain string
}

type ResPruneDB struct {
	Status string
}
