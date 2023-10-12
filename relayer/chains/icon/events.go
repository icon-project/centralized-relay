package icon

import (
	"encoding/hex"
	"strings"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/relayer/events"
)

// Events
var (
	// All the events
	EmitMessage = "EmitMessage(str,bytes)"
)

func MustConvertEventNameToBytes(eventName string) []byte {
	return []byte(eventName)
}

func ToEventLogBytes(evt types.EventLogStr) types.EventLog {
	indexed := make([][]byte, 0)

	for _, idx := range evt.Indexed {
		indexed = append(indexed, []byte(idx))
	}

	data := make([][]byte, 0)

	for _, d := range evt.Data {
		if isHexString(d) {
			filtered, _ := hex.DecodeString(strings.TrimPrefix(d, "0x"))
			data = append(data, filtered)
			continue
		}
		data = append(data, []byte(d))
	}

	return types.EventLog{
		Addr:    evt.Addr,
		Indexed: indexed,
		Data:    data,
	}

}

var EventTypesToName = map[string]string{
	EmitMessage: events.EmitMessage,
}

var MonitorEvents []string = []string{

	//TODO: list all the events to monitor
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
