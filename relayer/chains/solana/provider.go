package solana

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/icon-project/centralized-relay/relayer/chains/solana/types"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

type Provider struct {
	log    *zap.Logger
	cfg    *Config
	client IClient
	wallet *solana.Wallet
	kms    kms.KMS
	txmut  *sync.Mutex
}

func (p *Provider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	return p.client.GetLatestBlockHeight(ctx)
}

func (p *Provider) NID() string {
	return p.cfg.NID
}

func (p *Provider) Name() string {
	return p.cfg.ChainName
}

func (p *Provider) Init(ctx context.Context, homePath string, kms kms.KMS) error {
	p.kms = kms
	return nil
}

// Type returns chain-type
func (p *Provider) Type() string {
	return types.ChainType
}

func (p *Provider) Config() provider.Config {
	return p.cfg
}

// FinalityBlock returns the number of blocks the chain has to advance from current block inorder to
// consider it as final. In Solana blocks once published are final.
// So Solana doesn't need to be checked for block finality.
func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	return 0
}

func (p *Provider) GenerateMessages(ctx context.Context, messageKey *relayertypes.MessageKeyWithMessageHeight) ([]*relayertypes.Message, error) {
	return nil, fmt.Errorf("method not implemented")
}

// SetAdmin transfers the ownership of sui connection module to new address
func (p *Provider) SetAdmin(ctx context.Context, adminAddr string) error {
	return nil
}

func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	return nil
}

func (p *Provider) GetFee(ctx context.Context, networkID string, responseFee bool) (uint64, error) {
	return 0, nil
}

func (p *Provider) SetFee(ctx context.Context, networkID string, msgFee, resFee uint64) error {
	return nil
}

func (p *Provider) ClaimFee(ctx context.Context) error {
	return nil
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*relayertypes.Coin, error) {
	return nil, nil
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, messagekey *relayertypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, messageKey *relayertypes.Message) (bool, error) {
	return true, nil
}
