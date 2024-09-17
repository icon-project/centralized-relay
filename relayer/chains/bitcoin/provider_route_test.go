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

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/multisig"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
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

func TestCreateBitcoinMultisigTx(t *testing.T) {
	provider := &Provider{
		logger: nil,
		cfg:    &Config{Mode: SlaveMode},
		db:     nil,
	}

	data := "+QEdAbkBGfkBFrMweDIuaWNvbi9jeGZjODZlZTc2ODdlMWJmNjgxYjU1NDhiMjY2Nzg0NDQ4NWMwZTcxOTK4PnRiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1gh6FAbhS+FCKV2l0aGRyYXdUb4MwOjC4PnRiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1ZPhIuEYweDIuYnRjL3RiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1"
	// Decode base64
	dataBytes, _ := base64.StdEncoding.DecodeString(data)

	chainParam := &chaincfg.TestNet3Params
	_, relayersMultisigInfo := multisig.RandomMultisigInfo(3, 3, chainParam, []int{0, 1, 2}, 0, 1)
	relayersMultisigWallet, _ := multisig.BuildMultisigWallet(relayersMultisigInfo)

	_, _, hexRawTx, _, err := provider.CreateBitcoinMultisigTx(dataBytes, 5000, relayersMultisigWallet, chainParam, UNISAT_DEFAULT_TESTNET)
	fmt.Println("err: ", err)
	fmt.Println("hexRawTx: ", hexRawTx)
}

func TestBuildAndPartSignBitcoinMessageTx(t *testing.T) {
	provider := &Provider{
		logger: nil,
		cfg:    &Config{Mode: SlaveMode},
		db:     nil,
	}

	data := "+QEdAbkBGfkBFrMweDIuaWNvbi9jeGZjODZlZTc2ODdlMWJmNjgxYjU1NDhiMjY2Nzg0NDQ4NWMwZTcxOTK4PnRiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1gh6FAbhS+FCKV2l0aGRyYXdUb4MwOjC4PnRiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1ZPhIuEYweDIuYnRjL3RiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1"
	// Decode base64
	dataBytes, _ := base64.StdEncoding.DecodeString(data)

	_, _, msgTx, _, err := provider.BuildAndPartSignBitcoinMessageTx(dataBytes, "0x2")
	fmt.Println("err: ", err)

	var rawTxBytes bytes.Buffer
	msgTx.Serialize(&rawTxBytes)

	hexRawTx := hex.EncodeToString(rawTxBytes.Bytes())
	fmt.Println("hexRawTx: ", hexRawTx)
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
			TxHash:        "4933e04e3d9320df6e9f046ff83cfc3e9f884d8811df0539af7aaca0218189aa",
			OutputIdx:     0,
			OutputAmount:  4000000,
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
		Amount:          1000000,
	})

	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su"
	msgTx, _, txSigHashes, _ := multisig.CreateMultisigTx(inputs, outputs, 1000, &multisig.MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)


	tapSigParams := multisig.TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
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
			TxHash:        "eeb8c9f79ecfe7c084b2af95bf82acebd130185a0d188283d78abb58d85eddff",
			OutputIdx:     4,
			OutputAmount:  2999000,
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


	tapSigParams := multisig.TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
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
			TxHash:        "0416795b227e1a6a64eeb7bf7542d15964d18ac4c4732675d3189cda8d38bed7",
			OutputIdx:     3,
			OutputAmount:  2998000,
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


	tapSigParams := multisig.TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
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
			TxHash:        "dc21f89436d9fbda2cc521ed9b8988c7cbf84cdc67d728b2b2709a5efe7e775a",
			OutputIdx:     4,
			OutputAmount:  2996000,
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


	tapSigParams := multisig.TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
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
			TxHash:        "1e29fa62942f92dd4cf688e219641b54229fbc2ec2dc74cc9c1c7f247c7172b2",
			OutputIdx:     4,
			OutputAmount:  2985000,
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


	tapSigParams := multisig.TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
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