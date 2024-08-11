package stacks

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/crypto/ripemd160"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

func Hash160(b []byte) []byte {
	h := sha256.Sum256(b)
	ripemd160Hasher := ripemd160.New()
	ripemd160Hasher.Write(h[:])
	return ripemd160Hasher.Sum(nil)
}

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

func VerifySignature(publicKey []byte, message []byte, signature []byte) bool {
	pubKey, err := btcec.ParsePubKey(publicKey)
	if err != nil {
		return false
	}

	// Hash the message if it's not already a 32-byte hash
	var messageHash []byte
	if len(message) != 32 {
		h := sha256.Sum256(message)
		messageHash = h[:]
	} else {
		messageHash = message
	}

	// Adjust the signature to ECDSA compact format
	compactSignature := make([]byte, 65)
	copy(compactSignature, signature)
	compactSignature[64] -= 27 // Adjust recovery ID

	sig, err := ecdsa.ParseSignature(compactSignature[:64])
	if err != nil {
		return false
	}

	return sig.Verify(messageHash, pubKey)
}

func RecoverPublicKey(message []byte, signature []byte) ([]byte, error) {
	// Hash the message if it's not already a 32-byte hash
	var messageHash []byte
	if len(message) != 32 {
		h := sha256.Sum256(message)
		messageHash = h[:]
	} else {
		messageHash = message
	}

	// Adjust the signature to ECDSA compact format
	compactSignature := make([]byte, 65)
	copy(compactSignature, signature)
	compactSignature[64] -= 27 // Adjust recovery ID

	pubKey, _, err := ecdsa.RecoverCompact(compactSignature, messageHash)
	if err != nil {
		return nil, errors.New("failed to recover public key: " + err.Error())
	}

	return pubKey.SerializeCompressed(), nil
}

func (t *TokenTransferTransaction) Sign(privateKey []byte) error {
	// 1. Clear the other spending condition fields. If this is a single-signature spending condition, then set the fee and nonce to 0, and set the signature bytes to 0 (note that the address is preserved).
	t.Auth.OriginAuth.Signature = [RecoverableECDSASigLengthBytes]byte{}

	// 2. Serialize the transaction into a byte sequence, and hash it to form an initial sighash.
	serializedTx, err := t.Serialize()
	if err != nil {
		return err
	}
	sighash := sha256.Sum256(serializedTx)

	// 3. Calculate the presign-sighash over the sighash by hashing the sighash with the authorization type byte (0x04 or 0x05), the fee (as an 8-byte big-endian value), and the nonce (as an 8-byte big-endian value).
	var buf bytes.Buffer
	buf.WriteByte(byte(t.Auth.AuthType))
	binary.Write(&buf, binary.BigEndian, t.Auth.OriginAuth.Fee)
	binary.Write(&buf, binary.BigEndian, t.Auth.OriginAuth.Nonce)
	buf.Write(sighash[:])
	presignSighash := sha256.Sum256(buf.Bytes())

	// 4. Calculate the ECDSA signature over the presign-sighash by treating this hash as the message digest. Note that the signature must be a libsecp256k1 recoverable signature.
	privKey, pubKey := btcec.PrivKeyFromBytes(privateKey)
	sig, err := ecdsa.SignCompact(privKey, presignSighash[:], true)
	if err != nil {
		return err
	}

	// 5. Calculate the postsign-sighash over the resulting signature and public key by hashing the presign-sighash hash, the signing key's public key encoding byte, and the signature from step 4 to form the next sighash. Store the message signature and public key encoding byte as a signature auth field.
	copy(t.Auth.OriginAuth.Signature[:], sig)
	t.Auth.OriginAuth.KeyEncoding = PubKeyEncodingCompressed
	hash160 := Hash160(pubKey.SerializeCompressed())
	copy(t.Auth.OriginAuth.Signer[:], hash160[:])

	fmt.Printf("Signer calculated in Sign method: %x\n", t.Auth.OriginAuth.Signer)

	// For single-signature spending conditions, the only data the signing algorithm needs to return is the public key encoding byte and message signature.
	return nil
}

func (t *TokenTransferTransaction) Verify() error {
	// 1. Serialize the transaction without the signature
	originalSignature := t.Auth.OriginAuth.Signature
	t.Auth.OriginAuth.Signature = [RecoverableECDSASigLengthBytes]byte{} // Clear signature for serialization
	serializedTx, err := t.Serialize()
	if err != nil {
		return err
	}
	t.Auth.OriginAuth.Signature = originalSignature // Restore signature

	// 2. Calculate the initial sighash
	sighash := sha256.Sum256(serializedTx)

	// 3. Calculate the presign sighash
	var buf bytes.Buffer
	buf.WriteByte(byte(t.Auth.AuthType))
	binary.Write(&buf, binary.BigEndian, t.Auth.OriginAuth.Fee)
	binary.Write(&buf, binary.BigEndian, t.Auth.OriginAuth.Nonce)
	buf.Write(sighash[:])
	presignSighash := sha256.Sum256(buf.Bytes())

	// 4. Recover the public key from the signature
	pubKey, _, err := ecdsa.RecoverCompact(t.Auth.OriginAuth.Signature[:], presignSighash[:])
	if err != nil {
		return err
	}

	// 5. Verify that the recovered public key matches the signer
	hash160 := Hash160(pubKey.SerializeCompressed())
	if !bytes.Equal(hash160[:], t.Auth.OriginAuth.Signer[:]) {
		return errors.New("recovered public key does not match signer")
	}

	return nil
}

// func RecoverPublicKey(tx *TokenTransferTransaction) ([]byte, error) {
// 	// 1. Serialize the transaction without the signature
// 	originalSignature := tx.Auth.OriginAuth.Signature
// 	tx.Auth.OriginAuth.Signature = [65]byte{} // Clear signature for serialization
// 	serializedTx, err := tx.Serialize()
// 	if err != nil {
// 		return nil, err
// 	}
// 	tx.Auth.OriginAuth.Signature = originalSignature // Restore signature

// 	// 2. Calculate the initial sighash
// 	sighash := sha256.Sum256(serializedTx)

// 	// 3. Calculate the presign sighash
// 	var buf []byte
// 	buf = append(buf, byte(tx.Auth.AuthType))
// 	feeBytes := make([]byte, 8)
// 	binary.BigEndian.PutUint64(feeBytes, tx.Auth.OriginAuth.Fee)
// 	buf = append(buf, feeBytes...)
// 	nonceBytes := make([]byte, 8)
// 	binary.BigEndian.PutUint64(nonceBytes, tx.Auth.OriginAuth.Nonce)
// 	buf = append(buf, nonceBytes...)
// 	buf = append(buf, sighash[:]...)
// 	presignSighash := sha256.Sum256(buf)

// 	// 4. Recover the public key from the signature
// 	// The signature in Stacks includes the recovery ID as the last byte
// 	compactSignature := make([]byte, 65)
// 	copy(compactSignature[:64], tx.Auth.OriginAuth.Signature[:64])
// 	compactSignature[64] = tx.Auth.OriginAuth.Signature[64] - 27 // Adjust recovery ID

// 	pubKey, _, err := ecdsa.RecoverCompact(compactSignature, presignSighash[:])
// 	if err != nil {
// 		return nil, errors.New("failed to recover public key: " + err.Error())
// 	}

// 	return pubKey.SerializeCompressed(), nil
// }
