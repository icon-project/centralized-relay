package icon

import (
	"context"
	"errors"
	"fmt"
	"math/big"
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

func (ip *IconProvider) prepareCallParams(methodName string, param map[string]interface{}, options ...CallParamOption) *types.CallParam {
	callData := &types.CallData{
		Method: methodName,
		Params: param,
	}

	callParam := &types.CallParam{
		FromAddress: types.Address(fmt.Sprintf("hx%s", strings.Repeat("0", 40))),
		ToAddress:   types.Address(ip.cfg.Contracts[providerTypes.ConnectionContract]),
		DataType:    "call",
		Data:        callData,
	}

	for _, option := range options {
		option(callParam)
	}

	return callParam
}

func (ip *IconProvider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	block, err := ip.client.GetLastBlock()
	if err != nil {
		return 0, err
	}
	return uint64(block.Height), nil
}

func (ip *IconProvider) ShouldReceiveMessage(ctx context.Context, messagekey providerTypes.Message) (bool, error) {
	return true, nil
}

func (ip *IconProvider) ShouldSendMessage(ctx context.Context, messageKey providerTypes.Message) (bool, error) {
	return true, nil
}

func (ip *IconProvider) MessageReceived(ctx context.Context, messageKey providerTypes.MessageKey) (bool, error) {
	callParam := ip.prepareCallParams(MethodGetReceipts, map[string]interface{}{
		"srcNetwork": messageKey.Src,
		"_connSn":    types.NewHexInt(int64(messageKey.Sn)),
	})

	var status types.HexInt
	err := ip.client.Call(callParam, &status)
	if err != nil {
		return false, fmt.Errorf("MessageReceived: %v", err)
	}

	if status == types.NewHexInt(1) {
		return true, nil
	}

	return false, nil
}

func (ip *IconProvider) QueryBalance(ctx context.Context, addr string) (*providerTypes.Coin, error) {
	param := types.AddressParam{
		Address: types.Address(addr),
	}
	balance, err := ip.client.GetBalance(&param)
	if err != nil {
		return nil, err
	}
	coin := providerTypes.NewCoin("ICX", balance.Uint64())
	return &coin, nil
}

func (ip *IconProvider) GenerateMessage(ctx context.Context, key *providerTypes.MessageKeyWithMessageHeight) (*providerTypes.Message, error) {
	ip.log.Info("generating message", zap.Any("messagekey", key))
	if key == nil {
		return nil, errors.New("GenerateMessage: message key cannot be nil")
	}

	block, err := ip.client.GetBlockByHeight(&types.BlockHeightParam{
		Height: types.NewHexInt(int64(key.MsgHeight)),
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

		for _, el := range txResult.EventLogs {
			switch el.Indexed[0] {
			case EmitMessage:
				if el.Addr != types.Address(ip.cfg.Contracts[providerTypes.ConnectionContract]) &&
					len(el.Indexed) != 3 && len(el.Data) != 1 {
					continue
				}
			case CallMessage:
				if el.Addr != types.Address(ip.cfg.Contracts[providerTypes.XcallContract]) &&
					len(el.Indexed) != 3 && len(el.Data) != 1 {
					continue
				}
			}

			dst := el.Indexed[1]
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

			return &providerTypes.Message{
				MessageHeight: key.MsgHeight,
				EventType:     key.EventType,
				Dst:           dst,
				Src:           key.Src,
				Data:          dataValue,
				Sn:            sn.Uint64(),
			}, nil
		}
	}

	return nil, fmt.Errorf("error generating message: %v", key)
}

// QueryTransactionReceipt ->
// TxHash should be in hex string
func (icp *IconProvider) QueryTransactionReceipt(ctx context.Context, txHash string) (*providerTypes.Receipt, error) {
	res, err := icp.client.GetTransactionResult(&types.TransactionHashParam{
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

	receipt := providerTypes.Receipt{
		TxHash: txHash,
		Height: height.Uint64(),
	}
	if status == 1 {
		receipt.Status = true
	}
	return &receipt, nil
}

// SetAdmin sets the admin address of the bridge contract
func (ip *IconProvider) SetAdmin(ctx context.Context, admin string) error {
	callParam := map[string]interface{}{
		"_relayer": admin,
	}
	message := ip.NewIconMessage(callParam, "setAdmin")

	data, err := ip.SendTransaction(ctx, message)
	if err != nil {
		return fmt.Errorf("SetAdmin: %v", err)
	}
	ip.log.Info("SetAdmin: tx sent", zap.String("txHash", string(data)))
	return nil
}
