package events

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"container/ring"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/interfaces"
	blockchainApiClient "github.com/icon-project/stacks-go-sdk/pkg/stacks_blockchain_api_client"
	"github.com/icon-project/stacks-go-sdk/pkg/websocket"
	"go.uber.org/zap"
)

type EventListener struct {
	wsURL           string
	wsClient        *websocket.Client
	eventChan       chan *Event
	processChan     chan *Event
	backlog         *ring.Ring
	maxBufferSize   int
	log             *zap.Logger
	ctx             context.Context
	cancel          context.CancelFunc
	contractAddress string
	client          interfaces.IClient
}

func NewEventListener(ctx context.Context, wsURL string, bufferSize int, log *zap.Logger, contractAddress string, client interfaces.IClient) *EventListener {
	ctx, cancel := context.WithCancel(ctx)
	return &EventListener{
		wsURL:           wsURL,
		eventChan:       make(chan *Event, bufferSize),
		processChan:     make(chan *Event, bufferSize),
		backlog:         ring.New(bufferSize),
		maxBufferSize:   bufferSize,
		log:             log,
		ctx:             ctx,
		cancel:          cancel,
		contractAddress: contractAddress,
		client:          client,
	}
}

func (l *EventListener) Start() error {
	var err error
	l.wsClient, err = websocket.NewClient(l.wsURL)
	if err != nil {
		return fmt.Errorf("failed to create websocket client: %w", err)
	}

	addressTxChan, err := l.wsClient.SubscribeAddressTransactions(l.ctx, l.contractAddress)
	if err != nil {
		return fmt.Errorf("failed to subscribe to address transactions: %w", err)
	}

	l.log.Info("Subscribed",
		zap.String("address", l.contractAddress),
	)

	go l.handleAddressTransactions(addressTxChan)
	go l.bufferEvents()

	l.log.Info("EventListener started successfully",
		zap.String("wsURL", l.wsURL),
		zap.String("contractAddress", l.contractAddress))

	return nil
}

func (l *EventListener) Stop() {
	l.cancel()
	if l.wsClient != nil {
		l.wsClient.Close()
	}
	close(l.eventChan)
	close(l.processChan)
}

func (l *EventListener) handleAddressTransactions(txChan <-chan websocket.AddressTxUpdateEvent) {
	for {
		select {
		case <-l.ctx.Done():
			return
		case tx := <-txChan:
			l.log.Debug("Received address transaction event",
				zap.String("txID", tx.Params.TxID),
				zap.String("txType", tx.Params.TxType),
				zap.String("status", tx.Params.TxStatus),
				zap.String("contractID", tx.Params.Tx.ContractCall.ContractID))

			if tx.Params.TxType != "contract_call" {
				l.log.Debug("Ignoring non-contract-call transaction")
				continue
			}

			if tx.Params.TxStatus != "success" {
				l.log.Debug("Ignoring unsuccessful transaction",
					zap.String("txID", tx.Params.TxID),
					zap.String("status", tx.Params.TxStatus))
				continue
			}

			if err := l.processContractCall(&tx.Params); err != nil {
				l.log.Error("Failed to process contract call",
					zap.Error(err),
					zap.String("txID", tx.Params.TxID))
				continue
			}
		}
	}
}

func (l *EventListener) processContractCall(tx *websocket.AddressTxUpdate) error {
	time.Sleep(5 * time.Second) // occasionally the tx is a mempool tx

	l.log.Debug("Processing transaction",
		zap.String("txID", tx.TxID),
		zap.String("txType", tx.TxType),
		zap.Any("contractCall", tx.Tx.ContractCall))

	eventCount := 0
	if tx.Tx.Events != nil {
		eventCount = len(tx.Tx.Events)
	}
	l.log.Debug("Transaction event details",
		zap.Int("eventCount", eventCount),
		zap.Int("txEventCount", tx.Tx.EventCount),
		zap.Any("tx.Tx.Events", tx.Tx.Events))

	if len(tx.Tx.Events) == 0 {
		l.log.Debug("Events array is empty, fetching full transaction details")

		fullTx, err := l.client.GetTransactionById(l.ctx, tx.TxID)
		if err != nil {
			l.log.Error("Failed to fetch transaction by ID", zap.Error(err), zap.String("txID", tx.TxID))
			return err
		}

		var contractCallTx *blockchainApiClient.ContractCallTransaction

		if fullTx.GetTransactionList200ResponseResultsInner != nil && fullTx.GetTransactionList200ResponseResultsInner.ContractCallTransaction != nil {
			contractCallTx = fullTx.GetTransactionList200ResponseResultsInner.ContractCallTransaction
		}

		if contractCallTx == nil {
			l.log.Debug("Transaction is not a ContractCallTransaction or events are unavailable", zap.String("txID", tx.TxID))
			return nil
		}

		txEvents := contractCallTx.Events
		if len(txEvents) == 0 {
			l.log.Debug("No events found in full transaction details", zap.String("txID", tx.TxID))
			return nil
		}

		for i, event := range txEvents {
			l.log.Debug("Processing event from ContractCallTransaction",
				zap.Int("eventIndex", i),
				zap.Any("event", event))

			if event.SmartContractLogTransactionEvent != nil {
				contractLog := event.SmartContractLogTransactionEvent.ContractLog
				if contractLog.Topic == "print" {
					repr := contractLog.Value.Repr
					l.log.Debug("Found print event log",
						zap.String("repr", repr),
						zap.String("topic", contractLog.Topic))

					if strings.Contains(repr, "CallMessageSent") {
						l.log.Debug("Found CallMessageSent event in print log",
							zap.String("full_event", repr))

						eventData := parseClarityTuple(repr)

						snStr := eventData["sn"].(string)
						sn, err := strconv.ParseUint(strings.TrimPrefix(snStr, "u"), 10, 64)
						if err != nil {
							return fmt.Errorf("failed to parse sn: %w", err)
						}

						sentData := &CallMessageSentData{
							From:         eventData["from"].(string),
							To:           eventData["to"].(string),
							Sn:           sn,
							Data:         eventData["data"].(string),
							Sources:      eventData["sources"].([]string),
							Destinations: eventData["destinations"].([]string),
						}

						event := &Event{
							ID:          fmt.Sprintf("%s-%s-%d", "CallMessageSent", tx.TxID, time.Now().UnixNano()),
							Type:        "CallMessageSent",
							Data:        sentData,
							BlockHeight: uint64(tx.Tx.BlockHeight),
							Timestamp:   time.Now(),
							Raw:         []byte(repr),
						}

						l.eventChan <- event
						l.log.Debug("Processed and sent event",
							zap.String("type", event.Type),
							zap.String("id", event.ID))
					} else if strings.Contains(repr, "Message") {
						l.log.Debug("Found Message event in print log",
							zap.String("full_event", repr))

						eventData := parseClarityTuple(repr)
						if eventData == nil {
							return fmt.Errorf("failed to parse Message: %w", err)
						}

						snStr := eventData["sn"].(string)
						sn, err := strconv.ParseInt(snStr, 10, 64)
						if err != nil {
							return fmt.Errorf("failed to parse sn: %w", err)
						}

						messageData := &MessageData{
							From: "stacks_testnet",
							To:   eventData["to"].(string),
							Sn:   sn,
							Data: eventData["msg"].(string),
						}
						event := &Event{
							ID:          fmt.Sprintf("%s-%s-%d", "Message", tx.TxID, time.Now().UnixNano()),
							Type:        "Message",
							Data:        messageData,
							BlockHeight: uint64(tx.Tx.BlockHeight),
							Timestamp:   time.Now(),
							Raw:         []byte(repr),
						}
						l.eventChan <- event
					}
				}
			}
		}
	}
	return nil
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

func parseClarityTuple(repr string) map[string]interface{} {
	data := make(map[string]interface{})

	tupleContent := strings.TrimPrefix(strings.TrimSuffix(repr, ")"), "(tuple ")

	var fields []string
	var currentField strings.Builder
	parenthesesCount := 0

	for _, char := range tupleContent {
		if char == '(' {
			parenthesesCount++
		} else if char == ')' {
			parenthesesCount--
		}

		if char == ' ' && parenthesesCount == 0 && currentField.Len() > 0 {
			fields = append(fields, currentField.String())
			currentField.Reset()
		} else {
			currentField.WriteRune(char)
		}
	}
	if currentField.Len() > 0 {
		fields = append(fields, currentField.String())
	}

	for _, field := range fields {
		if len(field) == 0 {
			continue
		}

		parts := strings.SplitN(field, " ", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.Trim(parts[0], "()")
		value := strings.Trim(parts[1], "()")

		if strings.HasPrefix(value, "list") {
			listStr := strings.TrimPrefix(value, "list ")
			listItems := strings.Split(strings.Trim(listStr, "\""), "\" \"")
			data[key] = listItems
			continue
		}

		if strings.HasPrefix(value, "u") {
			data[key] = strings.TrimPrefix(value, "u")
		} else if strings.HasPrefix(value, "0x") {
			data[key] = value
		} else {
			data[key] = strings.Trim(value, "\"'")
		}
	}

	return data
}
