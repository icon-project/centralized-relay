package wasm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	abiTypes "github.com/cometbft/cometbft/abci/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/client"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	relayerEvents "github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/concurrency"
	"github.com/icon-project/centralized-relay/utils/sorter"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Provider struct {
	logger  *zap.Logger
	config  ProviderConfig
	client  client.IClient
	txMutex sync.Mutex
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
	return p.config.NID
}

func (p *Provider) ChainName() string {
	return p.config.ChainName
}

func (p *Provider) Init(ctx context.Context) error {
	return nil
}

func (p *Provider) Type() string {
	return types.ChainType
}

func (p *Provider) ProviderConfig() provider.ProviderConfig {
	return p.config
}

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfoChan chan relayTypes.BlockInfo) error {
	latestHeight, err := p.QueryLatestHeight(ctx)
	if err != nil {
		p.logger.Error("failed to get latest block height: ", zap.Error(err))
		return err
	}

	startHeight, err := p.getStartHeight(latestHeight, lastSavedHeight)
	if err != nil {
		p.logger.Error("failed to determine start height: ", zap.Error(err))
		return err
	}

	blockInterval, err := time.ParseDuration(p.config.BlockInterval)
	if err != nil {
		p.logger.Error("failed to parse block interval: ", zap.Error(err))
		return err
	}

	blockIntervalTicker := time.NewTicker(blockInterval)
	defer blockIntervalTicker.Stop()

	p.logger.Info("start querying from height", zap.Uint64("start-height", startHeight))

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
					break
				}
				p.logger.Error("failed to query latest height", zap.Error(err))
				time.Sleep(500 * time.Millisecond)
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
	return nil
}

func (p *Provider) Route(ctx context.Context, message *relayTypes.Message, callback relayTypes.TxResponseFunc) error {
	rawMsg, err := p.getRawContractMessage(message)
	if err != nil {
		return err
	}
	msg := wasmTypes.MsgExecuteContract{
		Sender:   p.client.Context().FromAddress.String(),
		Contract: p.config.ContractAddress,
		Msg:      rawMsg,
	}

	msgs := []sdkTypes.Msg{&msg}

	txf, err := p.buildTxFactory()
	if err != nil {
		return err
	}

	if txf.SimulateAndExecute() {
		_, adjusted, err := tx.CalculateGas(p.client.Context(), txf, msgs...)
		if err != nil {
			return err
		}
		txf = txf.WithGas(adjusted)
	}

	if txf.Gas() == 0 {
		return fmt.Errorf("gas amount cannot be zero")
	}

	if p.config.MinGasAmount > 0 && txf.Gas() < p.config.MinGasAmount {
		return fmt.Errorf("gas amount %d is too low; the minimum allowed gas amount is %d", txf.Gas(), p.config.MinGasAmount)
	}

	if p.config.MaxGasAmount > 0 && txf.Gas() > p.config.MaxGasAmount {
		return fmt.Errorf("gas amount %d exceeds the maximum allowed limit of %d", txf.Gas(), p.config.MaxGasAmount)
	}

	res, err := p.client.SendTx(ctx, txf, msgs)
	if err != nil || res.Code != types.CodeTypeOK {
		if err == nil {
			err = fmt.Errorf("failed to send tx: %v", res.RawLog)
		}
		callback(message.MessageKey(), relayTypes.TxResponse{}, err)
		return err
	}

	go p.waitForTxResult(ctx, message.MessageKey(), res.TxHash, callback)

	return nil
}

func (p *Provider) logTxFailed(err error, txHash string) {
	p.logger.Error("transaction failed: ",
		zap.Error(err),
		zap.String("chain-id", p.config.ChainID),
		zap.String("tx-hash", txHash),
	)
}

func (p *Provider) logTxSuccess(height uint64, txHash string) {
	p.logger.Error("transaction success: ",
		zap.Uint64("block-height", height),
		zap.String("chain-id", p.config.ChainID),
		zap.String("tx-hash", txHash),
	)
}

func (p *Provider) waitForTxResult(ctx context.Context, mk relayTypes.MessageKey, txHash string, callback relayTypes.TxResponseFunc) {
	client, err := p.client.HTTP(p.config.RpcUrl)
	if err != nil {
		p.logTxFailed(err, txHash)
		callback(mk, relayTypes.TxResponse{}, err)
		return
	}
	if err := client.Start(); err != nil {
		p.logTxFailed(err, txHash)
		callback(mk, relayTypes.TxResponse{}, err)
		return
	}
	defer client.Stop()

	timeOutInterval := types.TxConfirmationIntervalDefault
	if p.config.TxConfirmationInterval != "" {
		timeOutInterval, err = time.ParseDuration(p.config.TxConfirmationInterval)
		if err != nil {
			p.logTxFailed(err, txHash)
			callback(mk, relayTypes.TxResponse{}, err)
			return
		}
	}
	ctx, cancel := context.WithTimeout(ctx, timeOutInterval)
	defer cancel()

	query := fmt.Sprintf("tm.event = 'Tx' AND tx.hash = '%s'", txHash)
	resultEventChan, err := client.Subscribe(ctx, "tx-result-waiter", query)
	if err != nil {
		p.logTxFailed(err, txHash)
		callback(mk, relayTypes.TxResponse{}, err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			p.logTxFailed(err, txHash)
			callback(mk, relayTypes.TxResponse{}, ctx.Err())
			return
		case e := <-resultEventChan:
			eventDataJSON, err := json.Marshal(e.Data)
			if err != nil {
				p.logTxFailed(err, txHash)
				callback(mk, relayTypes.TxResponse{}, ctx.Err())
				return
			}

			var txWaitRes types.TxResultWaitResponse
			err = json.Unmarshal(eventDataJSON, &txWaitRes)
			if err != nil {
				p.logTxFailed(err, txHash)
				callback(mk, relayTypes.TxResponse{}, ctx.Err())
				return
			}

			if uint32(txWaitRes.Result.Code) != types.CodeTypeOK {
				p.logTxFailed(err, txHash)
				callback(mk, relayTypes.TxResponse{}, errors.New("something went wrong"))
				return
			}

			p.logTxSuccess(uint64(txWaitRes.Height), txHash)
			callback(mk, relayTypes.TxResponse{
				Height:    txWaitRes.Height,
				TxHash:    txHash,
				Codespace: txWaitRes.Result.Codespace,
				Code:      relayTypes.ResponseCode(txWaitRes.Result.Code),
				Data:      string(txWaitRes.Result.Data), //Todo this need to be confirmed.
			}, nil)
			return
		}
	}
}

func (p *Provider) MessageReceived(ctx context.Context, key relayTypes.MessageKey) (bool, error) {
	queryMsg := types.QueryReceiptMsg{
		GetReceipt: types.GetReceiptMsg{
			SrcNetwork: key.Src,
			ConnSn:     strconv.Itoa(int(key.Sn)),
		},
	}

	rawQueryMsg, err := json.Marshal(queryMsg)
	if err != nil {
		return false, err
	}

	res, err := p.client.QuerySmartContract(ctx, p.config.ContractAddress, rawQueryMsg)
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
	coin, err := p.client.GetBalance(ctx, addr, p.config.Denomination)
	if err != nil {
		p.logger.Error("failed to query balance: ", zap.Error(err))
		return nil, err
	}
	return &relayTypes.Coin{
		Denom:  coin.Denom,
		Amount: coin.Amount.Uint64(),
	}, nil
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, message relayTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, message relayTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) GenerateMessage(ctx context.Context, messageKey *relayTypes.MessageKeyWithMessageHeight) (*relayTypes.Message, error) {
	return nil, nil
}

func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	return 0
}

func (p *Provider) getStartHeight(latestHeight, lastSavedHeight uint64) (uint64, error) {
	startHeight := lastSavedHeight
	if p.config.StartHeight > 0 {
		startHeight = p.config.StartHeight
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
	go func() {
		defer close(blockInfoStream)
		for {
			select {
			case <-done:
				return
			case height, ok := <-heightStream:
				if ok {
					for {
						messages, err := p.fetchBlockMessages(height)
						if err != nil {
							p.logger.Error("failed to fetch block messages: ", zap.Error(err), zap.Uint64("block-height", height))
							time.Sleep(500 * time.Millisecond)
						} else {
							blockInfoStream <- relayTypes.BlockInfo{
								Height:   height,
								Messages: messages,
							}
							break
						}
					}
				}
			}
		}
	}()
	return blockInfoStream
}

func (p *Provider) fetchBlockMessages(height uint64) ([]*relayTypes.Message, error) {
	eventFilters := sdkTypes.Events{
		{
			Type: EventTypeWasmMessage,
			Attributes: []abiTypes.EventAttribute{
				{Key: EventAttrKeyContractAddress, Value: fmt.Sprintf("'%s'", p.config.ContractAddress)},
			},
		},
		//Todo add custom event type in contract for specific events and filter here
	}

	searchParam := types.TxSearchParam{
		BlockHeight: height,
		Events:      eventFilters,
	}

	res, err := p.client.TxSearch(context.Background(), searchParam)
	if err != nil {
		return nil, err
	}

	return p.getMessagesFromTxList(res.Txs)
}

func (p *Provider) getMessagesFromTxList(resultTx []*coretypes.ResultTx) ([]*relayTypes.Message, error) {
	var messages []*relayTypes.Message
	for _, tx := range resultTx {
		var eventsList []EventsList
		err := json.Unmarshal([]byte(tx.TxResult.Log), &eventsList)
		if err != nil {
			return nil, err
		}

		if len(eventsList) > 0 {
			for _, events := range eventsList {
				message, err := ParseMessageFromEvents(events.Events)
				if err != nil {
					return nil, err
				}
				message.MessageHeight = uint64(tx.Height)
				message.Src = p.NID()
				message.EventType = relayerEvents.EmitMessage

				if message.Dst != "" {
					p.logger.Info("detected event log ", zap.Uint64("height", message.MessageHeight),
						zap.String("target-network", message.Dst),
						zap.Uint64("sn", message.Sn),
						zap.String("event-type", message.EventType),
					)
					messages = append(messages, &message)
				}
			}
		}
	}
	return messages, nil
}

func (p *Provider) buildTxFactory() (tx.Factory, error) {
	txf, err := tx.NewFactoryCLI(p.client.Context(), &pflag.FlagSet{})
	if err != nil {
		return tx.Factory{}, err
	}

	senderAccount, err := p.client.GetAccountInfo(context.Background(), p.client.Context().FromAddress.String())
	if err != nil {
		return tx.Factory{}, err
	}

	txf = txf.
		WithAccountNumber(senderAccount.GetAccountNumber()).WithSequence(senderAccount.GetSequence()).
		WithTxConfig(p.client.Context().TxConfig).
		WithKeybase(p.client.Context().Keyring).
		WithFeePayer(p.client.Context().FeePayer).
		WithChainID(p.client.Context().ChainID).
		WithSimulateAndExecute(p.client.Context().Simulate).
		WithGasPrices(p.config.GasPrices).
		WithGasAdjustment(p.config.GasAdjustment)

	return txf, nil
}

func (p *Provider) getRawContractMessage(message *relayTypes.Message) (wasmTypes.RawContractMessage, error) {
	switch message.EventType {
	case relayerEvents.EmitMessage:
		rcvMsg := types.NewExecRecvMsg(message)
		return json.Marshal(rcvMsg)
	default:
		return nil, fmt.Errorf("unknown event type: %s ", message.EventType)
	}
}

func (p *Provider) getNumOfPipelines(startHeight, latestHeight uint64) int {
	diff := latestHeight - startHeight + 1 //since both heights are inclusive
	if int(diff) < runtime.NumCPU() {
		return int(diff)
	}
	return runtime.NumCPU()
}

func (p *Provider) runBlockQuery(blockInfoChan chan relayTypes.BlockInfo, fromHeight, toHeight uint64) {
	done := make(chan interface{})
	defer close(done)

	heightStream := p.getHeightStream(done, fromHeight, toHeight)

	numOfPipelines := p.getNumOfPipelines(fromHeight, toHeight)
	pipelines := make([]<-chan interface{}, numOfPipelines)

	for i := 0; i < numOfPipelines; i++ {
		pipelines[i] = p.getBlockInfoStream(done, heightStream)
	}

	var blockInfoList []relayTypes.BlockInfo
	for bn := range concurrency.Take(done, concurrency.FanIn(done, pipelines...), int(toHeight-fromHeight+1)) {
		block := bn.(relayTypes.BlockInfo)
		blockInfoList = append(blockInfoList, block)
	}

	sorter.Sort(blockInfoList, func(p1, p2 relayTypes.BlockInfo) bool {
		return p1.Height < p2.Height //ascending order
	})

	for _, blockInfo := range blockInfoList {
		blockInfoChan <- relayTypes.BlockInfo{
			Height: blockInfo.Height, Messages: blockInfo.Messages,
		}
	}
}
