package solana

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/icon-project/centralized-relay/relayer/chains/solana/types"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/near/borsh-go"
	"go.uber.org/zap"
)

type Provider struct {
	log      *zap.Logger
	cfg      *Config
	client   IClient
	wallet   *solana.Wallet
	kms      kms.KMS
	txmut    *sync.Mutex
	xcallIdl *IDL
	connIdl  *IDL

	pdaRegistry *types.PDARegistry
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

// SetAdmin transfers the ownership of solana connection module to new address
func (p *Provider) SetAdmin(ctx context.Context, adminAddr string) error {
	discriminator, err := p.connIdl.GetInstructionDiscriminator(types.MethodSetAdmin)
	if err != nil {
		return err
	}

	newAdmin, err := solana.PublicKeyFromBase58(adminAddr)
	if err != nil {
		return err
	}

	newAdminBytes, err := borsh.Serialize(newAdmin)
	if err != nil {
		return err
	}

	instructionData := append(discriminator, newAdminBytes...)

	payerAccount := solana.AccountMeta{
		PublicKey:  p.wallet.PublicKey(),
		IsWritable: true,
		IsSigner:   true,
	}

	connConfigAddr, err := p.pdaRegistry.ConnConfig.GetAddress()
	if err != nil {
		return err
	}

	progID, err := p.connIdl.GetProgramID()
	if err != nil {
		return err
	}

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID: progID,
			AccountValues: solana.AccountMetaSlice{
				&payerAccount,
				&solana.AccountMeta{
					PublicKey:  connConfigAddr,
					IsWritable: true,
				},
			},
			DataBytes: instructionData,
		},
	}

	signers := []solana.PrivateKey{p.wallet.PrivateKey}

	tx, err := p.prepareAndSimulateTx(ctx, instructions, signers)
	if err != nil {
		return fmt.Errorf("failed to prepare and simulate tx: %w", err)
	}

	txSign, err := p.client.SendTx(ctx, tx, nil)
	if err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	if _, err := p.waitForTxConfirmation(3*time.Second, txSign); err != nil {
		return err
	}

	p.log.Info("set admin successful")

	return nil
}

func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	return nil
}

func (p *Provider) GetFee(ctx context.Context, networkID string, responseFee bool) (uint64, error) {
	return 0, nil
}

func (p *Provider) SetFee(ctx context.Context, networkID string, msgFee, resFee uint64) error {
	discriminator, err := p.connIdl.GetInstructionDiscriminator(types.MethodSetFee)
	if err != nil {
		return err
	}

	netIDBytes, err := borsh.Serialize(networkID)
	if err != nil {
		return err
	}

	msgFeeBytes, err := borsh.Serialize(msgFee)
	if err != nil {
		return err
	}

	resFeeBytes, err := borsh.Serialize(resFee)
	if err != nil {
		return err
	}

	instructionData := append(discriminator, netIDBytes...)
	instructionData = append(instructionData, msgFeeBytes...)
	instructionData = append(instructionData, resFeeBytes...)

	networkFeeAddr, err := p.pdaRegistry.ConnNetworkFee.GetAddress(networkID)
	if err != nil {
		return err
	}

	connConfigAddr, err := p.pdaRegistry.ConnConfig.GetAddress()
	if err != nil {
		return err
	}

	progID, err := p.connIdl.GetProgramID()
	if err != nil {
		return err
	}

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID: progID,
			AccountValues: solana.AccountMetaSlice{
				&solana.AccountMeta{
					PublicKey:  p.wallet.PublicKey(),
					IsWritable: true,
					IsSigner:   true,
				},
				&solana.AccountMeta{
					PublicKey:  connConfigAddr,
					IsWritable: true,
				},
				&solana.AccountMeta{
					PublicKey:  networkFeeAddr,
					IsWritable: true,
				},
				&solana.AccountMeta{
					PublicKey: solana.SystemProgramID,
				},
			},
			DataBytes: instructionData,
		},
	}

	signers := []solana.PrivateKey{p.wallet.PrivateKey}

	tx, err := p.prepareAndSimulateTx(ctx, instructions, signers)
	if err != nil {
		return fmt.Errorf("failed to prepare and simulate tx: %w", err)
	}

	txSign, err := p.client.SendTx(ctx, tx, nil)
	if err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	if _, err := p.waitForTxConfirmation(3*time.Second, txSign); err != nil {
		return err
	}

	p.log.Info("set fee successful")

	return nil
}

func (p *Provider) ClaimFee(ctx context.Context) error {
	discriminator, err := p.connIdl.GetInstructionDiscriminator(types.MethodClaimFees)
	if err != nil {
		return err
	}

	instructionData := discriminator

	claimFeeAddr, err := p.pdaRegistry.ConnClaimFees.GetAddress()
	if err != nil {
		return err
	}

	connConfigAddr, err := p.pdaRegistry.ConnConfig.GetAddress()
	if err != nil {
		return err
	}

	progID, err := p.connIdl.GetProgramID()
	if err != nil {
		return err
	}

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID: progID,
			AccountValues: solana.AccountMetaSlice{
				&solana.AccountMeta{
					PublicKey:  p.wallet.PublicKey(),
					IsWritable: true,
					IsSigner:   true,
				},
				&solana.AccountMeta{
					PublicKey:  connConfigAddr,
					IsWritable: true,
				},
				&solana.AccountMeta{
					PublicKey:  claimFeeAddr,
					IsWritable: true,
				},
			},
			DataBytes: instructionData,
		},
	}

	signers := []solana.PrivateKey{p.wallet.PrivateKey}

	tx, err := p.prepareAndSimulateTx(ctx, instructions, signers)
	if err != nil {
		return fmt.Errorf("failed to prepare and simulate tx: %w", err)
	}

	txSign, err := p.client.SendTx(ctx, tx, nil)
	if err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	if _, err := p.waitForTxConfirmation(3*time.Second, txSign); err != nil {
		return err
	}

	p.log.Info("claim fees successful")

	return nil
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*relayertypes.Coin, error) {
	accAddr, err := solana.PublicKeyFromBase58(addr)
	if err != nil {
		return nil, err
	}

	res, err := p.client.GetBalance(ctx, accAddr)
	if err != nil {
		return nil, err
	}

	return &relayertypes.Coin{
		Denom:  types.SolanaDenom,
		Amount: res.Value,
	}, nil
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, messagekey *relayertypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, messageKey *relayertypes.Message) (bool, error) {
	return true, nil
}
