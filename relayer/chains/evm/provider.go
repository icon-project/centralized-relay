package evm

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/centralized-relay/relayer/store"
	"github.com/icon-project/icon-bridge/common/wallet"

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
	GasPrice        int64  `json:"gas-price" yaml:"gas-price"`
	GasLimit        uint64 `json:"gas-limit" yaml:"gas-limit"`
	EVMChainID      int64  `json:"evm-chain-id" yaml:"evm-chain-id"`
	ContractAddress string `json:"contract-address" yaml:"contract-address"`
}

type EVMProvider struct {
	sync.Mutex
	client IClient
	store.BlockStore
	log         *zap.Logger
	cfg         *EVMProviderConfig
	StartHeight uint64
	BlockReq    ethereum.FilterQuery
	wallet      wallet.EvmWallet
}

func (p *EVMProviderConfig) NewProvider(log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	client, err := newClient(p.RPCUrl, p.ContractAddress, log)
	if err != nil {
		return nil, err
	}

	// get wallet

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
	// Contract address check
	// gas limit mandatory
	// keystore
	return nil
}

func (p *EVMProvider) Init(context.Context) error {
	p.BlockReq = ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(p.cfg.ContractAddress)},
	}
	return nil
}

func (p *EVMProvider) EVMChainID() int64 {
	return p.cfg.EVMChainID
}
