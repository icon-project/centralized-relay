package socket

import (
	"net"

	"github.com/icon-project/centralized-relay/relayer/store"
)

type Event string

const (
	EventGetBlock       Event = "GetBlock"
	EventGetMessageList Event = "GetMessageList"
	EventRelayMessage   Event = "RelayMessage"
)

type MessageList struct {
	Chain      string
	Pagination *store.Pagination
}

type GetBlock struct {
	Chain string
	All   bool
}

type RelayMessage struct {
	Chain  string
	Sn     uint64
	Height uint64
}

func ConnectSocket() (net.Conn, error) {
	return net.Dial(network, unixSocketPath)
}
