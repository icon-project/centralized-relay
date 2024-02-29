package icon

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func (p *IconProvider) parseMessagesFromEventlogs(log *zap.Logger, eventlogs []*types.EventLog, height uint64) []*providerTypes.Message {
	msgs := make([]*providerTypes.Message, 0)
	for _, el := range eventlogs {
		message, ok := p.parseMessageFromEvent(log, el, height)
		if ok {
			msgs = append(msgs, message)
		}
	}
	return msgs
}

func (p *IconProvider) parseMessageFromEvent(log *zap.Logger, event *types.EventLog, height uint64) (*providerTypes.Message, bool) {
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

func (p *IconProvider) parseEmitMessage(e *types.EventLog, eventType string, height uint64) (*providerTypes.Message, error) {
	if indexdedLen, dataLen := len(e.Indexed), len(e.Data); indexdedLen != 3 && dataLen != 1 {
		return nil, fmt.Errorf("expected indexed: 3 & data: 1, got: %d indexed & %d", indexdedLen, dataLen)
	}

	dst := string(e.Indexed[1][:])
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

func (p *IconProvider) parseCallMessage(e *types.EventLog, eventType string, height uint64) (*providerTypes.Message, error) {
	if indexdedLen, dataLen := len(e.Indexed), len(e.Data); indexdedLen != 4 && dataLen != 2 {
		return nil, fmt.Errorf("expected indexed: 3 & data: 1, got: %d indexed & %d", indexdedLen, dataLen)
	}

	src := strings.SplitN(string(e.Indexed[1][:]), "/", 2)
	sn := big.NewInt(0).SetBytes(e.Indexed[2]).Uint64()
	reqID := big.NewInt(0).SetBytes(e.Indexed[3]).Uint64()

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
