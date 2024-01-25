package icon

import (
	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

// All the events
const (
	EmitMessage = "Message(str,int,bytes)"
	CallMessage = "CallMessage(str,str,int,int,bytes)"
)

func EventSigToEventType(sigContract map[string]string) map[string]string {
	eventMap := make(map[string]string, len(sigContract))
	for sig, contract := range sigContract {
		switch contract {
		case providerTypes.ConnectionContract:
			eventMap[sig] = events.EmitMessage
		case providerTypes.XcallContract:
			eventMap[sig] = events.CallMessage
		}
	}
	return eventMap
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
