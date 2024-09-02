package bitcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"path/filepath"

	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/multisig"
	"github.com/icon-project/goloop/common/codec"
	jsoniter "github.com/json-iterator/go"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

//var _ provider.ChainProvider = (*Provider)(nil)

var (
	BTCToken      = "0:0"
	MethodDeposit = "Deposit"
	MasterMode    = "master"
	SlaveMode     = "slave"
	BtcSlaveDB    = "btc_slave.db"
)

var chainIdToName = map[uint8]string{
	1: "0x1.icon",
	2: "0x1.btc",
	3: "0x2.icon",
	4: "0x2.btc",
	// Add more mappings as needed
}

type MessageDecoded struct {
	Action       string
	TokenAddress string
	To           string
	Amount       []byte
}

type slaveRequestParams struct {
	MsgSn		*big.Int		`json:"msgSn"`
}

type Provider struct {
	logger               *zap.Logger
	cfg                  *Config
	client               IClient
	LastSavedHeightFunc  func() uint64
	LastSerialNumFunc    func() *big.Int
	multisigAddrScript   []byte
	assetManagerAddrIcon string
	bearToken            string
	httpServer           chan struct{}
	db                   *leveldb.DB
}

type Config struct {
	provider.CommonConfig `json:",inline" yaml:",inline"`
	OpCode                int      `json:"op-code" yaml:"op-code"`
	UniSatURL             string   `json:"unisat-url" yaml:"unisat-url"`
	Type                  string   `json:"type" yaml:"type"`
	User                  string   `json:"rpc-user" yaml:"rpc-user"`
	Password              string   `json:"rpc-password" yaml:"rpc-password"`
	Protocals             []string `yaml:"protocals"`
	Mode                  string   `json:"mode" yaml:"mode"`
	SlaveServer1          string   `json:"slave-server-1" yaml:"slave-server-1"`
	SlaveServer2          string   `json:"slave-server-2" yaml:"slave-server-2"`
	Port                  string   `json:"port" yaml:"port"`
	ApiKey                string   `json:"api-key" yaml:"api-key"`
}

// NewProvider returns new Icon provider
func (c *Config) NewProvider(ctx context.Context, log *zap.Logger, homepath string, debug bool, chainName string) (provider.ChainProvider, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if err := c.sanitize(); err != nil {
		return nil, err
	}

	// Create the database file path
	dbPath := filepath.Join(homepath+"/data", BtcSlaveDB)

	// Open the database, creating it if it doesn't exist
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open or create database: %v", err)
	}

	client, err := newClient(ctx, c.RPCUrl, c.User, c.Password, true, false, log)
	if err != nil {
		db.Close() // Close the database if client creation fails
		return nil, fmt.Errorf("failed to create new client: %v", err)
	}

	c.ChainName = chainName
	c.HomeDir = homepath

	p := &Provider{
		logger:            log.With(zap.Stringp("nid", &c.NID), zap.Stringp("name", &c.ChainName)),
		cfg:               c,
		client:            client,
		LastSerialNumFunc: func() *big.Int { return big.NewInt(0) },
		httpServer:        make(chan struct{}),
		db:                db, // Add the database to the Provider
	}
	// Run an http server to help btc interact each others
	go func() {
		if c.Mode == MasterMode {
			startMaster(c)
		} else {
			startSlave(c)
		}
		close(p.httpServer)
	}()

	return p, nil
}

func (p *Provider) CallSlaves(slaveRequestData []byte) [][][]byte {
	resultChan := make(chan [][][]byte)
	go func() {
		responses := make(chan [][]byte, 2)
		var wg sync.WaitGroup
		wg.Add(2)

		go requestPartialSign(p.cfg.SlaveServer1, slaveRequestData, responses, &wg)
		go requestPartialSign(p.cfg.SlaveServer2, slaveRequestData, responses, &wg)

		go func() {
			wg.Wait()
			close(responses)
		}()

		var results [][][]byte
		for res := range responses {
			results = append(results, res)
		}
		resultChan <- results
	}()

	return <-resultChan
}

func (p *Provider) QueryLatestHeight(ctx context.Context) (uint64, error) {
	return p.client.GetLatestBlockHeight(ctx)
}

// todo: fill up the result
func (p *Provider) QueryTransactionReceipt(ctx context.Context, txHash string) (*relayTypes.Receipt, error) {
	res, err := p.client.GetTransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}
	return &relayTypes.Receipt{
		TxHash: res.Txid,
		// Height: uint64(res.TxResponse.Height),
		// Status: types.CodeTypeOK == res.TxResponse.Code,
	}, nil
}

func (p *Provider) NID() string {
	return p.cfg.NID
}

func (p *Provider) Name() string {
	return p.cfg.ChainName
}

func (p *Provider) Init(ctx context.Context, homePath string, kms kms.KMS) error {
	//if err := p.cfg.Contracts.Validate(); err != nil {
	//	return err
	//}
	//p.kms = kms
	return nil
}

// Wallet returns the wallet of the provider
func (p *Provider) Wallet() error {
	return nil
}

func (p *Provider) Type() string {
	return p.cfg.ChainName
}

func (p *Provider) Config() provider.Config {
	return nil
}

func (p *Provider) Listener(ctx context.Context, lastSavedHeight uint64, blockInfoChan chan *relayTypes.BlockInfo) error {
	// run http server to help btc interact each others
	latestHeight, err := p.QueryLatestHeight(ctx)
	if err != nil {
		p.logger.Error("failed to get latest block height", zap.Error(err))
		return err
	}

	startHeight, err := p.getStartHeight(latestHeight, lastSavedHeight)
	if err != nil {
		p.logger.Error("failed to determine start height", zap.Error(err))
		return err
	}

	pollHeightTicker := time.NewTicker(time.Second * 60) // do scan each 2 mins
	pollHeightTicker.Stop()

	p.logger.Info("Start from height", zap.Uint64("height", startHeight), zap.Uint64("finality block", p.FinalityBlock(ctx)))

	for {
		select {
		case <-pollHeightTicker.C:
			//pollHeightTicker.Stop()
			//startHeight = p.GetLastSavedHeight()
			latestHeight, err = p.QueryLatestHeight(ctx)
			if err != nil {
				p.logger.Error("failed to get latest block height", zap.Error(err))
				continue
			}
		default:
			if startHeight < latestHeight {
				p.logger.Debug("Query started", zap.Uint64("from-height", startHeight), zap.Uint64("to-height", latestHeight))
				startHeight = p.runBlockQuery(ctx, blockInfoChan, startHeight, latestHeight)
			}
		}
	}
}

func decodeWithdrawToMessage(input []byte) (*MessageDecoded, []byte, error) {
	withdrawInfoWrapper := CSMessage{}
	_, err := codec.RLP.UnmarshalFromBytes(input, &withdrawInfoWrapper)
	if err != nil {
		log.Fatal(err.Error())
	}

	// withdraw info data
	withdrawInfoWrapperV2 := CSMessageRequestV2{}
	_, err = codec.RLP.UnmarshalFromBytes(withdrawInfoWrapper.Payload, &withdrawInfoWrapperV2)
	if err != nil {
		log.Fatal(err.Error())
	}
	// withdraw info
	withdrawInfo := &MessageDecoded{}
	_, err = codec.RLP.UnmarshalFromBytes(withdrawInfoWrapperV2.Data, &withdrawInfo)
	if err != nil {
		log.Fatal(err.Error())
	}

	return withdrawInfo, withdrawInfoWrapperV2.Data, nil
}

func GetBitcoinUTXOs(server, address string, amountRequired uint64, timeout uint) ([]*multisig.UTXO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	// TODO: loop query until sastified amountRequired
	resp, err := GetBtcUtxo(ctx, UNISAT_DEFAULT_TESTNET, "60b7bf52654454f19d8553e1b6427fb9fd2c722ea8dc6822bdf1dd7615b4b35d", address, 0, 32)
	if err != nil {
		return nil , fmt.Errorf("failed to query bitcoin UTXOs from unisat: %v", err)
	}
	outputs := []*multisig.UTXO{}
	var totalAmount uint64
	for _, utxo := range resp.Data.Utxo {
		if totalAmount >= amountRequired {
			break
		}
		outputAmount := utxo.Satoshi.Uint64()
		outputs = append(outputs, &multisig.UTXO{
			IsRelayersMultisig: true,
			TxHash:        utxo.TxId,
			OutputIdx:     uint32(utxo.Vout),
			OutputAmount:  outputAmount,
		})
		totalAmount += outputAmount
	}

	return outputs, nil
}

func GetRuneUTXOs(server, address, runeId string, amountRequired uint64, timeout uint) ([]*multisig.UTXO, uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	// TODO: loop query until sastified amountRequired
	resp, err := GetRuneUtxo(ctx, UNISAT_DEFAULT_TESTNET, address, runeId)
	if err != nil {
		return nil , 0, fmt.Errorf("failed to query rune UTXOs from unisat: %v", err)
	}
	outputs := []*multisig.UTXO{}
	var totalAmount uint64
	for _, utxo := range resp.Data.Utxo {
		if totalAmount >= amountRequired {
			break
		}
		outputs = append(outputs, &multisig.UTXO{
			IsRelayersMultisig: true,
			TxHash:        utxo.TxId,
			OutputIdx:     uint32(utxo.Vout),
			OutputAmount:  utxo.Satoshi.Uint64(),
		})
		if len(utxo.Runes) != 1 {
			return nil , 0, fmt.Errorf("number of runes in the utxo is not 1, but: %v", err)
		}
		runeAmount, _ := strconv.ParseUint(utxo.Runes[0].Amount, 10, 64)
		totalAmount += runeAmount
	}

	return outputs, totalAmount - amountRequired, nil
}

func CreateBitcoinMultisigTx(
	messageData []byte,
	feePerOutput uint64,
	relayersMultisigWallet *multisig.MultisigWallet,
	chainParam *chaincfg.Params,
	server string,
) ([]*multisig.UTXO, *wire.MsgTx, string, *txscript.TxSigHashes, error) {
	// ----- BUILD OUTPUTS -----
	outputs := []*multisig.OutputTx{}
	var bitcoinAmountRequired, runeAmountRequired uint64
	var runeRequired multisig.Rune
	decodedData, opreturnData, err := decodeWithdrawToMessage(messageData)
	// add bridge message to output
	scripts, _ := multisig.CreateBridgeMessageScripts(opreturnData, 76)
	for _, script := range scripts {
		outputs = append(outputs, &multisig.OutputTx{
			OpReturnScript: script,
		})
	}
	// add withdraw output
	if err == nil {
		amount := new(big.Int).SetBytes(decodedData.Amount).Uint64()
		if decodedData.Action == "WithdrawTo" {
			if decodedData.TokenAddress == "0:0" {
				// transfer btc
				outputs = append(outputs, &multisig.OutputTx{
					ReceiverAddress: decodedData.To,
					Amount: amount,
				})

				bitcoinAmountRequired = amount
			} else {
				// transfer rune
				runeAddress := strings.Split(decodedData.TokenAddress, ":")
				blockNumber, _ := strconv.ParseUint(runeAddress[0], 10, 64)
				txIndex, _ := strconv.ParseUint(runeAddress[0], 10, 32)
				runeRequired = multisig.Rune{
					BlockNumber: blockNumber,
					TxIndex: uint32(txIndex),
				}
				runeScript, _ := multisig.CreateRuneTransferScript(
					runeRequired,
					new(big.Int).SetUint64(amount),
					uint64(len(outputs)+2),
				)
				outputs = append(outputs, &multisig.OutputTx{
					OpReturnScript: runeScript,
				})
				outputs = append(outputs, &multisig.OutputTx{
					ReceiverAddress: decodedData.To,
					Amount: 1000,
				})

				bitcoinAmountRequired = 1000
				runeAmountRequired = amount
			}
		}
	}

	// ----- BUILD INPUTS -----
	relayersMultisigAddress, err := multisig.AddressOnChain(chainParam, relayersMultisigWallet)
	if err != nil {
		return nil, nil, "", nil, err
	}
	relayersMultisigAddressStr := relayersMultisigAddress.String()
	// add tx fee the the required bitcoin amount
	inputs := []*multisig.UTXO{}
	var inputsSatNeeded uint64
	if runeAmountRequired != 0 {
		// query rune UTXOs from unisat
		runeInputs, runeChangeAmount, err := GetRuneUTXOs(server, relayersMultisigAddressStr, decodedData.TokenAddress, runeAmountRequired, 3)
		if err != nil {
			return nil, nil, "", nil, err
		}
		inputs = append(inputs, runeInputs...)
		// add rune change to outputs
		runeScript, _ := multisig.CreateRuneTransferScript(
			runeRequired,
			new(big.Int).SetUint64(runeChangeAmount),
			uint64(len(outputs)+2),
		)
		outputs = append(outputs, &multisig.OutputTx{
			OpReturnScript: runeScript,
		})
		outputs = append(outputs, &multisig.OutputTx{
			ReceiverAddress: relayersMultisigAddressStr,
			Amount: 1000,
		})
		inputsSatNeeded = 1000
	}
	txFee := uint64(len(outputs)) * feePerOutput
	inputsSatNeeded += bitcoinAmountRequired + txFee
	// todo: cover case rune UTXOs have big enough dust amount to cover inputsSatNeeded, can store rune and bitcoin in the same utxo
	// query bitcoin UTXOs from unisat
	bitcoinInputs, err := GetBitcoinUTXOs(server, relayersMultisigAddressStr, inputsSatNeeded, 3)
	if err != nil {
		return nil, nil, "", nil, err
	}
	inputs = append(inputs, bitcoinInputs...)
	// create raw tx
	msgTx, hexRawTx, txSigHashes, err := multisig.CreateMultisigTx(inputs, outputs, txFee, relayersMultisigWallet, nil, chainParam, relayersMultisigAddressStr, 0)

	return inputs, msgTx, hexRawTx, txSigHashes, err
}

func (p *Provider) Route(ctx context.Context, message *relayTypes.Message, callback relayTypes.TxResponseFunc) error {
	p.logger.Info("starting to route message", zap.Any("message", message))

	if (message.Src == "0x2.icon" && message.Dst == "0x2.btc") {
		// build unsigned tx
		// TODO: Build Multisig Wallet from config
		chainParam := &chaincfg.TestNet3Params
		_, relayersMultisigInfo := multisig.RandomMultisigInfo(3, 3, chainParam, []int{0, 1, 2}, 0, 1)
		relayersMultisigWallet, _ := multisig.BuildMultisigWallet(relayersMultisigInfo)

		inputs, msgTx, _, txSigHashes, _ := CreateBitcoinMultisigTx(message.Data, 10000, relayersMultisigWallet, chainParam, UNISAT_DEFAULT_TESTNET)
		tapSigParams := multisig.TapSigParams {
			TxSigHashes: txSigHashes,
			RelayersPKScript: relayersMultisigWallet.PKScript,
			RelayersTapLeaf: relayersMultisigWallet.TapLeaves[0],
			UserPKScript: []byte{},
			UserTapLeaf: txscript.TapLeaf{},
		}
		// TODO: Load priv key from config
		relayerPrivKey := "RELAYER_PRIV_KEY"
		// relayer sign tx
		relayerSigs, err := multisig.PartSignOnRawExternalTx(relayerPrivKey, msgTx, inputs, tapSigParams, chainParam, false)
		if err != nil {
			fmt.Println("err sign: ", err)
		}
		if p.cfg.Mode == SlaveMode {
			// Store the sigs in LevelDB
			key := []byte(fmt.Sprintf("bitcoin_message_%s", message.Sn.String()))
			value, _ := json.Marshal(relayerSigs)
			err := p.db.Put(key, value, nil)
			if err != nil {
				return fmt.Errorf("failed to store message in LevelDB: %v", err)
			}

			p.logger.Info("Message stored in LevelDB", zap.String("key", string(key)))
			return nil
		} else if p.cfg.Mode == MasterMode {
			totalSigs := [][][]byte{relayerSigs}
			// send unsigned raw tx and message sn to 2 slave relayers to get sign
			rsi := slaveRequestParams{
				MsgSn: message.Sn,
			}
			slaveRequestData, _ := json.Marshal(rsi)
			slaveSigs := p.CallSlaves(slaveRequestData)
			totalSigs = append(totalSigs, slaveSigs...)
			// combine sigs
			signedMsgTx, err := multisig.CombineMultisigSigs(msgTx, inputs, relayersMultisigWallet, 0, relayersMultisigWallet, 0, totalSigs)
			if err != nil {
				fmt.Println("err combine tx: ", err)
			}
			// broadcast tx
			fmt.Println("signedMsgTx:", signedMsgTx)
		}
	}

	// MasterMode logic
	// txid := "example_txid" // Replace with actual txid generation logic
	// signatures := p.CallSlaves(txid)

	// TODO: Implement logic to build bitcoin transaction using signatures

	// TODO: Broadcast transaction to bitcoin network

	// After successful broadcast, remove the message from LevelDB if it exists

	// TODO: Implement proper callback handling
	// callback(message.MessageKey(), txResponse, nil)

	return nil
}

// call the smart contract to send the message
func (p *Provider) call(ctx context.Context, message *relayTypes.Message) (string, error) {
	// rawMsg, err := p.getRawContractMessage(message)
	// if err != nil {
	// 	return nil, err
	// }

	// var contract string

	// switch message.EventType {
	// case events.EmitMessage, events.RevertMessage, events.SetAdmin, events.ClaimFee, events.SetFee:
	// 	//contract = p.cfg.Contracts[relayTypes.ConnectionContract]
	// case events.CallMessage, events.RollbackMessage:
	// 	//contract = p.cfg.Contracts[relayTypes.XcallContract]
	// default:
	// 	return nil, fmt.Errorf("unknown event type: %s ", message.EventType)
	// }

	// msg := &wasmTypes.MsgExecuteContract{
	// 	Sender:   p.Wallet().String(),
	// 	Contract: contract,
	// 	Msg:      rawMsg,
	// }

	// msgs := []sdkTypes.Msg{msg}

	// res, err := p.sendMessage(ctx, msgs...)
	// if err != nil {
	// 	if strings.Contains(err.Error(), errors.ErrWrongSequence.Error()) {
	// 		if mmErr := p.handleSequence(ctx); mmErr != nil {
	// 			return res, fmt.Errorf("failed to handle sequence mismatch error: %v || %v", mmErr, err)
	// 		}
	// 		return p.sendMessage(ctx, msgs...)
	// 	}
	// }
	return "", nil
}

func (p *Provider) sendMessage(ctx context.Context, msgs ...sdkTypes.Msg) (string, error) {

	// return p.prepareAndPushTxToMemPool(ctx, p.wallet.GetAccountNumber(), p.wallet.GetSequence(), msgs...)
	return "", nil
}

func (p *Provider) handleSequence(ctx context.Context) error {
	return nil
}

func (p *Provider) logTxFailed(err error, txHash string, code uint8) {
	p.logger.Error("transaction failed",
		zap.Error(err),
		zap.String("tx_hash", txHash),
		zap.Uint8("code", code),
	)
}

func (p *Provider) logTxSuccess(height uint64, txHash string) {
	p.logger.Info("successful transaction",
		zap.Uint64("block_height", height),
		//zap.String("chain_id", p.cfg.ChainID),
		zap.String("tx_hash", txHash),
	)
}

func (p *Provider) prepareAndPushTxToMemPool(ctx context.Context, acc, seq uint64, msgs ...sdkTypes.Msg) (*sdkTypes.TxResponse, error) {

	return nil, nil
}

func (p *Provider) waitForTxResult(ctx context.Context, mk *relayTypes.MessageKey, txHash string, callback relayTypes.TxResponseFunc) {
	//for txWaitRes := range p.subscribeTxResultStream(ctx, txHash, p.cfg.TxConfirmationInterval) {
	//	if txWaitRes.Error != nil && txWaitRes.Error != context.DeadlineExceeded {
	//		p.logTxFailed(txWaitRes.Error, txHash, uint8(txWaitRes.TxResult.Code))
	//		callback(mk, txWaitRes.TxResult, txWaitRes.Error)
	//		return
	//	}
	//	p.logTxSuccess(uint64(txWaitRes.TxResult.Height), txHash)
	//	callback(mk, txWaitRes.TxResult, nil)
	//}
}

func (p *Provider) pollTxResultStream(ctx context.Context, txHash string, maxWaitInterval time.Duration) <-chan *types.TxResult {
	txResChan := make(chan *types.TxResult)
	//startTime := time.Now()
	//go func(txChan chan *types.TxResult) {
	//	defer close(txChan)
	//	for range time.NewTicker(p.cfg.TxConfirmationInterval).C {
	//		res, err := p.client.GetTransactionReceipt(ctx, txHash)
	//		if err == nil {
	//			txChan <- &types.TxResult{
	//				TxResult: &relayTypes.TxResponse{
	//					Height:    res.TxResponse.Height,
	//					TxHash:    res.TxResponse.TxHash,
	//					Codespace: res.TxResponse.Codespace,
	//					Code:      relayTypes.ResponseCode(res.TxResponse.Code),
	//					Data:      res.TxResponse.Data,
	//				},
	//			}
	//			return
	//		} else if time.Since(startTime) > maxWaitInterval {
	//			txChan <- &types.TxResult{
	//				Error: err,
	//			}
	//			return
	//		}
	//	}
	//}(txResChan)
	return txResChan
}

func (p *Provider) subscribeTxResultStream(ctx context.Context, txHash string, maxWaitInterval time.Duration) <-chan *types.TxResult {
	txResChan := make(chan *types.TxResult)
	//go func(txRes chan *types.TxResult) {
	//	defer close(txRes)
	//
	//	newCtx, cancel := context.WithTimeout(ctx, maxWaitInterval)
	//	defer cancel()
	//
	//	query := fmt.Sprintf("tm.event = 'Tx' AND tx.hash = '%s'", txHash)
	//	resultEventChan, err := p.client.Subscribe(newCtx, "tx-result-waiter", query)
	//	if err != nil {
	//		txRes <- &types.TxResult{
	//			TxResult: &relayTypes.TxResponse{
	//				TxHash: txHash,
	//			},
	//			Error: err,
	//		}
	//		return
	//	}
	//	defer p.client.Unsubscribe(newCtx, "tx-result-waiter", query)
	//
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			return
	//		case e := <-resultEventChan:
	//			eventDataJSON, err := jsoniter.Marshal(e.Data)
	//			if err != nil {
	//				txRes <- &types.TxResult{
	//					TxResult: &relayTypes.TxResponse{
	//						TxHash: txHash,
	//					}, Error: err,
	//				}
	//				return
	//			}
	//
	//			txWaitRes := new(types.TxResultWaitResponse)
	//			if err := jsoniter.Unmarshal(eventDataJSON, txWaitRes); err != nil {
	//				txRes <- &types.TxResult{
	//					TxResult: &relayTypes.TxResponse{
	//						TxHash: txHash,
	//					}, Error: err,
	//				}
	//				return
	//			}
	//			if uint32(txWaitRes.Result.Code) != types.CodeTypeOK {
	//				txRes <- &types.TxResult{
	//					Error: fmt.Errorf(txWaitRes.Result.Log),
	//					TxResult: &relayTypes.TxResponse{
	//						Height:    txWaitRes.Height,
	//						TxHash:    txHash,
	//						Codespace: txWaitRes.Result.Codespace,
	//						Code:      relayTypes.ResponseCode(txWaitRes.Result.Code),
	//						Data:      string(txWaitRes.Result.Data),
	//					},
	//				}
	//				return
	//			}
	//
	//			txRes <- &types.TxResult{
	//				TxResult: &relayTypes.TxResponse{
	//					Height:    txWaitRes.Height,
	//					TxHash:    txHash,
	//					Codespace: txWaitRes.Result.Codespace,
	//					Code:      relayTypes.ResponseCode(txWaitRes.Result.Code),
	//					Data:      string(txWaitRes.Result.Data),
	//				},
	//			}
	//			return
	//		}
	//	}
	//}(txResChan)
	return txResChan
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayTypes.MessageKey) (bool, error) {
	//queryMsg := &types.QueryReceiptMsg{
	//	GetReceipt: &types.GetReceiptMsg{
	//		SrcNetwork: key.Src,
	//		ConnSn:     strconv.FormatUint(key.Sn, 10),
	//	},
	//}
	//rawQueryMsg, err := jsoniter.Marshal(queryMsg)
	//if err != nil {
	//	return false, err
	//}
	//
	//res, err := p.client.QuerySmartContract(ctx, p.cfg.Contracts[relayTypes.ConnectionContract], rawQueryMsg)
	//if err != nil {
	//	p.logger.Error("failed to check if message is received: ", zap.Error(err))
	//	return false, err
	//}
	//
	//receiptMsgRes := types.QueryReceiptMsgResponse{}
	//return receiptMsgRes.Status, jsoniter.Unmarshal(res.Data, &receiptMsgRes.Status)

	return false, nil
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*relayTypes.Coin, error) {
	//coin, err := p.client.GetBalance(ctx, addr, p.cfg.Denomination)
	//if err != nil {
	//	p.logger.Error("failed to query balance: ", zap.Error(err))
	//	return nil, err
	//}
	//return &relayTypes.Coin{
	//	Denom:  coin.Denom,
	//	Amount: coin.Amount.BigInt().Uint64(),
	//}, nil
	return nil, nil
}

func (p *Provider) ShouldReceiveMessage(ctx context.Context, message *relayTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) ShouldSendMessage(ctx context.Context, message *relayTypes.Message) (bool, error) {
	return true, nil
}

func (p *Provider) GenerateMessages(ctx context.Context, messageKey *relayTypes.MessageKeyWithMessageHeight) ([]*relayTypes.Message, error) {
	blocks, err := p.fetchBlockMessages(ctx, &HeightRange{messageKey.Height, messageKey.Height})
	if err != nil {
		return nil, err
	}
	var messages []*relayTypes.Message
	for _, block := range blocks {
		messages = append(messages, block.Messages...)
	}
	return messages, nil
}

func (p *Provider) FinalityBlock(ctx context.Context) uint64 {
	//return p.cfg.FinalityBlock
	return 0
}

func (p *Provider) RevertMessage(ctx context.Context, sn *big.Int) error {
	msg := &relayTypes.Message{
		Sn:        sn,
		EventType: events.RevertMessage,
	}
	_, err := p.call(ctx, msg)
	return err
}

// SetFee
func (p *Provider) SetFee(ctx context.Context, networkdID string, msgFee, resFee *big.Int) error {
	msg := &relayTypes.Message{
		Src:       networkdID,
		Sn:        msgFee,
		ReqID:     resFee,
		EventType: events.SetFee,
	}
	_, err := p.call(ctx, msg)
	return err
}

// ClaimFee
func (p *Provider) ClaimFee(ctx context.Context) error {
	msg := &relayTypes.Message{
		EventType: events.ClaimFee,
	}
	_, err := p.call(ctx, msg)
	return err
}

// GetFee returns the fee for the given networkID
// responseFee is used to determine if the fee should be returned
func (p *Provider) GetFee(ctx context.Context, networkID string, responseFee bool) (uint64, error) {
	//getFee := types.NewExecGetFee(networkID, responseFee)
	//data, err := jsoniter.Marshal(getFee)
	//if err != nil {
	//	return 0, err
	//}
	//return p.client.GetFee(ctx, p.cfg.Contracts[relayTypes.ConnectionContract], data)

	return 0, nil
}

func (p *Provider) SetAdmin(ctx context.Context, address string) error {
	msg := &relayTypes.Message{
		Src:       address,
		EventType: events.SetAdmin,
	}
	_, err := p.call(ctx, msg)
	return err
}

// ExecuteRollback
func (p *Provider) ExecuteRollback(ctx context.Context, sn *big.Int) error {
	return nil
}

// todo:
func (p *Provider) getStartHeight(latestHeight, lastSavedHeight uint64) (uint64, error) {
	startHeight := lastSavedHeight
	if p.cfg.StartHeight > 0 && p.cfg.StartHeight < latestHeight {
		return p.cfg.StartHeight, nil
	}

	if startHeight > latestHeight {
		return 0, fmt.Errorf("last saved height cannot be greater than latest height")
	}

	if startHeight != 0 && startHeight < latestHeight {
		return startHeight, nil
	}

	return latestHeight, nil
}

func (p *Provider) getHeightStream(done <-chan bool, fromHeight, toHeight uint64) <-chan *HeightRange {
	heightChan := make(chan *HeightRange)
	go func(fromHeight, toHeight uint64, heightChan chan *HeightRange) {
		defer close(heightChan)
		for fromHeight < toHeight {
			select {
			case <-done:
				return
			case heightChan <- &HeightRange{Start: fromHeight, End: fromHeight + 2}:
				fromHeight += 2
			}
		}
	}(fromHeight, toHeight, heightChan)
	return heightChan
}

func (p *Provider) getBlockInfoStream(ctx context.Context, done <-chan bool, heightStreamChan <-chan *HeightRange) <-chan interface{} {
	blockInfoStream := make(chan interface{})
	go func(blockInfoChan chan interface{}, heightChan <-chan *HeightRange) {
		defer close(blockInfoChan)
		for {
			select {
			case <-done:
				return
			case height, ok := <-heightChan:
				if ok {
					for {
						messages, err := p.fetchBlockMessages(ctx, height)
						if err != nil {
							p.logger.Error("failed to fetch block messages", zap.Error(err), zap.Any("height", height))
							time.Sleep(time.Second * 3)
						} else {
							for _, message := range messages {
								blockInfoChan <- message
							}
							break
						}
					}
				}
			}
		}
	}(blockInfoStream, heightStreamChan)
	return blockInfoStream
}

func (p *Provider) fetchBlockMessages(ctx context.Context, heightInfo *HeightRange) ([]*relayTypes.BlockInfo, error) {

	var (
		// todo: query from provide.config
		multisigAddress = p.cfg.Address
		preFixOP        = p.cfg.OpCode
	)

	multiSigScript, err := p.client.DecodeAddress(multisigAddress)
	if err != nil {
		return nil, err
	}

	searchParam := TxSearchParam{
		StartHeight:    heightInfo.Start,
		EndHeight:      heightInfo.End,
		BitcoinScript:  multiSigScript,
		OPReturnPrefix: preFixOP,
	}

	messages, err := p.client.TxSearch(context.Background(), searchParam)

	if err != nil {
		return nil, err
	}

	return p.getMessagesFromTxList(messages)
}

func (p *Provider) parseMessageFromTx(tx *TxSearchRes) (*relayTypes.Message, error) {
	// handle for bitcoin bridge
	// decode message from OP_RETURN
	p.logger.Info("parseMessageFromTx",
		zap.Uint64("height", tx.Height))

	bridgeMessage, err := multisig.ReadBridgeMessage(tx.Tx)
	if err != nil {
		return nil, err
	}
	messageInfo := bridgeMessage.Message

	if messageInfo.Action == MethodDeposit && messageInfo.To == p.assetManagerAddrIcon {
		// maybe get this function name from cf file
		// todo verify transfer amount match in calldata if it
		// call 3rd to check rune amount
		tokenId := messageInfo.TokenAddress
		amount := big.NewInt(0)
		amount.SetBytes(messageInfo.Amount)
		destContract := messageInfo.To

		fmt.Println(tokenId)
		fmt.Println(amount.String())
		fmt.Println(destContract)

		// call api to verify the data
		// https://docs.unisat.io/dev/unisat-developer-center/runes/get-utxo-runes-balance
		verified := false
		for i, out := range tx.Tx.TxOut {
			if bytes.Compare(out.PkScript, p.multisigAddrScript) != 0 {
				continue
			}

			if messageInfo.TokenAddress == BTCToken {
				if amount.Cmp(big.NewInt(out.Value)) == 0 {
					verified = true
					break
				}
			} else {
				// https://open-api.unisat.io/v1/indexer/runes/utxo
				runes, err := GetRuneTxIndex(p.cfg.UniSatURL, "GET", p.bearToken, tx.Tx.TxHash().String(), i)
				if err != nil {
					return nil, err
				}

				if len(runes.Data) == 0 {
					continue
				}

				for _, runeOut := range runes.Data {
					runeTokenBal, ok := big.NewInt(0).SetString(runeOut.Amount, 10)
					if !ok {
						return nil, fmt.Errorf("rune amount out invalid")
					}

					if amount.Cmp(runeTokenBal) == 0 && runeOut.RuneId == messageInfo.TokenAddress {
						verified = true
						break
					}
				}
			}
		}

		if !verified {
			return nil, fmt.Errorf("failed to verify transaction %v", tx.Tx.TxHash().String())
		}
	}

	// todo: verify bridge fee

	// parse message

	// todo: handle for rad fi

	// TODO: call xcallformat and then replace to data

	from := p.cfg.NID + "/" + p.cfg.Address
	decodeMessage, _ := codec.RLP.MarshalToBytes(messageInfo)
	data, _ := XcallFormat(decodeMessage, from, bridgeMessage.Address, uint(tx.Height), p.cfg.Protocals)

	return &relayTypes.Message{
		// TODO:
		Dst: "0x2.icon",
		// Dst: chainIdToName[bridgeMessage.ChainId],
		Src: p.NID(),
		// Sn:            p.LastSerialNumFunc(),
		Sn:            new(big.Int).SetUint64(tx.Height<<32 + tx.TxIndex),
		Data:          data,
		MessageHeight: tx.Height,
		EventType:     events.EmitMessage,
	}, nil
}

func (p *Provider) getMessagesFromTxList(resultTxList []*TxSearchRes) ([]*relayTypes.BlockInfo, error) {
	var messages []*relayTypes.BlockInfo
	for _, resultTx := range resultTxList {
		msg, err := p.parseMessageFromTx(resultTx)
		if err != nil {
			return nil, err
		}

		msg.MessageHeight = resultTx.Height
		p.logger.Info("Detected eventlog",
			zap.Uint64("height", msg.MessageHeight),
			zap.String("target_network", msg.Dst),
			zap.Uint64("sn", msg.Sn.Uint64()),
			zap.String("event_type", msg.EventType),
		)
		messages = append(messages, &relayTypes.BlockInfo{
			Height:   resultTx.Height,
			Messages: []*relayTypes.Message{msg},
		})
	}
	return messages, nil
}

func (p *Provider) getRawContractMessage(message *relayTypes.Message) (wasmTypes.RawContractMessage, error) {
	switch message.EventType {
	case events.EmitMessage:
		rcvMsg := types.NewExecRecvMsg(message)
		return jsoniter.Marshal(rcvMsg)
	case events.CallMessage:
		execMsg := types.NewExecExecMsg(message)
		return jsoniter.Marshal(execMsg)
	case events.RevertMessage:
		revertMsg := types.NewExecRevertMsg(message)
		return jsoniter.Marshal(revertMsg)
	case events.SetAdmin:
		setAdmin := types.NewExecSetAdmin(message.Dst)
		return jsoniter.Marshal(setAdmin)
	case events.ClaimFee:
		claimFee := types.NewExecClaimFee()
		return jsoniter.Marshal(claimFee)
	case events.SetFee:
		setFee := types.NewExecSetFee(message.Src, message.Sn, message.ReqID)
		return jsoniter.Marshal(setFee)
	case events.RollbackMessage:
		executeRollback := types.NewExecExecuteRollback(message.Sn)
		return jsoniter.Marshal(executeRollback)
	default:
		return nil, fmt.Errorf("unknown event type: %s ", message.EventType)
	}
}

func (p *Provider) getNumOfPipelines(diff int) int {
	if diff <= runtime.NumCPU() {
		return diff
	}
	return runtime.NumCPU() / 2
}

func (p *Provider) runBlockQuery(ctx context.Context, blockInfoChan chan *relayTypes.BlockInfo, fromHeight, toHeight uint64) uint64 {
	done := make(chan bool)
	defer close(done)

	heightStream := p.getHeightStream(done, fromHeight, toHeight)

	diff := int(toHeight-fromHeight) / 2

	numOfPipelines := p.getNumOfPipelines(diff)
	wg := &sync.WaitGroup{}
	for i := 0; i < numOfPipelines; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, heightStream <-chan *HeightRange) {
			defer wg.Done()
			for heightRange := range heightStream {
				blockInfo, err := p.fetchBlockMessages(ctx, heightRange)
				if err != nil {
					p.logger.Error("failed to fetch block messages", zap.Error(err))
					continue
				}
				var messages []*relayTypes.Message
				for _, block := range blockInfo {
					messages = append(messages, block.Messages...)
				}
				blockInfoChan <- &relayTypes.BlockInfo{
					Height:   heightRange.End,
					Messages: messages,
				}
			}
		}(wg, heightStream)
	}
	wg.Wait()
	return toHeight + 1
}

// SubscribeMessageEvents subscribes to the message events
// Expermental: Allows to subscribe to the message events realtime without fully syncing the chain
func (p *Provider) SubscribeMessageEvents(ctx context.Context, blockInfoChan chan *relayTypes.BlockInfo, opts *types.SubscribeOpts, resetFunc func()) error {
	return nil
}

// SetLastSavedHeightFunc sets the function to save the last saved height
func (p *Provider) SetLastSavedHeightFunc(f func() uint64) {
	p.LastSavedHeightFunc = f
}

// GetLastSavedHeight returns the last saved height
func (p *Provider) GetLastSavedHeight() uint64 {
	return p.LastSavedHeightFunc()
}

func (p *Provider) SetSerialNumberFunc(f func() *big.Int) {
	p.LastSerialNumFunc = f
}

func (p *Provider) GetSerialNumber() *big.Int {
	return p.LastSerialNumFunc()
}

func (p *Config) sanitize() error {
	// TODO:
	return nil
}

func (c *Config) Validate() error {
	if c.RPCUrl == "" {
		return fmt.Errorf("bitcoin provider rpc endpoint is empty")
	}
	return nil
}
