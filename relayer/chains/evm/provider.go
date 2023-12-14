package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/centralized-relay/relayer/store"

	"go.uber.org/zap"
)

var _ provider.ProviderConfig = &EVMProviderConfig{}

type EVMProviderConfig struct {
	ChainID         string `json:"chain-id" yaml:"chain-id"`
	Name            string `json:"name" yaml:"name"`
	RPCUrl          string `json:"rpc-url" yaml:"rpc-url"`
	VerifierRPCUrl  string `json:"verifier-rpc-url" yaml:"verifier-rpc-url"`
	StartHeight     uint64 `json:"start-height" yaml:"start-height"`
	Keystore        string `json:"keystore" yaml:"keystore"`
	Password        string `json:"password" yaml:"password"`
	GasPrice        int64  `json:"gas-price" yaml:"gas-price"`
	GasLimit        uint64 `json:"gas-limit" yaml:"gas-limit"`
	ContractAddress string `json:"contract-address" yaml:"contract-address"`
	Concurrency     uint64 `json:"concurrency" yaml:"concurrency"`
	FinalityBlock   uint64 `json:"finality-block" yaml:"finality-block"`
}

type EVMProvider struct {
	client   IClient
	verifier IClient
	store.BlockStore
	log         *zap.Logger
	cfg         *EVMProviderConfig
	StartHeight uint64
	blockReq    ethereum.FilterQuery
	wallet      *keystore.Key
}

func (p *EVMProviderConfig) NewProvider(log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	client, err := newClient(p.RPCUrl, p.ContractAddress, log)
	if err != nil {
		return nil, err
	}

	var verifierClient IClient

	if p.VerifierRPCUrl != "" {
		var err error
		verifierClient, err = newClient(p.VerifierRPCUrl, p.ContractAddress, log)
		if err != nil {
			return nil, err
		}
	} else {
		verifierClient = client // default to same client
	}

	// setting default finality block
	if p.FinalityBlock == 0 {
		p.FinalityBlock = uint64(DefaultFinalityBlock)
	}

	return &EVMProvider{
		cfg:      p,
		log:      log.With(zap.String("chain_id", p.ChainID)),
		client:   client,
		blockReq: getEventFilterQuery(p.ContractAddress),
		verifier: verifierClient,
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
	wallet, err := RestoreKey(p.cfg.Keystore, p.cfg.Password)
	if err != nil {
		return fmt.Errorf("failed to restore evm wallet %v", err)
	}
	p.wallet = wallet
	return nil
}

func (p *EVMProvider) Wallet() *keystore.Key {
	return p.wallet
}

func (p *EVMProvider) FinalityBlock(ctx context.Context) uint64 {
	return p.cfg.FinalityBlock
}

func (p *EVMProvider) WaitForResults(ctx context.Context, txHash common.Hash) (txr *ethTypes.Receipt, err error) {
	const DefaultGetTransactionResultPollingInterval = 1500 * time.Millisecond // 1.5sec
	ticker := time.NewTicker(time.Duration(DefaultGetTransactionResultPollingInterval) * time.Nanosecond)
	retryLimit := 10
	retryCounter := 0
	for {
		defer ticker.Stop()
		select {
		case <-ctx.Done():
			err = errors.New("Context Cancelled. ResultWait Exiting ")
			return
		case <-ticker.C:
			if retryCounter >= retryLimit {
				err = errors.New("Retry Limit Exceeded while waiting for results of transaction")
				return
			}
			retryCounter++
			txr, err = p.client.TransactionReceipt(context.Background(), txHash)
			if err != nil && err == ethereum.NotFound {
				err = nil
				continue
			}
			return
		}
	}
}

func (r *EVMProvider) transferBalance(senderKey, recepientAddress string, amount *big.Int) (txnHash common.Hash, err error) {
	from, err := crypto.HexToECDSA(senderKey)
	if err != nil {
		return common.Hash{}, err
	}

	fromAddress := crypto.PubkeyToAddress(from.PublicKey)

	nonce, err := r.client.NonceAt(context.TODO(), fromAddress, nil)
	if err != nil {
		err = errors.Wrap(err, "PendingNonceAt ")
		return common.Hash{}, err
	}
	gasPrice, err := r.client.SuggestGasPrice(context.Background())
	if err != nil {
		err = errors.Wrap(err, "SuggestGasPrice ")
		return common.Hash{}, err
	}
	chainID := r.client.GetChainID()

	tx := types.NewTransaction(nonce, common.HexToAddress(recepientAddress), amount, 30000000, gasPrice, []byte{})
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), from)
	if err != nil {
		err = errors.Wrap(err, "SignTx ")
		return common.Hash{}, err
	}

	if err = r.client.SendTransaction(context.TODO(), signedTx); err != nil {
		err = errors.Wrap(err, "SendTransaction ")
		return
	}
	txnHash = signedTx.Hash()
	return
}

func (p *EVMProvider) GetTransationOpts(ctx context.Context) (*bind.TransactOpts, error) {
	newTransactOpts := func(w *keystore.Key) (*bind.TransactOpts, error) {
		txo, err := bind.NewKeyedTransactorWithChainID(w.PrivateKey, p.client.GetChainID())
		if err != nil {
			return nil, err
		}
		ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
		defer cancel()
		txo.GasPrice, _ = p.client.SuggestGasPrice(ctx)
		txo.GasLimit = uint64(p.cfg.GasLimit)
		return txo, nil
	}

	txOpts, err := newTransactOpts(p.wallet)
	if err != nil {
		return nil, err
	}
	txOpts.Context = ctx
	if p.cfg.GasLimit > 0 {
		txOpts.GasLimit = p.cfg.GasLimit
	}
	txOpts.GasPrice = big.NewInt(p.cfg.GasPrice)

	return txOpts, nil
}
