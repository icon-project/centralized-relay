package kms

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go/aws"
)

type KMS interface {
	Init(context.Context) (*string, error)
	Encrypt(context.Context, string, []byte) ([]byte, error)
	Decrypt(context.Context, string, []byte) ([]byte, error)
}

type KMSConfig struct {
	client *kms.Client
}

func NewKMSConfig(ctx context.Context, profile string) (KMS, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, err
	}
	return &KMSConfig{client: kms.NewFromConfig(cfg)}, nil
}

// Init creates a kms key for decrptying and encrypting data
func (k *KMSConfig) Init(ctx context.Context) (*string, error) {
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
func (k *KMSConfig) Encrypt(ctx context.Context, keyID string, data []byte) ([]byte, error) {
	input := &kms.EncryptInput{
		KeyId:     &keyID,
		Plaintext: data,
	}
	output, err := k.client.Encrypt(ctx, input)
	if err != nil {
		return nil, err
	}
	return output.CiphertextBlob, nil
}

// Decrypt
func (k *KMSConfig) Decrypt(ctx context.Context, keyID string, data []byte) ([]byte, error) {
	input := &kms.DecryptInput{
		KeyId:          &keyID,
		CiphertextBlob: data,
	}
	output, err := k.client.Decrypt(ctx, input)
	if err != nil {
		return nil, err
	}
	return output.Plaintext, nil
}
