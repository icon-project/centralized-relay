package sui

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"sync"

	"github.com/coming-chat/go-sui/v2/account"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/icon-project/centralized-relay/relayer/chains/sui/types"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

var (
	MethodClaimFee      = "claim_fees"
	MethodGetReceipt    = "get_receipt"
	MethodSetFee        = "set_fee"
	MethodGetFee        = "get_fee"
	MethodRevertMessage = "revert_message"
	MethodSetAdmin      = "set_admin"
	MethodGetAdmin      = "get_admin"
	MethodRecvMessage   = "receive_message"
	MethodExecuteCall   = "execute_call"

	ConnectionModule = "centralized_connection"
	EntryModule      = "centralized_entry"
	XcallModule      = "xcall"
	DappModule       = "mock_dapp"

	suiCurrencyDenom = "SUI"
	suiBaseFee       = 1000
)

type Provider struct {
	log    *zap.Logger
	cfg    *Config
	client IClient
	wallet *account.Account
	kms    kms.KMS
	txmut  *sync.Mutex
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
	return nil, fmt.Errorf("method not implemented")
}

// SetAdmin transfers the ownership of sui connection module to new address
func (p *Provider) SetAdmin(ctx context.Context, adminAddr string) error {
	suiMessage := p.NewSuiMessage([]SuiCallArg{
		{Type: CallArgObject, Val: p.cfg.XcallStorageID},
		{Type: CallArgPure, Val: adminAddr},
	}, p.cfg.XcallPkgID, EntryModule, MethodSetAdmin)

	txBytes, err := p.prepareTxMoveCall(suiMessage)
	if err != nil {
		return err
	}
	res, err := p.SendTransaction(ctx, txBytes)
	if err != nil {
		return err
	}
	p.log.Info("set fee txn successful",
		zap.String("tx-hash", res.Digest.String()),
	)
	return nil
}

func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	suiMessage := p.NewSuiMessage([]SuiCallArg{
		{Type: CallArgPure, Val: sn},
	}, p.cfg.XcallPkgID, EntryModule, MethodRevertMessage)
	txBytes, err := p.prepareTxMoveCall(suiMessage)
	if err != nil {
		return err
	}
	res, err := p.SendTransaction(ctx, txBytes)
	if err != nil {
		return err
	}
	p.log.Info("revert message txn successful",
		zap.String("tx-hash", res.Digest.String()),
	)
	return nil
}

func (p *Provider) GetAdmin(ctx context.Context, networkID string, responseFee bool) (uint64, error) {
	suiMessage := p.NewSuiMessage([]SuiCallArg{
		{Type: CallArgObject, Val: p.cfg.XcallStorageID},
	}, p.cfg.XcallPkgID, EntryModule, "get_admin")
	var adminAddr move_types.AccountAddress
	wallet, err := p.Wallet()
	if err != nil {
		return 0, err
	}
	txBytes, err := p.preparePTB(suiMessage)
	if err != nil {
		return 0, err
	}
	if err := p.client.QueryContract(ctx, wallet.Address, txBytes, &adminAddr); err != nil {
		return 0, err
	}

	return 0, nil
}

func (p *Provider) GetFee(ctx context.Context, networkID string, responseFee bool) (uint64, error) {
	suiMessage := p.NewSuiMessage([]SuiCallArg{
		{Type: CallArgObject, Val: p.cfg.XcallStorageID},
		{Type: CallArgPure, Val: networkID},
		{Type: CallArgPure, Val: responseFee},
	}, p.cfg.XcallPkgID, EntryModule, MethodGetFee)
	var fee uint64
	wallet, err := p.Wallet()
	if err != nil {
		return fee, err
	}
	txBytes, err := p.preparePTB(suiMessage)
	if err != nil {
		return fee, err
	}
	if err := p.client.QueryContract(ctx, wallet.Address, txBytes, &fee); err != nil {
		return fee, err
	}
	return fee, nil
}

func (p *Provider) SetFee(ctx context.Context, networkID string, msgFee, resFee uint64) error {
	suiMessage := p.NewSuiMessage([]SuiCallArg{
		{Type: CallArgObject, Val: p.cfg.XcallStorageID},
		{Type: CallArgPure, Val: networkID},
		{Type: CallArgPure, Val: strconv.Itoa(int(msgFee))},
		{Type: CallArgPure, Val: strconv.Itoa(int(resFee))},
	}, p.cfg.XcallPkgID, EntryModule, MethodSetFee)
	txBytes, err := p.prepareTxMoveCall(suiMessage)
	if err != nil {
		return err
	}
	res, err := p.SendTransaction(ctx, txBytes)
	if err != nil {
		return err
	}
	p.log.Info("set fee txn successful",
		zap.String("network-id", networkID),
		zap.String("tx-hash", res.Digest.String()),
	)
	return nil
}

func (p *Provider) ClaimFee(ctx context.Context) error {
	suiMessage := p.NewSuiMessage([]SuiCallArg{
		{Type: CallArgObject, Val: p.cfg.XcallStorageID},
	},
		p.cfg.XcallPkgID, EntryModule, MethodClaimFee)
	txBytes, err := p.prepareTxMoveCall(suiMessage)
	if err != nil {
		return err
	}
	res, err := p.SendTransaction(ctx, txBytes)
	if err != nil {
		return err
	}
	p.log.Info("claim fee txn successful",
		zap.String("tx-hash", res.Digest.String()),
	)
	return nil
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*relayertypes.Coin, error) {
	balance, err := p.client.GetTotalBalance(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &relayertypes.Coin{Amount: balance, Denom: suiCurrencyDenom}, nil
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, messagekey *relayertypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, messageKey *relayertypes.Message) (bool, error) {
	return true, nil
}
