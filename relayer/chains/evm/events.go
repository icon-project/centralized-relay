package evm

import (
	"github.com/icon-project/centralized-relay/relayer/chains/evm/types"
	"github.com/icon-project/centralized-relay/relayer/events"
)

// All the events
var EmitMessage = "EmitMessage(str,bytes)"

func ToEventLogBytes(evt types.EventLogStr) types.EventLog {
	indexed := make([][]byte, 0)

	for _, idx := range evt.Indexed {
		indexed = append(indexed, []byte(idx))
	}

	data := make([][]byte, 0)

	for _, d := range evt.Data {
		data = append(data, []byte(d))
	}

	return types.EventLog{
		Addr:    evt.Addr,
		Indexed: indexed,
		Data:    data,
	}
}

var EventNameToType = map[string]string{
	EmitMessage: events.EmitMessage,
}

var MonitorEvents []string = []string{
	// TODO: list all the events to monitor
	EmitMessage,
}

func GetMonitorEventFilters(address string) []*types.EventFilter {
	filters := []*types.EventFilter{}
	if address == "" {
		return filters
	}

	for _, event := range MonitorEvents {
		filters = append(filters, &types.EventFilter{
			Addr:      types.Address(address),
			Signature: event,
		})
	}
	return filters
}
