package solana

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/icon-project/centralized-relay/relayer/chains/solana/alt"
	"github.com/icon-project/centralized-relay/relayer/chains/solana/types"
	relayerevents "github.com/icon-project/centralized-relay/relayer/events"
	relayertypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/near/borsh-go"
	"go.uber.org/zap"
)

const (
	defaultTxConfirmationTime = 2 * time.Second
	accountsQueryMaxLimit     = 30
)

func (p *Provider) Route(ctx context.Context, message *relayertypes.Message, callback relayertypes.TxResponseFunc) error {
	if err := p.RestoreKeystore(ctx); err != nil {
		return err
	}

	p.log.Info("starting to route message",
		zap.String("src", message.Src),
		zap.String("dst", message.Dst),
		zap.Any("sn", message.Sn),
		zap.Any("req-id", message.ReqID),
		zap.String("event-type", message.EventType),
		zap.String("data", hex.EncodeToString(message.Data)),
	)

	instructions, signers, err := p.MakeCallInstructions(message)
	if err != nil {
		return fmt.Errorf("failed to create call instructions: %w", err)
	}

	// addressTables, err := p.prepareAddressTablesForInstructions(instructions)
	// if err != nil {
	// 	return fmt.Errorf("failed to prepare address lookup tables for instructions: %w", err)
	// }

	opts := []solana.TransactionOption{
		solana.TransactionPayer(p.wallet.PublicKey()),
		solana.TransactionAddressTables(p.staticAlts),
	}

	tx, err := p.prepareTx(ctx, instructions, signers, opts...)
	if err != nil {
		return fmt.Errorf("failed to prepare and simulate tx: %w", err)
	}

	txSign, err := p.client.SendTx(ctx, tx, nil)
	if err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	go p.executeRouteCallback(txSign, message, callback)

	return nil
}

func (p *Provider) prepareTx(
	ctx context.Context,
	instructions []solana.Instruction,
	signers []solana.PrivateKey,
	opts ...solana.TransactionOption,
) (*solana.Transaction, error) {
	latestBlockHash, err := p.client.GetLatestBlockHash(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block hash: %w", err)
	}

	tx, err := solana.NewTransaction(instructions, *latestBlockHash, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new tx: %w", err)
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			for _, signer := range signers {
				if signer.PublicKey() == key {
					return &signer
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign tx: %w", err)
	}

	return tx, nil
}

func (p *Provider) waitForTxConfirmation(timeout time.Duration, sign solana.Signature) (*solrpc.SignatureStatusesResult, error) {
	startTime := time.Now()
	for range time.NewTicker(500 * time.Millisecond).C {
		txStatus, err := p.client.GetSignatureStatus(context.TODO(), false, sign)
		if err == nil && txStatus != nil && (txStatus.ConfirmationStatus == solrpc.ConfirmationStatusConfirmed || txStatus.ConfirmationStatus == solrpc.ConfirmationStatusFinalized) {
			return txStatus, nil
		} else if time.Since(startTime) > timeout {
			var cbErr error
			if err != nil {
				cbErr = err
			} else if txStatus != nil && txStatus.Err != nil {
				cbErr = fmt.Errorf("failed to get tx signature status: %v", txStatus.Err)
			} else {
				cbErr = fmt.Errorf("failed to finalize tx signature")
			}
			return nil, cbErr
		}
	}
	return nil, fmt.Errorf("request timeout")
}

func (p *Provider) executeRouteCallback(
	sign solana.Signature,
	msg *relayertypes.Message,
	callback relayertypes.TxResponseFunc,
) {
	txResult, err := p.waitForTxConfirmation(defaultTxConfirmationTime, sign)
	if err != nil {
		callback(
			msg.MessageKey(),
			&relayertypes.TxResponse{
				TxHash: sign.String(),
			},
			err,
		)
	} else {
		callback(
			msg.MessageKey(),
			&relayertypes.TxResponse{
				Height: int64(txResult.Slot),
				TxHash: sign.String(),
				Code:   relayertypes.Success,
			},
			nil,
		)
	}
}

func (p *Provider) CreateLookupTableAccount(ctx context.Context) (*solana.PublicKey, error) {
	recentSlot, err := p.client.GetLatestBlockHeight(ctx)
	if err != nil {
		return nil, err
	}
	altCreateInstruction, accountAddr, err := alt.CreateLookupTable(
		p.wallet.PublicKey(),
		p.wallet.PublicKey(),
		recentSlot,
	)
	if err != nil {
		return nil, err
	}

	signers := []solana.PrivateKey{p.wallet.PrivateKey}

	tx, err := p.prepareTx(
		context.Background(),
		[]solana.Instruction{altCreateInstruction},
		signers,
		solana.TransactionPayer(p.wallet.PublicKey()),
	)
	if err != nil {
		return nil, err
	}

	txSign, err := p.client.SendTx(ctx, tx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send tx: %w", err)
	}

	_, err = p.waitForTxConfirmation(defaultTxConfirmationTime, txSign)
	if err != nil {
		return nil, err
	}

	return &accountAddr, nil
}

func (p *Provider) ExtendLookupTableAccount(ctx context.Context, acTableAddr solana.PublicKey, addresses solana.PublicKeySlice) error {
	payer := p.wallet.PublicKey()
	altExtendInstruction := alt.ExtendLookupTable(
		acTableAddr,
		p.wallet.PublicKey(),
		&payer,
		addresses,
	)

	signers := []solana.PrivateKey{p.wallet.PrivateKey}

	tx, err := p.prepareTx(
		context.Background(),
		[]solana.Instruction{altExtendInstruction},
		signers,
		solana.TransactionPayer(p.wallet.PublicKey()),
	)
	if err != nil {
		return err
	}

	txSign, err := p.client.SendTx(ctx, tx, nil)
	if err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	_, err = p.waitForTxConfirmation(defaultTxConfirmationTime, txSign)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) GetLookupTableAccount(accountID solana.PublicKey) (*alt.LookupTableAccount, error) {
	acInfo, err := p.client.GetAccountInfoRaw(context.Background(), accountID)
	if err != nil {
		return nil, err
	}

	account, err := alt.DeserializeLookupTable(acInfo.Data.GetBinary())
	if err != nil {
		return nil, err
	}

	if account.ProgramState == alt.ProgramStateUninitialized {
		return nil, fmt.Errorf("account program not initialized")
	}

	if !account.IsActive() {
		return nil, fmt.Errorf("account deactivated")
	}

	return account, nil
}

func (p *Provider) initStaticAlts() error {
	xcallProgID, err := solana.PublicKeyFromBase58(p.cfg.XcallProgram)
	if err != nil {
		return err
	}

	addresses := solana.PublicKeySlice{solana.SystemProgramID, xcallProgID}

	connections := append([]string{p.cfg.ConnectionProgram}, p.cfg.OtherConnections...)
	for _, conn := range connections {
		connProgID, err := solana.PublicKeyFromBase58(conn)
		if err != nil {
			return err
		}

		connConfigAddr, err := types.GetPDA(connProgID, types.PrefixConfig)
		if err != nil {
			return err
		}

		connClaimFeeAddr, err := types.GetPDA(connProgID, types.PrefixClaimFees)
		if err != nil {
			return err
		}

		nidFees := solana.PublicKeySlice{}
		for _, net := range p.cfg.CpNIDs {
			netFeeAddr, err := types.GetPDA(connProgID, types.PrefixNetworkFee, []byte(net))
			if err != nil {
				return err
			}
			nidFees = append(nidFees, netFeeAddr)
		}

		addresses = append(addresses, solana.PublicKeySlice{
			connProgID, connConfigAddr, connClaimFeeAddr,
		}...)
		addresses = append(addresses, nidFees...)
	}

	xcallConfigAddr, err := p.pdaRegistry.XcallConfig.GetAddress()
	if err != nil {
		return err
	}

	addresses = append(addresses, solana.PublicKeySlice{
		xcallConfigAddr,
	}...)

	altPubKey, err := solana.PublicKeyFromBase58(p.cfg.AltAddress)
	if err != nil {
		return err
	}
	altAc, err := p.GetLookupTableAccount(altPubKey)
	if err != nil {
		return err
	}

	addressesToExtend := solana.PublicKeySlice{}
	for _, addr := range addresses {
		if !altAc.Addresses.Contains(addr) {
			addressesToExtend = append(addressesToExtend, addr)
		}
	}

	if len(addressesToExtend) > 0 {
		if err := p.ExtendLookupTableAccount(context.Background(), altPubKey, addressesToExtend); err != nil {
			return err
		}
	}

	p.staticAlts[altPubKey] = append(altAc.Addresses, addressesToExtend...)

	return nil
}

func (p *Provider) MakeCallInstructions(msg *relayertypes.Message) ([]solana.Instruction, []solana.PrivateKey, error) {
	switch msg.EventType {
	case relayerevents.EmitMessage:
		instructions, signers, err := p.getRecvMessageIntruction(msg)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get recv message instructions")
		}
		return instructions, signers, nil
	case relayerevents.CallMessage:
		instructions, signers, err := p.getExecuteCallInstruction(msg)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get execute call instructions: %w", err)
		}
		return instructions, signers, nil
	case relayerevents.RollbackMessage:
		instructions, signers, err := p.getExecuteRollbackInstruction(msg)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get execute rollback instructions: %w", err)
		}
		return instructions, signers, nil
	default:
		return nil, nil, fmt.Errorf("invalid event type in message")
	}
}

func (p *Provider) getRecvMessageIntruction(msg *relayertypes.Message) ([]solana.Instruction, []solana.PrivateKey, error) {
	discriminator, err := p.connIdl.GetInstructionDiscriminator(types.MethodRecvMessage)
	if err != nil {
		return nil, nil, err
	}

	srcArg, err := borsh.Serialize(msg.Src)
	if err != nil {
		return nil, nil, err
	}

	connSnArg, err := borsh.Serialize(*msg.Sn)
	if err != nil {
		return nil, nil, err
	}

	dataArg, err := borsh.Serialize(msg.Data)
	if err != nil {
		return nil, nil, err
	}

	csMessage, err := p.decodeCsMessage(context.Background(), msg.Data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode cs message: %w", err)
	}

	var sn big.Int
	if csMessage.MessageType == types.CsMessageRequest {
		sn = csMessage.Request.SequenceNo
	} else {
		sn = csMessage.Result.SequenceNo
	}

	snArg, err := borsh.Serialize(sn)
	if err != nil {
		return nil, nil, err
	}

	instructionData := append(discriminator, srcArg...)
	instructionData = append(instructionData, connSnArg...)
	instructionData = append(instructionData, dataArg...)
	instructionData = append(instructionData, snArg...)

	accounts := solana.AccountMetaSlice{
		&solana.AccountMeta{
			PublicKey:  p.wallet.PublicKey(),
			IsWritable: true,
			IsSigner:   true,
		},
	}

	xcallConfigAddr, err := p.pdaRegistry.XcallConfig.GetAddress()
	if err != nil {
		return nil, nil, err
	}

	xcallConfigAc := types.XcallConfigAccount{}
	if err := p.client.GetAccountInfo(context.Background(), xcallConfigAddr, &xcallConfigAc); err != nil {
		return nil, nil, err
	}

	recvMessageAccounts, err := p.fetchRecvMessageAccounts(msg, sn, csMessage.MessageType)
	if err != nil {
		return nil, nil, err
	}

	accounts = append(accounts, recvMessageAccounts...)

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID:        p.connIdl.GetProgramID(),
			AccountValues: accounts,
			DataBytes:     instructionData,
		},
	}

	return instructions, []solana.PrivateKey{p.wallet.PrivateKey}, nil
}

func (p *Provider) getExecuteCallInstruction(msg *relayertypes.Message) ([]solana.Instruction, []solana.PrivateKey, error) {
	discriminator, err := p.xcallIdl.GetInstructionDiscriminator(types.MethodExecuteCall)
	if err != nil {
		return nil, nil, err
	}

	reqIdBytes, err := borsh.Serialize(*msg.ReqID)
	if err != nil {
		return nil, nil, err
	}

	dataBytes, err := borsh.Serialize(msg.Data)
	if err != nil {
		return nil, nil, err
	}

	srcNIDBytes, err := borsh.Serialize(msg.Src)
	if err != nil {
		return nil, nil, err
	}

	instructionData := append(discriminator, reqIdBytes...)
	instructionData = append(instructionData, dataBytes...)
	instructionData = append(instructionData, srcNIDBytes...)

	accounts := solana.AccountMetaSlice{
		&solana.AccountMeta{
			PublicKey:  p.wallet.PublicKey(),
			IsWritable: true,
			IsSigner:   true,
		},
	}

	executeCallAccounts, err := p.fetchExecuteCallAccounts(msg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch execute call accounts: %w", err)
	}
	accounts = append(accounts, executeCallAccounts...)

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID:        p.xcallIdl.GetProgramID(),
			AccountValues: accounts,
			DataBytes:     instructionData,
		},
	}

	return instructions, []solana.PrivateKey{p.wallet.PrivateKey}, nil
}

func (p *Provider) getExecuteRollbackInstruction(msg *relayertypes.Message) ([]solana.Instruction, []solana.PrivateKey, error) {
	discriminator, err := p.xcallIdl.GetInstructionDiscriminator(types.MethodExecuteRollback)
	if err != nil {
		return nil, nil, err
	}

	snBytes, err := borsh.Serialize(*msg.Sn)
	if err != nil {
		return nil, nil, err
	}

	instructionData := append(discriminator, snBytes...)

	accounts := solana.AccountMetaSlice{
		&solana.AccountMeta{
			PublicKey:  p.wallet.PublicKey(),
			IsWritable: true,
			IsSigner:   true,
		},
	}

	executeRollbackAccounts, err := p.fetchExecuteRollbackAccounts(msg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch execute rollback accounts: %w", err)
	}
	accounts = append(accounts, executeRollbackAccounts...)

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID:        p.xcallIdl.GetProgramID(),
			AccountValues: accounts,
			DataBytes:     instructionData,
		},
	}

	return instructions, []solana.PrivateKey{p.wallet.PrivateKey}, nil
}

func (p *Provider) fetchRecvMessageAccounts(
	msg *relayertypes.Message,
	xcallSn big.Int,
	msgType types.CsMessageType,
) (solana.AccountMetaSlice, error) {
	accounts := []solana.AccountMeta{}
	page := uint8(1)
	limit := uint8(accountsQueryMaxLimit)

	res, err := p.queryRecvMessageAccounts(msg, xcallSn, msgType, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recv message accounts: %w", err)
	}
	accounts = append(accounts, res.Accounts...)

	for res.HasNextPage {
		page++
		res, err = p.queryRecvMessageAccounts(msg, xcallSn, msgType, page, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to query recv message accounts: %w", err)
		}
		accounts = append(accounts, res.Accounts...)
	}

	acMetaSlice := solana.AccountMetaSlice{}
	for _, acMeta := range accounts {
		acMetaSlice = append(acMetaSlice, &acMeta)
	}

	return acMetaSlice, nil
}

func (p *Provider) queryRecvMessageAccounts(
	msg *relayertypes.Message,
	xcallSn big.Int,
	msgType types.CsMessageType,
	page uint8,
	limit uint8,
) (*types.QueryAccountsResponse, error) {
	discriminator, err := p.connIdl.GetInstructionDiscriminator(types.MethodQueryRecvMessageAccounts)
	if err != nil {
		return nil, err
	}

	srcArg, err := borsh.Serialize(msg.Src)
	if err != nil {
		return nil, err
	}

	connSnArg, err := borsh.Serialize(*msg.Sn)
	if err != nil {
		return nil, err
	}

	dataArg, err := borsh.Serialize(msg.Data)
	if err != nil {
		return nil, err
	}

	snArg, err := borsh.Serialize(xcallSn)
	if err != nil {
		return nil, err
	}

	pageBytes, err := borsh.Serialize(page)
	if err != nil {
		return nil, err
	}

	limitBytes, err := borsh.Serialize(limit)
	if err != nil {
		return nil, err
	}

	instructionData := append(discriminator, srcArg...)
	instructionData = append(instructionData, connSnArg...)
	instructionData = append(instructionData, dataArg...)
	instructionData = append(instructionData, snArg...)
	instructionData = append(instructionData, pageBytes...)
	instructionData = append(instructionData, limitBytes...)

	connConfigAddr, err := p.pdaRegistry.ConnConfig.GetAddress()
	if err != nil {
		return nil, err
	}

	accounts := solana.AccountMetaSlice{
		&solana.AccountMeta{
			PublicKey:  connConfigAddr,
			IsWritable: false,
			IsSigner:   false,
		},
	}

	xcallConfigAddr, err := p.pdaRegistry.XcallConfig.GetAddress()
	if err != nil {
		return nil, err
	}

	accounts = append(accounts, &solana.AccountMeta{
		PublicKey:  xcallConfigAddr,
		IsWritable: false,
		IsSigner:   false,
	})

	if msgType == types.CsMessageResult {
		rollbackAddr, err := p.pdaRegistry.XcallRollback.GetAddress(xcallSn.FillBytes(make([]byte, 16)))
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &solana.AccountMeta{
			PublicKey:  rollbackAddr,
			IsWritable: false,
			IsSigner:   false,
		})
	}

	accounts = append(accounts, &solana.AccountMeta{
		PublicKey:  p.xcallIdl.GetProgramID(),
		IsWritable: false,
		IsSigner:   false,
	})

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID:        p.connIdl.GetProgramID(),
			AccountValues: accounts,
			DataBytes:     instructionData,
		},
	}

	signers := []solana.PrivateKey{p.wallet.PrivateKey}

	tx, err := p.prepareTx(
		context.Background(),
		instructions,
		signers,
		solana.TransactionPayer(p.wallet.PublicKey()),
	)
	if err != nil {
		return nil, err
	}

	txSign, err := p.client.SendTx(context.Background(), tx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send tx: %w", err)
	}

	_, err = p.waitForTxConfirmation(defaultTxConfirmationTime, txSign)
	if err != nil {
		return nil, err
	}

	txnres, err := p.client.GetTransaction(context.Background(), txSign, &solrpc.GetTransactionOpts{Commitment: solrpc.CommitmentConfirmed})
	if err != nil {
		return nil, fmt.Errorf("failed to get txn %s: %w", txSign.String(), err)
	}

	acRes := types.QueryAccountsResponse{}
	if err := parseReturnValueFromLogs(p.connIdl.GetProgramID().String(), txnres.Meta.LogMessages, &acRes); err != nil {
		return nil, fmt.Errorf("failed to parse return value: %w", err)
	}

	return &acRes, nil
}

func (p *Provider) fetchExecuteCallAccounts(msg *relayertypes.Message) (solana.AccountMetaSlice, error) {
	accounts := []solana.AccountMeta{}
	page := uint8(1)
	limit := uint8(accountsQueryMaxLimit)

	res, err := p.queryExecuteCallAccounts(msg, page, limit)
	if err != nil {
		return nil, err
	}
	accounts = append(accounts, res.Accounts...)

	for res.HasNextPage {
		page++
		res, err = p.queryExecuteCallAccounts(msg, page, limit)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, res.Accounts...)
	}

	acMetaSlice := solana.AccountMetaSlice{}
	for _, acMeta := range accounts {
		acMetaSlice = append(acMetaSlice, &acMeta)
	}

	return acMetaSlice, nil
}

func (p *Provider) queryExecuteCallAccounts(
	msg *relayertypes.Message,
	page uint8,
	limit uint8,
) (*types.QueryAccountsResponse, error) {
	discriminator, err := p.xcallIdl.GetInstructionDiscriminator(types.MethodQueryExecuteCallAccounts)
	if err != nil {
		return nil, err
	}

	reqIdBytes, err := borsh.Serialize(*msg.ReqID)
	if err != nil {
		return nil, err
	}

	dataBytes, err := borsh.Serialize(msg.Data)
	if err != nil {
		return nil, err
	}

	pageBytes, err := borsh.Serialize(page)
	if err != nil {
		return nil, err
	}

	limitBytes, err := borsh.Serialize(limit)
	if err != nil {
		return nil, err
	}

	instructionData := discriminator
	instructionData = append(instructionData, reqIdBytes...)
	instructionData = append(instructionData, dataBytes...)
	instructionData = append(instructionData, pageBytes...)
	instructionData = append(instructionData, limitBytes...)

	xcallConfigAddr, err := p.pdaRegistry.XcallConfig.GetAddress()
	if err != nil {
		return nil, err
	}

	xcallProxyReqAddr, err := p.pdaRegistry.XcallProxyRequest.GetAddress(msg.ReqID.FillBytes(make([]byte, 16)))
	if err != nil {
		return nil, err
	}

	accounts := solana.AccountMetaSlice{
		&solana.AccountMeta{
			PublicKey:  xcallConfigAddr,
			IsWritable: false,
			IsSigner:   false,
		},
		&solana.AccountMeta{
			PublicKey:  xcallProxyReqAddr,
			IsWritable: false,
			IsSigner:   false,
		},
	}

	xcallProxyReqAcc := types.ProxyRequestAccount{}
	if err := p.client.GetAccountInfo(context.Background(), xcallProxyReqAddr, &xcallProxyReqAcc); err != nil {
		return nil, err
	}

	programIds := []solana.PublicKey{}

	for _, connProgID := range xcallProxyReqAcc.ReqMessage.Protocols {
		connPubKey, err := solana.PublicKeyFromBase58(connProgID)
		if err != nil {
			return nil, err
		}
		connConfigAddr, err := types.GetPDA(connPubKey, "config")
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, &solana.AccountMeta{
			PublicKey:  connConfigAddr,
			IsWritable: true,
			IsSigner:   false,
		})

		programIds = append(programIds, connPubKey)

	}

	dappPubKey, err := solana.PublicKeyFromBase58(xcallProxyReqAcc.ReqMessage.To)
	if err != nil {
		return nil, err
	}
	dappConfigAddr, err := types.GetPDA(dappPubKey, "config")
	if err != nil {
		return nil, err
	}
	programIds = append(programIds, dappPubKey)

	accounts = append(accounts, &solana.AccountMeta{
		PublicKey:  dappConfigAddr,
		IsWritable: true,
		IsSigner:   false,
	})

	for _, progID := range programIds {
		accounts = append(accounts, &solana.AccountMeta{
			PublicKey:  progID,
			IsWritable: false,
			IsSigner:   false,
		})
	}

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID:        p.xcallIdl.GetProgramID(),
			AccountValues: accounts,
			DataBytes:     instructionData,
		},
	}

	signers := []solana.PrivateKey{p.wallet.PrivateKey}

	tx, err := p.prepareTx(
		context.Background(),
		instructions,
		signers,
		solana.TransactionPayer(p.wallet.PublicKey()),
	)
	if err != nil {
		return nil, err
	}

	txSign, err := p.client.SendTx(context.Background(), tx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send tx: %w", err)
	}

	_, err = p.waitForTxConfirmation(defaultTxConfirmationTime, txSign)
	if err != nil {
		return nil, err
	}

	txnres, err := p.client.GetTransaction(context.Background(), txSign, &solrpc.GetTransactionOpts{Commitment: solrpc.CommitmentConfirmed})
	if err != nil {
		return nil, fmt.Errorf("failed to get txn %s: %w", txSign.String(), err)
	}

	acRes := types.QueryAccountsResponse{}
	if err := parseReturnValueFromLogs(p.xcallIdl.GetProgramID().String(), txnres.Meta.LogMessages, &acRes); err != nil {
		return nil, fmt.Errorf("failed to parse return value: %w", err)
	}

	return &acRes, nil
}

func (p *Provider) fetchExecuteRollbackAccounts(msg *relayertypes.Message) (solana.AccountMetaSlice, error) {
	accounts := []solana.AccountMeta{}
	page := uint8(1)
	limit := uint8(accountsQueryMaxLimit)

	res, err := p.queryExecuteRollbackAccounts(msg, page, limit)
	if err != nil {
		return nil, err
	}
	accounts = append(accounts, res.Accounts...)

	for res.HasNextPage {
		page++
		res, err = p.queryExecuteRollbackAccounts(msg, page, limit)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, res.Accounts...)
	}

	acMetaSlice := solana.AccountMetaSlice{}
	for _, acMeta := range accounts {
		acMetaSlice = append(acMetaSlice, &acMeta)
	}

	return acMetaSlice, nil
}

func (p *Provider) queryExecuteRollbackAccounts(
	msg *relayertypes.Message,
	page uint8,
	limit uint8,
) (*types.QueryAccountsResponse, error) {
	discriminator, err := p.xcallIdl.GetInstructionDiscriminator(types.MethodQueryExecuteRollbackAccounts)
	if err != nil {
		return nil, err
	}

	snBytes, err := borsh.Serialize(*msg.Sn)
	if err != nil {
		return nil, err
	}

	pageBytes, err := borsh.Serialize(page)
	if err != nil {
		return nil, err
	}

	limitBytes, err := borsh.Serialize(limit)
	if err != nil {
		return nil, err
	}

	instructionData := discriminator
	instructionData = append(instructionData, snBytes...)
	instructionData = append(instructionData, pageBytes...)
	instructionData = append(instructionData, limitBytes...)

	xcallConfigAddr, err := p.pdaRegistry.XcallConfig.GetAddress()
	if err != nil {
		return nil, err
	}

	xcallRollbackAddr, err := p.pdaRegistry.XcallRollback.GetAddress(msg.Sn.FillBytes(make([]byte, 16)))
	if err != nil {
		return nil, err
	}

	accounts := solana.AccountMetaSlice{
		&solana.AccountMeta{
			PublicKey:  xcallConfigAddr,
			IsWritable: false,
			IsSigner:   false,
		},
		&solana.AccountMeta{
			PublicKey:  xcallRollbackAddr,
			IsWritable: false,
			IsSigner:   false,
		},
	}

	xcallRollbackAcc := types.XcallRollbackAccount{}
	if err := p.client.GetAccountInfo(context.Background(), xcallRollbackAddr, &xcallRollbackAcc); err != nil {
		return nil, err
	}

	dappPubKey := xcallRollbackAcc.Rollback.From
	dappConfigAddr, err := types.GetPDA(dappPubKey, "config")
	if err != nil {
		return nil, err
	}
	accounts = append(accounts, &solana.AccountMeta{
		PublicKey:  dappConfigAddr,
		IsWritable: true,
		IsSigner:   false,
	})

	accounts = append(accounts, &solana.AccountMeta{
		PublicKey:  dappPubKey,
		IsWritable: false,
		IsSigner:   false,
	})

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID:        p.xcallIdl.GetProgramID(),
			AccountValues: accounts,
			DataBytes:     instructionData,
		},
	}

	signers := []solana.PrivateKey{p.wallet.PrivateKey}

	tx, err := p.prepareTx(
		context.Background(),
		instructions,
		signers,
		solana.TransactionPayer(p.wallet.PublicKey()),
	)
	if err != nil {
		return nil, err
	}

	txSign, err := p.client.SendTx(context.Background(), tx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send tx: %w", err)
	}

	_, err = p.waitForTxConfirmation(defaultTxConfirmationTime, txSign)
	if err != nil {
		return nil, err
	}

	txnres, err := p.client.GetTransaction(context.Background(), txSign, &solrpc.GetTransactionOpts{Commitment: solrpc.CommitmentConfirmed})
	if err != nil {
		return nil, fmt.Errorf("failed to get txn %s: %w", txSign.String(), err)
	}

	acRes := types.QueryAccountsResponse{}
	if err := parseReturnValueFromLogs(p.xcallIdl.GetProgramID().String(), txnres.Meta.LogMessages, &acRes); err != nil {
		return nil, fmt.Errorf("failed to parse return value: %w", err)
	}

	return &acRes, nil
}

func (p *Provider) QueryTransactionReceipt(ctx context.Context, txSign string) (*relayertypes.Receipt, error) {
	txSignature, err := solana.SignatureFromBase58(txSign)
	if err != nil {
		return nil, err
	}

	txn, err := p.client.GetTransaction(ctx, txSignature, nil)
	if err != nil {
		return nil, err
	}

	return &relayertypes.Receipt{
		TxHash: txSign,
		Height: txn.Slot,
		Status: txn.Meta.Err == nil,
	}, nil
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayertypes.MessageKey) (bool, error) {
	receiptAc, err := p.pdaRegistry.ConnReceipt.GetAddress(key.Sn.FillBytes(make([]byte, 16)))
	if err != nil {
		return false, err
	}

	receipt := struct{}{}

	if err := p.client.GetAccountInfo(ctx, receiptAc, &receipt); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (p *Provider) decodeCsMessage(ctx context.Context, msg []byte) (*types.CsMessage, error) {
	if err := p.RestoreKeystore(ctx); err != nil {
		return nil, err
	}

	discriminator, err := p.xcallIdl.GetInstructionDiscriminator(types.MethodDecodeCsMessage)
	if err != nil {
		return nil, err
	}

	msgBorshBytes, err := borsh.Serialize(msg)
	if err != nil {
		return nil, err
	}

	instructionData := append(discriminator, msgBorshBytes...)

	instructions := []solana.Instruction{
		&solana.GenericInstruction{
			ProgID: p.xcallIdl.GetProgramID(),
			AccountValues: solana.AccountMetaSlice{
				&solana.AccountMeta{
					PublicKey:  solana.SystemProgramID,
					IsWritable: false,
					IsSigner:   false,
				},
			},
			DataBytes: instructionData,
		},
	}

	signers := []solana.PrivateKey{p.wallet.PrivateKey}

	tx, err := p.prepareTx(ctx, instructions, signers, solana.TransactionPayer(p.wallet.PublicKey()))
	if err != nil {
		return nil, fmt.Errorf("failed to prepare and simulate tx: %w", err)
	}

	txSign, err := p.client.SendTx(ctx, tx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send tx: %w", err)
	}

	if _, err := p.waitForTxConfirmation(3*time.Second, txSign); err != nil {
		return nil, fmt.Errorf("failed to confirm tx %s: %w", txSign.String(), err)
	}

	txnres, err := p.client.GetTransaction(ctx, txSign, &solrpc.GetTransactionOpts{Commitment: solrpc.CommitmentConfirmed})
	if err != nil {
		return nil, fmt.Errorf("failed to get txn %s: %w", txSign.String(), err)
	}

	for _, log := range txnres.Meta.LogMessages {
		xcallReturnPrefix := fmt.Sprintf("%s%s ", types.ProgramReturnPrefix, p.cfg.XcallProgram)
		if strings.HasPrefix(log, xcallReturnPrefix) {
			returnLog := strings.Replace(log, xcallReturnPrefix, "", 1)
			returnLogBytes, err := base64.StdEncoding.DecodeString(returnLog)
			if err != nil {
				return nil, err
			}

			csMsg := types.CsMessage{}
			if err := borsh.Deserialize(&csMsg, returnLogBytes); err != nil {
				return nil, err
			}

			return &csMsg, nil
		}
	}

	return nil, fmt.Errorf("failed to return value")
}

func parseReturnValueFromLogs(progID string, logs []string, dest interface{}) error {
	for _, log := range logs {
		xcallReturnPrefix := fmt.Sprintf("%s%s ", types.ProgramReturnPrefix, progID)
		if strings.HasPrefix(log, xcallReturnPrefix) {
			returnLog := strings.Replace(log, xcallReturnPrefix, "", 1)
			returnLogBytes, err := base64.StdEncoding.DecodeString(returnLog)
			if err != nil {
				return err
			}

			if err := borsh.Deserialize(dest, returnLogBytes); err != nil {
				return err
			}

			return nil
		}
	}
	return fmt.Errorf("logs is empty")
}
