package icon

import (
	"context"
	"fmt"
	"math/big"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	relayerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/goloop/module"
	"go.uber.org/zap"
)

type IconProviderConfig struct {
	ChainName     string                         `json:"-" yaml:"-"`
	RPCUrl        string                         `json:"rpc-url" yaml:"rpc-url"`
	KeyStore      string                         `json:"keystore" yaml:"keystore"`
	StartHeight   uint64                         `json:"start-height" yaml:"start-height"` // would be of highest priority
	Contracts     relayerTypes.ContractConfigMap `json:"contracts" yaml:"contracts"`
	NetworkID     uint                           `json:"network-id" yaml:"network-id"`
	FinalityBlock uint64                         `json:"finality-block" yaml:"finality-block"`
	NID           string                         `json:"nid" yaml:"nid"`
}

// NewProvider returns new Icon provider
func (c *IconProviderConfig) NewProvider(ctx context.Context, log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	c.ChainName = chainName

	return &IconProvider{
		log:       log.With(zap.String("nid ", c.NID), zap.String("chain", chainName)),
		client:    NewClient(ctx, c.RPCUrl, log),
		cfg:       c,
		contracts: c.eventMap(),
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
	log       *zap.Logger
	cfg       *IconProviderConfig
	wallet    module.Wallet
	client    *Client
	kms       kms.KMS
	homePath  string
	contracts map[string]providerTypes.EventMap
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
	return p.cfg.FinalityBlock
}

// ReverseMessage reverts a message
func (p *IconProvider) RevertMessage(ctx context.Context, sn *big.Int) error {
	params := map[string]interface{}{"sn": sn}
	message := p.NewIconMessage(types.Address(p.cfg.Contracts[providerTypes.ConnectionContract]), params, MethodRevertMessage)
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

// SetAdmin sets the admin address of the bridge contract
func (p *IconProvider) SetAdmin(ctx context.Context, admin string) error {
	callParam := map[string]interface{}{
		"_relayer": admin,
	}
	message := p.NewIconMessage(types.Address(p.cfg.Contracts[providerTypes.ConnectionContract]), callParam, MethodSetAdmin)

	data, err := p.SendTransaction(ctx, message)
	if err != nil {
		return fmt.Errorf("SetAdmin: %v", err)
	}
	txHash := types.HexBytes(data)
	p.log.Info("SetAdmin: waiting for tx result", zap.ByteString("txHash", data))
	_, txr, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: txHash})
	if err != nil {
		return fmt.Errorf("SetAdmin: WaitForResults: %v", err)
	}
	if txr.Status != types.NewHexInt(1) {
		return fmt.Errorf("SetAdmin: failed to set admin: %s", txr.TxHash)
	}
	return nil
}

// ExecuteCall executes a call to the bridge contract
func (p *IconProvider) ExecuteCall(ctx context.Context, reqID *big.Int, data []byte) ([]byte, error) {
	params := map[string]interface{}{"_reqId": reqID.Int64(), "_data": data}
	message := p.NewIconMessage(types.Address(p.cfg.Contracts[relayerTypes.XcallContract]), params, MethodExecuteCall)
	txHash, err := p.SendTransaction(ctx, message)
	if err != nil {
		return nil, err
	}
	_, txr, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: types.NewHexBytes(txHash)})
	if err != nil {
		return nil, err
	}
	if txr.Status != types.NewHexInt(1) {
		return nil, fmt.Errorf("failed: %s", txr.TxHash)
	}
	return txr.TxHash.Value()
}
