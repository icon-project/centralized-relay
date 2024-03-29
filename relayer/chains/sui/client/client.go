package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	suiwsconn "github.com/block-vision/sui-go-sdk/common/wsconn"
	suimodels "github.com/block-vision/sui-go-sdk/models"
	suisdk "github.com/block-vision/sui-go-sdk/sui"
	"github.com/gorilla/websocket"
	"github.com/icon-project/centralized-relay/relayer/chains/sui/types"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

type Client struct {
	rpc suisdk.ISuiAPI
	log *zap.Logger
}

func NewClient(rpcClient suisdk.ISuiAPI, l *zap.Logger) *Client {
	return &Client{
		rpc: rpcClient,
		log: l,
	}
}

func (c *Client) GetLatestCheckpointSeq(ctx context.Context) (uint64, error) {
	return c.rpc.SuiGetLatestCheckpointSequenceNumber(ctx)
}

func (c *Client) GetCheckpoints(ctx context.Context, req suimodels.SuiGetCheckpointsRequest) (suimodels.PaginatedCheckpointsResponse, error) {
	return c.rpc.SuiGetCheckpoints(ctx, req)
}

func (c *Client) GetBalance(ctx context.Context, addr string) ([]suimodels.CoinData, error) {
	result, err := c.rpc.SuiXGetAllCoins(ctx, suimodels.SuiXGetAllCoinsRequest{
		Owner: addr,
	})
	if err != nil {
		c.log.Error(fmt.Sprintf("error getting balance for address %s", addr), zap.Error(err))
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) SubscribeEventNotification(done chan interface{}, wsUrl string, eventFilters []interface{}) (<-chan types.EventNotification, error) {
	rpcReq := suimodels.JsonRPCRequest{
		JsonRPC: "2.0",
		ID:      time.Now().UnixMilli(),
		Method:  "suix_subscribeEvent",
		Params:  eventFilters,
	}

	reqBytes, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, errors.New("failed to json encode rpc request")
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		conn.Close()
		return nil, errors.Wrap(err, "failed to create ws connection")
	}

	err = conn.WriteMessage(websocket.TextMessage, reqBytes)
	if nil != err {
		conn.Close()
		return nil, errors.Wrap(err, "failed to send ws rpc request")
	}

	_, messageData, err := conn.ReadMessage()
	if nil != err {
		conn.Close()
		return nil, errors.Wrap(err, "failed to get ws rpc response")
	}

	var resp suiwsconn.SubscriptionResp
	if gjson.ParseBytes(messageData).Get("error").Exists() {
		conn.Close()
		return nil, fmt.Errorf(gjson.ParseBytes(messageData).Get("error").String())
	}

	if err = json.Unmarshal([]byte(gjson.ParseBytes(messageData).String()), &resp); err != nil {
		conn.Close()
		return nil, err
	}

	enStream := make(chan types.EventNotification)
	go func() {
		defer close(enStream)
		for {
			select {
			case <-done:
				conn.Close()
				return
			default:
				en := types.EventNotification{}
				if err := c.readWsConnMessage(conn, &en); err != nil {
					conn.Close()
					en.Error = errors.Wrap(err, "failed to read incoming event notification")
					enStream <- en
				} else if en.PackageId != "" {
					enStream <- en
				}
			}
		}
	}()

	return enStream, nil
}

func (c *Client) readWsConnMessage(conn *websocket.Conn, dest interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s panic: %v", types.WsConnReadError, r)
		}
	}()

	mt, messageData, readErr := conn.ReadMessage()
	if readErr != nil {
		return errors.Wrap(err, types.WsConnReadError)
	}

	if mt == websocket.TextMessage {
		if gjson.ParseBytes(messageData).Get("error").Exists() {
			return errors.New(gjson.ParseBytes(messageData).Get("error").String())
		}

		err := json.Unmarshal([]byte(gjson.ParseBytes(messageData).Get("params.result").String()), &dest)
		if err != nil {
			return err
		}
	}

	return nil
}
