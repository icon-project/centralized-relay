package icon

import (
	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

// All the events
const (
	EmitMessage = "Message(str,int,bytes)"
	CallMessage = "CallMessage(str,str,int,int,bytes)"
)

// EventSigToEventType converts event signature to event type
func (p *IconProviderConfig) eventMap() map[string]providerTypes.EventMap {
	eventMap := make(map[string]providerTypes.EventMap, len(p.Contracts))
	for contractName, addr := range p.Contracts {
		event := providerTypes.EventMap{ContractName: contractName, Address: addr}
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

// GetAddressNyEventType returns the address of the contract by event type
func (p *IconProviderConfig) GetAddressByEventType(eventType string) types.Address {
	for _, contract := range p.eventMap() {
		for _, name := range contract.SigType {
			if name == eventType {
				return types.Address(contract.Address)
			}
		}
	}
	return ""
}

func (p *IconProviderConfig) GetMonitorEventFilters() []*types.EventFilter {
	var filters []*types.EventFilter

	for addr, contract := range p.eventMap() {
		for sig := range contract.SigType {
			filters = append(filters, &types.EventFilter{
				Addr:      types.Address(addr),
				Signature: sig,
			})
		}
	}
	return filters
}

func (p *IconProviderConfig) GetEventName(sig string) string {
	for _, contract := range p.eventMap() {
		for s, name := range contract.SigType {
			if s == sig {
				return name
			}
		}
	}
	return ""
}
