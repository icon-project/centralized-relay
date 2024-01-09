package wasm

import (
	"context"
	"fmt"
	abiTypes "github.com/cometbft/cometbft/abci/types"
	"github.com/icon-project/centralized-relay/relayer/chains/wasm/client"
	wasmTypes "github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/concurrency"
	"go.uber.org/zap"
	"runtime"
	"sync"
	"time"
)

const (
	ChainType string = "wasm"
)

type Provider struct {
	logger *zap.Logger
	config *ProviderConfig
	client client.IClient
	txMu   sync.Mutex
}

type ProviderConfig struct {
	ChainName string `json:"-" yaml:"-"`
	ChainID   string `json:"chain_id" yaml:"chain-id"`
	NID       string `json:"nid" yaml:"nid"`

	KeyringBackend  string `json:"keyring_backend" yaml:"keyring-backend"`
	KeyringFilePath string `json:"keyring_file_path" yaml:"keyring-file-path"`
	KeyName         string `json:"key_name" yaml:"key-name"`

	RPCUrl string `json:"rpc-url" yaml:"rpc-url"`

	ContractAddress string `json:"contract-address" yaml:"contract-address"`

	AccountPrefix string `json:"account-prefix" yaml:"account-prefix"`

	GasAdjustment float64 `json:"gas-adjustment" yaml:"gas-adjustment"`
	GasPrices     string  `json:"gas-prices" yaml:"gas-prices"`
	MinGasAmount  uint64  `json:"min-gas-amount" yaml:"min-gas-amount"`
	MaxGasAmount  uint64  `json:"max-gas-amount" yaml:"max-gas-amount"`

	BlockInterval string `json:"block_interval" yaml:"block-interval"`

	SignModeStr      string `json:"sign-mode" yaml:"sign-mode"`
	SigningAlgorithm string `json:"signing-algorithm" yaml:"signing-algorithm"`

	Debug    bool   `json:"debug"`
	HomePath string `json:"home_path"`
}

func (pc ProviderConfig) NewProvider(logger *zap.Logger, homePath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := pc.Validate(); err != nil {
		return nil, err
	}

	pc.ChainName = chainName

	cp := &Provider{
		logger: logger,
	}

	return cp, nil
}

func (pc ProviderConfig) Validate() error {
	if _, err := time.ParseDuration(pc.BlockInterval); err != nil {
		return fmt.Errorf("invalid block-interval: %w", err)
	}
	return nil
}

func (p *Provider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	return p.client.GetLatestBlockHeight(ctx)
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error) {
	res, err := p.client.GetTransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}
	return &types.Receipt{
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
	return *p.config
}

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfo chan types.BlockInfo) error {
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

				numOfPipelines := runtime.NumCPU()

				pipelines := make([]<-chan interface{}, numOfPipelines)

				for i := 0; i < numOfPipelines; i++ {
					pipelines[i] = p.getBlockInfoStream(done, heightStream)
				}

				for bn := range concurrency.FanIn(done, pipelines...) {
					block, ok := bn.(types.BlockInfo)
					if !ok {
						// Todo handle this
					}
					if !block.HasError() {
						blockInfo <- types.BlockInfo{
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

func (p *Provider) Route(ctx context.Context, message *types.Message, callback types.TxResponseFunc) error {
	return nil
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, message types.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, message types.Message) (bool, error) {
	return true, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key types.MessageKey) (bool, error) {

	return false, nil
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*types.Coin, error) {
	coin, err := p.client.GetBalance(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &types.Coin{
		Denom:  coin.Denom,
		Amount: coin.Amount.Uint64(),
	}, nil
}

func (p *Provider) GenerateMessage(ctx context.Context, messageKey *types.MessageKeyWithMessageHeight) (*types.Message, error) {
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
				searchParam := wasmTypes.TxSearchParam{}
				messages, err := p.client.GetMessages(context.Background(), searchParam)
				blockInfoStream <- types.BlockInfo{
					Height:   height,
					Messages: messages,
					Error:    err,
				}
			}
		}
	}()
	return blockInfoStream
}
