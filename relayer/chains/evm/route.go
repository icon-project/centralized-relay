package evm

import (
	"context"
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
	globalRouteLock.Lock()

	p.log.Info("starting to route message", zap.Any("message", message))

	opts, err := p.GetTransationOpts(ctx)
	if err != nil {
		return fmt.Errorf("routing failed: %w", err)
	}

	messageKey := message.MessageKey()

	tx, err := p.SendTransaction(ctx, opts, message)
	globalRouteLock.Unlock()
	if err != nil {
		return fmt.Errorf("routing failed: %w", err)
	}
	p.log.Info("transaction sent", zap.String("tx_hash", tx.Hash().String()), zap.Any("message", messageKey))
	return p.WaitForTxResult(ctx, tx, messageKey, callback)
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

	if gasLimit > p.cfg.GasLimit {
		return nil, fmt.Errorf("gas limit exceeded: %d", gasLimit)
	}

	if gasLimit < p.cfg.GasMin {
		return nil, fmt.Errorf("gas price less than minimum: %d", gasLimit)
	}

	opts.GasLimit = gasLimit + (gasLimit * p.cfg.GasAdjustment / 100)

	p.log.Info("gas info",
		zap.Uint64("gas_cap", opts.GasFeeCap.Uint64()),
		zap.Uint64("gas_tip", opts.GasTipCap.Uint64()),
		zap.Uint64("estimated_limit", gasLimit),
		zap.Uint64("adjusted_limit", opts.GasLimit),
		zap.Uint64("nonce", opts.Nonce.Uint64()),
		zap.String("event_type", message.EventType),
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
	case events.ExecuteRollback:
		tx, err = p.client.ExecuteRollback(opts, message.Sn)
	default:
		return nil, fmt.Errorf("unknown event type: %s", message.EventType)
	}
	if err != nil {
		switch p.parseErr(err) {
		case ErrNonceTooLow, ErrNonceTooHigh, ErrorLessGas:
			nonce, err := p.client.PendingNonceAt(ctx, p.wallet.Address, nil)
			if err != nil {
				return nil, err
			}
			p.log.Info("nonce mismatch", zap.Uint64("tx", opts.Nonce.Uint64()), zap.Uint64("current", nonce.Uint64()), zap.Error(err))
			p.NonceTracker.Set(p.wallet.Address, nonce)
		default:
			return nil, err
		}
	}
	return tx, err
}

func (p *Provider) WaitForTxResult(ctx context.Context, tx *types.Transaction, m *providerTypes.MessageKey, callback providerTypes.TxResponseFunc) error {
	if callback == nil {
		// no point to wait for result if callback is nil
		return nil
	}

	res := &providerTypes.TxResponse{
		TxHash: tx.Hash().String(),
	}

	txReceipts, err := p.WaitForResults(ctx, tx)
	if err != nil {
		p.log.Error("failed to get tx result", zap.String("hash", res.TxHash), zap.Any("message", m), zap.Error(err))
		callback(m, res, err)
		return err
	}

	res.Height = txReceipts.BlockNumber.Int64()

	if txReceipts.Status != types.ReceiptStatusSuccessful {
		err = fmt.Errorf("transaction failed to execute")
		callback(m, res, err)
		p.LogFailedTx(m, txReceipts, err)
		return err
	}
	res.Code = providerTypes.Success
	callback(m, res, nil)
	p.LogSuccessTx(m, txReceipts)
	return nil
}

func (p *Provider) LogSuccessTx(message *providerTypes.MessageKey, receipt *types.Receipt) {
	p.log.Info("successful transaction",
		zap.Any("message-key", message),
		zap.String("tx_hash", receipt.TxHash.String()),
		zap.Int64("height", receipt.BlockNumber.Int64()),
		zap.Uint64("gas_used", receipt.GasUsed),
	)
}

func (p *Provider) LogFailedTx(messageKey *providerTypes.MessageKey, result *types.Receipt, err error) {
	p.log.Info("failed transaction",
		zap.String("tx_hash", result.TxHash.String()),
		zap.Int64("height", result.BlockNumber.Int64()),
		zap.Uint64("gas_used", result.GasUsed),
		zap.Uint("tx_index", result.TransactionIndex),
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
