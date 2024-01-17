package wasm

import (
	"fmt"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/hexstr"
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
					data, err := hexstr.NewFromString(attr.Value).ToByte()
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
