package evm

import (
	"fmt"
	"sync"
	"time"

	"github.com/icon-project/centralized-relay/relayer/common"
	"github.com/icon-project/centralized-relay/relayer/store"

	"go.uber.org/zap"
)

type EVMProviderConfig struct {
	ChainID         string            `json:"chain-id" yaml:"chain-id"`
	Name            string            `json:"name" yaml:"name"`
	RPCURL          string            `json:"rpc-url" yaml:"rpc-url"`
	StartHeight     uint64            `json:"start-height" yaml:"start-height"`
	BlockDelay      int64             `json:"block-delay" yaml:"block-delay"`
	BlockInterval   int64             `json:"block-interval" yaml:"block-interval"`
	Type            common.ModuleType `json:"type" yaml:"type"`
	Keystore        string            `json:"keystore" yaml:"keystore"`
	Password        string            `json:"password" yaml:"password"`
	ContractAddress string            `json:"contract-address" yaml:"contract-address"`
	Timeout         string            `json:"timeout" yaml:"timeout"`
}

type EVMProvider struct {
	sync.Mutex
	Client
	store.BlockStore
	log *zap.Logger
	cfg *EVMProviderConfig
}

func (p *EVMProvider) NewProvider() (*EVMProvider, error) {
	if err := p.cfg.Validate(); err != nil {
		return nil, err
	}
	return &EVMProvider{
		log:    zap.NewNop(),
		Client: Client{},
	}, nil
}

func (p *EVMProvider) ChainId() string {
	return p.cfg.ChainID
}

func (p *EVMProviderConfig) Validate() error {
	if _, err := time.ParseDuration(p.Timeout); err != nil {
		return fmt.Errorf("invalid Timeout: %w", err)
	}

	if p.BlockInterval == 0 {
		return fmt.Errorf("Block interval cannot be zero")
	}

	return nil
}

func (p *EVMProvider) Init() error {
	return nil
}
