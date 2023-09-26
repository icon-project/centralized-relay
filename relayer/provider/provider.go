package provider

import (
	"context"

	"go.uber.org/zap"
)

var (
	messageMaxRetry = 5
)

type ProviderConfig interface {
	NewProvider(log *zap.Logger, homepath string, debug bool, chainName string) (ChainProvider, error)
	Validate() error
}

type ChainProvider interface {
	ChainId() string
	Init(ctx context.Context) error
	Listener(ctx context.Context, blockInfo <-chan BlockInfo) error
	Route(ctx context.Context, message RouteMessage, callback func(response ExecuteMessageResponse)) error
}

type BlockInfo struct {
	Height   uint64
	Messages []RelayMessage
}

type RelayMessage struct {
	Target string
	Src    string
	Sn     uint64
	Data   []byte
}

type RouteMessage struct {
	RelayMessage
	retry uint64
}

type ExecuteMessageResponse struct {
	RouteMessage
	TxResponse
}

type TxResponse struct {
	Height    int64
	TxHash    string
	Codespace string
	Code      uint32
	Data      string
}
