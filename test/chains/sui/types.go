package sui

import (
	suisdkClient "github.com/coming-chat/go-sui/v2/client"
	"github.com/docker/docker/client"
	ibcLocal "github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"go.uber.org/zap"
)

type SuiRemotenet struct {
	cfg          ibcLocal.ChainConfig
	filepath     map[string]string
	IBCAddresses map[string]string     `json:"addresses"`
	Wallets      map[string]ibc.Wallet `json:"wallets"`
	log          *zap.Logger
	DockerClient *client.Client
	Network      string
	testconfig   *testconfig.Chain
	testName     string
	client       *suisdkClient.Client
}

func (c *SuiRemotenet) OverrideConfig(key string, value any) {
	if value == nil {
		return
	}
	c.cfg.ConfigFileOverrides[key] = value
}

type MoveTomlConfig struct {
	Package         map[string]string     `toml:"package"`
	Dependencies    map[string]Dependency `toml:"dependencies"`
	Addresses       map[string]string     `toml:"addresses"`
	DevDependencies map[string]Dependency `toml:"dev-dependencies"`
	DevAddresses    map[string]string     `toml:"dev-addresses"`
}

type Dependency struct {
	Git    string `toml:"git,omitempty"`
	Subdir string `toml:"subdir,omitempty"`
	Rev    string `toml:"rev,omitempty"`
	Local  string `toml:"local,omitempty"`
}

type DepoymentInfo struct {
	PackageId string
	AdminCap  string
	Storage   string
	Witness   string
}

type PackageInfo struct {
	Modules      []string `json:"modules"`
	Dependencies []string `json:"dependencies"`
	Digest       []int    `json:"digest"`
}

type FieldFilter struct {
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type MoveEventModule struct {
	Package string `json:"package"`
	Module  string `json:"module"`
}

type MoveEvent struct {
	MoveEventModule MoveEventModule `json:"MoveEventModule"`
}
