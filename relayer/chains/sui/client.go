package sui

import (
	"context"
	"fmt"
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

	QueryTxBlocks(
		ctx context.Context,
		query types.SuiTransactionBlockResponseQuery,
		cursor *sui_types.TransactionDigest,
		limit *uint,
		descendingOrder bool,
	) (*types.TransactionBlocksPage, error)

	GetEvents(
		ctx context.Context,
		txDigest sui_types.TransactionDigest,
	) ([]types.SuiEvent, error)
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
					Checkpoint: txRes.Checkpoint,
				})
			}
		}
	}

	return events, nil
}

func (c *Client) QueryTxBlocks(
	ctx context.Context,
	query types.SuiTransactionBlockResponseQuery,
	cursor *sui_types.TransactionDigest,
	limit *uint,
	descendingOrder bool,
) (*types.TransactionBlocksPage, error) {
	return c.rpc.QueryTransactionBlocks(ctx, query, cursor, limit, descendingOrder)
}

func (c *Client) GetEvents(
	ctx context.Context,
	txDigest sui_types.TransactionDigest,
) ([]types.SuiEvent, error) {
	return c.rpc.GetEvents(ctx, txDigest)
}
