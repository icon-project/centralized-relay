package stacks

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

type EmitMessageEvent struct {
	TargetNetwork string `json:"targetNetwork"`
	Sn            string `json:"sn"`
	Msg           string `json:"msg"`
}

type CallMessageEvent struct {
	ReqID string `json:"req_id"`
	Sn    string `json:"sn"`
	Data  string `json:"data"`
}

type RollbackMessageEvent struct {
	Sn string `json:"sn"`
}

func (p *Provider) getRelayMessageFromEvent(eventType string, data interface{}) (*providerTypes.Message, error) {
	switch eventType {
	case EmitMessage:
		return p.parseEmitMessageEvent(data)
	case CallMessage:
		return p.parseCallMessageEvent(data)
	case RollbackMessage:
		return p.parseRollbackMessageEvent(data)
	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}
}

func (p *Provider) parseEmitMessageEvent(data interface{}) (*providerTypes.Message, error) {
	var emitMsg EmitMessageEvent

	dataBytes, err := json.Marshal(data)
	if err != nil {
		p.log.Error("Failed to marshal EmitMessageEvent data", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal EmitMessageEvent data: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &emitMsg); err != nil {
		p.log.Error("Failed to unmarshal EmitMessageEvent", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal EmitMessageEvent: %w", err)
	}

	sn, ok := new(big.Int).SetString(emitMsg.Sn, 10)
	if !ok {
		p.log.Error("Invalid Sn in EmitMessageEvent", zap.String("Sn", emitMsg.Sn))
		return nil, fmt.Errorf("invalid Sn in EmitMessageEvent: %s", emitMsg.Sn)
	}

	msg := &providerTypes.Message{
		Dst:           emitMsg.TargetNetwork,
		Src:           p.NID(),
		Sn:            sn,
		MessageHeight: 0,
		EventType:     events.EmitMessage,
		Data:          []byte(emitMsg.Msg),
	}

	return msg, nil
}

func (p *Provider) parseCallMessageEvent(data interface{}) (*providerTypes.Message, error) {
	var callMsg CallMessageEvent

	dataBytes, err := json.Marshal(data)
	if err != nil {
		p.log.Error("Failed to marshal CallMessageEvent data", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal CallMessageEvent data: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &callMsg); err != nil {
		p.log.Error("Failed to unmarshal CallMessageEvent", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal CallMessageEvent: %w", err)
	}

	sn, ok := new(big.Int).SetString(callMsg.Sn, 10)
	if !ok {
		p.log.Error("Invalid Sn in CallMessageEvent", zap.String("Sn", callMsg.Sn))
		return nil, fmt.Errorf("invalid Sn in CallMessageEvent: %s", callMsg.Sn)
	}

	reqID, ok := new(big.Int).SetString(callMsg.ReqID, 10)
	if !ok {
		p.log.Error("Invalid ReqID in CallMessageEvent", zap.String("ReqID", callMsg.ReqID))
		return nil, fmt.Errorf("invalid ReqID in CallMessageEvent: %s", callMsg.ReqID)
	}

	msg := &providerTypes.Message{
		Dst:           p.NID(),
		Src:           p.NID(),
		Sn:            sn,
		MessageHeight: 0,
		EventType:     events.CallMessage,
		Data:          []byte(callMsg.Data),
		ReqID:         reqID,
	}

	return msg, nil
}

func (p *Provider) parseRollbackMessageEvent(data interface{}) (*providerTypes.Message, error) {
	var rollbackMsg RollbackMessageEvent

	dataBytes, err := json.Marshal(data)
	if err != nil {
		p.log.Error("Failed to marshal RollbackMessageEvent data", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal RollbackMessageEvent data: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &rollbackMsg); err != nil {
		p.log.Error("Failed to unmarshal RollbackMessageEvent", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal RollbackMessageEvent: %w", err)
	}

	sn, ok := new(big.Int).SetString(rollbackMsg.Sn, 10)
	if !ok {
		p.log.Error("Invalid Sn in RollbackMessageEvent", zap.String("Sn", rollbackMsg.Sn))
		return nil, fmt.Errorf("invalid Sn in RollbackMessageEvent: %s", rollbackMsg.Sn)
	}

	msg := &providerTypes.Message{
		Dst:           p.NID(),
		Src:           p.NID(),
		Sn:            sn,
		MessageHeight: 0,
		EventType:     events.RollbackMessage,
	}

	return msg, nil
}
