package stacks

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"go.uber.org/zap"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/interfaces"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

type Config struct {
	provider.CommonConfig `json:",inline" yaml:",inline"`
}

type Provider struct {
	cfg                 *Config
	client              interfaces.IClient
	log                 *zap.Logger
	kms                 kms.KMS
	privateKey          []byte
	contracts           map[string]providerTypes.EventMap
	LastSavedHeightFunc func() uint64
	routerMutex         sync.Mutex
}

func (c *Config) NewProvider(ctx context.Context, log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	client, err := NewClient(c.RPCUrl, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stacks client: %v", err)
	}

	c.ChainName = chainName
	c.HomeDir = homepath

	return &Provider{
		cfg:       c,
		client:    client,
		log:       log.With(zap.Stringp("nid", &c.NID), zap.Stringp("name", &c.ChainName)),
		contracts: c.eventMap(),
	}, nil
}

func (c *Config) Validate() error {
	if c.RPCUrl == "" {
		return fmt.Errorf("stacks provider rpc endpoint is empty")
	}

	if err := c.Contracts.Validate(); err != nil {
		return fmt.Errorf("contracts are not valid: %s", err)
	}

	return nil
}

func (c *Config) Enabled() bool {
	return !c.Disabled
}

func (c *Config) SetWallet(addr string) {
	c.Address = addr
}

func (c *Config) GetWallet() string {
	return c.Address
}

func (p *Provider) NID() string {
	return p.cfg.NID
}

func (p *Provider) Name() string {
	return p.cfg.ChainName
}

func (p *Provider) Type() string {
	return "stacks"
}

func (p *Provider) Config() provider.Config {
	return p.cfg
}

func (p *Provider) Init(ctx context.Context, homepath string, kms kms.KMS) error {
	p.kms = kms
	return nil
}

func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	return p.cfg.FinalityBlock
}

func (p *Provider) SetLastSavedHeightFunc(f func() uint64) {
	p.LastSavedHeightFunc = f
}

func (p *Provider) GetLastSavedBlockHeight() uint64 {
	return p.LastSavedHeightFunc()
}

func (p *Provider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	latestBlock, err := p.client.GetLatestBlock(ctx)
	if err != nil {
		return 0, err
	}
	if latestBlock == nil {
		return 0, fmt.Errorf("no blocks found")
	}
	return uint64(latestBlock.Height), nil
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txHash string) (*providerTypes.Receipt, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p *Provider) Route(ctx context.Context, message *providerTypes.Message, callback providerTypes.TxResponseFunc) error {
	return fmt.Errorf("not implemented")
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, message *providerTypes.Message) (bool, error) {
	return false, fmt.Errorf("not implemented")
}

func (p *Provider) ShouldSendMessage(ctx context.Context, message *providerTypes.Message) (bool, error) {
	return false, fmt.Errorf("not implemented")
}

func (p *Provider) MessageReceived(ctx context.Context, key *providerTypes.MessageKey) (bool, error) {
	return false, fmt.Errorf("not implemented")
}

func (p *Provider) SetAdmin(ctx context.Context, admin string) error {
	return fmt.Errorf("not implemented")
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*providerTypes.Coin, error) {
	balance, err := p.client.GetAccountBalance(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &providerTypes.Coin{
		Amount: balance.Uint64(),
		Denom:  "STX",
	}, nil
}

func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	return fmt.Errorf("not implemented")
}

func (p *Provider) GetFee(ctx context.Context, networkID string, responseFee bool) (uint64, error) {
	return 0, fmt.Errorf("not implemented")
}

func (p *Provider) SetFee(ctx context.Context, networkID string, msgFee, resFee *big.Int) error {
	return fmt.Errorf("not implemented")
}

func (p *Provider) ClaimFee(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

func (p *Provider) GenerateMessages(ctx context.Context, messageKey *providerTypes.MessageKeyWithMessageHeight) ([]*providerTypes.Message, error) {
	return nil, fmt.Errorf("not implemented")
}
