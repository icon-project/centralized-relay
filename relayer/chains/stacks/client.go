package stacks

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/icon-project/centralized-relay/relayer/chains/stacks/interfaces"
	rpcClient "github.com/icon-project/stacks-go-sdk/pkg/rpc_client"
	blockchainApiClient "github.com/icon-project/stacks-go-sdk/pkg/stacks_blockchain_api_client"

	"go.uber.org/zap"
)

var _ interfaces.IClient = (*Client)(nil)

type Client struct {
	apiClient       blockchainApiClient.APIClient
	rpcApiClient    rpcClient.APIClient
	log             *zap.Logger
	mtx             sync.Mutex
	wsConn          *websocket.Conn
	idCursor        int64
	pendingRequests map[int64]chan *json.RawMessage
}

func NewClient(apiBaseURL string, logger *zap.Logger) (*Client, error) {
	cfg := blockchainApiClient.NewConfiguration()
	cfg.Servers = blockchainApiClient.ServerConfigurations{
		{
			URL:         apiBaseURL,
			Description: "Custom API Server",
		},
	}

	apiClient := blockchainApiClient.NewAPIClient(cfg)

	rpcCfg := rpcClient.NewConfiguration()
	rpcCfg.Servers = rpcClient.ServerConfigurations{
		{
			URL:         apiBaseURL,
			Description: "Custom API Server",
		},
	}
	rpcApiClient := rpcClient.NewAPIClient(rpcCfg)

	return &Client{
		apiClient:       *apiClient,
		rpcApiClient:    *rpcApiClient,
		log:             logger,
		pendingRequests: make(map[int64]chan *json.RawMessage),
	}, nil
}

func (c *Client) Log() *zap.Logger {
	return c.log
}

func (c *Client) GetAccountBalance(ctx context.Context, address string) (*big.Int, error) {
	principal := blockchainApiClient.GetFilteredEventsAddressParameter{
		String: &address,
	}

	resp, _, err := c.apiClient.AccountsAPI.GetAccountBalance(ctx, principal).Execute()
	if err != nil {
		return nil, err
	}

	balanceStr := resp.Stx.Balance
	balance, ok := new(big.Int).SetString(balanceStr, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse balance: %s", balanceStr)
	}

	return balance, nil
}

func (c *Client) GetAccountNonce(ctx context.Context, address string) (uint64, error) {
	principal := blockchainApiClient.GetFilteredEventsAddressParameter{
		String: &address,
	}

	resp, _, err := c.apiClient.AccountsAPI.GetAccountNonces(ctx, principal).Execute()
	if err != nil {
		return 0, err
	}

	return uint64(resp.PossibleNextNonce), nil
}

func (c *Client) GetBlockByHeightOrHash(ctx context.Context, height uint64) (*blockchainApiClient.GetBlocks200ResponseResultsInner, error) {
	heightOrHash := blockchainApiClient.GetBlockHeightOrHashParameter{
		Uint64: &height,
	}

	resp, _, err := c.apiClient.BlocksAPI.GetBlock(ctx, heightOrHash).Execute()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetLatestBlock(ctx context.Context) (*blockchainApiClient.GetBlocks200ResponseResultsInner, error) {
	resp, _, err := c.apiClient.BlocksAPI.GetBlocks(ctx).Limit(1).Execute()
	if err != nil {
		return nil, err
	}
	return &resp.Results[0], nil
}

func (c *Client) CallReadOnlyFunction(ctx context.Context, contractAddress string, contractName string, functionName string, functionArgs []string) (*string, error) {
	fa := rpcClient.NewReadOnlyFunctionArgsschema(
		"ST1PQHQKV0RJXZFY1DGX8MNSNYVE3VGZJSRTPGZGM",
		functionArgs,
	)

	resp, _, err := c.rpcApiClient.SmartContractsAPI.CallReadOnlyFunction(ctx, contractAddress, contractName, functionName).ReadOnlyFunctionArgsschema(*fa).Execute()
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

func (c *Client) SubscribeToEvents(ctx context.Context, eventTypes []string, callback interfaces.EventCallback) error {
	wsURL, _ := c.apiClient.GetConfig().Servers.URL(0, make(map[string]string))

	if strings.HasPrefix(wsURL, "http://") {
		wsURL = strings.Replace(wsURL, "http://", "ws://", 1)
	} else if strings.HasPrefix(wsURL, "https://") {
		wsURL = strings.Replace(wsURL, "https://", "wss://", 1)
	}
	if !strings.HasSuffix(wsURL, "/") {
		wsURL += "/"
	}
	wsURL += "extended/v1/ws"

	var err error
	c.wsConn, _, err = websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}

	go c.listenToMessages(ctx, callback)

	for _, eventType := range eventTypes {
		if err := c.subscribe(eventType); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) sendRPCRequest(method string, params interface{}) (int64, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.idCursor++
	id := c.idCursor

	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"method":  method,
		"params":  params,
	}

	message, err := json.Marshal(request)
	if err != nil {
		return 0, err
	}

	c.pendingRequests[id] = make(chan *json.RawMessage, 1)

	err = c.wsConn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		delete(c.pendingRequests, id)
		return 0, err
	}

	return id, nil
}

func (c *Client) subscribe(eventType string) error {
	params := map[string]interface{}{
		"event": eventType,
	}
	id, err := c.sendRPCRequest("subscribe", params)
	if err != nil {
		return err
	}

	responseChan := c.pendingRequests[id]
	select {
	case response := <-responseChan:
		var resp struct {
			Jsonrpc string                 `json:"jsonrpc"`
			Id      int64                  `json:"id"`
			Result  map[string]interface{} `json:"result"`
			Error   map[string]interface{} `json:"error"`
		}
		err = json.Unmarshal(*response, &resp)
		if err != nil {
			return err
		}
		if resp.Error != nil {
			return fmt.Errorf("subscription error: %v", resp.Error)
		}
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("subscription timeout")
	}
}

func (c *Client) listenToMessages(ctx context.Context, callback interfaces.EventCallback) {
	defer c.wsConn.Close()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, message, err := c.wsConn.ReadMessage()
			if err != nil {
				c.log.Error("WebSocket read error", zap.Error(err))
				return
			}

			var parsedMessage map[string]interface{}
			if err := json.Unmarshal(message, &parsedMessage); err != nil {
				c.log.Error("Failed to parse message", zap.Error(err))
				continue
			}

			c.handleMessage(parsedMessage, callback)
		}
	}
}

func (c *Client) handleMessage(message map[string]interface{}, callback interfaces.EventCallback) {
	if idVal, ok := message["id"]; ok {
		idFloat, ok := idVal.(float64)
		if !ok {
			c.log.Error("Invalid id in response")
			return
		}
		id := int64(idFloat)

		c.mtx.Lock()
		responseChan, ok := c.pendingRequests[id]
		c.mtx.Unlock()
		if !ok {
			c.log.Error("Received response for unknown request id", zap.Int64("id", id))
			return
		}

		rawMessage, err := json.Marshal(message)
		if err != nil {
			c.log.Error("Failed to marshal message", zap.Error(err))
			return
		}
		responseChan <- (*json.RawMessage)(&rawMessage)

		c.mtx.Lock()
		delete(c.pendingRequests, id)
		c.mtx.Unlock()
	} else {
		if method, ok := message["method"]; ok {
			methodStr, _ := method.(string)
			params := message["params"]
			if err := callback(methodStr, params); err != nil {
				c.log.Error("Callback error", zap.Error(err))
			}
		} else {
			c.log.Error("Unknown message format", zap.Any("message", message))
		}
	}
}
