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
	DefaultTxRetry            = 5
	SaveHeightMaxAfter        = 1000

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

	go relayer.StartChainListeners(ctx, flushInterval, fresh, errorChan)
	go relayer.StartBlockProcessors(ctx, errorChan)

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
	flushInterval time.Duration,
	fresh bool,
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

	for srcChain, chainRuntime := range r.chains {
		listener := chainRuntime.listenerChan
		srcChain := srcChain
		eg.Go(func() error {
			for {
				select {
				case blockInfo, ok := <-listener:
					if !ok {
						return nil
					}
					r.processBlockInfo(ctx, srcChain, blockInfo)
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})
	}

	if err := eg.Wait(); err != nil {
		errorChan <- err // Report the error to the error channel.
	}
}

// processBlockInfo performs these operations
// save block height to database
// send messages to destionation chain
func (r *Relayer) processBlockInfo(ctx context.Context, srcChain string, blockInfo types.BlockInfo) {
	err := r.SaveBlockHeight(ctx, srcChain, blockInfo.Height, len(blockInfo.Messages))
	if err != nil {
		r.log.Error("unable to save height", zap.Error(err))
	}

	go r.RouteMessages(ctx, blockInfo)
}

func (r *Relayer) RouteMessages(ctx context.Context, info types.BlockInfo) {

	if len(info.Messages) == 0 {
		return
	}

	// TODO: check if the message is already processed then discard it
	// add messages to cache
	callback := func(response types.ExecuteMessageResponse) {
		if response.Code == types.Success {
			if response.GetRetry() > 0 {
				//TODO: remove from DB too
			}
			r.log.Info("Successfully relayed message:",
				zap.String("src chain", response.Src),
				zap.String("dst chain", response.Target),
				zap.Uint64("Sn number", response.Sn),
				zap.Any("Tx hash", response.TxHash),
			)
			return
		}

		if response.GetRetry() >= uint64(DefaultTxRetry) {
			//TODO remove message from db
			r.log.Error("failed to send message",
				zap.String("src chain", response.Src),
				zap.String("dst chain", response.Target),
				zap.Uint64("Sn number", response.Sn),
			)
		}
		//TODO: save message to db
		return
	}

	for _, m := range info.Messages {
		targerchain, ok := r.chains[m.Target]
		if !ok {
			r.log.Error("target chain not present: ", zap.Any("message", m))
			continue
		}

		// should relayMessage
		err := targerchain.Provider.Route(ctx, types.NewRouteMessage(m), callback)
		if err != nil {
			continue
		}
	}
}

func (r *Relayer) SaveBlockHeight(ctx context.Context, srcChain string, height uint64, messageCount int) error {
	r.log.Debug("saving height:", zap.String("srcChain", srcChain), zap.Uint64("height", height))

	srcChainRuntime, ok := r.chains[srcChain]
	if !ok {
		return fmt.Errorf("unable to find source chain")
	}

	if messageCount > 0 || (height-srcChainRuntime.LastSavedHeight) > uint64(SaveHeightMaxAfter) {
		srcChainRuntime.LastSavedHeight = height
		// save height to db
		err := r.blockStore.StoreBlock(height, srcChainRuntime.Provider.ChainId())
		if err != nil {
			return fmt.Errorf("error while saving height of chain:%s %v", srcChainRuntime.Provider.ChainId(), err)
		}

	}
	return nil
}
