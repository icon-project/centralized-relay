package relayer

import (
	"context"
	"fmt"
	"time"

	"github.com/icon-project/centralized-relay/relayer/provider"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	DefaultFlushInterval      = 5 * time.Minute
	listenerChannelBufferSize = 1000
	DefaultTxRetry            = 5
)

// main start loop
func Start(
	ctx context.Context,
	log *zap.Logger,
	chains map[string]*Chain,
	flushInterval time.Duration,
	fresh bool,
) chan error {
	errorChan := make(chan error, 1)
	relayer := NewRelayer(chains, log)
	go relayer.StartListeners(ctx, flushInterval, fresh, errorChan)
	go relayer.StartBlockProcessors(ctx, errorChan)
	return errorChan
}

type Relayer struct {
	chains          map[string]*Chain
	log             *zap.Logger
	listenerChans   map[string]chan provider.BlockInfo
	lastsavedBlocks map[string]uint64
}

func NewRelayer(chains map[string]*Chain, log *zap.Logger) *Relayer {
	listenerChans := make(map[string]chan provider.BlockInfo, len(chains))
	for chainID := range chains {
		listenerChans[chainID] = make(chan provider.BlockInfo, listenerChannelBufferSize)
	}

	return &Relayer{
		chains:        chains,
		log:           log,
		listenerChans: listenerChans,
	}
}

// processBlockInfo performs these operations
// save block height to database
// send messages to destionation chain
func (r *Relayer) processBlockInfo(ctx context.Context, srcChain string, blockInfo provider.BlockInfo) {

	// saving should not the thread dependent
	// if message > 0 or after certain block

	if len(blockInfo.Messages) > 0 {
		err := r.SaveBlockHeight(ctx, srcChain, blockInfo.Height)
		if err != nil {
			r.log.Error("unable to save height", zap.Error(err))
		}
	}
	go r.RouteMessages(ctx, blockInfo)

}

func (r *Relayer) RouteMessages(ctx context.Context, info provider.BlockInfo) {
	callback := func(response provider.ExecuteMessageResponse) {
		if response.Code == provider.Success {
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
		err := targerchain.ChainProvider.Route(ctx, provider.NewRouteMessage(m), callback)
		if err != nil {
			continue
		}
	}
}

func (r *Relayer) SaveBlockHeight(ctx context.Context, srcChain string, height uint64) error {
	r.log.Debug("saving height:", zap.String("srcChain", srcChain), zap.Uint64("height", height))

	return nil
}

func (r *Relayer) StartBlockProcessors(ctx context.Context, errorChan chan error) {
	var eg errgroup.Group

	for chainID, chainChan := range r.listenerChans {
		srcChain, chainChan := chainID, chainChan // Avoid closure variable capture issue
		eg.Go(func() error {
			for {
				select {
				case blockInfo, ok := <-chainChan:
					fmt.Println("block info received", blockInfo)
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

func (r *Relayer) StartListeners(
	ctx context.Context,
	flushInterval time.Duration,
	fresh bool,
	errCh chan error,
) {
	var eg errgroup.Group

	for chainId, chain := range r.chains {
		chain := chain
		listnerChan := r.listenerChans[chainId]
		eg.Go(func() error {
			err := chain.ChainProvider.Listener(ctx, listnerChan)
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		errCh <- err
	}
}
