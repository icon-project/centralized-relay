package stacks

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/clarity"
)

func TestMakeSTXTokenTransfer(t *testing.T) {
	recipient := "ST3YJD5Y1WTMC8R09ZKR3HJF562R3NM8HHXW2S2R9"
	amount := big.NewInt(1000000) // 1 STX
	memo := "Test transfer"
	network := NewStacksTestnet()
	senderAddress := "ST2CY5V39NHDPWSXMW9QDT3HC3GD6Q6XX4CFRK9AG"
	senderKeyHex := "c1d5bb638aa70862621667f9997711fce692cad782694103f8d9561f62e9f19701"
	senderKey, _ := hex.DecodeString(senderKeyHex)

	tx, err := MakeSTXTokenTransfer(recipient, *amount, memo, *network, senderAddress, senderKey, nil, nil)

	if err != nil {
		t.Fatalf("makeSTXTokenTransfer failed: %v", err)
	}

	expectedRecipient, err := clarity.StringToPrincipal(recipient)
	if err != nil {
		t.Fatalf("failed to cast recipient string to clarity principal: %v", err)
	}

	recipientPrincipal, _ := tx.Payload.Recipient.(*clarity.StandardPrincipal)
	expectedPrincipal, _ := expectedRecipient.(*clarity.StandardPrincipal)

	if !bytes.Equal(recipientPrincipal.Hash160[:], expectedPrincipal.Hash160[:]) {
		t.Errorf("Expected recipient %s, got %s", expectedRecipient, tx.Payload.Recipient)
	}

	if tx.Payload.Amount != amount.Uint64() {
		t.Errorf("Expected amount %d, got %d", amount.Uint64(), tx.Payload.Amount)
	}

	if tx.Payload.Memo != memo {
		t.Errorf("Expected memo %s, got %s", memo, tx.Payload.Memo)
	}

	if tx.Auth.OriginAuth.Fee == 0 {
		t.Error("Fee was not estimated")
	}

	if tx.Auth.OriginAuth.Nonce == 0 {
		t.Error("Nonce was not fetched")
	}

	if len(tx.Auth.OriginAuth.Signature) == 0 {
		t.Error("Transaction was not signed")
	}

	t.Logf("Estimated fee: %d", tx.Auth.OriginAuth.Fee)
	t.Logf("Fetched nonce: %d", tx.Auth.OriginAuth.Nonce)

	specifiedFee := big.NewInt(180)
	specifiedNonce := big.NewInt(3)
	tx2, err := MakeSTXTokenTransfer(recipient, *amount, memo, *network, senderAddress, senderKey, specifiedFee, specifiedNonce)

	if err != nil {
		t.Fatalf("makeSTXTokenTransfer with specified fee and nonce failed: %v", err)
	}

	if tx2.Auth.OriginAuth.Fee != specifiedFee.Uint64() {
		t.Errorf("Expected fee %d, got %d", specifiedFee.Uint64(), tx2.Auth.OriginAuth.Fee)
	}

	if tx2.Auth.OriginAuth.Nonce != specifiedNonce.Uint64() {
		t.Errorf("Expected nonce %d, got %d", specifiedNonce.Uint64(), tx2.Auth.OriginAuth.Nonce)
	}
}
