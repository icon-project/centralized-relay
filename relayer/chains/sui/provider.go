package sui

import (
	"context"
	"math/big"

	"github.com/coming-chat/go-sui/v2/account"
	"github.com/icon-project/centralized-relay/relayer/chains/sui/types"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

var (
	MethodClaimFee      = "methodClaimFee"
	MethodGetReceipt    = "methodGetReceipt"
	MethodSetFee        = "methodSetFee"
	MethodGetFee        = "methodGetFee"
	MethodRevertMessage = "methodRevertMessage"
	MethodSetAdmin      = "methodSetAdmin"
	connectionPackageId = "connectionPackageId"
	connectionModule    = "donuts_with_events"
)

type Provider struct {
	log    *zap.Logger
	cfg    *Config
	client IClient
	wallet *account.Account
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

func (p *Provider) Wallet() (*account.Account, error) {
	if p.wallet == nil {
		if err := p.RestoreKeystore(context.Background()); err != nil {
			return nil, err
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
func (p *Provider) SetAdmin(ctx context.Context, admin string) error {
	//Todo
	suiMessage := p.NewSuiMessage([]interface{}{
		admin,
	}, connectionPackageId, connectionModule, MethodSetAdmin)
	_, err := p.SendTransaction(ctx, suiMessage)
	return err
}

func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	suiMessage := p.NewSuiMessage([]interface{}{
		sn,
	}, connectionPackageId, connectionModule, MethodRevertMessage)
	_, err := p.SendTransaction(ctx, suiMessage)
	return err
}

func (p *Provider) GetFee(ctx context.Context, networkID string, responseFee bool) (uint64, error) {
	suiMessage := p.NewSuiMessage([]interface{}{
		networkID,
		responseFee,
	}, p.cfg.Contracts["connection"], connectionModule, MethodGetFee)
	fee, err := p.GetReturnValuesFromCall(ctx, suiMessage)
	if err != nil {
		return 0, err
	}
	return fee.(uint64), nil
}

func (p *Provider) SetFee(ctx context.Context, networkID string, msgFee, resFee uint64) error {
	suiMessage := p.NewSuiMessage([]interface{}{
		networkID,
		msgFee,
		resFee,
	}, p.cfg.Contracts["connection"], connectionModule, MethodSetFee)
	_, err := p.SendTransaction(ctx, suiMessage)
	return err
}

func (p *Provider) ClaimFee(ctx context.Context) error {
	suiMessage := p.NewSuiMessage([]interface{}{},
		p.cfg.Contracts["connection"], connectionModule, MethodClaimFee)
	_, err := p.SendTransaction(ctx, suiMessage)
	return err
}
