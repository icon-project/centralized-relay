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
	listenerChan    chan types.BlockInfo
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
		listenerChan: make(chan types.BlockInfo, listenerChannelBufferSize),
		MessageCache: types.NewMessageCache(),
	}, nil
}

func (r *ChainRuntime) mergeMessages(ctx context.Context, messages []*types.Message) {
	if len(messages) == 0 {
		return
	}

	for _, m := range messages {
		routeMessage := types.NewRouteMessage(m)
		r.MessageCache.Add(routeMessage)
	}
}

func (r *ChainRuntime) clearMessageFromCache(msgs []types.MessageKey) {
	for _, m := range msgs {
		r.MessageCache.Remove(m)
	}
}

func (dst *ChainRuntime) shouldSendMessage(ctx context.Context, routeMessage *types.RouteMessage, src *ChainRuntime) bool {
	if routeMessage == nil {
		return false
	}

	if routeMessage.GetIsProcessing() {
		return false
	}

	ok, _ := dst.Provider.ShouldReceiveMessage(ctx, routeMessage.Message)
	if !ok {
		return false
	}

	ok, _ = src.Provider.ShouldSendMessage(ctx, routeMessage.Message)
	if !ok {
		return false
	}

	return true
}
