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
	"github.com/icon-project/centralized-relay/relayer/events"
	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/multisig"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

const (
	// TESTNET
	TESTNET_USER_WALLET_ADDRESS         = "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su"
	TESTNET_RELAYER_ADDRESS             = "tb1pf0atpt2d3zel6udws38pkrh2e49vqd3c5jcud3a82srphnmpe55q0ecrzk"
	TESTNET_ASSET_MANAGER_ADDRESS       = "cx8c9a213cd5dcebfb539c30e1af6d77990b200ca4"
	TESTNET_ASSET_MANAGER_ADDRESS_WRONG = "cx8c9a213cd5dcebfb539c30e1af6d77990b200ca5"
	TESTNET_CONNECTION_ADDRESS          = "cxc7a77b874eddfe3fb5434effaf35375b697496ca"
	TESTNET_TX_FEE                      = 20000
	TESTNET_BTCTOKEN                    = "0:1"
	TESTNET_RUNETOKEN                   = "2904354:3119"
	TESTNET_ICON_RECEIVER_ADDRESS       = "0x2.icon/hx1493794ba31fa3372bf7903f04030497e7d14800"
)

func initBtcProviderTestnet() (*Provider, string) {
	tempDir, _ := os.MkdirTemp("", "bitcoin_provider_test")

	dbPath := filepath.Join(tempDir, "test.db")
	db, _ := leveldb.OpenFile(dbPath, nil)

	logger, _ := zap.NewDevelopment()

	config := &Config{}
	config = &Config{
		UniSatURL: "https://open-api.unisat.io",
		UniSatKey: "YOUR_UNISAT_API_KEY",

		MasterPubKey:     "02fe44ec9f26b97ed30bd33898cf22de726e05389bde632d3aa6ad6746e15221d2",
		Slave1PubKey:     "0230edd881db1bc32b94f83ea5799c2e959854e0f99427d07c211206abd876d052",
		Slave2PubKey:     "021e83d56728fde393b41b74f2b859381661025f2ecec567cf392da7372de47833",
		RelayerPrivKey:   "cTYRscQxVhtsGjHeV59RHQJbzNnJHbf3FX4eyX5JkpDhqKdhtRvy",
		RecoveryLockTime: 1234,
		OpCode:           0x5e,
		MempoolURL:       "https://mempool.space/api/v1",
		SlaveServer1:     "http://localhost:8081",
		SlaveServer2:     "http://localhost:8082",
		Port:             "8080",
		RequestTimeout:   1000,
		ApiKey:           "key",
		Connections:      []string{TESTNET_CONNECTION_ADDRESS},
	}
	config.NID = "0x2.btc"
	config.Address = TESTNET_RELAYER_ADDRESS
	config.RPCUrl = ""
	config.User = "123"
	config.Password = "123"

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

func buildUserMultisigWalletTestnet(chainParam *chaincfg.Params) ([]string, *multisig.MultisigWallet, error) {
	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)
	return userPrivKeys, userMultisigWallet, nil
}

// parse message from icon and build withdraw btc tx
func TestBuildWithdrawBtcTxMessageTestnet(t *testing.T) {
	btcProvider, termpDir := initBtcProviderTestnet()
	defer os.Remove(termpDir)

	// decode message from icon
	data, _ := base64.StdEncoding.DecodeString("+QEfAbkBG/kBGLMweDIuaWNvbi9jeGZjODZlZTc2ODdlMWJmNjgxYjU1NDhiMjY2Nzg0NDQ4NWMwZTcxOTK4PnRiMXBmMGF0cHQyZDN6ZWw2dWR3czM4cGtyaDJlNDl2cWQzYzVqY3VkM2E4MnNycGhubXBlNTVxMGVjcnprgi5kAbhU+FKKV2l0aGRyYXdUb4MwOjG4PnRiMXBnZTh0aHUzdTBreXF3NXY1dmxoamd0eWR6MzJtbWtkeGRnanRsOTlqcjVmNTlxczB5YXhzNTZ3a3l6gicQ+Ei4RjB4Mi5idGMvdGIxcGd6eDg4MHlmcjdxOGRnejhkcWh3NTBzbmN1NGY0aG13NWNuMzgwMDM1NHR1emN5OWp4NXNodnY3c3U=")
	message := &relayTypes.Message{
		Src:           "0x2.icon",
		Dst:           "0x2.btc",
		Sn:            big.NewInt(11),
		MessageHeight: 1000000,
		EventType:     events.EmitMessage,
		Data:          data,
	}
	feeRate := 10

	inputs, msWallet, msgTx, relayerSigs, feeRate, err := btcProvider.HandleBitcoinMessageTx(message, feeRate, []slaveRequestInput{})
	if err != nil {
		btcProvider.logger.Error("failed to handle bitcoin message tx: %v", zap.Error(err))
		// return
	}
	totalSigs := [][][]byte{relayerSigs}
	// send unsigned raw tx and message sn to 2 slave relayers to get sign
	rsi := slaveRequestParams{
		MsgSn:   message.Sn.String(),
		FeeRate: feeRate,
	}

	slaveRequestData, _ := json.Marshal(rsi)
	slaveSigs := btcProvider.CallSlaves(slaveRequestData, "")
	btcProvider.logger.Info("Slave signatures", zap.Any("slave sigs", slaveSigs))
	totalSigs = append(totalSigs, slaveSigs...)
	// combine sigs
	signedMsgTx, err := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, msWallet, 0)

	if err != nil {
		btcProvider.logger.Error("err combine tx: ", zap.Error(err))
	}

	var buf bytes.Buffer
	err = signedMsgTx.Serialize(&buf)

	if err != nil {
		btcProvider.logger.Error("err combine tx: ", zap.Error(err))
	}

	signedMsgTxHex := hex.EncodeToString(buf.Bytes())
	btcProvider.logger.Info("signedMsgTxHex", zap.String("transaction_hex", signedMsgTxHex))

	btcProvider.cacheSpentUTXOs(inputs)

	// Broadcast transaction to bitcoin network
	txHash, err := btcProvider.client.SendRawTransaction(btcProvider.cfg.MempoolURL, []json.RawMessage{json.RawMessage(signedMsgTxHex)})
	if err != nil {
		btcProvider.removeCachedSpentUTXOs(inputs)
		btcProvider.logger.Error("failed to send raw transaction", zap.Error(err))
		// return
	}

	btcProvider.logger.Info("txHash", zap.String("transaction_hash", txHash))
}

// parse message from icon and build withdraw rune tx
func TestBuildWithdrawRunesTxMessageTestnet(t *testing.T) {
	btcProvider, termpDir := initBtcProviderTestnet()
	defer os.Remove(termpDir)
	// decode message from icon
	data, _ := base64.StdEncoding.DecodeString("+QEoAbkBJPkBIbMweDIuaWNvbi9jeDhjOWEyMTNjZDVkY2ViZmI1MzljMzBlMWFmNmQ3Nzk5MGIyMDBjYTS4PnRiMXBmMGF0cHQyZDN6ZWw2dWR3czM4cGtyaDJlNDl2cWQzYzVqY3VkM2E4MnNycGhubXBlNTVxMGVjcnprgi6wArhd+FuKV2l0aGRyYXdUb4wyOTA0MzU0OjMxMTm4PnRiMXB5OHZ4ZTdlY3dqdGRsdmNjcmVrZ3E1dTkwcnpzZzV3djZhamg0YXZjcDRxNXBrdjYyOWRxeDZlZGg5gg+g+Ei4RjB4Mi5idGMvdGIxcGYwYXRwdDJkM3plbDZ1ZHdzMzhwa3JoMmU0OXZxZDNjNWpjdWQzYTgyc3JwaG5tcGU1NXEwZWNyems=")
	message := &relayTypes.Message{
		Src:           "0x2.icon",
		Dst:           "0x2.btc",
		Sn:            big.NewInt(12),
		MessageHeight: 46166601,
		EventType:     events.EmitMessage,
		Data:          data,
	}
	feeRate := 10
	inputs, msWallet, msgTx, relayerSigs, feeRate, err := btcProvider.HandleBitcoinMessageTx(message, feeRate, []slaveRequestInput{})
	if err != nil {
		btcProvider.logger.Error("failed to handle bitcoin message tx: %v", zap.Error(err))
		// return
	}
	totalSigs := [][][]byte{relayerSigs}
	// send unsigned raw tx and message sn to 2 slave relayers to get sign
	rsi := slaveRequestParams{
		MsgSn:   message.Sn.String(),
		FeeRate: feeRate,
	}

	slaveRequestData, _ := json.Marshal(rsi)
	slaveSigs := btcProvider.CallSlaves(slaveRequestData, "")
	btcProvider.logger.Info("Slave signatures", zap.Any("slave sigs", slaveSigs))
	totalSigs = append(totalSigs, slaveSigs...)
	// combine sigs
	signedMsgTx, err := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, msWallet, 0)

	if err != nil {
		btcProvider.logger.Error("err combine tx: ", zap.Error(err))
	}

	var buf bytes.Buffer
	err = signedMsgTx.Serialize(&buf)

	if err != nil {
		btcProvider.logger.Error("err combine tx: ", zap.Error(err))
	}

	signedMsgTxHex := hex.EncodeToString(buf.Bytes())
	btcProvider.logger.Info("signedMsgTxHex", zap.String("transaction_hex", signedMsgTxHex))

	btcProvider.cacheSpentUTXOs(inputs)

	// Broadcast transaction to bitcoin network
	txHash, err := btcProvider.client.SendRawTransaction(btcProvider.cfg.MempoolURL, []json.RawMessage{json.RawMessage(signedMsgTxHex)})
	if err != nil {
		btcProvider.removeCachedSpentUTXOs(inputs)
		btcProvider.logger.Error("failed to send raw transaction", zap.Error(err))
	}

	btcProvider.logger.Info("txHash", zap.String("transaction_hash", txHash))
}

func TestBuildRefundBtcTxMessageTestnet(t *testing.T) {
	// Create a mock Provider
	btcProvider, termpDir := initBtcProviderTestnet()
	defer os.Remove(termpDir)
	// relayer multisig
	decodedAddr, _ := btcutil.DecodeAddress(TESTNET_RELAYER_ADDRESS, btcProvider.chainParam)
	relayerPkScript, _ := txscript.PayToAddrScript(decodedAddr)
	// user key
	userPrivKeys, userMultisigWallet, _ := buildUserMultisigWalletTestnet(btcProvider.chainParam)

	bridgeMsg := multisig.BridgeDecodedMsg{
		Message: &multisig.XCallMessage{
			MessageType:  1,
			Action:       "Deposit",
			TokenAddress: TESTNET_BTCTOKEN,
			To:           TESTNET_ICON_RECEIVER_ADDRESS, // user icon address
			From:         TESTNET_USER_WALLET_ADDRESS,   // user bitcoin address
			Amount:       new(big.Int).SetUint64(5000).Bytes(),
			Data:         []byte(""),
		},
		ChainId:  1,
		Receiver: TESTNET_ASSET_MANAGER_ADDRESS, // asset manager
		Connectors: []string{
			TESTNET_CONNECTION_ADDRESS, // connector contract
		},
	}

	inputs := []*multisig.Input{
		{
			TxHash:       "89095d016b50644a328667cd5543b0f29c0f2a81242094ea7318bded49cf30a8",
			OutputIdx:    8,
			OutputAmount: 202763,
			PkScript:     userMultisigWallet.PKScript,
		},
	}

	// create tx
	msgTx, err := multisig.CreateBridgeTxSendBitcoin(
		&bridgeMsg,
		inputs,
		userMultisigWallet.PKScript,
		relayerPkScript,
		TESTNET_TX_FEE,
	)
	if err != nil {
		fmt.Println("err: ", err)
	}

	// add 10 to output amount to make wrong amount
	msgTx.TxOut[0].Value = msgTx.TxOut[0].Value + 10

	// sign tx
	totalSigs := [][][]byte{}
	for _, privKey := range userPrivKeys {
		// user key 1 sign tx
		userSigs, _ := multisig.SignTapMultisig(privKey, msgTx, inputs, userMultisigWallet, 0)
		totalSigs = append(totalSigs, userSigs)
	}
	// COMBINE SIGN
	signedMsgTx, _ := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	fmt.Println("hexSignedTx: ", hexSignedTx)

	txHash, err := btcProvider.client.SendRawTransaction(btcProvider.cfg.MempoolURL, []json.RawMessage{json.RawMessage(hexSignedTx)})
	fmt.Println("txHash: ", txHash)
	fmt.Println("err: ", err)

}

func TestDepositBitcoinToIconTestnet(t *testing.T) {
	btcProvider, termpDir := initBtcProviderTestnet()
	defer os.Remove(termpDir)
	// relayer multisig
	decodedAddr, _ := btcutil.DecodeAddress(TESTNET_RELAYER_ADDRESS, btcProvider.chainParam)
	relayerPkScript, _ := txscript.PayToAddrScript(decodedAddr)
	// user key
	userPrivKeys, userMultisigWallet, _ := buildUserMultisigWalletTestnet(btcProvider.chainParam)

	bridgeMsg := multisig.BridgeDecodedMsg{
		Message: &multisig.XCallMessage{
			MessageType:  1,
			Action:       "Deposit",
			TokenAddress: TESTNET_BTCTOKEN,
			To:           TESTNET_ICON_RECEIVER_ADDRESS,
			From:         TESTNET_USER_WALLET_ADDRESS,
			Amount:       new(big.Int).SetUint64(1000).Bytes(),
			Data:         []byte(""),
		},
		ChainId:  1,
		Receiver: TESTNET_ASSET_MANAGER_ADDRESS,
		Connectors: []string{
			TESTNET_CONNECTION_ADDRESS,
		},
	}

	inputs := []*multisig.Input{
		{
			TxHash:       "e57619adbb62c2e8add13eb2694010f3ba337cf5007ca3674224c272feffb097",
			OutputIdx:    6,
			OutputAmount: 271530,
			PkScript:     userMultisigWallet.PKScript,
		},
	}

	// create tx
	msgTx, err := multisig.CreateBridgeTxSendBitcoin(
		&bridgeMsg,
		inputs,
		userMultisigWallet.PKScript,
		relayerPkScript,
		TESTNET_TX_FEE,
	)

	// log hex of unsigned tx

	// sign tx
	totalSigs := [][][]byte{}
	for _, privKey := range userPrivKeys {
		// user key 1 sign tx
		userSigs, _ := multisig.SignTapMultisig(privKey, msgTx, inputs, userMultisigWallet, 0)
		totalSigs = append(totalSigs, userSigs)
	}
	// COMBINE SIGN

	signedMsgTx, _ := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	fmt.Println("hexSignedTx: ", hexSignedTx)

	txHash, err := btcProvider.client.SendRawTransaction(btcProvider.cfg.MempoolURL, []json.RawMessage{json.RawMessage(hexSignedTx)})
	fmt.Println("txHash: ", txHash)
	fmt.Println("err: ", err)
}
func TestDepositBitcoinToIconFailTestnet(t *testing.T) {
	btcProvider, termpDir := initBtcProviderTestnet()
	defer os.Remove(termpDir)
	// relayer multisig
	decodedAddr, _ := btcutil.DecodeAddress(TESTNET_RELAYER_ADDRESS, btcProvider.chainParam)
	relayerPkScript, _ := txscript.PayToAddrScript(decodedAddr)
	// user key
	userPrivKeys, userMultisigWallet, _ := buildUserMultisigWalletTestnet(btcProvider.chainParam)

	bridgeMsg := multisig.BridgeDecodedMsg{
		Message: &multisig.XCallMessage{
			MessageType:  1,
			Action:       "Deposit",
			TokenAddress: TESTNET_BTCTOKEN,
			To:           TESTNET_ICON_RECEIVER_ADDRESS,
			From:         TESTNET_USER_WALLET_ADDRESS,
			Amount:       new(big.Int).SetUint64(3000).Bytes(),
			Data:         []byte(""),
		},
		ChainId:  1,
		Receiver: TESTNET_ASSET_MANAGER_ADDRESS_WRONG,
		Connectors: []string{
			TESTNET_CONNECTION_ADDRESS,
		},
	}

	inputs := []*multisig.Input{
		{
			TxHash:       "37ba284f4664c517af1c7fbca8b6ce19b293fe5e4125ac42e35f8b6671968d64",
			OutputIdx:    6,
			OutputAmount: 234779,
			PkScript:     userMultisigWallet.PKScript,
		},
	}

	// create tx
	msgTx, err := multisig.CreateBridgeTxSendBitcoin(
		&bridgeMsg,
		inputs,
		userMultisigWallet.PKScript,
		relayerPkScript,
		TESTNET_TX_FEE,
	)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	for _, privKey := range userPrivKeys {
		// user key 1 sign tx
		userSigs, _ := multisig.SignTapMultisig(privKey, msgTx, inputs, userMultisigWallet, 0)
		totalSigs = append(totalSigs, userSigs)
	}
	// COMBINE SIGN
	signedMsgTx, _ := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	fmt.Println("hexSignedTx: ", hexSignedTx)

	txHash, err := btcProvider.client.SendRawTransaction(btcProvider.cfg.MempoolURL, []json.RawMessage{json.RawMessage(hexSignedTx)})
	fmt.Println("txHash: ", txHash)
	fmt.Println("err: ", err)
}

func TestDepositRuneToIconTestnet(t *testing.T) {
	btcProvider, termpDir := initBtcProviderTestnet()
	defer os.Remove(termpDir)
	// relayer multisig
	decodedAddr, _ := btcutil.DecodeAddress(TESTNET_RELAYER_ADDRESS, btcProvider.chainParam)
	relayerPkScript, _ := txscript.PayToAddrScript(decodedAddr)
	// user key
	userPrivKeys, userMultisigWallet, _ := buildUserMultisigWalletTestnet(btcProvider.chainParam)

	bridgeMsg := multisig.BridgeDecodedMsg{
		Message: &multisig.XCallMessage{
			MessageType:  1,
			Action:       "Deposit",
			TokenAddress: TESTNET_RUNETOKEN,
			To:           TESTNET_ICON_RECEIVER_ADDRESS,
			From:         TESTNET_USER_WALLET_ADDRESS,
			Amount:       new(big.Int).SetUint64(1).Bytes(),
			Data:         []byte(""),
		},
		ChainId:  1,
		Receiver: TESTNET_ASSET_MANAGER_ADDRESS,
		Connectors: []string{
			TESTNET_CONNECTION_ADDRESS,
		},
	}

	inputs := []*multisig.Input{
		// user rune UTXOs to spend
		{
			TxHash:       "e6107730017b2b84e187add5c3fed5a71edcda0b85b982297282137fe0480234",
			OutputIdx:    1,
			OutputAmount: 546,
			PkScript:     userMultisigWallet.PKScript,
		},
		// user bitcoin UTXOs to pay tx fee
		{
			TxHash:       "9a9d955dff45c6cef6f4e41a12052dde21179069a2e17fe8f381f6c75e112b6a",
			OutputIdx:    6,
			OutputAmount: 261795,
			PkScript:     userMultisigWallet.PKScript,
		},
	}

	// create tx
	msgTx, err := multisig.CreateBridgeTxSendRune(
		&bridgeMsg,
		inputs,
		userMultisigWallet.PKScript,
		relayerPkScript,
		TESTNET_TX_FEE,
	)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	for _, privKey := range userPrivKeys {
		// user key 1 sign tx
		userSigs, _ := multisig.SignTapMultisig(privKey, msgTx, inputs, userMultisigWallet, 0)
		totalSigs = append(totalSigs, userSigs)
	}
	// COMBINE SIGN
	signedMsgTx, _ := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	fmt.Println("hexSignedTx: ", hexSignedTx)

	txHash, err := btcProvider.client.SendRawTransaction(btcProvider.cfg.MempoolURL, []json.RawMessage{json.RawMessage(hexSignedTx)})
	fmt.Println("txHash: ", txHash)
	fmt.Println("err: ", err)
}

func TestDepositRuneToIconFailTestnet(t *testing.T) {
	btcProvider, termpDir := initBtcProviderTestnet()
	defer os.Remove(termpDir)
	// relayer multisig
	decodedAddr, _ := btcutil.DecodeAddress(TESTNET_RELAYER_ADDRESS, btcProvider.chainParam)
	relayerPkScript, _ := txscript.PayToAddrScript(decodedAddr)
	// user key
	userPrivKeys, userMultisigWallet, _ := buildUserMultisigWalletTestnet(btcProvider.chainParam)

	bridgeMsg := multisig.BridgeDecodedMsg{
		Message: &multisig.XCallMessage{
			MessageType:  1,
			Action:       "Deposit",
			TokenAddress: TESTNET_RUNETOKEN,
			To:           TESTNET_ICON_RECEIVER_ADDRESS,
			From:         TESTNET_USER_WALLET_ADDRESS,
			Amount:       new(big.Int).SetUint64(1).Bytes(),
			Data:         []byte(""),
		},
		ChainId:  1,
		Receiver: TESTNET_ASSET_MANAGER_ADDRESS_WRONG,
		Connectors: []string{
			TESTNET_CONNECTION_ADDRESS,
		},
	}

	inputs := []*multisig.Input{
		// user rune UTXOs to spend
		{
			TxHash:       "924c7c6bd13f465b0b50cb8ad883544b22bfe54fae42e2ecfc9f9609a1b616f7",
			OutputIdx:    1,
			OutputAmount: 546,
			PkScript:     userMultisigWallet.PKScript,
		},
		// user bitcoin UTXOs to pay tx fee
		{
			TxHash:       "b84060ce292dd61f8490bae54f8354caa8642de730e5f409b72d67b05617dcb0",
			OutputIdx:    6,
			OutputAmount: 226044,
			PkScript:     userMultisigWallet.PKScript,
		},
	}

	// create tx
	msgTx, err := multisig.CreateBridgeTxSendRune(
		&bridgeMsg,
		inputs,
		userMultisigWallet.PKScript,
		relayerPkScript,
		TESTNET_TX_FEE,
	)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	for _, privKey := range userPrivKeys {
		// user key 1 sign tx
		userSigs, _ := multisig.SignTapMultisig(privKey, msgTx, inputs, userMultisigWallet, 0)
		totalSigs = append(totalSigs, userSigs)
	}
	// COMBINE SIGN
	signedMsgTx, _ := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	fmt.Println("hexSignedTx: ", hexSignedTx)

	txHash, err := btcProvider.client.SendRawTransaction(btcProvider.cfg.RPCUrl, []json.RawMessage{json.RawMessage(`"` + hexSignedTx + `"`)})
	fmt.Println("txHash: ", txHash)
	fmt.Println("err: ", err)
}

// Run to build Relayer Multisig Wallet
func TestBuildRelayerMultisigWalletTestnet(t *testing.T) {
	btcProvider, termpDir := initBtcProviderTestnet()
	defer os.Remove(termpDir)
	msWallet, err := btcProvider.buildMultisigWallet()
	fmt.Println("err: ", err)
	fmt.Println("msWallet TapScriptTree: ", msWallet.TapScriptTree)
	fmt.Println("msWallet TapLeaves: ", msWallet.TapLeaves)
	fmt.Println("msWallet PKScript: ", hex.EncodeToString(msWallet.PKScript))
	fmt.Println("msWallet SharedPublicKey: ", msWallet.SharedPublicKey.ToECDSA())
	// log msWallet address
	addr, err := multisig.AddressOnChain(btcProvider.chainParam, msWallet)
	fmt.Println("msWallet Address: ", addr.String())
	fmt.Println("err: ", err)
}

// Run to build User Multisig Wallet
func TestBuildUserMultisigWalletTestnet(t *testing.T) {
	btcProvider, termpDir := initBtcProviderTestnet()
	defer os.Remove(termpDir)
	userPrivKeys, userMultisigWallet, _ := buildUserMultisigWalletTestnet(btcProvider.chainParam)
	fmt.Println("userPrivKeys: ", userPrivKeys)
	fmt.Println("userMultisigWallet: ", userMultisigWallet)
	addr, err := multisig.AddressOnChain(btcProvider.chainParam, userMultisigWallet)
	fmt.Println("userMultisigWallet Address: ", addr.String())
	fmt.Println("err: ", err)
}
