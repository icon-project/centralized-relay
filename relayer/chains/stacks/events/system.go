package events

import (
	"context"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/interfaces"
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

func NewEventSystem(ctx context.Context, wsURL string, log *zap.Logger, client interfaces.IClient, senderAddress string, senderKey []byte, contractAddress string) *EventSystem {
	ctx, cancel := context.WithCancel(ctx)

	store := NewMemoryEventStore()
	listener := NewEventListener(ctx, wsURL, 1000, log, contractAddress, client)
	processor := NewEventProcessor(ctx, store, listener.processChan, 5, log, client, senderAddress, senderKey)

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
	s.log.Info("Starting event system")

	if err := s.listener.Start(); err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	s.processor.Start()

	s.log.Info("Event system started successfully")
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
