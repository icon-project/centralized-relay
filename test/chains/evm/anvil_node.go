package evm

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerclient "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/icon-project/centralized-relay/test/chains"
	"github.com/icon-project/centralized-relay/test/interchaintest"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/blockdb"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/dockerutil"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer/centralized"
	icontypes "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon/types"
	"github.com/icza/dyno"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	rpcPort              = "8545/tcp"
	GOLOOP_IMAGE_ENV     = "GOLOOP_IMAGE"
	GOLOOP_IMAGE         = "iconloop/goloop-icon"
	GOLOOP_IMAGE_TAG_ENV = "GOLOOP_IMAGE_TAG"
	GOLOOP_IMAGE_TAG     = "latest"
)

type AnvilNode struct {
	VolumeName   string
	Index        int
	Chain        chains.Chain
	NetworkID    string
	DockerClient *dockerclient.Client
	Client       ethclient.Client
	TestName     string
	Image        ibc.DockerImage
	log          *zap.Logger
	ContainerID  string
	// Ports set during StartContainer.
	HostRPCPort string
	Validator   bool
	lock        sync.Mutex
	Address     string
	ContractABI map[string]abi.ABI
}

type HardhatNodes []*AnvilNode

// Name of the test node container
func (an *AnvilNode) Name() string {
	var nodeType string
	if an.Validator {
		nodeType = "val"
	} else {
		nodeType = "fn"
	}
	return fmt.Sprintf("%s-%s-%d-%s", an.Chain.Config().ChainID, nodeType, an.Index, dockerutil.SanitizeContainerName(an.TestName))
}

// Create Node Container with ports exposed and published for host to communicate with
func (an *AnvilNode) CreateNodeContainer(ctx context.Context, additionalGenesisWallets ...ibc.WalletAmount) error {
	imageRef := an.Image.Ref()
	//testBasePath := os.Getenv(chains.BASE_PATH)
	home := an.HomeDir()
	containerConfig := &types.ContainerCreateConfig{
		Config: &container.Config{
			Image:    imageRef,
			Hostname: an.HostName(),
			Env: []string{
				"BLOCK_TIME=2",
				fmt.Sprintf("GENESIS_PATH=%s/genesis.json", home),
			},
			Entrypoint: []string{"/bin/sh", fmt.Sprintf("%s/entrypoint.sh", home)},
			Labels:     map[string]string{dockerutil.CleanupLabel: an.TestName},
			ExposedPorts: nat.PortSet{
				"8545/tcp": struct{}{},
			},
		},

		HostConfig: &container.HostConfig{
			Binds:           an.Bind(),
			PublishAllPorts: true,
			AutoRemove:      false,
			DNS:             []string{},
		},
		NetworkingConfig: &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				an.NetworkID: {},
			},
		},
	}

	cc, err := an.DockerClient.ContainerCreate(ctx, containerConfig.Config, containerConfig.HostConfig, containerConfig.NetworkingConfig, an.Name())
	if err != nil {
		an.log.Error("Failed to create container", zap.Error(err))
		return err
	}
	an.ContainerID = cc.ID
	err = an.copyStartScript(ctx)
	if err != nil {
		an.log.Error("Failed to add entry script", zap.Error(err))
		return err
	}
	err = an.modifyGenesisToAddGenesisAccount(ctx, additionalGenesisWallets...)
	if err != nil {
		an.log.Error("Failed to add faucet account an container", zap.Error(err))
		return err
	}

	return nil
}

func (an *AnvilNode) copyStartScript(ctx context.Context) error {
	fileName := fmt.Sprintf("%s/test/chains/evm/foundry-docker/entrypoint.sh", os.Getenv(chains.BASE_PATH))

	config, err := interchaintest.GetLocalFileContent(fileName)

	err = an.WriteFile(ctx, config, "entrypoint.sh")
	return err
}

func (an *AnvilNode) modifyGenesisToAddGenesisAccount(ctx context.Context, additionalGenesisWallets ...ibc.WalletAmount) error {
	g := make(map[string]interface{})
	fileName := fmt.Sprintf("%s/test/chains/evm/foundry-docker/genesis.json", os.Getenv(chains.BASE_PATH))

	genbz, err := interchaintest.GetLocalFileContent(fileName)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(genbz, &g); err != nil {
		return fmt.Errorf("failed to unmarshal genesis file: %w", err)
	}

	for _, wallet := range additionalGenesisWallets {
		//g["alloc"][wallet.Address] = map[string]string{
		//	"balance": "0xd3c21bcecceda1000000", // 1_000_000*10**18
		//}
		//g["alloc"].(map[string]map[string]string)[wallet.Address] = map[string]string{
		//	"balance": "0xd3c21bcecceda1000000", // 1_000_000*10**18
		//}
		amount := map[string]string{
			"balance": "0xd3c21bcecceda1000000", // 1_000_000*10**18
		}
		if err := dyno.Set(g, amount, "alloc", wallet.Address); err != nil {
			return fmt.Errorf("failed to add genesis accounts an genesis json: %w", err)
		}
	}
	result, _ := json.Marshal(g)

	err = an.WriteFile(ctx, result, "genesis.json")
	return err
}

func (an *AnvilNode) CopyFileToContainer(ctx context.Context, content []byte, target string) error {
	header := ctx.Value("file-header").(map[string]string)
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	err := tw.WriteHeader(&tar.Header{
		Name: header["name"],
		Mode: 0644,
		Size: int64(len(content)),
	})
	_, err = tw.Write(content)
	if err != nil {
		return err
	}
	err = tw.Close()
	if err != nil {
		return err
	}
	if err := an.DockerClient.CopyToContainer(context.Background(), an.ContainerID, target, &buf, types.CopyToContainerOptions{}); err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

//
//func (in *AnvilNode) addFaucetAccount(ctx context.Context) error {
//
//	wallet, err := in.Chain.BuildWallet(ctx, interchaintest.FaucetAccountKeyName, "")
//	if err != nil {
//		return err
//	}
//
//	var accounts []map[string]interface{}
//	fileName := fmt.Sprintf("%s/test/chains/evm/hardhat-docker/accounts.json", os.Getenv(chains.BASE_PATH))
//
//	_accounts, err := interchaintest.GetLocalFileContent(fileName)
//	if err != nil {
//		return err
//	}
//
//	if err := json.Unmarshal(_accounts, &accounts); err != nil {
//		return fmt.Errorf("failed to unmarshal account file: %w", err)
//	}
//	fmt.Printf("new faucet %s\n", wallet.FormattedAddress())
//
//	accounts = append(accounts, map[string]interface{}{
//		"name":       "faucet",
//		"privateKey": wallet.Mnemonic(),
//		"address":    wallet.FormattedAddress(),
//		"balance":    "1000000000000000000000",
//	})
//	result, _internal := json.Marshal(accounts)
//
//	err = in.WriteFile(ctx, result, in.Chain.HomeDir()+"/accounts.json")
//	header := map[string]string{
//		"name": "accounts.json",
//	}
//	err = in.CopyFileToContainer(context.WithValue(ctx, "file-header", header), result, "/usr/src/app")
//	return err
//}

func (an *AnvilNode) HostName() string {
	return dockerutil.CondenseHostName(an.Name())
}

func (an *AnvilNode) Bind() []string {
	return []string{fmt.Sprintf("%s:%s", an.VolumeName, an.HomeDir())}
}

func (an *AnvilNode) HomeDir() string {
	return path.Join("/var/evm-chain", an.Chain.Config().Name)
}

func (an *AnvilNode) StartContainer(ctx context.Context) error {
	if err := dockerutil.StartContainer(ctx, an.DockerClient, an.ContainerID); err != nil {
		return err
	}

	c, err := an.DockerClient.ContainerInspect(ctx, an.ContainerID)
	if err != nil {
		return err
	}
	an.HostRPCPort = dockerutil.GetHostPort(c, rpcPort)
	an.logger().Info("EVM chain node started", zap.String("container", an.Name()), zap.String("rpc_port", an.HostRPCPort))

	uri := "http://" + an.HostRPCPort

	client, err := ethclient.Dial(uri)
	if err != nil {
		return err
	}
	an.Client = *client
	return nil
}

func (an *AnvilNode) logger() *zap.Logger {
	return an.log.With(
		zap.String("chain_id", an.Chain.Config().ChainID),
		zap.String("test", an.TestName),
	)
}

func (an *AnvilNode) Exec(ctx context.Context, cmd []string, env []string) ([]byte, []byte, error) {
	job := dockerutil.NewImage(an.logger(), an.DockerClient, an.NetworkID, an.TestName, an.Image.Repository, an.Image.Version)

	opts := dockerutil.ContainerOptions{
		Env:   env,
		Binds: an.Bind(),
	}
	res := job.Run(ctx, cmd, opts)
	return res.Stdout, res.Stderr, res.Err
}

func (an *AnvilNode) BinCommand(command ...string) []string {
	command = append([]string{an.Chain.Config().Bin}, command...)
	return command
}

func (an *AnvilNode) ExecBin(ctx context.Context, command ...string) ([]byte, []byte, error) {
	return an.Exec(ctx, an.BinCommand(command...), nil)
}

func (an *AnvilNode) GetBlockByHeight(ctx context.Context, height int64) (string, error) {
	an.lock.Lock()
	defer an.lock.Unlock()
	uri := "http://" + an.HostRPCPort + "/api/v3"
	block, _, err := an.ExecBin(ctx,
		"rpc", "blockbyheight", fmt.Sprint(height),
		"--uri", uri,
	)
	return string(block), err
}

func (an *AnvilNode) FindTxs(ctx context.Context, height uint64) ([]blockdb.Tx, error) {
	//var flag = true
	//if flag {
	//	time.Sleep(3 * time.Second)
	//	flag = false
	//}
	//
	//time.Sleep(2 * time.Second)
	//blockHeight := icontypes.BlockHeightParam{Height: icontypes.NewHexInt(int64(height))}
	//res, err := an.Client.GetBlockByHeight(&blockHeight)
	//if err != nil {
	//	return make([]blockdb.Tx, 0, 0), nil
	//}
	//txs := make([]blockdb.Tx, 0, len(res.NormalTransactions)+2)
	//var newTx blockdb.Tx
	//for _internal, tx := range res.NormalTransactions {
	//	newTx.Data = []byte(fmt.Sprintf(`{"data":"%s"}`, tx.Data))
	//}

	// ToDo Add events from block if any to newTx.Events.
	// Event is an alternative representation of tendermint/abci/types.Event
	//return txs, nil
	return nil, nil
}

func (an *AnvilNode) Height(ctx context.Context) (uint64, error) {
	header, err := an.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

func (an *AnvilNode) GetBalance(ctx context.Context, address string) (int64, error) {
	//addr := icontypes.AddressParam{Address: icontypes.Address(address)}
	//bal, err := an.Client.GetBalance(&addr)
	//return bal.Int64(), err
	return int64(0), nil
}

func (an *AnvilNode) DeployContract(ctx context.Context, contractPath, key string, params ...interface{}) (common.Address, error) {

	bytecode, contractABI, err := an.loadABI(contractPath)
	if err != nil {
		return common.Address{}, err
	}

	privateKey, fromAddress, err := getPrivateKey(key)

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

func (an *AnvilNode) loadABI(contractPath string) ([]byte, abi.ABI, error) {

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

func (an *AnvilNode) getNonceAndGasPrice(fromAddress common.Address) (uint64, *big.Int, error) {
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

func getPrivateKey(key string) (*ecdsa.PrivateKey, common.Address, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, common.Address{}, err
	}
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	return privateKey, fromAddress, nil
}

// Get Transaction result when hash is provided after executing a transaction
func (an *AnvilNode) TransactionResult(ctx context.Context, hash string) (*icontypes.TransactionResult, error) {
	uri := fmt.Sprintf("http://%s:9080/api/v3", an.Name()) //"http://" + an.HostRPCPort + "/api/v3"
	out, _, err := an.ExecBin(ctx, "rpc", "txresult", hash, "--uri", uri)
	if err != nil {
		return nil, err
	}
	var result = new(icontypes.TransactionResult)
	return result, json.Unmarshal(out, result)
}

// ExecTx executes a transaction, waits for 2 blocks if successful, then returns the tx hash.

// TxCommand is a helper to retrieve a full command for broadcasting a tx
// with the chain node binary.
func (an *AnvilNode) TxCommand(ctx context.Context, initMessage, filePath, keystorePath string, command ...string) []string {
	// get password from pathname as pathname will have the password prefixed. ex - Alice.Json
	_, key := filepath.Split(keystorePath)
	fileName := strings.Split(key, ".")
	password := fileName[0]

	command = append([]string{"run", "sendtx", "deploy", filePath}, command...)
	command = append(command,
		"--key_store", keystorePath,
		"--key_password", password,
		"--step_limit", "5000000000",
		"--content_type", "application/java",
	)

	return an.NodeCommand(command...)
}

// NodeCommand is a helper to retrieve a full command for a chain node binary.
// when interactions with the RPC endpoint are necessary.
// For example, if chain node binary is `gaiad`, and desired command is `gaiad keys show key1`,
// pass ("keys", "show", "key1") for command to return the full command.
// Will include additional flags for node URL, home directory, and chain ID.
func (an *AnvilNode) NodeCommand(command ...string) []string {
	command = an.BinCommand(command...)
	return append(command,
		"--rpc-url", fmt.Sprintf("http://%s", an.HostRPCPort),
	)
}

// CopyFile adds a file from the host filesystem to the docker filesystem
// relPath describes the location of the file in the docker volume relative to
// the home directory
func (an *AnvilNode) CopyFile(ctx context.Context, srcPath, dstPath string) error {
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	return an.WriteFile(ctx, content, dstPath)
}

// WriteFile accepts file contents in a byte slice and writes the contents to
// the docker filesystem. relPath describes the location of the file in the
// docker volume relative to the home directory
func (an *AnvilNode) WriteFile(ctx context.Context, content []byte, relPath string) error {
	fw := dockerutil.NewFileWriter(an.logger(), an.DockerClient, an.TestName)
	return fw.WriteFile(ctx, an.VolumeName, relPath, content)
}

func (an *AnvilNode) QueryContract(ctx context.Context, scoreAddress, methodName, params string) ([]byte, error) {
	uri := fmt.Sprintf("http://%s:9080/api/v3", an.Name())
	var args = []string{"rpc", "call", "--to", scoreAddress, "--method", methodName, "--uri", uri}
	if params != "" {
		var paramName = "--param"
		if strings.HasPrefix(params, "{") && strings.HasSuffix(params, "}") {
			paramName = "--raw"
		}
		args = append(args, paramName, params)
	}
	out, _, err := an.ExecBin(ctx, args...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (an *AnvilNode) RestoreKeystore(ctx context.Context, ks []byte, keyName string) error {
	return an.WriteFile(ctx, ks, keyName+".keystore")
}

func (an *AnvilNode) GetChainConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error) {
	//chainID, _internal := an.Client.NetworkID(context.Background())
	gasPrice, _ := an.Client.SuggestGasPrice(context.Background())

	config := &centralized.EVMRelayerChainConfig{
		Type: "evm",
		Value: centralized.EVMRelayerChainConfigValue{
			NID:             an.Chain.Config().ChainID,
			RPCURL:          an.Chain.GetRPCAddress(),
			StartHeight:     0,
			Keystore:        fmt.Sprintf("%s/keys/%s/%s", rlyHome, an.Chain.Config().ChainID, keyName),
			Password:        keyName,
			GasPrice:        gasPrice.Int64(),
			GasLimit:        2000000,
			ContractAddress: an.Chain.GetContractAddress("connection"), //cfg.ConfigFileOverrides["xcall-connection"].(string),
		},
	}
	return yaml.Marshal(config)
}

func (an *AnvilNode) ExecCallTx(ctx context.Context, contractAddress, methodName, pKey string, params ...interface{}) (*eth_types.Receipt, error) {
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
			an.Chain.GetRPCAddress(),
		}
		out, _, _ := an.ExecBin(ctx, cmd...)
		return receipt, fmt.Errorf("error on trasaction :: %s\n%v", methodName, string(out))
	}
	return receipt, nil

}

func (an *AnvilNode) ExecCallTxCommand(contractAddress, methodName, pKey string, params ...interface{}) []string {
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
