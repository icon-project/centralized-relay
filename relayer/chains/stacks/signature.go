// relayer/chains/stacks/signature.go
package stacks

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/crypto/ripemd160"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

type MessageSignature struct {
	Type StacksMessageType
	Data string
}

type StacksMessageType int

const (
	Address StacksMessageType = iota
	MessageSignatureType
)

type StacksPrivateKey struct {
	Data       []byte
	Compressed bool
}

func SignWithKey(privateKey []byte, messageHash string) (MessageSignature, error) {
	privKey, _ := btcec.PrivKeyFromBytes(privateKey[:32])
	messageHashBytes, err := hex.DecodeString(messageHash)
	if err != nil {
		return MessageSignature{}, err
	}
	signature, err := ecdsa.SignCompact(privKey, messageHashBytes, true)
	if err != nil {
		return MessageSignature{}, err
	}
	recoveryID := signature[0] - 27 - 4
	vrsSignature := fmt.Sprintf("%02x%s", recoveryID, hex.EncodeToString(signature[1:]))
	return CreateMessageSignature(vrsSignature)
}

func CreateMessageSignature(signature string) (MessageSignature, error) {
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return MessageSignature{}, err
	}
	if len(signatureBytes) != RecoverableECDSASigLengthBytes {
		return MessageSignature{}, errors.New("invalid signature")
	}
	return MessageSignature{
		Type: MessageSignatureType,
		Data: signature,
	}, nil
}

func GetPublicKeyFromPrivate(privateKey []byte) []byte {
	_, pubKey := btcec.PrivKeyFromBytes(privateKey)
	return pubKey.SerializeCompressed()
}

func VerifySignature(messageHash string, signature MessageSignature, publicKey []byte) (bool, error) {
	messageHashBytes, err := hex.DecodeString(messageHash)
	if err != nil {
		return false, errors.New("invalid message hash")
	}

	signatureBytes, err := hex.DecodeString(signature.Data)
	if err != nil {
		return false, errors.New("invalid signature")
	}

	// The signature is in [R || S] format, where R and S are 32 bytes each
	if len(signatureBytes) != 65 {
		return false, errors.New("invalid signature length")
	}

	// Parse the public key
	pubKey, err := btcec.ParsePubKey(publicKey)
	if err != nil {
		return false, errors.New("failed to parse public key")
	}

	// Create a new ECDSA signature from R and S components
	r := new(btcec.ModNScalar)
	r.SetByteSlice(signatureBytes[1:33])
	s := new(btcec.ModNScalar)
	s.SetByteSlice(signatureBytes[33:])
	ecdsaSignature := ecdsa.NewSignature(r, s)

	// Verify the signature
	return ecdsaSignature.Verify(messageHashBytes, pubKey), nil
}

func Hash160(b []byte) []byte {
	h := sha256.Sum256(b)
	ripemd160Hasher := ripemd160.New()
	ripemd160Hasher.Write(h[:])
	return ripemd160Hasher.Sum(nil)
}

func calculateSighash(serializedTx []byte) []byte {
	hash := sha256.Sum256(serializedTx)
	return hash[:]
}

func calculatePresignSighash(sighash []byte, authType AuthType, fee uint64, nonce uint64) []byte {
	data := make([]byte, 0, len(sighash)+1+8+8)
	data = append(data, sighash...)
	data = append(data, byte(authType))
	feeBytes := make([]byte, 8)
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(feeBytes, fee)
	binary.BigEndian.PutUint64(nonceBytes, nonce)
	data = append(data, feeBytes...)
	data = append(data, nonceBytes...)

	h := sha512.New512_256()
	h.Write(data)
	return h.Sum(nil)
}

func SignTransaction(tx *TokenTransferTransaction, privateKey []byte) error {
	// 1. Clear the signature in the spending condition
	tx.Auth.OriginAuth.Signature = [65]byte{}

	// 2. Serialize the transaction
	serializedTx, err := tx.Serialize()
	if err != nil {
		return err
	}

	// 3. Calculate the initial sighash
	sighash := calculateSighash(serializedTx)

	// 4. Calculate the presign-sighash
	presignSighash := calculatePresignSighash(sighash, tx.Auth.AuthType, tx.Auth.OriginAuth.Fee, tx.Auth.OriginAuth.Nonce)

	// 5. Sign the presign-sighash
	signature, err := SignWithKey(privateKey, hex.EncodeToString(presignSighash))
	if err != nil {
		return err
	}

	// 6. Set the signature in the spending condition
	signatureBytes, _ := hex.DecodeString(signature.Data)
	copy(tx.Auth.OriginAuth.Signature[:], signatureBytes)

	return nil
}

func VerifyTransaction(tx *TokenTransferTransaction, publicKey []byte) (bool, error) {
	txCopy := *tx

	// 1. Extract the signature
	signature := txCopy.Auth.OriginAuth.Signature

	// 2. Clear the signature in the spending condition
	txCopy.Auth.OriginAuth.Signature = [RecoverableECDSASigLengthBytes]byte{}

	// 3. Serialize the transaction
	serializedTx, err := txCopy.Serialize()
	if err != nil {
		return false, err
	}

	// 4. Calculate the initial sighash
	sighash := calculateSighash(serializedTx)

	// 5. Calculate the presign-sighash
	presignSighash := calculatePresignSighash(sighash, txCopy.Auth.AuthType, txCopy.Auth.OriginAuth.Fee, txCopy.Auth.OriginAuth.Nonce)

	// 6. Verify the signature
	messageSignature := MessageSignature{
		Type: MessageSignatureType,
		Data: hex.EncodeToString(signature[:]),
	}
	return VerifySignature(hex.EncodeToString(presignSighash), messageSignature, publicKey)
}
