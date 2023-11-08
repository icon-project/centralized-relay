package evm

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stretchr/testify/assert"
)

func TestParseEventLog(t *testing.T) {

	expected := providerTypes.Message{
		Dst:           "icon",
		Src:           "eth",
		Sn:            10,
		Data:          []byte("check"),
		MessageHeight: 4061,
		EventType:     events.EmitMessage,
	}

	data, _ := hex.DecodeString("7b2261646472657373223a22307830313635383738613539346361323535333338616466613464343834343966363932343265623866222c22746f70696373223a5b22307836646262623563383331383936373065303636643238316466633337643964656435313332616635643634303163666338333163373439396562373735663364225d2c2264617461223a22307830303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303630303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030613030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030613030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303034363936333666366530303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303536333638363536333662303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030222c22626c6f636b4e756d626572223a223078666464222c227472616e73616374696f6e48617368223a22307838633134616639663665333533346637643739383131323239353431383838653539346237373462373735656637303962633230623537353338373466616262222c227472616e73616374696f6e496e646578223a22307830222c22626c6f636b48617368223a22307831623635346561336663373063366538343532346261633238636139373838623761653461333563623934333363326334336537316633396264346134623337222c226c6f67496e646578223a22307830222c2272656d6f766564223a66616c73657d")
	var log ethTypes.Log
	err := json.Unmarshal(data, &log)
	assert.NoError(t, err)

	pro, err := MockEvmProvider("")
	assert.NoError(t, err)

	msg, err := pro.getRelayMessageFromLog(log)
	assert.NoError(t, err)
	assert.Equal(t, expected, msg)
	fmt.Println(msg)

}

// func TestSendMessageTest(t *testing.T) {
// 	// sending the transaction

// 	pro, err := MockEvmProvider("0x0165878A594ca255338adfa4d48449f69242Eb8F")
// 	assert.NoError(t, err)
// 	ctx := context.Background()
// 	opts, err := pro.GetTransationOpts(ctx)
// 	assert.NoError(t, err)

// 	tx, err := pro.client.SendMessage(opts, "icon", "--", big.NewInt(10), []byte("check"))
// 	assert.NoError(t, err)

// 	receipt, err := pro.WaitForResults(context.TODO(), tx.Hash())
// 	assert.NoError(t, err)

// 	for _, m := range receipt.Logs {
// 		fmt.Println("transaction log ", m)
// 		msg, err := pro.client.ParseMessage(*m)
// 		if err != nil {
// 			fmt.Println("show the error ", err)
// 			continue
// 		}
// 		fmt.Println("the message is ", string(msg.Msg))
// 	}

// }

// func TestFilterLog(t *testing.T) {

// 	pro, err := MockEvmProvider("0x0165878A594ca255338adfa4d48449f69242Eb8F")
// 	assert.NoError(t, err)
// 	height := big.NewInt(4061)
// 	header, _ := pro.client.GetHeaderByHeight(context.TODO(), height)
// 	fmt.Println(header.Number)

// 	blockReq := getEventFilterQuery("0x0165878A594ca255338adfa4d48449f69242Eb8F")
// 	blockReq.FromBlock = height
// 	blockReq.ToBlock = height
// 	log, err := pro.client.FilterLogs(context.TODO(), blockReq)
// 	fmt.Println(len(log))

// 	for _, m := range log {
// 		fmt.Println("transaction log ", m)
// 		b, _ := json.Marshal(m)

// 		fmt.Printf("log data: %x", b)
// 		msg, err := pro.client.ParseMessage(m)
// 		if err != nil {
// 			fmt.Println("show the error ", err)
// 			continue
// 		}
// 		fmt.Println("the message is ", string(msg.Msg))
// 	}
// }
