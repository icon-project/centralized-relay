package mockchain

import (
	"context"
	"reflect"
	"time"

	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

type MockProviderConfig struct {
	ChainId         string
	BlockDuration   time.Duration
	SendMessages    map[types.MessageKey]types.Message
	ReceiveMessages map[types.MessageKey]types.Message
	StartHeight     uint64
}

// NewProvider should provide a new Mock provider
func (pp *MockProviderConfig) NewProvider(log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {

	if err := pp.Validate(); err != nil {
		return nil, err
	}
	return &MockProvider{
		log:    log.With(zap.String("chain_id", pp.ChainId), zap.String("chain_name", chainName)),
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

func (icp *MockProvider) ChainId() string {
	return icp.PCfg.ChainId
}
func (icp *MockProvider) Init(ctx context.Context) error {
	return nil
}

func (icp *MockProvider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	return icp.Height, nil
}

func (icp *MockProvider) Listener(ctx context.Context, lastSavedHeight uint64, incoming chan types.BlockInfo) error {

	ticker := time.NewTicker(3 * time.Second)

	if icp.Height == 0 {
		if lastSavedHeight != 0 {
			icp.Height = lastSavedHeight
		}
	}
	icp.log.Info("listening to mock provider from height", zap.Uint64("Height", icp.Height))

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:

			height, _ := icp.QueryLatestHeight(ctx)
			msgs := icp.FindMessages()
			d := types.BlockInfo{
				Height:   uint64(height),
				Messages: msgs,
			}
			incoming <- d
			icp.Height += 1
		}

	}
}

func (icp *MockProvider) Route(ctx context.Context, message *types.RouteMessage, callback types.TxResponseFunc) error {

	icp.log.Info("message received", zap.Any("message", message))
	messageKey := message.MessageKey()

	icp.DeleteMessage(message)
	callback(messageKey, types.TxResponse{
		Code: types.Success,
	}, nil)
	return nil
}

func (icp *MockProvider) FindMessages() []types.Message {
	messages := make([]types.Message, 0)
	for _, m := range icp.PCfg.SendMessages {
		if m.MessageHeight == icp.Height {
			messages = append(messages, m)
		}

	}
	return messages
}

func (icp *MockProvider) DeleteMessage(routeMsg *types.RouteMessage) {

	if routeMsg == nil {
		return
	}
	var deleteKey *types.MessageKey

	for key, m := range icp.PCfg.ReceiveMessages {
		fromRouteMessage := routeMsg.GetMessage()

		if reflect.DeepEqual(fromRouteMessage, m) {
			deleteKey = &key
			break
		}
	}

	if deleteKey != nil {
		delete(icp.PCfg.ReceiveMessages, *deleteKey)
	}

}

func (icp *MockProvider) ShouldReceiveMessage(ctx context.Context, messagekey types.Message) (bool, error) {
	return true, nil

}
func (icp *MockProvider) ShouldSendMessage(ctx context.Context, messageKey types.Message) (bool, error) {
	return true, nil
}

func (icp *MockProvider) QueryBalance(ctx context.Context, addr string) (*types.Coin, error) {
	return nil, nil
}
