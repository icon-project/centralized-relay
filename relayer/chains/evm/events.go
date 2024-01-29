package evm

import (
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

// All the events
var (
	EmitMessage = "Message(string,uint256,bytes)"
	CallMessage = "CallMessage(str,str,uint256,uint256,bytes)"
)

// EventSigToEventType converts event signature to event type
func (p *EVMProviderConfig) eventMap() map[string]providerTypes.EventMap {
	eventMap := make(map[string]providerTypes.EventMap, len(p.Contracts))
	for contractName, addr := range p.Contracts {
		event := providerTypes.EventMap{ContractName: contractName}
		switch contractName {
		case providerTypes.XcallContract:
			event.SigType = map[string]string{CallMessage: events.CallMessage}
		case providerTypes.ConnectionContract:
			event.SigType = map[string]string{EmitMessage: events.EmitMessage}
		}
		eventMap[addr] = event
	}
	return eventMap
}

func (p *EVMProviderConfig) GetMonitorEventFilters() *ethereum.FilterQuery {
	filter := new(ethereum.FilterQuery)

	for addr, contract := range p.eventMap() {
		for sig := range contract.SigType {
			filter.Addresses = append(filter.Addresses, common.HexToAddress(addr))
			filter.Topics = [][]common.Hash{{crypto.Keccak256Hash([]byte(sig))}}
		}
	}
	return filter
}

func (p *EVMProviderConfig) GetEventName(sig string) string {
	for _, contract := range p.eventMap() {
		for s, name := range contract.SigType {
			if s == sig {
				return name
			}
		}
	}
	return ""
}
