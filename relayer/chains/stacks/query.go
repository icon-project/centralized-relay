package stacks

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

const eventListeningTimeout = 30 * time.Second

func (p *Provider) ShouldReceiveMessage(ctx context.Context, message *providerTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, message *providerTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key *providerTypes.MessageKey) (bool, error) {
	switch key.EventType {
	case events.EmitMessage:
		return p.client.GetReceipt(ctx, p.cfg.Contracts[providerTypes.ConnectionContract], key.Src, key.Sn)
	case events.CallMessage:
		return false, nil
	case events.RollbackMessage:
		return false, nil
	default:
		return true, fmt.Errorf("unknown event type")
	}
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, id string) (*providerTypes.Receipt, error) {
	res, err := p.client.GetTransactionById(ctx, id)
	if err != nil {
		p.log.Error("Failed to query transaction receipt", zap.String("txHash", id), zap.Error(err))
		return nil, err
	}

	if mempoolResp := res.GetMempoolTransactionList200ResponseResultsInner; mempoolResp != nil {
		receipt, err := GetReceipt(mempoolResp)
		if err != nil {
			return nil, fmt.Errorf("failed to extract mempool transaction: %w", err)
		}
		return receipt, nil
	}

	if confirmedResp := res.GetTransactionList200ResponseResultsInner; confirmedResp != nil {
		receipt, err := GetReceipt(confirmedResp)
		if err != nil {
			return nil, fmt.Errorf("failed to extract confirmed transaction: %w", err)
		}
		return receipt, nil
	}

	return nil, fmt.Errorf("failed to query transaction: %w", err)
}

func (p *Provider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	latestBlock, err := p.client.GetLatestBlock(ctx)
	if err != nil {
		return 0, err
	}
	if latestBlock == nil {
		return 0, fmt.Errorf("no blocks found")
	}
	return uint64(latestBlock.Height), nil
}

func GetReceipt(tx interface{}) (*providerTypes.Receipt, error) {
	val := reflect.ValueOf(tx).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i).Name

		if !field.IsNil() {
			txStruct := field.Elem()

			txTypeField := txStruct.FieldByName("TxType")
			if !txTypeField.IsValid() || txTypeField.Kind() != reflect.String {
				return nil, fmt.Errorf("TxType field is missing or not a string in %s", fieldType)
			}

			txIdField := txStruct.FieldByName("TxId")
			if !txIdField.IsValid() || txIdField.Kind() != reflect.String {
				return nil, fmt.Errorf("TxId field is missing or not a string in %s", fieldType)
			}
			txId := txIdField.String()

			var blockHeight uint64 = 0
			blockHeightField := txStruct.FieldByName("BlockHeight")
			if blockHeightField.IsValid() && blockHeightField.Kind() == reflect.Int32 {
				blockHeight = uint64(blockHeightField.Int())
			}

			var status bool = false
			txStatusField := txStruct.FieldByName("TxStatus")
			if txStatusField.IsValid() && txStatusField.Kind() == reflect.String {
				status = txStatusField.String() == "success"
			}

			height := blockHeight
			if txTypeField.Kind() == reflect.String && blockHeight == 0 {
				height = 0
			}

			return &providerTypes.Receipt{
				TxHash: txId,
				Height: height,
				Status: status,
			}, nil
		}
	}

	return nil, fmt.Errorf("no non-nil transaction field found")
}

func (p *Provider) GenerateMessages(ctx context.Context, key *providerTypes.MessageKeyWithMessageHeight) ([]*providerTypes.Message, error) {
	p.log.Info("Generating messages", zap.Any("messagekey", key))
	if key == nil {
		return nil, fmt.Errorf("GenerateMessages: message key cannot be nil")
	}

	eventTypes := p.getSubscribedEventTypes()

	var messages []*providerTypes.Message
	errChan := make(chan error, 1)

	callback := func(eventType string, data interface{}) error {
		msg, err := p.getRelayMessageFromEvent(eventType, data)
		if err != nil {
			p.log.Error("Failed to parse relay message from event", zap.Error(err))
			return err
		}

		msg.MessageHeight = key.Height

		messages = append(messages, msg)
		return nil
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, eventListeningTimeout)
	defer cancel()

	err := p.client.SubscribeToEvents(ctxWithTimeout, eventTypes, callback)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to events: %w", err)
	}

	select {
	case err := <-errChan:
		return nil, fmt.Errorf("error occurred while processing events: %w", err)
	case <-ctxWithTimeout.Done():
		if len(messages) == 0 {
			return nil, fmt.Errorf("no messages generated within the timeout period")
		}
	}

	return messages, nil
}
