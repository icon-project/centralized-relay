package icon

import (
	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/relayer/events"
)

// All the events
var EmitMessage = "Message(str,int,bytes)"

var EventNameToType = map[string]string{
	EmitMessage: events.EmitMessage,
}

var MonitorEventsList []string = []string{
	// TODO: list all the events to monitor
	EmitMessage,
}

func GetMonitorEventFilters(address string, eventsList []string) []*types.EventFilter {
	if address == "" {
		return nil
	}

	filters := []*types.EventFilter{}

	for _, event := range eventsList {
		filters = append(filters, &types.EventFilter{
			Addr:      types.Address(address),
			Signature: event,
		})
	}
	return filters
}
