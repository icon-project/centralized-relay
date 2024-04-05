package icon

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func (p *Provider) parseMessagesFromEventlogs(log *zap.Logger, eventlogs []*types.EventLog, height uint64) []*providerTypes.Message {
	msgs := make([]*providerTypes.Message, 0)
	for _, el := range eventlogs {
		message, ok := p.parseMessageFromEvent(log, el, height)
		if ok {
			msgs = append(msgs, message)
		}
	}
	return msgs
}

func (p *Provider) parseMessageFromEvent(log *zap.Logger, event *types.EventLog, height uint64) (*providerTypes.Message, bool) {
	eventName := string(event.Indexed[0][:])
	eventType := p.GetEventName(eventName)
	switch eventName {
	case EmitMessage:
		m, err := p.parseEmitMessage(event, eventType, height)
		if err != nil {
			log.Error("invalid event", zap.Error(err))
			return nil, false
		}
		return m, true
	case CallMessage:
		m, err := p.parseCallMessage(event, eventType, height)
		if err != nil {
			log.Error("invalid event", zap.Error(err))
			return nil, false
		}
		return m, true
	default:
		log.Error("unknown event", zap.String("event", eventName))
		return nil, false
	}
}

func (p *Provider) parseEmitMessage(e *types.EventLog, eventType string, height uint64) (*providerTypes.Message, error) {
	if indexdedLen, dataLen := len(e.Indexed), len(e.Data); indexdedLen != 3 && dataLen != 1 {
		return nil, fmt.Errorf("expected indexed: 3 & data: 1, got: %d indexed & %d", indexdedLen, dataLen)
	}

	dst := string(e.Indexed[1])
	sn := big.NewInt(0).SetBytes(e.Indexed[2]).Uint64()

	return &providerTypes.Message{
		MessageHeight: height,
		EventType:     eventType,
		Dst:           dst,
		Data:          e.Data[0],
		Sn:            sn,
		Src:           p.NID(),
	}, nil
}

func (p *Provider) parseCallMessage(e *types.EventLog, eventType string, height uint64) (*providerTypes.Message, error) {
	if indexdedLen, dataLen := len(e.Indexed), len(e.Data); indexdedLen != 4 && dataLen != 2 {
		return nil, fmt.Errorf("expected indexed: 4 & data: 2, got: %d indexed & %d", indexdedLen, dataLen)
	}

	src := strings.SplitN(string(e.Indexed[1][:]), "/", 2)
	sn := big.NewInt(0).SetBytes(e.Indexed[2]).Uint64()
	reqID := big.NewInt(0).SetBytes(e.Data[0]).Uint64()

	return &providerTypes.Message{
		MessageHeight: height,
		ReqID:         reqID,
		EventType:     eventType,
		Dst:           p.NID(),
		Data:          e.Data[1],
		Sn:            sn,
		Src:           src[0],
	}, nil
}

// Parse Event
func (p *Provider) parseMessageEvent(notifications *types.EventNotification) ([]*providerTypes.Message, error) {
	height, err := notifications.Height.BigInt()
	if err != nil {
		return nil, err
	}
	var messages []*providerTypes.Message
	for _, event := range notifications.Logs {
		switch event.Indexed[0] {
		case EmitMessage:
			msg, err := p.parseEmitMessageEvent(height.Uint64(), event)
			if err != nil {
				return nil, err
			}
			messages = append(messages, msg)
		case CallMessage:
			msg, err := p.parseCallMessageEvent(height.Uint64(), event)
			if err != nil {
				return nil, err
			}
			messages = append(messages, msg)
		}
	}
	return messages, nil
}

// parseEmitMessage parses EmitMessage event
func (p *Provider) parseEmitMessageEvent(height uint64, e *types.EventNotificationLog) (*providerTypes.Message, error) {
	if indexdedLen, dataLen := len(e.Indexed), len(e.Data); indexdedLen != 3 && dataLen != 1 {
		return nil, fmt.Errorf("expected indexed: 3 & data: 1, got: %d indexed & %d", indexdedLen, dataLen)
	}

	dst := e.Indexed[1]
	sn, err := types.HexInt(e.Indexed[2]).BigInt()
	if err != nil {
		return nil, fmt.Errorf("failed to parse sn: %s", e.Indexed[2])
	}
	data := e.Data[0]

	return &providerTypes.Message{
		MessageHeight: height,
		EventType:     p.GetEventName(e.Indexed[0]),
		Dst:           dst,
		Data:          []byte(data),
		Sn:            sn.Uint64(),
		Src:           p.NID(),
	}, nil
}

// parseCallMessage parses CallMessage event
func (p *Provider) parseCallMessageEvent(height uint64, e *types.EventNotificationLog) (*providerTypes.Message, error) {
	if indexdedLen, dataLen := len(e.Indexed), len(e.Data); indexdedLen != 4 && dataLen != 2 {
		return nil, fmt.Errorf("expected indexed: 4 & data: 2, got: %d indexed & %d", indexdedLen, dataLen)
	}
	src := strings.SplitN(e.Indexed[1], "/", 2)
	sn, err := types.HexInt(e.Indexed[3]).BigInt()
	if err != nil {
		return nil, fmt.Errorf("failed to parse sn: %s", e.Indexed[2])
	}
	reqID, err := types.HexInt(e.Data[0]).BigInt()
	if err != nil {
		return nil, fmt.Errorf("failed to parse reqID: %s", e.Data[0])
	}
	data, err := types.HexBytes(e.Data[1]).Value()
	if err != nil {
		return nil, fmt.Errorf("failed to parse data: %s", e.Data[1])
	}

	return &providerTypes.Message{
		MessageHeight: height,
		ReqID:         reqID.Uint64(),
		EventType:     p.GetEventName(e.Indexed[0]),
		Dst:           p.NID(),
		Data:          data,
		Sn:            sn.Uint64(),
		Src:           src[0],
	}, nil
}
