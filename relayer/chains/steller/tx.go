package steller

import (
	"context"
	"fmt"

	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stellar/go/network"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

func (p *Provider) Route(ctx context.Context, message *relayertypes.Message, callback relayertypes.TxResponseFunc) error {
	contractHash, err := strkey.Decode(strkey.VersionByteContract, p.cfg.Contracts["connection"])
	if err != nil {
		return err
	}

	contractAddr, err := xdr.NewScAddress(xdr.ScAddressTypeScAddressTypeContract, xdr.Hash(contractHash))
	if err != nil {
		return err
	}

	dst, err := xdr.NewScVal(xdr.ScValTypeScvString, xdr.ScString("icon"))
	if err != nil {
		return err
	}

	data, err := xdr.NewScVal(xdr.ScValTypeScvBytes, xdr.ScBytes([]byte("hello")))
	if err != nil {
		return err
	}

	sourceAccount, err := p.client.AccountDetail(p.wallet.Address())
	if err != nil {
		return err
	}

	contractCallOp := txnbuild.InvokeHostFunction{
		SourceAccount: sourceAccount.AccountID,
		HostFunction: xdr.HostFunction{
			Type: xdr.HostFunctionTypeHostFunctionTypeInvokeContract,
			InvokeContract: &xdr.InvokeContractArgs{
				ContractAddress: contractAddr,
				FunctionName:    xdr.ScSymbol("new_message"),
				Args:            []xdr.ScVal{dst, data},
			},
		},
	}

	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &sourceAccount,
		IncrementSequenceNum: true,
		BaseFee:              txnbuild.MinBaseFee,
		Operations:           []txnbuild.Operation{&contractCallOp},
		Preconditions: txnbuild.Preconditions{
			TimeBounds: txnbuild.NewTimeout(300),
		},
	})
	if err != nil {
		return err
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, p.wallet)
	if err != nil {
		return err
	}

	txe, err := tx.Base64()
	if err != nil {
		return err
	}

	fmt.Println("envelope xdr: ", txe)

	// Send the transaction to the network
	resp, err := p.client.SimulateTransaction(txe)
	if err != nil {
		return err
	}

	fmt.Printf("\nSimulation Response: %+v\n", resp)

	return nil
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txDigest string) (*relayertypes.Receipt, error) {
	//Todo
	return nil, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayertypes.MessageKey) (bool, error) {
	//Todo
	return false, nil
}
