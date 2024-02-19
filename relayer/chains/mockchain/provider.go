package mockchain

import (
	"context"
	"time"

	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/centralized-relay/relayer/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

type MockProviderConfig struct {
	NId             string
	BlockDuration   time.Duration
	SendMessages    map[types.MessageKey]*types.Message
	ReceiveMessages map[types.MessageKey]*types.Message
	StartHeight     uint64
	chainName       string
}

// NewProvider should provide a new Mock provider
func (pp *MockProviderConfig) NewProvider(log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := pp.Validate(); err != nil {
		return nil, err
	}

	pp.chainName = chainName
	return &MockProvider{
		log:    log.With(zap.String("nid", pp.NId), zap.String("chain_name", chainName)),
		PCfg:   pp,
		Height: pp.StartHeight,
	}, nil
}

func (pp *MockProviderConfig) Validate() error {
	return nil
}

type MockProvider struct {
	log    *zap.Logger
	PCfg   *MockProviderConfig
	Height uint64
}

func (p *MockProvider) NID() string {
	return p.PCfg.NId
}

func (p *MockProvider) Init(ctx context.Context) error {
	return nil
}

func (p *MockProvider) FinalityBlock(ctx context.Context) uint64 {
	return 0
}

func (p *MockProvider) Type() string {
	return "evm"
}

func (p *MockProvider) ProviderConfig() provider.Config {
	return p.PCfg
}

func (p *MockProvider) ChainName() string {
	return p.PCfg.chainName
}

func (p *MockProvider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	return p.Height, nil
}

func (p *MockProvider) Listener(ctx context.Context, lastSavedHeight uint64, incoming chan types.BlockInfo) error {
	ticker := time.NewTicker(3 * time.Second)

	if p.Height == 0 {
		if lastSavedHeight != 0 {
			p.Height = lastSavedHeight
		}
	}
	p.log.Info("listening to mock provider from height", zap.Uint64("Height", p.Height))

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:

			height, _ := p.QueryLatestHeight(ctx)
			msgs := p.FindMessages()
			d := types.BlockInfo{
				Height:   uint64(height),
				Messages: msgs,
			}
			incoming <- d
			p.Height += 1
		}
	}
}

func (p *MockProvider) Route(ctx context.Context, message *types.Message, callback types.TxResponseFunc) error {
	p.log.Info("message received", zap.Any("message", message))
	messageKey := message.MessageKey()

	p.DeleteMessage(message)
	callback(messageKey, types.TxResponse{
		Code: types.Success,
	}, nil)
	return nil
}

func (p *MockProvider) FindMessages() []*types.Message {
	messages := make([]*types.Message, 0)
	for _, m := range p.PCfg.SendMessages {
		if m.MessageHeight == p.Height {
			messages = append(messages, m)
		}
	}
	return messages
}

func (p *MockProvider) DeleteMessage(msg *types.Message) {
	var deleteKey types.MessageKey

	for key := range p.PCfg.ReceiveMessages {
		if msg.MessageKey() == key {
			deleteKey = key
			break
		}
	}

	delete(p.PCfg.ReceiveMessages, deleteKey)
}

func (p *MockProvider) ShouldReceiveMessage(ctx context.Context, messagekey types.Message) (bool, error) {
	return true, nil
}

func (p *MockProvider) ShouldSendMessage(ctx context.Context, messageKey types.Message) (bool, error) {
	return true, nil
}

func (p *MockProvider) QueryBalance(ctx context.Context, addr string) (*types.Coin, error) {
	return nil, nil
}

func (p *MockProvider) QueryTransactionReceipt(ctx context.Context, txHash string) (*types.Receipt, error) {
	return nil, nil
}

func (ip *MockProvider) GenerateMessage(ctx context.Context, key *providerTypes.MessageKeyWithMessageHeight) (*providerTypes.Message, error) {
	return nil, nil
}

func (p *MockProvider) MessageReceived(ctx context.Context, key types.MessageKey) (bool, error) {
	return false, nil
}
