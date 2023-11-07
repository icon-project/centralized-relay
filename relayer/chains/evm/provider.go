package evm

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/centralized-relay/relayer/chains/evm/abi"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/centralized-relay/relayer/store"

	"go.uber.org/zap"
)

var _ provider.ProviderConfig = &EVMProviderConfig{}

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

func (p *EVMProviderConfig) NewProvider(log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	client, err := NewClient(p.RPCUrl, log)
	if err != nil {
		return nil, err
	}

	return &EVMProvider{
		cfg:    p,
		log:    log.With(zap.String("chain_id", p.ChainID)),
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

func (p *EVMProvider) Init(context.Context) error {
	p.BlockReq = ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(p.cfg.ContractAddress)},
	}
	abi, err := abi.NewStorage(common.HexToAddress(p.cfg.ContractAddress), p.client.eth)
	if err != nil {
		return err
	}
	p.client.abi = abi
	return nil
}
