package bitcoin

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/bxelab/runestone"
	"github.com/icon-project/centralized-relay/relayer/chains/icon"
	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/multisig"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

func initBtcProvider() (*Provider, string) {
	tempDir, _ := os.MkdirTemp("", "bitcoin_provider_test")

	dbPath := filepath.Join(tempDir, "test.db")
	db, _ := leveldb.OpenFile(dbPath, nil)

	logger, _ := zap.NewDevelopment()

	config := &Config{}
	config = &Config{
		Connections:      []string{"cx577f5e756abd89cbcba38a58508b60a12754d2f5"},
		UniSatURL:        "https://open-api-testnet.unisat.io",
		UniSatKey:        "60b7bf52654454f19d8553e1b6427fb9fd2c722ea8dc6822bdf1dd7615b4b35d",
		RecoveryLockTime: 1234,
		MasterPubKey:     "02fe44ec9f26b97ed30bd33898cf22de726e05389bde632d3aa6ad6746e15221d2",
		Slave1PubKey:     "0230edd881db1bc32b94f83ea5799c2e959854e0f99427d07c211206abd876d052",
		Slave2PubKey:     "021e83d56728fde393b41b74f2b859381661025f2ecec567cf392da7372de47833",
		MempoolURL:       "https://mempool.space/testnet/api/v1",
		SlaveServer1:     "http://18.234.252.46:8081",
		SlaveServer2:     "http://3.81.55.23:8082",
		Port:             "8080",
		RequestTimeout:   1000,
		ApiKey:           "key",
	}
	config.NID = "0x2.btc"
	config.Address = "tb1pf0atpt2d3zel6udws38pkrh2e49vqd3c5jcud3a82srphnmpe55q0ecrzk"
	// config.RPCUrl = "stylish-white-paper.btc-testnet.quiknode.pro/f87ecf84754d9daf8cf45b0334ba3fee9b356db6/"
	config.RPCUrl = "http://lydialabs:lydialabsq23%21%40@54.92.152.113:18332"
	config.User = "lydialabs"
	config.Password = "lydialabsq23!@"
	config.MasterPubKey = "02fe44ec9f26b97ed30bd33898cf22de726e05389bde632d3aa6ad6746e15221d2"
	config.Slave1PubKey = "0230edd881db1bc32b94f83ea5799c2e959854e0f99427d07c211206abd876d052"
	config.Slave2PubKey = "021e83d56728fde393b41b74f2b859381661025f2ecec567cf392da7372de47833"
	config.RelayerPrivKey = "cTYRscQxVhtsGjHeV59RHQJbzNnJHbf3FX4eyX5JkpDhqKdhtRvy"
	config.RecoveryLockTime = 1234
	config.OpCode = 0x5e

	connConfig := &rpcclient.ConnConfig{
		Host:         config.RPCUrl,
		User:         config.User,
		Pass:         config.Password,
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	client, err := rpcclient.New(connConfig, nil)
	if err != nil {
		fmt.Println("err: ", err)
		return nil, ""
	}

	provider := &Provider{
		logger:     logger,
		db:         db,
		cfg:        config,
		client:     &Client{client: client, log: logger},
		chainParam: &chaincfg.TestNet3Params,
	}

	msPubkey, _ := btcutil.DecodeAddress(provider.cfg.Address, provider.chainParam)
	multisigAddressScript, _ := txscript.PayToAddrScript(msPubkey)
	provider.multisigAddrScript = multisigAddressScript
	return provider, tempDir
}

func TestParseMessageFromTx(t *testing.T) {
	// Create a mock Provider
	chainParam := &chaincfg.TestNet3Params
	logger, _ := zap.NewDevelopment()
	provider := &Provider{
		logger: logger,
		cfg: &Config{
			Connections: []string{"cx577f5e756abd89cbcba38a58508b60a12754d2f5"},
		},
		chainParam: &chaincfg.TestNet3Params,
	}
	provider.cfg.NID = "0x2.btc"

	provider.cfg.Address = "tb1pf0atpt2d3zel6udws38pkrh2e49vqd3c5jcud3a82srphnmpe55q0ecrzk"
	msPubkey, _ := btcutil.DecodeAddress(provider.cfg.Address, chainParam)
	multisigAddressScript, _ := txscript.PayToAddrScript(msPubkey)
	provider.multisigAddrScript = multisigAddressScript
	// Create a mock TxSearchRes
	// relayer multisig
	decodedAddr, _ := btcutil.DecodeAddress(RELAYER_MULTISIG_ADDRESS, chainParam)
	relayerPkScript, _ := txscript.PayToAddrScript(decodedAddr)
	// user key
	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	bridgeMsg := multisig.BridgeDecodedMsg{
		Message: &multisig.XCallMessage{
			MessageType:  0,
			Action:       MethodDeposit,
			TokenAddress: "0:1",
			To:           "0x2.icon/hx1493794ba31fa3372bf7903f04030497e7d14800",            // user icon address
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su", // user bitcoin address
			Amount:       new(big.Int).SetUint64(15000).Bytes(),
			Data:         []byte(""),
		},
		ChainId:  3,
		Receiver: "cxfc86ee7687e1bf681b5548b2667844485c0e7192", // asset manager
		Connectors: []string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5", // connector contract
		},
	}

	inputs := []*multisig.Input{
		{
			TxHash:       "af88dfc8907cc9867e83c05ec061472cdff8cc9b9f951282a5e0f6a0b0b7141b",
			OutputIdx:    1,
			OutputAmount: 909000,
			PkScript:     userMultisigWallet.PKScript,
		},
	}

	// create tx
	msgTx, err := multisig.CreateBridgeTxSendBitcoin(
		&bridgeMsg,
		inputs,
		userMultisigWallet.PKScript,
		relayerPkScript,
		TX_FEE,
	)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	// user key 1 sign tx
	userSigs1, _ := multisig.SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs1)
	// user key 2 sign tx
	userSigs2, _ := multisig.SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs2)
	// COMBINE SIGN
	signedMsgTx, _ := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	bridgeMessage, err := multisig.ReadBridgeMessage(signedMsgTx)
	fmt.Println("bridgeMessage: ", bridgeMessage)
	fmt.Println("err: ", err)
	txSearchRes := &TxSearchRes{
		Tx:            signedMsgTx,
		Height:        3181075,
		TxIndex:       21,
		BridgeMessage: bridgeMessage,
	}
	relayerMessage, err := provider.parseMessageFromTx(txSearchRes)
	fmt.Println("relayerMessage Src: ", relayerMessage.Src)
	fmt.Println("relayerMessage Sn: ", relayerMessage.Sn)
	fmt.Println("relayerMessage Data: ", relayerMessage.Data)
	fmt.Println("err: ", err)

	msg := &types.RecvMessage{
		SrcNID: relayerMessage.Src,
		ConnSn: types.NewHexInt(relayerMessage.Sn.Int64()),
		Msg:    types.NewHexBytes(relayerMessage.Data),
	}
	iconProvider := &icon.Provider{}

	iconMessage := iconProvider.NewIconMessage(iconProvider.GetAddressByEventType(relayerMessage.EventType), msg, icon.MethodRecvMessage)
	fmt.Println("iconMessage Method: ", iconMessage.Method)
	fmt.Println("iconMessage Params: ", iconMessage.Params)
}

func TestParseRuneMessageFromTx(t *testing.T) {
	btcProvider, termpDir := initBtcProvider()
	defer os.Remove(termpDir)
	// relayer multisig
	decodedAddr, _ := btcutil.DecodeAddress(RELAYER_MULTISIG_ADDRESS, btcProvider.chainParam)
	relayerPkScript, _ := txscript.PayToAddrScript(decodedAddr)
	// user key
	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, btcProvider.chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	bridgeMsg := multisig.BridgeDecodedMsg{
		Message: &multisig.XCallMessage{
			MessageType:  1,
			Action:       "Deposit",
			TokenAddress: "2904354:3119",
			To:           "0x2.icon/hx1493794ba31fa3372bf7903f04030497e7d14800",
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su",
			Amount:       new(big.Int).SetUint64(200000).Bytes(),
			Data:         []byte(""),
		},
		ChainId:  1,
		Receiver: "cxfc86ee7687e1bf681b5548b2667844485c0e7192",
		Connectors: []string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	}

	inputs := []*multisig.Input{
		// user rune UTXOs to spend
		{
			TxHash:       "69deba39f5a0700cc713f67fe8cb5ed1e35a9f0d4a3a437d839103c6e26cb947",
			OutputIdx:    2,
			OutputAmount: 546,
			PkScript:     userMultisigWallet.PKScript,
		},
		// user bitcoin UTXOs to pay tx fee
		{
			TxHash:       "073af76df4d7503446ccf704469926ba584ed6a85df1ffb702d3a148a083bd87",
			OutputIdx:    2,
			OutputAmount: 2874285,
			PkScript:     userMultisigWallet.PKScript,
		},
	}

	// create tx
	msgTx, err := multisig.CreateBridgeTxSendRune(
		&bridgeMsg,
		inputs,
		userMultisigWallet.PKScript,
		relayerPkScript,
		TX_FEE,
	)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	// user key 1 sign tx
	userSigs1, _ := multisig.SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs1)
	// user key 2 sign tx
	userSigs2, _ := multisig.SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs2)
	// COMBINE SIGN
	signedMsgTx, _ := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

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

	bridgeMessage, err := multisig.ReadBridgeMessage(signedMsgTx)
	fmt.Println("bridgeMessage: ", bridgeMessage)
	fmt.Println("err: ", err)
	txSearchRes := &TxSearchRes{
		Tx:            signedMsgTx,
		Height:        3181071,
		TxIndex:       21,
		BridgeMessage: bridgeMessage,
	}
	relayerMessage, err := btcProvider.parseMessageFromTx(txSearchRes)
	fmt.Println("relayerMessage Src: ", relayerMessage)
	// fmt.Println("relayerMessage Sn: ", relayerMessage.Sn)
	// fmt.Println("relayerMessage Data: ", relayerMessage.Data)
	fmt.Println("err: ", err)

	msg := &types.RecvMessage{
		SrcNID: relayerMessage.Src,
		ConnSn: types.NewHexInt(relayerMessage.Sn.Int64()),
		Msg:    types.NewHexBytes(relayerMessage.Data),
	}
	iconProvider := &icon.Provider{}

	iconMessage := iconProvider.NewIconMessage(iconProvider.GetAddressByEventType(relayerMessage.EventType), msg, icon.MethodRecvMessage)
	fmt.Println("iconMessage Method: ", iconMessage.Method)
	fmt.Println("iconMessage Params: ", iconMessage.Params)

}

// parse message from icon and build withdraw btc tx
func TestBuildWithdrawBtcTxMessage(t *testing.T) {
	btcProvider, termpDir := initBtcProvider()
	defer os.Remove(termpDir)

	data, _ := base64.StdEncoding.DecodeString("+QEfAbkBG/kBGLMweDIuaWNvbi9jeGZjODZlZTc2ODdlMWJmNjgxYjU1NDhiMjY2Nzg0NDQ4NWMwZTcxOTK4PnRiMXBmMGF0cHQyZDN6ZWw2dWR3czM4cGtyaDJlNDl2cWQzYzVqY3VkM2E4MnNycGhubXBlNTVxMGVjcnprgi5kAbhU+FKKV2l0aGRyYXdUb4MwOjG4PnRiMXBnZTh0aHUzdTBreXF3NXY1dmxoamd0eWR6MzJtbWtkeGRnanRsOTlqcjVmNTlxczB5YXhzNTZ3a3l6gicQ+Ei4RjB4Mi5idGMvdGIxcGd6eDg4MHlmcjdxOGRnejhkcWh3NTBzbmN1NGY0aG13NWNuMzgwMDM1NHR1emN5OWp4NXNodnY3c3U=")
	message := &relayTypes.Message{
		Src:           "0x2.icon",
		Dst:           "0x2.btc",
		Sn:            big.NewInt(11),
		MessageHeight: 1000000,
		EventType:     events.EmitMessage,
		Data:          data,
	}

	msWallet, err := btcProvider.buildMultisigWallet()
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	inputs, tx, err := btcProvider.buildTxMessage(message, 10, msWallet)
	fmt.Println("err: ", err)
	fmt.Println("inputs: ", inputs)
	fmt.Println("tx: ", tx)
	var signedTx bytes.Buffer
	tx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	fmt.Println("hexSignedTx: ", hexSignedTx)
}

// parse message from icon and build withdraw rune tx
func TestBuildWithdrawRunesTxMessage(t *testing.T) {
	btcProvider, termpDir := initBtcProvider()
	defer os.Remove(termpDir)
	data, _ := base64.StdEncoding.DecodeString("+QEpAbkBJfkBIrMweDIuaWNvbi9jeGZjODZlZTc2ODdlMWJmNjgxYjU1NDhiMjY2Nzg0NDQ4NWMwZTcxOTK4PnRiMXBmMGF0cHQyZDN6ZWw2dWR3czM4cGtyaDJlNDl2cWQzYzVqY3VkM2E4MnNycGhubXBlNTVxMGVjcnprgi5uAbhe+FyKV2l0aGRyYXdUb4wyOTA0MzU0OjMxMTm4PnRiMXBneng8ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1")
	message := &relayTypes.Message{
		Src:           "0x2.icon",
		Dst:           "0x2.btc",
		Sn:            big.NewInt(11),
		MessageHeight: 1000000,
		EventType:     events.EmitMessage,
		Data:          data,
	}

	msWallet, err := btcProvider.buildMultisigWallet()
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	inputs, tx, err := btcProvider.buildTxMessage(message, 10, msWallet)
	fmt.Println("err: ", err)
	fmt.Println("inputs: ", inputs)
	fmt.Println("tx: ", tx)
	var signedTx bytes.Buffer
	tx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	fmt.Println("hexSignedTx: ", hexSignedTx)
}

func TestBuildRollbackBtcTxMessage(t *testing.T) {
	// Create a mock Provider
	btcProvider, termpDir := initBtcProvider()
	defer os.Remove(termpDir)
	// Create a mock TxSearchRes
	// relayer multisig
	decodedAddr, _ := btcutil.DecodeAddress(RELAYER_MULTISIG_ADDRESS, btcProvider.chainParam)
	relayerPkScript, _ := txscript.PayToAddrScript(decodedAddr)
	// user key
	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, btcProvider.chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	bridgeMsg := multisig.BridgeDecodedMsg{
		Message: &multisig.XCallMessage{
			MessageType:  0,
			Action:       MethodDeposit,
			TokenAddress: "0:1",
			To:           "0x2.icon/hx1493794ba31fa3372bf7903f04030497e7d14800",            // user icon address
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su", // user bitcoin address
			Amount:       new(big.Int).SetUint64(15000).Bytes(),
			Data:         []byte(""),
		},
		ChainId:  3,
		Receiver: "cxfc86ee7687e1bf681b5548b2667844485c0e7192", // asset manager
		Connectors: []string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5", // connector contract
		},
	}

	inputs := []*multisig.Input{
		{
			TxHash:       "0c95bf20150345ae6e726c986edb8c3019c325ea1049490755191085e3c453fa",
			OutputIdx:    6,
			OutputAmount: 2607546,
			PkScript:     userMultisigWallet.PKScript,
		},
	}

	// create tx
	msgTx, err := multisig.CreateBridgeTxSendBitcoin(
		&bridgeMsg,
		inputs,
		userMultisigWallet.PKScript,
		relayerPkScript,
		TX_FEE,
	)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	// user key 1 sign tx
	userSigs1, _ := multisig.SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs1)
	// user key 2 sign tx
	userSigs2, _ := multisig.SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs2)
	// COMBINE SIGN
	signedMsgTx, _ := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	bridgeMessage, err := multisig.ReadBridgeMessage(signedMsgTx)
	fmt.Println("bridgeMessage: ", bridgeMessage)
	fmt.Println("err: ", err)
	txSearchRes := &TxSearchRes{
		Tx:            signedMsgTx,
		Height:        3181075,
		TxIndex:       21,
		BridgeMessage: bridgeMessage,
	}
	relayerMessage, err := btcProvider.parseMessageFromTx(txSearchRes)
	if err != nil {
		fmt.Println("error: ", err)
	}
	inputs, msWallet, msgTx, relayerSigs, err := btcProvider.HandleBitcoinMessageTx(relayerMessage)

	if err != nil {
		btcProvider.logger.Error("err combine tx: ", zap.Error(err))
	}
	totalSigs = [][][]byte{relayerSigs}
	// send unsigned raw tx and message sn to 2 slave relayers to get sign
	// rsi := slaveRequestParams{
	// 	MsgSn: relayerMessage.Sn.String(),
	// }

	// slaveRequestData, _ := json.Marshal(rsi)
	// slaveSigs := btcProvider.CallSlaves(slaveRequestData)

	// totalSigs = append(totalSigs, slaveSigs...)
	// combine sigs
	signedMsgTx, err = multisig.CombineTapMultisig(totalSigs, msgTx, inputs, msWallet, 0)

	if err != nil {
		btcProvider.logger.Error("err combine tx: ", zap.Error(err))
	}
	btcProvider.logger.Info("signedMsgTx", zap.Any("transaction", signedMsgTx))
	var buf bytes.Buffer
	err = signedMsgTx.Serialize(&buf)

	if err != nil {
		btcProvider.logger.Error("error", zap.Error(err))
	}

	txSize := len(buf.Bytes())
	btcProvider.logger.Info("--------------------txSize--------------------", zap.Int("size", txSize))
	signedMsgTxHex := hex.EncodeToString(buf.Bytes())
	btcProvider.logger.Info("signedMsgTxHex", zap.String("transaction_hex", signedMsgTxHex))

	txHash, err := btcProvider.client.SendRawTransactionV2(btcProvider.cfg.RPCUrl, signedMsgTxHex)
	fmt.Println("txHash: ", txHash)
	fmt.Println("err: ", err)
}

func TestUint64ToBytes(t *testing.T) {
	amount := uint64(15010)
	bytesAmount := []byte(fmt.Sprintf("%d", amount))
	fmt.Println("bytesAmount: ", bytesAmount)
	bigIntAmount := big.NewInt(0).SetBytes(bytesAmount)
	fmt.Println("bigIntAmount: ", bigIntAmount.Uint64(), int64(bigIntAmount.Uint64()))
}
