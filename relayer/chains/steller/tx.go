package steller

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/icon-project/centralized-relay/relayer/chains/steller/types"
	evtypes "github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	xdr3 "github.com/stellar/go-xdr/xdr3"
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
	responseCode := relayertypes.Failed
	if txRes != nil && txRes.Successful {
		responseCode = relayertypes.Success
	}
	if txRes != nil {
		cbTxResp.Height = int64(txRes.Ledger)
		cbTxResp.TxHash = txRes.Hash
		cbTxResp.Code = responseCode
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
	var auth []xdr.SorobanAuthorizationEntry
	for _, res := range simres.Results {
		var decodedRes xdr.ScVal
		err := xdr.SafeUnmarshalBase64(res.Xdr, &decodedRes)
		if err != nil {
			return nil, err
		}
		for _, authBase64 := range res.Auth {
			var authEntry xdr.SorobanAuthorizationEntry
			err = xdr.SafeUnmarshalBase64(authBase64, &authEntry)
			if err != nil {
				return nil, err
			}
			auth = append(auth, authEntry)
		}
	}
	callOp.Auth = auth
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
				stellerMsg.ScvSrc(),
				stellerMsg.ScvSn(),
				stellerMsg.ScvData(),
			},
		}, nil
	case evtypes.CallMessage:
		acc := xdr.MustAddressPtr(p.cfg.Address)
		return &xdr.InvokeContractArgs{
			ContractAddress: *scXcallAddr,
			FunctionName:    xdr.ScSymbol("execute_call"),
			Args: []xdr.ScVal{
				{
					Type: xdr.ScValTypeScvAddress,
					Address: &xdr.ScAddress{
						AccountId: acc,
					},
				},
				stellerMsg.ScvReqID(),
				stellerMsg.ScvData(),
			},
		}, nil
	default:
		return nil, fmt.Errorf("invalid message type")
	}
}

func (p *Provider) newMiscContractCallArgs(msg relayertypes.Message, params ...interface{}) (*xdr.InvokeContractArgs, error) {
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
	case evtypes.GetFee:
		rollBackVal := params[0].(bool)
		rollbackParam := xdr.ScVal{
			Type: xdr.ScValTypeScvBool,
			B:    &rollBackVal,
		}
		src := &xdr.ScVec{}
		sources := xdr.ScVal{
			Type: xdr.ScValTypeScvVec,
			Vec:  &src,
		}
		return &xdr.InvokeContractArgs{
			ContractAddress: *scXcallAddr,
			FunctionName:    xdr.ScSymbol("get_fee"),
			Args: []xdr.ScVal{
				stellerMsg.ScvSrc(),
				rollbackParam,
				sources,
			},
		}, nil
	case evtypes.SetFee:
		return &xdr.InvokeContractArgs{
			ContractAddress: *scConnAddr,
			FunctionName:    xdr.ScSymbol("set_fee"),
			Args: []xdr.ScVal{
				stellerMsg.ScvSrc(),
				stellerMsg.ScvSn(),
				stellerMsg.ScvReqID(),
			},
		}, nil
	case evtypes.SetAdmin:
		adminAddress := stellerMsg.Src
		acc := xdr.MustAddressPtr(adminAddress)
		adminAddressVal := xdr.ScVal{
			Type: xdr.ScValTypeScvAddress,
			Address: &xdr.ScAddress{
				AccountId: acc,
			},
		}
		return &xdr.InvokeContractArgs{
			ContractAddress: *scXcallAddr,
			FunctionName:    xdr.ScSymbol("set_admin"),
			Args: []xdr.ScVal{
				adminAddressVal,
			},
		}, nil
	case evtypes.ClaimFee:
		return &xdr.InvokeContractArgs{
			ContractAddress: *scConnAddr,
			FunctionName:    xdr.ScSymbol("claim_fees"),
		}, nil
	case evtypes.RevertMessage:
		return &xdr.InvokeContractArgs{
			ContractAddress: *scXcallAddr,
			FunctionName:    xdr.ScSymbol("revert_message"),
			Args: []xdr.ScVal{
				stellerMsg.ScvSn(),
			},
		}, nil
	default:
		return nil, fmt.Errorf("invalid message type")
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

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txHash string) (*relayertypes.Receipt, error) {
	tx, err := p.client.GetTransaction(txHash)
	if err != nil {
		return nil, err
	}
	return &relayertypes.Receipt{
		TxHash: txHash,
		Height: uint64(tx.Ledger),
		Status: tx.Successful,
	}, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayertypes.MessageKey) (bool, error) {
	connAddr, err := p.scContractAddr(p.cfg.Contracts[relayertypes.ConnectionContract])
	if err != nil {
		return false, err
	}
	stellerMsg := types.StellerMsg{
		Message: relayertypes.Message{
			Sn:  key.Sn,
			Src: key.Src,
		},
	}
	callArgs := xdr.InvokeContractArgs{
		ContractAddress: *connAddr,
		FunctionName:    xdr.ScSymbol("get_receipt"),
		Args: []xdr.ScVal{
			stellerMsg.ScvSrc(),
			stellerMsg.ScvSn(),
		},
	}

	var isReceived types.ScvBool
	if err := p.queryContract(callArgs, &isReceived); err != nil {
		return false, err
	}

	return bool(isReceived), nil
}

func (p *Provider) queryContract(callArgs xdr.InvokeContractArgs, dest types.ScValConverter) error {
	sourceAccount, err := p.client.AccountDetail(p.wallet.Address())
	if err != nil {
		return err
	}

	callOp := txnbuild.InvokeHostFunction{
		SourceAccount: sourceAccount.AccountID,
		HostFunction: xdr.HostFunction{
			Type:           xdr.HostFunctionTypeHostFunctionTypeInvokeContract,
			InvokeContract: &callArgs,
		},
	}

	txParam := txnbuild.TransactionParams{
		SourceAccount: &sourceAccount,
		Operations:    []txnbuild.Operation{&callOp},
		Preconditions: txnbuild.Preconditions{
			TimeBounds: txnbuild.NewTimeout(300),
		},
	}

	queryTx, err := txnbuild.NewTransaction(txParam)
	if err != nil {
		return err
	}
	queryTxe, err := queryTx.Base64()
	if err != nil {
		return err
	}
	queryRes, err := p.client.SimulateTransaction(queryTxe)
	if err != nil {
		return err
	}
	for _, callResult := range queryRes.Results {
		resBytes, err := base64.StdEncoding.DecodeString(callResult.Xdr)
		if err != nil {
			return err
		}
		var scVal xdr.ScVal
		if _, err := xdr3.Unmarshal(bytes.NewReader(resBytes), &scVal); err != nil {
			return err
		} else {
			if err := dest.Convert(scVal); err != nil {
				return err
			}
			break
		}
	}

	return nil
}
