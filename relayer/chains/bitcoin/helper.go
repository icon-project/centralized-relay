package bitcoin

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	bitcoinABI "github.com/icon-project/centralized-relay/relayer/chains/bitcoin/abi"
	"github.com/icon-project/centralized-relay/utils/multisig"
	"github.com/icon-project/icon-bridge/common/codec"
)

func GetRuneTxIndex(endpoint, method, bearToken, txId string, index int) (*RuneTxIndexResponse, error) {
	client := &http.Client{}
	endpoint = endpoint + "/runes/utxo/" + txId + "/" + strconv.FormatUint(uint64(index), 10) + "/balance"
	fmt.Println(endpoint)
	req, err := http.NewRequest(method, endpoint, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Authorization", bearToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var resp *RuneTxIndexResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp, nil
}

func ToXCallMessage(data interface{}, from, to string, sn uint, protocols []string, requester, token0, token1 common.Address) ([]byte, error) {
	var res []byte

	//bitcoinStateAbi, _ := abi.JSON(strings.NewReader(bitcoinABI.BitcoinStateMetaData.ABI))
	nonfungibleABI, _ := abi.JSON(strings.NewReader(bitcoinABI.InonfungibleTokenMetaData.ABI))
	addressTy, _ := abi.NewType("address", "", nil)
	bytes, _ := abi.NewType("bytes", "", nil)

	arguments := abi.Arguments{
		{
			Type: addressTy,
		},
		{
			Type: bytes,
		},
	}

	amount0, _ := big.NewInt(0).SetString("18999999999999999977305673", 10)

	switch data.(type) {
	case multisig.RadFiProvideLiquidityMsg:
		dataMint := data.(multisig.RadFiProvideLiquidityMsg)
		mintParams := bitcoinABI.INonfungiblePositionManagerMintParams{
			Token0:         token0,
			Token1:         token1,
			Fee:            big.NewInt(int64(dataMint.Detail.Fee) * 100),
			TickLower:      big.NewInt(int64(dataMint.Detail.LowerTick)),
			TickUpper:      big.NewInt(int64(dataMint.Detail.UpperTick)),
			Amount0Desired: amount0,
			Amount1Desired: big.NewInt(539580403982610478),
			Recipient:      common.HexToAddress(to),
			Deadline:       big.NewInt(1000000000),
		}

		mintParams.Amount0Min = mulDiv(mintParams.Amount0Desired, big.NewInt(int64(dataMint.Detail.Min0)), big.NewInt(1e4))
		mintParams.Amount1Min = mulDiv(mintParams.Amount1Desired, big.NewInt(int64(dataMint.Detail.Min1)), big.NewInt(1e4))

		// encode
		// todo: for init pool
		//initPoolCalldata, err := bitcoinStateAbi.Pack("initPool", mintParams, "btc", "rad", 1e0)
		//if err != nil {
		//	return nil, err
		//}

		provideLiquidity, err := nonfungibleABI.Pack("mint", mintParams)
		if err != nil {
			return nil, err
		}

		// encode with requester
		provideLiquidity, err = arguments.Pack(requester, provideLiquidity)
		if err != nil {
			return nil, err
		}

		//from := "0x3.BTC/bc1qvqkshkdj67uwvlwschyq8wja6df4juhewkg5fg"

		// encode to xcall format
		res, err = XcallFormat(provideLiquidity, from, to, sn, protocols)
		if err != nil {
			return nil, err
		}

	case multisig.RadFiWithdrawLiquidityMsg:

	default:
		return nil, fmt.Errorf("not supported")
	}
	return res, nil
}

func XcallFormat(callData []byte, from, to string, sn uint, protocols []string) ([]byte, error) {
	//
	csV2 := CSMessageRequestV2{
		From:        from,
		To:          to,
		Sn:          big.NewInt(int64(sn)).Bytes(),
		MessageType: uint8(CALL_MESSAGE_TYPE),
		Data:        callData,
		Protocols:   protocols,
	}

	//
	cvV2EncodeMsg, err := codec.RLP.MarshalToBytes(csV2)
	if err != nil {
		return nil, err
	}

	message := CSMessage{
		MsgType: big.NewInt(int64(CS_REQUEST)).Bytes(),
		Payload: cvV2EncodeMsg,
	}

	//
	finalMessage, err := codec.RLP.MarshalToBytes(message)
	if err != nil {
		return nil, err
	}

	return finalMessage, nil
}

func mulDiv(a, nNumerator, nDenominator *big.Int) *big.Int {
	return big.NewInt(0).Div(big.NewInt(0).Mul(a, nNumerator), nDenominator)
}
