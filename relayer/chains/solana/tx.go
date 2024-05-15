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
)

func (p *Provider) Route(ctx context.Context, message *relayertypes.Message, callback relayertypes.TxResponseFunc) error {
	instructions, err := p.MakeCallInstructions(message)
	if err != nil {
		return fmt.Errorf("failed to create call instructions: %w", err)
	}

	latestBlockHash, err := p.client.GetLatestBlockHash(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block hash: %w", err)
	}

	tx, err := solana.NewTransaction(instructions, *latestBlockHash)
	if err != nil {
		return fmt.Errorf("failed to create new tx: %w", err)
	}

	simres, err := p.client.SimulateTx(ctx, tx, nil)
	if err != nil {
		return fmt.Errorf("failed to simulate tx: %w", err)
	}

	if p.cfg.GasLimit != 0 && p.cfg.GasLimit < *simres.UnitsConsumed {
		return fmt.Errorf("budget requirement is too high: %d greater than allowed limit: %d", *simres.UnitsConsumed, p.cfg.GasLimit)
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if p.wallet.PublicKey().Equals(key) {
				return &p.wallet.PrivateKey
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed to sign tx: %w", err)
	}

	txSign, err := p.client.SendTx(ctx, tx, nil)
	if err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	go p.executeRouteCallback(ctx, txSign, message, callback)

	return nil
}

func (p *Provider) executeRouteCallback(
	ctx context.Context,
	sign solana.Signature,
	msg *relayertypes.Message,
	callback relayertypes.TxResponseFunc,
) {
	startTime := time.Now()
	for range time.NewTicker(500 * time.Millisecond).C {
		txStatus, err := p.client.GetSignatureStatus(ctx, false, sign)
		if err == nil && txStatus.ConfirmationStatus == solrpc.ConfirmationStatusFinalized {
			callback(
				msg.MessageKey(),
				&relayertypes.TxResponse{
					Height: int64(txStatus.Slot),
					TxHash: sign.String(),
				},
				nil,
			)
			return
		} else if time.Since(startTime) > 2*time.Second {
			var cbErr error
			if err != nil {
				cbErr = err
			} else if txStatus.Err != nil {
				cbErr = fmt.Errorf("failed to get tx signature status: %v", txStatus.Err)
			} else {
				cbErr = fmt.Errorf("failed to finalize tx signature")
			}
			callback(
				msg.MessageKey(),
				&relayertypes.TxResponse{
					TxHash: sign.String(),
				},
				cbErr,
			)
			return
		}
	}
}

func (p *Provider) MakeCallInstructions(msg *relayertypes.Message) ([]solana.Instruction, error) {
	switch msg.EventType {
	case relayerevents.EmitMessage:
		return p.getRecvMessageIntruction(msg)
	case relayerevents.CallMessage:
		return p.getExecuteCallIntruction(msg)
	case "sendMessage":
		return p.getSendMessageIntruction(msg)
	default:
		return nil, fmt.Errorf("invalid event type in message")
	}
}

func (p *Provider) getSendMessageIntruction(msg *relayertypes.Message) ([]solana.Instruction, error) {
	sendMsgParams := types.SendMessageParams{
		To:   msg.Dst,
		Data: msg.Data,
	}
	paramBytes, err := BorshEncode(sendMsgParams)
	if err != nil {
		return nil, err
	}

	progID, err := solana.PublicKeyFromBase58(p.cfg.ConnectionAddress)
	if err != nil {
		return nil, err
	}

	payerAccount := solana.AccountMeta{
		PublicKey:  p.wallet.PublicKey(),
		IsWritable: true,
		IsSigner:   true,
	}

	instructionData := append([]byte{types.MethodSendMessage}, paramBytes...)

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID:        progID,
			AccountValues: solana.AccountMetaSlice{&payerAccount},
			DataBytes:     instructionData,
		},
	}

	return instructions, nil
}

func (p *Provider) getRecvMessageIntruction(msg *relayertypes.Message) ([]solana.Instruction, error) {
	recvMsgParams := types.RecvMessageParams{
		Sn:   msg.Sn,
		Src:  msg.Src,
		Data: msg.Data,
	}
	paramBytes, err := BorshEncode(recvMsgParams)
	if err != nil {
		return nil, err
	}

	progID, err := solana.PublicKeyFromBase58(p.cfg.ConnectionAddress)
	if err != nil {
		return nil, err
	}

	payerAccount := solana.AccountMeta{
		PublicKey:  p.wallet.PublicKey(),
		IsWritable: true,
		IsSigner:   true,
	}

	instructionData := append([]byte{types.MethodRecvMessage}, paramBytes...)

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID:        progID,
			AccountValues: solana.AccountMetaSlice{&payerAccount},
			DataBytes:     instructionData,
		},
	}

	return instructions, nil
}

func (p *Provider) getExecuteCallIntruction(msg *relayertypes.Message) ([]solana.Instruction, error) {
	executeCallParams := types.ExecuteCallParams{
		ReqId: msg.ReqID,
		Data:  msg.Data,
	}
	paramBytes, err := BorshEncode(executeCallParams)
	if err != nil {
		return nil, err
	}

	progID, err := solana.PublicKeyFromBase58(p.cfg.XcallAddress)
	if err != nil {
		return nil, err
	}

	payerAccount := solana.AccountMeta{
		PublicKey:  p.wallet.PublicKey(),
		IsWritable: true,
		IsSigner:   true,
	}

	instructionData := append([]byte{types.MethodExecuteCall}, paramBytes...)

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID:        progID,
			AccountValues: solana.AccountMetaSlice{&payerAccount},
			DataBytes:     instructionData,
		},
	}

	return instructions, nil
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txDigest string) (*relayertypes.Receipt, error) {
	return nil, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayertypes.MessageKey) (bool, error) {
	return false, nil
}
