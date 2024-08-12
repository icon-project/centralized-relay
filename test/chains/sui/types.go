package sui

import (
	suisdkClient "github.com/coming-chat/go-sui/v2/client"
	"github.com/docker/docker/client"
	"github.com/icon-project/centralized-relay/test/chains"
	"github.com/icon-project/centralized-relay/test/testsuite/testconfig"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"go.uber.org/zap"
)

type SuiRemotenet struct {
	cfg          chains.ChainConfig
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
	PackageId  string
	AdminCap   string
	UpgradeCap string
	Storage    string
	Witness    string
	IdCap      string
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

type MoveModule struct {
	Package string `json:"package"`
	Module  string `json:"module"`
}

type MoveEventRequest struct {
	MoveModule MoveModule `json:"MoveModule"`
}

type ObjectResult struct {
	XcallCap struct {
		Fields struct {
			ID struct {
				ID string `json:"id,omitempty"`
			} `json:"id,omitempty"`
		} `json:"fields,omitempty"`
	} `json:"xcall_cap,omitempty"`
}
