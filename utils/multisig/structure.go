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
	RecoveryPubKey     []byte
	RecoveryBlockHeight uint64
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

type TapSigParams struct {
	TxSigHashes			*txscript.TxSigHashes `bson:"tx_sig_hashes" json:"txSigHashes"`
	RelayersPKScript	[]byte `bson:"relayers_PK_script" json:"relayersPKScript"`
	RelayersTapLeaf		txscript.TapLeaf `bson:"relayers_tap_leaf" json:"relayersTapLeaf"`
	UserPKScript		[]byte `bson:"user_PK_script" json:"userPKScript"`
	UserTapLeaf			txscript.TapLeaf `bson:"user_tap_leaf" json:"userTapLeaf"`
}
