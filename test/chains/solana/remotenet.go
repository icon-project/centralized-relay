package solana

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/icon-project/centralized-relay/test/chains/solana/alt"
	"github.com/icon-project/centralized-relay/test/interchaintest/_internal/dockerutil"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer/centralized"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"github.com/near/borsh-go"
	"gopkg.in/yaml.v3"

	"github.com/docker/docker/client"
	"github.com/icon-project/centralized-relay/test/chains"

	"go.uber.org/zap"
)

const (
	defaultTxConfirmationTime = 2 * time.Second
	accountsQueryMaxLimit     = 30
)

type SolanaRemoteNet struct {
	log           *zap.Logger
	testName      string
	cfg           chains.ChainConfig
	scorePaths    map[string]string
	IBCAddresses  map[string]string `json:"addresses"`
	appIds        map[string]IDL
	pdaAcs        map[string]solana.PublicKey
	Wallets       map[string]ibc.Wallet `json:"wallets"`
	Client        *client.Client
	Network       string
	testconfig    *testconfig.Chain
	rpcClient     *Client
	walletPrivKey solana.PrivateKey
	CNids         []string
}

const (
	xcallProgName       = "xcall"
	xcallIdlFile        = "xcall.json"
	connectionProgName  = "centralized_connection"
	connectionIdlFile   = "centralized_connection.json"
	mockAppProgName     = "mock_dapp_multi"
	mockDappIdlFilePath = "mock_dapp_multi.json"
)

func NewSolanaRemoteNet(testName string, log *zap.Logger, chainConfig chains.ChainConfig, client *client.Client, network string, testconfig *testconfig.Chain) chains.Chain {
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
		pdaAcs:        make(map[string]solana.PublicKey),
		walletPrivKey: walletKey,
	}
}

// Config fetches the chain configuration.
func (sn *SolanaRemoteNet) Config() chains.ChainConfig {
	return sn.cfg
}

// Exec runs an arbitrary command using Chain's docker environment.
// Whether the invoked command is run in a one-off container or execing into an already running container
// is up to the chain implementation.
//
// "env" are environment variables in the format "MY_ENV_VAR=value"
func (sn *SolanaRemoteNet) Exec(ctx context.Context, cmd []string, env []string) (stdout []byte, stderr []byte, err error) {
	cmd = append([]string{}, cmd...)
	job := dockerutil.NewImage(sn.log, sn.Client, sn.Network, sn.testName, sn.cfg.Images.Repository, sn.cfg.Images.Version)
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

// GetRPCAddress retrieves the rpc address that can be reached by other containers in the docker network.
func (sn *SolanaRemoteNet) GetRPCAddress() string {
	return sn.testconfig.RPCUri
}

func (sn *SolanaRemoteNet) GetRelayConfig(ctx context.Context, rlyHome string, keyName string) ([]byte, error) {
	testcase := ctx.Value("testcase").(string)
	config := &centralized.SolanaRelayerChainConfig{
		Type: "solana",
		Value: centralized.SolanaRelayerChainConfigValue{
			NID:               sn.Config().ChainID,
			RPCUrl:            sn.GetRPCAddress(),
			StartHeight:       0,
			Address:           sn.testconfig.RelayWalletAddress,
			XcallProgram:      sn.GetContractAddress("xcall"),
			ConnectionProgram: sn.GetContractAddress("connection"),
			AltAddress:        sn.pdaAcs["altAddress"].String(),
			CpNIDs:            sn.CNids,
			Dapps: []centralized.Dapp{
				{
					Name:      "mock_dapp",
					ProgramID: sn.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)],
				},
			},
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

// GetAddress fetches the bech32 address for a test key on the "user" node (either the first fullnode or the first validator if no fullnodes).
func (sn *SolanaRemoteNet) GetAddress(ctx context.Context, keyName string) ([]byte, error) {
	addrInByte, err := json.Marshal(keyName)
	if err != nil {
		return nil, err
	}
	return addrInByte, nil
}

// Height returns the current block height or an error if unable to get current height.
func (sn *SolanaRemoteNet) Height(ctx context.Context) (uint64, error) {
	res, err := sn.rpcClient.GetLatestBlockHeight(ctx, solrpc.CommitmentFinalized)
	return res, err
}

func (sn *SolanaRemoteNet) populateIDL(ctx context.Context, idlFilePath,
	programName string, pushIdl bool) error {
	progIdl := IDL{}
	filePath := sn.testconfig.ContractsPath + "/target/idl/" + idlFilePath
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)
	if err = json.Unmarshal(bytes, &progIdl); err != nil {
		return err
	}
	if pushIdl {
		programId, err := progIdl.GetProgramID()
		if err != nil {
			return err
		}
		err = sn.syncIdl(ctx, "/workdir/target/idl/"+idlFilePath, programId.String())
		if err != nil {
			return err
		}
	}
	sn.appIds[programName] = progIdl
	return nil
}

func (sn *SolanaRemoteNet) setCnFee(ctx context.Context, network string) error {
	var params []interface{}
	params = append(params, network)
	params = append(params, uint64(0))
	params = append(params, uint64(0))

	payerAccount := solana.AccountMeta{
		PublicKey:  sn.walletPrivKey.PublicKey(),
		IsWritable: true,
		IsSigner:   true,
	}
	seeds := [][]byte{
		[]byte("config"),
	}
	configAc, _, err := solana.FindProgramAddress(seeds,
		solana.MustPublicKeyFromBase58(sn.IBCAddresses["connection"]))
	if err != nil {
		log.Fatalf("Failed to find program address: %v", err)
	}
	seeds = [][]byte{
		[]byte("fee"),
		[]byte(network),
	}
	networkFeeAc, _, err := solana.FindProgramAddress(seeds,
		solana.MustPublicKeyFromBase58(sn.IBCAddresses["connection"]))
	if err != nil {
		log.Fatalf("Failed to find program address: %v", err)
	}
	accountsMeta := solana.AccountMetaSlice{
		&payerAccount,
		&solana.AccountMeta{
			PublicKey: solana.SystemProgramID,
		},
		{
			PublicKey:  networkFeeAc,
			IsWritable: true,
		},
		{
			PublicKey:  configAc,
			IsWritable: true,
		},
	}
	_, _, err = sn.executeContract(ctx, connectionProgName, "set_fee", params, accountsMeta)
	return err
}
func (sn *SolanaRemoteNet) initConnection(ctx context.Context, connectionProgId string) error {
	sn.log.Info("initializing connection..")
	var params []interface{}
	params = append(params, solana.MustPublicKeyFromBase58(sn.GetContractAddress(xcallProgName)))
	params = append(params, sn.walletPrivKey.PublicKey())

	seeds := [][]byte{
		[]byte("config"),
	}
	configAc, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(connectionProgId))
	if err != nil {
		log.Fatalf("Failed to find program address: %v", err)
	}

	sn.pdaAcs["connconfig"] = configAc
	payerAccount := solana.AccountMeta{
		PublicKey:  sn.walletPrivKey.PublicKey(),
		IsWritable: true,
		IsSigner:   true,
	}
	seeds = [][]byte{
		[]byte("connection_authority"),
	}
	connAuthAc, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(connectionProgId))
	if err != nil {
		log.Fatalf("Failed to find program address: %v", err)
	}
	accountsMeta := solana.AccountMetaSlice{
		&payerAccount,
		&solana.AccountMeta{
			PublicKey: solana.SystemProgramID,
		},
		{
			PublicKey:  configAc,
			IsWritable: true,
		},
		{
			PublicKey:  connAuthAc,
			IsWritable: true,
		},
	}
	_, _, err = sn.executeContract(ctx, connectionProgName, "initialize", params, accountsMeta)
	if err != nil {
		return err
	}
	time.Sleep(30 * time.Second)
	return err
}

func (sn *SolanaRemoteNet) SetupConnection(ctx context.Context, target chains.Chain) error {
	if sn.testconfig.Environment == "preconfigured" {
		return nil
	}
	if sn.IBCAddresses["connection"] != "" {
		return nil
	}
	if err := sn.syncProgramId(ctx, connectionProgName); err != nil {
		return err
	}
	connection, err := sn.DeployContractRemote(ctx, connectionProgName)
	if err != nil {
		return err
	}
	sn.IBCAddresses["connection"] = connection
	sn.log.Info("Connection deployed at ", zap.String("connection", connection))

	if err := sn.populateIDL(ctx, connectionIdlFile, connectionProgName, true); err != nil {
		return err
	}
	time.Sleep(30 * time.Second)
	return sn.initConnection(ctx, connection)

}

func (sn *SolanaRemoteNet) initXcall(ctx context.Context, xcallProgId string) error {
	sn.log.Info("initializing xcall..")
	var params []interface{}
	params = append(params, sn.testconfig.ChainConfig.ChainID)

	seeds := [][]byte{
		[]byte("config"),
	}
	configAc, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(xcallProgId))
	if err != nil {
		log.Fatalf("Failed to find program address: %v", err)
	}
	sn.pdaAcs["configAcXcall"] = configAc
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
		&payerAccount,
		&solana.AccountMeta{
			PublicKey: solana.SystemProgramID,
		},
	}
	_, _, err = sn.executeContract(ctx, xcallProgName, "initialize", params, accountsMeta)
	return err
}

type XcallConfigAccount struct {
	Admin       solana.PublicKey
	FeeHandler  solana.PublicKey
	NetworkID   string
	ProtocolFee uint64
	SequenceNo  big.Int
	LastReqID   big.Int
	Bump        uint8
}

type DappConfigAccount struct {
	Sn           big.Int
	XcallAddress solana.PublicKey
	Bump         uint8
}

type ConnAc struct {
	Connections []Conn
}

type Conn struct {
	SrcEndpoint string
	DstEndpoint string
}

func (sn *SolanaRemoteNet) SetupXCall(ctx context.Context) error {
	if sn.testconfig.Environment == "preconfigured" {
		testcase := ctx.Value("testcase").(string)
		xcallAccount := "ATcxAHKzczdHSRhiqJ9u6qDKoup2ah2ptfdrb1jmF6Yd"
		connAccount := "2USWPDrA4Yfm6jsgDPT8p3uYVFrq86hnSU82hgENhzEC"
		dappAccount := "5Xiwde2Jdxt5xCxbMW7xdt66mx5Gkuhd7BrnM6SMQxbU"
		sn.IBCAddresses["xcall"] = xcallAccount
		sn.IBCAddresses["connection"] = connAccount
		sn.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = dappAccount
		if err := sn.populateIDL(ctx, xcallIdlFile, xcallProgName, false); err != nil {
			return err
		}
		if err := sn.populateIDL(ctx, connectionIdlFile, connectionProgName, false); err != nil {
			return err
		}
		if err := sn.populateIDL(ctx, mockDappIdlFilePath, mockAppProgName, false); err != nil {
			return err
		}
		seeds := [][]byte{
			[]byte("config"),
		}
		dappConfigAc, _, err := solana.FindProgramAddress(seeds,
			solana.MustPublicKeyFromBase58(dappAccount))
		if err != nil {
			log.Fatalf("Failed to find program address: %v", err)
			return err
		}
		sn.pdaAcs["configAcDapp"] = dappConfigAc

		seeds = [][]byte{
			[]byte("dapp_authority"),
		}
		dappAuthorityAc, _, err := solana.FindProgramAddress(seeds,
			solana.MustPublicKeyFromBase58(dappAccount))
		if err != nil {
			log.Fatalf("Failed to find program address: %v", err)
			return err
		}
		sn.pdaAcs["dappAuthorityAc"] = dappAuthorityAc

		seeds = [][]byte{
			[]byte("config"),
		}
		configAc, _, err := solana.FindProgramAddress(seeds,
			solana.MustPublicKeyFromBase58(xcallAccount))
		if err != nil {
			log.Fatalf("Failed to find program address: %v", err)
		}
		sn.pdaAcs["configAcXcall"] = configAc

		xcallConfigAc := XcallConfigAccount{}
		sn.rpcClient.GetAccountInfo(ctx, configAc.String(), &xcallConfigAc)
		sn.pdaAcs["feeHandler"] = xcallConfigAc.FeeHandler

		seeds = [][]byte{
			[]byte("config"),
		}
		configAc, _, err = solana.FindProgramAddress(seeds,
			solana.MustPublicKeyFromBase58(connAccount))
		if err != nil {
			log.Fatalf("Failed to find program address: %v", err)
		}

		sn.pdaAcs["connconfig"] = configAc

		altPkey, err := sn.createLookupTableAccount(ctx)
		if err != nil {
			return err
		}
		sn.pdaAcs["altAddress"] = *altPkey
		time.Sleep(30 * time.Second)
		return nil
	}

	if err := sn.syncProgramId(ctx, xcallProgName); err != nil {
		return err
	}
	xcallProgId, err := sn.DeployContractRemote(ctx, "xcall")
	if err != nil {
		return err
	}
	sn.IBCAddresses["xcall"] = xcallProgId
	sn.log.Info("Xcall deployed at", zap.String("address", xcallProgId))

	//populate IDL for xcall
	if err := sn.populateIDL(ctx, xcallIdlFile, xcallProgName, true); err != nil {
		return err
	}
	time.Sleep(30 * time.Second)
	//init xcall
	return sn.initXcall(ctx, xcallProgId)
}

func (sn *SolanaRemoteNet) DeployXCallMockApp(ctx context.Context, keyName string, connections []chains.XCallConnection) error {
	if sn.testconfig.Environment == "preconfigured" {
		return nil
	}
	altPkey, err := sn.createLookupTableAccount(ctx)
	if err != nil {
		return err
	}
	sn.pdaAcs["altAddress"] = *altPkey

	configAc := sn.pdaAcs["configAcXcall"]

	xcallConfigAc := XcallConfigAccount{}
	sn.rpcClient.GetAccountInfo(ctx, configAc.String(), &xcallConfigAc)
	sn.pdaAcs["feeHandler"] = xcallConfigAc.FeeHandler

	rollbackSd := xcallConfigAc.SequenceNo.Add(&xcallConfigAc.SequenceNo, big.NewInt(1))
	rollbackSeeds := [][]byte{
		[]byte("rollback"),
		[]byte(rollbackSd.FillBytes(make([]byte, 16))),
	}
	rollbackAc, _, err := solana.FindProgramAddress(rollbackSeeds,
		solana.MustPublicKeyFromBase58(sn.IBCAddresses["xcall"]))
	if err != nil {
		log.Fatalf("Failed to find program address: %v", err)
	}
	sn.pdaAcs["rollbackAc"] = rollbackAc
	testcase := ctx.Value("testcase").(string)
	if err := sn.syncProgramId(ctx, mockAppProgName); err != nil {
		return err
	}
	dapp, err := sn.DeployContractRemote(ctx, mockAppProgName)
	if err != nil {
		return err
	}
	sn.IBCAddresses["dapp"] = dapp
	sn.log.Info("Dapp deployed at ", zap.String("address", dapp))
	if err := sn.populateIDL(ctx, mockDappIdlFilePath, mockAppProgName, true); err != nil {
		return err
	}
	time.Sleep(30 * time.Second)
	//init dapp
	xCall := solana.MustPublicKeyFromBase58(sn.IBCAddresses["xcall"])

	var params []interface{}
	params = append(params, xCall)
	payerAccount := solana.AccountMeta{
		PublicKey:  sn.walletPrivKey.PublicKey(),
		IsWritable: true,
		IsSigner:   true,
	}

	seeds := [][]byte{
		[]byte("config"),
	}
	configAc, _, err = solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(dapp))
	if err != nil {
		log.Fatalf("Failed to find program address: %v", err)
		return err
	}
	sn.pdaAcs["configAcDapp"] = configAc

	seeds = [][]byte{
		[]byte("dapp_authority"),
	}
	dappAuthorityAc, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(dapp))
	if err != nil {
		log.Fatalf("Failed to find program address: %v", err)
		return err
	}
	sn.pdaAcs["dappAuthorityAc"] = dappAuthorityAc
	accountsMeta := solana.AccountMetaSlice{
		{
			PublicKey:  configAc,
			IsWritable: true,
		}, {
			PublicKey:  dappAuthorityAc,
			IsWritable: true,
		},
		&payerAccount,
		&solana.AccountMeta{
			PublicKey: solana.SystemProgramID,
		},
	}
	_, _, err = sn.executeContract(ctx, mockAppProgName, "initialize", params, accountsMeta)
	if err != nil {
		sn.log.Warn("error occurred", zap.Error(err))
		return err
	}
	time.Sleep(30 * time.Second)
	dappConfigAcX := DappConfigAccount{}
	sn.rpcClient.GetAccountInfo(ctx, configAc.String(), &dappConfigAcX)
	sn.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)] = dapp
	sn.log.Info("Adding connections...")
	for _, connection := range connections {
		var params []interface{}
		params = append(params, connection.Nid) //
		params = append(params, sn.IBCAddresses[connection.Connection])
		params = append(params, connection.Destination)
		sn.CNids = append(sn.CNids, connection.Nid)
		seeds = [][]byte{
			[]byte("connections"),
			[]byte(connection.Nid),
		}
		connAc, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(dapp))
		if err != nil {
			log.Fatalf("Failed to find program address: %v", err)
			return err
		}
		accountsMeta = solana.AccountMetaSlice{
			{
				PublicKey:  connAc,
				IsWritable: true,
			},
			&payerAccount,
			&solana.AccountMeta{
				PublicKey: solana.SystemProgramID,
			},
		}
		_, _, err = sn.executeContract(ctx, mockAppProgName, "add_connection",
			params, accountsMeta)
		if err != nil {
			sn.log.Error("Unable to add connection",
				zap.Error(err),
				zap.String("nid", connection.Nid),
				zap.String("source", sn.IBCAddresses[connection.Connection]),
				zap.String("destination", connection.Destination),
			)
		}

		sn.log.Info("setting fee")
		err = sn.setCnFee(ctx, connection.Nid)
		if err != nil {
			sn.log.Error("Unable to set fee",
				zap.Error(err),
				zap.String("nid", connection.Nid),
			)
		}
		time.Sleep(10 * time.Second)
	}
	time.Sleep(30 * time.Second)
	return nil
}

func (sn *SolanaRemoteNet) GetContractAddress(key string) string {
	value, exist := sn.IBCAddresses[key]
	if !exist {
		panic(fmt.Sprintf(`IBC address not exist %s`, key))
	}
	return value
}

type NetworkAddress struct {
	Address string
}

func (sn *SolanaRemoteNet) SendPacketXCall(ctx context.Context, keyName, _to string, data, rollback []byte) (context.Context, error) {
	var params []interface{}
	params = append(params, NetworkAddress{Address: _to})
	params = append(params, data)
	networkId := strings.Split(_to, "/")[0]

	var rollbackAc solana.AccountMeta
	if rollback != nil {
		msgType := uint32(1)
		params = append(params, msgType)
		params = append(params, rollback)
		xcallConfigAc := XcallConfigAccount{}
		configAc := sn.pdaAcs["configAcXcall"]
		sn.rpcClient.GetAccountInfo(ctx, configAc.String(), &xcallConfigAc)
		sn.pdaAcs["feeHandler"] = xcallConfigAc.FeeHandler

		rollbackSd := xcallConfigAc.SequenceNo.Add(&xcallConfigAc.SequenceNo, big.NewInt(1))
		seeds := [][]byte{
			[]byte("rollback"),
			[]byte(rollbackSd.FillBytes(make([]byte, 16))),
		}
		rollbackAcc, _, err := solana.FindProgramAddress(seeds,
			solana.MustPublicKeyFromBase58(sn.IBCAddresses["xcall"]))
		if err != nil {
			return ctx, err
		}
		rollbackAc = solana.AccountMeta{
			PublicKey:  rollbackAcc,
			IsWritable: true,
		}
	} else {
		msgType := uint32(2)
		rollbackData := make([]byte, 0)
		params = append(params, msgType)
		params = append(params, rollbackData)
		rollbackAc = solana.AccountMeta{
			PublicKey:  solana.MustPublicKeyFromBase58(sn.IBCAddresses["xcall"]),
			IsWritable: true,
		}
	}
	payerAccount := solana.AccountMeta{
		PublicKey:  sn.walletPrivKey.PublicKey(),
		IsWritable: true,
		IsSigner:   true,
	}
	seeds := [][]byte{
		[]byte("fee"),
		[]byte(networkId),
	}
	networkFeeAc, _, err := solana.FindProgramAddress(seeds,
		solana.MustPublicKeyFromBase58(sn.IBCAddresses["connection"]))
	if err != nil {
		log.Fatalf("Failed to find program address: %v", err)
	}
	seeds = [][]byte{
		[]byte("connections"),
		[]byte(networkId),
	}
	testcase := ctx.Value("testcase").(string)
	dappAccount := sn.IBCAddresses[fmt.Sprintf("dapp-%s", testcase)]
	connAc, _, err := solana.FindProgramAddress(seeds,
		solana.MustPublicKeyFromBase58(dappAccount))
	if err != nil {
		log.Fatalf("Failed to find program address: %v", err)
		return ctx, err
	}
	accountsMeta := solana.AccountMetaSlice{
		{
			PublicKey:  sn.pdaAcs["configAcDapp"],
			IsWritable: true,
		},
		{
			PublicKey: sn.pdaAcs["dappAuthorityAc"],
		},
		{
			PublicKey:  connAc,
			IsWritable: true,
		},
		&payerAccount,
		&solana.AccountMeta{
			PublicKey: solana.SystemProgramID,
		},
		&solana.AccountMeta{
			PublicKey: solana.SysVarInstructionsPubkey,
		},
		{
			PublicKey:  sn.pdaAcs["configAcXcall"],
			IsWritable: true,
		},
		{
			PublicKey:  sn.pdaAcs["feeHandler"],
			IsWritable: true,
		},
		&rollbackAc,
		{
			PublicKey:  solana.MustPublicKeyFromBase58(sn.IBCAddresses["connection"]),
			IsWritable: true,
		},
		{
			PublicKey:  sn.pdaAcs["connconfig"],
			IsWritable: true,
		},
		{
			PublicKey:  networkFeeAc,
			IsWritable: true,
		},
		{
			PublicKey:  solana.MustPublicKeyFromBase58(sn.IBCAddresses["xcall"]),
			IsWritable: true,
		},
		{
			PublicKey:  solana.MustPublicKeyFromBase58(sn.IBCAddresses["connection"]),
			IsWritable: true,
		},
	}

	_, txSign, err := sn.executeContract(ctx, mockAppProgName, "send_call_message", params, accountsMeta)
	if err != nil {
		return nil, err
	}

	txnres, err := sn.rpcClient.GetTransaction(context.Background(), *txSign, &solrpc.GetTransactionOpts{Commitment: solrpc.CommitmentConfirmed})
	if err != nil {
		return nil, fmt.Errorf("failed to get txn %s: %w", txSign.String(), err)
	}

	events, err := sn.parseReturnValueFromLogs(
		txnres.Meta.LogMessages, EventCallMessageSent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse return value: %w", err)
	}
	return context.WithValue(ctx, "sn", events["sn"]), nil
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
	height, err := targetChain.Height(ctx)
	if err != nil {
		return nil, err
	}
	ctx, err = sn.SendPacketXCall(ctx, keyName, to, data, rollback)
	if err != nil {
		return nil, err
	}
	return sn.FindTargetXCallMessage(ctx, targetChain, height, strings.Split(to, "/")[1])
}

func (sn *SolanaRemoteNet) FindCallMessage(ctx context.Context, startHeight uint64, from, to, sno string) (string, string, error) {
	var allEvents []IdlEvent
	allEvents = append(
		allEvents,
		sn.appIds[mockAppProgName].Events...,
	)
	allEvents = append(
		allEvents,
		sn.appIds[xcallProgName].Events...,
	)
	event, err := sn.FindEvent(ctx, startHeight, allEvents, "xcall", "CallMessage", sno)
	if err != nil {
		return "", "", err
	}
	reqId := event.ValueDecoded["reqId"].(string)
	data := event.ValueDecoded["data"].([]byte)
	return reqId, string(data), nil
}

func (sn *SolanaRemoteNet) FindCallResponse(ctx context.Context, startHeight uint64, sno string) (string, error) {
	var allEvents []IdlEvent
	allEvents = append(
		allEvents,
		sn.appIds[mockAppProgName].Events...,
	)
	allEvents = append(
		allEvents,
		sn.appIds[xcallProgName].Events...,
	)
	event, err := sn.FindEvent(ctx, startHeight, allEvents, "xcall", "ResponseMessage", sno)
	if err != nil {
		return "", err
	}
	code := event.ValueDecoded["code"].(string)
	return code, nil
}

func (sn *SolanaRemoteNet) FindEvent(ctx context.Context, startHeight uint64,
	allEvents []IdlEvent, contract, signature, sno string) (*EventResponseEvent, error) {
	timeout := time.After(120 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("failed to find eventLog")
		case <-ticker.C:
			data, err := sn.getEvent(ctx, allEvents, sno, signature, contract)
			if err != nil {
				continue
			}
			return data, nil
		}
	}
}

func (sn *SolanaRemoteNet) getEvent(ctx context.Context,
	allEvents []IdlEvent, sno, signature, contractName string) (*EventResponseEvent, error) {
	progIdl := sn.appIds[contractName]
	progID, _ := progIdl.GetProgramID()
	return sn.rpcClient.GetEvent(ctx, progID,
		allEvents, signature, sno)
}

// executeContract implements chains.Chain
func (sn *SolanaRemoteNet) executeContract(ctx context.Context, contractProgName,
	methodName string, params []interface{}, accountValues solana.AccountMetaSlice) (context.Context, *solana.Signature, error) {
	progIdl := sn.appIds[contractProgName]
	discriminator, err := progIdl.GetInstructionDiscriminator(methodName)
	if err != nil {
		return nil, nil, err
	}
	progID, err := progIdl.GetProgramID()
	if err != nil {
		return nil, nil, err
	}
	instructionData := discriminator
	for _, param := range params {
		serParam, err := borsh.Serialize(param)
		if err != nil {
			return nil, nil, err
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
	txSign, err := sn.executeTx(ctx, sn.walletPrivKey, instructions, signers)
	return ctx, txSign, err
}

func (sn *SolanaRemoteNet) executeTx(ctx context.Context, payerPrivKey solana.PrivateKey, instructions []solana.Instruction, signers []solana.PrivateKey) (*solana.Signature, error) {
	latestBlockHash, err := sn.rpcClient.GetLatestBlockHash(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block hash: %w", err)
	}
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

	// sn.log.Info(tx.ToBase64())
	txSign, err := sn.rpcClient.SendTx(ctx, tx, nil)
	if err != nil {
		sn.log.Error("failed to send tx:", zap.Error(err))
		return nil, err
	}

	_, err = sn.waitForTxConfirmation(10*time.Second, txSign)
	if err != nil {
		sn.log.Error("error waiting for tx confirmation:", zap.Error(err))
	}
	sn.log.Info("send message successful", zap.String("tx-hash", txSign.String()))
	return &txSign, nil
}

// GetBlockByHeight implements chains.Chain
func (sn *SolanaRemoteNet) GetClientBlockByHeight(ctx context.Context, height int64) (string, error) {
	return "", nil
}

// GetLastBlock implements chains.Chain
func (sn *SolanaRemoteNet) GetLastBlock(ctx context.Context) (context.Context, error) {
	res, err := sn.rpcClient.GetLatestBlockHeight(ctx, solrpc.CommitmentFinalized)
	return context.WithValue(ctx, chains.LastBlock{}, res), err
}

// QueryContract implements chains.Chain
func (sn *SolanaRemoteNet) QueryContract(ctx context.Context, contractAddress, methodName string, params map[string]interface{}) (context.Context, error) {
	return ctx, nil
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
	// stdout, _, err := sn.Exec(ctx, sn.TxUpgradeCommand(ctx, contractName, command...), nil)
	// if err != nil {
	stdout, _, err := sn.Exec(ctx, sn.TxCommand(ctx, contractName, command...), nil)
	if err != nil {
		return "", err
	}
	// }
	return getProgramIdFromDeployment(string(stdout))
}

// TxCommand is a helper to retrieve a full command for broadcasting a tx
// with the chain node binary.
func (sn *SolanaRemoteNet) TxCommand(ctx context.Context, contractName string, command ...string) []string {
	command = append([]string{"deploy", "--provider.cluster", sn.testconfig.RPCUri, "--program-name", contractName}, command...)
	return sn.NodeCommand(command...)
}

func (sn *SolanaRemoteNet) syncIdl(ctx context.Context, filepath, programId string) error {
	command := []string{"idl", "init", "--provider.cluster", sn.testconfig.RPCUri, "--filepath", filepath, programId}
	_, _, err := sn.Exec(ctx, sn.NodeCommand(command...), nil)
	return err
}
func (sn *SolanaRemoteNet) syncProgramId(ctx context.Context, contractName string) error {
	//build

	command := []string{"build", "--provider.cluster", sn.testconfig.RPCUri, "--program-name", contractName}
	_, _, err := sn.Exec(ctx, sn.NodeCommand(command...), nil)
	if err != nil {
		return err
	}

	//sync the key
	command = []string{"keys", "sync", "--program-name", contractName}
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
	var allEvents []IdlEvent
	allEvents = append(
		allEvents,
		sn.appIds[xcallProgName].Events...,
	)
	allEvents = append(
		allEvents,
		sn.appIds[connectionProgName].Events...,
	)
	event, err := sn.FindEvent(ctx, startHeight, allEvents, "xcall", "RollbackExecuted", sno)
	if err != nil {
		return "", err
	}
	fsno := event.ValueDecoded["sn"].(string)
	return fsno, nil
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

func (sn *SolanaRemoteNet) createLookupTableAccount(ctx context.Context) (*solana.PublicKey, error) {
	recentSlot, err := sn.rpcClient.GetLatestBlockHeight(ctx, solrpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}

	recentSlot = recentSlot - 1

	altCreateInstruction, accountAddr, err := alt.CreateLookupTable(
		sn.walletPrivKey.PublicKey(),
		sn.walletPrivKey.PublicKey(),
		recentSlot,
	)
	if err != nil {
		return nil, err
	}

	signers := []solana.PrivateKey{sn.walletPrivKey}

	latestBlockHash, err := sn.rpcClient.GetLatestBlockHash(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block hash: %w", err)
	}
	tx, err := solana.NewTransaction(
		[]solana.Instruction{altCreateInstruction},
		*latestBlockHash,
		solana.TransactionPayer(sn.walletPrivKey.PublicKey()))
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

	txSign, err := sn.rpcClient.SendTx(ctx, tx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send tx: %w", err)
	}

	_, err = sn.waitForTxConfirmation(defaultTxConfirmationTime, txSign)
	if err != nil {
		return nil, err
	}

	return &accountAddr, nil
}

func (sn *SolanaRemoteNet) parseReturnValueFromLogs(logs []string, event string) (map[string]interface{}, error) {
	events := make(map[string]interface{})
	for _, log := range logs {
		if strings.HasPrefix(log, EventLogPrefix) {
			eventLog := strings.Replace(log, EventLogPrefix, "", 1)
			eventLogBytes, err := base64.StdEncoding.DecodeString(eventLog)
			if err != nil {
				return nil, err
			}

			if len(eventLogBytes) < 8 {
				return nil, fmt.Errorf("decoded bytes too short to contain discriminator: %v", eventLogBytes)
			}

			discriminator := eventLogBytes[:8]
			eventBytes := eventLogBytes[8:]
			var allEvents []IdlEvent
			allEvents = append(
				allEvents,
				sn.appIds[xcallProgName].Events...,
			)
			allEvents = append(
				allEvents,
				sn.appIds[connectionProgName].Events...,
			)

			for _, ev := range allEvents {
				if slices.Equal(ev.Discriminator, discriminator) {
					if ev.Name == event {
						if event == EventCallMessageSent {
							smEvent := CallMessageSent{}
							if err := borsh.Deserialize(&smEvent, eventBytes); err != nil {
								return nil, fmt.Errorf("failed to decode send message event: %w", err)
							}
							events["sn"] = strconv.Itoa(int(smEvent.Sn.Int64()))
							events["from"] = smEvent.From
							events["to"] = smEvent.To
							return events, nil
						} else if event == EventRollbackExecuted {
							smEvent := RollbackExecuted{}
							if err := borsh.Deserialize(&smEvent, eventBytes); err != nil {
								return nil, fmt.Errorf("failed to decode send message event: %w", err)
							}
							events["sn"] = strconv.Itoa(int(smEvent.Sn.Int64()))
							return events, nil
						} else if event == EventResponseMessage {
							smEvent := ResponseMessage{}
							if err := borsh.Deserialize(&smEvent, eventBytes); err != nil {
								return nil, fmt.Errorf("failed to decode send message event: %w", err)
							}
							events["sn"] = strconv.Itoa(int(smEvent.Sn.Int64()))
							events["code"] = strconv.Itoa(int(smEvent.Code))
							return events, nil
						}
					}
				}
			}
		}

	}
	return events, nil
}
