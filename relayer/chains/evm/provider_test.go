package evm

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func MockEvmProvider() (*EVMProvider, error) {

	evm := EVMProviderConfig{
		ChainID:     "eth",
		Name:        "eth",
		RPCUrl:      "http://localhost:8545",
		StartHeight: 0,
		Keystore:    testKeyStore,
		Password:    testKeyPassword,
		GasPrice:    1000565528,
		// GasLimit:        200_000_000,
		ContractAddress: "0x0165878A594ca255338adfa4d48449f69242Eb8F",
	}
	log := zap.NewNop()
	pro, err := evm.NewProvider(log, "/Users/viveksharmapoudel/my_work_bench/ibriz/ibc-related/centralized-relay", true, "evm-1")
	if err != nil {
		return nil, err
	}
	p, ok := pro.(*EVMProvider)
	if !ok {
		return nil, fmt.Errorf("failed to create mock evmprovider")
	}

	p.Init(context.TODO())
	return p, nil
}

func TestGetBalance(t *testing.T) {
	pro, err := MockEvmProvider()
	assert.NoError(t, err)

	header, _ := pro.client.GetHeaderByHeight(context.TODO(), big.NewInt(117))
	fmt.Println(header.GasLimit)
	txhash, err := pro.transferBalance(
		"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		pro.wallet.Address.Hex(), big.NewInt(100_000_000_000_000_000_0))

	assert.NoError(t, err)

	r, err := pro.WaitForResults(context.TODO(), txhash)
	assert.NoError(t, err)

	fmt.Println("status of the transaction ", r.Status)
	fmt.Println("transaction hash", r.TxHash)

}
