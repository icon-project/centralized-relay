package stellar

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/dockerutil"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer/centralized"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"github.com/stellar/go/clients/horizonclient"
	"gopkg.in/yaml.v3"

	//chantypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/docker/docker/client"
	"github.com/icon-project/centralized-relay/test/chains"

	"go.uber.org/zap"
)

type StellarRemotenet struct {
	log           *zap.Logger
	testName      string
	cfg           chains.ChainConfig
	scorePaths    map[string]string
	IBCAddresses  map[string]string     `json:"addresses"`
	Wallets       map[string]ibc.Wallet `json:"wallets"`
	Client        *client.Client
	Network       string
	testconfig    *testconfig.Chain
	horizonClient *horizonclient.Client
	sorobanClient *Client
}

const xcall = "xcall"
const connection = "connection"

func NewStellarRemotenet(testName string, log *zap.Logger, chainConfig chains.ChainConfig,
	client *client.Client, network string, testconfig *testconfig.Chain) chains.Chain {
	httpClient := http.Client{}
	stellarHorizonClient := &horizonclient.Client{
		HorizonURL: testconfig.WebsocketUrl,
		HTTP:       &httpClient,
		AppName:    "centralized-relay-e2e",
	}
	stellarSorobanClient := &Client{
		idCounter:  0,
		httpClient: &httpClient,
		rpcUrl:     testconfig.RPCUri,
	}
	return &StellarRemotenet{
		testName:      testName,
		cfg:           chainConfig,
		log:           log,
		scorePaths:    testconfig.Contracts,
		Wallets:       map[string]ibc.Wallet{},
		IBCAddresses:  make(map[string]string),
		Client:        client,
		testconfig:    testconfig,
		Network:       network,
		horizonClient: stellarHorizonClient,
		sorobanClient: stellarSorobanClient,
	}
}

// Config fetches the chain configuration.
func (sn *StellarRemotenet) Config() chains.ChainConfig {
	return sn.cfg
}

// Exec runs an arbitrary command using Chain's docker environment.
// Whether the invoked command is run in a one-off container or execing into an already running container
// is up to the chain implementation.
//
// "env" are environment variables in the format "MY_ENV_VAR=value"
func (sn *StellarRemotenet) Exec(ctx context.Context, cmd []string, env []string) (stdout []byte, stderr []byte, err error) {
	cmd = append([]string{}, cmd...)
	job := dockerutil.NewImage(sn.log, sn.Client, sn.Network, sn.testName, sn.cfg.Images.Repository, sn.cfg.Images.Version)
	bindPaths := []string{
		sn.testconfig.ContractsPath + ":/contracts",
	}
	opts := dockerutil.ContainerOptions{
		Binds: bindPaths,
	}
	res := job.Run(ctx, cmd, opts)
	return res.Stdout, res.Stderr, res.Err
}

// ExportState exports the chain state at specific height.
func (sn *StellarRemotenet) ExportState(ctx context.Context, height int64) (string, error) {
	block, err := sn.GetClientBlockByHeight(ctx, height)
	return block, err
}

// GetRPCAddress retrieves the rpc address that can be reached by other containers in the docker network.
func (sn *StellarRemotenet) GetRPCAddress() string {
	return sn.testconfig.RPCUri
}

func (sn *StellarRemotenet) GetRelayConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error) {
	contracts := make(map[string]string)
	contracts["xcall"] = sn.GetContractAddress("xcall")
	contracts["connection"] = sn.GetContractAddress("connection")
	config := &centralized.StellarRelayerChainConfig{
		Type: "stellar",
		Value: centralized.StellarRelayerChainConfigValue{
			NID:               sn.Config().ChainID,
			SorobanUrl:        sn.GetRPCAddress(),
			HorizonUrl:        sn.GetGRPCAddress(),
			StartHeight:       0,
			Contracts:         contracts,
			BlockInterval:     "6s",
			Address:           sn.testconfig.RelayWalletAddress,
			FinalityBlock:     uint64(0),
			NetworkPassphrase: sn.testconfig.KeystorePassword,
			MaxInclusionFee:   200,
		},
	}
	return yaml.Marshal(config)
}

// GetGRPCAddress retrieves the grpc address that can be reached by other containers in the docker network.
// Not Applicable for Icon
func (sn *StellarRemotenet) GetGRPCAddress() string {
	return sn.testconfig.WebsocketUrl
}

// GetHostRPCAddress returns the rpc address that can be reached by processes on the host machine.
// Note that this will not return a valid value until after Start returns.
func (sn *StellarRemotenet) GetHostRPCAddress() string {
	return sn.testconfig.RPCUri
}

// GetHostGRPCAddress returns the grpc address that can be reached by processes on the host machine.
// Note that this will not return a valid value until after Start returns.
// Not applicable for Icon
func (sn *StellarRemotenet) GetHostGRPCAddress() string {
	return sn.testconfig.RPCUri
}

// GetAddress fetches the bech32 address for a test key on the "user" node (either the first fullnode or the first validator if no fullnodes).
func (sn *StellarRemotenet) GetAddress(ctx context.Context, keyName string) ([]byte, error) {
	addrInByte, err := json.Marshal(keyName)
	if err != nil {
		return nil, err
	}
	return addrInByte, nil
}

// Height returns the current block height or an error if unable to get current height.
func (sn *StellarRemotenet) Height(ctx context.Context) (uint64, error) {
	res, err := sn.sorobanClient.GetLatestLedger(ctx)
	return uint64(res.Sequence), err
}

// GetGasFeesInNativeDenom gets the fees in native denom for an amount of spent gas.
func (sn *StellarRemotenet) GetGasFeesInNativeDenom(gasPaid int64) int64 {
	gasPrice, _ := strconv.ParseFloat(strings.Replace(sn.cfg.GasPrices, sn.cfg.Denom, "", 1), 64)
	fees := float64(gasPaid) * gasPrice
	return int64(fees)
}

func (sn *StellarRemotenet) SetupConnection(ctx context.Context, target chains.Chain) error {
	if sn.testconfig.Environment == "preconfigured" {
		return nil
	}
	connection, err := sn.DeployContractRemote(ctx, sn.scorePaths["connection"])
	if err != nil {
		return err
	}
	params := make(map[string]string)
	params["msg"] = `{"xcall_address":"` + sn.GetContractAddress("xcall") + `", "native_token":"CDLZFC3SYJYDZT7K67VZ75HPJVIEUVNIXF47ZG2FB2RMQQVU2HHGCYSC", "relayer":"` + sn.testconfig.RelayWalletAddress + `"}`

	_, err = sn.executeContract(context.Background(), connection, "initialize", params)
	if err != nil {
		return err
	}
	sn.IBCAddresses["connection"] = connection
	return nil
}

func (sn *StellarRemotenet) SetupXCall(ctx context.Context) error {
	if sn.testconfig.Environment == "preconfigured" {
		testcase := ctx.Value("testcase").(string)
		sn.IBCAddresses["xcall"] = "CAKBVQ2QLJWQ5AO4HPND3XEDLGCLDW7UV5WATLDLWBKP535EGDEVKSD4"
		sn.IBCAddresses["connection"] = "CC4WN23OV5MRA5DKRBP5SZB6XUNXXWKNOUQ3TKTPCW3EM3KGIE54YANI"
		sn.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = "CA7VNLUTPQZEKAXJHXHCS2YJHPELODRZKB26MBQQISD4BPHG2PGAUPMM"
		return nil
	}
	xcall, err := sn.DeployContractRemote(ctx, sn.scorePaths["xcall"])
	if err != nil {
		return err
	}
	sn.IBCAddresses["xcall"] = xcall

	//init xcall
	params := make(map[string]string)
	params["msg"] = `{"network_id":"` + sn.Config().ChainID + `", "sender":"` + sn.testconfig.RelayWalletAddress + `", "native_token":"CDLZFC3SYJYDZT7K67VZ75HPJVIEUVNIXF47ZG2FB2RMQQVU2HHGCYSC"}`
	_, err = sn.executeContract(context.Background(), xcall, "initialize", params)
	return err
}

func (sn *StellarRemotenet) DeployXCallMockApp(ctx context.Context, keyName string, connections []chains.XCallConnection) error {
	if sn.testconfig.Environment == "preconfigured" {
		return nil
	}
	testcase := ctx.Value("testcase").(string)

	dapp, err := sn.DeployContractRemote(ctx, sn.scorePaths["dapp"])
	if err != nil {
		return err
	}

	//init dapp
	xCall := sn.IBCAddresses["xcall"]
	params := make(map[string]string)
	params["xcall_address"] = xCall
	sn.executeContract(ctx, dapp, "init", params)

	sn.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = dapp
	for _, connection := range connections {
		params := make(map[string]string)
		params["src_endpoint"] = sn.IBCAddresses[connection.Connection]
		params["network_id"] = connection.Nid
		params["dst_endpoint"] = connection.Destination
		_, err = sn.executeContract(context.Background(), dapp, "add_connection", params)
		if err != nil {
			sn.log.Error("Unable to add connection",
				zap.Error(err),
				zap.String("nid", connection.Nid),
				zap.String("source", sn.IBCAddresses[connection.Connection]),
				zap.String("destination", connection.Destination),
			)
		}
	}

	return nil
}

func (sn *StellarRemotenet) GetContractAddress(key string) string {
	value, exist := sn.IBCAddresses[key]
	if !exist {
		panic(fmt.Sprintf(`IBC address not exist %s`, key))
	}
	return value
}

func (sn *StellarRemotenet) BackupConfig() ([]byte, error) {
	panic("not implemented")
}

func (sn *StellarRemotenet) RestoreConfig(backup []byte) error {
	panic("not implemented")
}

func (sn *StellarRemotenet) SendPacketXCall(ctx context.Context, keyName, _to string, data, rollback []byte) (context.Context, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	// TODO: send fees
	params := make(map[string]string)
	params["to"] = `{"0": "` + _to + `"}`
	params["data"] = hex.EncodeToString(data)
	params["sender"] = sn.testconfig.RelayWalletAddress
	params["msg_type"] = "0" //or 2 for message Persisted
	if rollback != nil {
		params["rollback"] = hex.EncodeToString(rollback)
		params["msg_type"] = "1"
	}

	ctx, err := sn.executeContract(ctx, sn.IBCAddresses[dappKey], "send_call_message", params)
	if err != nil {
		return nil, err
	}
	sno := ctx.Value("sno").(string)
	return context.WithValue(ctx, "sn", sno), nil
}

// HasPacketReceipt returns the receipt of the packet sent to the target chain
func (sn *StellarRemotenet) IsPacketReceived(ctx context.Context, params map[string]interface{}, order ibc.Order) bool {
	panic("not implemented")
}

// FindTargetXCallMessage returns the request id and the data of the message sent to the target chain
func (sn *StellarRemotenet) FindTargetXCallMessage(ctx context.Context, target chains.Chain, height uint64, to string) (*chains.XCallResponse, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	sno := ctx.Value("sn").(string)
	reqId, destData, err := target.FindCallMessage(ctx, height, sn.cfg.ChainID+"/"+sn.IBCAddresses[dappKey], to, sno)
	return &chains.XCallResponse{SerialNo: sno, RequestID: reqId, Data: destData}, err
}

func (sn *StellarRemotenet) XCall(ctx context.Context, targetChain chains.Chain, keyName, to string, data, rollback []byte) (*chains.XCallResponse, error) {
	height, err := targetChain.Height(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: send fees
	ctx, err = sn.SendPacketXCall(ctx, keyName, to, data, rollback)
	if err != nil {
		return nil, err
	}
	return sn.FindTargetXCallMessage(ctx, targetChain, height, strings.Split(to, "/")[1])
}

func (sn *StellarRemotenet) FindCallMessage(ctx context.Context, startHeight uint64, from, to, sno string) (string, string, error) {

	event, err := sn.FindEvent(ctx, startHeight, "xcall", "CallMessage", sno)
	if err != nil {
		return "", "", err
	}
	reqId := event.ValueDecoded["reqId"].(uint64)
	data := event.ValueDecoded["data"].([]byte)
	return strconv.FormatUint(reqId, 10), string(data), nil
}

func (sn *StellarRemotenet) FindCallResponse(ctx context.Context, startHeight uint64, sno string) (string, error) {
	event, err := sn.FindEvent(ctx, startHeight, "xcall", "ResponseMessage", sno)
	if err != nil {
		return "", err
	}
	code := event.ValueDecoded["code"].(uint64)
	return strconv.FormatUint(code, 10), nil
}

func (sn *StellarRemotenet) FindEvent(ctx context.Context, startHeight uint64, contract, signature, sno string) (*EventResponseEvent, error) {
	timeout := time.After(120 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("failed to find eventLog")
		case <-ticker.C:
			data, err := sn.getEvent(ctx, startHeight, sno, signature, sn.GetContractAddress(contract))
			if err != nil {
				continue
			}
			return data, nil
		}
	}
}

func (sn *StellarRemotenet) getEvent(ctx context.Context, startHeight uint64, sno, signature, contractId string) (*EventResponseEvent, error) {
	return sn.sorobanClient.GetEvent(ctx, startHeight, sno, contractId, signature)
}

// Remote implements chains.Chain
func (sn *StellarRemotenet) DeployContract(ctx context.Context, keyName string) (context.Context, error) {
	return ctx, nil
}

// executeContract implements chains.Chain
func (sn *StellarRemotenet) executeContract(ctx context.Context, contractAddress, methodName string, params map[string]string) (context.Context, error) {
	var output string
	stdout, err := sn.ExecCallTx(ctx, contractAddress, methodName, params)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(stdout), &output)
	return context.WithValue(ctx, "sno", output), err
}

// GetBlockByHeight implements chains.Chain
func (sn *StellarRemotenet) GetClientBlockByHeight(ctx context.Context, height int64) (string, error) {
	return "", nil
}

// GetLastBlock implements chains.Chain
func (sn *StellarRemotenet) GetLastBlock(ctx context.Context) (context.Context, error) {
	res, err := sn.sorobanClient.GetLatestLedger(ctx)
	h := uint64(res.Sequence)
	return context.WithValue(ctx, chains.LastBlock{}, h), err
}

// QueryContract implements chains.Chain
func (sn *StellarRemotenet) QueryContract(ctx context.Context, contractAddress, methodName string, params map[string]interface{}) (context.Context, error) {
	return ctx, nil
}

func (sn *StellarRemotenet) NodeCommand(command ...string) []string {
	command = sn.BinCommand(command...)
	return append(command,
		"--rpc-url", sn.GetRPCAddress(),
		"--source-account", sn.testconfig.KeystoreFile,
		"--network-passphrase", sn.testconfig.KeystorePassword,
	)
}

func (sn *StellarRemotenet) BinCommand(command ...string) []string {
	command = append([]string{sn.Config().Bin}, command...)
	return command
}

func (sn *StellarRemotenet) ExecBin(ctx context.Context, command ...string) ([]byte, []byte, error) {
	return sn.Exec(ctx, sn.BinCommand(command...), nil)
}

func (sn *StellarRemotenet) DeployContractRemote(ctx context.Context, contractPath string) (string, error) {
	_, score := filepath.Split(contractPath)
	// Deploy the contract
	contractId, err := sn.ExecTx(ctx, "/contracts/"+score)
	if err != nil {
		return "", err
	}
	return contractId, nil

}

func (sn *StellarRemotenet) ExecTx(ctx context.Context, filePath string, command ...string) (string, error) {
	stdout, _, err := sn.Exec(ctx, sn.TxCommand(ctx, filePath, command...), nil)
	return strings.Split(string(stdout), "\n")[0], err
}

// TxCommand is a helper to retrieve a full command for broadcasting a tx
// with the chain node binary.
func (sn *StellarRemotenet) TxCommand(ctx context.Context, filePath string, command ...string) []string {
	command = append([]string{"contract", "deploy", "--wasm", filePath}, command...)
	return sn.NodeCommand(command...)
}

func (sn *StellarRemotenet) ExecCallTx(ctx context.Context, scoreAddress, methodName string, params map[string]string) (string, error) {
	stdout, _, err := sn.Exec(ctx, sn.ExecCallTxCommand(ctx, scoreAddress, methodName, params), nil)
	return strings.Split(string(stdout), "\n")[0], err
}

func (sn *StellarRemotenet) ExecCallTxCommand(ctx context.Context, scoreAddress, methodName string, params map[string]string) []string {
	command := []string{"contract", "invoke"}

	command = append(command,
		"--id", scoreAddress,
		"--rpc-url", sn.GetRPCAddress(),
		"--source-account", sn.testconfig.KeystoreFile,
		"--network-passphrase", sn.testconfig.KeystorePassword,
		"--", methodName,
	)

	for key, value := range params {
		command = append(command, "--"+key, value)
	}
	command = sn.BinCommand(command...)
	return command
}

func (sn *StellarRemotenet) FindRollbackExecutedMessage(ctx context.Context, startHeight uint64, sno string) (string, error) {
	event, err := sn.FindEvent(ctx, startHeight, "xcall", "RollbackExecuted", sno)
	if err != nil {
		return "", err
	}
	fsno := event.ValueDecoded["sn"].(uint64)
	return strconv.FormatUint(fsno, 10), nil
}
