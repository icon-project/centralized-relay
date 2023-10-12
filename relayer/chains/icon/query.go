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

func (ip *IconProvider) QueryBalance(ctx context.Context, Address string) (uint64, error) {
	return 0, nil
}
