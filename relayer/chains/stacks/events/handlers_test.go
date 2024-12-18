package events

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/mocks"
	"github.com/icon-project/stacks-go-sdk/pkg/clarity"
	"github.com/icon-project/stacks-go-sdk/pkg/transaction"
)

func TestEventProcessor_handleCallMessageSentEvent(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	mockClient := new(mocks.MockClient)
	store := NewMemoryEventStore()

	processor := NewEventProcessor(ctx, store, make(chan *Event, 10), 1, log, mockClient, "senderAddress", []byte("senderKey"))

	tests := []struct {
		name        string
		event       *Event
		setupMocks  func()
		expectError bool
	}{
		{
			name: "successful call message sent",
			event: &Event{
				ID:   "event1",
				Type: CallMessageSent,
				Data: CallMessageSentData{
					From:         "fromAddr",
					To:           "toAddr",
					Sn:           123,
					Data:         "test data",
					Sources:      []string{"ST1.contract1", "ST2.contract2"},
					Destinations: []string{"dest1"},
				},
			},
			setupMocks: func() {
				mockTx := &transaction.ContractCallTransaction{}
				mockClient.On("MakeContractCall",
					mock.Anything,
					"ST1",
					"contract1",
					"send-message",
					mock.MatchedBy(func(args []clarity.ClarityValue) bool {
						return len(args) == 3
					}),
					"senderAddress",
					[]byte("senderKey"),
				).Return(mockTx, nil).Once()

				mockClient.On("MakeContractCall",
					mock.Anything,
					"ST2",
					"contract2",
					"send-message",
					mock.MatchedBy(func(args []clarity.ClarityValue) bool {
						return len(args) == 3
					}),
					"senderAddress",
					[]byte("senderKey"),
				).Return(mockTx, nil).Once()

				mockClient.On("BroadcastTransaction",
					mock.Anything,
					mockTx,
				).Return("tx-id", nil).Twice()
			},
			expectError: false,
		},
		{
			name: "invalid source contract format",
			event: &Event{
				ID:   "event2",
				Type: CallMessageSent,
				Data: CallMessageSentData{
					From:         "fromAddr",
					To:           "toAddr",
					Sn:           123,
					Data:         "test data",
					Sources:      []string{"invalid-format"},
					Destinations: []string{"dest1"},
				},
			},
			setupMocks:  func() {},
			expectError: true,
		},
		{
			name: "invalid data type",
			event: &Event{
				ID:   "event3",
				Type: CallMessageSent,
				Data: "invalid data type",
			},
			setupMocks:  func() {},
			expectError: true,
		},
		{
			name: "contract call fails",
			event: &Event{
				ID:   "event4",
				Type: CallMessageSent,
				Data: CallMessageSentData{
					From:         "fromAddr",
					To:           "toAddr",
					Sn:           123,
					Data:         "test data",
					Sources:      []string{"ST1.contract1"},
					Destinations: []string{"dest1"},
				},
			},
			setupMocks: func() {
				mockTx := &transaction.ContractCallTransaction{}
				mockClient.On("MakeContractCall",
					mock.Anything,
					"ST1",
					"contract1",
					"send-message",
					mock.MatchedBy(func(args []clarity.ClarityValue) bool {
						return len(args) == 3
					}),
					"senderAddress",
					[]byte("senderKey"),
				).Return(mockTx, fmt.Errorf("contract call failed")).Once()
			},
			expectError: true,
		},
		{
			name: "broadcast fails",
			event: &Event{
				ID:   "event5",
				Type: CallMessageSent,
				Data: CallMessageSentData{
					From:         "fromAddr",
					To:           "toAddr",
					Sn:           123,
					Data:         "test data",
					Sources:      []string{"ST1.contract1"},
					Destinations: []string{"dest1"},
				},
			},
			setupMocks: func() {
				mockTx := &transaction.ContractCallTransaction{}
				mockClient.On("MakeContractCall",
					mock.Anything,
					"ST1",
					"contract1",
					"send-message",
					mock.MatchedBy(func(args []clarity.ClarityValue) bool {
						return len(args) == 3
					}),
					"senderAddress",
					[]byte("senderKey"),
				).Return(mockTx, nil).Once()

				mockClient.On("BroadcastTransaction",
					mock.Anything,
					mockTx,
				).Return("", fmt.Errorf("broadcast failed")).Once()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			mockClient.Calls = nil

			tt.setupMocks()

			err := processor.handleCallMessageSentEvent(tt.event)

			if tt.expectError {
				assert.Error(t, err, "Expected an error but got nil")
			} else {
				assert.NoError(t, err, "Expected no error but got: %v", err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestEventProcessor_handleCallMessageEvent(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	mockClient := new(mocks.MockClient)
	store := NewMemoryEventStore()

	processor := NewEventProcessor(ctx, store, make(chan *Event, 10), 1, log, mockClient, "senderAddress", []byte("senderKey"))

	tests := []struct {
		name        string
		event       *Event
		expectError bool
	}{
		{
			name: "successful call message",
			event: &Event{
				ID:   "event1",
				Type: CallMessage,
				Data: CallMessageData{
					From:  "fromAddr",
					To:    "toAddr",
					Sn:    123,
					ReqID: 456,
					Data:  "test data",
				},
			},
			expectError: false,
		},
		{
			name: "invalid data type",
			event: &Event{
				ID:   "event2",
				Type: CallMessage,
				Data: "invalid data type",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processor.handleCallMessageEvent(tt.event)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEventProcessor_handleResponseMessageEvent(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	mockClient := new(mocks.MockClient)
	store := NewMemoryEventStore()

	processor := NewEventProcessor(ctx, store, make(chan *Event, 10), 1, log, mockClient, "senderAddress", []byte("senderKey"))

	tests := []struct {
		name        string
		event       *Event
		expectError bool
	}{
		{
			name: "successful response message",
			event: &Event{
				ID:   "event1",
				Type: ResponseMessage,
				Data: ResponseMessageData{
					Sn:   123,
					Code: 0,
					Msg:  "success",
				},
			},
			expectError: false,
		},
		{
			name: "invalid data type",
			event: &Event{
				ID:   "event2",
				Type: ResponseMessage,
				Data: "invalid data type",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processor.handleResponseMessageEvent(tt.event)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEventProcessor_handleRollbackMessageEvent(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	mockClient := new(mocks.MockClient)
	store := NewMemoryEventStore()

	processor := NewEventProcessor(ctx, store, make(chan *Event, 10), 1, log, mockClient, "senderAddress", []byte("senderKey"))

	tests := []struct {
		name        string
		event       *Event
		expectError bool
	}{
		{
			name: "successful rollback message",
			event: &Event{
				ID:   "event1",
				Type: RollbackMessage,
				Data: RollbackMessageData{
					Sn: 123,
				},
			},
			expectError: false,
		},
		{
			name: "invalid data type",
			event: &Event{
				ID:   "event2",
				Type: RollbackMessage,
				Data: "invalid data type",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processor.handleRollbackMessageEvent(tt.event)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEventProcessor_callSendMessageFunction(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	mockClient := new(mocks.MockClient)
	store := NewMemoryEventStore()

	processor := NewEventProcessor(ctx, store, make(chan *Event, 10), 1, log, mockClient, "senderAddress", []byte("senderKey"))

	tests := []struct {
		name           string
		sourceContract string
		to             string
		sn             uint64
		msg            string
		setupMocks     func()
		expectError    bool
	}{
		{
			name:           "successful contract call",
			sourceContract: "ST1.contract1",
			to:             "toAddr",
			sn:             123,
			msg:            "test data",
			setupMocks: func() {
				mockTx := &transaction.ContractCallTransaction{}
				mockClient.On("MakeContractCall",
					mock.Anything,
					"ST1",
					"contract1",
					"send-message",
					mock.MatchedBy(func(args []clarity.ClarityValue) bool {
						return len(args) == 3
					}),
					"senderAddress",
					[]byte("senderKey"),
				).Return(mockTx, nil).Once()

				mockClient.On("BroadcastTransaction",
					mock.Anything,
					mockTx,
				).Return("tx-id", nil).Once()
			},
			expectError: false,
		},
		{
			name:           "invalid contract format",
			sourceContract: "invalid-format",
			to:             "toAddr",
			sn:             123,
			msg:            "test data",
			setupMocks:     func() {},
			expectError:    true,
		},
		{
			name:           "contract call fails",
			sourceContract: "ST1.contract1",
			to:             "toAddr",
			sn:             123,
			msg:            "test data",
			setupMocks: func() {
				mockTx := &transaction.ContractCallTransaction{}
				mockClient.On("MakeContractCall",
					mock.Anything,
					"ST1",
					"contract1",
					"send-message",
					mock.MatchedBy(func(args []clarity.ClarityValue) bool {
						return len(args) == 3
					}),
					"senderAddress",
					[]byte("senderKey"),
				).Return(mockTx, fmt.Errorf("contract call failed")).Once()
			},
			expectError: true,
		},
		{
			name:           "broadcast fails",
			sourceContract: "ST1.contract1",
			to:             "toAddr",
			sn:             123,
			msg:            "test data",
			setupMocks: func() {
				mockTx := &transaction.ContractCallTransaction{}
				mockClient.On("MakeContractCall",
					mock.Anything,
					"ST1",
					"contract1",
					"send-message",
					mock.MatchedBy(func(args []clarity.ClarityValue) bool {
						return len(args) == 3
					}),
					"senderAddress",
					[]byte("senderKey"),
				).Return(mockTx, nil).Once()

				mockClient.On("BroadcastTransaction",
					mock.Anything,
					mockTx,
				).Return("", fmt.Errorf("broadcast failed")).Once()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			mockClient.Calls = nil

			tt.setupMocks()

			err := processor.callSendMessageFunction(tt.sourceContract, tt.to, tt.sn, tt.msg)

			if tt.expectError {
				assert.Error(t, err, "Expected an error but got nil")
				if err == nil {
					t.Logf("Expected error for test case: %s", tt.name)
				}
			} else {
				assert.NoError(t, err, "Expected no error but got: %v", err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
