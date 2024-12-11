package keys

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

type SignatureScheme int

const (
	Secp256k1 SignatureScheme = iota
)

func Sha256Hash(msg []byte) []byte {
	hash := sha256.Sum256(msg)
	return hash[:]
}

func Keccak256Hash(msg []byte) []byte {
	return ethcrypto.Keccak256Hash(msg).Bytes()
}

type PublicKey []byte

func (pbk PublicKey) String() string {
	return hex.EncodeToString(pbk)
}

type KeyPair interface {
	Scheme() SignatureScheme
	PrivateKey() []byte
	PublicKey() PublicKey
	// Returns a compact signature with a format:
	// <1-byte compact sig recovery code><32-byte R><32-byte S>
	Sign(msg []byte) (signature []byte)
	VerifySignature(msg []byte, signature []byte) (pubkey []byte, err error)
}

type secp256k1Keypair struct {
	pk *btcec.PrivateKey
}

func NewKeyPairFromPrivateKeyBytes(scheme SignatureScheme, privateKeyBytes []byte) (KeyPair, error) {
	switch scheme {
	case Secp256k1:
		priv, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
		if priv == nil {
			return nil, fmt.Errorf("invalid privatekey bytes")
		}
		return &secp256k1Keypair{pk: priv}, nil

	default:
		return nil, fmt.Errorf("unsupported signature scheme")
	}
}

func NewKeyPair(scheme SignatureScheme) (KeyPair, error) {
	switch scheme {
	case Secp256k1:
		priv, err := btcec.NewPrivateKey()
		if err != nil {
			return nil, err
		}
		return &secp256k1Keypair{pk: priv}, nil

	default:
		return nil, fmt.Errorf("unsupported signature scheme")
	}
}

func (pair *secp256k1Keypair) PrivateKey() []byte {
	return pair.pk.Serialize()
}

func (pair *secp256k1Keypair) PublicKey() PublicKey {
	return PublicKey(pair.pk.PubKey().SerializeUncompressed())
}

func (pair *secp256k1Keypair) Scheme() SignatureScheme {
	return Secp256k1
}

func (pair *secp256k1Keypair) Sign(msg []byte) []byte {
	return ecdsa.SignCompact(pair.pk, msg, false)
}

func (pair *secp256k1Keypair) VerifySignature(msg []byte, signature []byte) (pubkey []byte, err error) {
	publicKey, _, err := ecdsa.RecoverCompact(signature, msg)
	if err != nil {
		return nil, err
	}

	if !slices.Equal(pair.PublicKey(), publicKey.SerializeUncompressed()) {
		return nil, fmt.Errorf("signature verification failed: wrong signer")
	}

	return publicKey.SerializeUncompressed(), nil
}

func VerifySignaturee(msg []byte, signature []byte) error {
	publicKey, _, err := ecdsa.RecoverCompact(signature, msg)
	if err != nil {
		return err
	}

	fmt.Println("PUbKey: ", hex.EncodeToString(publicKey.SerializeUncompressed()))
	return nil
}
