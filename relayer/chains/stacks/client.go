// relayer/chains/stacks/client.go
package stacks

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/go-resty/resty/v2"
)

const stacksNodeURL = "https://api.testnet.hiro.so/"

func sendTransaction(tx *TokenTransferTransaction, privateKey []byte) error {
	err := SignTransaction(tx, privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	serializedTx, err := tx.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize transaction: %w", err)
	}

	deserializedTx, err := DeserializeTransaction(serializedTx)
	if err != nil {
		return fmt.Errorf("failed to serialize transaction: %w", err)
	}

	_, publicKey := btcec.PrivKeyFromBytes(privateKey)
	compressedPubKey := publicKey.SerializeCompressed()

	pubKeyHash160 := Hash160(compressedPubKey)

	if !bytes.Equal(pubKeyHash160, deserializedTx.(*TokenTransferTransaction).Auth.OriginAuth.Signer[:]) {
		return fmt.Errorf("senderKey does not correspond to signerArray")
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/octet-stream").
		SetBody(serializedTx).
		Post(stacksNodeURL + "/v2/transactions")

	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("transaction submission failed with status code: %d, body: %s", resp.StatusCode(), string(resp.Body()))
	}

	var txResponse struct {
		TxID string
	}
	err = json.Unmarshal(resp.Body(), &txResponse)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	fmt.Printf("Transaction submitted successfully. TxID: %s\n", txResponse.TxID)

	return nil
}
