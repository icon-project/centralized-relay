package relayer

import (
	"context"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

type ChainRuntime struct {
	Provider        provider.ChainProvider
	listenerChan    chan *types.BlockInfo
	log             *zap.Logger
	LastBlockHeight uint64
	LastSavedHeight uint64
	MessageCache    *types.MessageCache
}

func NewChainRuntime(log *zap.Logger, chain *Chain) (*ChainRuntime, error) {
	if chain == nil {
		return nil, fmt.Errorf("failed to construct chain runtime")
	}
	return &ChainRuntime{
		log:          log.With(zap.String("nid ", chain.NID())),
		Provider:     chain.ChainProvider,
		listenerChan: make(chan *types.BlockInfo, listenerChannelBufferSize),
		MessageCache: types.NewMessageCache(),
	}, nil
}

func (r *ChainRuntime) mergeMessages(ctx context.Context, messages []*types.Message) {
	for _, m := range messages {
		routeMessage := types.NewRouteMessage(m)
		r.MessageCache.Add(routeMessage)
	}
}

func (r *ChainRuntime) clearMessageFromCache(msgs []*types.MessageKey) {
	for _, m := range msgs {
		r.MessageCache.Remove(m)
	}
}

func (dst *ChainRuntime) shouldSendMessage(ctx context.Context, routeMessage *types.RouteMessage, src *ChainRuntime) bool {
	if routeMessage == nil {
		return false
	}

	if routeMessage.IsProcessing() {
		return false
	}

	if ok, err := dst.Provider.ShouldReceiveMessage(ctx, routeMessage.Message); !ok || err != nil {
		return false
	}

	if ok, err := src.Provider.ShouldSendMessage(ctx, routeMessage.Message); !ok || err != nil {
		return false
	}

	return true
}

func (r *ChainRuntime) shouldExecuteCall(ctx context.Context, msg *types.RouteMessage) bool {
	return !msg.IsProcessing()
}
