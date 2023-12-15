package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

func (p *EVMProvider) getRelayMessageFromLog(log types.Log) (*providerTypes.Message, error) {
	if len(log.Topics) != 1 {
		return nil, fmt.Errorf("topic lenght mismatch")
	}
	topic := log.Topics[0]

	if topic == crypto.Keccak256Hash([]byte(EmitMessageSig)) {
		msg, err := p.client.ParseMessage(log)
		if err != nil {
			return nil, fmt.Errorf("error parsing message:%v ", err)
		}
		return &providerTypes.Message{
			Dst:           msg.TargetNetwork,
			Src:           p.NID(),
			Sn:            msg.Sn.Uint64(),
			MessageHeight: log.BlockNumber,
			EventType:     eventSigToEventType[topic],
			Data:          msg.Msg,
		}, nil
	}
	return nil, fmt.Errorf("failed to match eventname")
}
