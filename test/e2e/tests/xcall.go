package tests

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/icon-project/centralized-relay/test/chains"
	"github.com/icon-project/centralized-relay/test/testsuite"
	"github.com/stretchr/testify/assert"
)

const (
	CS_RESP_FAILURE = "0"
	CS_RESP_SUCCESS = "1"
)

type XCallTestSuite struct {
	*testsuite.E2ETestSuite
	T *testing.T
}

func (x *XCallTestSuite) TextXCall() {
	testcase := "xcall"
	ctx := context.WithValue(context.Background(), "testcase", testcase)
	createdChains := x.GetChains()
	testChains(ctx, createdChains, x)

}

func isInList(processedList []string, item string) bool {
	for _, val := range processedList {
		if val == item {
			return true
		}
	}
	return false
}

func testChains(ctx context.Context, createdChains []chains.Chain, x *XCallTestSuite) {
	var processedList []string
	for index, chain := range createdChains {
		for innerIndex, innerChain := range createdChains {
			if index != innerIndex {
				chainFlowIdentifier := chain.Config().Name + "-" + innerChain.Config().Name
				if !isInList(processedList, chainFlowIdentifier) {
					processedList = append(processedList, chainFlowIdentifier)
					chainFlowName := chain.Config().Name + "->" + innerChain.Config().Name
					x.T.Run("xcall one way message chainA-chainB "+chainFlowName, func(t *testing.T) {
						fmt.Println("Sending message from src to dst", chain.Config().Name, innerChain.Config().Name)
						err := x.testOneWayMessage(ctx, t, chain, innerChain)
						assert.NoErrorf(t, err, "fail xCall one way message chainA-chainB( %s) ::%v\n ", chainFlowIdentifier, err)
					})
					x.T.Run("xcall test rollback chainA-chainB "+chainFlowName, func(t *testing.T) {
						fmt.Println("Sending rollback message from src to dst", chain.Config().Name, innerChain.Config().Name)
						err := x.testRollback(ctx, t, chain, innerChain)
						assert.NoErrorf(t, err, "fail xCall rollback message chainA-chainB( %s) ::%v\n ", chainFlowIdentifier, err)
					})

					x.T.Run("xcall test rollback data chainA-chainB without rollback "+chainFlowName, func(t *testing.T) {
						fmt.Println("Sending rollback message from src to dst", chain.Config().Name, innerChain.Config().Name)
						err := x.testRollbackDataWithoutRollback(ctx, t, chain, innerChain)
						assert.NoErrorf(t, err, "fail xCall rollback message chainA-chainB( %s) ::%v\n ", chainFlowIdentifier, err)
					})

					x.T.Run("xcall test rollback data reply chainA-chainB without rollback "+chainFlowName, func(t *testing.T) {
						fmt.Println("Sending rollback message from src to dst", chain.Config().Name, innerChain.Config().Name)
						err := x.testRollbackDataReplyWithoutRollback(ctx, t, chain, innerChain)
						assert.NoErrorf(t, err, "fail xCall rollback message chainA-chainB( %s) ::%v\n ", chainFlowIdentifier, err)
					})

					x.T.Run("xcall test send maxSize Data: 2048 bytes A-> B "+chainFlowName, func(t *testing.T) {
						fmt.Println("Sending allowed size data from src to dst", chain.Config().Name, innerChain.Config().Name)
						x.testOneWayMessageWithSize(ctx, t, 1300, chain, innerChain)
					})

					x.T.Run("xcall test send maxSize Data: 2049bytes  "+chainFlowName, func(t *testing.T) {
						fmt.Println("Sending more than max  size data from src to dst", chain.Config().Name, innerChain.Config().Name)
						x.testOneWayMessageWithSizeExpectingError(ctx, t, 2000, chain, innerChain)
					})

				}
				reverseChainFlowIdentifier := innerChain.Config().Name + "-" + chain.Config().Name
				if !isInList(processedList, reverseChainFlowIdentifier) {
					processedList = append(processedList, reverseChainFlowIdentifier)
					reverseChainFlowName := innerChain.Config().Name + "->" + chain.Config().Name
					x.T.Run("xcall one way message chainB-chainA "+reverseChainFlowName, func(t *testing.T) {
						fmt.Println("Sending message from src to dst", innerChain.Config().Name, chain.Config().Name)
						err := x.testOneWayMessage(ctx, t, innerChain, chain)
						assert.NoErrorf(t, err, "fail xCall one way message chainB-chainA (%s) ::%v  \n ", reverseChainFlowIdentifier, err)
					})

					x.T.Run("xcall test rollback chainB-chainA"+reverseChainFlowName, func(t *testing.T) {
						fmt.Println("Sending rollback message from src to dst", chain.Config().Name, innerChain.Config().Name)
						err := x.testRollback(ctx, t, innerChain, chain)
						assert.NoErrorf(t, err, "fail xCall rollback message chainB-chainA( %s) ::%v\n ", reverseChainFlowIdentifier, err)
					})

					x.T.Run("xcall test rollback data chainB-chainA without rollback "+reverseChainFlowName, func(t *testing.T) {
						fmt.Println("Sending rollback message from src to dst", chain.Config().Name, innerChain.Config().Name)
						err := x.testRollbackDataWithoutRollback(ctx, t, innerChain, chain)
						assert.NoErrorf(t, err, "fail xCall rollback message chainB-chainA( %s) ::%v\n ", reverseChainFlowIdentifier, err)
					})

					x.T.Run("xcall test rollback data reply data chainB-chainA without rollback "+reverseChainFlowName, func(t *testing.T) {
						fmt.Println("Sending rollback message from src to dst", chain.Config().Name, innerChain.Config().Name)
						err := x.testRollbackDataReplyWithoutRollback(ctx, t, innerChain, chain)
						assert.NoErrorf(t, err, "fail xCall rollback message chainB-chainA( %s) ::%v\n ", reverseChainFlowIdentifier, err)
					})

					x.T.Run("xcall test send maxSize Data: 2048 bytes B-> A "+reverseChainFlowName, func(t *testing.T) {
						fmt.Println("Sending allowed size data from src to dst", innerChain.Config().Name, chain.Config().Name)
						x.testOneWayMessageWithSize(ctx, t, 1300, innerChain, chain)
					})

					x.T.Run("xcall test send maxSize Data: 2049bytes "+reverseChainFlowName, func(t *testing.T) {
						fmt.Println("ending more than max  size data from src to dst", innerChain.Config().Name, chain.Config().Name)
						x.testOneWayMessageWithSizeExpectingError(ctx, t, 2000, innerChain, chain)
					})
				}

			}
		}
	}
}

func handlePanicAndGetContractAddress(chain chains.Chain, contractName, fallbackContractName string) (address string) {
	defer func() {
		if r := recover(); r != nil {
			address = chain.GetContractAddress(fallbackContractName)
			return
		}
	}()
	address = chain.GetContractAddress(contractName)
	return address
}

func (x *XCallTestSuite) testOneWayMessage(ctx context.Context, t *testing.T, chainA, chainB chains.Chain) error {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	msg := "MessageTransferTestingWithoutRollback"
	dAppAddress := handlePanicAndGetContractAddress(chainB, dappKey+"-idcap", dappKey)
	dst := chainB.Config().ChainID + "/" + dAppAddress
	res, err := chainA.XCall(ctx, chainB, chainB.Config().Name, dst, []byte(msg), nil)
	result := assert.NoErrorf(t, err, "error on sending packet- %v", err)
	if !result {
		return err
	}
	dataOutput := x.ConvertToPlainString(res.Data)
	result = assert.NoErrorf(t, err, "error on converting res data as msg- %v", err)
	if !result {
		return err
	}
	result = assert.Equal(t, msg, dataOutput)
	if !result {
		return err
	}
	fmt.Println("Data Transfer Testing Without Rollback from " + chainA.Config().ChainID + " to " + chainB.Config().ChainID + " with data " + msg + " and Received:" + dataOutput + " PASSED")
	return nil
}

func (x *XCallTestSuite) testRollback(ctx context.Context, t *testing.T, chainA, chainB chains.Chain) error {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	msg := "rollback"
	rollback := "RollbackDataTesting"
	dAppAddress := handlePanicAndGetContractAddress(chainB, dappKey+"-idcap", dappKey)
	dst := chainB.Config().ChainID + "/" + dAppAddress
	res, err := chainA.XCall(ctx, chainB, chainB.Config().Name, dst, []byte(msg), []byte(rollback))
	isSuccess := assert.NoErrorf(t, err, "error on sending packet- %v", err)
	if !isSuccess {
		return err
	}
	height, err := chainA.Height(ctx)
	assert.NoErrorf(t, err, "error getting height %v", err)
	code, err := chainA.FindCallResponse(ctx, height, res.SerialNo)
	assert.NoErrorf(t, err, "no call response found %v", err)
	isSuccess = assert.Equal(t, CS_RESP_FAILURE, code)
	if !isSuccess {
		return err
	}
	_, err = chainA.FindRollbackExecutedMessage(ctx, height, res.SerialNo)
	assert.NoErrorf(t, err, "no rollback executed message found %v", err)
	fmt.Println("Data Transfer Testing With Rollback from " + chainA.Config().ChainID + " to " + chainB.Config().ChainID + " with data " + msg + " and rollback:" + rollback + " PASSED")
	return err
}

func (x *XCallTestSuite) testRollbackDataWithoutRollback(ctx context.Context, t *testing.T, chainA, chainB chains.Chain) error {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	msg := "MessageTransferTestingWithoutRollback"
	rollback := "rollbackData"
	dAppAddress := handlePanicAndGetContractAddress(chainB, dappKey+"-idcap", dappKey)
	dst := chainB.Config().ChainID + "/" + dAppAddress
	res, err := chainA.XCall(ctx, chainB, chainB.Config().Name, dst, []byte(msg), []byte(rollback))
	isSuccess := assert.NoErrorf(t, err, "error on sending packet- %v", err)
	if !isSuccess {
		return err
	}
	height, err := chainA.Height(ctx)
	assert.NoErrorf(t, err, "error getting height %v", err)
	code, err := chainA.FindCallResponse(ctx, height, res.SerialNo)
	assert.NoErrorf(t, err, "no call response found %v", err)
	isSuccess = assert.Equal(t, CS_RESP_SUCCESS, code)
	if !isSuccess {
		return err
	}
	fmt.Println("Data Transfer Testing Without Rollback from " + chainA.Config().ChainID + " to " + chainB.Config().ChainID + " with data " + msg + " and rollback:" + rollback + " PASSED")
	return nil
}

func (x *XCallTestSuite) testRollbackDataReplyWithoutRollback(ctx context.Context, t *testing.T, chainA, chainB chains.Chain) error {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	msg := "reply-reponse"
	rollback := "rollbackData"
	dAppAddress := handlePanicAndGetContractAddress(chainB, dappKey+"-idcap", dappKey)
	dst := chainB.Config().ChainID + "/" + dAppAddress
	res, err := chainA.XCall(ctx, chainB, chainB.Config().Name, dst, []byte(msg), []byte(rollback))
	isSuccess := assert.NoErrorf(t, err, "error on sending packet- %v", err)
	if !isSuccess {
		return err
	}
	height, err := chainA.Height(ctx)
	assert.NoErrorf(t, err, "error getting height %v", err)
	code, err := chainA.FindCallResponse(ctx, height, res.SerialNo)
	assert.NoErrorf(t, err, "no call response found %v", err)
	isSuccess = assert.Equal(t, CS_RESP_SUCCESS, code)
	if !isSuccess {
		return err
	}
	time.Sleep(3 * time.Second)
	fmt.Println("Data Transfer Testing Without Rollback from " + chainA.Config().ChainID + " to " + chainB.Config().ChainID + " with data " + msg + " and rollback:" + rollback + " PASSED")
	return err
}

func (x *XCallTestSuite) testOneWayMessageWithSize(ctx context.Context, t *testing.T, dataSize int, chainA, chainB chains.Chain) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	_msg := make([]byte, dataSize)

	dAppAddress := handlePanicAndGetContractAddress(chainB, dappKey+"-idcap", dappKey)
	dst := chainB.Config().ChainID + "/" + dAppAddress
	res, err := chainA.XCall(ctx, chainB, chainB.Config().Name, dst, _msg, nil)
	assert.NoErrorf(t, err, "error on sending packet- %v", err)
	assert.NotEmpty(t, res.RequestID, "retrieved requestId should not be empty")
	assert.NoError(t, err)
	fmt.Println("Data Transfer Testing With Message Size from " + chainA.Config().ChainID + " to " + chainB.Config().ChainID + " with data " + string(_msg) + " PASSED")
}

func (x *XCallTestSuite) testOneWayMessageWithSizeExpectingError(ctx context.Context, t *testing.T, dataSize int, chainA, chainB chains.Chain) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	_msg := make([]byte, dataSize)
	dAppAddress := handlePanicAndGetContractAddress(chainB, dappKey+"-idcap", dappKey)
	dst := chainB.Config().ChainID + "/" + dAppAddress
	_, err := chainA.XCall(ctx, chainB, chainB.Config().Name, dst, _msg, nil)
	result := assert.Errorf(t, err, "large data transfer should failed")
	if result {
		result = false
		if strings.Contains(err.Error(), "submessages:") {
			subStart := strings.Index(err.Error(), "submessages:") + len("submessages:")
			subEnd := strings.Index(err.Error(), ": execute")
			subMsg := err.Error()[subStart:subEnd]
			result = assert.ObjectsAreEqual(strings.TrimSpace(subMsg), "MaxDataSizeExceeded")
		} else if strings.Contains(err.Error(), "MaxDataSizeExceeded") {
			result = true
		} else {
			result = assert.ObjectsAreEqual(errors.New("UnknownFailure"), err)
		}
		if result {
			t.Logf("Test passed: %v", err)
			fmt.Println("Data Transfer Testing With Message Size expecting error from " + chainA.Config().ChainID + " to " + chainB.Config().ChainID + " with data " + string(_msg) + " PASSED")
		} else {
			t.Errorf("Test failed: %v", err)
		}
	}

}
