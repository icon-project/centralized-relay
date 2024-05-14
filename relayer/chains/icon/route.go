package icon

import (
	"context"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (p *Provider) Route(ctx context.Context, message *providerTypes.Message, callback providerTypes.TxResponseFunc) error {
	p.log.Info("starting to route message", zap.Any("message", message))
	iconMessage, err := p.MakeIconMessage(message)
	if err != nil {
		return err
	}
	messageKey := message.MessageKey()

	txhash, err := p.SendTransaction(ctx, iconMessage)
	if err != nil {
		return errors.Wrapf(err, "error occured while sending transaction")
	}
	p.WaitForTxResult(ctx, txhash, messageKey, iconMessage.Method, callback)
	return nil
}

func (p *Provider) MakeIconMessage(message *providerTypes.Message) (*IconMessage, error) {
	switch message.EventType {
	case events.EmitMessage:
		msg := &types.RecvMessage{
			SrcNID: message.Src,
			ConnSn: types.NewHexInt(int64(message.Sn)),
			Msg:    types.NewHexBytes(message.Data),
		}
		return p.NewIconMessage(p.GetAddressByEventType(message.EventType), msg, MethodRecvMessage), nil
	case events.CallMessage:
		msg := &types.ExecuteCall{
			ReqID: types.NewHexInt(int64(message.ReqID)),
			Data:  types.NewHexBytes(message.Data),
		}
		return p.NewIconMessage(p.GetAddressByEventType(message.EventType), msg, MethodExecuteCall), nil
	case events.SetAdmin:
		msg := &types.SetAdmin{
			Relayer: message.Src,
		}
		return p.NewIconMessage(p.GetAddressByEventType(message.EventType), msg, MethodSetAdmin), nil
	case events.RevertMessage:
		msg := &types.RevertMessage{
			Sn: types.NewHexInt(int64(message.Sn)),
		}
		return p.NewIconMessage(p.GetAddressByEventType(message.EventType), msg, MethodRevertMessage), nil
	case events.ClaimFee:
		return p.NewIconMessage(p.GetAddressByEventType(message.EventType), nil, MethodClaimFees), nil
	case events.SetFee:
		msg := &types.SetFee{
			NetworkID: message.Src,
			MsgFee:    types.NewHexInt(int64(message.Sn)),
			ResFee:    types.NewHexInt(int64(message.ReqID)),
		}
		return p.NewIconMessage(p.GetAddressByEventType(message.EventType), msg, MethodSetFee), nil
	}
	return nil, fmt.Errorf("can't generate message for unknown event type: %s ", message.EventType)
}

func (p *Provider) SendTransaction(ctx context.Context, msg *IconMessage) ([]byte, error) {
	wallet, err := p.Wallet()
	if err != nil {
		return nil, err
	}

	txParam := types.TransactionParam{
		Version:     types.NewHexInt(JsonrpcApiVersion),
		FromAddress: types.NewAddress(wallet.Address().Bytes()),
		ToAddress:   msg.Address,
		NetworkID:   types.NewHexInt(p.cfg.NetworkID),
		DataType:    "call",
		Data: types.CallData{
			Method: msg.Method,
			Params: msg.Params,
		},
	}

	step, err := p.client.EstimateStep(txParam)
	if err != nil {
		return nil, fmt.Errorf("failed estimating step: %w", err)
	}

	stepVal, err := step.Int()
	if err != nil {
		return nil, err
	}
	// increase step limit by step buffer

	steps := int64(stepVal + stepVal*p.cfg.StepAdjustment/100)

	if steps > p.cfg.StepLimit {
		return nil, fmt.Errorf("step limit is too high: %d", steps)
	}

	if steps < p.cfg.StepMin {
		return nil, fmt.Errorf("step limit is too low: %d", steps)
	}

	txParam.StepLimit = types.NewHexInt(steps)

	if err := p.client.SignTransaction(wallet, &txParam); err != nil {
		return nil, err
	}

	_, err = p.client.SendTransaction(&txParam)
	if err != nil {
		return nil, err
	}
	return txParam.TxHash.Value()
}

// TODO: review try to remove wait for Tx from packet-transfer and only use this for client and connection creation
func (p *Provider) WaitForTxResult(
	ctx context.Context,
	txHash []byte,
	messageKey *providerTypes.MessageKey,
	method string,
	callback providerTypes.TxResponseFunc,
) {
	if callback == nil {
		// no point to wait for result if callback is nil
		return
	}

	txhash := types.NewHexBytes(txHash)
	res := &providerTypes.TxResponse{
		TxHash: string(txhash),
	}

	txRes, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: txhash})
	if err != nil {
		p.log.Error("get txn result failed", zap.String("txHash", string(txhash)), zap.String("method", method), zap.Error(err))
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
		err = fmt.Errorf("error: %s", err)
		callback(messageKey, res, err)
		p.LogFailedTx(method, txRes, err)
		return
	}
	res.Code = providerTypes.Success
	callback(messageKey, res, nil)
	p.LogSuccessTx(method, txRes)
}

func (p *Provider) LogSuccessTx(method string, result *types.TransactionResult) {
	stepUsed, err := result.StepUsed.Value()
	if err != nil {
		p.log.Error("failed to get step used", zap.Error(err))
	}
	height, err := result.BlockHeight.Value()
	if err != nil {
		p.log.Error("failed to get block height", zap.Error(err))
	}

	p.log.Info("transaction success",
		zap.String("chain_id", p.NID()),
		zap.String("method", method),
		zap.String("tx_hash", string(result.TxHash)),
		zap.Int64("height", height),
		zap.Int64("step_used", stepUsed),
		zap.Int64p("step_limit", &p.cfg.StepLimit),
	)
}

func (p *Provider) LogFailedTx(method string, result *types.TransactionResult, err error) {
	stepUsed, _ := result.StepUsed.Value()
	height, _ := result.BlockHeight.Value()

	p.log.Info("transaction failed",
		zap.String("method", method),
		zap.String("tx_hash", string(result.TxHash)),
		zap.Int64("height", height),
		zap.Int64("step_used", stepUsed),
		zap.Int64p("step_limit", &p.cfg.StepLimit),
		zap.Error(err),
	)
}
