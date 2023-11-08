package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

// this will be executed in go route
func (p *EVMProvider) Route(ctx context.Context, message providerTypes.Message, callback providerTypes.TxResponseFunc) error {
	p.log.Info("starting to route message", zap.Any("message", message))

	opts, err := p.GetTransationOpts(ctx)
	if err != nil {
		return fmt.Errorf("routing failed: %w", err)
	}
	messageKey := message.MessageKey()

	tx, err := p.SendTransaction(ctx, opts, message)
	if err != nil {
		return fmt.Errorf("routing failed: %w", err)
	}
	p.WaitForTxResult(ctx, tx, messageKey, callback)
	return nil
}

func (p *EVMProvider) SendTransaction(ctx context.Context, opts *bind.TransactOpts, message providerTypes.Message) (*types.Transaction, error) {

	switch message.EventType {
	// TODO: estimate and throw error if failed
	case events.EmitMessage:
		tx, err := p.client.ReceiveMessage(opts, message.Src, big.NewInt(int64(message.Sn)), message.Data)
		if err != nil {
			return nil, fmt.Errorf("routing failed: %w ", err)
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
	p.log.Info("Successful Transaction",
		zap.Any("message-key", messageKey),
		zap.String("tx_hash", receipt.TxHash.String()),
		zap.Int64("height", receipt.BlockNumber.Int64()),
	)
}

func (p *EVMProvider) LogFailedTx(messageKey providerTypes.MessageKey, result *types.Receipt, err error) {
	p.log.Info("Failed Transaction",
		zap.String("tx_hash", result.TxHash.String()),
		zap.Int64("height", result.BlockNumber.Int64()),
		zap.Error(err),
	)
}
