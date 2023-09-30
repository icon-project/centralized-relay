package mockchain

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

type MockProviderConfig struct {
	ChainId       string
	BlockDuration time.Duration
	TargetChains  []string
}

// NewProvider should provide a new Mock provider
func (pp *MockProviderConfig) NewProvider(log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {

	if err := pp.Validate(); err != nil {
		return nil, err
	}
	return &MockProvider{
		log:    log.With(zap.String("chain_id", pp.ChainId), zap.String("chain_name", chainName)),
		PCfg:   pp,
		Height: 10,
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
	sn := 1
	src := icp.ChainId()

	icp.log.Info("listening to mock provider")

	for {
		select {
		case <-ticker.C:
			// getting target random
			target := ""
			if len(icp.PCfg.TargetChains) > 0 {
				randomIndex := rand.Intn(len(icp.PCfg.TargetChains))
				target = icp.PCfg.TargetChains[randomIndex]
			}

			message := types.RelayMessage{
				Target: target,
				Src:    src,
				Sn:     uint64(sn),
				Data:   []byte(fmt.Sprintf("message from %s", src)),
			}
			height, _ := icp.QueryLatestHeight(ctx)
			fmt.Printf("found block %d of chain %s  \n", height, icp.ChainId())

			d := types.BlockInfo{
				Height:   uint64(height),
				Messages: []types.RelayMessage{message},
			}

			incoming <- d
			sn += 1
			icp.Height += 1
		}

	}
}

func (icp *MockProvider) Route(ctx context.Context, message *types.RouteMessage, callback func(response types.ExecuteMessageResponse)) error {

	icp.log.Info("message received", zap.Any("message", message))

	callback(types.ExecuteMessageResponse{
		RouteMessage: *message,
		TxResponse: types.TxResponse{
			Code: 0,
		},
	})
	return nil
}
