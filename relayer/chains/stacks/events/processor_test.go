package events

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/mocks"
	"github.com/icon-project/stacks-go-sdk/pkg/clarity"
	"github.com/icon-project/stacks-go-sdk/pkg/transaction"
)

func TestEventProcessor_processEvent(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryEventStore()
	processChan := make(chan *Event, 10)
	log := zap.NewNop()
	mockClient := new(mocks.MockClient)

	processor := NewEventProcessor(ctx, store, processChan, 1, log, mockClient, "senderAddress", []byte("senderKey"))

	tests := []struct {
		name        string
		event       *Event
		setupMocks  func()
		expectError bool
	}{
		{
			name: "CallMessageSent event",
			event: &Event{
				ID:   "event1",
				Type: CallMessageSent,
				Data: CallMessageSentData{
					From:         "fromAddress",
					To:           "toAddress",
					Sn:           123,
					Data:         "test data",
					Sources:      []string{"ST1.contract1"},
					Destinations: []string{"dest1"},
				},
				BlockHeight: 1,
				Timestamp:   time.Now(),
			},
			setupMocks: func() {
				mockTx := &transaction.ContractCallTransaction{}

				mockClient.On("MakeContractCall",
					mock.Anything,
					"ST1",
					"contract1",
					"send-message",
					mock.MatchedBy(func(args []clarity.ClarityValue) bool {
						if len(args) != 3 {
							return false
						}
						_, isString := args[0].(*clarity.StringASCII)
						_, isUint := args[1].(*clarity.UInt)
						_, isBuffer := args[2].(*clarity.Buffer)
						return isString && isUint && isBuffer
					}),
					"senderAddress",
					[]byte("senderKey"),
				).Return(mockTx, nil)

				mockClient.On("BroadcastTransaction",
					mock.Anything,
					mockTx,
				).Return("tx-id", nil)
			},
			expectError: false,
		},
		{
			name: "CallMessage event",
			event: &Event{
				ID:   "event2",
				Type: CallMessage,
				Data: CallMessageData{
					From:  "fromAddress",
					To:    "toAddress",
					Sn:    123,
					ReqID: 456,
					Data:  "test data",
				},
				BlockHeight: 2,
				Timestamp:   time.Now(),
			},
			setupMocks:  func() {},
			expectError: false,
		},
		{
			name: "ResponseMessage event",
			event: &Event{
				ID:   "event3",
				Type: ResponseMessage,
				Data: ResponseMessageData{
					Sn:   123,
					Code: 0,
					Msg:  "success",
				},
				BlockHeight: 3,
				Timestamp:   time.Now(),
			},
			setupMocks:  func() {},
			expectError: false,
		},
		{
			name: "RollbackMessage event",
			event: &Event{
				ID:   "event4",
				Type: RollbackMessage,
				Data: RollbackMessageData{
					Sn: 123,
				},
				BlockHeight: 4,
				Timestamp:   time.Now(),
			},
			setupMocks:  func() {},
			expectError: false,
		},
		{
			name: "Invalid event type",
			event: &Event{
				ID:          "event5",
				Type:        "InvalidType",
				Data:        nil,
				BlockHeight: 5,
				Timestamp:   time.Now(),
			},
			setupMocks:  func() {},
			expectError: true,
		},
		{
			name: "Invalid data type for CallMessageSent",
			event: &Event{
				ID:          "event6",
				Type:        CallMessageSent,
				Data:        "invalid data type",
				BlockHeight: 6,
				Timestamp:   time.Now(),
			},
			setupMocks:  func() {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			processor.processEvent(tt.event)

			events, err := store.GetEvents(tt.event.BlockHeight)
			assert.NoError(t, err)
			found := false
			for _, e := range events {
				if e.ID == tt.event.ID {
					found = true
					break
				}
			}
			assert.True(t, found, "Event should be found in store")

			if !tt.expectError {
				store.mu.RLock()
				processed := store.processedEvents[tt.event.ID]
				store.mu.RUnlock()
				assert.True(t, processed, "Event should be marked as processed")
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestEventProcessor_Start_Stop(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryEventStore()
	processChan := make(chan *Event, 10)
	log := zap.NewNop()
	mockClient := new(mocks.MockClient)

	processor := NewEventProcessor(ctx, store, processChan, 2, log, mockClient, "senderAddress", []byte("senderKey"))

	processor.Start()

	event := &Event{
		ID:   "test-event",
		Type: CallMessage,
		Data: CallMessageData{
			From:  "fromAddress",
			To:    "toAddress",
			Sn:    123,
			ReqID: 456,
			Data:  "test data",
		},
		BlockHeight: 1,
		Timestamp:   time.Now(),
	}

	processChan <- event

	time.Sleep(100 * time.Millisecond)

	processor.Stop()

	store.mu.RLock()
	processed := store.processedEvents["test-event"]
	store.mu.RUnlock()
	assert.True(t, processed)
}

func TestEventProcessor_AddHandler(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryEventStore()
	processChan := make(chan *Event, 10)
	log := zap.NewNop()
	mockClient := new(mocks.MockClient)

	processor := NewEventProcessor(ctx, store, processChan, 1, log, mockClient, "senderAddress", []byte("senderKey"))

	handledEvents := make([]*Event, 0)
	var mu sync.Mutex

	handler := func(event *Event) error {
		mu.Lock()
		handledEvents = append(handledEvents, event)
		mu.Unlock()
		return nil
	}

	processor.AddHandler(handler)

	event := &Event{
		ID:   "test-event",
		Type: CallMessage,
		Data: CallMessageData{
			From:  "fromAddress",
			To:    "toAddress",
			Sn:    123,
			ReqID: 456,
			Data:  "test data",
		},
		BlockHeight: 1,
		Timestamp:   time.Now(),
	}

	processor.processEvent(event)

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	assert.Equal(t, 1, len(handledEvents), "Handler should have been called exactly once")
	if len(handledEvents) > 0 {
		assert.Equal(t, event, handledEvents[0], "Handler should have received the correct event")
	}
	mu.Unlock()

	events, err := store.GetEvents(event.BlockHeight)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(events), "Event should be saved in store")

	store.mu.RLock()
	processed := store.processedEvents[event.ID]
	store.mu.RUnlock()
	assert.True(t, processed, "Event should be marked as processed")
}
