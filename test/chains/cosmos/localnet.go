package cosmos

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	dockerClient "github.com/docker/docker/client"
	interchaintest "github.com/icon-project/centralized-relay/test"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/blockdb"
	ibcLocal "github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer/centralized"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"github.com/icza/dyno"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"

	"github.com/avast/retry-go/v4"

	"github.com/icon-project/centralized-relay/test/chains"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/dockerutil"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	libclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"go.uber.org/zap"
)

type QueryContractResponse struct {
	Contracts []string `json:"contracts"`
}

type CosmosTx struct {
	TxHash string `json:"txhash"`
	Code   int    `json:"code"`
	RawLog string `json:"raw_log"`
}

type CodeInfo struct {
	CodeID string `json:"code_id"`
}

type CodeInfosResponse struct {
	CodeInfos []CodeInfo `json:"code_infos"`
}

var contracts = chains.ContractKey{
	ContractAddress: make(map[string]string),
	ContractOwner:   make(map[string]string),
}

func NewCosmosRemotenet(testName string, log *zap.Logger, chainConfig ibc.ChainConfig, client *dockerClient.Client, network string, testconfig *testconfig.Chain) (chains.Chain, error) {
	chain := cosmos.NewCosmosChain(testName, chainConfig, 0, 0, log)
	httpClient, err := libclient.DefaultHTTPClient(testconfig.RPCUri)
	if err != nil {
		return nil, err
	}

	httpClient.Timeout = 10 * time.Second
	rpcClient, err := rpchttp.NewWithClient(testconfig.RPCUri, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}
	grpcConn, err := grpc.Dial(
		testconfig.GRPCUri, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %w", err)
	}
	return &CosmosRemotenet{
		CosmosChain:  chain,
		log:          log,
		cfg:          toInterchantestConfig(chain.Config()),
		filepath:     testconfig.Contracts,
		Wallets:      map[string]ibc.Wallet{},
		IBCAddresses: make(map[string]string),
		DockerClient: client,
		testconfig:   testconfig,
		Network:      network,
		testName:     testName,
		Rpcclient:    rpcClient,
		GrpcConn:     grpcConn,
	}, nil
}

func (c *CosmosRemotenet) Start(testName string, ctx context.Context, additionalGenesisWallets ...ibcLocal.WalletAmount) error {
	wallets := []ibc.WalletAmount{}
	for index := range additionalGenesisWallets {
		wallet := ibc.WalletAmount{
			Address: additionalGenesisWallets[index].Address,
			Denom:   "arch",
			Amount:  math.NewInt(additionalGenesisWallets[index].Amount),
		}
		wallets = append(wallets, wallet)
	}

	return c.CosmosChain.Start(testName, ctx, wallets...)
}

func (c *CosmosRemotenet) GetBalance(ctx context.Context, address, denom string) (int64, error) {
	balance, err := c.CosmosChain.GetBalance(ctx, address, denom)
	return balance.Int64(), err
}

func (c *CosmosRemotenet) FindTxs(ctx context.Context, height uint64) ([]blockdb.Tx, error) {
	// return c.CosmosChain.FindTxs(ctx, height)
	return nil, nil
}

func (c *CosmosRemotenet) GetContractAddress(key string) string {
	value, exist := c.IBCAddresses[key]
	if !exist {
		panic(fmt.Sprintf(`IBC address not exist %s`, key))
	}
	return value
}

func (c *CosmosRemotenet) GetRelayConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error) {
	contracts := make(map[string]string)
	contracts["xcall"] = c.GetContractAddress("xcall")
	contracts["connection"] = c.GetContractAddress("connection")
	config := &centralized.CosmosRelayerChainConfig{
		Type: "cosmos",
		Value: centralized.CosmosRelayerChainConfigValue{
			NID:                    c.cfg.ChainID,
			RPCURL:                 c.GetRPCAddress(),
			GrpcUrl:                c.GetGRPCAddress(),
			StartHeight:            0,
			GasPrice:               "900000000000" + c.cfg.Denom,
			GasAdjustment:          1.5,
			BlockInterval:          "6s",
			GasLimit:               2000000,
			MinGasAmount:           200000,
			Contracts:              contracts,
			TxConfirmationInterval: "5s",
			ChainName:              c.Config().Name,
			KeyringBackend:         "memory",
			Address:                c.testconfig.RelayWalletAddress,
			AccountPrefix:          c.cfg.Bech32Prefix,
			Denomination:           c.cfg.Denom,
			ChainID:                c.cfg.ChainID,
			BroadcastMode:          "sync",
			SignModeStr:            "SIGN_MODE_DIRECT",
			MaxGasAmount:           2000000,
			Simulate:               true,
			FinalityBlock:          10,
		},
	}
	return yaml.Marshal(config)
}

func (c *CosmosRemotenet) SetupConnection(ctx context.Context, target chains.Chain) error {
	xcall := c.IBCAddresses["xcall"]
	denom := c.Config().Denom
	connectionCodeId, err := c.StoreContractRemote(ctx, c.filepath["connection"])
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	connectionAddress, err := c.InstantiateContractRemote(ctx, connectionCodeId, `{"denom":"`+denom+`","xcall_address":"`+xcall+`","relayer":"`+c.testconfig.RelayWalletAddress+`"}`, true, c.GetCommonArgs()...)

	if err != nil {
		return err
	}
	c.IBCAddresses["connection"] = connectionAddress
	time.Sleep(5 * time.Second)
	// methodName := "set_fee"
	// _, err = c.ExecuteContract(ctx, connectionAddress, keyName, methodName, map[string]interface{}{
	// 	"network_id":   target.Config().ChainID,
	// 	"message_fee":  "0x0",
	// 	"response_fee": "0x0",
	// },
	// )
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (c *CosmosRemotenet) BuildRelayerWallet(ctx context.Context, keyName string) (ibcLocal.Wallet, error) {
	wallet, err := c.CosmosChain.BuildRelayerWallet(ctx, keyName)
	c.Wallets[keyName] = wallet
	return wallet, err
}

func (c *CosmosRemotenet) InitEventListener(ctx context.Context, contract string) chains.EventListener {
	// listener := NewIconEventListener(c, contract)
	// return listener
	return nil
}

func (c *CosmosRemotenet) SendFunds(ctx context.Context, keyName string, amount ibcLocal.WalletAmount) error {
	amt := ibc.WalletAmount{
		Address: amount.Address,
		Denom:   amount.Denom,
		Amount:  math.NewInt(amount.Amount),
	}
	return c.getFullNode().BankSend(ctx, keyName, amt)
}

func (c *CosmosRemotenet) SetupIBC(ctx context.Context, keyName string) (context.Context, error) {
	var contracts chains.ContractKey
	time.Sleep(4 * time.Second)

	ibcCodeId, err := c.CosmosChain.StoreContract(ctx, keyName, c.filepath["ibc"])
	if err != nil {
		return nil, err
	}

	ibcAddress, err := c.CosmosChain.InstantiateContract(ctx, keyName, ibcCodeId, "{}", true, c.GetCommonArgs()...)
	if err != nil {
		return nil, err
	}

	clientCodeId, err := c.CosmosChain.StoreContract(ctx, keyName, c.filepath["client"])
	if err != nil {
		return ctx, err
	}

	// Parameters here will be empty in the future
	clientAddress, err := c.CosmosChain.InstantiateContract(ctx, keyName, clientCodeId, `{"ibc_host":"`+ibcAddress+`"}`, true, c.GetCommonArgs()...)
	if err != nil {
		return nil, err
	}

	contracts.ContractAddress = map[string]string{
		"ibc":    ibcAddress,
		"client": clientAddress,
	}
	_, err = c.executeContract(context.Background(), ibcAddress, keyName, "register_client", `{"client_type":"iconclient", "client_address":"`+clientAddress+`"}`)
	if err != nil {
		return nil, err
	}
	c.IBCAddresses = contracts.ContractAddress
	overrides := map[string]any{
		"ibc-handler-address": ibcAddress,
		"start-height":        0,
		"block-interval":      "6s",
	}

	cfg := c.cfg
	cfg.ConfigFileOverrides = overrides
	c.cfg = cfg

	return context.WithValue(ctx, chains.Mykey("contract Names"), chains.ContractKey{
		ContractAddress: contracts.ContractAddress,
		ContractOwner:   contracts.ContractOwner,
	}), nil
}

func (c *CosmosRemotenet) SetupXCall(ctx context.Context) error {
	denom := c.Config().Denom
	xCallCodeId, err := c.StoreContractRemote(ctx, c.filepath["xcall"])
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	xCallAddress, err := c.InstantiateContractRemote(ctx, xCallCodeId, `{"network_id": "`+c.Config().ChainID+`", "denom":"`+denom+`"}`, true, c.GetCommonArgs()...)
	if err != nil {
		return err
	}
	c.IBCAddresses["xcall"] = xCallAddress
	return nil
}

func (c *CosmosRemotenet) ConfigureBaseConnection(ctx context.Context, connection chains.XCallConnection) (context.Context, error) {
	panic("not implemented")
}

func (c *CosmosRemotenet) GetIBCAddress(key string) string {
	value, exist := c.IBCAddresses[key]
	if !exist {
		panic(fmt.Sprintf(`IBC address not exist %s`, key))
	}
	return value
}

func (c *CosmosRemotenet) DeployXCallMockApp(ctx context.Context, keyname string, connections []chains.XCallConnection) error {
	testcase := ctx.Value("testcase").(string)
	// connectionKey := fmt.Sprintf("connection-%s", testcase)
	// xCallKey := fmt.Sprintf("xcall-%s", testcase)
	xCall := c.IBCAddresses["xcall"]
	dappCodeId, err := c.StoreContractRemote(ctx, c.filepath["dapp"])
	if err != nil {
		return err
	}
	time.Sleep(8 * time.Second)
	dappAddress, err := c.InstantiateContractRemote(ctx, dappCodeId, `{"address":"`+xCall+`"}`, true, c.GetCommonArgs()...)
	if err != nil {
		return err
	}
	time.Sleep(8 * time.Second)
	c.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = dappAddress
	for _, connection := range connections {
		_, err = c.ExecuteContractRemote(context.Background(), dappAddress, "add_connection", `{"src_endpoint":"`+c.IBCAddresses["connection"]+`", "dest_endpoint":"`+connection.Destination+`","network_id":"`+connection.Nid+`"}`)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CosmosRemotenet) CheckForTimeout(ctx context.Context, target chains.Chain, params map[string]interface{}, listener chains.EventListener) (context.Context, error) {
	panic("not implemented")
}

func (c *CosmosRemotenet) SendPacketXCall(ctx context.Context, keyName, _to string, data, rollback []byte) (context.Context, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)

	dataArray := strings.Join(strings.Fields(fmt.Sprintf("%d", data)), ",")
	rollbackArray := strings.Join(strings.Fields(fmt.Sprintf("%d", rollback)), ",")
	params := fmt.Sprintf(`{"to":"%s", "data":%s, "rollback":%s}`, _to, dataArray, rollbackArray)
	txRes, err := c.ExecuteContractRemote(ctx, c.IBCAddresses[dappKey], "send_call_message", params)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, "sn", c.findSn(txRes, "wasm-CallMessageSent")), nil
}

// FindTargetXCallMessage returns the request id and the data of the message sent to the target chain
func (c *CosmosRemotenet) FindTargetXCallMessage(ctx context.Context, target chains.Chain, height uint64, to string) (*chains.XCallResponse, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	sn := ctx.Value("sn").(string)
	reqId, destData, err := target.FindCallMessage(ctx, height, c.cfg.ChainID+"/"+c.IBCAddresses[dappKey], to, sn)
	return &chains.XCallResponse{SerialNo: sn, RequestID: reqId, Data: destData}, err
}

func (c *CosmosRemotenet) XCall(ctx context.Context, targetChain chains.Chain, keyName, to string, data, rollback []byte) (*chains.XCallResponse, error) {
	height, err := targetChain.(ibcLocal.Chain).Height(ctx)
	if err != nil {
		return nil, err
	}
	ctx, err = c.SendPacketXCall(ctx, keyName, to, data, rollback)
	if err != nil {
		return nil, err
	}
	return c.FindTargetXCallMessage(ctx, targetChain, height, strings.Split(to, "/")[1])
}

func (c *CosmosRemotenet) EOAXCall(ctx context.Context, targetChain chains.Chain, keyName, _to string, data []byte, sources, destinations []string) (string, string, string, error) {
	dataArray := strings.Join(strings.Fields(fmt.Sprintf("%d", data)), ",")
	params := fmt.Sprintf(`{"to":"%s", "data":%s}`, _to, dataArray)
	height, _ := targetChain.(ibcLocal.Chain).Height(ctx)
	ctx, err := c.executeContract(context.Background(), c.IBCAddresses["xcall"], keyName, "send_call_message", params)
	if err != nil {
		return "", "", "", err
	}

	tx := ctx.Value("txResult").(*TxResul)
	sn := c.findSn(tx, "wasm-CallMessageSent")
	reqId, destData, err := targetChain.FindCallMessage(ctx, height, c.cfg.ChainID+"/"+c.IBCAddresses["dapp"], strings.Split(_to, "/")[1], sn)
	return sn, reqId, destData, err
}

func (c *CosmosRemotenet) findSn(tx *TxResul, eType string) string {
	// find better way to parse events
	for _, event := range tx.Events {
		if event.Type == eType {
			for _, attribute := range event.Attributes {
				keyName, _ := base64.StdEncoding.DecodeString(attribute.Key)
				if string(keyName) == "sn" {
					sn, _ := base64.StdEncoding.DecodeString(attribute.Value)
					return string(sn)
				}
			}
		}
	}
	return ""
}

// IsPacketReceived returns the receipt of the packet sent to the target chain
func (c *CosmosRemotenet) IsPacketReceived(ctx context.Context, params map[string]interface{}, order ibc.Order) bool {
	if order == ibc.Ordered {
		sequence := params["sequence"].(uint64)
		ctx, err := c.QueryContract(ctx, c.IBCAddresses["ibc"], chains.GetNextSequenceReceive, params)
		if err != nil {
			fmt.Printf("Error--%v\n", err)
			return false
		}
		response := ctx.Value("query-result").(map[string]interface{})
		fmt.Printf("response[\"data\"]----%v", response["data"])
		return sequence < uint64(response["data"].(float64))
	}
	ctx, err := c.QueryContract(ctx, c.IBCAddresses["ibc"], chains.HasPacketReceipt, params)
	if err != nil {
		fmt.Printf("Error--%v\n", err)
		return false
	}
	response := ctx.Value("query-result").(map[string]interface{})
	return response["data"].(bool)
}

func (c *CosmosRemotenet) ExecuteCall(ctx context.Context, reqId, data string) (context.Context, error) {
	// testcase := ctx.Value("testcase").(string)
	// xCallKey := fmt.Sprintf("xcall-%s", testcase)
	return c.executeContract(ctx, c.IBCAddresses["xcall"], interchaintest.FaucetAccountKeyName, "execute_call", `{"request_id":"`+reqId+`", "data":`+data+`}`)
}

func (c *CosmosRemotenet) ExecuteRollback(ctx context.Context, sn string) (context.Context, error) {
	// testcase := ctx.Value("testcase").(string)
	// xCallKey := fmt.Sprintf("xcall-%s", testcase)
	ctx, err := c.executeContract(context.Background(), c.IBCAddresses["xcall"], interchaintest.FaucetAccountKeyName, "execute_rollback", `{"sequence_no":"`+sn+`"}`)
	if err != nil {
		return nil, err
	}
	tx := ctx.Value("txResult").(*TxResul)
	sequence := c.findSn(tx, "wasm-RollbackExecuted")
	return context.WithValue(ctx, "IsRollbackEventFound", sequence == sn), nil
}

func (c *CosmosRemotenet) FindCallMessage(ctx context.Context, startHeight uint64, from, to, sn string) (string, string, error) {
	// testcase := ctx.Value("testcase").(string)
	xCallKey := "xcall" //fmt.Sprintf("xcall-%s", testcase)
	index := strings.Join([]string{
		fmt.Sprintf("wasm-CallMessage.from CONTAINS '%s'", from),
		fmt.Sprintf("wasm-CallMessage.to CONTAINS '%s'", to),
		fmt.Sprintf("wasm-CallMessage.sn CONTAINS '%s'", sn),
	}, " AND ")
	event, err := c.FindEvent(ctx, startHeight, xCallKey, index)
	if err != nil {
		return "", "", err
	}

	return event.Events["wasm-CallMessage.reqId"][0], event.Events["wasm-CallMessage.data"][0], nil

}

func (c *CosmosRemotenet) FindCallResponse(ctx context.Context, startHeight uint64, sn string) (string, error) {
	// testcase := ctx.Value("testcase").(string)
	xCallKey := "xcall" //fmt.Sprintf("xcall-%s", testcase)
	index := fmt.Sprintf("wasm-ResponseMessage.sn CONTAINS '%s'", sn)
	event, err := c.FindEvent(ctx, startHeight, xCallKey, index)
	if err != nil {
		return "", err
	}

	return event.Events["wasm-ResponseMessage.code"][0], nil
}

func (c *CosmosRemotenet) FindEvent(ctx context.Context, startHeight uint64, contract, index string) (*ctypes.ResultEvent, error) {
	endpoint := c.GetHostRPCAddress()
	client, err := rpchttp.New(endpoint, "/websocket")
	if err != nil {
		return nil, err
	}

	err = client.Start()
	if err != nil {
		return nil, err
	}
	defer client.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	query := strings.Join([]string{"tm.event = 'Tx'",
		fmt.Sprintf("tx.height >= %d ", startHeight),
		"message.module = 'wasm'",
		fmt.Sprintf("wasm._contract_address = '%s'", c.IBCAddresses["xcall"]),
		index,
	}, " AND ")
	channel, err := client.Subscribe(ctx, "wasm-client", query)
	if err != nil {
		fmt.Println("error subscribint to channel")
		return nil, err
	}

	select {
	case event := <-channel:
		return &event, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("failed to find eventLog")
	}
}

func (c *CosmosRemotenet) DeployContract(ctx context.Context, keyName string) (context.Context, error) {
	// Fund user to deploy contract
	wallet, _ := c.GetAndFundTestUser(ctx, keyName, "", int64(100_000_000))

	contractOwner := wallet.FormattedAddress()
	// Get contract Name from context
	ctxValue := ctx.Value(chains.ContractName{}).(chains.ContractName)
	initMsg := ctx.Value(chains.InitMessageKey("init-msg")).(chains.InitMessage)

	contractName := strings.ToLower(ctxValue.ContractName)
	codeId, err := c.CosmosChain.StoreContract(ctx, contractOwner, c.filepath[contractName])
	if err != nil {
		return ctx, err
	}

	initMessage := c.getInitParams(ctx, contractName, initMsg.Message)
	address, err := c.CosmosChain.InstantiateContract(ctx, contractOwner, codeId, initMessage, true, c.GetCommonArgs()...)
	if err != nil {
		return nil, err
	}

	testcase := ctx.Value("testcase").(string)
	contract := fmt.Sprintf("%s-%s", contractName, testcase)
	c.IBCAddresses[contract] = address

	return context.WithValue(ctx, chains.Mykey("contract Names"), contracts), err
}

func (c *CosmosRemotenet) QueryContract(ctx context.Context, contractAddress, methodName string, params map[string]interface{}) (context.Context, error) {
	// wait for few blocks after executing before querying
	time.Sleep(2 * time.Second)

	// get query msg
	query := c.GetQueryParam(methodName, params)
	chains.Response = ""
	err := c.CosmosChain.QueryContract(ctx, contractAddress, query, &chains.Response)
	fmt.Printf("Response is : %s \n", chains.Response)
	return context.WithValue(ctx, "query-result", chains.Response), err
	//return context.WithValue(ctx, "txResult", chains.Response.(map[string]interface{})["data"]), nil

}

func (c *CosmosRemotenet) executeContract(ctx context.Context, contractAddress, keyName, methodName, param string) (context.Context, error) {
	txHash, err := c.ExecTx(ctx,
		"wasm", "execute", contractAddress, `{"`+methodName+`":`+param+`}`, "--gas", "auto")
	if err != nil || txHash == "" {
		return ctx, err
	}
	tx, err := c.getTransaction(txHash)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, "txResult", tx), nil
}

func (c *CosmosRemotenet) ExecuteContract(ctx context.Context, contractAddress, keyName, methodName string, params map[string]interface{}) (context.Context, error) {
	execMethodName, execParams := c.getExecuteParam(ctx, methodName, params)
	return c.executeContract(ctx, contractAddress, keyName, execMethodName, execParams)
}

func (c *CosmosRemotenet) getTransaction(txHash string) (*TxResul, error) {
	// Retry because sometimes the tx is not committed to state yet.
	var result TxResul

	err := retry.Do(func() error {
		var err error
		stdout, _, err := c.ExecQuery(context.Background(), "tx", txHash)
		err = json.Unmarshal(stdout, &result)
		return err
	},
		// retry for total of 6 seconds
		retry.Attempts(30),
		retry.Delay(200*time.Millisecond),
		retry.DelayType(retry.FixedDelay),
		retry.LastErrorOnly(true),
	)
	return &result, err
}

func (c *CosmosRemotenet) getFullNode() *cosmos.ChainNode {
	panic("not implemented")
}

func (c *CosmosRemotenet) GetLastBlock(ctx context.Context) (context.Context, error) {
	h, err := c.CosmosChain.Height(ctx)
	return context.WithValue(ctx, chains.LastBlock{}, h), err
}

func (c *CosmosRemotenet) GetBlockByHeight(ctx context.Context) (context.Context, error) {
	panic("not implemented") // TODO: Implement
}

func (c *CosmosRemotenet) BuildWallets(ctx context.Context, keyName string) (ibcLocal.Wallet, error) {
	return c.GetAndFundTestUser(ctx, keyName, "", int64(100_000_000))
}

func (c *CosmosRemotenet) BuildWallet(ctx context.Context, keyName string, mnemonic string) (ibcLocal.Wallet, error) {
	if mnemonic == "" {
		wallet, _ := c.BuildRelayerWallet(ctx, keyName)
		mnemonic = wallet.Mnemonic()
	}

	if err := c.CosmosChain.RecoverKey(ctx, keyName, mnemonic); err != nil {
		return nil, fmt.Errorf("failed to recover key with name %q on chain %s: %w", keyName, c.cfg.Name, err)
	}

	addrBytes, err := c.CosmosChain.GetAddress(ctx, keyName)
	if err != nil {
		return nil, fmt.Errorf("failed to get account address for key %q on chain %s: %w", keyName, c.cfg.Name, err)
	}
	wallet := cosmos.NewWallet(keyName, addrBytes, mnemonic, toIbcConfig(c.cfg))
	c.Wallets[keyName] = wallet
	return wallet, nil
}

func (c *CosmosRemotenet) GetCommonArgs() []string {
	return []string{"--gas", "auto"}
}

func (c *CosmosRemotenet) GetClientName(suffix int) string {
	return fmt.Sprintf("iconclient-%d", suffix)
}

// GetClientsCount returns the next sequence number for the client
func (c *CosmosRemotenet) GetClientsCount(ctx context.Context) (int, error) {
	var err error
	ctx, err = c.QueryContract(ctx, c.GetIBCAddress("ibc"), chains.GetNextClientSequence, map[string]interface{}{})
	if err != nil {
		return 0, err
	}
	res := ctx.Value("query-result").(map[string]interface{})
	var data = res["data"].(float64)
	return int(data), nil
}

// GetNextConnectionSequence returns the next sequence number for the client
func (c *CosmosRemotenet) GetNextConnectionSequence(ctx context.Context) (int, error) {
	params := map[string]interface{}{}
	var err error
	ctx, err = c.QueryContract(ctx, c.GetIBCAddress("ibc"), chains.GetNextConnectionSequence, params)
	if err != nil {
		return 0, err
	}
	res := ctx.Value("query-result").(map[string]interface{})

	count := res["data"].(float64)
	return int(count), err
}

// GetNextChannelSequence returns the next sequence number for the client
func (c *CosmosRemotenet) GetNextChannelSequence(ctx context.Context) (int, error) {
	params := map[string]interface{}{}
	var err error
	ctx, err = c.QueryContract(ctx, c.GetIBCAddress("ibc"), chains.GetNextChannelSequence, params)
	if err != nil {
		return 0, err
	}
	res := ctx.Value("query-result").(map[string]interface{})

	count := res["data"].(float64)
	return int(count), err
}

// PauseNode halts a node
func (c *CosmosRemotenet) PauseNode(ctx context.Context) error {
	return nil
}

// UnpauseNode restarts a node
func (c *CosmosRemotenet) UnpauseNode(ctx context.Context) error {
	return nil
}

func (c *CosmosRemotenet) BackupConfig() ([]byte, error) {
	wallets := make(map[string]interface{})
	for key, value := range c.Wallets {
		wallets[key] = map[string]string{
			"mnemonic":         value.Mnemonic(),
			"address":          hex.EncodeToString(value.Address()),
			"formattedAddress": value.FormattedAddress(),
		}
	}
	backup := map[string]interface{}{
		"addresses": c.IBCAddresses,
		"wallets":   wallets,
	}
	return json.MarshalIndent(backup, "", "\t")
}

func (c *CosmosRemotenet) RestoreConfig(backup []byte) error {
	result := make(map[string]interface{})
	err := json.Unmarshal(backup, &result)
	if err != nil {
		return err
	}
	fmt.Println("Restoring config")
	c.IBCAddresses = result["addresses"].(map[string]string)
	wallets := make(map[string]ibc.Wallet)
	for key, value := range result["wallets"].(map[string]interface{}) {
		_value := value.(map[string]string)
		mnemonic := _value["mnemonic"]
		address, _ := hex.DecodeString(_value["address"])
		wallets[key] = cosmos.NewWallet(key, address, mnemonic, toIbcConfig(c.Config()))
	}
	c.Wallets = wallets
	return nil
}

func toIbcConfig(config ibcLocal.ChainConfig) ibc.ChainConfig {
	images := []ibc.DockerImage{
		ibc.NewDockerImage(
			config.Images[0].Repository, config.Images[0].Version, config.Images[0].UidGid),
	}
	decimals := int64(6)
	return ibc.ChainConfig{
		Type:           config.Type,
		Name:           config.Name,
		ChainID:        config.ChainID,
		Bin:            config.Bin,
		Bech32Prefix:   config.Bech32Prefix,
		Denom:          config.Denom,
		Images:         images,
		SkipGenTx:      config.SkipGenTx,
		CoinType:       config.CoinType,
		GasPrices:      config.GasPrices,
		GasAdjustment:  config.GasAdjustment,
		TrustingPeriod: config.TrustingPeriod,
		NoHostMount:    config.NoHostMount,
		CoinDecimals:   &decimals,
		ModifyGenesis: func(config ibc.ChainConfig, bytes []byte) ([]byte, error) {
			if config.Type == "icon" {
				return bytes, nil
			}

			g := make(map[string]interface{})
			if err := json.Unmarshal(bytes, &g); err != nil {
				return nil, fmt.Errorf("failed to unmarshal genesis file: %w", err)
			}

			if config.SkipGenTx {
				//add minimum gas fee
				if err := modifyGenesisMinGasPrice(g, config); err != nil {
					return nil, err
				}
			}

			out, err := json.Marshal(&g)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal genesis bytes to json: %w", err)
			}
			return out, nil

		},
	}
}

func toInterchantestConfig(config ibc.ChainConfig) ibcLocal.ChainConfig {
	images := []ibcLocal.DockerImage{
		{
			Repository: config.Images[0].Repository,
			Version:    config.Images[0].Version,
			UidGid:     config.Images[0].UidGid,
		},
	}
	return ibcLocal.ChainConfig{
		Type:           config.Type,
		Name:           config.Name,
		ChainID:        config.ChainID,
		Bin:            config.Bin,
		Bech32Prefix:   config.Bech32Prefix,
		Denom:          config.Denom,
		SkipGenTx:      config.SkipGenTx,
		CoinType:       config.CoinType,
		GasPrices:      config.GasPrices,
		GasAdjustment:  config.GasAdjustment,
		TrustingPeriod: config.TrustingPeriod,
		NoHostMount:    config.NoHostMount,
		Images:         images,
		ModifyGenesis: func(config ibcLocal.ChainConfig, bytes []byte) ([]byte, error) {
			if config.Type == "icon" {
				return bytes, nil
			}

			g := make(map[string]interface{})
			if err := json.Unmarshal(bytes, &g); err != nil {
				return nil, fmt.Errorf("failed to unmarshal genesis file: %w", err)
			}

			if config.SkipGenTx {
				//add minimum gas fee
				if err := modifyGenesisMinGasPriceLocal(g, config); err != nil {
					return nil, err
				}
			}

			out, err := json.Marshal(&g)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal genesis bytes to json: %w", err)
			}
			return out, nil

		},
	}
}

type MinimumGasPriceEntity struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

func modifyGenesisMinGasPrice(g map[string]interface{}, config ibc.ChainConfig) error {
	fmt.Println("modify Genesis config")
	minGasPriceEntities := []MinimumGasPriceEntity{
		{
			Denom:  config.Denom,
			Amount: "0",
		},
	}
	if err := dyno.Set(g, minGasPriceEntities, "app_state", "globalfee", "params", "minimum_gas_prices"); err != nil {
		return fmt.Errorf("failed to set params minimum gas price in genesis json: %w", err)
	}

	return nil
}

func modifyGenesisMinGasPriceLocal(g map[string]interface{}, config ibcLocal.ChainConfig) error {
	fmt.Println("modify local config")
	minGasPriceEntities := []MinimumGasPriceEntity{
		{
			Denom:  config.Denom,
			Amount: "0",
		},
	}
	if err := dyno.Set(g, minGasPriceEntities, "app_state", "globalfee", "params", "minimum_gas_prices"); err != nil {
		return fmt.Errorf("failed to set params minimum gas price in genesis json: %w", err)
	}

	return nil
}

func (an *CosmosRemotenet) Exec(ctx context.Context, cmd []string, env []string) (stdout []byte, stderr []byte, err error) {
	job := dockerutil.NewImage(an.log, an.DockerClient, an.Network, an.testName, an.cfg.Images[0].Repository, an.cfg.Images[0].Version)
	bindPaths := []string{
		an.testconfig.ContractsPath + ":/contracts",
		an.testconfig.ConfigPath + ":/root/.archway",
	}
	if an.testconfig.CertPath != "" {
		bindPaths = append(bindPaths, an.testconfig.CertPath+":/etc/ssl/certs/")
	}
	opts := dockerutil.ContainerOptions{
		Binds: bindPaths,
	}
	res := job.Run(ctx, cmd, opts)
	return res.Stdout, res.Stderr, res.Err
}

// func (c *CosmosLocalnet) Height(ctx context.Context) (uint64, error) {
// 	res, err := c.Client.Status(ctx)
// 	if err != nil {
// 		return 0, fmt.Errorf("tendermint rpc client status: %w", err)
// 	}
// 	height := res.SyncInfo.LatestBlockHeight
// 	return uint64(height), nil
// }

func (c *CosmosRemotenet) StoreContractRemote(ctx context.Context, fileName string, extraExecTxArgs ...string) (string, error) {
	_, file := filepath.Split(fileName)

	cmd := []string{"wasm", "store", path.Join("/contracts/", file), "--gas", "auto"}
	cmd = append(cmd, extraExecTxArgs...)
	// ht, _ := c.Height(ctx)
	// fmt.Println("Current height is ", ht)
	if _, err := c.ExecTx(ctx, cmd...); err != nil {
		return "", err
	}

	// err := testutil.WaitForBlocks(ctx, 5, c.CosmosChain)
	// if err != nil {
	// 	return "", fmt.Errorf("wait for blocks: %w", err)
	// }

	time.Sleep(9 * time.Second)
	// ht, _ = c.Height(ctx)
	// fmt.Println("Current New height is ", ht)
	stdout, _, err := c.ExecQuery(ctx, "wasm", "list-code", "--reverse")
	if err != nil {
		return "", err
	}
	res := CodeInfosResponse{}
	if err := json.Unmarshal([]byte(stdout), &res); err != nil {
		return "", err
	}

	return res.CodeInfos[0].CodeID, nil
}

// InstantiateContract takes a code id for a smart contract and initialization message and returns the instantiated contract address.
func (c *CosmosRemotenet) InstantiateContractRemote(ctx context.Context, codeID string, initMessage string, needsNoAdminFlag bool, extraExecTxArgs ...string) (string, error) {
	command := []string{"wasm", "instantiate", codeID, initMessage, "--label", "wasm-contract"}
	command = append(command, extraExecTxArgs...)
	if needsNoAdminFlag {
		command = append(command, "--no-admin")
	}
	_, err := c.ExecTx(ctx, command...)
	if err != nil {
		return "", err
	}

	// txResp, err := c.GetTransaction( txHash)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to get transaction %s: %w", txHash, err)
	// }
	// if txResp.Code != 0 {
	// 	return "", fmt.Errorf("error in transaction (code: %d): %s", txResp.Code, txResp.RawLog)
	// }
	time.Sleep(8 * time.Second)

	stdout, _, err := c.ExecQuery(ctx, "wasm", "list-contract-by-code", codeID)
	if err != nil {
		return "", err
	}
	contactsRes := QueryContractResponse{}
	if err := json.Unmarshal([]byte(stdout), &contactsRes); err != nil {
		return "", err
	}
	contractAddress := contactsRes.Contracts[len(contactsRes.Contracts)-1]
	return contractAddress, nil
}

// ExecuteContract executes a contract transaction with a message using it's address.
func (c *CosmosRemotenet) ExecuteContractRemote(ctx context.Context, contractAddress string, methodName, message string, extraExecTxArgs ...string) (res *TxResul, err error) {
	msg := `{"` + methodName + `":` + message + `}`
	cmd := []string{"wasm", "execute", contractAddress, msg}
	cmd = append(cmd, extraExecTxArgs...)
	if len(extraExecTxArgs) == 0 {
		cmd = append(cmd, "--gas", "auto")
	}

	txHash, err := c.ExecTx(ctx, cmd...)
	if err != nil {
		return nil, err
	}

	txResp, err := c.getTransaction(txHash)
	if err != nil {
		return txResp, fmt.Errorf("failed to get transaction %s: %w", txHash, err)
	}

	if txResp.Code != 0 {
		return txResp, fmt.Errorf("error in transaction (code: %d): %s", txResp.Code, txResp.RawLog)
	}

	return txResp, nil
}

// QueryContract performs a smart query, taking in a query struct and returning a error with the response struct populated.
func (c *CosmosRemotenet) QueryContractRemote(ctx context.Context, contractAddress string, queryMsg any, response any) error {
	var query []byte
	var err error

	if q, ok := queryMsg.(string); ok {
		var jsonMap map[string]interface{}
		if err := json.Unmarshal([]byte(q), &jsonMap); err != nil {
			return err
		}

		query, err = json.Marshal(jsonMap)
		if err != nil {
			return err
		}
	} else {
		query, err = json.Marshal(queryMsg)
		if err != nil {
			return err
		}
	}

	stdout, _, err := c.ExecQuery(ctx, "wasm", "contract-state", "smart", contractAddress, string(query))
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(stdout), response)
	return err
}

func (c *CosmosRemotenet) BinCommand(command ...string) []string {
	command = append([]string{c.CosmosChain.Config().Bin}, command...)
	return command
}

// ExecBin is a helper to execute a command for a chain node binary.
// For example, if chain node binary is `gaiad`, and desired command is `gaiad keys show key1`,
// pass ("keys", "show", "key1") for command to execute the command against the node.
// Will include additional flags for home directory and chain ID.
func (c *CosmosRemotenet) ExecBin(ctx context.Context, command ...string) ([]byte, []byte, error) {
	return c.Exec(ctx, c.BinCommand(command...), nil)
}

// QueryCommand is a helper to retrieve the full query command. For example,
// if chain node binary is gaiad, and desired command is `gaiad query gov params`,
// pass ("gov", "params") for command to return the full command with all necessary
// flags to query the specific node.
func (c *CosmosRemotenet) QueryCommand(command ...string) []string {
	command = append([]string{"query"}, command...)
	return c.NodeCommand(append(command,
		"--output", "json",
	)...)
}

// ExecQuery is a helper to execute a query command. For example,
// if chain node binary is gaiad, and desired command is `gaiad query gov params`,
// pass ("gov", "params") for command to execute the query against the node.
// Returns response in json format.
func (c *CosmosRemotenet) ExecQuery(ctx context.Context, command ...string) ([]byte, []byte, error) {
	return c.Exec(ctx, c.QueryCommand(command...), nil)
}

func (c *CosmosRemotenet) NodeCommand(command ...string) []string {
	command = c.BinCommand(command...)
	return append(command,
		"--node", c.testconfig.RPCUri,
	)
}

func (c *CosmosRemotenet) ExecTx(ctx context.Context, command ...string) (string, error) {

	stdout, _, err := c.Exec(ctx, c.TxCommand(command...), nil)
	if err != nil {
		return "", err
	}
	output := CosmosTx{}
	err = json.Unmarshal([]byte(stdout), &output)
	if err != nil {
		return "", err
	}
	if output.Code != 0 {
		return output.TxHash, fmt.Errorf("transaction failed with code %d: %s", output.Code, output.RawLog)
	}
	// if err := testutil.WaitForBlocks(ctx, 2, tn); err != nil {
	// 	return "", err
	// }
	return output.TxHash, nil
}

// TxHashToResponse returns the sdk transaction response struct for a given transaction hash.
func (c *CosmosRemotenet) TxHashToResponse(ctx context.Context, txHash string) (*sdk.TxResponse, error) {
	stdout, stderr, err := c.ExecQuery(ctx, "tx", txHash)
	if err != nil {
		fmt.Println("TxHashToResponse err: ", err.Error()+" "+string(stderr))
	}

	i := &sdk.TxResponse{}

	// ignore the error since some types do not unmarshal (ex: height of int64 vs string)
	_ = json.Unmarshal(stdout, &i)
	return i, nil
}

func (c *CosmosRemotenet) TxCommand(command ...string) []string {
	command = append([]string{"tx"}, command...)
	var gasPriceFound, gasAdjustmentFound, feesFound = false, false, false
	for i := 0; i < len(command); i++ {
		if command[i] == "--gas-prices" {
			gasPriceFound = true
		}
		if command[i] == "--gas-adjustment" {
			gasAdjustmentFound = true
		}
		if command[i] == "--fees" {
			feesFound = true
		}
	}
	if !gasPriceFound && !feesFound {
		command = append(command, "--gas-prices", c.CosmosChain.Config().GasPrices)
	}
	if !gasAdjustmentFound {
		command = append(command, "--gas-adjustment", strconv.FormatFloat(c.CosmosChain.Config().GasAdjustment, 'f', -1, 64))
	}
	return c.NodeCommand(append(command,
		"--from", c.testconfig.KeystoreFile,
		"--keyring-backend", keyring.BackendTest,
		"--output", "json",
		"-y",
		"--chain-id", c.CosmosChain.Config().ChainID,
	)...)
}

func (c *CosmosRemotenet) CliContext() client.Context {
	cfg := c.CosmosChain.Config()
	return client.Context{
		Client:            nil,
		GRPCClient:        nil,
		ChainID:           cfg.ChainID,
		InterfaceRegistry: cfg.EncodingConfig.InterfaceRegistry,
		Input:             os.Stdin,
		Output:            os.Stdout,
		OutputFormat:      "json",
		LegacyAmino:       cfg.EncodingConfig.Amino,
		TxConfig:          cfg.EncodingConfig.TxConfig,
	}
}

func (c *CosmosRemotenet) GetRPCAddress() string {
	return c.testconfig.RPCUri
}

// Implements Chain interface
func (c *CosmosRemotenet) GetAPIAddress() string {
	return c.testconfig.RPCUri
}

// Implements Chain interface
func (c *CosmosRemotenet) GetGRPCAddress() string {
	return c.testconfig.GRPCUri
}

// GetHostRPCAddress returns the address of the RPC server accessible by the host.
// This will not return a valid address until the chain has been started.
func (c *CosmosRemotenet) GetHostRPCAddress() string {
	return c.testconfig.RPCUri
}

func (c *CosmosRemotenet) Height(ctx context.Context) (uint64, error) {
	res, err := c.Rpcclient.Status(ctx)
	if err != nil {
		return 0, fmt.Errorf("tendermint rpc client status: %w", err)
	}
	height := res.SyncInfo.LatestBlockHeight
	return uint64(height), nil
}
