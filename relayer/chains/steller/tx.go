package steller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/icon-project/centralized-relay/relayer/chains/steller/types"
	evtypes "github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

func (p *Provider) Route(ctx context.Context, message *relayertypes.Message, callback relayertypes.TxResponseFunc) error {
	callArgs, err := p.newContractCallArgs(*message)
	if err != nil {
		return err
	}

	txRes, err := p.sendCallTransaction(*callArgs)

	cbTxResp := &relayertypes.TxResponse{}
	if txRes != nil {
		cbTxResp.Height = int64(txRes.Ledger)
		cbTxResp.TxHash = txRes.Hash
	}

	var cbErr error
	if err != nil {
		cbErr = err
	} else if txRes != nil && !txRes.Successful {
		cbErr = fmt.Errorf("failed to send call transaction")
	}

	callback(message.MessageKey(), cbTxResp, cbErr)

	return nil
}

func (p *Provider) sendCallTransaction(callArgs xdr.InvokeContractArgs) (*horizon.Transaction, error) {
	p.txmut.Lock()
	defer p.txmut.Unlock()

	callOp := txnbuild.InvokeHostFunction{
		HostFunction: xdr.HostFunction{
			Type:           xdr.HostFunctionTypeHostFunctionTypeInvokeContract,
			InvokeContract: &callArgs,
		},
	}

	sourceAccount, err := p.client.AccountDetail(p.wallet.Address())
	if err != nil {
		return nil, err
	}
	if _, err := sourceAccount.IncrementSequenceNumber(); err != nil {
		return nil, err
	}

	txParam := txnbuild.TransactionParams{
		SourceAccount:        &sourceAccount,
		IncrementSequenceNum: false,
		Operations:           []txnbuild.Operation{&callOp},
		BaseFee:              txnbuild.MinBaseFee,
		Preconditions: txnbuild.Preconditions{
			TimeBounds: txnbuild.NewTimeout(300),
		},
	}

	simtx, err := txnbuild.NewTransaction(txParam)
	if err != nil {
		return nil, err
	}
	simtxe, err := simtx.Base64()
	if err != nil {
		return nil, err
	}
	simres, err := p.client.SimulateTransaction(simtxe)
	if err != nil {
		return nil, err
	}

	var sorobanTxnData xdr.SorobanTransactionData
	if err := xdr.SafeUnmarshalBase64(simres.TransactionDataXDR, &sorobanTxnData); err != nil {
		return nil, err
	}

	callOp.Ext = xdr.TransactionExt{
		V:           1,
		SorobanData: &sorobanTxnData,
	}

	minResourceFee, err := strconv.Atoi(simres.MinResourceFee)
	if err != nil {
		return nil, err
	}

	txParam.BaseFee = int64(minResourceFee) + int64(p.cfg.MaxInclusionFee)

	tx, err := txnbuild.NewTransaction(txParam)
	if err != nil {
		return nil, err
	}
	tx, err = tx.Sign(p.cfg.NetworkPassphrase, p.wallet)
	if err != nil {
		return nil, err
	}
	txe, err := tx.Base64()
	if err != nil {
		return nil, err
	}
	txRes, err := p.client.SubmitTransactionXDR(txe)

	return &txRes, err
}

func (p *Provider) newContractCallArgs(msg relayertypes.Message) (*xdr.InvokeContractArgs, error) {
	scXcallAddr, err := p.scContractAddr(p.cfg.Contracts[relayertypes.XcallContract])
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

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txDigest string) (*relayertypes.Receipt, error) {
	//Todo
	return nil, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayertypes.MessageKey) (bool, error) {
	//Todo
	return false, nil
}
