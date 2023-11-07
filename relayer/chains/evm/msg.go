package evm

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/icon-bridge/common/wallet"
)

const defaultStepLimit = 13610920010

type Message struct {
	transactionOpt *bind.TransactOpts
	message        []byte
	pendingTx      *ethtypes.Transaction
	method         string
}

func (m *Message) Type() string {
	return m.method
}

func (m *Message) MsgBytes() ([]byte, error) {
	return json.Marshal(m.message)
}

func (p *EVMProvider) NewMessage(ctx context.Context, msg providerTypes.Message, method string) (*Message, error) {
	newTransactOpts := func(w *wallet.EvmWallet) (*bind.TransactOpts, error) {
		txo, err := bind.NewKeyedTransactorWithChainID(w.Skey, big.NewInt(p.EVMChainID()))
		if err != nil {
			return nil, err
		}
		ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
		defer cancel()
		txo.GasPrice, _ = p.client.SuggestGasPrice(ctx)
		txo.GasLimit = uint64(p.cfg.GasLimit)
		return txo, nil
	}

	txOpts, err := newTransactOpts(&p.wallet)
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
