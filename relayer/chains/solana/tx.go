package solana

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/icon-project/centralized-relay/relayer/chains/solana/types"
	relayerevents "github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/near/borsh-go"
)

func (p *Provider) Route(ctx context.Context, message *relayertypes.Message, callback relayertypes.TxResponseFunc) error {
	if err := p.RestoreKeystore(ctx); err != nil {
		return err
	}

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
	connSnArg, err := borsh.Serialize(msg.Sn)
	if err != nil {
		return nil, nil, err
	}
	dataArg, err := borsh.Serialize(msg.Data)
	if err != nil {
		return nil, nil, err
	}

	msgHash := sha256.Sum256(msg.Data)

	csMessage, err := p.decodeCsMessage(context.Background(), msg.Data)
	if err != nil {
		return nil, nil, err
	}

	sn := new(big.Int)
	if csMessage.MessageType == types.CsMessageRequest {
		sn = &csMessage.Request.SequenceNo
	} else {
		sn = &csMessage.Result.SequenceNo
	}

	snArg, err := borsh.Serialize(sn)
	if err != nil {
		return nil, nil, err
	}

	instructionData := append(discriminator, srcArg...)
	instructionData = append(instructionData, connSnArg...)
	instructionData = append(instructionData, dataArg...)
	instructionData = append(instructionData, snArg...)

	connConfigAddr, err := p.pdaRegistry.ConnConfig.GetAddress()
	if err != nil {
		return nil, nil, err
	}

	connReceiptAddr, err := p.pdaRegistry.ConnReceipt.GetAddress(new(big.Int).SetUint64(msg.Sn).Bytes())
	if err != nil {
		return nil, nil, err
	}

	accounts := solana.AccountMetaSlice{
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
			PublicKey:  connReceiptAddr,
			IsWritable: true,
		},
		&solana.AccountMeta{
			PublicKey: solana.SystemProgramID,
		},
	}

	xcallConfigAddr, err := p.pdaRegistry.XcallConfig.GetAddress()
	if err != nil {
		return nil, nil, err
	}

	xcallConfigAc := types.XcallConfigAccount{}
	if err := p.client.GetAccountInfo(context.Background(), xcallConfigAddr.String(), &xcallConfigAddr); err != nil {
		return nil, nil, err
	}

	nextReqID := xcallConfigAc.LastReqID.Add(&xcallConfigAc.LastReqID, new(big.Int).SetUint64(1))

	xcallProxyReqAddr, err := p.pdaRegistry.XcallProxyRequest.GetAddress(nextReqID.Bytes())
	if err != nil {
		return nil, nil, err
	}

	xcallDefaultConn, err := p.pdaRegistry.XcallDefaultConn.GetAddress([]byte(msg.Src))
	if err != nil {
		return nil, nil, err
	}

	remainingAccounts := solana.AccountMetaSlice{
		&solana.AccountMeta{PublicKey: xcallConfigAddr, IsWritable: true},
		&solana.AccountMeta{PublicKey: xcallProxyReqAddr, IsWritable: true},
		&solana.AccountMeta{PublicKey: xcallDefaultConn, IsWritable: true},
	}

	if csMessage.MessageType == types.CsMessageRequest {
		pendingRequest := &solana.AccountMeta{PublicKey: p.xcallIdl.GetProgramID(), IsWritable: true}
		pendingRequestCreator := &solana.AccountMeta{PublicKey: p.xcallIdl.GetProgramID(), IsWritable: true}
		if len(csMessage.Request.Protocols) > 0 {
			pendingReqAddr, err := p.pdaRegistry.XcallPendingRequest.GetAddress(msgHash[:])
			if err != nil {
				return nil, nil, err
			}
			pendingRequest.PublicKey = pendingReqAddr
			pendingRequestCreator.PublicKey = p.wallet.PublicKey()
		}

		remainingAccounts = append(
			remainingAccounts,
			solana.AccountMetaSlice{
				pendingRequest,
				pendingRequestCreator,
				&solana.AccountMeta{PublicKey: p.xcallIdl.GetProgramID(), IsWritable: true}, //pending response
				&solana.AccountMeta{PublicKey: p.xcallIdl.GetProgramID(), IsWritable: true}, //pending response creator

				&solana.AccountMeta{PublicKey: p.xcallIdl.GetProgramID(), IsWritable: true},
				&solana.AccountMeta{PublicKey: p.xcallIdl.GetProgramID(), IsWritable: true},
				&solana.AccountMeta{PublicKey: p.xcallIdl.GetProgramID(), IsWritable: true},
			}...,
		)

	} else {
		successResAddr, err := p.pdaRegistry.XcallSuccessRes.GetAddress(csMessage.Result.SequenceNo.Bytes())
		if err != nil {
			return nil, nil, err
		}

		rollbackAddr, err := p.pdaRegistry.XcallRollback.GetAddress(csMessage.Result.SequenceNo.Bytes())
		if err != nil {
			return nil, nil, err
		}

		rollbackAc := types.RollbackAccount{}
		if err := p.client.GetAccountInfo(context.Background(), rollbackAddr.String(), &rollbackAc); err != nil {
			return nil, nil, err
		}

		pendingResponse := &solana.AccountMeta{PublicKey: p.xcallIdl.GetProgramID(), IsWritable: true}
		pendingResponseCreator := &solana.AccountMeta{PublicKey: p.xcallIdl.GetProgramID(), IsWritable: true}
		if len(rollbackAc.Rollback.Protocols) > 0 {
			pendingReqAddr, err := p.pdaRegistry.XcallPendingRequest.GetAddress(msgHash[:])
			if err != nil {
				return nil, nil, err
			}
			pendingResponse.PublicKey = pendingReqAddr
			pendingResponseCreator.PublicKey = p.wallet.PublicKey()
		}

		accounts = append(
			accounts,
			solana.AccountMetaSlice{
				&solana.AccountMeta{PublicKey: p.xcallIdl.GetProgramID(), IsWritable: true}, //pending request
				&solana.AccountMeta{PublicKey: p.xcallIdl.GetProgramID(), IsWritable: true}, //pending request creator
				pendingResponse,        //pending response
				pendingResponseCreator, //pending response creator

				&solana.AccountMeta{PublicKey: successResAddr, IsWritable: true},
				&solana.AccountMeta{PublicKey: rollbackAddr, IsWritable: true},
				&solana.AccountMeta{PublicKey: p.wallet.PublicKey(), IsWritable: true},
			}...,
		)

	}

	accounts = append(accounts, remainingAccounts...)
	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID:        p.xcallIdl.GetProgramID(),
			AccountValues: accounts,
			DataBytes:     instructionData,
		},
	}

	return instructions, []solana.PrivateKey{p.wallet.PrivateKey}, nil
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txSign string) (*relayertypes.Receipt, error) {
	txSignature, err := solana.SignatureFromBase58(txSign)
	if err != nil {
		return nil, err
	}

	txn, err := p.client.GetTransaction(ctx, txSignature, nil)
	if err != nil {
		return nil, err
	}

	return &relayertypes.Receipt{
		TxHash: txSign,
		Height: txn.Slot,
		Status: txn.Meta.Err == nil,
	}, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayertypes.MessageKey) (bool, error) {
	receiptAc, err := p.pdaRegistry.ConnReceipt.GetAddress(new(big.Int).SetUint64(key.Sn).Bytes())
	if err != nil {
		return false, err
	}

	receipt := struct{}{}

	if err := p.client.GetAccountInfo(ctx, receiptAc.String(), &receipt); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (p *Provider) decodeCsMessage(ctx context.Context, msg []byte) (*types.CsMessage, error) {
	if err := p.RestoreKeystore(ctx); err != nil {
		return nil, err
	}

	discriminator, err := p.xcallIdl.GetInstructionDiscriminator(types.MethodDecodeCsMessage)
	if err != nil {
		return nil, err
	}

	msgBorshBytes, err := borsh.Serialize(msg)
	if err != nil {
		return nil, err
	}

	instructionData := append(discriminator, msgBorshBytes...)

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID: p.xcallIdl.GetProgramID(),
			AccountValues: solana.AccountMetaSlice{
				&solana.AccountMeta{
					PublicKey:  p.wallet.PublicKey(),
					IsWritable: true,
					IsSigner:   true,
				},
			},
			DataBytes: instructionData,
		},
	}

	signers := []solana.PrivateKey{p.wallet.PrivateKey}

	tx, err := p.prepareAndSimulateTx(ctx, instructions, signers)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare and simulate tx: %w", err)
	}

	txSign, err := p.client.SendTx(ctx, tx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send tx: %w", err)
	}

	if _, err := p.waitForTxConfirmation(3*time.Second, txSign); err != nil {
		return nil, fmt.Errorf("failed to confirm tx %s: %w", txSign.String(), err)
	}

	txnres, err := p.client.GetTransaction(ctx, txSign, &solrpc.GetTransactionOpts{Commitment: solrpc.CommitmentConfirmed})
	if err != nil {
		return nil, fmt.Errorf("failed to get txn %s: %w", txSign.String(), err)
	}

	for _, log := range txnres.Meta.LogMessages {
		if strings.HasPrefix(log, types.ProgramReturnPrefix) {
			returnLog := strings.Replace(log, types.ProgramReturnPrefix, "", 1)
			returnLogBytes, err := base64.StdEncoding.DecodeString(returnLog)
			if err != nil {
				return nil, err
			}

			csMsg := types.CsMessage{}
			if err := borsh.Deserialize(&csMsg, returnLogBytes); err != nil {
				return nil, err
			}

			return &csMsg, nil
		}
	}

	return nil, fmt.Errorf("failed to return value")
}
