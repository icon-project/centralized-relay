package icon

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

type CallParamOption func(*types.CallParam)

func callParamsWithHeight(height types.HexInt) CallParamOption {
	return func(cp *types.CallParam) {
		cp.Height = height
	}
}

func (p *Provider) prepareCallParams(methodName string, address string, param map[string]interface{}, options ...CallParamOption) *types.CallParam {
	callData := &types.CallData{
		Method: methodName,
		Params: param,
	}

	callParam := &types.CallParam{
		FromAddress: types.Address(p.cfg.Address),
		ToAddress:   types.Address(address),
		DataType:    "call",
		Data:        callData,
	}

	for _, option := range options {
		option(callParam)
	}

	return callParam
}

func (ip *Provider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	block, err := ip.client.GetLastBlock()
	if err != nil {
		return 0, err
	}
	return uint64(block.Height), nil
}

func (ip *Provider) ShouldReceiveMessage(ctx context.Context, messagekey *providerTypes.Message) (bool, error) {
	return true, nil
}

func (ip *Provider) ShouldSendMessage(ctx context.Context, messageKey *providerTypes.Message) (bool, error) {
	return true, nil
}

func (ip *Provider) QueryBalance(ctx context.Context, addr string) (*providerTypes.Coin, error) {
	param := types.AddressParam{
		Address: types.Address(addr),
	}
	balance, err := ip.client.GetBalance(&param)
	if err != nil {
		return nil, err
	}
	return providerTypes.NewCoin("ICX", balance.Uint64()), nil
}

func (p *Provider) GenerateMessages(ctx context.Context, key *providerTypes.MessageKeyWithMessageHeight) ([]*providerTypes.Message, error) {
	p.log.Info("generating message", zap.Any("messagekey", key))
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
				if el.Addr != types.Address(p.cfg.Contracts[providerTypes.ConnectionContract]) || len(el.Indexed) != 3 || len(el.Data) != 1 {
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
					Src:           key.Src,
					Data:          dataValue,
					Sn:            sn,
				}
				messages = append(messages, msg)
			case CallMessage:
				if el.Addr != types.Address(p.cfg.Contracts[providerTypes.XcallContract]) || len(el.Indexed) != 4 || len(el.Data) != 2 {
					continue
				}
				dst = p.NID()
				src := strings.SplitN(string(el.Indexed[1][:]), "/", 2)
				sn, err := types.HexInt(el.Indexed[3]).BigInt()
				if err != nil {
					return nil, fmt.Errorf("failed to parse sn: %s", el.Indexed[2])
				}
				requestID, err := types.HexInt(el.Data[0]).BigInt()
				if err != nil {
					return nil, fmt.Errorf("failed to parse reqID: %s", el.Data[0])
				}
				data, err := types.HexBytes(el.Data[1]).Value()
				if err != nil {
					p.log.Error("GenerateMessage: error decoding data ", zap.Error(err))
					continue
				}
				msg := &providerTypes.Message{
					MessageHeight: height.Uint64(),
					EventType:     p.GetEventName(el.Indexed[0]),
					Dst:           dst,
					Src:           src[0],
					Data:          data,
					Sn:            sn,
					ReqID:         requestID,
				}
				messages = append(messages, msg)
			case RollbackMessage:
				if el.Addr != types.Address(p.cfg.Contracts[providerTypes.XcallContract]) || len(el.Indexed) != 4 || len(el.Data) != 2 {
					continue
				}
				sn, err := types.HexInt(el.Indexed[3]).BigInt()
				if err != nil {
					return nil, fmt.Errorf("failed to parse sn: %s", el.Indexed[2])
				}
				msg := &providerTypes.Message{
					MessageHeight: height.Uint64(),
					EventType:     p.GetEventName(el.Indexed[0]),
					Dst:           p.NID(),
					Src:           p.NID(),
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

// QueryTransactionReceipt ->
// TxHash should be in hex string
func (p *Provider) QueryTransactionReceipt(ctx context.Context, txHash string) (*providerTypes.Receipt, error) {
	res, err := p.client.GetTransactionResult(&types.TransactionHashParam{
		Hash: types.HexBytes(txHash),
	})
	if err != nil {
		return nil, fmt.Errorf("QueryTransactionReceipt: GetTransactionResult: %v", err)
	}

	height, err := res.BlockHeight.BigInt()
	if err != nil {
		return nil, fmt.Errorf("QueryTransactionReceipt: bigIntConversion %v", err)
	}

	status, err := res.Status.Int()
	if err != nil {
		return nil, fmt.Errorf("QueryTransactionReceipt: bigIntConversion %v", err)
	}

	receipt := &providerTypes.Receipt{
		TxHash: txHash,
		Height: height.Uint64(),
		Status: status == 1,
	}
	return receipt, nil
}
