package icon

import (
	"context"
	"fmt"
	"testing"

	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	relayerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func GetMockIconProvider() (*Provider, error) {
	pc := Config{
		NID:             "0x2.icon",
		NetworkID:       2,
		Address:         testKeyAddr,
		RPCUrl:          "https://lisbon.net.solidwallet.io/api/v3/",
		Password:        testKeyPassword,
		StartHeight:     0,
		ContractAddress: "cxb2b31a5252bfcc9be29441c626b8b918d578a58b",
	}
	logger := zap.NewNop()
	prov, err := pc.NewProvider(logger, "", true, "icon-1")
	if err != nil {
		return nil, err
	}

	iconProvider, ok := prov.(*Provider)
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

func TestReceiveMessage(t *testing.T) {
	pro, err := GetMockIconProvider()
	assert.NoError(t, err)
	key := relayerTypes.NewMessageKey(24, "0x13881.mumbai", "0x2.icon", "ss")
	receipt, err := pro.MessageReceived(context.TODO(), key)
	assert.NoError(t, err)
	fmt.Println(receipt)
}

func TestGenerateMessage(t *testing.T) {
	pro, err := GetMockIconProvider()
	assert.NoError(t, err)
	msg, err := pro.GenerateMessage(context.TODO(), &relayerTypes.MessageKeyWithMessageHeight{
		MessageKey: relayerTypes.MessageKey{Sn: 45, Src: "0x2.icon", Dst: "0x13881.mumbai", EventType: "emitMessage"}, Height: 31969244,
	})
	assert.NoError(t, err)
	fmt.Println("message is ", msg)

	// 42 0x2.icon 0x13881.mumbai emitMessage} 31968628

	// {"Sn":45,"Src":"0x2.icon","Dst":"0x13881.mumbai","EventType":"emitMessage","MsgHeight":31969244}
}
