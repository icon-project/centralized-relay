package icon

import (
	"context"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/goloop/module"
	"go.uber.org/zap"
)

type IconProviderConfig struct {
	ChainName       string `json:"-" yaml:"-"`
	RPCUrl          string `json:"rpc-url" yaml:"rpc-url"`
	KeyStore        string `json:"keystore" yaml:"keystore"`
	Password        string `json:"password" yaml:"password"`
	StartHeight     uint64 `json:"start-height" yaml:"start-height"` // would be of highest priority
	ContractAddress string `json:"contract-address" yaml:"contract-address"`
	NetworkID       uint   `json:"network-id" yaml:"network-id"`
	NID             string `json:"nid" yaml:"nid"`
}

// NewProvider returns new Icon provider
func (pp *IconProviderConfig) NewProvider(log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := pp.Validate(); err != nil {
		return nil, err
	}

	pp.ChainName = chainName

	return &IconProvider{
		log:    log.With(zap.String("nid ", pp.NID)),
		client: NewClient(pp.RPCUrl, log),
		PCfg:   pp,
	}, nil
}

func (pp *IconProviderConfig) Validate() error {
	if pp.RPCUrl == "" {
		return fmt.Errorf("icon provider rpc endpoint is empty")
	}

	// TODO: validation for keystore
	// TODO: contractaddress validation
	// TODO: account should have some balance no balance then use another accoutn

	return nil
}

type IconProvider struct {
	log    *zap.Logger
	PCfg   *IconProviderConfig
	client *Client
}

func (ip *IconProvider) NID() string {
	return ip.PCfg.NID
}

func (ip *IconProvider) Init(ctx context.Context) error {
	return nil
}

func (p *IconProvider) Type() string {
	return "icon"
}

func (p *IconProvider) ProviderConfig() provider.ProviderConfig {
	return p.PCfg
}

func (p *IconProvider) ChainName() string {
	return p.PCfg.ChainName
}

func (cp *IconProvider) Wallet() (module.Wallet, error) {
	return cp.RestoreIconKeyStore()
}

func (cp *IconProvider) GetWalletAddress() (address string, err error) {
	return getAddrFromKeystore(cp.PCfg.KeyStore)
}

func (icp *IconProvider) FinalityBlock(ctx context.Context) uint64 {
	return 0
}
