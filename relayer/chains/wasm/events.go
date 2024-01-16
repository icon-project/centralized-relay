package wasm

import (
	"fmt"
	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"strconv"
)

const (
	EventTypeWasmMessage         string = "wasm-Message"
	EventTypeWasmCallMessageSent string = "wasm-CallMessageSent"

	EventAttrKeyMsg           string = "msg"
	EventAttrKeyTargetNetwork string = "targetNetwork"
	EventAttrKeySn            string = "sn"

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

func ParseMessageFromEvents(events []Event) (relayertypes.Message, error) {
	message := relayertypes.Message{}
	for _, ev := range events {
		switch ev.Type {
		case EventTypeWasmMessage:
			for _, attr := range ev.Attributes {
				if attr.Key == EventAttrKeyMsg {
					data, err := types.HexBytes(attr.Value).Value()
					if err != nil {
						return message, fmt.Errorf("failed to parse msg data from event: %v", err)
					}
					message.Data = data
				} else if attr.Key == EventAttrKeyTargetNetwork {
					message.Dst = attr.Value
				}
			}
		case EventTypeWasmCallMessageSent:
			for _, attr := range ev.Attributes {
				if attr.Key == EventAttrKeySn {
					sn, err := strconv.Atoi(attr.Value)
					if err != nil {
						return message, fmt.Errorf("failed to parse sn from event")
					}
					message.Sn = uint64(sn)
				}
			}
		}
	}
	return message, nil
}
