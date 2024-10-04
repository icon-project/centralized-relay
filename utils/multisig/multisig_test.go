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
	"github.com/bxelab/runestone"
)

const (
	TX_FEE                   = 10000
	RELAYER_MULTISIG_ADDRESS = "tb1pv5j5j0dmq2c8d0vnehrlsgrwr9g95m849dl5v0tal8chfdgzqxfskv0w8u"
	USER_MULTISIG_ADDRESS    = "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su"
)

func TestGenerateKeys(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	for i := 0; i < 3; i++ {
		privKey := GeneratePrivateKeyFromSeed([]byte{byte(i)}, chainParam)
		wif, _ := btcutil.NewWIF(privKey, chainParam, true)
		pubKey := wif.SerializePubKey()
		witnessProg := btcutil.Hash160(pubKey)
		p2wpkh, _ := btcutil.NewAddressWitnessPubKeyHash(witnessProg, chainParam)

		fmt.Printf("Account %v:\n Private Key: %v\n Public Key: %v\n Address: %v\n", i, wif.String(), hex.EncodeToString(pubKey), p2wpkh)
	}
}

func TestLoadWalletFromPrivateKey(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	wif, _ := btcutil.DecodeWIF("cTYRscQxVhtsGjHeV59RHQJbzNnJHbf3FX4eyX5JkpDhqKdhtRvy")
	pubKey := wif.SerializePubKey()
	witnessProg := btcutil.Hash160(pubKey)
	p2wpkh, _ := btcutil.NewAddressWitnessPubKeyHash(witnessProg, chainParam)

	fmt.Printf("Account:\n Private Key: %v\n Public Key: %v\n Address: %v\n", string(wif.String()), hex.EncodeToString(pubKey), p2wpkh)
}

func TestRandomKeys(t *testing.T) {
	randomKeys(3, &chaincfg.TestNet3Params, []int{0, 1, 2})
}

func TestBuildMultisigTapScript(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	relayerPrivKeys, relayersMultisigInfo := RandomMultisigInfo(3, 3, chainParam, []int{0, 1, 2}, 0, 1)
	relayersMultisigWallet, _ := BuildMultisigWallet(relayersMultisigInfo)
	_, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	relayersMultisigAddress, err := AddressOnChain(chainParam, relayersMultisigWallet)
	fmt.Println("relayersMultisigAddress, err : ", relayersMultisigAddress, err)
	fmt.Println("relayersPubKey Master : ", hex.EncodeToString(relayersMultisigInfo.PubKeys[0]))
	fmt.Println("relayersPubKey Slave 1 : ", hex.EncodeToString(relayersMultisigInfo.PubKeys[1]))
	fmt.Println("relayersPubKey Slave 2 : ", hex.EncodeToString(relayersMultisigInfo.PubKeys[2]))
	fmt.Println("relayersPrivKey Master : ", relayerPrivKeys[0])
	fmt.Println("relayersPrivKey Slave 1 : ", relayerPrivKeys[1])
	fmt.Println("relayersPrivKey Slave 2 : ", relayerPrivKeys[2])

	userMultisigAddress, err := AddressOnChain(chainParam, userMultisigWallet)
	fmt.Println("userMultisigAddress, err : ", userMultisigAddress, err)
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

func TestTransferBitcoinWithBridgeMessage(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params
	// relayer multisig
	_, relayersMultisigInfo := RandomMultisigInfo(3, 2, chainParam, []int{0, 1, 2}, 0, 1)
	relayersMultisigWallet, _ := BuildMultisigWallet(relayersMultisigInfo)
	// user key
	userPrivKeys, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	bridgeMsg := BridgeDecodedMsg {
		&XCallMessage{
			MessageType:  1,
			Action:       "Deposit",
			TokenAddress: "0:0",                                                 // bitcoin address
			To:           "0x2.icon/hx39eddef484f6bb08072c59cc0a206e6713fc6d7d", //(receiver)
			From:         USER_MULTISIG_ADDRESS,                                 //user
			Amount:       new(big.Int).SetUint64(10000).Bytes(),
			Data:         []byte(""),
		},
		1,
		"cx0000000000000000000000000000000000000009",
		[]string{
			"hx0000000000000000000000000000000000000001",
			"hx0000000000000000000000000000000000000002",
			"hx0000000000000000000000000000000000000003",
		},
	}

	inputs := []*Input{
		// user bitcoin UTXO to transfer and pay tx fee
		{
			TxHash:			"d316231a8aa1f74472ed9cc0f1ed0e36b9b290254cf6b2c377f0d92b299868bf",
			OutputIdx:		4,
			OutputAmount:	1929000,
			PkScript:		userMultisigWallet.PKScript,
		},
	}

	// create tx
	msgTx, err := CreateBridgeTxSendBitcoin(
		&bridgeMsg,
		inputs,
		userMultisigWallet.PKScript,
		relayersMultisigWallet.PKScript,
		TX_FEE,
	)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	// user key 1 sign tx
	userSigs1, _ := SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs1)
	// user key 2 sign tx
	userSigs2, _ := SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs2)
	// COMBINE SIGN
	signedMsgTx, _ := CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

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
	fmt.Println("Message.MessageType: ", decodedMessage.Message.MessageType)
	fmt.Println("Message.Action: ", decodedMessage.Message.Action)
	fmt.Println("Message.TokenAddress: ", decodedMessage.Message.TokenAddress)
	fmt.Println("Message.From: ", decodedMessage.Message.From)
	fmt.Println("Message.To: ", decodedMessage.Message.To)
	fmt.Println("Message.Amount: ", decodedMessage.Message.Amount)
	fmt.Println("Message.Data: ", decodedMessage.Message.Data)
	fmt.Println("ChainId: ", decodedMessage.ChainId)
	fmt.Println("Address: ", decodedMessage.Receiver)
	fmt.Println("Connectors: ", decodedMessage.Connectors)
}

func TestTransferRuneWithBridgeMessage(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params
	// relayer multisig
	_, relayersMultisigInfo := RandomMultisigInfo(3, 2, chainParam, []int{0, 1, 2}, 0, 1)
	relayersMultisigWallet, _ := BuildMultisigWallet(relayersMultisigInfo)
	// user key
	userPrivKeys, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	bridgeMsg := BridgeDecodedMsg {
		&XCallMessage{
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
	}

	inputs := []*Input{
		// userrune UTXO used to transfer
		{
			TxHash:			"647a499a394bdb2a477f29b9f0515ed186e57a469a732be362a172cde4ea67a5",
			OutputIdx:		0,
			OutputAmount:	DUST_UTXO_AMOUNT,
			PkScript:		userMultisigWallet.PKScript,
		},
		// user bitcoin UTXO to pay tx fee
		{
			TxHash:			"d316231a8aa1f74472ed9cc0f1ed0e36b9b290254cf6b2c377f0d92b299868bf",
			OutputIdx:		4,
			OutputAmount:	1929000,
			PkScript:		userMultisigWallet.PKScript,
		},
	}

	// create tx
	msgTx, err := CreateBridgeTxSendRune(
		&bridgeMsg,
		inputs,
		userMultisigWallet.PKScript,
		relayersMultisigWallet.PKScript,
		TX_FEE,
	)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	// user key 1 sign tx
	userSigs1, _ := SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs1)
	// user key 2 sign tx
	userSigs2, _ := SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs2)
	// COMBINE SIGN
	signedMsgTx, _ := CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

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

	// Decode bridge message
	decodedMessage, err := ReadBridgeMessage(signedMsgTx)
	fmt.Println("err decode: ", err)
	fmt.Println("-----------decodedMessage---------- ")
	fmt.Println("Message.MessageType: ", decodedMessage.Message.MessageType)
	fmt.Println("Message.Action: ", decodedMessage.Message.Action)
	fmt.Println("Message.TokenAddress: ", decodedMessage.Message.TokenAddress)
	fmt.Println("Message.From: ", decodedMessage.Message.From)
	fmt.Println("Message.To: ", decodedMessage.Message.To)
	fmt.Println("Message.Amount: ", decodedMessage.Message.Amount)
	fmt.Println("Message.Data: ", decodedMessage.Message.Data)
	fmt.Println("ChainId: ", decodedMessage.ChainId)
	fmt.Println("Address: ", decodedMessage.Receiver)
	fmt.Println("Connectors: ", decodedMessage.Connectors)
}
