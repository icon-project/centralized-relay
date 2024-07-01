package solana

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/icon-project/centralized-relay/relayer/chains/solana/types"
	relayerevents "github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/near/borsh-go"
)

func (p *Provider) Route(ctx context.Context, message *relayertypes.Message, callback relayertypes.TxResponseFunc) error {
	instructions, signers, err := p.MakeCallInstructions(message)
	if err != nil {
		return fmt.Errorf("failed to create call instructions: %w", err)
	}

	tx, err := p.prepareAndSimulateTx(ctx, instructions, signers)
	if err != nil {
		return fmt.Errorf("failed to prepare and simulate tx: %w", err)
	}

	txSign, err := p.client.SendTx(ctx, tx, nil)
	if err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	go p.executeRouteCallback(txSign, message, callback)

	return nil
}

func (p *Provider) prepareAndSimulateTx(ctx context.Context, instructions []solana.Instruction, signers []solana.PrivateKey) (*solana.Transaction, error) {
	latestBlockHash, err := p.client.GetLatestBlockHash(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block hash: %w", err)
	}

	tx, err := solana.NewTransaction(instructions, *latestBlockHash, solana.TransactionPayer(p.wallet.PublicKey()))
	if err != nil {
		return nil, fmt.Errorf("failed to create new tx: %w", err)
	}

	// simres, err := p.client.SimulateTx(ctx, tx, nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to simulate tx: %w", err)
	// }

	// if p.cfg.GasLimit != 0 && p.cfg.GasLimit < *simres.UnitsConsumed {
	// 	return nil, fmt.Errorf("budget requirement is too high: %d greater than allowed limit: %d", *simres.UnitsConsumed, p.cfg.GasLimit)
	// }

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			for _, signer := range signers {
				if signer.PublicKey() == key {
					return &signer
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign tx: %w", err)
	}

	return tx, nil
}

func (p *Provider) waitForTxConfirmation(timeout time.Duration, sign solana.Signature) (*solrpc.SignatureStatusesResult, error) {
	startTime := time.Now()
	for range time.NewTicker(500 * time.Millisecond).C {
		txStatus, err := p.client.GetSignatureStatus(context.TODO(), false, sign)
		if err == nil && txStatus != nil && (txStatus.ConfirmationStatus == solrpc.ConfirmationStatusConfirmed || txStatus.ConfirmationStatus == solrpc.ConfirmationStatusFinalized) {
			return txStatus, nil
		} else if time.Since(startTime) > timeout {
			var cbErr error
			if err != nil {
				cbErr = err
			} else if txStatus != nil && txStatus.Err != nil {
				cbErr = fmt.Errorf("failed to get tx signature status: %v", txStatus.Err)
			} else {
				cbErr = fmt.Errorf("failed to finalize tx signature")
			}
			return nil, cbErr
		}
	}
	return nil, fmt.Errorf("request timeout")
}

func (p *Provider) executeRouteCallback(
	sign solana.Signature,
	msg *relayertypes.Message,
	callback relayertypes.TxResponseFunc,
) {
	txResult, err := p.waitForTxConfirmation(3*time.Second, sign)
	if err != nil {
		callback(
			msg.MessageKey(),
			&relayertypes.TxResponse{
				TxHash: sign.String(),
			},
			err,
		)
	} else {
		callback(
			msg.MessageKey(),
			&relayertypes.TxResponse{
				Height: int64(txResult.Slot),
				TxHash: sign.String(),
			},
			nil,
		)
	}
}

func (p *Provider) MakeCallInstructions(msg *relayertypes.Message) ([]solana.Instruction, []solana.PrivateKey, error) {
	switch msg.EventType {
	case relayerevents.EmitMessage:
		return p.getRecvMessageIntruction(msg)
	default:
		return nil, nil, fmt.Errorf("invalid event type in message")
	}
}

func (p *Provider) getRecvMessageIntruction(msg *relayertypes.Message) ([]solana.Instruction, []solana.PrivateKey, error) {
	discriminator, err := p.xcallIdl.GetInstructionDiscriminator(types.MethodRecvMessage)
	if err != nil {
		return nil, nil, err
	}

	srcArg, err := borsh.Serialize(msg.Src)
	if err != nil {
		return nil, nil, err
	}
	snArg, err := borsh.Serialize(msg.Sn)
	if err != nil {
		return nil, nil, err
	}
	dataArg, err := borsh.Serialize(msg.Data)
	if err != nil {
		return nil, nil, err
	}

	progID, err := p.xcallIdl.GetProgramID()
	if err != nil {
		return nil, nil, err
	}

	payerAccount := solana.AccountMeta{
		PublicKey:  p.wallet.PublicKey(),
		IsWritable: true,
		IsSigner:   true,
	}

	xcallStatePubKey, err := solana.PublicKeyFromBase58(p.cfg.XcallStateAccount)
	if err != nil {
		return nil, nil, err
	}

	xcallStateAccount := solana.AccountMeta{
		PublicKey:  xcallStatePubKey,
		IsWritable: true,
	}

	instructionData := append(discriminator, srcArg...)
	instructionData = append(instructionData, snArg...)
	instructionData = append(instructionData, dataArg...)

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID:        progID,
			AccountValues: solana.AccountMetaSlice{&payerAccount, &xcallStateAccount},
			DataBytes:     instructionData,
		},
	}

	return instructions, []solana.PrivateKey{p.wallet.PrivateKey}, nil
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txDigest string) (*relayertypes.Receipt, error) {
	return nil, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayertypes.MessageKey) (bool, error) {
	return false, nil
}
