package steller

import (
	"context"
	"encoding/hex"
	"slices"

	"github.com/icon-project/centralized-relay/relayer/chains/steller/sorobanclient"
	"github.com/icon-project/centralized-relay/relayer/chains/steller/types"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

type IClient interface {
	SimulateTransaction(txXDR string) (*sorobanclient.TxSimulationResult, error)

	SubmitTransactionXDR(txXDR string) (horizon.Transaction, error)

	SubmitTransaction(tx *txnbuild.Transaction) (horizon.Transaction, error)

	AccountDetail(addr string) (horizon.Account, error)

	GetLatestLedger(ctx context.Context) (*sorobanclient.LatestLedgerResponse, error)

	FetchEvents(ctx context.Context, eventFilter types.EventFilter) ([]types.Event, error)
}

type Client struct {
	horizon *horizonclient.Client
	soroban *sorobanclient.Client
}

func NewClient(hClient *horizonclient.Client, srbClient *sorobanclient.Client) IClient {
	return &Client{horizon: hClient, soroban: srbClient}
}

func (cl *Client) SimulateTransaction(txXDR string) (*sorobanclient.TxSimulationResult, error) {
	return cl.soroban.SimulateTransaction(txXDR, nil)
}

func (cl *Client) SubmitTransactionXDR(txXDR string) (horizon.Transaction, error) {
	return cl.horizon.SubmitTransactionXDR(txXDR)
}

func (cl *Client) SubmitTransaction(tx *txnbuild.Transaction) (horizon.Transaction, error) {
	return cl.horizon.SubmitTransaction(tx)
}

func (cl *Client) AccountDetail(addr string) (horizon.Account, error) {
	return cl.horizon.AccountDetail(horizonclient.AccountRequest{AccountID: addr})
}

func (cl *Client) GetLatestLedger(ctx context.Context) (*sorobanclient.LatestLedgerResponse, error) {
	return cl.soroban.GetLatestLedger(ctx)
}

func (cl *Client) FetchEvents(ctx context.Context, eventFilter types.EventFilter) ([]types.Event, error) {
	req := horizonclient.TransactionRequest{
		ForLedger: uint(eventFilter.LedgerSeq),
	}
	txnPage, err := cl.horizon.Transactions(req)
	if err != nil {
		return nil, err
	}

	var events []types.Event
	for _, txn := range txnPage.Embedded.Records {
		var txnMeta xdr.TransactionMeta
		if err := xdr.SafeUnmarshalBase64(txn.ResultMetaXdr, &txnMeta); err != nil {
			return nil, err
		}
		if txnMeta.V3 == nil || txnMeta.V3.SorobanMeta == nil {
			continue
		}
		for _, ev := range txnMeta.V3.SorobanMeta.Events {
			hexBytes, err := hex.DecodeString(ev.ContractId.HexString())
			if err != nil {
				break
			}
			contractID, err := strkey.Encode(strkey.VersionByteContract, hexBytes)
			if err != nil {
				return nil, err
			}
			if slices.Contains(eventFilter.ContractIds, contractID) {
				for _, topic := range ev.Body.V0.Topics {
					if slices.Contains(eventFilter.Topics, topic.String()) {
						events = append(events, types.Event{
							ContractEvent: ev,
							LedgerSeq:     uint64(txn.Ledger),
						})
						break
					}
				}
			}
		}
	}

	return events, nil
}
