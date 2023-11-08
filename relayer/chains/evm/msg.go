package evm

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

const defaultStepLimit = 13610920010

type Message struct {
	transactionOpt *bind.TransactOpts
	message        []byte
	Sn             int64
	pendingTx      *ethtypes.Transaction
	method         string
}

func (p *EVMProvider) NewMessage(ctx context.Context, msg providerTypes.Message, method string) (*Message, error) {
	newTransactOpts := func(w *keystore.Key) (*bind.TransactOpts, error) {
		txo, err := bind.NewKeyedTransactorWithChainID(w.PrivateKey, p.client.GetChainID())
		if err != nil {
			return nil, err
		}
		ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
		defer cancel()
		txo.GasPrice, _ = p.client.SuggestGasPrice(ctx)
		txo.GasLimit = uint64(p.cfg.GasLimit * 2)
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

	return &Message{
		message:        msg.Data,
		transactionOpt: txOpts,
	}, nil
}
