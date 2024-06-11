package sui

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"cosmossdk.io/errors"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/move_types"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/fardream/go-bcs/bcs"
	"github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/hexstr"
	"go.uber.org/zap"
)

func (p *Provider) Route(ctx context.Context, message *relayertypes.Message, callback relayertypes.TxResponseFunc) error {
	p.log.Info("starting to route message",
		zap.String("src", message.Src),
		zap.String("event-type", message.EventType),
		zap.String("data", hex.EncodeToString(message.Data)))

	suiMessage, err := p.MakeSuiMessage(message)
	if err != nil {
		return err
	}

	p.txmut.Lock()
	defer p.txmut.Unlock()

	txBytes, err := p.prepareTxMoveCall(suiMessage)
	if err != nil {
		return err
	}

	txRes, err := p.SendTransaction(ctx, txBytes)
	go p.executeRouteCallBack(txRes, message.MessageKey(), suiMessage.Method, callback, err)
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
			{Type: CallArgObject, Val: p.cfg.ConnectionCapID},
			{Type: CallArgPure, Val: message.Src},
			{Type: CallArgPure, Val: snU128},
			{Type: CallArgPure, Val: "0x" + hex.EncodeToString(message.Data)},
		}
		return p.NewSuiMessage([]string{}, callParams, p.xcallPkgIDLatest(), ModuleEntry, MethodRecvMessage), nil
	case events.CallMessage:
		if _, err := p.Wallet(); err != nil {
			return nil, err
		}

		coins, err := p.client.GetCoins(context.Background(), p.wallet.Address)
		if err != nil {
			return nil, err
		}
		_, coin, err := coins.PickSUICoinsWithGas(big.NewInt(int64(suiBaseFee)), p.cfg.GasLimit, types.PickBigger)
		if err != nil {
			return nil, err
		}

		module, err := p.getModule(func(mod DappModule) bool {
			return hexstr.NewFromString(mod.CapID) == hexstr.NewFromString(message.DappModuleCapID)
		})
		if err != nil {
			return nil, err
		}

		var callParams []SuiCallArg
		typeArgs := []string{}

		switch module.Name {
		case ModuleMockDapp:
			callParams = []SuiCallArg{
				{Type: CallArgObject, Val: module.ConfigID},
				{Type: CallArgObject, Val: p.cfg.XcallStorageID},
				{Type: CallArgObject, Val: coin.CoinObjectId.String()},
				{Type: CallArgPure, Val: strconv.Itoa(int(message.ReqID))},
				{Type: CallArgPure, Val: "0x" + hex.EncodeToString(message.Data)},
			}
		case ModuleXcallManager:
			callParams = []SuiCallArg{
				{Type: CallArgObject, Val: module.ConfigID},
				{Type: CallArgObject, Val: p.cfg.XcallStorageID},
				{Type: CallArgObject, Val: coin.CoinObjectId.String()},
				{Type: CallArgPure, Val: strconv.Itoa(int(message.ReqID))},
				{Type: CallArgPure, Val: "0x" + hex.EncodeToString(message.Data)},
			}
		case ModuleAssetManager:
			xcallManagerModule, err := p.getModule(func(mod DappModule) bool {
				return mod.Name == ModuleXcallManager
			})
			if err != nil {
				return nil, fmt.Errorf("failed to find xcall manager module: %w", err)
			}

			withdrawTokenType, err := p.getWithdrawTokentype(context.Background(), message)
			if err != nil {
				return nil, fmt.Errorf("failed to get withdraw token type: %w", err)
			}

			typeArgs = append(typeArgs, *withdrawTokenType)

			callParams = []SuiCallArg{
				{Type: CallArgObject, Val: module.ConfigID},
				{Type: CallArgObject, Val: xcallManagerModule.ConfigID},
				{Type: CallArgObject, Val: p.cfg.XcallStorageID},
				{Type: CallArgObject, Val: coin.CoinObjectId.String()},
				{Type: CallArgObject, Val: suiClockObjectId},
				{Type: CallArgPure, Val: strconv.Itoa(int(message.ReqID))},
				{Type: CallArgPure, Val: "0x" + hex.EncodeToString(message.Data)},
			}
		case ModuleBalancedDollar:
			xcallManagerModule, err := p.getModule(func(mod DappModule) bool {
				return mod.Name == ModuleXcallManager
			})
			if err != nil {
				return nil, fmt.Errorf("failed to find xcall manager module: %w", err)
			}
			callParams = []SuiCallArg{
				{Type: CallArgObject, Val: module.ConfigID},
				{Type: CallArgObject, Val: xcallManagerModule.ConfigID},
				{Type: CallArgObject, Val: p.cfg.XcallStorageID},
				{Type: CallArgObject, Val: coin.CoinObjectId.String()},
				{Type: CallArgPure, Val: strconv.Itoa(int(message.ReqID))},
				{Type: CallArgPure, Val: "0x" + hex.EncodeToString(message.Data)},
			}

		default:
			return nil, fmt.Errorf("received unknown dapp module cap id: %s", message.DappModuleCapID)
		}

		return p.NewSuiMessage(typeArgs, callParams, p.cfg.DappPkgID, module.Name, MethodExecuteCall), nil

	case events.ExecuteRollback:
		module, err := p.getModule(func(mod DappModule) bool {
			return hexstr.NewFromString(mod.CapID) == hexstr.NewFromString(message.DappModuleCapID)
		})
		if err != nil {
			return nil, err
		}

		snU128, err := bcs.NewUint128FromBigInt(bcs.NewBigIntFromUint64(message.Sn))
		if err != nil {
			return nil, err
		}

		var callParams []SuiCallArg
		typeArgs := []string{}

		switch module.Name {
		case ModuleMockDapp:
			callParams = []SuiCallArg{
				{Type: CallArgObject, Val: module.ConfigID},
				{Type: CallArgObject, Val: p.cfg.XcallStorageID},
				{Type: CallArgPure, Val: snU128},
			}
		case ModuleAssetManager:
			withdrawTokenType, err := p.getWithdrawTokentype(context.Background(), message)
			if err != nil {
				return nil, fmt.Errorf("failed to get withdraw token type: %w", err)
			}

			typeArgs = append(typeArgs, *withdrawTokenType)

			callParams = []SuiCallArg{
				{Type: CallArgObject, Val: module.ConfigID},
				{Type: CallArgObject, Val: p.cfg.XcallStorageID},
				{Type: CallArgPure, Val: snU128},
				{Type: CallArgObject, Val: suiClockObjectId},
			}
		case ModuleBalancedDollar:
			callParams = []SuiCallArg{
				{Type: CallArgObject, Val: module.ConfigID},
				{Type: CallArgObject, Val: p.cfg.XcallStorageID},
				{Type: CallArgPure, Val: snU128},
			}

		default:
			return nil, fmt.Errorf("received unknown dapp module cap id: %s", message.DappModuleCapID)
		}

		return p.NewSuiMessage(typeArgs, callParams, p.cfg.DappPkgID, module.Name, MethodExecuteRollback), nil

	default:
		return nil, fmt.Errorf("can't generate message for unknown event type: %s ", message.EventType)
	}
}

func (p *Provider) getModule(condition func(module DappModule) bool) (*DappModule, error) {
	for _, mod := range p.cfg.DappModules {
		if condition(mod) {
			return &mod, nil
		}
	}
	return nil, fmt.Errorf("module not found")
}

func (p *Provider) preparePTB(msg *SuiMessage) (lib.Base64Data, error) {
	builder := sui_types.NewProgrammableTransactionBuilder()
	packageId, err := move_types.NewAccountAddressHex(msg.PackageId)
	if err != nil {
		return nil, err
	}

	callArgs, err := p.paramsToCallArgs(msg.Params)
	if err != nil {
		return nil, err
	}

	err = builder.MoveCall(
		*packageId,
		move_types.Identifier(msg.Module),
		move_types.Identifier(msg.Method),
		[]move_types.TypeTag{},
		callArgs,
	)
	if err != nil {
		return nil, err
	}
	bcsBytes, err := bcs.Marshal(builder.Finish())
	if err != nil {
		return nil, err
	}
	txBytes := append([]byte{0}, bcsBytes...)
	b64Data, err := lib.NewBase64Data(base64.StdEncoding.EncodeToString(txBytes))
	if err != nil {
		return nil, err
	}
	return *b64Data, nil
}

func (p *Provider) prepareTxMoveCall(msg *SuiMessage) (lib.Base64Data, error) {
	if _, err := p.Wallet(); err != nil {
		return nil, err
	}
	accountAddress, err := move_types.NewAccountAddressHex(p.wallet.Address)
	if err != nil {
		return nil, fmt.Errorf("error getting account address sender: %w", err)
	}
	packageId, err := move_types.NewAccountAddressHex(msg.PackageId)
	if err != nil {
		return nil, fmt.Errorf("invalid packageId: %w", err)
	}

	var args []interface{}
	for _, param := range msg.Params {
		args = append(args, param.Val)
	}

	res, err := p.client.MoveCall(
		context.Background(),
		*accountAddress,
		*packageId,
		msg.Module,
		msg.Method,
		msg.TypeArgs,
		args,
		nil,
		types.NewSafeSuiBigInt(p.cfg.GasLimit),
	)
	if err != nil {
		return nil, err
	}
	return res.TxBytes, nil
}

func (p *Provider) paramsToCallArgs(params []SuiCallArg) ([]sui_types.CallArg, error) {
	var callArgs []sui_types.CallArg
	for _, param := range params {
		switch param.Type {
		case CallArgObject:
			arg, err := p.getCallArgObject(param.Val.(string))
			if err != nil {
				return nil, err
			}
			callArgs = append(callArgs, *arg)
		case CallArgPure:
			arg, err := p.getCallArgPure(param.Val)
			if err != nil {
				return nil, err
			}
			callArgs = append(callArgs, *arg)
		default:
			return nil, fmt.Errorf("invalid call arg type")
		}
	}
	return callArgs, nil
}

func (p *Provider) getCallArgPure(arg interface{}) (*sui_types.CallArg, error) {
	byteParam, err := bcs.Marshal(arg)
	if err != nil {
		return nil, err
	}
	return &sui_types.CallArg{
		Pure: &byteParam,
	}, nil
}

func (p *Provider) getCallArgObject(arg string) (*sui_types.CallArg, error) {
	objectId, err := sui_types.NewAddressFromHex(arg)
	if err != nil {
		return nil, err
	}
	object, err := p.client.GetObject(context.Background(), *objectId, &types.SuiObjectDataOptions{
		ShowType:  true,
		ShowOwner: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %v", err)
	}

	if object.Data.Owner != nil && object.Data.Owner.Shared != nil {
		return &sui_types.CallArg{
			Object: &sui_types.ObjectArg{
				SharedObject: &struct {
					Id                   sui_types.ObjectID
					InitialSharedVersion sui_types.SequenceNumber
					Mutable              bool
				}{
					Id:                   object.Data.ObjectId,
					InitialSharedVersion: *object.Data.Owner.Shared.InitialSharedVersion,
					Mutable:              true,
				},
			},
		}, nil
	}

	objRef := object.Data.Reference()

	return &sui_types.CallArg{
		Object: &sui_types.ObjectArg{
			ImmOrOwnedObject: &objRef,
		},
	}, nil
}

func (p *Provider) SendTransaction(ctx context.Context, txBytes lib.Base64Data) (*types.SuiTransactionBlockResponse, error) {
	wallet, err := p.Wallet()
	if err != nil {
		return nil, err
	}

	dryRunResp, gasRequired, err := p.client.SimulateTx(ctx, txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed simulating tx: %w", err)
	}
	if gasRequired > int64(p.cfg.GasLimit) {
		return nil, fmt.Errorf("gas requirement is too high: %d", gasRequired)
	}
	if !dryRunResp.Effects.Data.IsSuccess() {
		return nil, fmt.Errorf(dryRunResp.Effects.Data.V1.Status.Error)
	}
	signature, err := wallet.SignSecureWithoutEncode(txBytes, sui_types.DefaultIntent())
	if err != nil {
		return nil, err
	}
	signatures := []any{signature}
	txnResp, err := p.client.ExecuteTx(ctx, wallet, txBytes, signatures)

	return txnResp, err
}

func (p *Provider) executeRouteCallBack(txRes *types.SuiTransactionBlockResponse, messageKey *relayertypes.MessageKey, method string, callback relayertypes.TxResponseFunc, err error) {
	// if error occurred before txn processing
	if err != nil || txRes == nil || txRes.Digest == nil {
		if err == nil {
			err = fmt.Errorf("txn execution failed; received empty tx digest")
		}
		callback(messageKey, &relayertypes.TxResponse{}, err)
		p.log.Error("failed to execute transaction", zap.Error(err), zap.String("method", method))
		return
	}

	res := &relayertypes.TxResponse{
		TxHash: txRes.Digest.String(),
	}

	txnData, err := p.client.GetTransaction(context.Background(), txRes.Digest.String())
	if err != nil {
		callback(messageKey, res, err)
		p.log.Error("failed to execute transaction", zap.Error(err), zap.String("method", method), zap.String("tx_hash", txRes.Digest.String()))
		return
	}

	if txnData.Checkpoint == nil {
		time.Sleep(1 * time.Second) //time to wait until tx is included in some checkpoint
		txnData, err = p.client.GetTransaction(context.Background(), txRes.Digest.String())
		if err != nil {
			callback(messageKey, res, err)
			p.log.Error("failed to execute transaction", zap.Error(err), zap.String("method", method), zap.String("tx_hash", txRes.Digest.String()))
			return
		}
	}

	// assign tx successful height
	res.Height = txnData.Checkpoint.Int64()
	success := txRes.Effects.Data.IsSuccess()
	if !success {
		err = fmt.Errorf("error: %s", txRes.Effects.Data.V1.Status.Error)
		callback(messageKey, res, err)
		p.log.Info("failed transaction",
			zap.Any("message-key", messageKey),
			zap.String("method", method),
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
		zap.String("method", method),
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
	suiMessage := p.NewSuiMessage(
		[]string{},
		[]SuiCallArg{
			{Type: CallArgObject, Val: p.cfg.XcallStorageID},
			{Type: CallArgPure, Val: p.cfg.ConnectionID},
			{Type: CallArgPure, Val: key.Src},
			{Type: CallArgPure, Val: snU128},
		}, p.xcallPkgIDLatest(), ModuleEntry, MethodGetReceipt)
	var msgReceived bool
	wallet, err := p.Wallet()
	if err != nil {
		return msgReceived, err
	}

	txBytes, err := p.preparePTB(suiMessage)
	if err != nil {
		return msgReceived, err
	}

	if err := p.client.QueryContract(ctx, wallet.Address, txBytes, &msgReceived); err != nil {
		return msgReceived, err
	}

	return msgReceived, nil
}

func (p *Provider) getWithdrawTokentype(ctx context.Context, message *relayertypes.Message) (*string, error) {
	suiMessage := p.NewSuiMessage(
		[]string{},
		[]SuiCallArg{
			{Type: CallArgPure, Val: message.Data},
		}, p.cfg.DappPkgID, ModuleAssetManager, MethodGetWithdrawTokentype)
	var tokenType string
	wallet, err := p.Wallet()
	if err != nil {
		return nil, err
	}

	txBytes, err := p.preparePTB(suiMessage)
	if err != nil {
		return nil, err
	}

	if err := p.client.QueryContract(ctx, wallet.Address, txBytes, &tokenType); err != nil {
		return nil, err
	}

	return &tokenType, nil
}
