package icon

import (
	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

// All the events
const (
	EmitMessage          = "Message(str,int,bytes)"
	CallMessage          = "CallMessage(str,str,int,int,bytes)"
	RollbackMessage      = "RollbackMessage(int)"
	AcknowledgeMessage   = "PacketAcknowledged(int,str,int,bytes,str)"
	TransmitreadyMessage = "TransactionReadyToExecute(str,int,bytes,str,str)"
)

// EventSigToEventType converts event signature to event type
func (p *Config) eventMap() map[string]providerTypes.EventMap {
	eventMap := make(map[string]providerTypes.EventMap, len(p.Contracts))
	for contractName, addr := range p.Contracts {
		event := providerTypes.EventMap{ContractName: contractName, Address: addr}
		switch contractName {
		case providerTypes.XcallContract:
			event.SigType = map[string]string{
				CallMessage:     events.CallMessage,
				RollbackMessage: events.RollbackMessage,
			}
		case providerTypes.ConnectionContract:
			event.SigType = map[string]string{EmitMessage: events.EmitMessage}
		case providerTypes.AggregationContract:
			event.SigType = map[string]string{
				AcknowledgeMessage:   events.AcknowledgeMessage,
				TransmitreadyMessage: events.TransmitreadyMessage,
			}
		}
		eventMap[addr] = event
	}
	return eventMap
}

// GetAddressNyEventType returns the address of the contract by event type
func (p *Provider) GetAddressByEventType(eventType string) types.Address {
	for _, contract := range p.contracts {
		for _, name := range contract.SigType {
			if name == eventType {
				return types.Address(contract.Address)
			}
		}
	}
	return ""
}

func (p *Provider) GetMonitorEventFilters() []*types.EventFilter {
	var filters []*types.EventFilter

	for addr, contract := range p.contracts {
		for sig := range contract.SigType {
			filters = append(filters, &types.EventFilter{
				Addr:      types.Address(addr),
				Signature: sig,
			})
		}
	}
	return filters
}

func (p *Provider) GetEventName(sig string) string {
	for _, contract := range p.contracts {
		for s, name := range contract.SigType {
			if s == sig {
				return name
			}
		}
	}
	return ""
}
