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
	utilKeystore "github.com/icon-project/centralized-relay/utils/keystore"
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
	kmsEncryptedKeystore, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	keystoreData, err := p.kms.Decrypt(ctx, kmsEncryptedKeystore)
	if err != nil {
		return err
	}

	passFile, err := os.ReadFile(path + ".pass")
	if err != nil {
		return err
	}
	pass, err := p.kms.Decrypt(ctx, passFile)
	if err != nil {
		return err
	}

	privateKey, _, err := utilKeystore.DecryptFromJSONKeystore(keystoreData, string(pass))
	if err != nil {
		return err
	}

	p.wallet, err = fetchKeyPair(string(privateKey))
	if err != nil {
		return err
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

	keystoreData, err := utilKeystore.EncryptToJSONKeystore(privateKeyWithFlag, account.Address, password)
	if err != nil {
		return "", err
	}

	keyStoreContent, err := p.kms.Encrypt(context.Background(), keystoreData)
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
	jsonKeystoreData, err := os.ReadFile(keyPath)
	if err != nil {
		return "", fmt.Errorf("error importing key: %w", err)
	}

	encrypedKeystoreBytes, addr, err := utilKeystore.DecryptFromJSONKeystore(jsonKeystoreData, passphrase)
	if err != nil {
		return "", err
	}

	keyStoreContent, err := p.kms.Encrypt(ctx, encrypedKeystoreBytes)
	if err != nil {
		return "", fmt.Errorf("error importing key: %w", err)
	}
	passphraseCipher, err := p.kms.Encrypt(ctx, []byte(passphrase))
	if err != nil {
		return "", fmt.Errorf("error importing key: %w", err)
	}
	path := p.keystorePath(addr)
	if err = os.WriteFile(path, keyStoreContent, 0o644); err != nil {
		return "", fmt.Errorf("error importing key: %w", err)
	}
	if err = os.WriteFile(path+".pass", passphraseCipher, 0o644); err != nil {
		return "", fmt.Errorf("error importing key: %w", err)
	}
	return addr, nil
}

// keystorePath is the path to the keystore file
func (p *Provider) keystorePath(addr string) string {
	return path.Join(p.cfg.HomeDir, "keystore", p.NID(), addr)
}

// Convert private to keystore
func (p *Provider) ConvertPrivateKey(ctx context.Context, keyPath, passphrase string) (string, error) {
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
	firstKeyPairString := string(ksData[0])
	keyPair, err := fetchKeyPair(firstKeyPairString)
	if err != nil {
		return "", err
	}

	pkey, err := base64.StdEncoding.DecodeString(firstKeyPairString)
	if err != nil {
		return "", err
	}

	jsonKeystoreBytes, err := utilKeystore.EncryptToJSONKeystore(pkey, keyPair.Address, passphrase)
	if err != nil {
		return "", err
	}

	path := p.keystorePath(keyPair.Address) + "_keystore.json"
	if err = os.WriteFile(path, jsonKeystoreBytes, 0o644); err != nil {
		return "", err
	}

	return path, nil
}
