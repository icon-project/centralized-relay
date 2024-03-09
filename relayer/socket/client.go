package socket

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net"

	"github.com/icon-project/centralized-relay/relayer/store"
)

const (
	EventGetBlock       Event = "GetBlock"
	EventGetMessageList Event = "GetMessageList"
	EventRelayMessage   Event = "RelayMessage"
	EventMessageRemove  Event = "MessageRemove"
	EventPruneDB        Event = "PruneDB"
	EventRevertMessage  Event = "RevertMessage"
	EventError          Event = "Error"
	EventGetFee         Event = "GetFee"
	EventSetFee         Event = "SetFee"
	EventClaimFee       Event = "ClaimFee"
)

var (
	ErrUnknownEvent    = fmt.Errorf("unknown event")
	ErrSocketClosed    = fmt.Errorf("socket closed")
	ErrInvalidResponse = func(err error) error {
		return fmt.Errorf("invalid response: %v", err)
	}
	ErrUnknown = fmt.Errorf("unknown error")
)

type Client struct {
	conn net.Conn
}

func NewClient() (*Client, error) {
	conn, err := net.Dial(network, SocketPath)
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
	buf := make([]byte, 1024*10)
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
		var res []*ResGetBlock
		if err := json.Unmarshal(msg.Data, &res); err != nil {
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
	case EventMessageRemove:
		res := new(ResMessageRemove)
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
	case EventError:
		return nil, ErrUnknown
	default:
		return nil, ErrUnknownEvent
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// GetBlock sends GetBlock event to socket
func (c *Client) GetBlock(chain string) ([]*ResGetBlock, error) {
	req := &ReqGetBlock{Chain: chain, All: chain == ""}
	if err := c.send(EventGetBlock, req); err != nil {
		return nil, err
	}
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	res, ok := data.([]*ResGetBlock)
	if !ok {
		return nil, ErrInvalidResponse(err)
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
		return nil, ErrInvalidResponse(err)
	}
	return res, nil
}

// RelayMessage sends RelayMessage event to socket
func (c *Client) RelayMessage(chain string, height, sn uint64) (*ResRelayMessage, error) {
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
		return nil, ErrInvalidResponse(err)
	}
	return res, nil
}

// MessageRemove sends MessageRemove event to socket
func (c *Client) MessageRemove(chain string, sn uint64) (*ResMessageRemove, error) {
	req := &ReqMessageRemove{Chain: chain, Sn: sn}
	if err := c.send(EventMessageRemove, req); err != nil {
		return nil, err
	}
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	res, ok := data.(*ResMessageRemove)
	if !ok {
		return nil, ErrInvalidResponse(err)
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
		return nil, ErrInvalidResponse(err)
	}
	return res, nil
}

// RevertMessage sends RevertMessage event to socket
func (c *Client) RevertMessage(chain string, sn uint64) (*ResRevertMessage, error) {
	req := &ReqRevertMessage{Chain: chain, Sn: sn}
	if err := c.send(EventRevertMessage, req); err != nil {
		return nil, err
	}
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	res, ok := data.(*ResRevertMessage)
	if !ok {
		return nil, ErrInvalidResponse(err)
	}
	return res, nil
}

// GetFee sends GetFee event to socket
func (c *Client) GetFee(chain string, isReponse bool) (*ResGetFee, error) {
	req := &ReqGetFee{Chain: chain, Response: isReponse}
	if err := c.send(EventGetFee, req); err != nil {
		return nil, err
	}
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	res, ok := data.(*ResGetFee)
	if !ok {
		return nil, ErrInvalidResponse(err)
	}
	return res, nil
}

// SetFee sends SetFee event to socket
func (c *Client) SetFee(chain string, msgFee, resFee *big.Int) error {
	req := &ReqSetFee{Chain: chain, MsgFee: msgFee, ResFee: resFee}
	if err := c.send(EventSetFee, req); err != nil {
		return err
	}
	return nil
}

// ClaimFee sends ClaimFee event to socket
func (c *Client) ClaimFee(chain string) error {
	req := &ReqClaimFee{Chain: chain}
	if err := c.send(EventClaimFee, req); err != nil {
		return err
	}
	return nil
}
