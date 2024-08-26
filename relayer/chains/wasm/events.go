package wasm

import (
	"fmt"
	"math/big"

	abiTypes "github.com/cometbft/cometbft/abci/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	relayerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/hexstr"
	"go.uber.org/zap"
)

const (
	EventTypeWasmMessage         = "wasm-Message"
	EventTypeWasmCallMessage     = "wasm-CallMessage"
	EventTypeWasmRollbackMessage = "wasm-RollbackMessage"

	// Attr keys for connection contract events
	EventAttrKeyMsg                  = "msg"
	EventAttrKeyTargetNetwork string = "targetNetwork"
	EventAttrKeyConnSn        string = "connSn"

	// Attr keys for xcall contract events
	EventAttrKeyReqID string = "reqId"
	EventAttrKeyData  string = "data"
	EventAttrKeyTo    string = "to"
	EventAttrKeyFrom  string = "from"
	EventAttrKeySn    string = "sn"

	EventAttrKeyContractAddress string = "_contract_address"
)

func (p *Provider) ParseMessageFromEvents(eventsList []abiTypes.Event) ([]*relayerTypes.Message, error) {
	var messages []*relayerTypes.Message
	for _, ev := range eventsList {
		switch ev.Type {
		case EventTypeWasmMessage:
			msg := &relayerTypes.Message{
				EventType: events.EmitMessage,
				Src:       p.NID(),
			}
			for _, attr := range ev.Attributes {
				switch attr.Key {
				case EventAttrKeyMsg:
					data, err := hexstr.NewFromString(attr.Value).ToByte()
					if err != nil {
						return nil, fmt.Errorf("failed to parse msg data from event: %v", err)
					}
					msg.Data = data
				case EventAttrKeyConnSn:
					sn, ok := new(big.Int).SetString(attr.Value, 10)
					if !ok {
						return nil, fmt.Errorf("failed to parse connSn from event")
					}
					msg.Sn = sn
				case EventAttrKeyTargetNetwork:
					msg.Dst = attr.Value
				case EventAttrKeyFrom:
					msg.Src = attr.Value
				}
			}
			messages = append(messages, msg)
		case EventTypeWasmCallMessage:
			msg := &relayerTypes.Message{
				EventType: events.CallMessage,
				Dst:       p.NID(),
			}
			for _, attr := range ev.Attributes {
				switch attr.Key {
				case EventAttrKeyReqID:
					reqID, ok := new(big.Int).SetString(attr.Value, 10)
					if !ok {
						return nil, fmt.Errorf("failed to parse connSn from event")
					}
					msg.ReqID = reqID
				case EventAttrKeyData:
					msg.Data = []byte(attr.Value)
				case EventAttrKeyFrom:
					msg.Src = attr.Value
				case EventAttrKeySn:
					sn, ok := new(big.Int).SetString(attr.Value, 10)
					if !ok {
						return nil, fmt.Errorf("failed to parse connSn from event")
					}
					msg.Sn = sn
				}
			}
			messages = append(messages, msg)
		case EventTypeWasmRollbackMessage:
			msg := &relayerTypes.Message{
				EventType: events.RollbackMessage,
				Src:       p.NID(),
				Dst:       p.NID(),
			}
			for _, attr := range ev.Attributes {
				switch attr.Key {
				case EventAttrKeySn:
					sn, ok := new(big.Int).SetString(attr.Value, 10)
					if !ok {
						return nil, fmt.Errorf("failed to parse connSn from event")
					}
					msg.Sn = sn
				}
			}
			messages = append(messages, msg)
		default:
			p.logger.Debug("unknown event type", zap.String("type", ev.Type))
		}
	}
	return messages, nil
}

// EventSigToEventType converts event signature to event type
func (p *Config) eventMap() map[string]relayerTypes.EventMap {
	eventMap := make(map[string]relayerTypes.EventMap, len(p.Contracts))
	for contractName, addr := range p.Contracts {
		event := relayerTypes.EventMap{ContractName: contractName, Address: addr, SigType: make(map[string]string)}
		switch contractName {
		case relayerTypes.XcallContract:
			event.SigType[EventTypeWasmCallMessage] = events.CallMessage
			event.SigType[EventTypeWasmRollbackMessage] = events.RollbackMessage
		case relayerTypes.ConnectionContract:
			event.SigType[EventTypeWasmMessage] = events.EmitMessage
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

func (p *Config) GetMonitorEventFilters(eventMap map[string]relayerTypes.EventMap) []sdkTypes.Event {
	var eventList []sdkTypes.Event

	for addr, contract := range eventMap {
		for _, eventType := range contract.SigType {
			var wasmMessggeType string
			switch eventType {
			case events.EmitMessage:
				wasmMessggeType = EventTypeWasmMessage
			case events.CallMessage:
				wasmMessggeType = EventTypeWasmCallMessage
			case events.RollbackMessage:
				wasmMessggeType = EventTypeWasmRollbackMessage
			}
			eventList = append(eventList, sdkTypes.Event{
				Type: wasmMessggeType,
				Attributes: []abiTypes.EventAttribute{
					{
						Key:   EventAttrKeyContractAddress,
						Value: fmt.Sprintf("'%s'", addr),
					},
				},
			})
		}
	}
	return eventList
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
