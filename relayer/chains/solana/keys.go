package solana

import (
	"context"
	"os"
	"path"

	"github.com/gagliardetto/solana-go"
)

func (p *Provider) RestoreKeystore(ctx context.Context) error {
	encryptedPrivateKey, err := os.ReadFile(p.keystorePath(p.cfg.Address))
	if err != nil {
		return err
	}

	rawPrivateKey, err := p.kms.Decrypt(ctx, encryptedPrivateKey)
	if err != nil {
		return err
	}

	wallet := solana.Wallet{
		PrivateKey: rawPrivateKey,
	}

	p.wallet = &wallet

	return nil
}

func (p *Provider) NewKeystore(password string) (string, error) {
	wallet := solana.NewWallet()

	privateKeyEncrypted, err := p.kms.Encrypt(context.Background(), wallet.PrivateKey)
	if err != nil {
		return "", err
	}

	passphraseEncrypted, err := p.kms.Encrypt(context.Background(), []byte(password))
	if err != nil {
		return "", err
	}

	walletAddr := wallet.PublicKey().String()

	path := p.keystorePath(walletAddr)

	if err = os.WriteFile(path, privateKeyEncrypted, 0o644); err != nil {
		return "", err
	}
	if err = os.WriteFile(path+".pass", passphraseEncrypted, 0o644); err != nil {
		return "", err
	}

	return walletAddr, nil
}

func (p *Provider) ImportKeystore(ctx context.Context, keyPath, passphrase string) (string, error) {
	privateKey, err := solana.PrivateKeyFromSolanaKeygenFile(keyPath)
	if err != nil {
		return "", err
	}

	privateKeyEncrypted, err := p.kms.Encrypt(ctx, privateKey)
	if err != nil {
		return "", err
	}

	passphraseEncrypted, err := p.kms.Encrypt(ctx, []byte(passphrase))
	if err != nil {
		return "", err
	}

	wallet := solana.Wallet{
		PrivateKey: privateKey,
	}

	walletAddr := wallet.PublicKey().String()

	path := p.keystorePath(walletAddr)

	if err = os.WriteFile(path, privateKeyEncrypted, 0o644); err != nil {
		return "", err
	}
	if err = os.WriteFile(path+".pass", passphraseEncrypted, 0o644); err != nil {
		return "", err
	}

	return walletAddr, nil
}

func (p *Provider) keystorePath(addr string) string {
	return path.Join(p.cfg.HomeDir, "keystore", p.NID(), addr)
}
