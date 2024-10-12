package stacks

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/interfaces"
	"github.com/icon-project/centralized-relay/relayer/chains/stacks/mocks"
	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/provider"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

func setupProvider(t *testing.T, mockClient *mocks.MockClient, cfg *Config) *Provider {
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	p, err := cfg.NewProvider(context.Background(), logger, "/tmp/relayer", false, "stacks_testnet")
	assert.NoError(t, err)
	assert.NotNil(t, p)

	providerInstance := p.(*Provider)
	providerInstance.client = mockClient
	providerInstance.contracts = cfg.eventMap()
	return providerInstance
}

func TestListener_Success(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockClient := new(mocks.MockClient)

	emitMsgEvent := EmitMessageEvent{
		TargetNetwork: "stacks_testnet",
		Sn:            "12345",
		Msg:           "Hello, Stacks!",
	}

	callMsgEvent := CallMessageEvent{
		ReqID: "67890",
		Sn:    "54321",
		Data:  "0xabcdef",
	}

	mockClient.On("SubscribeToEvents", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		callback := args.Get(2).(interfaces.EventCallback)

		go func() {
			err := callback("message_event", emitMsgEvent)
			assert.NoError(t, err, "Callback for EmitMessageEvent should not error")

			err = callback("call_message_event", callMsgEvent)
			assert.NoError(t, err, "Callback for CallMessageEvent should not error")
		}()
	}).Return(nil)

	cfg := &Config{
		CommonConfig: provider.CommonConfig{
			ChainName: "stacks_testnet",
			RPCUrl:    "https://stacks-node-api.example.com",
			Contracts: providerTypes.ContractConfigMap{
				"xcall":      "ST000000000000000000002AMW42H",
				"connection": "ST000000000000000000002AMW42H",
			},
			NID: "stacks_testnet",
		},
	}

	providerInstance := setupProvider(t, mockClient, cfg)

	blockInfoChan := make(chan *providerTypes.BlockInfo, 2)

	go func() {
		err := providerInstance.Listener(ctx, providerTypes.LastProcessedTx{}, blockInfoChan)
		if err != nil {
			t.Logf("Listener exited with error: %v", err)
		}
	}()

	expectedMessages := map[string]*providerTypes.Message{
		events.EmitMessage: {
			Dst:           "stacks_testnet",
			Src:           "stacks_testnet",
			Sn:            big.NewInt(12345),
			MessageHeight: 0, // todo: update
			EventType:     events.EmitMessage,
			Data:          []byte("Hello, Stacks!"),
		},
		events.CallMessage: {
			Dst:           "stacks_testnet",
			Src:           "stacks_testnet",
			Sn:            big.NewInt(54321),
			MessageHeight: 0, // todo: update
			EventType:     events.CallMessage,
			Data:          []byte("0xabcdef"),
			ReqID:         big.NewInt(67890),
		},
	}

	receivedMessages := make(map[string]*providerTypes.Message)

	timeout := time.After(2 * time.Second)
	for i := 0; i < 2; i++ {
		select {
		case blockInfo := <-blockInfoChan:
			assert.NotNil(t, blockInfo, "Received BlockInfo should not be nil")
			assert.Equal(t, uint64(0), blockInfo.Height, "Block height should be 0")
			assert.Len(t, blockInfo.Messages, 1, "BlockInfo should contain one message")

			msg := blockInfo.Messages[0]
			expectedMsg, exists := expectedMessages[msg.EventType]
			assert.True(t, exists, "Unexpected event type received: %s", msg.EventType)

			assert.Equal(t, expectedMsg.Dst, msg.Dst, "Dst field mismatch")
			assert.Equal(t, expectedMsg.Src, msg.Src, "Src field mismatch")
			assert.True(t, msg.Sn.Cmp(expectedMsg.Sn) == 0, "Sn field mismatch")
			assert.Equal(t, expectedMsg.EventType, msg.EventType, "EventType field mismatch")
			assert.Equal(t, expectedMsg.Data, msg.Data, "Data field mismatch")

			if msg.EventType == events.CallMessage {
				assert.True(t, msg.ReqID.Cmp(expectedMsg.ReqID) == 0, "ReqID field mismatch for CallMessageEvent")
			}

			receivedMessages[msg.EventType] = msg

		case <-timeout:
			t.Fatal("Timed out waiting for BlockInfo messages")
		}
	}

	for eventType, expectedMsg := range expectedMessages {
		receivedMsg, exists := receivedMessages[eventType]
		assert.True(t, exists, "Expected to receive event type: %s", eventType)
		assert.Equal(t, expectedMsg.Dst, receivedMsg.Dst, "Dst field mismatch for event type: %s", eventType)
		assert.Equal(t, expectedMsg.Src, receivedMsg.Src, "Src field mismatch for event type: %s", eventType)
		assert.True(t, receivedMsg.Sn.Cmp(expectedMsg.Sn) == 0, "Sn field mismatch for event type: %s", eventType)
		assert.Equal(t, expectedMsg.EventType, receivedMsg.EventType, "EventType field mismatch for event type: %s", eventType)
		assert.Equal(t, expectedMsg.Data, receivedMsg.Data, "Data field mismatch for event type: %s", eventType)

		if eventType == events.CallMessage {
			assert.True(t, receivedMsg.ReqID.Cmp(expectedMsg.ReqID) == 0, "ReqID field mismatch for CallMessageEvent")
		}
	}

	cancel()

	time.Sleep(100 * time.Millisecond)

	select {
	case blockInfo := <-blockInfoChan:
		t.Errorf("Received unexpected BlockInfo after cancellation: %+v", blockInfo)
	default:
	}

	mockClient.AssertNumberOfCalls(t, "SubscribeToEvents", 1)
}

func TestListener_CallbackError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockClient := new(mocks.MockClient)

	invalidEmitMsgEvent := EmitMessageEvent{
		TargetNetwork: "stacks_testnet",
		Sn:            "invalid_sn",
		Msg:           "Hello, Stacks!",
	}

	mockClient.On("SubscribeToEvents", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		callback := args.Get(2).(interfaces.EventCallback)

		go func() {
			err := callback("message_event", invalidEmitMsgEvent)
			assert.Error(t, err, "Callback should return an error due to invalid SN")
		}()
	}).Return(nil)

	cfg := &Config{
		CommonConfig: provider.CommonConfig{
			ChainName: "stacks_testnet",
			RPCUrl:    "https://stacks-node-api.example.com",
			Contracts: providerTypes.ContractConfigMap{
				"xcall":      "ST000000000000000000002AMW42H",
				"connection": "ST000000000000000000002AMW42H",
			},
			NID: "stacks_testnet",
		},
	}

	providerInstance := setupProvider(t, mockClient, cfg)

	blockInfoChan := make(chan *providerTypes.BlockInfo, 1)

	go func() {
		err := providerInstance.Listener(ctx, providerTypes.LastProcessedTx{}, blockInfoChan)
		if err != nil {
			t.Logf("Listener exited with error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	select {
	case blockInfo := <-blockInfoChan:
		t.Errorf("Received unexpected BlockInfo due to callback error: %+v", blockInfo)
	default:
	}

	mockClient.AssertNumberOfCalls(t, "SubscribeToEvents", 1)
}

func TestListener_InvalidEvent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockClient := new(mocks.MockClient)

	unknownEventType := "unknown_event_type"
	unknownEventData := map[string]interface{}{
		"some_field": "some_value",
	}

	mockClient.On("SubscribeToEvents", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		callback := args.Get(2).(interfaces.EventCallback)

		go func() {
			err := callback(unknownEventType, unknownEventData)
			assert.Error(t, err, "Callback should return an error for unknown event type")
		}()
	}).Return(nil)

	cfg := &Config{
		CommonConfig: provider.CommonConfig{
			ChainName: "stacks_testnet",
			RPCUrl:    "https://stacks-node-api.example.com",
			Contracts: providerTypes.ContractConfigMap{
				"xcall":      "ST000000000000000000002AMW42H",
				"connection": "ST000000000000000000002AMW42H",
			},
			NID: "stacks_testnet",
		},
	}

	providerInstance := setupProvider(t, mockClient, cfg)

	blockInfoChan := make(chan *providerTypes.BlockInfo, 1)

	go func() {
		err := providerInstance.Listener(ctx, providerTypes.LastProcessedTx{}, blockInfoChan)
		if err != nil {
			t.Logf("Listener exited with error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	select {
	case blockInfo := <-blockInfoChan:
		t.Errorf("Received unexpected BlockInfo for unknown event type: %+v", blockInfo)
	default:
	}

	mockClient.AssertNumberOfCalls(t, "SubscribeToEvents", 1)
}

func TestListener_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	mockClient := new(mocks.MockClient)

	emitMsgEvent := EmitMessageEvent{
		TargetNetwork: "stacks_testnet",
		Sn:            "12345",
		Msg:           "Hello, Stacks!",
	}

	mockClient.On("SubscribeToEvents", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		callback := args.Get(2).(interfaces.EventCallback)

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					callback("message_event", emitMsgEvent)
					time.Sleep(50 * time.Millisecond)
				}
			}
		}()
	}).Return(nil)

	cfg := &Config{
		CommonConfig: provider.CommonConfig{
			ChainName: "stacks_testnet",
			RPCUrl:    "https://stacks-node-api.example.com",
			Contracts: providerTypes.ContractConfigMap{
				"xcall":      "ST000000000000000000002AMW42H",
				"connection": "ST000000000000000000002AMW42H",
			},
			NID: "stacks_testnet",
		},
	}

	providerInstance := setupProvider(t, mockClient, cfg)

	blockInfoChan := make(chan *providerTypes.BlockInfo, 10)

	go func() {
		err := providerInstance.Listener(ctx, providerTypes.LastProcessedTx{}, blockInfoChan)
		if err != nil {
			t.Logf("Listener exited with error: %v", err)
		}
	}()

	time.Sleep(200 * time.Millisecond)

	cancel()

	time.Sleep(100 * time.Millisecond)

	select {
	case blockInfo := <-blockInfoChan:
		t.Logf("Received BlockInfo before cancellation: %+v", blockInfo)
	default:
	}

	mockClient.AssertNumberOfCalls(t, "SubscribeToEvents", 1)
}
