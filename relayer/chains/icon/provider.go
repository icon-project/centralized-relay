package icon

import (
	"context"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/provider"
	goloopclient "github.com/icon-project/goloop/client"
	"go.uber.org/zap"
)

/*
 * The provider assumes the key is in
 */
type IconProviderConfig struct {
	ChainID  string `json:"chain-id" yaml:"chain-id"`
	KeyStore string `json:"key-store" yaml:"key-store"`
	RPCAddr  string `json:"rpc-addr" yaml:"rpc-addr"`
	Password string `json:"password" yaml:"password"`
}

// NewProvider should provide a new Icon provider
func (pp *IconProviderConfig) NewProvider(log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {

	if err := pp.Validate(); err != nil {
		return nil, err
	}

	return &IconProvider{
		log:    log.With(zap.String("chain_id", pp.ChainID)),
		client: goloopclient.NewClientV3(pp.RPCAddr),
	}, nil

}

func (pp *IconProviderConfig) Validate() error {
	if pp.RPCAddr == "" {
		return fmt.Errorf("icon provider rpc endpoint is empty")
	}

	// check keystore and password also matches

	return nil
}

type IconProvider struct {
	log    *zap.Logger
	PCfg   *IconProviderConfig
	client *goloopclient.ClientV3
}

func (icp *IconProvider) ChainId() string {
	return icp.PCfg.ChainID
}
func (icp *IconProvider) Init(ctx context.Context) error {
	return nil
}
