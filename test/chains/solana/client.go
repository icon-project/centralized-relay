package solana

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

type Client struct {
	rpc *solrpc.Client
}

func New(rpcUrl string, rpcCl *solrpc.Client) (*Client, error) {
	return &Client{rpc: rpcCl}, nil
}

func (cl Client) GetAccountInfoRaw(ctx context.Context, addr string) (*solrpc.Account, error) {
	res, err := cl.rpc.GetAccountInfo(ctx, solana.MustPublicKeyFromBase58(addr))
	if err != nil {
		return nil, err
	}
	return res.Value, nil
}

func (cl Client) GetSignaturesForAddress(
	ctx context.Context,
	account solana.PublicKey,
	opts *solrpc.GetSignaturesForAddressOpts,
) ([]*solrpc.TransactionSignature, error) {
	return cl.rpc.GetSignaturesForAddressWithOpts(ctx, account, opts)
}

func (cl Client) GetTransaction(
	ctx context.Context,
	signature solana.Signature,
	opts *solrpc.GetTransactionOpts,
) (*solrpc.GetTransactionResult, error) {
	return cl.rpc.GetTransaction(ctx, signature, opts)
}

func (cl Client) GetEvent(ctx context.Context, account solana.PublicKey,
	allEvents []IdlEvent, signature, sno string) (*EventResponseEvent, error) {
	limit := 1000
	opts := &solrpc.GetSignaturesForAddressOpts{
		Limit: &limit,
	}
	txSignstaure, err := cl.GetSignaturesForAddress(ctx, account, opts)
	if err != nil {
		return nil, err
	}
	txVersion := uint64(0)
	txOpts := &solrpc.GetTransactionOpts{MaxSupportedTransactionVersion: &txVersion}
	evRespEvents := EventResponseEvent{}
	evRespEvents.ValueDecoded = make(map[string]interface{})
	for _, txSignature := range txSignstaure {
		txResult, err := cl.GetTransaction(ctx, txSignature.Signature, txOpts)
		if err != nil {
			return nil, err
		}
		if txResult.Meta != nil && len(txResult.Meta.LogMessages) > 0 {
			event := SolEvent{Slot: txResult.Slot, Signature: txSignature.Signature, Logs: txResult.Meta.LogMessages}
			for _, log := range event.Logs {
				if strings.HasPrefix(log, EventLogPrefix) {
					eventLog := strings.Replace(log, EventLogPrefix, "", 1)
					eventLogBytes, err := base64.StdEncoding.DecodeString(eventLog)
					if err != nil {
						return nil, err
					}

					if len(eventLogBytes) < 8 {
						return nil, fmt.Errorf("decoded bytes too short to contain discriminator: %v", eventLogBytes)
					}

					discriminator := eventLogBytes[:8]
					eventBytes := eventLogBytes[8:]
					for _, ev := range allEvents {
						if slices.Equal(ev.Discriminator, discriminator) {
							if ev.Name == signature {
								if signature == EventCallMessage {
									smEvent := CallMessageEvent{}
									if err := borsh.Deserialize(&smEvent, eventBytes); err != nil {
										return nil, fmt.Errorf("failed to decode send message event: %w", err)
									}
									if strconv.Itoa(int(smEvent.Sn.Int64())) == sno {
										evRespEvents.ValueDecoded["sn"] = strconv.Itoa(int(smEvent.Sn.Int64()))
										evRespEvents.ValueDecoded["reqId"] = strconv.Itoa(int(smEvent.ReqId.Int64()))
										evRespEvents.ValueDecoded["data"] = smEvent.Data
										evRespEvents.ValueDecoded["to"] = smEvent.To
										return &evRespEvents, nil
									}
								} else if signature == EventRollbackExecuted {
									smEvent := RollbackExecuted{}
									if err := borsh.Deserialize(&smEvent, eventBytes); err != nil {
										return nil, fmt.Errorf("failed to decode send message event: %w", err)
									}
									if strconv.Itoa(int(smEvent.Sn.Int64())) == sno {
										evRespEvents.ValueDecoded["sn"] = strconv.Itoa(int(smEvent.Sn.Int64()))
										return &evRespEvents, nil
									}
								} else if signature == EventResponseMessage {
									smEvent := ResponseMessage{}
									if err := borsh.Deserialize(&smEvent, eventBytes); err != nil {
										return nil, fmt.Errorf("failed to decode send message event: %w", err)
									}
									if strconv.Itoa(int(smEvent.Sn.Int64())) == sno {
										evRespEvents.ValueDecoded["sn"] = strconv.Itoa(int(smEvent.Sn.Int64()))
										evRespEvents.ValueDecoded["code"] = strconv.Itoa(int(smEvent.Code))
										return &evRespEvents, nil
									}
								}
							}
						}
					}
				}
			}

		}
	}
	return nil, errors.New("event not found")

}
func (cl Client) GetAccountInfo(ctx context.Context, acAddr string, accPtr interface{}) error {
	ac, err := cl.GetAccountInfoRaw(ctx, acAddr)
	if err != nil {
		return err
	}
	data := ac.Data.GetBinary()[8:] //skip discriminator

	if err := borsh.Deserialize(accPtr, data); err != nil {
		return fmt.Errorf("failed to deserialize to account ptr: %w", err)
	}

	return nil
}

func (cl Client) GetLatestBlockHeight(ctx context.Context, ctype solrpc.CommitmentType) (uint64, error) {
	return cl.rpc.GetBlockHeight(ctx, ctype)
}

func (cl Client) GetLatestBlockHash(ctx context.Context) (*solana.Hash, error) {
	hashRes, err := cl.rpc.GetLatestBlockhash(ctx, solrpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	return &hashRes.Value.Blockhash, nil
}

func (cl Client) SimulateTx(
	ctx context.Context,
	tx *solana.Transaction,
	opts *solrpc.SimulateTransactionOpts,
) (*solrpc.SimulateTransactionResult, error) {
	res, err := cl.rpc.SimulateTransactionWithOpts(ctx, tx, opts)
	if err != nil {
		return nil, err
	}

	if res.Value.Err != nil {
		return nil, fmt.Errorf("failed to simulate tx: %v", res.Value.Err)
	}

	return res.Value, nil
}

func (cl Client) SendTx(
	ctx context.Context,
	tx *solana.Transaction,
	opts *solrpc.TransactionOpts,
) (solana.Signature, error) {
	if opts != nil {
		return cl.rpc.SendTransactionWithOpts(ctx, tx, *opts)
	}
	return cl.rpc.SendTransaction(ctx, tx)
}

func (cl Client) GetSignatureStatus(
	ctx context.Context,
	searchTxHistory bool,
	sign solana.Signature,
) (*solrpc.SignatureStatusesResult, error) {
	res, err := cl.rpc.GetSignatureStatuses(ctx, searchTxHistory, sign)
	if err != nil {
		return nil, err
	}
	if len(res.Value) > 0 {
		return res.Value[0], nil
	}
	return nil, fmt.Errorf("tx signature result not found")
}
