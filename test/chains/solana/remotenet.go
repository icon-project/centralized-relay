package solana

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/blockdb"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/dockerutil"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer/centralized"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"github.com/near/borsh-go"
	"gopkg.in/yaml.v3"

	//chantypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/docker/docker/client"
	"github.com/icon-project/centralized-relay/test/chains"

	"go.uber.org/zap"
)

type SolanaRemoteNet struct {
	log           *zap.Logger
	testName      string
	cfg           ibc.ChainConfig
	scorePaths    map[string]string
	IBCAddresses  map[string]string `json:"addresses"`
	appIds        map[string]IDL
	Wallets       map[string]ibc.Wallet `json:"wallets"`
	Client        *client.Client
	Network       string
	testconfig    *testconfig.Chain
	rpcClient     *Client
	walletPrivKey solana.PrivateKey
}

const xcallProgName = "xcall"
const connectionProgName = "centralized_connection"
const mockAppProgName = "mock_dapp"

func NewSolanaRemoteNet(testName string, log *zap.Logger, chainConfig ibc.ChainConfig, client *client.Client, network string, testconfig *testconfig.Chain) chains.Chain {
	solanaClient := &Client{
		rpc: solrpc.New(testconfig.RPCUri),
	}
	walletKey, err := solana.PrivateKeyFromSolanaKeygenFile(testconfig.ConfigPath + testconfig.KeystoreFile)
	if err != nil {
		return nil
	}
	return &SolanaRemoteNet{
		testName:      testName,
		cfg:           chainConfig,
		log:           log,
		scorePaths:    testconfig.Contracts,
		Wallets:       map[string]ibc.Wallet{},
		IBCAddresses:  make(map[string]string),
		Client:        client,
		testconfig:    testconfig,
		Network:       network,
		rpcClient:     solanaClient,
		appIds:        make(map[string]IDL),
		walletPrivKey: walletKey,
	}
}

func (sn *SolanaRemoteNet) CreateKey(context.Context, string) error {
	return nil
}

func (sn *SolanaRemoteNet) PauseNode(ctx context.Context) error {
	return nil
}

// UnpauseNode starts the paused node
func (sn *SolanaRemoteNet) UnpauseNode(ctx context.Context) error {
	return nil
}

// Config fetches the chain configuration.
func (sn *SolanaRemoteNet) Config() ibc.ChainConfig {
	return sn.cfg
}

func (sn *SolanaRemoteNet) OverrideConfig(key string, value any) {
	if value == nil {
		return
	}
	sn.cfg.ConfigFileOverrides[key] = value
}

// Initialize initializes node structs so that things like initializing keys can be done before starting the chain
func (sn *SolanaRemoteNet) Initialize(ctx context.Context, testName string, cli *client.Client, networkID string) error {
	return nil
}

// Start sets up everything needed (validators, gentx, fullnodes, peering, additional accounts) for chain to start from genesis.
func (sn *SolanaRemoteNet) Start(testName string, ctx context.Context, additionalGenesisWallets ...ibc.WalletAmount) error {
	return nil
}

// Exec runs an arbitrary command using Chain's docker environment.
// Whether the invoked command is run in a one-off container or execing into an already running container
// is up to the chain implementation.
//
// "env" are environment variables in the format "MY_ENV_VAR=value"
func (sn *SolanaRemoteNet) Exec(ctx context.Context, cmd []string, env []string) (stdout []byte, stderr []byte, err error) {
	cmd = append([]string{}, cmd...)
	job := dockerutil.NewImage(sn.log, sn.Client, sn.Network, sn.testName, sn.cfg.Images[0].Repository, sn.cfg.Images[0].Version)
	bindPaths := []string{
		sn.testconfig.ContractsPath + ":/workdir/",
		sn.testconfig.ConfigPath + ":/root/.config/solana/",
	}
	opts := dockerutil.ContainerOptions{
		Binds: bindPaths,
	}
	res := job.Run(ctx, cmd, opts)
	return res.Stdout, res.Stderr, res.Err
}

// ExportState exports the chain state at specific height.
func (sn *SolanaRemoteNet) ExportState(ctx context.Context, height int64) (string, error) {
	block, err := sn.GetClientBlockByHeight(ctx, height)
	return block, err
}

// GetRPCAddress retrieves the rpc address that can be reached by other containers in the docker network.
func (sn *SolanaRemoteNet) GetRPCAddress() string {
	return sn.testconfig.RPCUri
}

func (sn *SolanaRemoteNet) GetRelayConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error) {
	contracts := make(map[string]string)
	contracts["xcall"] = sn.GetContractAddress("xcall")
	contracts["connection"] = sn.GetContractAddress("connection")
	config := &centralized.SolanaRelayerChainConfig{
		Type: "solana",
		Value: centralized.SolanaRelayerChainConfigValue{
			NID:             sn.Config().ChainID,
			RPCUrl:          sn.GetRPCAddress(),
			StartHeight:     0,
			Contracts:       contracts,
			BlockInterval:   "6s",
			Address:         sn.testconfig.RelayWalletAddress,
			FinalityBlock:   uint64(0),
			MaxInclusionFee: 200,
		},
	}
	return yaml.Marshal(config)
}

// GetGRPCAddress retrieves the grpc address that can be reached by other containers in the docker network.
// Not Applicable for Icon
func (sn *SolanaRemoteNet) GetGRPCAddress() string {
	return sn.testconfig.WebsocketUrl
}

// GetHostRPCAddress returns the rpc address that can be reached by processes on the host machine.
// Note that this will not return a valid value until after Start returns.
func (sn *SolanaRemoteNet) GetHostRPCAddress() string {
	return sn.testconfig.RPCUri
}

// GetHostGRPCAddress returns the grpc address that can be reached by processes on the host machine.
// Note that this will not return a valid value until after Start returns.
// Not applicable for Icon
func (sn *SolanaRemoteNet) GetHostGRPCAddress() string {
	return sn.testconfig.RPCUri
}

// HomeDir is the home directory of a node running in a docker container. Therefore, this maps to
// the container's filesystem (not the host).
func (sn *SolanaRemoteNet) HomeDir() string {
	return ""
}

func (sn *SolanaRemoteNet) createKeystore(ctx context.Context, keyName string) (string, string, error) {
	return "", "", nil
}

// RecoverKey recovers an existing user from a given mnemonic.
func (sn *SolanaRemoteNet) RecoverKey(ctx context.Context, name string, mnemonic string) error {
	panic("not implemented") // TODO: Implement
}

// GetAddress fetches the bech32 address for a test key on the "user" node (either the first fullnode or the first validator if no fullnodes).
func (sn *SolanaRemoteNet) GetAddress(ctx context.Context, keyName string) ([]byte, error) {
	addrInByte, err := json.Marshal(keyName)
	if err != nil {
		return nil, err
	}
	return addrInByte, nil
}

// SendFunds sends funds to a wallet from a user account.
func (sn *SolanaRemoteNet) SendFunds(ctx context.Context, keyName string, amount ibc.WalletAmount) error {
	return nil
}

// Height returns the current block height or an error if unable to get current height.
func (sn *SolanaRemoteNet) Height(ctx context.Context) (uint64, error) {
	res, err := sn.rpcClient.GetLatestBlockHeight(ctx)
	return res, err
}

// GetGasFeesInNativeDenom gets the fees in native denom for an amount of spent gas.
func (sn *SolanaRemoteNet) GetGasFeesInNativeDenom(gasPaid int64) int64 {
	gasPrice, _ := strconv.ParseFloat(strings.Replace(sn.cfg.GasPrices, sn.cfg.Denom, "", 1), 64)
	fees := float64(gasPaid) * gasPrice
	return int64(fees)
}

// BuildRelayerWallet will return a chain-specific wallet populated with the mnemonic so that the wallet can
// be restored in the relayer node using the mnemonic. After it is built, that address is included in
// genesis with some funds.
func (sn *SolanaRemoteNet) BuildRelayerWallet(ctx context.Context, keyName string) (ibc.Wallet, error) {
	return sn.BuildWallet(ctx, keyName, "")
}

func (sn *SolanaRemoteNet) BuildWallet(ctx context.Context, keyName string, mnemonic string) (ibc.Wallet, error) {
	return nil, nil
}

func (sn *SolanaRemoteNet) FindTxs(ctx context.Context, height uint64) ([]blockdb.Tx, error) {
	return nil, nil
}

// GetBalance fetches the current balance for a specific account address and denom.
func (sn *SolanaRemoteNet) GetBalance(ctx context.Context, address string, denom string) (int64, error) {
	return 0, nil
}

func (sn *SolanaRemoteNet) populateConnIDL() error {
	connIdl := IDL{}
	filePath := sn.testconfig.ContractsPath + "/target/idl/centralized_connection.json"
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)
	if err = json.Unmarshal(bytes, &connIdl); err != nil {
		return err
	}
	sn.appIds[connectionProgName] = connIdl
	return nil
}

func (sn *SolanaRemoteNet) initConnection(ctx context.Context, connectionProgId string) error {

	params := make([]interface{}, 2)
	params = append(params, solana.MustPublicKeyFromBase58(sn.GetContractAddress(xcallProgName)))
	params = append(params, sn.walletPrivKey.PublicKey())

	seeds := [][]byte{
		[]byte("config"),
	}
	configAc, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(connectionProgId))
	if err != nil {
		log.Fatalf("Failed to find program address: %v", err)
	}
	seeds = [][]byte{
		[]byte("claim_fees"),
	}
	claimfeeAc, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(connectionProgId))
	if err != nil {
		log.Fatalf("Failed to find program address: %v", err)
	}
	if err != nil {
		return err
	}
	payerAccount := solana.AccountMeta{
		PublicKey:  sn.walletPrivKey.PublicKey(),
		IsWritable: true,
		IsSigner:   true,
	}
	accountsMeta := solana.AccountMetaSlice{
		{
			PublicKey:  configAc,
			IsWritable: true,
		},
		{
			PublicKey:  claimfeeAc,
			IsWritable: true,
		},
		&payerAccount,
		&solana.AccountMeta{
			PublicKey: solana.SystemProgramID,
		},
	}
	_, err = sn.executeContract(ctx, connectionProgName, "initialize", params, accountsMeta)
	return err
}

type Config struct {
	Admin solana.PublicKey
	Xcall solana.PublicKey
	Sn    uint64
	Bump  []byte
}

type XcConfig struct {
	Admin       [32]byte
	FeeHandler  solana.PublicKey
	NetworkId   [32]byte
	ProtocolFee uint64
	SequenceNo  uint64
	LastReqId   uint64
}

func (sn *SolanaRemoteNet) SetupConnection(ctx context.Context, target chains.Chain) error {
	if sn.testconfig.Environment == "preconfigured" {
		return nil
	}

	// if err := sn.syncProgramId(ctx, connectionProgName); err != nil {
	// 	return err
	// }
	connection, err := sn.DeployContractRemote(ctx, connectionProgName)
	if err != nil {
		return err
	}
	sn.IBCAddresses["connection"] = connection
	fmt.Println("Connection deployed at ", connection)
	time.Sleep(3 * time.Second)

	if err := sn.populateConnIDL(); err != nil {
		return err
	}
	//initialize cc
	return sn.initConnection(ctx, connection)

}

func (sn *SolanaRemoteNet) populateXcallIDL() error {
	xcallIdl := IDL{}
	filePath := sn.testconfig.ContractsPath + "/target/idl/xcall.json"
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)
	if err = json.Unmarshal(bytes, &xcallIdl); err != nil {
		return err
	}
	sn.appIds[xcallProgName] = xcallIdl
	return nil
}

func (sn *SolanaRemoteNet) initXcall(ctx context.Context, xcallProgId string) error {
	params := make([]interface{}, 1)
	params = append(params, sn.testconfig.ChainConfig.ChainID)

	seeds := [][]byte{
		[]byte("config"),
	}
	configAc, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(xcallProgId))
	if err != nil {
		log.Fatalf("Failed to find program address: %v", err)
	}
	seeds = [][]byte{
		[]byte("reply"),
	}
	replyPda, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(xcallProgId))
	if err != nil {
		log.Fatalf("Failed to find program address: %v", err)
	}
	if err != nil {
		return err
	}
	payerAccount := solana.AccountMeta{
		PublicKey:  sn.walletPrivKey.PublicKey(),
		IsWritable: true,
		IsSigner:   true,
	}
	accountsMeta := solana.AccountMetaSlice{
		{
			PublicKey:  configAc,
			IsWritable: true,
		},
		{
			PublicKey:  replyPda,
			IsWritable: true,
		},
		&payerAccount,
		&solana.AccountMeta{
			PublicKey: solana.SystemProgramID,
		},
	}
	_, err = sn.executeContract(ctx, xcallProgName, "initialize", params, accountsMeta)
	return err
}

func (sn *SolanaRemoteNet) SetupXCall(ctx context.Context) error {
	if true {
		sn.IBCAddresses["xcall"] = "4XhdhgaaWYfJHSSQ7RVbyCJPZx3LLxB8fz35X36pchLh"
		return nil
	}
	if sn.testconfig.Environment == "preconfigured" {
		testcase := ctx.Value("testcase").(string)
		sn.IBCAddresses["xcall"] = "13irkyf42gWeAULwsqmvo8yVkgJix5AP8aZrmed24AbD"
		sn.IBCAddresses["connection"] = "4oT1bBWP2jLxdDchNGmuVJhFDxkCkqjq9hSgDCcv8KWK"
		sn.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = "CA7VNLUTPQZEKAXJHXHCS2YJHPELODRZKB26MBQQISD4BPHG2PGAUPMM"
		return nil
	}
	// The build should be multiple times syncing the programid
	// if err := sn.syncProgramId(ctx, xcallProgName); err != nil {
	// 	return err
	// }
	xcallProgId, err := sn.DeployContractRemote(ctx, "xcall")
	if err != nil {
		return err
	}
	sn.IBCAddresses["xcall"] = xcallProgId
	fmt.Println("Xcall deployed at", xcallProgId)

	//populate IDL for xcall
	if err := sn.populateXcallIDL(); err != nil {
		return err
	}

	//init xcall
	return sn.initXcall(ctx, xcallProgId)
}

func (sn *SolanaRemoteNet) DeployXCallMockApp(ctx context.Context, keyName string, connections []chains.XCallConnection) error {
	if sn.testconfig.Environment == "preconfigured" {
		return nil
	}
	testcase := ctx.Value("testcase").(string)

	dapp, err := sn.DeployContractRemote(ctx, mockAppProgName)
	if err != nil {
		return err
	}
	sn.IBCAddresses["dapp"] = dapp
	fmt.Println("Dapp deployed at ", dapp)

	//init dapp
	xCall := sn.IBCAddresses["xcall"]
	params := make(map[string]string)
	params["xcall_address"] = xCall
	// sn.executeContract(ctx, dapp, "init", params)

	sn.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = dapp
	for _, connection := range connections {
		params := make(map[string]string)
		params["src_endpoint"] = sn.IBCAddresses[connection.Connection]
		params["network_id"] = connection.Nid
		params["dst_endpoint"] = connection.Destination
		// _, err = sn.executeContract(context.Background(), dapp, "add_connection", params)
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

func (sn *SolanaRemoteNet) GetContractAddress(key string) string {
	value, exist := sn.IBCAddresses[key]
	if !exist {
		panic(fmt.Sprintf(`IBC address not exist %s`, key))
	}
	return value
}

func (sn *SolanaRemoteNet) BackupConfig() ([]byte, error) {
	panic("not implemented")
}

func (sn *SolanaRemoteNet) RestoreConfig(backup []byte) error {
	panic("not implemented")
}

func (sn *SolanaRemoteNet) SendPacketXCall(ctx context.Context, keyName, _to string, data, rollback []byte) (context.Context, error) {
	// testcase := ctx.Value("testcase").(string)
	// dappKey := fmt.Sprintf("dapp-%s", testcase)
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

	// ctx, err := sn.executeContract(ctx, sn.IBCAddresses[dappKey], "send_call_message", params)
	// if err != nil {
	// 	return nil, err
	// }
	// sno := ctx.Value("sno").(string)
	// return context.WithValue(ctx, "sn", sno), nil
	return ctx, nil
}

// HasPacketReceipt returns the receipt of the packet sent to the target chain
func (sn *SolanaRemoteNet) IsPacketReceived(ctx context.Context, params map[string]interface{}, order ibc.Order) bool {
	panic("not implemented")
}

// FindTargetXCallMessage returns the request id and the data of the message sent to the target chain
func (sn *SolanaRemoteNet) FindTargetXCallMessage(ctx context.Context, target chains.Chain, height uint64, to string) (*chains.XCallResponse, error) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	sno := ctx.Value("sn").(string)
	reqId, destData, err := target.FindCallMessage(ctx, height, sn.cfg.ChainID+"/"+sn.IBCAddresses[dappKey], to, sno)
	return &chains.XCallResponse{SerialNo: sno, RequestID: reqId, Data: destData}, err
}

func (sn *SolanaRemoteNet) XCall(ctx context.Context, targetChain chains.Chain, keyName, to string, data, rollback []byte) (*chains.XCallResponse, error) {
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

func (sn *SolanaRemoteNet) ExecuteCall(ctx context.Context, reqId, data string) (context.Context, error) {
	panic("not required in e2e")
}

func (sn *SolanaRemoteNet) ExecuteRollback(ctx context.Context, sno string) (context.Context, error) {
	return ctx, nil

}

func (sn *SolanaRemoteNet) FindCallMessage(ctx context.Context, startHeight uint64, from, to, sno string) (string, string, error) {

	event, err := sn.FindEvent(ctx, startHeight, "xcall", "CallMessage", sno)
	if err != nil {
		return "", "", err
	}
	reqId := event.ValueDecoded["reqId"].(uint64)
	data := event.ValueDecoded["data"].([]byte)
	return strconv.FormatUint(reqId, 10), string(data), nil
}

func (sn *SolanaRemoteNet) FindCallResponse(ctx context.Context, startHeight uint64, sno string) (string, error) {
	event, err := sn.FindEvent(ctx, startHeight, "xcall", "ResponseMessage", sno)
	if err != nil {
		return "", err
	}
	code := event.ValueDecoded["code"].(uint64)
	return strconv.FormatUint(code, 10), nil
}

func (sn *SolanaRemoteNet) FindEvent(ctx context.Context, startHeight uint64, contract, signature, sno string) (*EventResponseEvent, error) {
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

func (sn *SolanaRemoteNet) getEvent(ctx context.Context, startHeight uint64, sno, signature, contractId string) (*EventResponseEvent, error) {
	return nil, nil
	// TODO add event get implementation
	// return sn.rpcClient.GetEvent(ctx, startHeight, sno, contractId, signature)
}

// Remote implements chains.Chain
func (sn *SolanaRemoteNet) DeployContract(ctx context.Context, keyName string) (context.Context, error) {
	return ctx, nil
}

// executeContract implements chains.Chain
func (sn *SolanaRemoteNet) executeContract(ctx context.Context, contractProgName, methodName string, params []interface{}, accountValues solana.AccountMetaSlice) (context.Context, error) {
	progIdl := sn.appIds[contractProgName]
	discriminator, err := progIdl.GetInstructionDiscriminator(methodName)
	if err != nil {
		return nil, err
	}
	progID, err := progIdl.GetProgramID()
	if err != nil {
		return nil, err
	}
	instructionData := discriminator
	for _, param := range params {
		serParam, err := borsh.Serialize(param)
		if err != nil {
			return nil, err
		}
		instructionData = append(instructionData, serParam...)
	}
	signers := []solana.PrivateKey{sn.walletPrivKey}
	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID:        progID,
			AccountValues: accountValues,
			DataBytes:     instructionData,
		},
	}
	_, err = sn.executeTx(ctx, sn.walletPrivKey, instructions, signers)
	return ctx, err
}

func (sn *SolanaRemoteNet) executeTx(ctx context.Context, payerPrivKey solana.PrivateKey, instructions []solana.Instruction, signers []solana.PrivateKey) (*solrpc.SignatureStatusesResult, error) {
	latestBlockHash, err := sn.rpcClient.GetLatestBlockHash(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block hash: %w", err)
	}
	fmt.Println("last block hash", latestBlockHash)

	tx, err := solana.NewTransaction(instructions, *latestBlockHash,
		solana.TransactionPayer(payerPrivKey.PublicKey()))
	if err != nil {
		return nil, fmt.Errorf("failed to create new tx: %w", err)
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			for _, signer := range signers {
				if signer.PublicKey() == key {
					return &signer
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign tx: %w", err)
	}
	fmt.Println(tx.ToBase64())
	txSign, err := sn.rpcClient.SendTx(ctx, tx, nil)
	if err != nil {
		fmt.Errorf("failed to send tx: %w", err)
	}

	fmt.Println("Tx Send Successful:", txSign)

	txResult, err := sn.waitForTxConfirmation(10*time.Second, txSign)
	if err != nil {
		fmt.Printf("error waiting for tx confirmation: %w", err)
	}
	sn.log.Info("send message successful", zap.String("tx-hash", txSign.String()))
	return txResult, nil
}

func (sn *SolanaRemoteNet) ExecuteContract(ctx context.Context, contractAddress, keyName, methodName string, params map[string]interface{}) (context.Context, error) {
	return nil, nil
}

func (sn *SolanaRemoteNet) GetBlockByHeight(context.Context) (context.Context, error) {
	panic("not implemented")
}

// GetBlockByHeight implements chains.Chain
func (sn *SolanaRemoteNet) GetClientBlockByHeight(ctx context.Context, height int64) (string, error) {
	return "", nil
}

// GetLastBlock implements chains.Chain
func (sn *SolanaRemoteNet) GetLastBlock(ctx context.Context) (context.Context, error) {
	res, err := sn.rpcClient.GetLatestBlockHeight(ctx)
	return context.WithValue(ctx, chains.LastBlock{}, res), err
}

func (sn *SolanaRemoteNet) InitEventListener(ctx context.Context, contract string) chains.EventListener {
	return nil
}

// QueryContract implements chains.Chain
func (sn *SolanaRemoteNet) QueryContract(ctx context.Context, contractAddress, methodName string, params map[string]interface{}) (context.Context, error) {
	return ctx, nil
}

func (sn *SolanaRemoteNet) BuildWallets(ctx context.Context, keyName string) (ibc.Wallet, error) {
	panic("not implemented")
}

func (sn *SolanaRemoteNet) NodeCommand(command ...string) []string {
	command = sn.BinCommand(command...)
	return command
}

func (sn *SolanaRemoteNet) BinCommand(command ...string) []string {
	command = append([]string{sn.Config().Bin}, command...)
	return command
}

func (sn *SolanaRemoteNet) ExecBin(ctx context.Context, command ...string) ([]byte, []byte, error) {
	return sn.Exec(ctx, sn.BinCommand(command...), nil)
}

func (sn *SolanaRemoteNet) DeployContractRemote(ctx context.Context, contractName string) (string, error) {
	// Deploy the contract
	contractId, err := sn.ExecTx(ctx, contractName)
	if err != nil {
		return "", err
	}
	return contractId, nil
}

func (sn *SolanaRemoteNet) ExecTx(ctx context.Context, contractName string, command ...string) (string, error) {
	stdout, _, err := sn.Exec(ctx, sn.TxCommand(ctx, contractName, command...), nil)
	if err != nil {
		return "", err
	}
	return getProgramIdFromDeployment(string(stdout))
}

// TxCommand is a helper to retrieve a full command for broadcasting a tx
// with the chain node binary.
func (sn *SolanaRemoteNet) TxCommand(ctx context.Context, contractName string, command ...string) []string {
	command = append([]string{"deploy", "--provider.cluster", sn.testconfig.RPCUri, "--program-name", contractName}, command...)
	return sn.NodeCommand(command...)
}

func (sn *SolanaRemoteNet) syncProgramId(ctx context.Context, contractName string) error {
	//build
	command := []string{"build", "--provider.cluster", sn.testconfig.RPCUri, "--program-name", contractName}
	_, _, err := sn.Exec(ctx, sn.NodeCommand(command...), nil)
	if err != nil {
		return err
	}

	//sync the key
	command = []string{"keys", "sync", "--provider.cluster", sn.testconfig.RPCUri, "--program-name", contractName}
	_, _, err = sn.Exec(ctx, sn.NodeCommand(command...), nil)
	if err != nil {
		return err
	}

	//build again
	command = []string{"build", "--provider.cluster", sn.testconfig.RPCUri, "--program-name", contractName}
	_, _, err = sn.Exec(ctx, sn.NodeCommand(command...), nil)
	if err != nil {
		return err
	}
	return nil
}

func getProgramIdFromDeployment(input string) (string, error) {
	re := regexp.MustCompile(`Program Id: (\w+)`)

	match := re.FindStringSubmatch(input)

	if len(match) > 1 {
		programId := match[1]
		return programId, nil
	}
	return "", errors.New("programid not found")
}

func (sn *SolanaRemoteNet) FindRollbackExecutedMessage(ctx context.Context, startHeight uint64, sno string) (string, error) {
	event, err := sn.FindEvent(ctx, startHeight, "xcall", "RollbackExecuted", sno)
	if err != nil {
		return "", err
	}
	fsno := event.ValueDecoded["sn"].(uint64)
	return strconv.FormatUint(fsno, 10), nil
}

// helper methods

func (sn *SolanaRemoteNet) waitForTxConfirmation(timeout time.Duration, sign solana.Signature) (*solrpc.SignatureStatusesResult, error) {
	startTime := time.Now()
	for range time.NewTicker(500 * time.Millisecond).C {
		txStatus, err := sn.rpcClient.GetSignatureStatus(context.TODO(), true, sign)
		if err == nil && txStatus != nil && (txStatus.ConfirmationStatus == solrpc.ConfirmationStatusConfirmed || txStatus.ConfirmationStatus == solrpc.ConfirmationStatusFinalized) {
			return txStatus, nil
		} else if time.Since(startTime) > timeout {
			var cbErr error
			if err != nil {
				cbErr = err
			} else if txStatus != nil && txStatus.Err != nil {
				cbErr = fmt.Errorf("failed to get tx signature status: %v", txStatus.Err)
			} else {
				cbErr = fmt.Errorf("failed to finalize tx signature")
			}
			return nil, cbErr
		}
	}
	return nil, fmt.Errorf("request timeout")
}
