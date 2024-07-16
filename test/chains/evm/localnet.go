package evm

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	interchaintest "github.com/icon-project/centralized-relay/test"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/dockerutil"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer/centralized"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"gopkg.in/yaml.v3"

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
	cfg           chains.ChainConfig
	numValidators int
	numFullNodes  int
	scorePaths    map[string]string
	IBCAddresses  map[string]string     `json:"addresses"`
	Wallets       map[string]ibc.Wallet `json:"wallets"`
	DockerClient  *client.Client
	ContractABI   map[string]abi.ABI
	Network       string
	testconfig    *testconfig.Chain
	Client        ethclient.Client
}

// func NewEVMRemotenet(testName string, log *zap.Logger, chainConfig ibc.ChainConfig, numValidators int, numFullNodes int, scorePaths map[string]string) chains.Chain {
func NewEVMRemotenet(testName string, log *zap.Logger, chainConfig chains.ChainConfig, client *client.Client, network string, testconfig *testconfig.Chain) chains.Chain {
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
func (an *EVMRemotenet) Config() chains.ChainConfig {
	return an.cfg
}

// Exec runs an arbitrary command using Chain's docker environment.
// Whether the invoked command is run in a one-off container or execing into an already running container
// is up to the chain implementation.
//
// "env" are environment variables in the format "MY_ENV_VAR=value"
func (an *EVMRemotenet) Exec(ctx context.Context, cmd []string, env []string) (stdout []byte, stderr []byte, err error) {
	job := dockerutil.NewImage(an.log, an.DockerClient, an.Network, an.testName, an.cfg.Images.Repository, an.cfg.Images.Version)

	opts := dockerutil.ContainerOptions{
		Binds: []string{},
	}
	res := job.Run(ctx, cmd, opts)
	return res.Stdout, res.Stderr, res.Err
	// return an.getFullNode().Exec(ctx, cmd, env)
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
			WebsocketUrl:  an.testconfig.WebsocketUrl,
			StartHeight:   0,
			GasPrice:      10000000,
			GasLimit:      20000000,
			Contracts:     contracts,
			BlockInterval: "6s",
			Address:       an.testconfig.RelayWalletAddress,
			FinalityBlock: 10,
		},
	}
	return yaml.Marshal(config)
}

// Height returns the current block height or an error if unable to get current height.
func (an *EVMRemotenet) Height(ctx context.Context) (uint64, error) {
	header, err := an.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

func (an *EVMRemotenet) SetupConnection(ctx context.Context, target chains.Chain) error {
	xcall := common.HexToAddress(an.IBCAddresses["xcall"])
	connection, err := an.DeployContractRemote(ctx, an.scorePaths["connection"], an.testconfig.KeystorePassword)
	if err != nil {
		return err
	}
	relayerAddress := common.HexToAddress(an.testconfig.RelayWalletAddress)

	_, err = an.ExecCallTx(ctx, connection.Hex(), "initialize", an.testconfig.KeystorePassword, relayerAddress, xcall)
	if err != nil {
		fmt.Println("fail to initialized xcall-adapter : %w\n", err)
		return err
	}

	_, err = an.ExecCallTx(ctx, connection.Hex(), "setFee", an.testconfig.KeystorePassword, target.Config().ChainID, big.NewInt(0), big.NewInt(0))
	if err != nil {
		fmt.Println("fail to initialized fee for xcall-adapter : %w\n", err)
		return err
	}

	an.IBCAddresses["connection"] = connection.Hex()
	return nil
}

func (an *EVMRemotenet) SetupXCall(ctx context.Context) error {
	nid := an.cfg.ChainID
	xcall, err := an.DeployContractRemote(ctx, an.scorePaths["xcall"], an.testconfig.KeystorePassword)
	if err != nil {
		return err
	}
	_, err = an.ExecCallTx(ctx, xcall.Hex(), "initialize", an.testconfig.KeystorePassword, nid)
	if err != nil {
		fmt.Println("fail to initialized xcall : %w\n", err)
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
		fmt.Printf("fail to initialized dapp : %s\n", err)
		return err
	}

	an.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = dapp.Hex()

	for _, connection := range connections {
		//connectionKey := fmt.Sprintf("%s-%s", connection.Connection, testcase)
		params := []interface{}{connection.Nid, an.IBCAddresses[connection.Connection], connection.Destination}
		_, err = an.executeContract(context.Background(), dapp.Hex(), "addConnection", params...)
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
	ctx, err := an.executeContract(ctx, an.IBCAddresses[dappKey], "sendMessage", params...)
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
	height, err := targetChain.Height(ctx)
	if err != nil {
		return nil, err
	}

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
	if err != nil {
		return nil, err
	}
	sequence := events["_sn"]
	return context.WithValue(ctx, "IsRollbackEventFound", fmt.Sprintf("%d", sequence) == sn), nil

}

func (an *EVMRemotenet) FindCallMessage(ctx context.Context, startHeight uint64, from, to, sn string) (string, string, error) {
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
	if err != nil {
		return nil, errors.New("error: fail to create eth client")
	}
	defer eClient.Close()
	_, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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

// executeContract implements chains.Chain
func (an *EVMRemotenet) executeContract(ctx context.Context, contractAddress, methodName string, params ...interface{}) (context.Context, error) {
	//an.CheckForKeyStore(ctx, keyName)

	receipt, err := an.ExecCallTx(ctx, contractAddress, methodName, an.testconfig.KeystorePassword, params...)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Transaction Hash for %s: %s\n", methodName, receipt.TxHash.String())

	return context.WithValue(ctx, "txResult", receipt), nil

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
		fmt.Println("error while deploying contract :: %w", err)
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
	address := common.HexToAddress(contractAddress)
	parsedABI := an.ContractABI[contractAddress]

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
	if err != nil {
		return nil
	}
	return nil

}

func (an *EVMRemotenet) BinCommand(command ...string) []string {
	command = append([]string{an.Config().Bin}, command...)
	return command
}

func (an *EVMRemotenet) ExecBin(ctx context.Context, command ...string) ([]byte, []byte, error) {
	return an.Exec(ctx, an.BinCommand(command...), nil)
}

func getPrivateKey(key string) (*ecdsa.PrivateKey, common.Address, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, common.Address{}, err
	}
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	return privateKey, fromAddress, nil
}

func (an *EVMRemotenet) FindRollbackExecutedMessage(ctx context.Context, startHeight uint64, sn string) (string, error) {
	_sn, _ := big.NewInt(0).SetString(sn, 10)
	topics := []common.Hash{common.HexToHash(RollbackExecuted.hash), common.BytesToHash(_sn.Bytes())}
	_, err := an.FindEvent(ctx, startHeight, RollbackExecuted, topics)
	if err != nil {
		fmt.Printf("Topics %v", topics)
		return "", err
	}

	return "0", nil
}
