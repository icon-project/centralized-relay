package events

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"container/ring"

	"github.com/cenkalti/backoff/v4"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type EventListener struct {
	wsURL         string
	conn          *websocket.Conn
	eventChan     chan *Event
	processChan   chan *Event
	backlog       *ring.Ring
	maxBufferSize int
	mu            sync.RWMutex
	log           *zap.Logger
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewEventListener(ctx context.Context, wsURL string, bufferSize int, log *zap.Logger) *EventListener {
	ctx, cancel := context.WithCancel(ctx)
	return &EventListener{
		wsURL:         wsURL,
		eventChan:     make(chan *Event, bufferSize),
		processChan:   make(chan *Event, bufferSize),
		backlog:       ring.New(bufferSize),
		maxBufferSize: bufferSize,
		log:           log,
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (l *EventListener) Start() error {
	go l.maintainConnection()
	go l.bufferEvents()
	return nil
}

func (l *EventListener) Stop() {
	l.cancel()
	if l.conn != nil {
		l.conn.Close()
	}
	close(l.eventChan)
	close(l.processChan)
}

func (l *EventListener) maintainConnection() {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 0 // Retry forever

	for {
		select {
		case <-l.ctx.Done():
			return
		default:
			if err := l.connect(); err != nil {
				l.log.Error("WebSocket connection failed", zap.Error(err))
				time.Sleep(b.NextBackOff())
				continue
			}
			b.Reset()
			l.readMessages()
		}
	}
}

func (l *EventListener) connect() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.conn != nil {
		l.conn.Close()
	}

	conn, _, err := websocket.DefaultDialer.Dial(l.wsURL, nil)
	if err != nil {
		return err
	}

	l.conn = conn

	// Subscribe to all event types
	for _, eventType := range []string{EmitMessage, CallMessage, RollbackMessage} {
		if err := l.subscribe(eventType); err != nil {
			l.conn.Close()
			return fmt.Errorf("failed to subscribe to %s: %w", eventType, err)
		}
	}

	return nil
}

func (l *EventListener) readMessages() {
	for {
		select {
		case <-l.ctx.Done():
			return
		default:
			_, message, err := l.conn.ReadMessage()
			if err != nil {
				l.log.Error("Failed to read WebSocket message", zap.Error(err))
				return
			}

			event, err := l.parseEvent(message)
			if err != nil {
				l.log.Error("Failed to parse event", zap.Error(err))
				continue
			}

			if event == nil {
				continue
			}

			l.eventChan <- event
		}
	}
}

func (l *EventListener) bufferEvents() {
	for {
		select {
		case event := <-l.eventChan:
			l.backlog.Value = event
			l.backlog = l.backlog.Next()
			l.processChan <- event
		case <-l.ctx.Done():
			return
		}
	}
}

func (l *EventListener) subscribe(eventType string) error {
	request := WSRequest{
		JSONRPC: "2.0",
		ID:      time.Now().UnixNano(),
		Method:  "subscribe",
		Params: map[string]interface{}{
			"event": eventType,
		},
	}

	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription request: %w", err)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to send subscription request: %w", err)
	}

	_, message, err := l.conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("failed to read subscription response: %w", err)
	}

	var response WSResponse
	if err := json.Unmarshal(message, &response); err != nil {
		return fmt.Errorf("failed to unmarshal subscription response: %w", err)
	}

	if response.Error != nil {
		return fmt.Errorf("subscription failed: %s", response.Error.Message)
	}

	l.log.Info("Successfully subscribed to event", zap.String("type", eventType))
	return nil
}

func (l *EventListener) parseEvent(message []byte) (*Event, error) {
	var wsMsg WSMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal WebSocket message: %w", err)
	}

	if wsMsg.Method != "event" {
		return nil, nil
	}

	var smartContractLog SmartContractLogEvent
	if err := json.Unmarshal(wsMsg.Params, &smartContractLog); err != nil {
		return nil, fmt.Errorf("failed to unmarshal smart contract log: %w", err)
	}

	if smartContractLog.EventType != "smart_contract_log" ||
		smartContractLog.ContractEvent.Topic != "print" {
		return nil, nil
	}

	var printValue struct {
		Event string          `json:"event"`
		Data  json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(smartContractLog.ContractEvent.Value, &printValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal print value: %w", err)
	}

	var eventData interface{}
	switch printValue.Event {
	case EmitMessage:
		var data EmitMessageData
		if err := json.Unmarshal(printValue.Data, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal emit message data: %w", err)
		}
		eventData = data

	case CallMessage:
		var data CallMessageData
		if err := json.Unmarshal(printValue.Data, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal call message data: %w", err)
		}
		eventData = data

	case RollbackMessage:
		var data RollbackMessageData
		if err := json.Unmarshal(printValue.Data, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal rollback message data: %w", err)
		}
		eventData = data

	default:
		l.log.Debug("Ignoring unknown event type", zap.String("type", printValue.Event))
		return nil, nil
	}

	event := &Event{
		ID:          fmt.Sprintf("%s-%d", printValue.Event, time.Now().UnixNano()),
		Type:        printValue.Event,
		Data:        eventData,
		BlockHeight: smartContractLog.BlockHeight,
		Timestamp:   time.Now(),
		Raw:         message,
	}

	return event, nil
}
