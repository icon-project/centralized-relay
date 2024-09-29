package bitcoin

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/bxelab/runestone"
	"lukechampine.com/uint128"

	"path/filepath"

	"github.com/icon-project/centralized-relay/relayer/chains/wasm/types"
	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/kms"
	"github.com/icon-project/centralized-relay/relayer/provider"
	relayTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/multisig"
	"github.com/icon-project/goloop/common/codec"

	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

//var _ provider.ChainProvider = (*Provider)(nil)

var (
	BTCToken                = "0:1"
	MethodDeposit           = "Deposit"
	MethodWithdrawTo        = "WithdrawTo"
	MasterMode              = "master"
	SlaveMode               = "slave"
	BtcDB                   = "btc.db"
	MinSatsRequired  uint64 = 1000
	WitnessSize             = 385
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

type CSMessageResult struct {
	Sn      *big.Int
	Code    uint8
	Message []byte
}

type slaveRequestParams struct {
	MsgSn string `json:"msgSn"`
}
type StoredMessageData struct {
	OriginalMessage  *relayTypes.Message
	TxHash           string
	OutputIndex      uint32
	Amount           uint64
	RecipientAddress string
	SenderAddress    string
	RuneId           string
	RuneAmount       uint64
	ActionMethod     string
	TokenAddress     string
}

type Provider struct {
	logger              *zap.Logger
	cfg                 *Config
	client              IClient
	LastSavedHeightFunc func() uint64
	LastSerialNumFunc   func() *big.Int
	multisigAddrScript  []byte
	bearToken           string
	httpServer          chan struct{}
	db                  *leveldb.DB
	chainParam          *chaincfg.Params
}

type Config struct {
	provider.CommonConfig `json:",inline" yaml:",inline"`
	OpCode                int      `json:"op-code" yaml:"op-code"`
	UniSatURL             string   `json:"unisat-url" yaml:"unisat-url"`
	UniSatKey             string   `json:"unisat-key" yaml:"unisat-key"`
	MempoolURL            string   `json:"mempool-url" yaml:"mempool-url"`
	Type                  string   `json:"type" yaml:"type"`
	User                  string   `json:"rpc-user" yaml:"rpc-user"`
	Password              string   `json:"rpc-password" yaml:"rpc-password"`
	Mode                  string   `json:"mode" yaml:"mode"`
	SlaveServer1          string   `json:"slave-server-1" yaml:"slave-server-1"`
	SlaveServer2          string   `json:"slave-server-2" yaml:"slave-server-2"`
	Port                  string   `json:"port" yaml:"port"`
	ApiKey                string   `json:"api-key" yaml:"api-key"`
	MasterPubKey          string   `json:"masterPubKey" yaml:"masterPubKey"`
	Slave1PubKey          string   `json:"slave1PubKey" yaml:"slave1PubKey"`
	Slave2PubKey          string   `json:"slave2PubKey" yaml:"slave2PubKey"`
	RelayerPrivKey        string   `json:"relayerPrivKey" yaml:"relayerPrivKey"`
	RecoveryLockTime      int      `json:"recoveryLockTime" yaml:"recoveryLockTime"`
	Connections           []string `json:"connections" yaml:"connections"`
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
	dbPath := filepath.Join(homepath+"/data", BtcDB)

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
	chainParam := &chaincfg.TestNet3Params
	if c.NID == "0x1.btc" {
		chainParam = &chaincfg.MainNetParams
	}
	c.HomeDir = homepath
	c.HomeDir = homepath

	msPubkey, err := client.DecodeAddress(c.Address)
	if err != nil {
		return nil, err
	}

	p := &Provider{
		logger:             log.With(zap.Stringp("nid", &c.NID), zap.Stringp("name", &c.ChainName)),
		cfg:                c,
		client:             client,
		LastSerialNumFunc:  func() *big.Int { return big.NewInt(0) },
		httpServer:         make(chan struct{}),
		db:                 db, // Add the database to the Provider
		chainParam:         chainParam,
		multisigAddrScript: msPubkey,
	}
	// Run an http server to help btc interact each others
	go func() {
		if c.Mode == MasterMode {
			startMaster(c)
		} else {
			startSlave(c, p)
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

		go requestPartialSign(p.cfg.ApiKey, p.cfg.SlaveServer1, slaveRequestData, responses, &wg)
		go requestPartialSign(p.cfg.ApiKey, p.cfg.SlaveServer2, slaveRequestData, responses, &wg)

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
func (p *Provider) Wallet() (*multisig.MultisigWallet, error) {
	return p.buildMultisigWallet()
}

func (p *Provider) Type() string {
	return p.cfg.ChainName
}

func (p *Provider) Config() provider.Config {
	return p.cfg
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

func (p *Provider) GetBitcoinUTXOs(server, address string, amountRequired uint64, timeout uint) ([]*multisig.UTXO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	// TODO: loop query until sastified amountRequired
	resp, err := GetBtcUtxo(ctx, server, p.cfg.UniSatKey, address, 0, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to query bitcoin UTXOs from unisat: %v", err)
	}
	outputs := []*multisig.UTXO{}
	var totalAmount uint64

	utxos := resp.Data.Utxo
	// sort utxos by amount in descending order
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Satoshi.Cmp(utxos[j].Satoshi) == 1
	})

	for _, utxo := range utxos {
		if totalAmount >= amountRequired {
			break
		}
		outputAmount := utxo.Satoshi.Uint64()
		outputs = append(outputs, &multisig.UTXO{
			IsRelayersMultisig: true,
			TxHash:             utxo.TxId,
			OutputIdx:          uint32(utxo.Vout),
			OutputAmount:       outputAmount,
		})
		totalAmount += outputAmount
	}

	return outputs, nil
}

func GetRuneUTXOs(server, address, runeId string, amountRequired uint64, timeout uint) ([]*multisig.UTXO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	// TODO: loop query until sastified amountRequired
	resp, err := GetRuneUtxo(ctx, server, address, runeId)
	if err != nil {
		return nil, fmt.Errorf("failed to query rune UTXOs from unisat: %v", err)
	}

	utxos := resp.Data.Utxo
	// sort utxos by amount in descending order
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Satoshi.Cmp(utxos[j].Satoshi) == 1
	})

	inputs := []*multisig.UTXO{}
	var totalAmount uint64
	for _, utxo := range utxos {
		if totalAmount >= amountRequired {
			break
		}
		inputs = append(inputs, &multisig.UTXO{
			IsRelayersMultisig: true,
			TxHash:             utxo.TxId,
			OutputIdx:          uint32(utxo.Vout),
			OutputAmount:       utxo.Satoshi.Uint64(),
		})
		if len(utxo.Runes) != 1 {
			return nil, fmt.Errorf("number of runes in the utxo is not 1, but: %v", err)
		}
		runeAmount, _ := strconv.ParseUint(utxo.Runes[0].Amount, 10, 64)
		totalAmount += runeAmount
	}

	return inputs, nil
}

func (p *Provider) CreateBitcoinMultisigTx(
	outputData []*multisig.OutputTx,

	feeRate uint64,
	decodedData *MessageDecoded,
	msWallet *multisig.MultisigWallet,
) ([]*multisig.UTXO, *wire.MsgTx, string, *txscript.TxSigHashes, error) {
	// ----- BUILD OUTPUTS -----
	outputs := []*multisig.OutputTx{}
	outputs = append(outputs, outputData...)
	var bitcoinAmountRequired uint64
	var runeAmountRequired uint64

	rlMsAddress, err := multisig.AddressOnChain(p.chainParam, msWallet)
	if err != nil {
		return nil, nil, "", nil, err
	}
	msAddressStr := rlMsAddress.String()

	// add withdraw output
	amount, _ := strconv.ParseUint(string(decodedData.Amount), 10, 64)
	if decodedData.Action == MethodWithdrawTo || decodedData.Action == MethodDeposit {
		if decodedData.TokenAddress == BTCToken {
			// transfer btc
			outputs = append(outputs, &multisig.OutputTx{
				ReceiverAddress: decodedData.To,
				Amount:          amount,
			})

			bitcoinAmountRequired = amount
		} else {
			// transfer rune
			runeRequired, _ := runestone.RuneIdFromString(decodedData.TokenAddress)
			changeOutputId := uint32(len(outputs) + 2)
			runeOutput := &runestone.Runestone{
				Edicts: []runestone.Edict{
					{
						ID:     *runeRequired,
						Amount: uint128.FromBytes(decodedData.Amount),
						Output: uint32(len(outputs) + 1),
					},
				},
				Pointer: &changeOutputId,
			}
			runeScript, _ := runeOutput.Encipher()
			// add runestone OP_RETURN
			outputs = append(outputs, &multisig.OutputTx{
				OpReturnScript: runeScript,
			})
			// add receiver output
			outputs = append(outputs, &multisig.OutputTx{
				ReceiverAddress: decodedData.To,
				Amount:          MinSatsRequired,
			})
			// add change output
			outputs = append(outputs, &multisig.OutputTx{
				ReceiverAddress: msAddressStr,
				Amount:          MinSatsRequired,
			})
			runeAmountRequired = amount
			bitcoinAmountRequired = MinSatsRequired * 2
		}
	}

	// ----- BUILD INPUTS -----
	inputs, estFee, err := p.computeTx(feeRate, bitcoinAmountRequired, runeAmountRequired, decodedData.TokenAddress, msAddressStr, outputs, msWallet)
	if err != nil {
		return nil, nil, "", nil, err
	}
	// create raw tx
	msgTx, hexRawTx, txSigHashes, err := multisig.CreateMultisigTx(inputs, outputs, estFee, msWallet, nil, p.chainParam, msAddressStr, 0)

	return inputs, msgTx, hexRawTx, txSigHashes, err
}

// calculateTxSize calculates the size of a transaction given the inputs, outputs, estimated fee, change address, chain parameters, and multisig wallet.
// It returns the size of the transaction in bytes and an error if any occurs during the process.
func (p *Provider) calculateTxSize(inputs []*multisig.UTXO, outputs []*multisig.OutputTx, estFee uint64, changeAddress string, msWallet *multisig.MultisigWallet) (int, error) {
	msgTx, _, _, err := multisig.CreateMultisigTx(inputs, outputs, estFee, msWallet, msWallet, p.chainParam, changeAddress, 0)
	if err != nil {
		return 0, err
	}
	var rawTxBytes bytes.Buffer
	err = msgTx.Serialize(&rawTxBytes)
	if err != nil {
		return 0, err
	}
	txSize := len(rawTxBytes.Bytes()) + WitnessSize
	return txSize, nil
}

func (p *Provider) computeTx(feeRate, satsToSend, runeToSend uint64, runeId, changeAddress string, outputs []*multisig.OutputTx, msWallet *multisig.MultisigWallet) ([]*multisig.UTXO, uint64, error) {

	outputsCopy := make([]*multisig.OutputTx, len(outputs))
	copy(outputsCopy, outputs)

	inputs, err := p.selectUnspentUTXOs(satsToSend, runeToSend, runeId, outputsCopy, changeAddress)
	sumSelectedOutputs := multisig.SumInputsSat(inputs)
	if err != nil {
		return nil, 0, err
	}

	txSize, err := p.calculateTxSize(inputs, outputsCopy, 0, changeAddress, msWallet)
	if err != nil {
		return nil, 0, err
	}

	estFee := uint64(txSize) * feeRate
	count := 0
	loopEntered := false
	var iterationOutputs []*multisig.OutputTx

	for sumSelectedOutputs < satsToSend+estFee {
		loopEntered = true
		// Create a fresh copy of outputs for each iteration
		iterationOutputs := make([]*multisig.OutputTx, len(outputs))
		copy(iterationOutputs, outputs)

		newSatsToSend := satsToSend + estFee
		var err error
		selectedUnspentOutputs, err := p.selectUnspentUTXOs(newSatsToSend, runeToSend, runeId, iterationOutputs, changeAddress)
		if err != nil {
			return nil, 0, err
		}

		sumSelectedOutputs = multisig.SumInputsSat(selectedUnspentOutputs)

		txSize, err := p.calculateTxSize(selectedUnspentOutputs, iterationOutputs, estFee, changeAddress, msWallet)
		if err != nil {
			return nil, 0, err
		}

		estFee = feeRate * uint64(txSize)

		count += 1
		if count > 500 {
			return nil, 0, fmt.Errorf("Not enough sats for fee")
		}
	}
	// Need to do that cause avoid the same outputs being used for multiple transactions
	outputs = outputsCopy
	if loopEntered {
		outputs = iterationOutputs
	}

	return inputs, estFee, nil
}

func (p *Provider) selectUnspentUTXOs(satToSend uint64, runeToSend uint64, runeId string, outputs []*multisig.OutputTx, changeAddress string) ([]*multisig.UTXO, error) {
	// add tx fee the the required bitcoin amount
	inputs := []*multisig.UTXO{}
	if runeToSend != 0 {
		// query rune UTXOs from unisat
		runeInputs, err := GetRuneUTXOs(p.cfg.UniSatURL, changeAddress, runeId, runeToSend, 3)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, runeInputs...)
	}

	// TODO: cover case rune UTXOs have big enough dust amount to cover inputsSatNeeded, can store rune and bitcoin in the same utxo
	// query bitcoin UTXOs from unisat
	bitcoinInputs, err := p.GetBitcoinUTXOs(p.cfg.UniSatURL, changeAddress, satToSend, 3)
	if err != nil {
		return nil, err
	}
	inputs = append(inputs, bitcoinInputs...)

	return inputs, nil
}

// add tx fee the the required bitcoin amount

func (p *Provider) buildMultisigWallet() (*multisig.MultisigWallet, error) {
	masterPubKey, _ := hex.DecodeString(p.cfg.MasterPubKey)
	slave1PubKey, _ := hex.DecodeString(p.cfg.Slave1PubKey)
	slave2PubKey, _ := hex.DecodeString(p.cfg.Slave2PubKey)
	relayersMultisigInfo := multisig.MultisigInfo{
		PubKeys:            [][]byte{masterPubKey, slave1PubKey, slave2PubKey},
		EcPubKeys:          nil,
		NumberRequiredSigs: 3,
		RecoveryPubKey:     masterPubKey,
		RecoveryLockTime:   uint64(p.cfg.RecoveryLockTime),
	}
	msWallet, err := multisig.BuildMultisigWallet(&relayersMultisigInfo)
	if err != nil {
		p.logger.Error("failed to build multisig wallet: %v", zap.Error(err))
		return nil, err
	}

	return msWallet, nil
}

func (p *Provider) partSignTx(msgTx *wire.MsgTx, inputs []*multisig.UTXO, msWallet *multisig.MultisigWallet, txSigHashes *txscript.TxSigHashes) ([]*multisig.UTXO, *multisig.MultisigWallet, *wire.MsgTx, [][]byte, error) {
	tapSigParams := multisig.TapSigParams{
		TxSigHashes:      txSigHashes,
		RelayersPKScript: msWallet.PKScript,
		RelayersTapLeaf:  msWallet.TapLeaves[0],
		UserPKScript:     []byte{},
		UserTapLeaf:      txscript.TapLeaf{},
	}
	// relayer sign tx
	relayerSigs, err := multisig.PartSignOnRawExternalTx(p.cfg.RelayerPrivKey, msgTx, inputs, tapSigParams, p.chainParam, false)

	return inputs, msWallet, msgTx, relayerSigs, err
}

func (p *Provider) HandleBitcoinMessageTx(message *relayTypes.Message) ([]*multisig.UTXO, *multisig.MultisigWallet, *wire.MsgTx, [][]byte, bool, *big.Int, error) {
	msWallet, err := p.buildMultisigWallet()
	if err != nil {
		return nil, nil, nil, nil, false, nil, err
	}
	feeRate, err := p.client.GetFeeFromMempool(p.cfg.MempoolURL + "/fees/recommended")
	if err != nil {
		p.logger.Error("failed to get recommended fee from mempool", zap.Error(err))
		feeRate = 0
	}
	inputs, msgTx, _, txSigHashes, isRollbackMessage, rollbackSn, err := p.buildTxMessage(message, feeRate, msWallet)
	if err != nil {
		p.logger.Error("failed to build tx message: %v", zap.Error(err))
		return nil, nil, nil, nil, false, nil, err
	}
	inputs, msWallet, msgTx, relayerSigs, err := p.partSignTx(msgTx, inputs, msWallet, txSigHashes)
	return inputs, msWallet, msgTx, relayerSigs, isRollbackMessage, rollbackSn, err
}

func (p *Provider) Route(ctx context.Context, message *relayTypes.Message, callback relayTypes.TxResponseFunc) error {
	p.logger.Info("starting to route message", zap.Any("message", message))

	messageDstDetail := strings.Split(message.Dst, ".")
	if strings.Split(message.Src, ".")[1] == "icon" && messageDstDetail[1] == "btc" {

		if p.cfg.Mode == SlaveMode {
			// store the message in LevelDB
			key := []byte(message.Sn.String())
			value, _ := json.Marshal(message)
			err := p.db.Put(key, value, nil)
			if err != nil {
				p.logger.Error("failed to store message in LevelDB: %v", zap.Error(err))
				return err
			}
			p.logger.Info("Message stored in LevelDB", zap.String("key", string(key)))
			return nil
		} else if p.cfg.Mode == MasterMode {

			inputs, msWallet, msgTx, relayerSigs, isRollbackMessage, rollbackSn, err := p.HandleBitcoinMessageTx(message)
			if err != nil {
				p.logger.Error("failed to handle bitcoin message tx: %v", zap.Error(err))
				return err
			}
			totalSigs := [][][]byte{relayerSigs}
			// send unsigned raw tx and message sn to 2 slave relayers to get sign
			rsi := slaveRequestParams{
				MsgSn: message.Sn.String(),
			}
			if isRollbackMessage {
				rsi.MsgSn = "RB" + rollbackSn.String()
			}

			slaveRequestData, _ := json.Marshal(rsi)
			slaveSigs := p.CallSlaves(slaveRequestData)
			p.logger.Info("Slave signatures", zap.Any("slave sigs", slaveSigs))
			totalSigs = append(totalSigs, slaveSigs...)
			// combine sigs
			signedMsgTx, err := multisig.CombineMultisigSigs(msgTx, inputs, msWallet, 0, msWallet, 0, totalSigs)

			if err != nil {
				p.logger.Error("err combine tx: ", zap.Error(err))
			}
			p.logger.Info("signedMsgTx", zap.Any("transaction", signedMsgTx))
			var buf bytes.Buffer
			err = signedMsgTx.Serialize(&buf)

			if err != nil {
				log.Fatal(err)
			}

			txSize := len(buf.Bytes())
			p.logger.Info("--------------------txSize--------------------", zap.Int("size", txSize))
			signedMsgTxHex := hex.EncodeToString(buf.Bytes())
			p.logger.Info("signedMsgTxHex", zap.String("transaction_hex", signedMsgTxHex))

			// Broadcast transaction to bitcoin network
			txHash, err := p.client.SendRawTransaction(ctx, p.cfg.MempoolURL, []json.RawMessage{json.RawMessage(`"` + signedMsgTxHex + `"`)})

			if err != nil {
				p.logger.Error("failed to send raw transaction", zap.Error(err))
				return err
			}

			p.logger.Info("txHash", zap.String("transaction_hash", txHash))
			// TODO: After successful broadcast, request slaves to remove the message from LevelDB if it exists
		}

	}
	return nil
}

// TODO: Implement proper callback handling
// callback(message.MessageKey(), txResponse, nil)

func (p *Provider) decodeMessage(message *relayTypes.Message) (CSMessageResult, error) {

	wrapperInfo := CSMessage{}
	_, err := codec.RLP.UnmarshalFromBytes(message.Data, &wrapperInfo)
	if err != nil {
		p.logger.Error("failed to unmarshal message: %v", zap.Error(err))
		return CSMessageResult{}, err
	}

	messageDecoded := CSMessageResult{}
	_, err = codec.RLP.UnmarshalFromBytes(wrapperInfo.Payload, &messageDecoded)
	if err != nil {
		p.logger.Error("failed to unmarshal message: %v", zap.Error(err))
		return CSMessageResult{}, err
	}

	return messageDecoded, nil
}

func (p *Provider) buildTxMessage(message *relayTypes.Message, feeRate uint64, msWallet *multisig.MultisigWallet) ([]*multisig.UTXO, *wire.MsgTx, string, *txscript.TxSigHashes, bool, *big.Int, error) {
	outputs := []*multisig.OutputTx{}
	decodedData := &MessageDecoded{}
	isRollbackMessage := false
	rollbackSn := new(big.Int)
	switch message.EventType {
	case events.EmitMessage:
		messageDecoded, err := p.decodeMessage(message)
		if err != nil {
			p.logger.Error("failed to decode message: %v", zap.Error(err))
		}
		// 0 is need to rollback
		if messageDecoded.Code == 0 {
			isRollbackMessage = true
			rollbackSn = new(big.Int).SetBytes(messageDecoded.Sn.Bytes())
			// Process RollbackMessage
			data, err := p.db.Get([]byte("RB"+messageDecoded.Sn.String()), nil)
			if err != nil {
				return nil, nil, "", nil, isRollbackMessage, nil, fmt.Errorf("failed to retrieve stored data: %v", err)
			}
			var storedData StoredMessageData
			err = json.Unmarshal(data, &storedData)
			if err != nil {
				return nil, nil, "", nil, isRollbackMessage, nil, fmt.Errorf("failed to unmarshal stored data: %v", err)
			}
			decodedData = &MessageDecoded{
				Action:       storedData.ActionMethod,
				To:           storedData.SenderAddress,
				TokenAddress: storedData.TokenAddress,
				Amount:       []byte(fmt.Sprintf("%d", storedData.Amount)),
			}
			if storedData.RuneId != "" {
				decodedData.Amount = []byte(fmt.Sprintf("%d", storedData.RuneAmount))
			}
		} else {
			// Perform WithdrawData
			data, opreturnData, err := decodeWithdrawToMessage(message.Data)
			decodedData = data
			if err != nil {
				p.logger.Error("failed to decode message: %v", zap.Error(err))
				return nil, nil, "", nil, isRollbackMessage, nil, err
			}
			scripts, _ := multisig.CreateBridgeMessageScripts(opreturnData, 76)
			for _, script := range scripts {
				outputs = append(outputs, &multisig.OutputTx{
					OpReturnScript: script,
				})
			}
		}
	default:
		return nil, nil, "", nil, isRollbackMessage, nil, fmt.Errorf("unknown event type: %s", message.EventType)
	}

	inputs, msgTx, hexRawTx, txSigHashes, err := p.CreateBitcoinMultisigTx(outputs, feeRate, decodedData, msWallet)
	return inputs, msgTx, hexRawTx, txSigHashes, isRollbackMessage, rollbackSn, err
}

// call the smart contract to send the message
func (p *Provider) call(ctx context.Context, message *relayTypes.Message) (string, error) {

	return "", nil
}

func (p *Provider) sendTx(ctx context.Context, signedMsg *wire.MsgTx) (string, error) {

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
		zap.String("chain_id", p.cfg.NID),
		zap.String("tx_hash", txHash),
	)
}

func (p *Provider) waitForTxResult(ctx context.Context, mk *relayTypes.MessageKey, txHash string, callback relayTypes.TxResponseFunc) {

}

func (p *Provider) pollTxResultStream(ctx context.Context, txHash string, maxWaitInterval time.Duration) <-chan *types.TxResult {
	txResChan := make(chan *types.TxResult)

	return txResChan
}

func (p *Provider) subscribeTxResultStream(ctx context.Context, txHash string, maxWaitInterval time.Duration) <-chan *types.TxResult {
	txResChan := make(chan *types.TxResult)

	return txResChan
}

func (p *Provider) MessageReceived(ctx context.Context, key *relayTypes.MessageKey) (bool, error) {

	return false, nil
}

func (p *Provider) QueryBalance(ctx context.Context, addr string) (*relayTypes.Coin, error) {

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
func (p *Provider) extractOutputReceiver(tx *wire.MsgTx) []string {
	receiverAddresses := []string{}
	for _, out := range tx.TxOut {
		receiverAddresses = append(receiverAddresses, p.getAddressesFromTx(out, p.chainParam)...)
	}
	return receiverAddresses
}

func (p *Provider) parseMessageFromTx(tx *TxSearchRes) (*relayTypes.Message, error) {
	receiverAddresses := []string{}
	runeId := ""
	runeAmount := big.NewInt(0)
	actionMethod := ""
	// handle for bitcoin bridge
	// decode message from OP_RETURN
	p.logger.Info("parseMessageFromTx",
		zap.Uint64("height", tx.Height))

	bridgeMessage, err := multisig.ReadBridgeMessage(tx.Tx)
	if err != nil {
		return nil, err
	}
	messageInfo := bridgeMessage.Message

	isValidConnector := false
	for _, connector := range bridgeMessage.Connectors {
		if slices.Contains(p.cfg.Connections, connector) {
			isValidConnector = true
			break
		}
	}

	if messageInfo.Action == MethodDeposit && isValidConnector {
		actionMethod = MethodDeposit
		// maybe get this function name from cf file
		// todo verify transfer amount match in calldata if it
		// call 3rd to check rune amount
		tokenId := messageInfo.TokenAddress
		amount := big.NewInt(0)
		amount.SetBytes(messageInfo.Amount)
		destContract := messageInfo.To

		p.logger.Info("tokenId", zap.String("tokenId", tokenId))
		p.logger.Info("amount", zap.String("amount", amount.String()))
		p.logger.Info("destContract", zap.String("destContract", destContract))

		// call api to verify the data
		// https://docs.unisat.io/dev/unisat-developer-center/runes/get-utxo-runes-balance
		verified := false
		receiverAddresses = p.extractOutputReceiver(tx.Tx)
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
						runeId = runeOut.RuneId
						runeAmount = runeTokenBal
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
	sn := new(big.Int).SetUint64(tx.Height<<32 + tx.TxIndex)

	from := p.cfg.NID + "/" + p.cfg.Address
	decodeMessage, _ := codec.RLP.MarshalToBytes(messageInfo)
	data, _ := XcallFormat(decodeMessage, from, bridgeMessage.Receiver, uint(sn.Uint64()), bridgeMessage.Connectors, uint8(CALL_MESSAGE_ROLLBACK_TYPE))

	relayMessage := &relayTypes.Message{
		Dst: "0x2.icon",
		// TODO:
		// Dst:           chainIdToName[bridgeMessage.ChainId],
		Src:           p.NID(),
		Sn:            sn,
		Data:          data,
		MessageHeight: tx.Height,
		EventType:     events.EmitMessage,
	}

	var senderAddress string
	// Find sender address to store in db
	for _, address := range receiverAddresses {
		if address != p.cfg.Address {
			senderAddress = address
			break
		}
	}

	// When storing the message
	storedData := StoredMessageData{
		OriginalMessage:  relayMessage,
		TxHash:           tx.Tx.TxHash().String(),
		OutputIndex:      uint32(tx.TxIndex),
		Amount:           big.NewInt(0).SetBytes(messageInfo.Amount).Uint64(),
		RecipientAddress: p.cfg.Address,
		SenderAddress:    senderAddress,
		RuneId:           runeId,
		RuneAmount:       runeAmount.Uint64(),
		ActionMethod:     actionMethod,
		TokenAddress:     messageInfo.TokenAddress,
	}

	data, err = json.Marshal(storedData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stored data: %v", err)
	}

	err = p.db.Put([]byte("RB"+sn.String()), data, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to store message data: %v", err)
	}
	return relayMessage, nil
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

func (p *Provider) getAddressesFromTx(txOut *wire.TxOut, chainParams *chaincfg.Params) []string {
	receiverAddresses := []string{}

	scriptClass, addresses, _, err := txscript.ExtractPkScriptAddrs(txOut.PkScript, chainParams)
	if err != nil {
		fmt.Printf("  Script: Unable to parse (possibly OP_RETURN)\n")
	} else {
		fmt.Printf("  Script Class: %s\n", scriptClass)
		if len(addresses) > 0 {
			fmt.Printf("  Receiver Address: %s\n", addresses[0].String())
			receiverAddresses = append(receiverAddresses, addresses[0].String())
		}
	}

	return receiverAddresses
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
