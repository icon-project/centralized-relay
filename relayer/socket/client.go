package socket

import (
	"fmt"
	"math/big"
	"net"

	jsoniter "github.com/json-iterator/go"
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
	EventRelayerInfo             Event = "RelayerInfo"
	EventMessageReceived         Event = "MessageReceived"
	EventGetBlockEvents          Event = "GetBlockEvents"
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
func (c *Client) read() (*Response, error) {
	buf := make([]byte, 1024*100)
	nr, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	res := new(Response)
	if err := jsoniter.Unmarshal(buf[:nr], res); err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, fmt.Errorf(res.Message)
	}

	return res, nil
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
	res, err := c.read()
	if err != nil {
		return nil, err
	}

	resData := []*ResGetBlock{}
	if err := parseResData(res.Data, &resData); err != nil {
		return nil, err
	}

	return resData, nil
}

// GetMessageList sends GetMessageList event to socket
func (c *Client) GetMessageList(chain string, limit uint) (*ResMessageList, error) {
	req := &ReqMessageList{Chain: chain, Limit: limit}
	if err := c.send(&Request{Event: EventGetMessageList, Data: req}); err != nil {
		return nil, err
	}
	res, err := c.read()
	if err != nil {
		return nil, err
	}

	resData := new(ResMessageList)
	if err := parseResData(res.Data, &resData); err != nil {
		return nil, err
	}

	return resData, nil
}

// RelayMessage sends RelayMessage event to socket
func (c *Client) RelayMessage(chain string, height uint64, sn *big.Int) (*ResRelayMessage, error) {
	req := &ReqRelayMessage{Chain: chain, FromHeight: height, ToHeight: height}
	if err := c.send(&Request{Event: EventRelayMessage, Data: req}); err != nil {
		return nil, err
	}
	res, err := c.read()
	if err != nil {
		return nil, err
	}

	resData := new(ResRelayMessage)
	if err := parseResData(res.Data, &resData); err != nil {
		return nil, err
	}

	return resData, nil
}

// MessageRemove sends MessageRemove event to socket
func (c *Client) MessageRemove(chain string, sn *big.Int) (*ResMessageRemove, error) {
	req := &ReqMessageRemove{Chain: chain, Sn: sn}
	if err := c.send(&Request{Event: EventMessageRemove, Data: req}); err != nil {
		return nil, err
	}
	res, err := c.read()
	if err != nil {
		return nil, err
	}

	resData := new(ResMessageRemove)
	if err := parseResData(res.Data, &resData); err != nil {
		return nil, err
	}

	return resData, nil
}

// PruneDB sends PruneDB event to socket
func (c *Client) PruneDB() (*ResPruneDB, error) {
	req := &ReqPruneDB{}
	if err := c.send(&Request{Event: EventPruneDB, Data: req}); err != nil {
		return nil, err
	}
	res, err := c.read()
	if err != nil {
		return nil, err
	}

	resData := new(ResPruneDB)
	if err := parseResData(res.Data, &resData); err != nil {
		return nil, err
	}

	return resData, nil
}

// RevertMessage sends RevertMessage event to socket
func (c *Client) RevertMessage(chain string, sn uint64) (*ResRevertMessage, error) {
	req := &ReqRevertMessage{Chain: chain, Sn: sn}
	if err := c.send(&Request{Event: EventRevertMessage, Data: req}); err != nil {
		return nil, err
	}
	res, err := c.read()
	if err != nil {
		return nil, err
	}

	resData := new(ResRevertMessage)
	if err := parseResData(res.Data, &resData); err != nil {
		return nil, err
	}

	return resData, nil
}

// GetFee sends GetFee event to socket
func (c *Client) GetFee(chain string, network string, isReponse bool) (*ResGetFee, error) {
	req := &ReqGetFee{Chain: chain, Network: network, Response: isReponse}
	if err := c.send(&Request{Event: EventGetFee, Data: req}); err != nil {
		return nil, err
	}
	res, err := c.read()
	if err != nil {
		return nil, err
	}

	resData := new(ResGetFee)
	if err := parseResData(res.Data, &resData); err != nil {
		return nil, err
	}

	return resData, nil
}

// SetFee sends SetFee event to socket
func (c *Client) SetFee(chain, network string, msgFee, resFee *big.Int) (*ResSetFee, error) {
	req := &ReqSetFee{Chain: chain, Network: network, MsgFee: msgFee, ResFee: resFee}
	if err := c.send(&Request{Event: EventSetFee, Data: req}); err != nil {
		return nil, err
	}
	res, err := c.read()
	if err != nil {
		return nil, err
	}

	resData := new(ResSetFee)
	if err := parseResData(res.Data, &resData); err != nil {
		return nil, err
	}

	return resData, nil
}

// ClaimFee sends ClaimFee event to socket
func (c *Client) ClaimFee(chain string) (*ResClaimFee, error) {
	req := &ReqClaimFee{Chain: chain}
	if err := c.send(&Request{Event: EventClaimFee, Data: req}); err != nil {
		return nil, err
	}
	res, err := c.read()
	if err != nil {
		return nil, err
	}

	resData := new(ResClaimFee)
	if err := parseResData(res.Data, &resData); err != nil {
		return nil, err
	}

	return resData, nil
}

func (c *Client) GetLatestHeight(chain string) (*ResChainHeight, error) {
	req := &ReqChainHeight{Chain: chain}
	if err := c.send(&Request{Event: EventGetLatestHeight, Data: req}); err != nil {
		return nil, err
	}
	res, err := c.read()
	if err != nil {
		return nil, err
	}

	resData := new(ResChainHeight)
	if err := parseResData(res.Data, &resData); err != nil {
		return nil, err
	}

	return resData, nil
}

func (c *Client) GetLatestProcessedBlock(chain string) (*ResProcessedBlock, error) {
	req := &ReqGetBlock{Chain: chain}
	if err := c.send(&Request{Event: EventGetLatestProcessedBlock, Data: req}); err != nil {
		return nil, err
	}
	res, err := c.read()
	if err != nil {
		return nil, err
	}

	resData := new(ResProcessedBlock)
	if err := parseResData(res.Data, &resData); err != nil {
		return nil, err
	}

	return resData, nil
}

func (c *Client) QueryBlockRange(chain string, fromHeight, toHeight uint64) (*ResRangeBlockQuery, error) {
	req := &ReqRangeBlockQuery{Chain: chain, FromHeight: fromHeight, ToHeight: toHeight}
	if err := c.send(&Request{Event: EventGetBlockRange, Data: req}); err != nil {
		return nil, err
	}
	res, err := c.read()
	if err != nil {
		return nil, err
	}

	resData := new(ResRangeBlockQuery)
	if err := parseResData(res.Data, &resData); err != nil {
		return nil, err
	}

	return resData, nil
}

func parseResData(data any, dest interface{}) error {
	jsonData, err := jsoniter.Marshal(data)
	if err != nil {
		return err
	}

	if err := jsoniter.Unmarshal(jsonData, dest); err != nil {
		return err
	}

	return nil
}
