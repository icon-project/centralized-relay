package bitcoin

import (
	"math/big"

	"github.com/icon-project/icon-bridge/common/codec"
)

func XcallFormat(callData []byte, from, to string, sn uint, protocols []string, messType uint8) ([]byte, error) {
	//
	csV2 := CSMessageRequestV2{
		From:        from,
		To:          to,
		Sn:          big.NewInt(int64(sn)).Bytes(),
		MessageType: messType,
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
