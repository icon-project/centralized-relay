package keystore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/scrypt"
)

const (
	kdf  = "scrypt"
	algo = "aes-256-gcm"
)

type EncryptedKey struct {
	Address    string `json:"address"`
	Ciphertext string `json:"ciphertext"`
	Salt       string `json:"salt"`
	IV         string `json:"iv"`
	KDF        string `json:"kdf"`
	Algorithm  string `json:"algorithm"`
}

func EncryptToJSONKeystore(privateKey []byte, address string, password string) ([]byte, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	// Derive a key using scrypt
	key, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	// Encrypt the private key using AES-GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate a random IV
	iv := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, iv, []byte(privateKey), nil)

	// Prepare the JSON data
	encryptedKey := EncryptedKey{
		Address:    address,
		Ciphertext: base64.URLEncoding.EncodeToString(ciphertext),
		Salt:       base64.URLEncoding.EncodeToString(salt),
		IV:         base64.URLEncoding.EncodeToString(iv),
		KDF:        kdf,
		Algorithm:  algo,
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(encryptedKey, "", "  ")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func DecryptFromJSONKeystore(encryptedJSONBytes []byte, password string) ([]byte, string, error) {
	var encryptedKey EncryptedKey
	if err := json.Unmarshal(encryptedJSONBytes, &encryptedKey); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Decode base64 values
	salt, err := base64.URLEncoding.DecodeString(encryptedKey.Salt)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode salt: %w", err)
	}

	iv, err := base64.URLEncoding.DecodeString(encryptedKey.IV)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode IV: %w", err)
	}

	ciphertext, err := base64.URLEncoding.DecodeString(encryptedKey.Ciphertext)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	// Derive the key using scrypt
	key, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, 32)
	if err != nil {
		return nil, "", fmt.Errorf("failed to derive key: %w", err)
	}

	// Create the AES-GCM cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt the ciphertext
	privateKey, err := aesGCM.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decrypt private key: %w", err)
	}

	return privateKey, encryptedKey.Address, nil
}
