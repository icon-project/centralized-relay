package icon

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (p *Provider) Route(ctx context.Context, message *providerTypes.Message, callback providerTypes.TxResponseFunc) error {
	p.log.Info("starting to route message",
		zap.String("src", message.Src),
		zap.String("dst", message.Dst),
		zap.Any("sn", message.Sn),
		zap.Any("req_id", message.ReqID),
		zap.String("event_type", message.EventType),
		zap.String("data", hex.EncodeToString(message.Data)),
	)

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
		if p.cfg.ClusterMode {
			msg := &types.SubmitPacket{
				Data:           types.NewHexBytes(message.Data),
				Sn:             types.NewHexInt(message.Sn.Int64()),
				Height:         types.NewHexInt(int64(message.MessageHeight)),
				Source:         message.Src,
				Destination:    message.Dst,
				DstConnAddress: message.DstConnAddress,
				SrcConnAddress: message.SrcConnAddress,
				SignedBytes:    types.NewHexBytes(message.SignedData),
			}
			return p.NewIconMessage(p.GetAddressByEventType(events.PacketRegistered), msg, MethodSubmitPacket), nil
		}
		msg := &types.RecvMessage{
			SrcNID: message.Src,
			ConnSn: types.NewHexInt(message.Sn.Int64()),
			Msg:    types.NewHexBytes(message.Data),
		}
		return p.NewIconMessage(p.GetAddressByEventType(message.EventType), msg, MethodRecvMessage), nil
	case events.CallMessage:
		msg := &types.ExecuteCall{
			ReqID: types.NewHexInt(message.ReqID.Int64()),
			Data:  types.NewHexBytes(message.Data),
		}
		return p.NewIconMessage(p.GetAddressByEventType(message.EventType), msg, MethodExecuteCall), nil
	case events.RollbackMessage:
		msg := &types.ExecuteRollback{
			Sn: types.NewHexInt(message.Sn.Int64()),
		}
		return p.NewIconMessage(p.GetAddressByEventType(message.EventType), msg, MethodExecuteRollback), nil
	case events.SetAdmin:
		msg := &types.SetAdmin{
			Relayer: message.Src,
		}
		return p.NewIconMessage(p.GetAddressByEventType(message.EventType), msg, MethodSetAdmin), nil
	case events.RevertMessage:
		msg := &types.RevertMessage{
			Sn: types.NewHexInt(message.Sn.Int64()),
		}
		return p.NewIconMessage(p.GetAddressByEventType(message.EventType), msg, MethodRevertMessage), nil
	case events.ClaimFee:
		return p.NewIconMessage(p.GetAddressByEventType(message.EventType), nil, MethodClaimFees), nil
	case events.SetFee:
		msg := &types.SetFee{
			NetworkID: message.Src,
			MsgFee:    types.NewHexInt(message.Sn.Int64()),
			ResFee:    types.NewHexInt(message.ReqID.Int64()),
		}
		return p.NewIconMessage(p.GetAddressByEventType(message.EventType), msg, MethodSetFee), nil
	case events.PacketRegistered:
		msg := &types.SubmitPacket{
			Data:           types.NewHexBytes(message.Data),
			Sn:             types.NewHexInt(message.Sn.Int64()),
			Height:         types.NewHexInt(int64(message.MessageHeight)),
			Source:         message.Src,
			Destination:    message.Dst,
			SrcConnAddress: message.SrcConnAddress,
			DstConnAddress: message.DstConnAddress,
			SignedBytes:    types.NewHexBytes(message.SignedData),
		}
		return p.NewIconMessage(p.GetAddressByEventType(events.PacketRegistered), msg, MethodSubmitPacket), nil
	case events.PacketAcknowledged:
		var sigs []types.HexBytes
		for _, sig := range message.Signatures {
			sigs = append(sigs, types.NewHexBytes(sig))
		}
		msg := &types.RecvMessageWithSignature{
			SrcNID:     message.Src,
			ConnSn:     types.NewHexInt(message.Sn.Int64()),
			Msg:        types.NewHexBytes(message.Data),
			Signatures: sigs,
		}
		return p.NewIconMessage(p.GetAddressByEventType(events.EmitMessage), msg, MethodRecvMessageWithSignature), nil
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
		NetworkID:   p.NetworkID(),
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

	steps, err := step.Int64()
	if err != nil {
		return nil, err
	}

	if steps > p.cfg.StepLimit {
		return nil, fmt.Errorf("step limit is too high: %d", steps)
	}

	if steps < p.cfg.StepMin {
		return nil, fmt.Errorf("step limit is too low: %d", steps)
	}

	steps += steps * p.cfg.StepAdjustment / 100

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

func (p *Provider) WaitForTxResult(
	ctx context.Context,
	txHash []byte,
	messageKey *providerTypes.MessageKey,
	method string,
	callback providerTypes.TxResponseFunc,
) {
	txhash := types.NewHexBytes(txHash)
	res := &providerTypes.TxResponse{
		TxHash: string(txhash),
	}

	txRes, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: txhash})
	if err != nil {
		callback(messageKey, res, err)
		return
	}

	height, err := txRes.BlockHeight.Value()
	if err != nil {
		callback(messageKey, res, err)
		return
	}

	res.Height = height

	if status, err := txRes.Status.Int(); status != 1 || err != nil {
		res.Code = providerTypes.Failed
		callback(messageKey, res, fmt.Errorf("transaction failed: %w", err))
	} else {
		res.Code = providerTypes.Success
		callback(messageKey, res, nil)
	}
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
