package multisig

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
)

type MultisigInfo struct {
	PubKeys            [][]byte
	EcPubKeys          []*btcutil.AddressPubKey
	NumberRequiredSigs int
}

type MultisigWallet struct {
	TapScriptTree *txscript.IndexedTapScriptTree
	TapLeaves     []txscript.TapLeaf
	PKScript        []byte
	SharedPublicKey *btcec.PublicKey
}

type OutputTx struct {
	ReceiverAddress string
	Amount          uint64
}

type UTXO struct {
	IsRelayersMultisig	bool `bson:"is_relayers_multisig" json:"isRelayersMultisig"`
	TxHash        		string `bson:"tx_hash" json:"txHash"`
	OutputIdx     		uint32 `bson:"output_idx" json:"outputIdx"`
	OutputAmount  		uint64 `bson:"output_amount" json:"outputAmount"`
}
