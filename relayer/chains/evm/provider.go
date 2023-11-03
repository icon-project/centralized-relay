package evm

import (
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/icon-project/centralized-relay/relayer/store"

	"go.uber.org/zap"
)

type EVMProviderConfig struct {
	ChainID         string `json:"chain-id" yaml:"chain-id"`
	Name            string `json:"name" yaml:"name"`
	RPCUrl          string `json:"rpc-url" yaml:"rpc-url"`
	StartHeight     uint64 `json:"start-height" yaml:"start-height"`
	Keystore        string `json:"keystore" yaml:"keystore"`
	Password        string `json:"password" yaml:"password"`
	ContractAddress string `json:"contract-address" yaml:"contract-address"`
}

type EVMProvider struct {
	sync.Mutex
	client *Client
	store.BlockStore
	log         *zap.Logger
	cfg         *EVMProviderConfig
	StartHeight uint64
	BlockReq    ethereum.FilterQuery
}

func (p *EVMProvider) NewProvider() (*EVMProvider, error) {
	if err := p.cfg.Validate(); err != nil {
		return nil, err
	}
	log := zap.NewNop()
	client, err := NewClient(p.cfg.RPCUrl, log)
	if err != nil {
		return nil, err
	}

	return &EVMProvider{
		log:    log,
		client: client,
	}, nil
}

func (p *EVMProvider) ChainId() string {
	return p.cfg.ChainID
}

func (p *EVMProviderConfig) Validate() error {
	// TODO:
	// add right validation
	return nil
}

func (p *EVMProvider) Init() error {
	return nil
}
