package icon

import (
	"context"
	"fmt"
	"strings"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
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
		ToAddress:   types.Address(ip.PCfg.ContractAddress),
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
