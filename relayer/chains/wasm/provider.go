package wasm

import (
	"context"
	"fmt"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	abiTypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/client"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/concurrency"
	"go.uber.org/zap"
	"runtime"
	"time"
)

const (
	ChainType string = "wasm"
)

type Provider struct {
	logger *zap.Logger
	config ProviderConfig
	client client.IClient
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
		Status: abiTypes.CodeTypeOK == res.TxResponse.Code,
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
	return ChainType
}

func (p *Provider) ProviderConfig() provider.ProviderConfig {
	return p.config
}

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfo chan relayTypes.BlockInfo) error {
	startHeight, err := p.getStartHeight(ctx, lastSavedHeight)
	if err != nil {
		return err
	}

	latestHeight, err := p.QueryLatestHeight(ctx)
	if err != nil {
		return err
	}

	blockInterval, err := time.ParseDuration(p.config.BlockInterval)
	if err != nil {
		return err
	}

	blockIntervalTicker := time.NewTicker(blockInterval)
	defer blockIntervalTicker.Stop()

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
					if !ok {
						// Todo handle this
					}
					if !block.HasError() {
						blockInfo <- relayTypes.BlockInfo{
							Height: block.Height, Messages: block.Messages,
						}
					}
					//Todo Handle Error
				}
			}()
		}
	}

	return nil
}

func (p *Provider) Route(ctx context.Context, message *relayTypes.Message, callback relayTypes.TxResponseFunc) error {
	txFactory := p.buildTxFactory()
	//Todo customize txFactory: update account number and sequence

	//Build message
	msg := p.getMsgExecuteContract(message)

	//Todo apply retry strategy here
	_, err := p.client.SendTx(ctx, txFactory, []sdkTypes.Msg{&msg})
	if err != nil {
		callback(message.MessageKey(), relayTypes.TxResponse{}, err)
		return err
	}

	return nil
}

func (p *Provider) MessageReceived(ctx context.Context, key relayTypes.MessageKey) (bool, error) {
	_, err := p.client.QuerySmartContract(ctx, p.config.ContractAddress, []byte("hello"))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*relayTypes.Coin, error) {
	coin, err := p.client.GetBalance(ctx, addr, "denomination")
	if err != nil {
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

func (p *Provider) getStartHeight(ctx context.Context, lastSavedHeight uint64) (uint64, error) {
	latestHeight, err := p.client.GetLatestBlockHeight(ctx)
	if err != nil {
		return 0, err
	}

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

func (p *Provider) buildTxFactory() tx.Factory {
	return tx.Factory{}
}

func (p *Provider) getMsgExecuteContract(message *relayTypes.Message) wasmTypes.MsgExecuteContract {
	return wasmTypes.MsgExecuteContract{
		Sender:   p.client.Context().FromAddress.String(),
		Contract: p.config.ContractAddress,
		Msg:      []byte("msg here"),
	}
}
