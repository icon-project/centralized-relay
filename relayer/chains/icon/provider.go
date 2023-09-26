package icon

import (
	"context"

	"github.com/icon-project/centralized-relay/relayer/provider"
	"go.uber.org/zap"
)

/*
 * The provider assumes the key is in
 */
type IconProviderConfig struct {
	ChainID string `json:"chain-id" yaml:"chain-id"`
}

// NewProvider should provide a new Icon provider
func (pp *IconProviderConfig) NewProvider(log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {

	if err := pp.Validate(); err != nil {
		return nil, err
	}
	return &IconProvider{
		log: log.With(zap.String("chain_id", pp.ChainID)),
	}, nil

}

func (pp *IconProviderConfig) Validate() error {
	return nil
}

type IconProvider struct {
	log  *zap.Logger
	PCfg *IconProviderConfig
}

func (icp *IconProvider) ChainId() string {
	return icp.PCfg.ChainID
}
func (icp *IconProvider) Init(ctx context.Context) error {
	return nil
}
