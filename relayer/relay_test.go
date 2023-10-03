package relayer

import (
	"context"
	"fmt"
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
	levelDbName = "./testdb"
)

type RelayTestSuite struct {
	suite.Suite

	logger *zap.Logger
	db     *lvldb.LVLDB
	relay  Relayer
}

func TestRunTestRelaySuite(t *testing.T) {
	suite.Run(t, new(RelayTestSuite))
}

func GetMockMessages(srcChainId, dstchainId string, srcStartHeight uint64) map[types.MessageKey]types.Message {
	messages := []types.Message{
		{
			Src:           srcChainId,
			Dst:           dstchainId,
			Data:          []byte(fmt.Sprintf("from message %s", srcChainId)),
			MessageHeight: uint64(srcStartHeight + 3),
			Sn:            1,
			EventType:     "emitMessage",
		},
		{
			Src:           srcChainId,
			Dst:           dstchainId,
			Data:          []byte(fmt.Sprintf("from message %s", srcChainId)),
			MessageHeight: uint64(srcStartHeight + 5),
			Sn:            2,
			EventType:     "emitMessage",
		},
		{
			Src:           srcChainId,
			Dst:           dstchainId,
			Data:          []byte(fmt.Sprintf("from message %s", srcChainId)),
			MessageHeight: uint64(srcStartHeight + 7),
			Sn:            3,
			EventType:     "emitMessage",
		},
	}
	sendMockMessageMap := make(map[types.MessageKey]types.Message, 0)
	for _, m := range messages {
		sendMockMessageMap[m.MessageKey()] = m
	}
	return sendMockMessageMap
}

func GetMockChainProvider(log *zap.Logger, blockDuration time.Duration, srcChainId string, dstchainId string, srcStartHeight uint64, dstStartHeight uint64) (provider.ChainProvider, error) {
	sendMessages := GetMockMessages(srcChainId, dstchainId, srcStartHeight)
	receiveMessage := GetMockMessages(dstchainId, srcChainId, dstStartHeight)
	mock1ProviderConfig := mockchain.MockProviderConfig{
		ChainId:         srcChainId,
		BlockDuration:   blockDuration,
		StartHeight:     srcStartHeight,
		SendMessages:    sendMessages,
		ReceiveMessages: receiveMessage,
	}
	return mock1ProviderConfig.NewProvider(log, "empty", false, srcChainId)

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

func (s *RelayTestSuite) TearDownSuite() {
	// clearing all the db
	s.db.ClearStore()
	defer s.db.Close()

}

func (s *RelayTestSuite) TestListener() {

	mock1 := "mock-1"
	dstMock2 := "mock-2"
	srcStartHeight := uint64(10)
	mockProvider, err := GetMockChainProvider(s.logger, 500*time.Millisecond, mock1, dstMock2, srcStartHeight, 10)
	if err != nil {
		s.Fail("fail to create mockProvider", err)
	}

	mockMessages := GetMockMessages(mock1, dstMock2, srcStartHeight)

	chains := make(map[string]*Chain, 0)

	chains[mock1] = NewChain(s.logger, mockProvider, true)

	relayer, err := NewRelayer(s.logger, s.db, chains, true)
	if err != nil {
		s.Fail("failed to create relayer ")
	}

	errorchan := make(chan error, 1)
	ctx := context.Background()

	go relayer.StartChainListeners(ctx, errorchan)

	runtime, err := relayer.FindChainRuntime(mock1)
	if err != nil {
		s.Fail("failed to get chain runtime ")
	}

	listenerchan := runtime.listenerChan
loop:
	for {
		select {
		case err := <-errorchan:
			s.Fail("error occured ", err)

		case blockInfo := <-listenerchan:
			for _, m := range blockInfo.Messages {
				delete(mockMessages, m.MessageKey())
			}
			fmt.Println("mockmessage length ", len(mockMessages))
			if len(mockMessages) == 0 {
				break loop
			}

		}
	}
	defer close(listenerchan)
}

func (s *RelayTestSuite) TestRelay() {

	chains := make(map[string]*Chain, 0)

	logger, _ := zap.NewProduction()

	mock1ChainId := "mock-1"
	mock2ChainId := "mock-2"
	mock1StartHeight := 10
	mock2StartHeight := 20

	mock1Provider, err := GetMockChainProvider(s.logger, 500*time.Millisecond, mock1ChainId, mock2ChainId, uint64(mock1StartHeight), uint64(mock2StartHeight))
	if err != nil {
		s.Fail("fail to create mockProvider", err)
	}
	chains[mock1ChainId] = NewChain(logger, mock1Provider, true)

	mock2Provider, err := GetMockChainProvider(s.logger, 500*time.Millisecond, mock2ChainId, mock1ChainId, uint64(mock2StartHeight), uint64(mock1StartHeight))
	if err != nil {
		s.Fail("fail to create mockProvider", err)
	}

	chains[mock2ChainId] = NewChain(logger, mock2Provider, true)

	ctx := context.Background()
	errorchan, err := Start(ctx, s.logger, chains, 3*time.Second, true, s.db)
	if err != nil {
		s.Fail("unable to start the relayer ", err)
	}

	provider1 := mock1Provider.(*mockchain.MockProvider)
	provider2 := mock2Provider.(*mockchain.MockProvider)

	receivedTimer := time.NewTicker(5 * time.Second)
	failedReceived := time.NewTicker(1 * time.Minute)
loop:
	for {
		select {
		case err := <-errorchan:
			s.Fail("error occured when starting the relay", err)
			break

		case <-receivedTimer.C:

			if len(provider1.PCfg.ReceiveMessages) == 0 && len(provider2.PCfg.ReceiveMessages) == 0 {
				break loop
			}
		case <-failedReceived.C:
			s.Fail(" failed to receive all the messeages")
			return
		}

	}

}
