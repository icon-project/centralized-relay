package retry

import (
	"context"
	"math"
	"time"

	"go.uber.org/zap"
)

const (
	BaseRPCRetryDelay   = time.Second
	MaxRPCRetryDelay    = 30 * time.Second
	RPCMaxRetryAttempts = 5
	RetryPower          = 3
)

func Retry(ctx context.Context, logger *zap.Logger, operation func() error, zapFields []zap.Field) error {
	var retryCount uint8
	for retryCount < RPCMaxRetryAttempts {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := operation()
		if err == nil {
			return nil
		}

		retryCount++
		if retryCount >= RPCMaxRetryAttempts {
			logger.Error("operation failed", append(zapFields, zap.Uint8("attempt", retryCount), zap.Error(err))...)
			return err
		}

		delay := time.Duration(math.Pow(RetryPower, float64(retryCount))) * BaseRPCRetryDelay
		if delay > MaxRPCRetryDelay {
			delay = MaxRPCRetryDelay
		}
		logger.Warn("operation failed, retrying...", append(zapFields, zap.Uint8("attempt", retryCount), zap.Duration("retrying_in", delay), zap.Error(err))...)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return nil
}
