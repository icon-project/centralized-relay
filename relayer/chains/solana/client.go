package solana

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

type IClient interface {
	GetLatestBlockHeight(ctx context.Context, ctype solrpc.CommitmentType) (uint64, error)
	GetLatestSlot(ctx context.Context, ctype solrpc.CommitmentType) (uint64, error)
	GetLatestBlockHash(ctx context.Context) (*solana.Hash, error)
	GetBlock(ctx context.Context, slot uint64) (*solrpc.GetBlockResult, error)

	GetAccountInfoRaw(ctx context.Context, addr solana.PublicKey) (*solrpc.Account, error)
	GetAccountInfo(ctx context.Context, acAddr solana.PublicKey, accPtr interface{}) error

	GetBalance(ctx context.Context, accAddr solana.PublicKey) (*solrpc.GetBalanceResult, error)

	FetchIDL(ctx context.Context, idlAddress string, idlPtr interface{}) error

	GetSignatureStatus(
		ctx context.Context,
		searchTxHistory bool,
		sign solana.Signature,
	) (*solrpc.SignatureStatusesResult, error)

	GetMinBalanceForRentExemption(
		ctx context.Context,
		dataSize uint64,
	) (uint64, error)

	SimulateTx(
		ctx context.Context,
		tx *solana.Transaction,
	) (*solrpc.SimulateTransactionResult, error)

	SendTx(
		ctx context.Context,
		tx *solana.Transaction,
	) (solana.Signature, error)

	GetSignaturesForAddress(
		ctx context.Context,
		account solana.PublicKey,
		opts *solrpc.GetSignaturesForAddressOpts,
	) ([]*solrpc.TransactionSignature, error)

	GetTransaction(
		ctx context.Context,
		signature solana.Signature,
		opts *solrpc.GetTransactionOpts,
	) (*solrpc.GetTransactionResult, error)
}

type Client struct {
	rpc *solrpc.Client
}

func NewClient(rpcCl *solrpc.Client) IClient {
	return Client{rpc: rpcCl}
}

func (cl Client) GetAccountInfoRaw(ctx context.Context, addr solana.PublicKey) (*solrpc.Account, error) {
	res, err := cl.rpc.GetAccountInfo(ctx, addr)
	if err != nil {
		return nil, err
	}
	return res.Value, nil
}

func (cl Client) GetAccountInfo(ctx context.Context, acAddr solana.PublicKey, accPtr interface{}) error {
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

func (cl Client) GetBalance(ctx context.Context, accAddr solana.PublicKey) (*solrpc.GetBalanceResult, error) {
	return cl.rpc.GetBalance(ctx, accAddr, solrpc.CommitmentFinalized)
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

	idlPubkey, err := solana.PublicKeyFromBase58(idlAddress)
	if err != nil {
		return err
	}

	idlAccount, err := cl.GetAccountInfoRaw(context.Background(), idlPubkey)
	if err != nil {
		return err
	}

	data := idlAccount.Data.GetBinary()[8:] //skip discriminator

	idlAcInfo := struct {
		Authority solana.PublicKey
		DataLen   uint32
	}{}
	if err := borsh.Deserialize(&idlAcInfo, data); err != nil {
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

func (cl Client) GetLatestBlockHeight(ctx context.Context, ctype solrpc.CommitmentType) (uint64, error) {
	return cl.rpc.GetBlockHeight(ctx, ctype)
}

func (cl Client) GetLatestSlot(ctx context.Context, ctype solrpc.CommitmentType) (uint64, error) {
	return cl.rpc.GetSlot(ctx, ctype)
}

func (cl Client) GetBlock(ctx context.Context, slot uint64) (*solrpc.GetBlockResult, error) {
	return cl.rpc.GetBlock(ctx, slot)
}

func (cl Client) GetLatestBlockHash(ctx context.Context) (*solana.Hash, error) {
	hashRes, err := cl.rpc.GetLatestBlockhash(ctx, solrpc.CommitmentConfirmed)
	if err != nil {
		return nil, err
	}
	return &hashRes.Value.Blockhash, nil
}

func (cl Client) SimulateTx(
	ctx context.Context,
	tx *solana.Transaction,
) (*solrpc.SimulateTransactionResult, error) {
	res, err := cl.rpc.SimulateTransactionWithOpts(ctx, tx, &solrpc.SimulateTransactionOpts{
		Commitment: solrpc.CommitmentConfirmed,
	})
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
) (solana.Signature, error) {
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
