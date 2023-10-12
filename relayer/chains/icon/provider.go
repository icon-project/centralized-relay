package icon

import (
	"context"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/provider"
	"go.uber.org/zap"
)

type IconProviderConfig struct {
	ChainID         string `json:"chain-id" yaml:"chain-id"`
	KeyStore        string `json:"key-store" yaml:"key-store"`
	RPCAddr         string `json:"rpc-addr" yaml:"rpc-addr"`
	Password        string `json:"password" yaml:"password"`
	StartHeight     uint64 `json:"start-height" yaml:"start-height"` //would be of highest priority
	ContractAddress string `json:"contract-address" yaml:"contract-address"`
}

// NewProvider returns new Icon provider
func (pp *IconProviderConfig) NewProvider(log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {

	if err := pp.Validate(); err != nil {
		return nil, err
	}

	return &IconProvider{
		log:    log.With(zap.String("chain_id", pp.ChainID)),
		client: NewClient(pp.RPCAddr, log),
	}, nil

}

func (pp *IconProviderConfig) Validate() error {
	if pp.RPCAddr == "" {
		return fmt.Errorf("icon provider rpc endpoint is empty")
	}

	// TODO: validation for keystore
	// TODO: contractaddress validation
	// TODO: account should have balance

	return nil
}

type IconProvider struct {
	log    *zap.Logger
	PCfg   *IconProviderConfig
	client *Client
}

func (icp *IconProvider) ChainId() string {
	return icp.PCfg.ChainID
}
func (icp *IconProvider) Init(ctx context.Context) error {
	return nil
}
