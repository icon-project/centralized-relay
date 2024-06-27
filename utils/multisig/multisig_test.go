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
)

func TestGenerateKeys(t *testing.T) {
	chainParam := &chaincfg.RegressionNetParams

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
	chainParam := &chaincfg.RegressionNetParams

	wif, _ := btcutil.DecodeWIF("cTYRscQxVhtsGjHeV59RHQJbzNnJHbf3FX4eyX5JkpDhqKdhtRvy")
	pubKey := wif.SerializePubKey();
	witnessProg := btcutil.Hash160(pubKey)
	p2wpkh, _ := btcutil.NewAddressWitnessPubKeyHash(witnessProg, chainParam)

	fmt.Printf("Account:\n Private Key: %v\n Public Key: %v\n Address: %v\n", string(wif.String()), hex.EncodeToString(pubKey), p2wpkh)
}

func TestRandomKeys(t *testing.T) {
	randomKeys(3, &chaincfg.RegressionNetParams)
}

func TestBuildMultisigTapScript(t *testing.T) {
	// 2/3: bcrt1phdyt24adauupp7tawuu9ksl7gvtflr70raj3f2dzwzn06q5vhyhq0l43lz
	// 3/3: bcrt1py04eh93ae0e6dpps2ufxt58wjnvesj0ffzddcckmru3tyrhzsslstaag7h
	totalSigs := 3
	numSigsRequired := 2
	chainParam := &chaincfg.RegressionNetParams
	// 3 for multisig vault, 1 for recovery key
	_, pubKeys, ECPubKeys := randomKeys(totalSigs, chainParam)

	fmt.Printf("Len pub key: %v\n", len(pubKeys[0]))

	multisigInfo := &MultisigInfo{
		PubKeys:            pubKeys,
		EcPubKeys:          ECPubKeys,
		NumberRequiredSigs: numSigsRequired,
	}
	multisigWallet, _ := GenerateMultisigWallet(multisigInfo)
	multisigAddress, err := AddressOnChain(chainParam, multisigWallet)
	fmt.Println("address, err : ", multisigAddress, err)
}

func TestCreateTx(t *testing.T) {
	chainParam := &chaincfg.RegressionNetParams
	privKeys, multisigInfo := randomMultisigInfo(3, 2, chainParam)

	inputs := []*UTXO{
		// 2/3 - empty data
		{
			TxHash:        "62d19039c9d0eec493f3a1440f0fab65c525b1426b675445b01f26ddf1d8fa42",
			OutputIdx:     0,
			OutputAmount:  10000,
		},
		{
			TxHash:        "8f476a9a520f548e7b60512f5c14c5c6253a289dde02d146a02ca22892a2877a",
			OutputIdx:     0,
			OutputAmount:  20000,
		},
	}

	outputs := []*OutputTx{
		{
			ReceiverAddress: "bcrt1p65j57tzjufnjmt4fgx5xexfry6f3f87sggl02gl7fcxuky4x34fscyjejf",
			Amount:          8000,
		},
	}

	multisigWallet, _ := GenerateMultisigWallet(multisigInfo)

	changeReceiverAddress := "bcrt1phdyt24adauupp7tawuu9ksl7gvtflr70raj3f2dzwzn06q5vhyhq0l43lz"
	msgTx, prevOuts, err := CreateMultisigTx(inputs, outputs, 1000, chainParam, changeReceiverAddress, multisigWallet.PKScript)
	fmt.Println("msgTx: ", msgTx)
	fmt.Println("prevOuts: ", prevOuts)
	fmt.Println("err: ", err)

	// validators sign tx
	totalSigs := [][][]byte{} // index private key -> index input -> sig
	for _, privKey := range privKeys {
		sigs, err := PartSignOnRawExternalTx(privKey, msgTx, inputs, multisigWallet, 0, chainParam)
		if err != nil {
			fmt.Println("err sign: ", err)
		}

		totalSigs = append(totalSigs, sigs)
	}

	signedMsgTx, err := CombineMultisigSigs(multisigInfo, multisigWallet, msgTx, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)
}

func TestMultiRelayers(t *testing.T) {
	chainParam := &chaincfg.RegressionNetParams
	privKeys, multisigInfo := randomMultisigInfo(3, 2, chainParam)

	inputs := []*UTXO{
		{
			TxHash:        "62d19039c9d0eec493f3a1440f0fab65c525b1426b675445b01f26ddf1d8fa42",
			OutputIdx:     0,
			OutputAmount:  10000,
		},
		{
			TxHash:        "8f476a9a520f548e7b60512f5c14c5c6253a289dde02d146a02ca22892a2877a",
			OutputIdx:     0,
			OutputAmount:  20000,
		},
	}

	outputs := []*OutputTx{
		{
			ReceiverAddress: "bcrt1p65j57tzjufnjmt4fgx5xexfry6f3f87sggl02gl7fcxuky4x34fscyjejf",
			Amount:          8000,
		},
	}

	multisigWallet, _ := GenerateMultisigWallet(multisigInfo)

	changeReceiverAddress := "bcrt1phdyt24adauupp7tawuu9ksl7gvtflr70raj3f2dzwzn06q5vhyhq0l43lz"
	msgTx, _, _ := CreateMultisigTx(inputs, outputs, 1000, chainParam, changeReceiverAddress, multisigWallet.PKScript)

	totalSigs := [][][]byte{}
	// MATSTER RELAYER SIGN TX
	sigs, err := PartSignOnRawExternalTx(privKeys[0], msgTx, inputs, multisigWallet, 0, chainParam)
	if err != nil {
		fmt.Println("err sign: ", err)
	}
	totalSigs = append(totalSigs, sigs)

	router := SetUpRouter()
	// create post body using an instance of the requestSignInput struct
	rsi := requestSignInput{
		MsgTx: msgTx,
		Inputs:  inputs,
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
	signedMsgTx, err := CombineMultisigSigs(multisigInfo, multisigWallet, msgTx, totalSigs)

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
	// hexSignedTx := "02000000000104f6ce922d4b636e81b5fe301d541f14a07ca8b5ee9e7c9637479fbb77ae76ac9c0000000000ffffffffd9a54f6610fbc9b56cea421230b79edf5c48d92052c13c166656591f2d9ed1fd0000000000ffffffff7b02c2d6a2c8d14d3a1c717f4093306a2c485b4a0cff672ed26715555dff1fde0000000000fffffffff6ce922d4b636e81b5fe301d541f14a07ca8b5ee9e7c9637479fbb77ae76ac9c0100000000ffffffff05f82a000000000000225120863583d69b5e4ce5edd053d72148074b6fe8968dc9175d9e1610282faf0ff3cd3c3d0900000000002251203cb3f6132189f271733e94476de372e022ea84dc2cb914c7ff691c453ddce563bcb1000000000000160014a355d136171b816b2b2f08d531cf94555d796751e803000000000000225120863583d69b5e4ce5edd053d72148074b6fe8968dc9175d9e1610282faf0ff3cda22d020000000000225120863583d69b5e4ce5edd053d72148074b6fe8968dc9175d9e1610282faf0ff3cd01401bdbfd310cbd31807024070825fc5a3fec6e7b37eb1d359dd3027727e5719c52e68b178f438d18c9218b13e3d95b8e3bd8efe1b422778fb88a59a8e5d4a7ba1601418f67cd9c4799e32ac9aed36c14a60e0a2442c4497cddb176ae118a9aeec099f0eaa68901681b0d6df6cdf45c062b70b73171d002d2c50345ea3c53bc51b6e8198301410315713b00ab6f47392e947e90aa4c67c37b77f9c76d1c5cb80a0df5f725dde2b883177dda0abe8dcbc7e06453c2720fa0eefea78e4b905a333929ee03e065bf830140430cc6f8e554c98930c7e67e9512bdc3c73a6944a60a58033607a9b195c857d0426ef570fb6640b07234ef33402a2ce22c66199ec57cbda27ce6028d3258f53100000000"

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
