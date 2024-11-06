package sui

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"sync"

	"github.com/coming-chat/go-sui/v2/account"
	suitypes "github.com/coming-chat/go-sui/v2/types"
	"github.com/icon-project/centralized-relay/relayer/chains/sui/types"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
)

var (
	MethodClaimFee                 = "claim_fees"
	MethodGetReceipt               = "get_receipt"
	MethodSetFee                   = "set_fee"
	MethodGetFee                   = "get_fee"
	MethodRevertMessage            = "revert_message"
	MethodSetAdmin                 = "set_admin"
	MethodGetAdmin                 = "get_admin"
	MethodRecvMessage              = "receive_message"
	MethodExecuteCall              = "execute_call"
	MethodExecuteRollback          = "execute_rollback"
	MethodGetWithdrawTokentype     = "get_withdraw_token_type"
	MethodGetExecuteCallParams     = "get_execute_params"
	MethodGetExecuteRollbackParams = "get_rollback_params"

	suiCurrencyDenom = "SUI"
	suiBaseFee       = 1000
)

type Provider struct {
	log                 *zap.Logger
	cfg                 *Config
	client              IClient
	wallet              *account.Account
	kms                 kms.KMS
	txmut               *sync.Mutex
	LastSavedHeightFunc func() uint64
}

func (p *Provider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	return p.client.GetLatestCheckpointSeq(ctx)
}

// SetLastSavedBlockHeightFunc sets the function to save the last saved block height
func (p *Provider) SetLastSavedHeightFunc(f func() uint64) {
	p.LastSavedHeightFunc = f
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

func (p *Provider) GenerateMessages(ctx context.Context, fromHeight, toHeight uint64) ([]*relayertypes.Message, error) {
	var messages []*relayertypes.Message
	for h := fromHeight; h <= toHeight; h++ {
		checkpoint, err := p.client.GetCheckpoint(ctx, h)
		if err != nil {
			p.log.Error("failed to fetch checkpoint", zap.Error(err))
			return nil, err
		}

		eventResponse, err := p.client.GetEventsFromTxBlocks(
			ctx,
			checkpoint.Transactions,
			func(stbr *suitypes.SuiTransactionBlockResponse) bool {
				for _, t := range stbr.ObjectChanges {
					if t.Data.Mutated != nil && t.Data.Mutated.ObjectId.String() == p.cfg.XcallStorageID {
						return true
					}
				}
				return false
			},
		)
		if err != nil {
			p.log.Error("failed to query events", zap.Error(err))
			return nil, err
		}

		blockInfoList, err := p.parseMessagesFromEvents(eventResponse)
		if err != nil {
			p.log.Error("failed to parse messages from events", zap.Error(err))
			return nil, err
		}

		for _, bi := range blockInfoList {
			messages = append(messages, bi.Messages...)
		}
	}

	return messages, nil
}

func (p *Provider) FetchTxMessages(ctx context.Context, txHash string) ([]*relayertypes.Message, error) {
	eventResponse, err := p.client.GetEventsFromTxBlocks(
		ctx,
		[]string{txHash},
		func(stbr *suitypes.SuiTransactionBlockResponse) bool {
			for _, t := range stbr.ObjectChanges {
				if t.Data.Mutated != nil && t.Data.Mutated.ObjectId.String() == p.cfg.XcallStorageID {
					return true
				}
			}
			return false
		},
	)
	if err != nil {
		return nil, err
	}

	blockInfoList, err := p.parseMessagesFromEvents(eventResponse)
	if err != nil {
		return nil, err
	}

	var messages []*relayertypes.Message
	for _, bi := range blockInfoList {
		messages = append(messages, bi.Messages...)
	}

	return messages, nil
}

// SetAdmin transfers the ownership of sui connection module to new address
func (p *Provider) SetAdmin(ctx context.Context, adminAddr string) error {
	// implementation not needed in sui
	return fmt.Errorf("set_admin is not implmented in sui contract")
}

func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	// implementation not needed in sui
	return fmt.Errorf("revert_message is not implemented in sui contract")
}

func (p *Provider) GetFee(ctx context.Context, networkID string, responseFee bool) (uint64, error) {
	suiMessage := p.NewSuiMessage(
		[]string{},
		[]SuiCallArg{
			{Type: CallArgObject, Val: p.cfg.XcallStorageID},
			{Type: CallArgPure, Val: p.cfg.ConnectionID},
			{Type: CallArgPure, Val: networkID},
			{Type: CallArgPure, Val: responseFee},
		}, p.cfg.XcallPkgID, p.cfg.ConnectionModule, MethodGetFee)
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

func (p *Provider) SetFee(ctx context.Context, networkID string, msgFee, resFee *big.Int) error {
	suiMessage := p.NewSuiMessage(
		[]string{},
		[]SuiCallArg{
			{Type: CallArgObject, Val: p.cfg.XcallStorageID},
			{Type: CallArgObject, Val: p.cfg.ConnectionCapID},
			{Type: CallArgPure, Val: networkID},
			{Type: CallArgPure, Val: strconv.Itoa(int(msgFee.Int64()))},
			{Type: CallArgPure, Val: strconv.Itoa(int(resFee.Int64()))},
		}, p.cfg.XcallPkgID, p.cfg.ConnectionModule, MethodSetFee)
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
	suiMessage := p.NewSuiMessage(
		[]string{},
		[]SuiCallArg{
			{Type: CallArgObject, Val: p.cfg.XcallStorageID},
			{Type: CallArgObject, Val: p.cfg.ConnectionCapID},
		},
		p.cfg.XcallPkgID, p.cfg.ConnectionModule, MethodClaimFee)
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

// TODO: not sure if this is the right way to sign the message
func (p *Provider) GetSignMessage(message []byte) []byte {
	hash := sha3.New256()
	hash.Write(message)
	return hash.Sum(nil)
}
