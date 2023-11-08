package icon

import (
	"testing"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/stretchr/testify/assert"
)

func TestGetMonitorEventFilters(t *testing.T) {

	eventlist := []string{
		"test-event",
	}
	address := "cx000"

	assert.Equal(
		t,
		[]*types.EventFilter{
			{
				Addr:      types.Address(address),
				Signature: eventlist[0],
			},
		},
		GetMonitorEventFilters(address, eventlist),
	)

}
