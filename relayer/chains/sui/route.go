package sui

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func (p *Provider) Route(ctx context.Context, message *providerTypes.Message, callback providerTypes.TxResponseFunc) error {
	p.log.Info("starting to route message", zap.Any("message", message))
	suiMessage, err := p.MakeSuiMessage(message)
	if err != nil {
		return err
	}
	messageKey := message.MessageKey()
	txRes, err := p.SendTransaction(ctx, suiMessage)
	go p.executeRouteCallBack(*txRes, messageKey, suiMessage.Method, callback, err)
	if err != nil {
		return errors.Wrapf(err, "error occured while sending transaction in sui")
	}
	return nil
}

func (p *Provider) MakeSuiMessage(message *providerTypes.Message) (*SuiMessage, error) {
	switch message.EventType {
	//TODO generate appropriate callparams for the call
	case events.EmitMessage:
		callParams := []interface{}{
			message.Src,
			message.Sn,
			message.Data,
		}
		return p.NewSuiMessage(callParams, "packageID", "module", "MethodRecvMessage"), nil
	}
	return nil, fmt.Errorf("can't generate message for unknown event type: %s ", message.EventType)
}

func (p *Provider) GetReturnValuesFromCall(ctx context.Context, msg *SuiMessage) (any, error) {
	wallet, err := p.Wallet()
	if err != nil {
		return &types.SuiTransactionBlockResponse{}, err
	}
	return p.client.ExecuteContractAndReturnVal(ctx, msg, wallet.Address, p.cfg.GasLimit)

}

func (p *Provider) SendTransaction(ctx context.Context, msg *SuiMessage) (*types.SuiTransactionBlockResponse, error) {
	wallet, err := p.Wallet()
	if err != nil {
		return &types.SuiTransactionBlockResponse{}, err
	}
	txnMetadata, err := p.client.ExecuteContract(ctx, msg, wallet.Address, p.cfg.GasLimit)
	if err != nil {
		return &types.SuiTransactionBlockResponse{}, err
	}
	dryRunResp, gasRequired, err := p.client.EstimateGas(ctx, txnMetadata.TxBytes)
	if err != nil {
		return &types.SuiTransactionBlockResponse{}, fmt.Errorf("failed estimating gas: %w", err)
	}
	if gasRequired > int64(p.cfg.GasLimit) {
		return &types.SuiTransactionBlockResponse{}, fmt.Errorf("gas requirement is too high: %d", gasRequired)
	}
	if gasRequired < int64(p.cfg.GasMin) {
		return &types.SuiTransactionBlockResponse{}, fmt.Errorf("gas requirement is too low: %d", gasRequired)
	}
	if !dryRunResp.Effects.Data.IsSuccess() {
		return &types.SuiTransactionBlockResponse{}, fmt.Errorf(dryRunResp.Effects.Data.V1.Status.Error)
	}
	signature, err := wallet.SignSecureWithoutEncode(txnMetadata.TxBytes, sui_types.DefaultIntent())
	if err != nil {
		return nil, err
	}
	signatures := []any{signature}
	txnResp, err := p.client.CommitTx(ctx, wallet, txnMetadata.TxBytes, signatures)
	return txnResp, err
}

func (p *Provider) executeRouteCallBack(txRes types.SuiTransactionBlockResponse, messageKey *providerTypes.MessageKey, method string, callback providerTypes.TxResponseFunc, err error) {
	// if error occurred before txn processing
	if err != nil {
		return
	}

	res := &providerTypes.TxResponse{
		TxHash: txRes.Digest.String(),
	}
	intHeight := txRes.Checkpoint.Int64()

	// assign tx successful height
	res.Height = intHeight
	success := txRes.Effects.Data.IsSuccess()
	if !success {
		err = fmt.Errorf("error: %s", txRes.Effects.Data.V1.Status.Error)
		callback(messageKey, res, err)
		p.LogFailedTx(method, messageKey, txRes, err)
		return
	}
	res.Code = providerTypes.Success
	callback(messageKey, res, nil)
	p.LogSuccessTx(method, messageKey, txRes)
}

func (p *Provider) LogSuccessTx(method string, message *providerTypes.MessageKey, txRes types.SuiTransactionBlockResponse) {
	p.log.Info("successful transaction",
		zap.Any("message-key", message),
		zap.String("tx_hash", txRes.Digest.String()),
		zap.Int64("height", txRes.Checkpoint.Int64()),
	)
}

func (p *Provider) LogFailedTx(method string, messageKey *providerTypes.MessageKey, txRes types.SuiTransactionBlockResponse, err error) {
	p.log.Info("failed transaction",
		zap.String("tx_hash", txRes.Digest.String()),
		zap.Int64("height", txRes.Checkpoint.Int64()),
		zap.Error(err),
	)
}
