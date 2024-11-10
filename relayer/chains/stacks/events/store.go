package events

import (
	"sync"
)

type MemoryEventStore struct {
	events              map[string]*Event
	processedEvents     map[string]bool
	lastProcessedHeight uint64
	mu                  sync.RWMutex
}

func NewMemoryEventStore() *MemoryEventStore {
	return &MemoryEventStore{
		events:          make(map[string]*Event),
		processedEvents: make(map[string]bool),
	}
}

func (s *MemoryEventStore) SaveEvent(event *Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event
	return nil
}

func (s *MemoryEventStore) GetEvents(fromHeight uint64) ([]*Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []*Event
	for _, event := range s.events {
		if event.BlockHeight >= fromHeight {
			events = append(events, event)
		}
	}
	return events, nil
}

func (s *MemoryEventStore) MarkProcessed(eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event, exists := s.events[eventID]; exists {
		s.processedEvents[eventID] = true
		if event.BlockHeight > s.lastProcessedHeight {
			s.lastProcessedHeight = event.BlockHeight
		}
	}
	return nil
}

func (s *MemoryEventStore) GetLastProcessedHeight() (uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lastProcessedHeight, nil
}
