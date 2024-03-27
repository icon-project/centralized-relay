package sui

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"path"

	"github.com/block-vision/sui-go-sdk/common/sui_error"
	"github.com/block-vision/sui-go-sdk/models"
	"golang.org/x/crypto/blake2b"
)

type KeyPair byte

const (
	Ed25519Flag            KeyPair = 0
	Secp256k1Flag          KeyPair = 1
	suiAddressLength               = 64
	ed25519PublicKeyLength         = 32
)

func encodeBase64(value []byte) string {
	return base64.StdEncoding.EncodeToString(value)
}

func fromPublicKeyBytesToAddress(publicKey []byte, scheme byte) string {
	if scheme != byte(Ed25519Flag) && scheme != byte(Secp256k1Flag) {
		return ""
	}
	tmp := []byte{scheme}
	tmp = append(tmp, publicKey...)
	hexHash := blake2b.Sum256(tmp)
	return "0x" + hex.EncodeToString(hexHash[:])[:suiAddressLength]
}

// fetches Ed25519 keypair from the privatekey with flag
func fetchKeyPair(privateKey []byte) (models.SuiKeyPair, error) {
	switch privateKey[0] {
	case byte(Ed25519Flag):
		privKey := ed25519.NewKeyFromSeed(privateKey[1:])
		publicKey := privKey.Public().(ed25519.PublicKey)
		sk := privKey[:ed25519PublicKeyLength]
		pbInBase64 := encodeBase64(publicKey)
		return models.SuiKeyPair{
			Flag:            byte(Ed25519Flag),
			PrivateKey:      sk,
			PublicKeyBase64: pbInBase64,
			PublicKey:       publicKey,
			Address:         fromPublicKeyBytesToAddress(publicKey, byte(Ed25519Flag)),
		}, nil
	default:
		return models.SuiKeyPair{}, sui_error.ErrUnknownSignatureScheme
	}
}

// Restores the addres configured
func (p *Provider) RestoreKeystore(ctx context.Context) error {
	path := p.keystorePath(p.cfg.Address)
	keystore, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	privateKey, err := p.kms.Decrypt(ctx, keystore)
	if err != nil {
		return err
	}
	p.wallet, err = fetchKeyPair(privateKey)
	if err != nil {
		return err
	}
	return nil
}

// Creates new Ed25519 key
func (p *Provider) NewKeystore(password string) (string, error) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", err
	}
	privateKeyWithFlag := append([]byte{byte(Ed25519Flag)}, priv[:ed25519PublicKeyLength]...)
	keyPair, err := fetchKeyPair(privateKeyWithFlag)
	if err != nil {
		return "", err
	}
	keyStoreContent, err := p.kms.Encrypt(context.Background(), privateKeyWithFlag)
	if err != nil {
		return "", err
	}
	passphraseCipher, err := p.kms.Encrypt(context.Background(), []byte(password))
	if err != nil {
		return "", err
	}
	path := p.keystorePath(keyPair.Address)
	if err = os.WriteFile(path, keyStoreContent, 0o644); err != nil {
		return "", err
	}
	if err = os.WriteFile(path+".pass", passphraseCipher, 0o644); err != nil {
		return "", err
	}

	return keyPair.Address, nil
}

// Imports first ed25519 key pair from the keystore
func (p *Provider) ImportKeystore(ctx context.Context, keyPath, passphrase string) (string, error) {
	privFile, err := os.ReadFile(keyPath)
	if err != nil {
		return "", err
	}
	var ksData []string
	err = json.Unmarshal(privFile, &ksData)
	if err != nil {
		return "", err
	}
	//decode base64 for first key
	privateKey, err := base64.StdEncoding.DecodeString(string(ksData[0]))
	if err != nil {
		return "", err
	}
	keyPair, err := fetchKeyPair(privateKey)
	if err != nil {
		return "", err
	}
	keyStoreContent, err := p.kms.Encrypt(ctx, privateKey)
	if err != nil {
		return "", err
	}
	passphraseCipher, err := p.kms.Encrypt(ctx, []byte(passphrase))
	if err != nil {
		return "", err
	}
	path := p.keystorePath(keyPair.Address)
	if err = os.WriteFile(path, keyStoreContent, 0o644); err != nil {
		return "", err
	}
	if err = os.WriteFile(path+".pass", passphraseCipher, 0o644); err != nil {
		return "", err
	}
	return keyPair.Address, nil
}

// keystorePath is the path to the keystore file
func (p *Provider) keystorePath(addr string) string {
	return path.Join(p.cfg.HomeDir, "keystore", p.NID(), addr)
}
