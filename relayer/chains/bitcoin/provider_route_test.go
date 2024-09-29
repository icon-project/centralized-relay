package bitcoin

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/bxelab/runestone"
	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/multisig"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
	"lukechampine.com/uint128"
)

func TestDecodeWithdrawToMessage(t *testing.T) {
	input := "+QEdAbkBGfkBFrMweDIuaWNvbi9jeGZjODZlZTc2ODdlMWJmNjgxYjU1NDhiMjY2Nzg0NDQ4NWMwZTcxOTK4PnRiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1gh6FAbhS+FCKV2l0aGRyYXdUb4MwOjC4PnRiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1ZPhIuEYweDIuYnRjL3RiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1"
	// Decode base64
	decodedBytes, _ := base64.StdEncoding.DecodeString(input)

	result, data, err := decodeWithdrawToMessage(decodedBytes)

	fmt.Println("data", data)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su", result.To)
	assert.Equal(t, big.NewInt(100).Bytes(), result.Amount)
	assert.Equal(t, "WithdrawTo", result.Action)
	assert.Equal(t, "0:0", result.TokenAddress)
}

func TestProvider_Route(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "bitcoin_provider_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := leveldb.OpenFile(dbPath, nil)
	assert.NoError(t, err)
	defer db.Close()

	logger, _ := zap.NewDevelopment()
	provider := &Provider{
		logger: logger,
		cfg:    &Config{Mode: SlaveMode},
		db:     db,
	}

	// Create a test message
	testMessage := &types.Message{
		Dst:           "destination",
		Src:           "source",
		Sn:            big.NewInt(123),
		Data:          []byte("test data"),
		MessageHeight: 456,
		EventType:     events.EmitMessage,
	}

	// Test storing the message
	err = provider.Route(context.Background(), testMessage, nil)
	assert.NoError(t, err)

	// Test retrieving the message
	key := []byte(fmt.Sprintf("bitcoin_message_%s", testMessage.Sn.String()))
	storedData, err := db.Get(key, nil)
	assert.NoError(t, err)

	var retrievedMessage types.Message
	err = json.Unmarshal(storedData, &retrievedMessage)
	assert.NoError(t, err)

	assert.Equal(t, testMessage.Dst, retrievedMessage.Dst)
	assert.Equal(t, testMessage.Src, retrievedMessage.Src)
	assert.Equal(t, testMessage.Sn.String(), retrievedMessage.Sn.String())
	assert.Equal(t, testMessage.Data, retrievedMessage.Data)
	assert.Equal(t, testMessage.MessageHeight, retrievedMessage.MessageHeight)
	assert.Equal(t, testMessage.EventType, retrievedMessage.EventType)

	// Test deleting the message
	err = db.Delete(key, nil)
	assert.NoError(t, err)

	_, err = db.Get(key, nil)
	assert.Error(t, err) // Should return an error as the key no longer exists
}

func TestDepositBitcoinToIcon(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*multisig.UTXO{
		{
			IsRelayersMultisig: false,
			TxHash:             "e57c10e27f75dbf0856163ca5f825b5af7ffbb3874f606b31330464ddd9df9a1",
			OutputIdx:          4,
			OutputAmount:       2974000,
		},
	}

	outputs := []*multisig.OutputTx{}

	// Add Bridge Message
	payload, _ := multisig.CreateBridgePayload(
		&multisig.XCallMessage{
			MessageType:  1,
			Action:       "Deposit",
			TokenAddress: "0:1",
			To:           "0x2.icon/hx452e235f9f1fd1006b1941ed1ad19ef51d1192f6",
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su",
			Amount:       new(big.Int).SetUint64(100000).Bytes(),
			Data:         []byte(""),
		},
		1,
		"cxfc86ee7687e1bf681b5548b2667844485c0e7192",
		[]string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	)
	scripts, _ := multisig.CreateBridgeMessageScripts(payload, 76)
	for i, script := range scripts {
		fmt.Println("OP_RETURN ", i, " script ", script)
		outputs = append(outputs, &multisig.OutputTx{
			OpReturnScript: script,
		})
	}

	// Add transfering bitcoin to relayer multisig
	outputs = append(outputs, &multisig.OutputTx{
		ReceiverAddress: "tb1pf0atpt2d3zel6udws38pkrh2e49vqd3c5jcud3a82srphnmpe55q0ecrzk",
		Amount:          100000,
	})

	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)
	rlMsAddress, _ := multisig.AddressOnChain(chainParam, userMultisigWallet)

	msAddressStr := rlMsAddress.String()
	fmt.Printf(msAddressStr)

	changeReceiverAddress := "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su"
	msgTx, _, txSigHashes, _ := multisig.CreateMultisigTx(inputs, outputs, 100000, &multisig.MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

	tapSigParams := multisig.TapSigParams{
		TxSigHashes:      txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf:  txscript.TapLeaf{},
		UserPKScript:     userMultisigWallet.PKScript,
		UserTapLeaf:      userMultisigWallet.TapLeaves[1],
	}

	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := multisig.PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)
	// COMBINE SIGNS
	signedMsgTx, err := multisig.CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// TODO: test the signedMsgTx
}

func TestDepositBitcoinToIconFail1(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*multisig.UTXO{
		{
			IsRelayersMultisig: false,
			TxHash:             "eeb8c9f79ecfe7c084b2af95bf82acebd130185a0d188283d78abb58d85eddff",
			OutputIdx:          4,
			OutputAmount:       2999000,
		},
	}

	outputs := []*multisig.OutputTx{}

	// Add Bridge Message
	payload, _ := multisig.CreateBridgePayload(
		&multisig.XCallMessage{
			Action:       "Deposit",
			TokenAddress: "0:1",
			To:           "0x2.icon/hx452e235f9f1fd1006b1941ed1ad19ef51d1192f6",
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su",
			Amount:       new(big.Int).SetUint64(1000000).Bytes(),
			Data:         []byte(""),
		},
		1,
		"cx8b52dfea0aa1e548288102df15ad7159f7266106",
		[]string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	)
	scripts, _ := multisig.CreateBridgeMessageScripts(payload, 76)
	for i, script := range scripts {
		fmt.Println("OP_RETURN ", i, " script ", script)
		outputs = append(outputs, &multisig.OutputTx{
			OpReturnScript: script,
		})
	}

	// Add transfering bitcoin to relayer multisig
	// outputs = append(outputs, &multisig.OutputTx{
	// 	ReceiverAddress: "tb1pf0atpt2d3zel6udws38pkrh2e49vqd3c5jcud3a82srphnmpe55q0ecrzk",
	// 	Amount:          1000000,
	// })

	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su"
	msgTx, _, txSigHashes, _ := multisig.CreateMultisigTx(inputs, outputs, 1000, &multisig.MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

	tapSigParams := multisig.TapSigParams{
		TxSigHashes:      txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf:  txscript.TapLeaf{},
		UserPKScript:     userMultisigWallet.PKScript,
		UserTapLeaf:      userMultisigWallet.TapLeaves[1],
	}

	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := multisig.PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)
	// COMBINE SIGNS
	signedMsgTx, err := multisig.CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// TODO: test the signedMsgTx
}

func TestDepositBitcoinToIconFail2(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*multisig.UTXO{
		{
			IsRelayersMultisig: false,
			TxHash:             "0416795b227e1a6a64eeb7bf7542d15964d18ac4c4732675d3189cda8d38bed7",
			OutputIdx:          3,
			OutputAmount:       2998000,
		},
	}

	outputs := []*multisig.OutputTx{}

	// Add Bridge Message
	payload, _ := multisig.CreateBridgePayload(
		&multisig.XCallMessage{
			Action:       "Deposit",
			TokenAddress: "0:1",
			To:           "0x2.icon/hx452e235f9f1fd1006b1941ed1ad19ef51d1192f6",
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su",
			Amount:       new(big.Int).SetUint64(1000000).Bytes(),
			Data:         []byte(""),
		},
		1,
		"cx8b52dfea0aa1e548288102df15ad7159f7266106",
		[]string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	)
	scripts, _ := multisig.CreateBridgeMessageScripts(payload, 76)
	for i, script := range scripts {
		fmt.Println("OP_RETURN ", i, " script ", script)
		outputs = append(outputs, &multisig.OutputTx{
			OpReturnScript: script,
		})
	}

	// Add transfering bitcoin to relayer multisig
	outputs = append(outputs, &multisig.OutputTx{
		ReceiverAddress: "tb1pf0atpt2d3zel6udws38pkrh2e49vqd3c5jcud3a82srphnmpe55q0ecrzk",
		Amount:          1000,
	})

	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su"
	msgTx, _, txSigHashes, _ := multisig.CreateMultisigTx(inputs, outputs, 1000, &multisig.MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

	tapSigParams := multisig.TapSigParams{
		TxSigHashes:      txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf:  txscript.TapLeaf{},
		UserPKScript:     userMultisigWallet.PKScript,
		UserTapLeaf:      userMultisigWallet.TapLeaves[1],
	}

	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := multisig.PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)
	// COMBINE SIGNS
	signedMsgTx, err := multisig.CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// TODO: test the signedMsgTx
}

func TestDepositBitcoinToIconFail3(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*multisig.UTXO{
		{
			IsRelayersMultisig: false,
			TxHash:             "dc21f89436d9fbda2cc521ed9b8988c7cbf84cdc67d728b2b2709a5efe7e775a",
			OutputIdx:          4,
			OutputAmount:       2996000,
		},
	}

	outputs := []*multisig.OutputTx{}

	// Add Bridge Message
	payload, _ := multisig.CreateBridgePayload(
		&multisig.XCallMessage{
			Action:       "Deposit",
			TokenAddress: "0:1",
			To:           "0x2.icon/hx452e235f9f1fd1006b1941ed1ad19ef51d1192f6",
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su",
			Amount:       new(big.Int).SetUint64(1000).Bytes(),
			Data:         []byte(""),
		},
		1,
		"cx8b52dfea0aa1e548288102df15ad7159f7266106",
		[]string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	)
	scripts, _ := multisig.CreateBridgeMessageScripts(payload, 76)
	for i, script := range scripts {
		fmt.Println("OP_RETURN ", i, " script ", script)
		outputs = append(outputs, &multisig.OutputTx{
			OpReturnScript: script,
		})
	}

	// Add transfering bitcoin to relayer multisig
	outputs = append(outputs, &multisig.OutputTx{
		ReceiverAddress: "tb1pf0atpt2d3zel6udws38pkrh2e49vqd3c5jcud3a82srphnmpe55q0ecrzk",
		Amount:          10000,
	})

	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su"
	msgTx, _, txSigHashes, _ := multisig.CreateMultisigTx(inputs, outputs, 1000, &multisig.MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

	tapSigParams := multisig.TapSigParams{
		TxSigHashes:      txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf:  txscript.TapLeaf{},
		UserPKScript:     userMultisigWallet.PKScript,
		UserTapLeaf:      userMultisigWallet.TapLeaves[1],
	}

	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := multisig.PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)
	// COMBINE SIGNS
	signedMsgTx, err := multisig.CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// TODO: test the signedMsgTx
}

func TestDepositBitcoinToIconFail4(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*multisig.UTXO{
		{
			IsRelayersMultisig: false,
			TxHash:             "1e29fa62942f92dd4cf688e219641b54229fbc2ec2dc74cc9c1c7f247c7172b2",
			OutputIdx:          4,
			OutputAmount:       2985000,
		},
	}

	outputs := []*multisig.OutputTx{}

	// Add Bridge Message
	payload, _ := multisig.CreateBridgePayload(
		&multisig.XCallMessage{
			Action:       "Withdraw",
			TokenAddress: "0:1",
			To:           "0x2.icon/hx452e235f9f1fd1006b1941ed1ad19ef51d1192f6",
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su",
			Amount:       new(big.Int).SetUint64(1000).Bytes(),
			Data:         []byte(""),
		},
		1,
		"cx8b52dfea0aa1e548288102df15ad7159f7266106",
		[]string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	)
	scripts, _ := multisig.CreateBridgeMessageScripts(payload, 76)
	for i, script := range scripts {
		fmt.Println("OP_RETURN ", i, " script ", script)
		outputs = append(outputs, &multisig.OutputTx{
			OpReturnScript: script,
		})
	}

	// Add transfering bitcoin to relayer multisig
	outputs = append(outputs, &multisig.OutputTx{
		ReceiverAddress: "tb1pf0atpt2d3zel6udws38pkrh2e49vqd3c5jcud3a82srphnmpe55q0ecrzk",
		Amount:          10000,
	})

	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su"
	msgTx, _, txSigHashes, _ := multisig.CreateMultisigTx(inputs, outputs, 1000, &multisig.MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

	tapSigParams := multisig.TapSigParams{
		TxSigHashes:      txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf:  txscript.TapLeaf{},
		UserPKScript:     userMultisigWallet.PKScript,
		UserTapLeaf:      userMultisigWallet.TapLeaves[1],
	}

	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := multisig.PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)
	// COMBINE SIGNS
	signedMsgTx, err := multisig.CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// TODO: test the signedMsgTx
}

// ... existing code ...

func TestDecodeBitcoinTransaction(t *testing.T) {
	// The raw transaction hex string
	rawTxHex := "02000000000101b272717c247f1c9ccc74dcc22ebc9f22541b6419e288f64cdd922f9462fa291e040000000001000000050000000000000000506a5e4c4cf88588576974686472617783303a31b83e74623170677a7838383079667237713864677a38647168773530736e6375346634686d7735636e3338303033353474757a6379396a7835736876760000000000000000506a5e4c4c377375b33078322e69636f6e2f6878343532653233356639663166643130303662313934316564316164313965663531643131393266368203e88001738b52dfea0aa1e548288102df15ad7100000000000000001d6a5e1a59f726610673577f5e756abd89cbcba38a58508b60a12754d2f510270000000000002251204bfab0ad4d88b3fd71ae844e1b0eeacd4ac03638a4b1c6c7a754061bcf61cd2830612d0000000000225120408c73bc891f8076a047682eea3e13c72a9adf6ea62713bdf1a557c1608591a903405a713aad72e1a2717cab16446e84a0de1c9f908a100c5903086a97910adc1327db65d21c8533c4fe087027195db4ebe7d8595624382e39b7aba829e3b29a29dc2551b275207303e7756826cf3fabc2f9c06978c542039effbb7493627cad22236e3ff10ee4ac41c18a23958fc9bf526c81a09bca529d685adc6900a04e2f520a26c63aa0b61a770f27c71b203eebab28e2e37b992abd8e6e5f9d887bdc2d5ab0efde76df79d3520400000000"

	// Decode the raw transaction
	txBytes, err := hex.DecodeString(rawTxHex)
	assert.NoError(t, err)

	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewReader(txBytes))
	assert.NoError(t, err)

	// Extract sender information (from the input)
	assert.Equal(t, 1, len(tx.TxIn), "Expected 1 input")
	prevOutHash := tx.TxIn[0].PreviousOutPoint.Hash.String()
	prevOutIndex := tx.TxIn[0].PreviousOutPoint.Index
	fmt.Printf("Sender (Input): %s:%d\n", prevOutHash, prevOutIndex)

	// Extract receiver information (from the outputs)
	assert.Equal(t, 5, len(tx.TxOut), "Expected 5 outputs")

	for i, out := range tx.TxOut {
		fmt.Printf("Output %d:\n", i)
		fmt.Printf("  Amount: %d satoshis\n", out.Value)

		// Attempt to parse the output script
		scriptClass, addresses, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, &chaincfg.TestNet3Params)
		if err != nil {
			fmt.Printf("  Script: Unable to parse (possibly OP_RETURN)\n")
		} else {
			fmt.Printf("  Script Class: %s\n", scriptClass)
			if len(addresses) > 0 {
				fmt.Printf("  Receiver Address: %s\n", addresses[0].String())
			}
		}

		// If it's an OP_RETURN output, print the data
		if scriptClass == txscript.NullDataTy {
			fmt.Printf("  OP_RETURN Data: %x\n", out.PkScript[2:])
		}

		fmt.Println()
	}

	// Add assertions for specific outputs
	assert.Equal(t, txscript.NullDataTy, txscript.GetScriptClass(tx.TxOut[0].PkScript))
	assert.Equal(t, txscript.NullDataTy, txscript.GetScriptClass(tx.TxOut[1].PkScript))
	assert.Equal(t, txscript.NullDataTy, txscript.GetScriptClass(tx.TxOut[2].PkScript))

	_, addresses, _, err := txscript.ExtractPkScriptAddrs(tx.TxOut[3].PkScript, &chaincfg.TestNet3Params)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(addresses))
	assert.Equal(t, "tb1pqqqqp399et2xygdj5xreqhjjvcmzhxw4aywxecjdzew6hylgvsesf3hn0c", addresses[0].String())

	_, addresses, _, err = txscript.ExtractPkScriptAddrs(tx.TxOut[4].PkScript, &chaincfg.TestNet3Params)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(addresses))
	assert.Equal(t, "tb1pqqqqp399et2xygdj5xreqhjjvcmzhxw4aywxecjdzew6hylgvsesf3hn0c", addresses[0].String())
}

func TestDepositRuneToIcon(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*multisig.UTXO{
		// user rune UTXOs to spend
		{
			IsRelayersMultisig: false,
			TxHash:             "d316231a8aa1f74472ed9cc0f1ed0e36b9b290254cf6b2c377f0d92b299868bf",
			OutputIdx:          0,
			OutputAmount:       1000,
		},
		// user bitcoin UTXOs to pay tx fee
		{
			IsRelayersMultisig: false,
			TxHash:             "4933e04e3d9320df6e9f046ff83cfc3e9f884d8811df0539af7aaca0218189aa",
			OutputIdx:          0,
			OutputAmount:       4000000,
		},
	}

	outputs := []*multisig.OutputTx{}

	// Add Bridge Message
	payload, _ := multisig.CreateBridgePayload(
		&multisig.XCallMessage{
			Action:       "Deposit",
			TokenAddress: "840000:3",
			To:           "0x2.icon/hx452e235f9f1fd1006b1941ed1ad19ef51d1192f6",
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su",
			Amount:       new(big.Int).SetUint64(1000000).Bytes(),
			Data:         []byte(""),
		},
		1,
		"cx8b52dfea0aa1e548288102df15ad7159f7266106",
		[]string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	)
	scripts, _ := multisig.CreateBridgeMessageScripts(payload, 76)
	for i, script := range scripts {
		fmt.Println("OP_RETURN ", i, " script ", script)
		outputs = append(outputs, &multisig.OutputTx{
			OpReturnScript: script,
		})
	}

	// Add transfering rune to relayer multisig
	runeId, _ := runestone.NewRuneId(840000, 3)
	changeReceiver := uint32(len(outputs) + 2)
	runeStone := &runestone.Runestone{
		Edicts: []runestone.Edict{
			{
				ID:     *runeId,
				Amount: uint128.From64(1000000000),
				Output: uint32(len(outputs) + 1),
			},
		},
		Pointer: &changeReceiver,
	}
	runeScript, _ := runeStone.Encipher()
	// Runestone OP_RETURN
	outputs = append(outputs, &multisig.OutputTx{
		OpReturnScript: runeScript,
	})
	// Rune UTXO send to relayer multisig
	outputs = append(outputs, &multisig.OutputTx{
		ReceiverAddress: "tb1pf0atpt2d3zel6udws38pkrh2e49vqd3c5jcud3a82srphnmpe55q0ecrzk",
		Amount:          1000,
	})
	// Rune change UTXO send back to user
	changeReceiverAddress := "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su"
	outputs = append(outputs, &multisig.OutputTx{
		ReceiverAddress: changeReceiverAddress,
		Amount:          1000,
	})

	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	msgTx, _, txSigHashes, _ := multisig.CreateMultisigTx(inputs, outputs, 1000, &multisig.MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

	tapSigParams := multisig.TapSigParams{
		TxSigHashes:      txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf:  txscript.TapLeaf{},
		UserPKScript:     userMultisigWallet.PKScript,
		UserTapLeaf:      userMultisigWallet.TapLeaves[1],
	}

	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := multisig.PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)
	// COMBINE SIGNS
	signedMsgTx, err := multisig.CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// Decipher runestone
	r := &runestone.Runestone{}
	artifact, err := r.Decipher(signedMsgTx)
	if err != nil {
		fmt.Println(err)
		return
	}
	a, _ := json.Marshal(artifact)
	fmt.Printf("Artifact: %s\n", string(a))
	// TODO: test the signedMsgTx
}

func TestCreateAndSignBitcoinTransaction(t *testing.T) {
	// Set up the network parameters (use TestNet3 for testing)
	chainParam := &chaincfg.TestNet3Params
	wif, _ := btcutil.DecodeWIF("your_private_key")
	// Get the private key from WIF
	privateKey := wif.PrivKey
	// Create a Taproot address
	internalKey := privateKey.PubKey()
	taprootKey := txscript.ComputeTaprootOutputKey(internalKey, nil)
	taprootAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(taprootKey), chainParam)
	assert.NoError(t, err)

	// Create the output script (witness program)
	pkScript, err := txscript.PayToAddrScript(taprootAddress)
	assert.NoError(t, err)

	prevOutputAmount := uint64(99800000)
	sendingAmount := uint64(100000)
	fees := uint64(200000)

	inputs := []*multisig.UTXO{
		{
			IsRelayersMultisig: false,
			TxHash:             "9dbd6f6f976f9f31895214c6c3034c80c567a38e1e816b8eb6bed972df0fdad9",
			OutputIdx:          4,
			OutputAmount:       prevOutputAmount,
		},
	}

	outputs := []*multisig.OutputTx{}

	// Add Bridge Message
	payload, _ := multisig.CreateBridgePayload(
		&multisig.XCallMessage{MessageType: 1,
			Action:       "Deposit",
			TokenAddress: "0:1",
			To:           "0x2.icon/hx452e235f9f1fd1006b1941ed1ad19ef51d1192f6",
			From:         "tb1peg65qks0qum848kq8udf3n3psvkpkaxsr7wq60ukr7w2symtt83qwd7cmx",
			Amount:       new(big.Int).SetUint64(sendingAmount).Bytes(),
			Data:         []byte(""),
		},
		1,
		"cx8b52dfea0aa1e548288102df15ad7159f7266106",
		[]string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	)
	scripts, _ := multisig.CreateBridgeMessageScripts(payload, 76)
	for i, script := range scripts {
		fmt.Println("OP_RETURN ", i, " script ", script)
		outputs = append(outputs, &multisig.OutputTx{
			OpReturnScript: script,
		})
	}
	// Add transfering bitcoin to relayer multisig
	outputs = append(outputs, &multisig.OutputTx{
		ReceiverAddress: "tb1pf0atpt2d3zel6udws38pkrh2e49vqd3c5jcud3a82srphnmpe55q0ecrzk",
		Amount:          sendingAmount,
	})

	changeReceiverAddress := "tb1peg65qks0qum848kq8udf3n3psvkpkaxsr7wq60ukr7w2symtt83qwd7cmx"
	msgTx, _, _, _ := multisig.CreateMultisigTx(inputs, outputs, fees, &multisig.MultisigWallet{}, &multisig.MultisigWallet{}, chainParam, changeReceiverAddress, 1)

	// Sign the transaction using Taproot
	for i, txIn := range msgTx.TxIn {
		// Create the signing hash
		prevOutputFetcher := txscript.NewCannedPrevOutputFetcher(pkScript, int64(prevOutputAmount))
		sigHashes := txscript.NewTxSigHashes(msgTx, prevOutputFetcher)

		// Create the Taproot signature
		witness, err := txscript.TaprootWitnessSignature(
			msgTx,
			sigHashes,
			i,
			int64(prevOutputAmount),
			pkScript,
			txscript.SigHashDefault,
			privateKey,
		)
		assert.NoError(t, err)

		// Set the witness data
		txIn.Witness = witness
	}

	// Serialize the transaction
	var signedTx bytes.Buffer
	err = msgTx.Serialize(&signedTx)
	assert.NoError(t, err)

	// Print the hex-encoded transaction
	fmt.Printf("Signed Taproot transaction: %x\n", signedTx.Bytes())
}
