// Package rly provides an interface to the cosmos relayer running in a Docker container.
package centralized

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/interchaintest/relayer"

	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

const (
	RlyDefaultUidGid = "100:1000"
)

// ICONRelayer is the ibc.Relayer implementation for github.com/cosmos/relayer.
type ICONRelayer struct {
	// Embedded DockerRelayer so commands just work.
	*relayer.DockerRelayer
}

func NewCentralizedRelayer(log *zap.Logger, testName string, cli *client.Client, networkID string, options ...relayer.RelayerOption) *ICONRelayer {
	c := commander{log: log}
	for _, opt := range options {
		switch o := opt.(type) {
		case relayer.RelayerOptionExtraStartFlags:
			c.extraStartFlags = o.Flags
		}
	}
	dr, err := relayer.NewDockerRelayer(context.TODO(), log, testName, cli, networkID, c, options...)
	if err != nil {
		panic(err) // TODO: return
	}
	return &ICONRelayer{DockerRelayer: dr}
}

type ICONRelayerChainConfigValue struct {
	NID           string            `yaml:"nid"`
	RPCURL        string            `yaml:"rpc-url"`
	StartHeight   int               `yaml:"start-height"`
	NetworkID     int               `yaml:"network-id"`
	Contracts     map[string]string `yaml:"contracts"`
	BlockInterval string            `yaml:"block-interval"`
	Address       string            `yaml:"address"`
	FinalityBlock uint64            `yaml:"finality-block"`
	StepMin       int64             `yaml:"step-min"`
	StepLimit     int64             `yaml:"step-limit"`
}

type SUIRelayerChainConfigValue struct {
	NID             string    `yaml:"nid"`
	RPCURL          string    `yaml:"rpc-url"`
	WebsocketUrl    string    `yaml:"ws-url"`
	StartHeight     int       `yaml:"start-height"`
	XcallPkgId      string    `yaml:"xcall-package-id"`
	ConnectionId    string    `yaml:"connection-id"`
	ConnectionCapId string    `yaml:"connection-cap-id"`
	XcallStorageId  string    `yaml:"xcall-storage-id"`
	NetworkID       int       `yaml:"network-id"`
	BlockInterval   string    `yaml:"block-interval"`
	Address         string    `yaml:"address"`
	FinalityBlock   uint64    `yaml:"finality-block"`
	GasPrice        int64     `yaml:"gas-price"`
	GasLimit        int       `yaml:"gas-limit"`
	Dapps           []SuiDapp `yaml:"dapps"`
}

type SuiDappModule struct {
	Name     string `yaml:"name" json:"name"`
	CapID    string `yaml:"cap-id" json:"cap-id"`
	ConfigID string `yaml:"config-id" json:"config-id"`
}

type SuiDapp struct {
	PkgID string `json:"package-id" yaml:"package-id"`

	Modules []SuiDappModule `json:"modules" yaml:"modules"`
}

type StellarRelayerChainConfigValue struct {
	NID               string            `yaml:"nid"`
	SorobanUrl        string            `yaml:"soroban-url"`
	HorizonUrl        string            `yaml:"horizon-url"`
	StartHeight       int               `yaml:"start-height"`
	NetworkID         int               `yaml:"network-id"`
	Contracts         map[string]string `yaml:"contracts"`
	BlockInterval     string            `yaml:"block-interval"`
	Address           string            `yaml:"address"`
	FinalityBlock     uint64            `yaml:"finality-block"`
	MaxInclusionFee   int64             `yaml:"max-inclusion-fee"`
	NetworkPassphrase string            `yaml:"network-passphrase"`
}

type EVMRelayerChainConfigValue struct {
	NID           string            `yaml:"nid"`
	RPCURL        string            `yaml:"rpc-url"`
	WebsocketUrl  string            `yaml:"websocket-url"`
	StartHeight   int               `yaml:"start-height"`
	GasPrice      int64             `yaml:"gas-price"`
	GasLimit      int               `yaml:"gas-limit"`
	Contracts     map[string]string `yaml:"contracts"`
	BlockInterval string            `yaml:"block-interval"`
	Address       string            `yaml:"address"`
	FinalityBlock uint64            `yaml:"finality-block"`
}

type CosmosRelayerChainConfigValue struct {
	NID                    string            `yaml:"nid"`
	RPCURL                 string            `yaml:"rpc-url"`
	GrpcUrl                string            `yaml:"grpc-url"`
	StartHeight            int               `yaml:"start-height"`
	GasPrice               string            `yaml:"gas-price"`
	GasLimit               int               `yaml:"gas-limit"`
	Contracts              map[string]string `yaml:"contracts"`
	BlockInterval          string            `yaml:"block-interval"`
	Address                string            `yaml:"address"`
	KeyringBackend         string            `yaml:"keyring-backend"`
	TxConfirmationInterval string            `yaml:"tx-confirmation-interval"`
	ChainName              string            `yaml:"chain-name"`
	MinGasAmount           uint64            `yaml:"min-gas-amount"`
	AccountPrefix          string            `yaml:"account-prefix"`
	Denomination           string            `yaml:"denomination"`
	ChainID                string            `yaml:"chain-id"`
	BroadcastMode          string            `yaml:"broadcast-mode"` // sync, async and block. Recommended: sync
	SignModeStr            string            `yaml:"sign-mode"`
	MaxGasAmount           uint64            `yaml:"max-gas-amount"`
	Simulate               bool              `yaml:"simulate"`
	GasAdjustment          float64           `yaml:"gas-adjustment"`
	FinalityBlock          uint64            `yaml:"finality-block"`
}
type Dapp struct {
	Name      string `yaml:"name"`
	ProgramID string `yaml:"program-id"`
}
type SolanaRelayerChainConfigValue struct {
	Disabled          bool     `yaml:"disabled" json:"disabled"`
	ChainName         string   `yaml:"-"`
	RPCUrl            string   `yaml:"rpc-url"`
	Address           string   `yaml:"address"`
	XcallProgram      string   `yaml:"xcall-program"`
	ConnectionProgram string   `yaml:"connection-program"`
	OtherConnections  []string `yaml:"other-connections"`
	Dapps             []Dapp   `yaml:"dapps"`
	CpNIDs            []string `yaml:"cp-nids"`     //counter party NIDs Eg: ["0x2.icon", "0x7.icon"]
	AltAddress        string   `yaml:"alt-address"` // address lookup table address
	NID               string   `yaml:"nid"`
	HomeDir           string   `yaml:"home-dir"`
	GasLimit          uint64   `yaml:"gas-limit"`
	StartHeight       uint64   `yaml:"start-height"`
}

type StacksRelayerChainConfigValue struct {
	NID           string            `yaml:"nid" json:"nid"`
	RPCURL        string            `yaml:"rpc-url"`
	StartHeight   int               `yaml:"start-height"`
	Contracts     map[string]string `yaml:"contracts"`
	BlockInterval string            `yaml:"block-interval"`
	Address       string            `yaml:"address"`
	FinalityBlock uint64            `yaml:"finality-block"`
}

type ICONRelayerChainConfig struct {
	Type  string                      `json:"type"`
	Value ICONRelayerChainConfigValue `json:"value"`
}

type EVMRelayerChainConfig struct {
	Type  string                     `json:"type"`
	Value EVMRelayerChainConfigValue `json:"value"`
}

type CosmosRelayerChainConfig struct {
	Type  string                        `json:"type"`
	Value CosmosRelayerChainConfigValue `json:"value"`
}

type SUIRelayerChainConfig struct {
	Type  string                     `json:"type"`
	Value SUIRelayerChainConfigValue `json:"value"`
}

type SolanaRelayerChainConfig struct {
	Type  string                        `json:"type"`
	Value SolanaRelayerChainConfigValue `json:"value"`
}
type StellarRelayerChainConfig struct {
	Type  string                         `json:"type"`
	Value StellarRelayerChainConfigValue `json:"value"`
}

type StacksRelayerChainConfig struct {
	Type  string                        `yaml:"type"`
	Value StacksRelayerChainConfigValue `yaml:"value"`
}

const (
	DefaultContainerImage   = "centralized-relay"
	DefaultContainerVersion = "latest"
)

// Capabilities returns the set of capabilities of the Cosmos relayer.
//
// Note, this API may change if the rly package eventually needs
// to distinguish between multiple rly versions.
func Capabilities() map[relayer.Capability]bool {
	// RC1 matches the full set of capabilities as of writing.
	return relayer.FullCapabilities()
}

// commander satisfies relayer.RelayerCommander.
type commander struct {
	log             *zap.Logger
	extraStartFlags []string
}

func (commander) Name() string {
	return "centralized-relay"
}

func (commander) DockerUser() string {
	return RlyDefaultUidGid // docker run -it --rm --entrypoint echo ghcr.io/cosmos/relayer "$(id -u):$(id -g)"
}

func (commander) Flush(pathName, channelID, homeDir string) []string {
	cmd := []string{"centralized-relay", "tx", "flush"}
	if pathName != "" {
		cmd = append(cmd, pathName)
		if channelID != "" {
			cmd = append(cmd, channelID)
		}
	}
	cmd = append(cmd, "--home", homeDir)
	return cmd
}

func (commander) RestoreKey(chainID, keyName, coinType, mnemonic, homeDir string) []string {
	return []string{
		"centralized-relay", "keys", "restore", chainID, keyName, mnemonic,
		"--coin-type", fmt.Sprint(coinType),
	}
}

func (c commander) RelayerExecutable() string {
	return "centralized-relay"
}

func (c commander) RelayerCommand(command string, params ...interface{}) []string {
	cmd := []string{
		c.RelayerExecutable(),
	}
	switch command {
	case "stale":
		cmd = append(cmd, "database", "messages", "list")
	}
	return cmd
}

func (c commander) StartRelayer(homeDir string, pathNames ...string) []string {
	cmd := []string{
		"centralized-relay", "start", "--debug", "--flush-interval", "40s",
	}
	cmd = append(cmd, c.extraStartFlags...)
	cmd = append(cmd, pathNames...)
	return cmd
}

func (commander) DefaultContainerImage() string {
	return DefaultContainerImage
}

func (commander) DefaultContainerVersion() string {
	return DefaultContainerVersion
}

func (commander) ParseAddKeyOutput(stdout, stderr string) (ibc.Wallet, error) {
	var wallet WalletModel
	err := json.Unmarshal([]byte(stdout), &wallet)
	rlyWallet := NewWallet("", wallet.Address, wallet.Mnemonic)
	return rlyWallet, err
}

func (commander) ParseRestoreKeyOutput(stdout, stderr string) string {
	return strings.Replace(stdout, "\n", "", 1)
}

func (commander) Init(homeDir string) []string {
	return []string{
		"centralized-relay", "config", "init",
	}
}

func (c commander) CreateWallet(keyName, address, mnemonic string) ibc.Wallet {
	return NewWallet(keyName, address, mnemonic)
}
