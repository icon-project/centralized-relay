package icon

import (
	"context"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (p *Provider) SubmitClusterMessage(ctx context.Context, message *providerTypes.Message, callback providerTypes.TxResponseFunc) error {
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
