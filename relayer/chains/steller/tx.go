package steller

import (
	"context"
	"fmt"

	"github.com/icon-project/centralized-relay/relayer/chains/steller/sorobanclient"
	"github.com/icon-project/centralized-relay/relayer/chains/steller/types"
	evtypes "github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

func (p *Provider) scContractAddr(addr string) (*xdr.ScAddress, error) {
	contractHash, err := strkey.Decode(strkey.VersionByteContract, addr)
	if err != nil {
		return nil, err
	}

	scContractAddr, err := xdr.NewScAddress(xdr.ScAddressTypeScAddressTypeContract, xdr.Hash(contractHash))
	if err != nil {
		return nil, err
	}

	return &scContractAddr, nil
}

func (p *Provider) newContractCallArgs(msg relayertypes.Message) (*xdr.InvokeContractArgs, error) {
	scXcallAddr, err := p.scContractAddr(p.cfg.Contracts[relayertypes.ConnectionContract]) //Todo change to xcall contract
	if err != nil {
		return nil, err
	}
	scConnAddr, err := p.scContractAddr(p.cfg.Contracts[relayertypes.ConnectionContract])
	if err != nil {
		return nil, err
	}

	stellerMsg := types.StellerMsg{Message: msg}

	switch msg.EventType {
	case evtypes.EmitMessage:
		return &xdr.InvokeContractArgs{
			ContractAddress: *scConnAddr,
			FunctionName:    xdr.ScSymbol("recv_message"),
			Args: []xdr.ScVal{
				stellerMsg.ScvSn(),
				stellerMsg.ScvSrc(),
				stellerMsg.ScvData(),
			},
		}, nil
	case evtypes.CallMessage:
		return &xdr.InvokeContractArgs{
			ContractAddress: *scXcallAddr,
			FunctionName:    xdr.ScSymbol("execute_call"),
			Args: []xdr.ScVal{
				stellerMsg.ScvReqID(),
			},
		}, nil
	default:
		return &xdr.InvokeContractArgs{ //temporarily used for testing
			ContractAddress: *scConnAddr,
			FunctionName:    xdr.ScSymbol("new_message"),
			Args: []xdr.ScVal{
				stellerMsg.ScvDst(),
				stellerMsg.ScvData(),
			},
		}, nil
		// return nil, fmt.Errorf("invalid message type")
	}

}

func (p *Provider) simulateTransaction(op txnbuild.Operation) (*sorobanclient.TxSimulationResult, error) {
	sourceAccount, err := p.client.AccountDetail(p.wallet.Address())
	if err != nil {
		return nil, err
	}
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &sourceAccount,
		IncrementSequenceNum: true,
		BaseFee:              txnbuild.MinBaseFee,
		Operations:           []txnbuild.Operation{op},
		Preconditions: txnbuild.Preconditions{
			TimeBounds: txnbuild.NewTimeout(300),
		},
	})
	if err != nil {
		return nil, err
	}

	txe, err := tx.Base64()
	if err != nil {
		return nil, err
	}

	return p.client.SimulateTransaction(txe)
}

func (p *Provider) sendTransaction(fn func(param *txnbuild.TransactionParams)) (*horizon.Transaction, error) {
	sourceAccount, err := p.client.AccountDetail(p.wallet.Address())
	if err != nil {
		return nil, err
	}

	txParam := txnbuild.TransactionParams{
		SourceAccount:        &sourceAccount,
		IncrementSequenceNum: true,
		BaseFee:              txnbuild.MinBaseFee,
		Preconditions: txnbuild.Preconditions{
			TimeBounds: txnbuild.NewTimeout(300),
		},
	}

	fn(&txParam)

	tx, err := txnbuild.NewTransaction(txParam)
	if err != nil {
		return nil, err
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, p.wallet)
	if err != nil {
		return nil, err
	}

	txe, err := tx.Base64()
	if err != nil {
		return nil, err
	}

	txRes, err := p.client.SubmitTransactionXDR(txe)
	if err != nil {
		return nil, err
	}

	return &txRes, nil
}

func (p *Provider) Route(ctx context.Context, message *relayertypes.Message, callback relayertypes.TxResponseFunc) error {
	callArgs, err := p.newContractCallArgs(*message)
	if err != nil {
		return err
	}
	callOp := txnbuild.InvokeHostFunction{
		HostFunction: xdr.HostFunction{
			Type:           xdr.HostFunctionTypeHostFunctionTypeInvokeContract,
			InvokeContract: callArgs,
		},
	}

	simres, err := p.simulateTransaction(&callOp)
	if err != nil {
		return err
	}

	var sorobanTxnData xdr.SorobanTransactionData
	if err := xdr.SafeUnmarshalBase64(simres.TransactionDataXDR, &sorobanTxnData); err != nil {
		return err
	}

	callOp.Ext = xdr.TransactionExt{
		V:           1,
		SorobanData: &sorobanTxnData,
	}

	txRes, err := p.sendTransaction(func(param *txnbuild.TransactionParams) {
		param.Operations = []txnbuild.Operation{&callOp}
		param.BaseFee = 4074165
	})
	if err != nil {
		return err
	}

	fmt.Printf("\nTx Resp: %+v\n", txRes)

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
