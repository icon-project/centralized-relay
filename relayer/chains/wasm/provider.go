package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/client"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/concurrency"
	"go.uber.org/zap"
	"runtime"
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

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfo chan relayTypes.BlockInfo) error {
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

	for {
		select {
		case <-blockIntervalTicker.C:
			func() {
				done := make(chan interface{})
				defer close(done)

				heightStream := p.getHeightStream(done, startHeight, latestHeight)

				numOfPipelines := runtime.NumCPU() //Todo tune or configure this

				pipelines := make([]<-chan interface{}, numOfPipelines)

				for i := 0; i < numOfPipelines; i++ {
					pipelines[i] = p.getBlockInfoStream(done, heightStream)
				}

				for bn := range concurrency.FanIn(done, pipelines...) {
					block, ok := bn.(relayTypes.BlockInfo)
					if !ok || block.HasError() {
						if !block.HasError() {
							block.Error = fmt.Errorf("received invalid block type -> required: %T, got: %T", relayTypes.BlockInfo{}, bn)
						}
						p.logger.Error("error receiving block: ", zap.Error(block.Error))
						continue
					}
					blockInfo <- relayTypes.BlockInfo{
						Height: block.Height, Messages: block.Messages,
					}
				}
			}()
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
		p.logger.Error("failed to route message: ", zap.Error(err))
		callback(message.MessageKey(), relayTypes.TxResponse{}, err)
		return err
	}

	callback(message.MessageKey(), relayTypes.TxResponse{
		Height:    res.Height,
		TxHash:    res.TxHash,
		Codespace: res.Codespace,
		Code:      relayTypes.ResponseCode(res.Code),
		Data:      res.Data,
	}, nil)

	return nil
}

func (p *Provider) MessageReceived(ctx context.Context, key relayTypes.MessageKey) (bool, error) {
	queryMsg := types.QueryReceiptMsg{
		GetReceipt: types.GetReceiptMsg{
			SrcNetwork: key.Src,
			ConnSn:     key.Sn,
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
	if err := json.Unmarshal(res.Data, &receiptMsgRes); err != nil {
		return false, err
	}

	if receiptMsgRes.Status == 1 {
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
	if lastSavedHeight > latestHeight {
		return 0, fmt.Errorf("last saved height cannot be greater than latest height")
	}

	if lastSavedHeight != 0 && lastSavedHeight < latestHeight {
		return lastSavedHeight, nil
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
			case height := <-heightStream:
				searchParam := types.TxSearchParam{}
				messages, err := p.client.GetMessages(context.Background(), searchParam)
				blockInfoStream <- relayTypes.BlockInfo{
					Height:   height,
					Messages: messages,
					Error:    err,
				}
			}
		}
	}()
	return blockInfoStream
}

func (p *Provider) buildTxFactory() (tx.Factory, error) {
	signMode, ok := signing.SignMode_value[p.client.Context().SignModeStr]
	if !ok {
		return tx.Factory{}, fmt.Errorf("invalid value for sign-mode-str")
	}

	txf := tx.Factory{}.
		WithKeybase(p.client.Context().Keyring).
		WithFeePayer(p.client.Context().FeePayer).
		WithChainID(p.client.Context().ChainID).
		WithSimulateAndExecute(p.client.Context().Simulate).
		WithGasPrices(p.config.GasPrices).
		WithGasAdjustment(p.config.GasAdjustment).
		WithSignMode(signing.SignMode(signMode))

	return txf, nil
}

func (p *Provider) getRawContractMessage(message *relayTypes.Message) (wasmTypes.RawContractMessage, error) {
	switch message.EventType {
	case events.EmitMessage:
		rcvMsg := types.ExecRecvMsg{
			RecvMessage: types.ReceiveMessage{
				SrcNetwork: message.Src,
				ConnSn:     message.Sn,
				Msg:        message.Data,
			},
		}
		rcvMsgByte, err := json.Marshal(rcvMsg)
		if err != nil {
			return nil, err
		}
		return rcvMsgByte, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s ", message.EventType)
	}
}
