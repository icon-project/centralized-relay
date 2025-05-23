package evm

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

const (
	ErrorLessGas          = "transaction underpriced"
	ErrorLimitLessThanGas = "max fee per gas less than block base fee"
	ErrUnKnown            = "unknown"
	ErrNonceTooLow        = "nonce too low"
	ErrNonceTooHigh       = "nonce too high"
)

// this will be executed in go route
func (p *Provider) Route(ctx context.Context, message *providerTypes.Message, callback providerTypes.TxResponseFunc) error {
	// lock here to prevent transcation replacement
	p.routerMutex.Lock()
	defer p.routerMutex.Unlock()

	p.log.Info("starting to route message",
		zap.String("src", message.Src),
		zap.String("dst", message.Dst),
		zap.Any("sn", message.Sn),
		zap.Any("req_id", message.ReqID),
		zap.String("event_type", message.EventType),
		zap.String("data", hex.EncodeToString(message.Data)),
	)

	opts, err := p.GetTransationOpts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction options: %w", err)
	}

	tx, err := p.SendTransaction(ctx, opts, message)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	p.WaitForTxResult(ctx, tx, message.MessageKey(), callback)
	return nil
}

func (p *Provider) SendTransaction(ctx context.Context, opts *bind.TransactOpts, message *providerTypes.Message) (*types.Transaction, error) {
	var (
		tx  *types.Transaction
		err error
	)

	gasLimit, err := p.EstimateGas(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}

	if p.cfg.GasLimit > 0 && gasLimit > p.cfg.GasLimit {
		return nil, fmt.Errorf("gas limit exceeded: %d", gasLimit)
	}

	opts.GasLimit = gasLimit + (gasLimit * p.cfg.GasAdjustment / 100)

	p.log.Info("transaction info",
		zap.Any("gas_price", opts.GasPrice),
		zap.Any("gas_cap", opts.GasFeeCap),
		zap.Any("gas_tip", opts.GasTipCap),
		zap.Uint64("estimated_limit", gasLimit),
		zap.Uint64("adjusted_limit", opts.GasLimit),
		zap.Uint64("nonce", opts.Nonce.Uint64()),
		zap.String("event_type", message.EventType),
		zap.String("src", message.Src),
		zap.Uint64("sn", message.Sn.Uint64()),
	)

	switch message.EventType {
	case events.EmitMessage:
		tx, err = p.client.ReceiveMessage(opts, message.Src, message.Sn, message.Data)
	case events.CallMessage:
		tx, err = p.client.ExecuteCall(opts, message.ReqID, message.Data)
	case events.SetAdmin:
		addr := common.HexToAddress(message.Src)
		tx, err = p.client.SetAdmin(opts, addr)
	case events.RevertMessage:
		tx, err = p.client.RevertMessage(opts, message.Sn)
	case events.ClaimFee:
		tx, err = p.client.ClaimFee(opts)
	case events.SetFee:
		tx, err = p.client.SetFee(opts, message.Src, message.Sn, message.ReqID)
	case events.RollbackMessage:
		tx, err = p.client.ExecuteRollback(opts, message.Sn)
	case events.PacketAcknowledged:
		tx, err = p.client.ReceiveMessageWithSignature(opts, message.Src, message.Sn, message.Data, message.Signatures)
	default:
		return nil, fmt.Errorf("unknown event type: %s", message.EventType)
	}
	return tx, err
}

func (p *Provider) WaitForTxResult(ctx context.Context, tx *types.Transaction, m *providerTypes.MessageKey, callback providerTypes.TxResponseFunc) {
	res := &providerTypes.TxResponse{
		TxHash: tx.Hash().String(),
	}

	txReceipts, err := p.WaitForResults(ctx, tx)
	if err != nil {
		callback(m, res, fmt.Errorf("error waiting for tx result: %w", err))
		return
	}

	res.Height = txReceipts.BlockNumber.Int64()

	if txReceipts.Status != types.ReceiptStatusSuccessful {
		res.Code = providerTypes.Failed
		callback(m, res, fmt.Errorf("transaction failed to execute: %+v", txReceipts.Logs))
	} else {
		res.Code = providerTypes.Success
		callback(m, res, nil)
	}
}

func (p *Provider) LogSuccessTx(message *providerTypes.MessageKey, receipt *types.Receipt) {
	p.log.Info("successful transaction",
		zap.Any("message-key", message),
		zap.String("tx_hash", receipt.TxHash.String()),
		zap.Int64("height", receipt.BlockNumber.Int64()),
		zap.Uint64("gas_used", receipt.GasUsed),
		zap.String("contract_address", receipt.ContractAddress.Hex()),
	)
}

func (p *Provider) LogFailedTx(messageKey *providerTypes.MessageKey, result *types.Receipt, err error) {
	p.log.Error("failed transaction",
		zap.String("tx_hash", result.TxHash.String()),
		zap.Int64("height", result.BlockNumber.Int64()),
		zap.Uint64("gas_used", result.GasUsed),
		zap.Uint("tx_index", result.TransactionIndex),
		zap.String("event_type", messageKey.EventType),
		zap.Uint64("sn", messageKey.Sn.Uint64()),
		zap.String("src", messageKey.Src),
		zap.String("contract_address", result.ContractAddress.Hex()),
		zap.Error(err),
	)
}

func (p *Provider) parseErr(err error) string {
	msg := err.Error()
	switch {
	case strings.Contains(msg, ErrorLimitLessThanGas):
		return ErrorLimitLessThanGas
	case strings.Contains(msg, ErrorLessGas):
		return ErrorLessGas
	case strings.Contains(msg, ErrNonceTooLow):
		return ErrNonceTooLow
	case strings.Contains(msg, ErrNonceTooHigh):
		return ErrNonceTooHigh
	default:
		return ErrUnKnown
	}
}
