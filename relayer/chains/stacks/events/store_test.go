package events

import (
	"testing"
	"time"
)

func TestMemoryEventStore_SaveEvent(t *testing.T) {
	store := NewMemoryEventStore()
	event := &Event{
		ID:          "event1",
		Type:        "TestEvent",
		Data:        "test data",
		BlockHeight: 1,
		Timestamp:   time.Now(),
	}

	err := store.SaveEvent(event)
	if err != nil {
		t.Errorf("SaveEvent returned error: %v", err)
	}

	store.mu.RLock()
	defer store.mu.RUnlock()
	savedEvent, exists := store.events["event1"]
	if !exists {
		t.Errorf("Event not saved")
	}
	if savedEvent != event {
		t.Errorf("Saved event does not match the original")
	}
}

func TestMemoryEventStore_GetEvents(t *testing.T) {
	store := NewMemoryEventStore()
	event1 := &Event{
		ID:          "event1",
		Type:        "TestEvent",
		Data:        "data1",
		BlockHeight: 1,
		Timestamp:   time.Now(),
	}
	event2 := &Event{
		ID:          "event2",
		Type:        "TestEvent",
		Data:        "data2",
		BlockHeight: 2,
		Timestamp:   time.Now(),
	}

	store.SaveEvent(event1)
	store.SaveEvent(event2)

	events, err := store.GetEvents(1)
	if err != nil {
		t.Errorf("GetEvents returned error: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}

	events, err = store.GetEvents(2)
	if err != nil {
		t.Errorf("GetEvents returned error: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}
}

func TestMemoryEventStore_MarkProcessed(t *testing.T) {
	store := NewMemoryEventStore()
	event := &Event{
		ID:          "event1",
		Type:        "TestEvent",
		Data:        "test data",
		BlockHeight: 5,
		Timestamp:   time.Now(),
	}

	store.SaveEvent(event)
	err := store.MarkProcessed("event1")
	if err != nil {
		t.Errorf("MarkProcessed returned error: %v", err)
	}

	store.mu.RLock()
	defer store.mu.RUnlock()
	if !store.processedEvents["event1"] {
		t.Errorf("Event not marked as processed")
	}
	if store.lastProcessedHeight != 5 {
		t.Errorf("Incorrect last processed height, expected 5, got %d", store.lastProcessedHeight)
	}
}

func TestMemoryEventStore_GetLastProcessedHeight(t *testing.T) {
	store := NewMemoryEventStore()

	height, err := store.GetLastProcessedHeight()
	if err != nil {
		t.Errorf("GetLastProcessedHeight returned error: %v", err)
	}
	if height != 0 {
		t.Errorf("Expected height 0, got %d", height)
	}

	store.lastProcessedHeight = 10
	height, err = store.GetLastProcessedHeight()
	if err != nil {
		t.Errorf("GetLastProcessedHeight returned error: %v", err)
	}
	if height != 10 {
		t.Errorf("Expected height 10, got %d", height)
	}
}
