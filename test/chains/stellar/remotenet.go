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

	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/blockdb"
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
	cfg           ibc.ChainConfig
	numValidators int
	numFullNodes  int
	keystorePath  string
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

func NewStellarRemotenet(testName string, log *zap.Logger, chainConfig ibc.ChainConfig, client *client.Client, network string, testconfig *testconfig.Chain) chains.Chain {
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
func (sn *StellarRemotenet) Config() ibc.ChainConfig {
	return sn.cfg
}

func (sn *StellarRemotenet) OverrideConfig(key string, value any) {
	if value == nil {
		return
	}
	sn.cfg.ConfigFileOverrides[key] = value
}

// Initialize initializes node structs so that things like initializing keys can be done before starting the chain
func (sn *StellarRemotenet) Initialize(ctx context.Context, testName string, cli *client.Client, networkID string) error {
	return nil
}

// Start sets up everything needed (validators, gentx, fullnodes, peering, additional accounts) for chain to start from genesis.
func (sn *StellarRemotenet) Start(testName string, ctx context.Context, additionalGenesisWallets ...ibc.WalletAmount) error {
	return nil
}

// Exec runs an arbitrary command using Chain's docker environment.
// Whether the invoked command is run in a one-off container or execing into an already running container
// is up to the chain implementation.
//
// "env" are environment variables in the format "MY_ENV_VAR=value"
func (sn *StellarRemotenet) Exec(ctx context.Context, cmd []string, env []string) (stdout []byte, stderr []byte, err error) {
	cmd = append([]string{}, cmd...)
	job := dockerutil.NewImage(sn.log, sn.Client, sn.Network, sn.testName, sn.cfg.Images[0].Repository, sn.cfg.Images[0].Version)
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

// HomeDir is the home directory of a node running in a docker container. Therefore, this maps to
// the container's filesystem (not the host).
func (sn *StellarRemotenet) HomeDir() string {
	return ""
}

func (sn *StellarRemotenet) createKeystore(ctx context.Context, keyName string) (string, string, error) {
	return "", "", nil
}

// RecoverKey recovers an existing user from a given mnemonic.
func (sn *StellarRemotenet) RecoverKey(ctx context.Context, name string, mnemonic string) error {
	panic("not implemented") // TODO: Implement
}

// GetAddress fetches the bech32 address for a test key on the "user" node (either the first fullnode or the first validator if no fullnodes).
func (sn *StellarRemotenet) GetAddress(ctx context.Context, keyName string) ([]byte, error) {
	addrInByte, err := json.Marshal(keyName)
	if err != nil {
		return nil, err
	}
	return addrInByte, nil
}

// SendFunds sends funds to a wallet from a user account.
func (sn *StellarRemotenet) SendFunds(ctx context.Context, keyName string, amount ibc.WalletAmount) error {
	return nil
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

// BuildRelayerWallet will return a chain-specific wallet populated with the mnemonic so that the wallet can
// be restored in the relayer node using the mnemonic. After it is built, that address is included in
// genesis with some funds.
func (sn *StellarRemotenet) BuildRelayerWallet(ctx context.Context, keyName string) (ibc.Wallet, error) {
	return sn.BuildWallet(ctx, keyName, "")
}

func (sn *StellarRemotenet) BuildWallet(ctx context.Context, keyName string, mnemonic string) (ibc.Wallet, error) {
	return nil, nil
}

func (sn *StellarRemotenet) FindTxs(ctx context.Context, height uint64) ([]blockdb.Tx, error) {
	return nil, nil
}

// GetBalance fetches the current balance for a specific account address and denom.
func (sn *StellarRemotenet) GetBalance(ctx context.Context, address string, denom string) (int64, error) {
	return 0, nil
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
		sn.IBCAddresses["xcall"] = "CASJ76AJJKK6BFMMB4DLLYV5J6OBTBIO25FDWXWZY5BIXFVCNK7XACNG"
		sn.IBCAddresses["connection"] = "CBLG6CNVXWSCF7X5XXM6B5JZ2ULFIPV7KRENAIB2HB4FAKMGIMD63BJU"
		sn.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = "CA2QFEESKXOWQZJRQ7LUIGTDEEWNOLR4KRYL35GUWNRK5XYZ3FDP3L7W"
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
	height, err := targetChain.(ibc.Chain).Height(ctx)
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

func (sn *StellarRemotenet) ExecuteCall(ctx context.Context, reqId, data string) (context.Context, error) {
	panic("not required in e2e")
}

func (sn *StellarRemotenet) ExecuteRollback(ctx context.Context, sno string) (context.Context, error) {
	params := make(map[string]string)
	params["sequence_no"] = sno
	ctx, err := sn.executeContract(ctx, sn.IBCAddresses["xcall"], "execute_rollback", params)
	if err != nil {
		return nil, err
	}
	height, _ := sn.Height(ctx)
	_, err = sn.FindEvent(ctx, height-20, "xcall", "RollbackExecuted", sno)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, "IsRollbackEventFound", true), nil

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
	timeout := time.After(60 * time.Second)
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

func (sn *StellarRemotenet) ExecuteContract(ctx context.Context, contractAddress, keyName, methodName string, params map[string]interface{}) (context.Context, error) {
	return nil, nil
}

func (sn *StellarRemotenet) GetBlockByHeight(context.Context) (context.Context, error) {
	panic("not implemented")
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

func (sn *StellarRemotenet) InitEventListener(ctx context.Context, contract string) chains.EventListener {
	return nil
}

// QueryContract implements chains.Chain
func (sn *StellarRemotenet) QueryContract(ctx context.Context, contractAddress, methodName string, params map[string]interface{}) (context.Context, error) {
	return ctx, nil
}

func (sn *StellarRemotenet) BuildWallets(ctx context.Context, keyName string) (ibc.Wallet, error) {
	panic("not implemented")
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
