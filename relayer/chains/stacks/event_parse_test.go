package stacks

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/provider"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

func getTestLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func TestParseEmitMessageEvent(t *testing.T) {
	logger := getTestLogger()
	nid := "stacks"

	cfg := &Config{
		CommonConfig: provider.CommonConfig{
			RPCUrl: "https://stacks-node-api.example.com",
			Contracts: providerTypes.ContractConfigMap{
				"XcallContract":      "ST000000000000000000002AMW42H",
				"ConnectionContract": "ST000000000000000000002AMW42H",
			},
			NID: nid,
		},
	}

	p, err := cfg.NewProvider(context.Background(), logger, "/tmp/relayer", false, nid)
	assert.NoError(t, err, "Expected no error during provider initialization")
	assert.NotNil(t, p, "Expected Provider to be initialized")

	providerConcrete, ok := p.(*Provider)
	assert.True(t, ok, "Expected ChainProvider to be of type *Provider")

	tests := []struct {
		name        string
		eventData   interface{}
		expectedMsg *providerTypes.Message
		wantErr     bool
	}{
		{
			name: "Valid EmitMessageEvent",
			eventData: EmitMessageEvent{
				TargetNetwork: "stacks_testnet",
				Sn:            "12345",
				Msg:           "Hello, Stacks!",
			},
			expectedMsg: &providerTypes.Message{
				Dst:           "stacks_testnet",
				Src:           "stacks",
				Sn:            big.NewInt(12345),
				MessageHeight: 0, // To be set based on block information
				EventType:     events.EmitMessage,
				Data:          []byte("Hello, Stacks!"),
			},
			wantErr: false,
		},
		{
			name: "Invalid Sn in EmitMessageEvent",
			eventData: map[string]interface{}{
				"targetNetwork": "stacks_testnet",
				"sn":            "not_a_number",
				"msg":           "Hello, Stacks!",
			},
			expectedMsg: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := providerConcrete.getRelayMessageFromEvent("message_event", tt.eventData)
			if tt.wantErr {
				assert.Error(t, err, "Expected an error")
				assert.Nil(t, msg, "Expected Message to be nil")
			} else {
				assert.NoError(t, err, "Expected no error")
				assert.NotNil(t, msg, "Expected a Message object")

				assert.Equal(t, tt.expectedMsg.Dst, msg.Dst, "Dst should match")
				assert.Equal(t, tt.expectedMsg.Src, msg.Src, "Src should match")
				assert.Equal(t, tt.expectedMsg.Sn, msg.Sn, "Sn should match")
				assert.Equal(t, tt.expectedMsg.EventType, msg.EventType, "EventType should match")
				assert.Equal(t, tt.expectedMsg.Data, msg.Data, "Data should match")
			}
		})
	}
}

func TestParseCallMessageEvent(t *testing.T) {
	logger := getTestLogger()
	nid := "stacks"

	cfg := &Config{
		CommonConfig: provider.CommonConfig{
			RPCUrl: "https://stacks-node-api.example.com",
			Contracts: providerTypes.ContractConfigMap{
				"XcallContract":      "ST000000000000000000002AMW42H",
				"ConnectionContract": "ST000000000000000000002AMW42H",
			},
			NID: nid,
		},
	}

	p, err := cfg.NewProvider(context.Background(), logger, "/tmp/relayer", false, nid)
	assert.NoError(t, err, "Expected no error during provider initialization")
	assert.NotNil(t, p, "Expected Provider to be initialized")

	providerConcrete, ok := p.(*Provider)
	assert.True(t, ok, "Expected ChainProvider to be of type *Provider")

	tests := []struct {
		name        string
		eventData   interface{}
		expectedMsg *providerTypes.Message
		wantErr     bool
	}{
		{
			name: "Valid CallMessageEvent",
			eventData: CallMessageEvent{
				ReqID: "67890",
				Sn:    "54321",
				Data:  "0xabcdef",
			},
			expectedMsg: &providerTypes.Message{
				Dst:           "stacks",
				Src:           "stacks",
				Sn:            big.NewInt(54321),
				MessageHeight: 0, // To be set based on block information
				EventType:     events.CallMessage,
				Data:          []byte("0xabcdef"),
				ReqID:         big.NewInt(67890),
			},
			wantErr: false,
		},
		{
			name: "Invalid Sn in CallMessageEvent",
			eventData: map[string]interface{}{
				"req_id": "invalid_req_id",
				"sn":     "invalid_sn",
				"data":   "0xabcdef",
			},
			expectedMsg: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := providerConcrete.getRelayMessageFromEvent("call_message_event", tt.eventData)
			if tt.wantErr {
				assert.Error(t, err, "Expected an error")
				assert.Nil(t, msg, "Expected Message to be nil")
			} else {
				assert.NoError(t, err, "Expected no error")
				assert.NotNil(t, msg, "Expected a Message object")

				assert.Equal(t, tt.expectedMsg.Dst, msg.Dst, "Dst should match")
				assert.Equal(t, tt.expectedMsg.Src, msg.Src, "Src should match")
				assert.Equal(t, tt.expectedMsg.Sn, msg.Sn, "Sn should match")
				assert.Equal(t, tt.expectedMsg.ReqID, msg.ReqID, "ReqID should match")
				assert.Equal(t, tt.expectedMsg.EventType, msg.EventType, "EventType should match")
				assert.Equal(t, tt.expectedMsg.Data, msg.Data, "Data should match")
			}
		})
	}
}

func TestParseRollbackMessageEvent(t *testing.T) {
	logger := getTestLogger()
	nid := "stacks"

	cfg := &Config{
		CommonConfig: provider.CommonConfig{
			RPCUrl: "https://stacks-node-api.example.com",
			Contracts: providerTypes.ContractConfigMap{
				"XcallContract":      "ST000000000000000000002AMW42H",
				"ConnectionContract": "ST000000000000000000002AMW42H",
			},
			NID: nid,
		},
	}

	p, err := cfg.NewProvider(context.Background(), logger, "/tmp/relayer", false, nid)
	assert.NoError(t, err, "Expected no error during provider initialization")
	assert.NotNil(t, p, "Expected Provider to be initialized")

	providerConcrete, ok := p.(*Provider)
	assert.True(t, ok, "Expected ChainProvider to be of type *Provider")

	tests := []struct {
		name        string
		eventData   interface{}
		expectedMsg *providerTypes.Message
		wantErr     bool
	}{
		{
			name: "Valid RollbackMessageEvent",
			eventData: RollbackMessageEvent{
				Sn: "112233",
			},
			expectedMsg: &providerTypes.Message{
				Dst:           "stacks",
				Src:           "stacks",
				Sn:            big.NewInt(112233),
				MessageHeight: 0,
				EventType:     events.RollbackMessage,
				Data:          nil,
			},
			wantErr: false,
		},
		{
			name: "Invalid Sn in RollbackMessageEvent",
			eventData: map[string]interface{}{
				"sn": "invalid_sn",
			},
			expectedMsg: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := providerConcrete.getRelayMessageFromEvent("rollback_message_event", tt.eventData)
			if tt.wantErr {
				assert.Error(t, err, "Expected an error")
				assert.Nil(t, msg, "Expected Message to be nil")
			} else {
				assert.NoError(t, err, "Expected no error")
				assert.NotNil(t, msg, "Expected a Message object")

				assert.Equal(t, tt.expectedMsg.Dst, msg.Dst, "Dst should match")
				assert.Equal(t, tt.expectedMsg.Src, msg.Src, "Src should match")
				assert.Equal(t, tt.expectedMsg.Sn, msg.Sn, "Sn should match")
				assert.Equal(t, tt.expectedMsg.EventType, msg.EventType, "EventType should match")
				assert.Nil(t, msg.Data, "Data should be nil for RollbackMessageEvent")
			}
		})
	}
}

func TestGetRelayMessageFromEvent_UnknownEvent(t *testing.T) {
	logger := getTestLogger()
	nid := "stacks"

	cfg := &Config{
		CommonConfig: provider.CommonConfig{
			RPCUrl: "https://stacks-node-api.example.com",
			Contracts: providerTypes.ContractConfigMap{
				"XcallContract":      "ST000000000000000000002AMW42H",
				"ConnectionContract": "ST000000000000000000002AMW42H",
			},
			NID: nid,
		},
	}

	p, err := cfg.NewProvider(context.Background(), logger, "/tmp/relayer", false, nid)
	assert.NoError(t, err, "Expected no error during provider initialization")
	assert.NotNil(t, p, "Expected Provider to be initialized")

	providerConcrete, ok := p.(*Provider)
	assert.True(t, ok, "Expected ChainProvider to be of type *Provider")

	tests := []struct {
		name        string
		eventType   string
		eventData   interface{}
		expectedMsg *providerTypes.Message
		wantErr     bool
	}{
		{
			name:      "Unknown Event Type",
			eventType: "unknown_event",
			eventData: map[string]interface{}{
				"some_field": "some_value",
			},
			expectedMsg: nil,
			wantErr:     true,
		},
		{
			name:      "Empty Event Type",
			eventType: "",
			eventData: map[string]interface{}{
				"another_field": "another_value",
			},
			expectedMsg: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := providerConcrete.getRelayMessageFromEvent(tt.eventType, tt.eventData)
			if tt.wantErr {
				assert.Error(t, err, "Expected an error")
				assert.Nil(t, msg, "Expected Message to be nil")
			} else {
				assert.NoError(t, err, "Expected no error")
				assert.NotNil(t, msg, "Expected a Message object")
			}
		})
	}
}
