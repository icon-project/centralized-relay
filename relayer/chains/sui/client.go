package sui

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/coming-chat/go-sui/v2/account"
	suisdkClient "github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/fardream/go-bcs/bcs"
	"github.com/gorilla/websocket"
	suitypes "github.com/icon-project/centralized-relay/relayer/chains/sui/types"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

const (
	suiCurrencyType                           = "0x2::sui::SUI"
	suiStringType                             = "0x1::string::String"
	suiU64                                    = "u64"
	suiBool                                   = "bool"
	moveCall        suisdkClient.UnsafeMethod = "moveCall"

	CallArgPure   = "pure"
	CallArgObject = "object"
)

type IClient interface {
	GetLatestCheckpointSeq(ctx context.Context) (uint64, error)
	GetTotalBalance(ctx context.Context, addr string) (uint64, error)
	SimulateTx(ctx context.Context, txBytes lib.Base64Data) (*types.DryRunTransactionBlockResponse, int64, error)
	ExecuteTx(ctx context.Context, wallet *account.Account, txBytes lib.Base64Data, signatures []any) (*types.SuiTransactionBlockResponse, error)
	GetTransaction(ctx context.Context, txDigest string) (*types.SuiTransactionBlockResponse, error)
	QueryContract(ctx context.Context, senderAddr string, txBytes lib.Base64Data, resPtr interface{}) error

	GetCheckpoint(ctx context.Context, checkpoint uint64) (*suitypes.CheckpointResponse, error)
	GetCheckpoints(ctx context.Context, req suitypes.SuiGetCheckpointsRequest) (*suitypes.PaginatedCheckpointsResponse, error)
	GetEventsFromTxBlocks(ctx context.Context, allowedEventTypes []string, digests []string) ([]suitypes.EventResponse, error)

	GetObject(ctx context.Context, objID sui_types.ObjectID, options *types.SuiObjectDataOptions) (*types.SuiObjectResponse, error)

	GetCoins(ctx context.Context, accountAddress string) (types.Coins, error)

	MoveCall(
		ctx context.Context,
		signer move_types.AccountAddress,
		packageId move_types.AccountAddress,
		module, function string,
		typeArgs []string,
		arguments []any,
		gas *move_types.AccountAddress,
		gasBudget types.SafeSuiBigInt[uint64],
	) (*types.TransactionBytes, error)

	SubscribeEventNotification(
		done chan interface{},
		wsUrl string,
		eventFilters interface{},
	) (<-chan suitypes.EventNotification, error)

	QueryEvents(
		ctx context.Context,
		req suitypes.EventQueryRequest,
	) (*suitypes.EventQueryResponse, error)
}

type Client struct {
	rpc *suisdkClient.Client
	log *zap.Logger
}

func NewClient(rpcClient *suisdkClient.Client, l *zap.Logger) *Client {
	return &Client{
		rpc: rpcClient,
		log: l,
	}
}

func (c Client) MoveCall(
	ctx context.Context,
	signer move_types.AccountAddress,
	packageId move_types.AccountAddress,
	module, function string,
	typeArgs []string,
	arguments []any,
	gas *move_types.AccountAddress,
	gasBudget types.SafeSuiBigInt[uint64],
) (*types.TransactionBytes, error) {
	return c.rpc.MoveCall(ctx, signer, packageId, module, function, typeArgs, arguments, gas, gasBudget)
}

func (c Client) GetObject(ctx context.Context, objID sui_types.ObjectID, options *types.SuiObjectDataOptions) (*types.SuiObjectResponse, error) {
	return c.rpc.GetObject(ctx, objID, options)
}

func (c Client) GetCoins(ctx context.Context, addr string) (types.Coins, error) {
	accountAddress, err := move_types.NewAccountAddressHex(addr)
	if err != nil {
		return nil, err
	}
	return c.rpc.GetSuiCoinsOwnedByAddress(ctx, *accountAddress)
}

func (c Client) GetLatestCheckpointSeq(ctx context.Context) (uint64, error) {
	checkPoint, err := c.rpc.GetLatestCheckpointSequenceNumber(ctx)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(checkPoint, 10, 64)
}

func (c *Client) GetTotalBalance(ctx context.Context, addr string) (uint64, error) {
	accountAddress, err := move_types.NewAccountAddressHex(addr)
	if err != nil {
		return 0, fmt.Errorf("error getting balance: %w", err)
	}
	res, err := c.rpc.GetBalance(ctx, *accountAddress, suiCurrencyType)
	if err != nil {
		return 0, fmt.Errorf("error getting balance: %w", err)
	}
	return res.TotalBalance.BigInt().Uint64(), nil
}

func (cl *Client) SimulateTx(ctx context.Context, txBytes lib.Base64Data) (*types.DryRunTransactionBlockResponse, int64, error) {
	dryrunResult, err := cl.rpc.DryRunTransaction(ctx, txBytes)
	return dryrunResult, dryrunResult.Effects.Data.GasFee(), err
}

func (cl *Client) ExecuteTx(ctx context.Context, wallet *account.Account, txBytes lib.Base64Data, signatures []any) (*types.SuiTransactionBlockResponse, error) {
	return cl.rpc.ExecuteTransactionBlock(ctx, txBytes, signatures, &types.SuiTransactionBlockResponseOptions{
		ShowEffects: true,
		ShowEvents:  true,
	}, types.TxnRequestTypeWaitForLocalExecution)
}

func (cl *Client) GetTransaction(ctx context.Context, txDigest string) (*types.SuiTransactionBlockResponse, error) {
	b58Digest, err := lib.NewBase58(txDigest)
	if err != nil {
		return nil, err
	}
	txBlock, err := cl.rpc.GetTransactionBlock(ctx, *b58Digest, types.SuiTransactionBlockResponseOptions{
		ShowEffects: true,
	})
	return txBlock, err
}

func (cl *Client) QueryContract(ctx context.Context, senderAddr string, txBytes lib.Base64Data, resPtr interface{}) error {
	senderAddress, err := move_types.NewAccountAddressHex(senderAddr)
	if err != nil {
		return err
	}

	res, err := cl.rpc.DevInspectTransactionBlock(context.Background(), *senderAddress, txBytes, nil, nil)
	if err != nil {
		return err
	}

	if res.Error != nil {
		return fmt.Errorf("error occurred while calling sui contract: %s", *res.Error)
	}
	if len(res.Results) > 0 && len(res.Results[0].ReturnValues) > 0 {
		returnValues := res.Results[0].ReturnValues[0]
		returnResult := returnValues.([]interface{})[0]

		if _, ok := returnResult.([]byte); ok {
			if _, err := bcs.Unmarshal([]byte(returnResult.([]byte)), resPtr); err != nil {
				return err
			}
			return nil
		}

		resultBytes := []byte{}
		for _, el := range returnResult.([]interface{}) {
			resultBytes = append(resultBytes, byte(el.(float64)))
		}

		if _, err := bcs.Unmarshal(resultBytes, resPtr); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("got empty result")
}

func (c *Client) GetCheckpoints(ctx context.Context, req suitypes.SuiGetCheckpointsRequest) (*suitypes.PaginatedCheckpointsResponse, error) {
	paginatedRes := suitypes.PaginatedCheckpointsResponse{}
	if err := c.rpc.CallContext(
		ctx,
		&paginatedRes,
		suitypes.SuiMethod("sui_getCheckpoints"),
		req.Cursor,
		req.Limit,
		req.DescendingOrder,
	); err != nil {
		return nil, err
	}

	return &paginatedRes, nil
}

func (c *Client) GetCheckpoint(ctx context.Context, checkpoint uint64) (*suitypes.CheckpointResponse, error) {
	checkpointRes := suitypes.CheckpointResponse{}
	if err := c.rpc.CallContext(
		ctx,
		&checkpointRes,
		suitypes.SuiMethod("sui_getCheckpoint"),
		strconv.Itoa(int(checkpoint)),
	); err != nil {
		return nil, err
	}

	return &checkpointRes, nil
}

func (c *Client) GetEventsFromTxBlocks(ctx context.Context, allowedEventTypes []string, digests []string) ([]suitypes.EventResponse, error) {
	txnBlockResponses := []*types.SuiTransactionBlockResponse{}

	if err := c.rpc.CallContext(
		ctx,
		&txnBlockResponses,
		suitypes.SuiMethod("sui_multiGetTransactionBlocks"),
		digests,
		types.SuiTransactionBlockResponseOptions{ShowEvents: true},
	); err != nil {
		return nil, err
	}

	var events []suitypes.EventResponse
	for _, txRes := range txnBlockResponses {
		for _, ev := range txRes.Events {
			if slices.Contains(allowedEventTypes, ev.Type) {
				events = append(events, suitypes.EventResponse{
					SuiEvent:   ev,
					Checkpoint: txRes.Checkpoint.Uint64(),
				})
			}
		}
	}

	return events, nil
}

func (c *Client) readWsConnMessage(conn *websocket.Conn, dest interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s panic: %v", suitypes.WsConnReadError, r)
		}
	}()

	mt, messageData, readErr := conn.ReadMessage()
	if readErr != nil {
		return fmt.Errorf("%s: %w", suitypes.WsConnReadError, err)
	}

	if mt == websocket.TextMessage {
		if gjson.ParseBytes(messageData).Get("error").Exists() {
			return fmt.Errorf(gjson.ParseBytes(messageData).Get("error").String())
		}

		err := json.Unmarshal([]byte(gjson.ParseBytes(messageData).Get("params.result").String()), &dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) SubscribeEventNotification(done chan interface{}, wsUrl string, eventFilters interface{}) (<-chan suitypes.EventNotification, error) {
	rpcReq := suitypes.JsonRPCRequest{
		Version: "2.0",
		ID:      time.Now().UnixMilli(),
		Method:  "suix_subscribeEvent",
		Params: []interface{}{
			eventFilters,
		},
	}

	reqBytes, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to json encode rpc request")
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create ws connection: %w", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, reqBytes)
	if nil != err {
		conn.Close()
		return nil, fmt.Errorf("failed to send ws rpc request: %w", err)
	}

	_, messageData, err := conn.ReadMessage()
	if nil != err {
		conn.Close()
		return nil, fmt.Errorf("failed to get ws rpc response: %w", err)
	}

	var resp suitypes.WsSubscriptionResp
	if gjson.ParseBytes(messageData).Get("error").Exists() {
		conn.Close()
		return nil, fmt.Errorf(gjson.ParseBytes(messageData).Get("error").String())
	}

	if err = json.Unmarshal([]byte(gjson.ParseBytes(messageData).String()), &resp); err != nil {
		conn.Close()
		return nil, err
	}

	enStream := make(chan suitypes.EventNotification)
	go func() {
		defer close(enStream)
		for {
			select {
			case <-done:
				conn.Close()
				return
			default:
				en := suitypes.EventNotification{}
				if err := c.readWsConnMessage(conn, &en); err != nil {
					conn.Close()
					en.Error = fmt.Errorf("failed to read incoming event notification: %w", err)
					enStream <- en
				} else if en.PackageId.String() != "" {
					enStream <- en
				}
			}
		}
	}()

	return enStream, nil
}

func (c *Client) QueryEvents(ctx context.Context, req suitypes.EventQueryRequest) (*suitypes.EventQueryResponse, error) {
	events := suitypes.EventQueryResponse{}
	if err := c.rpc.CallContext(
		ctx,
		&events,
		suitypes.SuiMethod("suix_queryEvents"),
		req.EventFilter,
		req.Cursor,
		req.Limit,
		req.Descending,
	); err != nil {
		return nil, err
	}

	return &events, nil
}
