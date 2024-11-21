package stacks

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/mocks"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestRestoreKeystore_Success(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "keystore_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	encryptedKey := []byte("encrypted_private_key")
	address := "ST1PQHQKV0RJXZFY1DGX8MNSNYVE3VGZJSRTPGZGM"
	nid := "stacks"
	keystorePath := filepath.Join(tempDir, "keystore", nid, address)

	err = os.MkdirAll(filepath.Dir(keystorePath), 0700)
	assert.NoError(t, err)

	err = os.WriteFile(keystorePath, encryptedKey, 0600)
	assert.NoError(t, err)

	mockKMS := new(mocks.MockKMS)
	decryptedKey := []byte("decrypted_private_key")
	mockKMS.On("Decrypt", mock.Anything, encryptedKey).Return(decryptedKey, nil)

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	cfg := &Config{
		CommonConfig: provider.CommonConfig{
			HomeDir: tempDir,
			NID:     nid,
			Address: address,
		},
	}

	provider := &Provider{
		cfg:        cfg,
		log:        logger,
		kms:        mockKMS,
		privateKey: nil,
	}

	err = provider.RestoreKeystore(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, decryptedKey, provider.privateKey)
	mockKMS.AssertExpectations(t)
}

func TestRestoreKeystore_ReadFileError(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "keystore_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	address := "ST1PQHQKV0RJXZFY1DGX8MNSNYVE3VGZJSRTPGZGM"
	nid := "stacks"
	mockKMS := new(mocks.MockKMS)

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	cfg := &Config{
		CommonConfig: provider.CommonConfig{
			HomeDir: tempDir,
			NID:     nid,
			Address: address,
		},
	}

	provider := &Provider{
		cfg:        cfg,
		log:        logger,
		kms:        mockKMS,
		privateKey: nil,
	}

	err = provider.RestoreKeystore(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read keystore file")

	mockKMS.AssertNotCalled(t, "Decrypt", mock.Anything, mock.Anything)
}

func TestRestoreKeystore_DecryptError(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "keystore_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	encryptedKey := []byte("encrypted_private_key")
	address := "ST1PQHQKV0RJXZFY1DGX8MNSNYVE3VGZJSRTPGZGM"
	nid := "stacks"
	keystorePath := filepath.Join(tempDir, "keystore", nid, address)

	err = os.MkdirAll(filepath.Dir(keystorePath), 0700)
	assert.NoError(t, err)

	err = os.WriteFile(keystorePath, encryptedKey, 0600)
	assert.NoError(t, err)

	mockKMS := new(mocks.MockKMS)
	mockKMS.On("Decrypt", mock.Anything, encryptedKey).Return([]byte(nil), errors.New("decryption failed"))

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	cfg := &Config{
		CommonConfig: provider.CommonConfig{
			HomeDir: tempDir,
			NID:     nid,
			Address: address,
		},
	}

	provider := &Provider{
		cfg:        cfg,
		log:        logger,
		kms:        mockKMS,
		privateKey: nil,
	}

	err = provider.RestoreKeystore(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decrypt keystore")

	mockKMS.AssertExpectations(t)
}
