package evm

import (
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/icon-project/centralized-relay/relayer/events"
)

// All the events
var (
	EmitMessageSig = "Message(string,int256,bytes)"
)

var eventSigToEventType = map[common.Hash]string{
	crypto.Keccak256Hash([]byte(EmitMessageSig)): events.EmitMessage,
}

var MonitorEvents []common.Hash = []common.Hash{
	//TODO: list all the events to monitor
	crypto.Keccak256Hash([]byte(EmitMessageSig)),
}

func getEventFilterQuery(contractAddress string) ethereum.FilterQuery {
	return ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(contractAddress)},
		Topics: [][]common.Hash{
			MonitorEvents,
		},
	}

}
