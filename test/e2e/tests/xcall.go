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

type XCallTestSuite struct {
	*testsuite.E2ETestSuite
	T *testing.T
}

func (x *XCallTestSuite) TextXCall() {
	testcase := "xcall"
	portId := "transfer"
	ctx := context.WithValue(context.TODO(), "testcase", testcase)
	x.Require().NoError(x.DeployXCallMockApp(ctx, portId), "fail to deploy xcall demo dapp")
	createdChains := x.GetChains()
	if len(createdChains) == 3 {
		test3Chains(ctx, createdChains, x)
	}
	if len(createdChains) == 2 {
		test2Chains(ctx, createdChains, x)
	}

}
func test3Chains(ctx context.Context, createdChains []chains.Chain, x *XCallTestSuite) {
	chainA, chainB, chainC := createdChains[0], createdChains[1], createdChains[2]
	fmt.Println("ChainA", chainA.Config().Name)
	fmt.Println("ChainB", chainB.Config().Name)
	fmt.Println("ChainC", chainC.Config().Name)
	x.T.Run("xcall one way message chainA-chainB", func(t *testing.T) {
		fmt.Println("Sending message from src to dst", chainA.Config().Name, chainB.Config().Name)
		err := x.testOneWayMessage(ctx, t, chainA, chainB)
		assert.NoErrorf(t, err, "fail xCall one way message chainA-chainB ::%v\n ", err)
	})

	x.T.Run("xcall one way message chainB-chainA", func(t *testing.T) {
		fmt.Println("Sending message from src to dst", chainB.Config().Name, chainA.Config().Name)
		err := x.testOneWayMessage(ctx, t, chainB, chainA)
		assert.NoErrorf(t, err, "fail xCall one way message chainB-chainA ::%v\n ", err)
	})
	x.T.Run("xcall one way message chainB-chainC", func(t *testing.T) {
		fmt.Println("Sending message from src to dst", chainB.Config().Name, chainC.Config().Name)
		err := x.testOneWayMessage(ctx, t, chainB, chainC)
		assert.NoErrorf(t, err, "fail xCall one way message chainB-chainc ::%v\n ", err)
	})
	x.T.Run("xcall one way message chainC-chainB", func(t *testing.T) {
		fmt.Println("Sending message from src to dst", chainC.Config().Name, chainB.Config().Name)
		err := x.testOneWayMessage(ctx, t, chainC, chainB)
		assert.NoErrorf(t, err, "fail xCall one way message chainC-chainB ::%v\n ", err)
	})
	x.T.Run("xcall one way message chainA-chainC", func(t *testing.T) {
		fmt.Println("Sending message from src to dst", chainA.Config().Name, chainC.Config().Name)
		err := x.testOneWayMessage(ctx, t, chainA, chainC)
		assert.NoErrorf(t, err, "fail xCall one way message chainA-chainC ::%v\n ", err)
	})
	x.T.Run("xcall one way message chainC-chainA", func(t *testing.T) {
		fmt.Println("Sending message from src to dst", chainC.Config().Name, chainA.Config().Name)
		err := x.testOneWayMessage(ctx, t, chainC, chainA)
		assert.NoErrorf(t, err, "fail xCall one way message chainC-chainA ::%v\n ", err)
	})

	x.T.Run("xcall test rollback chainA-chainB", func(t *testing.T) {
		err := x.testRollback(ctx, t, chainA, chainB)
		assert.NoErrorf(t, err, "fail xCall rollback message chainA-chainB ::%v\n ", err)

	})

	x.T.Run("2xcall test rollback chainA-chainC", func(t *testing.T) {
		err := x.testRollback(ctx, t, chainA, chainC)
		assert.NoErrorf(t, err, "fail xCall rollback message chainA-chainC ::%v\n ", err)
	})

	x.T.Run("xcall test rollback chainB-chainA", func(t *testing.T) {
		err := x.testRollback(ctx, t, chainB, chainA)
		assert.NoErrorf(t, err, "fail xcCll rollback message chainB-chainA ::%v\n ", err)
	})
	x.T.Run("2xcall test rollback chainB-chainC", func(t *testing.T) {
		err := x.testRollback(ctx, t, chainB, chainC)
		assert.NoErrorf(t, err, "fail xcCll rollback message chainB-chainC ::%v\n ", err)
	})

	x.T.Run("xcall test rollback chainC-chainA", func(t *testing.T) {
		err := x.testRollback(ctx, t, chainC, chainA)
		assert.NoErrorf(t, err, "fail xcCll rollback message chainC-chainA ::%v\n ", err)
	})

	x.T.Run("2xcall test rollback chainC-chainB", func(t *testing.T) {
		err := x.testRollback(ctx, t, chainC, chainB)
		assert.NoErrorf(t, err, "fail xcCll rollback message chainC-chainB ::%v\n ", err)
	})

	x.T.Run("xcall test send maxSize Data: 2048 bytes", func(t *testing.T) {
		x.T.Run("xcall test send maxSize Data: 2048 bytes A->B", func(t *testing.T) {
			x.testOneWayMessageWithSize(ctx, t, 1200, chainA, chainB)
		})
		x.T.Run("xcall test send maxSize Data: 2048 bytes B->A", func(t *testing.T) {
			x.testOneWayMessageWithSize(ctx, t, 1200, chainB, chainA)
		})
		x.T.Run("xcall test send maxSize Data: 2048 bytes C->A", func(t *testing.T) {
			x.testOneWayMessageWithSize(ctx, t, 1200, chainC, chainA)
		})
		x.T.Run("xcall test send maxSize Data: 2048 bytes A->C", func(t *testing.T) {
			x.testOneWayMessageWithSize(ctx, t, 1200, chainA, chainC)
		})
	})

	x.T.Run("xcall test send maxSize Data: 2049bytes", func(t *testing.T) {
		x.T.Run("xcall test send maxSize Data: 2049 bytes B->A", func(t *testing.T) {
			x.testOneWayMessageWithSizeExpectingError(ctx, t, 2000, chainB, chainA)
		})
		x.T.Run("xcall test send maxSize Data: 2049 bytes A->B", func(t *testing.T) {
			x.testOneWayMessageWithSizeExpectingError(ctx, t, 2100, chainA, chainB)
		})
		x.T.Run("xcall test send maxSize Data: 2049 bytes C->A", func(t *testing.T) {
			x.testOneWayMessageWithSizeExpectingError(ctx, t, 2100, chainC, chainA)
		})
		x.T.Run("xcall test send maxSize Data: 2049 bytes A->C", func(t *testing.T) {
			x.testOneWayMessageWithSizeExpectingError(ctx, t, 2100, chainA, chainC)
		})
	})
}

func test2Chains(ctx context.Context, createdChains []chains.Chain, x *XCallTestSuite) {
	chainA, chainB := createdChains[0], createdChains[1]
	fmt.Println("ChainA", chainA.Config().Name)
	fmt.Println("ChainB", chainB.Config().Name)
	x.T.Run("xcall one way message chainA-chainB", func(t *testing.T) {
		fmt.Println("Sending message from src to dst", chainA.Config().Name, chainB.Config().Name)
		err := x.testOneWayMessage(ctx, t, chainA, chainB)
		assert.NoErrorf(t, err, "fail xCall one way message chainA-chainB ::%v\n ", err)
	})
	x.T.Run("xcall one way message chainB-chainA", func(t *testing.T) {

		fmt.Println("Sending message from src to dst", chainB.Config().Name, chainA.Config().Name)
		err := x.testOneWayMessage(ctx, t, chainB, chainA)
		assert.NoErrorf(t, err, "fail xCall one way message chainB-chainA ::%v\n ", err)
	})
	x.T.Run("2xcall test rollback chainA-chainB", func(t *testing.T) {
		fmt.Println("Sending rollback message from src to dst", chainA.Config().Name, chainB.Config().Name)
		err := x.testRollback(ctx, t, chainA, chainB)
		assert.NoErrorf(t, err, "fail xCall rollback message chainA-chainB ::%v\n ", err)
	})

	x.T.Run("xcall test rollback chainB-chainA", func(t *testing.T) {
		err := x.testRollback(ctx, t, chainB, chainA)
		assert.NoErrorf(t, err, "fail xcCll rollback message chainB-chainA ::%v\n ", err)
	})

	x.T.Run("xcall test send maxSize Data: 2048 bytes A-> B", func(t *testing.T) {
		x.testOneWayMessageWithSize(ctx, t, 1300, chainA, chainB)
	})

	x.T.Run("xcall test send maxSize Data: 2048 bytes B-> A", func(t *testing.T) {
		x.testOneWayMessageWithSize(ctx, t, 1300, chainB, chainA)
	})

	x.T.Run("xcall test send maxSize Data: 2049bytes", func(t *testing.T) {
		x.testOneWayMessageWithSizeExpectingError(ctx, t, 2000, chainB, chainA)
		x.testOneWayMessageWithSizeExpectingError(ctx, t, 2100, chainA, chainB)
	})
}

func (x *XCallTestSuite) testOneWayMessage(ctx context.Context, t *testing.T, chainA, chainB chains.Chain) error {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	msg := "MessageTransferTestingWithoutRollback"
	dst := chainB.Config().ChainID + "/" + chainB.GetContractAddress(dappKey)

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
	dst := chainB.Config().ChainID + "/" + chainB.GetContractAddress(dappKey)
	res, err := chainA.XCall(ctx, chainB, chainB.Config().Name, dst, []byte(msg), []byte(rollback))
	isSuccess := assert.NoErrorf(t, err, "error on sending packet- %v", err)
	if !isSuccess {
		return err
	}
	height, err := chainA.Height(ctx)
	assert.NoErrorf(t, err, "error getting height %v", err)
	code, err := chainA.FindCallResponse(ctx, height, res.SerialNo)
	assert.NoErrorf(t, err, "no call response found %v", err)
	isSuccess = assert.Equal(t, "0", code)
	if !isSuccess {
		return err
	}
	time.Sleep(3 * time.Second)
	_, err = chainA.FindRollbackExecutedMessage(ctx, height, res.SerialNo)
	assert.NoErrorf(t, err, "error on excute rollback- %w", err)
	fmt.Println("Data Transfer Testing Without Rollback from " + chainA.Config().ChainID + " to " + chainB.Config().ChainID + " with data " + msg + " and rollback:" + rollback + " PASSED")
	return err
}

func (x *XCallTestSuite) testOneWayMessageWithSize(ctx context.Context, t *testing.T, dataSize int, chainA, chainB chains.Chain) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	_msg := make([]byte, dataSize)
	dst := chainB.Config().ChainID + "/" + chainB.GetContractAddress(dappKey)
	_, err := chainA.XCall(ctx, chainB, chainB.Config().Name, dst, _msg, nil)
	assert.NoError(t, err)
}

func (x *XCallTestSuite) testOneWayMessageWithSizeExpectingError(ctx context.Context, t *testing.T, dataSize int, chainA, chainB chains.Chain) {
	testcase := ctx.Value("testcase").(string)
	dappKey := fmt.Sprintf("dapp-%s", testcase)
	_msg := make([]byte, dataSize)
	dst := chainB.Config().ChainID + "/" + chainB.GetContractAddress(dappKey)
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
		} else {
			t.Errorf("Test failed: %v", err)
		}
	}

}
