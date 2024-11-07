package icon

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/dockerutil"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer/centralized"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	iconlog "github.com/icon-project/icon-bridge/common/log"
	"gopkg.in/yaml.v3"

	//chantypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/icon-project/centralized-relay/test/chains"
	iconclient "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon"
	icontypes "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon/types"

	"go.uber.org/zap"
)

type IconRemotenet struct {
	log           *zap.Logger
	testName      string
	cfg           chains.ChainConfig
	numValidators int
	numFullNodes  int
	keystorePath  string
	scorePaths    map[string]string
	IBCAddresses  map[string]string     `json:"addresses"`
	Wallets       map[string]ibc.Wallet `json:"wallets"`
	Client        *client.Client
	Network       string
	testconfig    *testconfig.Chain
	IconClient    iconclient.Client
}

func NewIconRemotenet(testName string, log *zap.Logger, chainConfig chains.ChainConfig, client *client.Client, network string, testconfig *testconfig.Chain) chains.Chain {
	uri := testconfig.RPCUri
	var l iconlog.Logger
	return &IconRemotenet{
		testName:      testName,
		cfg:           chainConfig,
		numValidators: 0,
		numFullNodes:  0,
		log:           log,
		scorePaths:    testconfig.Contracts,
		Wallets:       map[string]ibc.Wallet{},
		IBCAddresses:  make(map[string]string),
		Client:        client,
		testconfig:    testconfig,
		Network:       network,
		IconClient:    *iconclient.NewClient(uri, l),
	}
}

// Config fetches the chain configuration.
func (in *IconRemotenet) Config() chains.ChainConfig {
	return in.cfg
}

// Exec runs an arbitrary command using Chain's docker environment.
// Whether the invoked command is run in a one-off container or execing into an already running container
// is up to the chain implementation.
//
// "env" are environment variables in the format "MY_ENV_VAR=value"
func (in *IconRemotenet) Exec(ctx context.Context, cmd []string, env []string) (stdout []byte, stderr []byte, err error) {
	cmd = append([]string{}, cmd...)

	job := dockerutil.NewImage(in.log, in.Client, in.Network, in.testName, in.cfg.Images.Repository, in.cfg.Images.Version)
	var ContainerEnvs = [9]string{
		"GOCHAIN_CONFIG=/goloop/data/config.json",
		"GOCHAIN_GENESIS=/goloop/data/genesis.json",
		"GOCHAIN_DATA=/goloop/chain/iconee",
		"GOCHAIN_LOGFILE=/goloop/chain/iconee.log",
		"GOCHAIN_DB_TYPE=rocksdb",
		"GOCHAIN_CLEAN_DATA=true",
		"JAVAEE_BIN=/goloop/execman/bin/execman",
		"PYEE_VERIFY_PACKAGE=true",
		"ICON_CONFIG=/goloop/data/icon_config.json",
	}
	bindPaths := []string{
		in.testconfig.ContractsPath + ":/contracts",
		in.testconfig.ConfigPath + ":/goloop/data",
	}
	if in.testconfig.CertPath != "" {
		bindPaths = append(bindPaths, in.testconfig.CertPath+":/etc/ssl/certs/")
	}
	opts := dockerutil.ContainerOptions{
		Binds: bindPaths,
		Env:   ContainerEnvs[:],
	}
	res := job.Run(ctx, cmd, opts)
	return res.Stdout, res.Stderr, res.Err
}

// GetRPCAddress retrieves the rpc address that can be reached by other containers in the docker network.
func (in *IconRemotenet) GetRPCAddress() string {
	return in.testconfig.RPCUri
}

func (in *IconRemotenet) GetRelayConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error) {
	contracts := make(map[string]string)
	contracts["xcall"] = in.GetContractAddress("xcall")
	contracts["connection"] = in.GetContractAddress("connection")
	config := &centralized.ICONRelayerChainConfig{
		Type: "icon",
		Value: centralized.ICONRelayerChainConfigValue{
			NID:           in.Config().ChainID,
			RPCURL:        in.GetRPCAddress(),
			StartHeight:   0,
			NetworkID:     0x3,
			Contracts:     contracts,
			BlockInterval: "6s",
			Address:       in.testconfig.RelayWalletAddress,
			FinalityBlock: uint64(10),
			StepMin:       25000,
			StepLimit:     2500000,
		},
	}
	return yaml.Marshal(config)
}

// GetGRPCAddress retrieves the grpc address that can be reached by other containers in the docker network.
// Not Applicable for Icon
func (in *IconRemotenet) GetGRPCAddress() string {
	return in.testconfig.RPCUri
}

// GetHostRPCAddress returns the rpc address that can be reached by processes on the host machine.
// Note that this will not return a valid value until after Start returns.
func (in *IconRemotenet) GetHostRPCAddress() string {
	return in.testconfig.RPCUri
}

// GetHostGRPCAddress returns the grpc address that can be reached by processes on the host machine.
// Note that this will not return a valid value until after Start returns.
// Not applicable for Icon
func (in *IconRemotenet) GetHostGRPCAddress() string {
	return in.testconfig.RPCUri
}

// HomeDir is the home directory of a node running in a docker container. Therefore, this maps to
// the container's filesystem (not the host).
func (in *IconRemotenet) HomeDir() string {
	return ""
}

// Height returns the current block height or an error if unable to get current height.
func (in *IconRemotenet) Height(ctx context.Context) (uint64, error) {
	res, err := in.IconClient.GetLastBlock()
	return uint64(res.Height), err
}

// GetGasFeesInNativeDenom gets the fees in native denom for an amount of spent gas.
func (in *IconRemotenet) GetGasFeesInNativeDenom(gasPaid int64) int64 {
	gasPrice, _ := strconv.ParseFloat(strings.Replace(in.cfg.GasPrices, in.cfg.Denom, "", 1), 64)
	fees := float64(gasPaid) * gasPrice
	return int64(fees)
}

func (in *IconRemotenet) SetupConnection(ctx context.Context, target chains.Chain) error {
	if in.testconfig.Environment == "preconfigured" {
		return nil
	}
	xcall := in.IBCAddresses["xcall"]

	connection, err := in.DeployContractRemote(ctx, in.scorePaths["connection"], in.keystorePath, `{"_xCall":"`+xcall+`","_relayer":"`+in.testconfig.RelayWalletAddress+`"}`)
	if err != nil {
		return err
	}

	params := `{"networkId":"` + target.Config().ChainID + `", "messageFee":"0x0", "responseFee":"0x0"}`
	_, err = in.executeContract(context.Background(), connection, "setFee", params)
	if err != nil {
		return err
	}
	in.IBCAddresses["connection"] = connection
	return nil
}

func (in *IconRemotenet) SetupXCall(ctx context.Context) error {
	if in.testconfig.Environment == "preconfigured" {
		testcase := ctx.Value("testcase").(string)
		in.IBCAddresses["xcall"] = "cx07c425640148e54b07cd92ad0835d83b64ad25a8"
		in.IBCAddresses["connection"] = "cx97b13c438086271a61f611f396e3834b381cbb91"
		in.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = "cxf18fc12e7fc1d34deb9e85f318640d433ddbc01f"
		return nil
	}
	nid := in.cfg.ChainID
	xcall, err := in.DeployContractRemote(ctx, in.scorePaths["xcall"], in.keystorePath, `{"networkId":"`+nid+`"}`)
	if err != nil {
		return err
	}

	in.IBCAddresses["xcall"] = xcall
	return nil
}

func (in *IconRemotenet) DeployXCallMockApp(ctx context.Context, keyName string, connections []chains.XCallConnection) error {
	if in.testconfig.Environment == "preconfigured" {
		return nil
	}
	testcase := ctx.Value("testcase").(string)

	xCall := in.IBCAddresses["xcall"]
	params := `{"_callService":"` + xCall + `"}`
	dapp, err := in.DeployContractRemote(ctx, in.scorePaths["dapp"], in.keystorePath, params)
	if err != nil {
		return err
	}

	in.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = dapp
	for _, connection := range connections {
		params = `{"nid":"` + connection.Nid + `", "source":"` + in.IBCAddresses[connection.Connection] + `", "destination":"` + connection.Destination + `"}`
		_, err = in.executeContract(context.Background(), dapp, "addConnection", params)
		if err != nil {
			in.log.Error("Unable to add connection",
				zap.Error(err),
				zap.String("nid", connection.Nid),
				zap.String("source", in.IBCAddresses[connection.Connection]),
				zap.String("destination", connection.Destination),
			)
		}
	}

	return nil
}

func (in *IconRemotenet) GetContractAddress(key string) string {
	value, exist := in.IBCAddresses[key]
	if !exist {
		panic(fmt.Sprintf(`IBC address not exist %s`, key))
	}
	return value
}

func (in *IconRemotenet) BackupConfig() ([]byte, error) {
	panic("not implemented")
}

func (in *IconRemotenet) RestoreConfig(backup []byte) error {
	panic("not implemented")
}

func (in *IconRemotenet) SendPacketXCall(ctx context.Context, keyName, _to string, data, rollback []byte) (context.Context, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	// TODO: send fees
	var params = `{"_to":"` + _to + `", "_data":"` + hex.EncodeToString(data) + `"}`
	if rollback != nil {
		params = `{"_to":"` + _to + `", "_data":"` + hex.EncodeToString(data) + `", "_rollback":"` + hex.EncodeToString(rollback) + `"}`
	}
	ctx, err := in.executeContract(ctx, in.IBCAddresses[dappKey], "sendMessage", params)
	if err != nil {
		return nil, err
	}
	txn := ctx.Value("txResult").(*icontypes.TransactionResult)
	return context.WithValue(ctx, "sn", getSn(txn)), nil
}

// HasPacketReceipt returns the receipt of the packet sent to the target chain
func (in *IconRemotenet) IsPacketReceived(ctx context.Context, params map[string]interface{}, order ibc.Order) bool {
	panic("not implemented")
}

// FindTargetXCallMessage returns the request id and the data of the message sent to the target chain
func (in *IconRemotenet) FindTargetXCallMessage(ctx context.Context, target chains.Chain, height uint64, to string) (*chains.XCallResponse, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	sn := ctx.Value("sn").(string)
	reqId, destData, err := target.FindCallMessage(ctx, height, in.cfg.ChainID+"/"+in.IBCAddresses[dappKey], to, sn)
	return &chains.XCallResponse{SerialNo: sn, RequestID: reqId, Data: destData}, err
}

func (in *IconRemotenet) XCall(ctx context.Context, targetChain chains.Chain, keyName, to string, data, rollback []byte) (*chains.XCallResponse, error) {
	height, err := targetChain.Height(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: send fees
	ctx, err = in.SendPacketXCall(ctx, keyName, to, data, rollback)
	if err != nil {
		return nil, err
	}
	return in.FindTargetXCallMessage(ctx, targetChain, height, strings.Split(to, "/")[1])
}

func getSn(tx *icontypes.TransactionResult) string {
	for _, log := range tx.EventLogs {
		if string(log.Indexed[0]) == "CallMessageSent(Address,str,int)" {
			sn, _ := strconv.ParseInt(log.Indexed[3], 0, 64)
			return strconv.FormatInt(sn, 10)
		}
	}
	return ""
}

func (in *IconRemotenet) FindCallMessage(ctx context.Context, startHeight uint64, from, to, sn string) (string, string, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	index := []*string{&from, &to, &sn}
	event, err := in.FindEvent(ctx, startHeight, "xcall", "CallMessage(str,str,int,int,bytes)", index)
	if err != nil {
		return "", "", err
	}

	intHeight, _ := event.Height.Int()
	block, _ := in.IconClient.GetBlockByHeight(&icontypes.BlockHeightParam{Height: icontypes.NewHexInt(int64(intHeight - 1))})
	i, _ := event.Index.Int()
	tx := block.NormalTransactions[i]
	trResult, _ := in.TransactionResult(ctx, string(tx.TxHash))
	eventIndex, _ := event.Events[0].Int()
	reqId := trResult.EventLogs[eventIndex].Data[0]
	data := trResult.EventLogs[eventIndex].Data[1]
	return reqId, data, nil
}

func (in *IconRemotenet) FindCallResponse(ctx context.Context, startHeight uint64, sn string) (string, error) {
	index := []*string{&sn}
	event, err := in.FindEvent(ctx, startHeight, "xcall", "ResponseMessage(int,int)", index)
	if err != nil {
		return "", err
	}
	intHeight, _ := event.Height.Int()
	block, _ := in.IconClient.GetBlockByHeight(&icontypes.BlockHeightParam{Height: icontypes.NewHexInt(int64(intHeight - 1))})
	i, _ := event.Index.Int()
	tx := block.NormalTransactions[i]
	trResult, _ := in.TransactionResult(ctx, string(tx.TxHash))
	eventIndex, _ := event.Events[0].Int()
	code, _ := strconv.ParseInt(trResult.EventLogs[eventIndex].Data[0], 0, 64)

	return strconv.FormatInt(code, 10), nil
}

func (in *IconRemotenet) FindEvent(ctx context.Context, startHeight uint64, contract, signature string, index []*string) (*icontypes.EventNotification, error) {
	filter := icontypes.EventFilter{
		Addr:      icontypes.Address(in.IBCAddresses["xcall"]),
		Signature: signature,
		Indexed:   index,
	}
	// Create a context with a timeout of 16 seconds.
	_ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	// Create an event request with the given filter and start height.
	req := &icontypes.EventRequest{
		EventFilter: filter,
		Height:      icontypes.NewHexInt(int64(startHeight)),
	}
	channel := make(chan *icontypes.EventNotification)
	response := func(_ *websocket.Conn, v *icontypes.EventNotification) error {
		channel <- v
		return nil
	}
	errRespose := func(conn *websocket.Conn, err error) {}
	go func(ctx context.Context, req *icontypes.EventRequest, response func(*websocket.Conn, *icontypes.EventNotification) error, errRespose func(*websocket.Conn, error)) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered: %v", err)
			}
		}()
		if err := in.IconClient.MonitorEvent(ctx, req, response, errRespose); err != nil {
			log.Printf("MonitorEvent error: %v", err)
		}
		defer in.IconClient.CloseAllMonitor()
	}(ctx, req, response, errRespose)

	select {
	case v := <-channel:
		return v, nil
	case <-_ctx.Done():
		latestHeight, _ := in.Height(ctx)
		return nil, errors.New(fmt.Sprintf("timeout : Event %s not found after %d block to %d block ", signature, startHeight, latestHeight))
	}
}

// executeContract implements chains.Chain
func (in *IconRemotenet) executeContract(ctx context.Context, contractAddress, methodName, params string) (context.Context, error) {
	hash, err := in.ExecuteRemoteContract(ctx, contractAddress, methodName, in.keystorePath, params)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Transaction Hash: %s\n", hash)

	txHashByte, err := hex.DecodeString(strings.TrimPrefix(hash, "0x"))
	if err != nil {
		return nil, fmt.Errorf("error when executing contract %v ", err)
	}
	_, res, err := in.IconClient.WaitForResults(ctx, &icontypes.TransactionHashParam{Hash: icontypes.NewHexBytes(txHashByte)})
	if err != nil {
		return nil, err
	}
	if res.Status == "0x1" {
		return context.WithValue(ctx, "txResult", res), nil
	}
	//TODO add debug flag to print trace
	trace, err := in.GetDebugTrace(ctx, icontypes.NewHexBytes(txHashByte))
	if err == nil {
		logs, _ := json.Marshal(trace.Logs)
		fmt.Printf("---------debug trace start-----------\n%s\n---------debug trace end-----------\n", string(logs))
	}
	return ctx, fmt.Errorf("%s", res.Failure.MessageValue)
}

func (in *IconRemotenet) NodeCommand(command ...string) []string {
	command = in.BinCommand(command...)
	return append(command,
		"--uri", in.GetRPCAddress(), //fmt.Sprintf("http://%s/api/v3", in.HostRPCPort),
		"--nid", "0x3",
	)
}

func (in *IconRemotenet) BinCommand(command ...string) []string {
	command = append([]string{in.Config().Bin}, command...)
	return command
}

func (in *IconRemotenet) ExecBin(ctx context.Context, command ...string) ([]byte, []byte, error) {
	return in.Exec(ctx, in.BinCommand(command...), nil)
}

func (in *IconRemotenet) TransactionResult(ctx context.Context, hash string) (*icontypes.TransactionResult, error) {
	uri := in.testconfig.RPCUri
	out, _, err := in.ExecBin(ctx, "rpc", "txresult", hash, "--uri", uri)
	if err != nil {
		return nil, err
	}
	var result = new(icontypes.TransactionResult)
	return result, json.Unmarshal(out, result)
}

func (in *IconRemotenet) GetDebugTrace(ctx context.Context, hash icontypes.HexBytes) (*DebugTrace, error) {
	uri := in.testconfig.RPCUri
	uri = strings.Replace(uri, "v3", "v3d", 1)
	out, _, err := in.ExecBin(ctx, "debug", "trace", string(hash), "--uri", uri)
	if err != nil {
		return nil, err
	}
	var result = new(DebugTrace)
	return result, json.Unmarshal(out, result)

}

func (in *IconRemotenet) DeployContractRemote(ctx context.Context, contractPath, keystorePath, initMessage string) (string, error) {
	_, score := filepath.Split(contractPath)
	// Deploy the contract
	hash, err := in.ExecTx(ctx, initMessage, "/contracts/"+score, keystorePath)
	if err != nil {
		return "", err
	}

	//wait for few blocks
	time.Sleep(3 * time.Second)

	// Get Score Address
	trResult, err := in.TransactionResult(ctx, hash)

	if err != nil {
		return "", err
	}

	return string(trResult.SCOREAddress), nil

}

func (in *IconRemotenet) ExecTx(ctx context.Context, initMessage string, filePath string, keystorePath string, command ...string) (string, error) {
	var output string
	stdout, _, err := in.Exec(ctx, in.TxCommand(ctx, initMessage, filePath, keystorePath, command...), nil)
	if err != nil {
		return "", err
	}
	return output, json.Unmarshal(stdout, &output)
}

// TxCommand is a helper to retrieve a full command for broadcasting a tx
// with the chain node binary.
func (in *IconRemotenet) TxCommand(ctx context.Context, initMessage, filePath, keystorePath string, command ...string) []string {
	// get password from pathname as pathname will have the password prefixed. ex - Alice.Json
	// _, key := filepath.Split(keystorePath)
	// fileName := strings.Split(key, ".")
	// password := fileName[0]

	command = append([]string{"rpc", "sendtx", "deploy", filePath}, command...)
	command = append(command,
		"--key_store", "/goloop/data/"+in.testconfig.KeystoreFile,
		"--key_password", "gochain",
		"--step_limit", "5000000000",
		"--content_type", "application/java",
	)
	if initMessage != "" && initMessage != "{}" {
		if strings.HasPrefix(initMessage, "{") {
			command = append(command, "--params", initMessage)
		} else {
			command = append(command, "--param", initMessage)
		}
	}

	return in.NodeCommand(command...)
}

func (in *IconRemotenet) ExecuteRemoteContract(ctx context.Context, scoreAddress, methodName, keyStorePath, params string) (string, error) {
	return in.ExecCallTx(ctx, scoreAddress, methodName, keyStorePath, params)
}

func (in *IconRemotenet) ExecCallTx(ctx context.Context, scoreAddress, methodName, keystorePath, params string) (string, error) {
	var output string
	stdout, _, err := in.Exec(ctx, in.ExecCallTxCommand(ctx, scoreAddress, methodName, keystorePath, params), nil)
	if err != nil {
		return "", err
	}
	return output, json.Unmarshal(stdout, &output)
}

func (in *IconRemotenet) ExecCallTxCommand(ctx context.Context, scoreAddress, methodName, keystorePath, params string) []string {
	// get password from pathname as pathname will have the password prefixed. ex - Alice.Json
	// _, key := filepath.Split(keystorePath)
	// fileName := strings.Split(key, ".")
	// password := fileName[0]
	command := []string{"rpc", "sendtx", "call"}

	command = append(command,
		"--to", scoreAddress,
		"--method", methodName,
		"--key_store", "/goloop/data/godwallet.json",
		"--key_password", "gochain",
		"--step_limit", "5000000000",
	)

	if params != "" && params != "{}" {
		if strings.HasPrefix(params, "{") {
			command = append(command, "--params", params)
		} else {
			command = append(command, "--param", params)
		}
	}

	if methodName == "registerPRep" {
		command = append(command, "--value", "2000000000000000000000")
	}

	return in.NodeCommand(command...)
}

func (in *IconRemotenet) FindRollbackExecutedMessage(ctx context.Context, startHeight uint64, sn string) (string, error) {
	index := []*string{&sn}
	_, err := in.FindEvent(ctx, startHeight, "xcall", "RollbackExecuted(int)", index)
	if err != nil {
		return "", err
	}
	return "0", nil
}
