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
	"github.com/stellar/go/xdr"
)

type IClient interface {
	SimulateTransaction(txXDR string) (*sorobanclient.TxSimulationResult, error)

	SubmitTransactionXDR(txXDR string) (horizon.Transaction, error)

	GetTransaction(txHash string) (horizon.Transaction, error)

	AccountDetail(addr string) (horizon.Account, error)

	GetLatestLedger(ctx context.Context) (*sorobanclient.LatestLedgerResponse, error)

	FetchEvents(ctx context.Context, eventFilter types.EventFilter) ([]types.Event, error)

	StreamEvents(ctx context.Context, eventFilter types.EventFilter, eventChannel chan<- types.Event)

	ParseTxnEvents(txn *horizon.Transaction, fl types.EventFilter) ([]types.Event, error)
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

func (cl *Client) GetTransaction(txHash string) (horizon.Transaction, error) {
	return cl.horizon.TransactionDetail(txHash)
}

func (cl *Client) AccountDetail(addr string) (horizon.Account, error) {
	return cl.horizon.AccountDetail(horizonclient.AccountRequest{AccountID: addr})
}

func (cl *Client) GetLatestLedger(ctx context.Context) (*sorobanclient.LatestLedgerResponse, error) {
	return cl.soroban.GetLatestLedger(ctx)
}

func (cl *Client) StreamEvents(ctx context.Context, eventFilter types.EventFilter, eventChannel chan<- types.Event) {
	ledger, err := cl.horizon.LedgerDetail(uint32(eventFilter.LedgerSeq))
	if err != nil {
		ctx.Done()
		return
	}
	ledgerCursor := ledger.PagingToken()
	trRequest := horizonclient.TransactionRequest{
		Cursor:        ledgerCursor,
		Order:         horizonclient.OrderAsc,
		IncludeFailed: false,
	}
	txnHandler := func(txn horizon.Transaction) {
		var txnMeta xdr.TransactionMeta
		if err := xdr.SafeUnmarshalBase64(txn.ResultMetaXdr, &txnMeta); err != nil {
			return
		}
		if txnMeta.V3 == nil || txnMeta.V3.SorobanMeta == nil {
			//update the processed height
			eventChannel <- types.Event{
				LedgerSeq: uint64(txn.Ledger),
			}
			return
		}
		if len(txnMeta.V3.SorobanMeta.Events) == 0 {
			eventChannel <- types.Event{
				LedgerSeq: uint64(txn.Ledger),
			}
		}
		for _, ev := range txnMeta.V3.SorobanMeta.Events {
			hexBytes, err := hex.DecodeString(ev.ContractId.HexString())
			if err != nil {
				break
			}
			contractID, err := strkey.Encode(strkey.VersionByteContract, hexBytes)
			if err != nil {
				return
			}
			if slices.Contains(eventFilter.ContractIds, contractID) {
				for _, topic := range ev.Body.V0.Topics {
					if slices.Contains(eventFilter.Topics, topic.String()) {
						eventChannel <- types.Event{
							ContractEvent: &ev,
							LedgerSeq:     uint64(txn.Ledger),
						}
						break
					}
				}
			} else {
				eventChannel <- types.Event{
					LedgerSeq: uint64(txn.Ledger),
				}
			}
		}
	}
	err = cl.horizon.StreamTransactions(ctx, trRequest, txnHandler)
	if err != nil {
		ctx.Done()
		return
	}
}

func (cl *Client) FetchEvents(ctx context.Context, eventFilter types.EventFilter) ([]types.Event, error) {
	req := horizonclient.TransactionRequest{
		ForLedger:     uint(eventFilter.LedgerSeq),
		IncludeFailed: false,
	}
	txnPage, err := cl.horizon.Transactions(req)
	if err != nil {
		return nil, err
	}

	var allEvents []types.Event
	for _, txn := range txnPage.Embedded.Records {
		events, err := cl.ParseTxnEvents(&txn, eventFilter)
		if err != nil {
			return allEvents, err
		}
		allEvents = append(allEvents, events...)
	}

	return allEvents, nil
}

func (cl *Client) ParseTxnEvents(txn *horizon.Transaction, fl types.EventFilter) ([]types.Event, error) {
	var events []types.Event
	var txnMeta xdr.TransactionMeta
	if err := xdr.SafeUnmarshalBase64(txn.ResultMetaXdr, &txnMeta); err != nil {
		return nil, err
	}
	if txnMeta.V3 == nil || txnMeta.V3.SorobanMeta == nil {
		return events, nil
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
		if slices.Contains(fl.ContractIds, contractID) {
			for _, topic := range ev.Body.V0.Topics {
				if slices.Contains(fl.Topics, topic.String()) {
					events = append(events, types.Event{
						ContractEvent: &ev,
						LedgerSeq:     uint64(txn.Ledger),
					})
					break
				}
			}
		}
	}

	return events, nil
}
