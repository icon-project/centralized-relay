package wasm

import (
	"context"
	"os"
	"path"
)

func (p *Provider) RestoreKeystore(ctx context.Context) error {
	filePath := path.Join(p.cfg.HomeDir, "keystore", p.NID(), p.cfg.Address)
	privFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	priv, err := p.kms.Decrypt(ctx, privFile)
	if err != nil {
		return err
	}
	passFile, err := os.ReadFile(filePath + ".pass")
	if err != nil {
		return err
	}
	pass, err := p.kms.Decrypt(ctx, passFile)
	if err != nil {
		return err
	}
	if err := p.client.LoadArmor(ctx, p.NID(), string(priv), string(pass)); err != nil {
		return err
	}
	return nil
}

func (p *Provider) NewKeystore(passphrase string) (string, error) {
	armor, addr, err := p.client.CreateAccount(p.NID(), passphrase)
	if err != nil {
		return "", err
	}
	encryptedArmor, err := p.kms.Encrypt(context.Background(), []byte(armor))
	if err != nil {
		return "", err
	}
	filePath := path.Join(p.cfg.HomeDir, "keystore", p.NID(), addr)
	if err = os.WriteFile(filePath, encryptedArmor, 0o644); err != nil {
		return "", err
	}
	encryptedPassphrase, err := p.kms.Encrypt(context.Background(), []byte(passphrase))
	if err != nil {
		return "", err
	}
	if err = os.WriteFile(filePath+".pass", encryptedPassphrase, 0o644); err != nil {
		return "", err
	}
	return addr, nil
}
