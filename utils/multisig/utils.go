package multisig

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
)

func GeneratePrivateKeyFromSeed(seed []byte, chainParam *chaincfg.Params) *btcec.PrivateKey {
	seedDigest := sha256.Sum256([]byte(seed))
	masterKey, _ := hdkeychain.NewMaster(seedDigest[:], chainParam)
	for _, childNum := range []uint32{1, 1, 1} {
		var err error
		masterKey, err = masterKey.Derive(hdkeychain.HardenedKeyStart + childNum)
		if err != nil {
			panic(err)
		}
	}
	privateKey, _ := masterKey.ECPrivKey()
	return privateKey
}

func randomKeys(n int, chainParam *chaincfg.Params) ([]string, [][]byte, []*btcutil.AddressPubKey) {

	privKeys := []string{}
	pubKeys := [][]byte{}
	ECPubKeys := []*btcutil.AddressPubKey{}

	for i := 0; i < n; i++ {
		privKey := GeneratePrivateKeyFromSeed([]byte{byte(i)}, chainParam)
		wif, _ := btcutil.NewWIF(privKey, chainParam, true)

		ECPubKey, _ := btcutil.NewAddressPubKey(wif.SerializePubKey(), chainParam)

		pubKey := privKey.PubKey().SerializeCompressed()
		pubKey2 := wif.SerializePubKey()

		fmt.Printf("Account %v:\n%v\n%v\n%v\n", i, wif.String(), hex.EncodeToString(wif.SerializePubKey()), wif.PrivKey.Serialize())
		privKeys = append(privKeys, wif.String())
		pubKeys = append(pubKeys, wif.SerializePubKey())
		ECPubKeys = append(ECPubKeys, ECPubKey)

		fmt.Printf("randomKeys pubKey %v - pubKey 2 %v\n ", pubKey, pubKey2)
	}
	fmt.Printf("len privateKey: %v\n", len(privKeys))
	fmt.Printf("len publicKey: %v\n", len(pubKeys))

	return privKeys, pubKeys, ECPubKeys
}


func randomMultisigInfo(n int, k int, chainParam *chaincfg.Params) ([]string, *MultisigInfo) {
	privKeys, pubKeys, EcPubKeys := randomKeys(n, chainParam)
	vaultInfo := MultisigInfo{
		PubKeys:            pubKeys,
		EcPubKeys:          EcPubKeys,
		NumberRequiredSigs: k,
	}

	return privKeys, &vaultInfo

}