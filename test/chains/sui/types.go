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

type ResultEvent struct {
	Events map[string][]string `json:"events"`
}
