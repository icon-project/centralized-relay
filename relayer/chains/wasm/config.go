package wasm

import (
	"fmt"
	"github.com/cometbft/cometbft/rpc/client/http"
	sdkClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/client"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"os"
	"path/filepath"
	"time"
)

type ProviderConfig struct {
	ChainName string `json:"-" yaml:"-"`
	RpcUrl    string `json:"rpc-url" yaml:"rpc-url"`
	ChainID   string `json:"chain_id" yaml:"chain-id"`
	NID       string `json:"nid" yaml:"nid"`

	HomeDir string `json:"home_dir" yaml:"home-dir"`

	KeyringBackend string                `json:"keyring_backend" yaml:"keyring-backend"`
	KeyringDir     string                `json:"-" yaml:"-"`
	KeyName        string                `json:"key_name" yaml:"key-name"`
	KeyPassword    types.KeyringPassword `json:"key_password" yaml:"key-password"`

	OutputFormat string `json:"output_format" yaml:"output-format"`
	Output       io.Writer
	Input        io.Reader

	AccountPrefix string `json:"account-prefix" yaml:"account-prefix"`

	ContractAddress string `json:"contract-address" yaml:"contract-address"`

	Denomination string `json:"denomination" yaml:"denomination"`

	GasPrices     string  `json:"gas-prices" yaml:"gas-prices"`
	GasAdjustment float64 `json:"gas-adjustment" yaml:"gas-adjustment"`
	MinGasAmount  uint64  `json:"min-gas-amount" yaml:"min-gas-amount"`
	MaxGasAmount  uint64  `json:"max-gas-amount" yaml:"max-gas-amount"`

	BlockInterval string `json:"block_interval" yaml:"block-interval"`

	BroadcastMode string `json:"broadcast_mode" yaml:"broadcast-mode"` //sync or async
	SignModeStr   string `json:"sign-mode" yaml:"sign-mode"`

	Simulate bool `json:"simulate" yaml:"simulate"`

	Debug bool `json:"-" yaml:"-"`
}

func (pc ProviderConfig) NewProvider(logger *zap.Logger, homePath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if chainName != "" {
		pc.ChainName = chainName
	}
	if pc.HomeDir == "" {
		pc.HomeDir = homePath
	}

	pc.Debug = debug

	pc.Input = os.Stdin
	pc.Output = os.Stdout

	pc.KeyringDir = filepath.Join(pc.HomeDir, ".config", pc.ChainName, "keys")

	if err := pc.Validate(); err != nil {
		return nil, err
	}

	clientContext, err := newClientContext(pc)
	if err != nil {
		return nil, err
	}

	cp := &Provider{
		logger: logger,
		config: pc,
		client: client.New(clientContext),
	}

	return cp, nil
}

func (pc ProviderConfig) Validate() error {
	if _, err := time.ParseDuration(pc.BlockInterval); err != nil {
		return fmt.Errorf("invalid block-interval: %w", err)
	}

	if pc.ChainName == "" {
		return fmt.Errorf("chain-name cannot be empty")
	}

	if pc.HomeDir == "" {
		return fmt.Errorf("home-dir cannot be empty")
	}
	return nil
}

func newClientContext(pc ProviderConfig) (sdkClient.Context, error) {
	clientContext := sdkClient.Context{}

	ifr := getInterfaceRegistry()
	protoCodec := codec.NewProtoCodec(ifr)

	keyRing, err := keyring.New(pc.ChainName, pc.KeyringBackend, pc.KeyringDir, pc.KeyPassword, protoCodec)
	if err != nil {
		return clientContext, err
	}

	keyRecord, err := keyRing.Key(pc.KeyName)
	if err != nil {
		return clientContext, err
	}
	fromAddress, err := keyRecord.GetAddress()
	if err != nil {
		return clientContext, err
	}

	cometRPCClient, err := http.New(pc.RpcUrl, "/websocket")
	if err != nil {
		return clientContext, err
	}

	grpcClient, err := grpc.Dial(pc.RpcUrl) //Todo use secured rpc channel for production
	if err != nil {
		return clientContext, err
	}

	clientContext.
		WithInterfaceRegistry(ifr).
		WithCodec(protoCodec).
		WithFromAddress(fromAddress).
		WithFrom(keyRecord.Name).
		WithFromName(keyRecord.Name).
		WithKeyring(keyRing).
		WithKeyringDir(pc.KeyringDir).
		WithNodeURI(pc.RpcUrl).
		WithChainID(pc.ChainID).
		WithHomeDir(pc.HomeDir).
		WithBroadcastMode(pc.BroadcastMode).
		WithSignModeStr(pc.SignModeStr).
		WithSimulation(pc.Simulate).
		WithFeePayerAddress(fromAddress).
		WithFeeGranterAddress(fromAddress).
		WithClient(cometRPCClient).
		WithGRPCClient(grpcClient).
		WithOutputFormat(pc.OutputFormat).
		WithOutput(pc.Output).
		WithInput(pc.Input)

	return clientContext, nil
}
