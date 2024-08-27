// relayer/chains/stacks/client_test.go

package stacks

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
)

func checkSenderKeyMatchesSigner(privateKeyHex string, signerArray [20]byte) error {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to decode private key: %v", err)
	}

	_, publicKey := btcec.PrivKeyFromBytes(privateKeyBytes)

	compressedPubKey := publicKey.SerializeCompressed()

	pubKeyHash160 := Hash160(compressedPubKey)

	if !bytes.Equal(pubKeyHash160, signerArray[:]) {
		return fmt.Errorf("senderKey does not correspond to signerArray")
	}

	return nil
}

func TestTokenTransferTransaction(t *testing.T) {
	mnemonic := "vapor unhappy gather snap project ball gain puzzle comic error avocado bounce letter anxiety wheel provide canyon promote sniff improve figure daughter mansion baby"
	expectedPrivateKey := "c1d5bb638aa70862621667f9997711fce692cad782694103f8d9561f62e9f19701"

	privateKey, err := DeriveStxPrivateKey(mnemonic, 0)
	hexPrivateKeyString := hex.EncodeToString(privateKey)
	if err != nil || hexPrivateKeyString != expectedPrivateKey {
		t.Fatalf("Failed to derive private key: %v", err)
	}

	senderPublicKey := GetPublicKeyFromPrivate(privateKey)
	var signerArray [20]byte
	copy(signerArray[:], Hash160(senderPublicKey))

	err = checkSenderKeyMatchesSigner(expectedPrivateKey, signerArray)
	if err != nil {
		t.Fatalf("Sender key does not match signer: %v", err)
	}

	recipient := "ST3YJD5Y1WTMC8R09ZKR3HJF562R3NM8HHXW2S2R9"
	amount := uint64(1000000) // 1 STX
	memo := "Test transfer"
	tx, err := NewTokenTransferTransaction(recipient, amount, memo, TransactionVersionTestnet, ChainIDTestnet, signerArray, 1, 180, AnchorModeOnChainOnly, PostConditionModeDeny)
	if err != nil {
		t.Fatalf("Failed to serialize payload: %v", err)
	}

	err = sendTransaction(tx, privateKey)
	if err != nil {
		t.Fatalf("Failed to send transaction: %v\n", err)
	}
}
