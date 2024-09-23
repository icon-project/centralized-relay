package sui

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
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
		zap.Any("sn", message.Sn),
		zap.Any("req_id", message.ReqID),
		zap.String("src", message.Src),
		zap.String("event_type", message.EventType))

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
		snU128, err := bcs.NewUint128FromBigInt(message.Sn)
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

		typeArgs, _, err := p.getRecvParams(context.Background(), message, p.cfg.XcallPkgID, p.cfg.ConnectionModule)
		if err == nil {
			callParams = []SuiCallArg{
				{Type: CallArgObject, Val: p.cfg.XcallStorageID},
				{Type: CallArgPure, Val: message.Src},
				{Type: CallArgPure, Val: snU128},
				{Type: CallArgPure, Val: "0x" + hex.EncodeToString(message.Data)},
			}
		} else {
			typeArgs = []string{}
		}

		return p.NewSuiMessage(typeArgs, callParams, p.cfg.XcallPkgID, p.cfg.ConnectionModule, MethodRecvMessage), nil
	case events.CallMessage:
		if _, err := p.Wallet(); err != nil {
			return nil, err
		}

		dapp, module, err := p.getModule(message.DappModuleCapID)
		if err != nil {
			return nil, fmt.Errorf("module cap id %s not found: %w", message.DappModuleCapID, err)
		}

		typeArgs, callArgs, err := p.getExecuteParams(context.Background(), message, dapp, module, MethodGetExecuteCallParams)
		if err != nil {
			return nil, fmt.Errorf("failed to get execute params: %w", err)
		}

		return p.NewSuiMessage(typeArgs, callArgs, dapp.PkgID, module.Name, MethodExecuteCall), nil

	case events.RollbackMessage:
		dapp, module, err := p.getModule(message.DappModuleCapID)
		if err != nil {
			return nil, fmt.Errorf("failed to get module cap id %s: %w", message.DappModuleCapID, err)
		}
		typeArgs, callArgs, err := p.getExecuteParams(context.Background(), message, dapp, module, MethodGetExecuteRollbackParams)
		if err != nil {
			return nil, fmt.Errorf("failed to get execute params: %w", err)
		}
		return p.NewSuiMessage(typeArgs, callArgs, dapp.PkgID, module.Name, MethodExecuteRollback), nil

	default:
		return nil, fmt.Errorf("can't generate message for unknown event type: %s ", message.EventType)
	}
}

func (p *Provider) getModule(moduleCapID string) (*Dapp, *DappModule, error) {
	for _, dapp := range p.cfg.Dapps {
		for _, mod := range dapp.Modules {
			if hexstr.NewFromString(mod.CapID) == hexstr.NewFromString(moduleCapID) {
				return &dapp, &mod, nil
			}
		}
	}

	return nil, nil, fmt.Errorf("module not found")
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
		p.log.Error("failed to get transaction details after execution", zap.Error(err), zap.String("method", method), zap.String("tx_hash", txRes.Digest.String()))
		return
	}

	for txnData.Checkpoint == nil {
		p.log.Warn("transaction not included in checkpoint", zap.String("tx-digest", txRes.Digest.String()))
		time.Sleep(3 * time.Second) //time to wait until tx is included in some checkpoint
		txnData, err = p.client.GetTransaction(context.Background(), txRes.Digest.String())
		if err != nil {
			callback(messageKey, res, err)
			p.log.Error("failed to get transaction details due to nil checkpoint after execution", zap.Error(err), zap.String("method", method), zap.String("tx_hash", txRes.Digest.String()))
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

func (p *Provider) MessageReceived(ctx context.Context, messageKey *relayertypes.MessageKey) (bool, error) {
	switch messageKey.EventType {
	case events.EmitMessage:
		snU128, err := bcs.NewUint128FromBigInt(messageKey.Sn)
		if err != nil {
			return false, err
		}
		suiMessage := p.NewSuiMessage(
			[]string{},
			[]SuiCallArg{
				{Type: CallArgObject, Val: p.cfg.XcallStorageID},
				{Type: CallArgPure, Val: p.cfg.ConnectionID},
				{Type: CallArgPure, Val: messageKey.Src},
				{Type: CallArgPure, Val: snU128},
			}, p.cfg.XcallPkgID, p.cfg.ConnectionModule, MethodGetReceipt)
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
	case events.CallMessage:
		return false, nil
	case events.RollbackMessage:
		return false, nil
	default:
		return true, fmt.Errorf("unknown event type")
	}
}

func (p *Provider) getExecuteParams(
	ctx context.Context,
	message *relayertypes.Message,
	dapp *Dapp,
	dappModule *DappModule,
	method string,
) (typeArgs []string, callArgs []SuiCallArg, err error) {
	suiMessage := p.NewSuiMessage(
		[]string{},
		[]SuiCallArg{
			{Type: CallArgObject, Val: dappModule.ConfigID},
			{Type: CallArgPure, Val: message.Data},
		},
		dapp.PkgID,
		dappModule.Name,
		method,
	)

	args := struct {
		TypeArgs []string
		Args     []string
	}{}

	wallet, err := p.Wallet()
	if err != nil {
		return nil, nil, err
	}

	txBytes, err := p.preparePTB(suiMessage)
	if err != nil {
		return nil, nil, err
	}

	if err := p.client.QueryContract(ctx, wallet.Address, txBytes, &args); err != nil {
		return nil, nil, err
	}

	typeArgs = args.TypeArgs

	for _, arg := range args.Args {
		if strings.HasPrefix(arg, "0x") {
			callArgs = append(callArgs, SuiCallArg{
				Type: CallArgObject, Val: arg,
			})
		} else {
			switch arg {
			case "coin":
				coins, err := p.client.GetCoins(context.Background(), p.wallet.Address)
				if err != nil {
					return nil, nil, err
				}
				_, coin, err := coins.PickSUICoinsWithGas(big.NewInt(int64(suiBaseFee)), p.cfg.GasLimit, types.PickBigger)
				if err != nil {
					return nil, nil, err
				}
				callArgs = append(callArgs, SuiCallArg{
					Type: CallArgObject, Val: coin.CoinObjectId.String(),
				})
			case "sn":
				snU128, err := bcs.NewUint128FromBigInt(message.Sn)
				if err != nil {
					return nil, nil, err
				}
				callArgs = append(callArgs, SuiCallArg{
					Type: CallArgPure, Val: snU128,
				})
			case "request_id":
				callArgs = append(callArgs, SuiCallArg{
					Type: CallArgPure, Val: strconv.Itoa(int(message.ReqID.Int64())),
				})
			case "data":
				callArgs = append(callArgs, SuiCallArg{
					Type: CallArgPure, Val: "0x" + hex.EncodeToString(message.Data),
				})
			}
			if keyObjID, ok := dapp.Constants[arg]; ok {
				callArgs = append(callArgs, SuiCallArg{
					Type: CallArgObject, Val: keyObjID,
				})
			}
		}
	}

	return
}

func (p *Provider) getRecvParams(
	ctx context.Context,
	message *relayertypes.Message,
	pkgID string,
	module string,
) (typeArgs []string, callArgs []SuiCallArg, err error) {
	suiMessage := p.NewSuiMessage(
		[]string{},
		[]SuiCallArg{
			{Type: CallArgObject, Val: p.cfg.XcallStorageID},
			{Type: CallArgPure, Val: message.Data},
		},
		pkgID,
		module,
		"get_receive_msg_args",
	)

	args := struct {
		TypeArgs []string
		Args     []string
	}{}

	wallet, err := p.Wallet()
	if err != nil {
		return nil, nil, err
	}

	txBytes, err := p.preparePTB(suiMessage)
	if err != nil {
		return nil, nil, err
	}

	if err := p.client.QueryContract(ctx, wallet.Address, txBytes, &args); err != nil {
		return nil, nil, err
	}

	typeArgs = args.TypeArgs

	return
}
