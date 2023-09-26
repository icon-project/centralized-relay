package provider

import (
	"context"

	"go.uber.org/zap"
)

type ProviderConfig interface {
	NewProvider(log *zap.Logger, homepath string, debug bool, chainName string) (ChainProvider, error)
	Validate() error
}

type ChainProvider interface {
	ChainId() string
	Init(ctx context.Context) error
}
