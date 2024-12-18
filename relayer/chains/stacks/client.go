package stacks

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/interfaces"
	"github.com/icon-project/stacks-go-sdk/pkg/clarity"
	rpcClient "github.com/icon-project/stacks-go-sdk/pkg/rpc_client"
	"github.com/icon-project/stacks-go-sdk/pkg/stacks"
	blockchainApiClient "github.com/icon-project/stacks-go-sdk/pkg/stacks_blockchain_api_client"
	"github.com/icon-project/stacks-go-sdk/pkg/transaction"

	"go.uber.org/zap"
)

var _ interfaces.IClient = (*Client)(nil)

type Client struct {
	apiClient       blockchainApiClient.APIClient
	rpcApiClient    rpcClient.APIClient
	log             *zap.Logger
	pendingRequests map[int64]chan *json.RawMessage
	network         *stacks.StacksNetwork
}

func NewClient(logger *zap.Logger, network *stacks.StacksNetwork) (*Client, error) {
	cfg := blockchainApiClient.NewConfiguration()
	cfg.Servers = blockchainApiClient.ServerConfigurations{
		{
			URL: network.CoreAPIURL,
		},
	}

	apiClient := blockchainApiClient.NewAPIClient(cfg)

	rpcCfg := rpcClient.NewConfiguration()
	rpcCfg.Servers = rpcClient.ServerConfigurations{
		{
			URL: network.CoreAPIURL,
		},
	}
	rpcApiClient := rpcClient.NewAPIClient(rpcCfg)

	return &Client{
		apiClient:       *apiClient,
		rpcApiClient:    *rpcApiClient,
		log:             logger,
		pendingRequests: make(map[int64]chan *json.RawMessage),
		network:         network,
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

func (c *Client) GetContractById(ctx context.Context, contractId string) (*blockchainApiClient.SmartContract, error) {
	resp, httpResp, err := c.apiClient.SmartContractsAPI.GetContractById(ctx, contractId).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return nil, nil // Contract doesn't exist
		}
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetContractEvents(ctx context.Context, contractId string, limit, offset int32) (*blockchainApiClient.GetContractEventsById200Response, error) {
	req := c.apiClient.SmartContractsAPI.GetContractEventsById(ctx, contractId)

	if limit > 0 {
		req = req.Limit(limit)
	}

	if offset > 0 {
		req = req.Offset(offset)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		if httpResp != nil {
			return nil, fmt.Errorf("failed to get contract events (status %d): %w", httpResp.StatusCode, err)
		}
		return nil, fmt.Errorf("failed to get contract events: %w", err)
	}

	return resp, nil
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

func (c *Client) GetTransactionById(ctx context.Context, id string) (*blockchainApiClient.GetTransactionById200Response, error) {
	response, httpResponse, err := c.apiClient.TransactionsAPI.GetTransactionById(ctx, id).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by ID: %w", err)
	}

	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response: %d", httpResponse.StatusCode)
	}

	return response, nil
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
	req := c.apiClient.BlocksAPI.GetBlocks(ctx).Limit(1)
	resp, httpResp, err := req.Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 200 {
			var rawMap map[string]interface{}
			if err := json.NewDecoder(httpResp.Body).Decode(&rawMap); err != nil {
				return nil, fmt.Errorf("failed to decode response: %w", err)
			}

			results, ok := rawMap["results"].([]interface{})
			if !ok || len(results) == 0 {
				return nil, fmt.Errorf("no blocks found")
			}

			firstBlock, ok := results[0].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid block format")
			}

			height, ok := firstBlock["height"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid height format")
			}

			return &blockchainApiClient.GetBlocks200ResponseResultsInner{
				Height: int32(height),
			}, nil
		}
		return nil, err
	}

	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("no blocks found")
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

func (c *Client) MakeContractCall(
	ctx context.Context,
	contractAddress string,
	contractName string,
	functionName string,
	args []clarity.ClarityValue,
	senderAddress string,
	senderKey []byte,
) (*transaction.ContractCallTransaction, error) {
	tx, err := transaction.MakeContractCall(
		contractAddress,
		contractName,
		functionName,
		args,
		*c.network,
		senderAddress,
		senderKey,
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract call transaction: %w", err)
	}

	return tx, nil
}

func (c *Client) BroadcastTransaction(ctx context.Context, tx transaction.StacksTransaction) (string, error) {
	txID, err := transaction.BroadcastTransaction(tx, c.network)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	return txID, nil
}

func (c *Client) GetCurrentImplementation(ctx context.Context, contractAddress string) (clarity.ClarityValue, error) {
	functionName := "get-current-implementation"

	result, err := c.CallReadOnlyFunction(ctx, contractAddress, "xcall-proxy", functionName, []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to get current implementation: %w", err)
	}

	fmt.Printf("result: %v", result)

	byteResult, err := hex.DecodeString(strings.TrimPrefix(*result, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to hex decode current implementation: %w", err)
	}

	impl, err := clarity.DeserializeClarityValue(byteResult)
	if err != nil {
		return nil, fmt.Errorf("unexpected type for implementation principal")
	}

	return impl, nil
}

func (c *Client) SetAdmin(ctx context.Context, contractAddress string, newAdmin string, currentImplementation clarity.ClarityValue, senderAddress string, senderKey []byte) (string, error) {
	functionName := "set-admin"

	newAdminPrincipal, err := clarity.StringToPrincipal(newAdmin)
	if err != nil {
		return "", fmt.Errorf("invalid new admin address: %w", err)
	}

	currentImplementation_, err := clarity.StringToPrincipal("ST15C893XJFJ6FSKM020P9JQDB5T7X6MQTXMBPAVH.xcall-impl")
	if err != nil {
		return "", fmt.Errorf("invalid new admin address: %w", err)
	}

	args := []clarity.ClarityValue{newAdminPrincipal, currentImplementation_}

	tx, err := c.MakeContractCall(
		ctx,
		contractAddress,
		"xcall-proxy",
		functionName,
		args,
		senderAddress,
		senderKey,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create contract call transaction: %w", err)
	}

	txID, err := c.BroadcastTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	if txID == "" {
		return "", fmt.Errorf("got empty transaction ID after broadcasting")
	}

	return txID, nil
}

func (c *Client) GetReceipt(ctx context.Context, contractAddress string, srcNetwork string, connSnIn *big.Int) (bool, error) {
	srcNetworkArg, err := clarity.NewStringASCII(srcNetwork)
	if err != nil {
		return false, fmt.Errorf("failed to create srcNetwork argument: %w", err)
	}
	encodedSrcNetwork, err := srcNetworkArg.Serialize()
	if err != nil {
		return false, fmt.Errorf("failed to serialize srcNetwork argument: %w", err)
	}
	hexEncodedSrcNetwork := hex.EncodeToString(encodedSrcNetwork)

	connSnInArg, err := clarity.NewInt(connSnIn)
	if err != nil {
		return false, fmt.Errorf("failed to create connSnIn argument: %w", err)
	}
	encodedConnSnIn, err := connSnInArg.Serialize()
	if err != nil {
		return false, fmt.Errorf("failed to serialize connSnIn argument: %w", err)
	}
	hexEncodedConnSnInArg := hex.EncodeToString(encodedConnSnIn)

	result, err := c.CallReadOnlyFunction(
		ctx,
		contractAddress,
		"centralized-connection",
		"get-receipt",
		[]string{hexEncodedSrcNetwork, hexEncodedConnSnInArg},
	)
	if err != nil {
		return false, fmt.Errorf("failed to call get-receipt: %w", err)
	}

	var response struct {
		Ok bool `json:"ok"`
	}
	if err := json.Unmarshal([]byte(*result), &response); err != nil {
		return false, fmt.Errorf("failed to parse get-receipt response: %w", err)
	}

	return response.Ok, nil
}

func (c *Client) ClaimFee(ctx context.Context, contractAddress string, senderAddress string, senderKey []byte) (string, error) {
	args := []clarity.ClarityValue{}
	tx, err := c.MakeContractCall(
		ctx,
		contractAddress,
		"centralized-connection",
		"claim-fees",
		args,
		senderAddress,
		senderKey,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create claim-fees transaction: %w", err)
	}

	txID, err := c.BroadcastTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast claim-fees transaction: %w", err)
	}

	return txID, nil
}

func (c *Client) SetFee(ctx context.Context, contractAddress string, networkID string, messageFee *big.Int, responseFee *big.Int, senderAddress string, senderKey []byte) (string, error) {
	networkIDArg, err := clarity.NewStringASCII(networkID)
	if err != nil {
		return "", fmt.Errorf("failed to create networkID argument: %w", err)
	}

	messageFeeArg, err := clarity.NewUInt(messageFee.String())
	if err != nil {
		return "", fmt.Errorf("failed to create messageFee argument: %w", err)
	}

	responseFeeArg, err := clarity.NewUInt(responseFee.String())
	if err != nil {
		return "", fmt.Errorf("failed to create responseFee argument: %w", err)
	}

	args := []clarity.ClarityValue{networkIDArg, messageFeeArg, responseFeeArg}
	tx, err := c.MakeContractCall(
		ctx,
		contractAddress,
		"centralized-connection",
		"set-fee",
		args,
		senderAddress,
		senderKey,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create set-fee transaction: %w", err)
	}

	txID, err := c.BroadcastTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast set-fee transaction: %w", err)
	}

	return txID, nil
}

func (c *Client) GetFee(ctx context.Context, contractAddress string, networkID string, responseFee bool) (uint64, error) {
	networkIDArg, err := clarity.NewStringASCII(networkID)
	if err != nil {
		return 0, fmt.Errorf("failed to create networkID argument: %w", err)
	}
	encodedNetworkID, err := networkIDArg.Serialize()
	if err != nil {
		return 0, fmt.Errorf("failed to serialize networkID argument: %w", err)
	}
	hexEncodedNetworkID := hex.EncodeToString(encodedNetworkID)

	responseFeeArg := clarity.NewBool(responseFee)
	encodedResponseFee, err := responseFeeArg.Serialize()
	if err != nil {
		return 0, fmt.Errorf("failed to serialize responseFee argument: %w", err)
	}
	hexEncodedResponseFee := hex.EncodeToString(encodedResponseFee)

	result, err := c.CallReadOnlyFunction(
		ctx,
		contractAddress,
		"centralized-connection",
		"get-fee",
		[]string{hexEncodedNetworkID, hexEncodedResponseFee},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to call get-fee: %w", err)
	}

	var response struct {
		Ok uint64 `json:"ok"`
	}
	if err := json.Unmarshal([]byte(*result), &response); err != nil {
		return 0, fmt.Errorf("failed to parse get-fee response: %w", err)
	}

	return response.Ok, nil
}

func (c *Client) SendCallMessage(ctx context.Context, contractAddress string, args []clarity.ClarityValue, senderAddress string, senderKey []byte) (string, error) {
	tx, err := c.MakeContractCall(
		ctx,
		contractAddress,
		"xcall-proxy",
		"send-call-message",
		args,
		senderAddress,
		senderKey,
	)
	if err != nil {
		return "", err
	}

	txID, err := c.BroadcastTransaction(ctx, tx)
	if err != nil {
		return "", err
	}

	return txID, nil
}

func (c *Client) ExecuteCall(ctx context.Context, contractAddress string, args []clarity.ClarityValue, senderAddress string, senderKey []byte) (string, error) {
	tx, err := c.MakeContractCall(
		ctx,
		contractAddress,
		"xcall-proxy",
		"execute-call",
		args,
		senderAddress,
		senderKey,
	)
	if err != nil {
		return "", err
	}

	txID, err := c.BroadcastTransaction(ctx, tx)
	if err != nil {
		return "", err
	}

	return txID, nil
}

func (c *Client) ExecuteRollback(ctx context.Context, contractAddress string, args []clarity.ClarityValue, senderAddress string, senderKey []byte) (string, error) {
	tx, err := c.MakeContractCall(
		ctx,
		contractAddress,
		"xcall-proxy",
		"execute-rollback",
		args,
		senderAddress,
		senderKey,
	)
	if err != nil {
		return "", err
	}

	txID, err := c.BroadcastTransaction(ctx, tx)
	if err != nil {
		return "", err
	}

	return txID, nil
}

func (c *Client) GetWebSocketURL() string {
	baseURL, _ := c.apiClient.GetConfig().Servers.URL(0, make(map[string]string))

	wsURL := baseURL
	if strings.HasPrefix(wsURL, "http://") {
		wsURL = strings.Replace(wsURL, "http://", "ws://", 1)
	} else if strings.HasPrefix(wsURL, "https://") {
		wsURL = strings.Replace(wsURL, "https://", "wss://", 1)
	}

	if !strings.HasSuffix(wsURL, "/") {
		wsURL += "/"
	}

	wsURL += "extended/v1/ws"

	return wsURL
}
