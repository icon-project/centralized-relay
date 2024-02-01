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

func (p *EVMProviderConfig) GetMonitorEventFilters() ethereum.FilterQuery {
	var (
		addresses []common.Address
		topics    []common.Hash
	)
	for addr, contract := range p.eventMap() {
		for sig := range contract.SigType {
			addresses = append(addresses, common.HexToAddress(addr))
			topics = append(topics, crypto.Keccak256Hash([]byte(sig)))
		}
	}
	return ethereum.FilterQuery{
		Addresses: addresses,
		Topics:    [][]common.Hash{topics},
	}
}

func (p *EVMProviderConfig) GetEventName(sig string) string {
	for _, contract := range p.eventMap() {
		if eventName, ok := contract.SigType[sig]; ok {
			return eventName
		}
	}
	return ""
}
