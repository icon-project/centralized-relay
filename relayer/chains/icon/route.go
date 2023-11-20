package icon

import (
	"context"
	"fmt"
	"time"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	defaultBroadcastWaitTimeout = 10 * time.Minute
)

func (icp *IconProvider) Route(ctx context.Context, message providerTypes.Message, callback providerTypes.TxResponseFunc) error {

	iconMessage, err := icp.MakeIconMessage(message)
	if err != nil {
		return err
	}

	messageKey := message.MessageKey()
	txhash, err := icp.SendTransaction(ctx, iconMessage)
	if err != nil {
		return errors.Wrapf(err, "error occured while sending transaction")
	}

	go icp.WaitForTxResult(ctx, txhash, messageKey, iconMessage.Method, callback)

	return nil
}

func (icp *IconProvider) MakeIconMessage(message providerTypes.Message) (IconMessage, error) {

	switch message.EventType {
	case events.EmitMessage:
		msg := types.RecvMessage{
			SrcNetwork: message.Src,
			Sn:         message.Sn,
			Msg:        types.HexBytes(message.Data),
		}
		return icp.NewIconMessage(msg, MethodRecvMessage), nil

	}
	return IconMessage{}, fmt.Errorf("can't generate message for unknown event type: %s ", message.EventType)
}

func (icp *IconProvider) SendTransaction(
	ctx context.Context,
	msg IconMessage) ([]byte, error) {
	wallet, err := icp.Wallet()
	if err != nil {
		return nil, err
	}

	txParamEst := &types.TransactionParamForEstimate{
		Version:     types.NewHexInt(JsonrpcApiVersion),
		FromAddress: types.Address(wallet.Address().String()),
		ToAddress:   types.Address(icp.PCfg.ContractAddress),
		NetworkID:   types.NewHexInt(icp.PCfg.ICONNetworkID),
		DataType:    "call",
		Data: types.CallData{
			Method: msg.Method,
			Params: msg.Params,
		},
	}

	step, err := icp.client.EstimateStep(txParamEst)
	if err != nil {
		return nil, fmt.Errorf("failed estimating step: %w", err)
	}
	stepVal, err := step.Int()
	if err != nil {
		return nil, err
	}
	stepLimit := types.NewHexInt(int64(stepVal + 200_000))

	txParam := &types.TransactionParam{
		Version:     types.NewHexInt(JsonrpcApiVersion),
		FromAddress: types.Address(wallet.Address().String()),
		ToAddress:   types.Address(icp.PCfg.ContractAddress),
		NetworkID:   types.NewHexInt(icp.PCfg.ICONNetworkID),
		StepLimit:   stepLimit,
		DataType:    "call",
		Data: types.CallData{
			Method: msg.Method,
			Params: msg.Params,
		},
	}

	if err := icp.client.SignTransaction(wallet, txParam); err != nil {
		return nil, err
	}
	_, err = icp.client.SendTransaction(txParam)
	if err != nil {
		return nil, err
	}
	return txParam.TxHash.Value()
}

// TODO: review try to remove wait for Tx from packet-transfer and only use this for client and connection creation
func (icp *IconProvider) WaitForTxResult(
	ctx context.Context,
	txHash []byte,
	messageKey providerTypes.MessageKey,
	method string,
	callback providerTypes.TxResponseFunc,
) {
	if callback == nil {
		// no point to wait for result if callback is nil
		return
	}

	txhash := types.NewHexBytes(txHash)
	res := providerTypes.TxResponse{}
	res.TxHash = string(txHash)

	_, txRes, err := icp.client.WaitForResults(ctx, &types.TransactionHashParam{Hash: txhash})
	if err != nil {
		icp.log.Error("Failed to get txn result", zap.String("txHash", string(txhash)), zap.String("method", method), zap.Error(err))
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
		err = fmt.Errorf("Transaction Failed to Execute")
		callback(messageKey, res, err)
		icp.LogFailedTx(method, txRes, err)
		return
	}
	res.Code = providerTypes.Success
	callback(messageKey, res, nil)
	icp.LogSuccessTx(method, txRes)
}

func (icp *IconProvider) LogSuccessTx(method string, result *types.TransactionResult) {
	stepUsed, _ := result.StepUsed.Value()
	height, _ := result.BlockHeight.Value()

	icp.log.Info("Successful Transaction",
		zap.String("chain_id", icp.ChainId()),
		zap.String("method", method),
		zap.String("tx_hash", string(result.TxHash)),
		zap.Int64("height", height),
		zap.Int64("step_used", stepUsed),
	)
}

func (icp *IconProvider) LogFailedTx(method string, result *types.TransactionResult, err error) {
	stepUsed, _ := result.StepUsed.Value()
	height, _ := result.BlockHeight.Value()

	icp.log.Info("Failed Transaction",
		zap.String("chain_id", icp.ChainId()),
		zap.String("method", method),
		zap.String("tx_hash", string(result.TxHash)),
		zap.Int64("height", height),
		zap.Int64("step_used", stepUsed),
		zap.Error(err),
	)
}
