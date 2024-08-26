package multisig

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/bxelab/runestone"
	"github.com/holiman/uint256"
)

const(
	TX_FEE = 10000
	RELAYER_MULTISIG_ADDRESS = "tb1pv5j5j0dmq2c8d0vnehrlsgrwr9g95m849dl5v0tal8chfdgzqxfskv0w8u"
	USER_MULTISIG_ADDRESS = "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su"
)

func TestGenerateKeys(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	for i := 0; i < 3; i++ {
		privKey := GeneratePrivateKeyFromSeed([]byte{byte(i)}, chainParam)
		wif, _ := btcutil.NewWIF(privKey, chainParam, true)
		pubKey := wif.SerializePubKey();
		witnessProg := btcutil.Hash160(pubKey)
		p2wpkh, _ := btcutil.NewAddressWitnessPubKeyHash(witnessProg, chainParam)

		fmt.Printf("Account %v:\n Private Key: %v\n Public Key: %v\n Address: %v\n", i, wif.String(), hex.EncodeToString(pubKey), p2wpkh)
	}
}

func TestLoadWalletFromPrivateKey(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	wif, _ := btcutil.DecodeWIF("cTYRscQxVhtsGjHeV59RHQJbzNnJHbf3FX4eyX5JkpDhqKdhtRvy")
	pubKey := wif.SerializePubKey();
	witnessProg := btcutil.Hash160(pubKey)
	p2wpkh, _ := btcutil.NewAddressWitnessPubKeyHash(witnessProg, chainParam)

	fmt.Printf("Account:\n Private Key: %v\n Public Key: %v\n Address: %v\n", string(wif.String()), hex.EncodeToString(pubKey), p2wpkh)
}

func TestRandomKeys(t *testing.T) {
	randomKeys(3, &chaincfg.TestNet3Params, []int{0, 1, 2})
}

func TestBuildMultisigTapScript(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	_, relayersMultisigInfo := randomMultisigInfo(3, 3, chainParam, []int{0, 1, 2}, 0, 1)
	relayersMultisigWallet, _ := BuildMultisigWallet(relayersMultisigInfo)
	_, userMultisigInfo := randomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	relayersMultisigAddress, err := AddressOnChain(chainParam, relayersMultisigWallet)
	fmt.Println("relayersMultisigAddress, err : ", relayersMultisigAddress, err)

	userMultisigAddress, err := AddressOnChain(chainParam, userMultisigWallet)
	fmt.Println("userMultisigAddress, err : ", userMultisigAddress, err)
}

func TestMultisigUserClaimLiquidity(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: true,
			TxHash:        "9ed822adb7c3623fcc6776bc93dadb030bf3b887e36975521d540c2a49510e27",
			OutputIdx:     1,
			OutputAmount:  3901,
		},
	}

	outputs := []*OutputTx{
		{
			ReceiverAddress: "tb1pfhttx6vvskhgv6h0w9rss3k63r0zy8vnmwrap0jvraqm5wme6vtsglfta8",
			Amount:          1000,
		},
	}

	relayerPrivKeys, relayersMultisigInfo := randomMultisigInfo(3, 3, chainParam, []int{0, 1, 2}, 0, 1)
	relayersMultisigWallet, _ := BuildMultisigWallet(relayersMultisigInfo)

	changeReceiverAddress := "tb1py04eh93ae0e6dpps2ufxt58wjnvesj0ffzddcckmru3tyrhzsslsxyhwtd"
	msgTx, hexRawTx, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, relayersMultisigWallet, &MultisigWallet{}, chainParam, changeReceiverAddress, 0)
	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: relayersMultisigWallet.PKScript,
		RelayersTapLeaf: relayersMultisigWallet.TapLeaves[0],
		UserPKScript: []byte{},
		UserTapLeaf: txscript.TapLeaf{},
	}

	totalSigs := [][][]byte{}
	// MATSTER RELAYER SIGN TX
	sigs, err := PartSignOnRawExternalTx(relayerPrivKeys[0], msgTx, inputs, tapSigParams, chainParam, true)
	if err != nil {
		fmt.Println("err sign: ", err)
	}
	totalSigs = append(totalSigs, sigs)

	router := SetUpRouter()
	// create post body using an instance of the requestSignInput struct
	rsi := requestSignInput{
		MsgTx:	hexRawTx,
		UTXOs:	inputs,
		TapSigInfo:	tapSigParams,
	}
	requestJson, _ := json.Marshal(rsi)

	// SLAVE RELAYER 1 SIGN TX
	sigs1 := requestSign("/requestSign1", requestJson, router)
	fmt.Println("resp: ", sigs1)
	totalSigs = append(totalSigs, sigs1)

	// SLAVE RELAYER 2 SIGN TX
	sigs2 := requestSign("/requestSign2", requestJson, router)
	fmt.Println("resp: ", sigs2)
	totalSigs = append(totalSigs, sigs2)

	// MATSTER RELAYER COMBINE SIGNS
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, relayersMultisigWallet, 0, relayersMultisigWallet, 0, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)
}

func TestMultisigUserSwap(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: false,
			TxHash:        "374702601b446e0a5c247d15bc6ea049a1266c29ef119ab03801d490ad223bd2",
			OutputIdx:     0,
			OutputAmount:  1501,
		},
		{
			IsRelayersMultisig: true,
			TxHash:        "374702601b446e0a5c247d15bc6ea049a1266c29ef119ab03801d490ad223bd2",
			OutputIdx:     1,
			OutputAmount:  4234,
		},
	}

	outputs := []*OutputTx{
		{
			ReceiverAddress: "tb1pfhttx6vvskhgv6h0w9rss3k63r0zy8vnmwrap0jvraqm5wme6vtsglfta8",
			Amount:          1834,
		},
	}

	relayerPrivKeys, relayersMultisigInfo := randomMultisigInfo(3, 3, chainParam, []int{0, 1, 2}, 0, 1)
	relayersMultisigWallet, _ := BuildMultisigWallet(relayersMultisigInfo)
	userPrivKeys, userMultisigInfo := randomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := "tb1py04eh93ae0e6dpps2ufxt58wjnvesj0ffzddcckmru3tyrhzsslsxyhwtd"
	msgTx, hexRawTx, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, relayersMultisigWallet, userMultisigWallet, chainParam, changeReceiverAddress, 0)
	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: relayersMultisigWallet.PKScript,
		RelayersTapLeaf: relayersMultisigWallet.TapLeaves[0],
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[0],
	}

	totalSigs := [][][]byte{}
	// MATSTER RELAYER SIGN TX
	sigs, err := PartSignOnRawExternalTx(relayerPrivKeys[0], msgTx, inputs, tapSigParams, chainParam, true)
	if err != nil {
		fmt.Println("err sign: ", err)
	}
	totalSigs = append(totalSigs, sigs)

	router := SetUpRouter()
	// create post body using an instance of the requestSignInput struct
	rsi := requestSignInput{
		MsgTx:	hexRawTx,
		UTXOs:	inputs,
		TapSigInfo:	tapSigParams,
	}
	requestJson, _ := json.Marshal(rsi)

	// SLAVE RELAYER 1 SIGN TX
	sigs1 := requestSign("/requestSign1", requestJson, router)
	fmt.Println("resp: ", sigs1)
	totalSigs = append(totalSigs, sigs1)

	// SLAVE RELAYER 2 SIGN TX
	sigs2 := requestSign("/requestSign2", requestJson, router)
	fmt.Println("resp: ", sigs2)
	totalSigs = append(totalSigs, sigs2)

	// USER SIGN TX
	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)

	// add user sign to total sigs
	for i := range msgTx.TxIn {
		if (!inputs[i].IsRelayersMultisig) {
			totalSigs[1][i] = userSigs[i]
		}
	}
	fmt.Println("--------totalSig: ", totalSigs)

	// MATSTER RELAYER COMBINE SIGNS
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, relayersMultisigWallet, 0, userMultisigWallet, 0, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)
}

func TestGenSharedInternalPubKey(t *testing.T) {
	b := make([]byte, 32)
	rand.Read(b)
	bHex := hex.EncodeToString(b)
	fmt.Printf("bHex: %v\n", bHex)
	sharedRandom := new(big.Int).SetBytes(b)
	genSharedInternalPubKey(sharedRandom)
}

func TestParseTx(t *testing.T) {
	hexSignedTx := "01000000000101bcbbb24bd5953d424debb9a24c8009298771eecd3ac0d3c4b219d906a319dfa80000000000e803000002e803000000000000225120d5254f2c52e2672daea941a86c99232693149fd0423ef523fe4e0dcb12a68d53401f000000000000225120d5254f2c52e2672daea941a86c99232693149fd0423ef523fe4e0dcb12a68d530540f4085e4f85eb81b8bd6afd77f728ea75716108cb29cd02aa031def6be65e97e98db40430554669d7b64476d76fd9ae6646529b7abfeee1ac4ad67de0bce9608040f3fc057a9ad0e4a0132040826e2c8e3ca0678ebd515146b8825f527f31195e1966d8424cdf9e963b7335178cab820534e1bd4ede4e8addf47c1bc449a764cec400962c7b22626173655661756c7441646472657373223a22222c22726563656976657241646472657373223a22227d7520fe44ec9f26b97ed30bd33898cf22de726e05389bde632d3aa6ad6746e15221d2ac2030edd881db1bc32b94f83ea5799c2e959854e0f99427d07c211206abd876d052ba201e83d56728fde393b41b74f2b859381661025f2ecec567cf392da7372de47833ba529c21c0636e6671d0135074f83177c5e456191043de9bd54744423b88d6b1ab4751650f00000000"
	msgTx, err := ParseTx(hexSignedTx)
	if err != nil {
		fmt.Printf("Err parse tx: %v", err)
		return
	}

	for _, txIn := range msgTx.TxIn {
		fmt.Printf("txIn: %+v\n ", txIn)
	}

	for _, txOut := range msgTx.TxOut {
		fmt.Printf("txOut: %+v\n ", txOut)
	}
}

func TestUserRecoveryTimeLock(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: false,
			TxHash:        "ae9f43a77d861d5076ebdb1af0d76af033843b784766a1d07a78a68fe845c012",
			OutputIdx:     1,
			OutputAmount:  3808,
		},
	}

	outputs := []*OutputTx{
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
	}

	userPrivKeys, userMultisigInfo := randomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := USER_MULTISIG_ADDRESS
	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)
	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
	}

	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)

	// COMBINE SIGNS
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)
}

func TestTransferRune(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: false,
			TxHash:        "ae9f43a77d861d5076ebdb1af0d76af033843b784766a1d07a78a68fe845c012",
			OutputIdx:     1,
			OutputAmount:  3808,
		},
	}

	outputs := []*OutputTx{
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
	}

	userPrivKeys, userMultisigInfo := randomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := USER_MULTISIG_ADDRESS
	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)
	// Add Rune transfering
	// rune id 840000:3, amount 10000 (5 decimals), to output id 0
	script1, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(1000000000), 0)
	msgTx.AddTxOut(wire.NewTxOut(0, script1))

	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
	}

	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)

	// COMBINE SIGNS
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

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
}

func TestTransferBitcoinWithBridgeMessage(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: false,
			TxHash:        "88cce50b4c9ff0d1c3f58f769fdb8b2e2bfbd3a637664c78b1a38c0f1e492b8f",
			OutputIdx:     0,
			OutputAmount:  2000000,
		},
	}

	outputs := []*OutputTx{
		{
			ReceiverAddress: USER_MULTISIG_ADDRESS,
			Amount:          10000,
		},
	}

	// Add Bridge Message
	payload, _ := CreateBridgePayload(
		1,
		"18403e24f5f5f7cf7a7bae8885d3e1d27ce5fd10",
		&XCallMessage{
			Action:       "Deposit",
			TokenAddress: "0:0",                                                            // bitcoin address
			To:           "0x2.icon/hx39eddef484f6bb08072c59cc0a206e6713fc6d7d",            //(receiver)
			From:         USER_MULTISIG_ADDRESS, //user
			Amount:       new(big.Int).SetUint64(1000).Bytes(),
			Data:         []byte(""),
		})
	scripts, _ := CreateBridgeMessageScripts(payload, 76)
	for i, script := range scripts {
		fmt.Println("OP_RETURN ", i, " script ", script)
		outputs = append(outputs, &OutputTx{
			OpReturnScript: script,
		})
	}

	userPrivKeys, userMultisigInfo := randomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := USER_MULTISIG_ADDRESS
	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)


	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
	}

	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)

	// COMBINE SIGNS
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// Decode bridge message
	decodedMessage, err := ReadBridgeMessage(signedMsgTx)
	fmt.Println("err decode: ", err)
	fmt.Println("-----------decodedMessage---------- ")
	fmt.Println("ChainId: ", decodedMessage.ChainId)
	fmt.Println("Address: ", decodedMessage.Address)
	fmt.Println("Message.Action: ", decodedMessage.Message.Action)
	fmt.Println("Message.TokenAddress: ", decodedMessage.Message.TokenAddress)
	fmt.Println("Message.From: ", decodedMessage.Message.From)
	fmt.Println("Message.To: ", decodedMessage.Message.To)
	fmt.Println("Message.Amount: ", decodedMessage.Message.Amount)
	fmt.Println("Message.Data: ", decodedMessage.Message.Data)
}

func TestRadFiInitPool(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: false,
			TxHash:        "d316231a8aa1f74472ed9cc0f1ed0e36b9b290254cf6b2c377f0d92b299868bf",
			OutputIdx:     4,
			OutputAmount:  1929000,
		},
	}

	// OP_RETURN RadFi Init pool Message
	radfiMsg := RadFiProvideLiquidityMsg {
		Detail: &RadFiProvideLiquidityDetail{
			Fee:		30,
			UpperTick:	12345,
			LowerTick: 	-12345,
			Min0:		0,
			Min1:		10000,
		},
		InitPrice:  uint256.MustFromDecimal("123456789"),
	}
	script, _ := CreateProvideLiquidityScript(&radfiMsg)

	outputs := []*OutputTx{
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
		{
			OpReturnScript: script,
		},
	}

	userPrivKeys, userMultisigInfo := randomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := USER_MULTISIG_ADDRESS
	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
	}
	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)

	// COMBINE SIGN
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// Decode Radfi message
	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)
	fmt.Println("err decode: ", err)
	fmt.Println("decoded message - Flag     : ", decodedRadFiMessage.Flag)
	fmt.Println("decoded message - Fee      : ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.Fee)
	fmt.Println("decoded message - UpperTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.UpperTick)
	fmt.Println("decoded message - LowerTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.LowerTick)
	fmt.Println("decoded message - Min0     : ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.Min0)
	fmt.Println("decoded message - Min1     : ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.Min1)
	fmt.Println("decoded message - InitPrice: ", decodedRadFiMessage.ProvideLiquidityMsg.InitPrice)
}

func TestRadFiProvideLiquidity(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: true,
			TxHash:        "647a499a394bdb2a477f29b9f0515ed186e57a469a732be362a172cde4ea67a5",
			OutputIdx:     0,
			OutputAmount:  1000,
		},
		{
			IsRelayersMultisig: false,
			TxHash:        "647a499a394bdb2a477f29b9f0515ed186e57a469a732be362a172cde4ea67a5",
			OutputIdx:     2,
			OutputAmount:  1918000,
		},
		{
			IsRelayersMultisig: false,
			TxHash:        "96f43d22e98f9abf78422b27bd31412912b7c5e3a2b3f706d61148fc8d0f6550",
			OutputIdx:     0,
			OutputAmount:  1000,
		},
	}

	// OP_RETURN RadFi Provive Liquidity Message
	radfiMsg := RadFiProvideLiquidityMsg {
		Detail: &RadFiProvideLiquidityDetail{
			Fee:		30,
			UpperTick:	12345,
			LowerTick: 	-12345,
			Min0:		0,
			Min1:		10000,
		},
	}
	script1, _ := CreateProvideLiquidityScript(&radfiMsg)
	// OP_RETURN Rune transfering
	// rune id 840000:3, amount 10000 (5 decimals), to output id 4
	script2, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(1000000000), 4)

	outputs := []*OutputTx{
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
		{
			OpReturnScript: script1,
		},
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          20000,
		},
		{
			OpReturnScript: script2,
		},
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
	}

	userPrivKeys, userMultisigInfo := randomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := USER_MULTISIG_ADDRESS
	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
	}
	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)

	// COMBINE SIGNS
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// Decode Radfi message
	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)
	fmt.Println("err decode: ", err)
	fmt.Println("decoded message - Flag     : ", decodedRadFiMessage.Flag)
	fmt.Println("decoded message - Fee      : ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.Fee)
	fmt.Println("decoded message - UpperTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.UpperTick)
	fmt.Println("decoded message - LowerTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.LowerTick)
	fmt.Println("decoded message - Min0     : ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.Min0)
	fmt.Println("decoded message - Min1     : ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.Min1)
	fmt.Println("decoded message - InitPrice: ", decodedRadFiMessage.ProvideLiquidityMsg.InitPrice)

	// get list of Relayer bitcoin an rune UTXOs
	relayerScriptAddress := signedMsgTx.TxOut[0].PkScript
	sequenceNumberUTXO, bitcoinUTXOs, runeUTXOs, err := GetRelayerReceivedUTXO(signedMsgTx, 0, relayerScriptAddress)
	fmt.Println("-------------GetRelayerReceivedUTXO:")
	fmt.Println("relayerScriptAddress: ", relayerScriptAddress)
	fmt.Println("err: ", err)
	fmt.Println("sequenceNumberUTXO: ", sequenceNumberUTXO.IsRelayersMultisig, sequenceNumberUTXO.OutputIdx, sequenceNumberUTXO.OutputAmount, sequenceNumberUTXO.TxHash)
	fmt.Println("bitcoinUTXOs: ")
	for _, utxo := range bitcoinUTXOs {
		fmt.Println("bitcoinUTXO: ", utxo.IsRelayersMultisig, utxo.OutputIdx, utxo.OutputAmount, utxo.TxHash)
	}
	fmt.Println("runeUTXOs: ")
	for _, utxo := range runeUTXOs {
		fmt.Println("runeUTXO: ", utxo.edict, utxo.edictUTXO.IsRelayersMultisig, utxo.edictUTXO.OutputIdx, utxo.edictUTXO.OutputAmount, utxo.edictUTXO.TxHash)
	}
}

func TestRadFiSwap(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: true,
			TxHash:        "4e2b225d732b4123443e8d737e292921ea4a891ddf3a5c9842a4409f2c82805d",
			OutputIdx:     0,
			OutputAmount:  1000,
		},
		{
			IsRelayersMultisig: false,
			TxHash:        "4e2b225d732b4123443e8d737e292921ea4a891ddf3a5c9842a4409f2c82805d",
			OutputIdx:     5,
			OutputAmount:  1888000,
		},
		{
			IsRelayersMultisig: false,
			TxHash:        "4e2b225d732b4123443e8d737e292921ea4a891ddf3a5c9842a4409f2c82805d",
			OutputIdx:     4,
			OutputAmount:  1000,
		},
	}

	// OP_RETURN RadFi Swap Message
	radfiMsg := RadFiSwapMsg {
		IsExactInOut: true,
		TokenOutIndex: 2,
	}
	script1, _ := CreateSwapScript(&radfiMsg)
	// OP_RETURN Rune transfering
	// rune id 840000:3, amount 4000 (5 decimals), to output id 3
	script2, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(400000000), 3)
	// rune id 840000:3, amount 6000 (5 decimals), to output id 6
	script3, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(600000000), 6)

	outputs := []*OutputTx{
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
		{
			OpReturnScript: script1,
		},
		{
			OpReturnScript: script2,
		},
		{
			ReceiverAddress: USER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          10000,
		},
		{
			OpReturnScript: script3,
		},
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
	}

	userPrivKeys, userMultisigInfo := randomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := USER_MULTISIG_ADDRESS
	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
	}
	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)

	// COMBINE SIGNS
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// Decode Radfi message
	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)

	fmt.Println("err decode: ", err)
	fmt.Println("decoded message - Flag                  : ", decodedRadFiMessage.Flag)
	fmt.Println("decoded message - IsExactInOut          : ", decodedRadFiMessage.SwapMsg.IsExactInOut)
	fmt.Println("decoded message - TokenOutIndex         : ", decodedRadFiMessage.SwapMsg.TokenOutIndex)
}

func TestRadFiWithdrawLiquidity(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: true,
			TxHash:        "3aa9c4b9a71fe19560c467cdddce932eae8a10e28123a09acb27701123f2a8f7",
			OutputIdx:     0,
			OutputAmount:  1000,
		},
		{
			IsRelayersMultisig: true,
			TxHash:        "3aa9c4b9a71fe19560c467cdddce932eae8a10e28123a09acb27701123f2a8f7",
			OutputIdx:     4,
			OutputAmount:  10000,
		},
		{
			IsRelayersMultisig: true,
			TxHash:        "3aa9c4b9a71fe19560c467cdddce932eae8a10e28123a09acb27701123f2a8f7",
			OutputIdx:     6,
			OutputAmount:  1000,
		},
		{
			IsRelayersMultisig: false,
			TxHash:        "4e2b225d732b4123443e8d737e292921ea4a891ddf3a5c9842a4409f2c82805d",
			OutputIdx:     5,
			OutputAmount:  1867000,
		},
	}

	// OP_RETURN RadFi Withdraw Liquidity Message
	radfiMsg := RadFiWithdrawLiquidityMsg {
		RecipientIndex:	2,
		LiquidityValue: uint256.MustFromDecimal("123456"),
		NftId:			uint256.MustFromDecimal("123456789"),
	}
	script1, _ := CreateWithdrawLiquidityScript(&radfiMsg)
	// OP_RETURN Rune transfering
	// rune id 840000:3, amount 2500 (5 decimals), to output id 3
	script2, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(250000000), 4)
	// rune id 840000:3, amount 3500 (5 decimals), to output id 7
	script3, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(350000000), 7)

	outputs := []*OutputTx{
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
		{
			OpReturnScript: script1,
		},
		{
			ReceiverAddress: USER_MULTISIG_ADDRESS,
			Amount:          5000,
		},
		{
			OpReturnScript: script2,
		},
		{
			ReceiverAddress: USER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          5000,
		},
		{
			OpReturnScript: script3,
		},
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
	}
	userPrivKeys, userMultisigInfo := randomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := USER_MULTISIG_ADDRESS
	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
	}
	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)

	// COMBINE SIGNS
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// Decode Radfi message
	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)

	fmt.Println("err decode: ", err)
	fmt.Println("decoded message - Flag           : ", decodedRadFiMessage.Flag)
	fmt.Println("decoded message - RecipientIndex : ", decodedRadFiMessage.WithdrawLiquidityMsg.RecipientIndex)
	fmt.Println("decoded message - LiquidityValue : ", decodedRadFiMessage.WithdrawLiquidityMsg.LiquidityValue)
	fmt.Println("decoded message - NftId          : ", decodedRadFiMessage.WithdrawLiquidityMsg.NftId)
}

func TestRadFiCollectFees(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: true,
			TxHash:        "b508b193dd9712406c8d712af00e6d3e34b8d2ad52809b72c63f6463c1555703",
			OutputIdx:     0,
			OutputAmount:  1000,
		},
		{
			IsRelayersMultisig: true,
			TxHash:        "b508b193dd9712406c8d712af00e6d3e34b8d2ad52809b72c63f6463c1555703",
			OutputIdx:     5,
			OutputAmount:  5000,
		},
		{
			IsRelayersMultisig: true,
			TxHash:        "b508b193dd9712406c8d712af00e6d3e34b8d2ad52809b72c63f6463c1555703",
			OutputIdx:     7,
			OutputAmount:  1000,
		},
		{
			IsRelayersMultisig: false,
			TxHash:        "b508b193dd9712406c8d712af00e6d3e34b8d2ad52809b72c63f6463c1555703",
			OutputIdx:     8,
			OutputAmount:  1856000,
		},
	}

	// OP_RETURN RadFi Collect Fees
	radfiMsg := RadFiCollectFeesMsg {
		RecipientIndex:	2,
		NftId:			uint256.MustFromDecimal("123456789"),
	}
	script1, _ := CreateCollectFeesScript(&radfiMsg)
	// OP_RETURN Rune transfering
	// rune id 840000:3, amount 500 (5 decimals), to output id 4
	script2, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(50000000), 4)
	// rune id 840000:3, amount 3000 (5 decimals), to output id 7
	script3, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(300000000), 7)

	outputs := []*OutputTx{
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
		{
			OpReturnScript: script1,
		},
		{
			ReceiverAddress: USER_MULTISIG_ADDRESS,
			Amount:          2000,
		},
		{
			OpReturnScript: script2,
		},
		{
			ReceiverAddress: USER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          3000,
		},
		{
			OpReturnScript: script3,
		},
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
	}

	userPrivKeys, userMultisigInfo := randomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := USER_MULTISIG_ADDRESS
	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
	}
	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)

	// COMBINE SIGNS
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// Decode Radfi message
	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)

	fmt.Println("err decode: ", err)
	fmt.Println("decoded message - Flag           : ", decodedRadFiMessage.Flag)
	fmt.Println("decoded message - RecipientIndex : ", decodedRadFiMessage.CollectFeesMsg.RecipientIndex)
	fmt.Println("decoded message - NftId          : ", decodedRadFiMessage.CollectFeesMsg.NftId)
}

func TestRadFiIncreaseLiquidity(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: true,
			TxHash:        "546fd5aa7336a82ec31dda54de0ec0e81c9285e42f75768f5de058b390126e51",
			OutputIdx:     0,
			OutputAmount:  1000,
		},
		{
			IsRelayersMultisig: false,
			TxHash:        "546fd5aa7336a82ec31dda54de0ec0e81c9285e42f75768f5de058b390126e51",
			OutputIdx:     8,
			OutputAmount:  1845000,
		},
		{
			IsRelayersMultisig: false,
			TxHash:        "546fd5aa7336a82ec31dda54de0ec0e81c9285e42f75768f5de058b390126e51",
			OutputIdx:     0,
			OutputAmount:  1000,
		},
	}

	// OP_RETURN RadFi Increase Liquidity Message
	radfiMsg := RadFiIncreaseLiquidityMsg {
		Min0:	0,
		Min1:	10000,
		NftId:	uint256.MustFromDecimal("123456789"),
	}
	script1, _ := CreateIncreaseLiquidityScript(&radfiMsg)
	// OP_RETURN Rune transfering
	// rune id 840000:3, amount 500 (5 decimals), to output id 4
	script2, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(50000000), 4)

	outputs := []*OutputTx{
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
		{
			OpReturnScript: script1,
		},
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          20000,
		},
		{
			OpReturnScript: script2,
		},
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
	}

	userPrivKeys, userMultisigInfo := randomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := USER_MULTISIG_ADDRESS
	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
	}
	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)

	// COMBINE SIGNS
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// Decode Radfi message
	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)

	fmt.Println("err decode: ", err)
	fmt.Println("decoded message - Flag           : ", decodedRadFiMessage.Flag)
	fmt.Println("decoded message - Min0           : ", decodedRadFiMessage.IncreaseLiquidityMsg.Min0)
	fmt.Println("decoded message - Min1           : ", decodedRadFiMessage.IncreaseLiquidityMsg.Min1)
	fmt.Println("decoded message - NftId          : ", decodedRadFiMessage.IncreaseLiquidityMsg.NftId)
}
