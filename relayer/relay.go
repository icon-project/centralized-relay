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
	listenerChannelBufferSize = 1000 * 5

	HeightSaveInterval         = time.Minute * 5
	maxFlushMessage       uint = 10
	FinalityInterval           = 30 * time.Second
	DeleteExpiredInterval      = 6 * time.Hour
	MessageExpiration          = 24 * time.Hour

	prefixMessageStore  = "message"
	prefixBlockStore    = "block"
	prefixFinalityStore = "finality"
)

// main start loop
func (r *Relayer) Start(ctx context.Context, flushInterval time.Duration, fresh bool) (chan error, error) {
	errorChan := make(chan error, 1)

	// once flush completes then only start processing
	if fresh {
		// flush all the packet and then continue
		r.flushMessages(ctx)
	}

	// // create ctx -> with cancel function and senc cancel function to all -> ctx.done():
	// // start all the chain listeners
	go r.StartChainListeners(ctx, errorChan)

	// // start all the block processor
	go r.StartBlockProcessors(ctx, errorChan)

	// responsible to relaying  messages
	go r.StartRouter(ctx, flushInterval)

	// responsible for checking finality
	go r.StartFinalityProcessor(ctx)

	return errorChan, nil
}

type Relayer struct {
	log           *zap.Logger
	db            store.Store
	chains        map[string]*ChainRuntime
	messageStore  *store.MessageStore
	blockStore    *store.BlockStore
	finalityStore *store.FinalityStore
}

func NewRelayer(log *zap.Logger, db store.Store, chains map[string]*Chain, fresh bool) (*Relayer, error) {
	// if fresh clearing db
	if fresh {
		if err := db.ClearStore(); err != nil {
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
		chainRuntime.Provider.SetLastSavedHeightFunc(func() uint64 {
			return chainRuntime.LastSavedHeight
		})
	}

	return &Relayer{
		log:           log,
		db:            db,
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

func (r *Relayer) StartChainListeners(ctx context.Context, errCh chan error) {
	var eg errgroup.Group

	for _, chainRuntime := range r.chains {
		eg.Go(func() error {
			return chainRuntime.Provider.Listener(ctx, chainRuntime.LastSavedHeight, chainRuntime.listenerChan)
		})
	}
	if err := eg.Wait(); err != nil {
		errCh <- err
	}
}

func (r *Relayer) StartBlockProcessors(ctx context.Context, errorChan chan error) {
	var eg errgroup.Group

	for _, chainRuntime := range r.chains {
		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case blockInfo, ok := <-chainRuntime.listenerChan:
					if !ok {
						return fmt.Errorf("listener channel closed")
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
	routeTimer := time.NewTicker(types.RouteDuration)
	flushTimer := time.NewTicker(1 * time.Second)
	heightTimer := time.NewTicker(HeightSaveInterval)
	cleanMessageTimer := time.NewTicker(1 * time.Second)
	resetTimer := time.NewTicker(3 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-flushTimer.C:
			// flushMessage gets all the message from DB
			go r.flushMessages(ctx)
		case <-routeTimer.C:
			// processMessage starting working on all the runtime Messages
			r.processMessages(ctx)
		case <-heightTimer.C:
			go r.SaveChainsBlockHeight(ctx)
		case <-cleanMessageTimer.C:
			go r.cleanExpiredMessages(ctx)
		case <-resetTimer.C:
			resetTimer.Stop()
			flushTimer.Reset(flushInterval)
			cleanMessageTimer.Reset(DeleteExpiredInterval)
		}
	}
}

func (r *Relayer) flushMessages(ctx context.Context) {
	r.log.Debug("flushing messages from db to cache")
	for _, chain := range r.chains {
		nId := chain.Provider.NID()
		messages, err := r.getActiveMessagesFromStore(nId, maxFlushMessage)
		if err != nil {
			chain.log.Warn("error occured when query messagesFromStore", zap.Error(err))
			continue
		}
		chain.log.Debug("flushing messages", zap.Int("count", len(messages)))
		// adding message to messageCache
		// TODO: message with no txHash

		for _, m := range messages {
			chain.MessageCache.Add(m)
		}
	}
}

// TODO: optimize the logic
func (r *Relayer) getActiveMessagesFromStore(nId string, maxMessages uint) ([]*types.RouteMessage, error) {
	var activeMessages []*types.RouteMessage

	p := store.NewPagination().WithLimit(maxMessages)
	msgs, err := r.messageStore.GetMessages(nId, p)
	if err != nil {
		return nil, err
	}
	for _, m := range msgs {
		if !m.IsStale() {
			activeMessages = append(activeMessages, m)
		}
	}
	return activeMessages, nil
}

func (r *Relayer) processMessages(ctx context.Context) {
	for _, src := range r.chains {
		for key, message := range src.MessageCache.Messages {
			dst, err := r.FindChainRuntime(message.Dst)
			if err != nil {
				r.log.Error("dst chain nid not found", zap.String("nid", message.Dst))
				r.ClearMessages(ctx, []*types.MessageKey{&key}, src)
				continue
			}

			if ok := dst.shouldSendMessage(ctx, message, src); !ok {
				r.log.Debug("processing", zap.Any("message", message))
				continue
			}

			message.ToggleProcessing()

			// if message reached delete the message
			messageReceived, err := dst.Provider.MessageReceived(ctx, &key)
			if err != nil {
				dst.log.Error("error occured when checking message received", zap.String("src", message.Src), zap.Uint64("sn", message.Sn.Uint64()), zap.Error(err))
				message.ToggleProcessing()
				continue
			}

			// if message is received we can remove the message from db
			if messageReceived {
				dst.log.Info("message already received", zap.String("src", message.Src), zap.Uint64("sn", message.Sn.Uint64()))
				r.ClearMessages(ctx, []*types.MessageKey{&key}, src)
				continue
			}
			go r.RouteMessage(ctx, message, dst, src)
		}
	}
}

// processBlockInfo->
// save block height to database
// & merge message to src cache
func (r *Relayer) processBlockInfo(ctx context.Context, src *ChainRuntime, blockInfo *types.BlockInfo) {
	src.LastBlockHeight = blockInfo.Height

	for _, msg := range blockInfo.Messages {
		msg := types.NewRouteMessage(msg)
		src.MessageCache.Add(msg)
		if err := r.messageStore.StoreMessage(msg); err != nil {
			r.log.Error("failed to store a message in db", zap.Error(err))
		}
	}
}

func (r *Relayer) SaveBlockHeight(ctx context.Context, chainRuntime *ChainRuntime, height uint64) error {
	r.log.Debug("saving height:", zap.String("srcChain", chainRuntime.Provider.NID()), zap.Uint64("height", height))
	chainRuntime.LastSavedHeight = height
	chainRuntime.LastBlockHeight = height
	return r.blockStore.StoreBlock(height, chainRuntime.Provider.NID())
}

func (r *Relayer) FindChainRuntime(nId string) (*ChainRuntime, error) {
	if chain, ok := r.chains[nId]; ok {
		return chain, nil
	}
	return nil, fmt.Errorf("chain runtime not found, nId:%s ", nId)
}

func (r *Relayer) GetAllChainsRuntime() []*ChainRuntime {
	var chains []*ChainRuntime
	for _, chainRuntime := range r.chains {
		chains = append(chains, chainRuntime)
	}
	return chains
}

// callback function
func (r *Relayer) callback(ctx context.Context, src, dst *ChainRuntime, key *types.MessageKey) types.TxResponseFunc {
	return func(key *types.MessageKey, response *types.TxResponse, err error) {
		routeMessage, ok := src.MessageCache.Get(key)
		if !ok {
			r.log.Error("key not found in messageCache", zap.Any("key", &key))
			return
		}
		if response.Code == types.Success {
			dst.log.Info("message relayed successfully",
				zap.String("src", src.Provider.NID()),
				zap.String("dst", dst.Provider.NID()),
				zap.String("event_type", routeMessage.EventType),
				zap.Uint64("sn", key.Sn.Uint64()),
				zap.String("tx_hash", response.TxHash),
			)

			// cannot clear incase of finality block
			if dst.Provider.FinalityBlock(ctx) > 0 {
				txObj := types.NewTransactionObject(types.NewMessagekeyWithMessageHeight(key, routeMessage.MessageHeight), response.TxHash, uint64(response.Height))
				r.log.Info("storing txhash to check finality later", zap.Any("txObj", txObj))
				if err := r.finalityStore.StoreTxObject(txObj); err != nil {
					r.log.Error("error occured: while storing transaction object in db", zap.Error(err))
					return
				}
			}
			// if success remove message from everywhere
			if err := r.ClearMessages(ctx, []*types.MessageKey{key}, src); err != nil {
				r.log.Error("error occured when clearing successful message", zap.Error(err))
			}
		}
	}
}

func (r *Relayer) RouteMessage(ctx context.Context, m *types.RouteMessage, dst, src *ChainRuntime) {
	m.IncrementRetry()
	if err := dst.Provider.Route(ctx, m.Message, r.callback(ctx, src, dst, m.MessageKey())); err != nil {
		dst.log.Error("message routing failed", zap.String("src", m.Src), zap.String("event_type", m.EventType), zap.Error(err))
		r.HandleMessageFailed(m, dst, src)
	}
}

func (r *Relayer) HandleMessageFailed(routeMessage *types.RouteMessage, dst, src *ChainRuntime) {
	routeMessage.ToggleProcessing()
	if routeMessage.Retry >= types.MaxTxRetry {
		if err := r.messageStore.StoreMessage(routeMessage); err != nil {
			r.log.Error("error occured when storing the message after max retry", zap.Error(err))
			return
		}

		// removed message from messageCache
		src.MessageCache.Remove(routeMessage.MessageKey())

		dst.log.Error("message relay failed",
			zap.String("src", routeMessage.Src),
			zap.String("dst", routeMessage.Dst),
			zap.Uint64("sn", routeMessage.Sn.Uint64()),
			zap.String("event_type", routeMessage.EventType),
			zap.Uint8("count", routeMessage.Retry),
		)
	}
	return
}

// PruneDB removes all the messages from db
func (r *Relayer) PruneDB() error {
	return r.db.ClearStore()
}

func (r *Relayer) ClearMessages(ctx context.Context, msgs []*types.MessageKey, srcChain *ChainRuntime) error {
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
	for nid, c := range r.chains {
		// check for the finality only if finalityblock is provided by the chain
		finalityBlock := c.Provider.FinalityBlock(ctx)
		latestHeight := c.LastBlockHeight
		if finalityBlock > 0 {
			pagination := store.NewPagination().WithLimit(10)
			txObjects, err := r.finalityStore.GetTxObjects(nid, pagination)
			if err != nil {
				r.log.Warn("finality processor: retrive message from store",
					zap.String("nid", nid),
					zap.Error(err),
				)
				continue
			}

			for _, txObject := range txObjects {
				r.log.Debug("checking finality for tx object", zap.Any("txobj", txObjects), zap.Uint64("latest height", latestHeight))
				if txObject == nil {
					continue
				}
				if txObject.TxHeight == 0 {
					r.log.Warn("stored  transaction height of txObject cannot be 0 ",
						zap.String("nid", c.Provider.NID()),
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
					if err := r.finalityStore.DeleteTxObject(txObject.MessageKey); err != nil {
						r.log.Error("finality processor: deleteTxObject ",
							zap.Any("message key", txObject.MessageKey),
							zap.Error(err))
					}
					r.log.Debug("finality processor: transaction still exist after finalized block, deleting txObject")
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
				if err := r.finalityStore.DeleteTxObject(txObject.MessageKey); err != nil {
					r.log.Error("finality processor: deleteTxObject ",
						zap.Any("message key", txObject.MessageKey),
						zap.Error(err))
					continue
				}

				// generateMessage
				messages, err := srcChainRuntime.Provider.GenerateMessages(ctx, txObject.MessageKeyWithMessageHeight)
				if err != nil {
					r.log.Error("finality processor: generateMessage",
						zap.Any("message key", txObject.MessageKey),
						zap.Error(err),
					)
					continue
				}

				// merging message to srcChainRuntime
				srcChainRuntime.mergeMessages(ctx, messages)
			}
		}
	}
}

// SaveBlockHeight for all chains
func (r *Relayer) SaveChainsBlockHeight(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	for nid, chain := range r.chains {
		height, err := chain.Provider.QueryLatestHeight(ctx)
		if err != nil {
			r.log.Error("error occured when querying latest height", zap.String("nid", nid), zap.Error(err))
			continue
		}
		if err := r.SaveBlockHeight(ctx, chain, height); err != nil {
			r.log.Error("error occured when saving block height", zap.String("nid", nid), zap.Error(err))
			continue
		}
	}
}

// cleanExpiredMessages
func (r *Relayer) cleanExpiredMessages(ctx context.Context) {
	for nid, chain := range r.chains {
		p := store.NewPagination().WithLimit(maxFlushMessage)
		messages, err := r.messageStore.GetMessages(nid, p)
		if err != nil {
			r.log.Error("error occured when fetching messages from db", zap.Error(err))
			continue
		}

		for _, m := range messages {
			if m.IsElasped(MessageExpiration) {
				if err := r.ClearMessages(ctx, []*types.MessageKey{m.MessageKey()}, chain); err != nil {
					r.log.Error("error occured when clearing expired message", zap.Error(err))
				}
			}
		}
	}
}
