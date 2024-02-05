package icon

import (
	"context"
	"fmt"
	"math/big"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/goloop/module"
	"go.uber.org/zap"
)

type IconProviderConfig struct {
	ChainName   string                         `json:"-" yaml:"-"`
	RPCUrl      string                         `json:"rpc-url" yaml:"rpc-url"`
	KeyStore    string                         `json:"keystore" yaml:"keystore"`
	StartHeight uint64                         `json:"start-height" yaml:"start-height"` // would be of highest priority
	Contracts   relayerTypes.ContractConfigMap `json:"contracts" yaml:"contracts"`
	NetworkID   uint                           `json:"network-id" yaml:"network-id"`
	NID         string                         `json:"nid" yaml:"nid"`
}

// NewProvider returns new Icon provider
func (c *IconProviderConfig) NewProvider(log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	c.ChainName = chainName

	return &IconProvider{
		log:    log.With(zap.String("nid ", c.NID)),
		client: NewClient(c.RPCUrl, log),
		cfg:    c,
	}, nil
}

func (c *IconProviderConfig) Validate() error {
	if c.RPCUrl == "" {
		return fmt.Errorf("icon provider rpc endpoint is empty")
	}

	if err := c.Contracts.Validate(); err != nil {
		return fmt.Errorf("contracts are not valid: %s", err)
	}

	// TODO: validation for keystore
	// TODO: contractaddress validation
	// TODO: account should have some balance no balance then use another accoutn

	return nil
}

func (p *IconProviderConfig) SetWallet(addr string) {
	p.KeyStore = addr
}

func (p *IconProviderConfig) GetWallet() string {
	return p.KeyStore
}

type IconProvider struct {
	log      *zap.Logger
	cfg      *IconProviderConfig
	wallet   module.Wallet
	client   *Client
	kms      kms.KMS
	homePath string
}

func (p *IconProvider) NID() string {
	return p.cfg.NID
}

func (p *IconProvider) Init(ctx context.Context, homepath string, kms kms.KMS) error {
	p.kms = kms
	p.homePath = homepath
	return nil
}

func (p *IconProvider) Type() string {
	return "icon"
}

func (p *IconProvider) ProviderConfig() provider.ProviderConfig {
	return p.cfg
}

func (p *IconProvider) ChainName() string {
	return p.cfg.ChainName
}

func (p *IconProvider) Wallet() (module.Wallet, error) {
	if p.wallet == nil {
		if err := p.RestoreKeyStore(context.Background(), p.homePath, p.kms); err != nil {
			return nil, err
		}
	}
	return p.wallet, nil
}

func (p *IconProvider) FinalityBlock(ctx context.Context) uint64 {
	return 0
}

// ReverseMessage reverts a message
func (p *IconProvider) RevertMessage(ctx context.Context, sn *big.Int) error {
	params := map[string]interface{}{"sn": sn}
	message := p.NewIconMessage(params, "revertMessage")
	txHash, err := p.SendTransaction(ctx, message)
	if err != nil {
		return err
	}
	_, txr, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: types.NewHexBytes(txHash)})
	if err != nil {
		return err
	}
	if txr.Status != types.NewHexInt(1) {
		return fmt.Errorf("failed: %s", txr.TxHash)
	}
	return nil
}

// ExecuteCall executes a call to the bridge contract
func (p *IconProvider) ExecuteCall(ctx context.Context, reqID *big.Int, data []byte) error {
	params := map[string]interface{}{"reqID": reqID, "data": data}
	message := p.NewIconMessage(params, "executeCall")
	txHash, err := p.SendTransaction(ctx, message)
	if err != nil {
		return err
	}
	_, txr, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: types.NewHexBytes(txHash)})
	if err != nil {
		return err
	}
	if txr.Status != types.NewHexInt(1) {
		return fmt.Errorf("failed: %s", txr.TxHash)
	}
	return nil
}
