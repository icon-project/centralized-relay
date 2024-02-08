package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

func (p *EVMProvider) getRelayMessageFromLog(log types.Log) (*providerTypes.Message, error) {
	if len(log.Topics) < 1 {
		return nil, fmt.Errorf("topic length mismatch")
	}
	topic := log.Topics[0]

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
			EventType:     p.cfg.GetEventName(EmitMessage),
			Data:          msg.Msg,
		}, nil
	case crypto.Keccak256Hash([]byte(CallMessage)):
		msg, err := p.client.ParseXcallMessage(log)
		if err != nil {
			return nil, fmt.Errorf("error parsing message:%v ", err)
		}
		return &providerTypes.Message{
			Dst:           msg.To.Hex(),
			Src:           msg.From.Hex(),
			Sn:            msg.Sn.Uint64(),
			MessageHeight: log.BlockNumber,
			EventType:     p.cfg.GetEventName(CallMessage),
			Data:          msg.Data,
			ReqID:         msg.ReqId.Uint64(),
		}, nil
	}
	return nil, fmt.Errorf("unknown topic")
}
