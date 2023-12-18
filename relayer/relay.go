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

	SaveHeightMaxAfter = 10
	RouteDuration      = 1 * time.Second
	maxFlushMessage    = 10
	FinalityInterval   = 5 * time.Second

	prefixMessageStore  = "message"
	prefixBlockStore    = "block"
	prefixFinalityStore = "finality"
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

	// once flush completes then only start processing
	if !fresh {
		// flush all the packet and then continue
		relayer.flushMessages(ctx)
	}

	// // create ctx -> with cancel function and senc cancel function to all -> ctx.done():
	// // start all the chain listeners
	go relayer.StartChainListeners(ctx, errorChan)

	// // start all the block processor
	go relayer.StartBlockProcessors(ctx, errorChan)

	// responsible to relaying  messages
	go relayer.StartRouter(ctx, flushInterval)

	// responsible for checking finality
	go relayer.StartFinalityProcessor(ctx)

	return errorChan, nil
}

type Relayer struct {
	log           *zap.Logger
	chains        map[string]*ChainRuntime
	messageStore  *store.MessageStore
	blockStore    *store.BlockStore
	finalityStore *store.FinalityStore
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

	// finality store
	finalityStore := store.NewFinalityStore(db, prefixFinalityStore)

	chainRuntimes := make(map[string]*ChainRuntime, len(chains))
	for _, chain := range chains {
		chainRuntime, err := NewChainRuntime(log, chain)
		if err != nil {
			return nil, err
		}

		lastSavedHeight, err := blockStore.GetLastStoredBlock(chain.NID())
		if err == nil {
			// successfully fetched last savedBlock
			chainRuntime.LastSavedHeight = lastSavedHeight
		}
		chainRuntimes[chain.NID()] = chainRuntime

	}

	return &Relayer{
		log:           log,
		chains:        chainRuntimes,
		messageStore:  messageStore,
		blockStore:    blockStore,
		finalityStore: finalityStore,
	}, nil
}

// GetBlockStore returns the block store
func (r *Relayer) GetBlockStore() *store.BlockStore {
	return r.blockStore
}

// GetBlockStore returns the block store
func (r *Relayer) GetMessageStore() *store.MessageStore {
	return r.messageStore
}

func (r *Relayer) StartChainListeners(
	ctx context.Context,
	errCh chan error,
) {
	var eg errgroup.Group

	for _, chainRuntime := range r.chains {
		chainRuntime := chainRuntime

		eg.Go(func() error {
			// listening to the block
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

func (r *Relayer) StartRouter(ctx context.Context, flushInterval time.Duration) {
	routeTimer := time.NewTicker(RouteDuration)
	flushTimer := time.NewTicker(flushInterval)

	for {
		select {
		case <-flushTimer.C:
			// flushMessage gets all the message from DB
			r.flushMessages(ctx)
		case <-routeTimer.C:
			// processMessage starting working on all the runtime Messages
			r.processMessages(ctx)
		}
	}
}

func (r *Relayer) flushMessages(ctx context.Context) {
	r.log.Info("starting flush logic by adding messages to the messageCache")

	count, err := r.messageStore.TotalCount()
	if err != nil {
		r.log.Warn("error occured when querying total failed delivery message")
	}
	if count == 0 {
		r.log.Debug("no message to flushout")
		return
	}

	for _, chain := range r.chains {
		nId := chain.Provider.NID()
		messages, err := r.getActiveMessagesFromStore(nId, maxFlushMessage)
		if err != nil {
			r.log.Warn("error occured when query messagesFromStore", zap.String("nid", nId), zap.Error(err))
			continue
		}

		if len(messages) == 0 {
			continue
		}
		r.log.Debug(" flushing messages ", zap.String("nid", nId), zap.Int("message count", len(messages)))
		// adding message to messageCache
		// TODO: message with no txHash

		for _, m := range messages {
			chain.MessageCache.Add(m)
		}
	}
}

// TODO: optimize the logic
func (r *Relayer) getActiveMessagesFromStore(nId string, maxMessages int) ([]*types.RouteMessage, error) {
	activeMessages := make([]*types.RouteMessage, 0)

	p := store.NewPagination().GetAll()
	msgs, err := r.messageStore.GetMessages(nId, p)
	if err != nil {
		return nil, err
	}
	for _, m := range msgs {
		if !m.IsStale() {
			activeMessages = append(activeMessages, m)
		}
		if len(activeMessages) > maxMessages {
			break
		}
	}
	return activeMessages, nil
}

func (r *Relayer) processMessages(ctx context.Context) {
	for _, srcChainRuntime := range r.chains {
		for _, routeMessage := range srcChainRuntime.MessageCache.Messages {
			dstChainRuntime, err := r.FindChainRuntime(routeMessage.Dst)
			if err != nil {
				r.log.Error("dst chain runtime not found ", zap.String("dst chain", routeMessage.Dst))
				// remove message if src runtime if dst not found
				r.ClearMessages(ctx, []types.MessageKey{routeMessage.MessageKey()}, srcChainRuntime)
				continue
			}
			if ok := dstChainRuntime.shouldSendMessage(ctx, routeMessage, srcChainRuntime); !ok {
				continue
			}
			go r.RouteMessage(ctx, routeMessage, dstChainRuntime, srcChainRuntime)
		}
	}
}

// processBlockInfo->
// save block height to database
// & merge message to src cache
func (r *Relayer) processBlockInfo(ctx context.Context, srcChainRuntime *ChainRuntime, blockInfo types.BlockInfo) {
	err := r.SaveBlockHeight(ctx, srcChainRuntime, blockInfo.Height, len(blockInfo.Messages))
	if err != nil {
		r.log.Error("unable to save height", zap.Error(err))
	}

	go srcChainRuntime.mergeMessages(ctx, blockInfo.Messages)
}

func (r *Relayer) SaveBlockHeight(ctx context.Context, chainRuntime *ChainRuntime, height uint64, messageCount int) error {

	if messageCount > 0 || (height-chainRuntime.LastSavedHeight) > uint64(SaveHeightMaxAfter) {
		r.log.Debug("saving height:", zap.String("srcChain", chainRuntime.Provider.NID()), zap.Uint64("height", height))
		chainRuntime.LastSavedHeight = height
		err := r.blockStore.StoreBlock(height, chainRuntime.Provider.NID())
		if err != nil {
			return fmt.Errorf("error while saving height of chain:%s %v", chainRuntime.Provider.NID(), err)
		}
	}
	return nil
}

func (r *Relayer) FindChainRuntime(nId string) (*ChainRuntime, error) {
	if chainRuntime, ok := r.chains[nId]; ok {
		return chainRuntime, nil
	}
	return nil, fmt.Errorf("chain runtime not found, nId:%s ", nId)
}

func (r *Relayer) RouteMessage(ctx context.Context, m *types.RouteMessage, dst, src *ChainRuntime) {
	callback := func(key types.MessageKey, response types.TxResponse, err error) {
		// note: it is ok if err is not checked
		if response.Code == types.Success {
			dst.log.Info("successfully relayed message:",
				zap.String("src chain", src.Provider.NID()),
				zap.String("dst chain", dst.Provider.NID()),
				zap.Uint64("Sn number", key.Sn),
				zap.Any("Tx hash", response.TxHash),
			)

			// cannot clear incase of finality block
			if dst.Provider.FinalityBlock(ctx) > 0 {
				routeMessage, ok := src.MessageCache.Messages[key]
				if !ok {
					r.log.Error("message of key not found in messageCache", zap.Any("message key", key))
					return
				}
				txObj := types.NewTransactionObject(*types.NewMessagekeyWithMessageHeight(*key, routeMessage.MessageHeight), response.TxHash, uint64(response.Height))
				if err := r.finalityStore.StoreTxObject(txObj); err != nil {
					r.log.Error("error occured: while storing transaction object in db", zap.Error(err))
				}
				return
			}

			// if success remove message from everywhere
			if err := r.ClearMessages(ctx, []types.MessageKey{key}, src); err != nil {
				r.log.Error("error occured when clearing successful message", zap.Error(err))
			}
			return
		}

		routeMessage, ok := src.MessageCache.Messages[key]
		if !ok {
			r.log.Error("message of key not found in messageCache", zap.Any("key", key))
			return
		}

		r.HandleMessageFailed(routeMessage, dst, src)
	}

	// setting before message is processed
	m.SetIsProcessing(true)
	m.IncrementRetry()

	err := dst.Provider.Route(ctx, m.Message, callback)
	if err != nil {
		dst.log.Error("error occured during message route", zap.Error(err))
		r.HandleMessageFailed(m, dst, src)
	}
}

func (r *Relayer) HandleMessageFailed(routeMessage *types.RouteMessage, dst, src *ChainRuntime) {
	routeMessage.SetIsProcessing(false)

	if routeMessage.GetRetry() != 0 && routeMessage.GetRetry()%uint64(types.DefaultTxRetry) == 0 {
		// save to db
		if err := r.messageStore.StoreMessage(routeMessage); err != nil {
			r.log.Error("error occured when storing the message after max retry", zap.Error(err))
			return
		}

		// removed message from messageCache
		src.MessageCache.Remove(routeMessage.MessageKey())

		dst.log.Error("failed to send message saving to database",
			zap.String("src chain", routeMessage.Src),
			zap.String("dst chain", routeMessage.Dst),
			zap.Uint64("Sn number", routeMessage.Sn),
		)
		return
	}
}

func (r *Relayer) ClearMessages(ctx context.Context, msgs []types.MessageKey, srcChain *ChainRuntime) error {
	// clear from cache
	srcChain.clearMessageFromCache(msgs)

	for _, m := range msgs {
		if err := r.messageStore.DeleteMessage(m); err != nil {
			r.log.Error("error occured when deleting message from db ", zap.Error(err))
			return err
		}
	}
	return nil
}

func (r *Relayer) StartFinalityProcessor(ctx context.Context) {
	ticker := time.NewTicker(FinalityInterval)

	for {
		select {
		case <-ticker.C:
			r.CheckFinality(ctx)
		}
	}

}

func (r *Relayer) CheckFinality(ctx context.Context) {

	for _, c := range r.chains {
		// check for the finality only if finalityblock is provided by the chain
		finalityBlock := c.Provider.FinalityBlock(ctx)
		latestHeight := c.LastBlockHeight
		if finalityBlock > 0 {
			pagination := store.NewPagination().GetAll()
			txObjects, err := r.finalityStore.GetTxObjects(c.Provider.ChainId(), pagination)
			if err != nil {
				r.log.Warn("finality processor: retrive message from store",
					zap.String("chain id ", c.Provider.ChainId()),
					zap.Error(err),
				)
				continue
			}

			for _, txObject := range txObjects {
				if txObject == nil {
					continue
				}
				if txObject.TxHeight == 0 {
					r.log.Warn(" stored  transaction height of txObject cannot be 0 ",
						zap.String("chain-id", c.Provider.ChainId()),
						zap.Any("message key", txObject.MessageKey))
					continue
				}
				// hasn't reached finality
				if txObject.TxHeight+finalityBlock > latestHeight {
					continue
				}

				// check if the txReceipt still exist
				receipt, err := c.Provider.QueryTransactionReceipt(ctx, txObject.TxHash)
				if err != nil {
					r.log.Error("finality processor: queryTransactionReceipt ",
						zap.Any("message key", txObject.MessageKey),
						zap.Error(err))
					continue
				}

				// Transaction Still exist so can be pruned
				if receipt.Status {
					if err := r.finalityStore.DeleteTxObject(&txObject.MessageKey); err != nil {
						r.log.Error("finality processor: deleteTxObject ",
							zap.Any("message key", txObject.MessageKey),
							zap.Error(err))
					}
					continue
				}

				r.log.Info("Transaction Receipt doesn't exist after finalized block, regenerating message",
					zap.Any("message-key", txObject.MessageKey),
					zap.String("tx hash on destination chain", txObject.TxHash))

				// if receipt donot exist generate message again and send to src chain
				srcChainRuntime, ok := r.chains[txObject.Src]
				if !ok {
					r.log.Error("finality processor:  ",
						zap.Any("message key", txObject.MessageKey),
						zap.Error(err))
					continue
				}

				// removing tx object
				if err := r.finalityStore.DeleteTxObject(&txObject.MessageKey); err != nil {
					r.log.Error("finality processor: deleteTxObject ",
						zap.Any("message key", txObject.MessageKey),
						zap.Error(err))
					continue
				}

				// generateMessage
				message, err := srcChainRuntime.Provider.GenerateMessage(ctx, &txObject.MessageKeyWithMessageHeight)
				if err != nil {
					r.log.Error("finality processor: generateMessage",
						zap.Any("message key", txObject.MessageKey),
					)
				}

				// merging message to srcChainRuntime
				srcChainRuntime.mergeMessages(ctx, []*types.Message{message})
			}
		}
	}
}
