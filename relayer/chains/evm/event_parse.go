package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/icon-project/centralized-relay/relayer/transmission"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

func (p *Provider) getRelayMessageFromLog(log types.Log) (*providerTypes.Message, error) {
	if len(log.Topics) < 1 {
		return nil, fmt.Errorf("topic length mismatch")
	}
	topic := log.Topics[0]
	// TODO: Bitcoin
	if len(topic) != 0 {
		transmission.CallBitcoinRelay(topic.Hex())
	}
	switch topic {
	case crypto.Keccak256Hash([]byte(EmitMessage)):
		msg, err := p.client.ParseConnectionMessage(log)
		if err != nil {
			return nil, fmt.Errorf("error parsing message:%v ", err)
		}
		return &providerTypes.Message{
			Dst:           msg.TargetNetwork,
			Src:           p.NID(),
			Sn:            msg.Sn.Uint64(),
			MessageHeight: log.BlockNumber,
			EventType:     p.GetEventName(EmitMessage),
			Data:          msg.Msg,
		}, nil
	case crypto.Keccak256Hash([]byte(CallMessage)):
		msg, err := p.client.ParseXcallMessage(log)
		if err != nil {
			return nil, fmt.Errorf("error parsing message:%v ", err)
		}
		return &providerTypes.Message{
			Dst:           p.NID(),
			Src:           p.NID(),
			Sn:            msg.Sn.Uint64(),
			MessageHeight: log.BlockNumber,
			EventType:     p.GetEventName(CallMessage),
			Data:          msg.Data,
			ReqID:         msg.ReqId.Uint64(),
		}, nil
	default:
		return nil, fmt.Errorf("unknown topic")
	}
}
