package stacks

import (
	"context"
	"fmt"
	"math/big"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/interfaces"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/centralized-relay/relayer/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/stacks-go-sdk/pkg/stacks"
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

	var network *stacks.StacksNetwork
	if c.NID == "stacks" {
		network = stacks.NewStacksMainnet()
	} else if c.NID == "stacks_testnet" {
		network = stacks.NewStacksTestnet()
	} else {
		return nil, fmt.Errorf("no network found for nid: %v", c.NID)
	}

	xcallAbiPath := filepath.Join("abi", "xcall-proxy-abi.json")
	client, err := NewClient(log, network, xcallAbiPath)
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

func (p *Provider) SetAdmin(ctx context.Context, newAdmin string) error {
	if newAdmin == "" {
		return fmt.Errorf("new admin address cannot be empty")
	}

	contractAddress := p.cfg.Contracts[providerTypes.XcallContract]
	if contractAddress == "" {
		return fmt.Errorf("xcall contract address is not set")
	}

	currentImplementation, err := p.client.GetCurrentImplementation(ctx, contractAddress)
	if err != nil {
		return fmt.Errorf("failed to get current implementation: %w", err)
	}

	txID, err := p.client.SetAdmin(ctx, contractAddress, newAdmin, currentImplementation, p.cfg.Address, p.privateKey)
	if err != nil {
		return fmt.Errorf("failed to set new admin: %w", err)
	}

	p.log.Info("SetAdmin transaction broadcasted", zap.String("txID", txID))

	receipt, err := p.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return fmt.Errorf("failed to confirm SetAdmin transaction: %w", err)
	}

	if !receipt.Status {
		return fmt.Errorf("SetAdmin transaction failed: %s", txID)
	}

	p.log.Info("SetAdmin transaction confirmed",
		zap.String("txID", txID),
		zap.Uint64("blockHeight", receipt.Height))

	return nil
}

func (p *Provider) waitForTransactionConfirmation(ctx context.Context, txID string) (*types.Receipt, error) {
	timeout := time.After(2 * time.Minute)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			receipt, err := p.QueryTransactionReceipt(ctx, txID)
			if err != nil {
				p.log.Warn("Failed to query transaction receipt", zap.Error(err))
				continue
			}
			if receipt.Status {
				p.log.Info("Transaction confirmed", zap.String("txID", txID))
				return receipt, nil
			}
			p.log.Debug("Transaction not yet confirmed", zap.String("txID", txID))
		case <-timeout:
			return nil, fmt.Errorf("transaction confirmation timed out after 2 minutes")
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
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
	p.log.Debug("Getting fee", zap.String("networkID", networkID), zap.Bool("responseFee", responseFee))

	fee, err := p.client.GetFee(
		ctx,
		p.cfg.Contracts[providerTypes.ConnectionContract],
		networkID,
		responseFee,
	)
	if err != nil {
		p.log.Error("Failed to get fee", zap.Error(err))
		return 0, fmt.Errorf("failed to get fee: %w", err)
	}

	p.log.Debug("Fee retrieved successfully", zap.Uint64("fee", fee))
	return fee, nil
}

func (p *Provider) SetFee(ctx context.Context, networkID string, msgFee, resFee *big.Int) error {
	p.log.Debug("Setting fees", zap.String("networkID", networkID), zap.String("messageFee", msgFee.String()), zap.String("responseFee", resFee.String()))

	txID, err := p.client.SetFee(
		ctx,
		p.cfg.Contracts[providerTypes.ConnectionContract],
		networkID,
		msgFee,
		resFee,
		p.cfg.Address,
		p.privateKey,
	)
	if err != nil {
		p.log.Error("Failed to set fees", zap.Error(err))
		return fmt.Errorf("failed to set fees: %w", err)
	}

	p.log.Info("Fees set successfully", zap.String("txID", txID))
	return nil
}

func (p *Provider) ClaimFee(ctx context.Context) error {
	p.log.Debug("Claiming fees")

	txID, err := p.client.ClaimFee(
		ctx,
		p.cfg.Contracts[providerTypes.ConnectionContract],
		p.cfg.Address,
		p.privateKey,
	)
	if err != nil {
		p.log.Error("Failed to claim fees", zap.Error(err))
		return fmt.Errorf("failed to claim fees: %w", err)
	}

	p.log.Info("Fees claimed successfully", zap.String("txID", txID))
	return nil
}
