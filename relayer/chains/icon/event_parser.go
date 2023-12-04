package icon

import (
	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func parseMessagesFromEventlogs(log *zap.Logger, eventlogs []types.EventLog, height uint64) []*providerTypes.Message {
	msgs := make([]*providerTypes.Message, 0)
	for _, el := range eventlogs {
		message, ok := parseMessageFromEvent(log, el, height)
		if ok {
			msgs = append(msgs, message)
		}
	}
	return msgs
}

func parseMessageFromEvent(
	log *zap.Logger,
	event types.EventLog,
	height uint64,
) (*providerTypes.Message, bool) {
	eventName := string(event.Indexed[0][:])
	eventType := EventNameToType[eventName]

	switch eventName {
	case EmitMessage:
		m, err := parseEmitMessage(event, eventType, height)
		if err != nil {
			return nil, false
		}
		return m, true
	}
	return nil, false
}

func parseEmitMessage(e types.EventLog, eventType string, height uint64) (*providerTypes.Message, error) {
	if len(e.Indexed) != 2 && len(e.Data) != 2 {
		panic("Icon processor, emitMessage event is not correct")
	}

	dst := string(e.Indexed[1][:])
	// TODO: temporary soln should be something permanent
	sn := e.Data[0][0]

	return &providerTypes.Message{
		MessageHeight: height,
		EventType:     eventType,
		Dst:           dst,
		Data:          e.Data[1][:],
		Sn:            uint64(sn),
	}, nil
}
