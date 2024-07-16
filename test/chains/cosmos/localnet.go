package cosmos

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	dockerClient "github.com/docker/docker/client"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer/centralized"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"

	"github.com/avast/retry-go/v4"

	"github.com/icon-project/centralized-relay/test/chains"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/dockerutil"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	libclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	ibcv8 "github.com/strangelove-ventures/interchaintest/v8/ibc"
	"go.uber.org/zap"
)

func NewCosmosRemotenet(testName string, log *zap.Logger, chainConfig chains.ChainConfig, client *dockerClient.Client, network string, testconfig *testconfig.Chain) (chains.Chain, error) {
	// chain := cosmos.NewCosmosChain(testName, chainConfig, 0, 0, log)
	httpClient, err := libclient.DefaultHTTPClient(testconfig.RPCUri)
	if err != nil {
		return nil, err
	}

	httpClient.Timeout = 10 * time.Second
	rpcClient, err := rpchttp.NewWithClient(testconfig.RPCUri, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}
	grpcConn, err := grpc.NewClient(
		testconfig.GRPCUri, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %w", err)
	}
	return &CosmosRemotenet{
		log:          log,
		cfg:          chainConfig,
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
	time.Sleep(5 * time.Second)
	connectionAddress, err := c.InstantiateContractRemote(ctx, connectionCodeId, `{"denom":"`+denom+`","xcall_address":"`+xcall+`","relayer":"`+c.testconfig.RelayWalletAddress+`"}`, true, c.GetCommonArgs()...)

	if err != nil {
		return err
	}
	c.IBCAddresses["connection"] = connectionAddress
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

func (c *CosmosRemotenet) SendPacketXCall(ctx context.Context, keyName, _to string, data, rollback []byte) (context.Context, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)

	dataArray := strings.Join(strings.Fields(fmt.Sprintf("%d", data)), ",")
	params := fmt.Sprintf(`{"to":"%s", "data":%s}`, _to, dataArray)
	if rollback != nil {
		rollbackArray := strings.Join(strings.Fields(fmt.Sprintf("%d", rollback)), ",")
		params = fmt.Sprintf(`{"to":"%s", "data":%s, "rollback":%s}`, _to, dataArray, rollbackArray)
	}
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
	height, err := targetChain.Height(ctx)
	if err != nil {
		return nil, err
	}
	ctx, err = c.SendPacketXCall(ctx, keyName, to, data, rollback)
	if err != nil {
		return nil, err
	}
	return c.FindTargetXCallMessage(ctx, targetChain, height, strings.Split(to, "/")[1])
}

func (c *CosmosRemotenet) findSn(tx *TxResul, eType string) string {
	// find better way to parse events
	for _, event := range tx.Events {
		if event.Type == eType {
			for _, attribute := range event.Attributes {
				keyName, _ := base64.StdEncoding.DecodeString(attribute.Key)
				if attribute.Key == "sn" {
					return attribute.Value
				}
				if string(keyName) == "sn" {
					sn, _ := base64.StdEncoding.DecodeString(attribute.Value)
					return string(sn)

				}
			}
		}
	}
	return ""
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

func (c *CosmosRemotenet) getTransaction(txHash string) (*TxResul, error) {
	// Retry because sometimes the tx is not committed to state yet.
	var result TxResul

	err := retry.Do(func() error {
		var err error
		stdout, _, _ := c.ExecQuery(context.Background(), "tx", txHash)
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

func (c *CosmosRemotenet) GetCommonArgs() []string {
	return []string{"--gas", "auto"}
}

func toInterchantestConfig(config chains.ChainConfig) ibcv8.ChainConfig {
	images := []ibcv8.DockerImage{
		{
			Repository: config.Images.Repository,
			Version:    config.Images.Version,
			UidGid:     config.Images.UidGid,
		},
	}
	return ibcv8.ChainConfig{
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
	}
}

type MinimumGasPriceEntity struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

func (an *CosmosRemotenet) Exec(ctx context.Context, cmd []string, env []string) (stdout []byte, stderr []byte, err error) {
	job := dockerutil.NewImage(an.log, an.DockerClient, an.Network, an.testName, an.cfg.Images.Repository, an.cfg.Images.Version)
	bindPaths := []string{
		an.testconfig.ContractsPath + ":/contracts",
		an.testconfig.ConfigPath + ":/root/.archway",
	}

	opts := dockerutil.ContainerOptions{
		Binds: bindPaths,
	}
	res := job.Run(ctx, cmd, opts)
	return res.Stdout, res.Stderr, res.Err
}

func (c *CosmosRemotenet) StoreContractRemote(ctx context.Context, fileName string, extraExecTxArgs ...string) (string, error) {
	_, file := filepath.Split(fileName)

	cmd := []string{"wasm", "store", path.Join("/contracts/", file), "--gas", "auto"}
	cmd = append(cmd, extraExecTxArgs...)
	if _, err := c.ExecTx(ctx, cmd...); err != nil {
		return "", err
	}

	time.Sleep(9 * time.Second)
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
	command = append([]string{c.Config().Bin}, command...)
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
		command = append(command, "--gas-prices", c.Config().GasPrices)
	}
	if !gasAdjustmentFound {
		command = append(command, "--gas-adjustment", strconv.FormatFloat(c.Config().GasAdjustment, 'f', -1, 64))
	}
	return c.NodeCommand(append(command,
		"--from", c.testconfig.KeystoreFile,
		"--keyring-backend", keyring.BackendTest,
		"--output", "json",
		"-y",
		"--chain-id", c.Config().ChainID,
	)...)
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

func (c *CosmosRemotenet) FindRollbackExecutedMessage(ctx context.Context, startHeight uint64, sn string) (string, error) {
	xCallKey := "xcall" //fmt.Sprintf("xcall-%s", testcase)
	index := fmt.Sprintf("wasm-RollbackExecuted.sn CONTAINS '%s'", sn)
	_, err := c.FindEvent(ctx, startHeight, xCallKey, index)
	if err != nil {
		return "", err
	}

	return "0", nil
}
