package steller

import (
	"context"
	"math/big"
	"strconv"
	"sync"

	"github.com/icon-project/centralized-relay/relayer/chains/steller/types"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stellar/go/keypair"
	"go.uber.org/zap"
)

type Provider struct {
	log    *zap.Logger
	cfg    *Config
	client IClient
	kms    kms.KMS
	wallet *keypair.Full
	txmut  *sync.Mutex
}

func (p *Provider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	latestLedger, err := p.client.GetLatestLedger(ctx)
	if err != nil {
		return 0, err
	}
	return latestLedger.Sequence, nil
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
// consider it as final. In Steller ledgers are analogues to blocks and ledgers once published are
// final. So Steller doesn't need to be checked for block/ledger finality.
func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	return 0
}

func (p *Provider) GenerateMessages(ctx context.Context, messageKey *relayertypes.MessageKeyWithMessageHeight) ([]*relayertypes.Message, error) {
	//Todo
	return nil, nil
}

func (p *Provider) SetAdmin(ctx context.Context, admin string) error {
	//Todo
	return nil
}

func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	//Todo
	return nil
}

func (p *Provider) GetFee(ctx context.Context, networkID string, responseFee bool) (uint64, error) {
	//Todo
	return 0, nil
}

func (p *Provider) SetFee(ctx context.Context, networkID string, msgFee, resFee uint64) error {
	//Todo
	return nil
}

func (p *Provider) ClaimFee(ctx context.Context) error {
	//Todo
	return nil
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*relayertypes.Coin, error) {
	account, err := p.client.AccountDetail(addr)
	if err != nil {
		return nil, err
	}

	var amt uint64
	for _, bal := range account.Balances {
		balance, err := strconv.Atoi(bal.Balance)
		if err != nil {
			return nil, err
		} else {
			amt = uint64(balance)
			break
		}
	}

	return &relayertypes.Coin{
		Denom:  "XLM",
		Amount: amt,
	}, nil
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, messagekey *relayertypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, messageKey *relayertypes.Message) (bool, error) {
	return true, nil
}
