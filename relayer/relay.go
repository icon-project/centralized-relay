package relayer

import (
	"context"
	"sync"
	"time"

	"github.com/icon-project/centralized-relay/relayer/provider"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	DefaultFlushInterval      = 5 * time.Minute
	listenerChannelBufferSize = 1000
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

	go relayer.StartListeners(
		ctx,
		flushInterval,
		fresh,
		errorChan,
	)
	go relayer.StartBlockProcessor(ctx, errorChan)
	return errorChan
}

type Relayer struct {
	chains        map[string]*Chain
	log           *zap.Logger
	listenerChans map[string]chan provider.BlockInfo
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

func (r *Relayer) StartBlockProcessor(ctx context.Context, errorChan chan error) {
	var wg sync.WaitGroup

	for chainID, chainChan := range r.listenerChans {
		wg.Add(1)
		go func(id string, ch <-chan provider.BlockInfo) {
			defer wg.Done() // Ensure WaitGroup is decremented when goroutine exits.
			for {
				// Continuously listen to the channel
				select {
				case blockInfo, ok := <-ch:
					if !ok {
						// The channel has been closed, break the loop
						return
					}
					r.processBlockInfo(blockInfo)
				}
			}
		}(chainID, chainChan)
	}

	// Wait for all processing goroutines to finish.
	wg.Wait()

	close(errorChan)
}

func (r *Relayer) StartListeners(
	ctx context.Context,
	flushInterval time.Duration,
	fresh bool,
	errCh chan error,
) {
	var eg errgroup.Group

	runCtx, runCtxCancel := context.WithCancel(ctx)

	for chainId, chain := range r.chains {
		chain := chain
		listnerChan := r.listenerChans[chainId]
		eg.Go(func() error {
			err := chain.ChainProvider.Listener(runCtx, listnerChan)
			runCtxCancel()
			return err
		})
	}

	err := eg.Wait()
	runCtxCancel()
	errCh <- err
}

func (r *Relayer) processBlockInfo(blockInfo provider.BlockInfo) {

}
