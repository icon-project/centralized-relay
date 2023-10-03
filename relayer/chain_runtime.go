package relayer

// import "github.com/icon-project/centralized-relay/relayer/store"
import (
	"context"
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
	MessageCache    types.MessageCache
}

func NewChainRuntime(log *zap.Logger, chain *Chain) (*ChainRuntime, error) {
	if chain == nil {
		return nil, fmt.Errorf("failed to construct chain runtime")
	}
	return &ChainRuntime{
		log:          log.With(zap.String("chain_id", chain.ChainID())),
		Provider:     chain.ChainProvider,
		listenerChan: make(chan types.BlockInfo, listenerChannelBufferSize),
		MessageCache: make(types.MessageCache),
	}, nil

}

func (r *ChainRuntime) mergeMessages(ctx context.Context, info types.BlockInfo) {
	if len(info.Messages) == 0 {
		return
	}

	for _, m := range info.Messages {
		routeMessage := types.NewRouteMessage(m)
		r.MessageCache.Add(routeMessage)
	}
}

func (r *ChainRuntime) shouldSendMessage(ctx context.Context, routeMessage *types.RouteMessage, src *ChainRuntime) bool {
	if routeMessage == nil {
		return false
	}

	if routeMessage.GetIsProcessing() {
		return false
	}

	// TODO: filter from the src

	// TODO: filter for dst

	// TODO: query if the message is received to destination
	// if received remove from the destination

	return true
}
