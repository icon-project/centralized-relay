package keys

import (
	"context"
	"os"
	"path/filepath"

	"github.com/icon-project/centralized-relay/relayer/kms"
)

func LoadKeypairFromFile(keypath string, scheme SignatureScheme, kms kms.KMS) (KeyPair, error) {
	cipherKey, err := os.ReadFile(keypath)
	if err != nil {
		return nil, err
	}
	privateKey, err := kms.Decrypt(context.Background(), cipherKey)
	if err != nil {
		return nil, err
	}

	return NewKeyPairFromPrivateKeyBytes(scheme, privateKey)
}

func GetClusterKeyPath(homedir string, pubkey string) string {
	return filepath.Join(GetClusterKeyDir(homedir), pubkey)
}

func GetClusterKeyDir(homedir string) string {
	return filepath.Join(homedir, "keystore", "cluster")
}
