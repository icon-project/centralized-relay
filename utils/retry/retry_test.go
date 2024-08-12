package retry

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"
)

func TestRetry(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()

	// Test case 1: Operation succeeds on the first attempt
	attempts := 0
	operation := func() error {
		attempts++
		return nil
	}
	err := Retry(ctx, logger, operation, nil)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got: %d", attempts)
	}

	// Test case 2: Operation fails on all attempts
	attempts = 0
	expectedErr := errors.New("operation failed")
	operation = func() error {
		attempts++
		return expectedErr
	}
	err = Retry(ctx, logger, operation, nil)
	if err != expectedErr {
		t.Errorf("Expected error: %v, got: %v", expectedErr, err)
	}
	if attempts != RPCMaxRetryAttempts {
		t.Errorf("Expected %d attempts, got: %d", RPCMaxRetryAttempts, attempts)
	}

	// Test case 3: Operation succeeds on the third attempt
	attempts = 0
	operation = func() error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary failure")
		}
		return nil
	}
	err = Retry(ctx, logger, operation, nil)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got: %d", attempts)
	}

	// Test case 4: Context canceled
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	attempts = 0
	operation = func() error {
		attempts++
		return errors.New("operation failed")
	}
	err = Retry(ctx, logger, operation, nil)
	if err != ctx.Err() {
		t.Errorf("Expected error: %v, got: %v", ctx.Err(), err)
	}
	if attempts != 0 {
		t.Errorf("Expected 0 attempts, got: %d", attempts)
	}

	// Test case 5: Operation fails on the first attempt
	ctx = context.Background()
	attempts = 0
	expectedErr = errors.New("operation failed on the first attempt")
	operation = func() error {
		attempts++
		return expectedErr
	}
	err = Retry(ctx, logger, operation, nil)
	if err != expectedErr {
		t.Errorf("Expected error: %v, got: %v", expectedErr, err)
	}

	// Test case 6: Operation fails on all attempts
	ctx = context.Background()
	attempts = 0
	expectedErr = errors.New("operation failed on all attempts")
	operation = func() error {
		attempts++
		return expectedErr
	}
	err = Retry(ctx, logger, operation, nil)
	if err != expectedErr {
		t.Errorf("Expected error: %v, got: %v", expectedErr, err)
	}
	if attempts != RPCMaxRetryAttempts {
		t.Errorf("Expected %d attempts, got: %d", RPCMaxRetryAttempts, attempts)
	}
}
