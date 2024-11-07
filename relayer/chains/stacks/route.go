package stacks

import (
	"context"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/stacks-go-sdk/pkg/clarity"
	"go.uber.org/zap"
)

func (p *Provider) Route(ctx context.Context, message *types.Message, callback types.TxResponseFunc) error {
	p.log.Info("Starting to route message",
		zap.Any("sn", message.Sn),
		zap.Any("req_id", message.ReqID),
		zap.String("src", message.Src),
		zap.String("event_type", message.EventType))

	var txID string
	var err error

	switch message.EventType {
	case events.EmitMessage:
		txID, err = p.handleEmitMessage(ctx, message)
	case events.CallMessage:
		txID, err = p.handleCallMessage(ctx, message)
	case events.RollbackMessage:
		txID, err = p.handleRollbackMessage(ctx, message)
	default:
		return fmt.Errorf("unknown event type: %s", message.EventType)
	}

	if err != nil {
		return fmt.Errorf("failed to handle %s: %w", message.EventType, err)
	}

	p.log.Info("Transaction sent", zap.String("txID", txID))

	receipt, err := p.waitForTransactionConfirmation(ctx, txID)
	if err != nil {
		return fmt.Errorf("failed to confirm transaction: %w", err)
	}

	response := &types.TxResponse{
		TxHash: txID,
		Height: int64(receipt.Height),
		Code:   types.Success,
	}

	callback(message.MessageKey(), response, nil)

	return nil
}

func (p *Provider) handleEmitMessage(ctx context.Context, message *types.Message) (string, error) {
	contractAddress := p.cfg.Contracts[types.ConnectionContract]

	srcNetworkArg, err := clarity.NewStringASCII(message.Src)
	if err != nil {
		return "", fmt.Errorf("failed to create srcNetwork argument: %w", err)
	}

	connSnArg, err := clarity.NewUInt(message.Sn.String())
	if err != nil {
		return "", fmt.Errorf("failed to create connSn argument: %w", err)
	}

	msgArg, err := clarity.NewStringASCII(string(message.Data))
	if err != nil {
		return "", fmt.Errorf("failed to create msg argument: %w", err)
	}

	args := []clarity.ClarityValue{srcNetworkArg, connSnArg, msgArg}

	txID, err := p.client.SendCallMessage(ctx, contractAddress, args, p.cfg.Address, p.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	return txID, nil
}

func (p *Provider) handleCallMessage(ctx context.Context, message *types.Message) (string, error) {
	contractAddress := p.cfg.Contracts[types.XcallContract]

	reqIDArg, err := clarity.NewUInt(message.ReqID.String())
	if err != nil {
		return "", fmt.Errorf("failed to create reqID argument: %w", err)
	}

	dataArg, err := clarity.NewStringASCII(string(message.Data))
	if err != nil {
		return "", fmt.Errorf("failed to create data argument: %w", err)
	}

	args := []clarity.ClarityValue{reqIDArg, dataArg}

	txID, err := p.client.ExecuteCall(ctx, contractAddress, args, p.cfg.Address, p.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to execute call: %w", err)
	}

	return txID, nil
}

func (p *Provider) handleRollbackMessage(ctx context.Context, message *types.Message) (string, error) {
	contractAddress := p.cfg.Contracts[types.XcallContract]

	snArg, err := clarity.NewUInt(message.Sn.String())
	if err != nil {
		return "", fmt.Errorf("failed to create sn argument: %w", err)
	}

	args := []clarity.ClarityValue{snArg}

	txID, err := p.client.ExecuteRollback(ctx, contractAddress, args, p.cfg.Address, p.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to execute rollback: %w", err)
	}

	return txID, nil
}
