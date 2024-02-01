package icon

import (
	"fmt"
	"math/big"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"go.uber.org/zap"
)

func (icp *IconProvider) parseMessagesFromEventlogs(log *zap.Logger, eventlogs []*types.EventLog, height uint64) []*providerTypes.Message {
	msgs := make([]*providerTypes.Message, 0)
	for _, el := range eventlogs {
		message, ok := icp.parseMessageFromEvent(log, el, height)
		if ok {
			msgs = append(msgs, message)
		}
	}
	return msgs
}

func (icp *IconProvider) parseMessageFromEvent(log *zap.Logger, event *types.EventLog, height uint64) (*providerTypes.Message, bool) {
	eventName := string(event.Indexed[0][:])
	eventType := icp.PCfg.GetEventName(eventName)

	switch eventName {
	case EmitMessage:
		m, err := icp.parseEmitMessage(event, eventType, height)
		if err != nil {
			log.Error("invalid event", zap.Error(err))
			return nil, false
		}
		return m, true
	}
	return nil, false
}

func (icp *IconProvider) parseEmitMessage(e *types.EventLog, eventType string, height uint64) (*providerTypes.Message, error) {
	if indexdedLen, dataLen := len(e.Indexed), len(e.Data); indexdedLen != 3 && dataLen != 1 {
		return nil, fmt.Errorf("expected indexed: 3 & data: 1, got: %d indexed & %d", indexdedLen, dataLen)
	}

	dst := string(e.Indexed[1][:])
	// TODO: temporary soln should be something permanent
	sn := big.NewInt(0).SetBytes(e.Indexed[2]).Uint64()

	return &providerTypes.Message{
		MessageHeight: height,
		EventType:     eventType,
		Dst:           dst,
		Data:          e.Data[0],
		Sn:            sn,
		Src:           icp.NID(),
	}, nil
}

func (p *IconProvider) parseCallMessage(e *types.EventLog, eventType string, height uint64) (*providerTypes.Message, error) {
	if indexdedLen, dataLen := len(e.Indexed), len(e.Data); indexdedLen != 3 && dataLen != 1 {
		return nil, fmt.Errorf("expected indexed: 3 & data: 1, got: %d indexed & %d", indexdedLen, dataLen)
	}

	p.log.Debug("detected eventlog ", zap.Uint64("height", height))
	return nil, nil
}
