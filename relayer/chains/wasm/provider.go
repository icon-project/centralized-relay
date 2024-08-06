package wasm

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"runtime"
	"strings"
	"sync"
	"time"

	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	coreTypes "github.com/cometbft/cometbft/rpc/core/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

var _ provider.ChainProvider = (*Provider)(nil)

type Provider struct {
	logger              *zap.Logger
	cfg                 *Config
	client              IClient
	kms                 kms.KMS
	wallet              sdkTypes.AccountI
	contracts           map[string]relayTypes.EventMap
	eventList           []sdkTypes.Event
	LastSavedHeightFunc func() uint64
	routerMutex         *sync.Mutex
}

func (p *Provider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	return p.client.GetLatestBlockHeight(ctx)
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txHash string) (*relayTypes.Receipt, error) {
	res, err := p.client.GetTransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}
	return &relayTypes.Receipt{
		TxHash: txHash,
		Height: uint64(res.TxResponse.Height),
		Status: types.CodeTypeOK == res.TxResponse.Code,
	}, nil
}

func (p *Provider) NID() string {
	return p.cfg.NID
}

func (p *Provider) Name() string {
	return p.cfg.ChainName
}

func (p *Provider) Init(ctx context.Context, homePath string, kms kms.KMS) error {
	if err := p.cfg.Contracts.Validate(); err != nil {
		return err
	}
	p.kms = kms
	return nil
}

// Wallet returns the wallet of the provider
func (p *Provider) Wallet() sdkTypes.AccAddress {
	ctx := context.Background()
	done := p.SetSDKContext()
	defer done()
	if p.wallet == nil || p.wallet.GetAddress().Empty() {
		if err := p.RestoreKeystore(ctx); err != nil {
			p.logger.Error("failed to restore keystore", zap.Error(err))
			return nil
		}
		account, err := p.client.GetAccountInfo(ctx, p.cfg.GetWallet())
		if err != nil {
			p.logger.Error("failed to get account info", zap.Error(err))
			return nil
		}
		p.wallet = account
		return p.client.SetAddress(account.GetAddress())
	}
	return p.wallet.GetAddress()
}

func (p *Provider) Type() string {
	return types.ChainType
}

func (p *Provider) Config() provider.Config {
	return p.cfg
}

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfoChan chan *relayTypes.BlockInfo) error {
	latestHeight, err := p.QueryLatestHeight(ctx)
	if err != nil {
		p.logger.Error("failed to get latest block height", zap.Error(err))
		return err
	}

	startHeight, err := p.getStartHeight(latestHeight, lastSavedHeight)
	if err != nil {
		p.logger.Error("failed to determine start height", zap.Error(err))
		return err
	}

	subscribeStarter := time.NewTicker(time.Second * 1)
	pollHeightTicker := time.NewTicker(time.Second * 1)
	pollHeightTicker.Stop()

	resetFunc := func() {
		subscribeStarter.Reset(time.Second * 3)
		pollHeightTicker.Reset(time.Second * 2)
	}

	p.logger.Info("Start from height", zap.Uint64("height", startHeight), zap.Uint64("finality block", p.FinalityBlock(ctx)))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-subscribeStarter.C:
			subscribeStarter.Stop()
			for _, event := range p.contracts {
				for msgType := range event.SigType {
					go p.SubscribeMessageEvents(ctx, blockInfoChan, &types.SubscribeOpts{
						Address: event.Address,
						Method:  msgType,
						Height:  latestHeight,
					}, resetFunc)
				}
			}
			if startHeight < latestHeight {
				p.logger.Info("Syncing", zap.Uint64("from-height", startHeight),
					zap.Uint64("to-height", latestHeight), zap.Uint64("delta", latestHeight-startHeight))
				startHeight = p.runBlockQuery(ctx, blockInfoChan, startHeight, latestHeight)
			}
		case <-pollHeightTicker.C:
			pollHeightTicker.Stop()
			startHeight = p.GetLastSavedHeight()
			if startHeight == 0 {
				startHeight = latestHeight
			}
			latestHeight, err = p.QueryLatestHeight(ctx)
			if err != nil {
				p.logger.Error("failed to get latest block height", zap.Error(err))
				pollHeightTicker.Reset(time.Second * 3)
			}
		}
	}
}

func (p *Provider) Route(ctx context.Context, message *relayTypes.Message, callback relayTypes.TxResponseFunc) error {
	p.logger.Info("starting to route message", zap.Any("message", message))
	res, err := p.call(ctx, message)
	if err != nil {
		return err
	}
	seq := p.wallet.GetSequence() + 1
	if err := p.wallet.SetSequence(seq); err != nil {
		p.logger.Error("failed to set sequence", zap.Error(err))
	}
	return p.waitForTxResult(ctx, message.MessageKey(), res, callback)
}

// call the smart contract to send the message
func (p *Provider) call(ctx context.Context, message *relayTypes.Message) (*sdkTypes.TxResponse, error) {
	rawMsg, err := p.getRawContractMessage(message)
	if err != nil {
		return nil, err
	}

	var contract string

	switch message.EventType {
	case events.EmitMessage, events.RevertMessage, events.SetAdmin, events.ClaimFee, events.SetFee:
		contract = p.cfg.Contracts[relayTypes.ConnectionContract]
	case events.CallMessage, events.RollbackMessage:
		contract = p.cfg.Contracts[relayTypes.XcallContract]
	default:
		return nil, fmt.Errorf("unknown event type: %s ", message.EventType)
	}

	msg := &wasmTypes.MsgExecuteContract{
		Sender:   p.Wallet().String(),
		Contract: contract,
		Msg:      rawMsg,
	}

	msgs := []sdkTypes.Msg{msg}

	res, err := p.sendMessage(ctx, msgs...)
	if err != nil {
		if strings.Contains(err.Error(), errors.ErrWrongSequence.Error()) {
			if mmErr := p.handleSequence(ctx); mmErr != nil {
				return res, fmt.Errorf("failed to handle sequence mismatch error: %v || %v", mmErr, err)
			}
			return p.sendMessage(ctx, msgs...)
		}
	}
	return res, err
}

func (p *Provider) sendMessage(ctx context.Context, msgs ...sdkTypes.Msg) (*sdkTypes.TxResponse, error) {
	p.routerMutex.Lock()
	defer p.routerMutex.Unlock()
	return p.prepareAndPushTxToMemPool(ctx, p.wallet.GetAccountNumber(), p.wallet.GetSequence(), msgs...)
}

func (p *Provider) handleSequence(ctx context.Context) error {
	acc, err := p.client.GetAccountInfo(ctx, p.Wallet().String())
	if err != nil {
		return err
	}
	return p.wallet.SetSequence(acc.GetSequence())
}

func (p *Provider) logTxFailed(err error, tx *sdkTypes.TxResponse) {
	p.logger.Error("transaction failed",
		zap.String("tx_hash", tx.TxHash),
		zap.String("codespace", tx.Codespace),
		zap.Error(err),
	)
}

func (p *Provider) logTxSuccess(res *types.TxResult) {
	p.logger.Info("successful transaction",
		zap.Int64("block_height", res.TxResult.Height),
		zap.String("tx_hash", res.TxResult.TxHash),
	)
}

func (p *Provider) prepareAndPushTxToMemPool(ctx context.Context, acc, seq uint64, msgs ...sdkTypes.Msg) (*sdkTypes.TxResponse, error) {
	txf, err := p.client.BuildTxFactory()
	if err != nil {
		return nil, err
	}

	txf = txf.
		WithGasPrices(p.cfg.GasPrices).
		WithGasAdjustment(p.cfg.GasAdjustment).
		WithAccountNumber(acc).
		WithSequence(seq)

	if txf.SimulateAndExecute() {
		_, adjusted, err := p.client.EstimateGas(txf, msgs...)
		if err != nil {
			return nil, err
		}
		txf = txf.WithGas(adjusted)
	}

	if txf.Gas() < p.cfg.MinGasAmount {
		return nil, fmt.Errorf("gas amount %d is too low; the minimum allowed gas amount is %d", txf.Gas(), p.cfg.MinGasAmount)
	}

	if txf.Gas() > p.cfg.MaxGasAmount {
		return nil, fmt.Errorf("gas amount %d exceeds the maximum allowed limit of %d", txf.Gas(), p.cfg.MaxGasAmount)
	}

	txBytes, err := p.client.PrepareTx(ctx, txf, msgs...)
	if err != nil {
		return nil, err
	}

	res, err := p.client.BroadcastTx(txBytes)
	if err != nil || res.Code != types.CodeTypeOK {
		if err == nil {
			err = fmt.Errorf("failed to send tx: %v", res.RawLog)
		}
		return nil, err
	}

	return res, nil
}

func (p *Provider) waitForTxResult(ctx context.Context, mk *relayTypes.MessageKey, tx *sdkTypes.TxResponse, callback relayTypes.TxResponseFunc) error {
	res, err := p.subscribeTxResult(ctx, tx, p.cfg.TxConfirmationInterval)
	if err != nil {
		callback(mk, res.TxResult, err)
		p.logTxFailed(err, tx)
		return err
	}
	callback(mk, res.TxResult, nil)
	p.logTxSuccess(res)
	return nil
}

func (p *Provider) pollTxResultStream(ctx context.Context, txHash string, maxWaitInterval time.Duration) <-chan *types.TxResult {
	txResChan := make(chan *types.TxResult)
	startTime := time.Now()
	go func(txChan chan *types.TxResult) {
		defer close(txChan)
		for range time.NewTicker(p.cfg.TxConfirmationInterval).C {
			res, err := p.client.GetTransactionReceipt(ctx, txHash)
			if err == nil {
				txChan <- &types.TxResult{
					TxResult: &relayTypes.TxResponse{
						Height:    res.TxResponse.Height,
						TxHash:    res.TxResponse.TxHash,
						Codespace: res.TxResponse.Codespace,
						Code:      relayTypes.ResponseCode(res.TxResponse.Code),
						Data:      res.TxResponse.Data,
					},
				}
				return
			} else if time.Since(startTime) > maxWaitInterval {
				txChan <- &types.TxResult{
					Error: err,
				}
				return
			}
		}
	}(txResChan)
	return txResChan
}

func (p *Provider) subscribeTxResult(ctx context.Context, tx *sdkTypes.TxResponse, maxWaitInterval time.Duration) (*types.TxResult, error) {
	newCtx, cancel := context.WithTimeout(ctx, maxWaitInterval)
	defer cancel()

	query := fmt.Sprintf("tm.event = 'Tx' AND tx.hash = '%s'", tx.TxHash)
	resultEventChan, err := p.client.Subscribe(newCtx, "tx-result-waiter", query)
	if err != nil {
		return &types.TxResult{
			Error: err,
			TxResult: &relayTypes.TxResponse{
				TxHash: tx.TxHash,
			},
		}, err
	}
	defer p.client.Unsubscribe(newCtx, "tx-result-waiter", query)

	for {
		select {
		case <-ctx.Done():
			return &types.TxResult{TxResult: &relayTypes.TxResponse{TxHash: tx.TxHash}}, ctx.Err()
		case e := <-resultEventChan:
			eventDataJSON, err := jsoniter.Marshal(e.Data)
			if err != nil {
				return &types.TxResult{TxResult: &relayTypes.TxResponse{TxHash: tx.TxHash}}, err
			}

			txRes := new(types.TxResultResponse)
			if err := jsoniter.Unmarshal(eventDataJSON, txRes); err != nil {
				return &types.TxResult{TxResult: &relayTypes.TxResponse{TxHash: tx.TxHash}}, err
			}

			res := &types.TxResult{
				TxResult: &relayTypes.TxResponse{
					Height:    txRes.Height,
					TxHash:    tx.TxHash,
					Codespace: txRes.Result.Codespace,
					Data:      string(txRes.Result.Data),
				},
			}
			if uint32(txRes.Result.Code) != types.CodeTypeOK {
				return res, fmt.Errorf("transaction failed with error: %v", txRes.Result.Log)
			}
			res.TxResult.Code = relayTypes.Success
			return res, nil
		}
	}
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayTypes.MessageKey) (bool, error) {
	switch key.EventType {
	case events.EmitMessage:
		queryMsg := &types.QueryReceiptMsg{
			GetReceipt: &types.GetReceiptMsg{
				SrcNetwork: key.Src,
				ConnSn:     key.Sn.String(),
			},
		}
		rawQueryMsg, err := jsoniter.Marshal(queryMsg)
		if err != nil {
			return false, err
		}

		res, err := p.client.QuerySmartContract(ctx, p.cfg.Contracts[relayTypes.ConnectionContract], rawQueryMsg)
		if err != nil {
			p.logger.Error("failed to check if message is received: ", zap.Error(err))
			return false, err
		}

		receiptMsgRes := types.QueryReceiptMsgResponse{}
		return receiptMsgRes.Status, jsoniter.Unmarshal(res.Data, &receiptMsgRes.Status)
	case events.CallMessage:
		return false, nil
	case events.RollbackMessage:
		return false, nil
	default:
		return true, fmt.Errorf("unknown event type")
	}
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*relayTypes.Coin, error) {
	coin, err := p.client.GetBalance(ctx, addr, p.cfg.Denomination)
	if err != nil {
		p.logger.Error("failed to query balance: ", zap.Error(err))
		return nil, err
	}
	return &relayTypes.Coin{
		Denom:  coin.Denom,
		Amount: coin.Amount.BigInt().Uint64(),
	}, nil
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, message *relayTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, message *relayTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) GenerateMessages(ctx context.Context, messageKey *relayTypes.MessageKeyWithMessageHeight) ([]*relayTypes.Message, error) {
	blocks, err := p.fetchBlockMessages(ctx, &types.HeightRange{messageKey.Height, messageKey.Height})
	if err != nil {
		return nil, err
	}
	var messages []*relayTypes.Message
	for _, block := range blocks {
		messages = append(messages, block.Messages...)
	}
	return messages, nil
}

func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	return p.cfg.FinalityBlock
}

func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	msg := &relayTypes.Message{
		Sn:        sn,
		EventType: events.RevertMessage,
	}
	_, err := p.call(ctx, msg)
	return err
}

// SetFee
func (p *Provider) SetFee(ctx context.Context, networkdID string, msgFee, resFee *big.Int) error {
	msg := &relayTypes.Message{
		Src:       networkdID,
		Sn:        msgFee,
		ReqID:     resFee,
		EventType: events.SetFee,
	}
	_, err := p.call(ctx, msg)
	return err
}

// ClaimFee
func (p *Provider) ClaimFee(ctx context.Context) error {
	msg := &relayTypes.Message{
		EventType: events.ClaimFee,
	}
	_, err := p.call(ctx, msg)
	return err
}

// GetFee returns the fee for the given networkID
// responseFee is used to determine if the fee should be returned
func (p *Provider) GetFee(ctx context.Context, networkID string, responseFee bool) (uint64, error) {
	getFee := types.NewExecGetFee(networkID, responseFee)
	data, err := jsoniter.Marshal(getFee)
	if err != nil {
		return 0, err
	}
	return p.client.GetFee(ctx, p.cfg.Contracts[relayTypes.ConnectionContract], data)
}

func (p *Provider) SetAdmin(ctx context.Context, address string) error {
	msg := &relayTypes.Message{
		Src:       address,
		EventType: events.SetAdmin,
	}
	_, err := p.call(ctx, msg)
	return err
}

// ExecuteRollback
func (p *Provider) ExecuteRollback(ctx context.Context, sn *big.Int) error {
	msg := &relayTypes.Message{
		Sn:        sn,
		EventType: events.RollbackMessage,
	}
	_, err := p.call(ctx, msg)
	return err
}

func (p *Provider) getStartHeight(latestHeight, lastSavedHeight uint64) (uint64, error) {
	startHeight := lastSavedHeight
	if p.cfg.StartHeight > 0 && p.cfg.StartHeight < latestHeight {
		return p.cfg.StartHeight, nil
	}

	if startHeight > latestHeight {
		return 0, fmt.Errorf("last saved height cannot be greater than latest height")
	}

	if startHeight != 0 && startHeight < latestHeight {
		return startHeight, nil
	}

	return latestHeight, nil
}

func (p *Provider) getHeightStream(done <-chan bool, fromHeight, toHeight uint64) <-chan *types.HeightRange {
	heightChan := make(chan *types.HeightRange)
	go func(fromHeight, toHeight uint64, heightChan chan *types.HeightRange) {
		defer close(heightChan)
		for fromHeight <= toHeight {
			select {
			case <-done:
				return
			case heightChan <- &types.HeightRange{Start: fromHeight, End: fromHeight + p.cfg.BlockBatchSize - 1}:
				fromHeight += p.cfg.BlockBatchSize
			}
		}
	}(fromHeight, toHeight, heightChan)
	return heightChan
}

func (p *Provider) getBlockInfoStream(ctx context.Context, done <-chan bool, heightStreamChan <-chan *types.HeightRange) <-chan interface{} {
	blockInfoStream := make(chan interface{})
	go func(blockInfoChan chan interface{}, heightChan <-chan *types.HeightRange) {
		defer close(blockInfoChan)
		for {
			select {
			case <-done:
				return
			case height, ok := <-heightChan:
				if ok {
					for {
						messages, err := p.fetchBlockMessages(ctx, height)
						if err != nil {
							p.logger.Error("failed to fetch block messages", zap.Error(err), zap.Any("height", height))
							time.Sleep(time.Second * 3)
						} else {
							for _, message := range messages {
								blockInfoChan <- message
							}
							break
						}
					}
				}
			}
		}
	}(blockInfoStream, heightStreamChan)
	return blockInfoStream
}

// fetchBlockMessages fetches block messages from the chain
// TODO: remove retry deplicated code
func (p *Provider) fetchBlockMessages(ctx context.Context, heightInfo *types.HeightRange) ([]*relayTypes.BlockInfo, error) {
	perPage := 25
	searchParam := types.TxSearchParam{
		StartHeight: heightInfo.Start,
		EndHeight:   heightInfo.End,
		PerPage:     &perPage,
	}

	var (
		wg           sync.WaitGroup
		messages     coreTypes.ResultTxSearch
		messagesChan = make(chan *coreTypes.ResultTxSearch)
		errorChan    = make(chan error)
	)

	for _, event := range p.eventList {
		wg.Add(1)
		go func(event sdkTypes.Event) {
			defer wg.Done()

			var (
				localSearchParam = searchParam
				localMessages    = new(coreTypes.ResultTxSearch)
				retryCount       uint8
			)
			localSearchParam.Events = append(localSearchParam.Events, event)

			for retryCount < types.RPCMaxRetryAttempts {
				p.logger.Info("fetching block messages", zap.Uint64("start_height", localSearchParam.StartHeight), zap.Uint64("end_height", localSearchParam.EndHeight))
				msgs, err := p.client.TxSearch(ctx, localSearchParam)
				if err == nil {
					p.logger.Info("fetched block messages", zap.Uint64("start_height", localSearchParam.StartHeight), zap.Uint64("end_height", localSearchParam.EndHeight))
					localMessages = msgs
					break
				}
				retryCount++
				delay := time.Duration(math.Pow(2, float64(retryCount))) * types.BaseRPCRetryDelay
				if delay > types.MaxRPCRetryDelay {
					delay = types.MaxRPCRetryDelay
				}
				if retryCount >= types.RPCMaxRetryAttempts {
					p.logger.Error("fetchBlockMessages failed", zap.Uint8("attempt", retryCount), zap.Error(err))
					errorChan <- err
					break
				}
				p.logger.Error("failed to fetch block messages", zap.Uint64("start_height", localSearchParam.StartHeight), zap.Uint64("end_height", localSearchParam.EndHeight), zap.Error(err))
				time.Sleep(delay)
			}

			if localMessages.TotalCount > perPage {
				for i := 2; i <= int(localMessages.TotalCount/perPage)+1; i++ {
					p.logger.Info("fetching block messages", zap.Uint64("start_height", localSearchParam.StartHeight), zap.Uint64("end_height", localSearchParam.EndHeight), zap.Int("page", i))
					localSearchParam.Page = &i
					resNext, err := p.client.TxSearch(ctx, localSearchParam)
					if err != nil {
						var attempts uint8
						for attempts < types.RPCMaxRetryAttempts {
							p.logger.Error("failed to fetch block messages with page", zap.Uint64("start_height", localSearchParam.StartHeight), zap.Uint64("end_height", localSearchParam.EndHeight), zap.Int("page", i), zap.Error(err))
							attempts++
							delay := time.Duration(math.Pow(2, float64(attempts))) * types.BaseRPCRetryDelay
							if delay > types.MaxRPCRetryDelay {
								delay = types.MaxRPCRetryDelay
							}
							if attempts >= types.RPCMaxRetryAttempts {
								p.logger.Error("fetchBlockMessages failed", zap.Uint8("attempt", attempts), zap.Duration("retrying_in", delay), zap.Uint64("start_height", localSearchParam.StartHeight), zap.Uint64("end_height", localSearchParam.EndHeight))
								errorChan <- err
								break
							}
							p.logger.Warn("fetchBlockMessages failed, retrying...", zap.Uint8("attempt", attempts), zap.Duration("retrying_in", delay))
							time.Sleep(delay)
						}
					}
					localMessages.Txs = append(localMessages.Txs, resNext.Txs...)
				}
			}
			messagesChan <- localMessages
		}(event)
	}

	go func() {
		wg.Wait()
		close(messagesChan)
		close(errorChan)
	}()

	var errors []error
	for {
		select {
		case msgs, ok := <-messagesChan:
			if !ok {
				messagesChan = nil
			} else {
				messages.Txs = append(messages.Txs, msgs.Txs...)
				messages.TotalCount += msgs.TotalCount
			}
		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
			} else {
				errors = append(errors, err)
			}
		}
		if messagesChan == nil && errorChan == nil {
			break
		}
	}

	if len(errors) > 0 {
		p.logger.Error("Errors occurred while fetching block messages", zap.Errors("errors", errors))
		return nil, fmt.Errorf("errors occurred while fetching block messages: %v", errors)
	}

	return p.getMessagesFromTxList(messages.Txs)
}

func (p *Provider) getMessagesFromTxList(resultTxList []*coreTypes.ResultTx) ([]*relayTypes.BlockInfo, error) {
	var messages []*relayTypes.BlockInfo
	for _, resultTx := range resultTxList {
		msgs, err := p.ParseMessageFromEvents(resultTx.TxResult.GetEvents())
		if err != nil {
			return nil, err
		}
		for _, msg := range msgs {
			msg.MessageHeight = uint64(resultTx.Height)
			p.logger.Info("Detected eventlog",
				zap.Uint64("height", msg.MessageHeight),
				zap.String("target_network", msg.Dst),
				zap.Uint64("sn", msg.Sn.Uint64()),
				zap.String("event_type", msg.EventType),
			)
		}
		messages = append(messages, &relayTypes.BlockInfo{
			Height:   uint64(resultTx.Height),
			Messages: msgs,
		})
	}
	return messages, nil
}

func (p *Provider) getRawContractMessage(message *relayTypes.Message) (wasmTypes.RawContractMessage, error) {
	switch message.EventType {
	case events.EmitMessage:
		rcvMsg := types.NewExecRecvMsg(message)
		return jsoniter.Marshal(rcvMsg)
	case events.CallMessage:
		execMsg := types.NewExecExecMsg(message)
		return jsoniter.Marshal(execMsg)
	case events.RevertMessage:
		revertMsg := types.NewExecRevertMsg(message)
		return jsoniter.Marshal(revertMsg)
	case events.SetAdmin:
		setAdmin := types.NewExecSetAdmin(message.Dst)
		return jsoniter.Marshal(setAdmin)
	case events.ClaimFee:
		claimFee := types.NewExecClaimFee()
		return jsoniter.Marshal(claimFee)
	case events.SetFee:
		setFee := types.NewExecSetFee(message.Src, message.Sn, message.ReqID)
		return jsoniter.Marshal(setFee)
	case events.RollbackMessage:
		executeRollback := types.NewExecExecuteRollback(message.Sn)
		return jsoniter.Marshal(executeRollback)
	default:
		return nil, fmt.Errorf("unknown event type: %s ", message.EventType)
	}
}

func (p *Provider) getNumOfPipelines(diff int) int {
	if diff <= runtime.NumCPU() {
		return diff
	}
	return runtime.NumCPU()
}

func (p *Provider) runBlockQuery(ctx context.Context, blockInfoChan chan *relayTypes.BlockInfo, fromHeight, toHeight uint64) uint64 {
	done := make(chan bool)
	defer close(done)

	heightStream := p.getHeightStream(done, fromHeight, toHeight)

	for heightRange := range heightStream {
		var attempts uint8
		for {
			blockInfo, err := p.fetchBlockMessages(ctx, heightRange)
			if err == nil {
				for _, block := range blockInfo {
					blockInfoChan <- block
				}
				break
			}
			attempts++
			if attempts >= types.RPCMaxRetryAttempts {
				p.logger.Error("fetchBlockMessages failed", zap.Uint8("attempt", attempts), zap.Error(err))
				break
			}
			delay := time.Duration(math.Pow(2, float64(attempts))) * types.BaseRPCRetryDelay
			if delay > types.MaxRPCRetryDelay {
				delay = types.MaxRPCRetryDelay
			}
			p.logger.Warn("fetchBlockMessages failed, retrying...", zap.Uint8("attempt", attempts), zap.Duration("retrying_in", delay), zap.Uint64("start_height", heightRange.Start), zap.Uint64("end_height", heightRange.End), zap.Error(err))
			time.Sleep(delay)
		}
	}
	return toHeight + 1
}

// SubscribeMessageEvents subscribes to the message events
// Expermental: Allows to subscribe to the message events realtime without fully syncing the chain
func (p *Provider) SubscribeMessageEvents(ctx context.Context, blockInfoChan chan *relayTypes.BlockInfo, opts *types.SubscribeOpts, resetFunc func()) error {
	query := strings.Join([]string{
		"tm.event = 'Tx'",
		fmt.Sprintf("tx.height >= %d ", opts.Height),
		fmt.Sprintf("%s._contract_address = '%s'", opts.Method, opts.Address),
	}, " AND ")

	resultEventChan, err := p.client.Subscribe(ctx, "tx-result-waiter", query)
	if err != nil {
		p.logger.Error("event subscription failed", zap.Error(err))
		resetFunc()
		return err
	}
	defer p.client.Unsubscribe(ctx, opts.Address, query)
	p.logger.Info("event subscription started", zap.String("contract_address", opts.Address), zap.String("method", opts.Method))

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("event subscription stopped")
			return ctx.Err()
		case e := <-resultEventChan:
			eventDataJSON, err := jsoniter.Marshal(e.Data)
			if err != nil {
				p.logger.Error("failed to marshal event data", zap.Error(err))
				continue
			}
			var res types.TxResultResponse
			if err := jsoniter.Unmarshal(eventDataJSON, &res); err != nil {
				p.logger.Error("failed to unmarshal event data", zap.Error(err))
				continue
			}
			var messages []*relayTypes.Message
			msgs, err := p.ParseMessageFromEvents(res.Result.Events)
			if err != nil {
				p.logger.Error("failed to parse message from events", zap.Error(err))
				continue
			}
			messages = append(messages, msgs...)
			blockInfo := &relayTypes.BlockInfo{
				Height:   uint64(res.Height),
				Messages: messages,
			}
			blockInfoChan <- blockInfo
			opts.Height = blockInfo.Height
			for _, msg := range blockInfo.Messages {
				p.logger.Info("Detected eventlog",
					zap.Int64("height", res.Height),
					zap.String("target_network", msg.Dst),
					zap.Uint64("sn", msg.Sn.Uint64()),
					zap.String("event_type", msg.EventType),
				)
			}
		case <-time.After(2 * time.Minute):
			if !p.client.IsConnected() {
				p.logger.Warn("http client stopped")
				if err := p.client.Reconnect(); err != nil {
					p.logger.Warn("failed to reconnect", zap.Error(err))
					time.Sleep(time.Second * 1)
					continue
				}
				p.logger.Info("http client reconnected")
				resetFunc()
				return err
			}
		}
	}
}

// SetLastSavedHeightFunc sets the function to save the last saved height
func (p *Provider) SetLastSavedHeightFunc(f func() uint64) {
	p.LastSavedHeightFunc = f
}

// GetLastSavedHeight returns the last saved height
func (p *Provider) GetLastSavedHeight() uint64 {
	return p.LastSavedHeightFunc()
}
