package evm

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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
	"io"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/icon-project/centralized-relay/test/chains"
	icontypes "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon/types"

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
	CallMessageSent = Event{"CallMessageSent", "0x69e53ea70fdf945f6d035b3979748bc999151691fb1dc69d66f8017f8840ae28", "xcall"}
	CallMessage     = Event{"CallMessage", "0x6ff9bfdb841175019a45e50c5c186a051194be6f51d46e2c901839550c9d413d", "xcall"}
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

func (c *EVMLocalnet) CheckForTimeout(ctx context.Context, src chains.Chain, params map[string]interface{}, listener chains.EventListener) (context.Context, error) {
	//TODO implement me
	panic("implement me")
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

	// Specify the path to the Dockerfile and the context directory
	//dockerfilePath := "Dockerfile"
	//contextDir := fmt.Sprintf("%s/test/chains/evm/hardhat-docker", os.Getenv(chains.BASE_PATH))

	//_internal = c.buildDockerImage(ctx, cli, contextDir, dockerfilePath, image.Repository+":"+image.Version)

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
	connection, err := c.getFullNode().DeployContract(ctx, c.scorePaths["connection"], c.privateKey)
	if err != nil {
		return err
	}
	relayerAddress := common.HexToAddress(c.Wallets[fmt.Sprintf("relayer-%s", c.Config().Name)].FormattedAddress())

	_, err = c.getFullNode().ExecCallTx(ctx, connection.Hex(), "initialize", c.privateKey, relayerAddress, xcall)
	if err != nil {
		fmt.Printf("fail to initialized xcall-adapter : %w\n", err)
		return err
	}
	//ctx, err = c.executeContract(context.Background(), ibcAddress, interchaintest.IBCOwnerAccount, "bindPort", `{"portId":"`+portId+`", "moduleAddress":"`+connection+`"}`)
	//c.IBCAddresses[fmt.Sprintf("xcall-%s", testcase)] = xcall
	//c.IBCAddresses[fmt.Sprintf("connection-%s", testcase)] = connection
	c.IBCAddresses["xcall"] = xcall.Hex()
	c.IBCAddresses["connection"] = connection.Hex()
	return err
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

// HasPacketReceipt returns the receipt of the packet sent to the target chain
func (c *EVMLocalnet) IsPacketReceived(ctx context.Context, params map[string]interface{}, order ibc.Order) bool {
	if order == ibc.Ordered {
		sequence := params["sequence"].(uint64) //2
		ctx, err := c.QueryContract(ctx, c.IBCAddresses["ibc"], chains.GetNextSequenceReceive, params)
		if err != nil {
			fmt.Printf("Error--%v\n", err)
			return false
		}
		response, err := formatHexNumberFromResponse(ctx.Value("query-result").([]byte))

		if err != nil {
			fmt.Printf("Error--%v\n", err)
			return false
		}
		fmt.Printf("response[\"data\"]----%v", response)
		return sequence < response
	}
	ctx, _ = c.QueryContract(ctx, c.IBCAddresses["ibc"], chains.HasPacketReceipt, params)

	response, err := formatHexNumberFromResponse(ctx.Value("query-result").([]byte))
	if err != nil {
		fmt.Printf("Error--%v\n", err)
		return false
	}
	return response == 1
}

func formatHexNumberFromResponse(value []byte) (uint64, error) {
	pattern := `0x[0-9a-fA-F]+`
	regex := regexp.MustCompile(pattern)
	result := regex.FindString(string(value))
	if result == "" {
		return 0, fmt.Errorf("number not found")

	}

	response, err := strconv.ParseInt(result, 0, 64)
	if err != nil {
		return 0, err
	}
	return uint64(response), nil
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
	if rollback == nil {
		rollback = make([]byte, 0)
	}
	ctx, err = c.SendPacketXCall(ctx, keyName, to, data, rollback)
	if err != nil {
		return nil, err
	}
	return c.FindTargetXCallMessage(ctx, targetChain, height, strings.Split(to, "/")[1])
}

func parseEventLogs(contractABI abi.ABI, logs []*eth_types.Log) ([]LogEvent, error) {
	//0x69e53ea70fdf945f6d035b3979748bc999151691fb1dc69d66f8017f8840ae28
	//0x6dbbb5c83189670e066d281dfc37d9ded5132af5d6401cfc831c7499eb775f3d
	var logEvents []LogEvent
	for _, _log := range logs {
		for _, topic := range _log.Topics {
			event, err := contractABI.EventByID(topic)
			if err != nil {
				continue
			}

			//var decodedData = make(map[string]interface{})
			decodedData, err := event.Inputs.UnpackValues(_log.Data)
			if err != nil {
				return nil, err
			}

			logEvent := LogEvent{
				Event:   event.Name,
				Message: decodedData[0].(string),
			}

			logEvents = append(logEvents, logEvent)

		}
	}

	return logEvents, nil
}

type LogEvent struct {
	Event   string
	Message string
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

func (c *EVMLocalnet) ExecuteCall(ctx context.Context, reqId, data string) (context.Context, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)

	return c.executeContract(ctx, c.IBCAddresses["xcall"], interchaintest.UserAccount, "executeCall", reqId, data)
}

func (c *EVMLocalnet) ExecuteRollback(ctx context.Context, sn string) (context.Context, error) {
	testcase := ctx.Value("testcase").(string)
	xCallKey := fmt.Sprintf("xcall-%s", testcase)
	ctx, err := c.executeContract(ctx, c.IBCAddresses[xCallKey], interchaintest.UserAccount, "executeRollback", sn)
	if err != nil {
		return nil, err
	}
	txn := ctx.Value("txResult").(*icontypes.TransactionResult)
	sequence, err := icontypes.HexInt(txn.EventLogs[0].Indexed[1]).Int()
	return context.WithValue(ctx, "IsRollbackEventFound", fmt.Sprintf("%d", sequence) == sn), nil

}

func (c *EVMLocalnet) FindCallMessage(ctx context.Context, startHeight uint64, from, to, sn string) (string, string, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	index := []common.Hash{common.BytesToHash([]byte(from)), common.BytesToHash([]byte(to)), common.BytesToHash([]byte(sn))}
	event, err := c.FindEvent(ctx, startHeight, CallMessage, index)
	if err != nil {
		return "", "", err
	}
	return event["_reqId"].(string), string(event["_data"].(byte)), nil
	//intHeight, _internal := event.Height.Int()
	//block, _internal := c.getFullNode().Client.GetBlockByHeight(&icontypes.BlockHeightParam{Height: icontypes.NewHexInt(int64(intHeight - 1))})
	//i, _internal := event.Index.Int()
	//tx := block.NormalTransactions[i]
	//trResult, _internal := c.getFullNode().TransactionResult(ctx, string(tx.TxHash))
	//eventIndex, _internal := event.Events[0].Int()
	//reqId := trResult.EventLogs[eventIndex].Data[0]
	//data := trResult.EventLogs[eventIndex].Data[1]
	//return reqId, data, nil
	//return "", "", nil
}

func (c *EVMLocalnet) FindCallResponse(ctx context.Context, startHeight uint64, sn string) (string, error) {
	//testcase := ctx.Value("testcase").(string)
	//xCallKey := fmt.Sprintf("xcall-%s", testcase)
	//index := []*string{&sn}
	//event, err := c.FindEvent(ctx, startHeight, xCallKey, "ResponseMessage(int,int)", index)
	//if err != nil {
	//	return "", err
	//}
	//
	//intHeight, _internal := event.Height.Int()
	//block, _internal := c.getFullNode().Client.GetBlockByHeight(&icontypes.BlockHeightParam{Height: icontypes.NewHexInt(int64(intHeight - 1))})
	//i, _internal := event.Index.Int()
	//tx := block.NormalTransactions[i]
	//trResult, _internal := c.getFullNode().TransactionResult(ctx, string(tx.TxHash))
	//eventIndex, _internal := event.Events[0].Int()
	//code, _internal := strconv.ParseInt(trResult.EventLogs[eventIndex].Data[0], 0, 64)
	//
	//return strconv.FormatInt(code, 10), nil
	return "", nil
}

func (c *EVMLocalnet) FindEvent(ctx context.Context, startHeight uint64, event Event, index []common.Hash) (map[string]interface{}, error) {
	//eventSignature := []byte(event.hash)
	eClient, err := ethclient.Dial(fmt.Sprintf("ws://%s", c.getFullNode().HostRPCPort))
	defer eClient.Close()
	if err != nil {
		return nil, errors.New("error: fail to create eth client")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Starting block number
	fromBlock := new(big.Int).SetUint64(startHeight)
	address := common.HexToAddress(c.IBCAddresses[event.contract])

	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		Addresses: []common.Address{address},
		Topics: [][]common.Hash{
			{common.HexToHash(event.hash)},
			index,
		},
	}

	logs := make(chan eth_types.Log)

	sub, err := eClient.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		return nil, err
	}
	defer sub.Unsubscribe()

	ev := make(map[string]interface{})
	select {
	case <-ctx.Done():
		return nil, errors.New("timeout: event not found within specified duration")
	case err := <-sub.Err():
		return nil, err
	case lg := <-logs:
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

}

func (c *EVMLocalnet) FindEvent__(ctx context.Context, startHeight uint64, event Event, index []common.Hash) (map[string]interface{}, error) {
	//eventSignature := []byte(event.hash)

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Second)
	defer cancel()

	// Starting block number
	fromBlock := new(big.Int).SetUint64(startHeight)
	address := common.HexToAddress(c.IBCAddresses[event.contract])

	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		Addresses: []common.Address{address},
		Topics: [][]common.Hash{
			{common.HexToHash(event.hash)},
			index,
		},
	}

	//logs := make(chan eth_types.Log)

	logs, err := c.getFullNode().Client.FilterLogs(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(logs) == 0 {
		return nil, fmt.Errorf("No event found")
	}

	ev := make(map[string]interface{})
	for _, lg := range logs {
		contractABI := c.getFullNode().ContractABI[address.Hex()]
		if len(lg.Data) > 0 {
			if err := contractABI.UnpackIntoMap(ev, event.name, lg.Data); err != nil {
				return ev, err
			}
		}

		var indexed abi.Arguments
		for _, arg := range contractABI.Events[event.name].Inputs {
			if arg.Indexed {
				indexed = append(indexed, arg)
			}
		}

		err := abi.ParseTopicsIntoMap(ev, indexed, lg.Topics[1:])
		if err != nil {
			return ev, err
		}
	}

	return ev, nil
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
	time.Sleep(2 * time.Second)
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

func (c *EVMLocalnet) SendPacketMockDApp(ctx context.Context, targetChain chains.Chain, keyName string, params map[string]interface{}) (chains.PacketTransferResponse, error) {
	//listener := targetChain.InitEventListener(ctx, "ibc")
	//response := chains.PacketTransferResponse{}
	//testcase := ctx.Value("testcase").(string)
	//dappKey := fmt.Sprintf("mockdapp-%s", testcase)
	//execMethodName, execParams := c.getExecuteParam(ctx, chains.SendMessage, params)
	//ctx, err := c.executeContract(ctx, c.IBCAddresses[dappKey], keyName, execMethodName, execParams)
	//if err != nil {
	//	return response, err
	//}
	////txn := ctx.Value("txResult").(*icontypes.TransactionResult)
	//response.IsPacketSent = true

	//var packet = chantypes.Packet{}
	//var protoPacket = icontypes.HexBytes(txn.EventLogs[1].Indexed[1])
	//_ = chains.HexBytesToProtoUnmarshal(protoPacket, &packet)
	//response.Packet = packet
	//filter := map[string]interface{}{
	//	"wasm-recv_packet.packet_sequence":    fmt.Sprintf("%d", packet.Sequence),
	//	"wasm-recv_packet.packet_src_port":    packet.SourcePort,
	//	"wasm-recv_packet.packet_src_channel": packet.SourceChannel,
	//}
	//event, err := listener.FindEvent(filter)
	//response.IsPacketReceiptEventFound = event != nil
	//return response, err
	return chains.PacketTransferResponse{IsPacketSent: true}, nil
}
