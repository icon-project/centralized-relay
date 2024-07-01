package solana

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (cl Client) GetEvent(ctx context.Context, account solana.PublicKey, sno string) (string, error) {
	limit := 1000
	opts := &solrpc.GetSignaturesForAddressOpts{
		Limit: &limit,
	}
	txSignstaure, err := cl.GetSignaturesForAddress(ctx, account, opts)
	if err != nil {
		return "", err
	}
	txOpts := &solrpc.GetTransactionOpts{}
	for _, txSignature := range txSignstaure {
		txResult, err := cl.GetTransaction(ctx, txSignature.Signature, txOpts)
		if err != nil {
			return "", err
		}
		if txResult.Meta != nil && len(txResult.Meta.LogMessages) > 0 {
			event := SolEvent{Slot: txResult.Slot, Signature: txSignature.Signature, Logs: txResult.Meta.LogMessages}
			for _, log := range event.Logs {
				if strings.HasPrefix(log, EventLogPrefix) {
					eventLog := strings.Replace(log, EventLogPrefix, "", 1)
					eventLogBytes, err := base64.StdEncoding.DecodeString(eventLog)
					if err != nil {
						return "", err
					}

					if len(eventLogBytes) < 8 {
						return "", fmt.Errorf("decoded bytes too short to contain discriminator: %v", eventLogBytes)
					}

					discriminator := eventLogBytes[:8]
					eventBytes := eventLogBytes[8:]
					_ = discriminator
					_ = eventBytes
					// TODO : get xcallIdl and check for passed sn
					// for _, ev := range p.xcallIdl.Events {
					// 	if slices.Equal(ev.Discriminator, discriminator) {
					// 		switch ev.Name {
					// 		case EventSendMessage:
					// 			smEvent := SendMessageEvent{}
					// 			if err := borsh.Deserialize(&smEvent, eventBytes); err != nil {
					// 				return "", fmt.Errorf("failed to decode send message event: %w", err)
					// 			}

					// 		}
					// 	}
					// }
				}
			}

		}
	}
	return "", errors.New("event not found")

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

func idlAddrFromProgID(progID string) (string, error) {
	progPubkey, err := solana.PublicKeyFromBase58(progID)
	if err != nil {
		return "", err
	}

	basePubkey, _, err := solana.FindProgramAddress([][]byte{}, progPubkey)
	if err != nil {
		return "", err
	}

	calculatedIdlAddr, err := solana.CreateWithSeed(basePubkey, "anchor:idl", progPubkey)
	if err != nil {
		return "", err
	}

	return calculatedIdlAddr.String(), nil
}

func (cl Client) FetchIDL(ctx context.Context, progID string, idlPtr interface{}) error {

	idlAddress, err := idlAddrFromProgID(progID)
	if err != nil {
		return err
	}
	fmt.Println("addr", idlAddress)

	idlAccount, err := cl.GetAccountInfoRaw(context.Background(), progID)
	if err != nil {
		return err
	}
	data := idlAccount.Data.GetBinary()[8:] //skip discriminator
	fmt.Println("new Data is", data)
	idlAcInfo := struct {
		Authority solana.PublicKey
		DataLen   uint32
	}{}
	var idl interface{}
	if err := borsh.Deserialize(&idl, data); err != nil {
		return err
	}

	compressedBytes := data[36 : 36+idlAcInfo.DataLen] //skip authority and unwanted trailing bytes

	decompressedBytes, err := decompress(compressedBytes)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(decompressedBytes, idlPtr); err != nil {
		return err
	}

	return nil
}

func decompress(compressedData []byte) ([]byte, error) {
	// Create a new bytes reader from the compressed data
	b := bytes.NewReader(compressedData)

	// Create a new zlib reader
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// Create a buffer to hold the decompressed data
	var out bytes.Buffer

	// Copy the decompressed data into the buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return nil, err
	}

	// Return the decompressed data
	return out.Bytes(), nil
}

func (cl Client) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	return cl.rpc.GetBlockHeight(ctx, solrpc.CommitmentFinalized)
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

func (cl Client) GetMinBalanceForRentExemption(
	ctx context.Context,
	dataSize uint64,
) (uint64, error) {
	return cl.rpc.GetMinimumBalanceForRentExemption(ctx, dataSize, "")
}

// func createTxn() {
// 	discriminator, err := p.xcallIdl.GetInstructionDiscriminator(types.MethodRecvMessage)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	toArg, err := borsh.Serialize(msg.Dst)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	msgArg, err := borsh.Serialize(msg.Data)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	progID, err := p.xcallIdl.GetProgramID()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	payerAccount := solana.AccountMeta{
// 		PublicKey:  p.wallet.PublicKey(),
// 		IsWritable: true,
// 		IsSigner:   true,
// 	}

// 	xcallStatePubKey, err := solana.PublicKeyFromBase58(p.cfg.XcallStateAccount)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	xcallStateAccount := solana.AccountMeta{
// 		PublicKey:  xcallStatePubKey,
// 		IsWritable: true,
// 	}

// 	instructionData := append(discriminator, toArg...)
// 	instructionData = append(instructionData, msgArg...)

// 	instructions := []solana.Instruction{
// 		&solana.GenericInstruction{
// 			ProgID:        progID,
// 			AccountValues: solana.AccountMetaSlice{&payerAccount, &xcallStateAccount},
// 			DataBytes:     instructionData,
// 		},
// 	}

// 	return instructions, []solana.PrivateKey{p.wallet.PrivateKey}, nil
// }
