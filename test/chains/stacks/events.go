package stacks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/icon-project/centralized-relay/test/chains"
	blockchainApiClient "github.com/icon-project/stacks-go-sdk/pkg/stacks_blockchain_api_client"
)

func (s *StacksLocalnet) FindEvent(ctx context.Context, startHeight uint64, contract, signature string, index []string) (*blockchainApiClient.SmartContractLogTransactionEvent, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while finding event %s", signature)
		default:
			events, err := s.client.GetContractEvents(ctx, contract, 50, 0)
			if err != nil {
				return nil, fmt.Errorf("failed to get contract events: %w", err)
			}

			for _, event := range events.Results {
				if event.SmartContractLogTransactionEvent != nil &&
					event.SmartContractLogTransactionEvent.ContractLog.Topic == signature {
					return event.SmartContractLogTransactionEvent, nil
				}
			}

			time.Sleep(BLOCK_TIME)
		}
	}
}

func (s *StacksLocalnet) FindCallMessage(ctx context.Context, startHeight uint64, from, to, sn string) (string, string, error) {
	for {
		select {
		case <-ctx.Done():
			return "", "", fmt.Errorf("context cancelled while finding call message with sn %s", sn)
		default:
			events, err := s.client.GetContractEvents(ctx, s.IBCAddresses["xcall-proxy"], 50, 0)
			if err != nil {
				return "", "", fmt.Errorf("failed to get contract events: %w", err)
			}

			for _, event := range events.Results {
				if event.SmartContractLogTransactionEvent != nil {
					log := event.SmartContractLogTransactionEvent.ContractLog
					if log.Topic == "print" && strings.Contains(log.Value.Repr, "CallMessage") {
						eventSn := extractSnFromEvent(log.Value.Repr)
						if eventSn == sn {
							reqId, data := extractCallMessageData(log.Value.Repr)
							if reqId != "" && data != "" {
								return reqId, data, nil
							}
						}
					}
				}
			}

			time.Sleep(BLOCK_TIME)
		}
	}
}

func (s *StacksLocalnet) FindCallResponse(ctx context.Context, startHeight uint64, sn string) (string, error) {
	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("context cancelled while finding call response with sn %s", sn)
		default:
			events, err := s.client.GetContractEvents(ctx, s.IBCAddresses["xcall-proxy"], 50, 0)
			if err != nil {
				return "", fmt.Errorf("failed to get contract events: %w", err)
			}

			for _, event := range events.Results {
				if event.SmartContractLogTransactionEvent != nil {
					log := event.SmartContractLogTransactionEvent.ContractLog
					if log.Topic == "print" && strings.Contains(log.Value.Repr, "CallResponse") {
						eventSn := extractSnFromEvent(log.Value.Repr)
						if eventSn == sn {
							return event.SmartContractLogTransactionEvent.TxId, nil
						}
					}
				}
			}

			time.Sleep(BLOCK_TIME)
		}
	}
}

func (s *StacksLocalnet) FindRollbackExecutedMessage(ctx context.Context, startHeight uint64, sn string) (string, error) {
	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("context cancelled while finding rollback message with sn %s", sn)
		default:
			events, err := s.client.GetContractEvents(ctx, s.IBCAddresses["xcall-proxy"], 50, 0)
			if err != nil {
				return "", fmt.Errorf("failed to get contract events: %w", err)
			}

			for _, event := range events.Results {
				if event.SmartContractLogTransactionEvent != nil {
					log := event.SmartContractLogTransactionEvent.ContractLog
					if log.Topic == "print" && strings.Contains(log.Value.Repr, "RollbackExecuted") {
						eventSn := extractSnFromEvent(log.Value.Repr)
						if eventSn == sn {
							return sn, nil
						}
					}
				}
			}

			time.Sleep(BLOCK_TIME)
		}
	}
}

func (s *StacksLocalnet) FindTargetXCallMessage(ctx context.Context, target chains.Chain, height uint64, to string) (*chains.XCallResponse, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while finding target xcall message")
		default:
			events, err := s.client.GetContractEvents(ctx, s.IBCAddresses["xcall-proxy"], 50, 0)
			if err != nil {
				return nil, fmt.Errorf("failed to get contract events: %w", err)
			}

			for _, event := range events.Results {
				if event.SmartContractLogTransactionEvent != nil {
					log := event.SmartContractLogTransactionEvent.ContractLog
					if log.Topic == "print" {
						if strings.Contains(log.Value.Repr, "EmitMessage") {
							sn, msg, targetNetwork := extractEmitMessageData(log.Value.Repr)
							if targetNetwork == to {
								return &chains.XCallResponse{
									SerialNo: sn,
									Data:     msg,
								}, nil
							}
						} else if strings.Contains(log.Value.Repr, "CallMessage") {
							sn, reqId, data := extractFullCallMessageData(log.Value.Repr)
							return &chains.XCallResponse{
								SerialNo:  sn,
								RequestID: reqId,
								Data:      data,
							}, nil
						}
					}
				}
			}

			time.Sleep(BLOCK_TIME)
		}
	}
}

func extractSnFromEvent(repr string) string {
	startIdx := strings.Index(repr, "(sn u")
	if startIdx != -1 {
		startIdx += 5 // Move past "(sn u"
		endIdx := strings.Index(repr[startIdx:], ")")
		if endIdx != -1 {
			return repr[startIdx : startIdx+endIdx]
		}
	}
	return ""
}

func extractCallMessageData(repr string) (reqId, data string) {
	reqIdStart := strings.Index(repr, "reqId u")
	if reqIdStart != -1 {
		reqIdStart += 7
		reqIdEnd := strings.Index(repr[reqIdStart:], " ")
		if reqIdEnd != -1 {
			reqId = repr[reqIdStart : reqIdStart+reqIdEnd]
		}
	}

	dataStart := strings.Index(repr, "data 0x")
	if dataStart != -1 {
		dataStart += 7
		dataEnd := strings.Index(repr[dataStart:], ")")
		if dataEnd != -1 {
			data = repr[dataStart : dataStart+dataEnd]
		}
	}

	return reqId, data
}

func extractEmitMessageData(repr string) (sn, msg, targetNetwork string) {
	sn = extractSnFromEvent(repr)

	msgStart := strings.Index(repr, "msg 0x")
	if msgStart != -1 {
		msgStart += 6
		msgEnd := strings.Index(repr[msgStart:], " ")
		if msgEnd != -1 {
			msg = repr[msgStart : msgStart+msgEnd]
		}
	}

	networkStart := strings.Index(repr, "network \"")
	if networkStart != -1 {
		networkStart += 9
		networkEnd := strings.Index(repr[networkStart:], "\"")
		if networkEnd != -1 {
			targetNetwork = repr[networkStart : networkStart+networkEnd]
		}
	}

	return sn, msg, targetNetwork
}

func extractFullCallMessageData(repr string) (sn, reqId, data string) {
	sn = extractSnFromEvent(repr)
	reqId, data = extractCallMessageData(repr)
	return sn, reqId, data
}
