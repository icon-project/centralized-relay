package icon

import (
	"context"
	"errors"
	"fmt"
	"math/big"

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

func (ip *Provider) GenerateMessages(ctx context.Context, key *providerTypes.MessageKeyWithMessageHeight) ([]*providerTypes.Message, error) {
	ip.log.Info("generating message", zap.Any("messagekey", key))
	if key == nil {
		return nil, errors.New("GenerateMessage: message key cannot be nil")
	}

	block, err := ip.client.GetBlockByHeight(&types.BlockHeightParam{
		Height: types.NewHexInt(int64(key.Height)),
	})
	if err != nil {
		return nil, fmt.Errorf("GenerateMessage:GetBlockByHeight %v", err)
	}

	for _, res := range block.NormalTransactions {
		txResult, err := ip.client.GetTransactionResult(&types.TransactionHashParam{
			Hash: res.TxHash,
		})
		if err != nil {
			return nil, fmt.Errorf("GenerateMessage:GetTransactionResult %v", err)
		}

		var messages []*providerTypes.Message

		for _, el := range txResult.EventLogs {
			var dst string
			switch el.Indexed[0] {
			case EmitMessage:
				if el.Addr != types.Address(ip.cfg.Contracts[providerTypes.ConnectionContract]) &&
					len(el.Indexed) != 3 && len(el.Data) != 1 {
					continue
				}
				dst = el.Indexed[1]
			case CallMessage:
				if el.Addr != types.Address(ip.cfg.Contracts[providerTypes.XcallContract]) &&
					len(el.Indexed) != 4 && len(el.Data) != 1 {
					continue
				}
				dst = ip.NID()
			}

			sn, ok := big.NewInt(0).SetString(el.Indexed[2], 0)
			if !ok {
				ip.log.Error("GenerateMessage: error decoding int value ")
				continue
			}

			data := types.HexBytes(el.Data[0])
			dataValue, err := data.Value()
			if err != nil {
				ip.log.Error("GenerateMessage: error decoding data ", zap.Error(err))
				continue
			}

			msg := &providerTypes.Message{
				MessageHeight: key.Height,
				EventType:     key.EventType,
				Dst:           dst,
				Src:           key.Src,
				Data:          dataValue,
				Sn:            sn.Uint64(),
			}
			messages = append(messages, msg)
		}
		return messages, nil
	}

	return nil, fmt.Errorf("error generating message: %v", key)
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
