package stellar

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"sync/atomic"

	"github.com/icon-project/goloop/common/log"
	"github.com/stellar/go/xdr"
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

func (c *Client) GetLatestLedger(ctx context.Context) (*LatestLedgerResponse, error) {
	ledgerRes := &LatestLedgerResponse{}
	if err := c.CallContext(ctx, ledgerRes, "getLatestLedger", nil); err != nil {
		return nil, err
	}
	return ledgerRes, nil
}

func (c *Client) GetEvent(ctx context.Context, height uint64, inputSn, contractId, eventSignature string) (*EventResponseEvent, error) {
	eventResponse := &EventResponse{}
	params := EventQueryFilter{
		StartLedger: height - 500,
		Pagination: Pagination{
			Limit: 100,
		},
		Filters: []Filter{
			{
				Type: "contract",
				ContractIDS: []string{
					contractId,
				},
			},
		},
	}
	if err := c.CallContext(ctx, eventResponse, "getEvents", params); err != nil {
		return nil, err
	}
	for _, event := range eventResponse.Events {
		for _, topic := range event.Topic {
			decodedTopic, _ := decodeTopic(topic)
			if decodedTopic == eventSignature {
				// check For Value
				var xdrValue xdr.ScVal
				bytesD, err := base64.StdEncoding.DecodeString(event.Value)
				if err != nil {
					return nil, err
				}
				err = xdr.SafeUnmarshal(bytesD, &xdrValue)
				if err != nil {
					return nil, err
				}
				eventValues, ok := xdrValue.GetMap()
				if !ok {
					return nil, fmt.Errorf("error geting map from values")
				}
				for _, mapItem := range *eventValues {
					if mapItem.Key.String() == "sn" || mapItem.Key.String() == "reqId" {
						sn, ok := mapItem.Val.GetU128()
						if !ok {
							return nil, fmt.Errorf("failed to decode sn")
						}
						if strconv.FormatUint(uint64(sn.Lo), 10) == inputSn {
							decodecMap := convertScMapToMap(*eventValues)
							event.ValueDecoded = decodecMap
							return &event, nil
						}
					}
				}

			}
		}
	}
	return nil, fmt.Errorf("event not found")
}

func convertScMapToMap(scMap xdr.ScMap) map[string]interface{} {
	normalMap := make(map[string]interface{})

	for _, entry := range scMap {
		key := entry.Key.String()
		var value interface{}
		switch entry.Val.Type {
		case xdr.ScValTypeScvBool:
			valueBool, _ := entry.Val.GetB()
			value = bool(valueBool)
		case xdr.ScValTypeScvBytes:
			valueBytes, _ := entry.Val.GetBytes()
			value = []byte(valueBytes)
		case xdr.ScValTypeScvU64:
			value64, _ := entry.Val.GetU64()
			value = uint64(value64)
		case xdr.ScValTypeScvU128:
			value128, _ := entry.Val.GetU128()
			value = uint64(value128.Lo)
		case xdr.ScValTypeScvString:
			valueString, _ := entry.Val.GetStr()
			value = string(valueString)
		case xdr.ScValTypeScvU32:
			value32, _ := entry.Val.GetU32()
			value = uint64(value32)
		default:
			log.Info("Encountered unmatched type ", entry.Val.Type)
			value = entry.Val
		}
		normalMap[key] = value
	}

	return normalMap
}

func decodeTopic(topic string) (string, error) {
	bytesD, err := base64.StdEncoding.DecodeString(topic)
	if err != nil {
		return "", err
	}
	var xdrValue xdr.ScVal
	err = xdr.SafeUnmarshal(bytesD, &xdrValue)
	if err != nil {
		return "", err
	}
	return xdrValue.String(), nil
}

func (c *Client) CallContext(ctx context.Context, result interface{}, method string, params interface{}) error {
	if result != nil && reflect.TypeOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("call result parameter must be pointer or nil interface: %v", result)
	}

	msg, err := c.newMessage(method, params)
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
	return json.Unmarshal(respmsg.Result, &result)
}

func (c *Client) newMessage(method string, paramsIn interface{}) (*jsonRPCRequest, error) {
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
