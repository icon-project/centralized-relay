package evm

import (
	"context"
	"fmt"
	"math/big"
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
	ErrMaxTried           = "max tried"
	ErrNonceTooLow        = "nonce too low"
	ErrNonceTooHigh       = "nonce too high"
)

// this will be executed in go route
func (p *Provider) Route(ctx context.Context, message *providerTypes.Message, callback providerTypes.TxResponseFunc) error {
	p.log.Info("starting to route message", zap.Any("message", message))

	opts, err := p.GetTransationOpts(ctx)
	if err != nil {
		return fmt.Errorf("routing failed: %w", err)
	}

	messageKey := message.MessageKey()

	tx, err := p.SendTransaction(ctx, opts, message, MaxTxFixtures)
	if err != nil {
		return fmt.Errorf("routing failed: %w", err)
	}
	p.NonceTracker.Inc(p.wallet.Address)
	return p.WaitForTxResult(ctx, tx, messageKey, callback)
}

func (p *Provider) SendTransaction(ctx context.Context, opts *bind.TransactOpts, message *providerTypes.Message, maxRetry uint8) (*types.Transaction, error) {
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

	p.log.Info("gas info", zap.Uint64("gas_price", opts.GasPrice.Uint64()), zap.Uint64("original_limit", gasLimit), zap.Uint64("adjusted_limit", opts.GasLimit))

	switch message.EventType {
	case events.EmitMessage:
		tx, err = p.client.ReceiveMessage(opts, message.Src, new(big.Int).SetUint64(message.Sn), message.Data)
	case events.CallMessage:
		tx, err = p.client.ExecuteCall(opts, new(big.Int).SetUint64(message.ReqID), message.Data)
	case events.SetAdmin:
		addr := common.HexToAddress(message.Src)
		tx, err = p.client.SetAdmin(opts, addr)
	case events.RevertMessage:
		tx, err = p.client.RevertMessage(opts, new(big.Int).SetUint64(message.Sn))
	case events.ClaimFee:
		tx, err = p.client.ClaimFee(opts)
	case events.SetFee:
		tx, err = p.client.SetFee(opts, message.Src, new(big.Int).SetUint64(message.Sn), new(big.Int).SetUint64(message.ReqID))
	case events.ExecuteRollback:
		tx, err = p.client.ExecuteRollback(opts, new(big.Int).SetUint64(message.Sn))
	default:
		return nil, fmt.Errorf("unknown event type: %s", message.EventType)
	}
	if err != nil {
		switch p.parseErr(err, maxRetry > 0) {
		case ErrorLessGas, ErrorLimitLessThanGas:
			p.log.Info("gasfee low", zap.Error(err))
			gasPrice, err := p.client.SuggestGasPrice(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get gas price: %w", err)
			}
			opts.GasPrice = gasPrice
		case ErrNonceTooLow, ErrNonceTooHigh:
			p.log.Info("nonce mismatch", zap.Uint64("nonce", opts.Nonce.Uint64()), zap.Error(err))
			nonce, err := p.client.NonceAt(ctx, p.wallet.Address, nil)
			if err != nil {
				return nil, err
			}
			opts.Nonce = nonce
			p.NonceTracker.Set(p.wallet.Address, nonce)
		default:
			return nil, err
		}
		p.log.Info("adjusted", zap.Uint64("nonce", opts.Nonce.Uint64()), zap.Uint64("gas_price", opts.GasPrice.Uint64()), zap.Any("message", message))
		return p.SendTransaction(ctx, opts, message, maxRetry-1)
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

func (p *Provider) parseErr(err error, shouldParse bool) string {
	msg := err.Error()
	switch {
	case !shouldParse:
		return ErrMaxTried
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
