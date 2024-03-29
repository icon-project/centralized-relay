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
)

// this will be executed in go route
func (p *Provider) Route(ctx context.Context, message *providerTypes.Message, callback providerTypes.TxResponseFunc) error {
	p.log.Info("starting to route message", zap.Any("message", message))

	opts, err := p.GetTransationOpts(ctx)
	if err != nil {
		return fmt.Errorf("routing failed: %w", err)
	}

	messageKey := message.MessageKey()

	tx, err := p.SendTransaction(ctx, opts, message, MaxGasPriceInceremtRetry)
	if err != nil {
		return fmt.Errorf("routing failed: %w", err)
	}
	go p.WaitForTxResult(ctx, tx, messageKey, callback)
	p.NonceTracker.Inc(p.wallet.Address)
	return nil
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
	opts.GasLimit = gasLimit

	switch message.EventType {
	// check estimated gas and gas price
	case events.EmitMessage:
		tx, err = p.client.ReceiveMessage(opts, message.Src, big.NewInt(int64(message.Sn)), message.Data)
	case events.CallMessage:
		tx, err = p.client.ExecuteCall(opts, big.NewInt(0).SetUint64(message.ReqID), message.Data)
	case events.SetAdmin:
		addr := common.HexToAddress(message.Src)
		tx, err = p.client.SetAdmin(opts, addr)
	case events.RevertMessage:
		tx, err = p.client.RevertMessage(opts, big.NewInt(int64(message.Sn)))
	case events.ClaimFee:
		tx, err = p.client.ClaimFee(opts)
	case events.SetFee:
		tx, err = p.client.SetFee(opts, message.Src, big.NewInt(int64(message.Sn)), big.NewInt(int64(message.ReqID)))
	case events.ExecuteRollback:
		tx, err = p.client.ExecuteRollback(opts, big.NewInt(0).SetUint64(message.Sn))
	default:
		return nil, fmt.Errorf("unknown event type: %s", message.EventType)
	}
	if err != nil {
		switch p.parseErr(err, maxRetry > 0) {
		case ErrorLessGas:
			p.log.Info(ErrorLessGas, zap.Uint64("gas_price", opts.GasPrice.Uint64()))
			// GasPriceRatio (10%) percent increment
			opts.GasPrice = big.NewInt(0).Add(opts.GasPrice, big.NewInt(0).Div(opts.GasPrice, big.NewInt(GasPriceRatio)))
		case ErrorLimitLessThanGas:
			p.log.Info("gasfee low", zap.Uint64("gas_price", opts.GasPrice.Uint64()))
			// get gas price parsing error message
			startIndex := strings.Index(err.Error(), "baseFee: ")
			endIndex := strings.Index(err.Error(), "(supplied gas")
			baseGasPrice := err.Error()[startIndex+len("baseFee: ") : endIndex-1]
			gasPrice, ok := big.NewInt(0).SetString(baseGasPrice, 0)
			if !ok {
				gasPrice, err = p.client.SuggestGasPrice(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to get gas price: %w", err)
				}
			}
			opts.GasPrice = gasPrice
		case ErrNonceTooLow:
			p.log.Info("nonce too low", zap.Uint64("nonce", opts.Nonce.Uint64()))
			nonce, err := p.client.NonceAt(ctx, p.wallet.Address, nil)
			if err != nil {
				return nil, err
			}
			opts.Nonce = nonce
		default:
			return nil, err
		}
		p.log.Info("adjusted", zap.Uint64("nonce", opts.Nonce.Uint64()), zap.Uint64("gas_price", opts.GasPrice.Uint64()), zap.Any("message", message))
		return p.SendTransaction(ctx, opts, message, maxRetry-1)
	}
	return tx, err
}

func (p *Provider) WaitForTxResult(
	ctx context.Context,
	tx *types.Transaction,
	message *providerTypes.MessageKey,
	callback providerTypes.TxResponseFunc,
) {
	if callback == nil {
		// no point to wait for result if callback is nil
		return
	}

	res := &providerTypes.TxResponse{
		TxHash: tx.Hash().String(),
	}

	txReceipts, err := p.WaitForResults(ctx, tx.Hash())
	if err != nil {
		p.log.Error("failed to get tx result",
			zap.String("hash", res.TxHash),
			zap.Any("message", message),
			zap.Error(err))
		callback(message, res, err)
		return
	}

	res.Height = txReceipts.BlockNumber.Int64()

	status := txReceipts.Status
	if status != 1 {
		err = fmt.Errorf("transaction failed to execute")
		callback(message, res, err)
		p.LogFailedTx(message, txReceipts, err)
		return
	}
	res.Code = providerTypes.Success
	callback(message, res, nil)
	p.LogSuccessTx(message, txReceipts)
}

func (p *Provider) LogSuccessTx(message *providerTypes.MessageKey, receipt *types.Receipt) {
	p.log.Info("successful transaction",
		zap.Any("message-key", message),
		zap.String("tx_hash", receipt.TxHash.String()),
		zap.Int64("height", receipt.BlockNumber.Int64()),
	)
}

func (p *Provider) LogFailedTx(messageKey *providerTypes.MessageKey, result *types.Receipt, err error) {
	p.log.Info("failed transaction",
		zap.String("tx_hash", result.TxHash.String()),
		zap.Int64("height", result.BlockNumber.Int64()),
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
	default:
		return ErrUnKnown
	}
}
