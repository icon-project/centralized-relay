package sui

import (
	"context"
	"encoding/base64"
	"encoding/hex"
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
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/dockerutil"
	ibcLocal "github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer/centralized"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const (
	suiCurrencyType                             = "0x2::sui::SUI"
	pickMethod                                  = 1
	baseSuiFee                                  = 1000
	suiStringType                               = "0x1::string::String"
	suiU64                                      = "u64"
	suiBool                                     = "bool"
	moveCall          suisdkClient.UnsafeMethod = "moveCall"
	publish           suisdkClient.UnsafeMethod = "publish"
	queryEvents       suisdkClient.SuiXMethod   = "queryEvents"
	callGasBudget                               = 500000000
	deployGasBudget                             = "500000000"
	xcallAdmin                                  = "xcall-admin"
	xcallStorage                                = "xcall-storage"
	sui_rlp_path                                = "libs/sui_rlp"
	adminCap                                    = "AdminCap"
	upgradeCap                                  = "UpgradeCap"
	connectionCap                               = "ConnCap"
	IdCapSuffix                                 = "-idcap"
	StateSuffix                                 = "-state"
	WitnessSuffix                               = "-witness"
	witnessCarrier                              = "WitnessCarrier"
	storage                                     = "Storage"
	CentralConnModule                           = "centralized_entry"
	MockAppModule                               = "mock_dapp"
	RegisterXcall                               = "register_xcall"
	CallArgObject                               = "object"
	CallArgPure                                 = "pure"
	connectionName                              = "centralized-1"
)

func NewSuiRemotenet(testName string, log *zap.Logger, chainConfig chains.ChainConfig, client *client.Client, network string, testconfig *testconfig.Chain) chains.Chain {
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

func (an *SuiRemotenet) Config() chains.ChainConfig {
	return an.cfg
}

// DeployContract implements chains.Chain.
func (an *SuiRemotenet) DeployContract(ctx context.Context, keyName string) (context.Context, error) {
	filePath := "/xcall/" + keyName
	stdout, _, err := an.ExecBin(ctx, "sui", "client", "publish", filePath, "--skip-dependency-verification", "--gas-budget", deployGasBudget, "--json")
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
		if changes.Data.Created != nil && strings.Contains(changes.Data.Created.ObjectType, adminCap) {
			depoymentInfo.AdminCap = changes.Data.Created.ObjectId.String()
		}
		if changes.Data.Created != nil && strings.Contains(changes.Data.Created.ObjectType, upgradeCap) {
			depoymentInfo.UpgradeCap = changes.Data.Created.ObjectId.String()
		}
		if changes.Data.Created != nil && strings.Contains(changes.Data.Created.ObjectType, storage) {
			depoymentInfo.Storage = changes.Data.Created.ObjectId.String()
		}
		if changes.Data.Created != nil && strings.Contains(changes.Data.Created.ObjectType, witnessCarrier) {
			depoymentInfo.Witness = changes.Data.Created.ObjectId.String()
		}
	}
	return context.WithValue(ctx, "objId", depoymentInfo), nil
}

// DeployXCallMockApp implements chains.Chain.
func (an *SuiRemotenet) DeployXCallMockApp(ctx context.Context, keyName string, connections []chains.XCallConnection) error {
	if an.testconfig.Environment == "preconfigured" {
		return nil
	}
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	ctx, err := an.DeployContract(ctx, MockAppModule)
	if err != nil {
		return err
	}
	deploymentInfo := ctx.Value("objId").(DepoymentInfo)
	an.IBCAddresses[dappKey] = deploymentInfo.PackageId
	an.IBCAddresses[dappKey+WitnessSuffix] = deploymentInfo.Witness
	an.log.Info("setup Dapp completed ", zap.Any("packageId", deploymentInfo.PackageId), zap.Any("witness", deploymentInfo.Witness))

	// register xcall
	params := []SuiCallArg{
		{Type: CallArgObject, Val: an.IBCAddresses[xcallStorage]},
		{Type: CallArgObject, Val: an.IBCAddresses[dappKey+WitnessSuffix]},
	}

	msg := an.NewSuiMessage(params, an.IBCAddresses[dappKey], MockAppModule, RegisterXcall)
	resp, err := an.callContract(ctx, msg)
	for _, changes := range resp.ObjectChanges {
		if changes.Data.Created != nil && strings.Contains(changes.Data.Created.ObjectType, "DappState") {
			an.IBCAddresses[dappKey+StateSuffix] = changes.Data.Created.ObjectId.String()
			time.Sleep(2 * time.Second)
			response, err := an.client.GetObject(ctx, changes.Data.Created.ObjectId, &types.SuiObjectDataOptions{
				ShowContent: true,
			})
			if err != nil {
				return err
			}
			fields := response.Data.Content.Data.MoveObject.Fields.(map[string]interface{})
			js, _ := json.Marshal(fields)
			var objRes ObjectResult
			json.Unmarshal(js, &objRes)
			an.IBCAddresses[dappKey+IdCapSuffix] = objRes.XcallCap.Fields.ID.ID
		}
	}

	an.log.Info("register xcall completed ", zap.Any("dapp-state", an.IBCAddresses[dappKey+StateSuffix]), zap.Any("dapp-state-idcap", an.IBCAddresses[dappKey+IdCapSuffix]))
	if err != nil {
		return err
	}
	// add connections
	for _, connection := range connections {
		params = []SuiCallArg{
			{Type: CallArgObject, Val: an.IBCAddresses[dappKey+StateSuffix]},
			{Type: CallArgPure, Val: connection.Nid},
			{Type: CallArgPure, Val: []string{connectionName}},
			{Type: CallArgPure, Val: []string{connection.Destination}},
		}

		msg = an.NewSuiMessage(params, an.IBCAddresses[dappKey], MockAppModule, "add_connection")
		resp, err = an.callContract(ctx, msg)
		for _, changes := range resp.ObjectChanges {
			if changes.Data.Created != nil && strings.Contains(changes.Data.Created.ObjectType, "DappState") {
				an.IBCAddresses[dappKey+connection.Connection+StateSuffix] = changes.Data.Created.ObjectId.String()
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
	job := dockerutil.NewImage(an.log, an.DockerClient, an.Network, an.testName, an.cfg.Images.Repository, an.cfg.Images.Version)

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

// FindCallMessage implements chains.Chain.
func (an *SuiRemotenet) FindCallMessage(ctx context.Context, startHeight uint64, from string, to string, sn string) (string, string, error) {
	index := sn
	event, err := an.FindEvent(ctx, startHeight, "xcall", index, "::main::CallMessage", CentralConnModule)
	if err != nil {
		return "", "", err
	}
	jsonData := (event.ParsedJson.(map[string]interface{}))
	data := jsonData["data"].([]interface{})
	valueSlice := make([]byte, len(data))
	for i, v := range data {
		valueSlice[i] = byte(v.(float64))
	}
	return jsonData["sn"].(string), string(valueSlice), nil
}

// FindRollbackExecutedMessage implements chains.Chain.
func (an *SuiRemotenet) FindRollbackExecutedMessage(ctx context.Context, startHeight uint64, sn string) (string, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	index := sn
	event, err := an.FindEvent(ctx, startHeight, dappKey, index, "::main::RollbackExecuted", MockAppModule)
	if err != nil {
		return "", err
	}
	jsonData := (event.ParsedJson.(map[string]interface{}))
	return jsonData["sn"].(string), nil
}

func (an *SuiRemotenet) getEvent(ctx context.Context, sn, eventType, module, packageKey string) (*types.SuiEvent, error) {
	limit := uint(100)
	query := MoveEventRequest{
		MoveModule: MoveModule{
			Package: an.IBCAddresses[packageKey],
			Module:  module,
		},
	}
	var resp types.EventPage
	err := an.client.CallContext(ctx, &resp, queryEvents, query, nil, limit, true)

	if err != nil {
		return nil, err
	}
	for _, event := range resp.Data {
		jsonData := (event.ParsedJson.(map[string]interface{}))
		if jsonData["sn"] == sn && strings.Contains(event.Type, eventType) {
			return &event, nil
		}
	}
	return nil, fmt.Errorf("event not found")
}

func (an *SuiRemotenet) FindEvent(ctx context.Context, startHeight uint64, packageKey, index, eventType, module string) (*types.SuiEvent, error) {
	timeout := time.After(60 * time.Second)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("failed to find eventLog")
		case <-ticker.C:
			data, err := an.getEvent(ctx, index, eventType, module, packageKey)
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
	index := sn
	event, err := an.FindEvent(ctx, startHeight, "xcall", index, "::main::ResponseMessage", CentralConnModule)
	if err != nil {
		return "", err
	}
	jsonData := (event.ParsedJson.(map[string]interface{}))
	responseCode := jsonData["response_code"].(float64)
	return strconv.FormatFloat(responseCode, 'f', -1, 64), nil

}

// FindTargetXCallMessage implements chains.Chain.
func (an *SuiRemotenet) FindTargetXCallMessage(ctx context.Context, target chains.Chain, height uint64, to string) (*chains.XCallResponse, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	sn := ctx.Value("sn").(string)
	reqId, destData, err := target.FindCallMessage(ctx, height, an.cfg.ChainID+"/"+an.IBCAddresses[dappKey+IdCapSuffix], to, sn)
	return &chains.XCallResponse{SerialNo: sn, RequestID: reqId, Data: destData}, err
}

// GetContractAddress implements chains.Chain.
func (an *SuiRemotenet) GetContractAddress(key string) string {
	if key == "connection" {
		return connectionName
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

// GetHostRPCAddress implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).GetHostRPCAddress of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) GetHostRPCAddress() string {
	return an.testconfig.RPCUri
}

// GetRPCAddress implements chains.Chain.
// Subtle: this method shadows the method (*CosmosChain).GetRPCAddress of SuiRemotenet.CosmosChain.
func (an *SuiRemotenet) GetRPCAddress() string {
	return an.testconfig.RPCUri
}

// GetRelayConfig implements chains.Chain.
func (an *SuiRemotenet) GetRelayConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	contracts := make(map[string]string)
	contracts["xcall"] = an.GetContractAddress("xcall")
	dappModule := centralized.SuiDappModule{
		Name:     MockAppModule,
		CapId:    an.GetContractAddress(dappKey + IdCapSuffix)[2:],
		ConfigId: an.GetContractAddress(dappKey + StateSuffix),
	}
	config := &centralized.SUIRelayerChainConfig{
		Type: "sui",
		Value: centralized.SUIRelayerChainConfigValue{
			NID:             an.Config().ChainID,
			RPCURL:          an.GetRPCAddress(),
			WebsocketUrl:    an.testconfig.WebsocketUrl,
			XcallPkgIds:     []string{an.GetContractAddress("xcall")},
			XcallStorageId:  an.GetContractAddress(xcallStorage),
			ConnectionId:    an.GetContractAddress("connection"),
			ConnectionCapId: an.GetContractAddress("connectionCap"),
			DappPkgId:       an.GetContractAddress(dappKey),
			DappModules: []centralized.SuiDappModule{
				dappModule,
			},
			StartHeight:   0,
			BlockInterval: "0s",
			Address:       an.testconfig.RelayWalletAddress,
			FinalityBlock: uint64(0),
			GasPrice:      4000000,
			GasLimit:      50000000,
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
	}
	gasFeeCoin := an.getGasCoinId(ctx, an.testconfig.RelayWalletAddress, callGasBudget).CoinObjectId
	coinId := an.getAnotherGasCoinId(ctx, an.testconfig.RelayWalletAddress, callGasBudget, gasFeeCoin)
	params := []SuiCallArg{
		{Type: CallArgObject, Val: an.IBCAddresses[dappKey+StateSuffix]},
		{Type: CallArgObject, Val: an.IBCAddresses[xcallStorage]},
		{Type: CallArgObject, Val: coinId.CoinObjectId},
		{Type: CallArgPure, Val: _to},
		{Type: CallArgPure, Val: "0x" + hex.EncodeToString(data)},
		{Type: CallArgPure, Val: "0x" + hex.EncodeToString(rollback)},
	}
	msg := an.NewSuiMessage(params, an.IBCAddresses[dappKey], MockAppModule, "send_message")
	resp, err := an.callContract(ctx, msg)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, "sn", an.findSn(resp, "::main::CallMessageSent")), nil
}

func (an *SuiRemotenet) findSn(tx *types.SuiTransactionBlockResponse, eType string) string {
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
	if an.testconfig.Environment == "preconfigured" {
		return nil
	}
	return nil
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
	if !resp.Effects.Data.IsSuccess() {
		if strings.Contains(resp.Effects.Data.V1.Status.Error, "send_call_inner") {
			return resp, fmt.Errorf("MaxDataSizeExceeded")
		}
		return resp, fmt.Errorf("txn execution failed")
	}
	return resp, nil
}

// SetupXCall implements chains.Chain.
func (an *SuiRemotenet) SetupXCall(ctx context.Context) error {
	if an.testconfig.Environment == "preconfigured" {
		testcase := ctx.Value("testcase").(string)
		an.IBCAddresses["xcall"] = "0x024bb2ef51d49bace743297d18a7ce17ee0a39d541832826f9d1ed9516b31a95"
		an.IBCAddresses[xcallAdmin] = "0xe1c8d2e988ffc4ba8080dce6e490192c63717d25899025d97eab99e8548d8862"
		an.IBCAddresses[xcallStorage] = "0xbc74501f4a771b9126e67d4eec781c13dcc6c790d54af789e96f9264b1eed2d8"
		an.IBCAddresses["connectionCap"] = "0xba6d692ff0e2ef4deb2676f24fda5bced79abac2026b4b1f70c51f3651da02d2"
		dappKey := fmt.Sprintf("dapp-%s", testcase)
		an.IBCAddresses[dappKey] = "0x42c6f06edd37db92d68a6cd08713972fae7be8e7e358678400cb10d3657fa210"
		an.IBCAddresses[dappKey+WitnessSuffix] = "0xa43e19bc8e56e236e4756e23ed011d1d65980b835778ac2d4f7c77d21ab0c2c9"
		an.IBCAddresses[dappKey+StateSuffix] = "0xa28f7dda0671e405e0684867de55c9f3fdc7bdf50bba2e9d485ab2ad081d6a9a"
		an.IBCAddresses[dappKey+IdCapSuffix] = "0x33c9150bca7c5f05b3e20ac7c82efdff0ceb4a73d7c7ed907067f4787b5534b1"
		return nil
	}
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
	an.log.Info("setup xcall completed ", zap.Any("packageId", deploymentInfo.PackageId),
		zap.Any("admin", deploymentInfo.AdminCap), zap.Any("storage", deploymentInfo.Storage),
		zap.Any("upgradeCap", deploymentInfo.UpgradeCap))
	//configuing nid
	//init
	params := []SuiCallArg{
		{Type: CallArgObject, Val: an.IBCAddresses[xcallStorage]},
		{Type: CallArgObject, Val: deploymentInfo.UpgradeCap},
		{Type: CallArgPure, Val: "sui"},
	}
	msg := an.NewSuiMessage(params, an.IBCAddresses["xcall"], "main", "configure_nid")
	_, err = an.callContract(ctx, msg)
	if err != nil {
		return err
	}
	//init
	params = []SuiCallArg{
		{Type: CallArgObject, Val: an.IBCAddresses[xcallStorage]},
		{Type: CallArgObject, Val: an.IBCAddresses[xcallAdmin]},
		{Type: CallArgPure, Val: connectionName},
		{Type: CallArgPure, Val: an.testconfig.RelayWalletAddress},
	}
	msg = an.NewSuiMessage(params, an.IBCAddresses["xcall"], "main", "register_connection_admin")
	resp, err := an.callContract(ctx, msg)
	if err != nil {
		return err
	}
	for _, changes := range resp.ObjectChanges {
		if changes.Data.Created != nil && strings.Contains(changes.Data.Created.ObjectType, connectionCap) {
			an.IBCAddresses["connectionCap"] = changes.Data.Created.ObjectId.String()
		}

	}
	an.log.Info("connection registered", zap.Any("connectionCap", an.IBCAddresses["connectionCap"]))
	return err
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
	height, err := targetChain.Height(ctx)
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
	var args []interface{}
	for _, param := range suiMessage.Params {
		args = append(args, param.Val)
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
		args,
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
