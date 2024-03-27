package sui

import (
	"context"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/icon-project/centralized-relay/relayer/chains/sui/types"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

type Provider struct {
	log    *zap.Logger
	cfg    *Config
	client types.IClient
	wallet models.SuiKeyPair
	kms    kms.KMS
}

func (p *Provider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	return p.client.GetLatestCheckpointSeq(ctx)
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

func (p *Provider) Wallet() (models.SuiKeyPair, error) {
	if p.wallet.PrivateKey == nil {
		if err := p.RestoreKeystore(context.Background()); err != nil {
			return models.SuiKeyPair{}, err
		}
	}
	return p.wallet, nil
}

// FinalityBlock returns the number of blocks the chain has to advance from current block inorder to
// consider it as final. In Sui checkpoints are analogues to blocks and checkpoints once published are
// final. So Sui doesn't need to be checked for block/checkpoint finality.
func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	return 0
}

func (p *Provider) GenerateMessages(ctx context.Context, messageKey *relayertypes.MessageKeyWithMessageHeight) ([]*relayertypes.Message, error) {
	//Todo
	return nil, nil
}

// SetAdmin transfers the ownership of sui connection module to new address
func (p *Provider) SetAdmin(context.Context, string) error {
	//Todo
	return nil
}

func (p *Provider) RevertMessage(context.Context, *big.Int) error {
	//Todo
	return nil
}

func (p *Provider) GetFee(context.Context, string, bool) (uint64, error) {
	//Todo
	return 0, nil
}

func (p *Provider) SetFee(context.Context, string, uint64, uint64) error {
	//Todo
	return nil
}

func (p *Provider) ClaimFee(context.Context) error {
	//Todo
	return nil
}
