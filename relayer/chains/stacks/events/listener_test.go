package events

import (
	"encoding/json"
	"testing"

	"go.uber.org/zap"
)

func TestEventListener_parseEvent(t *testing.T) {
	log := zap.NewNop()
	listener := &EventListener{
		log: log,
	}

	messageStruct := WSMessage{
		JSONRPC: "2.0",
		Method:  "event",
		Params: json.RawMessage(`{
			"event_type": "smart_contract_log",
			"contract_event": {
				"contract_id": "ST1234.contract",
				"topic": "print",
				"value": {
					"event": "CallMessageSent",
					"data": {
						"from": "fromAddress",
						"to": "toAddress",
						"sn": 123,
						"data": "message data",
						"sources": ["source1"],
						"destinations": ["dest1"]
					}
				}
			},
			"tx_id": "0xabc",
			"block_height": 100
		}`),
	}

	messageBytes, err := json.Marshal(messageStruct)
	if err != nil {
		t.Fatalf("Failed to marshal test message: %v", err)
	}

	event, err := listener.parseEvent(messageBytes)
	if err != nil {
		t.Fatalf("parseEvent returned error: %v", err)
	}
	if event == nil {
		t.Fatal("parseEvent returned nil event")
	}

	if event.Type != CallMessageSent {
		t.Errorf("Expected event type %s, got %s", CallMessageSent, event.Type)
	}

	data, ok := event.Data.(CallMessageSentData)
	if !ok {
		t.Errorf("Event data type is not CallMessageSentData")
	} else {
		if data.From != "fromAddress" {
			t.Errorf("Expected data.From to be 'fromAddress', got '%s'", data.From)
		}
		if data.To != "toAddress" {
			t.Errorf("Expected data.To to be 'toAddress', got '%s'", data.To)
		}
		if data.Sn != 123 {
			t.Errorf("Expected data.Sn to be 123, got %d", data.Sn)
		}
		if data.Data != "message data" {
			t.Errorf("Expected data.Data to be 'message data', got '%s'", data.Data)
		}
		if len(data.Sources) != 1 || data.Sources[0] != "source1" {
			t.Errorf("Expected data.Sources to be ['source1'], got %v", data.Sources)
		}
		if len(data.Destinations) != 1 || data.Destinations[0] != "dest1" {
			t.Errorf("Expected data.Destinations to be ['dest1'], got %v", data.Destinations)
		}
	}
}

func TestEventListener_parseEvent_InvalidMethod(t *testing.T) {
	log := zap.NewNop()
	listener := &EventListener{
		log: log,
	}

	messageStruct := WSMessage{
		JSONRPC: "2.0",
		Method:  "non_event",
		Params:  json.RawMessage(`{}`),
	}

	messageBytes, err := json.Marshal(messageStruct)
	if err != nil {
		t.Fatalf("Failed to marshal test message: %v", err)
	}

	event, err := listener.parseEvent(messageBytes)
	if err != nil {
		t.Errorf("parseEvent returned error: %v", err)
	}
	if event != nil {
		t.Errorf("Expected nil event for non 'event' method")
	}
}

func TestEventListener_parseEvent_InvalidEventType(t *testing.T) {
	log := zap.NewNop()
	listener := &EventListener{
		log: log,
	}

	messageStruct := WSMessage{
		JSONRPC: "2.0",
		Method:  "event",
		Params: json.RawMessage(`{
			"event_type": "other_event",
			"contract_event": {}
		}`),
	}

	messageBytes, err := json.Marshal(messageStruct)
	if err != nil {
		t.Fatalf("Failed to marshal test message: %v", err)
	}

	event, err := listener.parseEvent(messageBytes)
	if err != nil {
		t.Errorf("parseEvent returned error: %v", err)
	}
	if event != nil {
		t.Errorf("Expected nil event for invalid event_type")
	}
}
