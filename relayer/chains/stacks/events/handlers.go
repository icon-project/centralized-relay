package events

import (
	"fmt"
	"math/big"
	"strings"

	"go.uber.org/zap"

	"github.com/icon-project/stacks-go-sdk/pkg/clarity"
)

func (p *EventProcessor) handleCallMessageSentEvent(event *Event) error {
	data, ok := event.Data.(*CallMessageSentData)
	if !ok {
		return fmt.Errorf("invalid event data type for CallMessageSent")
	}

	p.log.Info("Processing CallMessageSent event",
		zap.String("from", data.From),
		zap.String("to", data.To),
		zap.Uint64("sn", data.Sn),
		zap.Strings("sources", data.Sources),
		zap.Strings("destinations", data.Destinations),
	)

	for _, source := range data.Sources {
		if err := p.callSendMessageFunction(source, data.To, data.Sn, data.Data); err != nil {
			p.log.Error("Failed to call send-message", zap.Error(err), zap.String("source", source))
			return err
		}
	}

	return nil
}

func (p *EventProcessor) handleMessageEvent(event *Event) error {
	data, ok := event.Data.(*MessageData)
	if !ok {
		return fmt.Errorf("invalid event data type for Message")
	}

	p.log.Info("Processing Message event",
		zap.String("from", data.From),
		zap.String("to", data.To),
		zap.Int64("sn", data.Sn),
	)

	return nil
}

func (p *EventProcessor) handleCallMessageEvent(event *Event) error {
	data, ok := event.Data.(CallMessageData)
	if !ok {
		return fmt.Errorf("invalid event data type for CallMessage")
	}

	p.log.Info("Processing CallMessage event",
		zap.String("from", data.From),
		zap.String("to", data.To),
		zap.Uint64("sn", data.Sn),
		zap.Uint64("req-id", data.ReqID),
	)

	return nil
}

func (p *EventProcessor) handleResponseMessageEvent(event *Event) error {
	data, ok := event.Data.(ResponseMessageData)
	if !ok {
		return fmt.Errorf("invalid event data type for ResponseMessage")
	}

	p.log.Info("Processing ResponseMessage event",
		zap.Uint64("sn", data.Sn),
		zap.Uint64("code", data.Code),
		zap.String("msg", data.Msg),
	)

	return nil
}

func (p *EventProcessor) handleRollbackMessageEvent(event *Event) error {
	data, ok := event.Data.(RollbackMessageData)
	if !ok {
		return fmt.Errorf("invalid event data type for RollbackMessage")
	}

	p.log.Info("Processing RollbackMessage event",
		zap.Uint64("sn", data.Sn),
	)

	return nil
}

func (p *EventProcessor) callSendMessageFunction(sourceContract string, to string, sn uint64, msg string) error {
	parts := strings.Split(sourceContract, ".")
	if len(parts) != 2 {
		return fmt.Errorf("invalid source contract format: %s", sourceContract)
	}
	sourceContractAddress := parts[0]
	sourceContractName := parts[1]

	p.log.Info("Calling send-message on",
		zap.String("sourceContract", sourceContract),
	)

	cvTo, err := clarity.NewStringASCII(to)
	if err != nil {
		return fmt.Errorf("failed to create to address clarity value: %w", err)
	}

	cvSn, err := clarity.NewInt(big.NewInt(int64(sn)))
	if err != nil {
		return fmt.Errorf("failed to create sn clarity value: %w", err)
	}

	cvBuffer := clarity.NewBuffer([]byte(msg))
	cvSvc, err := clarity.NewStringASCII("")
	if err != nil {
		return fmt.Errorf("failed to create svc clarity value: %w", err)
	}

	args := []clarity.ClarityValue{
		cvTo,
		cvSvc,
		cvSn,
		cvBuffer,
	}

	tx, err := p.client.MakeContractCall(
		p.ctx,
		sourceContractAddress,
		sourceContractName,
		"send-message",
		args,
		p.senderAddress,
		p.senderKey,
	)
	if err != nil {
		return fmt.Errorf("failed to create contract call transaction: %w", err)
	}

	txID, err := p.client.BroadcastTransaction(p.ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	p.log.Info("send-message transaction sent", zap.String("txID", txID), zap.String("source", sourceContract))

	return nil
}
