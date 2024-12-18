package stacks

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/icon-project/stacks-go-sdk/pkg/crypto"
	"github.com/icon-project/stacks-go-sdk/pkg/stacks"
	"go.uber.org/zap"
)

func (p *Provider) RestoreKeystore(ctx context.Context) error {
	keystorePath := p.keystorePath(p.cfg.Address)

	encryptedKey, err := os.ReadFile(keystorePath)
	if err != nil {
		p.log.Error("Failed to read keystore file", zap.String("path", keystorePath), zap.Error(err))
		return fmt.Errorf("failed to read keystore file: %w", err)
	}

	privateKey, err := p.kms.Decrypt(ctx, encryptedKey)
	if err != nil {
		p.log.Error("Failed to decrypt keystore", zap.Error(err))
		return fmt.Errorf("failed to decrypt keystore: %w", err)
	}

	p.privateKey = privateKey

	p.log.Info("Keystore restored successfully", zap.String("address", p.cfg.Address))
	return nil
}

func (p *Provider) NewKeystore(passphrase string) (string, error) {
	newKey, err := btcec.NewPrivateKey()
	if err != nil {
		p.log.Error("Failed to generate new private key", zap.Error(err))
		return "", fmt.Errorf("failed to generate new private key: %w", err)
	}
	privateKeyBytes := append(newKey.Serialize(), byte(0x01))

	network, err := MapNIDToChainID(p.cfg.NID)
	if err != nil {
		p.log.Error("Chain id not found. Check the NID config", zap.Error(err))
		return "", fmt.Errorf("chain id not found: %w", err)
	}

	address, err := crypto.GetAddressFromPrivateKey(privateKeyBytes, network)
	if err != nil {
		p.log.Error("Failed to derive address from private key", zap.Error(err))
		return "", fmt.Errorf("failed to derive address: %w", err)
	}

	encryptedKey, err := p.kms.Encrypt(context.Background(), privateKeyBytes)
	if err != nil {
		p.log.Error("Failed to encrypt private key", zap.Error(err))
		return "", fmt.Errorf("failed to encrypt private key: %w", err)
	}

	encryptedPass, err := p.kms.Encrypt(context.Background(), []byte(passphrase))
	if err != nil {
		p.log.Error("Failed to encrypt passphrase", zap.Error(err))
		return "", fmt.Errorf("failed to encrypt passphrase: %w", err)
	}

	keystorePath := p.keystorePath(address)
	passPath := keystorePath + ".pass"

	if err := os.MkdirAll(filepath.Dir(keystorePath), 0700); err != nil {
		p.log.Error("Failed to create keystore directory", zap.Error(err))
		return "", fmt.Errorf("failed to create keystore directory: %w", err)
	}

	if err := os.WriteFile(keystorePath, encryptedKey, 0600); err != nil {
		p.log.Error("Failed to write keystore file", zap.String("path", keystorePath), zap.Error(err))
		return "", fmt.Errorf("failed to write keystore file: %w", err)
	}

	if err := os.WriteFile(passPath, encryptedPass, 0600); err != nil {
		p.log.Error("Failed to write passphrase file", zap.String("path", passPath), zap.Error(err))
		return "", fmt.Errorf("failed to write passphrase file: %w", err)
	}

	p.privateKey = privateKeyBytes
	p.cfg.Address = address

	p.log.Info("New keystore created successfully", zap.String("address", address))
	return address, nil
}

func (p *Provider) ImportKeystore(ctx context.Context, keyPath, passphrase string) (string, error) {
	encryptedKey, err := os.ReadFile(keyPath)
	if err != nil {
		p.log.Error("Failed to read imported keystore file", zap.String("path", keyPath), zap.Error(err))
		return "", fmt.Errorf("failed to read imported keystore file: %w", err)
	}

	privateKey, err := p.kms.Decrypt(ctx, encryptedKey)
	if err != nil {
		p.log.Error("Failed to decrypt imported keystore", zap.Error(err))
		return "", fmt.Errorf("failed to decrypt imported keystore: %w", err)
	}

	network, err := MapNIDToChainID(p.cfg.NID)
	if err != nil {
		p.log.Error("Chain id not found. Check the NID config", zap.Error(err))
		return "", fmt.Errorf("chain id not found: %w", err)
	}

	address, err := crypto.GetAddressFromPrivateKey(privateKey, network)
	if err != nil {
		p.log.Error("Failed to derive address from imported private key", zap.Error(err))
		return "", fmt.Errorf("failed to derive address: %w", err)
	}

	encryptedStoredKey, err := p.kms.Encrypt(ctx, privateKey)
	if err != nil {
		p.log.Error("Failed to encrypt imported private key for storage", zap.Error(err))
		return "", fmt.Errorf("failed to encrypt imported private key: %w", err)
	}

	encryptedPass, err := p.kms.Encrypt(ctx, []byte(passphrase))
	if err != nil {
		p.log.Error("Failed to encrypt passphrase", zap.Error(err))
		return "", fmt.Errorf("failed to encrypt passphrase: %w", err)
	}

	destKeystorePath := p.keystorePath(address)
	destPassPath := destKeystorePath + ".pass"

	if err := os.MkdirAll(filepath.Dir(destKeystorePath), 0700); err != nil {
		p.log.Error("Failed to create keystore directory", zap.Error(err))
		return "", fmt.Errorf("failed to create keystore directory: %w", err)
	}

	if err := os.WriteFile(destKeystorePath, encryptedStoredKey, 0600); err != nil {
		p.log.Error("Failed to write imported keystore file", zap.String("path", destKeystorePath), zap.Error(err))
		return "", fmt.Errorf("failed to write imported keystore file: %w", err)
	}

	if err := os.WriteFile(destPassPath, encryptedPass, 0600); err != nil {
		p.log.Error("Failed to write imported passphrase file", zap.String("path", destPassPath), zap.Error(err))
		return "", fmt.Errorf("failed to write imported passphrase file: %w", err)
	}

	p.privateKey = privateKey
	p.cfg.Address = address

	p.log.Info("Keystore imported successfully", zap.String("address", address))
	return address, nil
}

func (p *Provider) keystorePath(addr string) string {
	return filepath.Join(p.cfg.HomeDir, "keystore", p.NID(), addr)
}

func MapNIDToChainID(nid string) (stacks.ChainID, error) {
	switch nid {
	case "stacks":
		return stacks.ChainIDMainnet, nil
	case "stacks_testnet":
		return stacks.ChainIDTestnet, nil
	default:
		return 0, fmt.Errorf("unsupported NID: %s", nid)
	}
}
