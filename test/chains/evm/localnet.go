package evm

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	interchaintest "github.com/icon-project/centralized-relay/test"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/blockdb"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/dockerutil"

	dockertypes "github.com/docker/docker/api/types"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/icon-project/centralized-relay/test/chains"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

type EVMLocalnet struct {
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
}

func (c *EVMLocalnet) CreateKey(ctx context.Context, keyName string) error {
	//TODO implement me
	panic("implement me")
}

func NewEVMLocalnet(testName string, log *zap.Logger, chainConfig ibc.ChainConfig, numValidators int, numFullNodes int, scorePaths map[string]string) chains.Chain {
	return &EVMLocalnet{
		testName:      testName,
		cfg:           chainConfig,
		numValidators: numValidators,
		numFullNodes:  numFullNodes,
		log:           log,
		scorePaths:    scorePaths,
		Wallets:       map[string]ibc.Wallet{},
		IBCAddresses:  make(map[string]string),
	}
}

// Config fetches the chain configuration.
func (c *EVMLocalnet) Config() ibc.ChainConfig {
	return c.cfg
}

// Initialize initializes node structs so that things like initializing keys can be done before starting the chain
func (c *EVMLocalnet) Initialize(ctx context.Context, testName string, cli *client.Client, networkID string) error {
	chainCfg := c.Config()
	c.pullImages(ctx, cli)
	image := chainCfg.Images[0]

	newFullNodes := make(HardhatNodes, c.numFullNodes)
	copy(newFullNodes, c.FullNodes)

	eg, egCtx := errgroup.WithContext(ctx)
	for i := len(c.FullNodes); i < c.numFullNodes; i++ {
		i := i
		eg.Go(func() error {
			fn, err := c.NewChainNode(egCtx, testName, cli, networkID, image, false)
			if err != nil {
				return err
			}
			fn.Index = i
			newFullNodes[i] = fn
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	c.findTxMu.Lock()
	defer c.findTxMu.Unlock()
	c.FullNodes = newFullNodes
	return nil
}

func (c *EVMLocalnet) pullImages(ctx context.Context, cli *client.Client) {
	for _, image := range c.Config().Images {
		rc, err := cli.ImagePull(
			ctx,
			image.Repository+":"+image.Version,
			dockertypes.ImagePullOptions{},
		)
		if err != nil {
			c.log.Error("Failed to pull image",
				zap.Error(err),
				zap.String("repository", image.Repository),
				zap.String("tag", image.Version),
			)
		} else {
			_, _ = io.Copy(io.Discard, rc)
			_ = rc.Close()
		}
	}
}

func (c *EVMLocalnet) NewChainNode(
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
		log:          c.log,
		Chain:        c,
		DockerClient: cli,
		NetworkID:    networkID,
		TestName:     testName,
		Image:        image,
		ContractABI:  make(map[string]abi.ABI),
	}

	v, err := cli.VolumeCreate(ctx, volumetypes.VolumeCreateBody{
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
		Log: c.log,

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
func (c *EVMLocalnet) Start(testName string, ctx context.Context, additionalGenesisWallets ...ibc.WalletAmount) error {
	c.findTxMu.Lock()
	defer c.findTxMu.Unlock()
	eg, egCtx := errgroup.WithContext(ctx)
	for _, n := range c.FullNodes {
		n := n
		eg.Go(func() error {

			if err := n.CreateNodeContainer(egCtx, additionalGenesisWallets...); err != nil {
				return err
			}
			// All (validators, gentx, fullnodes, peering, additional accounts) are included in the image itself.
			return n.StartContainer(ctx)
		})
	}
	return eg.Wait()
}

// Exec runs an arbitrary command using Chain's docker environment.
// Whether the invoked command is run in a one-off container or execing into an already running container
// is up to the chain implementation.
//
// "env" are environment variables in the format "MY_ENV_VAR=value"
func (c *EVMLocalnet) Exec(ctx context.Context, cmd []string, env []string) (stdout []byte, stderr []byte, err error) {
	return c.getFullNode().Exec(ctx, cmd, env)
}

// ExportState exports the chain state at specific height.
func (c *EVMLocalnet) ExportState(ctx context.Context, height int64) (string, error) {
	block, err := c.getFullNode().GetBlockByHeight(ctx, height)
	return block, err
}

// GetRPCAddress retrieves the rpc address that can be reached by other containers in the docker network.
func (c *EVMLocalnet) GetRPCAddress() string {
	return fmt.Sprintf("http://%s:8545", c.getFullNode().HostName())
}

// GetGRPCAddress retrieves the grpc address that can be reached by other containers in the docker network.
// Not Applicable for Icon
func (c *EVMLocalnet) GetGRPCAddress() string {
	return ""
}

// GetHostRPCAddress returns the rpc address that can be reached by processes on the host machine.
// Note that this will not return a valid value until after Start returns.
func (c *EVMLocalnet) GetHostRPCAddress() string {
	return "http://" + c.getFullNode().HostRPCPort
}

// GetHostGRPCAddress returns the grpc address that can be reached by processes on the host machine.
// Note that this will not return a valid value until after Start returns.
// Not applicable for Icon
func (c *EVMLocalnet) GetHostGRPCAddress() string {
	return ""
}

// HomeDir is the home directory of a node running in a docker container. Therefore, this maps to
// the container's filesystem (not the host).
func (c *EVMLocalnet) HomeDir() string {
	return c.FullNodes[0].HomeDir()
}

func (c *EVMLocalnet) createKeystore(ctx context.Context, keyName string) (string, string, error) {
	keydir := keystore.NewKeyStore("/tmp", keystore.StandardScryptN, keystore.StandardScryptP)

	account, err := keydir.NewAccount(keyName)

	if err != nil {
		c.log.Fatal("unable to create new account", zap.Error(err))
		return "", "", err
	}

	keyJSON, err := keydir.Export(account, keyName, keyName)
	err = c.FullNodes[0].RestoreKeystore(ctx, keyJSON, keyName)
	if err != nil {
		c.log.Error("fail to restore keystore", zap.Error(err))
		return "", "", err
	}
	key, err := keystore.DecryptKey(keyJSON, keyName)
	privateKey := crypto.FromECDSA(key.PrivateKey)
	return account.Address.Hex(), hex.EncodeToString(privateKey), nil
}

// RecoverKey recovers an existing user from a given mnemonic.
func (c *EVMLocalnet) RecoverKey(ctx context.Context, name string, mnemonic string) error {
	panic("not implemented") // TODO: Implement
}

// GetAddress fetches the bech32 address for a test key on the "user" node (either the first fullnode or the first validator if no fullnodes).
func (c *EVMLocalnet) GetAddress(ctx context.Context, keyName string) ([]byte, error) {
	addrInByte, err := json.Marshal(keyName)
	if err != nil {
		return nil, err
	}
	return addrInByte, nil
}

func (c *EVMLocalnet) GetRelayConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error) {
	return c.FullNodes[0].GetChainConfig(ctx, rlyHome, keyName)
}

// SendFunds sends funds to a wallet from a user account.
func (c *EVMLocalnet) SendFunds(ctx context.Context, keyName string, amount ibc.WalletAmount) error {
	c.CheckForKeyStore(ctx, keyName)

	privateKey, err := crypto.HexToECDSA(c.Wallets[interchaintest.FaucetAccountKeyName].Mnemonic())
	if err != nil {
		return err
	}
	ethClient := c.getFullNode().Client
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
	signedTx, err := eth_types.SignTx(tx, eth_types.NewEIP155Signer(networkID), privateKey)
	err = ethClient.SendTransaction(context.Background(), signedTx)

	return err

}

// Height returns the current block height or an error if unable to get current height.
func (c *EVMLocalnet) Height(ctx context.Context) (uint64, error) {
	return c.getFullNode().Height(ctx)
}

// GetGasFeesInNativeDenom gets the fees in native denom for an amount of spent gas.
func (c *EVMLocalnet) GetGasFeesInNativeDenom(gasPaid int64) int64 {
	gasPrice, _ := strconv.ParseFloat(strings.Replace(c.cfg.GasPrices, c.cfg.Denom, "", 1), 64)
	fees := float64(gasPaid) * gasPrice
	return int64(fees)
}

// BuildRelayerWallet will return a chain-specific wallet populated with the mnemonic so that the wallet can
// be restored in the relayer node using the mnemonic. After it is built, that address is included in
// genesis with some funds.
func (c *EVMLocalnet) BuildRelayerWallet(ctx context.Context, keyName string) (ibc.Wallet, error) {
	return c.BuildWallet(ctx, keyName, "")
}

func (c *EVMLocalnet) BuildWallet(ctx context.Context, keyName string, mnemonic string) (ibc.Wallet, error) {
	address, privateKey, err := c.createKeystore(ctx, keyName)
	if err != nil {
		return nil, fmt.Errorf("failed to create key with name %q on chain %s: %w", keyName, c.cfg.Name, err)

	}

	w := NewWallet(keyName, []byte(address), privateKey, c.cfg)
	c.Wallets[keyName] = w
	return w, nil
}

func (c *EVMLocalnet) getFullNode() *AnvilNode {
	c.findTxMu.Lock()
	defer c.findTxMu.Unlock()
	if len(c.FullNodes) > 0 {
		// use first full node
		return c.FullNodes[0]
	}
	return c.FullNodes[0]
}

func (c *EVMLocalnet) FindTxs(ctx context.Context, height uint64) ([]blockdb.Tx, error) {
	fn := c.getFullNode()
	return fn.FindTxs(ctx, height)
}

// GetBalance fetches the current balance for a specific account address and denom.
func (c *EVMLocalnet) GetBalance(ctx context.Context, address string, denom string) (int64, error) {
	return c.getFullNode().GetBalance(ctx, address)
}

func (c *EVMLocalnet) SetupConnection(ctx context.Context, keyName string, target chains.Chain) error {
	//testcase := ctx.Value("testcase").(string)
	xcall := common.HexToAddress(c.IBCAddresses["xcall"])
	_ = c.CheckForKeyStore(ctx, keyName)
	connection, err := c.getFullNode().DeployContract(ctx, c.scorePaths["connection"], c.privateKey)
	if err != nil {
		return err
	}
	relayerKey := fmt.Sprintf("relayer-%s", c.Config().Name)
	relayerAddress := common.HexToAddress(c.Wallets[relayerKey].FormattedAddress())

	_, err = c.getFullNode().ExecCallTx(ctx, connection.Hex(), "initialize", c.privateKey, relayerAddress, xcall)
	if err != nil {
		fmt.Printf("fail to initialized xcall-adapter : %w\n", err)
		return err
	}

	_ = c.CheckForKeyStore(ctx, relayerKey)

	_, err = c.getFullNode().ExecCallTx(ctx, connection.Hex(), "setFee", c.privateKey, target.Config().ChainID, big.NewInt(0), big.NewInt(0))
	if err != nil {
		fmt.Printf("fail to initialized fee for xcall-adapter : %w\n", err)
		return err
	}

	c.IBCAddresses["connection"] = connection.Hex()
	return nil
}

func (c *EVMLocalnet) SetupXCall(ctx context.Context, keyName string) error {
	//testcase := ctx.Value("testcase").(string)
	nid := c.cfg.ChainID
	//ibcAddress := c.IBCAddresses["ibc"]
	_ = c.CheckForKeyStore(ctx, keyName)
	xcall, err := c.getFullNode().DeployContract(ctx, c.scorePaths["xcall"], c.privateKey)
	if err != nil {
		return err
	}

	_, err = c.getFullNode().ExecCallTx(ctx, xcall.Hex(), "initialize", c.privateKey, nid)
	if err != nil {
		fmt.Printf("fail to initialized xcall : %w\n", err)
		return err
	}

	c.IBCAddresses["xcall"] = xcall.Hex()
	return nil
}

func (c *EVMLocalnet) DeployXCallMockApp(ctx context.Context, keyName string, connections []chains.XCallConnection) error {
	testcase := ctx.Value("testcase").(string)
	c.CheckForKeyStore(ctx, keyName)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	xCall := c.GetContractAddress("xcall")
	dapp, err := c.getFullNode().DeployContract(ctx, c.scorePaths["dapp"], c.privateKey)
	if err != nil {
		return err
	}
	//_, err = c.executeContract(ctx,dapp,keyName,"initialize",xCall)
	_, err = c.getFullNode().ExecCallTx(ctx, dapp.Hex(), "initialize", c.privateKey, common.HexToAddress(xCall))
	if err != nil {
		fmt.Printf("fail to initialized dapp : %w\n", err)
		return err
	}

	c.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = dapp.Hex()

	for _, connection := range connections {
		//connectionKey := fmt.Sprintf("%s-%s", connection.Connection, testcase)
		params := []interface{}{connection.Nid, c.IBCAddresses[connection.Connection], connection.Destination}
		ctx, err = c.executeContract(context.Background(), dapp.Hex(), keyName, "addConnection", params...)
		if err != nil {
			c.log.Error("Unable to add connection",
				zap.Error(err),
				zap.String("nid", connection.Nid),
				zap.String("source", c.IBCAddresses[connection.Connection]),
				zap.String("destination", connection.Destination),
			)
		}
	}

	return nil
}

func (c *EVMLocalnet) GetContractAddress(key string) string {
	value, exist := c.IBCAddresses[key]
	if !exist {
		panic(fmt.Sprintf(`IBC address not exist %s`, key))
	}
	return value
}

func (c *EVMLocalnet) BackupConfig() ([]byte, error) {
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

func (c *EVMLocalnet) RestoreConfig(backup []byte) error {
	result := make(map[string]interface{})
	err := json.Unmarshal(backup, &result)
	if err != nil {
		return err
	}
	c.IBCAddresses = result["addresses"].(map[string]string)
	wallets := make(map[string]ibc.Wallet)

	for key, value := range result["wallets"].(map[string]interface{}) {
		_value := value.(map[string]string)
		mnemonic := _value["mnemonic"]
		address, _ := hex.DecodeString(_value["address"])
		wallets[key] = NewWallet(key, address, mnemonic, c.Config())
	}
	c.Wallets = wallets
	return nil
}

func (c *EVMLocalnet) SendPacketXCall(ctx context.Context, keyName, _to string, data, rollback []byte) (context.Context, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	if rollback == nil {
		rollback = make([]byte, 0)
	}
	params := []interface{}{_to, data, rollback}
	// TODO: send fees
	ctx, err := c.executeContract(ctx, c.IBCAddresses[dappKey], keyName, "sendMessage", params...)
	if err != nil {
		return nil, err
	}
	txn := ctx.Value("txResult").(*eth_types.Receipt)

	events, err := c.ParseEvent(CallMessageSent, txn.Logs)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, "sn", events["_sn"].(*big.Int).String()), nil
}

func (c *EVMLocalnet) SendNewPacketXCall(ctx context.Context, keyName, _to string, data []byte, messageType *big.Int, rollback []byte) (context.Context, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	if rollback == nil {
		rollback = make([]byte, 0)
	}
	params := []interface{}{_to, data, messageType, rollback}
	// TODO: send fees
	ctx, err := c.executeContract(ctx, c.IBCAddresses[dappKey], keyName, "sendNewMessage", params...)
	if err != nil {
		return nil, err
	}
	txn := ctx.Value("txResult").(*eth_types.Receipt)
	events, err := c.ParseEvent(CallMessageSent, txn.Logs)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, "sn", events["_sn"].(*big.Int).String()), nil
}

// FindTargetXCallMessage returns the request id and the data of the message sent to the target chain
func (c *EVMLocalnet) FindTargetXCallMessage(ctx context.Context, target chains.Chain, height uint64, to string) (*chains.XCallResponse, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	sn := ctx.Value("sn").(string)
	reqId, destData, err := target.FindCallMessage(ctx, height, c.cfg.ChainID+"/"+c.IBCAddresses[dappKey], to, sn)
	return &chains.XCallResponse{SerialNo: sn, RequestID: reqId, Data: destData}, err
}

func (c *EVMLocalnet) XCall(ctx context.Context, targetChain chains.Chain, keyName, to string, data, rollback []byte) (*chains.XCallResponse, error) {
	height, err := targetChain.(ibc.Chain).Height(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: send fees

	ctx, err = c.SendPacketXCall(ctx, keyName, to, data, rollback)
	if err != nil {
		return nil, err
	}
	return c.FindTargetXCallMessage(ctx, targetChain, height, strings.Split(to, "/")[1])
}

func (c *EVMLocalnet) NewXCall(ctx context.Context, targetChain chains.Chain, keyName, to string, data []byte, messageType *big.Int, rollback []byte) (*chains.XCallResponse, error) {
	height, err := targetChain.(ibc.Chain).Height(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: send fees

	ctx, err = c.SendNewPacketXCall(ctx, keyName, to, data, messageType, rollback)
	if err != nil {
		return nil, err
	}
	return c.FindTargetXCallMessage(ctx, targetChain, height, strings.Split(to, "/")[1])
}

func (c *EVMLocalnet) ExecuteCall(ctx context.Context, reqId, data string) (context.Context, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	_reqId, _ := big.NewInt(0).SetString(reqId, 10)
	return c.executeContract(ctx, c.IBCAddresses["xcall"], interchaintest.UserAccount, "executeCall", _reqId, []byte(data))
}

func (c *EVMLocalnet) ExecuteRollback(ctx context.Context, sn string) (context.Context, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	_sn, _ := big.NewInt(0).SetString(sn, 10)
	ctx, err := c.executeContract(ctx, c.IBCAddresses["xcall"], interchaintest.UserAccount, "executeRollback", _sn)
	if err != nil {
		return nil, err
	}
	txn := ctx.Value("txResult").(*eth_types.Receipt)
	events, err := c.ParseEvent(RollbackExecuted, txn.Logs)
	sequence := events["_sn"]
	return context.WithValue(ctx, "IsRollbackEventFound", fmt.Sprintf("%d", sequence) == sn), nil

}

func (c *EVMLocalnet) FindCallMessage(ctx context.Context, startHeight uint64, from, to, sn string) (string, string, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	fmt.Printf("%s,--%s,--%s\n", from, to, sn)
	_sn, _ := big.NewInt(0).SetString(sn, 10)
	topics := []common.Hash{common.HexToHash(CallMessage.hash), crypto.Keccak256Hash([]byte(from)), crypto.Keccak256Hash([]byte(to)), common.BytesToHash(_sn.Bytes())}
	event, err := c.FindEvent(ctx, startHeight, CallMessage, topics)

	if err != nil {
		fmt.Printf("Topics %v\n", topics)
		return "", "", err
	}
	return event["_reqId"].(*big.Int).String(), string(event["_data"].([]byte)), nil
}

func (c *EVMLocalnet) FindCallResponse(ctx context.Context, startHeight uint64, sn string) (string, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	_sn, _ := big.NewInt(0).SetString(sn, 10)
	topics := []common.Hash{common.HexToHash(ResponseMessage.hash), common.BytesToHash(_sn.Bytes())}

	event, err := c.FindEvent(ctx, startHeight, ResponseMessage, topics)
	if err != nil {
		fmt.Printf("Topics %v", topics)
		return "", err
	}

	return event["_code"].(*big.Int).String(), nil

}

func (c *EVMLocalnet) FindEvent(ctx context.Context, startHeight uint64, event Event, topics []common.Hash) (map[string]interface{}, error) {
	//eventSignature := []byte(event.hash)
	eClient, err := ethclient.Dial(fmt.Sprintf("ws://%s", c.getFullNode().HostRPCPort))
	defer eClient.Close()
	if err != nil {
		return nil, errors.New("error: fail to create eth client")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Starting block number
	fromBlock := new(big.Int).SetUint64(startHeight - 1)
	address := common.HexToAddress(c.IBCAddresses[event.contract])
	fmt.Printf("address :: %v\n", address)
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
			contractABI := c.getFullNode().ContractABI[address.Hex()]
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

func (c *EVMLocalnet) ParseEvent(event Event, logs []*eth_types.Log) (map[string]interface{}, error) {
	contractABI := c.getFullNode().ContractABI[c.IBCAddresses[event.contract]]
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
func (c *EVMLocalnet) DeployContract(ctx context.Context, keyName string) (context.Context, error) {
	// Get contract Name from context
	ctxValue := ctx.Value(chains.ContractName{}).(chains.ContractName)
	contractName := ctxValue.ContractName

	// Get Init Message from context
	ctxVal := ctx.Value(chains.InitMessageKey("init-msg")).(chains.InitMessage)

	initMessage := c.getInitParams(ctx, contractName, ctxVal.Message)

	var contracts chains.ContractKey

	// Check if 6123f953784d27e0729bc7a640d6ad8f04ed6710.keystore is alreadry available for given keyName
	ownerAddr := c.CheckForKeyStore(ctx, keyName)
	if ownerAddr != nil {
		contracts.ContractOwner = map[string]string{
			keyName: ownerAddr.FormattedAddress(),
		}
	}

	// Get ScoreAddress
	scoreAddress, err := c.getFullNode().DeployContract(ctx, c.scorePaths[contractName], c.privateKey, initMessage)

	contracts.ContractAddress = map[string]string{
		contractName: scoreAddress.Hex(),
	}

	testcase := ctx.Value("testcase").(string)
	contract := fmt.Sprintf("%s-%s", contractName, testcase)
	c.IBCAddresses[contract] = scoreAddress.Hex()
	return context.WithValue(ctx, chains.Mykey("contract Names"), chains.ContractKey{
		ContractAddress: contracts.ContractAddress,
		ContractOwner:   contracts.ContractOwner,
	}), err
}

// executeContract implements chains.Chain
func (c *EVMLocalnet) executeContract(ctx context.Context, contractAddress, keyName, methodName string, params ...interface{}) (context.Context, error) {
	c.CheckForKeyStore(ctx, keyName)

	receipt, err := c.getFullNode().ExecCallTx(ctx, contractAddress, methodName, c.privateKey, params...)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Transaction Hash for %s: %s\n", methodName, receipt.TxHash.String())

	return context.WithValue(ctx, "txResult", receipt), nil

}

func (c *EVMLocalnet) ExecuteContract(ctx context.Context, contractAddress, keyName, methodName string, params map[string]interface{}) (context.Context, error) {
	execMethodName, execParams := c.getExecuteParam(ctx, methodName, params)
	return c.executeContract(ctx, contractAddress, keyName, execMethodName, execParams...)
}

// GetBlockByHeight implements chains.Chain
func (c *EVMLocalnet) GetBlockByHeight(ctx context.Context) (context.Context, error) {
	panic("unimplemented")
}

// GetLastBlock implements chains.Chain
func (c *EVMLocalnet) GetLastBlock(ctx context.Context) (context.Context, error) {
	h, err := c.getFullNode().Height(ctx)
	return context.WithValue(ctx, chains.LastBlock{}, h), err
}

func (c *EVMLocalnet) InitEventListener(ctx context.Context, contract string) chains.EventListener {
	//listener := NewIconEventListener(c, contract)
	return nil
}

// QueryContract implements chains.Chain
func (c *EVMLocalnet) QueryContract(ctx context.Context, contractAddress, methodName string, params map[string]interface{}) (context.Context, error) {
	time.Sleep(2 * time.Second)

	// get query msg
	query := c.GetQueryParam(methodName, params)
	_params, _ := json.Marshal(query.Value)
	output, err := c.getFullNode().QueryContract(ctx, contractAddress, query.MethodName, string(_params))

	chains.Response = output
	fmt.Printf("Response is : %s \n", output)
	return context.WithValue(ctx, "query-result", chains.Response), err

}

func (c *EVMLocalnet) BuildWallets(ctx context.Context, keyName string) (ibc.Wallet, error) {
	w := c.CheckForKeyStore(ctx, keyName)
	if w == nil {
		return nil, fmt.Errorf("error keyName already exists")
	}

	amount := ibc.WalletAmount{
		Address: w.FormattedAddress(),
		Amount:  10000,
	}
	var err error

	err = c.SendFunds(ctx, interchaintest.FaucetAccountKeyName, amount)
	return w, err
}

// PauseNode pauses the node
func (c *EVMLocalnet) PauseNode(ctx context.Context) error {
	return c.getFullNode().DockerClient.ContainerPause(ctx, c.getFullNode().ContainerID)
}

// UnpauseNode starts the paused node
func (c *EVMLocalnet) UnpauseNode(ctx context.Context) error {
	return c.getFullNode().DockerClient.ContainerUnpause(ctx, c.getFullNode().ContainerID)
}
