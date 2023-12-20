package evm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

var (
	ErrGasPriceTooHigh = fmt.Errorf("gas price is too high")
	ErrGasLimitTooHigh = fmt.Errorf("gas limit is too high")
)

// this will be executed in go route
func (p *EVMProvider) Route(ctx context.Context, message *providerTypes.Message, callback providerTypes.TxResponseFunc) error {
	p.log.Info("starting to route message", zap.Any("message", message))

	opts, err := p.GetTransationOpts(ctx)
	if err != nil {
		return fmt.Errorf("routing failed: %w", err)
	}

	messageKey := message.MessageKey()

	tx, err := p.SendTransaction(ctx, opts, message, MaxGasPriceInceremtRetry)
	if err != nil {
		if errors.Is(err, ErrGasPriceTooHigh) || errors.Is(err, ErrGasLimitTooHigh) {
			p.log.Info("failed to send transaction", zap.Error(err))
			return nil
		}
		return fmt.Errorf("routing failed: %w", err)
	}
	p.WaitForTxResult(ctx, tx, messageKey, callback)
	return nil
}

func (p *EVMProvider) SendTransaction(ctx context.Context, opts *bind.TransactOpts, message *providerTypes.Message, maxRetry uint8) (*types.Transaction, error) {
	switch message.EventType {
	// check estimated gas and gas price
	case events.EmitMessage:
		tx, err := p.client.ReceiveMessage(opts, message.Src, big.NewInt(int64(message.Sn)), message.Data)
		if err != nil {
			if strings.HasPrefix(err.Error(), "transaction underpriced") && maxRetry > 0 {
				p.log.Info("transaction underpriced", zap.Uint64("gas_price", opts.GasPrice.Uint64()), zap.Uint64("gas_limit", opts.GasLimit))
				gasRatio := float64(GasPriceRatio) / 100 * float64(p.cfg.GasPrice) // 10% of gas price
				gas := big.NewFloat(gasRatio)
				gasPrice, _ := gas.Int(nil)
				opts.GasPrice = big.NewInt(0).Add(opts.GasPrice, gasPrice)
				p.log.Info("adjusted", zap.Uint64("gas_price", opts.GasPrice.Uint64()))
				opts.GasLimit = 0
				return p.SendTransaction(ctx, opts, message, maxRetry-1)
			}
			return nil, err
		}
		return tx, nil
	}
	return nil, fmt.Errorf("contract method missing for eventtype: %s", message.EventType)
}

func (p *EVMProvider) WaitForTxResult(
	ctx context.Context,
	tx *types.Transaction,
	messageKey providerTypes.MessageKey,
	callback providerTypes.TxResponseFunc,
) {
	if callback == nil {
		// no point to wait for result if callback is nil
		return
	}

	res := providerTypes.TxResponse{}
	res.TxHash = tx.Hash().String()

	txReceipts, err := p.WaitForResults(ctx, tx.Hash())
	if err != nil {
		p.log.Error("failed to get txn result",
			zap.String("txHash", res.TxHash),
			zap.Any("messagekey ", messageKey),
			zap.Error(err))
		callback(messageKey, res, err)
		return
	}

	res.Height = txReceipts.BlockNumber.Int64()

	status := txReceipts.Status
	if status != 1 {
		err = fmt.Errorf("transaction failed to execute")
		callback(messageKey, res, err)
		p.LogFailedTx(messageKey, txReceipts, err)
		return
	}
	res.Code = providerTypes.Success
	callback(messageKey, res, nil)
	p.LogSuccessTx(messageKey, txReceipts)
}

func (p *EVMProvider) LogSuccessTx(messageKey providerTypes.MessageKey, receipt *types.Receipt) {
	p.log.Info("successful transaction",
		zap.Any("message-key", messageKey),
		zap.String("tx_hash", receipt.TxHash.String()),
		zap.Int64("height", receipt.BlockNumber.Int64()),
	)
}

func (p *EVMProvider) LogFailedTx(messageKey providerTypes.MessageKey, result *types.Receipt, err error) {
	p.log.Info("failed transaction",
		zap.String("tx_hash", result.TxHash.String()),
		zap.Int64("height", result.BlockNumber.Int64()),
		zap.Error(err),
	)
}
