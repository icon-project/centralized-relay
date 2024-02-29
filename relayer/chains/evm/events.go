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
	CallMessage = "CallMessage(string,string,uint256,uint256,bytes)"
)

// EventSigToEventType converts event signature to event type
func (p *EVMProviderConfig) eventMap() map[string]providerTypes.EventMap {
	eventMap := make(map[string]providerTypes.EventMap, len(p.Contracts))
	for contractName, addr := range p.Contracts {
		event := providerTypes.EventMap{ContractName: contractName, Address: addr}
		sig := make(map[string]string)
		switch contractName {
		case providerTypes.XcallContract:
			sig[CallMessage] = events.CallMessage
		case providerTypes.ConnectionContract:
			sig[EmitMessage] = events.EmitMessage
		}
		event.SigType = sig
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

func (p *EVMProvider) GetEventName(sig string) string {
	for _, contract := range p.contracts {
		if eventName, ok := contract.SigType[sig]; ok {
			return eventName
		}
	}
	return ""
}

func (p *EVMProvider) GetAddressByEventType(eventType string) *common.Address {
	for _, contract := range p.contracts {
		for _, eventName := range contract.SigType {
			if eventName == eventType {
				addr := common.HexToAddress(contract.Address)
				return &addr
			}
		}
	}
	return nil
}
