package stacks

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

type MessageData struct {
	From string `json:"from"`
	To   string `json:"to"`
	Sn   int64  `json:"sn"`
	Data string `json:"data"`
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

func extractNetworkID(fullAddress string) string {
	parts := strings.Split(fullAddress, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return fullAddress
}

func (p *Provider) parseBaseMessage(data interface{}) (*providerTypes.Message, error) {
	var msgData MessageData

	dataBytes, err := json.Marshal(data)
	if err != nil {
		p.log.Error("Failed to marshal event data", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal event data: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &msgData); err != nil {
		p.log.Error("Failed to unmarshal to MessageData", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal to MessageData: %w", err)
	}

	sn := new(big.Int).SetInt64(msgData.Sn)
	dstNetworkID := extractNetworkID(msgData.To)

	return &providerTypes.Message{
		Dst:           dstNetworkID,
		Src:           p.NID(),
		Sn:            sn,
		Data:          []byte(msgData.Data),
		MessageHeight: 0, // Set by Listener after parsing
	}, nil
}

func (p *Provider) parseEmitMessageEvent(data interface{}) (*providerTypes.Message, error) {
	msg, err := p.parseBaseMessage(data)
	if err != nil {
		return nil, err
	}
	msg.EventType = events.EmitMessage
	return msg, nil
}

func (p *Provider) parseCallMessageEvent(data interface{}) (*providerTypes.Message, error) {
	msg, err := p.parseBaseMessage(data)
	if err != nil {
		return nil, err
	}

	msg.Dst = p.NID()
	msg.Src = p.NID()
	msg.EventType = events.CallMessage

	return msg, nil
}

func (p *Provider) parseRollbackMessageEvent(data interface{}) (*providerTypes.Message, error) {
	msg, err := p.parseBaseMessage(data)
	if err != nil {
		return nil, err
	}

	msg.Dst = p.NID()
	msg.Src = p.NID()
	msg.Data = nil // RollbackMessage doesn't need data
	msg.EventType = events.RollbackMessage

	return msg, nil
}
