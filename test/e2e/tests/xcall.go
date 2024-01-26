package tests

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	interchaintest "github.com/icon-project/centralized-relay/test"
	"github.com/icon-project/centralized-relay/test/chains"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/testsuite"
	"github.com/stretchr/testify/assert"
)

type XCallTestSuite struct {
	*testsuite.E2ETestSuite
	T *testing.T
}

func (x *XCallTestSuite) TextXCall() {
	testcase := "xcall"
	portId := "transfer"
	ctx := context.WithValue(context.TODO(), "testcase", testcase)
	x.Require().NoError(x.DeployXCallMockApp(ctx, portId), "fail to deploy xcall demo dapp")
	chainA, chainB := x.GetChains()
	x.T.Run("test xcall", func(t *testing.T) {
		x.T.Run("xcall one way message chainA-chainB", func(t *testing.T) {
			err := x.testOneWayMessage(ctx, t, chainA, chainB)
			assert.NoErrorf(t, err, "fail xCall one way message chainA-chainB ::%v\n ", err)
		})

		x.T.Run("xcall one way message chainB-chainA", func(t *testing.T) {
			err := x.testOneWayMessage(ctx, t, chainB, chainA)
			assert.NoErrorf(t, err, "fail xCall one way message chainB-chainA ::%v\n ", err)
		})

		x.T.Run("xcall test rollback chainA-chainB", func(t *testing.T) {
			err := x.testRollback(ctx, t, chainA, chainB)
			assert.NoErrorf(t, err, "fail xCall rollback message chainB-chainA ::%v\n ", err)

		})

		x.T.Run("xcall test rollback chainB-chainA", func(t *testing.T) {
			err := x.testRollback(ctx, t, chainB, chainA)
			assert.NoErrorf(t, err, "fail xcCll rollback message chainB-chainA ::%v\n ", err)

		})

		x.T.Run("xcall test send maxSize Data: 2048 bytes", func(t *testing.T) {
			x.testOneWayMessageWithSize(ctx, t, 1800, chainA, chainB)
			x.testOneWayMessageWithSize(ctx, t, 1800, chainB, chainA)
		})

		x.T.Run("xcall test send maxSize Data: 2049bytes", func(t *testing.T) {
			x.testOneWayMessageWithSizeExpectingError(ctx, t, 2049, chainB, chainA)
			x.testOneWayMessageWithSizeExpectingError(ctx, t, 2049, chainA, chainB)
		})
	})
	//TC for sendNewMessage Xcall

	x.T.Run("test Newxcall", func(t *testing.T) {
		x.T.Run("xcall one way new message chainA-chainB", func(t *testing.T) {
			err := x.testOneWayMessage(ctx, t, chainA, chainB, true)
			assert.NoErrorf(t, err, "fail xCall one way message chainA-chainB ::%v\n ", err)
		})

		x.T.Run("xcall one way new message chainB-chainA", func(t *testing.T) {
			err := x.testOneWayMessage(ctx, t, chainB, chainA, true)
			assert.NoErrorf(t, err, "fail xCall one way message chainB-chainA ::%v\n ", err)
		})

		x.T.Run("xcall test new rollback chainA-chainB", func(t *testing.T) {
			err := x.testRollback(ctx, t, chainA, chainB, true)
			assert.NoErrorf(t, err, "fail xCall rollback message chainB-chainA ::%v\n ", err)

		})

		x.T.Run("xcall test new rollback chainB-chainA", func(t *testing.T) {
			err := x.testRollback(ctx, t, chainB, chainA, true)
			assert.NoErrorf(t, err, "fail xcCll rollback message chainB-chainA ::%v\n ", err)

		})

		x.T.Run("xcall test newsend maxSize Data: <2048 bytes", func(t *testing.T) {
			x.testOneWayMessageWithSize(ctx, t, 1800, chainA, chainB, true)
			x.testOneWayMessageWithSize(ctx, t, 1800, chainB, chainA, true)
		})

		x.T.Run("xcall test newsend maxSize Data: 2049bytes", func(t *testing.T) {
			x.testOneWayMessageWithSizeExpectingError(ctx, t, 2049, chainB, chainA, true)
			x.testOneWayMessageWithSizeExpectingError(ctx, t, 2049, chainA, chainB, true) // failing need to modify error checks
		})
	})

}

func (x *XCallTestSuite) TestXCallFlush() {
	testcase := "packet-flush"
	portId := "transfer-1"
	ctx := context.WithValue(context.TODO(), "testcase", testcase)
	x.Require().NoError(x.DeployXCallMockApp(ctx, portId), "fail to deploy xcall dapp")
	chainA, chainB := x.GetChains()
	x.T.Run("test xcall packet flush", func(t *testing.T) {
		x.T.Run("xcall packet flush chainA-chainB", func(t *testing.T) {
			err := x.testPacketFlush(ctx, chainA, chainB)
			assert.NoErrorf(t, err, "xcall packet flush chainA-chainB ::%v\n ", err)

		})

		x.T.Run("xcall packet flush chainB-chainA", func(t *testing.T) {
			err := x.testPacketFlush(ctx, chainB, chainA)
			assert.NoErrorf(t, err, "xcall packet flush chainB-chainA ::%v\n ", err)
		})
	})

	//test flush for sendNewMessage with msgType
	x.T.Run("test Newxcall packet flush", func(t *testing.T) {
		x.T.Run("newxcall packet flush chainA-chainB", func(t *testing.T) {
			err := x.testPacketFlush(ctx, chainA, chainB, true)
			assert.NoErrorf(t, err, "xcall packet flush chainA-chainB ::%v\n ", err)

		})

		x.T.Run("newxcall packet flush chainB-chainA", func(t *testing.T) {
			err := x.testPacketFlush(ctx, chainB, chainA, true)
			assert.NoErrorf(t, err, "xcall packet flush chainB-chainA ::%v\n ", err)
		})
	})
}

func (x *XCallTestSuite) testPacketFlush(ctx context.Context, chainA, chainB chains.Chain, newFunctionCall ...bool) error {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	msg := "flush-msg"
	heightB, _ := chainB.Height(ctx)

	dst := chainB.(ibc.Chain).Config().ChainID + "/" + chainB.GetContractAddress(dappKey)

	err := chainB.PauseNode(ctx)
	if err != nil {
		return fmt.Errorf("failed to pause node %s - %v", chainB.Config().Name, err)
	}
	if len(newFunctionCall) > 0 && newFunctionCall[0] {
		msgType := big.NewInt(1)
		ctx, err = chainA.SendNewPacketXCall(ctx, interchaintest.UserAccount, dst, []byte(msg), msgType, nil)
	} else {
		ctx, err = chainA.SendPacketXCall(ctx, interchaintest.UserAccount, dst, []byte(msg), nil)
	}

	if err != nil {
		return fmt.Errorf("failed send xCall message find eventlog - %v", err)
	}
	sn := ctx.Value("sn").(string)
	fmt.Printf("sn-%s\n", sn)
	waitDuration := 90 * time.Second
	fmt.Printf("Wait for %v \n", waitDuration)
	// TODO: Wait for 1.5 mins (90 seconds)
	time.Sleep(waitDuration)

	err = chainB.UnpauseNode(ctx)
	if err != nil {
		return fmt.Errorf("failed to unpause node %s - %v", chainB.Config().Name, err)
	}

	//wait 90 sec
	fmt.Printf("Wait for %v after node unpause\n", waitDuration)
	time.Sleep(waitDuration)

	reqId, destData, err := chainB.FindCallMessage(ctx, heightB, chainA.Config().ChainID+"/"+chainA.GetContractAddress(dappKey), chainB.GetContractAddress(dappKey), sn)
	if err != nil {
		return fmt.Errorf("error on execute call packet req-id::%s- %v", reqId, err)
	}
	_, err = chainB.ExecuteCall(ctx, reqId, destData)
	if err != nil {
		return fmt.Errorf("error on execute call packet req-id::%s- %v", reqId, err)
	}

	return nil
}

func (x *XCallTestSuite) testOneWayMessage(ctx context.Context, t *testing.T, chainA, chainB chains.Chain, newFunctionCall ...bool) error {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	msg := "MessageTransferTestingWithoutRollback"
	dst := chainB.(ibc.Chain).Config().ChainID + "/" + chainB.GetContractAddress(dappKey)
	var res *chains.XCallResponse
	var err error
	if len(newFunctionCall) > 0 && newFunctionCall[0] {
		msgType := big.NewInt(1)
		res, err = chainA.NewXCall(ctx, chainB, interchaintest.UserAccount, dst, []byte(msg), msgType, nil)
	} else {
		res, err = chainA.XCall(ctx, chainB, interchaintest.UserAccount, dst, []byte(msg), nil)
	}

	result := assert.NoErrorf(t, err, "error on sending packet- %v", err)
	if !result {
		return err
	}
	_, err = chainB.ExecuteCall(ctx, res.RequestID, res.Data)
	result = assert.NoErrorf(t, err, "error on execute call packet- %v", err)
	if !result {
		return err
	}
	//x.Require().NoErrorf(err, "error on execute call packet- %v", err)
	dataOutput := x.ConvertToPlainString(res.Data)
	//x.Require().NoErrorf(err, "error on converting res data as msg- %v", err)
	result = assert.NoErrorf(t, err, "error on converting res data as msg- %v", err)
	if !result {
		return err
	}
	result = assert.Equal(t, msg, dataOutput)
	if !result {
		return err
	}
	fmt.Println("Data Transfer Testing Without Rollback from " + chainA.(ibc.Chain).Config().ChainID + " to " + chainB.(ibc.Chain).Config().ChainID + " with data " + msg + " and Received:" + dataOutput + " PASSED")
	return nil
}

func (x *XCallTestSuite) testRollback(ctx context.Context, t *testing.T, chainA, chainB chains.Chain, newFunctionCall ...bool) error {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	msg := "rollback"
	rollback := "RollbackDataTesting"
	dst := chainB.(ibc.Chain).Config().ChainID + "/" + chainB.GetContractAddress(dappKey)
	var res *chains.XCallResponse
	var err error
	if len(newFunctionCall) > 0 && newFunctionCall[0] {
		msgType := big.NewInt(2)
		res, err = chainA.NewXCall(ctx, chainB, interchaintest.UserAccount, dst, []byte(msg), msgType, []byte(rollback))
	} else {
		res, err = chainA.XCall(ctx, chainB, interchaintest.UserAccount, dst, []byte(msg), []byte(rollback))
	}
	isSuccess := assert.NoErrorf(t, err, "error on sending packet- %v", err)
	if !isSuccess {
		return err
	}
	height, err := chainA.(ibc.Chain).Height(ctx)
	assert.NoErrorf(t, err, "error on getting height- %w", err)
	_, err = chainB.ExecuteCall(ctx, res.RequestID, res.Data)
	assert.NoErrorf(t, err, "error on excute call- %w", err)
	code, err := chainA.FindCallResponse(ctx, height, res.SerialNo)
	assert.NoErrorf(t, err, "no call response found %v", err)
	isSuccess = assert.Equal(t, "0", code)
	if !isSuccess {
		return err
	}
	_, err = chainA.ExecuteRollback(ctx, res.SerialNo)
	assert.NoErrorf(t, err, "error on excute rollback- %w", err)
	return err
}

func (x *XCallTestSuite) testOneWayMessageWithSize(ctx context.Context, t *testing.T, dataSize int, chainA, chainB chains.Chain, newFunctionCall ...bool) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	_msg := make([]byte, dataSize)
	dst := chainB.(ibc.Chain).Config().ChainID + "/" + chainB.GetContractAddress(dappKey)
	var res *chains.XCallResponse
	var err error
	if len(newFunctionCall) > 0 && newFunctionCall[0] {
		msgType := big.NewInt(1)
		res, err = chainA.NewXCall(ctx, chainB, interchaintest.UserAccount, dst, _msg, msgType, nil)
	} else {
		res, err = chainA.XCall(ctx, chainB, interchaintest.UserAccount, dst, _msg, nil)
	}
	assert.NoError(t, err)

	_, err = chainB.ExecuteCall(ctx, res.RequestID, res.Data)
	assert.NoError(t, err)
}

func (x *XCallTestSuite) testOneWayMessageWithSizeExpectingError(ctx context.Context, t *testing.T, dataSize int, chainA, chainB chains.Chain, newFunctionCall ...bool) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	_msg := make([]byte, dataSize)
	dst := chainB.(ibc.Chain).Config().ChainID + "/" + chainB.GetContractAddress(dappKey)
	var err error
	if len(newFunctionCall) > 0 && newFunctionCall[0] {
		msgType := big.NewInt(1)
		_, err = chainA.NewXCall(ctx, chainB, interchaintest.UserAccount, dst, _msg, msgType, nil)
	} else {
		_, err = chainA.XCall(ctx, chainB, interchaintest.UserAccount, dst, _msg, nil)
	}
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
		} else if strings.Contains(err.Error(), "error on") {
			result = true
		} else {
			result = assert.ObjectsAreEqual(errors.New("UnknownFailure"), err)
		}
		if result {
			t.Logf("Test passed: %v", err)
		} else {
			t.Errorf("Test failed: %v", err)
		}
	}

}
