package relayer

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/icon-project/centralized-relay/relayer/chains/mockchain"
	"github.com/icon-project/centralized-relay/relayer/lvldb"
	"github.com/icon-project/centralized-relay/relayer/provider"
	"github.com/icon-project/centralized-relay/relayer/types"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

const (
	levelDbName = "/tmp/testdb1"
)

type RelayTestSuite struct {
	suite.Suite

	logger *zap.Logger
	db     *lvldb.LVLDB
}

func TestRunTestRelaySuite(t *testing.T) {
	suite.Run(t, new(RelayTestSuite))
}

func GetMockMessages(srcNId, dstNId string, srcStartHeight uint64) map[types.MessageKey]*types.Message {
	messages := []types.Message{
		{
			Src:           srcNId,
			Dst:           dstNId,
			Data:          []byte(fmt.Sprintf("from message %s", srcNId)),
			MessageHeight: uint64(srcStartHeight + 3),
			Sn:            big.NewInt(1),
			EventType:     "emitMessage",
		},
		{
			Src:           srcNId,
			Dst:           dstNId,
			Data:          []byte(fmt.Sprintf("from message %s", srcNId)),
			MessageHeight: uint64(srcStartHeight + 5),
			Sn:            big.NewInt(2),
			EventType:     "emitMessage",
		},
		{
			Src:           srcNId,
			Dst:           dstNId,
			Data:          []byte(fmt.Sprintf("from message %s", srcNId)),
			MessageHeight: uint64(srcStartHeight + 7),
			Sn:            big.NewInt(3),
			EventType:     "emitMessage",
		},
	}
	sendMockMessageMap := make(map[types.MessageKey]*types.Message, 0)
	for _, m := range messages {
		sendMockMessageMap[*m.MessageKey()] = &m
	}
	return sendMockMessageMap
}

func GetMockChainProvider(log *zap.Logger, blockDuration time.Duration, NId string, dstNId string, srcStartHeight uint64, dstStartHeight uint64) (provider.ChainProvider, error) {
	sendMessages := GetMockMessages(NId, dstNId, srcStartHeight)
	receiveMessage := GetMockMessages(dstNId, NId, dstStartHeight)
	mock1ProviderConfig := mockchain.MockProviderConfig{
		NId:             NId,
		BlockDuration:   blockDuration,
		StartHeight:     srcStartHeight,
		SendMessages:    sendMessages,
		ReceiveMessages: receiveMessage,
	}
	return mock1ProviderConfig.NewProvider(context.Background(), log, "empty", false, NId)
}

func (s *RelayTestSuite) SetupTest() {
	logger, _ := zap.NewProduction()
	db, err := lvldb.NewLvlDB(levelDbName)
	if err != nil {
		s.Fail("fail to create leveldb", err)
	}

	s.db = db
	s.logger = logger
}

func (s *RelayTestSuite) TestRelay() {
	chains := make(map[string]*Chain, 0)
	logger, _ := zap.NewProduction()

	mock1Nid := "mock-1"
	mock2Nid := "mock-2"
	mock1StartHeight := 10
	mock2StartHeight := 20

	mock1Provider, err := GetMockChainProvider(s.logger, 500*time.Millisecond, mock1Nid, mock2Nid, uint64(mock1StartHeight), uint64(mock2StartHeight))
	if err != nil {
		s.Fail("fail to create mockProvider", err)
	}
	chains[mock1Nid] = NewChain(logger, mock1Provider, true)

	mock2Provider, err := GetMockChainProvider(s.logger, 500*time.Millisecond, mock2Nid, mock1Nid, uint64(mock2StartHeight), uint64(mock1StartHeight))
	if err != nil {
		s.Fail("fail to create mockProvider", err)
	}

	chains[mock2Nid] = NewChain(logger, mock2Provider, true)

	ctx := context.Background()
	rly, err := NewRelayer(logger, s.db, chains, true, false)
	if err != nil {
		s.Fail("unable to start the relayer ", err)
	}
	errorchan, err := rly.Start(ctx, 1*time.Second, true)
	if err != nil {
		s.Fail("unable to start the relayer ", err)
	}

	provider1 := mock1Provider.(*mockchain.MockProvider)
	provider2 := mock2Provider.(*mockchain.MockProvider)
	receivedTimer := time.NewTicker(2 * time.Second)
	failedReceived := time.NewTicker(30 * time.Second)
loop:
	for {
		select {
		case err := <-errorchan:
			s.Fail("error occured when starting the relay", err)
		case <-receivedTimer.C:
			if len(provider1.PCfg.ReceiveMessages) == 0 && len(provider2.PCfg.ReceiveMessages) == 0 {
				break loop
			}
		case <-failedReceived.C:
			s.Fail(" failed to receive all the messeages")
			return
		}
	}

	s.T().Cleanup(func() {
		s.db.Close()
		s.db.RemoveDbFile(levelDbName)
	})
}
