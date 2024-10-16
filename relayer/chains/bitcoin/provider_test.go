package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/icon-project/centralized-relay/relayer/chains/icon"
	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	"github.com/icon-project/centralized-relay/utils/multisig"
	"go.uber.org/zap"
)

func TestParseMessageFromTx(t *testing.T) {
	// Create a mock Provider
	chainParam := &chaincfg.TestNet3Params
	logger, _ := zap.NewDevelopment()
	provider := &Provider{
		logger: logger,
		cfg: &Config{
			Connections: []string{"cx577f5e756abd89cbcba38a58508b60a12754d2f5"},
		},
		chainParam: &chaincfg.TestNet3Params,
	}
	provider.cfg.NID = "0x2.btc"

	provider.cfg.Address = "tb1pf0atpt2d3zel6udws38pkrh2e49vqd3c5jcud3a82srphnmpe55q0ecrzk"
	msPubkey, _ := btcutil.DecodeAddress(provider.cfg.Address, chainParam)
	multisigAddressScript, _ := txscript.PayToAddrScript(msPubkey)
	provider.multisigAddrScript = multisigAddressScript
	// Create a mock TxSearchRes
	// relayer multisig
	decodedAddr, _ := btcutil.DecodeAddress(RELAYER_MULTISIG_ADDRESS, chainParam)
	relayerPkScript, _ := txscript.PayToAddrScript(decodedAddr)
	// user key
	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	bridgeMsg := multisig.BridgeDecodedMsg{
		Message: &multisig.XCallMessage{
			MessageType:  0,
			Action:       "Deposit",
			TokenAddress: "0:1",
			To:           "0x2.icon/hx01ca85287d6342722fe733c25667676b9cf9f8a4",
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su",
			Amount:       new(big.Int).SetUint64(15000).Bytes(),
			Data:         []byte(""),
		},
		ChainId:  1,
		Receiver: "cxfc86ee7687e1bf681b5548b2667844485c0e7192",
		Connectors: []string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	}

	inputs := []*multisig.Input{
		{
			TxHash:       "16de7df933dacd95b0d3af7325a5a2e680a1b7dd447a97e7678d8dfa1ac750b4",
			OutputIdx:    4,
			OutputAmount: 445000,
			PkScript:     userMultisigWallet.PKScript,
		},
	}

	// create tx
	msgTx, err := multisig.CreateBridgeTxSendBitcoin(
		&bridgeMsg,
		inputs,
		userMultisigWallet.PKScript,
		relayerPkScript,
		TX_FEE,
	)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	// user key 1 sign tx
	userSigs1, _ := multisig.SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs1)
	// user key 2 sign tx
	userSigs2, _ := multisig.SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs2)
	// COMBINE SIGN
	signedMsgTx, _ := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	bridgeMessage, err := multisig.ReadBridgeMessage(signedMsgTx)
	fmt.Println("bridgeMessage: ", bridgeMessage)
	fmt.Println("err: ", err)
	txSearchRes := &TxSearchRes{
		Tx:            signedMsgTx,
		Height:        3181070,
		TxIndex:       21,
		BridgeMessage: bridgeMessage,
	}
	relayerMessage, err := provider.parseMessageFromTx(txSearchRes)
	fmt.Println("relayerMessage Src: ", relayerMessage.Src)
	fmt.Println("relayerMessage Sn: ", relayerMessage.Sn)
	fmt.Println("relayerMessage Data: ", relayerMessage.Data)
	fmt.Println("err: ", err)

	msg := &types.RecvMessage{
		SrcNID: relayerMessage.Src,
		ConnSn: types.NewHexInt(relayerMessage.Sn.Int64()),
		Msg:    types.NewHexBytes(relayerMessage.Data),
	}
	iconProvider := &icon.Provider{}

	iconMessage := iconProvider.NewIconMessage(iconProvider.GetAddressByEventType(relayerMessage.EventType), msg, icon.MethodRecvMessage)
	fmt.Println("iconMessage Method: ", iconMessage.Method)
	fmt.Println("iconMessage Params: ", iconMessage.Params)
	// iconProvider.SendTransaction(nil, iconMessage)
}
