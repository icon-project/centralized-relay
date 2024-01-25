package evm

import (
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/icon-project/centralized-relay/relayer/events"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
)

// All the events
var (
	EmitMessage = "Message(string,uint256,bytes)"
	CallMessage = "CallMessage(str,str,uint256,uint256,bytes)"
)

func eventSigToEventType(sigContract map[string]string) map[common.Hash]string {
	eventMap := make(map[common.Hash]string, len(sigContract))
	for sig, contract := range sigContract {
		switch contract {
		case providerTypes.ConnectionContract:
			eventMap[crypto.Keccak256Hash([]byte(sig))] = events.EmitMessage
		case providerTypes.XcallContract:
			eventMap[crypto.Keccak256Hash([]byte(sig))] = events.CallMessage
		}
	}
	return eventMap
}

func MonitorEventsList(sigContract map[string]string) []string {
	eventsList := []string{}
	for sig, contract := range sigContract {
		switch contract {
		case providerTypes.ConnectionContract:
			eventsList = append(eventsList, EmitMessage)
		case providerTypes.XcallContract:
			eventsList = append(eventsList, CallMessage)
		}
	}
	return eventsList
}

var MonitorEvents []common.Hash = []common.Hash{
	// TODO: list all the events to monitor
	crypto.Keccak256Hash([]byte(EmitMessage)),
}

func getEventFilterQuery(contractAddress string) ethereum.FilterQuery {
	return ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(contractAddress)},
		Topics: [][]common.Hash{
			MonitorEvents,
		},
	}
}
