package stacks

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/chains/stacks/clarity"
)

func NewStacksTransaction(
	version TransactionVersion,
	auth Authorization,
	payload Payload,
	postConditions []PostCondition,
	postConditionMode PostConditionMode,
	anchorMode AnchorMode,
	chainID uint32,
) *StacksTransaction {
	return &StacksTransaction{
		Version:           version,
		ChainID:           chainID,
		Auth:              auth,
		AnchorMode:        anchorMode,
		PostConditionMode: postConditionMode,
		PostConditions:    postConditions,
		Payload:           payload,
	}
}

func (tx *StacksTransaction) Serialize() ([]byte, error) {
	buffer := new(bytes.Buffer)

	if err := binary.Write(buffer, binary.BigEndian, tx.Version); err != nil {
		return nil, fmt.Errorf("failed to write version: %w", err)
	}

	if err := binary.Write(buffer, binary.BigEndian, tx.ChainID); err != nil {
		return nil, fmt.Errorf("failed to write chain ID: %w", err)
	}

	authBytes, err := tx.Auth.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize authorization: %w", err)
	}
	if _, err := buffer.Write(authBytes); err != nil {
		return nil, fmt.Errorf("failed to write authorization: %w", err)
	}

	if err := binary.Write(buffer, binary.BigEndian, tx.AnchorMode); err != nil {
		return nil, fmt.Errorf("failed to write anchor mode: %w", err)
	}

	if err := binary.Write(buffer, binary.BigEndian, tx.PostConditionMode); err != nil {
		return nil, fmt.Errorf("failed to write post-condition mode: %w", err)
	}

	postConditionsLPList := createLPList(tx.PostConditions, 4)
	postConditionsBytes, err := postConditionsLPList.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize post conditions: %w", err)
	}
	if _, err := buffer.Write(postConditionsBytes); err != nil {
		return nil, fmt.Errorf("failed to write post conditions: %w", err)
	}

	payloadBytes, err := tx.Payload.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize payload: %w", err)
	}
	if _, err := buffer.Write(payloadBytes); err != nil {
		return nil, fmt.Errorf("failed to write payload: %w", err)
	}

	return buffer.Bytes(), nil
}

func (a *StandardAuthorization) Serialize() ([]byte, error) {
	buffer := new(bytes.Buffer)

	if err := binary.Write(buffer, binary.BigEndian, AuthTypeStandard); err != nil {
		return nil, fmt.Errorf("failed to write auth type: %w", err)
	}

	scBytes, err := a.SpendingCondition.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize spending condition: %w", err)
	}
	if _, err := buffer.Write(scBytes); err != nil {
		return nil, fmt.Errorf("failed to write spending condition: %w", err)
	}

	return buffer.Bytes(), nil
}

func (a *SponsoredAuthorization) Serialize() ([]byte, error) {
	buffer := new(bytes.Buffer)

	if err := binary.Write(buffer, binary.BigEndian, AuthTypeSponsored); err != nil {
		return nil, fmt.Errorf("failed to write auth type: %w", err)
	}

	scBytes, err := a.SpendingCondition.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize spending condition: %w", err)
	}
	if _, err := buffer.Write(scBytes); err != nil {
		return nil, fmt.Errorf("failed to write spending condition: %w", err)
	}

	sponsorScBytes, err := a.SponsorSpendingCondition.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize sponsor spending condition: %w", err)
	}
	if _, err := buffer.Write(sponsorScBytes); err != nil {
		return nil, fmt.Errorf("failed to write sponsor spending condition: %w", err)
	}

	return buffer.Bytes(), nil
}

func (s *SingleSigSpendingCondition) Serialize() ([]byte, error) {
	buffer := new(bytes.Buffer)

	if err := binary.Write(buffer, binary.BigEndian, s.HashMode); err != nil {
		return nil, fmt.Errorf("failed to write hash mode: %w", err)
	}

	if _, err := buffer.Write(s.Signer); err != nil {
		return nil, fmt.Errorf("failed to write signer: %w", err)
	}

	if err := binary.Write(buffer, binary.BigEndian, s.Nonce); err != nil {
		return nil, fmt.Errorf("failed to write nonce: %w", err)
	}

	if err := binary.Write(buffer, binary.BigEndian, s.Fee); err != nil {
		return nil, fmt.Errorf("failed to write fee: %w", err)
	}

	if err := binary.Write(buffer, binary.BigEndian, s.KeyEncoding); err != nil {
		return nil, fmt.Errorf("failed to write key encoding: %w", err)
	}

	sigBytes, err := s.Signature.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signature: %w", err)
	}
	if _, err := buffer.Write(sigBytes); err != nil {
		return nil, fmt.Errorf("failed to write signature: %w", err)
	}

	return buffer.Bytes(), nil
}

func (m *MultiSigSpendingCondition) Serialize() ([]byte, error) {
	buffer := new(bytes.Buffer)

	if err := binary.Write(buffer, binary.BigEndian, m.HashMode); err != nil {
		return nil, fmt.Errorf("failed to write hash mode: %w", err)
	}

	if _, err := buffer.Write(m.Signer); err != nil {
		return nil, fmt.Errorf("failed to write signer: %w", err)
	}

	if err := binary.Write(buffer, binary.BigEndian, m.Nonce); err != nil {
		return nil, fmt.Errorf("failed to write nonce: %w", err)
	}

	if err := binary.Write(buffer, binary.BigEndian, m.Fee); err != nil {
		return nil, fmt.Errorf("failed to write fee: %w", err)
	}

	fieldsLPList := createLPList(m.Fields, 4)
	fieldBytes, err := fieldsLPList.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize fields: %w", err)
	}
	if _, err := buffer.Write(fieldBytes); err != nil {
		return nil, fmt.Errorf("failed to write fields: %w", err)
	}

	if err := binary.Write(buffer, binary.BigEndian, m.SignaturesRequired); err != nil {
		return nil, fmt.Errorf("failed to write signatures required: %w", err)
	}

	return buffer.Bytes(), nil
}

func CreateTransactionAuthField(pubKeyEncoding PubKeyEncoding, contents TransactionAuthFieldContents) *TransactionAuthField {
	return &TransactionAuthField{
		Type:           StacksMessageTypeTransactionAuthField,
		PubKeyEncoding: pubKeyEncoding,
		Contents:       contents,
	}
}

func (taf *TransactionAuthField) Serialize() ([]byte, error) {
	buffer := new(bytes.Buffer)

	switch contents := taf.Contents.(type) {
	case *StacksPublicKey:
		if taf.PubKeyEncoding == PubKeyEncodingCompressed {
			buffer.WriteByte(byte(AuthFieldTypePublicKeyCompressed))
		} else {
			buffer.WriteByte(byte(AuthFieldTypePublicKeyUncompressed))
		}
		contentBytes, err := contents.Serialize()
		if err != nil {
			return nil, fmt.Errorf("failed to serialize StacksPublicKey: %w", err)
		}
		buffer.Write(contentBytes)
	case *MessageSignature:
		if taf.PubKeyEncoding == PubKeyEncodingCompressed {
			buffer.WriteByte(byte(AuthFieldTypeSignatureCompressed))
		} else {
			buffer.WriteByte(byte(AuthFieldTypeSignatureUncompressed))
		}
		contentBytes, err := contents.Serialize()
		if err != nil {
			return nil, fmt.Errorf("failed to serialize MessageSignature: %w", err)
		}
		buffer.Write(contentBytes)
	default:
		return nil, fmt.Errorf("unknown TransactionAuthField contents type")
	}

	return buffer.Bytes(), nil
}

func (spk *StacksPublicKey) Serialize() ([]byte, error) {
	return spk.Data, nil
}

func (p *TokenTransferPayload) Serialize() ([]byte, error) {
	buffer := new(bytes.Buffer)

	// Write payload type
	if err := binary.Write(buffer, binary.BigEndian, p.Type()); err != nil {
		return nil, fmt.Errorf("failed to write payload type: %w", err)
	}

	// Serialize recipient
	recipientBytes, err := clarity.SerializeCV(p.Recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize recipient: %w", err)
	}
	buffer.Write(recipientBytes)

	// Write amount
	if err := binary.Write(buffer, binary.BigEndian, p.Amount); err != nil {
		return nil, fmt.Errorf("failed to write amount: %w", err)
	}

	// Serialize memo
	memoBytes, err := p.Memo.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize memo: %w", err)
	}
	buffer.Write(memoBytes)

	return buffer.Bytes(), nil
}

func SerializePayload(payload Payload) ([]byte, error) {
	return payload.Serialize()
}
