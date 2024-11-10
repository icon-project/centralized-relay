package events

import (
	"context"

	"go.uber.org/zap"
)

type EventSystem struct {
	listener  *EventListener
	processor *EventProcessor
	store     EventStore
	log       *zap.Logger
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewEventSystem(ctx context.Context, wsURL string, log *zap.Logger) *EventSystem {
	ctx, cancel := context.WithCancel(ctx)

	store := NewMemoryEventStore()
	listener := NewEventListener(ctx, wsURL, 1000, log)
	processor := NewEventProcessor(ctx, store, listener.processChan, 5, log)

	return &EventSystem{
		listener:  listener,
		processor: processor,
		store:     store,
		log:       log,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (s *EventSystem) Start() error {
	if err := s.listener.Start(); err != nil {
		return err
	}
	s.processor.Start()
	return nil
}

func (s *EventSystem) Stop() {
	s.cancel()
	s.listener.Stop()
	s.processor.Stop()
}

func (s *EventSystem) OnEvent(handler EventHandler) {
	s.processor.AddHandler(handler)
}

func (s *EventSystem) GetLastProcessedHeight() (uint64, error) {
	return s.store.GetLastProcessedHeight()
}
