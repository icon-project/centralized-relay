package stacks

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/mocks"
	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/provider"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/stacks-go-sdk/pkg/clarity"
	blockchainApiClient "github.com/icon-project/stacks-go-sdk/pkg/stacks_blockchain_api_client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupTestProvider(t *testing.T) (*Provider, *mocks.MockClient) {
	logger, _ := zap.NewDevelopment()
	mockClient := new(mocks.MockClient)

	cfg := &Config{
		CommonConfig: provider.CommonConfig{
			RPCUrl: "https://stacks-node-api.testnet.stacks.co",
			Contracts: providerTypes.ContractConfigMap{
				providerTypes.XcallContract:      "ST15C893XJFJ6FSKM020P9JQDB5T7X6MQTXMBPAVH.xcall-proxy",
				providerTypes.ConnectionContract: "ST15C893XJFJ6FSKM020P9JQDB5T7X6MQTXMBPAVH.centralized-connection",
			},
			NID:       "stacks_testnet",
			ChainName: "stacks_testnet",
		},
	}

	p, err := cfg.NewProvider(context.Background(), logger, "/tmp/relayer", false, "stacks_testnet")
	assert.NoError(t, err)

	provider := p.(*Provider)
	provider.client = mockClient

	return provider, mockClient
}

func TestProvider_Init(t *testing.T) {
	provider, _ := setupTestProvider(t)
	mockKMS := new(mocks.MockKMS)

	err := provider.Init(context.Background(), "/tmp/relayer", mockKMS)
	assert.NoError(t, err)
	assert.Equal(t, mockKMS, provider.kms)
}

func TestProvider_QueryLatestHeight(t *testing.T) {
	provider, mockClient := setupTestProvider(t)

	expectedHeight := uint64(1234)
	mockBlock := &blockchainApiClient.GetBlocks200ResponseResultsInner{
		Height: int32(expectedHeight),
	}

	mockClient.On("GetLatestBlock", mock.Anything).Return(mockBlock, nil)

	height, err := provider.QueryLatestHeight(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedHeight, height)

	mockClient.AssertExpectations(t)
}

func TestProvider_QueryBalance(t *testing.T) {
	provider, mockClient := setupTestProvider(t)

	address := "ST1PQHQKV0RJXZFY1DGX8MNSNYVE3VGZJSRTPGZGM"
	expectedBalance := big.NewInt(1000000)

	mockClient.On("GetAccountBalance", mock.Anything, address).Return(expectedBalance, nil)

	balance, err := provider.QueryBalance(context.Background(), address)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1000000), balance.Amount)
	assert.Equal(t, "STX", balance.Denom)

	mockClient.AssertExpectations(t)
}

func TestProvider_GetFee(t *testing.T) {
	provider, mockClient := setupTestProvider(t)

	networkID := "icon"
	expectedFee := uint64(1000)

	mockClient.On("GetFee", mock.Anything, provider.cfg.Contracts[providerTypes.ConnectionContract], networkID, true).
		Return(expectedFee, nil)

	fee, err := provider.GetFee(context.Background(), networkID, true)
	assert.NoError(t, err)
	assert.Equal(t, expectedFee, fee)

	mockClient.AssertExpectations(t)
}

func TestProvider_SetFee(t *testing.T) {
	provider, mockClient := setupTestProvider(t)

	networkID := "icon"
	msgFee := big.NewInt(1000)
	resFee := big.NewInt(500)
	expectedTxID := "0x123456789"

	mockClient.On("SetFee",
		mock.Anything,
		provider.cfg.Contracts[providerTypes.ConnectionContract],
		networkID,
		msgFee,
		resFee,
		provider.cfg.Address,
		provider.privateKey,
	).Return(expectedTxID, nil)

	err := provider.SetFee(context.Background(), networkID, msgFee, resFee)
	assert.NoError(t, err)

	mockClient.AssertExpectations(t)
}

func TestProvider_ClaimFee(t *testing.T) {
	provider, mockClient := setupTestProvider(t)

	expectedTxID := "0x123456789"

	mockClient.On("ClaimFee",
		mock.Anything,
		provider.cfg.Contracts[providerTypes.ConnectionContract],
		provider.cfg.Address,
		provider.privateKey,
	).Return(expectedTxID, nil)

	err := provider.ClaimFee(context.Background())
	assert.NoError(t, err)

	mockClient.AssertExpectations(t)
}

func TestProvider_SetAdmin(t *testing.T) {
	provider, mockClient := setupTestProvider(t)

	newAdmin := "ST1PQHQKV0RJXZFY1DGX8MNSNYVE3VGZJSRTPGZGM"
	expectedTxID := "0x123456789"

	var hash160 [20]byte
	copy(hash160[:], []byte("ST15C893XJFJ6FSKM020P9JQDB5T7X6MQTXMBPAVH"))
	currentImpl, err := clarity.NewContractPrincipal(0, hash160, "xcall-impl")
	assert.NoError(t, err)

	mockClient.On("GetCurrentImplementation", mock.Anything, provider.cfg.Contracts[providerTypes.XcallContract]).
		Return(currentImpl, nil).Once()

	mockClient.On("SetAdmin",
		mock.Anything,
		provider.cfg.Contracts[providerTypes.XcallContract],
		newAdmin,
		currentImpl,
		provider.cfg.Address,
		provider.privateKey,
	).Return(expectedTxID, nil).Once()

	successStatus := "success"
	mockResponse := &blockchainApiClient.GetTransactionById200Response{
		GetTransactionList200ResponseResultsInner: &blockchainApiClient.GetTransactionList200ResponseResultsInner{
			ContractCallTransaction: &blockchainApiClient.ContractCallTransaction{
				TxId:        expectedTxID,
				BlockHeight: 1234,
				TxStatus: blockchainApiClient.TokenTransferTransactionTxStatus{
					String: &successStatus,
				},
				Canonical: true,
			},
		},
	}
	mockClient.On("GetTransactionById", mock.Anything, expectedTxID).
		Return(mockResponse, nil).Once()

	err = provider.SetAdmin(context.Background(), newAdmin)
	assert.NoError(t, err, "SetAdmin should execute without error")

	mockClient.AssertExpectations(t)
}

func TestProvider_MessageReceived(t *testing.T) {
	provider, mockClient := setupTestProvider(t)

	tests := []struct {
		name      string
		key       *providerTypes.MessageKey
		mockSetup func()
		expected  bool
		expectErr bool
	}{
		{
			name: "EmitMessage Success",
			key: &providerTypes.MessageKey{
				EventType: events.EmitMessage,
				Src:       "icon",
				Sn:        big.NewInt(12345),
			},
			mockSetup: func() {
				mockClient.On("GetReceipt",
					mock.Anything,
					provider.cfg.Contracts[providerTypes.ConnectionContract],
					"icon",
					big.NewInt(12345),
				).Return(true, nil).Once()
			},
			expected:  true,
			expectErr: false,
		},
		{
			name: "CallMessage Returns False",
			key: &providerTypes.MessageKey{
				EventType: events.CallMessage,
				Src:       "icon",
				Sn:        big.NewInt(12345),
			},
			mockSetup: func() {},
			expected:  false,
			expectErr: false,
		},
		{
			name: "RollbackMessage Returns False",
			key: &providerTypes.MessageKey{
				EventType: events.RollbackMessage,
				Src:       "icon",
				Sn:        big.NewInt(12345),
			},
			mockSetup: func() {},
			expected:  false,
			expectErr: false,
		},
		{
			name: "Unknown Event Type",
			key: &providerTypes.MessageKey{
				EventType: "unknown",
				Src:       "icon",
				Sn:        big.NewInt(12345),
			},
			mockSetup: func() {},
			expected:  true,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			received, err := provider.MessageReceived(context.Background(), tt.key)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, received)
			}
		})
	}

	mockClient.AssertExpectations(t)
}

func TestProvider_RestoreKeystore(t *testing.T) {
	provider, _ := setupTestProvider(t)
	mockKMS := new(mocks.MockKMS)
	provider.kms = mockKMS

	tempDir := t.TempDir()
	provider.cfg.HomeDir = tempDir
	provider.cfg.Address = "ST1PQHQKV0RJXZFY1DGX8MNSNYVE3VGZJSRTPGZGM"

	keystorePath := filepath.Join(tempDir, "keystore", provider.NID(), provider.cfg.Address)
	err := os.MkdirAll(filepath.Dir(keystorePath), 0700)
	assert.NoError(t, err)

	encryptedKey := []byte("encrypted_key_data")
	err = os.WriteFile(keystorePath, encryptedKey, 0600)
	assert.NoError(t, err)

	decryptedKey := []byte("decrypted_key_data")
	mockKMS.On("Decrypt", mock.Anything, encryptedKey).Return(decryptedKey, nil)

	err = provider.RestoreKeystore(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, decryptedKey, provider.privateKey)

	mockKMS.AssertExpectations(t)
}

func TestProvider_NewKeystore(t *testing.T) {
	provider, _ := setupTestProvider(t)
	mockKMS := new(mocks.MockKMS)
	provider.kms = mockKMS

	passphrase := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	tempDir := t.TempDir()
	provider.cfg.HomeDir = tempDir

	mockKMS.On("Encrypt", mock.Anything, mock.AnythingOfType("[]uint8")).
		Return([]byte("encrypted_key"), nil)
	mockKMS.On("Encrypt", mock.Anything, []byte(passphrase)).
		Return([]byte("encrypted_passphrase"), nil)

	address, err := provider.NewKeystore(passphrase)
	assert.NoError(t, err, "NewKeystore should not error")
	assert.NotEmpty(t, address, "NewKeystore should return a non-empty address")

	assert.Equal(t, 33, len(provider.privateKey), "Private key should be 33 bytes")
	t.Logf("PrivateKey length: %d", len(provider.privateKey))
	t.Logf("PrivateKey bytes: %x", provider.privateKey)

	keystorePath := filepath.Join(tempDir, "keystore", provider.NID(), address)
	_, err = os.Stat(keystorePath)
	assert.NoError(t, err, "Keystore file should exist")

	_, err = os.Stat(keystorePath + ".pass")
	assert.NoError(t, err, "Passphrase file should exist")

	mockKMS.AssertExpectations(t)
}

func TestProvider_Config(t *testing.T) {
	provider, _ := setupTestProvider(t)

	cfg := provider.Config()
	assert.NotNil(t, cfg)
	assert.IsType(t, &Config{}, cfg)

	typedCfg := cfg.(*Config)
	assert.Equal(t, provider.cfg, typedCfg)
}

func TestProvider_Type(t *testing.T) {
	provider, _ := setupTestProvider(t)

	providerType := provider.Type()
	assert.Equal(t, "stacks", providerType)
}

func TestProvider_NID(t *testing.T) {
	provider, _ := setupTestProvider(t)

	nid := provider.NID()
	assert.Equal(t, provider.cfg.NID, nid)
}

func TestProvider_Name(t *testing.T) {
	provider, _ := setupTestProvider(t)

	name := provider.Name()
	assert.Equal(t, provider.cfg.ChainName, name)
}

func TestProvider_FinalityBlock(t *testing.T) {
	provider, _ := setupTestProvider(t)
	provider.cfg.FinalityBlock = 10

	finality := provider.FinalityBlock(context.Background())
	assert.Equal(t, uint64(10), finality)
}

func TestProvider_GetLastSavedBlockHeight(t *testing.T) {
	provider, _ := setupTestProvider(t)

	expectedHeight := uint64(1234)
	provider.LastSavedHeightFunc = func() uint64 {
		return expectedHeight
	}

	height := provider.GetLastSavedBlockHeight()
	assert.Equal(t, expectedHeight, height)
}

func TestProvider_RevertMessage(t *testing.T) {
	provider, _ := setupTestProvider(t)

	err := provider.RevertMessage(context.Background(), big.NewInt(12345))
	assert.Error(t, err)
	assert.Equal(t, "not implemented", err.Error())
}

func TestProviderConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *Config
		expectErr bool
	}{
		{
			name: "Valid Config",
			cfg: &Config{
				CommonConfig: provider.CommonConfig{
					RPCUrl: "https://example.com",
					Contracts: providerTypes.ContractConfigMap{
						"XcallContract": "ST1234",
					},
				},
			},
			expectErr: false,
		},
		{
			name: "Empty RPC URL",
			cfg: &Config{
				CommonConfig: provider.CommonConfig{
					RPCUrl: "",
					Contracts: providerTypes.ContractConfigMap{
						"XcallContract": "ST1234",
					},
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProviderConfig_Enabled(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *Config
		expected bool
	}{
		{
			name: "Enabled Config",
			cfg: &Config{
				CommonConfig: provider.CommonConfig{
					Disabled: false,
				},
			},
			expected: true,
		},
		{
			name: "Disabled Config",
			cfg: &Config{
				CommonConfig: provider.CommonConfig{
					Disabled: true,
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enabled := tt.cfg.Enabled()
			assert.Equal(t, tt.expected, enabled)
		})
	}
}

func TestProviderConfig_SetWallet(t *testing.T) {
	cfg := &Config{}
	address := "ST1PQHQKV0RJXZFY1DGX8MNSNYVE3VGZJSRTPGZGM"

	cfg.SetWallet(address)
	assert.Equal(t, address, cfg.Address)
}

func TestProviderConfig_GetWallet(t *testing.T) {
	cfg := &Config{
		CommonConfig: provider.CommonConfig{
			Address: "ST1PQHQKV0RJXZFY1DGX8MNSNYVE3VGZJSRTPGZGM",
		},
	}

	address := cfg.GetWallet()
	assert.Equal(t, cfg.Address, address)
}

func TestProvider_QueryTransactionReceipt(t *testing.T) {
	provider, mockClient := setupTestProvider(t)

	txID := "0x123456789"

	t.Run("Mempool Transaction", func(t *testing.T) {
		pendingStatus := "pending"
		response := &blockchainApiClient.GetTransactionById200Response{
			GetMempoolTransactionList200ResponseResultsInner: &blockchainApiClient.GetMempoolTransactionList200ResponseResultsInner{
				ContractCallMempoolTransaction1: &blockchainApiClient.ContractCallMempoolTransaction1{
					TxId: txID,
					TxStatus: blockchainApiClient.TokenTransferMempoolTransaction1TxStatus{
						String: &pendingStatus,
					},
				},
			},
		}

		mockClient.On("GetTransactionById", mock.Anything, txID).Return(response, nil).Once()

		receipt, err := provider.QueryTransactionReceipt(context.Background(), txID)
		assert.NoError(t, err)
		assert.NotNil(t, receipt)
		assert.Equal(t, txID, receipt.TxHash)
		assert.Equal(t, uint64(0), receipt.Height)
	})

	t.Run("Confirmed Transaction", func(t *testing.T) {
		successStatus := "success"
		response := &blockchainApiClient.GetTransactionById200Response{
			GetTransactionList200ResponseResultsInner: &blockchainApiClient.GetTransactionList200ResponseResultsInner{
				ContractCallTransaction: &blockchainApiClient.ContractCallTransaction{
					TxId:        txID,
					BlockHeight: 1234,
					TxStatus: blockchainApiClient.TokenTransferTransactionTxStatus{
						String: &successStatus,
					},
					Canonical: true,
				},
			},
		}

		mockClient.On("GetTransactionById", mock.Anything, txID).Return(response, nil).Once()

		receipt, err := provider.QueryTransactionReceipt(context.Background(), txID)
		assert.NoError(t, err)
		assert.NotNil(t, receipt)
		assert.Equal(t, txID, receipt.TxHash)
		assert.Equal(t, uint64(1234), receipt.Height)
		assert.True(t, receipt.Status)
	})

	t.Run("Error Case", func(t *testing.T) {
		mockClient.On("GetTransactionById", mock.Anything, txID).
			Return((*blockchainApiClient.GetTransactionById200Response)(nil), fmt.Errorf("transaction not found")).Once()

		receipt, err := provider.QueryTransactionReceipt(context.Background(), txID)
		assert.Error(t, err)
		assert.Nil(t, receipt)
	})

	mockClient.AssertExpectations(t)
}

func TestProvider_ShouldReceiveMessage(t *testing.T) {
	provider, _ := setupTestProvider(t)

	message := &providerTypes.Message{
		Dst: "stacks_testnet",
		Src: "icon",
		Sn:  big.NewInt(12345),
	}

	should, err := provider.ShouldReceiveMessage(context.Background(), message)
	assert.NoError(t, err)
	assert.True(t, should)
}

func TestProvider_ShouldSendMessage(t *testing.T) {
	provider, _ := setupTestProvider(t)

	message := &providerTypes.Message{
		Dst: "icon",
		Src: "stacks_testnet",
		Sn:  big.NewInt(12345),
	}

	should, err := provider.ShouldSendMessage(context.Background(), message)
	assert.NoError(t, err)
	assert.True(t, should)
}

func TestProvider_WaitForTransactionConfirmation(t *testing.T) {
	provider, mockClient := setupTestProvider(t)

	txID := "0x123456789"

	t.Run("Successful Confirmation", func(t *testing.T) {
		successStatus := "success"
		response := &blockchainApiClient.GetTransactionById200Response{
			GetTransactionList200ResponseResultsInner: &blockchainApiClient.GetTransactionList200ResponseResultsInner{
				ContractCallTransaction: &blockchainApiClient.ContractCallTransaction{
					TxId:        txID,
					BlockHeight: 1234,
					TxStatus: blockchainApiClient.TokenTransferTransactionTxStatus{
						String: &successStatus,
					},
					Canonical: true,
				},
			},
		}

		mockClient.On("GetTransactionById", mock.Anything, txID).Return(response, nil)

		receipt, err := provider.waitForTransactionConfirmation(context.Background(), txID, MAX_WAIT_TIME)
		assert.NoError(t, err)
		assert.NotNil(t, receipt)
		assert.Equal(t, txID, receipt.TxHash)
		assert.True(t, receipt.Status)
		assert.Equal(t, uint64(1234), receipt.Height)
	})

	t.Run("Context Cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := provider.waitForTransactionConfirmation(ctx, txID, MAX_WAIT_TIME)
		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})

	t.Run("Timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		pendingStatus := "pending"
		pendingResponse := &blockchainApiClient.GetTransactionById200Response{
			GetMempoolTransactionList200ResponseResultsInner: &blockchainApiClient.GetMempoolTransactionList200ResponseResultsInner{
				ContractCallMempoolTransaction1: &blockchainApiClient.ContractCallMempoolTransaction1{
					TxId: txID,
					TxStatus: blockchainApiClient.TokenTransferMempoolTransaction1TxStatus{
						String: &pendingStatus,
					},
				},
			},
		}

		mockClient.On("GetTransactionById", mock.Anything, txID).Return(pendingResponse, nil)

		timeoutDuration := 100 * time.Millisecond
		_, err := provider.waitForTransactionConfirmation(ctx, txID, timeoutDuration)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction confirmation timed out")
	})

	mockClient.AssertExpectations(t)
}

func TestProvider_ImportKeystore(t *testing.T) {
	provider, _ := setupTestProvider(t)
	mockKMS := new(mocks.MockKMS)
	provider.kms = mockKMS

	tempDir := t.TempDir()
	provider.cfg.HomeDir = tempDir

	importKeyPath := filepath.Join(tempDir, "import_keystore")
	encryptedKey := []byte("encrypted_import_key")
	err := os.WriteFile(importKeyPath, encryptedKey, 0600)
	assert.NoError(t, err)

	// Create a 33-byte private key (32 bytes for the key + 1 byte for compression flag)
	decryptedKey := make([]byte, 33)
	copy(decryptedKey, []byte("00000000000000000000000000000000")) // 32 bytes
	decryptedKey[32] = 0x01                                        // compression flag

	passphrase := "test_passphrase"

	mockKMS.On("Decrypt", mock.Anything, encryptedKey).Return(decryptedKey, nil)
	mockKMS.On("Encrypt", mock.Anything, decryptedKey).Return([]byte("new_encrypted_key"), nil)
	mockKMS.On("Encrypt", mock.Anything, []byte(passphrase)).Return([]byte("encrypted_passphrase"), nil)

	err = os.MkdirAll(filepath.Join(tempDir, "keystore", provider.NID()), 0700)
	assert.NoError(t, err)

	address, err := provider.ImportKeystore(context.Background(), importKeyPath, passphrase)
	assert.NoError(t, err)
	assert.NotEmpty(t, address)

	keystorePath := filepath.Join(tempDir, "keystore", provider.NID(), address)
	_, err = os.Stat(keystorePath)
	assert.NoError(t, err)

	_, err = os.Stat(keystorePath + ".pass")
	assert.NoError(t, err)

	assert.Equal(t, decryptedKey, provider.privateKey)
	assert.Equal(t, address, provider.cfg.Address)

	mockKMS.AssertExpectations(t)
}

func TestProvider_GetAddressByEventType(t *testing.T) {
	provider, _ := setupTestProvider(t)

	tests := []struct {
		name      string
		eventType string
		expected  string
	}{
		{
			name:      "EmitMessage Event",
			eventType: events.EmitMessage,
			expected:  provider.cfg.Contracts[providerTypes.ConnectionContract],
		},
		{
			name:      "CallMessage Event",
			eventType: events.CallMessage,
			expected:  provider.cfg.Contracts[providerTypes.XcallContract],
		},
		{
			name:      "RollbackMessage Event",
			eventType: events.RollbackMessage,
			expected:  provider.cfg.Contracts[providerTypes.XcallContract],
		},
		{
			name:      "Unknown Event",
			eventType: "unknown",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address := provider.GetAddressByEventType(tt.eventType)
			assert.Equal(t, tt.expected, address)
		})
	}
}
