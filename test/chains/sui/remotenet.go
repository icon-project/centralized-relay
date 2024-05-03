package sui

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/coming-chat/go-sui/v2/account"
	suisdkClient "github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/docker/docker/client"
	"github.com/fardream/go-bcs/bcs"
	"github.com/icon-project/centralized-relay/test/chains"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/blockdb"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/dockerutil"
	ibcLocal "github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer/centralized"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const (
	suiCurrencyType                           = "0x2::sui::SUI"
	pickMethod                                = 1
	baseSuiFee                                = 1000
	suiStringType                             = "0x1::string::String"
	suiU64                                    = "u64"
	suiBool                                   = "bool"
	moveCall        suisdkClient.UnsafeMethod = "moveCall"
	publish         suisdkClient.UnsafeMethod = "publish"
	queryEvents     suisdkClient.SuiXMethod   = "queryEvents"
	callGasBudget                             = 5000000
	deployGasBudget                           = "500000000"
	xcallAdmin                                = "xcall-admin"
	xcallStorage                              = "xcall-storage"
	sui_rlp_path                              = "libs/sui_rlp"
)

type MoveTomlConfig struct {
	Package         map[string]string     `toml:"package"`
	Dependencies    map[string]Dependency `toml:"dependencies"`
	Addresses       map[string]string     `toml:"addresses"`
	DevDependencies map[string]Dependency `toml:"dev-dependencies"`
	DevAddresses    map[string]string     `toml:"dev-addresses"`
}

type Dependency struct {
	Git    string `toml:"git,omitempty"`
	Subdir string `toml:"subdir,omitempty"`
	Rev    string `toml:"rev,omitempty"`
	Local  string `toml:"local,omitempty"`
}

type DepoymentInfo struct {
	PackageId string
	AdminCap  string
	Storage   string
	Witness   string
}

type PackageInfo struct {
	Modules      []string `json:"modules"`
	Dependencies []string `json:"dependencies"`
	Digest       []int    `json:"digest"`
}

func NewSuiRemotenet(testName string, log *zap.Logger, chainConfig ibcLocal.ChainConfig, client *client.Client, network string, testconfig *testconfig.Chain) chains.Chain {

	suiClient, err := suisdkClient.Dial(testconfig.RPCUri)
	if err != nil {
		panic("error connecting sui rpc")
	}

	return &SuiRemotenet{
		testName:     testName,
		cfg:          chainConfig,
		log:          log,
		IBCAddresses: make(map[string]string),
		filepath:     testconfig.Contracts,
		client:       suiClient,
		DockerClient: client,
		testconfig:   testconfig,
		Network:      network,
	}
}

func (an *SuiRemotenet) Config() ibcLocal.ChainConfig {
	return an.cfg
}

// Initialize initializes node structs so that things like initializing keys can be done before starting the chain
func (an *SuiRemotenet) Initialize(ctx context.Context, testName string, cli *client.Client, networkID string) error {
	return nil
}

func (an *SuiRemotenet) BuildRelayerWallet(ctx context.Context, keyName string) (ibcLocal.Wallet, error) {
	return nil, nil
}

func (an *SuiRemotenet) BuildWallet(ctx context.Context, keyName string, mnemonic string) (ibcLocal.Wallet, error) {
	return nil, nil
}

// CreateKey implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).CreateKey of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) CreateKey(ctx context.Context, keyName string) error {
	panic("unimplemented")
}

// DeployContract implements chains.Chain.
func (an *SuiRemotenet) DeployContract(ctx context.Context, keyName string) (context.Context, error) {
	filePath := "/xcall/" + keyName
	stdout, _, err := an.ExecBin(ctx, "sui", "client", "publish", filePath, "--gas-budget", deployGasBudget, "--json")
	if err != nil {
		return ctx, err
	}
	var resp *types.SuiTransactionBlockResponse
	err = json.Unmarshal(stdout, &resp)
	if err != nil {
		return ctx, err
	}
	an.log.Info("Deploy completed ", zap.Any("txDigest", resp.Digest), zap.Any("status", resp.Effects.Data.IsSuccess()))
	if resp.Effects.Data.V1.Status.Status != "success" {
		return nil, fmt.Errorf("error while committing tx : %s", resp.Effects.Data.V1.Status.Error)
	}
	depoymentInfo := DepoymentInfo{}
	for _, changes := range resp.ObjectChanges {
		if changes.Data.Published != nil {
			depoymentInfo.PackageId = changes.Data.Published.PackageId.String()
		}
		if changes.Data.Created != nil && strings.Contains(changes.Data.Created.ObjectType, "AdminCap") {
			depoymentInfo.AdminCap = changes.Data.Created.ObjectId.String()
		}
		if changes.Data.Created != nil && strings.Contains(changes.Data.Created.ObjectType, "Storage") {
			depoymentInfo.Storage = changes.Data.Created.ObjectId.String()
		}
		if changes.Data.Created != nil && strings.Contains(changes.Data.Created.ObjectType, "WitnessCarrier") {
			depoymentInfo.Witness = changes.Data.Created.ObjectId.String()
		}
	}
	return context.WithValue(ctx, "objId", depoymentInfo), nil
}

// DeployXCallMockApp implements chains.Chain.
func (an *SuiRemotenet) DeployXCallMockApp(ctx context.Context, keyName string, connections []chains.XCallConnection) error {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	ctx, err := an.DeployContract(ctx, "mock_dapp")
	if err != nil {
		return err
	}
	deploymentInfo := ctx.Value("objId").(DepoymentInfo)
	an.IBCAddresses[dappKey] = deploymentInfo.PackageId
	an.IBCAddresses[dappKey+"-witness"] = deploymentInfo.Witness
	an.log.Info("setup Dapp completed ", zap.Any("pacckageId", deploymentInfo.PackageId), zap.Any("witness", deploymentInfo.Witness))

	// register xcall
	params := []interface{}{
		an.IBCAddresses[xcallStorage],
		an.IBCAddresses[dappKey+"-witness"],
	}
	msg := an.NewSuiMessage(params, an.IBCAddresses[dappKey], "mock_dapp", "register_xcall")
	resp, err := an.callContract(ctx, msg)
	for _, changes := range resp.ObjectChanges {
		if changes.Data.Created != nil && strings.Contains(changes.Data.Created.ObjectType, "DappState") {
			an.IBCAddresses[dappKey+"-state"] = changes.Data.Created.ObjectId.String()
		}
	}
	an.log.Info("register xcall completed ", zap.Any("dapp-state", an.IBCAddresses[dappKey+"-state"]))
	if err != nil {
		return err
	}
	// add connections
	for _, connection := range connections {
		// connections passed are  {connection emv.local 0x47c05BCCA7d57c87083EB4e586007530eE4539e9 0}
		// Connection    string Nid           string Destination   string TimeoutHeight int `default:"100"`
		params = []interface{}{
			an.IBCAddresses[dappKey+"-state"],
			connection.Nid,
			"centralized",
			connection.Destination,
		}
		msg = an.NewSuiMessage(params, an.IBCAddresses[dappKey], "mock_dapp", "add_connection")
		resp, err = an.callContract(ctx, msg)
		for _, changes := range resp.ObjectChanges {
			if changes.Data.Created != nil && strings.Contains(changes.Data.Created.ObjectType, "DappState") {
				an.IBCAddresses[dappKey+connection.Connection+"-state"] = changes.Data.Created.ObjectId.String()
			}
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (an *SuiRemotenet) ExecBin(ctx context.Context, command ...string) ([]byte, []byte, error) {
	return an.Exec(ctx, command, nil)
}

// Exec implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).Exec of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) Exec(ctx context.Context, cmd []string, env []string) (stdout []byte, stderr []byte, err error) {
	job := dockerutil.NewImage(an.log, an.DockerClient, an.Network, an.testName, an.cfg.Images[0].Repository, an.cfg.Images[0].Version)

	bindPaths := []string{
		an.testconfig.ContractsPath + ":/xcall",
		an.testconfig.ConfigPath + ":/root/.sui/sui_config/",
	}
	opts := dockerutil.ContainerOptions{
		Binds: bindPaths,
	}
	res := job.Run(ctx, cmd, opts)
	return res.Stdout, res.Stderr, res.Err
}

// ExecuteCall implements chains.Chain.
func (an *SuiRemotenet) ExecuteCall(ctx context.Context, reqId string, data string) (context.Context, error) {
	panic("won't implement")
}

// ExecuteContractDevnet implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).ExecuteContractDevnet of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) ExecuteContract(ctx context.Context, contractAddress string, keyName string, methodName string, param map[string]interface{}) (context.Context, error) {
	panic("unimplemented")
}

// ExecuteRollback implements chains.Chain.
func (an *SuiRemotenet) ExecuteRollback(ctx context.Context, sn string) (context.Context, error) {

	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	params := []interface{}{
		an.IBCAddresses[dappKey+"-state"],
		an.IBCAddresses["xcall-storage"],
		sn,
	}
	msg := an.NewSuiMessage(params, an.IBCAddresses[dappKey], "mock_dapp", "execute_rollback")
	resp, err := an.callContract(ctx, msg)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, "IsRollbackEventFound", an.findSn(resp, "::main::RollbackExecuted") == sn), nil
}

// ExportState implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).ExportState of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) ExportState(ctx context.Context, height int64) (string, error) {
	panic("unimplemented")
}

// FindCallMessage implements chains.Chain.
func (an *SuiRemotenet) FindCallMessage(ctx context.Context, startHeight uint64, from string, to string, sn string) (string, string, error) {
	xCallKey := "xcall" //fmt.Sprintf("xcall-%s", testcase)
	index := ""
	event, err := an.FindEvent(ctx, startHeight, xCallKey, index)
	if err != nil {
		return "", "", err
	}
	jsonData := (event.ParsedJson.(map[string]interface{}))
	return jsonData["req_id"].(string), jsonData["data"].(string), nil
}

type FieldFilter struct {
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type MoveEventModule struct {
	Package string `json:"package"`
	Module  string `json:"module"`
}

type MoveEvent struct {
	MoveEventModule MoveEventModule `json:"MoveEventModule"`
}

func (an *SuiRemotenet) getEvent(ctx context.Context) (*types.SuiEvent, error) {
	limit := uint(100)
	query := MoveEvent{
		MoveEventModule: MoveEventModule{
			Package: an.IBCAddresses["xcall"],
			Module:  "main",
		},
	}
	var resp types.EventPage
	err := an.client.CallContext(ctx, &resp, queryEvents, query, nil, limit, true)

	if err != nil {
		return nil, err
	}
	for _, event := range resp.Data {
		jsonData := (event.ParsedJson.(map[string]interface{}))
		if jsonData["sn"] == "1" && strings.Contains(event.Type, "::main::CallMessage") {
			return &event, nil
		}
	}
	return nil, fmt.Errorf("event not found")
}

func (an *SuiRemotenet) FindEvent(ctx context.Context, startHeight uint64, contract, index string) (*types.SuiEvent, error) {
	timeout := time.After(60 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	// txDigest, _ := lib.NewBase58("GUP5Z9rAhQ3Li6b54gyCFni4LJKudeB1M5BTJ4J1j2Xt")

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("failed to find eventLog")
		case <-ticker.C:
			data, err := an.getEvent(ctx)
			if err != nil {
				continue
			}
			return data, nil
		}
	}
	// // wss not working in devnet/testnet due to limited wss connections
}

// FindCallResponse implements chains.Chain.
func (an *SuiRemotenet) FindCallResponse(ctx context.Context, startHeight uint64, sn string) (string, error) {
	panic("unimplemented")
}

// FindTargetXCallMessage implements chains.Chain.
func (an *SuiRemotenet) FindTargetXCallMessage(ctx context.Context, target chains.Chain, height uint64, to string) (*chains.XCallResponse, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	sn := ctx.Value("sn").(string)
	an.log.Info("Finding xcall msg ", zap.Any("sn", sn), zap.Any("target_chain", an.cfg.ChainID+"/"+an.IBCAddresses[dappKey]))
	reqId, destData, err := target.FindCallMessage(ctx, height, an.cfg.ChainID+"/"+an.IBCAddresses[dappKey], to, sn)
	return &chains.XCallResponse{SerialNo: sn, RequestID: reqId, Data: destData}, err
}

// FindTxs implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).FindTxs of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) FindTxs(ctx context.Context, height uint64) ([]blockdb.Tx, error) {
	return nil, nil
}

// GetAddress implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).GetAddress of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) GetAddress(ctx context.Context, keyName string) ([]byte, error) {
	panic("unimplemented")
}

// GetBalance implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).GetBalance of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) GetBalance(ctx context.Context, address string, denom string) (int64, error) {
	panic("unimplemented")
}

// GetBlockByHeight implements chains.Chain.
func (an *SuiRemotenet) GetBlockByHeight(ctx context.Context) (context.Context, error) {
	panic("unimplemented")
}

// GetContractAddress implements chains.Chain.
func (an *SuiRemotenet) GetContractAddress(key string) string {
	if key == "connection" {
		key = "xcall"
	}
	value, exist := an.IBCAddresses[key]
	if !exist {
		panic(fmt.Sprintf(`IBC address not exist %s`, key))
	}

	return value
}

// GetGRPCAddress implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).GetGRPCAddress of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) GetGRPCAddress() string {
	return an.testconfig.RPCUri
}

// GetGasFeesInNativeDenom implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).GetGasFeesInNativeDenom of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) GetGasFeesInNativeDenom(gasPaid int64) int64 {
	panic("unimplemented")
}

// GetHostGRPCAddress implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).GetHostGRPCAddress of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) GetHostGRPCAddress() string {
	panic("unimplemented")
}

// GetHostRPCAddress implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).GetHostRPCAddress of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) GetHostRPCAddress() string {
	return an.testconfig.RPCUri
}

// GetLastBlock implements chains.Chain.
func (an *SuiRemotenet) GetLastBlock(ctx context.Context) (context.Context, error) {
	panic("unimplemented")
}

// GetRPCAddress implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).GetRPCAddress of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) GetRPCAddress() string {
	return an.testconfig.RPCUri
}

// GetRelayConfig implements chains.Chain.
func (an *SuiRemotenet) GetRelayConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error) {
	contracts := make(map[string]string)
	contracts["xcall"] = an.GetContractAddress("xcall")
	config := &centralized.SUIRelayerChainConfig{
		Type: "icon",
		Value: centralized.SUIRelayerChainConfigValue{
			NID:           an.Config().ChainID,
			RPCURL:        an.GetRPCAddress(),
			StartHeight:   0,
			NetworkID:     0x3,
			Contracts:     contracts,
			BlockInterval: "6s",
			Address:       an.testconfig.RelayWalletAddress,
			FinalityBlock: uint64(10),
			StepMin:       25000,
			StepLimit:     2500000,
		},
	}
	return yaml.Marshal(config)
}

// Height implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).Height of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) Height(ctx context.Context) (uint64, error) {
	checkPoint, err := an.client.GetLatestCheckpointSequenceNumber(ctx)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(checkPoint, 10, 64)
}

// HomeDir implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).HomeDir of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) HomeDir() string {
	return ""
}

// InitEventListener implements chains.Chain.
func (an *SuiRemotenet) InitEventListener(ctx context.Context, contract string) chains.EventListener {
	panic("unimplemented")
}

// QueryContract implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).QueryContract of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) QueryContract(ctx context.Context, contractAddress string, methodName string, params map[string]interface{}) (context.Context, error) {
	panic("unimplemented")
}

// RecoverKey implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).RecoverKey of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) RecoverKey(ctx context.Context, name string, mnemonic string) error {
	panic("unimplemented")
}

// SendFunds implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).SendFunds of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) SendFunds(ctx context.Context, keyName string, amount ibcLocal.WalletAmount) error {
	panic("unimplemented")
}

// SendPacketXCall implements chains.Chain.
func (an *SuiRemotenet) SendPacketXCall(ctx context.Context, keyName string, _to string, data []byte, rollback []byte) (context.Context, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	if rollback == nil {
		rollback = make([]byte, 0)
		_ = rollback
	}
	gasFeeCoin := an.getGasCoinId(ctx, an.testconfig.RelayWalletAddress, callGasBudget).CoinObjectId
	coinId := an.getAnotherGasCoinId(ctx, an.testconfig.RelayWalletAddress, callGasBudget, gasFeeCoin)
	params := []interface{}{
		an.IBCAddresses[dappKey+"-state"],
		an.IBCAddresses[xcallStorage],
		coinId.CoinObjectId,
		_to,
		data,
	}
	msg := an.NewSuiMessage(params, an.IBCAddresses[dappKey], "mock_dapp", "send_message")
	resp, err := an.callContract(ctx, msg)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, "sn", an.findSn(resp, "::main::CallMessageSent")), nil
}

func (an *SuiRemotenet) findSn(tx *types.SuiTransactionBlockResponse, eType string) string {
	// find better way to parse events
	for _, event := range tx.Events {
		if event.Type == (an.IBCAddresses["xcall"] + eType) {
			jsonData := (event.ParsedJson.(map[string]interface{}))
			return jsonData["sn"].(string)
		}
	}
	return ""
}

// SetupConnection implements chains.Chain.
func (an *SuiRemotenet) SetupConnection(ctx context.Context, target chains.Chain) error {
	params := []interface{}{
		an.IBCAddresses[xcallStorage],
		an.IBCAddresses[xcallAdmin],
		"sui",
		"centralized",
	}
	msg := an.NewSuiMessage(params, an.IBCAddresses["xcall"], "main", "register_connection")
	_, err := an.callContract(ctx, msg)
	return err

}

func (an *SuiRemotenet) callContract(ctx context.Context, msg *SuiMessage) (*types.SuiTransactionBlockResponse, error) {
	txnMetadata, err := an.ExecuteContractRemote(ctx, msg, an.testconfig.RelayWalletAddress, uint64(callGasBudget))
	if err != nil {
		return nil, err
	}

	walletAccount, err := account.NewAccountWithKeystore(an.testconfig.KeystorePassword)
	if err != nil {
		return nil, err
	}

	signature, err := walletAccount.SignSecureWithoutEncode(txnMetadata.TxBytes, sui_types.DefaultIntent())
	if err != nil {
		return nil, err
	}
	signatures := []any{signature}

	resp, err := an.CommitTx(ctx, walletAccount, txnMetadata.TxBytes, signatures)
	if err != nil {
		return nil, err
	}
	an.log.Info("Txn created", zap.Any("ID", resp.Digest), zap.Any("status", resp.Effects.Data.IsSuccess()))
	return resp, nil
}

// SetupXCall implements chains.Chain.
func (an *SuiRemotenet) SetupXCall(ctx context.Context) error {
	//deploy rlp
	ctx, err := an.DeployContract(ctx, sui_rlp_path)
	if err != nil {
		return err
	}
	deploymentInfo := ctx.Value("objId").(DepoymentInfo)
	err = an.updateTomlFile(sui_rlp_path, deploymentInfo.PackageId)
	if err != nil {
		return err
	}

	// deploy xcall
	ctx, err = an.DeployContract(ctx, "xcall")
	if err != nil {
		return err
	}
	deploymentInfo = ctx.Value("objId").(DepoymentInfo)
	an.IBCAddresses["xcall"] = deploymentInfo.PackageId
	an.IBCAddresses[xcallAdmin] = deploymentInfo.AdminCap
	an.IBCAddresses[xcallStorage] = deploymentInfo.Storage
	err = an.updateTomlFile("xcall", deploymentInfo.PackageId)
	if err != nil {
		return err
	}
	an.log.Info("setup xcall completed ", zap.Any("packageId", deploymentInfo.PackageId), zap.Any("admin", deploymentInfo.AdminCap), zap.Any("storage", deploymentInfo.Storage))
	return nil
}

func (an *SuiRemotenet) updateTomlFile(keyName, deployedPackageId string) error {
	var cfg MoveTomlConfig
	filePath := an.testconfig.ContractsPath + "/" + keyName
	file, err := os.Open(filePath + "/Move.toml")
	if err != nil {
		return err
	}
	defer file.Close()
	moveConfig, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	err = toml.Unmarshal(moveConfig, &cfg)
	if err != nil {
		return err
	}
	pkgName := cfg.Package["name"]
	cfg.Addresses[pkgName] = deployedPackageId
	cfg.Package["published-at"] = deployedPackageId
	b, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath+"/Move.toml", b, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Start implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).Start of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) Start(testName string, ctx context.Context, additionalGenesisWallets ...ibcLocal.WalletAmount) error {
	panic("unimplemented")
}

// XCall implements chains.Chain.
func (an *SuiRemotenet) XCall(ctx context.Context, targetChain chains.Chain, keyName string, _to string, data []byte, rollback []byte) (*chains.XCallResponse, error) {
	height, err := targetChain.(ibcLocal.Chain).Height(ctx)
	if err != nil {
		return nil, err
	}
	ctx, err = an.SendPacketXCall(ctx, keyName, _to, data, rollback)
	if err != nil {
		return nil, err
	}
	return an.FindTargetXCallMessage(ctx, targetChain, height, strings.Split(_to, "/")[1])
}

func (an *SuiRemotenet) CommitTx(ctx context.Context, wallet *account.Account, txBytes lib.Base64Data, signatures []any) (*types.SuiTransactionBlockResponse, error) {
	return an.client.ExecuteTransactionBlock(ctx, txBytes, signatures, &types.SuiTransactionBlockResponseOptions{
		ShowEffects:       true,
		ShowEvents:        true,
		ShowObjectChanges: true,
	}, types.TxnRequestTypeWaitForLocalExecution)
}

func (an *SuiRemotenet) getGasCoinId(ctx context.Context, addr string, gasCost uint64) *types.Coin {
	accountAddress, err := move_types.NewAccountAddressHex(addr)
	if err != nil {
		an.log.Error(fmt.Sprintf("error getting account address sender %s", addr), zap.Error(err))
		return nil
	}
	result, err := an.client.GetSuiCoinsOwnedByAddress(ctx, *accountAddress)
	if err != nil {
		an.log.Error(fmt.Sprintf("error getting gas coins for address %s", addr), zap.Error(err))
		return nil
	}
	_, t, err := result.PickSUICoinsWithGas(big.NewInt(baseSuiFee), gasCost, pickMethod)
	if err != nil {
		an.log.Error(fmt.Sprintf("error getting gas coins with enough gas for address %s", addr), zap.Error(err))
		return nil
	}
	return t
}

func (an *SuiRemotenet) getAnotherGasCoinId(ctx context.Context, addr string, gasCost uint64, existingGasAddress move_types.AccountAddress) *types.Coin {
	accountAddress, err := move_types.NewAccountAddressHex(addr)
	if err != nil {
		an.log.Error(fmt.Sprintf("error getting account address sender %s", addr), zap.Error(err))
		return nil
	}
	coins, err := an.client.GetAllCoins(ctx, *accountAddress, nil, 1000)
	if err != nil {
		an.log.Error(fmt.Sprintf("error getting gas coins for address %s", addr), zap.Error(err))
		return nil
	}
	for _, coin := range coins.Data {
		if coin.Balance.Uint64() > gasCost && coin.CoinObjectId != existingGasAddress {
			return &coin
		}
	}
	return nil
}

func (an *SuiRemotenet) ExecuteContractRemote(ctx context.Context, suiMessage *SuiMessage, address string, gasBudget uint64) (*types.TransactionBytes, error) {
	accountAddress, err := move_types.NewAccountAddressHex(an.testconfig.RelayWalletAddress)
	if err != nil {
		return &types.TransactionBytes{}, fmt.Errorf("error getting account address sender: %w", err)
	}
	packageId, err := move_types.NewAccountAddressHex(suiMessage.PackageObjectId)
	if err != nil {
		return &types.TransactionBytes{}, fmt.Errorf("invalid packageId: %w", err)
	}
	coinId := an.getGasCoinId(ctx, an.testconfig.RelayWalletAddress, gasBudget)
	coinAddress, err := move_types.NewAccountAddressHex(coinId.CoinObjectId.String())
	if err != nil {
		return &types.TransactionBytes{}, fmt.Errorf("error getting gas coinid : %w", err)
	}
	typeArgs := []string{}
	var stringParams []interface{}
	for _, s := range suiMessage.Params {
		stringParams = append(stringParams, fmt.Sprint(s))
	}
	if stringParams == nil {
		stringParams = make([]interface{}, 0)
	}

	resp := types.TransactionBytes{}
	err = an.client.CallContext(
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
	return &resp, err
}

func (an *SuiRemotenet) QueryContractRemote(ctx context.Context, suiMessage *SuiMessage, address string, gasBudget uint64) (any, error) {
	builder := sui_types.NewProgrammableTransactionBuilder()
	packageId, err := move_types.NewAccountAddressHex(suiMessage.PackageObjectId)
	if err != nil {
		return nil, err
	}
	senderAddress, err := move_types.NewAccountAddressHex(address)
	if err != nil {
		return nil, err
	}
	callArgs, err := paramsToCallArgs(suiMessage)
	if err != nil {
		return nil, err
	}
	err = builder.MoveCall(
		*packageId,
		move_types.Identifier(suiMessage.Module),
		move_types.Identifier(suiMessage.Method),
		[]move_types.TypeTag{},
		callArgs,
	)
	if err != nil {
		return nil, err
	}
	transaction := builder.Finish()
	bcsBytes, err := bcs.Marshal(transaction)
	if err != nil {
		return nil, err
	}
	txBytes := append([]byte{0}, bcsBytes...)
	b64Data, err := lib.NewBase64Data(base64.StdEncoding.EncodeToString(txBytes))
	if err != nil {
		return nil, err
	}
	res, err := an.client.DevInspectTransactionBlock(context.Background(), *senderAddress, *b64Data, nil, nil)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, fmt.Errorf("error occurred while calling sui contract: %s", *res.Error)
	}
	result := (res.Results[0].ReturnValues[0]).([]interface{})
	resultType := result[1]
	byteSlice, ok := result[0].([]byte)
	if !ok {
		return nil, err
	}
	return extractResult(resultType, byteSlice, result[0])
}

func extractResult(resultType interface{}, byteSlice []byte, defResult interface{}) (any, error) {
	switch resultType {
	case suiU64:
		var u64Value uint64
		bcs.Unmarshal(byteSlice, &u64Value)
		return u64Value, nil
	case suiStringType:
		var strValue string
		bcs.Unmarshal(byteSlice, &strValue)
		return strValue, nil
	case suiBool:
		var booleanValue bool
		bcs.Unmarshal(byteSlice, &booleanValue)
		return booleanValue, nil
	default:
		return defResult, nil
	}
}

// convert native params to bcs encoded params
func paramsToCallArgs(suiMessage *SuiMessage) ([]sui_types.CallArg, error) {
	var callArgs []sui_types.CallArg
	for _, param := range suiMessage.Params {
		byteParam, err := bcs.Marshal(param)
		if err != nil {
			return nil, err
		}
		callArgs = append(callArgs, sui_types.CallArg{
			Pure: &byteParam,
		})
	}
	return callArgs, nil
}
