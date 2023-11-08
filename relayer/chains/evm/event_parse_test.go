package evm

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	pro, err := MockEvmProvider()
	assert.NoError(t, err)
	assert.NotNil(t, pro)

	fmt.Println(pro.wallet.Address)
}

func TestMessageTest(t *testing.T) {
	// sending the transaction

	pro, err := MockEvmProvider()
	assert.NoError(t, err)
	ctx := context.Background()
	// methodName := "receivemessage"
	// msg := providerTypes.Message{
	// 	Sn:   20,
	// 	Data: []byte("name"),
	// }
	opts, err := pro.GetTransationOpts(ctx)
	assert.NoError(t, err)

	// txhash, err := pro.transferBalance(
	// 	"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
	// 	pro.wallet.Address.Hex(), big.NewInt(100_000_000))

	// assert.NoError(t, err)

	tx, err := pro.client.SendMessage(opts, "icon", "--", big.NewInt(10), []byte("check"))
	assert.NoError(t, err)

	receipt, err := pro.WaitForResults(context.TODO(), tx.Hash())
	assert.NoError(t, err)

	fmt.Println("transaction status ", receipt.Status)
	fmt.Println("transaction hash", receipt.TxHash)
	fmt.Println("transaction height", receipt.BlockNumber)

	for _, m := range receipt.Logs {
		fmt.Println("transaction log ", m)
		msg, err := pro.client.ParseMessage(*m)
		if err != nil {
			fmt.Println("show the error ", err)
			continue
		}
		fmt.Println("the message is ", string(msg.Msg))
	}

	// gasPrice, err := pro.client.SuggestGasPrice(ctx)
	// assert.NoError(t, err)
	// fmt.Println("gas price is ", gasPrice)
}

func TestBlockInfo(t *testing.T) {

	pro, err := MockEvmProvider()
	assert.NoError(t, err)
	height := big.NewInt(4061)
	header, _ := pro.client.GetHeaderByHeight(context.TODO(), height)
	fmt.Println(header.Number)

	blockReq := getEventFilterQuery("0x0165878A594ca255338adfa4d48449f69242Eb8F")
	blockReq.FromBlock = height
	blockReq.ToBlock = height
	log, err := pro.client.FilterLogs(context.TODO(), blockReq)
	fmt.Println(len(log))
	for _, m := range log {
		fmt.Println("transaction log ", m)
		msg, err := pro.client.ParseMessage(m)
		if err != nil {
			fmt.Println("show the error ", err)
			continue
		}
		fmt.Println("the message is ", string(msg.Msg))
	}
}
