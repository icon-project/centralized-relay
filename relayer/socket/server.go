package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"

	"github.com/icon-project/centralized-relay/relayer"
	"github.com/icon-project/centralized-relay/relayer/types"
)

var (
	addr    = path.Join(os.TempDir(), "relayer.sock")
	network = "unix"
)

func NewSocket(rly *relayer.Relayer) (*dbServer, error) {
	l, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}
	return &dbServer{listener: l, rly: rly}, nil
}

// Listen to socket
func (s *dbServer) Listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.server(conn)
	}
}

// Send sends message to socket
func (s *dbServer) server(c net.Conn) {
	for {
		buf := make([]byte, 1024*2)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}
		message, err := s.parse(buf[:nr])
		if err != nil {
			return
		}
		if err := s.send(c, message); err != nil {
			return
		}
	}
}

// Parse message from socket
func (s *dbServer) parse(data []byte) ([]byte, error) {
	msg := new(Message)
	if err := json.Unmarshal(data, msg); err != nil {
		return nil, err
	}
	payload, err := s.parseEvent(msg)
	if err != nil {
		return nil, err
	}
	return json.Marshal(payload)
}

// Send message to socket
func (s *dbServer) send(conn net.Conn, data []byte) error {
	_, err := conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// parseEvent for the client to write to socket
func (s *dbServer) parseEvent(msg *Message) (*Message, error) {
	switch msg.Event {
	case EventGetBlock:
		req := new(ReqGetBlock)
		if err := json.Unmarshal(msg.Data, req); err != nil {
			return nil, err
		}
		var blocks []*ResGetBlock

		if req.All {
			for _, chain := range s.rly.GetAllChainsRuntime() {
				blocks = append(blocks, &ResGetBlock{chain.Provider.NID(), chain.LastSavedHeight})
			}
			data, err := json.Marshal(blocks)
			if err != nil {
				return nil, err
			}
			return &Message{EventGetBlock, data}, nil
		}

		store := s.rly.GetBlockStore()
		height, err := store.GetLastStoredBlock(req.Chain)
		fmt.Println("height", err)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, &ResGetBlock{req.Chain, height})
		data, err := json.Marshal(blocks)
		if err != nil {
			return nil, err
		}
		fmt.Println("blocks", blocks)
		return &Message{EventGetBlock, data}, nil
	case EventGetMessageList:
		req := new(ReqMessageList)
		if err := json.Unmarshal(msg.Data, req); err != nil {
			return nil, err
		}
		store := s.rly.GetMessageStore()
		messages, err := store.GetMessages(req.Chain, req.Pagination)
		if err != nil {
			return nil, err
		}
		total, err := store.TotalCountByChain(req.Chain)
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(&ResMessageList{messages, int(total)})
		if err != nil {
			return nil, err
		}
		return &Message{EventGetMessageList, data}, nil
	case EventRelayMessage:
		req := new(ReqRelayMessage)
		if err := json.Unmarshal(msg.Data, req); err != nil {
			return nil, err
		}

		src, err := s.rly.FindChainRuntime(req.Chain)
		if err != nil {
			return nil, err
		}

		if req.Height != 0 {
			// TODO: Find message by height
		}

		store := s.rly.GetMessageStore()
		key := types.MessageKey{Src: req.Chain, Sn: req.Sn}
		message, err := store.GetMessage(key)
		if err != nil {
			return nil, err
		}
		dst, err := s.rly.FindChainRuntime(message.Dst)
		if err != nil {
			return nil, err
		}
		message.SetIsProcessing(true)
		s.rly.RouteMessage(context.Background(), message, dst, src)
		data, err := json.Marshal(&ResRelayMessage{message, ""})
		if err != nil {
			return nil, err
		}
		return &Message{EventRelayMessage, data}, nil
	case EventPruneDB:
		if err := s.rly.PruneDB(); err != nil {
			return nil, err
		}
		data, err := json.Marshal(&ResPruneDB{"Success"})
		if err != nil {
			return nil, err
		}
		return &Message{EventPruneDB, data}, nil
	default:
		return nil, fmt.Errorf("invalid request")
	}
}

func (s *dbServer) Close() error {
	return s.listener.Close()
}

func (s *dbServer) IsClosed() bool {
	return false
}
