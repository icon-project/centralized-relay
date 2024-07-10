package multisig

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

func toXOnly(pubKey []byte) []byte {
	if len(pubKey) == 33 {
		return pubKey[1:33]
	}

	return pubKey
}

// use OP_CHECKSIGADD instead of OP_CHECKMULTISIG legacy
func buildMultisigTapScript(numSigsRequired int, pubKeys [][]byte) ([]byte, string, error) {
	builder := txscript.NewScriptBuilder()

	// the first pubkey
	builder.AddData(toXOnly(pubKeys[0]))
	builder.AddOp(txscript.OP_CHECKSIG)

	// the remaining pubkeys
	for i := 1; i < len(pubKeys); i++ {
		builder.AddData(toXOnly(pubKeys[i]))
		builder.AddOp(txscript.OP_CHECKSIGADD)
	}

	// add number of required sigs
	builder.AddInt64(int64(numSigsRequired))
	builder.AddOp(txscript.OP_NUMEQUAL)

	redeemScript, err := builder.Script()
	if err != nil {
		return []byte{}, "", fmt.Errorf("could not build script - Error %v", err)
	}

	return redeemScript, "", nil
}

// Point time lock contract script
// use OP_CHECKLOCKTIMEVERIFY
func buildPTLCTapScript(
	expiredBlkHeight uint64,
	pubKey []byte,
) ([]byte, string, error) {

	// <blockHeight> OP_CHECKLOCKTIMEVERIFY OP_DROP <PubKey> OP_CHECKSIG

	builder := txscript.NewScriptBuilder()
	// builder.AddData(new(big.Int).SetUint64(expiredBlkHeight).Bytes()) // TODO: need to fixed length?
	builder.AddInt64(int64(expiredBlkHeight))
	builder.AddOp(txscript.OP_CHECKLOCKTIMEVERIFY)
	builder.AddOp(txscript.OP_DROP)
	builder.AddData(toXOnly(pubKey))
	builder.AddOp(txscript.OP_CHECKSIG)

	scriptBytes, err := builder.Script()
	if err != nil {
		return []byte{}, "", fmt.Errorf("buildPTLCTapScript could not build script - Error %v", err)
	}

	return scriptBytes, "", nil
}

func buildMultisigWalletFromScripts(scripts [][]byte) (*MultisigWallet, error) {
	tapLeaves := []txscript.TapLeaf{}
	for _, script := range scripts {
		tapLeaf := txscript.NewBaseTapLeaf(script)
		tapLeaves = append(tapLeaves, tapLeaf)
	}

	tapScriptTree := txscript.AssembleTaprootScriptTree(tapLeaves...)
	tapScriptRootHash := tapScriptTree.RootNode.TapHash()

	sharedRandomBytes, _ := hex.DecodeString(SHARED_RANDOM_HEX)
	sharedRandom := new(big.Int).SetBytes(sharedRandomBytes)

	sharedPublicKey, _, err := genSharedInternalPubKey(sharedRandom)
	if err != nil {
		return nil, err
	}

	outputKey := txscript.ComputeTaprootOutputKey(
		sharedPublicKey, tapScriptRootHash[:],
	)

	pkScript, err := txscript.PayToTaprootScript(outputKey)
	if err != nil {
		return nil, fmt.Errorf("build taproot PK script err %v", err)
	}

	return &MultisigWallet{
		PKScript:        pkScript,
		SharedPublicKey: sharedPublicKey,

		TapScriptTree: tapScriptTree,
		TapLeaves:     tapLeaves,
	}, nil
}

func computeYCoordinate(x *big.Int) *big.Int {
	// secp256k1 curve parameters
	params := btcec.S256()

	// Compute y^2 = x^3 + 7 (mod p)
	xCubed := new(big.Int).Exp(x, big.NewInt(3), params.P)
	ySquared := new(big.Int).Add(xCubed, big.NewInt(7))
	ySquared.Mod(ySquared, params.P)

	// Compute y-coordinate using square root modulo p
	y := new(big.Int).ModSqrt(ySquared, params.P)
	return y
}

func genSharedInternalPubKey(sharedRandom *big.Int) (*btcec.PublicKey, []byte, error) {

	// P = H + rG
	uncompressGBytes := []byte{4}
	uncompressGBytes = append(uncompressGBytes, btcec.S256().CurveParams.Params().Gx.Bytes()...)
	uncompressGBytes = append(uncompressGBytes, btcec.S256().CurveParams.Params().Gy.Bytes()...)
	hashGBytes := sha256.Sum256(uncompressGBytes)

	xH := new(big.Int).SetBytes(hashGBytes[:])
	yH := computeYCoordinate(xH)
	isValidPoint := btcec.S256().IsOnCurve(xH, yH)
	if !isValidPoint {
		return nil, nil, fmt.Errorf("can not generate H point from hash of G")
	}

	xrG, yrG := btcec.S256().ScalarBaseMult(sharedRandom.Bytes())
	xP, yP := btcec.S256().Add(xH, yH, xrG, yrG)
	isValidPoint = btcec.S256().IsOnCurve(xP, yP)
	if !isValidPoint {
		return nil, nil, fmt.Errorf("can not generate P point")
	}
	xField := &btcec.FieldVal{}
	xField.SetBytes((*[32]byte)(xP.Bytes()))

	yField := &btcec.FieldVal{}
	yField.SetBytes((*[32]byte)(yP.Bytes()))

	publicKey := btcec.NewPublicKey(xField, yField)

	return publicKey, toXOnly(publicKey.SerializeCompressed()), nil

}

// create multisig struct contain multisig wallet detail
// input: multisig info (public keys, number of sigs required)
// output: multisig struct
func BuildMultisigWallet(
	multisigInfo *MultisigInfo,
) (*MultisigWallet, error) {
	// Taptree structure:
	// TapLeaf 1: <MULTISIG_SCRIPT>
	// TapLeaf 2: <PTLC>
	script1, _, err := buildMultisigTapScript(multisigInfo.NumberRequiredSigs, multisigInfo.PubKeys)
	if err != nil {
		return nil, fmt.Errorf("build script multisig err %v", err)
	}

	if multisigInfo.RecoveryBlockHeight == 0 {
		return buildMultisigWalletFromScripts([][]byte{script1})
	} else {
		script2, _, err := buildPTLCTapScript(multisigInfo.RecoveryBlockHeight, multisigInfo.RecoveryPubKey)
		if err != nil {
			return nil, fmt.Errorf("build script PTLC err %v", err)
		}

		return buildMultisigWalletFromScripts([][]byte{script1, script2})
	}
}

func AddressOnChain(
	chainParam *chaincfg.Params,
	multisigWallet *MultisigWallet,
) (*btcutil.AddressTaproot, error) {
	tapScriptRootHash := multisigWallet.TapScriptTree.RootNode.TapHash()

	outputKey := txscript.ComputeTaprootOutputKey(
		multisigWallet.SharedPublicKey, tapScriptRootHash[:],
	)

	address, err := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(outputKey), chainParam)
	if err != nil {
		return nil, fmt.Errorf("build address from script err %v", err)
	}

	return address, nil
}