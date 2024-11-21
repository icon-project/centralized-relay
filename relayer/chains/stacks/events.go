package stacks

import (
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

const (
	EmitMessage      = "Message"
	CallMessage      = "CallMessage"
	RollbackMessage  = "RollbackMessage"
	CallExecuted     = "CallExecuted"
	CallMessageSent  = "CallMessageSent"
	ResponseMessage  = "ResponseMessage"
	RollbackExecuted = "RollbackExecuted"
)

func (c *Config) eventMap() map[string]providerTypes.EventMap {
	eventMap := make(map[string]providerTypes.EventMap, len(c.Contracts))
	for contractName, addr := range c.Contracts {
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

func (p *Provider) GetAddressByEventType(eventType string) string {
	for _, contract := range p.contracts {
		for _, name := range contract.SigType {
			if name == eventType {
				return contract.Address
			}
		}
	}
	return ""
}
