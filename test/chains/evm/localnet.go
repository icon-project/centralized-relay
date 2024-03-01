package evm

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	interchaintest "github.com/icon-project/centralized-relay/test"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/blockdb"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/dockerutil"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer/centralized"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"gopkg.in/yaml.v3"

	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/icon-project/centralized-relay/test/chains"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"

	"go.uber.org/zap"
)

type Event struct {
	name     string
	hash     string
	contract string
}

var (
	CallMessageSent  = Event{"CallMessageSent", "0x69e53ea70fdf945f6d035b3979748bc999151691fb1dc69d66f8017f8840ae28", "xcall"}
	CallMessage      = Event{"CallMessage", "0x2cbc78425621c181f9f8a25fc06e44a0ac2b67cd6a31f8ed7918934187f8cc59", "xcall"}
	ResponseMessage  = Event{"ResponseMessage", "0xbeacafd006c5e60667f6f04aec3a498f81c8e94142b4e95b5a5a763de43ca0ab", "xcall"}
	RollbackExecuted = Event{"RollbackExecuted", "0x08f0ac7aef6da8bbe43bee8b1444a1883f1359566618bc379ce5abba44883837", "xcall"}
)

type EVMRemotenet struct {
	log           *zap.Logger
	testName      string
	cfg           ibc.ChainConfig
	numValidators int
	numFullNodes  int
	FullNodes     HardhatNodes
	findTxMu      sync.Mutex
	privateKey    string
	scorePaths    map[string]string
	IBCAddresses  map[string]string     `json:"addresses"`
	Wallets       map[string]ibc.Wallet `json:"wallets"`
	DockerClient  *client.Client
	ContractABI   map[string]abi.ABI
	Network       string
	testconfig    *testconfig.Chain
	Client        ethclient.Client
}

func (an *EVMRemotenet) CreateKey(ctx context.Context, keyName string) error {
	panic("implement me")
}

// func NewEVMRemotenet(testName string, log *zap.Logger, chainConfig ibc.ChainConfig, numValidators int, numFullNodes int, scorePaths map[string]string) chains.Chain {
func NewEVMRemotenet(testName string, log *zap.Logger, chainConfig ibc.ChainConfig, client *client.Client, network string, testconfig *testconfig.Chain) chains.Chain {
	ethclient, err := ethclient.Dial(testconfig.RPCUri)
	if err != nil {
		fmt.Println(err)
	}
	contractABI := make(map[string]abi.ABI)

	return &EVMRemotenet{
		testName:      testName,
		cfg:           chainConfig,
		numValidators: 0,
		numFullNodes:  0,
		log:           log,
		scorePaths:    testconfig.Contracts,
		Wallets:       map[string]ibc.Wallet{},
		IBCAddresses:  make(map[string]string),
		Client:        *ethclient,
		DockerClient:  client,
		ContractABI:   contractABI,
		testconfig:    testconfig,
		Network:       network,
	}
}

// Config fetches the chain configuration.
func (an *EVMRemotenet) Config() ibc.ChainConfig {
	return an.cfg
}

// Initialize initializes node structs so that things like initializing keys can be done before starting the chain
func (an *EVMRemotenet) Initialize(ctx context.Context, testName string, cli *client.Client, networkID string) error {
	return nil
}

func (an *EVMRemotenet) NewChainNode(
	ctx context.Context,
	testName string,
	cli *client.Client,
	networkID string,
	image ibc.DockerImage,
	validator bool,
) (*AnvilNode, error) {
	// Construct the ChainNode first so we can access its name.
	// The ChainNode's VolumeName cannot be set until after we create the volume.

	in := &AnvilNode{
		log:          an.log,
		Chain:        an,
		DockerClient: cli,
		NetworkID:    networkID,
		TestName:     testName,
		Image:        image,
		ContractABI:  make(map[string]abi.ABI),
	}

	v, err := cli.VolumeCreate(ctx, volumetypes.CreateOptions{
		Labels: map[string]string{
			dockerutil.CleanupLabel: testName,

			dockerutil.NodeOwnerLabel: in.Name(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("creating volume for chain node: %w", err)
	}
	in.VolumeName = v.Name

	if err := dockerutil.SetVolumeOwner(ctx, dockerutil.VolumeOwnerOptions{
		Log: an.log,

		Client: cli,

		VolumeName: v.Name,
		ImageRef:   image.Ref(),
		TestName:   testName,
		UidGid:     image.UidGid,
	}); err != nil {
		return nil, fmt.Errorf("set volume owner: %w", err)
	}
	return in, nil
}

// Start sets up everything needed (validators, gentx, fullnodes, peering, additional accounts) for chain to start from genesis.
func (an *EVMRemotenet) Start(testName string, ctx context.Context, additionalGenesisWallets ...ibc.WalletAmount) error {
	return nil
}

// Exec runs an arbitrary command using Chain's docker environment.
// Whether the invoked command is run in a one-off container or execing into an already running container
// is up to the chain implementation.
//
// "env" are environment variables in the format "MY_ENV_VAR=value"
func (an *EVMRemotenet) Exec(ctx context.Context, cmd []string, env []string) (stdout []byte, stderr []byte, err error) {
	job := dockerutil.NewImage(an.log, an.DockerClient, an.Network, an.testName, an.cfg.Images[0].Repository, an.cfg.Images[0].Version)

	opts := dockerutil.ContainerOptions{
		Binds: []string{},
	}
	res := job.Run(ctx, cmd, opts)
	return res.Stdout, res.Stderr, res.Err
	// return an.getFullNode().Exec(ctx, cmd, env)
}

// ExportState exports the chain state at specific height.
func (an *EVMRemotenet) ExportState(ctx context.Context, height int64) (string, error) {
	block, err := an.getFullNode().GetBlockByHeight(ctx, height)
	return block, err
}

// GetRPCAddress retrieves the rpc address that can be reached by other containers in the docker network.
func (an *EVMRemotenet) GetRPCAddress() string {
	return an.testconfig.RPCUri
}

// GetGRPCAddress retrieves the grpc address that can be reached by other containers in the docker network.
// Not Applicable for Icon
func (an *EVMRemotenet) GetGRPCAddress() string {
	return an.testconfig.RPCUri
}

// GetHostRPCAddress returns the rpc address that can be reached by processes on the host machine.
// Note that this will not return a valid value until after Start returns.
func (an *EVMRemotenet) GetHostRPCAddress() string {
	return an.testconfig.RPCUri
}

// GetHostGRPCAddress returns the grpc address that can be reached by processes on the host machine.
// Note that this will not return a valid value until after Start returns.
// Not applicable for Icon
func (an *EVMRemotenet) GetHostGRPCAddress() string {
	return ""
}

// HomeDir is the home directory of a node running in a docker container. Therefore, this maps to
// the container's filesystem (not the host).
func (an *EVMRemotenet) HomeDir() string {
	return an.FullNodes[0].HomeDir()
}

func (an *EVMRemotenet) createKeystore(ctx context.Context, keyName string) (string, string, error) {
	keydir := keystore.NewKeyStore("/tmp", keystore.StandardScryptN, keystore.StandardScryptP)

	account, err := keydir.NewAccount(keyName)

	if err != nil {
		an.log.Fatal("unable to create new account", zap.Error(err))
		return "", "", err
	}

	keyJSON, err := keydir.Export(account, keyName, keyName)
	err = an.FullNodes[0].RestoreKeystore(ctx, keyJSON, keyName)
	if err != nil {
		an.log.Error("fail to restore keystore", zap.Error(err))
		return "", "", err
	}
	key, err := keystore.DecryptKey(keyJSON, keyName)
	privateKey := crypto.FromECDSA(key.PrivateKey)
	return account.Address.Hex(), hex.EncodeToString(privateKey), nil
}

// RecoverKey recovers an existing user from a given mnemonic.
func (an *EVMRemotenet) RecoverKey(ctx context.Context, name string, mnemonic string) error {
	panic("not implemented") // TODO: Implement
}

// GetAddress fetches the bech32 address for a test key on the "user" node (either the first fullnode or the first validator if no fullnodes).
func (an *EVMRemotenet) GetAddress(ctx context.Context, keyName string) ([]byte, error) {
	addrInByte, err := json.Marshal(keyName)
	if err != nil {
		return nil, err
	}
	return addrInByte, nil
}

func (an *EVMRemotenet) GetRelayConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error) {
	// gasPrice, _ := an.Client.SuggestGasPrice(context.Background())
	contracts := make(map[string]string)
	contracts["xcall"] = an.GetContractAddress("xcall")
	contracts["connection"] = an.GetContractAddress("connection")
	config := &centralized.EVMRelayerChainConfig{
		Type: "evm",
		Value: centralized.EVMRelayerChainConfigValue{
			NID:           an.Config().ChainID,
			RPCURL:        an.GetRPCAddress(),
			StartHeight:   0,
			GasPrice:      10000000,
			GasLimit:      20000000,
			Contracts:     contracts,
			BlockInterval: "6s",
			Address:       an.testconfig.RelayWalletAddress,
		},
	}
	return yaml.Marshal(config)
}

// SendFunds sends funds to a wallet from a user account.
func (an *EVMRemotenet) SendFunds(ctx context.Context, keyName string, amount ibc.WalletAmount) error {
	an.CheckForKeyStore(ctx, keyName)

	privateKey, err := crypto.HexToECDSA(an.Wallets[interchaintest.FaucetAccountKeyName].Mnemonic())
	if err != nil {
		return err
	}
	ethClient := an.getFullNode().Client
	// Get the public key and address from the private key

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
	fmt.Printf("fromAddress %s\n", fromAddress)
	if err != nil {
		return err
	}
	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)                // in units
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	toAddress := common.HexToAddress(amount.Address)
	var data []byte
	txdata := &eth_types.LegacyTx{
		To:       &toAddress,
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		Value:    value,
		Data:     data,
	}
	tx := eth_types.NewTx(txdata)
	networkID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	signedTx, err := eth_types.SignTx(tx, eth_types.NewEIP155Signer(networkID), privateKey)
	if err != nil {
		fmt.Println(err)
	}
	err = ethClient.SendTransaction(context.Background(), signedTx)

	return err

}

// Height returns the current block height or an error if unable to get current height.
func (an *EVMRemotenet) Height(ctx context.Context) (uint64, error) {
	header, err := an.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

// GetGasFeesInNativeDenom gets the fees in native denom for an amount of spent gas.
func (an *EVMRemotenet) GetGasFeesInNativeDenom(gasPaid int64) int64 {
	gasPrice, _ := strconv.ParseFloat(strings.Replace(an.cfg.GasPrices, an.cfg.Denom, "", 1), 64)
	fees := float64(gasPaid) * gasPrice
	return int64(fees)
}

// BuildRelayerWallet will return a chain-specific wallet populated with the mnemonic so that the wallet can
// be restored in the relayer node using the mnemonic. After it is built, that address is included in
// genesis with some funds.
func (an *EVMRemotenet) BuildRelayerWallet(ctx context.Context, keyName string) (ibc.Wallet, error) {
	return an.BuildWallet(ctx, keyName, "")
}

func (an *EVMRemotenet) BuildWallet(ctx context.Context, keyName string, mnemonic string) (ibc.Wallet, error) {
	address, privateKey, err := an.createKeystore(ctx, keyName)
	if err != nil {
		return nil, fmt.Errorf("failed to create key with name %q on chain %s: %w", keyName, an.cfg.Name, err)

	}

	w := NewWallet(keyName, []byte(address), privateKey, an.cfg)
	an.Wallets[keyName] = w
	return w, nil
}

func (an *EVMRemotenet) getFullNode() *AnvilNode {
	an.findTxMu.Lock()
	defer an.findTxMu.Unlock()
	if len(an.FullNodes) > 0 {
		// use first full node
		return an.FullNodes[0]
	}
	return an.FullNodes[0]
}

func (an *EVMRemotenet) FindTxs(ctx context.Context, height uint64) ([]blockdb.Tx, error) {
	return nil, nil
}

// GetBalance fetches the current balance for a specific account address and denom.
func (an *EVMRemotenet) GetBalance(ctx context.Context, address string, denom string) (int64, error) {
	return an.getFullNode().GetBalance(ctx, address)
}

func (an *EVMRemotenet) SetupConnection(ctx context.Context, target chains.Chain) error {
	//testcase := ctx.Value("testcase").(string)
	xcall := common.HexToAddress(an.IBCAddresses["xcall"])
	// _ =an.CheckForKeyStore(ctx, keyName)
	connection, err := an.DeployContractRemote(ctx, an.scorePaths["connection"], an.testconfig.KeystorePassword)
	if err != nil {
		return err
	}
	// relayerKey := fmt.Sprintf("relayer-%s",an.Config().Name)
	// relayerAddress := common.HexToAddress(c.Wallets[relayerKey].FormattedAddress())
	relayerAddress := common.HexToAddress(an.testconfig.RelayWalletAddress)

	_, err = an.ExecCallTx(ctx, connection.Hex(), "initialize", an.testconfig.KeystorePassword, relayerAddress, xcall)
	if err != nil {
		fmt.Printf("fail to initialized xcall-adapter : %w\n", err)
		return err
	}

	// _ =an.CheckForKeyStore(ctx, relayerKey)

	// _, err =an.ExecCallTx(ctx, connection.Hex(), "setFee", privateKey, target.Config().ChainID, big.NewInt(0), big.NewInt(0))
	// if err != nil {
	// 	fmt.Printf("fail to initialized fee for xcall-adapter : %w\n", err)
	// 	return err
	// }

	an.IBCAddresses["connection"] = connection.Hex()
	return nil
}

func (an *EVMRemotenet) SetupXCall(ctx context.Context) error {
	//testcase := ctx.Value("testcase").(string)
	nid := an.cfg.ChainID
	//ibcAddress :=an.IBCAddresses["ibc"]
	// _ =an.CheckForKeyStore(ctx, keyName)
	xcall, err := an.DeployContractRemote(ctx, an.scorePaths["xcall"], an.testconfig.KeystorePassword)
	if err != nil {
		return err
	}
	_, err = an.ExecCallTx(ctx, xcall.Hex(), "initialize", an.testconfig.KeystorePassword, nid)
	if err != nil {
		fmt.Printf("fail to initialized xcall : %w\n", err)
		return err
	}

	an.IBCAddresses["xcall"] = xcall.Hex()
	return nil
}

func (an *EVMRemotenet) DeployXCallMockApp(ctx context.Context, keyName string, connections []chains.XCallConnection) error {
	testcase := ctx.Value("testcase").(string)
	//an.CheckForKeyStore(ctx, keyName)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	xCall := an.GetContractAddress("xcall")
	dapp, err := an.DeployContractRemote(ctx, an.scorePaths["dapp"], an.testconfig.KeystorePassword)
	if err != nil {
		return err
	}
	//_, err =an.executeContract(ctx,dapp,keyName,"initialize",xCall)
	_, err = an.ExecCallTx(ctx, dapp.Hex(), "initialize", an.testconfig.KeystorePassword, common.HexToAddress(xCall))
	if err != nil {
		fmt.Printf("fail to initialized dapp : %w\n", err)
		return err
	}

	an.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = dapp.Hex()

	for _, connection := range connections {
		//connectionKey := fmt.Sprintf("%s-%s", connection.Connection, testcase)
		params := []interface{}{connection.Nid, an.IBCAddresses[connection.Connection], connection.Destination}
		ctx, err = an.executeContract(context.Background(), dapp.Hex(), keyName, "addConnection", params...)
		if err != nil {
			an.log.Error("Unable to add connection",
				zap.Error(err),
				zap.String("nid", connection.Nid),
				zap.String("source", an.IBCAddresses[connection.Connection]),
				zap.String("destination", connection.Destination),
			)
		}
	}

	return nil
}

func (an *EVMRemotenet) GetContractAddress(key string) string {
	value, exist := an.IBCAddresses[key]
	if !exist {
		panic(fmt.Sprintf(`IBC address not exist %s`, key))
	}
	return value
}

func (an *EVMRemotenet) BackupConfig() ([]byte, error) {
	panic("not implemented")
}

func (an *EVMRemotenet) RestoreConfig(backup []byte) error {
	panic("not implemented")
}

func (an *EVMRemotenet) SendPacketXCall(ctx context.Context, keyName, _to string, data, rollback []byte) (context.Context, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	if rollback == nil {
		rollback = make([]byte, 0)
	}
	params := []interface{}{_to, data, rollback}
	// TODO: send fees
	ctx, err := an.executeContract(ctx, an.IBCAddresses[dappKey], keyName, "sendMessage", params...)
	if err != nil {
		return nil, err
	}
	txn := ctx.Value("txResult").(*eth_types.Receipt)

	events, err := an.ParseEvent(CallMessageSent, txn.Logs)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, "sn", events["_sn"].(*big.Int).String()), nil
}

// FindTargetXCallMessage returns the request id and the data of the message sent to the target chain
func (an *EVMRemotenet) FindTargetXCallMessage(ctx context.Context, target chains.Chain, height uint64, to string) (*chains.XCallResponse, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	sn := ctx.Value("sn").(string)
	reqId, destData, err := target.FindCallMessage(ctx, height, an.cfg.ChainID+"/"+an.IBCAddresses[dappKey], to, sn)
	return &chains.XCallResponse{SerialNo: sn, RequestID: reqId, Data: destData}, err
}

func (an *EVMRemotenet) XCall(ctx context.Context, targetChain chains.Chain, keyName, to string, data, rollback []byte) (*chains.XCallResponse, error) {
	height, err := targetChain.(ibc.Chain).Height(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: send fees

	ctx, err = an.SendPacketXCall(ctx, keyName, to, data, rollback)
	if err != nil {
		return nil, err
	}
	return an.FindTargetXCallMessage(ctx, targetChain, height, strings.Split(to, "/")[1])
}

func (an *EVMRemotenet) ExecuteCall(ctx context.Context, reqId, data string) (context.Context, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	_reqId, _ := big.NewInt(0).SetString(reqId, 10)
	return an.executeContract(ctx, an.IBCAddresses["xcall"], interchaintest.UserAccount, "executeCall", _reqId, []byte(data))
}

func (an *EVMRemotenet) ExecuteRollback(ctx context.Context, sn string) (context.Context, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	_sn, _ := big.NewInt(0).SetString(sn, 10)
	ctx, err := an.executeContract(ctx, an.IBCAddresses["xcall"], interchaintest.UserAccount, "executeRollback", _sn)
	if err != nil {
		return nil, err
	}
	txn := ctx.Value("txResult").(*eth_types.Receipt)
	events, err := an.ParseEvent(RollbackExecuted, txn.Logs)
	sequence := events["_sn"]
	return context.WithValue(ctx, "IsRollbackEventFound", fmt.Sprintf("%d", sequence) == sn), nil

}

func (an *EVMRemotenet) FindCallMessage(ctx context.Context, startHeight uint64, from, to, sn string) (string, string, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	fmt.Printf("%s,--%s,--%s\n", from, to, sn)
	_sn, _ := big.NewInt(0).SetString(sn, 10)
	topics := []common.Hash{common.HexToHash(CallMessage.hash), crypto.Keccak256Hash([]byte(from)), crypto.Keccak256Hash([]byte(to)), common.BytesToHash(_sn.Bytes())}
	event, err := an.FindEvent(ctx, startHeight, CallMessage, topics)

	if err != nil {
		fmt.Printf("Topics %v\n", topics)
		return "", "", err
	}
	return event["_reqId"].(*big.Int).String(), string(event["_data"].([]byte)), nil
}

func (an *EVMRemotenet) FindCallResponse(ctx context.Context, startHeight uint64, sn string) (string, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	_sn, _ := big.NewInt(0).SetString(sn, 10)
	topics := []common.Hash{common.HexToHash(ResponseMessage.hash), common.BytesToHash(_sn.Bytes())}

	event, err := an.FindEvent(ctx, startHeight, ResponseMessage, topics)
	if err != nil {
		fmt.Printf("Topics %v", topics)
		return "", err
	}

	return event["_code"].(*big.Int).String(), nil

}

func (an *EVMRemotenet) FindEvent(ctx context.Context, startHeight uint64, event Event, topics []common.Hash) (map[string]interface{}, error) {
	//eventSignature := []byte(event.hash)
	wsAddress := strings.Replace(an.testconfig.RPCUri, "http", "ws", 1)
	eClient, err := ethclient.Dial(wsAddress)
	defer eClient.Close()
	if err != nil {
		return nil, errors.New("error: fail to create eth client")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Starting block number
	fromBlock := new(big.Int).SetUint64(startHeight - 1)
	address := common.HexToAddress(an.IBCAddresses[event.contract])
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		Addresses: []common.Address{address},
		Topics: [][]common.Hash{
			topics,
		},
	}
	ev := make(map[string]interface{})
	maxIterations := 6
	iterations := 0
	for iterations < maxIterations {
		logs, err := eClient.FilterLogs(context.Background(), query)
		if err != nil {
			return ev, err
		}

		for _, lg := range logs {
			contractABI := an.ContractABI[address.Hex()]
			if len(lg.Data) > 0 {
				if err = contractABI.UnpackIntoMap(ev, event.name, lg.Data); err != nil {
					return ev, err
				}
			}

			var indexed abi.Arguments
			for _, arg := range contractABI.Events[event.name].Inputs {
				if arg.Indexed {
					indexed = append(indexed, arg)
				}
			}

			err = abi.ParseTopicsIntoMap(ev, indexed, lg.Topics[1:])
			return ev, err

		}
		iterations++
		time.Sleep(5 * time.Second)
	}
	return ev, errors.New("event not found after maximum iterations")
}

func (an *EVMRemotenet) ParseEvent(event Event, logs []*eth_types.Log) (map[string]interface{}, error) {
	contractABI := an.ContractABI[an.IBCAddresses[event.contract]]
	var ev = make(map[string]interface{})
	var err error
	for _, lg := range logs {
		if event.hash != lg.Topics[0].Hex() {
			continue
		}
		if len(lg.Data) > 0 {
			if err = contractABI.UnpackIntoMap(ev, event.name, lg.Data); err != nil {
				break
			}
		}
		var indexed abi.Arguments
		for _, arg := range contractABI.Events[event.name].Inputs {
			if arg.Indexed {
				indexed = append(indexed, arg)
			}
		}
		err = abi.ParseTopicsIntoMap(ev, indexed, lg.Topics[1:])

		break
	}
	return ev, err
}

// DeployContract implements chains.Chain
func (an *EVMRemotenet) DeployContract(ctx context.Context, keyName string) (context.Context, error) {
	// Get contract Name from context
	ctxValue := ctx.Value(chains.ContractName{}).(chains.ContractName)
	contractName := ctxValue.ContractName

	// Get Init Message from context
	ctxVal := ctx.Value(chains.InitMessageKey("init-msg")).(chains.InitMessage)

	initMessage := an.getInitParams(ctx, contractName, ctxVal.Message)

	var contracts chains.ContractKey

	// Check if 6123f953784d27e0729bc7a640d6ad8f04ed6710.keystore is alreadry available for given keyName
	ownerAddr := an.CheckForKeyStore(ctx, keyName)
	if ownerAddr != nil {
		contracts.ContractOwner = map[string]string{
			keyName: ownerAddr.FormattedAddress(),
		}
	}

	// Get ScoreAddress
	scoreAddress, err := an.getFullNode().DeployContract(ctx, an.scorePaths[contractName], an.privateKey, initMessage)

	contracts.ContractAddress = map[string]string{
		contractName: scoreAddress.Hex(),
	}

	testcase := ctx.Value("testcase").(string)
	contract := fmt.Sprintf("%s-%s", contractName, testcase)
	an.IBCAddresses[contract] = scoreAddress.Hex()
	return context.WithValue(ctx, chains.Mykey("contract Names"), chains.ContractKey{
		ContractAddress: contracts.ContractAddress,
		ContractOwner:   contracts.ContractOwner,
	}), err
}

// executeContract implements chains.Chain
func (an *EVMRemotenet) executeContract(ctx context.Context, contractAddress, keyName, methodName string, params ...interface{}) (context.Context, error) {
	//an.CheckForKeyStore(ctx, keyName)

	receipt, err := an.ExecCallTx(ctx, contractAddress, methodName, an.testconfig.KeystorePassword, params...)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Transaction Hash for %s: %s\n", methodName, receipt.TxHash.String())

	return context.WithValue(ctx, "txResult", receipt), nil

}

func (an *EVMRemotenet) ExecuteContract(ctx context.Context, contractAddress, keyName, methodName string, params map[string]interface{}) (context.Context, error) {
	execMethodName, execParams := an.getExecuteParam(ctx, methodName, params)
	return an.executeContract(ctx, contractAddress, keyName, execMethodName, execParams...)
}

// GetBlockByHeight implements chains.Chain
func (an *EVMRemotenet) GetBlockByHeight(ctx context.Context) (context.Context, error) {
	panic("unimplemented")
}

// GetLastBlock implements chains.Chain
func (an *EVMRemotenet) GetLastBlock(ctx context.Context) (context.Context, error) {
	h, err := an.Height(ctx)
	return context.WithValue(ctx, chains.LastBlock{}, h), err
}

func (an *EVMRemotenet) InitEventListener(ctx context.Context, contract string) chains.EventListener {
	//listener := NewIconEventListener(c, contract)
	return nil
}

// QueryContract implements chains.Chain
func (an *EVMRemotenet) QueryContract(ctx context.Context, contractAddress, methodName string, params map[string]interface{}) (context.Context, error) {
	time.Sleep(2 * time.Second)

	// get query msg
	query := an.GetQueryParam(methodName, params)
	_params, _ := json.Marshal(query.Value)
	output, err := an.getFullNode().QueryContract(ctx, contractAddress, query.MethodName, string(_params))

	chains.Response = output
	fmt.Printf("Response is : %s \n", output)
	return context.WithValue(ctx, "query-result", chains.Response), err

}

func (an *EVMRemotenet) BuildWallets(ctx context.Context, keyName string) (ibc.Wallet, error) {
	fmt.Println("I am building wallets", "Foundry")
	w := an.CheckForKeyStore(ctx, keyName)
	if w == nil {
		return nil, fmt.Errorf("error keyName already exists")
	}

	amount := ibc.WalletAmount{
		Address: w.FormattedAddress(),
		Amount:  10000,
	}
	var err error

	err = an.SendFunds(ctx, interchaintest.FaucetAccountKeyName, amount)
	return w, err
}

// PauseNode pauses the node
func (an *EVMRemotenet) PauseNode(ctx context.Context) error {
	return nil
}

// UnpauseNode starts the paused node
func (an *EVMRemotenet) UnpauseNode(ctx context.Context) error {
	return nil
}

func (an *EVMRemotenet) DeployContractRemote(ctx context.Context, contractPath, key string, params ...interface{}) (common.Address, error) {
	bytecode, contractABI, err := an.loadABI(contractPath)
	if err != nil {
		return common.Address{}, err
	}
	privateKey, fromAddress, err := getPrivateKey(an.testconfig.KeystorePassword)

	if err != nil {
		return common.Address{}, err
	}

	nonce, gasPrice, err := an.getNonceAndGasPrice(fromAddress)
	if err != nil {
		return common.Address{}, err
	}
	chainID, _ := an.Client.ChainID(ctx)
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(30000000)
	auth.GasPrice = gasPrice

	address, tx, _, err := bind.DeployContract(auth, contractABI, bytecode, &an.Client, params...)
	if err != nil {
		fmt.Println(auth, fromAddress)
		fmt.Printf("error while deploying contract :: %w", err)
		return common.Address{}, err
	}
	fmt.Println("Contract deployment transaction hash:", tx.Hash().Hex())
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second) // Set your desired timeout
	defer cancel()
	minedTx, err := bind.WaitMined(timeoutCtx, &an.Client, tx)
	if err != nil {
		return common.Address{}, err
	}
	fmt.Println("Contract deployment transaction status:", minedTx.Status)
	an.ContractABI[address.Hex()] = contractABI
	return address, err
}

func (an *EVMRemotenet) loadABI(contractPath string) ([]byte, abi.ABI, error) {

	_abi, err := os.Open(contractPath + ".abi.json")
	if err != nil {
		return nil, abi.ABI{}, err
	}
	defer _abi.Close()

	_bin, err := os.ReadFile(contractPath + ".bin")
	if err != nil {
		return nil, abi.ABI{}, err
	}

	bytecode, err := hex.DecodeString(string(_bin))
	if err != nil {
		return nil, abi.ABI{}, err
	}

	contractABI, err := abi.JSON(_abi)

	return bytecode, contractABI, err
}

func (an *EVMRemotenet) getNonceAndGasPrice(fromAddress common.Address) (uint64, *big.Int, error) {
	nonce, err := an.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return 0, nil, err
	}

	gasPrice, err := an.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return 0, nil, err
	}

	return nonce, gasPrice, nil
}

func (an *EVMRemotenet) ExecCallTx(ctx context.Context, contractAddress, methodName, pKey string, params ...interface{}) (*eth_types.Receipt, error) {
	address := common.HexToAddress(contractAddress)
	parsedABI := an.ContractABI[contractAddress]
	privateKey, fromAddress, _ := getPrivateKey(pKey)

	nonce, gasPrice, _ := an.getNonceAndGasPrice(fromAddress)

	data, err := parsedABI.Pack(methodName, params...)
	if err != nil {
		an.log.Error("Failed to pack abi", zap.Error(err), zap.String("method", methodName), zap.Any("params", params))
		return nil, err
	}
	txdata := &eth_types.LegacyTx{
		To:       &address,
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      uint64(2000000),
		Value:    big.NewInt(0x0),
		Data:     data,
	}
	tx := eth_types.NewTx(txdata)

	networkID, _ := an.Client.NetworkID(context.Background())
	signedTx, err := eth_types.SignTx(tx, eth_types.NewEIP155Signer(networkID), privateKey)
	if err != nil {
		return nil, err
	}
	err = an.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	receipt, err := bind.WaitMined(context.Background(), &an.Client, signedTx)
	if err != nil {
		return nil, err
	}
	if receipt.Status == 0 {
		cmd := []string{
			"run",
			receipt.TxHash.String(),
			"--rpc-url",
			an.GetRPCAddress(),
		}
		out, _, _ := an.ExecBin(ctx, cmd...)
		return receipt, fmt.Errorf("error on trasaction :: %s\n%v", methodName, string(out))
	}
	return receipt, nil

}

func (an *EVMRemotenet) ExecCallTxCommand(contractAddress, methodName, pKey string, params ...interface{}) []string {
	// get password from pathname as pathname will have the password prefixed. ex - Alice.Json
	address := common.HexToAddress(contractAddress)
	parsedABI := an.ContractABI[contractAddress]

	//if err != nil {
	//
	//}
	//msg := ethereum.CallMsg{
	//	To:   &address,
	//	Data: data,
	//}
	//result, err := an.Client.CallContract(context.Background(), msg, nil)

	privateKey, fromAddress, _ := getPrivateKey(pKey)

	nonce, gasPrice, _ := an.getNonceAndGasPrice(fromAddress)

	data, _ := parsedABI.Pack(methodName, params...)
	txdata := &eth_types.LegacyTx{
		To:       &address,
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      uint64(200000),
		Value:    big.NewInt(0x0),
		Data:     data,
	}
	tx := eth_types.NewTx(txdata)

	networkID, _ := an.Client.NetworkID(context.Background())
	signedTx, err := eth_types.SignTx(tx, eth_types.NewEIP155Signer(networkID), privateKey)
	if err != nil {
		return nil
	}
	err = an.Client.SendTransaction(context.Background(), signedTx)
	return nil

}

func (an *EVMRemotenet) BinCommand(command ...string) []string {
	command = append([]string{an.Config().Bin}, command...)
	return command
}

func (an *EVMRemotenet) ExecBin(ctx context.Context, command ...string) ([]byte, []byte, error) {
	return an.Exec(ctx, an.BinCommand(command...), nil)
}
