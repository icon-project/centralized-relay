package events

import (
	"context"

	"go.uber.org/zap"
)

type EventProcessor struct {
	handlers     []EventHandler
	store        EventStore
	processChan  chan *Event
	maxWorkers   int
	workerTokens chan struct{}
	log          *zap.Logger
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewEventProcessor(ctx context.Context, store EventStore, processChan chan *Event, maxWorkers int, log *zap.Logger) *EventProcessor {
	ctx, cancel := context.WithCancel(ctx)
	return &EventProcessor{
		store:        store,
		processChan:  processChan,
		maxWorkers:   maxWorkers,
		workerTokens: make(chan struct{}, maxWorkers),
		log:          log,
		ctx:          ctx,
		cancel:       cancel,
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
	if err := p.store.SaveEvent(event); err != nil {
		p.log.Error("Failed to save event", zap.Error(err))
		return
	}

	for _, handler := range p.handlers {
		if err := handler(event); err != nil {
			p.log.Error("Handler failed", zap.Error(err))
			continue
		}
	}

	if err := p.store.MarkProcessed(event.ID); err != nil {
		p.log.Error("Failed to mark event as processed", zap.Error(err))
	}
}
