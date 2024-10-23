package steller

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"sync"

	"github.com/icon-project/centralized-relay/relayer/chains/steller/types"
	evtypes "github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stellar/go/keypair"
	"go.uber.org/zap"
)

type Provider struct {
	log                 *zap.Logger
	cfg                 *Config
	client              IClient
	kms                 kms.KMS
	wallet              *keypair.Full
	txmut               *sync.Mutex
	LastSavedHeightFunc func() uint64
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
	if err := p.RestoreKeystore(context.Background()); err != nil {
		return nil
	}
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

func (p *Provider) GenerateMessages(ctx context.Context, fromHeight, toHeight uint64) ([]*relayertypes.Message, error) {
	var messages []*relayertypes.Message

	for h := fromHeight; h <= toHeight; h++ {
		msgs, err := p.fetchLedgerMessages(context.Background(), h)
		if err != nil {
			return nil, fmt.Errorf("failed to generate message for height %d: %w", h, err)
		}
		messages = append(messages, msgs...)
	}

	return messages, nil
}

func (p *Provider) FetchTxMessages(ctx context.Context, txHash string) ([]*relayertypes.Message, error) {
	tx, err := p.client.GetTransaction(txHash)
	if err != nil {
		return nil, err
	}

	events, err := p.client.ParseTxnEvents(&tx, p.getEventFilter(0))
	if err != nil {
		return nil, err
	}

	return p.parseMessagesFromEvents(events)
}

func (p *Provider) SetAdmin(ctx context.Context, admin string) error {
	message := &relayertypes.Message{
		EventType: evtypes.SetAdmin,
		Src:       admin,
	}
	callArgs, err := p.newMiscContractCallArgs(*message)
	if err != nil {
		return err
	}
	_, err = p.sendCallTransaction(*callArgs)
	return err
}

func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	message := &relayertypes.Message{
		EventType: evtypes.RevertMessage,
		Sn:        sn,
	}
	callArgs, err := p.newMiscContractCallArgs(*message)
	if err != nil {
		return err
	}
	_, err = p.sendCallTransaction(*callArgs)
	return err
}

func (p *Provider) GetFee(ctx context.Context, networkID string, responseFee bool) (uint64, error) {
	message := &relayertypes.Message{
		EventType: evtypes.GetFee,
		Src:       networkID,
	}
	callArgs, err := p.newMiscContractCallArgs(*message, responseFee)
	if err != nil {
		return 0, err
	}
	var fee types.ScvU64F128
	err = p.queryContract(*callArgs, &fee)
	return uint64(fee), err
}

func (p *Provider) SetFee(ctx context.Context, networkID string, msgFee, resFee *big.Int) error {
	message := &relayertypes.Message{
		EventType: evtypes.SetFee,
		Src:       networkID,
		Sn:        msgFee,
		ReqID:     resFee,
	}
	callArgs, err := p.newMiscContractCallArgs(*message)
	if err != nil {
		return err
	}
	_, err = p.sendCallTransaction(*callArgs)
	return err
}

func (p *Provider) ClaimFee(ctx context.Context) error {

	message := &relayertypes.Message{
		EventType: evtypes.ClaimFee,
	}
	callArgs, err := p.newMiscContractCallArgs(*message)
	if err != nil {
		return err
	}
	_, err = p.sendCallTransaction(*callArgs)
	return err
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

// SetLastSavedBlockHeightFunc sets the function to save the last saved block height
func (p *Provider) SetLastSavedHeightFunc(f func() uint64) {
	p.LastSavedHeightFunc = f
}
