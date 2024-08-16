package socket

import (
	"fmt"
	"math/big"
	"net"

	jsoniter "github.com/json-iterator/go"

	"github.com/icon-project/centralized-relay/relayer/store"
)

const (
	EventGetBlock                Event = "GetBlock"
	EventGetMessageList          Event = "GetMessageList"
	EventRelayMessage            Event = "RelayMessage"
	EventRelayRangeMessage       Event = "RelayRangeMessage"
	EventMessageRemove           Event = "MessageRemove"
	EventPruneDB                 Event = "PruneDB"
	EventRevertMessage           Event = "RevertMessage"
	EventError                   Event = "Error"
	EventGetFee                  Event = "GetFee"
	EventSetFee                  Event = "SetFee"
	EventClaimFee                Event = "ClaimFee"
	EventGetLatestHeight         Event = "GetLatestHeight"
	EventGetLatestProcessedBlock Event = "GetLatestProcessedBlock"
	EventGetBlockRange           Event = "GetBlockRange"
	EventGetConfig               Event = "GetConfig"
	EventListChainInfo           Event = "ListChainInfo"
	EventGetBalance              Event = "GetChainBalance"
	EventRelayInfo               Event = "RelayInfo"
	EventMessageReceived         Event = "MessageReceived"
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
func (c *Client) send(req interface{}) error {
	data, err := jsoniter.Marshal(req)
	if err != nil {
		return err
	}
	if _, err := c.conn.Write(data); err != nil {
		return err
	}
	return nil
}

// read and parse message from socket
func (c *Client) read() (interface{}, error) {
	buf := make([]byte, 1024*100)
	nr, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	msg := new(Response)
	return msg, jsoniter.Unmarshal(buf[:nr], msg)
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// GetBlock sends GetBlock event to socket
func (c *Client) GetBlock(chain string) ([]*ResGetBlock, error) {
	req := &ReqGetBlock{Chain: chain, All: chain == ""}
	if err := c.send(&Request{Event: EventGetBlock, Data: req}); err != nil {
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
	if err := c.send(&Request{Event: EventGetMessageList, Data: req}); err != nil {
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
func (c *Client) RelayMessage(chain string, height uint64, sn *big.Int) (*ResRelayMessage, error) {
	req := &ReqRelayMessage{Chain: chain, FromHeight: height, ToHeight: height}
	if err := c.send(&Request{Event: EventRelayMessage, Data: req}); err != nil {
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

func (c *Client) RelayRangeMessage(chain string, fromHeight, toHeight uint64) (*ResRelayRangeMessage, error) {
	req := &ReqRelayRangeMessage{Chain: chain, FromHeight: fromHeight, ToHeight: toHeight}
	if err := c.send(&Request{Event: EventRelayRangeMessage, Data: req}); err != nil {
		return nil, err
	}
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	res, ok := data.(*ResRelayRangeMessage)
	if !ok {
		return &ResRelayRangeMessage{}, nil
	}
	return res, nil
}

// MessageRemove sends MessageRemove event to socket
func (c *Client) MessageRemove(chain string, sn *big.Int) (*ResMessageRemove, error) {
	req := &ReqMessageRemove{Chain: chain, Sn: sn}
	if err := c.send(&Request{Event: EventMessageRemove, Data: req}); err != nil {
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
	if err := c.send(&Request{Event: EventPruneDB, Data: req}); err != nil {
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
	if err := c.send(&Request{Event: EventRevertMessage, Data: req}); err != nil {
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
func (c *Client) GetFee(chain string, network string, isReponse bool) (*ResGetFee, error) {
	req := &ReqGetFee{Chain: chain, Network: network, Response: isReponse}
	if err := c.send(&Request{Event: EventGetFee, Data: req}); err != nil {
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
func (c *Client) SetFee(chain, network string, msgFee, resFee *big.Int) (*ResSetFee, error) {
	req := &ReqSetFee{Chain: chain, Network: network, MsgFee: msgFee, ResFee: resFee}
	if err := c.send(&Request{Event: EventSetFee, Data: req}); err != nil {
		return nil, err
	}
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	res, ok := data.(*ResSetFee)
	if !ok {
		return nil, ErrInvalidResponse(err)
	}
	return res, nil
}

// ClaimFee sends ClaimFee event to socket
func (c *Client) ClaimFee(chain string) (*ResClaimFee, error) {
	req := &ReqClaimFee{Chain: chain}
	if err := c.send(&Request{Event: EventClaimFee, Data: req}); err != nil {
		return nil, err
	}
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	res, ok := data.(*ResClaimFee)
	if !ok {
		return nil, ErrInvalidResponse(err)
	}
	return res, nil
}

func (c *Client) GetLatestHeight(chain string) (*ResChainHeight, error) {
	req := &ReqChainHeight{Chain: chain}
	if err := c.send(&Request{Event: EventGetLatestHeight, Data: req}); err != nil {
		return nil, err
	}
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	res, ok := data.(*ResChainHeight)
	if !ok {
		return nil, ErrInvalidResponse(err)
	}
	return res, nil
}

func (c *Client) GetLatestProcessedBlock(chain string) (*ResProcessedBlock, error) {
	req := &ReqGetBlock{Chain: chain}
	if err := c.send(&Request{Event: EventGetLatestProcessedBlock, Data: req}); err != nil {
		return nil, err
	}
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	res, ok := data.(*ResProcessedBlock)
	if !ok {
		return nil, ErrInvalidResponse(err)
	}
	return res, nil
}

func (c *Client) QueryBlockRange(chain string, fromHeight, toHeight uint64) (*ResRangeBlockQuery, error) {
	req := &ReqRangeBlockQuery{Chain: chain, FromHeight: fromHeight, ToHeight: toHeight}
	if err := c.send(&Request{Event: EventGetBlockRange, Data: req}); err != nil {
		return nil, err
	}
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	res, ok := data.(*ResRangeBlockQuery)
	if !ok {
		return nil, ErrInvalidResponse(err)
	}
	return res, nil
}
