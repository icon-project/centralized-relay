package sui

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"slices"
	"strconv"

	"github.com/coming-chat/go-sui/v2/account"
	suisdkClient "github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/fardream/go-bcs/bcs"
	suitypes "github.com/icon-project/centralized-relay/relayer/chains/sui/types"
	"go.uber.org/zap"
)

const (
	suiCurrencyType                           = "0x2::sui::SUI"
	pickMethod                                = 1
	baseSuiFee                                = 1000
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
	EstimateGas(ctx context.Context, txBytes lib.Base64Data) (*types.DryRunTransactionBlockResponse, int64, error)
	ExecuteContract(ctx context.Context, suiMessage *SuiMessage, address string, gasBudget uint64) (*types.TransactionBytes, error)
	CommitTx(ctx context.Context, wallet *account.Account, txBytes lib.Base64Data, signatures []any) (*types.SuiTransactionBlockResponse, error)
	GetTransaction(ctx context.Context, txDigest string) (*types.SuiTransactionBlockResponse, error)
	QueryContract(ctx context.Context, suiMessage *SuiMessage, address string, gasBudget uint64, resPtr interface{}) error

	GetCheckpoints(ctx context.Context, req suitypes.SuiGetCheckpointsRequest) (*suitypes.PaginatedCheckpointsResponse, error)
	GetEventsFromTxBlocks(ctx context.Context, allowedEventTypes []string, digests []string) ([]suitypes.EventResponse, error)

	GetObject(ctx context.Context, objID sui_types.ObjectID, options *types.SuiObjectDataOptions) (*types.SuiObjectResponse, error)
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

func (c Client) GetObject(ctx context.Context, objID sui_types.ObjectID, options *types.SuiObjectDataOptions) (*types.SuiObjectResponse, error) {
	return c.rpc.GetObject(ctx, objID, options)
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

// Returns dry run result of txn with gas and status response
func (cl *Client) EstimateGas(ctx context.Context, txBytes lib.Base64Data) (*types.DryRunTransactionBlockResponse, int64, error) {
	dryrunResult, err := cl.rpc.DryRunTransaction(ctx, txBytes)
	return dryrunResult, dryrunResult.Effects.Data.GasFee(), err
}

func (cl *Client) ExecuteContract(ctx context.Context, suiMessage *SuiMessage, address string, gasBudget uint64) (*types.TransactionBytes, error) {
	accountAddress, err := move_types.NewAccountAddressHex(address)
	if err != nil {
		return &types.TransactionBytes{}, fmt.Errorf("error getting account address sender: %w", err)
	}
	packageId, err := move_types.NewAccountAddressHex(suiMessage.PackageObjectId)
	if err != nil {
		return &types.TransactionBytes{}, fmt.Errorf("invalid packageId: %w", err)
	}
	coinId := cl.getGasCoinId(ctx, address, gasBudget)
	coinAddress, err := move_types.NewAccountAddressHex(coinId.CoinObjectId.String())
	if err != nil {
		return &types.TransactionBytes{}, fmt.Errorf("error getting gas coinid : %w", err)
	}
	typeArgs := []string{}
	var stringParams []interface{}
	for _, s := range suiMessage.Params {
		stringParams = append(stringParams, fmt.Sprint(s.Val))
	}
	if stringParams == nil {
		stringParams = make([]interface{}, 0)
	}

	resp := types.TransactionBytes{}
	return &resp, cl.rpc.CallContext(
		ctx,
		&resp,
		moveCall,
		*accountAddress,
		packageId,
		suiMessage.Module,
		suiMessage.Method,
		typeArgs,
		stringParams,
		coinAddress,
		types.NewSafeSuiBigInt(gasBudget),
		"DevInspect",
	)
}

func (cl *Client) CommitTx(ctx context.Context, wallet *account.Account, txBytes lib.Base64Data, signatures []any) (*types.SuiTransactionBlockResponse, error) {
	return cl.rpc.ExecuteTransactionBlock(ctx, txBytes, signatures, &types.SuiTransactionBlockResponseOptions{
		ShowEffects: true,
		ShowEvents:  true,
	}, types.TxnRequestTypeWaitForLocalExecution)
}

func (c *Client) getGasCoinId(ctx context.Context, addr string, gasCost uint64) *types.Coin {
	accountAddress, err := move_types.NewAccountAddressHex(addr)
	if err != nil {
		c.log.Error(fmt.Sprintf("error getting account address sender %s", addr), zap.Error(err))
		return nil
	}
	result, err := c.rpc.GetSuiCoinsOwnedByAddress(ctx, *accountAddress)
	if err != nil {
		c.log.Error(fmt.Sprintf("error getting gas coins for address %s", addr), zap.Error(err))
		return nil
	}
	_, t, err := result.PickSUICoinsWithGas(big.NewInt(baseSuiFee), gasCost, pickMethod)
	if err != nil {
		c.log.Error(fmt.Sprintf("error getting gas coins with enough gas for address %s", addr), zap.Error(err))
		return nil
	}
	return t
}

func (cl *Client) GetTransaction(ctx context.Context, txDigest string) (*types.SuiTransactionBlockResponse, error) {
	txBlock, err := cl.rpc.GetTransactionBlock(ctx, lib.Base58(txDigest), types.SuiTransactionBlockResponseOptions{
		ShowEffects: true,
	})
	return txBlock, err
}

func (cl *Client) paramsToCallArgs(params []SuiCallArg) ([]sui_types.CallArg, error) {
	var callArgs []sui_types.CallArg
	for _, p := range params {
		switch p.Type {
		case CallArgObject:
			arg, err := cl.getCallArgObject(p.Val.(string))
			if err != nil {
				return nil, err
			}
			callArgs = append(callArgs, *arg)
		case CallArgPure:
			arg, err := cl.getCallArgPure(p.Val)
			if err != nil {
				return nil, err
			}
			callArgs = append(callArgs, *arg)
		default:
			return nil, fmt.Errorf("invalid call arg type")
		}
	}
	return callArgs, nil
}

func (cl *Client) getCallArgPure(arg interface{}) (*sui_types.CallArg, error) {
	byteParam, err := bcs.Marshal(arg)
	if err != nil {
		return nil, err
	}
	return &sui_types.CallArg{
		Pure: &byteParam,
	}, nil
}

func (cl *Client) getCallArgObject(arg string) (*sui_types.CallArg, error) {
	objectId, err := sui_types.NewAddressFromHex(arg)
	if err != nil {
		return nil, err
	}
	object, err := cl.GetObject(context.Background(), *objectId, &types.SuiObjectDataOptions{
		ShowType: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %v", err)
	}

	objRef := object.Data.Reference()

	// return &sui_types.CallArg{
	// 	Object: &sui_types.ObjectArg{
	// 		ImmOrOwnedObject: &objRef,
	// 	},
	// }, nil
	return &sui_types.CallArg{
		Object: &sui_types.ObjectArg{
			SharedObject: &struct {
				Id                   sui_types.ObjectID
				InitialSharedVersion sui_types.SequenceNumber
				Mutable              bool
			}{
				Id:                   objRef.ObjectId,
				InitialSharedVersion: 944734,
				Mutable:              true,
			},
		},
	}, nil
}

func (cl *Client) QueryContract(
	ctx context.Context,
	suiMessage *SuiMessage,
	address string,
	gasBudget uint64,
	resPtr interface{},
) error {
	builder := sui_types.NewProgrammableTransactionBuilder()
	packageId, err := move_types.NewAccountAddressHex(suiMessage.PackageObjectId)
	if err != nil {
		return err
	}
	senderAddress, err := move_types.NewAccountAddressHex(address)
	if err != nil {
		return err
	}

	callArgs, err := cl.paramsToCallArgs(suiMessage.Params)
	if err != nil {
		return err
	}

	err = builder.MoveCall(
		*packageId,
		move_types.Identifier(suiMessage.Module),
		move_types.Identifier(suiMessage.Method),
		[]move_types.TypeTag{},
		callArgs,
	)
	if err != nil {
		return err
	}
	bcsBytes, err := bcs.Marshal(builder.Finish())
	if err != nil {
		return err
	}
	txBytes := append([]byte{0}, bcsBytes...)
	b64Data, err := lib.NewBase64Data(base64.StdEncoding.EncodeToString(txBytes))
	if err != nil {
		return err
	}

	res, err := cl.rpc.DevInspectTransactionBlock(context.Background(), *senderAddress, *b64Data, nil, nil)
	if err != nil {
		return err
	}

	if res.Error != nil {
		return fmt.Errorf("error occurred while calling sui contract: %s", *res.Error)
	}
	if len(res.Results) > 0 && len(res.Results[0].ReturnValues) > 0 {
		returnVal := res.Results[0].ReturnValues[0]
		byteSlice, ok := returnVal.([]byte)
		if !ok {
			return err
		}
		if _, err := bcs.Unmarshal(byteSlice, resPtr); err != nil {
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
		suisdkClient.SuiMethod("getCheckpoints"),
		req.Cursor,
		req.Limit,
		req.DescendingOrder,
	); err != nil {
		return nil, err
	}

	return &paginatedRes, nil
}

func (c *Client) GetEventsFromTxBlocks(ctx context.Context, allowedEventTypes []string, digests []string) ([]suitypes.EventResponse, error) {
	txnBlockResponses := []*types.SuiTransactionBlockResponse{}

	if err := c.rpc.CallContext(
		ctx,
		&txnBlockResponses,
		suisdkClient.SuiMethod("multiGetTransactionBlocks"),
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
