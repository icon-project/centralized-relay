package wasm

import (
	"fmt"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/hexstr"
	"strconv"
)

const (
	EventTypeWasmMessage string = "wasm-Message"

	EventAttrKeyMsg           string = "msg"
	EventAttrKeyTargetNetwork string = "targetNetwork"
	EventAttrKeyConnSn        string = "connSn"

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
				switch attr.Key {
				case EventAttrKeyMsg:
					data, err := hexstr.NewFromString(attr.Value).ToByte()
					if err != nil {
						return message, fmt.Errorf("failed to parse msg data from event: %v", err)
					}
					message.Data = data
				case EventAttrKeyConnSn:
					connSn, err := strconv.Atoi(attr.Value)
					if err != nil {
						return message, fmt.Errorf("failed to parse connSn from event")
					}
					message.Sn = uint64(connSn)
				case EventAttrKeyTargetNetwork:
					message.Dst = attr.Value
				}
			}
		}
	}
	return message, nil
}
