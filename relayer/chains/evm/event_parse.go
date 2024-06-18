package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

func (p *Provider) getRelayMessageFromLog(log types.Log) (*providerTypes.Message, error) {
	if len(log.Topics) < 1 {
		return nil, fmt.Errorf("topic length mismatch")
	}
	topic := log.Topics[0]

	switch topic {
	case EmitMessageHash:
		msg, err := p.client.ParseConnectionMessage(log)
		if err != nil {
			return nil, fmt.Errorf("error parsing message:%v ", err)
		}
		return &providerTypes.Message{
			Dst:           msg.TargetNetwork,
			Src:           p.NID(),
			Sn:            msg.Sn,
			MessageHeight: log.BlockNumber,
			EventType:     p.GetEventName(EmitMessage),
			Data:          msg.Msg,
		}, nil
	case CallMessageHash:
		msg, err := p.client.ParseCallMessage(log)
		if err != nil {
			return nil, fmt.Errorf("error parsing message:%v ", err)
		}
		return &providerTypes.Message{
			Dst:           p.NID(),
			Src:           p.NID(),
			Sn:            msg.Sn,
			MessageHeight: log.BlockNumber,
			EventType:     p.GetEventName(CallMessage),
			Data:          msg.Data,
			ReqID:         msg.ReqId,
		}, nil
	case ExecuteRollbackHash:
		msg, err := p.client.ParseRollbackMessage(log)
		if err != nil {
			return nil, fmt.Errorf("error parsing message:%v ", err)
		}
		return &providerTypes.Message{
			Dst:           p.NID(),
			Src:           p.NID(),
			Sn:            msg.Sn,
			MessageHeight: log.BlockNumber,
			EventType:     p.GetEventName(ExecuteRollback),
		}, nil
	default:
		return nil, fmt.Errorf("unknown topic")
	}
}
