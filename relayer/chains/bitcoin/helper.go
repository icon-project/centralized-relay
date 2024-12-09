package bitcoin

import (
	"encoding/binary"
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

func uint64ToBytes(amount uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, amount)
	return bytes
}

// Helper function to get minimum of two uint64 values
func min(a, b uint64) uint64 {
	if a <= b {
		return a
	}
	return b
}
