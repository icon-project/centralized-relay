package relayer

import (
	"context"
	"fmt"
	"time"

	"github.com/icon-project/centralized-relay/relayer/store"
	"github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	DefaultFlushInterval      = 5 * time.Minute
	listenerChannelBufferSize = 1000
	DefaultTxRetry            = 2
	SaveHeightMaxAfter        = 1000
	RouteDuration             = 1 * time.Second

	prefixMessageStore = "message"
	prefixBlockStore   = "block"
)

// main start loop
func Start(
	ctx context.Context,
	log *zap.Logger,
	chains map[string]*Chain,
	flushInterval time.Duration,
	fresh bool,
	db store.Store,
) (chan error, error) {
	errorChan := make(chan error, 1)
	relayer, err := NewRelayer(log, db, chains, fresh)
	if err != nil {
		return nil, fmt.Errorf("error creating new relayer %v", err)
	}

	// start all the chain listeners
	go relayer.StartChainListeners(ctx, errorChan)

	// start all the block processor
	go relayer.StartBlockProcessors(ctx, errorChan)

	// responsible to relaying  messages
	go relayer.StartRouter(ctx, flushInterval, fresh)

	return errorChan, nil
}

type Relayer struct {
	log          *zap.Logger
	chains       map[string]*ChainRuntime
	messageStore *store.MessageStore
	blockStore   *store.BlockStore
}

func NewRelayer(log *zap.Logger, db store.Store, chains map[string]*Chain, fresh bool) (*Relayer, error) {

	// if fresh clearing db
	if fresh {
		err := db.ClearStore()
		if err != nil {
			return nil, err
		}
	}

	// initializing message store
	messageStore := store.NewMessageStore(db, prefixMessageStore)

	// blockStore store
	blockStore := store.NewBlockStore(db, prefixBlockStore)

	chainRuntimes := make(map[string]*ChainRuntime, len(chains))
	for _, chain := range chains {
		chainRuntime, err := NewChainRuntime(log, chain)
		if err != nil {
			return nil, err
		}

		lastSavedHeight, err := blockStore.GetLastStoredBlock(chain.ChainID())
		if err == nil {
			// successfully fetched last savedBlock
			chainRuntime.LastSavedHeight = lastSavedHeight
		}
		chainRuntimes[chain.ChainID()] = chainRuntime

	}

	return &Relayer{
		log:          log,
		chains:       chainRuntimes,
		messageStore: messageStore,
		blockStore:   blockStore,
	}, nil
}

func (r *Relayer) StartChainListeners(
	ctx context.Context,
	errCh chan error,
) {

	var eg errgroup.Group
	for _, chainRuntime := range r.chains {
		chainRuntime := chainRuntime

		eg.Go(func() error {

			//
			err := chainRuntime.Provider.Listener(ctx, chainRuntime.LastSavedHeight, chainRuntime.listenerChan)
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		errCh <- err
	}
}

func (r *Relayer) StartBlockProcessors(ctx context.Context, errorChan chan error) {
	var eg errgroup.Group

	for _, chainRuntime := range r.chains {
		chainRuntime := chainRuntime
		listener := chainRuntime.listenerChan
		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case blockInfo, ok := <-listener:
					if !ok {
						return nil
					}
					r.processBlockInfo(ctx, chainRuntime, blockInfo)
				}
			}
		})
	}

	if err := eg.Wait(); err != nil {
		errorChan <- err // Report the error to the error channel.
	}
}

func (r *Relayer) StartRouter(ctx context.Context, flushInterval time.Duration, fresh bool) {

	routeTimer := time.NewTicker(RouteDuration)
	for {
		select {
		case <-routeTimer.C:
			r.processMessages(ctx)
		}
	}
}

func (r *Relayer) processMessages(ctx context.Context) {
	for _, srcChainRuntime := range r.chains {
		for _, routeMessage := range srcChainRuntime.MessageCache {
			dstChainRuntime, err := r.FindChainRuntime(routeMessage.Dst)
			if err != nil {

				// TODO: remove current message as dst chain not found
				continue
			}
			if dstChainRuntime.shouldSendMessage(ctx, routeMessage, srcChainRuntime) {
				go r.RouteMessage(ctx, routeMessage, dstChainRuntime, srcChainRuntime)
			}

		}
	}
}

// processBlockInfo performs these operations
// save block height to database
// send messages to destionation chain
func (r *Relayer) processBlockInfo(ctx context.Context, srcChainRuntime *ChainRuntime, blockInfo types.BlockInfo) {
	err := r.SaveBlockHeight(ctx, srcChainRuntime, blockInfo.Height, len(blockInfo.Messages))
	if err != nil {
		r.log.Error("unable to save height", zap.Error(err))
	}

	go srcChainRuntime.mergeMessages(ctx, blockInfo)
}

func (r *Relayer) SaveBlockHeight(ctx context.Context, chainRuntime *ChainRuntime, height uint64, messageCount int) error {
	r.log.Debug("saving height:", zap.String("srcChain", chainRuntime.Provider.ChainId()), zap.Uint64("height", height))

	if messageCount > 0 || (height-chainRuntime.LastSavedHeight) > uint64(SaveHeightMaxAfter) {
		chainRuntime.LastSavedHeight = height
		// save height to db
		err := r.blockStore.StoreBlock(height, chainRuntime.Provider.ChainId())
		if err != nil {
			return fmt.Errorf("error while saving height of chain:%s %v", chainRuntime.Provider.ChainId(), err)
		}
	}
	return nil
}

func (r *Relayer) FindChainRuntime(chainId string) (*ChainRuntime, error) {
	var chainRuntime *ChainRuntime
	var ok bool

	if chainRuntime, ok = r.chains[chainId]; !ok {
		return nil, fmt.Errorf("chain runtime not found, chainId:%s ", chainId)
	}

	return chainRuntime, nil
}

func (r *Relayer) RouteMessage(ctx context.Context, m *types.RouteMessage, dst, src *ChainRuntime) {

	callback := func(response types.ExecuteMessageResponse) {

		// localization of the variables
		// src := src
		dst := dst
		routeMessage := m

		if response.Code == types.Success {
			// TODO: clearMessage
			dst.log.Info("Successfully relayed message:",
				zap.String("src chain", routeMessage.Src),
				zap.String("dst chain", routeMessage.Dst),
				zap.Uint64("Sn number", routeMessage.Sn),
				zap.Any("Tx hash", response.TxHash),
			)
			return
		}

		if routeMessage.GetRetry() >= uint64(DefaultTxRetry) {
			//TODO: saveMessagetoDB

			dst.log.Error("failed to send message",
				zap.String("src chain", routeMessage.Src),
				zap.String("dst chain", routeMessage.Dst),
				zap.Uint64("Sn number", routeMessage.Sn),
			)
		}

		return
	}

	err := dst.Provider.Route(ctx, m, callback)
	if err != nil {
		dst.log.Error("error occured during message route", zap.Error(err))
	}

}
