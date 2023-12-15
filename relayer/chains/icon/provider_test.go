package icon

import (
	"fmt"
	"testing"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func GetMockIconProvider() (*IconProvider, error) {
	pc := IconProviderConfig{
		NID:             "0x3.icon",
		KeyStore:        testKeyAddr,
		RPCUrl:          "http://localhost:9082/api/v3",
		Password:        testKeyPassword,
		StartHeight:     0,
		ContractAddress: "cxcacc844737024565cb56ac6ac8c1dab8fff1e2f7",
	}
	logger := zap.NewNop()
	prov, err := pc.NewProvider(logger, "", true, "icon-1")
	if err != nil {
		return nil, err
	}

	iconProvider, ok := prov.(*IconProvider)
	if !ok {
		return nil, fmt.Errorf("unbale to type case to icon chain provider")
	}
	return iconProvider, nil
}

// {�G7Ee�j�ڸ� [[77 101 115 115 97 103 101 40 115 116 114 44 105 110 116 44 98 121 116 101 115 41] [101 116 104]] [[1] [110 105 108 105 110]]}

func TestMessageFromEventlog(t *testing.T) {
	adr := "cxcacc844737024565cb56ac6ac8c1dab8fff1e2f7"
	pro, err := GetMockIconProvider()
	assert.NoError(t, err)
	eventlogs := &types.EventLog{
		Addr: types.NewAddress([]byte(adr)),
		Indexed: [][]byte{
			{77, 101, 115, 115, 97, 103, 101, 40, 115, 116, 114, 44, 105, 110, 116, 44, 98, 121, 116, 101, 115, 41},
			{101, 116, 104},
		},
		Data: [][]byte{
			{1},
			{110, 105, 108, 105, 110},
		},
	}
	logger := zap.NewNop()

	m, _ := pro.parseMessageFromEvent(logger, eventlogs, 20)
	fmt.Println("message", m)
}
