package events

import (
	"context"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/interfaces"
	"go.uber.org/zap"
)

type EventProcessor struct {
	handlers      []EventHandler
	store         EventStore
	processChan   chan *Event
	maxWorkers    int
	workerTokens  chan struct{}
	log           *zap.Logger
	ctx           context.Context
	cancel        context.CancelFunc
	client        interfaces.IClient
	senderAddress string
	senderKey     []byte
}

func NewEventProcessor(ctx context.Context, store EventStore, processChan chan *Event, maxWorkers int, log *zap.Logger, client interfaces.IClient, senderAddress string, senderKey []byte) *EventProcessor {
	ctx, cancel := context.WithCancel(ctx)
	return &EventProcessor{
		store:         store,
		processChan:   processChan,
		maxWorkers:    maxWorkers,
		workerTokens:  make(chan struct{}, maxWorkers),
		log:           log,
		client:        client,
		senderAddress: senderAddress,
		senderKey:     senderKey,
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (p *EventProcessor) AddHandler(handler EventHandler) {
	p.handlers = append(p.handlers, handler)
}

func (p *EventProcessor) Start() {
	for i := 0; i < p.maxWorkers; i++ {
		go p.worker()
	}
}

func (p *EventProcessor) Stop() {
	p.cancel()
}

func (p *EventProcessor) worker() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case event := <-p.processChan:
			p.workerTokens <- struct{}{}
			p.processEvent(event)
			<-p.workerTokens
		}
	}
}

func (p *EventProcessor) processEvent(event *Event) {
	p.log.Debug("Processing event",
		zap.String("type", event.Type),
		zap.Any("data", event.Data))

	if err := p.store.SaveEvent(event); err != nil {
		p.log.Error("Failed to save event", zap.Error(err))
		return
	}

	for _, handler := range p.handlers {
		if err := handler(event); err != nil {
			p.log.Error("Handler failed",
				zap.Error(err),
				zap.String("type", event.Type))
			continue
		}
	}

	var err error
	switch event.Type {
	case CallMessageSent:
		p.log.Debug("Processing CallMessageSent event")
		err = p.handleCallMessageSentEvent(event)
	case CallMessage:
		p.log.Debug("Processing CallMessage event")
		err = p.handleCallMessageEvent(event)
	case ResponseMessage:
		p.log.Debug("Processing ResponseMessage event")
		err = p.handleResponseMessageEvent(event)
	case RollbackMessage:
		p.log.Debug("Processing RollbackMessage event")
		err = p.handleRollbackMessageEvent(event)
	default:
		p.log.Warn("No handler for event type", zap.String("type", event.Type))
		return
	}

	if err != nil {
		p.log.Error("Handler failed",
			zap.Error(err),
			zap.String("eventType", event.Type))
		return
	}

	if err := p.store.MarkProcessed(event.ID); err != nil {
		p.log.Error("Failed to mark event as processed", zap.Error(err))
	}
}
