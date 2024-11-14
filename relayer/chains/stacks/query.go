package stacks

import (
	"context"
	"fmt"
	"reflect"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	blockchainApiClient "github.com/icon-project/stacks-go-sdk/pkg/stacks_blockchain_api_client"
	"go.uber.org/zap"
)

func (p *Provider) ShouldReceiveMessage(ctx context.Context, message *providerTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, message *providerTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key *providerTypes.MessageKey) (bool, error) {
	switch key.EventType {
	case events.CallMessageSent:
		return p.client.GetReceipt(ctx, p.cfg.Contracts[providerTypes.ConnectionContract], key.Src, key.Sn)
	case events.CallMessage:
		return false, nil
	case events.ResponseMessage:
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
	if response, ok := tx.(*blockchainApiClient.GetTransactionById200Response); ok {
		if mempool := response.GetMempoolTransactionList200ResponseResultsInner; mempool != nil {
			if mempool.ContractCallMempoolTransaction1 != nil {
				txStatus := mempool.ContractCallMempoolTransaction1.TxStatus.String
				if txStatus == nil {
					return nil, fmt.Errorf("nil tx status for contract call mempool transaction")
				}

				return &providerTypes.Receipt{
					TxHash: mempool.ContractCallMempoolTransaction1.TxId,
					Height: 0,
					Status: *txStatus == "success",
				}, nil
			}
			if mempool.SmartContractMempoolTransaction1 != nil {
				txStatus := mempool.SmartContractMempoolTransaction1.TxStatus.String
				if txStatus == nil {
					return nil, fmt.Errorf("nil tx status for smart contract mempool transaction")
				}

				return &providerTypes.Receipt{
					TxHash: mempool.SmartContractMempoolTransaction1.TxId,
					Height: 0,
					Status: *txStatus == "success",
				}, nil
			}
		}

		if confirmed := response.GetTransactionList200ResponseResultsInner; confirmed != nil {
			return getConfirmedReceipt(confirmed)
		}
	}

	return nil, fmt.Errorf("unsupported transaction response type")
}

func getConfirmedReceipt(tx interface{}) (*providerTypes.Receipt, error) {
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
			if txStatusField.IsValid() && txStatusField.Kind() == reflect.Struct {
				stringField := txStatusField.FieldByName("String")
				if stringField.IsValid() {
					if stringField.Kind() == reflect.Ptr {
						if !stringField.IsNil() && stringField.Elem().Kind() == reflect.String {
							status = stringField.Elem().String() == "success"
						} else {
							return nil, fmt.Errorf("string field in txstatus did not parse %s", fieldType)
						}
					} else if stringField.Kind() == reflect.String {
						status = stringField.String() == "success"
					} else {
						return nil, fmt.Errorf("string field in txstatus did not parse %s", fieldType)
					}
				} else {
					return nil, fmt.Errorf("string field in txstatus did not parse %s", fieldType)
				}
			} else {
				return nil, fmt.Errorf("TxStatus field is missing or not a struct in %s", fieldType)
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

func (p *Provider) GenerateMessages(ctx context.Context, fromHeight uint64, toHeight uint64) ([]*providerTypes.Message, error) {
	p.log.Info("Generating messages",
		zap.Uint64("fromHeight", fromHeight),
		zap.Uint64("toHeight", toHeight))

	ctx, cancel := context.WithTimeout(ctx, MAX_WAIT_TIME)
	defer cancel()

	var messages []*providerTypes.Message
	messageChan := make(chan *providerTypes.Message)
	errorChan := make(chan error, 1)

	wsURL := p.client.GetWebSocketURL()
	eventSystem := events.NewEventSystem(ctx, wsURL, p.log, p.client, p.cfg.GetWallet(), p.privateKey, p.cfg.Contracts[providerTypes.XcallContract])

	eventSystem.OnEvent(func(event *events.Event) error {
		if event.BlockHeight < fromHeight || event.BlockHeight > toHeight {
			return nil
		}

		msg, err := p.getRelayMessageFromEvent(event.Type, event.Data)
		if err != nil {
			return fmt.Errorf("failed to parse relay message from event: %w", err)
		}

		msg.MessageHeight = event.BlockHeight

		select {
		case messageChan <- msg:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	go func() {
		defer close(messageChan)

		if err := eventSystem.Start(); err != nil {
			errorChan <- fmt.Errorf("failed to start event system: %w", err)
			return
		}

		<-ctx.Done()
		eventSystem.Stop()

		if ctx.Err() == context.DeadlineExceeded {
			p.log.Info("Event collection completed due to timeout")
		}
	}()

	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				return messages, nil
			}
			messages = append(messages, msg)

		case err := <-errorChan:
			return nil, err

		case <-ctx.Done():
			if len(messages) == 0 {
				return nil, fmt.Errorf("no messages generated within the timeout period")
			}
			return messages, nil
		}
	}
}

func (p *Provider) FetchTxMessages(ctx context.Context, txHash string) ([]*providerTypes.Message, error) {
	p.log.Info("Fetching messages from transaction", zap.String("txHash", txHash))

	receipt, err := p.QueryTransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
	}

	return p.GenerateMessages(ctx, receipt.Height, receipt.Height)
}
