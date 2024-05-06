package sui

import (
	"context"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/errors"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/fardream/go-bcs/bcs"
	"github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func (p *Provider) Route(ctx context.Context, message *relayertypes.Message, callback relayertypes.TxResponseFunc) error {
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

func (p *Provider) MakeSuiMessage(message *relayertypes.Message) (*SuiMessage, error) {
	switch message.EventType {
	case events.EmitMessage:
		snU128, err := bcs.NewUint128FromBigInt(bcs.NewBigIntFromUint64(message.Sn))
		if err != nil {
			return nil, err
		}
		callParams := []SuiCallArg{
			{Type: CallArgObject, Val: p.cfg.XcallStorageID},
			{Type: CallArgPure, Val: message.Src},
			{Type: CallArgPure, Val: snU128},
			{Type: CallArgPure, Val: "0x" + hex.EncodeToString(message.Data)},
		}
		return p.NewSuiMessage(callParams, p.cfg.XcallPkgID, EntryModule, MethodRecvMessage), nil
	case events.CallMessage:
		reqIdU128, err := bcs.NewUint128FromBigInt(bcs.NewBigIntFromUint64(message.ReqID))
		if err != nil {
			return nil, err
		}
		callParams := []SuiCallArg{
			{Type: CallArgObject, Val: p.cfg.DappStateID},
			{Type: CallArgObject, Val: p.cfg.XcallStorageID},
			{Type: CallArgPure, Val: reqIdU128},
			{Type: CallArgPure, Val: "0x" + hex.EncodeToString(message.Data)},
		}
		return p.NewSuiMessage(callParams, p.cfg.DappPkgID, DappModule, MethodExecuteCall), nil
	default:
		return nil, fmt.Errorf("can't generate message for unknown event type: %s ", message.EventType)
	}
}

func (p *Provider) SendTransaction(ctx context.Context, msg *SuiMessage) (*types.SuiTransactionBlockResponse, error) {
	wallet, err := p.Wallet()
	if err != nil {
		return nil, err
	}
	txnMetadata, err := p.client.ExecuteContract(ctx, msg, wallet.Address, p.cfg.GasLimit)
	if err != nil {
		return nil, err
	}
	dryRunResp, gasRequired, err := p.client.EstimateGas(ctx, txnMetadata.TxBytes)
	if err != nil {
		return nil, fmt.Errorf("failed estimating gas: %w", err)
	}
	if gasRequired > int64(p.cfg.GasLimit) {
		return nil, fmt.Errorf("gas requirement is too high: %d", gasRequired)
	}
	if !dryRunResp.Effects.Data.IsSuccess() {
		return nil, fmt.Errorf(dryRunResp.Effects.Data.V1.Status.Error)
	}
	signature, err := wallet.SignSecureWithoutEncode(txnMetadata.TxBytes, sui_types.DefaultIntent())
	if err != nil {
		return nil, err
	}
	signatures := []any{signature}
	txnResp, err := p.client.CommitTx(ctx, wallet, txnMetadata.TxBytes, signatures)
	return txnResp, err
}

func (p *Provider) executeRouteCallBack(txRes types.SuiTransactionBlockResponse, messageKey *relayertypes.MessageKey, method string, callback relayertypes.TxResponseFunc, err error) {
	// if error occurred before txn processing
	if err != nil || txRes.Digest == nil {
		if err == nil {
			err = fmt.Errorf("txn execution failed; received empty tx digest")
		}
		callback(messageKey, nil, err)
		p.log.Error("failed to execute transaction", zap.Error(err))
		return
	}

	res := &relayertypes.TxResponse{
		TxHash: txRes.Digest.String(),
	}

	txnData, err := p.client.GetTransaction(context.Background(), txRes.Digest.String())
	if err != nil {
		callback(messageKey, res, err)
		p.log.Error("failed to execute transaction", zap.Error(err), zap.String("tx_hash", txRes.Digest.String()))
		return
	}

	// assign tx successful height
	res.Height = txnData.Checkpoint.Int64()
	success := txRes.Effects.Data.IsSuccess()
	if !success {
		err = fmt.Errorf("error: %s", txRes.Effects.Data.V1.Status.Error)
		callback(messageKey, res, err)
		p.log.Info("failed transaction",
			zap.Any("message-key", messageKey),
			zap.String("tx_hash", txRes.Digest.String()),
			zap.Int64("height", txnData.Checkpoint.Int64()),
			zap.Error(err),
		)
		return
	}
	res.Code = relayertypes.Success
	callback(messageKey, res, nil)
	p.log.Info("successful transaction",
		zap.Any("message-key", messageKey),
		zap.String("tx_hash", txRes.Digest.String()),
		zap.Int64("height", txnData.Checkpoint.Int64()),
	)
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txDigest string) (*relayertypes.Receipt, error) {
	txBlock, err := p.client.GetTransaction(ctx, txDigest)
	if err != nil {
		return nil, err
	}
	receipt := &relayertypes.Receipt{
		TxHash: txDigest,
		Height: txBlock.Checkpoint.Uint64(),
		Status: txBlock.Effects.Data.IsSuccess(),
	}
	return receipt, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayertypes.MessageKey) (bool, error) {
	snU128, err := bcs.NewUint128FromBigInt(bcs.NewBigIntFromUint64(key.Sn))
	if err != nil {
		return false, err
	}
	suiMessage := p.NewSuiMessage([]SuiCallArg{
		{Type: CallArgObject, Val: p.cfg.XcallStorageID},
		{Type: CallArgPure, Val: key.Src},
		{Type: CallArgPure, Val: snU128},
	}, p.cfg.XcallPkgID, EntryModule, MethodGetReceipt)
	var msgReceived bool
	wallet, err := p.Wallet()
	if err != nil {
		return msgReceived, err
	}
	if err := p.client.QueryContract(ctx, suiMessage, wallet.Address, p.cfg.GasLimit, &msgReceived); err != nil {
		return msgReceived, err
	}
	return msgReceived, nil
}
