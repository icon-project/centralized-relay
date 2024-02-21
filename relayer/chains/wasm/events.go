package wasm

import (
	"fmt"
	"strconv"

	abiTypes "github.com/cometbft/cometbft/abci/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/hexstr"
)

const (
	EventTypeWasmMessage string = "wasm-Message"

	EventAttrKeyMsg           string = "msg"
	EventAttrKeyTargetNetwork string = "targetNetwork"
	EventAttrKeyConnSn        string = "connSn"
	EventAttrKeyReqID         string = "reqId"

	EventAttrKeyContractAddress string = "_contract_address"
)

type Event struct {
	Type       string      `json:"type"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EventsList struct {
	Events []Event `json:"events"`
}

func (p *Provider) ParseMessageFromEvents(eventsList []Event) ([]*providerTypes.Message, error) {
	var messages []*providerTypes.Message
	for _, ev := range eventsList {
		switch ev.Type {
		case EventTypeWasmMessage:
			msg := new(providerTypes.Message)
			for _, attr := range ev.Attributes {
				switch attr.Key {
				case EventAttrKeyMsg:
					data, err := hexstr.NewFromString(attr.Value).ToByte()
					if err != nil {
						return nil, fmt.Errorf("failed to parse msg data from event: %v", err)
					}
					msg.Data = data
					msg.EventType = events.EmitMessage
				case EventAttrKeyConnSn:
					sn, err := strconv.ParseUint(attr.Value, 10, strconv.IntSize)
					if err != nil {
						return nil, fmt.Errorf("failed to parse connSn from event")
					}
					msg.Sn = sn
				case EventAttrKeyTargetNetwork:
					msg.Dst = attr.Value
				case EventAttrKeyReqID:
					reqID, err := strconv.ParseUint(attr.Value, 10, strconv.IntSize)
					if err != nil {
						return nil, fmt.Errorf("failed to parse connSn from event")
					}
					msg.ReqID = reqID
					msg.EventType = events.CallMessage
				}
			}
			messages = append(messages, msg)
		}
	}
	return messages, nil
}

// EventSigToEventType converts event signature to event type
func (p *ProviderConfig) eventMap() map[string]providerTypes.EventMap {
	eventMap := make(map[string]providerTypes.EventMap, len(p.Contracts))
	for contractName, addr := range p.Contracts {
		event := providerTypes.EventMap{ContractName: contractName, Address: addr}
		switch contractName {
		case providerTypes.XcallContract:
			event.SigType = map[string]string{addr: events.CallMessage}
		case providerTypes.ConnectionContract:
			event.SigType = map[string]string{addr: events.EmitMessage}
		}
		eventMap[addr] = event
	}
	return eventMap
}

// GetAddressNyEventType returns the address of the contract by event type
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

func (p *Provider) GetMonitorEventFilters() []abiTypes.EventAttribute {
	var filters []abiTypes.EventAttribute

	for addr, contract := range p.contracts {
		for range contract.SigType {
			filters = append(filters, abiTypes.EventAttribute{
				Key:   EventAttrKeyContractAddress,
				Value: "'" + addr + "'",
			})
		}
	}
	return filters
}

func (p *Provider) GetEventName(addr string) string {
	for _, contract := range p.contracts {
		for a, name := range contract.SigType {
			if a == addr {
				return name
			}
		}
	}
	return ""
}
