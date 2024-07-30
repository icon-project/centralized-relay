package stacks

import (
	"encoding/hex"
	"fmt"
)

type MessageSignature struct {
	Type StacksMessageType
	Data string
}

func CreateMessageSignature(signature string) (*MessageSignature, error) {
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, fmt.Errorf("invalid signature hex: %w", err)
	}
	if len(signatureBytes) != RECOVERABLE_ECDSA_SIG_LENGTH_BYTES {
		return nil, fmt.Errorf("invalid signature length: expected %d, got %d", RECOVERABLE_ECDSA_SIG_LENGTH_BYTES, len(signatureBytes))
	}
	return &MessageSignature{
		Type: StacksMessageTypeMessageSignature,
		Data: signature,
	}, nil
}

func (ms *MessageSignature) Serialize() ([]byte, error) {
	return hex.DecodeString(ms.Data)
}
