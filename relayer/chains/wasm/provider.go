package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	coreTypes "github.com/cometbft/cometbft/rpc/core/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/concurrency"
	"go.uber.org/zap"
)

type Provider struct {
	logger    *zap.Logger
	cfg       *ProviderConfig
	client    IClient
	kms       kms.KMS
	wallet    sdkTypes.AccountI
	contracts map[string]relayTypes.EventMap
	eventList []sdkTypes.Event
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
	if p.wallet == nil {
		if err := p.RestoreKeystore(context.Background()); err != nil {
			p.logger.Error("failed to restore keystore", zap.Error(err))
			return nil
		}
		account, err := p.client.GetAccountInfo(context.Background(), p.cfg.GetWallet())
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

	heightTicker := time.NewTicker(p.cfg.BlockInterval)
	heightPoller := time.NewTicker(time.Minute)
	sequenceTicker := time.NewTicker(3 * time.Minute)
	defer heightTicker.Stop()
	defer heightPoller.Stop()
	defer sequenceTicker.Stop()

	p.logger.Info("Start from height", zap.Uint64("height", startHeight), zap.Uint64("finality block", p.FinalityBlock(ctx)))

	for {
		select {
		case <-heightTicker.C:
			latestHeight++
		case <-heightPoller.C:
			height, err := p.QueryLatestHeight(ctx)
			if err != nil {
				p.logger.Error("failed to query latest height", zap.Error(err))
				heightPoller.Reset(time.Second * 3)
				continue
			}
			heightPoller.Reset(time.Minute)
			latestHeight = height
		case <-ctx.Done():
			return ctx.Err()
		case <-sequenceTicker.C:
			if err := p.handleSequence(ctx); err != nil {
				p.logger.Error("failed to update sequence", zap.Error(err))
			}
		default:
			if startHeight > latestHeight {
				continue
			}
			toHeight := latestHeight
			p.logger.Debug("Query started", zap.Uint64("from-height", startHeight), zap.Uint64("to-height", latestHeight))
			p.runBlockQuery(ctx, blockInfoChan, startHeight, toHeight)
			startHeight = toHeight + 1
		}
	}
}

func (p *Provider) Route(ctx context.Context, message *relayTypes.Message, callback relayTypes.TxResponseFunc) error {
	p.logger.Info("starting to route message", zap.Any("message", message))
	rawMsg, err := p.getRawContractMessage(message)
	if err != nil {
		return err
	}

	var contract string

	switch message.EventType {
	case events.CallMessage:
		contract = p.cfg.Contracts[relayTypes.XcallContract]
	case events.EmitMessage:
		contract = p.cfg.Contracts[relayTypes.ConnectionContract]
	default:
		return fmt.Errorf("unknown event type: %s ", message.EventType)
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
				return fmt.Errorf("failed to handle sequence mismatch error: %v || %v", mmErr, err)
			}
		}
		return err
	}

	go p.waitForTxResult(ctx, message.MessageKey(), res.TxHash, callback)
	seq := p.wallet.GetSequence() + 1
	if err := p.wallet.SetSequence(seq); err != nil {
		p.logger.Error("failed to set sequence", zap.Error(err))
	}

	return nil
}

func (p *Provider) sendMessage(ctx context.Context, msgs ...sdkTypes.Msg) (*sdkTypes.TxResponse, error) {
	res, err := p.prepareAndPushTxToMemPool(ctx, p.wallet.GetAccountNumber(), p.wallet.GetSequence(), msgs...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Provider) handleSequence(ctx context.Context) error {
	addr := p.Wallet()
	acc, err := p.client.GetAccountInfo(ctx, addr.String())
	if err != nil {
		return err
	}
	return p.wallet.SetSequence(acc.GetSequence())
}

func (p *Provider) logTxFailed(err error, txHash string) {
	p.logger.Error("transaction failed",
		zap.Error(err),
		zap.String("tx_hash", txHash),
	)
}

func (p *Provider) logTxSuccess(height uint64, txHash string) {
	p.logger.Info("successful transaction",
		zap.Uint64("block_height", height),
		zap.String("chain_id", p.cfg.ChainID),
		zap.String("tx_hash", txHash),
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

func (p *Provider) waitForTxResult(ctx context.Context, mk *relayTypes.MessageKey, txHash string, callback relayTypes.TxResponseFunc) {
	for txWaitRes := range p.subscribeTxResultStream(ctx, txHash, p.cfg.TxConfirmationInterval) {
		if txWaitRes.Error != nil {
			p.logTxFailed(txWaitRes.Error, txHash)
			callback(mk, txWaitRes.TxResult, txWaitRes.Error)
			return
		}
		p.logTxSuccess(uint64(txWaitRes.TxResult.Height), txHash)
		callback(mk, txWaitRes.TxResult, nil)
	}
}

func (p *Provider) pollTxResultStream(ctx context.Context, txHash string, maxWaitInterval time.Duration) <-chan *types.TxResultChan {
	txResChan := make(chan *types.TxResultChan)
	startTime := time.Now()
	go func(txChan chan *types.TxResultChan) {
		defer close(txChan)
		for range time.NewTicker(p.cfg.TxConfirmationInterval).C {
			res, err := p.client.GetTransactionReceipt(ctx, txHash)
			if err == nil {
				txChan <- &types.TxResultChan{
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
				txChan <- &types.TxResultChan{
					Error: err,
				}
				return
			}
		}
	}(txResChan)
	return txResChan
}

func (p *Provider) subscribeTxResultStream(ctx context.Context, txHash string, maxWaitInterval time.Duration) <-chan *types.TxResultChan {
	txResChan := make(chan *types.TxResultChan)
	go func(txRes chan *types.TxResultChan) {
		defer close(txRes)
		httpClient, err := p.client.HTTP(p.cfg.RpcUrl)
		if err != nil {
			txRes <- &types.TxResultChan{
				TxResult: nil, Error: err,
			}
			return
		}
		if err := httpClient.Start(); err != nil {
			txRes <- &types.TxResultChan{
				TxResult: nil, Error: err,
			}
			return
		}
		defer httpClient.Stop()

		newCtx, cancel := context.WithTimeout(ctx, maxWaitInterval)
		defer cancel()

		query := fmt.Sprintf("tm.event = 'Tx' AND tx.hash = '%s'", txHash)
		resultEventChan, err := httpClient.Subscribe(newCtx, "tx-result-waiter", query)
		if err != nil {
			txRes <- &types.TxResultChan{
				TxResult: nil, Error: err,
			}
			return
		}

		select {
		case <-ctx.Done():
			txRes <- &types.TxResultChan{
				TxResult: nil, Error: ctx.Err(),
			}
			return
		case e := <-resultEventChan:
			eventDataJSON, err := json.Marshal(e.Data)
			if err != nil {
				txRes <- &types.TxResultChan{
					TxResult: nil, Error: err,
				}
				return
			}

			txWaitRes := new(types.TxResultWaitResponse)
			if err := json.Unmarshal(eventDataJSON, txWaitRes); err != nil {
				txRes <- &types.TxResultChan{
					TxResult: nil, Error: err,
				}
				return
			}
			if uint32(txWaitRes.Result.Code) != types.CodeTypeOK {
				txRes <- &types.TxResultChan{
					Error: fmt.Errorf(txWaitRes.Result.Log),
					TxResult: &relayTypes.TxResponse{
						Height:    txWaitRes.Height,
						TxHash:    txHash,
						Codespace: txWaitRes.Result.Codespace,
						Code:      relayTypes.ResponseCode(txWaitRes.Result.Code),
						Data:      string(txWaitRes.Result.Data),
					},
				}
				return
			}

			txRes <- &types.TxResultChan{
				TxResult: &relayTypes.TxResponse{
					Height:    txWaitRes.Height,
					TxHash:    txHash,
					Codespace: txWaitRes.Result.Codespace,
					Code:      relayTypes.ResponseCode(txWaitRes.Result.Code),
					Data:      string(txWaitRes.Result.Data),
				},
			}
		}
	}(txResChan)
	return txResChan
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayTypes.MessageKey) (bool, error) {
	queryMsg := &types.QueryReceiptMsg{
		GetReceipt: &types.GetReceiptMsg{
			SrcNetwork: key.Src,
			ConnSn:     strconv.Itoa(int(key.Sn)),
		},
	}
	rawQueryMsg, err := json.Marshal(queryMsg)
	if err != nil {
		return false, err
	}

	res, err := p.client.QuerySmartContract(ctx, p.cfg.Contracts[relayTypes.ConnectionContract], rawQueryMsg)
	if err != nil {
		p.logger.Error("failed to check if message is received: ", zap.Error(err))
		return false, err
	}

	receiptMsgRes := types.QueryReceiptMsgResponse{}
	return receiptMsgRes.Status, json.Unmarshal(res.Data, &receiptMsgRes.Status)
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*relayTypes.Coin, error) {
	coin, err := p.client.GetBalance(ctx, addr, p.cfg.Denomination)
	if err != nil {
		p.logger.Error("failed to query balance: ", zap.Error(err))
		return nil, err
	}
	return &relayTypes.Coin{
		Denom:  coin.Denom,
		Amount: coin.Amount.Uint64(),
	}, nil
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, message *relayTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, message *relayTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) GenerateMessage(ctx context.Context, messageKey *relayTypes.MessageKeyWithMessageHeight) (*relayTypes.Message, error) {
	msgs, err := p.fetchBlockMessages(ctx, messageKey.Height)
	if err != nil {
		return nil, err
	}
	for _, msg := range msgs {
		if msg.Sn == messageKey.Sn {
			return msg, nil
		}
	}
	return nil, fmt.Errorf("message not found")
}

func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	return p.cfg.FinalityBlock
}

func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	execMsg := types.NewExecRevertMsg(&relayTypes.Message{
		Sn: sn.Uint64(),
	})
	rawMsg, err := json.Marshal(execMsg)
	if err != nil {
		return err
	}

	msg := &wasmTypes.MsgExecuteContract{
		Sender:   p.Wallet().String(),
		Contract: p.cfg.Contracts[relayTypes.ConnectionContract],
		Msg:      rawMsg,
	}

	msgs := []sdkTypes.Msg{msg}

	_, err = p.sendMessage(ctx, msgs...)
	return err
}

func (p *Provider) SetAdmin(ctx context.Context, address string) error {
	execMsg := types.NewExecSetAdmin(address)
	rawMsg, err := json.Marshal(execMsg)
	if err != nil {
		return err
	}
	msg := &wasmTypes.MsgExecuteContract{
		Sender:   p.Wallet().String(),
		Contract: p.cfg.Contracts[relayTypes.ConnectionContract],
		Msg:      rawMsg,
	}
	msgs := []sdkTypes.Msg{msg}
	_, err = p.sendMessage(ctx, msgs...)
	return err
}

func (p *Provider) getStartHeight(latestHeight, lastSavedHeight uint64) (uint64, error) {
	startHeight := lastSavedHeight
	if p.cfg.StartHeight > 0 {
		startHeight = p.cfg.StartHeight
	}

	if startHeight > latestHeight {
		return 0, fmt.Errorf("last saved height cannot be greater than latest height")
	}

	if startHeight != 0 && startHeight < latestHeight {
		return startHeight, nil
	}

	return latestHeight, nil
}

func (p *Provider) getHeightStream(done <-chan bool, fromHeight, toHeight uint64) <-chan uint64 {
	heightStream := make(chan uint64)
	go func() {
		defer close(heightStream)
		for i := fromHeight; i <= toHeight; i++ {
			select {
			case <-done:
				return
			case heightStream <- i:
			}
		}
	}()
	return heightStream
}

func (p *Provider) getBlockInfoStream(ctx context.Context, done <-chan bool, heightStream <-chan uint64) <-chan interface{} {
	blockInfoStream := make(chan interface{})
	go func(blockInfoChan chan interface{}, heightChan <-chan uint64) {
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
							p.logger.Error("failed to fetch block messages: ", zap.Error(err), zap.Uint64("block-height", height))
							time.Sleep(time.Second)
						} else {
							blockInfoChan <- &relayTypes.BlockInfo{
								Height:   height,
								Messages: messages,
							}
							break
						}
					}
				}
			}
		}
	}(blockInfoStream, heightStream)
	return blockInfoStream
}

func (p *Provider) fetchBlockMessages(ctx context.Context, height uint64) ([]*relayTypes.Message, error) {
	perPage := 100
	searchParam := types.TxSearchParam{
		BlockHeight: height,
		PerPage:     &perPage,
	}

	var (
		wg           sync.WaitGroup
		messages     coretypes.ResultTxSearch
		messagesChan = make(chan *coretypes.ResultTxSearch)
		errorChan    = make(chan error)
	)

	for _, event := range p.eventList {
		wg.Add(1)
		go func(wg *sync.WaitGroup, searchParam types.TxSearchParam, messagesChan chan *coretypes.ResultTxSearch, errorChan chan error) {
			defer wg.Done()
			searchParam.Events = append(searchParam.Events, event)
			res, err := p.client.TxSearch(ctx, searchParam)
			if err != nil {
				errorChan <- err
				return
			}
			messagesChan <- res
		}(&wg, searchParam, messagesChan, errorChan)
		select {
		case msgs := <-messagesChan:
			messages.Txs = append(messages.Txs, msgs.Txs...)
			messages.TotalCount += msgs.TotalCount
		case err := <-errorChan:
			return nil, err
		}
	}
	wg.Wait()
	return p.getMessagesFromTxList(messages.Txs)
}

func (p *Provider) getMessagesFromTxList(resultTxList []*coreTypes.ResultTx) ([]*relayTypes.Message, error) {
	var messages []*relayTypes.Message
	for _, resultTx := range resultTxList {
		var eventsList []*EventsList
		if err := json.Unmarshal([]byte(resultTx.TxResult.Log), &eventsList); err != nil {
			return nil, err
		}

		for _, event := range eventsList {
			msgs, err := p.ParseMessageFromEvents(event.Events)
			if err != nil {
				return nil, err
			}
			for _, msg := range msgs {
				msg.MessageHeight = uint64(resultTx.Height)
				p.logger.Info("Detected eventlog",
					zap.Uint64("height", msg.MessageHeight),
					zap.String("target_network", msg.Dst),
					zap.Uint64("sn", msg.Sn),
					zap.String("event_type", msg.EventType),
				)
				messages = append(messages, msg)
			}
		}
	}
	return messages, nil
}

func (p *Provider) getRawContractMessage(message *relayTypes.Message) (wasmTypes.RawContractMessage, error) {
	switch message.EventType {
	case events.EmitMessage:
		rcvMsg := types.NewExecRecvMsg(message)
		return json.Marshal(rcvMsg)
	case events.CallMessage:
		execMsg := types.NewExecExecMsg(message)
		return json.Marshal(execMsg)
	default:
		return nil, fmt.Errorf("unknown event type: %s ", message.EventType)
	}
}

func (p *Provider) getNumOfPipelines(startHeight, latestHeight uint64) int {
	diff := latestHeight - startHeight + 1 // since both heights are inclusive
	if int(diff) < runtime.NumCPU() {
		return int(diff)
	}
	return runtime.NumCPU()
}

func (p *Provider) runBlockQuery(ctx context.Context, blockInfoChan chan *relayTypes.BlockInfo, fromHeight, toHeight uint64) {
	done := make(chan bool)
	defer close(done)

	heightStream := p.getHeightStream(done, fromHeight, toHeight)

	numOfPipelines := p.getNumOfPipelines(fromHeight, toHeight)
	pipelines := make([]<-chan interface{}, numOfPipelines)

	for i := 0; i < numOfPipelines; i++ {
		pipelines[i] = p.getBlockInfoStream(ctx, done, heightStream)
	}

	for bn := range concurrency.Take(done, concurrency.FanIn(done, pipelines...), int(toHeight-fromHeight+1)) {
		block := bn.(*relayTypes.BlockInfo)
		blockInfoChan <- block
	}
}

// SubscribeMessageEvents subscribes to the message events
// Expermental: Allows to subscribe to the message events from the chain realtime without syncing the chain
func (p *Provider) SubscribeMessageEvents(ctx context.Context, contractAddress string, height uint64) error {
	httpClient, err := p.client.HTTP(p.cfg.RpcUrl)
	if err != nil {
		p.logger.Error("failed to create http client", zap.Error(err))
		return err
	}
	defer httpClient.Stop()
	if err := httpClient.Start(); err != nil {
		p.logger.Error("http client start failed", zap.Error(err))
		return err
	}
	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	query := strings.Join([]string{
		"tm.event = 'Tx'",
		fmt.Sprintf("tx.height >= %d ", height),
		fmt.Sprintf("wasm-Message._contract_address = '%s'", contractAddress),
	}, " AND ")
	resultEventChan, err := httpClient.Subscribe(newCtx, "event", query)
	if err != nil {
		p.logger.Error("event subscription failed", zap.Error(err))
		return err
	}
	for {
		select {
		case <-ctx.Done():
			p.logger.Info("event subscription stopped")
			return ctx.Err()
		case e := <-resultEventChan:
			eventDataJSON, err := json.Marshal(e.Data)
			if err != nil {
				p.logger.Error("failed to marshal event data", zap.Error(err))
				return err
			}
			p.logger.Info("event data", zap.ByteString("data", eventDataJSON))
		}
	}
}
