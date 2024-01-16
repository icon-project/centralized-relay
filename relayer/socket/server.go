package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
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

func NewSocket(rly *relayer.Relayer) (*Server, error) {
	l, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}
	return &Server{listener: l, rly: rly}, nil
}

// Listen to socket
func (s *Server) Listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.server(conn)
	}
}

// Send sends message to socket
func (s *Server) server(c net.Conn) {
	for {
		buf := make([]byte, 1024*2)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}
		message, err := s.parse(buf[:nr])
		if err != nil {
			message = makeError(err)
		}
		if err := s.send(c, message); err != nil {
			return
		}
	}
}

// Parse message from socket
func (s *Server) parse(data []byte) ([]byte, error) {
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

// makeError for the client to write to socket
func makeError(err error) []byte {
	message := &Message{EventError, []byte(fmt.Sprintf(`{"message": "%s"}`, err.Error()))}
	data, err := json.Marshal(message)
	if err != nil {
		return []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}
	return data
}

// Send message to socket
func (s *Server) send(conn net.Conn, data []byte) error {
	_, err := conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// parseEvent for the client to write to socket
func (s *Server) parseEvent(msg *Message) (*Message, error) {
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
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, &ResGetBlock{req.Chain, height})
		data, err := json.Marshal(blocks)
		if err != nil {
			return nil, err
		}
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
	case EventMessageRemove:
		req := new(ReqMessageRemove)
		if err := json.Unmarshal(msg.Data, req); err != nil {
			return nil, err
		}
		store := s.rly.GetMessageStore()
		key := types.MessageKey{Src: req.Chain, Sn: req.Sn}
		message, err := store.GetMessage(key)
		if err != nil {
			return nil, err
		}
		if err := store.DeleteMessage(key); err != nil {
			return nil, err
		}
		data, err := json.Marshal(&ResMessageRemove{req.Sn, req.Chain, message.Dst, message.MessageHeight, message.EventType})
		if err != nil {
			return nil, err
		}
		return &Message{EventMessageRemove, data}, nil
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
			// TODO: Find message by height and sn
			return nil, fmt.Errorf("not implemented")
		}

		store := s.rly.GetMessageStore()
		key := types.MessageKey{Src: req.Chain, Sn: req.Sn}
		message, err := store.GetMessage(key)
		if err != nil {
			return nil, err
		}
		src.MessageCache.Add(message)
		data, err := json.Marshal(&ResRelayMessage{message})
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
	case EventRevertMessage:
		req := new(ReqRevertMessage)
		if err := json.Unmarshal(msg.Data, req); err != nil {
			return nil, err
		}
		chain, err := s.rly.FindChainRuntime(req.Chain)
		if err != nil {
			return nil, err
		}
		if err := chain.Provider.RevertMessage(context.Background(), big.NewInt(0).SetUint64(req.Sn)); err != nil {
			return nil, err
		}
		data, err := json.Marshal(&ResRevertMessage{req.Sn})
		if err != nil {
			return nil, err
		}
		return &Message{EventRevertMessage, data}, nil
	default:
		return nil, fmt.Errorf("invalid request")
	}
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) IsClosed() bool {
	return s.listener == nil
}
