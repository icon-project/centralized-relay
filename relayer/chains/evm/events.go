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
	EmitMessage     = "Message(string,uint256,bytes)"
	CallMessage     = "CallMessage(string,string,uint256,uint256,bytes)"
	RollbackMessage = "RollbackMessage(uint256)"

	EmitMessageHash     = crypto.Keccak256Hash([]byte(EmitMessage))
	CallMessageHash     = crypto.Keccak256Hash([]byte(CallMessage))
	RollbackMessageHash = crypto.Keccak256Hash([]byte(RollbackMessage))
)

// EventSigToEventType converts event signature to event type
func (p *Config) eventMap() map[string]providerTypes.EventMap {
	eventMap := make(map[string]providerTypes.EventMap, len(p.Contracts))
	for contractName, addr := range p.Contracts {
		event := providerTypes.EventMap{ContractName: contractName, Address: addr}
		sig := make(map[string]string)
		switch contractName {
		case providerTypes.XcallContract:
			sig[CallMessage] = events.CallMessage
			sig[RollbackMessage] = events.RollbackMessage
		case providerTypes.ConnectionContract:
			sig[EmitMessage] = events.EmitMessage
		}
		event.SigType = sig
		eventMap[addr] = event
	}
	return eventMap
}

func (p *Config) GetMonitorEventFilters() ethereum.FilterQuery {
	var (
		addresses []common.Address
		topics    []common.Hash
	)
	for addr, contract := range p.eventMap() {
		addresses = append(addresses, common.HexToAddress(addr))
		for sig := range contract.SigType {
			topics = append(topics, crypto.Keccak256Hash([]byte(sig)))
		}
	}
	return ethereum.FilterQuery{
		Addresses: addresses,
		Topics:    [][]common.Hash{topics},
	}
}

func (p *Provider) GetEventName(sig string) string {
	for _, contract := range p.contracts {
		if eventName, ok := contract.SigType[sig]; ok {
			return eventName
		}
	}
	return ""
}

func (p *Provider) GetAddressByEventType(eventType string) *common.Address {
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
