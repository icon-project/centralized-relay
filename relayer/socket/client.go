package socket

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/icon-project/centralized-relay/relayer/store"
)

const (
	EventGetBlock       Event = "GetBlock"
	EventGetMessageList Event = "GetMessageList"
	EventRelayMessage   Event = "RelayMessage"
	EventPruneDB        Event = "PruneDB"
)

var (
	ErrUnknownEvent    = fmt.Errorf("unknown event")
	ErrSocketClosed    = fmt.Errorf("socket closed")
	ErrInvalidResponse = fmt.Errorf("invalid response")
)

type Client struct {
	conn net.Conn
}

func NewClient() (*Client, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, ErrSocketClosed
	}
	return &Client{conn: conn}, nil
}

// send sends message to socket
func (c *Client) send(event Event, req interface{}) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	msg := &Message{Event: event, Data: data}
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if _, err := c.conn.Write(payload); err != nil {
		return err
	}
	return nil
}

// read and parse message from socket
func (c *Client) read() (interface{}, error) {
	buf := make([]byte, 1024)
	nr, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	msg := new(Message)
	if err := json.Unmarshal(buf[:nr], msg); err != nil {
		return nil, err
	}
	return c.parseEvent(msg)
}

// parse event from message
func (c *Client) parseEvent(msg *Message) (interface{}, error) {
	switch msg.Event {
	case EventGetBlock:
		res := new([]*ResGetBlock)
		if err := json.Unmarshal(msg.Data, res); err != nil {
			return nil, err
		}
		return res, nil
	case EventGetMessageList:
		res := new(ResMessageList)
		if err := json.Unmarshal(msg.Data, res); err != nil {
			return nil, err
		}
		return res, nil
	case EventRelayMessage:
		res := new(ResRelayMessage)
		if err := json.Unmarshal(msg.Data, res); err != nil {
			return nil, err
		}
		return res, nil
	case EventPruneDB:
		res := new(ResPruneDB)
		if err := json.Unmarshal(msg.Data, res); err != nil {
			return nil, err
		}
		return res, nil
	default:
		return nil, ErrUnknownEvent
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// GetBlock sends GetBlock event to socket
func (c *Client) GetBlock(chain string, all bool) ([]*ResGetBlock, error) {
	req := &ReqGetBlock{Chain: chain, All: all}
	if err := c.send(EventGetBlock, req); err != nil {
		return nil, err
	}
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	res, ok := data.([]*ResGetBlock)
	if !ok {
		return nil, ErrInvalidResponse
	}
	return res, nil
}

// GetMessageList sends GetMessageList event to socket
func (c *Client) GetMessageList(chain string, pagination *store.Pagination) (*ResMessageList, error) {
	req := &ReqMessageList{Chain: chain, Pagination: pagination}
	if err := c.send(EventGetMessageList, req); err != nil {
		return nil, err
	}
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	res, ok := data.(*ResMessageList)
	if !ok {
		return nil, ErrInvalidResponse
	}
	return res, nil
}

// RelayMessage sends RelayMessage event to socket
func (c *Client) RelayMessage(chain string, sn uint64, height uint64) (*ResRelayMessage, error) {
	req := &ReqRelayMessage{Chain: chain, Sn: sn, Height: height}
	if err := c.send(EventRelayMessage, req); err != nil {
		return nil, err
	}
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	res, ok := data.(*ResRelayMessage)
	if !ok {
		return nil, ErrInvalidResponse
	}
	return res, nil
}

// PruneDB sends PruneDB event to socket
func (c *Client) PruneDB() (*ResPruneDB, error) {
	req := &ReqPruneDB{}
	if err := c.send(EventPruneDB, req); err != nil {
		return nil, err
	}
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	res, ok := data.(*ResPruneDB)
	if !ok {
		return nil, ErrInvalidResponse
	}
	return res, nil
}
