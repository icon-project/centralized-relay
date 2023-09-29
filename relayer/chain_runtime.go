package relayer

// import "github.com/icon-project/centralized-relay/relayer/store"
import (
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

type ChainRuntime struct {
	Provider        provider.ChainProvider
	listenerChan    chan types.BlockInfo
	log             *zap.Logger
	LastBlockHeight uint64
	LastSavedHeight uint64
}

func NewChainRuntime(log *zap.Logger, chain *Chain) (*ChainRuntime, error) {

	if chain == nil {
		return nil, fmt.Errorf("failed to construct chain runtime")
	}
	return &ChainRuntime{
		log:          log.With(zap.String("chain_id", chain.ChainID())),
		Provider:     chain.ChainProvider,
		listenerChan: make(chan types.BlockInfo, listenerChannelBufferSize),
	}, nil

}
