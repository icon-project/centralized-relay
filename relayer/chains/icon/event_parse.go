package icon

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	relayerevents "github.com/icon-project/centralized-relay/relayer/events"
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
	sn := new(big.Int).SetBytes(e.Indexed[2])

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
	sn := new(big.Int).SetBytes(e.Indexed[2])
	reqID := new(big.Int).SetBytes(e.Data[0])

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
	aboveHeight, err := notifications.Height.BigInt()
	if err != nil {
		return nil, err
	}
	height := aboveHeight.Sub(aboveHeight, big.NewInt(1))

	var messages []*providerTypes.Message
	for _, event := range notifications.Logs {
		msg, err := p.parseMessageFromEventLog(height.Uint64(), event)
		if err != nil {
			p.log.Warn("received invalid event", zap.Error(err))
		} else if msg != nil {
			messages = append(messages, msg)
		}
	}
	return messages, nil
}

func (p *Provider) parseMessageFromEventLog(height uint64, event *types.EventNotificationLog) (*providerTypes.Message, error) {
	switch event.Indexed[0] {
	case EmitMessage:
		return p.parseEmitMessageEvent(height, event)
	case CallMessage:
		return p.parseCallMessageEvent(height, event)
	case RollbackMessage:
		return p.parseRollbackMessageEvent(height, event)
	case PacketRegistered:
		return p.parsePacketRegisteredEvent(height, event)
	case PacketAcknowledged:
		return p.parsePacketAcknowledgedEvent(height, event)
	default:
		return nil, fmt.Errorf("invalid event: %s", event.Indexed[0])
	}
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
	data, err := types.HexBytes(e.Data[0]).Value()
	if err != nil {
		return nil, fmt.Errorf("failed to parse data: %s", e.Data[0])
	}

	return &providerTypes.Message{
		MessageHeight: height,
		EventType:     p.GetEventName(e.Indexed[0]),
		Dst:           dst,
		Data:          data,
		Sn:            sn,
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
		ReqID:         reqID,
		EventType:     p.GetEventName(e.Indexed[0]),
		Dst:           p.NID(),
		Data:          data,
		Sn:            sn,
		Src:           src[0],
	}, nil
}

// parseRollbackMessage parses RollbackMessage event
func (p *Provider) parseRollbackMessageEvent(height uint64, e *types.EventNotificationLog) (*providerTypes.Message, error) {
	if indexdedLen := len(e.Indexed); indexdedLen != 2 {
		return nil, fmt.Errorf("expected indexed: 2, got: %d indexed", indexdedLen)
	}
	sn, err := types.HexInt(e.Indexed[1]).BigInt()
	if err != nil {
		return nil, fmt.Errorf("failed to parse sn: %s", e.Indexed[1])
	}

	return &providerTypes.Message{
		MessageHeight: height,
		EventType:     p.GetEventName(e.Indexed[0]),
		Dst:           p.NID(),
		Src:           p.NID(),
		Sn:            sn,
	}, nil
}

func (p *Provider) parsePacketRegisteredEvent(height uint64, e *types.EventNotificationLog) (*providerTypes.Message, error) {
	if indexdedLen := len(e.Indexed); indexdedLen != 3 {
		return nil, fmt.Errorf("expected indexed: 3 got: %d indexed", indexdedLen)
	}
	srcNetwork := e.Indexed[1]
	srcConnAddress := e.Indexed[2]

	srcConnSn, err := types.HexInt(e.Data[0]).BigInt()
	if err != nil {
		return nil, fmt.Errorf("failed to parse srcConnSn: %s", e.Data[0])
	}

	srcHeight, err := types.HexInt(e.Data[1]).BigInt()
	if err != nil {
		return nil, fmt.Errorf("failed to parse srcHeight: %s", e.Data[1])
	}

	dstNetwork := e.Data[2]
	data, err := types.HexBytes(e.Data[3]).Value()
	if err != nil {
		return nil, fmt.Errorf("failed to parse data: %s", e.Data[3])
	}
	dstConnAddress := e.Data[4]

	return &providerTypes.Message{
		EventType:      relayerevents.EmitMessage,
		Src:            srcNetwork,
		Sn:             srcConnSn,
		SrcConnAddress: srcConnAddress,
		MessageHeight:  srcHeight.Uint64(),

		Dst:            dstNetwork,
		DstConnAddress: dstConnAddress,
		Data:           data,

		AggregatorHeight: height,
		AggregatorEvent:  p.GetEventName(e.Indexed[0]),
	}, nil
}

func (p *Provider) parsePacketAcknowledgedEvent(height uint64, e *types.EventNotificationLog) (*providerTypes.Message, error) {
	if indexdedLen := len(e.Indexed); indexdedLen != 3 {
		return nil, fmt.Errorf("expected indexed: 3 got: %d indexed", indexdedLen)
	}
	srcNetwork := e.Indexed[1]
	srcConnAddress := e.Indexed[2]

	srcConnSn, err := types.HexInt(e.Data[0]).BigInt()
	if err != nil {
		return nil, fmt.Errorf("failed to parse srcConnSn: %s", e.Data[0])
	}

	srcHeight, err := types.HexInt(e.Data[1]).BigInt()
	if err != nil {
		return nil, fmt.Errorf("failed to parse srcHeight: %s", e.Data[1])
	}

	dstNetwork := e.Data[2]
	data, err := types.HexBytes(e.Data[3]).Value()
	if err != nil {
		return nil, fmt.Errorf("failed to parse data: %s", e.Data[3])
	}
	dstConnAddress := e.Data[4]

	signaturesRlp, err := types.HexBytes(e.Data[5]).Value()
	if err != nil {
		return nil, fmt.Errorf("failed to parse signatures: %s", e.Data[5])
	}
	var signatures [][]byte
	if err := rlp.DecodeBytes(signaturesRlp, &signatures); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to parse rlp signatures: %s", e.Data[5])
	}

	return &providerTypes.Message{
		EventType:      relayerevents.EmitMessage,
		Src:            srcNetwork,
		Sn:             srcConnSn,
		SrcConnAddress: srcConnAddress,
		MessageHeight:  srcHeight.Uint64(),

		Dst:            dstNetwork,
		DstConnAddress: dstConnAddress,
		Data:           data,
		Signatures:     signatures,

		AggregatorHeight: height,
		AggregatorEvent:  p.GetEventName(e.Indexed[0]),
	}, nil
}
