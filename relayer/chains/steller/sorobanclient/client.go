package sorobanclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"sync/atomic"
)

const (
	jsonRPCVersion = "2.0"
)

type Client struct {
	idCounter  uint64
	httpClient *http.Client
	rpcUrl     string
}

func New(rpcUrl string, httpCl *http.Client) (*Client, error) {
	if _, err := url.Parse(rpcUrl); err != nil {
		return nil, err
	}

	if httpCl == nil {
		httpCl = &http.Client{}
	}

	return &Client{
		httpClient: httpCl,
		rpcUrl:     rpcUrl,
	}, nil
}

func (c *Client) SimulateTransaction(txnXdr string, resourceCfg *ResourceConfig) (*TxSimulationResult, error) {
	simResult := &TxSimulationResult{}
	if err := c.CallContext(
		context.Background(),
		simResult,
		"simulateTransaction",
		txnXdr,
		resourceCfg,
	); err != nil {
		return nil, err
	}

	return simResult, nil
}

func (c *Client) GetLatestLedger(ctx context.Context) (*LatestLedgerResponse, error) {
	ledgerRes := &LatestLedgerResponse{}
	if err := c.CallContext(ctx, ledgerRes, "getLatestLedger"); err != nil {
		return nil, err
	}
	return ledgerRes, nil
}

func (c *Client) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	if result != nil && reflect.TypeOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("call result parameter must be pointer or nil interface: %v", result)
	}

	msg, err := c.newMessage(method, args...)
	if err != nil {
		return err
	}

	respBody, err := c.doRequest(ctx, msg)
	if err != nil {
		return err
	}
	defer respBody.Close()

	var respmsg jsonRPCResponse
	if err := json.NewDecoder(respBody).Decode(&respmsg); err != nil {
		return err
	}
	if respmsg.Error != nil {
		return respmsg.Error
	}
	if len(respmsg.Result) == 0 {
		return fmt.Errorf("result is empty")
	}

	return json.Unmarshal(respmsg.Result, result)
}

func (c *Client) newMessage(method string, paramsIn ...interface{}) (*jsonRPCRequest, error) {
	msg := &jsonRPCRequest{Version: jsonRPCVersion, ID: c.nextID(), Method: method}
	if paramsIn != nil { // prevent sending "params":null
		var err error
		if msg.Params, err = json.Marshal(paramsIn); err != nil {
			return nil, err
		}
	}
	return msg, nil
}

func (c *Client) doRequest(ctx context.Context, msg interface{}) (io.ReadCloser, error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.rpcUrl, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return nil, err
	}
	req.ContentLength = int64(len(body))
	req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }

	req.Header.Set("Content-Type", "application/json")

	// do request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var buf bytes.Buffer
		var body []byte
		if _, err := buf.ReadFrom(resp.Body); err == nil {
			body = buf.Bytes()
		}

		return nil, HTTPError{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Body:       body,
		}
	}
	return resp.Body, nil
}

func (c *Client) nextID() json.RawMessage {
	id := atomic.AddUint64(&c.idCounter, 1)
	return strconv.AppendUint(nil, uint64(id), 10)
}
