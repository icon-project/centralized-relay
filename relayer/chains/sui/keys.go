package sui

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/cometbft/cometbft/libs/rand"
	"github.com/coming-chat/go-sui/v2/account"
	"github.com/coming-chat/go-sui/v2/sui_types"
)

type KeyPair byte

const (
	Ed25519Flag            KeyPair = 0
	Secp256k1Flag          KeyPair = 1
	suiAddressLength               = 64
	ed25519PublicKeyLength         = 32
)

// fetches Ed25519 keypair from the privatekey with flag
func fetchKeyPair(privateKey string) (*account.Account, error) {
	return account.NewAccountWithKeystore(privateKey)
}

// Restores the addres configured
func (p *Provider) RestoreKeystore(ctx context.Context) error {
	path := p.keystorePath(p.cfg.Address)
	keystore, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error restoring account: %w", err)
	}
	privateKey, err := p.kms.Decrypt(ctx, keystore)
	if err != nil {
		return fmt.Errorf("error restoring account: %w", err)
	}
	p.wallet, err = fetchKeyPair(string(privateKey))
	if err != nil {
		return fmt.Errorf("error restoring account: %w", err)
	}
	return nil
}

// Creates new Ed25519 key
func (p *Provider) NewKeystore(password string) (string, error) {
	signatureScheme, err := sui_types.NewSignatureScheme(byte(Ed25519Flag))
	if err != nil {
		return "", err
	}
	account := account.NewAccount(signatureScheme, rand.Bytes(32))
	privateKeyWithFlag := append([]byte{byte(Ed25519Flag)}, account.KeyPair.PrivateKey()[:ed25519PublicKeyLength]...)
	encodedPkey := base64.StdEncoding.EncodeToString([]byte(privateKeyWithFlag))
	keyStoreContent, err := p.kms.Encrypt(context.Background(), []byte(encodedPkey))
	if err != nil {
		return "", fmt.Errorf("error adding new account: %w", err)
	}
	passphraseCipher, err := p.kms.Encrypt(context.Background(), []byte(password))
	if err != nil {
		return "", fmt.Errorf("error adding new account: %w", err)
	}
	path := p.keystorePath(account.Address)
	if err = os.WriteFile(path, keyStoreContent, 0o644); err != nil {
		return "", fmt.Errorf("error adding new account: %w", err)
	}
	if err = os.WriteFile(path+".pass", passphraseCipher, 0o644); err != nil {
		return "", fmt.Errorf("error adding new account: %w", err)
	}
	return account.Address, nil
}

// Imports first ed25519 key pair from the keystore
func (p *Provider) ImportKeystore(ctx context.Context, keyPath, passphrase string) (string, error) {
	privFile, err := os.ReadFile(keyPath)
	if err != nil {
		return "", fmt.Errorf("error importing key: %w", err)
	}
	var ksData []string
	err = json.Unmarshal(privFile, &ksData)
	if err != nil {
		return "", fmt.Errorf("error importing key: %w", err)
	}
	// decode base64 for first key
	firstKeyPairString := string(ksData[0])
	keyPair, err := fetchKeyPair(firstKeyPairString)
	if err != nil {
		return "", fmt.Errorf("error importing key: %w", err)
	}
	keyStoreContent, err := p.kms.Encrypt(ctx, []byte(firstKeyPairString))
	if err != nil {
		return "", fmt.Errorf("error importing key: %w", err)
	}
	passphraseCipher, err := p.kms.Encrypt(ctx, []byte(passphrase))
	if err != nil {
		return "", fmt.Errorf("error importing key: %w", err)
	}
	path := p.keystorePath(keyPair.Address)
	if err = os.WriteFile(path, keyStoreContent, 0o644); err != nil {
		return "", fmt.Errorf("error importing key: %w", err)
	}
	if err = os.WriteFile(path+".pass", passphraseCipher, 0o644); err != nil {
		return "", fmt.Errorf("error importing key: %w", err)
	}
	return keyPair.Address, nil
}

// keystorePath is the path to the keystore file
func (p *Provider) keystorePath(addr string) string {
	return path.Join(p.cfg.HomeDir, "keystore", p.NID(), addr)
}
