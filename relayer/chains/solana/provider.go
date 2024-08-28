package solana

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
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
	staticAlts  types.AddressTables
}

func (p *Provider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	return p.client.GetLatestBlockHeight(ctx, solrpc.CommitmentFinalized)
}

func (p *Provider) NID() string {
	return p.cfg.NID
}

func (p *Provider) Name() string {
	return p.cfg.ChainName
}

func (p *Provider) Init(ctx context.Context, homePath string, kms kms.KMS) error {
	p.kms = kms
	if err := p.RestoreKeystore(ctx); err != nil {
		p.log.Warn("failed to load wallet", zap.Error(err))
		return nil
	}

	if p.cfg.AltAddress == "" {
		lutAddr, err := p.createLookupTableAccount(ctx)
		if err != nil {
			return err
		}
		p.log.Info("Look table account created successfully: ", zap.String("alt-address", lutAddr.String()))
		p.cfg.AltAddress = lutAddr.String()
	}

	if err := p.initStaticAlts(); err != nil {
		return fmt.Errorf("failed to init static address lookup tables: %w", err)
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
// consider it as final. In Solana blocks once published are final.
// So Solana doesn't need to be checked for block finality.
func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	return 0
}

func (p *Provider) GenerateMessages(ctx context.Context, messageKey *relayertypes.MessageKeyWithMessageHeight) ([]*relayertypes.Message, error) {
	blockRes, err := p.client.GetBlock(ctx, messageKey.Height)
	if err != nil {
		return nil, err
	}

	messages := []*relayertypes.Message{}

	for _, txn := range blockRes.Transactions {
		event := types.SolEvent{
			Slot:      txn.Slot,
			Signature: txn.MustGetTransaction().Signatures[0],
			Logs:      txn.Meta.LogMessages,
		}

		messages, err := p.parseMessagesFromEvent(event)
		if err != nil {
			return nil, fmt.Errorf("failed to parse messages from event [%+v]: %w", event, err)
		}
		for _, msg := range messages {
			p.log.Info("Detected event log: ",
				zap.Uint64("height", msg.MessageHeight),
				zap.String("event-type", msg.EventType),
				zap.Any("sn", msg.Sn),
				zap.Any("req-id", msg.ReqID),
				zap.String("src", msg.Src),
				zap.String("dst", msg.Dst),
				zap.Any("data", hex.EncodeToString(msg.Data)),
			)
			messages = append(messages, msg)
		}
	}
	return messages, nil
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

	connConfigAddr, err := p.pdaRegistry.ConnConfig.GetAddress()
	if err != nil {
		return err
	}

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID: p.connIdl.GetProgramID(),
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
			},
			DataBytes: instructionData,
		},
	}

	signers := []solana.PrivateKey{p.wallet.PrivateKey}

	tx, err := p.prepareTx(ctx, instructions, signers)
	if err != nil {
		return fmt.Errorf("failed to prepare and simulate tx: %w", err)
	}

	txSign, err := p.client.SendTx(ctx, tx, nil)
	if err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	if _, err := p.waitForTxConfirmation(defaultTxConfirmationTime, txSign); err != nil {
		return fmt.Errorf("failed to confirm tx %s: %w", txSign.String(), err)
	}

	p.log.Info("set admin successful", zap.String("tx-sign", txSign.String()))

	return nil
}

func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	discriminator, err := p.connIdl.GetInstructionDiscriminator(types.MethodRevertMessage)
	if err != nil {
		return err
	}

	snBytes, err := borsh.Serialize(*sn)
	if err != nil {
		return err
	}

	instructionData := append(discriminator, snBytes...)

	accounts := solana.AccountMetaSlice{
		&solana.AccountMeta{
			PublicKey:  p.wallet.PublicKey(),
			IsSigner:   true,
			IsWritable: true,
		},
	}

	res, err := p.queryRevertMessageAccounts(sn, 1, accountsQueryMaxLimit)
	if err != nil {
		return fmt.Errorf("failed to query revert message accounts: %w", err)
	}

	for _, ac := range res.Accounts {
		accounts = append(accounts, &ac)
	}

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID:        p.connIdl.GetProgramID(),
			AccountValues: accounts,
			DataBytes:     instructionData,
		},
	}

	signers := []solana.PrivateKey{p.wallet.PrivateKey}

	tx, err := p.prepareTx(ctx, instructions, signers)
	if err != nil {
		return fmt.Errorf("failed to prepare and simulate tx: %w", err)
	}

	txSign, err := p.client.SendTx(ctx, tx, nil)
	if err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	if _, err := p.waitForTxConfirmation(defaultTxConfirmationTime, txSign); err != nil {
		return fmt.Errorf("failed to confirm tx %s: %w", txSign.String(), err)
	}

	p.log.Info("revert message successful", zap.String("tx-sign", txSign.String()))

	return nil
}

func (p *Provider) GetFee(ctx context.Context, networkID string, responseFee bool) (uint64, error) {
	fee := struct {
		MessageFee  uint64
		ResponseFee uint64
		Bump        uint8
	}{}

	networkFeeAc, err := p.pdaRegistry.ConnNetworkFee.GetAddress([]byte(networkID))
	if err != nil {
		return 0, fmt.Errorf("failed to get network account fee address")
	}

	if err := p.client.GetAccountInfo(ctx, networkFeeAc, &fee); err != nil {
		return 0, fmt.Errorf("failed to get account info: %w", err)
	}

	if responseFee {
		return fee.MessageFee + fee.ResponseFee, nil
	}

	return fee.MessageFee, nil
}

func (p *Provider) SetFee(ctx context.Context, networkID string, msgFee, resFee *big.Int) error {
	discriminator, err := p.connIdl.GetInstructionDiscriminator(types.MethodSetFee)
	if err != nil {
		return err
	}

	netIDBytes, err := borsh.Serialize(networkID)
	if err != nil {
		return err
	}

	msgFeeBytes, err := borsh.Serialize(msgFee.Uint64())
	if err != nil {
		return err
	}

	resFeeBytes, err := borsh.Serialize(resFee.Uint64())
	if err != nil {
		return err
	}

	instructionData := append(discriminator, netIDBytes...)
	instructionData = append(instructionData, msgFeeBytes...)
	instructionData = append(instructionData, resFeeBytes...)

	networkFeeAddr, err := p.pdaRegistry.ConnNetworkFee.GetAddress([]byte(networkID))
	if err != nil {
		return err
	}

	connConfigAddr, err := p.pdaRegistry.ConnConfig.GetAddress()
	if err != nil {
		return err
	}

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID: p.connIdl.GetProgramID(),
			AccountValues: solana.AccountMetaSlice{
				&solana.AccountMeta{
					PublicKey:  p.wallet.PublicKey(),
					IsWritable: true,
					IsSigner:   true,
				},
				&solana.AccountMeta{
					PublicKey: solana.SystemProgramID,
				},
				&solana.AccountMeta{
					PublicKey:  networkFeeAddr,
					IsWritable: true,
				},
				&solana.AccountMeta{
					PublicKey:  connConfigAddr,
					IsWritable: true,
				},
			},
			DataBytes: instructionData,
		},
	}

	signers := []solana.PrivateKey{p.wallet.PrivateKey}

	tx, err := p.prepareTx(ctx, instructions, signers)
	if err != nil {
		return fmt.Errorf("failed to prepare and simulate tx: %w", err)
	}

	txSign, err := p.client.SendTx(ctx, tx, nil)
	if err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	if _, err := p.waitForTxConfirmation(defaultTxConfirmationTime, txSign); err != nil {
		return fmt.Errorf("failed to confirm tx %s: %w", txSign.String(), err)
	}

	p.log.Info("set fee successful", zap.String("tx-sign", txSign.String()))

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

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID: p.connIdl.GetProgramID(),
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

	tx, err := p.prepareTx(ctx, instructions, signers)
	if err != nil {
		return fmt.Errorf("failed to prepare and simulate tx: %w", err)
	}

	txSign, err := p.client.SendTx(ctx, tx, nil)
	if err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	if _, err := p.waitForTxConfirmation(defaultTxConfirmationTime, txSign); err != nil {
		return fmt.Errorf("failed to confirm tx %s: %w", txSign.String(), err)
	}

	p.log.Info("claim fees successful", zap.String("tx-sign", txSign.String()))

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

func (p *Provider) SetLastSavedHeightFunc(f func() uint64) {
	//TODO
}

func (p *Provider) queryRevertMessageAccounts(
	sn *big.Int,
	page uint8,
	limit uint8,
) (*types.QueryAccountsResponse, error) {
	discriminator, err := p.xcallIdl.GetInstructionDiscriminator(types.MethodQueryRecvMessageAccounts)
	if err != nil {
		return nil, err
	}

	snBytes, err := borsh.Serialize(*sn)
	if err != nil {
		return nil, err
	}

	pageBytes, err := borsh.Serialize(page)
	if err != nil {
		return nil, err
	}

	limitBytes, err := borsh.Serialize(limit)
	if err != nil {
		return nil, err
	}

	instructionData := discriminator
	instructionData = append(instructionData, snBytes...)
	instructionData = append(instructionData, pageBytes...)
	instructionData = append(instructionData, limitBytes...)

	connConfigAddr, err := p.pdaRegistry.ConnConfig.GetAddress()
	if err != nil {
		return nil, err
	}

	accounts := solana.AccountMetaSlice{
		&solana.AccountMeta{
			PublicKey:  connConfigAddr,
			IsWritable: false,
			IsSigner:   false,
		},
	}

	xcallConfigAddr, err := p.pdaRegistry.XcallConfig.GetAddress()
	if err != nil {
		return nil, err
	}
	accounts = append(accounts, &solana.AccountMeta{
		PublicKey:  xcallConfigAddr,
		IsWritable: false,
		IsSigner:   false,
	})

	rollbackAddr, err := p.pdaRegistry.XcallRollback.GetAddress(sn.FillBytes(make([]byte, 16)))
	if err != nil {
		return nil, err
	}
	accounts = append(accounts, &solana.AccountMeta{
		PublicKey:  rollbackAddr,
		IsWritable: false,
		IsSigner:   false,
	})

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID:        p.xcallIdl.GetProgramID(),
			AccountValues: accounts,
			DataBytes:     instructionData,
		},
	}

	signers := []solana.PrivateKey{p.wallet.PrivateKey}

	tx, err := p.prepareTx(
		context.Background(),
		instructions,
		signers,
		solana.TransactionPayer(p.wallet.PublicKey()),
	)
	if err != nil {
		return nil, err
	}

	txSign, err := p.client.SendTx(context.Background(), tx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send tx: %w", err)
	}

	_, err = p.waitForTxConfirmation(defaultTxConfirmationTime, txSign)
	if err != nil {
		return nil, err
	}

	txnres, err := p.client.GetTransaction(context.Background(), txSign, &solrpc.GetTransactionOpts{Commitment: solrpc.CommitmentConfirmed})
	if err != nil {
		return nil, fmt.Errorf("failed to get txn %s: %w", txSign.String(), err)
	}

	acRes := types.QueryAccountsResponse{}
	if err := parseReturnValueFromLogs(p.xcallIdl.GetProgramID().String(), txnres.Meta.LogMessages, &acRes); err != nil {
		return nil, fmt.Errorf("failed to parse return value: %w", err)
	}

	return &acRes, nil
}
