package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"runtime"
	"strconv"
	"strings"
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
	"github.com/icon-project/centralized-relay/utils/concurrency"
	"github.com/icon-project/centralized-relay/utils/sorter"
	"go.uber.org/zap"
)

type Provider struct {
	logger         *zap.Logger
	cfg            *ProviderConfig
	client         IClient
	seqTracker     *SequenceTracker
	memPoolTracker *MemPoolInfo
	kms            kms.KMS
	wallet         sdkTypes.AccountI
	contracts      map[string]relayTypes.EventMap
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
		account, err := p.client.GetAccountInfo(context.Background(), p.NID())
		if err != nil {
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

	blockIntervalTicker := time.NewTicker(p.cfg.BlockInterval)
	defer blockIntervalTicker.Stop()

	p.logger.Info("Start from height", zap.Uint64("height", startHeight))

	runningLatestHeight := latestHeight

	isFirstIter := true

	for {
		select {
		case <-blockIntervalTicker.C:
			for {
				newLatestHeight, err := p.QueryLatestHeight(ctx)
				if err == nil {
					if newLatestHeight > runningLatestHeight {
						runningLatestHeight = newLatestHeight
					}
					continue
				}
				p.logger.Error("failed to query latest height", zap.Error(err))
				time.Sleep(p.cfg.BlockInterval)
			}
		default:
			if isFirstIter || runningLatestHeight > latestHeight {
				isFirstIter = false
				latestHeight = runningLatestHeight
				p.logger.Debug("Query started.", zap.Uint64("from-height", startHeight), zap.Uint64("to-height", latestHeight))
				p.runBlockQuery(blockInfoChan, startHeight, latestHeight)
				startHeight = latestHeight + 1
			}
		}
	}
}

func (p *Provider) Route(ctx context.Context, message *relayTypes.Message, callback relayTypes.TxResponseFunc) error {
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

	res, err := p.sendMessages(ctx, msgs)
	if err != nil {
		if strings.Contains(err.Error(), errors.ErrWrongSequence.Error()) {
			if mmErr := p.handleAccountSequenceMismatchError(); mmErr != nil {
				return fmt.Errorf("failed to handle sequence mismatch error: %v || %v", mmErr, err)
			}
		}
		return err
	}

	go p.waitForTxResult(ctx, message.MessageKey(), res.TxHash, callback)

	return nil
}

func (p *Provider) sendMessages(ctx context.Context, msgs []sdkTypes.Msg) (*sdkTypes.TxResponse, error) {
	if p.seqTracker == nil {
		addr := p.Wallet()
		p.seqTracker = p.NewSeqTracker(addr)
	}
	p.seqTracker.Lock()
	p.memPoolTracker.Lock()
	defer p.seqTracker.Unlock()
	defer p.memPoolTracker.Unlock()

	var accountNumber, sequence uint64

	if p.memPoolTracker.IsBlocked() {
		accountNumber, sequence = p.wallet.GetAccountNumber(), p.wallet.GetSequence()
	} else {
		senderAccount, err := p.seqTracker.Get(p.Wallet().String())
		if err != nil {
			return nil, err
		}
		accountNumber, sequence = senderAccount.AccountNumber, senderAccount.Sequence
	}

	res, err := p.prepareAndPushTxToMemPool(ctx, accountNumber, sequence, msgs)
	if err != nil {
		return nil, err
	}

	if p.memPoolTracker.IsBlocked() {
		p.memPoolTracker.SetBlockedStatus(false)
	} else if err := p.seqTracker.IncrementSequence(p.Wallet().String()); err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Provider) handleAccountSequenceMismatchError() error {
	if err := p.seqTracker.Set(p.Wallet().String(), &AccountInfo{
		AccountNumber: p.wallet.GetAccountNumber(), Sequence: p.wallet.GetSequence(),
	}); err != nil {
		return err
	}
	return nil
}

func (p *Provider) logTxFailed(err error, txHash string) {
	p.logger.Error("transaction failed",
		zap.Error(err),
		zap.String("chain_id", p.cfg.ChainID),
		zap.String("tx_hash", txHash),
	)
}

func (p *Provider) logTxSuccess(height uint64, txHash string) {
	p.logger.Info("transaction success",
		zap.Uint64("block_height", height),
		zap.String("chain_id", p.cfg.ChainID),
		zap.String("tx_hash", txHash),
	)
}

func (p *Provider) prepareAndPushTxToMemPool(ctx context.Context, accountNumber, sequence uint64, msgs []sdkTypes.Msg) (*sdkTypes.TxResponse, error) {
	txf, err := p.client.BuildTxFactory()
	if err != nil {
		return nil, err
	}

	txf = txf.
		WithGasPrices(p.cfg.GasPrices).WithGasAdjustment(p.cfg.GasAdjustment).
		WithAccountNumber(accountNumber).
		WithSequence(sequence)

	if txf.SimulateAndExecute() {
		_, adjusted, err := p.client.CalculateGas(txf, msgs)
		if err != nil {
			return nil, err
		}
		txf = txf.WithGas(adjusted)
	}

	if txf.Gas() == 0 {
		return nil, fmt.Errorf("gas amount cannot be zero")
	}

	if p.cfg.MinGasAmount > 0 && txf.Gas() < p.cfg.MinGasAmount {
		return nil, fmt.Errorf("gas amount %d is too low; the minimum allowed gas amount is %d", txf.Gas(), p.cfg.MinGasAmount)
	}

	if p.cfg.MaxGasAmount > 0 && txf.Gas() > p.cfg.MaxGasAmount {
		return nil, fmt.Errorf("gas amount %d exceeds the maximum allowed limit of %d", txf.Gas(), p.cfg.MaxGasAmount)
	}

	txBytes, err := p.client.PrepareTx(ctx, txf, msgs)
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
			p.memPoolTracker.SetBlockedStatusWithLock(true)
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
		for range time.NewTicker(p.cfg.BlockInterval).C {
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
	queryMsg := types.QueryReceiptMsg{
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
	if err := json.Unmarshal(res.Data, &receiptMsgRes.Status); err != nil {
		return false, err
	}

	if receiptMsgRes.Status {
		return true, nil
	}

	return false, nil
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
	return nil, nil
}

func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	return 0
}

func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	return nil
}

func (p *Provider) SetAdmin(context.Context, string) error {
	return nil
}

// ExecuteCall executes a call to the bridge contract
func (p *Provider) ExecuteCall(ctx context.Context, reqID *big.Int, data []byte) ([]byte, error) {
	return nil, nil
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

func (p *Provider) getHeightStream(done <-chan interface{}, fromHeight, toHeight uint64) <-chan uint64 {
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

func (p *Provider) getBlockInfoStream(done <-chan interface{}, heightStream <-chan uint64) <-chan interface{} {
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
						messages, err := p.fetchBlockMessages(height)
						if err != nil {
							p.logger.Error("failed to fetch block messages: ", zap.Error(err), zap.Uint64("block-height", height))
							time.Sleep(500 * time.Millisecond)
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

func (p *Provider) fetchBlockMessages(height uint64) ([]*relayTypes.Message, error) {
	searchParam := types.TxSearchParam{
		BlockHeight: height,
		Events: sdkTypes.Events{
			{
				Type:       EventTypeWasmMessage,
				Attributes: p.GetMonitorEventFilters(),
			},
		},
	}

	res, err := p.client.TxSearch(context.Background(), searchParam)
	if err != nil {
		return nil, err
	}

	return p.getMessagesFromTxList(res.Txs)
}

func (p *Provider) getMessagesFromTxList(resultTxList []*coreTypes.ResultTx) ([]*relayTypes.Message, error) {
	var messages []*relayTypes.Message
	for _, resultTx := range resultTxList {
		var eventsList []*EventsList
		err := json.Unmarshal([]byte(resultTx.TxResult.Log), &eventsList)
		if err != nil {
			return nil, err
		}

		for _, event := range eventsList {
			messages, err := p.ParseMessageFromEvents(event.Events)
			if err != nil {
				return nil, err
			}
			for _, message := range messages {
				message.MessageHeight = uint64(resultTx.Height)
				if message.Dst != "" {
					p.logger.Info("Detected eventlog",
						zap.Uint64("height", message.MessageHeight),
						zap.String("target_network", message.Dst),
						zap.Uint64("sn", message.Sn),
						zap.String("event_type", message.EventType),
					)
					messages = append(messages, message)
				}
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

func (p *Provider) runBlockQuery(blockInfoChan chan *relayTypes.BlockInfo, fromHeight, toHeight uint64) {
	done := make(chan interface{})
	defer close(done)

	heightStream := p.getHeightStream(done, fromHeight, toHeight)

	numOfPipelines := p.getNumOfPipelines(fromHeight, toHeight)
	pipelines := make([]<-chan interface{}, numOfPipelines)

	for i := 0; i < numOfPipelines; i++ {
		pipelines[i] = p.getBlockInfoStream(done, heightStream)
	}

	var blockInfoList []*relayTypes.BlockInfo
	for bn := range concurrency.Take(done, concurrency.FanIn(done, pipelines...), int(toHeight-fromHeight+1)) {
		block := bn.(*relayTypes.BlockInfo)
		blockInfoList = append(blockInfoList, block)
	}

	sorter.Sort(blockInfoList, func(p1, p2 *relayTypes.BlockInfo) bool {
		return p1.Height < p2.Height // ascending order
	})

	for _, blockInfo := range blockInfoList {
		blockInfoChan <- &relayTypes.BlockInfo{
			Height: blockInfo.Height, Messages: blockInfo.Messages,
		}
	}
}
