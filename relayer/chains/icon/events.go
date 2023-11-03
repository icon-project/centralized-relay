package icon

import (
	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/relayer/events"
)

var (
	// All the events
	EmitMessage = "EmitMessage(str,bytes)"
)

var EventNameToType = map[string]string{
	EmitMessage: events.EmitMessage,
}

var MonitorEventsList []string = []string{
	//TODO: list all the events to monitor
	EmitMessage,
}

func GetMonitorEventFilters(address string, eventsList []string) []*types.EventFilter {

	filters := []*types.EventFilter{}
	if address == "" {
		return filters
	}

	for _, event := range eventsList {
		filters = append(filters, &types.EventFilter{
			Addr:      types.Address(address),
			Signature: event,
		})
	}
	return filters
}
