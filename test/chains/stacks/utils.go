package stacks

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	stacksClient "github.com/icon-project/centralized-relay/relayer/chains/stacks"
	"github.com/icon-project/stacks-go-sdk/pkg/crypto"
	"go.uber.org/zap"
)

func (s *StacksLocalnet) extractSnFromTransaction(ctx context.Context, txID string) (string, error) {
	tx, err := s.client.GetTransactionById(ctx, txID)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction by ID: %w", err)
	}

	if confirmed := tx.GetTransactionList200ResponseResultsInner; confirmed != nil {
		if contractCall := confirmed.ContractCallTransaction; contractCall != nil {
			for _, event := range contractCall.Events {
				if event.SmartContractLogTransactionEvent != nil {
					contractLog := event.SmartContractLogTransactionEvent.ContractLog
					if contractLog.Topic == "print" {
						repr := contractLog.Value.Repr
						s.log.Debug("Found event log",
							zap.String("repr", repr),
							zap.String("topic", contractLog.Topic))
						if strings.Contains(repr, "CallMessageSent") {
							s.log.Debug("Found CallMessageSent event",
								zap.String("full_event", repr))
							startIdx := strings.Index(repr, "(sn u")
							if startIdx != -1 {
								startIdx += 5 // Move past "(sn u"
								endIdx := strings.Index(repr[startIdx:], ")")
								if endIdx != -1 {
									sn := repr[startIdx : startIdx+endIdx]
									s.log.Info("Successfully extracted sn",
										zap.String("sn", sn))
									return sn, nil
								}
							}
						}
					}
				}
			}
		}
	}

	return "", fmt.Errorf("serial number 'sn' not found in transaction events")
}

func (s *StacksLocalnet) loadPrivateKey(keystoreFile, password string) ([]byte, string, error) {
	if s.kms == nil {
		return nil, "", fmt.Errorf("KMS not initialized")
	}

	encryptedKey, err := os.ReadFile(keystoreFile)
	if err != nil {
		s.log.Error("Failed to read keystore file",
			zap.String("path", keystoreFile),
			zap.Error(err))
		return nil, "", fmt.Errorf("failed to read keystore file: %w", err)
	}

	privateKey, err := s.kms.Decrypt(context.Background(), encryptedKey)
	if err != nil {
		s.log.Error("Failed to decrypt keystore", zap.Error(err))
		return nil, "", fmt.Errorf("failed to decrypt keystore: %w", err)
	}

	network, err := stacksClient.MapNIDToChainID(s.cfg.ChainID)
	if err != nil {
		s.log.Error("Chain id not found. Check the NID config", zap.Error(err))
		return nil, "", fmt.Errorf("chain id not found: %w", err)
	}

	address, err := crypto.GetAddressFromPrivateKey(privateKey, network)
	if err != nil {
		s.log.Error("Failed to derive address from private key", zap.Error(err))
		return nil, "", fmt.Errorf("failed to derive address: %w", err)
	}

	return privateKey, address, nil
}

func (s *StacksLocalnet) waitForTransactionConfirmation(ctx context.Context, txID string) error {
	timeout := time.After(MAX_WAIT_TIME)
	ticker := time.NewTicker(BLOCK_TIME)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			res, err := s.client.GetTransactionById(ctx, txID)
			if err != nil {
				s.log.Warn("Failed to query transaction receipt", zap.Error(err))
				continue
			}

			receipt, err := stacksClient.GetReceipt(res)
			if err != nil {
				s.log.Warn("Failed to extract transaction receipt", zap.Error(err))
				continue
			}

			if receipt.Status {
				s.log.Info("Transaction confirmed",
					zap.String("txID", txID),
					zap.Uint64("height", receipt.Height))
				return nil
			}
			s.log.Debug("Transaction not yet confirmed", zap.String("txID", txID))

		case <-timeout:
			return fmt.Errorf("transaction confirmation timed out after %v seconds", MAX_WAIT_TIME)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
