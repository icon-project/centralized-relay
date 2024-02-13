package icon

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	defaultBroadcastWaitTimeout = 10 * time.Minute
)

func (p *IconProvider) Route(ctx context.Context, message *providerTypes.Message, callback providerTypes.TxResponseFunc) error {
	iconMessage, err := p.MakeIconMessage(message)
	if err != nil {
		return err
	}
	messageKey := message.MessageKey()

	var txhash []byte

	switch message.EventType {
	case events.EmitMessage:
		txhash, err = p.SendTransaction(ctx, iconMessage)
		if err != nil {
			return errors.Wrapf(err, "error occured while sending transaction")
		}
	case events.CallMessage:
		txhash, err = p.ExecuteCall(ctx, big.NewInt(0).SetUint64(message.ReqID), message.Data)
		if err != nil {
			return errors.Wrapf(err, "error occured while executing call")
		}
	}
	go p.WaitForTxResult(ctx, txhash, messageKey, iconMessage.Method, callback)
	return nil
}

func (p *IconProvider) MakeIconMessage(message *providerTypes.Message) (*IconMessage, error) {
	switch message.EventType {
	case events.EmitMessage:
		msg := types.RecvMessage{
			SrcNID: message.Src,
			ConnSn: types.NewHexInt(int64(message.Sn)),
			Msg:    types.NewHexBytes(message.Data),
		}
		return p.NewIconMessage(p.GetAddressByEventType(message.EventType), msg, MethodRecvMessage), nil
	case events.CallMessage:
		msg := types.SendMessage{
			Msg:   types.NewHexBytes(message.Data),
			Sn:    message.Sn,
			ReqID: message.ReqID,
		}
		return p.NewIconMessage(p.GetAddressByEventType(message.EventType), msg, MethodExecuteCall), nil
	}
	return nil, fmt.Errorf("can't generate message for unknown event type: %s ", message.EventType)
}

func (p *IconProvider) SendTransaction(ctx context.Context, msg *IconMessage) ([]byte, error) {
	wallet, err := p.Wallet()
	if err != nil {
		return nil, err
	}

	txParamEst := &types.TransactionParamForEstimate{
		Version:     types.NewHexInt(JsonrpcApiVersion),
		FromAddress: types.Address(wallet.Address().String()),
		ToAddress:   msg.Address,
		NetworkID:   types.NewHexInt(int64(p.cfg.NetworkID)),
		DataType:    "call",
		Data: types.CallData{
			Method: msg.Method,
			Params: msg.Params,
		},
	}

	step, err := p.client.EstimateStep(txParamEst)
	if err != nil {
		return nil, fmt.Errorf("failed estimating step: %w", err)
	}

	stepVal, err := step.Int()
	if err != nil {
		return nil, err
	}
	stepLimit := types.NewHexInt(int64(stepVal + 200_000))

	txParam := &types.TransactionParam{
		Version:     types.NewHexInt(JsonrpcApiVersion),
		FromAddress: types.Address(wallet.Address().String()),
		ToAddress:   types.Address(p.cfg.Contracts[providerTypes.ConnectionContract]),
		NetworkID:   types.NewHexInt(int64(p.cfg.NetworkID)),
		StepLimit:   stepLimit,
		DataType:    "call",
		Data: types.CallData{
			Method: msg.Method,
			Params: msg.Params,
		},
	}

	if err := p.client.SignTransaction(wallet, txParam); err != nil {
		return nil, err
	}

	_, err = p.client.SendTransaction(txParam)
	if err != nil {
		return nil, err
	}
	return txParam.TxHash.Value()
}

// TODO: review try to remove wait for Tx from packet-transfer and only use this for client and connection creation
func (p *IconProvider) WaitForTxResult(
	ctx context.Context,
	txHash []byte,
	messageKey providerTypes.MessageKey,
	method string,
	callback providerTypes.TxResponseFunc,
) {
	if callback == nil {
		// no point to wait for result if callback is nil
		return
	}

	txhash := types.NewHexBytes(txHash)
	res := providerTypes.TxResponse{}
	res.TxHash = string(txHash)

	_, txRes, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: txhash})
	if err != nil {
		p.log.Error("Failed to get txn result", zap.String("txHash", string(txhash)), zap.String("method", method), zap.Error(err))
		callback(messageKey, res, err)
		return
	}

	height, err := txRes.BlockHeight.Value()
	if err != nil {
		callback(messageKey, res, err)
	}
	// assign tx successful height
	res.Height = height

	status, err := txRes.Status.Int()
	if status != 1 {
		err = fmt.Errorf("Transaction Failed to Execute")
		callback(messageKey, res, err)
		p.LogFailedTx(method, txRes, err)
		return
	}
	res.Code = providerTypes.Success
	callback(messageKey, res, nil)
	p.LogSuccessTx(method, txRes)
}

func (p *IconProvider) LogSuccessTx(method string, result *types.TransactionResult) {
	stepUsed, err := result.StepUsed.Value()
	if err != nil {
		p.log.Error("Failed to get step used", zap.Error(err))
	}
	height, err := result.BlockHeight.Value()
	if err != nil {
		p.log.Error("Failed to get block height", zap.Error(err))
	}

	p.log.Info("Successful Transaction",
		zap.String("chain_id", p.NID()),
		zap.String("method", method),
		zap.String("tx_hash", string(result.TxHash)),
		zap.Int64("height", height),
		zap.Int64("step_used", stepUsed),
	)
}

func (p *IconProvider) LogFailedTx(method string, result *types.TransactionResult, err error) {
	stepUsed, _ := result.StepUsed.Value()
	height, _ := result.BlockHeight.Value()

	p.log.Info("Failed Transaction",
		zap.String("chain_id", p.NID()),
		zap.String("method", method),
		zap.String("tx_hash", string(result.TxHash)),
		zap.Int64("height", height),
		zap.Int64("step_used", stepUsed),
		zap.Error(err),
	)
}
