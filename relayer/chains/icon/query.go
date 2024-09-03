package icon

import (
	"context"
	"fmt"
	"slices"

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

func (p *Provider) GenerateMessages(ctx context.Context, fromHeight, toHeight uint64) ([]*providerTypes.Message, error) {
	p.log.Info("generating message", zap.Uint64("fromHeight", fromHeight), zap.Uint64("toHeight", toHeight))
	return p.QueryBlockMessages(ctx, fromHeight, toHeight)
}

func (p *Provider) FetchTxMessages(ctx context.Context, txHash string) ([]*providerTypes.Message, error) {
	txResult, err := p.client.GetTransactionResult(&types.TransactionHashParam{
		Hash: types.HexBytes(txHash),
	})
	if err != nil {
		return nil, err
	}

	connectionContract := types.Address(p.cfg.Contracts[providerTypes.ConnectionContract])
	xcallContract := types.Address(p.cfg.Contracts[providerTypes.XcallContract])
	allowedAddresses := []types.Address{connectionContract, xcallContract}

	messages := []*providerTypes.Message{}
	for _, log := range txResult.EventLogs {
		if slices.Contains(allowedAddresses, log.Addr) {
			event := types.EventNotificationLog{
				Address: log.Addr,
				Indexed: log.Indexed,
				Data:    log.Data,
			}
			height, err := txResult.BlockHeight.Int64()
			if err != nil {
				return nil, err
			}
			msg, err := p.parseMessageFromEventLog(uint64(height), &event)
			if err != nil {
				p.log.Warn("received invalid event", zap.Error(err))
			} else if msg != nil {
				messages = append(messages, msg)
			}
		}
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
