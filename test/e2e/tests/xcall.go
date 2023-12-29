package tests

import (
	"context"
	"errors"
	"fmt"
	interchaintest "github.com/icon-project/centralized-relay/test"
	"github.com/icon-project/centralized-relay/test/chains"
	"github.com/icon-project/centralized-relay/test/interchaintest/ibc"
	"github.com/icon-project/centralized-relay/test/testsuite"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

type XCallTestSuite struct {
	*testsuite.E2ETestSuite
	T *testing.T
}

func (x *XCallTestSuite) TextXCall() {
	testcase := "xcall"
	portId := "transfer"
	ctx := context.WithValue(context.TODO(), "testcase", testcase)
	//x.Require().NoError(x.SetupXCall(ctx), "fail to setup xcall")
	x.Require().NoError(x.DeployXCallMockApp(ctx, portId), "fail to deploy xcall demo dapp")
	chainA, chainB := x.GetChains()
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
		x.testOneWayMessageWithSize(ctx, t, 1300, chainA, chainB)
		x.testOneWayMessageWithSize(ctx, t, 1300, chainB, chainA)
	})

	x.T.Run("xcall test send maxSize Data: 2049bytes", func(t *testing.T) {
		x.testOneWayMessageWithSizeExpectingError(ctx, t, 2000, chainB, chainA)
		x.testOneWayMessageWithSizeExpectingError(ctx, t, 2100, chainA, chainB)
	})

}

func (x *XCallTestSuite) TestXCallFlush() {
	testcase := "packet-flush"
	portId := "transfer-1"
	ctx := context.WithValue(context.TODO(), "testcase", testcase)
	//x.Require().NoError(x.SetupXCall(ctx), "fail to setup xcall")
	x.Require().NoError(x.DeployXCallMockApp(ctx, portId), "fail to deploy xcall dapp")
	chainA, chainB := x.GetChains()
	x.T.Run("xcall packet flush chainA-chainB", func(t *testing.T) {
		err := x.testPacketFlush(ctx, chainA, chainB)
		assert.NoErrorf(t, err, "xcall packet flush chainA-chainB ::%v\n ", err)

	})

	x.T.Run("xcall packet flush chainB-chainA", func(t *testing.T) {
		err := x.testPacketFlush(ctx, chainB, chainA)
		assert.NoErrorf(t, err, "xcall packet flush chainB-chainA ::%v\n ", err)
	})
}

func (x *XCallTestSuite) testPacketFlush(ctx context.Context, chainA, chainB chains.Chain) error {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	msg := "flush-msg"
	heightB, _ := chainB.Height(ctx)

	dst := chainB.(ibc.Chain).Config().ChainID + "/" + chainB.GetContractAddress(dappKey)

	err := chainB.PauseNode(ctx)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to pause node %s - %v", chainB.Config().Name, err))
	}

	ctx, err = chainA.SendPacketXCall(ctx, interchaintest.UserAccount, dst, []byte(msg), nil)

	if err != nil {
		return errors.New(fmt.Sprintf("failed send xCall message find eventlog - %v", err))
	}
	sn := ctx.Value("sn").(string)
	fmt.Printf("sn-%s\n", sn)
	waitDuration := 90 * time.Second
	fmt.Printf("Wait for %v \n", waitDuration)
	// TODO: Wait for 1.5 mins (90 seconds)
	time.Sleep(waitDuration)

	err = chainB.UnpauseNode(ctx)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to unpause node %s - %v", chainB.Config().Name, err))
	}

	//wait 90 sec
	fmt.Printf("Wait for %v after node unpause\n", waitDuration)
	time.Sleep(waitDuration)

	reqId, destData, err := chainB.FindCallMessage(ctx, heightB, chainA.Config().ChainID+"/"+chainA.GetContractAddress(dappKey), chainB.GetContractAddress(dappKey), sn)

	ctx, err = chainB.ExecuteCall(ctx, reqId, destData)
	if err != nil {
		return errors.New(fmt.Sprintf("error on execute call packet req-id::%s- %v", reqId, err))
	}

	return nil
}

func (x *XCallTestSuite) testOneWayMessage(ctx context.Context, t *testing.T, chainA, chainB chains.Chain) error {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	msg := "MessageTransferTestingWithoutRollback"
	dst := chainB.(ibc.Chain).Config().ChainID + "/" + chainB.GetContractAddress(dappKey)

	res, err := chainA.XCall(ctx, chainB, interchaintest.UserAccount, dst, []byte(msg), nil)
	result := assert.NoErrorf(t, err, "error on sending packet- %v", err)
	if !result {
		return err
	}
	ctx, err = chainB.ExecuteCall(ctx, res.RequestID, res.Data)
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

func (x *XCallTestSuite) testRollback(ctx context.Context, t *testing.T, chainA, chainB chains.Chain) error {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	msg := "rollback"
	rollback := "RollbackDataTesting"
	dst := chainB.(ibc.Chain).Config().ChainID + "/" + chainB.GetContractAddress(dappKey)
	res, err := chainA.XCall(ctx, chainB, interchaintest.UserAccount, dst, []byte(msg), []byte(rollback))
	isSuccess := assert.NoErrorf(t, err, "error on sending packet- %v", err)
	if !isSuccess {
		return err
	}
	height, err := chainA.(ibc.Chain).Height(ctx)
	_, err = chainB.ExecuteCall(ctx, res.RequestID, res.Data)
	code, err := chainA.FindCallResponse(ctx, height, res.SerialNo)
	isSuccess = assert.NoErrorf(t, err, "no call response found %v", err)
	isSuccess = assert.Equal(t, "0", code)
	if !isSuccess {
		return err
	}
	ctx, err = chainA.ExecuteRollback(ctx, res.SerialNo)
	assert.NoErrorf(t, err, "error on excute rollback- %w", err)
	return err
}

func (x *XCallTestSuite) testOneWayMessageWithSize(ctx context.Context, t *testing.T, dataSize int, chainA, chainB chains.Chain) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	_msg := make([]byte, dataSize)
	dst := chainB.(ibc.Chain).Config().ChainID + "/" + chainB.GetContractAddress(dappKey)
	res, err := chainA.XCall(ctx, chainB, interchaintest.UserAccount, dst, _msg, nil)
	assert.NoError(t, err)

	_, err = chainB.ExecuteCall(ctx, res.RequestID, res.Data)
	assert.NoError(t, err)
}

func (x *XCallTestSuite) testOneWayMessageWithSizeExpectingError(ctx context.Context, t *testing.T, dataSize int, chainA, chainB chains.Chain) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	_msg := make([]byte, dataSize)
	dst := chainB.(ibc.Chain).Config().ChainID + "/" + chainB.GetContractAddress(dappKey)
	_, err := chainA.XCall(ctx, chainB, interchaintest.UserAccount, dst, _msg, nil)
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
		} else {
			t.Errorf("Test failed: %v", err)
		}
	}

}
