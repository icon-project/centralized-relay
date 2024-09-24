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
	p.log.Info("starting to route message",
		zap.Any("sn", message.Sn),
		zap.Any("req_id", message.ReqID),
		zap.String("src", message.Src),
		zap.String("event_type", message.EventType))

	iconMessage, err := p.MakeIconMessage(message)
	if err != nil {
		return err
	}
	messageKey := message.MessageKey()

	txhash, err := p.SendTransaction(ctx, iconMessage)
	if err != nil {
		return errors.Wrapf(err, "error occured while sending transaction")
	}
	return p.WaitForTxResult(ctx, txhash, messageKey, iconMessage.Method, callback)
}

func (p *Provider) MakeIconMessage(message *providerTypes.Message) (*IconMessage, error) {
	switch message.EventType {
	case events.AcknowledgeMessage:
		msg := &types.AcknowledgePacket{
			Source:      message.Src,
			Sn:          types.NewHexInt(message.Sn.Int64()),
			SignedBytes: types.NewHexBytes(message.Data),
		}
		return p.NewIconMessage(p.GetAddressByEventType(events.AcknowledgeMessage), msg, MethodAcknowledgePacket), nil
	case events.EmitMessage:
		if p.cfg.ClusterMode {
			msg := &types.RegisterPacket{
				Data:        types.NewHexBytes(message.Data),
				Sn:          types.NewHexInt(message.Sn.Int64()),
				Height:      types.NewHexInt(int64(message.MessageHeight)),
				Source:      message.Src,
				Destination: message.Dst,
				ConnAddress: message.ConnAddress,
			}
			return p.NewIconMessage(p.GetAddressByEventType(events.AcknowledgeMessage), msg, MethodRegisterPacket), nil
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

// TODO: review try to remove wait for Tx from packet-transfer and only use this for client and connection creation
func (p *Provider) WaitForTxResult(
	ctx context.Context,
	txHash []byte,
	messageKey *providerTypes.MessageKey,
	method string,
	callback providerTypes.TxResponseFunc,
) error {
	if callback == nil {
		// no point to wait for result if callback is nil
		return nil
	}

	txhash := types.NewHexBytes(txHash)
	res := &providerTypes.TxResponse{
		TxHash: string(txhash),
	}

	txRes, err := p.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: txhash})
	if err != nil {
		p.log.Error("get txn result failed", zap.String("txHash", string(txhash)), zap.String("method", method), zap.Error(err))
		callback(messageKey, res, err)
		return err
	}

	height, err := txRes.BlockHeight.Value()
	if err != nil {
		callback(messageKey, res, err)
	}
	// assign tx successful height
	res.Height = height

	if status, err := txRes.Status.Int(); status != 1 || err != nil {
		err = fmt.Errorf("error: %s", err)
		callback(messageKey, res, err)
		p.LogFailedTx(method, txRes, err)
		return err
	}
	res.Code = providerTypes.Success
	callback(messageKey, res, nil)
	p.LogSuccessTx(method, txRes)
	return nil
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

func (p *Provider) RegisterClusterMessage(ctx context.Context, message *providerTypes.Message, callback providerTypes.TxResponseFunc) error {
	p.log.Info("starting to register message",
		zap.Any("sn", message.Sn),
		zap.Any("req_id", message.ReqID),
		zap.String("src", message.Src),
		zap.String("event_type", message.EventType))
	iconMessage, err := p.MakeIconMessage(message)
	if err != nil {
		return err
	}
	messageKey := message.MessageKey()
	txhash, err := p.SendTransaction(ctx, iconMessage)
	if err != nil {
		return errors.Wrapf(err, "error occured while sending transaction")
	}
	return p.WaitForTxResult(ctx, txhash, messageKey, iconMessage.Method, callback)
}

func (p *Provider) AcknowledgeClusterMessage(ctx context.Context, message *providerTypes.Message, callback providerTypes.TxResponseFunc) error {
	p.log.Info("starting to acknowledge message",
		zap.Any("sn", message.Sn),
		zap.Any("req_id", message.ReqID),
		zap.String("src", message.Src),
		zap.String("event_type", message.EventType))

	iconMessage, err := p.MakeIconMessage(message)
	if err != nil {
		return err
	}
	messageKey := message.MessageKey()
	txhash, err := p.SendTransaction(ctx, iconMessage)
	if err != nil {
		return errors.Wrapf(err, "error occured while sending transaction")
	}
	return p.WaitForTxResult(ctx, txhash, messageKey, iconMessage.Method, callback)
}

func (p *Provider) VerifyMessage(ctx context.Context, key *providerTypes.MessageKeyWithMessageHeight) ([]*providerTypes.Message, error) {
	p.log.Info("Verifying message", zap.Any("messagekey", key))
	if key == nil {
		return nil, errors.New("GenerateMessage: message key cannot be nil")
	}

	block, err := p.client.GetBlockByHeight(&types.BlockHeightParam{
		Height: types.NewHexInt(int64(key.Height)),
	})
	if err != nil {
		return nil, fmt.Errorf("GenerateMessage:GetBlockByHeight %v", err)
	}
	var messages []*providerTypes.Message
	for _, res := range block.NormalTransactions {
		txResult, err := p.client.GetTransactionResult(&types.TransactionHashParam{Hash: res.TxHash})
		if err != nil {
			return nil, fmt.Errorf("GenerateMessage:GetTransactionResult %v", err)
		}

		for _, el := range txResult.EventLogs {
			var (
				dst       string
				eventType = p.GetEventName(el.Indexed[0])
			)
			height, err := txResult.BlockHeight.BigInt()
			if err != nil {
				return nil, fmt.Errorf("GenerateMessage: bigIntConversion %v", err)
			}
			switch el.Indexed[0] {
			case EmitMessage:
				if len(el.Indexed) != 3 || len(el.Data) != 1 {
					continue
				}
				if len(el.Indexed) != 3 || len(el.Data) != 1 {
					continue
				}
				dst = el.Indexed[1]
				sn, err := types.HexInt(el.Indexed[2]).BigInt()
				if err != nil {
					p.log.Error("GenerateMessage: error decoding int value ")
					continue
				}
				data := types.HexBytes(el.Data[0])
				dataValue, err := data.Value()
				if err != nil {
					p.log.Error("GenerateMessage: error decoding data ", zap.Error(err))
					continue
				}
				msg := &providerTypes.Message{
					MessageHeight: height.Uint64(),
					EventType:     eventType,
					Dst:           dst,
					Src:           p.NID(),
					Data:          dataValue,
					Sn:            sn,
				}
				messages = append(messages, msg)
			}
		}
	}
	if len(messages) == 0 {
		return nil, errors.New("GenerateMessage: no messages found")
	}
	return messages, nil
}
