package kms

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

var ErrKeyAlreadyExists = fmt.Errorf("kms key already exists")

type KMS interface {
	Init(context.Context) (*string, error)
	Encrypt(context.Context, []byte) ([]byte, error)
	Decrypt(context.Context, []byte) ([]byte, error)
}

type KMSConfig struct {
	client *kms.Client
	key    *string
}

func NewKMSConfig(ctx context.Context, key *string) (KMS, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithDefaultRegion("us-east-1"))
	if err != nil {
		return nil, err
	}
	return &KMSConfig{kms.NewFromConfig(cfg), key}, nil
}

// Init creates a kms key for decryptying and encrypting data
func (k *KMSConfig) Init(ctx context.Context) (*string, error) {
	if len(*k.key) > 1 {
		return nil, ErrKeyAlreadyExists
	}
	input := &kms.CreateKeyInput{
		Description: aws.String("centralized-relay"),
	}
	output, err := k.client.CreateKey(ctx, input)
	if err != nil {
		return nil, err
	}
	return output.KeyMetadata.KeyId, nil
}

// Encrypt
func (k *KMSConfig) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	input := &kms.EncryptInput{
		KeyId:     k.key,
		Plaintext: data,
	}
	output, err := k.client.Encrypt(ctx, input)
	if err != nil {
		return nil, err
	}
	return output.CiphertextBlob, nil
}

// Decrypt
func (k *KMSConfig) Decrypt(ctx context.Context, data []byte) ([]byte, error) {
	input := &kms.DecryptInput{
		KeyId:          k.key,
		CiphertextBlob: data,
	}
	output, err := k.client.Decrypt(ctx, input)
	if err != nil {
		return nil, err
	}
	return output.Plaintext, nil
}
