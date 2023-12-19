package evm

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func MockEvmProvider(contractAddress string) (*EVMProvider, error) {
	evm := EVMProviderConfig{
		NID:             "0x13881.mumbai",
		ChainName:       "avalanche",
		RPCUrl:          "https://rpc-mumbai.maticvigil.com",
		StartHeight:     0,
		Keystore:        testKeyStore,
		Password:        testKeyPassword,
		GasPrice:        3000000015,
		GasLimit:        431877,
		ContractAddress: contractAddress,
	}
	log := zap.NewNop()
	pro, err := evm.NewProvider(log, "", true, "avalanche")
	if err != nil {
		return nil, err
	}
	p, ok := pro.(*EVMProvider)
	if !ok {
		return nil, fmt.Errorf("failed to create mock evmprovider")
	}

	p.Init(context.TODO())
	return p, nil
}

func TestTransferBalance(t *testing.T) {
	pro, err := MockEvmProvider("0x0165878A594ca255338adfa4d48449f69242Eb8F")
	assert.NoError(t, err)

	txhash, err := pro.transferBalance(
		"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		pro.wallet.Address.Hex(), big.NewInt(100_000_000_000_000_000_0))

	assert.NoError(t, err)

	r, err := pro.WaitForResults(context.TODO(), txhash)
	assert.NoError(t, err)
	fmt.Println("status of the transaction ", r.Status)
	fmt.Println("transaction hash", r.TxHash)
}

func TestRouteMessage(t *testing.T) {
	pro, err := MockEvmProvider("0x64FDC0B87019cEeA603f9DD559b9bAd31F1157b8")
	assert.NoError(t, err)

	d, _ := hex.DecodeString("f89501b892f890b33078322e69636f6e2f687862366235373931626530623565663637303633623363313062383430666238313531346462326664aa30783865383362316430613264656636646433633735346139316638633038636466386139313766343282019e0061ebaa307836344644433042383730313963456541363033663944443535396239624164333146313135376238")

	expected := &providerTypes.Message{
		Dst:           "0x13881.mumbai",
		Src:           "0x2.icon",
		Sn:            28,
		Data:          d,
		MessageHeight: 31965963,
		EventType:     events.EmitMessage,
	}

	var callback providerTypes.TxResponseFunc

	callback = func(key providerTypes.MessageKey, response providerTypes.TxResponse, err error) {
		fmt.Println("response obj", response)
		if response.Code != 1 {
			assert.Fail(t, "transaction failed")
		}
	}

	err = pro.Route(context.TODO(), expected, callback)
	assert.NoError(t, err)
}

func TestSendMessageTest(t *testing.T) {
	// sending the transaction

	pro, err := MockEvmProvider("e7f1725E7734CE288F8367e1Bb143E90bb3F0512")
	assert.NoError(t, err)
	ctx := context.Background()
	opts, err := pro.GetTransationOpts(ctx)
	assert.NoError(t, err)

	tx, err := pro.client.SendMessage(opts, "icon", "--", big.NewInt(19), []byte("check"))
	assert.NoError(t, err)

	receipt, err := pro.WaitForResults(context.TODO(), tx.Hash())
	assert.NoError(t, err)
	fmt.Println("receipt blocknumber  is:", receipt.BlockNumber)

	for _, m := range receipt.Logs {
		msg, err := pro.client.ParseMessage(*m)
		if err != nil {
			fmt.Println("show the error ", err)
			continue
		}
		fmt.Println("the message is ", string(msg.Msg))
	}
}

func TestEventLogReceived(t *testing.T) {
	mock, err := MockEvmProvider("0x64FDC0B87019cEeA603f9DD559b9bAd31F1157b8")

	assert.NoError(t, err)

	ht := big.NewInt(43587936)
	ht2 := big.NewInt(43587936)
	blockRequest := mock.blockReq
	blockRequest.ToBlock = ht2
	blockRequest.FromBlock = ht

	log, err := mock.client.FilterLogs(context.TODO(), blockRequest)
	assert.NoError(t, err)

	fmt.Println("logs is ", len(log))
	for _, log := range log {
		message, err := mock.getRelayMessageFromLog(log)
		assert.NoError(t, err)
		// p.log.Info("message received evm: ", zap.Uint64("height", lbn.Height.Uint64()),
		// 	zap.String("target-network", message.Dst),
		// 	zap.Uint64("sn", message.Sn),
		// 	zap.String("event-type", message.EventType),
		// )
		fmt.Println("message", message)
	}
}
