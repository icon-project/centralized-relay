// relayer/chains/stacks/signature_test.go
package stacks

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/clarity"
)

func privateKeyToBytes(privateKey interface{}) ([]byte, error) {
	var privateKeyBuffer []byte
	switch v := privateKey.(type) {
	case string:
		var err error
		privateKeyBuffer, err = hex.DecodeString(v)
		if err != nil {
			return nil, err
		}
	case []byte:
		privateKeyBuffer = v
	default:
		return nil, errors.New("privateKey must be a string or []byte")
	}

	if len(privateKeyBuffer) != 32 && len(privateKeyBuffer) != 33 {
		return nil, errors.New("improperly formatted private-key. Private-key byte length should be 32 or 33")
	}

	if len(privateKeyBuffer) == 33 && privateKeyBuffer[32] != 1 {
		return nil, errors.New("improperly formatted private-key. 33 bytes indicate compressed key, but the last byte must be == 01")
	}

	return privateKeyBuffer, nil
}

func createStacksPrivateKey(key interface{}) (StacksPrivateKey, error) {
	data, err := privateKeyToBytes(key)
	if err != nil {
		return StacksPrivateKey{}, err
	}
	compressed := len(data) == CompressedPubkeyLengthBytes
	return StacksPrivateKey{Data: data, Compressed: compressed}, nil
}

func TestSignWithKey(t *testing.T) {
	privateKey, err := createStacksPrivateKey("bcf62fdd286f9b30b2c289cce3189dbf3b502dcd955b2dc4f67d18d77f3e73c7")
	if err != nil {
		t.Fatalf("Failed to create private key: %v", err)
	}

	publicKey := GetPublicKeyFromPrivate(privateKey.Data)

	expectedMessageHash := "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e"
	expectedSignatureVrs := "00f540e429fc6e8a4c27f2782479e739cae99aa21e8cb25d4436f333577bc791cd1d9672055dd1604dd5194b88076e4f859dd93c834785ed589ec38291698d4142"

	hash := calculateSighash([]byte("Hello World"))
	messageHash := hex.EncodeToString(hash[:])
	if messageHash != expectedMessageHash {
		t.Fatalf("Message hash doesn't match expected. Got %s, want %s", messageHash, expectedMessageHash)
	}

	signature, err := SignWithKey(privateKey.Data, messageHash)
	if err != nil {
		t.Fatalf("Failed to sign: %v", err)
	}

	if signature.Data != expectedSignatureVrs {
		t.Fatalf("Signature doesn't match expected. Got %s, want %s", signature.Data, expectedSignatureVrs)
	}

	isValid, err := VerifySignature(messageHash, signature, publicKey)
	if err != nil {
		t.Fatalf("Error verifying signature: %v", err)
	}

	if !isValid {
		t.Fatalf("Signature verification failed: expected valid signature")
	}

	incorrectMessageHash := "0000000000000000000000000000000000000000000000000000000000000000"
	isValid, err = VerifySignature(incorrectMessageHash, signature, publicKey)
	if err != nil {
		t.Fatalf("Signature verification failed: %v", err)
	}

	if isValid {
		t.Errorf("Signature verification failed: expected invalid signature for incorrect message hash")
	}

	incorrectPublicKey := make([]byte, len(publicKey))
	copy(incorrectPublicKey, publicKey)
	incorrectPublicKey[0] ^= 0xFF // Flip bits in the first byte

	isValid, _ = VerifySignature(messageHash, signature, incorrectPublicKey)

	if isValid {
		t.Errorf("Signature verification failed: expected invalid signature for incorrect public key")
	}
}

func TestTransactionSignAndVerify(t *testing.T) {
	privateKey, err := createStacksPrivateKey("bcf62fdd286f9b30b2c289cce3189dbf3b502dcd955b2dc4f67d18d77f3e73c7")
	if err != nil {
		t.Fatalf("Failed to create private key: %v", err)
	}

	publicKey := GetPublicKeyFromPrivate(privateKey.Data)

	tx := createTestTransaction()

	tx.Auth.OriginAuth.Fee = 1000
	tx.Auth.OriginAuth.Nonce = 123

	err = SignTransaction(tx, privateKey.Data)
	if err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}

	isValid, err := VerifyTransaction(tx, publicKey)
	if err != nil {
		t.Fatalf("Failed to verify transaction: %v", err)
	}

	if !isValid {
		t.Fatalf("Transaction signature verification failed")
	}

	incorrectPublicKey := make([]byte, len(publicKey))
	copy(incorrectPublicKey, publicKey)
	incorrectPublicKey[0] ^= 0xFF // Flip bits in the first byte

	isValid, _ = VerifyTransaction(tx, incorrectPublicKey)
	if isValid {
		t.Fatalf("Transaction verification should fail with incorrect public key")
	}

	tx.Payload.Amount += 1
	isValid, _ = VerifyTransaction(tx, publicKey)
	if isValid {
		t.Fatalf("Transaction verification should fail with modified transaction data")
	}
}

func createTestTransaction() *TokenTransferTransaction {
	recipientPrincipal, _ := clarity.StringToPrincipal("SP3FGQ8Z7JY9BWYZ5WM53E0M9NK7WHJF0691NZ159")
	return &TokenTransferTransaction{
		BaseTransaction: BaseTransaction{
			Version: TransactionVersion(0),
			ChainID: ChainID(1),
			Auth: TransactionAuth{
				AuthType: AuthType(1),
				OriginAuth: SpendingCondition{
					HashMode:    AddressHashMode(1),
					Signer:      [20]byte{},
					Nonce:       0,
					Fee:         0,
					KeyEncoding: PubKeyEncoding(1),
					Signature:   [65]byte{},
				},
			},
			AnchorMode:        AnchorMode(1),
			PostConditionMode: PostConditionMode(1),
			PostConditions:    []PostCondition{},
		},
		Payload: TokenTransferPayload{
			Recipient: recipientPrincipal,
			Amount:    12345,
			Memo:      "test memo",
		},
	}
}

func TestCalculatePresignSighash(t *testing.T) {
	curSigHash := "2be5719a864803128bc90aba66a2c42a28e07431fb3c7362f43dadf3c6be6cb5"
	authType := AuthTypeStandard
	fee := uint64(1)
	nonce := uint64(1)

	expectedHash := "9047149584c3e1af556484afc14dd599351c04cf5caca37c6e4d438490cead7b"

	curSigHashBytes, _ := hex.DecodeString(curSigHash)
	result := calculatePresignSighash(curSigHashBytes, authType, fee, nonce)

	if hex.EncodeToString(result) != expectedHash {
		t.Errorf("calculatePresignSighash mismatch. Got %s, want %s", hex.EncodeToString(result), expectedHash)
	}
}

func TestSignWithKey2(t *testing.T) {
	privateKey := "c1d5bb638aa70862621667f9997711fce692cad782694103f8d9561f62e9f19701"
	messageHash := "9047149584c3e1af556484afc14dd599351c04cf5caca37c6e4d438490cead7b"

	expectedSignature := "001ca86b18448768a0eacc41ed749eceeb2c3073e424e9fc36b8e86481d92ad2244ac572e30782158b5a1d954a0f50ed3f6f02d1ff9c9dd6efe6d14c52878c52da"

	privateKeyBytes, _ := hex.DecodeString(privateKey)
	signature, err := SignWithKey(privateKeyBytes, messageHash)
	if err != nil {
		t.Fatalf("SignWithKey failed: %v", err)
	}

	if signature.Data != expectedSignature {
		t.Fatalf("SignWithKey signature mismatch. Got %s, want %s", signature.Data, expectedSignature)
	}
}

// func TestCalculatePostSignSighash(t *testing.T) {
// 	presignSighash := "9047149584c3e1af556484afc14dd599351c04cf5caca37c6e4d438490cead7b"
// 	publicKey := "0332fc778e5beb5f944c75b2b63c21dd12c40bdcdf99ba0663168ae0b2be880aef"
// 	signature := MessageSignature{
// 		Type: MessageSignatureType,
// 		Data: "001ca86b18448768a0eacc41ed749eceeb2c3073e424e9fc36b8e86481d92ad2244ac572e30782158b5a1d954a0f50ed3f6f02d1ff9c9dd6efe6d14c52878c52da",
// 	}

// 	expectedHash := "f3a8b86e09c962a61074207427d6236e85cb04cf0e3ff3a2f44be2eb7a51a4b6"

// 	presignSighashBytes, _ := hex.DecodeString(presignSighash)
// 	publicKeyBytes, _ := hex.DecodeString(publicKey)
// 	result := calculatePostSignSighash(presignSighashBytes, publicKeyBytes, signature)

// 	if hex.EncodeToString(result) != expectedHash {
// 		t.Errorf("calculatePostSignSighash mismatch. Got %s, want %s", hex.EncodeToString(result), expectedHash)
// 	}
// }
