package stacks

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
)

func getNonce(address string, network StacksNetwork) (*big.Int, error) {
	url := network.GetAccountAPIURL(address)

	resp, err := network.FetchFn(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching nonce: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error fetching nonce. Response %d: %s. URL: %s, Body: %s",
			resp.StatusCode, resp.Status, url, string(body))
	}

	var result struct {
		Nonce uint64 `json:"nonce"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %w", err)
	}

	nonce := big.NewInt(int64(result.Nonce))
	return nonce, nil
}

func estimateTransactionFeeWithFallback(tx TokenTransferTransaction, network StacksNetwork) (*big.Int, error) {
	fee, err := estimateTransaction(tx, network)
	if err == nil {
		return fee, nil
	}

	return estimateTransferUnsafe(tx, network)
}

func estimateTransaction(tx TokenTransferTransaction, network StacksNetwork) (*big.Int, error) {
	url := network.GetTransactionFeeEstimateAPIURL()
	serializedTx, err := tx.Serialize()
	if err != nil {
		return nil, fmt.Errorf("error serializing transaction: %w", err)
	}

	byteLength, err := estimateTransactionByteLength(tx)
	if err != nil {
		return nil, fmt.Errorf("error estimating transaction byte length: %w", err)
	}

	payload := map[string]interface{}{
		"transaction_payload": serializedTx,
		"estimated_len":       byteLength,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error sending estimation request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error estimating fee. Response %d: %s. URL: %s, Body: %s",
			resp.StatusCode, resp.Status, url, string(body))
	}

	var result struct {
		EstimatedCost struct {
			FeeRate string `json:"fee_rate"`
		} `json:"estimated_cost"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %w", err)
	}

	fee, success := new(big.Int).SetString(result.EstimatedCost.FeeRate, 10)
	if !success {
		return nil, fmt.Errorf("failed to parse fee as big.Int: %s", result.EstimatedCost.FeeRate)
	}

	return fee, nil
}

func estimateTransferUnsafe(tx TokenTransferTransaction, network StacksNetwork) (*big.Int, error) {
	url := network.GetTransferFeeEstimateAPIURL()

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching fee estimate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error fetching fee estimate. Response %d: %s. URL: %s, Body: %s",
			resp.StatusCode, resp.Status, url, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	feeRate, success := new(big.Int).SetString(string(body), 10)
	if !success {
		return nil, fmt.Errorf("failed to parse fee rate as big.Int: %s", string(body))
	}

	txBytes, err := tx.Serialize()
	if err != nil {
		return nil, fmt.Errorf("error serializing transaction: %w", err)
	}

	txBytesLen := big.NewInt(int64(len(txBytes)))
	fee := new(big.Int).Mul(feeRate, txBytesLen)

	return fee, nil
}

func MakeSTXTokenTransfer(
	recipient string,
	amount big.Int,
	memo string,
	network StacksNetwork,
	senderAddress string,
	senderKey []byte,
	fee *big.Int,
	nonce *big.Int,
) (*TokenTransferTransaction, error) {
	if recipient == "" || len(senderKey) == 0 {
		return nil, fmt.Errorf("invalid parameters: recipient or senderKey are empty")
	}

	senderPublicKey := GetPublicKeyFromPrivate(senderKey)
	var signer [20]byte
	copy(signer[:], Hash160(senderPublicKey))

	tx, err := NewTokenTransferTransaction(recipient, amount.Uint64(), memo, network.version, network.chainID, signer, 0, 0, AnchorModeOnChainOnly, PostConditionModeDeny)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if fee == nil {
		estimatedFee, err := estimateTransactionFeeWithFallback(*tx, network)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate fee: %w", err)
		}
		fee = estimatedFee
	}

	if nonce == nil {
		fetchedNonce, err := getNonce(senderAddress, network)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch nonce: %w", err)
		}
		nonce = fetchedNonce
	}

	tx.Auth.OriginAuth.Fee = fee.Uint64()
	tx.Auth.OriginAuth.Nonce = nonce.Uint64()

	err = SignTransaction(tx, senderKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return tx, nil
}

func estimateTransactionByteLength(tx TokenTransferTransaction) (int, error) {
	serializedTx, err := tx.Serialize()
	if err != nil {
		return 0, fmt.Errorf("error serializing transaction: %w", err)
	}

	return len(serializedTx), nil
}

type BroadcastResponse struct {
	TxId string `json:"txid"`
}

func BroadcastTransaction(tx *TokenTransferTransaction, network *StacksNetwork) (string, error) {
	serializedTx, err := tx.Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize transaction: %w", err)
	}

	url := network.GetBroadcastAPIURL()

	requestBody := bytes.NewBuffer(serializedTx)

	resp, err := http.Post(url, "application/octet-stream", requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("broadcast failed with status %d: %s", resp.StatusCode, string(body))
	}

	var broadcastResp BroadcastResponse
	err = json.Unmarshal(body, &broadcastResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if !isValidTransactionID(broadcastResp.TxId) {
		return "", fmt.Errorf("received invalid transaction ID: %s", broadcastResp.TxId)
	}

	return broadcastResp.TxId, nil
}

func isValidTransactionID(txID string) bool {
	if len(txID) != 64 {
		return false
	}

	_, err := hex.DecodeString(txID)
	return err == nil
}
